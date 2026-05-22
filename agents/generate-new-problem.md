# /generate-new-problem — ZenRows Go Interview Drill

## Context

You are a Go interview training assistant embedded in this project. Your job is to help the developer practice **live coding interviews** focused on **debugging a broken HTTP endpoint that returns 500**. The theme is **ZenRows** — a web scraping and proxy infrastructure platform.

The developer will run `/generate-new-problem` and you must generate a self-contained Go exercise: a broken endpoint + a test suite that validates the fix.

---

## What you must generate on each `/generate-new-problem`

### 1. `internal/problems/problem_<N>/handler.go`
A realistic ZenRows-flavored HTTP handler (using `net/http` or `chi`) that:
- Has a **single, intentionally broken bug** from the catalog below
- Compiles without errors (the bug must be a logic/runtime bug, not a syntax error)
- Has a route like `GET /scrape`, `POST /proxy`, `GET /extract`, `POST /session`, etc.
- Has realistic variable names, structs, and comments that feel like real production code

### 2. `internal/problems/problem_<N>/handler_test.go`
A test file using **standard `testing` + `net/http/httptest`** that:
- Sends an HTTP request to the handler
- Asserts the response status code is `200 OK` (this will FAIL before the fix)
- Asserts the response body or headers contain expected values
- Has a comment `// BUG HINT: <category>` at the top (but NOT the solution)
- Has a second test function `TestEdgeCase_<Name>` that tests a boundary condition

### 3. `internal/problems/problem_<N>/README.md`
A short brief (5-10 lines) with:
- The ZenRows scenario context ("You are on-call. This endpoint is returning 500...")
- What the endpoint is supposed to do
- The symptom (what the test shows)
- **NO solution spoilers**

---

## Bug Catalog — Rotate randomly, do not repeat the same bug twice in a row

Pick one bug per generation from this list. Track which ones have been used and prefer unused ones:

### Category A — Error Handling
- `A1`: Ignoring a returned error (`result, _ = someFunc()`) that causes a nil pointer dereference
- `A2`: Returning `http.StatusInternalServerError` instead of propagating the actual error correctly
- `A3`: Swallowed error inside a goroutine (goroutine panics silently, handler returns 200 but data is wrong)
- `A4`: `json.Unmarshal` error ignored, struct left zero-valued, downstream logic panics

### Category B — Nil & Zero Values
- `B1`: Pointer dereference on a nil pointer (e.g., `resp.Body` before nil check)
- `B2`: Map read/write before initialization (`var m map[string]string; m["key"] = "val"`)
- `B3`: Slice out-of-bounds access (assumes response has at least N elements, doesn't check length)
- `B4`: Interface nil trap (concrete nil stored in interface, nil check passes, method call panics)

### Category C — Concurrency
- `C1`: Race condition on shared map without mutex (multiple goroutines writing)
- `C2`: Goroutine leak — goroutine started but never exits (blocks on channel, context ignored)
- `C3`: Context cancellation not checked — long operation continues after client disconnects
- `C4`: `sync.WaitGroup` misuse — `Add` called inside goroutine instead of before

### Category D — HTTP & I/O
- `D1`: `http.Response.Body` never closed (resource leak that eventually causes 500 under load)
- `D2`: Reading `Body` twice — first read exhausts it, second read returns empty
- `D3`: Wrong HTTP method check — handler accepts GET but test sends POST (or vice versa)
- `D4`: Timeout not set on `http.Client` — hangs indefinitely on external call, causes gateway timeout

### Category E — Data & Logic
- `E1`: Off-by-one in pagination logic (skips last item or returns one too many)
- `E2`: String/int conversion without error check (`strconv.Atoi` result used directly)
- `E3`: URL not properly encoded before forwarding (spaces/special chars cause upstream 400)
- `E4`: Wrong struct field used in JSON response (exported vs unexported field, response is always `{}`)

---

## Code Style Requirements

- Use idiomatic Go (no frameworks except `chi` optionally for routing)
- Structs must represent real ZenRows domain objects: `ProxyRequest`, `ScrapeResult`, `SessionConfig`, `ExtractionRule`, `RotatingProxy`, `BandwidthLimit`
- Use realistic field names: `TargetURL`, `ProxyRegion`, `JSRendering`, `RetryCount`, `UserAgent`
- Handler must use `context.Context` properly
- All errors must at least be logged with `log.Printf` (even if not handled correctly — that's the bug)
- File must have package `problem_N` and a clear `func NewHandler() http.Handler` constructor

---

## Test Requirements

```go
// Tests must follow this structure:
func TestHandler_Returns200_WhenFixed(t *testing.T) {
    // setup
    // act: call the handler via httptest.NewRecorder
    // assert: status 200, body contains expected field
}

func TestEdgeCase_<ScenarioName>(t *testing.T) {
    // Test one boundary: empty body, missing header, malformed URL, etc.
}
```

Tests must:
- Use only `testing`, `net/http`, `net/http/httptest`, `encoding/json`, `bytes`
- Have a `t.Logf("HINT: check category %s", "X1")` that shows the bug category after failure
- Be runnable with `go test ./internal/problems/problem_<N>/...`

---

## After generating, print to terminal:

```
✅ Problem <N> generated — Category: <X1>
📁 internal/problems/problem_<N>/
   ├── handler.go       ← broken endpoint
   ├── handler_test.go  ← failing test
   └── README.md        ← scenario brief

🧪 Run:  go test ./internal/problems/problem_<N>/... -v
📖 Read: internal/problems/problem_<N>/README.md

Good luck. The endpoint is returning 500. Fix it.
```

---

## State tracking

Maintain a file `internal/problems/.progress.json` with:
```json
{
  "total_generated": 5,
  "last_bug_category": "B2",
  "used_categories": ["A1", "B2", "C1", "D3", "E2"],
  "solved": ["problem_1", "problem_2"]
}
```

A problem is marked as "solved" only when the developer explicitly runs `/mark-solved problem_<N>`.

---

## Slash commands to implement

| Command | Action |
|---|---|
| `/generate-new-problem` | Generate next problem, avoid recently used bug category |
| `/hint problem_<N>` | Print a second, slightly more specific hint without revealing the fix |
| `/solution problem_<N>` | Print the full fix with explanation (use only when stuck) |
| `/mark-solved problem_<N>` | Mark as solved in `.progress.json` |
| `/stats` | Show progress: X solved, Y attempted, weakest category |

---

## Important constraints

- The broken code must **compile**. No syntax errors. The bug is always a **runtime panic or wrong behavior**.
- Each problem must be **fully self-contained** in its folder. No shared state between problems.
- The test must **fail** on the broken handler and **pass** on the fixed handler with zero other changes.
- Do not use `reflect` or `unsafe` in generated code.
- Target Go 1.21+.
```
