package queue_test

import (
	"github.com/stretchr/testify/require"
	"go-labs/internal/queue"
	"testing"
	"time"
)

func TestNewTestQueue(t *testing.T) {
	q := queue.NewTestQueue()

	defer func() {
		err := q.Close()
		require.NoError(t, err)
	}()

	err := q.Enqueue(nil, queue.BasePayload{
		JobName:    "NewTestQueue",
		OnQueue:    queue.Low,
		RunAtTimes: []time.Time{time.Now()},
	})

	require.NoError(t, err)
	require.Equal(t, 1, q.Len())
	require.Equal(t, queue.Low, q.AtIndex(0).Queue())
	require.WithinDuration(t, time.Now(), q.AtIndex(0).RunAt()[0], time.Second)
}
