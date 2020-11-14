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

// CheckDaemonSetAnnotation will check the annotation of DaemonSet
func CheckDaemonSetAnnotation(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getDaemonSetLists(clientset, engine)
	if err != nil {
		return engine, err
	}
	engine, chaosEnabledDaemonSet, err := checkForChaosEnabledDaemonSet(targetAppList, engine)
	if err != nil {
		return engine, err
	}
	if chaosEnabledDaemonSet == 0 {
		return engine, errors.New("no chaos-candidate found")
	}
	return engine, nil
}

// getDaemonSetLists will list the daemonSets which having the chaos label
func getDaemonSetLists(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*appsV1.DaemonSetList, error) {
	targetAppList, err := clientset.AppsV1().DaemonSets(engine.AppInfo.Namespace).List(metaV1.ListOptions{
		LabelSelector: engine.Instance.Spec.Appinfo.Applabel,
		FieldSelector: ""})
	if err != nil {
		return nil, fmt.Errorf("error while listing daemonSets with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no daemonSets apps with matching labels %s", engine.Instance.Spec.Appinfo.Applabel)
	}
	return targetAppList, err
}

// checkForChaosEnabledDaemonSet will check and count the total chaos enabled application
func checkForChaosEnabledDaemonSet(targetAppList *appsV1.DaemonSetList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int, error) {
	chaosEnabledDaemonSet := 0
	for _, daemonSet := range targetAppList.Items {
		engine.AppName = daemonSet.ObjectMeta.Name
		engine.AppUUID = daemonSet.ObjectMeta.UID
		annotationValue := daemonSet.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledDaemonSet = CountTotalChaosEnabled(annotationValue, chaosEnabledDaemonSet)
		if chaosEnabledDaemonSet > 1 {
			return engine, chaosEnabledDaemonSet, errors.New("too many daemonsets with specified label are annotated for chaos, either provide unique labels or annotate only desired app for chaos")
		}
	}
	return engine, chaosEnabledDaemonSet, nil
}
