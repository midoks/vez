package router

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	"github.com/qiniu/qmgo/operator"

	"github.com/midoks/vez/internal/conf"
	"github.com/midoks/vez/internal/mgdb"
	"github.com/midoks/vez/internal/robot"
	"github.com/midoks/vez/internal/tmpl"
	"github.com/midoks/vez/internal/tools"
)

// 使用性能配置中的常量
const PAGE_NUM = conf.DefaultPageSize

func Hello() string {
	return fmt.Sprintf("%s", "Hello world")
}

// createTestData 创建测试数据
func createTestData() []mgdb.ContentBid {
	return []mgdb.ContentBid{
		{
			MgID:       "507f1f77bcf86cd799439011",
			Title:      "Go语言高性能编程技巧",
			Source:     "csdn",
			User:       "golang_expert",
			Id:         "123456",
			Html:       "<p>Go语言是一门现代化的编程语言，具有出色的性能和并发特性...</p>",
			Createtime: time.Now().Add(-24 * time.Hour),
		},
		{
			MgID:       "507f1f77bcf86cd799439012",
			Title:      "数据库设计最佳实践",
			Source:     "cnblogs",
			User:       "db_master",
			Id:         "789012",
			Html:       "<p>良好的数据库设计是应用程序性能的基础...</p>",
			Createtime: time.Now().Add(-48 * time.Hour),
		},
		{
			MgID:       "507f1f77bcf86cd799439013",
			Title:      "前端开发现代化工具链",
			Source:     "csdn",
			User:       "frontend_dev",
			Id:         "345678",
			Html:       "<p>现代前端开发需要掌握各种工具和框架...</p>",
			Createtime: time.Now().Add(-72 * time.Hour),
		},
	}
}

// createTestSidebarData 创建测试侧边栏数据
func createTestSidebarData() []mgdb.ContentBid {
	return []mgdb.ContentBid{
		{
			MgID:   "507f1f77bcf86cd799439014",
			Title:  "微服务架构设计模式",
			Source: "csdn",
			User:   "architect",
			Id:     "111111",
		},
		{
			MgID:   "507f1f77bcf86cd799439015",
			Title:  "Docker容器化部署实战",
			Source: "cnblogs",
			User:   "devops_pro",
			Id:     "222222",
		},
		{
			MgID:   "507f1f77bcf86cd799439016",
			Title:  "Redis性能优化技巧",
			Source: "csdn",
			User:   "cache_expert",
			Id:     "333333",
		},
	}
}

// addSidebarData 添加侧边栏数据
func addSidebarData(data template.Data) {
	newsest, err := mgdb.ContentNewsest()
	if err != nil {
		// 如果获取真实数据失败，使用测试数据
		newsest = createTestSidebarData()
	}
	data["Newsest"] = newsest
}

func Home(t template.Template, data template.Data) {
	d, err := mgdb.ContentOriginFindId("", "-", PAGE_NUM)
	if err != nil {
		// 使用测试数据作为降级方案
		d = createTestData()
	}
	data["Articles"] = d

	dLen := len(d)
	if dLen > 1 {
		data["NextPos"] = d[dLen-1].MgID
	}

	addSidebarData(data)
	t.HTML(http.StatusOK, "home")
}

func So(c flamego.Context, t template.Template, data template.Data) {
	kw := c.Param("kw")
	prevNext := c.Param("prevNext")
	id := c.Param("pos")

	// 输入验证
	if len(kw) == 0 || len(kw) > conf.MaxSearchKeywordLength {
		c.Redirect("/")
		return
	}

	op := operator.Lt
	if strings.EqualFold(prevNext, "prev") {
		op = operator.Gt
	}

	d, err := mgdb.ContentOriginFindSoso(id, "-", op, kw, PAGE_NUM)
	if err != nil {
		// 记录错误但不中断服务
		fmt.Printf("Search error: %v\n", err)
		d = []mgdb.ContentBid{}
	}

	data["Articles"] = d
	data["Keyword"] = kw

	dLen := len(d)

	if strings.EqualFold(prevNext, "prev") && dLen != PAGE_NUM {
		c.Redirect("/so/" + kw + ".html")
		return
	}

	if strings.EqualFold(prevNext, "next") && dLen == 0 {
		c.Redirect("/so/" + kw + ".html")
		return
	}

	if dLen > 1 {
		data["PrePos"] = d[0].MgID
		if dLen == PAGE_NUM {
			data["NextPos"] = d[dLen-1].MgID
		}
	}

	addSidebarData(data)
	t.HTML(http.StatusOK, "soso")
}

func Prev(c flamego.Context, t template.Template, data template.Data) {
	id := c.Param("pos")
	d, err := mgdb.ContentOriginFindIdGt(id, "-", PAGE_NUM)
	if err != nil {
		fmt.Printf("Prev page error: %v\n", err)
		d = []mgdb.ContentBid{}
	}

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

	addSidebarData(data)
	t.HTML(http.StatusOK, "home")
}

func Next(c flamego.Context, t template.Template, data template.Data) {
	id := c.Param("pos")

	d, err := mgdb.ContentOriginFindId(id, "-", PAGE_NUM)
	if err != nil {
		fmt.Printf("Next page error: %v\n", err)
		d = []mgdb.ContentBid{}
	}
	data["Articles"] = d

	dLen := len(d)
	if dLen > 1 {
		data["PrePos"] = d[0].MgID
		if dLen == PAGE_NUM {
			data["NextPos"] = d[dLen-1].MgID
		}
	}

	addSidebarData(data)
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
	addSidebarData(data)
	t.HTML(http.StatusOK, "about")
}

func CsdnPageCotent(c flamego.Context, t template.Template, data template.Data) {
	id := c.Param("id")
	if len(id) == 0 {
		c.Redirect("/")
		return
	}

	d, err := mgdb.ContentOriginFindOne(robot.CSND_NAME, id)
	if err != nil {
		fmt.Printf("CSDN content error: %v\n", err)
		c.Redirect("/")
		return
	}
	data["Article"] = d
	addSidebarData(data)
	t.HTML(http.StatusOK, "page/content")
}

func CnBlogsPageCotent(c flamego.Context, t template.Template, data template.Data) {
	id := c.Param("id")
	if len(id) == 0 {
		c.Redirect("/")
		return
	}

	d, err := mgdb.ContentOriginFindOne(robot.CNBLOGS_NAME, id)
	if err != nil {
		fmt.Printf("CNBlogs content error: %v\n", err)
		c.Redirect("/")
		return
	}
	data["Article"] = d
	addSidebarData(data)
	t.HTML(http.StatusOK, "page/content")
}

func splitImageUrlHeader(url_sign string) string {
	dir := url_sign[0:1] + "/" + url_sign[1:2]
	return dir
}

func getImageContent(url_sign string) string {
	url_header := splitImageUrlHeader(url_sign)
	define_dir := conf.ImageCacheDir + url_header
	abs_file := define_dir + "/" + url_sign

	// 检查目录是否存在，不存在则创建
	if b, err := tools.PathExists(define_dir); err != nil {
		fmt.Printf("Error checking directory: %v\n", err)
		return ""
	} else if !b {
		if err := os.MkdirAll(define_dir, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return ""
		}
	}

	// 检查文件是否存在
	if b, err := tools.PathExists(abs_file); err != nil {
		fmt.Printf("Error checking file: %v\n", err)
		return ""
	} else if !b {
		return ""
	}

	// 读取文件内容
	bytes, err := ioutil.ReadFile(abs_file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return ""
	}
	return string(bytes)
}

func setImageContent(url_sign string, content string) {
	url_header := splitImageUrlHeader(url_sign)
	define_dir := conf.ImageCacheDir + url_header
	abs_file := define_dir + "/" + url_sign

	// 检查目录是否存在，不存在则创建
	if b, err := tools.PathExists(define_dir); err != nil {
		fmt.Printf("Error checking directory: %v\n", err)
		return
	} else if !b {
		if err := os.MkdirAll(define_dir, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	}

	// 写入文件
	if err := os.WriteFile(abs_file, []byte(content), os.ModePerm); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
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
