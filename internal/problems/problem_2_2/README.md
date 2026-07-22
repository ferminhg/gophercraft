# Problem 2_2

**Context:** You are on-call. The scrape job stats endpoint (`POST /scrape/job/complete`) is panicking under concurrent load when multiple worker goroutines report job completions at the same time.

**Expected:** The endpoint should accept job completion events (`domain`, `success`) and update the global `jobStats` tracker with per-domain success/failure counts.

**Symptom:** Running the tests with the race detector (`go test -race`) or simulating concurrent traffic panics the server with a "concurrent map writes" fatal error.
