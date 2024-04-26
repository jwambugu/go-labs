package dummyjson

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-labs/internal/testutils"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestTodoApi_Get(t *testing.T) {
	var (
		testClient = testutils.NewTestHttpClient()
		ctx        = context.Background()
	)

	testClient.Request(baseURL+"/1", func() (code int, body string) {
		return http.StatusOK, `{
		  "id": 1,
		  "todo": "Do something nice for someone I care about",
		  "completed": true,
		  "userId": 26
		}`
	})

	todoAPI := NewTodoApi(testClient)

	todo, err := todoAPI.Get(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, todo)

	require.Equal(t, 1, todo.ID)
	require.Equal(t, 26, todo.UserID)

	testClient.Request("", func() (code int, body string) {
		return http.StatusNotFound, ""
	})

	todo, err = todoAPI.Get(ctx, 2)
	require.Error(t, err)
	require.Nil(t, todo)
}

func TestTodoApi_GetAll(t *testing.T) {
	var (
		testClient = testutils.NewTestHttpClient()
		ctx        = context.Background()
		limit      = 3
		skip       = 10
	)

	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("skip", strconv.Itoa(skip))

	uri := baseURL + "?" + params.Encode()

	testClient.Request(uri, func() (code int, body string) {
		return http.StatusOK, `
		{
		  "todos": [
			{
			  "id": 11,
			  "todo": "Text a friend I haven't talked to in a long time",
			  "completed": false,
			  "userId": 39
			},
			{
			  "id": 12,
			  "todo": "Organize pantry",
			  "completed": true,
			  "userId": 39
			},
			{
			  "id": 13,
			  "todo": "Buy a new house decoration",
			  "completed": false,
			  "userId": 16
			}
		  ],
		  "total": 150,
		  "skip": 10,
		  "limit": 3
		}`
	})

	todoAPI := NewTodoApi(testClient)

	resp, err := todoAPI.GetAll(ctx, limit, skip)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, limit, resp.Limit)
	require.Equal(t, skip, resp.Skip)
	require.Len(t, resp.Todos, limit)

	testClient.Request("", func() (code int, body string) {
		return http.StatusNotFound, ""
	})

	resp, err = todoAPI.GetAll(ctx, 1, 2)
	require.Error(t, err)
	require.Nil(t, resp)
}
