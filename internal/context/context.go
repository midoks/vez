package context

import (
	// "fmt"
	"time"

	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"github.com/midoks/vez/internal/mgdb"
)

// var once sync.Once
// func init() {}
var mg []mgdb.ContentBid

func Contexter() flamego.Handler {
	return func(c flamego.Context, t template.Template, d template.Data) {

		if len(mg) == 0 {
			mg, _ = mgdb.ContentRandNum(10)
		}

		go func() {
			for {
				mg = []mgdb.ContentBid{}
				mg, _ = mgdb.ContentRandNum(10)
				time.Sleep(time.Second * 30)
			}
		}()

		d["Newsest"] = mg
	}
}
