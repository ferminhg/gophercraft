# Problem 10

**Context:** You are on-call. Users are trying to paginate through their scraping results using `GET /scrape/results?page=N&limit=N`. Whenever they reach the final page of their data, the endpoint throws a 500 error.

**Expected:** The endpoint should safely slice the results and return the remaining items, or an empty list if the page is out of bounds.

**Symptom:** The test attempts to fetch the last page, causing the handler to panic with a "slice bounds out of range" error.
