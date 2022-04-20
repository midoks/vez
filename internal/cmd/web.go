package cmd

import (
	"github.com/urfave/cli"
)

var Service = cli.Command{
	Name:        "web",
	Usage:       "This command starts all services",
	Description: `Start Web services`,
	Action:      runWebService,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func runWebService(c *cli.Context) error {

	return nil
}
