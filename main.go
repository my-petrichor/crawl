package main

import (
	"os"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

const (
	filename = "院士.txt"
	url      = "http://www.casad.cas.cn/ysxx2017/ysmdyjj/qtysmd_124280/"
)

var (
	departments = make(map[string]int)
	names       = make([]string, 0)
	resumes     = make([]string, 0)
)

func main() {
	c := colly.NewCollector()
	c.AllowURLRevisit = true
	c.DetectCharset = true

	//The purpose of the reg is to get department_people_number
	reg := regexp.MustCompile("[1-9]+")

	c.OnHTML("dl#allNameBar", func(e *colly.HTMLElement) {
		e.ForEach("dt", func(_ int, element *colly.HTMLElement) {
			departmentName := element.Text
			department_people_number, _ := strconv.Atoi(reg.FindString(element.ChildText("em")))
			departments[departmentName] = department_people_number
		})

		e.ForEach("a[href]", func(_ int, ele *colly.HTMLElement) {
			name := ele.Text
			names = append(names, name)

			link := ele.Attr("href")

			c.Visit(e.Request.AbsoluteURL(link))
		})
	})

	c.OnHTML("div.contentBar", func(e *colly.HTMLElement) {
		if e.DOM.Find("p:contains(院士)").Text() != "" {
			resumes = append(resumes, e.DOM.Find("p:contains(院士)").Text())
		} else {
			resumes = append(resumes, e.DOM.Find("font:contains(院士)").Text())
		}
	})

	c.Visit(url)

	file, _ := os.Create(filename)
	defer file.Close()

	sum := 0
	max := 0
	for dep, number := range departments {
		file.WriteString(dep + "\n")

		max += number

		for sum < max {
			file.WriteString(names[sum] + ":" + " " + resumes[sum] + "\n")
			sum++
		}

		file.WriteString("\n")
	}
}
