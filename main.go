package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

const (
	filename = "院士.txt"
	url      = "http://www.casad.cas.cn/ysxx2017/ysmdyjj/qtysmd_124280/"
)

var (
	wg = sync.WaitGroup{}
)

type people struct {
	departments              []string
	department_people_number []int
	names                    []string
	resumes                  []string
}

func main() {
	c := colly.NewCollector()
	c.AllowURLRevisit = true
	c.DetectCharset = true

	var p people

	wg.Add(1)

	//The purpose of the reg is to get department_people_number
	reg := regexp.MustCompile("[1-9]+")

	u := time.Now()
	c.OnHTML("dl#allNameBar", func(e *colly.HTMLElement) {
		go func() {
			e.ForEach("dt", func(_ int, element *colly.HTMLElement) {
				p.departments = append(p.departments, element.Text)
				department_people_number, _ := strconv.Atoi(reg.FindString(element.ChildText("em")))
				p.department_people_number = append(p.department_people_number, department_people_number)
			})
			wg.Done()
		}()

		e.ForEach("a[href]", func(_ int, ele *colly.HTMLElement) {
			name := ele.Text
			p.names = append(p.names, name)

			link := ele.Attr("href")

			fmt.Println("visiting =>", link)

			c.Visit(e.Request.AbsoluteURL(link))
		})
	})

	c.OnHTML("div.contentBar", func(e *colly.HTMLElement) {
		if e.DOM.Find("p:contains(院士)").Text() != "" {
			p.resumes = append(p.resumes, e.DOM.Find("p:contains(院士)").Text())
		} else {
			p.resumes = append(p.resumes, e.DOM.Find("font:contains(院士)").Text())
		}
	})

	c.Visit(url)
	wg.Wait()
	t := time.Now()
	crawlDuration := t.Sub(u)
	fmt.Println(crawlDuration)

	file, _ := os.Create(filename)
	defer file.Close()

	sum := 0
	max := 0
	for i, dep := range p.departments {
		file.WriteString(dep + "\n")

		max += p.department_people_number[i]

		for sum < max {
			file.WriteString(p.names[sum] + ":" + " " + p.resumes[sum] + "\n")
			sum++
		}

		file.WriteString("\n")
	}
}
