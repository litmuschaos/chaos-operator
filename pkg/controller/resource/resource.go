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
	"os"
	"strings"

	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// Annotations on app to enable chaos on it
const (
	ChaosAnnotationValue      = "true"
	DefaultChaosAnnotationKey = "litmuschaos.io/chaos"
)

var (
	// ChaosAnnotationKey is global variable used as the Key for annotation check.
	ChaosAnnotationKey = GetAnnotationKey()
)

// GetAnnotationKey returns the annotation to be used while validating applications.
func GetAnnotationKey() string {

	annotationKey := os.Getenv("CUSTOM_ANNOTATION")
	if len(annotationKey) != 0 {
		return annotationKey
	}
	return DefaultChaosAnnotationKey

}

// CheckChaosAnnotation will check for the annotation of required resources
func CheckChaosAnnotation(engine *chaosTypes.EngineInfo, clientset kubernetes.Interface, dynamicClientSet dynamic.Interface) (*chaosTypes.EngineInfo, error) {

	switch strings.ToLower(engine.AppInfo.Kind) {
	case "deployment", "deployments":
		engine, err := CheckDeploymentAnnotation(clientset, engine)
		if err != nil {
			return engine, fmt.Errorf("resource type 'deployment', err: %+v", err)
		}
	case "statefulset", "statefulsets":
		engine, err := CheckStatefulSetAnnotation(clientset, engine)
		if err != nil {
			return engine, fmt.Errorf("resource type 'statefulset', err: %+v", err)
		}
	case "daemonset", "daemonsets":
		engine, err := CheckDaemonSetAnnotation(clientset, engine)
		if err != nil {
			return engine, fmt.Errorf("resource type 'daemonset', err: %+v", err)
		}
	case "deploymentconfig", "deploymentconfigs":
		engine, err := CheckDeploymentConfigAnnotation(dynamicClientSet, engine)
		if err != nil {
			return engine, fmt.Errorf("resource type 'deploymentconfig', err: %+v", err)
		}
	case "rollout", "rollouts":
		engine, err := CheckRolloutAnnotation(dynamicClientSet, engine)
		if err != nil {
			return engine, fmt.Errorf("resource type 'rollout', err: %+v", err)
		}
	default:
		return engine, fmt.Errorf("resource type '%s' not supported for induce chaos", engine.AppInfo.Kind)
	}
	chaosTypes.Log.Info("chaos candidate of", "kind:", engine.AppInfo.Kind, "appName: ", engine.AppName, "appUUID: ", engine.AppUUID)
	return engine, nil
}

// CountTotalChaosEnabled will count the number of chaos enabled applications
func CountTotalChaosEnabled(annotationValue string, chaosCandidates int) int {
	if annotationValue == ChaosAnnotationValue {
		chaosCandidates++
	}
	return chaosCandidates
}
