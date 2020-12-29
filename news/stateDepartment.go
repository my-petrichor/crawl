package news

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/my-Sakura/crawl/config"
)

var (
	stateDepartmentNewsURL     = "http://www.gov.cn/xinwen/yaowen.htm"
	stateDepartmentPoliciesURL = "http://www.gov.cn/zhengce/index.htm"
	stateDepartmentPoliciesMap = make(map[string]string)
	stateDepartmentNewsMap     = make(map[string]string)
	policiesTitle              = make(chan string)
	policiesContent            = make(chan string)
	policiesFinish             = make(chan interface{})
	mainNewsTitle              = make(chan string)
	mainNewsContent            = make(chan string)
	mainNewsFinish             = make(chan interface{})
)

func stateDepartmentNewsCrawl(date string) {
	var title string
	var content string

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
	)

	c.OnHTML("div.news_box", func(e *colly.HTMLElement) {
		e.ForEach("h4", func(_ int, element *colly.HTMLElement) {
			if element.ChildText("span.date") == date {
				link := element.ChildAttr("a[href]", "href")
				c.Visit(element.Request.AbsoluteURL(link))
			}
		})
	})

	c.OnHTML("div.content", func(e *colly.HTMLElement) {
		title = e.DOM.Find("h1").Text()
		mainNewsTitle <- title

		e.ForEach("p", func(_ int, element *colly.HTMLElement) {
			content = content + element.Text + "\n"
		})

		//avoid blank title
		if title != "" {
			mainNewsContent <- content
		}

		//clear content
		content = ""
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	c.Visit(stateDepartmentNewsURL)

	close(mainNewsFinish)
}

func stateDepartmentPoliciesCrawl(date string) {
	var title string
	var content string
	c := colly.NewCollector()

	c.OnHTML("div.list_left_con", func(e *colly.HTMLElement) {
		e.ForEach("div.latestPolicy_left_item", func(_ int, element *colly.HTMLElement) {
			if element.DOM.Find("span").Eq(0).Text() == date {
				link := element.ChildAttr("a[href]", "href")
				c.Visit(e.Request.AbsoluteURL(link))
			}
		})
	})

	c.OnHTML("div.article.oneColumn.pub_border", func(e *colly.HTMLElement) {
		title = e.DOM.Find("h1").Eq(0).Text()
		policiesTitle <- title

		e.ForEach("div#UCAP-CONTENT.pages_content", func(_ int, element *colly.HTMLElement) {
			content = content + element.Text + "\n"
		})

		policiesContent <- content

		//clear content
		content = ""
	})

	c.OnHTML("td.b12c#UCAP-CONTENT", func(e *colly.HTMLElement) {
		e.ForEach("strong", func(_ int, element *colly.HTMLElement) {
			title = title + element.ChildText("span") + "\n"
		})
		policiesTitle <- title

		e.ForEach("p", func(_ int, element *colly.HTMLElement) {
			content = content + element.Text + "\n"
		})

		policiesContent <- content

		//clear content
		content = ""
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	c.Visit(stateDepartmentPoliciesURL)

	close(policiesFinish)
}

func stateDepartmentNewsCreate(wg *sync.WaitGroup) {
	count := 1
	file, err := os.Create(config.Conf.StateDepartmentNews)
	defer file.Close()
	if err != nil {
		fmt.Println("file create err", err)
	}

	for {
		select {
		case title := <-mainNewsTitle:
			file.WriteString(strconv.Itoa(count) + "." + title + "\n")
			count++

		case content := <-mainNewsContent:
			file.WriteString(content + "\n\n")
		case <-mainNewsFinish:
			wg.Done()
			return
		}
	}
}

func stateDepartmentPoliciesCreate(wg *sync.WaitGroup) {
	count := 1
	file, err := os.Create(config.Conf.StateDepartmentPolicies)
	defer file.Close()
	if err != nil {
		fmt.Println("stateDepartmentPoliciesCreate error: ", err)
	}

	for {
		select {
		case title := <-policiesTitle:
			file.WriteString(strconv.Itoa(count) + "." + title + "\n")
		case content := <-policiesContent:
			file.WriteString(content + "\n\n")
		case <-policiesFinish:
			wg.Done()
			return
		}
	}
}

//StateDepartmentCrawlStart generate .txt file after crawl stateDepartment Policies and mainNews
func StateDepartmentCrawlStart(date string) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go stateDepartmentNewsCrawl(date)
	go stateDepartmentPoliciesCrawl(date)
	go stateDepartmentNewsCreate(&wg)
	go stateDepartmentPoliciesCreate(&wg)

	wg.Wait()
}
