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

	appsV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// CheckStatefulSetAnnotation check the annotation of the StatefulSet
func CheckStatefulSetAnnotation(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getStatefulSetLists(clientset, engine)
	if err != nil {
		return engine, err
	}
	engine, chaosEnabledStatefulSet := checkForChaosEnabledStatefulSet(targetAppList, engine)
	if chaosEnabledStatefulSet == 0 {
		return engine, errors.New("no statefulsets chaos-candidate found")
	}
	return engine, nil
}

// getStatefulSetLists returns a list of statefulsets that are found in the app namespace with specified label
func getStatefulSetLists(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*appsV1.StatefulSetList, error) {
	targetAppList, err := clientset.AppsV1().StatefulSets(engine.Instance.Spec.Appinfo.Appns).List(metaV1.ListOptions{
		LabelSelector: engine.Instance.Spec.Appinfo.Applabel,
	})
	if err != nil {
		return nil, fmt.Errorf("error while listing statefulsets with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no statefulset found with matching labels: %s, namespace: %s", engine.Instance.Spec.Appinfo.Applabel, engine.Instance.Spec.Appinfo.Appns)
	}
	return targetAppList, err
}

// checkForChaosEnabledStatefulSet check and count the total chaos enabled application
func checkForChaosEnabledStatefulSet(targetAppList *appsV1.StatefulSetList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int) {
	chaosEnabledStatefulSet := 0
	for _, statefulSet := range targetAppList.Items {
		annotationValue := statefulSet.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		if IsChaosEnabled(annotationValue) {
			chaosTypes.Log.Info("chaos candidate for statefulset", "kind:", engine.Instance.Spec.Appinfo.AppKind, "appName: ", statefulSet.ObjectMeta.Name, "appUUID: ", statefulSet.ObjectMeta.UID)
			chaosEnabledStatefulSet++
		}
	}
	return engine, chaosEnabledStatefulSet
}
