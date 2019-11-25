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

package cronjob

import (
	"errors"
	"fmt"

	batchv1beta1 "k8s.io/api/batch/v1beta1"

	jobtemplatespec "github.com/litmuschaos/kube-helper/kubernetes/jobtemplatespec"
)

// Builder is the builder object for CronJob
type Builder struct {
	cronjob *CronJob
	errs    []error
}

// NewBuilder returns new instance of Builder
func NewBuilder() *Builder {
	return &Builder{
		cronjob: &CronJob{
			object: &batchv1beta1.CronJob{},
		},
	}
}

// WithName sets the Name field of CronJob with provided value.
func (b *Builder) WithName(name string) *Builder {
	if len(name) == 0 {
		b.errs = append(
			b.errs,
			errors.New("Failed to build Job object: missing Job Name"),
		)
		return b
	}
	b.cronjob.object.Name = name
	return b
}

// WithNamespace sets the Namespace field of CronJob with provided value.
func (b *Builder) WithNamespace(namespace string) *Builder {
	if len(namespace) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Job object: missing namespace"),
		)
		return b
	}
	b.cronjob.object.Namespace = namespace
	return b
}

// WithLabels sets the labels field of CronJob with provided value
func (b *Builder) WithLabels(labels map[string]string) *Builder {
	if len(labels) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Job object: missing labels"),
		)
		return b
	}

	if b.cronjob.object.Labels == nil {
		b.cronjob.object.Labels = map[string]string{}
	}

	for key, value := range labels {
		b.cronjob.object.Labels[key] = value
	}
	return b
}

// WithSchedule sets the Schedule field of CronJob with provided value.
func (b *Builder) WithSchedule(schedule string) *Builder {
	if len(schedule) == 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Job object: missing schedule"),
		)
		return b
	}
	b.cronjob.object.Spec.Schedule = schedule
	return b
}

// WithSuccessfulJobHistoryLimit sets the SuccessfulJobHistoryLimit field of CronJob with provided value.
func (b *Builder) WithSuccessfulJobHistoryLimit(limit *int32) *Builder {
	if int(*limit) < 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Job: invalid successfulJobHistoryLimit "),
		)
		return b
	}
	b.cronjob.object.Spec.SuccessfulJobsHistoryLimit = limit
	return b
}

// WithFailedJobHistoryLimit sets the FailedJobHistoryLimit field of CronJob with provided value.
func (b *Builder) WithFailedJobHistoryLimit(limit *int32) *Builder {
	if int(*limit) < 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Job: invalid failedJobHistoryLimit "),
		)
		return b
	}
	b.cronjob.object.Spec.FailedJobsHistoryLimit = limit
	return b
}

// WithJobTemplateSpecBuilder sets the jobtemplate of this cronjob
func (b *Builder) WithJobTemplateSpecBuilder(tmplbuilder *jobtemplatespec.Builder) *Builder {
	if tmplbuilder == nil {
		b.errs = append(
			b.errs,
			errors.New("failed to build job: nil templatespecbuilder"),
		)
		return b
	}

	jobtemplatespecObj, err := tmplbuilder.Build()

	if err != nil {
		b.errs = append(
			b.errs,
			errors.New(
				"failed to build job"),
		)
		return b
	}
	b.cronjob.object.Spec.JobTemplate = *jobtemplatespecObj.Object
	return b
}

// Build returns a job object
func (b *Builder) Build() (*batchv1beta1.CronJob, error) {
	if len(b.errs) > 0 {
		return nil, fmt.Errorf("%+v", b.errs)
	}
	return b.cronjob.object, nil
}
