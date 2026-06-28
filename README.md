# gophercraft

[![codecov](https://codecov.io/github/ferminhg/gophercraft/graph/badge.svg?token=2RU5KXN4ZF)](https://codecov.io/github/ferminhg/gophercraft)



A Go project template built with care 🐹 It comes with Hexagonal Architecture, CQRS, and observability out of the box (structured logs and **Prometheus `/metrics`**) — so you can skip the boring setup and focus on what matters: writing great code ✨ Just clone, configure, and start crafting.

## Getting started

You need a recent **Go** toolchain (this repo targets **Go 1.26**). Clone the repository, then download modules and run the API entrypoint:

```bash
go mod download
make run
```

The API listens for HTTP (by default on **`:3000`**; override with `HTTP_ADDR`). A minimal **`GET /status`** route answers `200` with `{"status":"ok"}` so you can verify the stack quickly. **`GET /metrics`** exposes [**Prometheus**](https://prometheus.io/) text format for scraping or for wiring [**Grafana**](https://grafana.com/) with Prometheus as a datasource — see [Metrics (Prometheus)](#metrics-prometheus).

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

## Metrics (Prometheus)

The HTTP server exposes a **`GET /metrics`** endpoint in **Prometheus exposition format**, so a Prometheus server can scrape it and Grafana can chart request rates and latencies.

| Piece | Location | Role |
|-------|----------|------|
| **Port** | [`internal/domain/port/metrics_recorder.go`](internal/domain/port/metrics_recorder.go) | `MetricsRecorder` — records HTTP request samples (`method`, normalized route, status, duration in seconds) without tying the domain to Prometheus. |
| **Adapter** | [`internal/infrastructure/metrics/prometheus_recorder.go`](internal/infrastructure/metrics/prometheus_recorder.go) | **[prometheus/client_golang](https://github.com/prometheus/client_golang)** implementation: dedicated registry (not the global default), suitable for tests. |
| **No-op** | [`internal/infrastructure/metrics/noop_recorder.go`](internal/infrastructure/metrics/noop_recorder.go) | Drops samples — used in HTTP tests and anywhere you do not need metrics. |
| **HTTP** | [`internal/infrastructure/handler/http.go`](internal/infrastructure/handler/http.go) | `metricsMiddleware` emits per-request metrics; **`/metrics`** is served with `promhttp.HandlerFor` when the composition root supplies a `prometheus.Gatherer`. |

**Exported series** (application HTTP traffic; **`GET /metrics`** itself is not counted):

| Metric | Type | Labels |
|--------|------|--------|
| `http_requests_total` | Counter | `method`, `route`, `status_code` |
| `http_request_duration_seconds` | Histogram | `method`, `route` |

At startup, [`cmd/api/main.go`](cmd/api/main.go) constructs `PrometheusRecorder`, injects it into the server as both **`MetricsRecorder`** (middleware) and **`Gatherer`** (scraping). To try locally after `make run`:

```bash
curl -s http://localhost:3000/metrics | head
```

## Project structure

The layout follows **hexagonal architecture** (ports and adapters) with three main areas under `internal/`:

| Path | Role |
|------|------|
| `internal/domain` | **Domain**: entities, value objects, and **ports** (interfaces) that describe what the application needs from the outside world — including [`port.Logger`](internal/domain/port/logger.go) for structured logs and [`port.MetricsRecorder`](internal/domain/port/metrics_recorder.go) for HTTP metrics. |
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

There is a small **dummy** test under `internal/domain/model`, HTTP and middleware tests under `internal/infrastructure/handler`, zerolog adapter tests under `internal/infrastructure/logger`, and Prometheus recorder tests under `internal/infrastructure/metrics` so CI exercises the stack end to end.

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

### Monitoring stack (Prometheus + Grafana)

The [compose.yml](compose.yml) file also starts a small monitoring stack next to the API, so you can scrape and chart the **`/metrics`** endpoint without any extra setup:

| Service | Image | Host port | What it does |
|---------|-------|-----------|--------------|
| `prometheus` | `prom/prometheus:latest` | **9090** | Scrapes the API and stores the time series. It reads its config from [`deploy/prometheus.yml`](deploy/prometheus.yml), which is mounted read-only into the container. |
| `grafana` | `grafana/grafana:latest` | **4000** | Dashboards on top of Prometheus. Grafana also listens on `3000` inside the container, so it is mapped to **4000** on the host to avoid a clash with the API. Default login is `admin` / `admin`. |

Prometheus uses the scrape job defined in [`deploy/prometheus.yml`](deploy/prometheus.yml):

```yaml
global:
  scrape_interval: 5s

scrape_configs:
  - job_name: "gophercraft"
    metrics_path: /metrics
    static_configs:
      - targets: ["api:3000"]
```

Because every service shares the Compose network, Prometheus reaches the API by its service name (**`api:3000`**), not `localhost`. The `prometheus` and `grafana` services declare `depends_on` so they start after the API. Grafana data is kept in the `grafana-data` named volume, so dashboards and settings survive a restart.

After `docker compose up -d --build`, open:

- **API metrics**: http://localhost:3000/metrics
- **Prometheus** (check *Status → Targets* — the `gophercraft` target should be `UP`): http://localhost:9090
- **Grafana** (`admin` / `admin`, then add Prometheus as a datasource at `http://prometheus:9090`): http://localhost:4000

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

This template separates **commands** (`internal/application/command`) from **queries** (`internal/application/query`) so you can grow toward **CQRS** without rewriting the folder layout. **Structured logging** is wired for HTTP and fatal startup errors; extend it by injecting `port.Logger` into application handlers as you add behaviour. **Prometheus metrics** cover HTTP traffic at the adapter layer and can be extended behind `port.MetricsRecorder` the same way.
