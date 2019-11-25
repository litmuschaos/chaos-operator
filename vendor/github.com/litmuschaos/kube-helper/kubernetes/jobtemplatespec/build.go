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

package jobtemplatespec

import (
	"errors"
	"fmt"

	batchv1beta1 "k8s.io/api/batch/v1beta1"

	jobspec "github.com/litmuschaos/kube-helper/kubernetes/jobspec"
)

// Builder is the builder object for JobTemplateSpec
type Builder struct {
	jobtemplatespec *JobTemplateSpec
	errs            []error
}

// NewBuilder returns new instance of Builder
func NewBuilder() *Builder {
	return &Builder{
		jobtemplatespec: &JobTemplateSpec{
			Object: &batchv1beta1.JobTemplateSpec{},
		},
	}
}

// WithJobSpecBuilder sets the spec of this jobtemplate
func (b *Builder) WithJobSpecBuilder(tmplbuilder *jobspec.Builder) *Builder {
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
	b.jobtemplatespec.Object.Spec = *jobspecObj.Object
	return b
}

// Build returns a jobspec object
func (b *Builder) Build() (*JobTemplateSpec, error) {
	if len(b.errs) > 0 {
		return nil, fmt.Errorf("%+v", b.errs)
	}
	return b.jobtemplatespec, nil
}
