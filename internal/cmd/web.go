package cmd

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	"github.com/urfave/cli"

	"github.com/midoks/vez/internal/conf"
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
	f.Use(template.Templater())
	// f.Use(template.Templater(template.Options{FileSystem: fs}))

	return f
}

func setRouter(f *flamego.Flame) {
	f.Get("/", func(t template.Template, data template.Data) {
		t.HTML(http.StatusOK, "home")
	})

	f.Get("/csdn/{user}/{id}.html", func(t template.Template, data template.Data) {
		t.HTML(http.StatusOK, "page/content")
	})
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
