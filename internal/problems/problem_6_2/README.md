# Problem 6_2

**Context:** You are on-call again. The `POST /proxy/session/refresh` endpoint refreshes a rotating proxy session by handshaking with the upstream proxy provider. Under sustained traffic with slow or disconnecting clients, the API nodes slowly run out of memory.

**Expected:** The endpoint should start the refresh in the background, wait for it to finish, and abandon it cleanly (without leaking anything) if the client disconnects or times out before the refresh completes.

**Symptom:** The test simulates clients disconnecting early. Currently, this causes background goroutines to leak and hang forever instead of exiting.
