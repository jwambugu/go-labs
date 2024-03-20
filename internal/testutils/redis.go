package testutils

import (
	"context"
	"github.com/ory/dockertest"
	"github.com/redis/go-redis/v9"
	"log"
)

type CleanupFunc func() error

func NewRedisClient() (string, CleanupFunc) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("redis pool: %v", err)
	}

	if err = pool.Client.Ping(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}

	resource, err := pool.Run("redis", "alpine3.19", nil)
	if err != nil {
		log.Fatalf("redis run: %v", err)
	}

	var addr string
	if err = pool.Retry(func() error {
		addr = "localhost:" + resource.GetPort("6379/tcp")

		client := redis.NewClient(&redis.Options{Addr: addr})
		return client.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatalf("redis connect: %v", err)
	}

	return addr, func() error {
		return pool.Purge(resource)
	}
}
