package router

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	"github.com/qiniu/qmgo/operator"

	"github.com/midoks/vez/internal/mgdb"
	"github.com/midoks/vez/internal/robot"
	"github.com/midoks/vez/internal/tmpl"
	"github.com/midoks/vez/internal/tools"
)

const PAGE_NUM = 15

func Hello() string {
	return fmt.Sprintf("%s", "Hello world")
}

func Home(t template.Template, data template.Data) {

	d, _ := mgdb.ContentOriginFindId("", "-", PAGE_NUM)
	data["Articles"] = d

	dLen := len(d)
	if dLen > 1 {
		data["NextPos"] = d[dLen-1].MgID
	}

	t.HTML(http.StatusOK, "home")
}

func So(c flamego.Context, t template.Template, data template.Data) {
	kw := c.Param("kw")
	prevNext := c.Param("prevNext")
	id := c.Param("pos")

	op := operator.Lt
	if strings.EqualFold(prevNext, "prev") {
		op = operator.Gt
	}

	d, _ := mgdb.ContentOriginFindSoso(id, "-", op, kw, PAGE_NUM)

	data["Articles"] = d
	data["Keyword"] = kw

	dLen := len(d)

	if strings.EqualFold(prevNext, "prev") && dLen != PAGE_NUM {
		c.Redirect("/so/" + kw + ".html")
	}

	if strings.EqualFold(prevNext, "next") && dLen == 0 {
		c.Redirect("/so/" + kw + ".html")
	}

	if dLen > 1 {
		data["PrePos"] = d[0].MgID
		if dLen == PAGE_NUM {
			data["NextPos"] = d[dLen-1].MgID
		}
	}

	t.HTML(http.StatusOK, "soso")
}

func Prev(c flamego.Context, t template.Template, data template.Data) {
	id := c.Param("pos")
	d, _ := mgdb.ContentOriginFindIdGt(id, "-", PAGE_NUM)

	data["Articles"] = d

	dLen := len(d)
	if dLen != PAGE_NUM {
		c.Redirect("/")
		return
	}

	if dLen > 1 {
		data["PrePos"] = d[0].MgID
		data["NextPos"] = d[dLen-1].MgID
	}

	t.HTML(http.StatusOK, "home")
}

func Next(c flamego.Context, t template.Template, data template.Data) {
	id := c.Param("pos")

	d, _ := mgdb.ContentOriginFindId(id, "-", PAGE_NUM)
	data["Articles"] = d

	dLen := len(d)
	if dLen > 1 {
		data["PrePos"] = d[0].MgID
		if dLen == PAGE_NUM {
			data["NextPos"] = d[dLen-1].MgID
		}
	}

	t.HTML(http.StatusOK, "home")
}

func Rand(c flamego.Context, t template.Template, data template.Data) {
	d, err := mgdb.ContentRand()

	if err != nil {
		t.HTML(http.StatusOK, "home")
		return
	}

	url := "/" + d.Source + "/" + d.User + "/" + d.Id + ".html"
	c.Redirect(url)
}

func About(c flamego.Context, t template.Template, data template.Data) {
	t.HTML(http.StatusOK, "about")
}

func CsdnPageCotent(c flamego.Context, t template.Template, data template.Data) {
	d, _ := mgdb.ContentOriginFindOne(robot.CSND_NAME, c.Param("id"))
	data["Article"] = d
	t.HTML(http.StatusOK, "page/content")
}

func CnBlogsPageCotent(c flamego.Context, t template.Template, data template.Data) {
	d, _ := mgdb.ContentOriginFindOne(robot.CNBLOGS_NAME, c.Param("id"))
	data["Article"] = d
	t.HTML(http.StatusOK, "page/content")
}

func splitImageUrlHeader(url_sign string) string {
	dir := url_sign[0:1] + "/" + url_sign[1:2]
	return dir
}

func getImageContent(url_sign string) string {
	url_header := splitImageUrlHeader(url_sign)
	define_dir := "upload/image/" + url_header

	abs_file := define_dir + "/" + url_sign

	// fmt.Println(url_header, define_dir

	b, _ := tools.PathExists(define_dir)
	if !b {
		os.MkdirAll(define_dir, os.ModePerm)
	}

	b, _ = tools.PathExists(abs_file)
	if !b {
		return ""
	} else {
		bytes, _ := ioutil.ReadFile(abs_file)
		return string(bytes)
	}

	return ""
}

func setImageContent(url_sign string, content string) {
	url_header := splitImageUrlHeader(url_sign)
	define_dir := "upload/image/" + url_header

	abs_file := define_dir + "/" + url_sign

	b, _ := tools.PathExists(define_dir)
	if !b {
		os.MkdirAll(define_dir, os.ModePerm)
	}
	os.WriteFile(abs_file, []byte(content), os.ModePerm)
}

func Image(c flamego.Context, t template.Template, data template.Data) string {

	url := c.Param("id")
	decoded, err := base64.StdEncoding.DecodeString(url)

	if err == nil {
		// fmt.Println(tmpl.BytesToString(decoded), err)
		url = tmpl.BytesToString(decoded)

		url_sign := tools.Md5(url)
		// fmt.Println(url_sign)

		content := getImageContent(url_sign)
		if !strings.EqualFold(content, "") {
			return content
		}

		content, err := tools.GetHttpData(url)
		if err == nil {
			setImageContent(url_sign, content)
			return content
		}

		// fmt.Println(url)
	}
	return url
}
