package queue

import (
	"errors"
	"github.com/hibiken/asynq"
	"golang.org/x/net/context"
	"time"
)

const MaxRetry = 3

type Priority string

func (p Priority) String() string {
	return string(p)
}

const (
	PriorityHigh    = "critical"
	PriorityDefault = "default"
	PriorityLow     = "low"
)

const Workers = 10

// Queues is a list of queues to process with given priority value
var Queues = map[string]int{
	PriorityHigh:    6,
	PriorityDefault: 3,
	PriorityLow:     1,
}

var ErrRequiredIdentifier = errors.New("identifier is required")
var ErrRequiredPriority = errors.New("priority is required")

type TaskIdentifier string

func (t TaskIdentifier) String() string {
	return string(t)
}

type Worker interface {
	Key() TaskIdentifier
	Handler(ctx context.Context, task *asynq.Task) error
}

type Processor interface {
	Run(workers ...Worker) error
}

type Queuer interface {
	Close() error
	Enqueue(ctx context.Context, payload Payload) error
}

type Payload interface {
	Priority() Priority
	RunAt() []time.Time
	TaskIdentifier() TaskIdentifier
}

type BasePayload struct {
	Identifier    TaskIdentifier `json:"identifier,omitempty"`
	PriorityLevel Priority       `json:"on_priority,omitempty"`
	RunAtTimes    []time.Time    `json:"run_at,omitempty"`
}

func (b BasePayload) Priority() Priority {
	return b.PriorityLevel
}

func (b BasePayload) RunAt() []time.Time {
	return b.RunAtTimes
}

func (b BasePayload) TaskIdentifier() TaskIdentifier {
	return b.Identifier
}
