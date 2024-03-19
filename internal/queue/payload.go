package queue

import "time"

// Payload defines the execution of a job.
//
// # Job acts as an identifier for the task being run
//
// # Queue specifies the queue to run on
//
// RunAt defines the times to run the job.
type Payload interface {
	Job() Job
	Queue() Queue
	RunAt() []time.Time
}

// BasePayload implements Payload
type BasePayload struct {
	JobName    Job         `json:"job,omitempty"`
	OnQueue    Queue       `json:"on_queue,omitempty"`
	RunAtTimes []time.Time `json:"run_at,omitempty"`
}

func (b BasePayload) Job() Job {
	return b.JobName
}

func (b BasePayload) Queue() Queue {
	return b.OnQueue
}

func (b BasePayload) RunAt() []time.Time {
	return b.RunAtTimes
}
