package problem_7

import (
	"encoding/json"
	"errors"
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

var ErrAccountIDRequired = errors.New("account_id is required")
var ErrLimitConfigMissing = errors.New("limit_config is missing")
var ErrBadRequest = errors.New("bad request")

func NewParseBandwidthRequestBody(r *http.Request) (*BandwidthRequest, error) {
	var req BandwidthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrBadRequest
	}
	if req.AccountID == "" {
		return nil, ErrAccountIDRequired
	}
	if req.LimitConfig == nil {
		return nil, ErrLimitConfigMissing
	}
	return &req, nil
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req, err := NewParseBandwidthRequestBody(r)
		if err != nil {
			if errors.Is(err, ErrAccountIDRequired) {
				http.Error(w, "account_id is required", http.StatusBadRequest)
				return
			}
			if errors.Is(err, ErrLimitConfigMissing) {
				http.Error(w, "limit_config is missing", http.StatusOK)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Updating limits for %s to max bytes: %d", req.AccountID, req.LimitConfig.MaxBytes)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "updated"}`))
	}
}
