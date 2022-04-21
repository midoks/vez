package tmpl

import (
	// "fmt"
	"html/template"
	// "mime"
	// "path/filepath"
	// "strings"
	"strconv"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/midoks/vez/internal/conf"
)

var (
	funcMap     []template.FuncMap
	funcMapOnce sync.Once
)

// FuncMap returns a list of user-defined template functions.
func FuncMaps() []template.FuncMap {
	funcMapOnce.Do(func() {
		funcMap = []template.FuncMap{map[string]interface{}{
			"BuildCommit": func() string {
				if conf.BuildCommit != "" {
					return conf.BuildCommit
				}
				return strconv.FormatInt(time.Now().Unix(), 10)
			},
			"Year": func() int {
				return time.Now().Year()
			},
			"Safe":     Safe,
			"Sanitize": bluemonday.UGCPolicy().Sanitize,
		}}
	})
	return funcMap
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}
