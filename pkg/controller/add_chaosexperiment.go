package controller

import (
	"github.com/litmuschaos/chaos-operator/pkg/controller/chaosexperiment"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, chaosexperiment.Add)
}
