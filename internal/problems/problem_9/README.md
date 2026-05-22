# Problem 9

**Context:** You are on-call. The customer dashboard relies on `GET /proxy/regions` to display available proxy locations. When the internal system temporarily has no regions available, the endpoint returns a generic 500 error, which causes the dashboard to break entirely instead of showing an "Out of Stock" message.

**Expected:** If `fetchRegions()` returns the known `ErrNoRegionsAvailable` error, the endpoint should return a `404 Not Found` so the client can handle it gracefully.

**Symptom:** The test expects a 404 but receives a 500 because all errors are blanket-masked.
