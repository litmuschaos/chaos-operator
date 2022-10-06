/*
Copyright 2019 LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1"
	"github.com/litmuschaos/chaos-operator/pkg/analytics"
	dynamicclientset "github.com/litmuschaos/chaos-operator/pkg/client/dynamic"
	clientset "github.com/litmuschaos/chaos-operator/pkg/client/kubernetes"
	"github.com/litmuschaos/chaos-operator/pkg/resource"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/types"
	"github.com/litmuschaos/chaos-operator/pkg/utils"
	"github.com/litmuschaos/chaos-operator/pkg/utils/retry"
	"github.com/litmuschaos/elves/kubernetes/container"
	"github.com/litmuschaos/elves/kubernetes/pod"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"os"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
	"time"
)

const finalizer = "chaosengine.litmuschaos.io/finalizer"

// ChaosEngineReconciler reconciles a ChaosEngine object
type ChaosEngineReconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	// Used for serializing and deserializing API objects(group, version, and kind)
	Scheme *runtime.Scheme
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	Recorder record.EventRecorder
}

// reconcileEngine contains details of reconcileEngine
type reconcileEngine struct {
	r         *ChaosEngineReconciler
	reqLogger logr.Logger
}

// podEngineRunner contains the information of pod
type podEngineRunner struct {
	pod, engineRunner *corev1.Pod
	*reconcileEngine
}

//+kubebuilder:rbac:groups=litmuschaos.io,resources=chaosengines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=litmuschaos.io,resources=chaosengines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=litmuschaos.io,resources=chaosengines/finalizers,verbs=update

// Reconcile reads that state of the cluster for a ChaosEngine object and makes changes based on the state read
// and what is in the ChaosEngine.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue
func (r *ChaosEngineReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := startReqLogger(request)
	engine := &chaosTypes.EngineInfo{}

	if err := r.getChaosEngineInstance(engine, request); err != nil {
		if k8serrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Handle deletion of ChaosEngine
	if engine.Instance.ObjectMeta.GetDeletionTimestamp() != nil {
		return r.reconcileForDelete(engine, request)
	}

	// Start the reconcile by setting default values into ChaosEngine
	if err := r.initEngine(engine); err != nil {
		return reconcile.Result{}, err
	}

	// Handling of normal execution of ChaosEngine
	if engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateActive && engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusInitialized {
		return r.reconcileForCreationAndRunning(engine, reqLogger)
	}

	// Handling Graceful completion of ChaosEngine
	if engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateStop && engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusCompleted {
		return r.reconcileForComplete(engine, request)
	}

	// Handling forceful Abort of ChaosEngine
	if engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateStop && engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusInitialized {
		return r.reconcileForDelete(engine, request)
	}

	// Handling restarting of ChaosEngine post Abort
	if engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateActive && (engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusStopped) {
		return r.reconcileForRestartAfterAbort(engine, request)
	}

	// Handling restarting of ChaosEngine post Completion
	if engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateActive && (engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusCompleted) {
		return r.reconcileForRestartAfterComplete(engine, request)
	}

	return ctrl.Result{}, nil
}

// getChaosRunnerENV return the env required for chaos-runner
func getChaosRunnerENV(cr *litmuschaosv1alpha1.ChaosEngine, aExList []string, ClientUUID string) []corev1.EnvVar {
	appNS := cr.Spec.Appinfo.Appns
	if appNS == "" {
		appNS = cr.Namespace
	}

	var envDetails utils.ENVDetails
	envDetails.SetEnv("CHAOSENGINE", cr.Name).
		SetEnv("APP_LABEL", cr.Spec.Appinfo.Applabel).
		SetEnv("APP_KIND", cr.Spec.Appinfo.AppKind).
		SetEnv("APP_NAMESPACE", appNS).
		SetEnv("EXPERIMENT_LIST", fmt.Sprint(strings.Join(aExList, ","))).
		SetEnv("CHAOS_SVC_ACC", cr.Spec.ChaosServiceAccount).
		SetEnv("AUXILIARY_APPINFO", cr.Spec.AuxiliaryAppInfo).
		SetEnv("CLIENT_UUID", ClientUUID).
		SetEnv("CHAOS_NAMESPACE", cr.Namespace).
		SetEnv("ANNOTATION_CHECK", cr.Spec.AnnotationCheck).
		SetEnv("ANNOTATION_KEY", resource.GetAnnotationKey())

	return envDetails.ENV
}

// getChaosRunnerLabels return the labels required for chaos-runner
func getChaosRunnerLabels(cr *litmuschaosv1alpha1.ChaosEngine) map[string]string {
	return map[string]string{
		"app":                         cr.Name,
		"chaosUID":                    string(cr.UID),
		"app.kubernetes.io/component": "chaos-runner",
		"app.kubernetes.io/part-of":   "litmus",
	}
}

// newGoRunnerPodForCR defines a new go-based Runner Pod
func (r *ChaosEngineReconciler) newGoRunnerPodForCR(engine *chaosTypes.EngineInfo) (*corev1.Pod, error) {
	engine.VolumeOpts.VolumeOperations(engine.Instance.Spec.Components.Runner.ConfigMaps, engine.Instance.Spec.Components.Runner.Secrets)

	containerForRunner := container.NewBuilder().
		WithEnvsNew(getChaosRunnerENV(engine.Instance, engine.AppExperiments, analytics.ClientUUID)).
		WithName("chaos-runner").
		WithImage(engine.Instance.Spec.Components.Runner.Image).
		WithImagePullPolicy(corev1.PullIfNotPresent)

	if engine.Instance.Spec.Components.Runner.ImagePullPolicy != "" {
		containerForRunner.WithImagePullPolicy(engine.Instance.Spec.Components.Runner.ImagePullPolicy)
	}

	if engine.Instance.Spec.Components.Runner.Args != nil {
		containerForRunner.WithArgumentsNew(engine.Instance.Spec.Components.Runner.Args)
	}

	if engine.VolumeOpts.VolumeMounts != nil {
		containerForRunner.WithVolumeMountsNew(engine.VolumeOpts.VolumeMounts)
	}

	if engine.Instance.Spec.Components.Runner.Command != nil {
		containerForRunner.WithCommandNew(engine.Instance.Spec.Components.Runner.Command)
	}

	if !reflect.DeepEqual(engine.Instance.Spec.Components.Runner.Resources, corev1.ResourceRequirements{}) {
		containerForRunner.WithResourceRequirements(engine.Instance.Spec.Components.Runner.Resources)
	}

	podForRunner := pod.NewBuilder().
		WithName(engine.Instance.Name + "-runner").
		WithNamespace(engine.Instance.Namespace).
		WithAnnotations(engine.Instance.Spec.Components.Runner.RunnerAnnotation).
		WithLabels(getChaosRunnerLabels(engine.Instance)).
		WithServiceAccountName(engine.Instance.Spec.ChaosServiceAccount).
		WithRestartPolicy("OnFailure").
		WithContainerBuilder(containerForRunner)

	if engine.Instance.Spec.Components.Runner.Tolerations != nil {
		podForRunner.WithTolerations(engine.Instance.Spec.Components.Runner.Tolerations...)
	}

	if len(engine.Instance.Spec.Components.Runner.NodeSelector) != 0 {
		podForRunner.WithNodeSelector(engine.Instance.Spec.Components.Runner.NodeSelector)
	}

	if engine.VolumeOpts.VolumeBuilders != nil {
		podForRunner.WithVolumeBuilders(engine.VolumeOpts.VolumeBuilders)
	}

	if engine.Instance.Spec.Components.Runner.ImagePullSecrets != nil {
		podForRunner.WithImagePullSecrets(engine.Instance.Spec.Components.Runner.ImagePullSecrets)
	}

	runnerPod, err := podForRunner.Build()
	if err != nil {
		return nil, err
	}
	if err := controllerutil.SetControllerReference(engine.Instance, runnerPod, r.Scheme); err != nil {
		return nil, err
	}
	return runnerPod, nil
}

// initializeApplicationInfo to initialize application info
func initializeApplicationInfo(instance *litmuschaosv1alpha1.ChaosEngine, appInfo *chaosTypes.ApplicationInfo) (*chaosTypes.ApplicationInfo, error) {
	if instance == nil {
		return nil, errors.New("chaosengine is empty")
	}

	if instance.Spec.Appinfo.Applabel != "" {
		appInfo.Label = instance.Spec.Appinfo.Applabel
	}

	if instance.Spec.Appinfo.Appns != "" {
		appInfo.Namespace = instance.Spec.Appinfo.Appns
	} else {
		appInfo.Namespace = instance.Namespace
	}
	appInfo.Kind = instance.Spec.Appinfo.AppKind

	appInfo.ExperimentList = instance.Spec.Experiments
	appInfo.ServiceAccountName = instance.Spec.ChaosServiceAccount

	return appInfo, nil
}

// engineRunnerPod to Check if the engineRunner pod already exists, else create
func engineRunnerPod(runnerPod *podEngineRunner) error {
	if err := runnerPod.r.Client.Get(context.TODO(), types.NamespacedName{Name: runnerPod.engineRunner.Name, Namespace: runnerPod.engineRunner.Namespace}, runnerPod.pod); err != nil && k8serrors.IsNotFound(err) {
		runnerPod.reqLogger.Info("Creating a new engineRunner Pod", "Pod.Namespace", runnerPod.engineRunner.Namespace, "Pod.Name", runnerPod.engineRunner.Name)
		if err = runnerPod.r.Client.Create(context.TODO(), runnerPod.engineRunner); err != nil {
			return err
		}

		// Pod created successfully - don't reconcile
		runnerPod.reqLogger.Info("engineRunner Pod created successfully")
		return nil
	} else if err != nil {
		return err
	}
	runnerPod.reqLogger.Info("Skip reconcile: engineRunner Pod already exists", "Pod.Namespace", runnerPod.pod.Namespace, "Pod.Name", runnerPod.pod.Name)
	return nil
}

// Fetch the ChaosEngine instance
func (r *ChaosEngineReconciler) getChaosEngineInstance(engine *chaosTypes.EngineInfo, request reconcile.Request) error {
	instance := &litmuschaosv1alpha1.ChaosEngine{}
	if err := r.Client.Get(context.TODO(), request.NamespacedName, instance); err != nil {
		// Error reading the object - reconcile the request.
		return err
	}
	engine.Instance = instance
	return nil
}

// Get application details
func getApplicationDetail(engine *chaosTypes.EngineInfo) error {
	applicationInfo := &chaosTypes.ApplicationInfo{}

	appInfo, err := initializeApplicationInfo(engine.Instance, applicationInfo)
	if err != nil {
		return err
	}

	engine.AppInfo = appInfo
	var appExperiments []string
	for _, exp := range appInfo.ExperimentList {
		appExperiments = append(appExperiments, exp.Name)
	}
	engine.AppExperiments = appExperiments

	chaosTypes.Log.Info("App Label derived from Chaosengine is ", "appLabels", appInfo.Label)
	chaosTypes.Log.Info("App NS derived from Chaosengine is ", "appNamespace", appInfo.Namespace)
	chaosTypes.Log.Info("Exp list derived from chaosengine is ", "appExpirements", appExperiments)
	chaosTypes.Log.Info("Runner image derived from chaosengine is", "runnerImage", engine.Instance.Spec.Components.Runner.Image)
	chaosTypes.Log.Info("Annotation check is ", "annotationCheck", engine.Instance.Spec.AnnotationCheck)
	return nil
}

// Check if the engineRunner pod already exists, else create
func (r *ChaosEngineReconciler) checkEngineRunnerPod(engine *chaosTypes.EngineInfo, reqLogger logr.Logger) error {
	if len(engine.AppExperiments) == 0 {
		return errors.New("application experiment list is empty")
	}

	engineRunner, err := r.newGoRunnerPodForCR(engine)
	if err != nil {
		return err
	}

	// Create an object of engine reconcile.
	engineReconcile := &reconcileEngine{
		r:         r,
		reqLogger: reqLogger,
	}
	// Creates an object of engineRunner Pod
	runnerPod := &podEngineRunner{
		pod:             &corev1.Pod{},
		engineRunner:    engineRunner,
		reconcileEngine: engineReconcile,
	}

	return engineRunnerPod(runnerPod)
}

// setChaosResourceImage take the runner image from engine spec
// if it is not there then it will take from chaos-operator env
// at last if it is not able to find image in engine spec and operator env then it will take default images
func setChaosResourceImage(engine *chaosTypes.EngineInfo) {
	ChaosRunnerImage := os.Getenv("CHAOS_RUNNER_IMAGE")

	if engine.Instance.Spec.Components.Runner.Image == "" && ChaosRunnerImage == "" {
		engine.Instance.Spec.Components.Runner.Image = chaosTypes.DefaultChaosRunnerImage
	} else if engine.Instance.Spec.Components.Runner.Image == "" {
		engine.Instance.Spec.Components.Runner.Image = ChaosRunnerImage
	}
}

// getAnnotationCheck() checks for annotation on the application
func getAnnotationCheck(engine *chaosTypes.EngineInfo) error {
	if engine.Instance.Spec.AnnotationCheck == "" {
		engine.Instance.Spec.AnnotationCheck = chaosTypes.DefaultAnnotationCheck
	}

	if engine.Instance.Spec.AnnotationCheck != "true" && engine.Instance.Spec.AnnotationCheck != "false" {
		return fmt.Errorf("annotationCheck '%s', is not supported it should be either true or false", engine.Instance.Spec.AnnotationCheck)
	}

	return nil
}

// reconcileForDelete reconciles for deletion/force deletion of Chaos Engine
func (r *ChaosEngineReconciler) reconcileForDelete(engine *chaosTypes.EngineInfo, request reconcile.Request) (reconcile.Result, error) {
	patch := client.MergeFrom(engine.Instance.DeepCopy())

	chaosTypes.Log.Info("Checking if there are any chaos resources to be deleted for", "chaosengine", engine.Instance.Name)

	chaosPodList := &corev1.PodList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
		client.MatchingLabels{"chaosUID": string(engine.Instance.UID)},
	}
	if err := r.Client.List(context.TODO(), chaosPodList, opts...); err != nil {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos stop) Unable to list chaos experiment pods")
		return reconcile.Result{}, err
	}

	if len(chaosPodList.Items) != 0 {
		chaosTypes.Log.Info("Performing a force delete of chaos experiment pods", "chaosengine", engine.Instance.Name)
		err := r.forceRemoveChaosResources(engine, request)
		if err != nil {
			r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos stop) Unable to delete chaos experiment pods")
			return reconcile.Result{}, err
		}
	}

	// update the chaos status in result for abort cases
	if err := r.updateChaosStatus(engine, request); err != nil {
		return reconcile.Result{}, err
	}

	if engine.Instance.ObjectMeta.Finalizers != nil {
		engine.Instance.ObjectMeta.Finalizers = utils.RemoveString(engine.Instance.ObjectMeta.Finalizers, "chaosengine.litmuschaos.io/finalizer")
	}

	// Update ChaosEngine ExperimentStatuses, with aborted Status.
	updateExperimentStatusesForStop(engine)
	engine.Instance.Status.EngineStatus = litmuschaosv1alpha1.EngineStatusStopped

	if err := r.Client.Patch(context.TODO(), engine.Instance, patch); err != nil && !k8serrors.IsNotFound(err) {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos stop) Unable to update chaosengine")
		return reconcile.Result{}, fmt.Errorf("unable to remove finalizer from chaosEngine Resource, due to error: %v", err)
	}

	// we are repeating this condition/check here as we want the events for 'ChaosEngineStopped'
	// generated only after successful finalizer removal from the chaosengine resource
	if len(chaosPodList.Items) != 0 {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeNormal, "ChaosEngineStopped", "Chaos resources deleted successfully")
	}

	return reconcile.Result{}, nil
}

// forceRemoveChaosResources force removes all chaos-related pods
func (r *ChaosEngineReconciler) forceRemoveChaosResources(engine *chaosTypes.EngineInfo, request reconcile.Request) error {
	optsDelete := []client.DeleteAllOfOption{client.InNamespace(request.NamespacedName.Namespace), client.MatchingLabels{"chaosUID": string(engine.Instance.UID)}, client.PropagationPolicy(v1.DeletePropagationBackground)}
	if engine.Instance.Spec.TerminationGracePeriodSeconds != 0 {
		optsDelete = append(optsDelete, client.GracePeriodSeconds(engine.Instance.Spec.TerminationGracePeriodSeconds))
	}

	var (
		deleteEvent []string
		err         []error
	)

	if errJob := r.Client.DeleteAllOf(context.TODO(), &batchv1.Job{}, optsDelete...); errJob != nil {
		err = append(err, errJob)
		deleteEvent = append(deleteEvent, "Jobs, ")
	}

	if errPod := r.Client.DeleteAllOf(context.TODO(), &corev1.Pod{}, optsDelete...); errPod != nil {
		err = append(err, errPod)
		deleteEvent = append(deleteEvent, "Pods, ")
	}
	if err != nil {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos stop) Unable to delete chaos resources: %v allocated to chaosengine", strings.Join(deleteEvent, ""))
		return fmt.Errorf("unable to delete ChaosResources due to %v", err)
	}

	return nil
}

// updateEngineState updates Chaos Engine Status with given State
func (r *ChaosEngineReconciler) updateEngineState(engine *chaosTypes.EngineInfo, state litmuschaosv1alpha1.EngineState) error {
	patch := client.MergeFrom(engine.Instance.DeepCopy())
	engine.Instance.Spec.EngineState = state

	if err := r.Client.Patch(context.TODO(), engine.Instance, patch); err != nil {
		return fmt.Errorf("unable to patch state of chaosEngine Resource, due to error: %v", err)
	}

	return nil
}

// checkRunnerContainerCompletedStatus check for the runner pod's container status for Completed
func (r *ChaosEngineReconciler) checkRunnerContainerCompletedStatus(engine *chaosTypes.EngineInfo) (bool, error) {
	runnerPod := corev1.Pod{}
	isCompleted := false

	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: engine.Instance.Name + "-runner", Namespace: engine.Instance.Namespace}, &runnerPod)
	if err != nil {
		return isCompleted, err
	}

	if runnerPod.Status.Phase == corev1.PodRunning || runnerPod.Status.Phase == corev1.PodSucceeded {
		for _, container := range runnerPod.Status.ContainerStatuses {
			if container.Name == "chaos-runner" && container.State.Terminated != nil {
				if container.State.Terminated.Reason == "Completed" {
					isCompleted = !container.Ready
				}
			}
		}
	}

	return isCompleted, nil
}

// gracefullyRemoveDefaultChaosResources removes all chaos-resources gracefully
func (r *ChaosEngineReconciler) gracefullyRemoveDefaultChaosResources(engine *chaosTypes.EngineInfo, request reconcile.Request) (reconcile.Result, error) {
	if engine.Instance.Spec.JobCleanUpPolicy == litmuschaosv1alpha1.CleanUpPolicyDelete {
		if err := r.gracefullyRemoveChaosPods(engine, request); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// gracefullyRemoveChaosPods removes chaos default resources gracefully
func (r *ChaosEngineReconciler) gracefullyRemoveChaosPods(engine *chaosTypes.EngineInfo, request reconcile.Request) error {
	optsList := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace), client.MatchingLabels{"app": engine.Instance.Name, "chaosUID": string(engine.Instance.UID)},
	}

	var podList corev1.PodList
	if errList := r.Client.List(context.TODO(), &podList, optsList...); errList != nil {
		return errList
	}

	for _, v := range podList.Items {
		if errDel := r.Client.Delete(context.TODO(), &v, []client.DeleteOption{}...); errDel != nil {
			return errDel
		}
	}

	return nil
}

// reconcileForComplete reconciles for graceful completion of Chaos Engine
func (r *ChaosEngineReconciler) reconcileForComplete(engine *chaosTypes.EngineInfo, request reconcile.Request) (reconcile.Result, error) {
	if _, err := r.gracefullyRemoveDefaultChaosResources(engine, request); err != nil {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos completion) Unable to delete chaos pods upon chaos completion")
		return reconcile.Result{}, err
	}

	if err := r.updateEngineState(engine, litmuschaosv1alpha1.EngineStateStop); err != nil {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos completion) Unable to update chaosengine")
		return reconcile.Result{}, fmt.Errorf("unable to Update Engine State: %v", err)
	}

	return reconcile.Result{}, nil
}

// reconcileForRestartAfterAbort reconciles for restart of ChaosEngine after it was aborted previously
func (r *ChaosEngineReconciler) reconcileForRestartAfterAbort(engine *chaosTypes.EngineInfo, request reconcile.Request) (reconcile.Result, error) {
	if err := r.forceRemoveChaosResources(engine, request); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.updateEngineForRestart(engine); err != nil {
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil

}

// reconcileForRestartAfterComplete reconciles for restart of ChaosEngine after it has completed successfully
func (r *ChaosEngineReconciler) reconcileForRestartAfterComplete(engine *chaosTypes.EngineInfo, request reconcile.Request) (reconcile.Result, error) {
	patch := client.MergeFrom(engine.Instance.DeepCopy())

	if err := r.forceRemoveChaosResources(engine, request); err != nil {
		return reconcile.Result{}, err
	}

	engine.Instance.Status.EngineStatus = litmuschaosv1alpha1.EngineStatusInitialized
	engine.Instance.Status.Experiments = nil

	// finalizers have been retained in a completed chaosengine till this point (as chaos pods may be "retained")
	// as per the jobCleanUpPolicy. Stale finalizer is removed so that initEngine() generates the
	// ChaosEngineInitialized event and re-adds the finalizer before starting chaos.

	if engine.Instance.ObjectMeta.Finalizers != nil {
		engine.Instance.ObjectMeta.Finalizers = utils.RemoveString(engine.Instance.ObjectMeta.Finalizers, "chaosengine.litmuschaos.io/finalizer")
	}

	if err := r.Client.Patch(context.TODO(), engine.Instance, patch); err != nil {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos restart) Unable to update chaosengine")
		return reconcile.Result{}, fmt.Errorf("unable to patch state & remove stale finalizer in chaosEngine Resource, due to error: %v", err)
	}

	return reconcile.Result{}, nil
}

// initEngine initialize Chaos Engine, and add a finalizer to it.
func (r *ChaosEngineReconciler) initEngine(engine *chaosTypes.EngineInfo) error {
	if engine.Instance.Spec.EngineState == "" {
		engine.Instance.Spec.EngineState = litmuschaosv1alpha1.EngineStateActive
	}

	if engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateActive && engine.Instance.Status.EngineStatus == "" {
		engine.Instance.Status.EngineStatus = litmuschaosv1alpha1.EngineStatusInitialized
	}

	if engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusInitialized {
		if engine.Instance.ObjectMeta.Finalizers == nil {
			engine.Instance.ObjectMeta.Finalizers = append(engine.Instance.ObjectMeta.Finalizers, finalizer)
			if err := r.Client.Update(context.TODO(), engine.Instance, &client.UpdateOptions{}); err != nil {
				return fmt.Errorf("unable to initialize ChaosEngine, because of Update Error: %v", err)
			}
			// generate the ChaosEngineInitialized event once finalizer has been added
			r.Recorder.Eventf(engine.Instance, corev1.EventTypeNormal, "ChaosEngineInitialized", "Identifying app under test & launching %s", engine.Instance.Name+"-runner")
		}
	}

	return nil
}

// reconcileForCreationAndRunning reconciles for Chaos execution of Chaos Engine
func (r *ChaosEngineReconciler) reconcileForCreationAndRunning(engine *chaosTypes.EngineInfo, reqLogger logr.Logger) (reconcile.Result, error) {
	if err := r.validateAnnontatedApplication(engine); err != nil {
		if stopEngineWithAnnotationErrorMessage := r.updateEngineState(engine, litmuschaosv1alpha1.EngineStateStop); stopEngineWithAnnotationErrorMessage != nil {
			r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos stop) Unable to update chaosengine")
			return reconcile.Result{}, fmt.Errorf("unable to Update Engine State: %v", err)
		}
		return reconcile.Result{}, err
	}

	// Check if the engineRunner pod already exists, else create
	if err := r.checkEngineRunnerPod(engine, reqLogger); err != nil {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos start) Unable to get chaos resources")
		return reconcile.Result{}, err
	}

	isCompleted, err := r.checkRunnerContainerCompletedStatus(engine)
	if err != nil {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos running) Unable to check chaos status")
		return reconcile.Result{}, err
	}

	if isCompleted {
		if err := r.updateEngineForComplete(engine, isCompleted); err != nil {
			r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(chaos completed) Unable to update chaos engine")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// updateExperimentStatusesForStop updates ChaosEngine.Status.Experiment with Abort Status.
func updateExperimentStatusesForStop(engine *chaosTypes.EngineInfo) {
	for i := range engine.Instance.Status.Experiments {
		if engine.Instance.Status.Experiments[i].Status == litmuschaosv1alpha1.ExperimentStatusRunning || engine.Instance.Status.Experiments[i].Status == litmuschaosv1alpha1.ExperimentStatusWaiting {
			engine.Instance.Status.Experiments[i].Status = litmuschaosv1alpha1.ExperimentStatusAborted
			engine.Instance.Status.Experiments[i].Verdict = "Stopped"
			engine.Instance.Status.Experiments[i].LastUpdateTime = v1.Now()
		}
	}
}

func startReqLogger(request reconcile.Request) logr.Logger {
	reqLogger := chaosTypes.Log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ChaosEngine")

	return reqLogger
}

func (r *ChaosEngineReconciler) validateAnnontatedApplication(engine *chaosTypes.EngineInfo) error {
	// Get the image for runner pod from chaosengine spec,operator env or default values.
	setChaosResourceImage(engine)

	clientSet, err := clientset.CreateClientSet()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamicclientset.CreateClientSet()
	if err != nil {
		return err
	}

	// getAnnotationCheck fetch the annotationCheck from engine spec
	if err = getAnnotationCheck(engine); err != nil {
		return err
	}

	// Fetch the app details from ChaosEngine instance. Check if app is present
	// Also check, if the app is annotated for chaos & that the labels are unique
	if err = getApplicationDetail(engine); err != nil {
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(appinfo derivation) Unable to get chaosengine")
		return err
	}

	if engine.Instance.Spec.AnnotationCheck == "true" {
		if engine.AppInfo.Label == "" || engine.AppInfo.Namespace == "" || engine.AppInfo.Kind == "" {
			return errors.Errorf("incomplete AppInfo inside chaosengine")
		}
		// Determine whether apps with matching labels have chaos annotation set to true
		engine, err = resource.CheckChaosAnnotation(engine, clientSet, *dynamicClient)
		if err != nil {
			// using an event msg that indicates the app couldn't be identified. By this point in execution,
			// if the engine could not be found or accessed, it would already be caught in r.initEngine & getApplicationDetail
			r.Recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "ChaosResourcesOperationFailed", "(app indentification) Unable to filter app by specified info")
			chaosTypes.Log.Info("Annotation check failed with", "error:", err)

			return err
		}
	}

	return nil
}

func (r *ChaosEngineReconciler) updateEngineForComplete(engine *chaosTypes.EngineInfo, isCompleted bool) error {
	if engine.Instance.Status.EngineStatus != litmuschaosv1alpha1.EngineStatusCompleted {
		engine.Instance.Status.EngineStatus = litmuschaosv1alpha1.EngineStatusCompleted
		engine.Instance.Spec.EngineState = litmuschaosv1alpha1.EngineStateStop
		if err := r.Client.Update(context.TODO(), engine.Instance, &client.UpdateOptions{}); err != nil {
			return fmt.Errorf("unable to update ChaosEngine Status, due to update error: %v", err)
		}
		r.Recorder.Eventf(engine.Instance, corev1.EventTypeNormal, "ChaosEngineCompleted", "ChaosEngine completed, will delete or retain the resources according to jobCleanUpPolicy")
	}

	return nil
}

func (r *ChaosEngineReconciler) updateEngineForRestart(engine *chaosTypes.EngineInfo) error {
	r.Recorder.Eventf(engine.Instance, corev1.EventTypeNormal, "RestartInProgress", "ChaosEngine is restarted")
	engine.Instance.Status.EngineStatus = litmuschaosv1alpha1.EngineStatusInitialized
	engine.Instance.Status.Experiments = nil
	if err := r.Client.Update(context.TODO(), engine.Instance, &client.UpdateOptions{}); err != nil {
		return fmt.Errorf("unable to restart ChaosEngine, due to update error: %v", err)
	}

	return nil
}

// updateChaosStatus update the chaos status inside the chaosresult
func (r *ChaosEngineReconciler) updateChaosStatus(engine *chaosTypes.EngineInfo, request reconcile.Request) error {
	if err := r.waitForChaosPodTermination(engine, request); err != nil {
		return err
	}

	// skipping CRD validation for the namespace scoped operator
	if os.Getenv("WATCH_NAMESPACE") == "" {
		found, err := isResultCRDAvailable()
		if err != nil {
			return err
		}
		if !found {
			return nil
		}
	}

	return r.updateChaosResult(engine, request)
}

// updateChaosResult update the chaosstatus and annotation inside the chaosresult
func (r *ChaosEngineReconciler) updateChaosResult(engine *chaosTypes.EngineInfo, request reconcile.Request) error {
	chaosresultList := &litmuschaosv1alpha1.ChaosResultList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
		client.MatchingLabels{},
	}

	if err := r.Client.List(context.TODO(), chaosresultList, opts...); err != nil {
		return err
	}

	for _, result := range chaosresultList.Items {
		if result.Labels["chaosUID"] == string(engine.Instance.UID) {
			targetsList, annotations := getChaosStatus(result)
			result.Status.History.Targets = targetsList
			result.ObjectMeta.Annotations = annotations

			chaosTypes.Log.Info("updating chaos status inside chaosresult", "chaosresult", result.Name)
			return r.Client.Update(context.TODO(), &result, &client.UpdateOptions{})
		}
	}

	return nil
}

// waitForChaosPodTermination wait until the termination of chaos pod after abort
func (r *ChaosEngineReconciler) waitForChaosPodTermination(engine *chaosTypes.EngineInfo, request reconcile.Request) error {
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
		client.MatchingLabels{"chaosUID": string(engine.Instance.UID)},
	}

	return retry.
		Times(uint(180)).
		Wait(1 * time.Second).
		Try(func(attempt uint) error {
			chaosPodList := &corev1.PodList{}
			if err := r.Client.List(context.TODO(), chaosPodList, opts...); err != nil {
				return err
			}
			if len(chaosPodList.Items) != 0 {
				return errors.Errorf("chaos pods are not deleted yet")
			}
			return nil
		})
}

// getChaosStatus return the target application details along with their chaos status
func getChaosStatus(result litmuschaosv1alpha1.ChaosResult) ([]litmuschaosv1alpha1.TargetDetails, map[string]string) {
	annotations := result.ObjectMeta.Annotations

	targetsList := result.Status.History.Targets
	for k, v := range annotations {
		switch strings.ToLower(v) {
		case "injected", "reverted", "targeted":
			kind := strings.TrimSpace(strings.Split(k, "/")[0])
			name := strings.TrimSpace(strings.Split(k, "/")[1])
			if !updateTargets(name, v, &targetsList) {
				targetsList = append(targetsList, litmuschaosv1alpha1.TargetDetails{
					Name:        name,
					Kind:        kind,
					ChaosStatus: v,
				})
			}
			delete(annotations, k)
		}
	}

	return targetsList, annotations
}

// isResultCRDAvailable check the existence of chaosresult CRD inside cluster
func isResultCRDAvailable() (bool, error) {

	dynamicClient, err := dynamicclientset.CreateClientSet()
	if err != nil {
		return false, err
	}

	// defining the gvr for the requested resource
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	resultList, err := (*dynamicClient).Resource(gvr).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return false, err
	}

	// check the presence of chaosresult CRD inside cluster
	for _, crd := range resultList.Items {
		if crd.GetName() == chaosTypes.ResultCRDName {
			return true, nil
		}
	}

	return false, nil
}

// updates the chaos status of targets which is already present inside history.targets
func updateTargets(name, status string, data *[]litmuschaosv1alpha1.TargetDetails) bool {
	for i := range *data {
		if (*data)[i].Name == name {
			(*data)[i].ChaosStatus = status
			return true
		}
	}

	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChaosEngineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&litmuschaosv1alpha1.ChaosEngine{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
