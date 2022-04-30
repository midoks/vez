package robot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unsafe"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"

	"github.com/midoks/vez/internal/lazyregexp"
	"github.com/midoks/vez/internal/mgdb"
	// "github.com/gocolly/colly/v2/debug"
)

const (
	CSND_NAME    = "csdn"
	LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
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

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = LETTER_BYTES[rand.Intn(len(LETTER_BYTES))]
	}
	return BytesToString(b)
}

func CreateCSDNCollector() *colly.Collector {
	csdn := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
		// Attach a debugger to the collector
		// colly.Debugger(&debug.LogDebugger{}),
	)

	csdn.Limit(&colly.LimitRule{
		DomainGlob:  "*blog.csdn.net.*",
		Parallelism: 3,
		Delay:       5 * time.Second,
	})

	// Find and visit all links
	csdn.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		if isMatchCSDN_Article(url) {

			user, id := getMatchCSDN_User_ID(url)
			if user == "" {
				return
			}

			_, err := mgdb.ContentOriginFindOne(CSND_NAME, id)
			if err != nil {
				e.Request.Visit(url)
				return
			}
			// fmt.Println("repeat", url)
		}
	})

	csdn.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	// Set error handler
	csdn.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	csdn.OnScraped(func(r *colly.Response) {

		url := r.Request.URL.String()
		if isMatchCSDN_Article(url) {

			fmt.Println("match", url)

			doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
			if err != nil {
				return
			}

			contentBody := htmlquery.Find(doc, `//div[@id="article_content"]/div`)
			if len(contentBody) == 0 {
				contentBody = htmlquery.Find(doc, `//div[@id="content_views"]/div`)
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

			_, err = mgdb.ContentAdd(mgdb.Content{
				Url:    url,
				Source: CSND_NAME,
				User:   user,
				Id:     id,
				Title:  title,
				Html:   html,
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	})
	return csdn
}

func SpiderCSDNUrl(url string) {
	go func() {
		csdn := CreateCSDNCollector()
		csdn.Visit(url)
		csdn.Wait()
	}()
}

func RunCSDN() {

	// dStart, _ := mgdb.ContentOriginFindNewsestOne(CSND_NAME)

	csdn := CreateCSDNCollector()

	r, err := mgdb.ContentOneByOne(CSND_NAME)
	// r, err := mgdb.ContentRand()
	if err == nil {
		fmt.Println("rand visiting: ", r.Url)
		// csdn.Visit("https://blog.csdn.net/suguoliang/article/details/123956701")
		csdn.Visit(r.Url)
	} else {
		fmt.Println("visiting start")
		csdn.Visit("https://blog.csdn.net")
	}

	csdn.Wait()

	// dEnd, _ := mgdb.ContentOriginFindNewsestOne(CSND_NAME)

	// if dEnd.MgID == dStart.MgID {
	// 	fmt.Println("visiting restart")
	// 	csdn.Visit("https://blog.csdn.net")
	// 	csdn.Wait()
	// }

}
