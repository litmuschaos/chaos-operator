package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChaosEngineSpec defines the desired state of ChaosEngine
// +k8s:openapi-gen=true


// Describes a user-facing custom resource which is used by developers
// to create a chaos profile 
type ChaosEngineSpec struct {
        Appinfo        ApplicationParams    `json:"appinfo"`       //Appinfo contains deployment details of AUT
        Experiments    []ExperimentList     `json:"experiments"`   //Consists of experiments executed by the engine
        Schedule       ChaosSchedule        `json:"schedule"`      //Execution schedule of batch of chaos experiments
}

// ChaosEngineStatus defines the observed state of ChaosEngine
// +k8s:openapi-gen=true

// Derived information about status of experiments listed in the chaos engine
type ChaosEngineStatus struct {
        Experiments    []ExperimentStatuses   `json:"experiments"`   //Detailed status of individual experiments
}

// Information about Application-Under-Test (AUT) on the cluster
// Controller expects AUT to be annotated with litmuschaos.io/chaos: "true" to run chaos
type ApplicationParams struct {
        Appns          string               `json:"appns"`         //Namespace of the AUT 
        Applabel       string               `json:"applabel"`      //Unique label of the AUT
}

// Information about chaos experiments defined in the chaos engine
// These experiments are "pulled" as versioned charts from a "hub"
type ExperimentList struct {
        Name           string               `json:"name"`          //Name of the chaos experiment
        Spec           ExperimentAttributes `json:"spec"`          //Holds properties of an experiment listed in the engine 
}

// Information about schedule of chaos batch run
type ChaosSchedule struct {
        Interval       string               `json:"interval"`      //Period b/w two iterations of chaos experiments batch run 
        ExcludedTimes  string               `json:"excludedTimes"` //Time(s) of day when experiments batch run is not scheduled
        ExcludedDays   string               `json:"excludedDays"`  //Days of week when experiments batch run is not scheduled
        ConcurrencyPolicy string            `json:"concurrencyPolicy"` //Action upon schedule interval if older batch run is in progress 
}

// Detailed Information about experiments 
type ExperimentAttributes struct {
        Rank           uint32               `json:"rank"`          //Execution priority of the chaos experiment
        Component      ObjectUnderTest      `json:"component"`     //K8s, infra or app objects subjected to chaos
        Schedule       ExperimentSchedule   `json:"schedule"`      //Execution schedule of individual chaos experiment
}

// Detailed Information about component subjected to chaos in an experiment
// +optional
type ObjectUnderTest struct {
        Container      string               `json:"container"`     //Name of container under test in a pod 
        NWinterface    string               `json:"nwinterface"`   //Name of interface under test in a container
        Node           string               `json:"node"`          //Name of node under test in a K8s cluster
        PVC            string               `json:"pvc"`           //Name of persistent volume claim used by app
        Disk           string               `json:"disk"`          //Name of backend disk under test on a node
}

// Information about schedule of individual experiments
// +optional
type ExperimentSchedule struct {
        Interval       string               `json:"interval"`      //Period b/w two iterations of a specific experiment
        ExcludedTimes  string               `json:"excludedTimes"` //Time(s) of day when experiment is not scheduled
        ExcludedDays   string               `json:"excludedDays"`  //Days of week when experiment is not scheduled
        ConcurrencyPolicy string            `json:"concurrencyPolicy"` //Action upon schedule interval if older experiment is in progress 
}

// This information is immutable after engine has been created, fields are derived by kubernetes(operator)
type ExperimentStatuses struct {
       Name            string               `json:"name"`          //Name of experiment whose status is detailed
       Status          string               `json:"status"`        //Current state of chaos experiment 
       Verdict         string               `json:"verdict"`       //Result of a completed chaos experiment
       LastUpdateTime  metav1.Time          `json:"lastUpdateTime"`//Time of last state change of chaos experiment     
}

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
