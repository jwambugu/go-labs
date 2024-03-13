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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// ErrNotFound is returned when a page is not found.
var ErrNotFound = errors.New("page not found")

// alphanumericRegex is a regular expression to match non-alphanumeric characters.
var alphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// DestinationDir is the directory where fetched pages will be saved if none is provided.
const DestinationDir = "storage"

// Fetcher defines an interface for fetching web pages and extracting URLs.
//
// Fetch returns the body of URL and a slice of URLs found on that page.
type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

// HttpClient defines an interface for performing HTTP requests.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Crawler struct {
	destinationDir string
	httpClient     HttpClient

	mu           sync.Mutex
	visitedPages map[string]struct{}
}

type fetchResult struct {
	body string
	err  error
	urls []string
}

// Downloader downloads the content of a web page specified by the given URL.
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

func (c *Crawler) Save(filename string, contents []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}

	if _, err = file.Write(contents); err != nil {
		return fmt.Errorf("write: %v", err)
	}

	return nil
}

// GetLinks extracts the URLs from the HTML content provided.
func (c *Crawler) GetLinks(uri *url.URL, r io.Reader) (links []string) {
	var (
		foundLinks = make(map[string]struct{})
		tokenizer  = html.NewTokenizer(r)
	)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			delete(foundLinks, uri.String())

			for link := range foundLinks {
				links = append(links, link)
			}

			break
		}

		token := tokenizer.Token()

		if tokenType != html.StartTagToken || token.DataAtom != atom.A {
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

// Fetch fetches the content of the web page specified by the given URL.
func (c *Crawler) Fetch(rawURL string) (body string, urls []string, err error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return "", nil, fmt.Errorf("parse url: %v", err)
	}

	filename := rawURL
	filename = alphanumericRegex.ReplaceAllString(filename, "-") + ".html"
	filename = filepath.Join(c.destinationDir, filename)

	contents, err := os.ReadFile(filename)
	if err != nil && !errors.Is(err, io.EOF) {
		if !os.IsNotExist(err) {
			return "", nil, fmt.Errorf("exists: %v", err)
		}

		contents, err = c.Downloader(uri.String())
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return "", nil, err
			}
		}

		if err = c.Save(filename, contents); err != nil {
			return "", nil, err
		}
	}

	var (
		buffer = bytes.NewBuffer(contents)
		links  = c.GetLinks(uri, buffer)
	)

	return string(contents), links, nil
}

// Crawl crawls the web starting from the specified URL up to the given depth.
func (c *Crawler) Crawl(rawURL string, depth int, fetcher Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, visited := c.visitedPages[rawURL]; visited || depth <= 0 {
		return
	}

	c.visitedPages[rawURL] = struct{}{}

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

		fmt.Printf("-- %s, links %d\n", rawURL, len(fetchedResult.urls))

		for _, u := range fetchedResult.urls {
			wg.Add(1)
			go c.Crawl(u, depth-1, fetcher, wg)
		}
	}
}

// NewCrawler creates a new instance of Crawler with the specified destination directory and HTTP client.
// If the destination directory is empty, the default directory DestinationDir will be used.
func NewCrawler(destinationDir string, httpClient HttpClient) (*Crawler, error) {
	if destinationDir == "" {
		destinationDir = DestinationDir
	}

	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("mkdir: %v", err)
	}

	return &Crawler{
		destinationDir: destinationDir,
		httpClient:     httpClient,
		visitedPages:   make(map[string]struct{}),
	}, nil
}
