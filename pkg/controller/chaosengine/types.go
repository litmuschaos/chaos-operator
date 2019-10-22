package chaosengine

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/go-logr/logr"
	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
)

// To create logs for debugging or detailing, please follow this syntax.
// use function log.Info
// in parameters give the name of the log / error (string) ,
// with the variable name for the value(string)
// and then the value to log (any datatype)
// All values should be in key : value pairs only
// For eg. : log.Info("name_of_the_log","variable_name_for_the_value",value, ......)
// For eg. : log.Error(err,"error_statement","variable_name",value)
// For eg. : log.Printf
//("error statement %q other variables %s/%s",targetValue, object.Namespace, object.Name)
// For eg. : log.Errorf
//("unable to reconcile object %s/%s: %v", object.Namespace, object.Name, err)
// This logger uses a structured logging schema in JSON format, which will / can be used further
// to access the values in the logger.

var (
	appLabelKey         string
	appLabelValue       string
	log                                      = logf.Log.WithName("controller_chaosengine")
	_                   reconcile.Reconciler = &ReconcileChaosEngine{}
	defaultRunnerImage                       = "litmuschaos/ansible-runner:ci"
	defaultMonitorImage                      = "litmuschaos/chaos-exporter:ci"
)

// Annotations on app to enable chaos on it
const (
	chaosAnnotationKey   = "litmuschaos.io/chaos"
	chaosAnnotationValue = "true"
)

// ReconcileChaosEngine reconciles a ChaosEngine object
type ReconcileChaosEngine struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// applicationInfo contains the chaos details for target application
type applicationInfo struct {
	namespace          string
	label              map[string]string
	experimentList     []litmuschaosv1alpha1.ExperimentList
	serviceAccountName string
	kind               string
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

//engine Related information
type engineInfo struct {
	instance       *litmuschaosv1alpha1.ChaosEngine
	appInfo        *applicationInfo
	appExperiments []string
	appName        string
	appUUID        types.UID
}
