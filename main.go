package main

import (
	"sync"
	"time"

	"github.com/my-Sakura/crawl/data"
	"github.com/my-Sakura/crawl/news"
)

var (
	wg sync.WaitGroup
)

func main() {
	now := time.Now()
	date := now.Format("2006-01-02")
	wg.Add(3)

	go func() {
		news.StateDepartmentCrawlStart(date)
		wg.Done()
	}()
	go func() {
		data.EconomistFiftyCrawlStart()
		wg.Done()
	}()
	go func() {
		data.AcademicianCrawlStart()
		wg.Done()
	}()

	wg.Wait()
}
