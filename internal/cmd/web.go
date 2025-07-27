package cmd

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/flamego/brotli"
	"github.com/flamego/flamego"
	"github.com/flamego/gzip"
	"github.com/flamego/template"

	"github.com/midoks/vez/internal/assets/public"
	"github.com/midoks/vez/internal/assets/templates"
	"github.com/midoks/vez/internal/conf"
	appcontext "github.com/midoks/vez/internal/context"
	"github.com/midoks/vez/internal/router"
	"github.com/midoks/vez/internal/tmpl"
	"github.com/midoks/vez/internal/tools"
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

	f := flamego.New()

	if !conf.Web.DisableRouterLog {
		f.Use(flamego.Logger())
	}

	f.Use(flamego.Recovery())

	if conf.Web.EnableGzip {
		f.Use(gzip.Gzip(gzip.Options{
			CompressionLevel: 9, // 最优压缩
		}))
	}

	f.Use(brotli.Brotli())

	// public
	f.Use(flamego.Static(flamego.StaticOptions{
		Directory:     filepath.Join(conf.CustomDir(), "public"),
		EnableLogging: !conf.Web.DisableRouterLog,
	}))

	var publicFs http.FileSystem
	if !conf.Web.LoadAssetsFromDisk {
		publicFs = public.NewFileSystem()
	}

	f.Use(flamego.Static(flamego.StaticOptions{
		Directory:     filepath.Join(conf.WorkDir(), "public"),
		FileSystem:    publicFs,
		EnableLogging: !conf.Web.DisableRouterLog,
	}))

	// template
	renderOpt := template.Options{
		Directory:         filepath.Join(conf.WorkDir(), "templates"),
		AppendDirectories: []string{filepath.Join(conf.CustomDir(), "templates")},
		FuncMaps:          tmpl.FuncMaps(),
	}

	if !conf.Web.LoadAssetsFromDisk {
		renderOpt.FileSystem = templates.NewTemplateFileSystem("", renderOpt.AppendDirectories[0])
	}

	f.Use(template.Templater(renderOpt))

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
		f.Get("/image/{id}", router.Image)
	}, appcontext.Contexter())

}

func runWebService(c *cli.Context) error {
	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动图片服务器健康检查
	go startImageServerHealthCheck(ctx)

	// 启动pprof服务器（仅在非生产环境）
	if conf.App.RunMode != "prod" {
		go startPprofServer()
	}

	f := newFlamego()
	setRouter(f)
	f.Run(conf.Web.HttpPort)

	return nil
}

// startImageServerHealthCheck 启动图片服务器健康检查
func startImageServerHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkImageServer()
		}
	}
}

// checkImageServer 检查图片服务器状态
func checkImageServer() {
	pingUrl := conf.Image.Ping
	if pingUrl == "" {
		return
	}

	r, err := tools.GetHttpData(pingUrl)
	if err != nil {
		conf.Image.PingStatus = false
		return
	}

	conf.Image.PingStatus = strings.EqualFold(r, conf.Image.PingResponse)
}

// startPprofServer 启动pprof服务器
func startPprofServer() {
	port := ":" + conf.Debug.Port
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Pprof server error: %v\n", err)
	}
}
