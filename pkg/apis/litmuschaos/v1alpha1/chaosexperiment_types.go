package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChaosExperimentSpec defines the desired state of ChaosExperiment
// +k8s:openapi-gen=true
// An experiment is the definition of a chaos test and is listed as an item
// in the chaos engine to be run against a given app.
type ChaosExperimentSpec struct {
        // ChaosGraph refers to the resource carrying low-level chaos params
        Chaosgraph    string                `json:"chaosgraph"`
        Components    ComponentUnderTest    `json:"components"` 
}

// ChaosExperimentStatus defines the observed state of ChaosExperiment
// +k8s:openapi-gen=true
type ChaosExperimentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// ComponentUnderTest defines information about component subjected to chaos in an experiment
type ComponentUnderTest struct {
        //Name of container under test in a pod
        Container      string               `json:"container"`
        //Name of interface under test in a container 
        NWinterface    string               `json:"nwinterface"`
        //Name of node under test in a K8s cluster
        Node           string               `json:"node"`
        //Name of persistent volume claim used by app
        PVC            string               `json:"pvc"`
        //Name of backend disk under test on a node
        Disk           string               `json:"disk"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosExperiment is the Schema for the chaosexperiments API
// +k8s:openapi-gen=true
type ChaosExperiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChaosExperimentSpec   `json:"spec,omitempty"`
	Status ChaosExperimentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosExperimentList contains a list of ChaosExperiment
type ChaosExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosExperiment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChaosExperiment{}, &ChaosExperimentList{})
}
