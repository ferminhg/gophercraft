package problem_6

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type AsyncScrapeRequest struct {
	URL string `json:"url"`
}

func performHeavyScrape(url string, done chan<- bool) {
	// Simulating a long-running scrape
	time.Sleep(100 * time.Millisecond)
	log.Printf("Finished scraping %s", url)
	done <- true
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

		// BUG: Unbuffered channel. If the client disconnects or the context times out before
		// the goroutine finishes, the handler returns, but the goroutine blocks forever
		// trying to send on `done`.
		done := make(chan bool)

		go performHeavyScrape(req.URL, done)

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
