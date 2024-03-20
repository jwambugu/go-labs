package queue

// Job represents a task to be done on the queue
type Job string

func (j Job) String() string {
	return string(j)
}

const (
	JobSendWelcomeEmail Job = "send:welcome:email"
)
