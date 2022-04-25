package cmd

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/urfave/cli"

	"github.com/midoks/vez/internal/conf"
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

	// go tool pprof -http=:11113 --seconds 30 http://127.0.0.1:11013/debug/pprof/profile
	if conf.App.RunMode != "prod" {
		go func() {
			port := ":11013"
			http.ListenAndServe(port, nil)
		}()
	}

	robot.RunCSDN()
	return nil
}
