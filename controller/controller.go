package controller

import (
	"fmt"

	"github.com/mittz/extended-network-manager/provider"
)

type ENMController interface {
	Init(message string)
	Run(enwProvider provider.ENMProvider)
}

var (
	controllers map[string]ENMController
)

func GetController(name string) ENMController {
	if controller, ok := controllers[name]; ok {
		return controller
	}
	return nil
}

func RegisterController(name string, controller ENMController) error {
	if controllers == nil {
		controllers = make(map[string]ENMController)
	}
	if _, exists := controllers[name]; exists {
		return fmt.Errorf("controller: %s already registered", name)
	}
	controllers[name] = controller
	return nil
}
