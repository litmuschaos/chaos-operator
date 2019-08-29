package chaosengine

import (
	"context"
	"fmt"
	"strings"

	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	// Temp test purposes
	//"github.com/Sirupsen/logrus"
)

var log = logf.Log.WithName("controller_chaosengine")

// Annotations on app to enable chaos on it
const (
	chaosAnnotation = "litmuschaos.io/chaos"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

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
	// Create a new controller
	c, err := controller.New("chaosengine-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ChaosEngine
	err = c.Watch(&source.Kind{Type: &litmuschaosv1alpha1.ChaosEngine{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ChaosEngine
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &litmuschaosv1alpha1.ChaosEngine{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileChaosEngine{}

// ReconcileChaosEngine reconciles a ChaosEngine object
type ReconcileChaosEngine struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
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
		if errors.IsNotFound(err) {
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

	aLabelKeyValue := strings.Split(instance.Spec.Appinfo.Applabel, "=")
	lKey := aLabelKeyValue[0]
	lValue := aLabelKeyValue[1]
	aLabel := make(map[string]string)
	aLabel[lKey] = lValue
	aNamespace := instance.Spec.Appinfo.Appns
	aExpList := instance.Spec.Experiments

	var appExperiments []string
	for _, exp := range aExpList {
		appExperiments = append(appExperiments, exp.Name)
	}

	// Temp test purposes
	/*
	   logrus.Info("App Label derived from Chaosengine is ", aLabel)
	   logrus.Info("App NS derived from Chaosengine is ", aNamespace)
	   logrus.Info("Exp list derived from chaosengine is ", appExperiments)
	*/

	log.Info("App Label derived from Chaosengine is ", aLabel)
	log.Info("App NS derived from Chaosengine is ", aNamespace)
	log.Info("Exp list derived from chaosengine is ", appExperiments)

	// Use client-Go to obtain a list of apps w/ specified labels
	config, err := config.GetConfig()
	if err != nil {
		//logrus.Fatal(err.Error())
		log.Error(err, "unable to get kube config")
		return reconcile.Result{}, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		//logrus.Fatal(err.Error())
		log.Error(err, "unable to create clientset using kubeconfig")
		return reconcile.Result{}, err
	}

	cApp, err := clientset.AppsV1().Deployments(aNamespace).List(metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", lKey, lValue), FieldSelector: ""})
	if err != nil {
		//logrus.Fatal("Failed to list deployments. Error is ", err)
		log.Error(err, "unable to list apps matching labels")
		return reconcile.Result{}, err
	}

	var appName string
	var appUUID types.UID

	// Determine whether apps with matching labels have chaos annotation set to true
	chaosCandidates := 0
	if len(cApp.Items) > 0 {
		for _, app := range cApp.Items {
			appName = app.ObjectMeta.Name
			appUUID = app.ObjectMeta.UID
			appCaSts := metav1.HasAnnotation(app.ObjectMeta, chaosAnnotation)
			//if appCaSts == true {
			if appCaSts {
				//logrus.Info ("chaos candidate app: ", appName, appUUID)
				log.Info("chaos candidate app: ", appName, appUUID)
				chaosCandidates++
			}
		}
		if chaosCandidates == 0 {
			//logrus.Info("No chaos candidates found")
			log.Info("No chaos candidates found")
			return reconcile.Result{}, nil
		} else if chaosCandidates > 1 {
			//logrus.Info ("Too many chaos candidates with same label",
			log.Info("Too many chaos candidates with same label",
				"either provide unique labels or annotate only desired app for chaos")
			return reconcile.Result{}, nil
		}
	} else {
		//logrus.Info("No app deployments with matching labels")
		log.Info("No app deployments with matching labels")
		return reconcile.Result{}, nil
	}

	// Define an engine(ansible?)-runner pod which is secondary-resource #1
	engineRunner := newRunnerPodForCR(instance, appUUID, appExperiments)

	// Define the engine-monitor service which is secondary-resource #2
	engineMonitor := newMonitorServiceForCR(instance)

	// Set ChaosEngine instance as the owner and controller of engine-runner pod
	if err := controllerutil.SetControllerReference(instance, engineRunner, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set ChaosEngine instance as the owner and controller of engine-monitor service
	if err := controllerutil.SetControllerReference(instance, engineMonitor, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the engineRunner pod already exists, else create
	foundS1 := &corev1.Pod{} //secondary resource #1
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: engineRunner.Name, Namespace: engineRunner.Namespace}, foundS1)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new engineRunner Pod", "Pod.Namespace", engineRunner.Namespace, "Pod.Name", engineRunner.Name)
		err = r.client.Create(context.TODO(), engineRunner)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		reqLogger.Info("engineRunner Pod created successfully")
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: engineRunner Pod already exists", "Pod.Namespace", foundS1.Namespace, "Pod.Name", foundS1.Name)

	// Check if the engineMonitorservice already exists, else create
	foundS2 := &corev1.Service{} //secondary resource #2
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: engineMonitor.Name, Namespace: engineMonitor.Namespace}, foundS2)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new engineMonitor Service", "Service.Namespace", engineMonitor.Namespace, "Service.Name", engineMonitor.Name)
		err = r.client.Create(context.TODO(), engineMonitor)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Service created successfully - don't requeue
		return reconcile.Result{}, nil /*You can return now, both sec resources are created */
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Service already exists - don't requeue
	reqLogger.Info("Skip reconcile: engineMonitor Service already exists", "Service.Namespace", foundS2.Namespace, "Service.Name", foundS2.Name)
	return reconcile.Result{}, nil /*You can return now, both sec resources are existing */
}

// newRunnerPodForCR defines secondary resource #1 in same namespace as CR */
func newRunnerPodForCR(cr *litmuschaosv1alpha1.ChaosEngine, aUUID types.UID, aExList []string) *corev1.Pod {
	monitorstatus := cr.Spec.Monitorstatus
	labels := map[string]string{
		"app": cr.Name,
	}
	if monitorstatus {
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.Name + "-runner",
				Namespace: cr.Namespace,
				Labels:    labels,
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: "chaos-operator",
				Containers: []corev1.Container{
					{
						Name:    "chaos-runner",
						Image:   "openebs/ansible-runner:ci",
						Command: []string{"/bin/bash"},
						//TODO: Reconcile will restart tests. This has to be addressed
						Args: []string{"-c", "ansible-playbook ./executor/test.yml -i /etc/ansible/hosts -vv; exit 0"},
						Env: []corev1.EnvVar{
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
								Name: "EXPERIMENT_LIST",
								//Value: fmt.Sprintf(strings.Join(aExList,",")),
								Value: fmt.Sprint(strings.Join(aExList, ",")),
							},
						},
					},
					{
						Name:  "chaos-exporter",
						Image: "litmuschaos/chaos-exporter:ci",
						Env: []corev1.EnvVar{
							{
								Name:  "CHAOSENGINE",
								Value: cr.Name,
							},
							{
								Name:  "APP_UUID",
								Value: string(aUUID),
							},
						},
					},
				},
			},
		}
	} else {
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.Name + "-runner",
				Namespace: cr.Namespace,
				Labels:    labels,
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: "chaos-operator",
				Containers: []corev1.Container{
					{
						Name:    "chaos-runner",
						Image:   "openebs/ansible-runner:ci",
						Command: []string{"/bin/bash"},
						//TODO: Reconcile will restart tests. This has to be addressed
						Args: []string{"-c", "ansible-playbook ./executor/test.yml -i /etc/ansible/hosts -vv; exit 0"},
						Env: []corev1.EnvVar{
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
								Name: "EXPERIMENT_LIST",
								//Value: fmt.Sprintf(strings.Join(aExList,",")),
								Value: fmt.Sprint(strings.Join(aExList, ",")),
							},
						},
					},
				},
			},
		}

	}
}

// newMonitorServiceForCR defines secondary resource #2 in same namespace as CR */
func newMonitorServiceForCR(cr *litmuschaosv1alpha1.ChaosEngine) *corev1.Service {
	monitorstatus := cr.Spec.Monitorstatus
	labels := map[string]string{
		"app": cr.Name,
	}
	if monitorstatus {
		return &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.Name + "-monitor",
				Namespace: cr.Namespace,
				Labels:    labels,
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: "metrics",
						Port: 8080,
					},
				},
				Selector: map[string]string{
					"app": cr.Name,
				},
			},
		}
	} else {
		return nil
	}
}

