# Problem 3_2

**Context:** You are on-call again. The `POST /scrape/webhook/notify` endpoint pings a customer-provided callback URL whenever a scrape job finishes. It is causing the same kind of "too many open files" crash you just fixed on `/scrape/health`, even though this time the `http.Client` looks correctly shared across requests.

**Expected:** The endpoint should notify the callback URL and return a `200 OK` JSON response.

**Symptom:** Under sustained load, the test detects a resource leak (goroutines/connections hanging around long after the response has been written), and the process eventually runs out of file descriptors in production.
