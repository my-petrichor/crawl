package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

var (
	economistURL       = "http://www.50forum.org.cn/home/"
	economist          = "经济50人.csv"
	economistFiftyData = make(chan []string)
	economistFinish    = make(chan interface{})
)

func economistCrawl() {
	var name string
	var job string
	var academy string
	flag := true

	c := colly.NewCollector()

	c.OnHTML("div.f_people", func(e *colly.HTMLElement) {
		//there are two div.f_people element
		//so exclude the one
		if e.ChildAttr("a[href]", "href") == "/home/article/lists/category/help_qiyejia.html" {
			return
		}

		e.ForEach("a[href]", func(_ int, element *colly.HTMLElement) {
			link := element.Attr("href")

			c.Visit(element.Request.AbsoluteURL(link))
		})
	})

	c.OnHTML("div.people_intro", func(e *colly.HTMLElement) {
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
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	c.Visit(economistURL)

	close(economistFinish)
}

func economistCreate(wg *sync.WaitGroup) {
	file, err := os.Create(economist)
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
