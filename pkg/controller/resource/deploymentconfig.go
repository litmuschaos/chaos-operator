package resource

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var (
	dcGVR = schema.GroupVersionResource{
		Group:    "apps.openshift.io",
		Version:  "v1",
		Resource: "deploymentconfigs",
	}
)
// CheckDeploymentAnnotation will check the annotation of deployment
func CheckDeploymentConfigAnnotation(clientSet dynamic.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {

	targetAppList, err := getDeploymentConfigLists(clientSet, engine)
	if err != nil {
		return engine, err
	}

	engine, chaosEnabledDeployment, err := checkForChaosEnabledDeploymentConfig(targetAppList, engine)
	if err != nil {
		return engine, err
	}

	if chaosEnabledDeployment == 0 {
		return engine, errors.New("no chaos-candidate found")
	}
	chaosTypes.Log.Info("Deployment chaos candidate:", "appName: ", engine.AppName, " appUUID: ", engine.AppUUID)

	return engine, nil
}

func getDeploymentConfigLists(clientSet dynamic.Interface, engine *chaosTypes.EngineInfo) (*unstructured.UnstructuredList, error) {

	dyn := clientSet.Resource(dcGVR)
	targetAppList, err := dyn.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return targetAppList, err
}

// This will check and count the total chaos enabled application
func checkForChaosEnabledDeploymentConfig(targetAppList *unstructured.UnstructuredList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int, error) {

	chaosEnabledDeployment := 0
	for _, deploymentconfig := range targetAppList.Items {
		engine.AppName = deploymentconfig.GetName()
		engine.AppUUID = deploymentconfig.GetUID()
		annotationValue := deploymentconfig.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledDeployment = CountTotalChaosEnabled(annotationValue, chaosEnabledDeployment)
		if chaosEnabledDeployment > 1 {
			return engine, chaosEnabledDeployment, errors.New("too many deployments with specified label are annotated for chaos, either provide unique labels or annotate only desired app for chaos")
		}
	}
	return engine, chaosEnabledDeployment, nil
}