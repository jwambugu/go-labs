package main

import (
	"context"
	"go-labs/internal/dummyjson"
	"log"
	"net/http"
	"sync"
	"time"
)

func waitGroupImpl(ctx context.Context, api dummyjson.Todoer, todos []*dummyjson.Todo) {
	var wg sync.WaitGroup

	for _, todo := range todos {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			resp, err := api.Get(ctx, id)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("todo: %#+v \n", resp)
		}(todo.ID)
	}

	wg.Wait()
}

func channelImpl(ctx context.Context, api dummyjson.Todoer, todos []*dummyjson.Todo) {
	resultChan := make(chan *dummyjson.Todo)
	for _, todo := range todos {
		go func(id int) {
			resp, err := api.Get(ctx, id)
			if err != nil {
				log.Fatalln(err)
			}

			resultChan <- resp
		}(todo.ID)
	}

	for i := 0; i < len(todos); i++ {
		todo := <-resultChan
		log.Printf("todo: %#+v \n", todo)
	}
}

func main() {
	var (
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}

		todoApi = dummyjson.NewTodoApi(httpClient)
		ctx     = context.Background()
	)

	todos, err := todoApi.GetAll(ctx, 5, 2)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("todos: %#+v \n", todos)

	//waitGroupImpl(ctx, todoApi, todos.Todos)
	channelImpl(ctx, todoApi, todos.Todos)
}
