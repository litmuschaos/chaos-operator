package resource

import (
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// CheckDeploymentAnnotation will check the annotation of deployment
func CheckDeploymentAnnotation(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := clientSet.AppsV1().Deployments(ce.AppInfo.Namespace).List(metav1.ListOptions{LabelSelector: ce.Instance.Spec.Appinfo.Applabel, FieldSelector: ""})
	if err != nil {
		return ce, errors.New("unable to list apps matching labels")
	}
	chaosCandidates := 0
	if len(targetAppList.Items) == 0 {
		return ce, errors.New("no app deployments with matching labels")
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
