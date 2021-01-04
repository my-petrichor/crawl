package main

import (
	"github.com/my-Sakura/crawl/internal/pkg/data"
	"github.com/my-Sakura/crawl/internal/pkg/news"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()

	c.AddFunc("0 0 12 ? *", data.AcademicianCrawlStart)
	c.AddFunc("0 0 12 ? *", data.EconomistFiftyCrawlStart)
	c.AddFunc("0 0 8,10,12,14,16,18 ? *", news.StateDepartmentCrawlStart)

	c.Start()
	select {}
}
