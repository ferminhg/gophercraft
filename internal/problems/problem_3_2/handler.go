package problem_3_2

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type WebhookNotifyRequest struct {
	CallbackURL string `json:"callback_url"`
}

type WebhookNotifyResponse struct {
	Status string `json:"status"`
}

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
	}
}

func notifyCallback(client *http.Client, callbackURL string) (int, error) {
	log.Printf("Notifying callback: %s", callbackURL)
	resp, err := client.Get(callbackURL)
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

		var req WebhookNotifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if req.CallbackURL == "" {
			http.Error(w, "callback_url is required", http.StatusBadRequest)
			return
		}

		statusCode, err := notifyCallback(client, req.CallbackURL)
		if err != nil {
			log.Printf("Callback notify failed: %v", err)
			http.Error(w, "Callback unreachable", http.StatusBadGateway)
			return
		}
		if statusCode >= 400 {
			http.Error(w, "Callback returned error status", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(WebhookNotifyResponse{Status: "notified"})
	}
}
