package analytics

import (
	"crypto/rand"
	"fmt"
)

// UUIDGenerator creates a new UUID each time a new user triggers an event
func UUIDGenerator() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

var (
	// ClientUUID contains the UUID generated for the Google-Analytics
	ClientUUID = UUIDGenerator()
)
