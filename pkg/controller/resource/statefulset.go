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

	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// CheckStatefulSetAnnotation will check the annotation of StatefulSet
func CheckStatefulSetAnnotation(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getStatefulSetLists(clientSet, ce)
	if err != nil {
		return ce, err
	}
	ce, chaosEnabledStatefulSet, err := checkForChaosEnabledStatefulSet(targetAppList, ce)
	if err != nil {
		return ce, err
	}
	if chaosEnabledStatefulSet == 0 {
		return ce, errors.New("no chaos-candidate found")
	}
	chaosTypes.Log.Info("Statefulset chaos candidate:", "appName: ", ce.AppName, " appUUID: ", ce.AppUUID)
	return ce, nil
}

// getStatefulSetLists will list the statefulset which having the chaos label
func getStatefulSetLists(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*appsV1.StatefulSetList, error) {
	targetAppList, err := clientSet.AppsV1().StatefulSets(ce.AppInfo.Namespace).List(metaV1.ListOptions{
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

// This will check and count the total chaos enabled application
func checkForChaosEnabledStatefulSet(targetAppList *appsV1.StatefulSetList, ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int, error) {
	chaosEnabledStatefulSet := 0
	for _, statefulSet := range targetAppList.Items {
		ce.AppName = statefulSet.ObjectMeta.Name
		ce.AppUUID = statefulSet.ObjectMeta.UID
		annotationValue := statefulSet.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledStatefulSet = CountTotalChaosEnabled(annotationValue, chaosEnabledStatefulSet)
		if chaosEnabledStatefulSet > 1 {
			return ce, chaosEnabledStatefulSet, errors.New("too many statefulsets with specified label are annotated for chaos, either provide unique labels or annotate only desired app for chaos")
		}
	}
	return ce, chaosEnabledStatefulSet, nil
}
