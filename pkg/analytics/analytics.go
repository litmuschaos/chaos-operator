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

const (
	// ClientID contains TrackingID of the application
	clientID = "UA-92076314-21"

	// supported event categories

	// Category category notifies installation of a component of Litmus Infrastructure
	category = "Litmus-Infra"

	// supported event actions

	// Action is sent when the installation is triggered
	action = "Installation"

	// supported event labels

	// Label denotes event is associated to which Litmus component
	label = "Chaos-Operator"
)

// UUIDGenerator creates a new UUID each time a new user triggers an event
func UUIDGenerator() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}

// TriggerAnalytics is reponsible for sending out events
func TriggerAnalytics() error {

	client, err := ga.NewClient(clientID)
	if err != nil {
		return fmt.Errorf("new client generation failed, error : %s", err)
	}
	uuid, err := UUIDGenerator()
	if err != nil {
		return fmt.Errorf("uuid generation failed, error: %s", err)
	}
	client.ClientID(uuid)

	err = client.Send(ga.NewEvent(category, action).Label(label))
	if err != nil {
		return fmt.Errorf("analytics event sending failed, error: %s", err)
	}
	return nil
}
