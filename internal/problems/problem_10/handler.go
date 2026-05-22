package problem_10

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ScrapeResult struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var globalResults = []ScrapeResult{
	{ID: "1", Title: "Product A"},
	{ID: "2", Title: "Product B"},
	{ID: "3", Title: "Product C"},
	{ID: "4", Title: "Product D"},
	{ID: "5", Title: "Product E"},
}

type PagedResponse struct {
	Data []ScrapeResult `json:"data"`
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 2
		}

		log.Printf("Fetching results for page %d with limit %d", page, limit)

		start := (page - 1) * limit

		end := start + limit + 1

		if start > len(globalResults) {
			start = len(globalResults)
		}
		// Missing check to ensure 'end' doesn't exceed len(globalResults) before slicing

		pagedData := globalResults[start:end]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(PagedResponse{Data: pagedData})
	}
}
