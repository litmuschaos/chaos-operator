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
	gvrro = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "rollouts",
	}
)

// CheckRolloutAnnotation will check the annotation of argo rollout 
func CheckRolloutAnnotation(clientSet dynamic.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {

	rolloutList, err := getRolloutList(clientSet, engine)
	if err != nil {
		return engine, err
	}
	engine, chaosEnabledRollout, err := checkForChaosEnabledRollout(rolloutList, engine)
	if err != nil {
		return engine, err
	}

	if chaosEnabledRollout == 0 {
		return engine, errors.New("no argo rollout chaos-candidate found")
	}
	chaosTypes.Log.Info("argo rollout chaos candidate:", "appName: ", engine.AppName, " appUUID: ", engine.AppUUID)

	return engine, nil
}

// getRolloutList returns a list of argo rollout resources that are found in the app namespace with specified label
func getRolloutList(clientSet dynamic.Interface, engine *chaosTypes.EngineInfo) (*unstructured.UnstructuredList, error) {

	dynamicClient := clientSet.Resource(gvrro)

	rolloutList, err := dynamicClient.Namespace(engine.AppInfo.Namespace).List(metav1.ListOptions{
			LabelSelector: engine.Instance.Spec.Appinfo.Applabel})
	if err != nil {
		return nil, fmt.Errorf("error while listing argo rollouts with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	if len(rolloutList.Items) == 0 {
		return nil, fmt.Errorf("no argo rollouts with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	return rolloutList, err
}

// checkForChaosEnabledRollout  will check and count the total chaos enabled application
func checkForChaosEnabledRollout(rolloutList *unstructured.UnstructuredList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int, error) {

	chaosEnabledRollout := 0
	for _, rollout := range rolloutList.Items {
		engine.AppName = rollout.GetName()
		engine.AppUUID = rollout.GetUID()
		annotationValue := rollout.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledRollout = CountTotalChaosEnabled(annotationValue, chaosEnabledRollout)
		if chaosEnabledRollout > 1 {
			return engine, chaosEnabledRollout, errors.New("too many argo rollouts with specified label are annotated for chaos, either provide unique labels or annotate only desired rollout for chaos")
		}
	}
	return engine, chaosEnabledRollout, nil
}
