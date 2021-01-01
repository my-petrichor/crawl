package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/my-Sakura/crawl/config"
)

var (
	economistURL       = "http://www.50forum.org.cn/home/"
	economistFiftyData = make(chan []string)
	economistFinish    = make(chan interface{})
)

func economistCrawl() {
	var name string
	var job string
	var academy string
	flag := true

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
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

		data := []string{name, job, academy}
		economistFiftyData <- data

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

	time.Sleep(time.Second)
	close(economistFinish)
}

func economistCreate(wg *sync.WaitGroup) {
	file, err := os.Create(config.Conf.Economist)
	defer file.Close()
	if err != nil {
		fmt.Printf("fileCreateError: %v\n", err)
	}

	w := csv.NewWriter(file)
	w.Write([]string{"name", "job", "academy"})
	for {
		select {
		case data := <-economistFiftyData:
			w.Write(data)
		case <-economistFinish:
			w.Flush()
			err = w.Error()
			if err != nil {
				fmt.Println("flush error ", err)
			}

			wg.Done()
			return
		}
	}
}

func EconomistFiftyCrawlStart() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go economistCrawl()
	go economistCreate(&wg)

	wg.Wait()
}
