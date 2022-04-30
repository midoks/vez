package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/midoks/vez/internal/cmd"
	"github.com/midoks/vez/internal/conf"
)

const Version = "0.0.1"
const AppName = "vez"

func init() {
	conf.App.Version = Version
	conf.App.Name = AppName
}

func main() {

	app := cli.NewApp()
	app.Name = AppName
	app.Version = Version
	app.Usage = "A Vez service"
	app.Commands = []cli.Command{
		cmd.Service,
		cmd.Robot,
	}

	if err := app.Run(os.Args); err != nil {
		log.Println("Failed to start application: ", err)
	}
}
