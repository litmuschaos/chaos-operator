package dynamic


import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func CreateClientSet() (*dynamic.Interface, error){
	restConfig, err := config.GetConfig()
	if err != nil {
		fmt.Print(err)
	}
	clientSet, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		fmt.Print(err)
	}
	return &clientSet, nil
}