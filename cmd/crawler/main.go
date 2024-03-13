package main

import (
	"go-labs/internal/crawler"
	"log"
	"net/http"
	"sync"
)

func main() {
	crawl, err := crawler.NewCrawler("", http.DefaultClient)
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	crawl.Crawl("https://withkoa.com/", 4, crawl, &wg)
	wg.Wait()
}
