.PHONY: build test coverage lint run docker-build

build:
	go build ./...

test:
	go test -v -race -coverprofile=coverage.out ./...

coverage: test
	go tool cover -html=coverage.out

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { printf 'golangci-lint not installed; see https://golangci-lint.run/install\n'; exit 1; }
	golangci-lint run ./...

run:
	go run ./cmd/api

IMAGE ?= gophercraft:latest

docker-build:
	docker build -t $(IMAGE) .
