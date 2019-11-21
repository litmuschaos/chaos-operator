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

package jobspec

import (
	"errors"
	"fmt"

	templatespec "github.com/litmuschaos/kube-helper/kubernetes/podtemplatespec"
	batchv1 "k8s.io/api/batch/v1"
)

// Builder is the builder object for JobSpec
type Builder struct {
	jobspec *JobSpec
	errs    []error
}

// NewBuilder returns new instance of Builder
func NewBuilder() *Builder {
	return &Builder{
		jobspec: &JobSpec{
			Object: &batchv1.JobSpec{},
		},
	}
}

// WithBackOffLimit sets the number of retries before marking this job failed
func (b *Builder) WithBackOffLimit(backoff *int32) *Builder {
	if int(*backoff) < 0 {
		b.errs = append(
			b.errs,
			errors.New("failed to build Job: invalid backofflimit "),
		)
		return b
	}

	b.jobspec.Object.BackoffLimit = backoff
	return b
}

// WithPodTemplateSpecBuilder sets the template of pod to be created by this job
func (b *Builder) WithPodTemplateSpecBuilder(
	tmplbuilder *templatespec.Builder,
) *Builder {
	if tmplbuilder == nil {
		b.errs = append(
			b.errs,
			errors.New("failed to build job: nil templatespecbuilder"),
		)
		return b
	}

	templatespecObj, err := tmplbuilder.Build()

	if err != nil {
		b.errs = append(
			b.errs,
			errors.New(
				"failed to build job"),
		)
		return b
	}
	b.jobspec.Object.Template = *templatespecObj.Object
	return b
}

// Build returns a jobspec object
func (b *Builder) Build() (*JobSpec, error) {
	if len(b.errs) > 0 {
		return nil, fmt.Errorf("%+v", b.errs)
	}
	return b.jobspec, nil
}
