package problem_5

import (
	"encoding/json"
	"log"
	"net/http"
)

type ForwardResponse struct {
	UpstreamURL string `json:"upstream_url"`
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		targetURL := r.URL.Query().Get("target_url")
		if targetURL == "" {
			http.Error(w, "target_url query param is required", http.StatusBadRequest)
			return
		}

		upstreamReqURL := "http://internal-scraper.local/v1/fetch?url=" + targetURL

		log.Printf("Forwarding to: %s", upstreamReqURL)

		// Simulate the request breaking upstream (in reality client.Get would fail)
		// For the exercise, we will just parse it to show it breaks
		_, err := http.NewRequest(http.MethodGet, upstreamReqURL, nil)
		if err != nil {
			log.Printf("Failed to create upstream request: %v", err)
			http.Error(w, "Internal server error building request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ForwardResponse{UpstreamURL: upstreamReqURL})
	}
}
