# Problem 5

**Context:** You are on-call. Users are submitting requests to proxy certain URLs via the `GET /proxy/forward` endpoint. However, whenever a URL contains spaces or special query parameters, the internal system throws a 500.

**Expected:** The endpoint should safely forward the target URL to the internal scraper.

**Symptom:** The test sends a complex URL which currently causes the internal request builder to fail, returning a 500 instead of a 200.
