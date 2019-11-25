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
	"strings"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
	k8s "github.com/litmuschaos/chaos-operator/pkg/kubernetes"
)

// Annotations on app to enable chaos on it
const (
	ChaosAnnotationKey   = "litmuschaos.io/chaos"
	ChaosAnnotationValue = "true"
)

// CheckChaosAnnotation will check for the annotation of required resources
func CheckChaosAnnotation(ce *chaosTypes.EngineInfo) (*chaosTypes.EngineInfo, error) {
	// Use client-Go to obtain a list of apps w/ specified labels
	//var chaosEngine chaosTypes.EngineInfo
	clientSet, err := k8s.CreateClientSet()
	if err != nil {
		return ce, fmt.Errorf("clientset generation failed with error: %+v", err)
	}
	switch strings.ToLower(ce.AppInfo.Kind) {
	case "deployment", "deployments":
		ce, err = CheckDeploymentAnnotation(clientSet, ce)
		if err != nil {
			return ce, fmt.Errorf("resource type 'deployment', err: %+v", err)
		}
	case "statefulset", "statefulsets":
		ce, err = CheckStatefulSetAnnotation(clientSet, ce)
		if err != nil {
			return ce, fmt.Errorf("resource type 'statefulset', err: %+v", err)
		}
	case "daemonset", "daemonsets":
		ce, err = CheckDaemonSetAnnotation(clientSet, ce)
		if err != nil {
			return ce, fmt.Errorf("resource type 'daemonset', err: %+v", err)
		}
	default:
		return ce, fmt.Errorf("resource type '%s' not supported for induce chaos", ce.AppInfo.Kind)
	}
	chaosTypes.Log.Info("chaos candidate of", "kind:", ce.AppInfo.Kind, "appName: ", ce.AppName, "appUUID: ", ce.AppUUID)
	return ce, nil
}

// CountTotalChaosEnabled will count the number of chaos enabled applications
func CountTotalChaosEnabled(annotationValue string, chaosCandidates int) int {
	if annotationValue == ChaosAnnotationValue {
		chaosCandidates++
	}
	return chaosCandidates
}
