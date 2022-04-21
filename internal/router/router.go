package router

import (
	// "fmt"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"github.com/midoks/vez/internal/mgdb"
	"github.com/midoks/vez/internal/robot"
)

func Hello() string {
	return "Hello world"
}

func Home(t template.Template, data template.Data) {
	t.HTML(http.StatusOK, "home")
}

func Rand(c flamego.Context, t template.Template, data template.Data) {
	d, err := mgdb.ContentRand()

	if err != nil {
		t.HTML(http.StatusOK, "home")
		return
	}

	url := "/csdn/" + d.User + "/" + d.Id + ".html"
	c.Redirect(url)
}

func CsdnPageCotent(c flamego.Context, t template.Template, data template.Data) {

	// url := "https://blog.csdn.net/" + c.Param("user") + "/article/details/" + c.Param("id")
	// robot.SpiderCSDNUrl(url)

	d, _ := mgdb.ContentOriginFindOne(robot.CSND_NAME, c.Param("id"))
	data["Article"] = d
	t.HTML(http.StatusOK, "page/content")
}
