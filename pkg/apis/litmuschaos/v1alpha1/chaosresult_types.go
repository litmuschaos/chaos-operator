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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChaosResultSpec defines the desired state of ChaosResult
// +k8s:openapi-gen=true
// The chaosresult holds the status of a chaos experiment that is listed as an item
// in the chaos engine to be run against a given app.
type ChaosResultSpec struct {
	// EngineName defines the name of chaosEngine
	EngineName string `json:"engine,omitempty"`
	// ExperimentName defines the name of chaosexperiment
	ExperimentName string `json:"experiment"`
	// InstanceID defines the instance id
	InstanceID string `json:"instance,omitempty"`
}

// ChaosResultStatus defines the observed state of ChaosResult
// +k8s:openapi-gen=true
type ChaosResultStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// Definition carries low-level chaos options
	ExperimentStatus TestStatus `json:"experimentstatus"`
}

// TestStatus defines information about the status and results of a chaos experiment
type TestStatus struct {
	// Phase defines whether an experiment is running or completed
	Phase string `json:"phase"`
	// Verdict defines whether an experiment result is pass or fail
	Verdict string `json:"verdict"`
	// FailStep defines step where the experiments fails
	FailStep string `json:"failStep,omitempty"`
}

// +genclient
// +resource:path=chaosresult
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosResult is the Schema for the chaosresults API
// +k8s:openapi-gen=true
type ChaosResult struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChaosResultSpec   `json:"spec,omitempty"`
	Status ChaosResultStatus `json:"status,omitempty"`
}

// ChaosResultList contains a list of ChaosResult
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ChaosResultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosResult `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChaosResult{}, &ChaosResultList{})
}
