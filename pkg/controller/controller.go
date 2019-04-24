package controller

import (
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	k8s "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
	cmd "k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	chaosAnnotation  = "litmuschaos.io/chaos"
	engineAnnotation = "litmuschaos.io/engine"
)

var log = logf.Log.WithName("cmd")

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}

// EventListner ...
func EventListner() {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := cmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Info(err.Error())
	}
	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		log.Info(err.Error())
	}
	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(metav1.Object)
			_, chaos := mObj.GetAnnotations()[chaosAnnotation]
			_, engine := mObj.GetAnnotations()[engineAnnotation]
			if chaos && engine {
				log.Info(fmt.Sprintf("Annotation present: %s", mObj.GetAnnotations()))
			}
		},
	})
	informer.Run(stopper)
}
