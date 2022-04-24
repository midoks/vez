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

	strip "github.com/grokify/html-strip-tags-go"

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
			"HeadTitle": HeadTitle,
			"Safe":      Safe,
			"Sanitize":  bluemonday.UGCPolicy().Sanitize,
		}}
	})
	return funcMap
}

func Safe(original string) template.HTML {
	return template.HTML(original)
}

func HeadTitle(original string) string {
	stripped := strip.StripTags(original)
	strippedRune := []rune(stripped)
	sublen := 100
	orilen := len(strippedRune)
	if orilen < sublen {
		return string(strippedRune)
	}
	return string(strippedRune[0:sublen])
}
