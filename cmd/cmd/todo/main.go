package main

import (
	"context"
	"go-labs/internal/dummyjson"
	"log"
	"net/http"
	"sync"
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

	todos, err := todoApi.GetAll(ctx, 5, 2)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("todos: %#+v \n", todos)

	var wg sync.WaitGroup

	for _, t := range todos.Todos {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			todo, err := todoApi.Get(ctx, id)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("todo: %#+v \n", todo)
		}(t.ID)
	}

	wg.Wait()
}
