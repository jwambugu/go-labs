package crawler_test

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"go-labs/internal/crawler"
	"go-labs/internal/testutils"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

func TestCrawler_Download(t *testing.T) {
	var (
		httpClient = testutils.NewTestHttpClient()
		crw        = crawler.NewCrawler(httpClient)
	)

	httpClient.Request("https://localhost.com", func() (code int, body string) {
		return http.StatusOK, `
			<!DOCTYPE html>
				<html>
					<head>
						<title>Page Title</title>
					</head>
				<body>
					<h1>This is a Heading</h1>
					<p>This is a paragraph.</p>
				</body>
			</html>
		`
	})

	contents, err := crw.Downloader("https://localhost.com")
	require.NoError(t, err)
	require.NotNil(t, contents)

	contents, err = crw.Downloader("https://localghost.com")
	require.Error(t, err)
	require.EqualError(t, err, crawler.ErrNotFound.Error())
	require.Nil(t, contents)
}

func TestCrawler_GetLinks(t *testing.T) {
	var (
		httpClient = testutils.NewTestHttpClient()
		crw        = crawler.NewCrawler(httpClient)
		rawURL     = "https://localhost.com"
	)

	uri, err := url.Parse(rawURL)
	require.NoError(t, err)

	contents := bytes.NewBuffer([]byte(`
			<ul>
				<a href="/">Home</a>
				<a href="/advanced-features">Advance features</a>
				<a href="/pricing">Pricing</a>
				<a href="https://google.com"> External </a>
				<a href="mailto:someone@example.com">Send email</a>
				<a href="#">Go Home</a>
			</ul>`))

	links := crw.GetLinks(uri, contents)
	require.Len(t, links, 2)
}

func TestCrawler_Fetch(t *testing.T) {
	var (
		httpClient = testutils.NewTestHttpClient()
		crw        = crawler.NewCrawler(httpClient)
		rawURL     = "https://localhost.com"
	)

	httpClient.Request(rawURL, func() (code int, body string) {
		return http.StatusOK, `
			<ul>
				<a href="/">Home</a>
				<a href="/advanced-features">Advance features</a>
				<a href="/pricing">Pricing</a>
				<a href="https://google.com"> External </a>
				<a href="mailto:someone@example.com">Send email</a>
				<a href="#">Go Home</a>
			</ul>`
	})

	body, urls, err := crw.Fetch(rawURL)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Len(t, urls, 2)
}

func TestCrawler_Crawl(t *testing.T) {
	var (
		httpClient = testutils.NewTestHttpClient()
		rawURL     = "https://localhost.com"
	)

	httpClient.Request(rawURL, func() (code int, body string) {
		return http.StatusOK, `
			<ul>
				<a href="/">Home</a>
				<a href="/advanced-features">Advance features</a>
				<a href="/pricing">Pricing</a>
				<a href="https://google.com"> External </a>
				<a href="mailto:someone@example.com">Send email</a>
				<a href="#">Go Home</a>
			</ul>`
	})

	httpClient.Request(rawURL+"/advanced-features", func() (code int, body string) {
		return http.StatusOK, `
			<p>Advanced Features</p>
			<ul>
				<a href="/">Home</a>
			</ul>`
	})

	httpClient.Request(rawURL+"/pricing", func() (code int, body string) {
		return http.StatusOK, `
			<p>Pricing</p>
			<ul>
				<a href="/test">Home</a>
			</ul>`
	})

	crw := crawler.NewCrawler(httpClient)

	var wg sync.WaitGroup
	wg.Add(1)

	crw.Crawl(rawURL, 4, crw, &wg)
	wg.Wait()
}
