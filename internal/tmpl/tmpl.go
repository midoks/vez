package tmpl

import (
	"fmt"
	"html/template"
	// "mime"
	// "path/filepath"
	"encoding/base64"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
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
			"ParseTest": ParseTest,
			"ParseHtml": ParseHtml,
			"Sanitize":  bluemonday.UGCPolicy().Sanitize,
		}}
	})
	return funcMap
}

func ParseTest() string {
	return fmt.Sprintf("%s", "test")
}

func Safe(original string) template.HTML {
	return template.HTML(original)
}

func ParseHtml(original string) template.HTML {
	if conf.Image.PingStatus {

		prefix := conf.Image.Addr

		doc, _ := htmlquery.Parse(strings.NewReader(original))
		imgList := htmlquery.Find(doc, "//img")
		for _, img := range imgList {

			imagePath := htmlquery.SelectAttr(img, "src")
			if strings.EqualFold(imagePath, "") {
				continue
			}

			t := prefix + base64.StdEncoding.EncodeToString([]byte(imagePath))
			original = strings.Replace(original, imagePath, t, -1)
		}
	}

	return template.HTML(original)
}

func HeadTitle(original string) string {
	stripped := strip.StripTags(original)
	stripped = strings.TrimSpace(stripped)
	strippedRune := []rune(stripped)
	sublen := 100
	orilen := len(strippedRune)
	if orilen < sublen {
		return string(strippedRune)
	}
	return string(strippedRune[0:sublen])
}
