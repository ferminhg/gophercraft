# Problem 3

**Context:** You are on-call. The `/scrape/health` endpoint is causing the API nodes to crash with "too many open files" errors under load.

**Expected:** The endpoint should ping a target URL to check if it's reachable and return a `200 OK` JSON response.

**Symptom:** The test detects a resource leak (goroutines hanging onto open connections).

---

## Root Cause

The original code was creating a new `http.Client` on every request inside the handler closure. This is an anti-pattern with two concrete consequences:

**1. Connection pool is thrown away on every request.**
`http.Client` is only a thin wrapper. The actual connection pool lives inside `http.Transport`. Each `Transport` owns a map of idle connections keyed by `{scheme, host, port}`. When you create a new client per request, you create a new pool per request — empty, with no reusable connections. The pool is then garbage-collected after the request, making HTTP Keep-Alive and TCP connection reuse impossible by construction.

**2. Goroutine leak from orphaned `persistConn` workers.**
Every live TCP connection in a `Transport` is backed by two dedicated goroutines: a `writeLoop` and a `readLoop`. They stay alive as long as the connection is in the pool. When a new client is created inside the handler and the handler returns, those goroutines have nowhere to return their connection — the pool that owns them is unreachable. They keep running until the OS-level `TIME_WAIT` timeout (typically 60 seconds on Linux), consuming file descriptors and goroutine stack memory. Under load, this exhausts the process file descriptor limit and produces `too many open files` errors.

### How `http.Transport` connection pooling works

```
http.Client
    └── http.Transport
            ├── idleConn     map[connectMethodKey][]*persistConn  ← the pool
            ├── idleConnCh   map[connectMethodKey]chan *persistConn
            └── dialConnFor()  ← creates new TCP/TLS connections
```

The lifecycle of a connection on each request:

1. **Lookup idle conn** — `Transport.roundTrip()` checks `idleConn[{scheme, host}]`. If a free connection exists, it is reused immediately with no TCP handshake overhead.
2. **Dial if needed** — If no idle connection is available, `dialConnFor()` runs the TCP dial and TLS handshake in a goroutine. If another goroutine returns an idle connection in the meantime, the new dial is cancelled to avoid thundering herd.
3. **`persistConn` runs the request** — Each connection has a dedicated `writeLoop` goroutine that serializes outbound requests and a `readLoop` goroutine that reads responses from the TCP socket. Both are always running while the connection is alive.
4. **Body drain returns the conn** — When the response body is fully read (or explicitly drained with `io.Discard`) and `Close()` is called, the connection is returned to `idleConn` for reuse. If the body is not drained, the `Transport` cannot safely reuse the connection and closes it instead.

Key `Transport` parameters that control pool behavior:

```go
&http.Transport{
    MaxIdleConns:        100,            // total idle conns across all hosts
    MaxIdleConnsPerHost: 10,             // idle conns per destination (default: 2)
    MaxConnsPerHost:     0,              // 0 = unlimited active conns per host
    IdleConnTimeout:     90 * time.Second,
}
```

`MaxIdleConnsPerHost` has the most impact in practice. The default value of `2` is too low for services with moderate concurrency against a single upstream host.

## Solution Applied

The `http.Client` is created once in `NewHandler()` and captured by the returned closure. All concurrent requests share the same `Transport` and the same connection pool:

```go
func NewHandler() http.HandlerFunc {
    client := NewHttpClient()  // created once, shared across all requests
    return func(w http.ResponseWriter, r *http.Request) {
        pingUpstream(client, req.TargetURL)
    }
}
```

`pingUpstream` also explicitly drains the response body before closing it, which is required for the `Transport` to return the connection to the pool:

```go
defer resp.Body.Close()
_, _ = io.Copy(io.Discard, resp.Body)  // drain so the conn can be reused
```

`http.Client` is safe for concurrent use by design — it was built to be created once and shared.

## Known Remaining Gaps

- **No custom `Transport`** — the client uses the default `Transport` with `MaxIdleConnsPerHost: 2`. Under high concurrency against a single `target_url`, this creates new connections even when idle ones are available. A tuned `Transport` should be injected via `NewHttpClient()`.
- **`target_url` is fully user-controlled** — there is no allowlist, no scheme restriction, and no private IP block check. This opens a Server-Side Request Forgery (SSRF) vector where an attacker can use the endpoint to probe internal services.
- **No redirect policy** — the default client follows up to 10 redirects automatically, including cross-host redirects. A `CheckRedirect` function should be added to prevent redirect-based SSRF escalation.
