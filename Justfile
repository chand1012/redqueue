set dotenv-load

default:
  just --list --unsorted

tidy:
  go mod tidy

gen-test path prompt="Write unit tests for all functions in the given file.":
  #!/bin/bash
  NEW_FILE=$(echo {{path}} | sed 's/\.go/_test.go/')
  otto edit $NEW_FILE -c {{path}} -g "{{prompt}}"

add command:
  cobra-cli add {{command}}

clean:
  go clean -testcache
  go clean -modcache

start-keydb:
  docker run --rm -d --name keydb -p 6379:6379 eqalpha/keydb

stop-keydb:
  docker stop keydb

test: start-keydb
  sleep 1
  go test -v ./...
  just stop-keydb
