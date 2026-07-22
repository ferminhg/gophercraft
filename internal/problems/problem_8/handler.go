package problem_8

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SessionResponse struct {
	Status string `json:"status"`
}

func sessionIDByQueryParam(r *http.Request) (string, error) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		return "", errors.New("session_id query param is required")
	}
	return sessionID, nil
}

var ErrSessionIDRequired = errors.New("session_id query param is required")

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sesionResponse SessionResponse
		switch r.Method {
		case http.MethodDelete:
			_, err := sessionIDByQueryParam(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			sesionResponse = SessionResponse{Status: "terminated"}
		case http.MethodGet:
			_, err := sessionIDByQueryParam(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			sesionResponse = SessionResponse{Status: "active"}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(sesionResponse)
	}
}
