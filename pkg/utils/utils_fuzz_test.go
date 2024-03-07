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
		if key != "" && value != "" && edUpdated != nil && len(edUpdated.ENV) != 1 {
			t.Errorf("expected env to be available, len %d", len(edUpdated.ENV))
		}
	})
}
