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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChaosEngineSpec defines the desired state of ChaosEngine
// +k8s:openapi-gen=true
// ChaosEngineSpec describes a user-facing custom resource which is used by developers
// to create a chaos profile
type ChaosEngineSpec struct {
	//Appinfo contains deployment details of AUT
	Appinfo ApplicationParams `json:"appinfo"`
	//AnnotationCheck defines whether annotation check is allowed or not. It can be true or false
	AnnotationCheck string `json:"annotationCheck,omitempty"`
	//ChaosServiceAccount is the SvcAcc specified for chaos runner pods
	ChaosServiceAccount string `json:"chaosServiceAccount"`
	//Components contains the image of runnner and monitor pod
	Components ComponentParams `json:"components"`
	//Consists of experiments executed by the engine
	Experiments []ExperimentList `json:"experiments"`
	//Monitor Enable Status
	Monitoring bool `json:"monitoring,omitempty"`
	//JobCleanUpPolicy decides to retain or delete the jobs
	JobCleanUpPolicy CleanUpPolicy `json:"jobCleanUpPolicy,omitempty"`
	//AuxiliaryAppInfo contains details of dependent applications (infra chaos)
	AuxiliaryAppInfo string `json:"auxiliaryAppInfo,omitempty"`
	//EngineStatus is a requirement for validation
	EngineState EngineState `json:"engineState"`
}

// EngineState provides interface for all supported strings in spec.EngineState
type EngineState string

const (
	// EngineStateActive starts the reconcile call
	EngineStateActive EngineState = "active"
	// EngineStateStop stops the reconcile call
	EngineStateStop EngineState = "stop"
)

// EngineStatus provides interface for all supported strings in status.EngineStatus
type EngineStatus string

const (
	// EngineStatusInitialized is used for reconcile calls to start reconcile for creation
	EngineStatusInitialized EngineStatus = "initialized"
	// EngineStatusCompleted is used for reconcile calls to start reconcile for completion
	EngineStatusCompleted EngineStatus = "completed"
	// EngineStatusStopped is used for reconcile calls to start reconcile for delete
	EngineStatusStopped EngineStatus = "stopped"
)

// CleanUpPolicy defines the garbage collection method used by chaos-operator
type CleanUpPolicy string

const (
	//CleanUpPolicyDelete sets the garbage collection policy of chaos-operator to Delete Chaos Resources
	CleanUpPolicyDelete CleanUpPolicy = "delete"

	//CleanUpPolicyRetain sets the garbage collection policy of chaos-operator to Retain Chaos Resources
	CleanUpPolicyRetain CleanUpPolicy = "retain"
)

// ChaosEngineStatus defines the observed state of ChaosEngine
// +k8s:openapi-gen=true

// ChaosEngineStatus derives information about status of individual experiments
type ChaosEngineStatus struct {
	//
	EngineStatus EngineStatus `json:"engineStatus"`
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
	//kind of application
	AppKind string `json:"appkind"`
}

// ComponentParams defines information about the runner and monitor image
type ComponentParams struct {
	//Contains informations of the monitor pod
	Monitor MonitorInfo `json:"monitor"`
	//Contains informations of the the runner pod
	Runner RunnerInfo `json:"runner"`
}

// MonitorInfo defines the information of the monitor pod
type MonitorInfo struct {
	//Image of the monitor pod
	Image string `json:"image"`
}

// RunnerInfo defines the information of the runnerinfo pod
type RunnerInfo struct {
	//Image of the runner pod
	Image string `json:"image,omitempty"`
	//Type of runner
	Type string `json:"type,omitempty"`
	//Args of runner
	Args []string `json:"args,omitempty"`
	//Command for runner
	Command []string `json:"command,omitempty"`
	//ImagePullPolicy for runner pod
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

// ExperimentList defines information about chaos experiments defined in the chaos engine
// These experiments are "pulled" as versioned charts from a "hub"
type ExperimentList struct {
	//Name of the chaos experiment
	Name string `json:"name"`
	//Holds properties of an experiment listed in the engine
	Spec ExperimentAttributes `json:"spec"`
}

// ExperimentAttributes defines attributes of experiments
type ExperimentAttributes struct {
	//Execution priority of the chaos experiment
	Rank uint32 `json:"rank"`
	//Environment Varibles to override the default values in chaos-experiments
	Components ExperimentComponents `json:"components,omitempty"`
}

// ExperimentComponents contains ENV, Configmaps and Secrets
type ExperimentComponents struct {
	ENV        []ExperimentENV `json:"env,omitempty"`
	ConfigMaps []ConfigMap     `json:"configMaps,omitempty"`
	Secrets    []Secret        `json:"secrets,omitempty"`
}

// ExperimentENV varibles to override the default values in chaosexperiment
type ExperimentENV struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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

// ChaosEngineList contains a list of ChaosEngine
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ChaosEngineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosEngine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChaosEngine{}, &ChaosEngineList{})
}
