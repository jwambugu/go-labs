package queue

import (
	"errors"
	"github.com/hibiken/asynq"
	"golang.org/x/net/context"
)

const MaxRetry = 3

type Queue string

func (q Queue) String() string {
	return string(q)
}

const (
	High    Queue = "high"
	Default Queue = "default"
	Low     Queue = "low"
)

const Workers = 10

// Queues is a list of queues to process with given priority value
var Queues = map[string]int{
	High.String():    6,
	Default.String(): 3,
	Low.String():     1,
}

var ErrRequiredJobName = errors.New("job name is required")

// Worker executes the provided Job using the underlying Handler
type Worker interface {
	Key() Job
	Handler(ctx context.Context, task *asynq.Task) error
}

// Consumer consumes the queued Job using the registered Worker
type Consumer interface {
	Register(workers ...Worker) error
}

// Queuer enqueues jobs for processing
type Queuer interface {
	Close() error
	Enqueue(ctx context.Context, payload Payload) error
}

type TestQueue struct {
	payloads []Payload
}

func (t *TestQueue) Close() error {
	return nil
}

func (t *TestQueue) Enqueue(_ context.Context, payload Payload) error {
	t.payloads = append(t.payloads, payload)
	return nil
}

func (t *TestQueue) Len() int {
	return len(t.payloads)
}

func (t *TestQueue) AtIndex(idx int) Payload {
	return t.payloads[idx]
}

func NewTestQueue() *TestQueue {
	return &TestQueue{}
}
