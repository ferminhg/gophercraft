package problem_4_2

import (
	"encoding/json"
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

func parseAuthorization(authorization string) *ParsedToken {
	parts := strings.Split(authorization, " ")
	return &ParsedToken{Scheme: parts[0], Token: parts[1]}
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ScheduleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Mal formed json", http.StatusBadRequest)
			return
		}

		token := parseAuthorization(req.Authorization)

		log.Printf("Scheduling extraction for domain: %s using scheme: %s", req.Domain, token.Scheme)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "scheduled"}`))
	}
}
