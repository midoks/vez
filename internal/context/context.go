package context

import (
	// "fmt"

	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"github.com/midoks/vez/internal/mgdb"
)

func Contexter() flamego.Handler {
	return func(c flamego.Context, t template.Template, d template.Data) {
		n, _ := mgdb.ContentRandNum(5)

		d["Newsest"] = n
	}
}
