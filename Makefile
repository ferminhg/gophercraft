.PHONY: build test lint run docker-build

build:
	go build ./...

test:
	go test -v -race ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { printf 'golangci-lint not installed; see https://golangci-lint.run/install\n'; exit 1; }
	golangci-lint run ./...

run:
	go run ./cmd/api

IMAGE ?= gophercraft:latest

docker-build:
	docker build -t $(IMAGE) .
