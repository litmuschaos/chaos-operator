package chaosengine

import (
        "strings"
	"context"
        "fmt"
	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
        //appsv1 "k8s.io/api/apps/v1"
        //"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	corev1 "k8s.io/api/core/v1"
        "k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

        // @ksatchit: temp test purposes
        "github.com/Sirupsen/logrus"
)

var log = logf.Log.WithName("controller_chaosengine")

// @ksatchit: Annotations on app to enable chaos on it 
const (
      chaosAnnotation  = "litmuschaos.io/chaos"
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
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
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

	// Fetch the app details from ChaosEngine instance. Check if app is present, else fail
        // TODO: Get app kind from chaosengine spec as well. Using "deploy" for now
        // TODO: Establish label format in chaosengine OR "=" as a const 

        a_labelKeyValue := strings.Split(instance.Spec.Appinfo.Applabel, "=")
        a_label := make(map[string]string)
        lkey := a_labelKeyValue[0]; lvalue := a_labelKeyValue[1]
        a_label[lkey] = lvalue
        a_namespace := instance.Spec.Appinfo.Appns

        // @ksatchit: temp test purposes
        logrus.Info("App Label derived from Chaosengine is ", a_label)
        logrus.Info("App NS derived from Chaosengine is ", a_namespace)

        //found_a1 := &appsv1.DeploymentList{}
        //err = r.client.List(context.TODO(), client.InNamespace(a_namespace).MatchingLabels(a_label), found_a1)

        config, err := config.GetConfig()
        if err != nil {
	  logrus.Fatal(err.Error())
	}

        clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
	  logrus.Fatal(err.Error())
	}

        c_app, err := clientset.AppsV1().Deployments(a_namespace).List(metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", lkey, lvalue), FieldSelector: ""})
        if err != nil {
          logrus.Fatal("Failed to list deployments. Error is ", err)
        }

        var app_uuid string
        chaos_candidates := 0
        if len(c_app.Items) > 0 {
          for _, app := range c_app.Items {
            app_name := app.ObjectMeta.Name
            app_uuid := app.ObjectMeta.UID
            app_ca_sts := metav1.HasAnnotation(app.ObjectMeta, chaosAnnotation)
            if app_ca_sts == true {
              logrus.Info ("chaos candidate app: ", app_name, app_uuid)
              chaos_candidates++
            }
          }
          if chaos_candidates == 0 {
            logrus.Info("No chaos candidates found")
            return reconcile.Result{}, nil
          } else if chaos_candidates > 1 {
            logrus.Info ("Too many chaos candidates with same label",
              "either provide unique labels or annotate only desired app for chaos",)
            return reconcile.Result{}, nil
          }
        } else {
          logrus.Info("No app deployments with matching labels")
          return reconcile.Result{}, nil
        }

	// Define a new Pod object (OrigComment)
	//pod := newPodForCR(instance) // 
        /* @ksatchit: define an engine(ansible?)-runner pod which is secondary-resource #1 */
        engineRunner := newRunnerPodForCR(instance, app_uuid)
        engineMonitor := newMonitorServiceForCR(instance)

	// Set ChaosEngine instance as the owner and controller (OrigComment)
	if err := controllerutil.SetControllerReference(instance, engineRunner, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	if err := controllerutil.SetControllerReference(instance, engineMonitor, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

        /* @ksatchit: Check if the engineRunner pod already exists */
	found_s1 := &corev1.Pod{} //secondary resource #1
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: engineRunner.Name, Namespace: engineRunner.Namespace}, found_s1)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new engineRunner Pod", "Pod.Namespace", engineRunner.Namespace, "Pod.Name", engineRunner.Name)
		err = r.client.Create(context.TODO(), engineRunner)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
                reqLogger.Info("engineRunner Pod created successfully")
		// @ksatchit: return reconcile.Result{}, nil /*Dont return, go to service check */ 
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: engineRunner Pod already exists", "Pod.Namespace", found_s1.Namespace, "Pod.Name", found_s1.Name)
	// @ksatchit: return reconcile.Result{}, nil / *Dont return, go to service check */

        /* @ksatchit: Check if the engineMonitor pod already exists */
	found_s2 := &corev1.Service{} //secondary resource #2
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: engineMonitor.Name, Namespace: engineMonitor.Namespace}, found_s2)
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
	reqLogger.Info("Skip reconcile: engineMonitor Service already exists", "Service.Namespace", found_s2.Namespace, "Service.Name", found_s2.Name)
	return reconcile.Result{}, nil /*You can return now, both sec resources are existing */
}


/* @ksatchit: function defining pod as secondary resource #1 in same namespace as cr */
func newRunnerPodForCR(cr *litmuschaosv1alpha1.ChaosEngine, a_uuid string) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-runner",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "chaos-runner",
					Image:   "openebs/ansible-runner:ci",
                                        //TODO: Get exp list - stage#1
					Command: []string{"sleep", "3600"},
                                        Env:     []corev1.EnvVar{
                                             {
                                                 Name: "CHAOSENGINE",
                                                 Value: cr.Name,
                                             },
                                             {
                                                 Name: "APP_UUID",
                                                 Value: a_uuid,
                                             },
				         },
                                },
				{
					Name:    "chaos-exporter",
					Image:   "ksatchit/sample-chaos-exporter:ci",
                                        //TODO: Get exp list - stage#1
					Command: []string{"sleep", "3600"},
                                        Env:     []corev1.EnvVar{
                                            {
                                                 Name: "CHAOSENGINE",
                                                 Value: cr.Name,
                                            },
                                            {
                                                 Name: "APP_UUID",
                                                 Value: a_uuid,
                                            },
                                        },
				},
			},
		},
	}
}

/* @ksatchit: function defining yaml for secondary resource #1 */
func newMonitorServiceForCR(cr *litmuschaosv1alpha1.ChaosEngine) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-monitor",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:    "metrics",
					Port:    8080,
				},
			},
                        Selector: map[string]string{
                                "app": cr.Name,
                        },
		},
	}
}
