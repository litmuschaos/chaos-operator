package analytics

import (
	// "github.com/go-logr/logr"
	ga "github.com/jpillora/go-ogle-analytics"
)

// Test ashdihadihai
func Test() {

	client, err := ga.NewClient("UA-127388617-2")
	if err != nil {
		panic(err)
	}
	client.ClientID("6f460195-4c50-4150-a3ae-4683dda3ae23")
	err = client.Send(ga.NewEvent("Install", "ChaosOperator").Label("AppName"))
	if err != nil {
		panic(err)
	}
	println("Event fired!")

}
