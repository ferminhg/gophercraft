.PHONY: install build test test-watch coverage lint run dev docker-build load

install:
	go install gotest.tools/gotestsum
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint

build:
	go build ./...

test:
	@command -v gotestsum >/dev/null 2>&1 || { printf 'gotestsum not installed; run: go install gotest.tools/gotestsum@latest\n'; exit 1; }
	gotestsum --format-icons pkgname -- -race -coverprofile=coverage.out ./... 

test-watch:
	@command -v gotestsum >/dev/null 2>&1 || { printf 'gotestsum not installed; run: go install gotest.tools/gotestsum@latest\n'; exit 1; }
	gotestsum --watch --format pkgname

coverage: test
	go tool cover -html=coverage.out

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { printf 'golangci-lint not installed; see https://golangci-lint.run/install\n'; exit 1; }
	golangci-lint run ./...

run:
	go run ./cmd/api

dev:
	@command -v air >/dev/null 2>&1 || { printf 'air not installed; run: go install github.com/air-verse/air@latest\n'; exit 1; }
	air -c .air.toml

IMAGE ?= gophercraft:latest

docker-build:
	docker build -t $(IMAGE) .

DURATION    ?= 180s
CONNECTIONS ?= 10
THREADS     ?= 2
LOAD_URL    ?= http://api:3000/status

load:
	docker run --rm \
		--network gophercraft_default \
		alpine:3.20 \
		sh -c "apk add --no-cache wrk >/dev/null && wrk -t$(THREADS) -c$(CONNECTIONS) -d$(DURATION) $(LOAD_URL)"
