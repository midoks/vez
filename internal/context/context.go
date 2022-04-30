package context

import (
	// "fmt"
	"time"

	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"github.com/midoks/vez/internal/mgdb"
)

var mg []mgdb.ContentBid

func InitMG() {

	go func() {
		for {
			mg, _ = mgdb.ContentRandNum(10)
			time.Sleep(time.Second * 300)
		}
	}()

}

func Contexter() flamego.Handler {
	return func(c flamego.Context, t template.Template, d template.Data) {

		if len(mg) == 0 {
			mg, _ = mgdb.ContentRandNum(10)
		}

		d["Newsest"] = mg
	}
}
