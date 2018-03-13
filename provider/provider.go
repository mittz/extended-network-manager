package provider

import (
	"fmt"

	"github.com/mittz/extended-network-manager/config"
)

type ENMProvider interface {
	ApplyConfig(providerID string) error
	AddInterface(cnc config.ContainerNetworkConfig)
	DeleteInterface()
}

var (
	providers map[string]ENMProvider
)

func GetProvider(name string) ENMProvider {
	if provider, ok := providers[name]; ok {
		return provider
	}
	return nil
}

func RegisterProvider(name string, provider ENMProvider) error {
	if providers == nil {
		providers = make(map[string]ENMProvider)
	}
	if _, exists := providers[name]; exists {
		return fmt.Errorf("provider: %s already registered", name)
	}
	providers[name] = provider
	return nil
}
