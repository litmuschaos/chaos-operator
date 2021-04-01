package resource

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

var (
	gvrdc = schema.GroupVersionResource{
		Group:    "apps.openshift.io",
		Version:  "v1",
		Resource: "deploymentconfigs",
	}
)

// CheckDeploymentConfigAnnotation check the annotation of the deploymentConfig
func CheckDeploymentConfigAnnotation(clientSet dynamic.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {

	deploymentConfigList, err := getDeploymentConfigList(clientSet, engine)
	if err != nil {
		return engine, err
	}
	engine, chaosEnabledDeploymentConfig := checkForChaosEnabledDeploymentConfig(deploymentConfigList, engine)
	if chaosEnabledDeploymentConfig == 0 {
		return engine, errors.New("no DeploymentConfig chaos-candidate found")
	}

	return engine, nil
}

//getDeploymentConfigList returns a list of deploymentConfigs that are found in the app namespace with specified label
func getDeploymentConfigList(clientSet dynamic.Interface, engine *chaosTypes.EngineInfo) (*unstructured.UnstructuredList, error) {

	dynamicClient := clientSet.Resource(gvrdc)

	deploymentConfigList, err := dynamicClient.Namespace(engine.Instance.Spec.Appinfo.Appns).List(metav1.ListOptions{
		LabelSelector: engine.Instance.Spec.Appinfo.Applabel})
	if err != nil {
		return nil, fmt.Errorf("error while listing deploymentconfigs with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	if len(deploymentConfigList.Items) == 0 {
		return nil, fmt.Errorf("no deploymentconfigs found with matching labels: %s, namespace: %s", engine.Instance.Spec.Appinfo.Applabel, engine.Instance.Spec.Appinfo.Appns)
	}
	return deploymentConfigList, err
}

// checkForChaosEnabledDeploymentConfig check and count the total chaos enabled application
func checkForChaosEnabledDeploymentConfig(deploymentConfigList *unstructured.UnstructuredList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int) {

	chaosEnabledDeploymentConfig := 0
	for _, deploymentconfig := range deploymentConfigList.Items {
		annotationValue := deploymentconfig.GetAnnotations()[ChaosAnnotationKey]
		if IsChaosEnabled(annotationValue) {
			chaosTypes.Log.Info("chaos candidate for deploymentconfig", "kind:", engine.Instance.Spec.Appinfo.AppKind, "appName: ", deploymentconfig.GetName(), "appUUID: ", deploymentconfig.GetUID())
			chaosEnabledDeploymentConfig++
		}
	}
	return engine, chaosEnabledDeploymentConfig
}
