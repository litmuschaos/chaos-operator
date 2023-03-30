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
// ChaosEngineSpec describes a user-facing custom resource which is used by developers
// to create a chaos profile
type ChaosEngineSpec struct {
	//Appinfo contains the AUT details
	Appinfo ApplicationParams `json:"appinfo,omitempty"`
	//DefaultHealthCheck defines whether default health checks should be executed or not. It can be true or false
	// default value is true
	DefaultHealthCheck string `json:"defaultHealthCheck,omitempty"`
	//ChaosServiceAccount is the SvcAcc specified for chaos runner pods
	ChaosServiceAccount string `json:"chaosServiceAccount"`
	//Components contains the image, imagePullPolicy, arguments, and commands of runner
	Components ComponentParams `json:"components"`
	//Consists of experiments executed by the engine
	Experiments []ExperimentList `json:"experiments"`
	//JobCleanUpPolicy decides to retain or delete the jobs
	JobCleanUpPolicy CleanUpPolicy `json:"jobCleanUpPolicy,omitempty"`
	//AuxiliaryAppInfo contains details of dependent applications (infra chaos)
	AuxiliaryAppInfo string `json:"auxiliaryAppInfo,omitempty"`
	//EngineStatus is a requirement for validation
	EngineState EngineState `json:"engineState"`
	// TerminationGracePeriodSeconds contains terminationGracePeriod for the chaos resources
	TerminationGracePeriodSeconds int64 `json:"terminationGracePeriodSeconds,omitempty"`
	// Selectors contains the target application details
	Selectors *Selector `json:"selectors,omitempty"`
}

// EngineState provides interface for all supported strings in spec.EngineState
type EngineState string

const (
	// EngineStateActive starts the reconcile call
	EngineStateActive EngineState = "active"
	// EngineStateStop stops the reconcile call
	EngineStateStop EngineState = "stop"
)

// ExperimentStatus is typecasted to string for supporting the values below.
type ExperimentStatus string

const (
	// ExperimentStatusRunning is status of Experiment which is currently running
	ExperimentStatusRunning ExperimentStatus = "Running"
	// ExperimentStatusCompleted is status of Experiment which has been completed
	ExperimentStatusCompleted ExperimentStatus = "Completed"
	// ExperimentStatusWaiting is status of Experiment which will be executed via a Job
	ExperimentStatusWaiting ExperimentStatus = "Waiting for Job Creation"
	// ExperimentStatusNotFound is status of Experiment which is not found inside ChaosNamespace
	ExperimentStatusNotFound ExperimentStatus = "ChaosExperiment Not Found"
	// ExperimentStatusAborted is status of a Experiment is forcefully aborted
	ExperimentStatusAborted ExperimentStatus = "Forcefully Aborted"
	// ExperimentSkipped is status of Experiment which has been skipped
	ExperimentSkipped ExperimentStatus = "Skipped"
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
// ChaosEngineStatus derives information about status of individual experiments
type ChaosEngineStatus struct {
	//EngineStatus is a typed string to support limited values for ChaosEngine Status
	EngineStatus EngineStatus `json:"engineStatus"`
	//Detailed status of individual experiments
	Experiments []ExperimentStatuses `json:"experiments"`
}

// ApplicationParams defines information about Application-Under-Test (AUT) on the cluster
// Controller expects AUT to be annotated with litmuschaos.io/chaos: "true" to run chaos
type ApplicationParams struct {
	//Namespace of the AUT
	Appns string `json:"appns,omitempty"`
	//Unique label of the AUT
	Applabel string `json:"applabel,omitempty"`
	//kind of application
	AppKind string `json:"appkind,omitempty"`
}

type Selector struct {
	Workloads []Workload `json:"workloads,omitempty"`
	Pods      []Pod      `json:"pods,omitempty"`
}

type WorkloadKind string

const (
	WorkloadDeployment       WorkloadKind = "deployment"
	WorkloadStatefulSet      WorkloadKind = "statefulset"
	WorkloadDaemonSet        WorkloadKind = "daemonSet"
	WorkloadDeploymentConfig WorkloadKind = "deploymentconfig"
	WorkloadRollout          WorkloadKind = "rollout"
)

type Workload struct {
	Kind      WorkloadKind `json:"kind"`
	Namespace string       `json:"namespace"`
	Names     string       `json:"names,omitempty"`
	Labels    string       `json:"labels,omitempty"`
}

type Pod struct {
	Namespace string `json:"namespace"`
	Names     string `json:"names"`
}

// ComponentParams defines information about the runner
type ComponentParams struct {
	//Contains information of the runner pod
	Runner RunnerInfo `json:"runner"`
	// Contains information of the sidecar
	Sidecar []Sidecar `json:"sidecar,omitempty"`
}

type Sidecar struct {
	//Image of the sidecar container
	Image string `json:"image"`
	//ImagePullPolicy of the sidecar container
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Secrets for the sidecar container
	Secrets []Secret `json:"secrets,omitempty"`
	// EnvFrom for the sidecar container
	EnvFrom []corev1.EnvFromSource `json:"envFrom"`
	// ENV contains ENV passed to the sidecar container
	ENV []corev1.EnvVar `json:"env,omitempty"`
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
	//ImagePullSecrets for runner pod
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Runner Annotations that needs to be provided in the pod for pod that is getting created
	RunnerAnnotation map[string]string `json:"runnerAnnotations,omitempty"`
	// RunnerLabels labels for the runner pod
	RunnerLabels map[string]string `json:"runnerLabels,omitempty"`
	// NodeSelector for runner pod
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// ConfigMaps for runner pod
	ConfigMaps []ConfigMap `json:"configMaps,omitempty"`
	// Secrets for runner pod
	Secrets []Secret `json:"secrets,omitempty"`
	// Tolerations for runner pod
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Resource requirements for the runner pod
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
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
	// It contains env, configmaps, secrets, experimentImage, node selector, custom experiment annotation
	// which can be provided or overridden from the chaos engine
	Components ExperimentComponents `json:"components,omitempty"`
	// Probe contains details of probe, which can be applied on the experiments
	// Probe can be httpProbe, k8sProbe or cmdProbe
	Probe []ProbeAttributes `json:"probe,omitempty"`
}

// ProbeAttributes contains details of probe, which can be applied on the experiments
type ProbeAttributes struct {
	// Name of probe
	Name string `json:"name,omitempty"`
	// Type of probe
	Type string `json:"type,omitempty"`
	// inputs needed for the k8s probe
	K8sProbeInputs K8sProbeInputs `json:"k8sProbe/inputs,omitempty"`
	// inputs needed for the http probe
	HTTPProbeInputs HTTPProbeInputs `json:"httpProbe/inputs,omitempty"`
	// inputs needed for the cmd probe
	CmdProbeInputs CmdProbeInputs `json:"cmdProbe/inputs,omitempty"`
	// inputs needed for the prometheus probe
	PromProbeInputs PromProbeInputs `json:"promProbe/inputs,omitempty"`
	// inputs needed for the SLO probe
	SLOProbeInputs SLOProbeInputs `json:"sloProbe/inputs,omitempty"`
	// RunProperty contains timeout, retry and interval for the probe
	RunProperties RunProperty `json:"runProperties,omitempty"`
	// mode for k8s probe
	// it can be SOT, EOT, Edge
	Mode string `json:"mode,omitempty"`
	// Data contains the manifest/data for the resource, which need to be created
	// it supported for create operation only
	Data string `json:"data,omitempty"`
}

// K8sProbeInputs contains all the inputs required for k8s probe
type K8sProbeInputs struct {
	// group of the resource
	Group string `json:"group,omitempty"`
	// apiversion of the resource
	Version string `json:"version,omitempty"`
	// kind of resource
	Resource string `json:"resource,omitempty"`
	// ResourceNames to get the resources using their list of comma separated names
	ResourceNames string `json:"resourceNames,omitempty"`
	// namespace of the resource
	Namespace string `json:"namespace,omitempty"`
	// fieldselector to get the resource using fields selector
	FieldSelector string `json:"fieldSelector,omitempty"`
	// labelselector to get the resource using labels selector
	LabelSelector string `json:"labelSelector,omitempty"`
	// Operation performed by the k8s probe
	// it can be create, delete, present, absent
	Operation string `json:"operation,omitempty"`
}

// CmdProbeInputs contains all the inputs required for cmd probe
type CmdProbeInputs struct {
	// Command need to be executed for the probe
	Command string `json:"command,omitempty"`
	// Comparator check for the correctness of the probe output
	Comparator ComparatorInfo `json:"comparator,omitempty"`
	// The source where we have to run the command
	// It will run in inline(inside experiment itself) mode if source is nil
	Source SourceDetails `json:"source,omitempty"`
}

// SourceDetails contains source details of the cmdProbe
type SourceDetails struct {
	// Image for the source pod
	Image string `json:"image,omitempty"`
	// HostNetwork define the hostNetwork of the external pod
	// it supports boolean values and default value is false
	HostNetwork bool `json:"hostNetwork,omitempty"`
	// InheritInputs defined to inherit experiment pod attributes(ENV, volumes, and volumeMounts) into probe pod
	// it supports boolean values and default value is false
	InheritInputs bool `json:"inheritInputs,omitempty"`
	// Args for the source pod
	Args []string `json:"args,omitempty"`
	// ENVList contains ENV passed to the source pod
	ENVList []corev1.EnvVar `json:"env,omitempty"`
	// Labels for the source pod
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations for the source pod
	Annotations map[string]string `json:"annotations,omitempty"`
	// Command for the source pod
	Command []string `json:"command,omitempty"`
	// ImagePullPolicy for the source pod
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Privileged for the source pod
	Privileged bool `json:"privileged,omitempty"`
	// NodeSelector for the source pod
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Volumes for the source pod
	Volumes []corev1.Volume `json:"volumes,omitempty"`
	// VolumesMount for the source pod
	VolumesMount []corev1.VolumeMount `json:"volumeMount,omitempty"`
	//ImagePullSecrets for source pod
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// PromProbeInputs contains all the inputs required for prometheus probe
type PromProbeInputs struct {
	// Endpoint for the prometheus probe
	Endpoint string `json:"endpoint,omitempty"`
	// Query to get promethus metrics
	Query string `json:"query,omitempty"`
	// QueryPath contains filePath, which contains prometheus query
	QueryPath string `json:"queryPath,omitempty"`
	// Comparator check for the correctness of the probe output
	Comparator ComparatorInfo `json:"comparator,omitempty"`
}

// SLOProbeInputs contains all the inputs required for SLO probe
type SLOProbeInputs struct {
	// PlatformEndpoint for the monitoring service endpoint
	PlatformEndpoint string `json:"platformEndpoint,omitempty"`
	// SLOIdentifier for fetching the details of the SLO
	SLOIdentifier string `json:"sloIdentifier,omitempty"`
	// EvaluationWindow is the time period for which the metrics will be evaluated
	EvaluationWindow EvaluationWindow `json:"evaluationWindow,omitempty"`
	// SLOSourceMetadata consists of required metadata details to fetch metric data
	SLOSourceMetadata SLOSourceMetadata `json:"sloSourceMetadata,omitempty"`
	// Comparator check for the correctness of the probe output
	Comparator ComparatorInfo `json:"comparator,omitempty"`
}

// EvaluationWindow is the time period for which the SLO probe will work
type EvaluationWindow struct {
	// Start time of evaluation
	EvaluationStartTime int `json:"evaluationStartTime,omitempty"`
	// End time of evaluation
	EvaluationEndTime int `json:"evaluationEndTime,omitempty"`
}

type SLOSourceMetadata struct {
	// APITokenSecret for authenticating with the platform service
	APITokenSecret string `json:"apiTokenSecret,omitempty"`
	// SecretNamespace where the APITokenSecret is present
	SecretNamespace string `json:"secretNamespace,omitempty"`
	// Scope required for fetching details
	Scope Identifier `json:"scope,omitempty"`
}

// Identifier required for fetching details from the Platform APIs
type Identifier struct {
	// AccountIdentifier for account ID
	AccountIdentifier string `json:"accountIdentifier,omitempty"`
	// OrgIdentifier for organization ID
	OrgIdentifier string `json:"orgIdentifier,omitempty"`
	// ProjectIdentifier for project ID
	ProjectIdentifier string `json:"projectIdentifier,omitempty"`
}

// ComparatorInfo contains the comparator details
type ComparatorInfo struct {
	// Type of data
	// it can be int, float, string
	Type string `json:"type,omitempty"`
	// Criteria for matching data
	// it supports >=, <=, ==, >, <, != for int and float
	// it supports equal, notEqual, contains for string
	Criteria string `json:"criteria,omitempty"`
	// Value contains relative value for criteria
	Value string `json:"value,omitempty"`
}

// HTTPProbeInputs contains all the inputs required for http probe
type HTTPProbeInputs struct {
	// URL which needs to curl, to check the status
	URL string `json:"url,omitempty"`
	// InsecureSkipVerify flag to skip certificate checks
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
	// Method define the http method, it can be get or post
	Method HTTPMethod `json:"method,omitempty"`
}

// HTTPMethod define the http method details
type HTTPMethod struct {
	Get  GetMethod  `json:"get,omitempty"`
	Post PostMethod `json:"post,omitempty"`
}

// GetMethod define the http Get method
type GetMethod struct {
	// Criteria for matching data
	// it supports  == != operations
	Criteria string `json:"criteria,omitempty"`
	// Value contains relative value for criteria
	ResponseCode string `json:"responseCode,omitempty"`
}

// PostMethod define the http Post method
type PostMethod struct {
	// ContentType contains content type for http body data
	ContentType string `json:"contentType,omitempty"`
	// Body contains http body for post request
	Body string `json:"body,omitempty"`
	// BodyPath contains filePath, which contains http body
	BodyPath string `json:"bodyPath,omitempty"`
	// Criteria for matching data
	// it supports  == != operations
	Criteria string `json:"criteria,omitempty"`
	// Value contains relative value for criteria
	ResponseCode string `json:"responseCode,omitempty"`
}

// RunProperty contains timeout, retry and interval for the probe
type RunProperty struct {
	//ProbeTimeout contains timeout for the probe
	ProbeTimeout int `json:"probeTimeout,omitempty"`
	// Interval contains the interval for the probe
	Interval int `json:"interval,omitempty"`
	// Retry contains the retry count for the probe
	Retry int `json:"retry,omitempty"`
	// Attempt contains the total attempt count for the probe
	Attempt int `json:"attempt,omitempty"`
	//ProbePollingInterval contains time interval, for which continuous probe should be sleep
	// after each iteration
	ProbePollingInterval int `json:"probePollingInterval,omitempty"`
	//InitialDelaySeconds time interval for which probe will wait before run
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"`
	// StopOnFailure contains flag to stop/continue experiment execution, if probe fails
	// it will stop the experiment execution, if provided true
	// it will continue the experiment execution, if provided false or not provided(default case)
	StopOnFailure bool `json:"stopOnFailure,omitempty"`
}

// ExperimentComponents contains ENV, Configmaps and Secrets
type ExperimentComponents struct {
	ENV                        []corev1.EnvVar               `json:"env,omitempty"`
	ConfigMaps                 []ConfigMap                   `json:"configMaps,omitempty"`
	Secrets                    []Secret                      `json:"secrets,omitempty"`
	ExperimentAnnotations      map[string]string             `json:"experimentAnnotations,omitempty"`
	ExperimentImage            string                        `json:"experimentImage,omitempty"`
	ExperimentImagePullSecrets []corev1.LocalObjectReference `json:"experimentImagePullSecrets,omitempty"`
	NodeSelector               map[string]string             `json:"nodeSelector,omitempty"`
	StatusCheckTimeouts        StatusCheckTimeout            `json:"statusCheckTimeouts,omitempty"`
	Resources                  corev1.ResourceRequirements   `json:"resources,omitempty"`
	Tolerations                []corev1.Toleration           `json:"tolerations,omitempty"`
}

// StatusCheckTimeout contains Delay and timeouts for the status checks
type StatusCheckTimeout struct {
	Delay   int `json:"delay,omitempty"`
	Timeout int `json:"timeout,omitempty"`
}

// ExperimentStatuses defines information about status of individual experiments
// These fields are immutable, and are derived by kubernetes(operator)
type ExperimentStatuses struct {
	//Name of the chaos experiment
	Name string `json:"name"`
	//Name of chaos-runner pod managing this experiment
	Runner string `json:"runner"`
	//Name of experiment pod executing the chaos
	ExpPod string `json:"experimentPod"`
	//Current state of chaos experiment
	Status ExperimentStatus `json:"status"`
	//Result of a completed chaos experiment
	Verdict string `json:"verdict"`
	//Time of last state change of chaos experiment
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

// +genclient
// +resource:path=chaosengine
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChaosEngine is the Schema for the chaosengines API
type ChaosEngine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChaosEngineSpec   `json:"spec,omitempty"`
	Status ChaosEngineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChaosEngineList contains a list of ChaosEngine
type ChaosEngineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosEngine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChaosEngine{}, &ChaosEngineList{})
}
