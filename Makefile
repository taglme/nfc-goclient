export GO111MODULE=on


all: deps lint test build

build:
	go build -mod=vendor -o artifacts/svc

lint:
	golangci-lint run pkg/client

test:
	go test -count=1 -cover -v `go list ./...`

dockerise:
	docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" -f Dockerfile .

deps:
	go mod download
	go mod vendor

deps_check:
	@test -z "$(shell git status -s ./vendor ./go.mod ./go.sum)"