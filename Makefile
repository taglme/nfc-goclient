export GO111MODULE=on


all: deps lint test build

lint:
	golangci-lint run pkg/client

test:
	go test -count=1 -cover -v `go list ./...`

deps:
	go mod download
	go mod vendor

deps_check:
	@test -z "$(shell git status -s ./vendor ./go.mod ./go.sum)"