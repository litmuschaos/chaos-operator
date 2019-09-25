package chaosengine

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/litmuschaos/kube-helper/kubernetes/container"
	"github.com/litmuschaos/kube-helper/kubernetes/pod"
	"github.com/litmuschaos/kube-helper/kubernetes/service"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
)

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
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ChaosEngine")

	// Fetch the ChaosEngine instance
	instance := &litmuschaosv1alpha1.ChaosEngine{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	// Fetch the app details from ChaosEngine instance. Check if app is present
	// Also check, if the app is annotated for chaos & that the labels are unique

	// TODO: Get app kind from chaosengine spec as well. Using "deploy" for now
	// TODO: Freeze label format in chaosengine( "=" as a const)
	appInfo := &applicationInfo{}
	appInfo, err = appInfo.initializeApplicationInfo(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	var appExperiments []string
	for _, exp := range appInfo.experimentList {
		appExperiments = append(appExperiments, exp.Name)
	}

	log.Info("App key derived from chaosengine is ", "appLabelKey", appLabelKey)
	log.Info("App Label derived from Chaosengine is ", "appLabelValue", appLabelValue)
	log.Info("App NS derived from Chaosengine is ", "appNamespace", appInfo.namespace)
	log.Info("Exp list derived from chaosengine is ", "appExpirements", appExperiments)

	// Use client-Go to obtain a list of apps w/ specified labels
	restConfig, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get rest kube config")
		return reconcile.Result{}, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Error(err, "unable to create clientset using restconfig")
		return reconcile.Result{}, err
	}

	chaosAppList, err := clientset.AppsV1().Deployments(appInfo.namespace).List(metav1.ListOptions{LabelSelector: instance.Spec.Appinfo.Applabel, FieldSelector: ""})
	if err != nil {
		log.Error(err, "unable to list apps matching labels")
		return reconcile.Result{}, err
	}

	var appName string
	var appUUID types.UID

	// Determine whether apps with matching labels have chaos annotation set to true
	chaosCandidates := 0
	if len(chaosAppList.Items) > 0 {
		for _, app := range chaosAppList.Items {
			appName = app.ObjectMeta.Name
			appUUID = app.ObjectMeta.UID
			appCaSts := metav1.HasAnnotation(app.ObjectMeta, chaosAnnotation)
			if appCaSts {
				//Checks if the annotation is "true" / "false"
				var annotationFlag bool
				annotationFlag, err = strconv.ParseBool(app.ObjectMeta.GetAnnotations()[chaosAnnotation])
				//log.Info("Annotation Flag", "aflag", annotationFlag)
				if err != nil {
					// Unable to check the annotation
					// Would not add in the chaosCandidates
					log.Info("Unable to check the annotationFlag", "annotationFlag", annotationFlag)
				} else {
					if annotationFlag {
						// If annotationFlag is true
						// Add it to the Chaos Candidates, and log the details
						log.Info("chaos candidate : ", "appName", appName, "appUUID", appUUID)
						chaosCandidates++
					}
				}
			}
		}
		if chaosCandidates == 0 {
			log.Info("No chaos candidates found")
			return reconcile.Result{}, nil

		} else if chaosCandidates > 1 {
			log.Info("Too many chaos candidates with same label, either provide unique labels or annotate only desired app for chaos")
			return reconcile.Result{}, nil
		}
	} else {
		log.Info("No app deployments with matching labels")
		return reconcile.Result{}, nil
	}

	// Define an engineRunner pod which is secondary-resource #1
	engineRunner, err := newRunnerPodForCR(instance, appUUID, appExperiments)
	if err != nil {
		return reconcile.Result{}, err
	}
	// Define the engine-monitor service which is secondary-resource #2
	engineMonitor, err := newMonitorServiceForCR(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	// Define an engine-exporter pod which is secondary-resource #3
	engineExporter, err := newExporterPodForCR(instance, appUUID)
	if err != nil {
		return reconcile.Result{}, err
	}
	// Set ChaosEngine instance as the owner and controller of engine-runner pod
	if err := controllerutil.SetControllerReference(instance, engineRunner, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Set ChaosEngine instance as the owner and controller of engine-exporter pod
	if err := controllerutil.SetControllerReference(instance, engineExporter, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Set ChaosEngine instance as the owner and controller of engine-monitor service
	if err := controllerutil.SetControllerReference(instance, engineMonitor, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	engineReconcile := &reconcileEngine{
		r:         r,
		reqLogger: reqLogger,
	}

	runnerPod := &podEngineRunner{
		pod:             &corev1.Pod{},
		engineRunner:    engineRunner,
		reconcileEngine: engineReconcile,
	}

	monitorService := &serviceEngineMonitor{
		service:         &corev1.Service{},
		engineMonitor:   engineMonitor,
		reconcileEngine: engineReconcile,
	}

	// Check if the engineRunner pod already exists, else create
	err = engineRunnerPod(runnerPod)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the engineMonitorService already exists, else create
	err = engineMonitorService(monitorService)
	if err != nil {
		return reconcile.Result{}, err
	}
	// Check if the EngineExporterPod already exists, else create
	err = engineExporterPod(engineExporter, r, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// getChaosRunnerENV return the env required for chaos-runner
func getChaosRunnerENV(cr *litmuschaosv1alpha1.ChaosEngine, aExList []string) []corev1.EnvVar {
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
			Value: cr.Namespace,
		},
		{
			Name:  "EXPERIMENT_LIST",
			Value: fmt.Sprint(strings.Join(aExList, ",")),
		},
		{
			Name:  "CHAOS_SVC_ACC",
			Value: cr.Spec.ChaosServiceAccount,
		},
	}
}

// getChaosExporterENV return the env required for chaos-exporter
func getChaosExporterENV(cr *litmuschaosv1alpha1.ChaosEngine, aUUID types.UID) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "CHAOSENGINE",
			Value: cr.Name,
		},
		{
			Name:  "APP_UUID",
			Value: string(aUUID),
		},
		{
			Name:  "APP_NAMESPACE",
			Value: cr.Namespace,
		},
	}
}

// getMonitoring return env required for metrics
func getMonitoringENV() []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Name: "metrics",
			Port: 8080,
		},
	}
}

// newRunnerPodForCR defines secondary resource #1 in same namespace as CR */
func newRunnerPodForCR(cr *litmuschaosv1alpha1.ChaosEngine, aUUID types.UID, aExList []string) (*corev1.Pod, error) {
	if len(aExList) == 0 || aUUID == "" {
		return nil, errors.New("expected aExList not found")
	}
	labels := map[string]string{
		"app": cr.Name,
	}
	podObj, err := pod.NewBuilder().
		WithName(cr.Name + "-runner").
		WithNamespace(cr.Namespace).
		WithLabels(labels).
		WithServiceAccountName(cr.Spec.ChaosServiceAccount).
		WithRestartPolicy("OnFailure").
		WithContainerBuilder(
			container.NewBuilder().
				WithName("chaos-runner").
				WithImage("ksatchit/ansible-runner:trial7").
				WithImagePullPolicy(corev1.PullIfNotPresent).
				WithCommandNew([]string{"/bin/bash"}).
				WithArgumentsNew([]string{"-c", "ansible-playbook ./executor/test.yml -i /etc/ansible/hosts; exit 0"}).
				WithEnvsNew(getChaosRunnerENV(cr, aExList)),
		).
		Build()
	if err != nil {
		return nil, err
	}
	return podObj, nil
}

// newExporterPodForCR defines secondary resource #2 in same namespace as CR */
func newExporterPodForCR(cr *litmuschaosv1alpha1.ChaosEngine, aUUID types.UID) (*corev1.Pod, error) {
	if cr == nil {
		return nil, errors.New("chaosengine got nil")
	}
	labels := map[string]string{
		"app":         cr.Name,
		"exporterFor": cr.Name,
	}
	exporterPod, err := pod.NewBuilder().
		WithName(cr.Name + "-exporter").
		WithNamespace(cr.Namespace).
		WithLabels(labels).
		WithServiceAccountName(cr.Spec.ChaosServiceAccount).
		WithRestartPolicy("OnFailure").
		WithContainerBuilder(
			container.NewBuilder().
				WithName("chaos-exporter").
				WithImage("litmuschaos/chaos-exporter:ci").
				WithPortsNew([]corev1.ContainerPort{{ContainerPort: 8080, Protocol: "TCP", Name: "metrics"}}).
				WithEnvsNew(getChaosExporterENV(cr, aUUID)),
		).
		Build()

	if err != nil {
		return nil, err
	}
	return exporterPod, nil
}

// newMonitorServiceForCR defines secondary resource #2 in same namespace as CR */
func newMonitorServiceForCR(cr *litmuschaosv1alpha1.ChaosEngine) (*corev1.Service, error) {
	if cr == nil {
		return nil, errors.New("nil chaosengine object")
	}
	labels := map[string]string{
		"app":         cr.Name,
		"exporterFor": cr.Name,
	}
	serviceObj, err := service.NewBuilder().
		WithName(cr.Name + "-monitor").
		WithNamespace(cr.Namespace).
		WithLabels(labels).
		WithPorts(getMonitoringENV()).
		WithSelectorsNew(
			map[string]string{
				"app":         cr.Name,
				"exporterFor": cr.Name,
			}).
		Build()
	if err != nil {
		return nil, err
	}
	return serviceObj, nil
}

// initializeApplicationInfo to initialize application info
func (appInfo *applicationInfo) initializeApplicationInfo(instance *litmuschaosv1alpha1.ChaosEngine) (*applicationInfo, error) {
	if instance == nil {
		return nil, errors.New("empty chaosengine")
	}
	appLabel := strings.Split(instance.Spec.Appinfo.Applabel, "=")
	appLabelKey = appLabel[0]
	appLabelValue = appLabel[1]
	appInfo.label = make(map[string]string)
	appInfo.label[appLabelKey] = appLabelValue
	appInfo.namespace = instance.Spec.Appinfo.Appns
	appInfo.experimentList = instance.Spec.Experiments
	appInfo.serviceAccountName = instance.Spec.ChaosServiceAccount
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
	monitorService.reqLogger.Info("Skip reconcile: engineMonitor Service already exists", "Service.Namespace", monitorService.service.Namespace, "Service.Name", monitorService.service.Name)
	return nil /*You can return now, both sec resources are existing */
}

// engineExporterPod to Check if the engineExporter Pod is already exists, else create
func engineExporterPod(engineExporter *corev1.Pod, r *ReconcileChaosEngine, reqLogger logr.Logger) error {
	pod := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: engineExporter.Name, Namespace: engineExporter.Namespace}, pod)
	if err != nil && k8serrors.IsNotFound(err) {
		reqLogger.Info("Creating a new engineExporter Pod", "Pod.Namespace", engineExporter.Namespace, "Pod.Name", engineExporter.Name)
		err = r.client.Create(context.TODO(), engineExporter)
		if err != nil {
			return err
		}

		reqLogger.Info("engineExporter Pod created successfully")
	} else if err != nil {
		return err
	}
	reqLogger.Info("Skip reconcile: engineExporter Pod already exists", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	return nil
}
