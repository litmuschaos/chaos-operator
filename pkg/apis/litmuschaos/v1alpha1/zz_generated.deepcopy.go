// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/rbac/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationParams) DeepCopyInto(out *ApplicationParams) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationParams.
func (in *ApplicationParams) DeepCopy() *ApplicationParams {
	if in == nil {
		return nil
	}
	out := new(ApplicationParams)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosEngine) DeepCopyInto(out *ChaosEngine) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosEngine.
func (in *ChaosEngine) DeepCopy() *ChaosEngine {
	if in == nil {
		return nil
	}
	out := new(ChaosEngine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChaosEngine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosEngineList) DeepCopyInto(out *ChaosEngineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChaosEngine, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosEngineList.
func (in *ChaosEngineList) DeepCopy() *ChaosEngineList {
	if in == nil {
		return nil
	}
	out := new(ChaosEngineList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChaosEngineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosEngineSpec) DeepCopyInto(out *ChaosEngineSpec) {
	*out = *in
	out.Appinfo = in.Appinfo
	in.Components.DeepCopyInto(&out.Components)
	if in.Experiments != nil {
		in, out := &in.Experiments, &out.Experiments
		*out = make([]ExperimentList, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosEngineSpec.
func (in *ChaosEngineSpec) DeepCopy() *ChaosEngineSpec {
	if in == nil {
		return nil
	}
	out := new(ChaosEngineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosEngineStatus) DeepCopyInto(out *ChaosEngineStatus) {
	*out = *in
	if in.Experiments != nil {
		in, out := &in.Experiments, &out.Experiments
		*out = make([]ExperimentStatuses, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosEngineStatus.
func (in *ChaosEngineStatus) DeepCopy() *ChaosEngineStatus {
	if in == nil {
		return nil
	}
	out := new(ChaosEngineStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosExperiment) DeepCopyInto(out *ChaosExperiment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosExperiment.
func (in *ChaosExperiment) DeepCopy() *ChaosExperiment {
	if in == nil {
		return nil
	}
	out := new(ChaosExperiment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChaosExperiment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosExperimentList) DeepCopyInto(out *ChaosExperimentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChaosExperiment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosExperimentList.
func (in *ChaosExperimentList) DeepCopy() *ChaosExperimentList {
	if in == nil {
		return nil
	}
	out := new(ChaosExperimentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChaosExperimentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosExperimentSpec) DeepCopyInto(out *ChaosExperimentSpec) {
	*out = *in
	in.Definition.DeepCopyInto(&out.Definition)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosExperimentSpec.
func (in *ChaosExperimentSpec) DeepCopy() *ChaosExperimentSpec {
	if in == nil {
		return nil
	}
	out := new(ChaosExperimentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosExperimentStatus) DeepCopyInto(out *ChaosExperimentStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosExperimentStatus.
func (in *ChaosExperimentStatus) DeepCopy() *ChaosExperimentStatus {
	if in == nil {
		return nil
	}
	out := new(ChaosExperimentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosResult) DeepCopyInto(out *ChaosResult) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosResult.
func (in *ChaosResult) DeepCopy() *ChaosResult {
	if in == nil {
		return nil
	}
	out := new(ChaosResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChaosResult) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosResultList) DeepCopyInto(out *ChaosResultList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChaosResult, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosResultList.
func (in *ChaosResultList) DeepCopy() *ChaosResultList {
	if in == nil {
		return nil
	}
	out := new(ChaosResultList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChaosResultList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosResultSpec) DeepCopyInto(out *ChaosResultSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosResultSpec.
func (in *ChaosResultSpec) DeepCopy() *ChaosResultSpec {
	if in == nil {
		return nil
	}
	out := new(ChaosResultSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChaosResultStatus) DeepCopyInto(out *ChaosResultStatus) {
	*out = *in
	out.ExperimentStatus = in.ExperimentStatus
	if in.ProbeStatus != nil {
		in, out := &in.ProbeStatus, &out.ProbeStatus
		*out = make([]ProbeStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChaosResultStatus.
func (in *ChaosResultStatus) DeepCopy() *ChaosResultStatus {
	if in == nil {
		return nil
	}
	out := new(ChaosResultStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CmdProbeInputs) DeepCopyInto(out *CmdProbeInputs) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CmdProbeInputs.
func (in *CmdProbeInputs) DeepCopy() *CmdProbeInputs {
	if in == nil {
		return nil
	}
	out := new(CmdProbeInputs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentParams) DeepCopyInto(out *ComponentParams) {
	*out = *in
	in.Runner.DeepCopyInto(&out.Runner)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentParams.
func (in *ComponentParams) DeepCopy() *ComponentParams {
	if in == nil {
		return nil
	}
	out := new(ComponentParams)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMap) DeepCopyInto(out *ConfigMap) {
	*out = *in
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMap.
func (in *ConfigMap) DeepCopy() *ConfigMap {
	if in == nil {
		return nil
	}
	out := new(ConfigMap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ENVPair) DeepCopyInto(out *ENVPair) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ENVPair.
func (in *ENVPair) DeepCopy() *ENVPair {
	if in == nil {
		return nil
	}
	out := new(ENVPair)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExperimentAttributes) DeepCopyInto(out *ExperimentAttributes) {
	*out = *in
	in.Components.DeepCopyInto(&out.Components)
	if in.Probe != nil {
		in, out := &in.Probe, &out.Probe
		*out = make([]ProbeAttributes, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExperimentAttributes.
func (in *ExperimentAttributes) DeepCopy() *ExperimentAttributes {
	if in == nil {
		return nil
	}
	out := new(ExperimentAttributes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExperimentComponents) DeepCopyInto(out *ExperimentComponents) {
	*out = *in
	if in.ENV != nil {
		in, out := &in.ENV, &out.ENV
		*out = make([]ExperimentENV, len(*in))
		copy(*out, *in)
	}
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]ConfigMap, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]Secret, len(*in))
		copy(*out, *in)
	}
	if in.ExperimentAnnotations != nil {
		in, out := &in.ExperimentAnnotations, &out.ExperimentAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.StatusCheckTimeouts = in.StatusCheckTimeouts
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExperimentComponents.
func (in *ExperimentComponents) DeepCopy() *ExperimentComponents {
	if in == nil {
		return nil
	}
	out := new(ExperimentComponents)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExperimentDef) DeepCopyInto(out *ExperimentDef) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Permissions != nil {
		in, out := &in.Permissions, &out.Permissions
		*out = make([]v1.PolicyRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ENVList != nil {
		in, out := &in.ENVList, &out.ENVList
		*out = make([]ENVPair, len(*in))
		copy(*out, *in)
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]ConfigMap, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]Secret, len(*in))
		copy(*out, *in)
	}
	if in.HostFileVolumes != nil {
		in, out := &in.HostFileVolumes, &out.HostFileVolumes
		*out = make([]HostFile, len(*in))
		copy(*out, *in)
	}
	if in.ExperimentAnnotations != nil {
		in, out := &in.ExperimentAnnotations, &out.ExperimentAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.SecurityContext.DeepCopyInto(&out.SecurityContext)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExperimentDef.
func (in *ExperimentDef) DeepCopy() *ExperimentDef {
	if in == nil {
		return nil
	}
	out := new(ExperimentDef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExperimentENV) DeepCopyInto(out *ExperimentENV) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExperimentENV.
func (in *ExperimentENV) DeepCopy() *ExperimentENV {
	if in == nil {
		return nil
	}
	out := new(ExperimentENV)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExperimentList) DeepCopyInto(out *ExperimentList) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExperimentList.
func (in *ExperimentList) DeepCopy() *ExperimentList {
	if in == nil {
		return nil
	}
	out := new(ExperimentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExperimentStatuses) DeepCopyInto(out *ExperimentStatuses) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExperimentStatuses.
func (in *ExperimentStatuses) DeepCopy() *ExperimentStatuses {
	if in == nil {
		return nil
	}
	out := new(ExperimentStatuses)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPProbeInputs) DeepCopyInto(out *HTTPProbeInputs) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPProbeInputs.
func (in *HTTPProbeInputs) DeepCopy() *HTTPProbeInputs {
	if in == nil {
		return nil
	}
	out := new(HTTPProbeInputs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostFile) DeepCopyInto(out *HostFile) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostFile.
func (in *HostFile) DeepCopy() *HostFile {
	if in == nil {
		return nil
	}
	out := new(HostFile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *K8sCommand) DeepCopyInto(out *K8sCommand) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new K8sCommand.
func (in *K8sCommand) DeepCopy() *K8sCommand {
	if in == nil {
		return nil
	}
	out := new(K8sCommand)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *K8sProbeInputs) DeepCopyInto(out *K8sProbeInputs) {
	*out = *in
	out.Command = in.Command
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new K8sProbeInputs.
func (in *K8sProbeInputs) DeepCopy() *K8sProbeInputs {
	if in == nil {
		return nil
	}
	out := new(K8sProbeInputs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProbeAttributes) DeepCopyInto(out *ProbeAttributes) {
	*out = *in
	out.K8sProbeInputs = in.K8sProbeInputs
	out.HTTPProbeInputs = in.HTTPProbeInputs
	out.CmdProbeInputs = in.CmdProbeInputs
	out.RunProperties = in.RunProperties
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProbeAttributes.
func (in *ProbeAttributes) DeepCopy() *ProbeAttributes {
	if in == nil {
		return nil
	}
	out := new(ProbeAttributes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProbeStatus) DeepCopyInto(out *ProbeStatus) {
	*out = *in
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProbeStatus.
func (in *ProbeStatus) DeepCopy() *ProbeStatus {
	if in == nil {
		return nil
	}
	out := new(ProbeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunProperty) DeepCopyInto(out *RunProperty) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunProperty.
func (in *RunProperty) DeepCopy() *RunProperty {
	if in == nil {
		return nil
	}
	out := new(RunProperty)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunnerInfo) DeepCopyInto(out *RunnerInfo) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RunnerAnnotation != nil {
		in, out := &in.RunnerAnnotation, &out.RunnerAnnotation
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunnerInfo.
func (in *RunnerInfo) DeepCopy() *RunnerInfo {
	if in == nil {
		return nil
	}
	out := new(RunnerInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Secret) DeepCopyInto(out *Secret) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Secret.
func (in *Secret) DeepCopy() *Secret {
	if in == nil {
		return nil
	}
	out := new(Secret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecurityContext) DeepCopyInto(out *SecurityContext) {
	*out = *in
	in.PodSecurityContext.DeepCopyInto(&out.PodSecurityContext)
	in.ContainerSecurityContext.DeepCopyInto(&out.ContainerSecurityContext)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecurityContext.
func (in *SecurityContext) DeepCopy() *SecurityContext {
	if in == nil {
		return nil
	}
	out := new(SecurityContext)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatusCheckTimeout) DeepCopyInto(out *StatusCheckTimeout) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatusCheckTimeout.
func (in *StatusCheckTimeout) DeepCopy() *StatusCheckTimeout {
	if in == nil {
		return nil
	}
	out := new(StatusCheckTimeout)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestStatus) DeepCopyInto(out *TestStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestStatus.
func (in *TestStatus) DeepCopy() *TestStatus {
	if in == nil {
		return nil
	}
	out := new(TestStatus)
	in.DeepCopyInto(out)
	return out
}
