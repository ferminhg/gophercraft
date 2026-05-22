package problem_2

import (
	"encoding/json"
	"log"
	"net/http"
)

type ProxyStatus struct {
	ActiveCount int
	ErrorCount  int
}

// Global state for proxy statistics
var regionStats = make(map[string]ProxyStatus)

type RotateRequest struct {
	Region string `json:"region"`
	Success bool  `json:"success"`
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RotateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if req.Region == "" {
			http.Error(w, "region is required", http.StatusBadRequest)
			return
		}

		log.Printf("Recording proxy rotation for region %s", req.Region)

		// BUG: Concurrent writes to the global map without a mutex
		stats := regionStats[req.Region]
		if req.Success {
			stats.ActiveCount++
		} else {
			stats.ErrorCount++
		}
		regionStats[req.Region] = stats

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "recorded"}`))
	}
}
