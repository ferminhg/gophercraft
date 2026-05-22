package problem_4

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type ExtractionRequest struct {
	URL      string `json:"url"`
	Selector string `json:"selector"`
}

type ParsedRule struct {
	Tag   string
	Class string
}

func parseSelector(selector string) (*ParsedRule, error) {
	if selector == "" || !strings.Contains(selector, ".") {
		return nil, errors.New("invalid CSS selector format")
	}
	parts := strings.Split(selector, ".")
	return &ParsedRule{Tag: parts[0], Class: parts[1]}, nil
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ExtractionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		rule, _ := parseSelector(req.Selector)

		log.Printf("Preparing to extract tag: %s, class: %s from %s", rule.Tag, rule.Class, req.URL)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "rule_accepted"}`))
	}
}
