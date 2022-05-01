package cmd

import (
	// "fmt"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/flamego/brotli"
	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"github.com/midoks/vez/internal/conf"
	"github.com/midoks/vez/internal/context"
	"github.com/midoks/vez/internal/router"
	"github.com/midoks/vez/internal/tmpl"
	"github.com/midoks/vez/internal/tools"

	"github.com/midoks/vez/internal/assets/public"
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

	f.Use(flamego.Static(flamego.StaticOptions{Directory: filepath.Join(conf.CustomDir(), "public")}))

	var publicFs http.FileSystem
	if !conf.Web.LoadAssetsFromDisk {
		publicFs = public.NewFileSystem()
	}

	f.Use(flamego.Static(flamego.StaticOptions{
		Directory:     filepath.Join(conf.WorkDir(), "public"),
		FileSystem:    publicFs,
		EnableLogging: false,
	}))

	f.Use(template.Templater(template.Options{
		FuncMaps: tmpl.FuncMaps(),
	}))

	// f.Use(template.Templater(template.Options{FileSystem: fs}))
	f.Use(brotli.Brotli())
	return f
}

func setRouter(f *flamego.Flame) {

	f.Group("", func() {
		f.Get("/", router.Home)
		f.Get("/prev/{pos}", router.Prev)
		f.Get("/next/{pos}", router.Next)
		f.Get("/rand", router.Rand)
		f.Get("/about", router.About)

		f.Get("/so/{kw}.html", router.So)
		f.Get("/so/{kw}/{prevNext}/{pos}.html", router.So)
		f.Get("/csdn/{user}/{id}.html", router.CsdnPageCotent)
		f.Get("/cnblogs/{user}/{id}.html", router.CnBlogsPageCotent)
	}, context.Contexter())

}

func runWebService(c *cli.Context) error {

	//check image server
	go func() {
		for {
			pingUrl := conf.Image.Ping

			if pingUrl != "" {
				r, err := tools.GetHttpData(pingUrl)

				if err != nil {
					conf.Image.PingStatus = false
				}

				// fmt.Println("..", r, conf.Image.PingResponse, conf.Image.PingStatus)
				if strings.EqualFold(r, conf.Image.PingResponse) {
					conf.Image.PingStatus = true
				}
			}
			time.Sleep(time.Second * 25)
		}
	}()

	// go tool pprof -http=:11112 --seconds 30 http://127.0.0.1:11012/debug/pprof/profile
	if conf.App.RunMode != "prod" {
		go func() {
			port := ":" + conf.Debug.Port
			http.ListenAndServe(port, nil)
		}()
	}

	f := newFlamego()
	setRouter(f)
	f.Run(conf.Web.HttpPort)

	return nil
}
