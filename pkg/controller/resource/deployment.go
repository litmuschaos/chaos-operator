/*
Copyright 2019 LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resource

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// CheckDeploymentAnnotation will check the annotation of deployment
func CheckDeploymentAnnotation(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getDeploymentLists(clientSet, ce)
	if err != nil {
		return ce, err
	}
	chaosEnabledDeployment := 0
	for _, deployment := range targetAppList.Items {
		ce.AppName = deployment.ObjectMeta.Name
		ce.AppUUID = deployment.ObjectMeta.UID
		annotationValue := deployment.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledDeployment = CountTotalChaosEnabled(annotationValue, chaosEnabledDeployment)
	}
	err = ValidateTotalChaosEnabled(chaosEnabledDeployment)
	if err != nil {
		return ce, err
	}
	chaosTypes.Log.Info("Deployment chaos candidate:", "appName: ", ce.AppName, " appUUID: ", ce.AppUUID)
	return ce, nil
}

// getDeploymentLists will list the deployments which having the chaos label
func getDeploymentLists(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*v1.DeploymentList, error) {
	targetAppList, err := clientSet.AppsV1().Deployments(ce.AppInfo.Namespace).List(metav1.ListOptions{
		LabelSelector: ce.Instance.Spec.Appinfo.Applabel,
		FieldSelector: ""})
	if err != nil {
		return nil, fmt.Errorf("error while listing deployments with matching labels %s", ce.Instance.Spec.Appinfo.Applabel)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no deployments apps with matching labels %s", ce.Instance.Spec.Appinfo.Applabel)
	}
	return targetAppList, err
}
