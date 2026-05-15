# gophercraft

[![codecov](https://codecov.io/github/ferminhg/gophercraft/graph/badge.svg?token=2RU5KXN4ZF)](https://codecov.io/github/ferminhg/gophercraft)



A Go project template built with care 🐹 It comes with Hexagonal Architecture, CQRS, and observability out of the box — so you can skip the boring setup and focus on what matters: writing great code ✨ Just clone, configure, and start crafting.

## Getting started

You need a recent **Go** toolchain (this repo targets **Go 1.26**). Clone the repository, then download modules and run the API entrypoint:

```bash
go mod download
make run
```

The API listens for HTTP (by default on **`:3000`**; override with `HTTP_ADDR`). A minimal **`GET /status`** route answers `200` with `{"status":"ok"}` so you can verify the stack quickly.

## Logging

The service uses **structured logging** so you avoid ad-hoc `fmt`/`log` at the edges and keep output consistent for operators and log aggregators.

| Piece | Location | Role |
|-------|----------|------|
| **Port** | [`internal/domain/port/logger.go`](internal/domain/port/logger.go) | `Logger` — framework-agnostic API (`Info`, `Warn`, `Error`, `Debug`, `Fatal`, `With`). Args are alternating **key-value** pairs (same idea as `log/slog`). |
| **Adapter** | [`internal/infrastructure/logger`](internal/infrastructure/logger) | **[zerolog](https://github.com/rs/zerolog)** implementation: JSON lines by default, optional human-readable console output in development. |
| **HTTP** | [`internal/infrastructure/handler/middleware_logger.go`](internal/infrastructure/handler/middleware_logger.go) | Replaces Gin’s default request logger with a line per request (`method`, `path`, `status`, `latency_ms`). |

At startup, [`cmd/api/main.go`](cmd/api/main.go) builds the adapter and injects it into the HTTP server. Use cases can accept `port.Logger` later the same way.

**Environment variables**

| Variable | Default | Meaning |
|----------|---------|---------|
| `LOG_LEVEL` | `info` | Zerolog level: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`, or `disabled`. Invalid values fall back to `info`. |
| `LOG_PRETTY` | (off) | Set to `true` or `1` for **colored, multi-field console** lines instead of JSON. Use this for **local** or `docker compose` terminals; keep JSON (`LOG_PRETTY` unset) in production so **Loki, ELK, Cloud Logging**, etc. can parse one object per line. |
| `HTTP_ADDR` | `:3000` | Listen address for the HTTP server. |
| `OTEL_SERVICE_NAME` | — | If set, becomes **`service.name`** on every log line ([OpenTelemetry resource](https://opentelemetry.io/docs/specs/semconv/resource/#service)). |
| `SERVICE_NAME` | — | Used as **`service.name`** when `OTEL_SERVICE_NAME` is empty. |
| *(implicit default)* | `gophercraft` | **`service.name`** when neither of the above is set. |
| `DEPLOYMENT_ENVIRONMENT` | — | If set, becomes **`deployment.environment`** (OTel-style). |
| `ENV` / `APP_ENV` | — | Fallback for **`deployment.environment`** when `DEPLOYMENT_ENVIRONMENT` is empty. |

In Docker or Kubernetes, **JSON + `service.name` + `deployment.environment`** is the usual pattern so you can filter and correlate logs across services.

## Project structure

The layout follows **hexagonal architecture** (ports and adapters) with three main areas under `internal/`:

| Path | Role |
|------|------|
| `internal/domain` | **Domain**: entities, value objects, and **ports** (interfaces) that describe what the application needs from the outside world — including [`port.Logger`](internal/domain/port/logger.go) for structured logs. |
| `internal/application` | **Application**: use cases. Commands and queries are separated to keep a **CQRS**-friendly shape. |
| `internal/infrastructure` | **Infrastructure**: adapters — HTTP handlers (driving), logging (zerolog), and repositories (driven) that implement the domain ports. |

The `cmd/api` package is the **composition root**: it wires concrete adapters to handlers. Shared, stable helpers can live in `pkg/` when they are safe to import from other modules.

## Running tests

Run the full test suite with the race detector:

```bash
make test
```

Tests use **[testify](https://pkg.go.dev/github.com/stretchr/testify)** for assertions:

- **`github.com/stretchr/testify/require`** — fatal checks (`t.FailNow`). Use for preconditions and setup (for example constructing a server or decoding a response) so the test stops immediately on failure.
- **`github.com/stretchr/testify/assert`** — non-fatal checks. Use for comparing values (status codes, fields) when you still want clearer failures in one place.

Docs: [testify on pkg.go.dev](https://pkg.go.dev/github.com/stretchr/testify).

There is a small **dummy** test under `internal/domain/model`, HTTP and middleware tests under `internal/infrastructure/handler`, and zerolog adapter tests under `internal/infrastructure/logger` so CI exercises the stack end to end.

## Continuous integration (GitHub Actions)

This repository ships with a **GitHub Actions** workflow at [`.github/workflows/ci.yml`](.github/workflows/ci.yml). It runs on **pull requests** and on **pushes to `main`**, and includes:

- **Testing**: `go test -v -race ./...` (same idea as `make test`)
- **Linting**: **golangci-lint** using the checked-in [`.golangci.yml`](.golangci.yml)

You get the same checks locally with `make test` and `make lint`.

## Docker

Build the container image (tag defaults to `gophercraft:latest`; override with `IMAGE`):

```bash
make docker-build
```

Run the stack with **Docker Compose** (uses [compose.yml](compose.yml)). Create a `.env` file when you need overrides; you can start from **`.env.example`**:

```bash
cp .env.example .env   # optional
docker compose up -d --build
```

The `api` service maps **port 3000** on the host to **3000** in the container. Adjust `ports`, `environment`, or `.env` as needed.

Stop and remove containers for this project:

```bash
docker compose down
```

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | `go build ./...` |
| `make test` | `go test -v -race ./...` |
| `make lint` | `golangci-lint run ./...` (requires `golangci-lint` on your PATH) |
| `make run` | `go run ./cmd/api` |
| `make docker-build` | Builds the Docker image (`IMAGE` overrides the tag) |

## Architecture (short overview)

**Hexagonal architecture** keeps **domain rules** in the centre. **Driving** adapters (for example HTTP) call into the **application** layer. **Driven** adapters (for example a database or in-memory store) implement **ports** defined next to the domain. Dependencies point **inward**, so the domain does not know about frameworks or IO details.

This template separates **commands** (`internal/application/command`) from **queries** (`internal/application/query`) so you can grow toward **CQRS** without rewriting the folder layout. **Structured logging** is wired for HTTP and fatal startup errors; extend it by injecting `port.Logger` into application handlers as you add behaviour.
