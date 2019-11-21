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

package job

import (
	"errors"
	"fmt"

	jobspec "github.com/litmuschaos/kube-helper/kubernetes/jobspec"
	batchv1 "k8s.io/api/batch/v1"
)

// Builder is the builder object for Job
type Builder struct {
	job  *Job
	errs []error
}

// NewBuilder returns new instance of Builder
func NewBuilder() *Builder {
	return &Builder{
		job: &Job{
			object: &batchv1.Job{},
		},
	}
}

// WithName sets the Name field of Job with provided value.
func (b *Builder) WithName(name string) *Builder {
	if len(name) == 0 {
		b.errs = append(
			b.errs,
			errors.New("Failed to build Job object: missing Job Name"),
		)
		return b
	}
	b.job.object.Name = name
	return b
}

// WithNamespace sets the Namespace field of Job with provided value.
func (b *Builder) WithNamespace(namespace string) *Builder {
	if len(namespace) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Job object: missing namespace"),
		)
		return b
	}
	b.job.object.Namespace = namespace
	return b
}

// WithLabels sets the labels field of Job with provided value
func (b *Builder) WithLabels(labels map[string]string) *Builder {
	if len(labels) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Job object: missing labels"),
		)
		return b
	}

	if b.job.object.Labels == nil {
		b.job.object.Labels = map[string]string{}
	}

	for key, value := range labels {
		b.job.object.Labels[key] = value
	}
	return b
}

// WithJobSpecBuilder  sets the jobspec object to be used in this job
func (b *Builder) WithJobSpecBuilder(
	tmplbuilder *jobspec.Builder,
) *Builder {
	if tmplbuilder == nil {
		b.errs = append(
			b.errs,
			errors.New("failed to build job: nil templatespecbuilder"),
		)
		return b
	}

	jobspecObj, err := tmplbuilder.Build()

	if err != nil {
		b.errs = append(
			b.errs,
			errors.New(
				"failed to build job"),
		)
		return b
	}
	b.job.object.Spec = *jobspecObj.Object
	return b
}

// Build returns a job object
func (b *Builder) Build() (*batchv1.Job, error) {
	if len(b.errs) > 0 {
		return nil, fmt.Errorf("%+v", b.errs)
	}
	return b.job.object, nil
}
