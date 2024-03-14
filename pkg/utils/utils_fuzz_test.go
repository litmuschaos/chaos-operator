/*
Copyright 2024 LitmusChaos Authors

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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func FuzzSetEnv(f *testing.F) {
	kv := map[string]string{
		"KEY1": "VALUE1",
		"KEY2": "VALUE2",
	}
	for k, v := range kv {
		f.Add(k, v)
	}
	f.Fuzz(func(t *testing.T, key, value string) {
		ed := ENVDetails{
			ENV: make([]v1.EnvVar, 0),
		}
		edUpdated := ed.SetEnv(key, value)
		if edUpdated == nil {
			t.Error("nil object not expected")
		}
		if key == "" && edUpdated != nil {
			assert.Equal(t, 0, len(edUpdated.ENV))
		}
		if value == "" && edUpdated != nil {
			assert.Equal(t, 0, len(edUpdated.ENV))
		}
		if key != "" && value != "" && edUpdated != nil {
			assert.Equal(t, 1, len(edUpdated.ENV))
		}
		if key != "" && value != "" && edUpdated != nil && len(edUpdated.ENV) == 1 {
			env := edUpdated.ENV[0]
			assert.Equal(t, key, env.Name)
			assert.Equal(t, value, env.Value)
		}
	})
}
