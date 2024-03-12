package testutils

import (
	"bytes"
	"io"
	"net/http"
)

// testResponseFunc is a function type representing a test HTTP response.
type testResponseFunc func() (code int, body string)

// TestHttpClient is a mock implementation of http.Client for testing purposes.
type TestHttpClient struct {
	responses map[string]testResponseFunc // responses stores the URL-to-response function mappings.
}

// testHttpResponse creates a new http.Response with the specified status code and body.
func testHttpResponse(code int, body string) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer([]byte(body))),
		Header:     make(http.Header),
		Status:     http.StatusText(code),
		StatusCode: code,
	}
}

// Do is a method of TestHttpClient, implementing the http.RoundTripper interface.
// It performs a mock HTTP request and returns a mock HTTP response based on the registered URL-to-response mappings.
func (t *TestHttpClient) Do(req *http.Request) (*http.Response, error) {
	fn, ok := t.responses[req.URL.String()]
	if !ok || fn == nil {
		return testHttpResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound)), nil
	}

	return testHttpResponse(fn()), nil
}

// Request registers a URL-to-response function mapping in the TestHttpClient.
func (t *TestHttpClient) Request(url string, fn testResponseFunc) {
	t.responses[url] = fn
}

// NewTestHttpClient creates a new instance of TestHttpClient
func NewTestHttpClient() *TestHttpClient {
	return &TestHttpClient{
		responses: make(map[string]testResponseFunc),
	}
}
