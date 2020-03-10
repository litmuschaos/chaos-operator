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

package chaosengine

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/litmuschaos/elves/kubernetes/container"
	"github.com/litmuschaos/elves/kubernetes/pod"
	"github.com/litmuschaos/elves/kubernetes/service"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/litmuschaos/chaos-operator/pkg/analytics"
	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	"github.com/litmuschaos/chaos-operator/pkg/controller/resource"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
	"github.com/litmuschaos/chaos-operator/pkg/controller/utils"
	"github.com/litmuschaos/chaos-operator/pkg/controller/watcher"
)

const finalizer = "chaosengine.litmuschaos.io/finalizer"

var _ reconcile.Reconciler = &ReconcileChaosEngine{}

// ReconcileChaosEngine reconciles a ChaosEngine object
type ReconcileChaosEngine struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// reconcileEngine contains details of reconcileEngine
type reconcileEngine struct {
	r         *ReconcileChaosEngine
	reqLogger logr.Logger
}

//podEngineRunner contains the information of pod
type podEngineRunner struct {
	pod, engineRunner *corev1.Pod
	*reconcileEngine
}

//serviceEngineMonitor contains informatiom of service
type serviceEngineMonitor struct {
	service, engineMonitor *corev1.Service
	*reconcileEngine
	monitoring bool
}

//podEngineMonitor contains the information of pod
type podEngineMonitor struct {
	pod, engineMonitor *corev1.Pod
	*reconcileEngine
	monitoring bool
}

// Add creates a new ChaosEngine Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileChaosEngine{client: mgr.GetClient(), scheme: mgr.GetScheme(), recorder: mgr.GetEventRecorderFor("chaos-operator")}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("chaosengine-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	err = watchChaosResources(mgr.GetClient(), c)
	if err != nil {
		return err
	}
	return nil
}

var engine *chaosTypes.EngineInfo

// watchSecondaryResources watch's for changes in chaos resources
func watchChaosResources(clientSet client.Client, c controller.Controller) error {
	// Watch for Primary Chaos Resource
	err := c.Watch(&source.Kind{Type: &litmuschaosv1alpha1.ChaosEngine{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	//Watch for Secondary Chaos Resources
	err = watcher.WatchForRunnerPod(clientSet, c)
	if err != nil {
		return err
	}
	err = watcher.WatchForMonitorPod(clientSet, c)
	if err != nil {
		return err
	}
	err = watcher.WatchForMonitorService(clientSet, c)
	if err != nil {
		return err
	}
	return nil
}

// Reconcile reads that state of the cluster for a ChaosEngine object and makes changes based on the state read
// and what is in the ChaosEngine.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileChaosEngine) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := chaosTypes.Log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ChaosEngine")
	err := r.getChaosEngineInstance(request)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// At the start of this reconcile calls, if the status of chaos engine is empty, fill it up with active
	if err := r.initEngineState(engine); err != nil {
		reqLogger.Error(err, "Unable to update EngineState in ChaosEngine Resource, due to error: %v", err)
	}

	// Check if the chaos-runner pod is completed or not
	isCompleted := r.checkRunnerPodCompleted(engine)
	// If isCompleted is true, then proceed, and verify if this is the first call of this type,
	// by verifying that this engineStatus has never been seen before
	// If thats the case, then udpate the engineStatus to completed
	if isCompleted && engine.Instance.Status.EngineStatus != litmuschaosv1alpha1.EngineStatusCompleted {
		if err := r.updateStatus(engine, litmuschaosv1alpha1.EngineStatusCompleted); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Verify that the engineStatus is set to completed,
	// if thats the case, then reconcile for completed
	if checkEngineStatusForComplete(engine) {
		return r.reconcileForComplete(request)
	}

	// Verify that the engineState is set to stop
	// if true, set the engineStatus as stopped
	if checkEngineStateForStop(engine) {
		if err := r.updateStatus(engine, litmuschaosv1alpha1.EngineStatusStopped); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Verify that the engineStatus is set to stopped,
	// then reconcile for delete
	if checkEngineStatusForStopped(engine) {
		return r.reconcileForDelete(request)
	}

	// Verify that the engineState, and engineStatus to initalized chaos engine resources
	if checkEngineForCreation(engine) {
		if err = r.checkRunnerPodForCompletion(engine, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
		return r.reconcileForCreationAndRunning(engine, request)
	}

	return reconcile.Result{}, nil
}

//MonitorServiceAndPod checks if the EngineMonitorPod And EngineMonitorService already exist or not
func MonitorServiceAndPod(monitorService *serviceEngineMonitor, monitorPod *podEngineMonitor) error {
	// Check if the EngineMonitorService already exists, else create
	if err := engineMonitorService(monitorService); err != nil {
		return err
	}
	// Check if the EngineMonitorPod already exists, else create
	if err := engineMonitorPod(monitorPod); err != nil {
		return err
	}
	return nil
}

// Creates engineMonitor pod and engineMonitor Service
// Also reconciles those resources
func createMonitoringResources(engine chaosTypes.EngineInfo, recEngine *reconcileEngine) (reconcile.Result, error) {

	// Define the engine-monitor service which is secondary-resource #3
	engineMonitorSvc, err := newMonitorServiceForCR(engine)
	if err != nil {
		return reconcile.Result{}, err
	}
	//Define an engine-monitor pod which is secondary-resource #2
	engineMonitorPod, err := newMonitorPodForCR(engine)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Creates an object of monitorService
	monitorService := &serviceEngineMonitor{
		service:         &corev1.Service{},
		engineMonitor:   engineMonitorSvc,
		reconcileEngine: recEngine,
		monitoring:      engine.Instance.Spec.Monitoring,
	}
	// Creates an object of monitorPod
	monitorPod := &podEngineMonitor{
		pod:             &corev1.Pod{},
		engineMonitor:   engineMonitorPod,
		reconcileEngine: recEngine,
		monitoring:      engine.Instance.Spec.Monitoring,
	}

	// Check if the engineMonitorService already exists, else create
	if err = MonitorServiceAndPod(monitorService, monitorPod); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

// getChaosRunnerENV return the env required for chaos-runner
func getChaosRunnerENV(cr *litmuschaosv1alpha1.ChaosEngine, aExList []string, ClientUUID string) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "CHAOSENGINE",
			Value: cr.Name,
		},
		{
			Name:  "APP_LABEL",
			Value: cr.Spec.Appinfo.Applabel,
		},
		{
			Name:  "APP_NAMESPACE",
			Value: cr.Spec.Appinfo.Appns,
		},
		{
			Name:  "EXPERIMENT_LIST",
			Value: fmt.Sprint(strings.Join(aExList, ",")),
		},
		{
			Name:  "CHAOS_SVC_ACC",
			Value: cr.Spec.ChaosServiceAccount,
		},
		{
			Name:  "AUXILIARY_APPINFO",
			Value: cr.Spec.AuxiliaryAppInfo,
		},
		{
			Name:  "CLIENT_UUID",
			Value: ClientUUID,
		},
	}
}

// getChaosMonitorENV return the env required for chaos-Monitor
func getChaosMonitorENV(cr *litmuschaosv1alpha1.ChaosEngine, aUUID types.UID) []corev1.EnvVar {
	return []corev1.EnvVar{
		{Name: "CHAOSENGINE", Value: cr.Name},
		{Name: "APP_UUID", Value: string(aUUID)},
		{Name: "APP_NAMESPACE", Value: cr.Spec.Appinfo.Appns},
	}
}

// newRunnerPodForCR defines secondary resource #1 in same namespace as CR
func newRunnerPodForCR(ce chaosTypes.EngineInfo) (*corev1.Pod, error) {
	if (len(ce.AppExperiments) == 0 || ce.AppUUID == "") && ce.Instance.Spec.AnnotationCheck == "true" {
		return nil, errors.New("application experiment list or UUID is empty")
	}
	//Initiate the Engine Info, with the type of runner to be used
	if ce.Instance.Spec.Components.Runner.Type == "ansible" {
		return newAnsibleRunnerPodForCR(ce)
	}
	return newGoRunnerPodForCR(ce)
}

func newGoRunnerPodForCR(engine chaosTypes.EngineInfo) (*corev1.Pod, error) {
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

	if engine.Instance.Spec.Components.Runner.Command != nil {
		containerForRunner.WithCommandNew(engine.Instance.Spec.Components.Runner.Command)
	}

	return pod.NewBuilder().
		WithName(engine.Instance.Name + "-runner").
		WithNamespace(engine.Instance.Namespace).
		WithLabels(map[string]string{"app": engine.Instance.Name, "chaosUID": string(engine.Instance.UID)}).
		WithServiceAccountName(engine.Instance.Spec.ChaosServiceAccount).
		WithRestartPolicy("OnFailure").
		WithContainerBuilder(containerForRunner).Build()
}

func newAnsibleRunnerPodForCR(engine chaosTypes.EngineInfo) (*corev1.Pod, error) {
	containerForRunner := container.NewBuilder().
		WithName("chaos-runner").
		WithImage(engine.Instance.Spec.Components.Runner.Image).
		WithImagePullPolicy(corev1.PullIfNotPresent).
		WithCommandNew([]string{"/bin/bash"}).
		WithArgumentsNew([]string{"-c", "ansible-playbook ./executor/test.yml -i /etc/ansible/hosts; exit 0"}).
		WithEnvsNew(getChaosRunnerENV(engine.Instance, engine.AppExperiments, analytics.ClientUUID))

	if engine.Instance.Spec.Components.Runner.ImagePullPolicy != "" {
		containerForRunner.WithImagePullPolicy(engine.Instance.Spec.Components.Runner.ImagePullPolicy)
	}

	if engine.Instance.Spec.Components.Runner.Args != nil {
		containerForRunner.WithArgumentsNew(engine.Instance.Spec.Components.Runner.Args)
	}

	if engine.Instance.Spec.Components.Runner.Command != nil {
		containerForRunner.WithCommandNew(engine.Instance.Spec.Components.Runner.Command)
	}

	return pod.NewBuilder().
		WithName(engine.Instance.Name + "-runner").
		WithLabels(map[string]string{"app": engine.Instance.Name, "chaosUID": string(engine.Instance.UID)}).
		WithNamespace(engine.Instance.Namespace).
		WithRestartPolicy("OnFailure").
		WithServiceAccountName(engine.Instance.Spec.ChaosServiceAccount).
		WithContainerBuilder(containerForRunner).
		Build()
}

// newMonitorPodForCR defines secondary resource #2 in same namespace as CR */
func newMonitorPodForCR(engine chaosTypes.EngineInfo) (*corev1.Pod, error) {
	if engine.Instance == nil {
		return nil, errors.New("chaosengine got nil")
	}
	labels := map[string]string{
		"app":      engine.Instance.Name,
		"chaosUID": string(engine.Instance.UID),
	}
	monitorPod, err := pod.NewBuilder().
		WithName(engine.Instance.Name + "-monitor").
		WithNamespace(engine.Instance.Namespace).
		WithLabels(labels).
		WithServiceAccountName(engine.Instance.Spec.ChaosServiceAccount).
		WithRestartPolicy("OnFailure").
		WithContainerBuilder(
			container.NewBuilder().
				WithName("chaos-monitor").
				WithImage(engine.Instance.Spec.Components.Monitor.Image).
				WithPortsNew([]corev1.ContainerPort{{ContainerPort: 8080, Protocol: "TCP", Name: "metrics"}}).
				WithEnvsNew(getChaosMonitorENV(engine.Instance, engine.AppUUID))).Build()
	if err != nil {
		return nil, err
	}
	return monitorPod, nil
}

// newMonitorServiceForCR defines secondary resource #2 in same namespace as CR */
func newMonitorServiceForCR(engine chaosTypes.EngineInfo) (*corev1.Service, error) {
	if engine.Instance == nil {
		return nil, errors.New("nil chaosengine object")
	}
	serviceObj, err := service.NewBuilder().
		WithName(engine.Instance.Name + "-monitor").
		WithNamespace(engine.Instance.Namespace).
		WithLabels(map[string]string{"app": "chaos-exporter", "chaosUID": string(engine.Instance.UID)}).
		WithPorts([]corev1.ServicePort{{Name: "metrics", Port: 8080}}).
		WithSelectorsNew(map[string]string{"monitorFor": engine.Instance.Name}).Build()
	if err != nil {
		return nil, err
	}
	return serviceObj, nil
}

// initializeApplicationInfo to initialize application info
func initializeApplicationInfo(instance *litmuschaosv1alpha1.ChaosEngine, appInfo *chaosTypes.ApplicationInfo) (*chaosTypes.ApplicationInfo, error) {
	if instance == nil {
		return nil, errors.New("empty chaosengine")
	}
	appLabel := strings.Split(instance.Spec.Appinfo.Applabel, "=")
	chaosTypes.AppLabelKey = appLabel[0]
	chaosTypes.AppLabelValue = appLabel[1]
	appInfo.Label = make(map[string]string)
	appInfo.Label[chaosTypes.AppLabelKey] = chaosTypes.AppLabelValue
	appInfo.Namespace = instance.Spec.Appinfo.Appns
	appInfo.ExperimentList = instance.Spec.Experiments
	appInfo.ServiceAccountName = instance.Spec.ChaosServiceAccount
	appInfo.Kind = instance.Spec.Appinfo.AppKind
	return appInfo, nil
}

// engineRunnerPod to Check if the engineRunner pod already exists, else create
func engineRunnerPod(runnerPod *podEngineRunner) error {
	err := runnerPod.r.client.Get(context.TODO(), types.NamespacedName{Name: runnerPod.engineRunner.Name, Namespace: runnerPod.engineRunner.Namespace}, runnerPod.pod)
	if err != nil && k8serrors.IsNotFound(err) {
		runnerPod.reqLogger.Info("Creating a new engineRunner Pod", "Pod.Namespace", runnerPod.engineRunner.Namespace, "Pod.Name", runnerPod.engineRunner.Name)
		err = runnerPod.r.client.Create(context.TODO(), runnerPod.engineRunner)
		if err != nil {
			return err
		}

		// Pod created successfully - don't requeue
		runnerPod.reqLogger.Info("engineRunner Pod created successfully")
	} else if err != nil {
		return err
	}
	runnerPod.reqLogger.Info("Skip reconcile: engineRunner Pod already exists", "Pod.Namespace", runnerPod.pod.Namespace, "Pod.Name", runnerPod.pod.Name)
	return nil
}

// Check if the engineMonitorService already exists, else create
func engineMonitorService(monitorService *serviceEngineMonitor) error {
	err := monitorService.r.client.Get(context.TODO(), types.NamespacedName{Name: monitorService.engineMonitor.Name, Namespace: monitorService.engineMonitor.Namespace}, monitorService.service)
	if err != nil && k8serrors.IsNotFound(err) {
		monitorService.reqLogger.Info("Creating a new engineMonitor Service", "Service.Namespace", monitorService.engineMonitor.Namespace, "Service.Name", monitorService.engineMonitor.Name)
		err = monitorService.r.client.Create(context.TODO(), monitorService.engineMonitor)
		if err != nil {
			return err
		}

		// Service created successfully - don't requeue
	} else if err != nil {
		return err
	}
	monitorService.reqLogger.Info("Skip reconcile: engineMonitor Service already exists", "Service.Namespace", monitorService.engineMonitor.Namespace, "Service.Name", monitorService.engineMonitor.Name)
	return nil /*You can return now, both sec resources are existing */
}

// engineMonitorPod to Check if the engineMonitor Pod is already exists, else create
func engineMonitorPod(monitorPod *podEngineMonitor) error {
	coreV1Pod := &corev1.Pod{}
	err := monitorPod.r.client.Get(context.TODO(), types.NamespacedName{Name: monitorPod.engineMonitor.Name, Namespace: monitorPod.engineMonitor.Namespace}, coreV1Pod)
	if err != nil && k8serrors.IsNotFound(err) {
		monitorPod.reqLogger.Info("Creating a new engineMonitor Pod", "Pod.Namespace", monitorPod.engineMonitor.Namespace, "Pod.Name", monitorPod.engineMonitor.Name)
		if err = monitorPod.r.client.Create(context.TODO(), monitorPod.engineMonitor); err != nil {
			return err
		}

		monitorPod.reqLogger.Info("engineMonitor Pod created successfully")
	} else if err != nil {
		return err
	}
	monitorPod.reqLogger.Info("Skip reconcile: engineMonitor Pod already exists", "Pod.Namespace", coreV1Pod.Namespace, "Pod.Name", coreV1Pod.Name)
	return nil
}

// Fetch the ChaosEngine instance
func (r *ReconcileChaosEngine) getChaosEngineInstance(request reconcile.Request) error {
	instance := &litmuschaosv1alpha1.ChaosEngine{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		// Error reading the object - requeue the request.
		return err
	}
	engine = &chaosTypes.EngineInfo{
		Instance: instance,
	}
	return nil
}

// Get application details
func getApplicationDetail(ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	applicationInfo := &chaosTypes.ApplicationInfo{}
	appInfo, err := initializeApplicationInfo(ce.Instance, applicationInfo)
	if err != nil {
		return ce, err
	}
	ce.AppInfo = appInfo

	var appExperiments []string
	for _, exp := range appInfo.ExperimentList {
		appExperiments = append(appExperiments, exp.Name)
	}
	ce.AppExperiments = appExperiments

	chaosTypes.Log.Info("App key derived from chaosengine is ", "appLabelKey", chaosTypes.AppLabelKey)
	chaosTypes.Log.Info("App Label derived from Chaosengine is ", "appLabelValue", chaosTypes.AppLabelValue)
	chaosTypes.Log.Info("App NS derived from Chaosengine is ", "appNamespace", appInfo.Namespace)
	chaosTypes.Log.Info("Exp list derived from chaosengine is ", "appExpirements", appExperiments)
	chaosTypes.Log.Info("Monitoring Status derived from chaosengine is", "monitoringStatus", ce.Instance.Spec.Monitoring)
	chaosTypes.Log.Info("Runner image derived from chaosengine is", "runnerImage", ce.Instance.Spec.Components.Runner.Image)
	chaosTypes.Log.Info("exporter image derived from chaosengine is", "exporterImage", ce.Instance.Spec.Components.Monitor.Image)
	chaosTypes.Log.Info("Annotation check is ", "annotationCheck", ce.Instance.Spec.AnnotationCheck)
	return ce, nil
}

// Check if the engineRunner pod already exists, else create
func (r *ReconcileChaosEngine) checkEngineRunnerPod(reqLogger logr.Logger, engineRunner *corev1.Pod) (*reconcileEngine, error) {
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

	if err := engineRunnerPod(runnerPod); err != nil {
		return engineReconcile, err
	}
	return engineReconcile, nil
}

// check monitoring status
func checkMonitoring(engineReconcile *reconcileEngine, reqLogger logr.Logger) (reconcile.Result, error) {
	if engine.Instance.Spec.Monitoring {
		reconcileResult, err := createMonitoringResources(*engine, engineReconcile)
		if err != nil {
			return reconcileResult, err
		}
	} else {
		reqLogger.Info("Monitoring is disabled")
	}
	return reconcile.Result{}, nil
}

//setChaosResourceImage take the runner and monitor image from engine spec
//if it is not there then it will take from chaos-operator env
//at last if it is not able to find image in engine spec and operator env then it will take default images
func setChaosResourceImage() {

	ChaosMonitorImage := os.Getenv("CHAOS_MONITOR_IMAGE")
	ChaosRunnerImage := os.Getenv("CHAOS_RUNNER_IMAGE")

	if engine.Instance.Spec.Components.Monitor.Image == "" && ChaosMonitorImage == "" {
		engine.Instance.Spec.Components.Monitor.Image = chaosTypes.DefaultChaosMonitorImage
	} else if engine.Instance.Spec.Components.Monitor.Image == "" {
		engine.Instance.Spec.Components.Monitor.Image = ChaosMonitorImage
	}
	if engine.Instance.Spec.Components.Runner.Image == "" && ChaosRunnerImage == "" {
		engine.Instance.Spec.Components.Runner.Image = chaosTypes.DefaultChaosRunnerImage
	} else if engine.Instance.Spec.Components.Runner.Image == "" {
		engine.Instance.Spec.Components.Runner.Image = ChaosRunnerImage
	}
}

func getAnnotationCheck() error {

	if engine.Instance.Spec.AnnotationCheck == "" {
		engine.Instance.Spec.AnnotationCheck = chaosTypes.DefaultAnnotationCheck

	}
	if engine.Instance.Spec.AnnotationCheck != "true" && engine.Instance.Spec.AnnotationCheck != "false" {
		return fmt.Errorf("annotationCheck '%s', is not supported it should be true or false", engine.Instance.Spec.AnnotationCheck)
	}
	return nil
}

// reconcileForDelete
func (r *ReconcileChaosEngine) reconcileForDelete(request reconcile.Request) (reconcile.Result, error) {
	reconcileResult, err := r.removeChaosResources(engine, request)
	if err != nil {
		return reconcileResult, err
	}
	err = r.removeChaosServices(engine, request)
	if err != nil {
		return reconcile.Result{}, err
	}
	opts := client.UpdateOptions{}

	if engine.Instance.ObjectMeta.Finalizers != nil {
		engine.Instance.ObjectMeta.Finalizers = utils.RemoveString(engine.Instance.ObjectMeta.Finalizers, "chaosengine.litmuschaos.io/finalizer")
		r.recorder.Eventf(engine.Instance, corev1.EventTypeNormal, "ChaosEngineStopped", " Chaos resources deleted successfully")
	}
	if err := r.client.Update(context.TODO(), engine.Instance, &opts); err != nil {
		return reconcile.Result{}, fmt.Errorf("Unable to remove Finalizer from chaosEngine Resource, due to error: %v", err)
	}
	return reconcile.Result{}, nil

}

func (r *ReconcileChaosEngine) removeChaosServices(engine *chaosTypes.EngineInfo, request reconcile.Request) error {
	optsList := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
		client.MatchingLabels{"app": "chaos-exporter", "chaosUID": string(engine.Instance.UID)},
	}
	var serviceList corev1.ServiceList
	if errList := r.client.List(context.TODO(), &serviceList, optsList...); errList != nil {
		return errList
	}
	for _, v := range serviceList.Items {
		if errDel := r.client.Delete(context.TODO(), &v, []client.DeleteOption{}...); errDel != nil {
			return errDel
		}
	}
	return nil
}

func (r *ReconcileChaosEngine) removeChaosResources(engine *chaosTypes.EngineInfo, request reconcile.Request) (reconcile.Result, error) {
	optsDelete := []client.DeleteAllOfOption{
		client.InNamespace(request.NamespacedName.Namespace),
		client.MatchingLabels{"chaosUID": string(engine.Instance.UID)},
		client.PropagationPolicy(metav1.DeletePropagationBackground),
		client.GracePeriodSeconds(int64(0)),
	}
	var deleteEvent []string
	var err []error

	if errDeployment := r.client.DeleteAllOf(context.TODO(), &appsv1.Deployment{}, optsDelete...); errDeployment != nil {
		err = append(err, errDeployment)
		deleteEvent = append(deleteEvent, "Deployments, ")
	}

	if errDaemonSet := r.client.DeleteAllOf(context.TODO(), &appsv1.DaemonSet{}, optsDelete...); errDaemonSet != nil {
		err = append(err, errDaemonSet)
		deleteEvent = append(deleteEvent, "DaemonSets, ")
	}

	if errJob := r.client.DeleteAllOf(context.TODO(), &batchv1.Job{}, optsDelete...); errJob != nil {
		err = append(err, errJob)
		deleteEvent = append(deleteEvent, "Jobs, ")
	}

	if errPod := r.client.DeleteAllOf(context.TODO(), &corev1.Pod{}, optsDelete...); errPod != nil {
		err = append(err, errPod)
		deleteEvent = append(deleteEvent, "Pods, ")
	}
	if err != nil {
		r.recorder.Eventf(engine.Instance, corev1.EventTypeWarning, "Deletion Failed", "Unable to delete chaos resources: %v allocated to ChaosEngine: %v in Namespace: %v", strings.Join(deleteEvent, ""), engine.Instance.Name, engine.Instance.Namespace)
		return reconcile.Result{}, fmt.Errorf("Unable to delete ChaosResources due to %v", err)
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileChaosEngine) addFinalizerToEngine(engine *chaosTypes.EngineInfo, finalizer string) error {
	optsUpdate := client.UpdateOptions{}
	if engine.Instance.ObjectMeta.Finalizers == nil {
		engine.Instance.ObjectMeta.Finalizers = append(engine.Instance.ObjectMeta.Finalizers, finalizer)
		r.recorder.Eventf(engine.Instance, corev1.EventTypeNormal, "ChaosEngineInitialized", "%s created successfully", engine.Instance.Name+"-runner")
	}
	err := r.client.Update(context.TODO(), engine.Instance, &optsUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileChaosEngine) updateStatus(engine *chaosTypes.EngineInfo, status litmuschaosv1alpha1.EngineStatus) error {
	opts := client.UpdateOptions{}
	engine.Instance.Status.EngineStatus = status
	return r.client.Update(context.TODO(), engine.Instance, &opts)
}

func (r *ReconcileChaosEngine) updateState(engine *chaosTypes.EngineInfo, state litmuschaosv1alpha1.EngineState) error {
	opts := client.UpdateOptions{}
	engine.Instance.Spec.EngineState = state
	return r.client.Update(context.TODO(), engine.Instance, &opts)
}

func checkEngineStateForStop(engine *chaosTypes.EngineInfo) bool {
	deletetimeStamp := engine.Instance.ObjectMeta.GetDeletionTimestamp()
	if engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateStop ||
		deletetimeStamp != nil {
		return true
	}
	return false
}

func (r *ReconcileChaosEngine) checkRunnerPodCompleted(engine *chaosTypes.EngineInfo) bool {
	runnerPod := corev1.Pod{}
	r.client.Get(context.TODO(), types.NamespacedName{Name: engine.Instance.Name + "-runner", Namespace: engine.Instance.Namespace}, &runnerPod)
	return runnerPod.Status.Phase == corev1.PodSucceeded
}

func (r *ReconcileChaosEngine) removeDefaultChaosResources(request reconcile.Request) (reconcile.Result, error) {
	if err := r.removeChaosRunner(engine, request); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.removeChaosServices(engine, request); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileChaosEngine) removeChaosRunner(engine *chaosTypes.EngineInfo, request reconcile.Request) error {
	optsList := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
		client.MatchingLabels{"app": engine.Instance.Name, "chaosUID": string(engine.Instance.UID)},
	}
	var podList corev1.PodList
	if errList := r.client.List(context.TODO(), &podList, optsList...); errList != nil {
		return errList
	}
	for _, v := range podList.Items {
		if errDel := r.client.Delete(context.TODO(), &v, []client.DeleteOption{}...); errDel != nil {
			return errDel
		}
	}
	return nil
}

func (r *ReconcileChaosEngine) checkRunnerPodForCompletion(engine *chaosTypes.EngineInfo, reqLogger logr.Logger) error {
	isCompleted := r.checkRunnerPodCompleted(engine)
	if isCompleted {
		err := r.updateStatus(engine, litmuschaosv1alpha1.EngineStatusCompleted)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkEngineStatusForComplete(engine *chaosTypes.EngineInfo) bool {
	return engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusCompleted

}
func (r *ReconcileChaosEngine) reconcileForComplete(request reconcile.Request) (reconcile.Result, error) {
	_, err := r.removeDefaultChaosResources(request)
	if err != nil {
		return reconcile.Result{}, err
	}
	if engine.Instance.ObjectMeta.Finalizers != nil {
		engine.Instance.ObjectMeta.Finalizers = utils.RemoveString(engine.Instance.ObjectMeta.Finalizers, "chaosengine.litmuschaos.io/finalizer")
		r.recorder.Eventf(engine.Instance, corev1.EventTypeNormal, "ChaosEngine Stopped", "Removing all experiment resources allocated to ChaosEngine: %v in Namespace: %v", engine.Instance.Name, engine.Instance.Namespace)
	}
	err = r.updateState(engine, litmuschaosv1alpha1.EngineStateStop)
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileChaosEngine) initEngineState(engine *chaosTypes.EngineInfo) error {
	if engine.Instance.Spec.EngineState == "" {
		err := r.updateState(engine, litmuschaosv1alpha1.EngineStateActive)
		return err
	}
	if err := r.initEngineStatus(engine); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileChaosEngine) initEngineStatus(engine *chaosTypes.EngineInfo) error {
	if engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateActive && engine.Instance.Status.EngineStatus != litmuschaosv1alpha1.EngineStatusCompleted {

		err := r.updateStatus(engine, litmuschaosv1alpha1.EngineStatusInitialized)
		if err != nil {
			return err
		}

		err = r.addFinalizerToEngine(engine, finalizer)
		if err != nil {
			return err
		}

	}
	return nil
}

func checkEngineStatusForStopped(engine *chaosTypes.EngineInfo) bool {
	return (engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusStopped && engine.Instance.Status.EngineStatus != litmuschaosv1alpha1.EngineStatusCompleted)

}
func checkEngineForCreation(engine *chaosTypes.EngineInfo) bool {
	return engine.Instance.Status.EngineStatus == litmuschaosv1alpha1.EngineStatusInitialized && engine.Instance.Spec.EngineState == litmuschaosv1alpha1.EngineStateActive
}

func (r *ReconcileChaosEngine) reconcileForCreationAndRunning(engine *chaosTypes.EngineInfo, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := chaosTypes.Log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ChaosEngine")
	// Get the image for runner and monitor pod from chaosengine spec,operator env or default values.
	setChaosResourceImage()

	//getAnnotationCheck fetch the annotationCheck from engine spec
	err := getAnnotationCheck()
	if err != nil {
		return reconcile.Result{}, err
	}

	// Fetch the app details from ChaosEngine instance. Check if app is present
	// Also check, if the app is annotated for chaos & that the labels are unique
	engine, err = getApplicationDetail(engine)
	if err != nil {
		return reconcile.Result{}, err
	}

	if engine.Instance.Spec.AnnotationCheck == "true" {
		// Determine whether apps with matching labels have chaos annotation set to true
		engine, err = resource.CheckChaosAnnotation(engine)
		if err != nil {
			chaosTypes.Log.Info("Annotation check failed with", "error:", err)
			return reconcile.Result{}, nil
		}
	}
	// Define an engineRunner pod which is secondary-resource #1
	engineRunner, err := newRunnerPodForCR(*engine)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Check if the engineRunner pod already exists, else create
	engineReconcile, err := r.checkEngineRunnerPod(reqLogger, engineRunner)
	if err != nil {
		return reconcile.Result{}, err
	}
	// If monitoring is set to true,
	// Define an engineMonitor pod which is secondary-resource #2 and
	// Define an engineMonitor service which is secondary-resource #3
	// in the same namespace as CR
	reconcileResult, err := checkMonitoring(engineReconcile, reqLogger)
	if err != nil {
		return reconcileResult, err
	}
	return reconcile.Result{}, nil
}
