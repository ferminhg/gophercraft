# Problem 8

**Context:** You are on-call. Users are complaining that they cannot terminate their active sessions using the `DELETE /session/active` endpoint. They receive a 405 Method Not Allowed error.

**Expected:** The endpoint should accept `DELETE` requests and successfully terminate the session.

**Symptom:** The test sends a valid `DELETE` request but receives a 405 because the HTTP method validation is incorrect.
