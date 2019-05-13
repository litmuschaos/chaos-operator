package controller

import (
        /* @ksatchit: temp changes to facilitate operator-sdk build */
	//"github.com/litmuschaos/chaos-operator/pkg/controller/chaosengine"
	"github.com/ksatchit/chaos-operator/pkg/controller/chaosengine"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, chaosengine.Add)
}
