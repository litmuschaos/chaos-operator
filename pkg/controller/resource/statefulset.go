package resource

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// CheckStatefulSetAnnotation will check the annotation of StatefulSet
func CheckStatefulSetAnnotation(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := clientSet.AppsV1().StatefulSets(ce.AppInfo.Namespace).List(metav1.ListOptions{LabelSelector: ce.Instance.Spec.Appinfo.Applabel, FieldSelector: ""})
	if err != nil {
		return ce, fmt.Errorf("error while listing statefulsets with matching labels %s", ce.Instance.Spec.Appinfo.Applabel)
	}
	chaosCandidates := 0
	if len(targetAppList.Items) == 0 {
		return ce, fmt.Errorf("no statefulset apps with matching labels %s", ce.Instance.Spec.Appinfo.Applabel)
	}
	for _, app := range targetAppList.Items {
		ce.AppName = app.ObjectMeta.Name
		ce.AppUUID = app.ObjectMeta.UID
		annotationValue := app.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosCandidates, err = ValidateAnnotation(annotationValue, chaosCandidates)
		if err != nil {
			return ce, err
		}
		chaosTypes.Log.Info("chaos candidate : ", "appName", ce.AppName, "appUUID", ce.AppUUID)
	}
	return ce, nil
}
