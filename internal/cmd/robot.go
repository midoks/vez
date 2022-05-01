package cmd

import (
	"time"

	"github.com/urfave/cli"

	"github.com/midoks/vez/internal/robot"
)

var Robot = cli.Command{
	Name:        "robot",
	Usage:       "This Robot services",
	Description: `Start robot services`,
	Action:      runRobotService,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func runRobotService(c *cli.Context) error {

	robot.RunCSDN()

	time.Sleep(time.Second * 60)

	robot.RunCnBlogs()

	return nil
}
