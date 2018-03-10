package rancher

import (
	"log"

	"github.com/mittz/extended-network-manager/controller"
)

func init() {
	enmc, err := NewExtendedNetworkManagerController()

	if err != nil {
		log.Fatalf("%v", err)
	}

	controller.RegisterController(enmc.GetName(), enmc)
}

func (*ExntendedNetworkManagerController) Init(message string) {
	log.Println(message)
}

func (*ExntendedNetworkManagerController) Run(provider provider.ENMProvider) {
	provider.
}

type ExntendedNetworkManagerController struct {
}

func (*ExntendedNetworkManagerController) GetName() string {
	return "rancher"
}

func NewExtendedNetworkManagerController() (*ExntendedNetworkManagerController, error) {
	enmc := &ExntendedNetworkManagerController{}
	return enmc, nil
}
