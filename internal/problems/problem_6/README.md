# Problem 6

**Context:** You are on-call. The async scraping endpoint (`POST /scrape/async`) is slowly eating up memory on the servers until they OOM crash.

**Expected:** The endpoint should start a background job and wait for it, but gracefully handle client disconnects.

**Symptom:** The test simulates clients disconnecting early. Currently, this causes background goroutines to leak and hang forever.
