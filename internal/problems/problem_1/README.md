# Problem 1

**Context:** You are on-call. Customers are reporting intermittent 500 Internal Server Errors when trying to create new rotating proxy sessions via the `POST /proxy/session` endpoint.

**Expected:** The endpoint should accept a `SessionConfig` JSON payload, inject an internal tracking header, and return a `200 OK` with a session ID.

**Symptom:** The provided test sends a valid payload but the server crashes (panics). It seems to happen specifically when customers don't provide any `custom_headers`.
