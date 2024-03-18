package main

import (
	"context"
	"github.com/hibiken/asynq"
	"go-labs/cmd/queue/task"
	"go-labs/internal/queue"
	redis_queue "go-labs/internal/queue/redis-queue"
	"log"
	"time"
)

const redisAddr = "127.0.0.1:6379"

type SendEmailTask struct {
	queue.BasePayload

	UserID string `json:"user_id,omitempty"`
}

type emailTask struct {
}

func (e emailTask) Key() queue.TaskIdentifier {
	return "SendEmailTask"
}

func (e emailTask) Handler(ctx context.Context, task *asynq.Task) error {
	log.Printf("payload: %s", task.Payload())
	return nil
}

func NewSendEmailTaskPayload(userID string) *SendEmailTask {
	return &SendEmailTask{
		BasePayload: queue.BasePayload{
			Identifier: "SendEmailTask",
			RunAtTimes: []time.Time{
				time.Now().Add(15 * time.Second),
				time.Now().Add(30 * time.Second),
				time.Now().Add(45 * time.Second),
			},
		},
		UserID: userID,
	}
}

func NewSendEmailTask() *emailTask {
	return &emailTask{}
}

func main() {
	opts := asynq.RedisClientOpt{Addr: redisAddr}
	queuer := redis_queue.NewQueue(opts)
	dequeue := redis_queue.NewDequeue(opts)

	ctx := context.Background()

	if err := queuer.Enqueue(ctx, NewSendEmailTaskPayload("123")); err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(dequeue.Run(NewSendEmailTask()))

	client := asynq.NewClient(opts)
	defer func(client *asynq.Client) {
		_ = client.Close()
	}(client)

	//	Enqueue task to be processed immediately
	emailDeliveryTask, err := task.NewEmailDeliveryTask(1, "template:id")
	if err != nil {
		log.Fatalf("failed to create email delivery task: %v", err)
	}

	emailDeliveryTaskInfo, err := client.Enqueue(emailDeliveryTask)
	if err != nil {
		log.Fatalf("failed to queue email delivery task: %v", err)
	}

	log.Printf("enqueued email delivery task: id=%s queue=%s\n", emailDeliveryTaskInfo.ID, emailDeliveryTaskInfo.Queue)

	//	Schedule task to be processed in the future.
	emailDeliveryTaskInfo, err = client.Enqueue(emailDeliveryTask, asynq.ProcessIn(30*time.Second))
	if err != nil {
		log.Fatalf("failed to queue email delivery task: %v", err)
	}

	log.Printf("enqueued email delivery task: id=%s queue=%s\n", emailDeliveryTaskInfo.ID, emailDeliveryTaskInfo.Queue)

	//	Set other options to tune task processing behavior.
	imageResizeTask, err := task.NewImageResizeTask("https://example.com/myassets/image.jpg")
	if err != nil {
		log.Fatalf("failed to create image resize task: %v", err)
	}

	imageResizeTaskInfo, err := client.Enqueue(imageResizeTask, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("failed to queue image resize task: %v", err)
	}

	log.Printf("enqueued image resize task: id=%s queue=%s\n", imageResizeTaskInfo.ID, imageResizeTaskInfo.Queue)
}
