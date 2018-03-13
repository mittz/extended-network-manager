package pipework

import (
	"log"
	"os/exec"

	"github.com/mittz/extended-network-manager/config"
	"github.com/mittz/extended-network-manager/provider"
)

var (
	PipeworkID string
)

func init() {
	enmp, err := NewExtendedNetworkManagerProvider()

	if err != nil {
		log.Fatal("%v", err)
	}

	err = provider.RegisterProvider(enmp.GetName(), enmp)

	if err != nil {
		log.Fatal("%v", err)
	}
}

type ExtendedNetworkManagerProvider struct {
}

func (*ExtendedNetworkManagerProvider) GetName() string {
	return "pipework"
}

func NewExtendedNetworkManagerProvider() (*ExtendedNetworkManagerProvider, error) {
	enmp := &ExtendedNetworkManagerProvider{}
	return enmp, nil
}

func (*ExtendedNetworkManagerProvider) ApplyConfig(pipeworkID string) error {
	PipeworkID = pipeworkID
	return nil
}

func (*ExtendedNetworkManagerProvider) AddInterface(cnc config.ContainerNetworkConfig) {
	// docker exec -it ${pipework_container_id} pipework ${host_interface} -i ${container_interface} ${container_id} ${ipaddress} ${macaddress}
	out, err := exec.Command("docker", "exec", PipeworkID, "pipework", cnc.HostInterface, "-i", cnc.Interface, cnc.ID, cnc.IPAddress, cnc.MACAddress).Output()

	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%s\n", out)
}

func (*ExtendedNetworkManagerProvider) DeleteInterface() {

}
