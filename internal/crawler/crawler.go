package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var ErrNotFound = errors.New("page not found")

type Fetcher interface {
	// Fetch returns the body of URL and a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Crawler struct {
	httpClient HttpClient

	mu           sync.Mutex
	visitedPages map[string]struct{}
}

type fetchResult struct {
	body string
	err  error
	urls []string
}

func (c *Crawler) Downloader(uri string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %v", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	switch res.StatusCode {
	case http.StatusOK:
		contents, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %v", err)
		}

		return contents, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	}

	return nil, fmt.Errorf("request failed with status: %d", res.StatusCode)
}

func (c *Crawler) GetLinks(uri *url.URL, r io.Reader) (links []string) {
	var (
		foundLinks = make(map[string]struct{})
		tokenizer  = html.NewTokenizer(r)
	)

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			delete(foundLinks, uri.String()) // Already crawled, let's evict it.
			for link := range foundLinks {
				links = append(links, link)
			}

			break
		}

		token := tokenizer.Token()
		if tokenType != html.StartTagToken && token.DataAtom != atom.A {
			continue
		}

		for _, attr := range token.Attr {
			if attr.Key == "href" {
				link := attr.Val
				if link == "" || strings.Contains(link, "mailto") || strings.Contains(link, "#") {
					continue
				}

				linkURL, err := url.Parse(link)
				if err != nil {
					fmt.Printf("failed to parse %q: %v\n", link, err)
					continue
				}

				if linkURL.Host == uri.Host {
					link = strings.Trim(link, "/")
					foundLinks[link] = struct{}{}
					continue
				}

				if linkURL.Host == "" {
					link = uri.Scheme + "://" + uri.Host + link
					link = strings.Trim(link, "/")
					foundLinks[link] = struct{}{}
				}
			}
		}
	}
	return links
}

func (c *Crawler) Fetch(rawURL string) (body string, urls []string, err error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return "", nil, fmt.Errorf("parse url: %v", err)
	}

	webpage, err := c.Downloader(uri.String())
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return "", nil, err
		}
	}

	var (
		buffer = bytes.NewBuffer(webpage)
		links  = c.GetLinks(uri, buffer)
	)

	return string(webpage), links, nil
}

func (c *Crawler) Crawl(rawURL string, depth int, fetcher Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	c.mu.Lock()

	if _, visited := c.visitedPages[rawURL]; visited || depth <= 0 {
		c.mu.Unlock()
		return
	}

	c.visitedPages[rawURL] = struct{}{}
	c.mu.Unlock()

	fetchResultCh := make(chan *fetchResult, 1)

	go func() {
		body, urls, err := fetcher.Fetch(rawURL)

		fetchResultCh <- &fetchResult{
			body: body,
			err:  err,
			urls: urls,
		}
	}()

	select {
	case fetchedResult := <-fetchResultCh:
		if fetchedResult.err != nil {
			fmt.Println(fetchedResult.err)
			return
		}

		fmt.Printf("-- %s, %d depth %d\n", rawURL, len(fetchedResult.urls), depth)

		for _, u := range fetchedResult.urls {
			wg.Add(1)
			go c.Crawl(u, depth-1, fetcher, wg)
		}

	}
}

func NewCrawler(httpClient HttpClient) *Crawler {
	return &Crawler{
		httpClient:   httpClient,
		visitedPages: make(map[string]struct{}),
	}
}
