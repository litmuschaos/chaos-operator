package chaosengine

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Use client-Go to obtain a list of apps w/ specified labels
func createClientSet() (*kubernetes.Clientset, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get rest kube config")
		return &kubernetes.Clientset{}, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Error(err, "unable to create clientset using restconfig")
		return &kubernetes.Clientset{}, err
	}
	return clientSet, nil
}

// Determine whether apps with matching labels have chaos annotation set to true
func checkAnnotation() error {
	clientSet, err := createClientSet()
	if err != nil {
		log.Error(err, "Clientset generation failed with error: ")
	}
	switch strings.ToLower(engine.appInfo.kind) {
	case "deployment":
		{
			err = checkAnnotationDeployemet(clientSet)
			if err != nil {
				return fmt.Errorf("no deployement found with required annotation, err: %+v", err)
			}
		}
	case "statefulset":
		{
			err := checkAnnotationStatefulSet(clientSet)
			if err != nil {
				return fmt.Errorf("no statefulset found with required annotation, err: %+v", err)
			}
		}

	case "daemonset":
		{
			err := checkAnnotationDaemonSet(clientSet)
			if err != nil {
				return fmt.Errorf("no daemonset found with required annotation, err: %+v", err)
			}
		}
	}
	return nil
}
