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
	"errors"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// CheckDeploymentAnnotation will check the annotation of deployment
func CheckDeploymentAnnotation(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getDeploymentLists(clientset, engine)
	if err != nil {
		return engine, err
	}
	engine, chaosEnabledDeployment, err := checkForChaosEnabledDeployment(targetAppList, engine)
	if err != nil {
		return engine, err
	}
	if chaosEnabledDeployment == 0 {
		return engine, errors.New("no chaos-candidate found")
	}
	chaosTypes.Log.Info("Deployment chaos candidate:", "appName: ", engine.AppName, " appUUID: ", engine.AppUUID)
	return engine, nil
}

// getDeploymentLists will list the deployments which having the chaos label
func getDeploymentLists(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*v1.DeploymentList, error) {
	targetAppList, err := clientset.AppsV1().Deployments(engine.AppInfo.Namespace).List(metaV1.ListOptions{
		LabelSelector: engine.Instance.Spec.Appinfo.Applabel,
		FieldSelector: ""})
	if err != nil {
		return nil, fmt.Errorf("error while listing deployments with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no deployments apps with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	return targetAppList, err
}

// checkForChaosEnabledDeployment will check and count the total chaos enabled application
func checkForChaosEnabledDeployment(targetAppList *v1.DeploymentList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int, error) {
	chaosEnabledDeployment := 0
	for _, deployment := range targetAppList.Items {
		engine.AppName = deployment.ObjectMeta.Name
		engine.AppUUID = deployment.ObjectMeta.UID
		annotationValue := deployment.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledDeployment = CountTotalChaosEnabled(annotationValue, chaosEnabledDeployment)
		if chaosEnabledDeployment > 1 {
			return engine, chaosEnabledDeployment, errors.New("too many deployments with specified label are annotated for chaos, either provide unique labels or annotate only desired app for chaos")
		}
	}
	return engine, chaosEnabledDeployment, nil
}
