package resource

import (
	"errors"
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
	return ce, nil
}

// CountTotalChaosEnabled will count the number of chaos enabled applications
func CountTotalChaosEnabled(annotationValue string, chaosCandidates int) int {
	if annotationValue == ChaosAnnotationValue {
		chaosCandidates++
	}
	return chaosCandidates
}

// ValidateTotalChaosEnabled will validate the total chaos count
func ValidateTotalChaosEnabled(chaosCandidates int) error {
	if chaosCandidates > 1 {
		return errors.New("too many chaos candidates with same label, either provide unique labels or annotate only desired app for chaos")
	} else if chaosCandidates == 0 {
		return errors.New("no chaos-candidate found")
	}
	return nil
}