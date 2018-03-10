package provider

type ENMProvider interface {
	AddInterface()
	DeleteInterface()
}

func RegisterProvider(name string, provider ENMProvider) error {
	if providers == nil {
		providers = namke(map[string]ENMProvider)
	}
	if _, exists := providers[name]; exists {
		return
	}
}
