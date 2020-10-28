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

package analytics

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"
)

// UUIDGenerator creates a new UUID each time a new user triggers an event
func UUIDGenerator() string {
	uuid := ""
	if strings.ToUpper(os.Getenv("ANALYTICS")) != "FALSE" {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return ""
		}
		uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	}
	return uuid
}

// ClientUUID contains the UUID generated for the Google-Analytics
var ClientUUID = UUIDGenerator()
