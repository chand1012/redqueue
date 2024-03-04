package queue

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	rdb            *redis.Client
	name           string
	mainName       string
	processingName string
	ctx            context.Context
	cancel         context.CancelFunc
	lastTask       []byte
}

// creates a new queue from a redis client and a name
func FromClient(rdb *redis.Client, name string) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	return &Queue{
		rdb:            rdb,
		name:           name,
		mainName:       name + "_main",
		processingName: name + "_processing",
		ctx:            ctx,
		cancel:         cancel,
		lastTask:       nil,
	}
}

// New creates a new queue with the given name and redis options
func New(opts *redis.Options, name string) *Queue {
	rdb := redis.NewClient(opts)
	return FromClient(rdb, name)
}

// close the queue and the redis client
func (q *Queue) Close() error {
	q.cancel()
	return q.rdb.Close()
}

// takes a string, bytes, a map, or a struct. Will marshal struct or map to JSON. Errors if marshaling fails.
func (q *Queue) Push(data any) error {
	// the data should be bytes, a string, or marshalable to JSON
	t := reflect.TypeOf(data)
	// if its a map or a struct, marshal it to JSON
	if t.Kind() == reflect.Map || t.Kind() == reflect.Struct {
		// marshal to JSON
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return q.rdb.LPush(q.ctx, q.mainName, dataBytes).Err()
	}
	return q.rdb.LPush(q.ctx, q.mainName, data).Err()
}

// pop a task from the queue, marking it as processing
func (q *Queue) Process() ([]byte, error) {
	// pop a task from the main queue
	data, err := q.rdb.RPopLPush(q.ctx, q.mainName, q.processingName).Bytes()
	if err != nil {
		return nil, err
	}
	// store the last task
	q.lastTask = data
	return data, nil
}

// same as process, but takes in a pointer to a struct or map, and unmarshals the data into it
func (q *Queue) ProcessInto(v any) error {
	data, err := q.Process()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// removes a task from the processing queue
func (q *Queue) Finish() error {
	// if there is no last task, return an error
	if q.lastTask == nil {
		return errors.New("no task to finish")
	}
	// using LRem, remove the last task from the processing queue
	err := q.rdb.LRem(q.ctx, q.processingName, 1, q.lastTask).Err()
	if err != nil {
		return err
	}
	q.lastTask = nil
	return nil
}
