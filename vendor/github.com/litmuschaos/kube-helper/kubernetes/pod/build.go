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

package pod

import (
	"errors"
	"fmt"

	"github.com/litmuschaos/kube-helper/kubernetes/container"

	corev1 "k8s.io/api/core/v1"
)

// Builder is the builder object for Pod
type Builder struct {
	pod  *Pod
	errs []error
}

// NewBuilder returns new instance of Builder
func NewBuilder() *Builder {
	return &Builder{pod: &Pod{object: &corev1.Pod{}}}
}

// WithName sets the Name field of Pod with provided value.
func (b *Builder) WithName(name string) *Builder {
	if len(name) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Pod object: missing Pod name"),
		)
		return b
	}
	b.pod.object.Name = name
	return b
}

// WithNamespace sets the Namespace field of Pod with provided value.
func (b *Builder) WithNamespace(namespace string) *Builder {
	if len(namespace) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Pod object: missing namespace"),
		)
		return b
	}
	b.pod.object.Namespace = namespace
	return b
}

// WithLabels sets the labels field of Pod with provided value
func (b *Builder) WithLabels(labels map[string]string) *Builder {
	if len(labels) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build pod object: missing labels"),
		)
		return b
	}

	if b.pod.object.Labels == nil {
		b.pod.object.Labels = map[string]string{}
	}

	for key, value := range labels {
		b.pod.object.Labels[key] = value
	}
	return b
}

// WithServiceAccountName sets the serviceaccountname field of Pod with provided value
func (b *Builder) WithServiceAccountName(serviceaccountname string) *Builder {
	if len(serviceaccountname) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Pod object: missing Pod serviceaccountname"),
		)
		return b
	}
	b.pod.object.Spec.ServiceAccountName = serviceaccountname
	return b
}

// WithRestartPolicy sets the restartpolicy field of Pod spec with provided value
func (b *Builder) WithRestartPolicy(restartpolicy corev1.RestartPolicy) *Builder {
	if len(restartpolicy) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Pod object: missing Pod restartpolicy"),
		)
		return b
	}
	b.pod.object.Spec.RestartPolicy = restartpolicy
	return b
}

// WithContainerBuilder adds a container to this pod object.
//
// NOTE:
//   container details are present in the provided container
// builder object
func (b *Builder) WithContainerBuilder(
	containerBuilder *container.Builder,
) *Builder {
	containerObj, err := containerBuilder.Build()
	if err != nil {
		b.errs = append(b.errs, fmt.Errorf("failed to build pod %v", err))
		return b
	}
	b.pod.object.Spec.Containers = append(
		b.pod.object.Spec.Containers,
		containerObj,
	)
	return b
}

// Build returns the Pod API instance
func (b *Builder) Build() (*corev1.Pod, error) {
	if len(b.errs) > 0 {
		return nil, fmt.Errorf("%+v", b.errs)
	}
	return b.pod.object, nil
}
