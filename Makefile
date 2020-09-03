.PHONY: all test testfv build run-linter

all: test build

test:
	go test -count=1 ./internal

build:
	go build ./cmd/dir_cmd_line

testfv:
	go test -failfast -v -count=1 ./internal

run-linter:
	golangci-lint run ./internal