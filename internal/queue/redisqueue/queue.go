package redisqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"go-labs/internal/queue"
)

type publisher struct {
	client *asynq.Client
}

type consumer struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func (p *publisher) Close() error {
	return p.client.Close()
}

func (p *publisher) Enqueue(ctx context.Context, payload queue.Payload) error {
	if payload.Job() == "" {
		return queue.ErrRequiredJobName
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %v", err)
	}

	job := payload.Job().String()
	task := asynq.NewTask(job, b)

	priority := payload.Queue().String()
	if priority == "" {
		priority = queue.Default
	}

	for _, time := range payload.RunAt() {
		_, err = p.client.EnqueueContext(ctx, task,
			asynq.Queue(priority),
			asynq.MaxRetry(queue.MaxRetry),
			asynq.ProcessAt(time),
		)

		if err != nil {
			return fmt.Errorf("enqueue %q at %q: %v", job, time.String(), err)
		}
	}

	return nil
}

func NewPublisher(opts asynq.RedisClientOpt) queue.Queuer {
	return &publisher{
		client: asynq.NewClient(opts),
	}
}

func (c *consumer) Register(workers ...queue.Worker) error {
	for _, worker := range workers {
		c.mux.HandleFunc(worker.Key().String(), worker.Handler)
	}

	return c.srv.Run(c.mux)
}

func NewConsumer(opts asynq.RedisClientOpt) queue.Consumer {
	srv := asynq.NewServer(
		opts,
		asynq.Config{
			Concurrency: queue.Workers,
			Queues:      queue.Queues,
		},
	)

	return &consumer{
		srv: srv,
		mux: asynq.NewServeMux(),
	}
}
