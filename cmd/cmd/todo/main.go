package main

import (
	"context"
	"fmt"
	"go-labs/internal/dummyjson"
	"golang.org/x/sync/errgroup"
	"log"
	"runtime"
	"sync"
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

func limitGoroutines(ctx context.Context, api dummyjson.Todoer, todos []*dummyjson.Todo) {
	var (
		maxGoroutines = runtime.GOMAXPROCS(runtime.NumCPU())
		limiter       = make(chan struct{}, maxGoroutines)
	)

	for _, todo := range todos {
		limiter <- struct{}{}

		go func(id int) {
			defer func() {
				<-limiter
			}()

			resp, err := api.Get(ctx, id)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("%+#v\n", resp)
		}(todo.ID)
	}

	for i := 0; i < cap(limiter); i++ {
		limiter <- struct{}{}
	}
}

func errGroup(ctx context.Context, api dummyjson.Todoer, todos []*dummyjson.Todo) {
	errGrp := errgroup.Group{}

	for _, todo := range todos {
		todo := todo

		errGrp.Go(func() error {
			resp, err := api.Get(ctx, todo.ID)
			if err != nil {
				return err
			}

			log.Println(resp.ID)
			return nil
		})
	}

	if err := errGrp.Wait(); err != nil {
		log.Fatalln(err)
	}
}

func errGroupCtx(ctx context.Context, api dummyjson.Todoer, todos []*dummyjson.Todo) {
	errGrp, gCtx := errgroup.WithContext(ctx)

	for _, todo := range todos {
		todo := todo
		errGrp.Go(func() error {
			for {
				select {
				case <-gCtx.Done():
					return gCtx.Err()
				default:
					resp, err := api.Get(ctx, todo.ID)
					if err != nil {
						return err
					}

					log.Println(resp.ID)
					return nil
				}
			}
		})
	}

	if err := errGrp.Wait(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	//var (
	//	httpClient = &http.Client{
	//		Timeout: 10 * time.Second,
	//	}
	//
	//	todoApi = dummyjson.NewTodoApi(httpClient)
	//	ctx     = context.Background()
	//)
	//
	//todos, err := todoApi.GetAll(ctx, 10, 2)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("todos: %#+v \n", len(todos.Todos))

	//waitGroupImpl(ctx, todoApi, todos.Todos)
	//channelImpl(ctx, todoApi, todos.Todos)
	//limitGoroutines(ctx, todoApi, todos.Todos)
	//errGroup(ctx, todoApi, todos.Todos)
	//errGroupCtx(ctx, todoApi, todos.Todos)

	var x any
	x = 1

	switch v := x.(type) {
	case string:
		log.Println("string", v)
	case int:
		log.Println("int", v)
	}

	if val, ok := x.(int); ok {
		log.Println("int", val)
	}

	//sendOnly := make(chan<- struct{})
	//receiveOnly := make(<-chan struct{})

	numbers := make(chan int, 4)
	numbers <- 1
	numbers <- 2
	numbers <- 3
	numbers <- 4

	fmt.Println(<-numbers)
	fmt.Println(<-numbers)
	fmt.Println(<-numbers)
	fmt.Println(<-numbers)
}
