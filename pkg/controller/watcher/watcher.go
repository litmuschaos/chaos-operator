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

package watcher

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
)

func getPodchaosUIDLabel(podLabels map[string]string) string {
	var chaosUID string
	if _, ok := podLabels["chaosUID"]; ok {
		chaosUID = podLabels["chaosUID"]
	}
	return chaosUID
}

func createListOptionsInNamespace(namespace string) []client.ListOption {
	listOptions := []client.ListOption{
		client.InNamespace(namespace),
	}
	return listOptions
}

func getChaosEngineList(listOptions []client.ListOption, clientSet client.Client) (litmuschaosv1alpha1.ChaosEngineList, error) {
	var listChaosEngine litmuschaosv1alpha1.ChaosEngineList

	err := clientSet.List(context.TODO(), &listChaosEngine, listOptions...)
	if err != nil {
		return litmuschaosv1alpha1.ChaosEngineList{}, err
	}
	return listChaosEngine, nil
}
