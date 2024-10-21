package tmpl

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

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

// Byte to string, only read-only
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String to byte, only read-only
func StringToBytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

func ParseTest() string {
	return fmt.Sprintf("%s", "test")
}

func Safe(original string) template.HTML {
	return template.HTML(original)
}

func ParseHtml(original string) template.HTML {
	doc, _ := htmlquery.Parse(strings.NewReader(original))
	imgList := htmlquery.Find(doc, "//img")

	for _, img := range imgList {

		imagePath := htmlquery.SelectAttr(img, "src")
		if strings.EqualFold(imagePath, "") {
			continue
		}

		image_status := conf.Image.Status

		imagePath = strings.Trim(imagePath, " ")

		suffix := url.QueryEscape(base64.StdEncoding.EncodeToString(StringToBytes(imagePath)))
		t := ""
		if image_status {
			t = conf.Image.Addr + "/" + suffix
		} else {
			t = "/image/" + suffix
		}
		original = strings.Replace(original, imagePath, t, -1)
		// fmt.Println("image:", imagePath, t, suffix)
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
