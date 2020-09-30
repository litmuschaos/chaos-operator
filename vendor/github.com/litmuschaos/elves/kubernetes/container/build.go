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

package container

import (
	"errors"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
)

// Builder is the builder object for container
type Builder struct {
	con    *container // container instance
	errors []error    // errors found while building the container instance
}

// NewBuilder returns a new instance of builder
func NewBuilder() *Builder {
	return &Builder{
		con: &container{},
	}
}

// validate will run checks against container instance
func (b *Builder) validate() error {

	if len(b.errors) == 0 {
		return nil
	}
	return fmt.Errorf("Error while Validating Container Spec: %v", b.errors)
}

// Build returns the final kubernetes container
func (b *Builder) Build() (corev1.Container, error) {
	err := b.validate()
	if err != nil {
		return corev1.Container{}, err
	}
	return b.con.object, nil
}

// WithName sets the name of the container
func (b *Builder) WithName(name string) *Builder {
	if len(name) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: missing name"),
		)
		return b
	}
	b.con.object.Name = name
	return b
}

// WithImage sets the image of the container
func (b *Builder) WithImage(img string) *Builder {
	if len(img) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: missing image"),
		)
		return b
	}
	b.con.object.Image = img
	return b
}

// WithImagePullPolicy sets the image pull policy of the container
func (b *Builder) WithImagePullPolicy(policy corev1.PullPolicy) *Builder {
	if len(policy) == 0 {
		b.errors = append(
			b.errors,
			errors.New(
				"failed to build container object: missing imagepullpolicy",
			),
		)
		return b
	}

	b.con.object.ImagePullPolicy = policy
	return b
}

// WithCommandNew sets the command of the container
func (b *Builder) WithCommandNew(cmd []string) *Builder {
	if cmd == nil {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: nil command"),
		)
		return b
	}

	if len(cmd) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: missing command"),
		)
		return b
	}

	newcmd := []string{}
	newcmd = append(newcmd, cmd...)

	b.con.object.Command = newcmd
	return b
}

// WithArgumentsNew sets the command arguments of the container
func (b *Builder) WithArgumentsNew(args []string) *Builder {
	if args == nil {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: nil arguments"),
		)
		return b
	}

	if len(args) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: missing arguments"),
		)
		return b
	}

	newargs := []string{}
	newargs = append(newargs, args...)

	b.con.object.Args = newargs
	return b
}

// WithEnvsNew sets the envs of the container
func (b *Builder) WithEnvsNew(envs []corev1.EnvVar) *Builder {
	if envs == nil {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: nil envs"),
		)
		return b
	}

	if len(envs) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: missing envs"),
		)
		return b
	}

	newenvs := []corev1.EnvVar{}
	newenvs = append(newenvs, envs...)

	b.con.object.Env = newenvs
	return b
}

// WithPortsNew sets ports of the container
func (b *Builder) WithPortsNew(ports []corev1.ContainerPort) *Builder {
	if len(ports) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: missing ports"),
		)
		return b
	}

	newports := []corev1.ContainerPort{}
	newports = append(newports, ports...)

	b.con.object.Ports = newports
	return b
}

// WithVolumeMountsNew sets the command arguments of the container
func (b *Builder) WithVolumeMountsNew(volumeMounts []corev1.VolumeMount) *Builder {
	if volumeMounts == nil {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: nil volumemounts"),
		)
		return b
	}

	if len(volumeMounts) == 0 {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: missing volumemounts"),
		)
		return b
	}
	newvolumeMounts := []corev1.VolumeMount{}
	newvolumeMounts = append(newvolumeMounts, volumeMounts...)
	b.con.object.VolumeMounts = newvolumeMounts
	return b
}

// WithSecurityContext sets the security context of the container
func (b *Builder) WithSecurityContext(sc corev1.SecurityContext) *Builder {

	if reflect.DeepEqual(sc, corev1.SecurityContext{}) {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: empty security contexts"),
		)
		return b
	}

	b.con.object.SecurityContext = &sc
	return b
}

// WithResourceRequirements sets the resource requirements of the container
func (b *Builder) WithResourceRequirements(rr corev1.ResourceRequirements) *Builder {

	if reflect.DeepEqual(rr, corev1.ResourceRequirements{}) {
		b.errors = append(
			b.errors,
			errors.New("failed to build container object: empty resource requirements"),
		)
		return b
	}

	b.con.object.Resources = rr
	return b
}
