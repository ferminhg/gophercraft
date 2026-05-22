package problem_3

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthCheckRequest struct {
	TargetURL string `json:"target_url"`
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req HealthCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if req.TargetURL == "" {
			http.Error(w, "target_url is required", http.StatusBadRequest)
			return
		}

		log.Printf("Pinging upstream: %s", req.TargetURL)

		client := &http.Client{}
		resp, err := client.Get(req.TargetURL)
		if err != nil {
			log.Printf("Upstream ping failed: %v", err)
			http.Error(w, "Upstream unreachable", http.StatusBadGateway)
			return
		}

		if resp.StatusCode >= 400 {
			http.Error(w, "Upstream returned error status", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthCheckResponse{Status: "healthy"})
	}
}
