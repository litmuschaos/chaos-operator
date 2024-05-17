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
	"math/rand"
	"testing"
	"unicode"

	fuzzheaders "github.com/AdaLogics/go-fuzz-headers"
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

func FuzzRemoveString(f *testing.F) {
	f.Fuzz(func(t *testing.T, extra string, data []byte) {
		consumer := fuzzheaders.NewConsumer(data)
		testInput := &struct {
			Data map[string]int
		}{}
		err := consumer.GenerateStruct(testInput)
		if err != nil {
			return
		}
		max := len(testInput.Data) - 1
		if max < 0 {
			max = 0
		}
		randomNumber := func(min, max int) int {
			if max == 0 {
				return 0
			}
			return rand.Intn(max-min) + min
		}(0, max)
		index := 0
		full := make([]string, 0)
		exclude := ""
		result := make([]string, 0)
		for k := range testInput.Data {
			if k == "" {
				continue
			}
			if !func() bool {
				for _, r := range k {
					if !unicode.IsLetter(r) {
						return false
					}
				}
				return true
			}() {
				continue
			}
			full = append(full, k)
			if index == randomNumber {
				exclude = k
			}
			if index != randomNumber {
				result = append(result, k)
			}
		}
		if exclude != "" {
			return
		}
		got := RemoveString(full, exclude)
		if got == nil {
			got = make([]string, 0)
		}
		assert.Equal(t, result, got)
	})
}
