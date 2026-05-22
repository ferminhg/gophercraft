package problem_9

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var ErrNoRegionsAvailable = errors.New("no proxy regions currently available")

type RegionsResponse struct {
	Regions []string `json:"regions"`
}

func fetchRegions() ([]string, error) {
	// Simulate checking an upstream system
	// For the sake of this problem, it always returns an expected, domain-specific error
	return nil, ErrNoRegionsAvailable
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		regions, err := fetchRegions()
		if err != nil {
			log.Printf("Failed to fetch regions: %v", err)
			// BUG: Always returns 500, even for expected errors like ErrNoRegionsAvailable
			// which should probably be a 404 or a 200 with an empty list depending on the API design.
			// The test expects a 404 for this specific error.
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(RegionsResponse{Regions: regions})
	}
}
