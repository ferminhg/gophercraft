package problem_6

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type AsyncScrapeRequest struct {
	URL string `json:"url"`
}

func performHeavyScrape(ctx context.Context, url string, done chan<- bool) {
	select {
	case <-time.After(100 * time.Millisecond):
		log.Printf("Finished scraping %s", url)
		done <- true // no bloquea: el canal tiene buffer
	case <-ctx.Done():
		// El cliente se fue: abortamos sin escribir en el canal.
		log.Printf("Scrape aborted for %s: %v", url, ctx.Err())
	}
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AsyncScrapeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		done := make(chan bool, 1)

		go performHeavyScrape(r.Context(), req.URL, done)

		select {
		case <-done:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "completed"}`))
		case <-r.Context().Done():
			log.Printf("Client disconnected or timeout")
			// We return, but the goroutine is leaked!
			http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		}
	}
}
