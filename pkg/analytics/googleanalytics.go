package analytics

import (
	"crypto/rand"
	"fmt"

	ga "github.com/jpillora/go-ogle-analytics"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// UUIDGenerator creates a new UUID each time a new user triggers an event
func UUIDGenerator() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logf.Log.Error(err, "UUID cannot be generated")
	}

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// TriggerAnalytics is reponsible for sending out events
func TriggerAnalytics() {

	client, err := ga.NewClient(GAclientID)
	if err != nil {
		logf.Log.Error(err, "GA Client ID Error")
	}
	uuid := UUIDGenerator()
	client.ClientID(uuid)

	err = client.Send(ga.NewEvent(CategoryLI, ActionI).Label(LabelO))
	if err != nil {
		logf.Log.Info("Unable to send GA event")
	}
	logf.Log.Info("Successful GA event sent !")
}
