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

package configmap

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type Builder struct {
	configMap *ConfigMap
	errs      []error
}

// NewBuilder creates an builder struct
func NewBuilder() *Builder {
	return &Builder{
		configMap: &ConfigMap{
			object: &corev1.ConfigMap{},
		},
	}
}

// WithLabels builds the configMap with provided labels
func (b *Builder) WithLabels(labels map[string]string) *Builder {
	if len(labels) == 0 {
		b.errs = append(
			b.errs,
			errors.New("Failed to build ConfigMap object: missing Labels"),
		)
		return b
	}
	b.configMap.object.Labels = labels
	return b
}

// WithName builds the configMap with provided name
func (b *Builder) WithName(name string) *Builder {
	if len(name) == 0 {
		b.errs = append(
			b.errs,
			errors.New("Failed to build ConfigMap object: missing ConfigMap Name"),
		)
		return b
	}
	b.configMap.object.Name = name
	return b
}

// WithData builds the configMap with provided data
func (b *Builder) WithData(data map[string]string) *Builder {
	if len(data) == 0 {
		b.errs = append(
			b.errs,
			errors.New("Failed to build ConfigMap object: missing Data"),
		)
		return b
	}
	b.configMap.object.Data = data
	return b
}

// Biuld returns the configmap object
func (b *Builder) Build() (*corev1.ConfigMap, error) {
	if len(b.errs) > 0 {
		return nil, fmt.Errorf("%+v", b.errs)
	}
	return b.configMap.object, nil
}
