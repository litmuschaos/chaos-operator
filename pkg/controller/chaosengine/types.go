package chaosengine

import (
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
