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

// CheckStatefulSetAnnotation will check the annotation of StatefulSet
func CheckStatefulSetAnnotation(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getStatefulSetLists(clientSet, ce)
	if err != nil {
		return ce, err
	}
	chaosEnabledStatefulset := 0
	for _, statefulset := range targetAppList.Items {
		ce.AppName = statefulset.ObjectMeta.Name
		ce.AppUUID = statefulset.ObjectMeta.UID
		annotationValue := statefulset.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledStatefulset = CountTotalChaosEnabled(annotationValue, chaosEnabledStatefulset)
	}
	err = ValidateTotalChaosEnabled(chaosEnabledStatefulset)
	if err != nil {
		return ce, err
	}
	chaosTypes.Log.Info("Statefulset chaos candidate:", "appName: ", ce.AppName, " appUUID: ", ce.AppUUID)
	return ce, nil
}

// getStatefulSetLists will list the statefulset which having the chaos label
func getStatefulSetLists(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*v1.StatefulSetList, error) {
	targetAppList, err := clientSet.AppsV1().StatefulSets(ce.AppInfo.Namespace).List(metav1.ListOptions{
		LabelSelector: ce.Instance.Spec.Appinfo.Applabel,
		FieldSelector: ""})
	if err != nil {
		return nil, fmt.Errorf("error while listing statefulsets with matching labels %s", ce.Instance.Spec.Appinfo.Applabel)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no statefulset apps with matching labels %s", ce.Instance.Spec.Appinfo.Applabel)
	}
	return targetAppList, err
}
