package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"net/http"
	_ "net/http/pprof"

	"gopkg.in/ini.v1"

	"github.com/midoks/vez/internal/conf"
	"github.com/midoks/vez/internal/logs"
	"github.com/midoks/vez/internal/mgdb"
	// "github.com/midoks/vez/internal/render"
	"github.com/midoks/vez/internal/tools"
)

func autoMakeCustomConf(customConf string) error {

	if tools.IsExist(customConf) {
		return nil
	}

	// auto make custom conf file
	cfg := ini.Empty()
	if tools.IsFile(customConf) {
		if err := cfg.Append(customConf); err != nil {
			return err
		}
	}

	cfg.Section("").Key("app_name").SetValue("vez")
	cfg.Section("").Key("run_mode").SetValue("prod")

	cfg.Section("web").Key("http_port").SetValue("11011")
	cfg.Section("session").Key("provider").SetValue("memory")

	os.MkdirAll(filepath.Dir(customConf), os.ModePerm)
	if err := cfg.SaveTo(customConf); err != nil {
		return err
	}

	return nil
}

func Init(customConf string) error {
	var err error

	if customConf == "" {
		customConf = filepath.Join(conf.CustomDir(), "conf", "app.conf")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return fmt.Errorf("custom conf file get absolute path: %s", err)
		}
	}

	err = autoMakeCustomConf(customConf)
	if err != nil {
		return err
	}

	conf.Init(customConf)
	logs.Init()
	mgdb.Init()

	// renderOpt := render.Options{
	// 	Directory:         filepath.Join(conf.WorkDir(), "templates"),
	// 	AppendDirectories: []string{filepath.Join(conf.CustomDir(), "templates")},
	// 	Funcs:             template.FuncMap(),
	// 	IndentJSON:        true,
	// }

	// if !conf.Server.LoadAssetsFromDisk {
	// 	renderOpt.TemplateFileSystem = templates.NewTemplateFileSystem("", renderOpt.AppendDirectories[0])
	// }

	// render.Renderer(renderOpt)

	if conf.App.RunMode != "prod" {
		go func() {
			port := fmt.Sprintf(":%s", conf.Debug.Port)
			http.ListenAndServe(port, nil)
		}()
	}

	return nil
}

func init() {
	Init("")
}
