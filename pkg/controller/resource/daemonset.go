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
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// CheckDaemonSetAnnotation will check the annotation of DaemonSet
func CheckDaemonSetAnnotation(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getDaemonSetLists(clientSet, ce)
	if err != nil {
		return ce, err
	}
	ce, chaosEnabledDaemonSet := checkForEnabledChaos(targetAppList, ce)
	err = ValidateTotalChaosEnabled(chaosEnabledDaemonSet)
	if err != nil {
		return ce, err
	}
	chaosTypes.Log.Info("DaemonSet chaos candidate:", "appName: ", ce.AppName, " appUUID: ", ce.AppUUID)
	return ce, nil
}

// getDaemonSetLists will list the daemonSets which having the chaos label
func getDaemonSetLists(clientSet *kubernetes.Clientset, ce *chaosTypes.EngineInfo) (*v1.DaemonSetList, error) {
	targetAppList, err := clientSet.AppsV1().DaemonSets(ce.AppInfo.Namespace).List(metaV1.ListOptions{
		LabelSelector: ce.Instance.Spec.Appinfo.Applabel,
		FieldSelector: ""})
	if err != nil {
		return nil, fmt.Errorf("error while listing daemonSets with matching labels %s", ce.Instance.Spec.Appinfo.Applabel)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no daemonSets apps with matching labels %s", ce.Instance.Spec.Appinfo.Applabel)
	}
	return targetAppList, err
}

// This will check and count the total chaos enabled application
func checkForEnabledChaos(targetAppList *v1.DaemonSetList, ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int) {
	chaosEnabledDaemonSet := 0
	for _, daemonSet := range targetAppList.Items {
		ce.AppName = daemonSet.ObjectMeta.Name
		ce.AppUUID = daemonSet.ObjectMeta.UID
		annotationValue := daemonSet.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		chaosEnabledDaemonSet = CountTotalChaosEnabled(annotationValue, chaosEnabledDaemonSet)
	}
	return ce, chaosEnabledDaemonSet
}
