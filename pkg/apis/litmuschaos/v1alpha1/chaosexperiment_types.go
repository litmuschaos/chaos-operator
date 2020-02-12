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
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChaosExperimentSpec defines the desired state of ChaosExperiment
// +k8s:openapi-gen=true
// An experiment is the definition of a chaos test and is listed as an item
// in the chaos engine to be run against a given app.
type ChaosExperimentSpec struct {
	// Definition carries low-level chaos options
	Definition ExperimentDef `json:"definition"`
}

// ChaosExperimentStatus defines the observed state of ChaosExperiment
// +k8s:openapi-gen=true
type ChaosExperimentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// ConfigMap is an simpler implementation of corev1.ConfigMaps, needed for experiments
type ConfigMap struct {
	Data      map[string]string `json:"data,omitempty"`
	Name      string            `json:"name"`
	MountPath string            `json:"mountPath"`
}

// Secret is an simpler implementation of corev1.Secret
type Secret struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

// ExperimentDef defines information about nature of chaos & components subjected to it
type ExperimentDef struct {
	// Default labels of the runner pod
	// +optional
	Labels map[string]string `json:"labels"`
	// Image of the chaos executor
	Image string `json:"image"`
	//Scope specifies the service account scope (& thereby blast radius) of the experiment
	Scope string `json:"scope"`
	// List of Permission needed for a service account to execute experiment
	Permissions []rbacV1.PolicyRule `json:"permissions"`
	// List of ENV vars passed to executor pod
	ENVList []ENVPair `json:"env"`
	// Defines command to invoke experiment
	Command []string `json:"command"`
	// Defines arguments to runner's entrypoint command
	Args []string `json:"args"`
	// ConfigMaps contains a list of ConfigMaps
	ConfigMaps []ConfigMap `json:"configmaps,omitempty"`
	// Secrets contains a list of Secrets
	Secrets []Secret `json:"secrets,omitempty"`
}

// ENVPair defines env var list to hold chaos params
type ENVPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// +genclient
// +resource:path=chaosexperiment
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosExperiment is the Schema for the chaosexperiments API
// +k8s:openapi-gen=true
type ChaosExperiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChaosExperimentSpec   `json:"spec,omitempty"`
	Status ChaosExperimentStatus `json:"status,omitempty"`
}

// ChaosExperimentList contains a list of ChaosExperiment
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ChaosExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosExperiment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChaosExperiment{}, &ChaosExperimentList{})
}
