package resource

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// CheckDeploymentAnnotation will check the annotation of deployment
func CheckDeploymentAnnotation(clientSet *kubernetes.Clientset, engine chaosTypes.EngineInfo) error {
	targetAppList, err := clientSet.AppsV1().Deployments(engine.AppInfo.Namespace).List(metav1.ListOptions{LabelSelector: engine.Instance.Spec.Appinfo.Applabel, FieldSelector: ""})
	if err != nil {
		return errors.New("unable to list apps matching labels")
	}
	chaosCandidates := 0
	if len(targetAppList.Items) == 0 {
		return errors.New("no app deployments with matching labels")
	}
	for _, app := range targetAppList.Items {
		engine.AppName = app.ObjectMeta.Name
		engine.AppUUID = app.ObjectMeta.UID
		annotationValue := app.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosCandidates, err = ValidateAnnotation(annotationValue, chaosCandidates)
		if err != nil {
			return err
		}
		//chaosTypes.Log.Info("chaos candidate : ", "appName", engine.AppName, "appUUID", engine.AppUUID)
	}
	return nil
}
