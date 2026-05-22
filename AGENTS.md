# Agent Instructions

## Architecture (Hexagonal + CQRS)
- **Domain (`internal/domain`)**: Core entities, value objects, and **Ports** (interfaces). Must remain completely framework-agnostic.
- **Application (`internal/application`)**: Use cases, strictly separated into `command/` and `query/` to follow CQRS. Dependencies always point inward toward the Domain.
- **Infrastructure (`internal/infrastructure`)**: Adapters for HTTP (`gin`), logging (`zerolog`), metrics (`prometheus`), and repositories.
- **Composition Root (`cmd/api`)**: Responsible for wiring concrete adapters to application use cases.

## Commands & Workflows
- **Test Suite**: `make test` (runs with `-v -race -coverprofile`).
- **Single Test**: `go test -v -race -run ^TestName$ ./internal/path/to/pkg`
- **Lint**: `make lint` (runs `golangci-lint`). **Always run `make lint` and `make test` before completing a task** as CI enforces these.
- **Run Locally**: `make run` (listens on `:3000` by default).
- **Docker**: `docker compose up -d --build` (copy `.env.example` to `.env` first if overrides are needed).

## Conventions & Quirks
- **Testing (`testify`)**: Use `require` for fatal setup/precondition checks (halts test immediately). Use `assert` for value comparisons (allows test to continue and report multiple failures).
- **Logging**: Do not use `fmt` or `log` in business logic. Inject `port.Logger` and pass alternating key-value arguments (e.g., `logger.Info("msg", "key", "value")`). Set `LOG_PRETTY=1` in your local environment for readable console output instead of JSON.
- **Metrics**: Endpoints are automatically instrumented for Prometheus (`/metrics`). If domain events need metric tracking, extend `port.MetricsRecorder` and implement it in `infrastructure/metrics`.
