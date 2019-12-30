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

	ga "github.com/jpillora/go-ogle-analytics"
)

// UUIDGenerator creates a new UUID each time a new user triggers an event
func UUIDGenerator() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, err
}

// TriggerAnalytics is reponsible for sending out events
func TriggerAnalytics() (string, error) {

	// GAclientID contains TrackingID of the application
	GAclientID := "UA-92076314-21"

	// supported event categories

	// CategoryLI category notifies installation of a component of Litmus Infrastructure
	CategoryLI := "Litmus-Infra"

	// supported event actions

	// ActionI is sent when the installation is triggered
	ActionI := "Installation"

	// supported event labels

	// LabelO denotes event is associated to which Litmus component
	LabelO := "Chaos-Operator"

	client, err := ga.NewClient(GAclientID)
	if err != nil {
		return "GA Client ID Error", err
	}
	uuid, err := UUIDGenerator()
	if err != nil {
		return "UUID cannot be generated", err
	}
	client.ClientID(uuid)

	err = client.Send(ga.NewEvent(CategoryLI, ActionI).Label(LabelO))
	if err != nil {
		return "Unable to send GA event", err
	}
	l
}
