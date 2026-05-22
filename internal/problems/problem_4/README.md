# Problem 4

**Context:** You are on-call. Users are submitting new extraction rules via `POST /extract/rule`, but occasionally the endpoint crashes with a 500 error and no clear log message.

**Expected:** The endpoint should accept a valid CSS selector and URL. If the selector is invalid, it should return a 400 Bad Request.

**Symptom:** The test sends an invalid selector, which currently causes the handler to panic instead of returning a 400.
