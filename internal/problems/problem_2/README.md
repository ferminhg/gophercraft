# Problem 2

**Context:** You are on-call. The proxy rotation metrics endpoint (`POST /proxy/rotate`) is crashing heavily during peak traffic spikes, causing 500 errors and missing data.

**Expected:** The endpoint should accept rotation success/failure events and update the global `regionStats` tracker securely.

**Symptom:** Running the tests with the race detector (`go test -race`) or simulating concurrent traffic panics the server.
