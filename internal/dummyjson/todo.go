package dummyjson

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const baseURL = "https://dummyjson.com/todos"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Todoer interface {
	Get(ctx context.Context, id int) (*Todo, error)
	GetAll(ctx context.Context, limit int, skip int) (*TodosResp, error)
}

type todoApi struct {
	cl      HttpClient
	baseURL string
}

func (t *todoApi) Get(ctx context.Context, id int) (*Todo, error) {
	uri := fmt.Sprintf("%s/%d", t.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create request - %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do req: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if code := resp.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", code)
	}

	var todo *Todo
	if err = json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}

	return todo, nil
}

func (t *todoApi) GetAll(ctx context.Context, limit int, skip int) (*TodosResp, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request - %v", err)
	}

	queryParams := req.URL.Query()
	queryParams.Add("limit", strconv.Itoa(limit))
	queryParams.Add("skip", strconv.Itoa(skip))

	req.URL.RawQuery = queryParams.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do req: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if code := resp.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", code)
	}

	var todosResp *TodosResp
	if err = json.NewDecoder(resp.Body).Decode(&todosResp); err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}

	return todosResp, nil
}

func NewTodoApi(cl HttpClient) Todoer {
	return &todoApi{
		cl:      cl,
		baseURL: baseURL,
	}
}
