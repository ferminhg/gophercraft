# Problem 7

**Context:** You are on-call. Users can optionally update their bandwidth limits via `POST /bandwidth/limit`. However, if they just pass an account ID without the limit config to trigger a default reset, the server crashes.

**Expected:** The endpoint should handle both payloads with and without the optional `limit_config`.

**Symptom:** The test sends a payload without `limit_config`, causing a panic (nil pointer dereference).
