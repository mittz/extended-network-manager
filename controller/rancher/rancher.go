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

	log.Println("rancher.go init()")

	controller.RegisterController(enmc.GetName(), enmc)
}

func (*ExntendedNetworkManagerController) Init(message string) {
	log.Println(message)
}

func (*ExntendedNetworkManagerController) Run() {
	log.Println("rancher run")
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
