.PHONY: build test coverage lint run docker-build load

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

DURATION    ?= 180s
CONNECTIONS ?= 10
THREADS     ?= 2
LOAD_URL    ?= http://api:3000/status

load:
	docker run --rm \
		--network gophercraft_default \
		alpine:3.20 \
		sh -c "apk add --no-cache wrk >/dev/null && wrk -t$(THREADS) -c$(CONNECTIONS) -d$(DURATION) $(LOAD_URL)"
