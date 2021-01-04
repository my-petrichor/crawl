package data

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/my-Sakura/crawl/pkg/model/mysql"
)

var (
	economistURL       = "http://www.50forum.org.cn/home/"
	economistFiftyData = make(chan []string)
	economistFinish    = make(chan interface{})
)

func economistCrawl(wg *sync.WaitGroup) {
	var name string
	var job string
	var academy string
	flag := true

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
		colly.Async(true),
	)

	c1 := c.Clone()

	c.OnHTML("div.f_people", func(e *colly.HTMLElement) {
		//there are two div.f_people element
		//so exclude the one
		if e.ChildAttr("a[href]", "href") == "/home/article/lists/category/help_qiyejia.html" {
			return
		}

		e.ForEach("a[href]", func(_ int, element *colly.HTMLElement) {
			link := element.Attr("href")

			c1.Visit(element.Request.AbsoluteURL(link))
		})
	})

	c1.OnHTML("div.people_intro", func(e *colly.HTMLElement) {
		//get Index
		name = e.DOM.Find("p").Eq(0).Text()
		var basePoint int
		e.ForEach("p", func(_ int, element *colly.HTMLElement) {
			if flag {
				if element.Text == "" {
					basePoint = element.Index
					flag = false
				}
			}
		})

		for i := 1; i < basePoint; i++ {
			job = job + e.DOM.Find("p").Eq(i).Text() + "、"
		}

		academyIndex := basePoint + 1
		academy = e.DOM.Find("p").Eq(academyIndex).Text()

		economistFiftyData <- []string{name, job, academy}

		job = ""
		flag = true
	})

	c.OnRequest(func(r *colly.Request) {
		r.ProxyURL = "http://192.168.0.102:7890"
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	c1.OnRequest(func(r *colly.Request) {
		r.ProxyURL = "http://192.168.0.102:7890"
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	c.Visit(economistURL)

	economistFiftyData <- nil
	<-economistFinish
	wg.Done()
}

func economistCreate() {
	for data := range economistFiftyData {
		name := mysql.EconomistFiftyGet(mysql.C.DB, data)
		if name == "" {
			mysql.EconomistFiftyInsert(mysql.C.DB, data)
		}

		if data == nil {
			break
		}
	}

}

func EconomistFiftyCrawlStart() {
	mysql.EconomistFiftyTableCreate(mysql.C.DB)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go economistCrawl(&wg)
	go economistCreate()

	wg.Wait()
}
