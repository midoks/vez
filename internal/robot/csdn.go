package robot

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"

	"github.com/midoks/vez/internal/lazyregexp"
	"github.com/midoks/vez/internal/mgdb"
	// "github.com/gocolly/colly/v2/debug"
)

func isMatchCSDN_Article(url string) bool {
	return lazyregexp.New(`https:\/\/blog.csdn.net/\w+/article/details/\d+`).Regexp().Match([]byte(url))
}

func getMatchCSDN_User_ID(url string) (string, string) {
	m := lazyregexp.New(`https:\/\/blog.csdn.net/(\w+)/article/details/(\d+)`).Regexp().FindStringSubmatch(url)
	if len(m) == 3 {
		return m[1], m[2]
	}
	return "", ""
}

func CreateCSDNCollector() *colly.Collector {
	csdn := colly.NewCollector(
		colly.Async(true),
		// Attach a debugger to the collector
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// csdn.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*httpbin.*",
	// 	Parallelism: 3,
	// 	Delay:       5 * time.Second,
	// })

	// Find and visit all links
	csdn.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		if isMatchCSDN_Article(url) {

			user, id := getMatchCSDN_User_ID(url)
			if user == "" {
				return
			}

			_, err := mgdb.ContentOriginFindOne("csdn", id)
			if err != nil {
				e.Request.Visit(url)
				return
			}
			// fmt.Println("repeat", url)
		}
	})

	// tmp.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	csdn.OnScraped(func(r *colly.Response) {

		url := r.Request.URL.String()
		if isMatchCSDN_Article(url) {

			fmt.Println("match", url)

			doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
			if err != nil {
				return
			}

			contentBody := htmlquery.Find(doc, `//div[@id="content_views"]/div`)

			if len(contentBody) == 0 {
				contentBody = htmlquery.Find(doc, `//div[@id="article_content"]/div`)
				if len(contentBody) == 0 {
					return
				}
			}

			user, id := getMatchCSDN_User_ID(url)
			if user == "" {
				return
			}

			contentTitle := htmlquery.Find(doc, `//h1[@id="articleContentId"]`)
			if len(contentTitle) == 0 {
				return
			}

			title := htmlquery.OutputHTML(contentTitle[0], false)

			html := htmlquery.OutputHTML(contentBody[0], false)
			mgdb.ContentAdd(mgdb.Content{
				Url:    url,
				Source: "csdn",
				User:   user,
				Id:     id,
				Title:  title,
				Html:   html,
			})
		}
	})
	return csdn
}

func RunCSDN() {
	csdn.CreateCSDNCollector()

	csdn.Visit("https://blog.csdn.net/mhs12345")
	csdn.Wait()
}
