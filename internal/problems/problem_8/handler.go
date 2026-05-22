package problem_8

import (
	"encoding/json"
	"log"
	"net/http"
)

type SessionResponse struct {
	Status string `json:"status"`
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			http.Error(w, "session_id query param is required", http.StatusBadRequest)
			return
		}

		log.Printf("Terminating session: %s", sessionID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SessionResponse{Status: "terminated"})
	}
}
