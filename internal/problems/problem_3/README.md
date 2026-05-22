# Problem 3

**Context:** You are on-call. The `/scrape/health` endpoint is causing the API nodes to crash with "too many open files" errors under load.

**Expected:** The endpoint should ping a target URL to check if it's reachable and return a `200 OK` JSON response.

**Symptom:** The test detects a resource leak (goroutines hanging onto open connections).
