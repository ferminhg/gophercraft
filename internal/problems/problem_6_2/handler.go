package problem_6_2

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SessionRefreshRequest struct {
	ProxyRegion string `json:"proxy_region"`
}

type SessionRefreshResponse struct {
	Status string `json:"status"`
}

func refreshProxySession(ctx context.Context, proxyRegion string, done chan<- bool) {
	select {
	case <-time.After(100 * time.Millisecond): // en real: la llamada de red usando ctx
		done <- true
	case <-ctx.Done():
		log.Printf("Proxy session refresh cancelled for region: %s", proxyRegion)
	}
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SessionRefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if req.ProxyRegion == "" {
			http.Error(w, "proxy_region is required", http.StatusBadRequest)
			return
		}

		done := make(chan bool, 1)

		go refreshProxySession(r.Context(), req.ProxyRegion, done)

		select {
		case <-done:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(SessionRefreshResponse{Status: "refreshed"})
		case <-r.Context().Done():
			log.Printf("Client disconnected while refreshing proxy session")
			// The handler returns here, but the goroutine keeps running.
			http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		}
	}
}
