package redis_queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"go-labs/internal/queue"
)

type distributor struct {
	client *asynq.Client
}

type dequeue struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func (r *distributor) Close() error {
	return r.client.Close()
}

func (r *distributor) Enqueue(ctx context.Context, payload queue.Payload) error {
	if payload.TaskIdentifier() == "" {
		return queue.ErrRequiredIdentifier
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %v", err)
	}

	taskIdentifier := payload.TaskIdentifier().String()
	task := asynq.NewTask(taskIdentifier, b)

	priority := payload.Priority().String()
	if priority == "" {
		priority = queue.PriorityDefault
	}

	for _, time := range payload.RunAt() {
		_, err = r.client.EnqueueContext(ctx, task,
			asynq.Queue(priority),
			asynq.MaxRetry(queue.MaxRetry),
			asynq.ProcessAt(time),
		)

		if err != nil {
			return fmt.Errorf("enqueue %q at %q: %v", taskIdentifier, time.String(), err)
		}
	}

	return nil
}

func NewQueue(opts asynq.RedisClientOpt) queue.Queuer {
	return &distributor{
		client: asynq.NewClient(opts),
	}
}

func (d *dequeue) Run(workers ...queue.Worker) error {
	for _, worker := range workers {
		d.mux.HandleFunc(worker.Key().String(), worker.Handler)
	}

	return d.srv.Run(d.mux)
}

func NewDequeue(opts asynq.RedisClientOpt) queue.Processor {
	srv := asynq.NewServer(
		opts,
		asynq.Config{
			Concurrency: 10,
			Queues:      queue.Queues,
		},
	)

	return &dequeue{
		srv: srv,
		mux: asynq.NewServeMux(),
	}
}
