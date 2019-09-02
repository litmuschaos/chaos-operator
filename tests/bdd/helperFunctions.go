package bdd

import (
	"fmt"
	"os"
)

// HomeDir return the Home Directory of the environement
func HomeDir() (string, error) {
	if h := os.Getenv("HOME"); h != "" { // linux
		return h, nil
	} else if h := os.Getenv("USERPROFILE"); h != "" { // windows
		return h, nil
	}

	return "", fmt.Errorf("Not able to locate home directory")
}

// GetConfigPath returns the path of kubeconfig
func GetConfigPath() (kubeConfigPath string, err error) {
	home, err := HomeDir()
	if err != nil {
		return
	}
	kubeConfigPath = os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		// Parse the kube config path
		kubeConfigPath = home + "/.kube/config"
	}
	return kubeConfigPath, err
}
