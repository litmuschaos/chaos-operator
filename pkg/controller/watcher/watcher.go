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
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

// WatchForRunnerPod creates watcher for Chaos Runner Pod
func WatchForRunnerPod(client client.Client, c controller.Controller) error {

	runnerPodHandler := handlerForRunnerPod(client)

	return c.Watch(&source.Kind{Type: &corev1.Pod{}}, &runnerPodHandler)
}

// handlerForRunnerPod creates a event Handler for Chaos Runner Pod
func handlerForRunnerPod(clientSet client.Client) handler.EnqueueRequestsFromMapFunc {
	reqLogger := chaosTypes.Log.WithName("Chaos Resources Watch")
	var runnerPodRequest []reconcile.Request
	var err error

	handlerForRunner := handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			runnerNameCheck := strings.HasSuffix(a.Meta.GetName(), "-runner")
			if runnerNameCheck {
				runnerPodRequest, err = createHandlerRequestForEngine(a, clientSet)
				if err != nil {
					reqLogger.Error(err, "Unable to get the ChaosEngine Resources", "namespace", a.Meta.GetNamespace())
					return nil
				}
			}
			return runnerPodRequest
		}),
	}
	return handlerForRunner
}

// handlerRequestFromEngineList initilize a event Watcher filtering the chaosEngine from the list.
func handlerRequestFromEngineList(listChaosEngine litmuschaosv1alpha1.ChaosEngineList, chaosUID string) []reconcile.Request {
	for i := range listChaosEngine.Items {
		uuid := string(listChaosEngine.Items[i].GetUID())
		if chaosUID == uuid {
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      listChaosEngine.Items[i].GetName(),
					Namespace: listChaosEngine.Items[i].GetNamespace(),
				}},
			}
		}
	}
	return []reconcile.Request{
		{NamespacedName: types.NamespacedName{
			Name:      "",
			Namespace: "",
		}},
	}

}

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

func createHandlerRequestForEngine(a handler.MapObject, clientSet client.Client) ([]reconcile.Request, error) {
	chaosUID := getPodchaosUIDLabel(a.Meta.GetLabels())
	listChaosEngine, err := getChaosEngineList(createListOptionsInNamespace(a.Meta.GetNamespace()), clientSet)
	if err != nil {
		return nil, fmt.Errorf("Unable to get the ChaosEngine Resources in namespace: %v", a.Meta.GetNamespace())
	}
	chaosEngineListRequest := handlerRequestFromEngineList(listChaosEngine, chaosUID)
	return chaosEngineListRequest, nil
}
