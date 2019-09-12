package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChaosEngineSpec defines the desired state of ChaosEngine
// +k8s:openapi-gen=true
// ChaosEngineSpec describes a user-facing custom resource which is used by developers
// to create a chaos profile
type ChaosEngineSpec struct {
	//Appinfo contains deployment details of AUT
	Appinfo ApplicationParams `json:"appinfo"`
	//ChaosServiceAccount is the SvcAcc specified for chaos runner pods
	ChaosServiceAccount string `json:"chaosServiceAccount"`
	//Consists of experiments executed by the engine
	Experiments []ExperimentList `json:"experiments"`
	//Execution schedule of batch of chaos experiments
	Schedule ChaosSchedule `json:"schedule"`
}

// ChaosEngineStatus defines the observed state of ChaosEngine
// +k8s:openapi-gen=true
// ChaosEngineStatus derives information about status of individual experiments
type ChaosEngineStatus struct {
	//Detailed status of individual experiments
	Experiments []ExperimentStatuses `json:"experiments"`
}

// ApplicationParams defines information about Application-Under-Test (AUT) on the cluster
// Controller expects AUT to be annotated with litmuschaos.io/chaos: "true" to run chaos
type ApplicationParams struct {
	//Namespace of the AUT
	Appns string `json:"appns"`
	//Unique label of the AUT
	Applabel string `json:"applabel"`
}

// ExperimentList defines information about chaos experiments defined in the chaos engine
// These experiments are "pulled" as versioned charts from a "hub"
type ExperimentList struct {
	//Name of the chaos experiment
	Name string `json:"name"`
	//Holds properties of an experiment listed in the engine
	Spec ExperimentAttributes `json:"spec"`
}

// ChaosSchedule defines information about schedule of chaos batch run
type ChaosSchedule struct {
	//Period b/w two iterations of chaos experiments batch run
	Interval string `json:"interval"`
	//Time(s) of day when experiments batch run is not scheduled
	ExcludedTimes string `json:"excludedTimes"`
	//Days of week when experiments batch run is not scheduled
	ExcludedDays string `json:"excludedDays"`
	//Action upon schedule interval if older batch run is in progress
	ConcurrencyPolicy string `json:"concurrencyPolicy"`
}

// ExperimentAttributes defines attributes of experiments
type ExperimentAttributes struct {
	//Execution priority of the chaos experiment
	Rank uint32 `json:"rank"`
	//K8s, infra or app objects subjected to chaos
	Components ObjectUnderTest `json:"components"`
	//Execution schedule of individual chaos experiment
	Schedule ExperimentSchedule `json:"schedule"`
}

// ObjectUnderTest defines information about component subjected to chaos in an experiment
// +optional
type ObjectUnderTest struct {
	//Name of container under test in a pod
	Container string `json:"container"`
	//Name of interface under test in a container
	NWinterface string `json:"nwinterface"`
	//Name of node under test in a K8s cluster
	Node string `json:"node"`
	//Name of persistent volume claim used by app
	PVC string `json:"pvc"`
	//Name of backend disk under test on a node
	Disk string `json:"disk"`
}

// ExperimentSchedule defines information about schedule of individual experiments
// +optional
type ExperimentSchedule struct {
	//Period b/w two iterations of a specific experiment
	Interval string `json:"interval"`
	//Time(s) of day when experiment is not scheduled
	ExcludedTimes string `json:"excludedTimes"`
	//Days of week when experiment is not scheduled
	ExcludedDays string `json:"excludedDays"`
	//Action upon schedule interval if older experiment is in progress
	ConcurrencyPolicy string `json:"concurrencyPolicy"`
}

// ExperimentStatuses defines information about status of individual experiments
// These fields are immutable, and are derived by kubernetes(operator)
type ExperimentStatuses struct {
	//Name of experiment whose status is detailed
	Name string `json:"name"`
	//Current state of chaos experiment
	Status string `json:"status"`
	//Result of a completed chaos experiment
	Verdict string `json:"verdict"`
	//Time of last state change of chaos experiment
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

// +genclient
// +resource:path=chaosengine
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosEngine is the Schema for the chaosengines API
// +k8s:openapi-gen=true
type ChaosEngine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChaosEngineSpec   `json:"spec,omitempty"`
	Status ChaosEngineStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosEngineList contains a list of ChaosEngine
type ChaosEngineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosEngine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChaosEngine{}, &ChaosEngineList{})
}
