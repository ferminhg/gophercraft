package problem_1

import (
	"encoding/json"
	"log"
	"net/http"
)

type SessionConfig struct {
	TargetURL     string            `json:"target_url"`
	ProxyRegion   string            `json:"proxy_region"`
	JSRendering   bool              `json:"js_rendering"`
	CustomHeaders map[string]string `json:"custom_headers,omitempty"`
}

type SessionResponse struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var config SessionConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if config.TargetURL == "" {
			http.Error(w, "target_url is required", http.StatusBadRequest)
			return
		}

		// ZenRows internal: inject a tracking trace ID for this session
		traceID := "zr_trace_" + config.ProxyRegion
		
		// BUG: If custom_headers is omitted in JSON, this map is nil and assignment panics
		config.CustomHeaders["X-ZenRows-Trace"] = traceID

		log.Printf("Starting session for %s with trace %s", config.TargetURL, traceID)

		resp := SessionResponse{
			SessionID: "sess_12345",
			Status:    "created",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
