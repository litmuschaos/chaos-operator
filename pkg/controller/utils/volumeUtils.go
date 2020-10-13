package utils

import (
	"github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	volume "github.com/litmuschaos/elves/kubernetes/volume/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var (
	// hostpathTypeFile represents the hostpath type
	hostpathTypeFile = corev1.HostPathFile
)

// CreateVolumeBuilders build Volume needed in execution of experiments
func CreateVolumeBuilders(configMaps []v1alpha1.ConfigMap, secrets []v1alpha1.Secret) []*volume.Builder {
	volumeBuilderList := []*volume.Builder{}

	volumeBuilderForConfigMaps := BuildVolumeBuilderForConfigMaps(configMaps)
	volumeBuilderList = append(volumeBuilderList, volumeBuilderForConfigMaps...)

	volumeBuilderForSecrets := BuildVolumeBuilderForSecrets(secrets)
	volumeBuilderList = append(volumeBuilderList, volumeBuilderForSecrets...)

	return volumeBuilderList
}

// CreateVolumeMounts mounts Volume needed in execution of experiments
func CreateVolumeMounts(configMaps []v1alpha1.ConfigMap, secrets []v1alpha1.Secret) []corev1.VolumeMount {

	var volumeMountsList []corev1.VolumeMount

	volumeMountsListForConfigMaps := BuildVolumeMountsForConfigMaps(configMaps)
	volumeMountsList = append(volumeMountsList, volumeMountsListForConfigMaps...)

	volumeMountsListForSecrets := BuildVolumeMountsForSecrets(secrets)
	volumeMountsList = append(volumeMountsList, volumeMountsListForSecrets...)

	return volumeMountsList
}

// VolumeOperations filles up VolumeOpts strucuture
func (volumeOpts *VolumeOpts) VolumeOperations(configMaps []v1alpha1.ConfigMap, secrets []v1alpha1.Secret) {
	volumeOpts.VolumeBuilders = CreateVolumeBuilders(configMaps, secrets)
	volumeOpts.VolumeMounts = CreateVolumeMounts(configMaps, secrets)
}

// BuildVolumeMountsForConfigMaps builds VolumeMounts for ConfigMaps
func BuildVolumeMountsForConfigMaps(configMaps []v1alpha1.ConfigMap) []corev1.VolumeMount {
	var volumeMountsList []corev1.VolumeMount
	for _, v := range configMaps {
		var volumeMount corev1.VolumeMount
		volumeMount.Name = v.Name
		volumeMount.MountPath = v.MountPath
		volumeMountsList = append(volumeMountsList, volumeMount)
	}
	return volumeMountsList
}

// BuildVolumeMountsForSecrets builds VolumeMounts for Secrets
func BuildVolumeMountsForSecrets(secrets []v1alpha1.Secret) []corev1.VolumeMount {
	var volumeMountsList []corev1.VolumeMount
	for _, v := range secrets {
		var volumeMount corev1.VolumeMount
		volumeMount.Name = v.Name
		volumeMount.MountPath = v.MountPath
		volumeMountsList = append(volumeMountsList, volumeMount)
	}
	return volumeMountsList
}

// BuildVolumeBuilderForConfigMaps builds VolumeBuilders for ConfigMaps
func BuildVolumeBuilderForConfigMaps(configMaps []v1alpha1.ConfigMap) []*volume.Builder {
	volumeBuilderList := []*volume.Builder{}
	if configMaps == nil {
		return nil
	}
	for _, v := range configMaps {
		volumeBuilder := volume.NewBuilder().
			WithConfigMap(v.Name)
		volumeBuilderList = append(volumeBuilderList, volumeBuilder)
	}
	return volumeBuilderList
}

// BuildVolumeBuilderForSecrets builds VolumeBuilders for Secrets
func BuildVolumeBuilderForSecrets(secrets []v1alpha1.Secret) []*volume.Builder {
	volumeBuilderList := []*volume.Builder{}
	if secrets == nil {
		return nil
	}
	for _, v := range secrets {
		volumeBuilder := volume.NewBuilder().
			WithSecret(v.Name)
		volumeBuilderList = append(volumeBuilderList, volumeBuilder)
	}
	return volumeBuilderList
}
