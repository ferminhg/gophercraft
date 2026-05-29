package problem_3

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type HealthCheckRequest struct {
	TargetURL string `json:"target_url"`
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
	}
}

func pingUpstream(client *http.Client, targetURL string) (int, error) {
	log.Printf("Pinging upstream: %s", targetURL)
	resp, err := client.Get(targetURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

func NewHandler() http.HandlerFunc {
	client := NewHttpClient()
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

		statusCode, err := pingUpstream(client, req.TargetURL)
		if err != nil {
			log.Printf("Upstream ping failed: %v", err)
			http.Error(w, "Upstream unreachable", http.StatusBadGateway)
			return
		}
		if statusCode >= 400 {
			http.Error(w, "Upstream returned error status", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthCheckResponse{Status: "healthy"})
	}
}
