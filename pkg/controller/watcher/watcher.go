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
	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
)

// WatchForMonitorService creates a watcher for Chaos MonitorService
func WatchForMonitorService(client client.Client, c controller.Controller) error {

	monitorServiceHandler := handlerForMonitorService(client)

	return c.Watch(&source.Kind{Type: &corev1.Service{}}, &monitorServiceHandler)
}

// WatchForRunnerPod creates watcher for Chaos Runner Pod
func WatchForRunnerPod(client client.Client, c controller.Controller) error {

	runnerPodHandler := handlerForRunnerPod(client)

	return c.Watch(&source.Kind{Type: &corev1.Pod{}}, &runnerPodHandler)
}

// WatchForMonitorPod creates watcher for Chaos Monitor Pod
func WatchForMonitorPod(client client.Client, c controller.Controller) error {

	monitorPodHandler := handlerForMonitorPod(client)

	return c.Watch(&source.Kind{Type: &corev1.Pod{}}, &monitorPodHandler)
}

// handlerForMonitorService creates a event handler for Monitor Service
func handlerForMonitorService(clientSet client.Client) handler.EnqueueRequestsFromMapFunc {
	reqLogger := chaosTypes.Log.WithName("Chaos Resources Watch")
	var monitorServiceRequest []reconcile.Request

	handlerForMonitorService := handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			monitorServiceCheck := strings.HasSuffix(a.Meta.GetName(), "-monitor")
			if monitorServiceCheck {
				svcLabels := a.Meta.GetLabels()
				engineUID := getPodEngineUIDLabel(svcLabels)
				listOptions := createListOptionsInNamespace(a.Meta.GetNamespace())
				listChaosEngine, err := getChaosEngineList(listOptions, clientSet)
				if err != nil {
					reqLogger.Error(err, "Unable to get the ChaosEngine Resources in namespace: %v", a.Meta.GetNamespace())
					return nil
				}
				monitorServiceRequest = handlerRequestFromEngineList(listChaosEngine, engineUID)
			}
			return monitorServiceRequest

		}),
	}
	return handlerForMonitorService

}

// handlerForMonitorPod creates a event Handler for Chaos Monitor Pod
func handlerForMonitorPod(clientSet client.Client) handler.EnqueueRequestsFromMapFunc {
	reqLogger := chaosTypes.Log.WithName("Chaos Resources Watch")
	var monitorPodRequest []reconcile.Request

	handlerForMonitorPod := handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			monitorNameCheck := strings.HasSuffix(a.Meta.GetName(), "-monitor")
			if monitorNameCheck {
				podLabels := a.Meta.GetLabels()

				engineUID := getPodEngineUIDLabel(podLabels)
				listOptions := createListOptionsInNamespace(a.Meta.GetNamespace())
				listChaosEngine, err := getChaosEngineList(listOptions, clientSet)
				if err != nil {
					reqLogger.Error(err, "Unable to get the ChaosEngine Resources in namespace: %v", a.Meta.GetNamespace())
					return nil
				}

				monitorPodRequest = handlerRequestFromEngineList(listChaosEngine, engineUID)

			}
			return monitorPodRequest
		}),
	}
	return handlerForMonitorPod
}

// handlerForRunnerPod creates a event Handler for Chaos Runner Pod
func handlerForRunnerPod(clientSet client.Client) handler.EnqueueRequestsFromMapFunc {
	reqLogger := chaosTypes.Log.WithName("Chaos Resources Watch")
	var runnerPodRequest []reconcile.Request

	handlerForRunner := handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			runnerNameCheck := strings.HasSuffix(a.Meta.GetName(), "-runner")
			if runnerNameCheck {
				podLabels := a.Meta.GetLabels()
				engineUID := getPodEngineUIDLabel(podLabels)
				listOptions := createListOptionsInNamespace(a.Meta.GetNamespace())
				listChaosEngine, err := getChaosEngineList(listOptions, clientSet)
				if err != nil {
					reqLogger.Error(err, "Unable to get the ChaosEngine Resources in namespace: %v", a.Meta.GetNamespace())
					return nil
				}
				runnerPodRequest = handlerRequestFromEngineList(listChaosEngine, engineUID)
			}
			return runnerPodRequest
		}),
	}
	return handlerForRunner
}

// handlerRequestFromEngineList initilize a event Watcher filtering the chaosEngine from the list.
func handlerRequestFromEngineList(listChaosEngine litmuschaosv1alpha1.ChaosEngineList, engineUID string) []reconcile.Request {
	for i := range listChaosEngine.Items {
		uuid := string(listChaosEngine.Items[i].GetUID())
		if engineUID == uuid {
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

func getPodEngineUIDLabel(podLabels map[string]string) string {
	var engineUID string
	if _, ok := podLabels["engineUID"]; ok {
		engineUID = podLabels["engineUID"]
	}
	return engineUID
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
