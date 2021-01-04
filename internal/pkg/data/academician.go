package data

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/my-Sakura/crawl/config"
)

var (
	academicianURL        = "http://www.casad.cas.cn/ysxx2017/ysmdyjj/qtysmd_124280/"
	academicianDepartment = make(chan string)
	academicianName       = make(chan string)
	academicianContent    = make(chan string)
	academicianFinish     = make(chan interface{})
)

func academicianCrawl() {
	reg := regexp.MustCompile("[1-9]+")

	var count int
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
	)

	c1 := c.Clone()

	c.OnHTML("dl#allNameBar", func(e *colly.HTMLElement) {

		e.ForEach("dt", func(_ int, element *colly.HTMLElement) {
			department := element.Text
			academicianDepartment <- department

			department_people_number, _ := strconv.Atoi(reg.FindString(element.ChildText("em")))

			for i := count; i < department_people_number+count; i++ {
				link, _ := e.DOM.Find("a[href]").Eq(i).Attr("href")
				fmt.Println(link)
				c1.Visit(element.Request.AbsoluteURL(link))
			}
			count += department_people_number
		})
	})

	c1.OnHTML("div.contentBar", func(e *colly.HTMLElement) {
		name := e.DOM.Find("h1").Eq(0).Text()
		academicianName <- name

		content := e.DOM.Find("p:contains(院士)").Text()
		academicianContent <- content
	})

	c.OnRequest(func(r *colly.Request) {
		r.ProxyURL = "http://192.168.0.102:7890"
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	c1.OnRequest(func(r *colly.Request) {
		r.ProxyURL = "http://192.168.0.102:7890"
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	c.Visit(academicianURL)

	time.Sleep(time.Second)
	close(academicianFinish)
}

func academicianCreate(wg *sync.WaitGroup) {
	file, _ := os.Create(config.Conf.Academician)
	defer file.Close()

	for {
		select {
		case department := <-academicianDepartment:
			file.WriteString("\n" + department + "\n")
		case name := <-academicianName:
			file.WriteString(name + ": ")
		case content := <-academicianContent:
			file.WriteString(content + "\n")
		case <-academicianFinish:
			wg.Done()
			return
		}
	}
}

func AcademicianCrawlStart() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go academicianCrawl()
	go academicianCreate(&wg)

	wg.Wait()
}
