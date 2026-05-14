# gophercraft

A Go project template built with care 🐹 It comes with Hexagonal Architecture, CQRS, and observability out of the box — so you can skip the boring setup and focus on what matters: writing great code ✨ Just clone, configure, and start crafting.

## Getting started

You need a recent **Go** toolchain (this repo targets **Go 1.26**). Clone the repository, then download modules and run the API entrypoint:

```bash
go mod download
make run
```

The program prints a short startup message. The HTTP server is not fully wired yet; that is intentional so you can extend the driving adapters step by step.

## Project structure

The layout follows **hexagonal architecture** (ports and adapters) with three main areas under `internal/`:

| Path | Role |
|------|------|
| `internal/domain` | **Domain**: entities, value objects, and **ports** (interfaces) that describe what the application needs from the outside world. |
| `internal/application` | **Application**: use cases. Commands and queries are separated to keep a **CQRS**-friendly shape. |
| `internal/infrastructure` | **Infrastructure**: adapters — HTTP handlers (driving) and repositories (driven) that implement the domain ports. |

The `cmd/api` package is the **composition root**: it wires concrete adapters to handlers. Shared, stable helpers can live in `pkg/` when they are safe to import from other modules.

## Running tests

Run the full test suite with the race detector:

```bash
make test
```

There is a small **dummy** test under `internal/domain/model` so CI has something green from day one.

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

This template separates **commands** (`internal/application/command`) from **queries** (`internal/application/query`) so you can grow toward **CQRS** without rewriting the folder layout. Observability and production hardening can be added in infrastructure and at the edges as the service matures.
