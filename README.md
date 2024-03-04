# RedQueue

RedQueue is a simple Redis and KeyDB-backed queue implementation in Go. It provides an easy way to manage jobs using them as a persistent storage system. 

## How to Install

Firstly, you need to make sure that your environment have Go installed. Afterwards, you can download and install redqueue with the command:

```
go get -u github.com/chand1012/redqueue
```

## Usage

Let's take a look how to use the RedQueue in your Go application.

```go
package main

import (
	queue "github.com/chand1012/redqueue"
	"github.com/redis/go-redis/v9"
)

func main() {
	opts := &redis.Options{
		Addr:     "localhost:6379",
		DB:       0, // use default DB
		Password: "", // no password set
	}

	// Create a new queue named "queue_test"
	q := queue.New(opts, "queue_test")

	type testData struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	data := testData{Id: 1, Name: "test"}

	// Add data to the queue
	err := q.Push(data)
	if err != nil {
		panic(err)
	}

	// add more test data
	data = testData{Id: 2, Name: "test2"}
	err = q.Push(data)
	if err != nil {
		panic(err)
	}

	// Pulls data from the queue as bytes
	fetched, err := q.Process()
	if err != nil {
		panic(err)
	}

	// Unmarshal the pulled data into the a given struct
	var processed testData
	err = q.ProcessInto(fetched, &processed)
	if err != nil {
		panic(err)
	}

	// Mark the task as completed
	err = q.Finish()
	if err != nil {
		panic(err)
	}

	// Close the queue
	q.Close()
}
```

Please note that the address, database, and password for Redis server are passed to the queue with `redis.Options` .

## Tests

There are some comprehensive unit tests provided in this queue library. Make sure you have a local Redis or KeyDB server running, then you can run them with command `go test` .

## Contributing

We welcome contributors, please feel free to enhance this library.
