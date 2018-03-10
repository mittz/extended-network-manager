package pipework

import "log"

func init() {
	enmp, err := NewExtendedNetworkManagerProvider()

	if err != nil {
		log.Fatal("%v", err)
	}

	provider.RegisterProvider(enmp.GetName(), enmp)
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
