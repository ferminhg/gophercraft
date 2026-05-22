package problem_7

import (
	"encoding/json"
	"log"
	"net/http"
)

type LimitConfig struct {
	MaxBytes int64 `json:"max_bytes"`
	MaxReqs  int   `json:"max_reqs"`
}

type BandwidthRequest struct {
	AccountID   string       `json:"account_id"`
	LimitConfig *LimitConfig `json:"limit_config,omitempty"`
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req BandwidthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if req.AccountID == "" {
			http.Error(w, "account_id is required", http.StatusBadRequest)
			return
		}

		log.Printf("Updating limits for %s to max bytes: %d", req.AccountID, req.LimitConfig.MaxBytes)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "updated"}`))
	}
}
