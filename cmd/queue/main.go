package main

import (
	"context"
	"github.com/hibiken/asynq"
	"go-labs/internal/queue"
	"go-labs/internal/queue/redisqueue"
	"log"
	"time"
)

const redisAddr = "127.0.0.1:6379"

type SendEmailTask struct {
	queue.BasePayload

	UserID string `json:"user_id,omitempty"`
}

type sendWelcomeEmailJob struct {
}

func (s *sendWelcomeEmailJob) Key() queue.Job {
	return queue.JobSendWelcomeEmail
}

func (s *sendWelcomeEmailJob) Handler(ctx context.Context, task *asynq.Task) error {
	log.Printf("payload: %s", task.Payload())
	return nil
}

func NewSendEmailJobPayload(userID string) *SendEmailTask {
	return &SendEmailTask{
		BasePayload: queue.BasePayload{
			JobName: queue.JobSendWelcomeEmail,
			RunAtTimes: []time.Time{
				time.Now().Add(15 * time.Second),
				time.Now().Add(30 * time.Second),
				time.Now().Add(45 * time.Second),
			},
		},
		UserID: userID,
	}
}

func NewSendEmailTask() queue.Worker {
	return &sendWelcomeEmailJob{}
}

func main() {
	var (
		opts      = asynq.RedisClientOpt{Addr: redisAddr}
		publisher = redisqueue.NewPublisher(opts)
		consumer  = redisqueue.NewConsumer(opts)
	)

	ctx := context.Background()

	if err := publisher.Enqueue(ctx, NewSendEmailJobPayload("123")); err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(consumer.Register(NewSendEmailTask()))
}
