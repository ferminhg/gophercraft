# Problem 4_2

**Context:** You are on-call. Workers authenticate against the extraction API by sending an `Authorization` field with each job request via `POST /extract/schedule`, but the endpoint occasionally crashes with a 500 error and no clear log message.

**Expected:** The endpoint should accept a job payload with an `Authorization` field formatted as `"Bearer <token>"`. If the field is missing, empty, or not in that format, it should return a 400 Bad Request instead of crashing.

**Symptom:** The test sends a malformed `Authorization` value (no `Bearer` prefix / no token), which currently causes the handler to panic instead of returning a 400.
