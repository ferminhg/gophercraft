package problem_4_2

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type ScheduleRequest struct {
	Domain        string `json:"domain"`
	Authorization string `json:"authorization"`
}

type ParsedToken struct {
	Scheme string
	Token  string
}

var ErrAuthorizationRequired = errors.New("Authorization is required")
var ErrInvalidAuthorization = errors.New("Invalid authorization")

func NewParse

func parseAuthorization(authorization string) (*ParsedToken, error) {
	parts := strings.Split(authorization, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, ErrInvalidAuthorization
	}
	return &ParsedToken{Scheme: parts[0], Token: parts[1]}, nil
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ScheduleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Mal formed request", http.StatusBadRequest)
			return
		}
		if req.Authorization == "" {
			http.Error(w, "Authorization is required", http.StatusBadRequest)
			return
		}
		token, err := parseAuthorization(req.Authorization)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Scheduling extraction for domain: %s using scheme: %s", req.Domain, token.Scheme)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "scheduled"}`))
	}
}
