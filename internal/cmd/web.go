package cmd

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/urfave/cli"

	"github.com/flamego/brotli"
	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"github.com/midoks/vez/internal/conf"
	"github.com/midoks/vez/internal/router"
	"github.com/midoks/vez/internal/tmpl"
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

func newFlamego() *flamego.Flame {

	f := flamego.Classic()

	f.Use(flamego.Static(flamego.StaticOptions{Directory: "public"}))

	// fs, err := template.EmbedFS(Templates, "templates", []string{".tmpl"})
	// if err != nil {
	// 	panic(err)
	// }

	f.Use(template.Templater(template.Options{
		FuncMaps: tmpl.FuncMaps(),
	}))
	// f.Use(template.Templater(template.Options{FileSystem: fs}))
	f.Use(brotli.Brotli())
	return f
}

func setRouter(f *flamego.Flame) {
	f.Get("/", router.Home)

	f.Get("/rand", router.Rand)
	f.Get("/csdn/{user}/{id}.html", router.CsdnPageCotent)
}

func runWebService(c *cli.Context) error {

	if conf.App.RunMode != "prod" {
		go func() {
			port := fmt.Sprintf(":%s", conf.Debug.Port)
			http.ListenAndServe(port, nil)
		}()
	}

	f := newFlamego()
	setRouter(f)
	f.Run()

	return nil
}
