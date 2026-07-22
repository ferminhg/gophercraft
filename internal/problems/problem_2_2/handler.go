package problem_2_2

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type JobStats struct {
	SuccessCount int
	FailureCount int
}

type MutexJobStats struct {
	locker sync.Mutex
	stats map[string]JobStats
}

func NewMutexJobStats() *MutexJobStats {
	return &MutexJobStats{
		stats: make(map[string]JobStats),
	}
}

// Global state for scrape job statistics
var jobStats = NewMutexJobStats()

type JobCompleteRequest struct {
	Domain  string `json:"domain"`
	Success bool   `json:"success"`
}

func (m *MutexJobStats) RecordJobCompletion(domain string, success bool) JobStats {
	m.locker.Lock()
	defer m.locker.Unlock()
	stats := m.stats[domain]
	if success {
		stats.SuccessCount++
	} else {
		stats.FailureCount++
	}
	m.stats[domain] = stats
	return stats
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req JobCompleteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if req.Domain == "" {
			http.Error(w, "domain is required", http.StatusBadRequest)
			return
		}

		log.Printf("📊 Recording job completion for domain %s", req.Domain)

		stats := jobStats.RecordJobCompletion(req.Domain, req.Success)
		log.Printf("📊 Recorded job completion for domain %s: %+v", req.Domain, stats)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "recorded"}`))
	}
}
