package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChaosGraphSpec defines the desired state of ChaosGraph
// +k8s:openapi-gen=true
// ChaosGraph defines the low-level chaos options and the executor infrastructure
type ChaosGraphSpec struct {
        // Defines the chaos executor framework/infrastructure to run the experiment
        Executor     string              `json:"executor"`
        // Defines the chaos parameters and execution artifacts 
        Definition   ChaosRunDefinition  `json:"definition"`
}

// ChaosRunDefinition defines the information fed to the executor framework
// to enable successful chaos injection
type ChaosRunDefinition struct {
        // Default labels of the executor pod
        // +optional
        Labels       map[string]string   `json:"labels"`
        // List of ENV vars passed to executor pod
        ENVList      []ENVPair           `json:"env"`
        // Defines command to invoke experiment
        Command	     []string            `json:"command"`
        // Defines arguments to executor's entrypoint command 
        Args         []string            `json:"args"`
}

// EnvPair defines env var list to hold chaos params
type ENVPair struct {
        Name         string              `json:"name"`
        Value        string              `json:"value"`
}

// ChaosGraphStatus defines the observed state of ChaosGraph
// +k8s:openapi-gen=true
type ChaosGraphStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosGraph is the Schema for the chaostemplates API
// +k8s:openapi-gen=true
type ChaosGraph struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChaosGraphSpec   `json:"spec,omitempty"`
	Status ChaosGraphStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosGraphList contains a list of ChaosGraph
type ChaosGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosGraph `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChaosGraph{}, &ChaosGraphList{})
}
