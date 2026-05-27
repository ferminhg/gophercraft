# Problem 2

**Context:** You are on-call. The proxy rotation metrics endpoint (`POST /proxy/rotate`) is crashing heavily during peak traffic spikes, causing 500 errors and missing data.

**Expected:** The endpoint should accept rotation success/failure events and update the global `regionStats` tracker securely.

**Symptom:** Running the tests with the race detector (`go test -race`) or simulating concurrent traffic panics the server.

---

## 🎯 PR Description

Fixes a race condition on the `POST /proxy/rotate` endpoint that caused panics under concurrent traffic.

✨ **What changed**
- Introduced `MutexRegionStats` struct that encapsulates `sync.Mutex` alongside the `regionStats` map, replacing the previous bare global variables
- `RecordRotation` method now owns the lock/unlock cycle, making concurrent writes safe by design

🐛 **Bug fixed** — unprotected concurrent map writes triggered a data race, causing server panics during peak traffic spikes

🔧 **Technical** — mutex ownership moved from the handler call-site into the stats type itself; old global vars replaced with a properly initialised struct via `NewMutexRegionStats()`

📝 **Tests** — `TestHandler_Returns200_WhenFixed` hammers the handler with 100 concurrent goroutines under `-race`; `TestEdgeCase_MissingRegion` validates input validation

🚀 Eliminates 500 errors and missing data during traffic spikes with zero API surface changes
