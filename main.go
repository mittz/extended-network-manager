package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	service           string
	ENMControllerName string
	enmc              controller.ENMController
)

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
				ENMControllerName = c.String("controller")
				enmc = controller.GetController(ENMControllerName)
				enmc.Run()

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "controller",
					Value: "rancher",
					Usage: "controller name",
				},
			},
		},
	}

	app.Run(os.Args)
}
