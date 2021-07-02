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

// CheckDeploymentAnnotation check the annotation of the deployment
func CheckDeploymentAnnotation(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getDeploymentLists(clientset, engine)
	if err != nil {
		return engine, err
	}
	engine, chaosEnabledDeployment := checkForChaosEnabledDeployment(targetAppList, engine)
	if chaosEnabledDeployment == 0 {
		return engine, errors.New("no deployment chaos-candidate found")
	}
	return engine, nil
}

// getDeploymentLists returns a list of deployments that are found in the app namespace with specified label
func getDeploymentLists(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*v1.DeploymentList, error) {
	targetAppList, err := clientset.AppsV1().Deployments(engine.Instance.Spec.Appinfo.Appns).List(metaV1.ListOptions{
		LabelSelector: engine.Instance.Spec.Appinfo.Applabel,
	})
	if err != nil {
		return nil, fmt.Errorf("error while listing deployments with matching labels %s, namespace: %s", engine.Instance.Spec.Appinfo.Applabel, engine.Instance.Spec.Appinfo.Appns)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no deployment found with matching labels: %s, namespace: %s", engine.Instance.Spec.Appinfo.Applabel, engine.Instance.Spec.Appinfo.Appns)
	}
	return targetAppList, err
}

// checkForChaosEnabledDeployment check and count the total chaos enabled application
func checkForChaosEnabledDeployment(targetAppList *v1.DeploymentList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int) {
	chaosEnabledDeployment := 0
	for _, deployment := range targetAppList.Items {
		annotationValue := deployment.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		if IsChaosEnabled(annotationValue) {
			chaosTypes.Log.Info("chaos candidate for deployment", "kind:", engine.Instance.Spec.Appinfo.AppKind, "appName: ", deployment.ObjectMeta.Name)
			chaosEnabledDeployment++
		}
	}
	return engine, chaosEnabledDeployment
}
