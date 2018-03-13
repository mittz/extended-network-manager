package main

import (
	"log"
	"os"

	"github.com/mittz/extended-network-manager/controller"
	"github.com/mittz/extended-network-manager/provider"
	"github.com/urfave/cli"
)

var (
	service           string
	enmControllerName string
	enmProviderName   string
	enmc              controller.ENMController
	enmp              provider.ENMProvider
)

// TODO: Adopting to CNI for Rancher(Cattle) / Kubernetes
func main() {
	app := cli.NewApp()

	app.Name = "extended-network-manager"
	app.Usage = "Rancher Extended Network Manager"
	app.Version = "0.1.0"
	app.UsageText = "extended-network-manager [command] [options]"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run rancher extended-network-manager service",
			Action: func(c *cli.Context) error {
				enmControllerName = c.String("controller")
				enmProviderName = c.String("provider")

				enmc = controller.GetController(enmControllerName)

				if enmc == nil {
					log.Fatalf("Unable to find controller by name %s", enmControllerName)
				}

				enmp = provider.GetProvider(enmProviderName)

				if enmp == nil {
					log.Fatalf("Unable to find provider by name %s", enmProviderName)
				}

				enmc.Run(enmp)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "controller",
					Value: "rancher",
					Usage: "controller name",
				},
				cli.StringFlag{
					Name:  "provider",
					Value: "pipework",
					Usage: "provicer name",
				},
			},
		},
	}

	app.Run(os.Args)
}
