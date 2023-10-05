package dynamic

import (
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// CreateClientSet returns a Dynamic Kubernetes ClientSet
func CreateClientSet() (*dynamic.Interface, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return &dynamicClient, nil
}
