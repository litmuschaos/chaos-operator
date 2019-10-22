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
func CheckChaosAnnotation(engine chaosTypes.EngineInfo) error {
	// Use client-Go to obtain a list of apps w/ specified labels
	clientSet, err := k8s.CreateClientSet()
	if err != nil {
		return fmt.Errorf("clientset generation failed with error: %+v", err)
	}
	switch strings.ToLower(engine.AppInfo.Kind) {
	case "deployment", "deployments":
		err = CheckDeploymentAnnotation(clientSet, engine)
		if err != nil {
			return fmt.Errorf("no deployement found with required annotation, err: %+v", err)
		}
	default:
		return fmt.Errorf("resource type not supported for induce chaos")
	}
	return nil
}

// ValidateAnnotation will verify the validation require for induce chaos
func ValidateAnnotation(annotationValue string, chaosCandidates int) (int, error) {
	if annotationValue == ChaosAnnotationValue {
		chaosCandidates++
	} else if chaosCandidates > 1 {
		return chaosCandidates, errors.New("too many chaos candidates with same label, either provide unique labels or annotate only desired app for chaos")
	} else if chaosCandidates == 0 {
		return chaosCandidates, errors.New("no chaos-candidate found")
	}
	return chaosCandidates, nil
}
