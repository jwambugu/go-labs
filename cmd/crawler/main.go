package main

import (
	"go-labs/internal/crawler"
	"net/http"
	"sync"
)

func main() {
	crawl := crawler.NewCrawler(http.DefaultClient)
	var wg sync.WaitGroup
	wg.Add(1)

	crawl.Crawl("https://withkoa.com/", 4, crawl, &wg)
	wg.Wait()
}
