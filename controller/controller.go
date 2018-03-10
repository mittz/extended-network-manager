package controller

type ENMController interface {
	Init(message string)
	// Run(enwProvider provider.ENWProvider)
}

var (
	controllers map[string]ENMController
)

func GetController(name string) ENMController {
	if controller, ok := controllers[name]; ok {
		controller.Init("hello")
		return controller
	}
	return controllers["rancher"]
}
