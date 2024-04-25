package main

import (
	"context"
	"go-labs/internal/dummyjson"
	"log"
	"net/http"
	"time"
)

func main() {
	var (
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}

		todoApi = dummyjson.NewTodoApi(httpClient)
		ctx     = context.Background()
	)

	todo, err := todoApi.Get(ctx, 1)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("todo: %#+v \n", todo)

	todos, err := todoApi.GetAll(ctx, 5, 2)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("todos: %#+v \n", todos)
}
