export GO111MODULE=on


all: deps lint test build

build:
	go build ./...

lint:
	golangci-lint run pkg/client

test:
	go test -count=1 -cover -v `go list ./...`

deps:
	go mod download
	go mod tidy

deps_check:
	@test -z "$(shell git status -s ./go.mod ./go.sum)"