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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/litmuschaos/chaos-operator/pkg/analytics"
	"github.com/litmuschaos/kube-helper/kubernetes/container"
	"github.com/litmuschaos/kube-helper/kubernetes/pod"
	"github.com/litmuschaos/kube-helper/kubernetes/service"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	"github.com/litmuschaos/chaos-operator/pkg/controller/resource"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

var _ reconcile.Reconciler = &ReconcileChaosEngine{}

// ReconcileChaosEngine reconciles a ChaosEngine object
type ReconcileChaosEngine struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
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
	return &ReconcileChaosEngine{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("chaosengine-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	handlerForOwner := handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &litmuschaosv1alpha1.ChaosEngine{},
	}
	err = watchChaosResources(handlerForOwner, c)
	if err != nil {
		return err
	}
	return nil
}

var engine *chaosTypes.EngineInfo

// watchSecondaryResources watch's for changes in chaos resources
func watchChaosResources(handlerForOwner handler.EnqueueRequestForOwner, c controller.Controller) error {
	// Watch for Primary Resource
	err := c.Watch(&source.Kind{Type: &litmuschaosv1alpha1.ChaosEngine{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for Secondary Resources
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handlerForOwner)
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handlerForOwner)
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
	// Fetch the ChaosEngine instance
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

	// Get the image for runner and monitor pod from chaosengine spec,operator env or default values.
	setChaosResourceImage()

	//getChaosType fetch the chaosType from engine spec
	err = getChaosType()
	if err != nil {
		return reconcile.Result{}, err
	}

	// Fetch the app details from ChaosEngine instance. Check if app is present
	// Also check, if the app is annotated for chaos & that the labels are unique
	engine, err := getApplicationDetail(engine)
	if err != nil {
		return reconcile.Result{}, err
	}

	if engine.Instance.Spec.ChaosType == "app" {
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
	// Set ChaosEngine instance as the owner and controller of engine-runner pod
	if err := controllerutil.SetControllerReference(engine.Instance, engineRunner, r.scheme); err != nil {
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

// Set ChaosEngine instance as the owner and controller of engine-Monitor service
func setControllerReference(recEngine *reconcileEngine, engineMonitor *corev1.Pod, engineMonitorSvc *corev1.Service) error {
	if err := controllerutil.SetControllerReference(engine.Instance, engineMonitor, recEngine.r.scheme); err != nil {
		return err
	}
	if err := controllerutil.SetControllerReference(engine.Instance, engineMonitorSvc, recEngine.r.scheme); err != nil {
		return err
	}
	return nil
}

//MonitorServiceAndPod checks if the EngineMonitorPod And EngineMonitorService already exist or not
func MonitorServiceAndPod(monitorService *serviceEngineMonitor, monitorPod *podEngineMonitor) error {
	err := engineMonitorService(monitorService)
	if err != nil {
		return err
	}
	// Check if the EngineMonitorPod already exists, else create
	if err = engineMonitorPod(monitorPod); err != nil {
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
	// Set ChaosEngine instance as the owner and controller of engine-Monitor pod
	if err = setControllerReference(recEngine, engineMonitorPod, engineMonitorSvc); err != nil {
		return reconcile.Result{}, err
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
	if (len(ce.AppExperiments) == 0 || ce.AppUUID == "") && ce.Instance.Spec.ChaosType == "app" {
		return nil, errors.New("application experiment list or UUID is empty")
	}
	//Initiate the Engine Info, with the type of executor to be used
	if ce.Instance.Spec.Components.Runner.Type == "go" {
		return newGoRunnerPodForCR(ce)
	}
	return newAnsibleRunnerPodForCR(ce)
}

func newGoRunnerPodForCR(ce chaosTypes.EngineInfo) (*corev1.Pod, error) {
	return pod.NewBuilder().
		WithName(ce.Instance.Name + "-runner").
		WithNamespace(ce.Instance.Namespace).
		WithLabels(map[string]string{"app": ce.Instance.Name}).
		WithServiceAccountName(ce.Instance.Spec.ChaosServiceAccount).
		WithRestartPolicy("OnFailure").
		WithContainerBuilder(
			container.NewBuilder().
				WithEnvsNew(getChaosRunnerENV(ce.Instance, ce.AppExperiments, analytics.ClientUUID)).
				WithName("chaos-runner").
				WithImage(ce.Instance.Spec.Components.Runner.Image).
				WithImagePullPolicy(corev1.PullIfNotPresent)).Build()
}

func newAnsibleRunnerPodForCR(ce chaosTypes.EngineInfo) (*corev1.Pod, error) {
	return pod.NewBuilder().
		WithName(ce.Instance.Name + "-runner").
		WithLabels(map[string]string{"app": ce.Instance.Name}).
		WithNamespace(ce.Instance.Namespace).
		WithRestartPolicy("OnFailure").
		WithServiceAccountName(ce.Instance.Spec.ChaosServiceAccount).
		WithContainerBuilder(
			container.NewBuilder().
				WithName("chaos-runner").
				WithImage(ce.Instance.Spec.Components.Runner.Image).
				WithImagePullPolicy(corev1.PullIfNotPresent).
				WithCommandNew([]string{"/bin/bash"}).
				WithArgumentsNew([]string{"-c", "ansible-playbook ./executor/test.yml -i /etc/ansible/hosts; exit 0"}).
				WithEnvsNew(getChaosRunnerENV(ce.Instance, ce.AppExperiments, analytics.ClientUUID))).Build()
}

// newMonitorPodForCR defines secondary resource #2 in same namespace as CR */
func newMonitorPodForCR(engine chaosTypes.EngineInfo) (*corev1.Pod, error) {
	if engine.Instance == nil {
		return nil, errors.New("chaosengine got nil")
	}
	labels := map[string]string{
		"app":        "chaos-exporter",
		"monitorFor": engine.Instance.Name,
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
		WithLabels(map[string]string{"app": "chaos-exporter"}).
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
	chaosTypes.Log.Info("chaos type is ", "chaostype", ce.Instance.Spec.ChaosType)
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

func getChaosType() error {

	if engine.Instance.Spec.ChaosType == "" {
		engine.Instance.Spec.ChaosType = chaosTypes.DefaultChaosType

	}
	if engine.Instance.Spec.ChaosType != "app" && engine.Instance.Spec.ChaosType != "infra" {
		return fmt.Errorf("chaos type '%s', is not supported it should be app or infra", engine.Instance.Spec.ChaosType)
	}
	return nil
}
