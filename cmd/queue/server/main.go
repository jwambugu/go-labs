package main

import (
	"github.com/hibiken/asynq"
	"go-labs/cmd/queue/task"
	"log"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr}, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeEmailDelivery, task.HandleEmailDeliveryTask)

	imageProcessor := task.NewImageProcessor()
	mux.HandleFunc(task.TypeImageResize, imageProcessor.ProcessTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
