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

// CheckStatefulSetAnnotation will check the annotation of StatefulSet
func CheckStatefulSetAnnotation(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getStatefulSetLists(clientset, engine)
	if err != nil {
		return engine, err
	}
	engine, chaosEnabledStatefulSet, err := checkForChaosEnabledStatefulSet(targetAppList, engine)
	if err != nil {
		return engine, err
	}
	if chaosEnabledStatefulSet == 0 {
		return engine, errors.New("no chaos-candidate found")
	}
	chaosTypes.Log.Info("Statefulset chaos candidate:", "appName: ", engine.AppName, " appUUID: ", engine.AppUUID)
	return engine, nil
}

// getStatefulSetLists will list the statefulset which having the chaos label
func getStatefulSetLists(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*appsV1.StatefulSetList, error) {
	targetAppList, err := clientset.AppsV1().StatefulSets(engine.AppInfo.Namespace).List(metaV1.ListOptions{
		LabelSelector: engine.Instance.Spec.Appinfo.Applabel,
		FieldSelector: ""})
	if err != nil {
		return nil, fmt.Errorf("error while listing statefulsets with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no statefulset apps with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	return targetAppList, err
}

// checkForChaosEnabledStatefulSet will check and count the total chaos enabled application
func checkForChaosEnabledStatefulSet(targetAppList *appsV1.StatefulSetList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int, error) {
	chaosEnabledStatefulSet := 0
	for _, statefulSet := range targetAppList.Items {
		engine.AppName = statefulSet.ObjectMeta.Name
		engine.AppUUID = statefulSet.ObjectMeta.UID
		annotationValue := statefulSet.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledStatefulSet = CountTotalChaosEnabled(annotationValue, chaosEnabledStatefulSet)
		if chaosEnabledStatefulSet > 1 {
			return engine, chaosEnabledStatefulSet, errors.New("too many statefulsets with specified label are annotated for chaos, either provide unique labels or annotate only desired app for chaos")
		}
	}
	return engine, chaosEnabledStatefulSet, nil
}
