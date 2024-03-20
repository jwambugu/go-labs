package redisqueue_test

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
	"go-labs/internal/queue"
	"go-labs/internal/queue/redisqueue"
	"go-labs/internal/testutils"
	"log"
	"os"
	"testing"
	"time"
)

var redisAddr string

func testMain(m *testing.M) int {
	addr, cleanupFunc := testutils.NewRedisClient()
	defer func() {
		if err := cleanupFunc(); err != nil {
			log.Fatalln(err)
		}
	}()

	redisAddr = addr
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

type testJob struct {
	queue.BasePayload

	UserID uint8 `json:"user_id,omitempty"`
}

type testJobWorker struct {
}

func (t testJobWorker) Key() queue.Job {
	return "test:job"
}

func (t testJobWorker) Handler(ctx context.Context, task *asynq.Task) error {
	log.Printf("payload: %s\n", task.Payload())
	return nil
}

func TestNewPublisher_Consumer(t *testing.T) {
	var (
		ctx       = context.Background()
		opts      = asynq.RedisClientOpt{Addr: redisAddr}
		publisher = redisqueue.NewPublisher(opts)
		consumer  = redisqueue.NewConsumer(opts)
	)

	defer func() {
		err := publisher.Close()
		require.NoError(t, err)

	}()

	err := consumer.Register(testJobWorker{})
	require.NoError(t, err)

	err = publisher.Enqueue(ctx, testJob{
		BasePayload: queue.BasePayload{RunAtTimes: []time.Time{time.Now()}},
		UserID:      1,
	})

	require.EqualError(t, err, queue.ErrRequiredJobName.Error())

	err = publisher.Enqueue(ctx, testJob{
		BasePayload: queue.BasePayload{
			JobName:    "test:job",
			OnQueue:    queue.Default,
			RunAtTimes: []time.Time{time.Now()},
		},
		UserID: 1,
	})

	require.NoError(t, err)

	err = publisher.Enqueue(ctx, testJob{
		BasePayload: queue.BasePayload{
			JobName:    "test:job",
			RunAtTimes: []time.Time{time.Now()},
		},
		UserID: 1,
	})

	err = consumer.Close()
	require.NoError(t, err)
}
