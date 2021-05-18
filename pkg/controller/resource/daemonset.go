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

// CheckDaemonSetAnnotation check the annotation of the DaemonSet
func CheckDaemonSetAnnotation(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	targetAppList, err := getDaemonSetLists(clientset, engine)
	if err != nil {
		return engine, err
	}
	engine, chaosEnabledDaemonSet := checkForChaosEnabledDaemonSet(targetAppList, engine)
	if chaosEnabledDaemonSet == 0 {
		return engine, errors.New("no daemonsets chaos-candidate found")
	}
	return engine, nil
}

// getDaemonSetLists returns a list of daemonsets that are found in the app namespace with specified label
func getDaemonSetLists(clientset kubernetes.Interface, engine *chaosTypes.EngineInfo) (*appsV1.DaemonSetList, error) {
	targetAppList, err := clientset.AppsV1().DaemonSets(engine.Instance.Spec.Appinfo.Appns).List(metaV1.ListOptions{
		LabelSelector: engine.Instance.Spec.Appinfo.Applabel,
	})
	if err != nil {
		return nil, fmt.Errorf("error while listing daemonsets with matching labels %s, namespace: %s", engine.Instance.Spec.Appinfo.Applabel, engine.Instance.Spec.Appinfo.Appns)
	}
	if len(targetAppList.Items) == 0 {
		return nil, fmt.Errorf("no daemonsets found with matching labels: %s, namespace: %s", engine.Instance.Spec.Appinfo.Applabel, engine.Instance.Spec.Appinfo.Appns)
	}
	return targetAppList, err
}

// checkForChaosEnabledDaemonSet check and count the total chaos enabled application
func checkForChaosEnabledDaemonSet(targetAppList *appsV1.DaemonSetList, engine *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, int) {
	chaosEnabledDaemonSet := 0
	for _, daemonSet := range targetAppList.Items {
		annotationValue := daemonSet.ObjectMeta.GetAnnotations()[ChaosAnnotationKey]
		if IsChaosEnabled(annotationValue) {
			chaosTypes.Log.Info("chaos candidate for daemonset", "kind:", engine.Instance.Spec.Appinfo.AppKind, "appName: ", daemonSet.ObjectMeta.Name, "appUUID: ", daemonSet.ObjectMeta.UID)
			chaosEnabledDaemonSet++
		}
	}
	return engine, chaosEnabledDaemonSet
}
