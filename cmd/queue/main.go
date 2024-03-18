package main

import (
	"github.com/hibiken/asynq"
	"go-labs/cmd/queue/task"
	"log"
	"time"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
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
