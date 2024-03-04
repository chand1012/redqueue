package queue_test

import (
	"context"
	"testing"

	queue "github.com/chand1012/redqueue"
	"github.com/redis/go-redis/v9"
)

type testData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func TestQueue(t *testing.T) {
	opts := &redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	}
	name := "queue_test"

	// first check if redis is running
	// if not throw an error and skip the tests
	rdb := redis.NewClient(opts)
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Errorf("redis is not running: %v", err)
		t.Skip()
	}

	t.Run("New", func(t *testing.T) {
		q := queue.New(opts, name)
		if q == nil {
			t.Errorf("failed to create queue")
		}
	})

	t.Run("Push", func(t *testing.T) {
		q := queue.New(opts, name)
		data := testData{Id: 1, Name: "test"}
		err := q.Push(data)
		if err != nil {
			t.Errorf("failed to push data to queue: %v", err)
		}
	})

	t.Run("Process", func(t *testing.T) {
		q := queue.New(opts, name)
		_, err := q.Process()
		if err != nil {
			t.Errorf("failed to process data from queue: %v", err)
		}
	})

	t.Run("ProcessInto", func(t *testing.T) {
		q := queue.New(opts, name)
		// push the test data to the queue
		q.Push(testData{Id: 1, Name: "test"})
		data := &testData{}
		err := q.ProcessInto(data)
		if err != nil {
			t.Errorf("failed to process data into struct: %v", err)
		}
	})

	t.Run("Finish", func(t *testing.T) {
		q := queue.New(opts, name)
		// push some test data
		q.Push(testData{Id: 1, Name: "test"})
		// process the test data
		_, err := q.Process()
		if err != nil {
			t.Errorf("failed to process data from queue: %v", err)
		}
		// finish processing the task
		err = q.Finish()
		if err != nil {
			t.Errorf("failed to finish processing task: %v", err)
		}
	})

	t.Run("Close", func(t *testing.T) {
		q := queue.New(opts, name)
		err := q.Close()
		if err != nil {
			t.Errorf("failed to close queue: %v", err)
		}
	})
}
