package robot

import (
	"fmt"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"

	"github.com/midoks/vez/internal/lazyregexp"
	"github.com/midoks/vez/internal/mgdb"
)

const (
	CNBLOGS_NAME = "cnblogs"
)

func isMatchCnBlogs_Article(url string) bool {
	return lazyregexp.New(`https://www.cnblogs.com/(\w+)/p/([\w-]+).html`).Regexp().Match([]byte(url))
}

func isMatchCnBlogs_List(url string) bool {
	return lazyregexp.New(`https://www.cnblogs.com/(\w+)/`).Regexp().Match([]byte(url))
}

func isMatchCnBlogs_List_P(url string) bool {
	return lazyregexp.New(`https://www.cnblogs.com/(\w+)/default.html?page=(\d+)`).Regexp().Match([]byte(url))
}

func getMatchCnBlogs_User_ID(url string) (string, string) {
	m := lazyregexp.New(`https://www.cnblogs.com/(\w+)/p/([\w-]+).html`).Regexp().FindStringSubmatch(url)
	if len(m) == 3 {
		return m[1], m[2]
	}
	return "", ""
}

func CreateCnBlogsCollector() *colly.Collector {
	cnBlogs := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(4),
	)
	cnBlogs.Limit(&colly.LimitRule{
		DomainGlob:  "*www.cnblogs.com.*",
		Parallelism: 3,
		Delay:       5 * time.Second,
	})

	cnBlogs.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		// e.Request.Visit("https://www.cnblogs.com")

		if isMatchCnBlogs_List(url) {
			e.Request.Visit(url)
			return
		}

		if isMatchCnBlogs_List_P(url) {
			e.Request.Visit(url)
			return
		}

		if !isMatchCnBlogs_Article(url) {
			return
		}

		user, id := getMatchCnBlogs_User_ID(url)
		if user == "" {
			return
		}

		_, err := mgdb.ContentOriginFindOne(CNBLOGS_NAME, id)
		if err != nil {
			e.Request.Visit(url)
			return
		}
	})

	cnBlogs.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	// Set error handler
	cnBlogs.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "\nError:", err)
	})

	cnBlogs.OnScraped(func(r *colly.Response) {

		url := r.Request.URL.String()
		if isMatchCnBlogs_Article(url) {

			fmt.Println("match", url)
			doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
			if err != nil {
				return
			}

			user, id := getMatchCnBlogs_User_ID(url)
			if user == "" {
				return
			}

			contentBody := htmlquery.Find(doc, `//div[@id="cnblogs_post_body"]`)
			contentTitle := htmlquery.Find(doc, `//span[@role="heading"]`)
			if len(contentTitle) == 0 {
				return
			}

			title := htmlquery.OutputHTML(contentTitle[0], false)
			html := htmlquery.OutputHTML(contentBody[0], false)

			if strings.EqualFold(html, "") {
				fmt.Println("match", url, "data is empty")
				return
			}

			_, err = mgdb.ContentAdd(mgdb.Content{
				Url:    url,
				Source: CNBLOGS_NAME,
				User:   user,
				Id:     id,
				Title:  title,
				Html:   html,
			})
			if err != nil {
				fmt.Println("err:", err)
			}
		}
	})

	return cnBlogs
}
func RunCnBlogs() {

	app := CreateCnBlogsCollector()
	r, err := mgdb.ContentRandSource(CNBLOGS_NAME)
	if err == nil {
		fmt.Println("rand visiting: ", r.Url)
		app.Visit(r.Url)
		app.Visit("https://www.cnblogs.com")
	} else {
		fmt.Println("visiting start")
		app.Visit("https://www.cnblogs.com")
	}

	app.Wait()
}
