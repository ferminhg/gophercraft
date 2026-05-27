package problem_2

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type ProxyStatus struct {
	ActiveCount int
	ErrorCount  int
}

// Global state for proxy statistics
// Old implementation
// var regionStats = make(map[string]ProxyStatus)
// var regionStatsMutex sync.Mutex

type MutexRegionStats struct {
	locker sync.Mutex
	stats map[string]ProxyStatus
}

func NewMutexRegionStats() *MutexRegionStats {
	return &MutexRegionStats{
		stats: make(map[string]ProxyStatus),
	}
}

func (mr *MutexRegionStats) RecordRotation(region string, success bool) {
	mr.locker.Lock()
	defer mr.locker.Unlock()
	stats := mr.stats[region]
	if success {
		stats.ActiveCount++
	} else {
		stats.ErrorCount++
	}
	mr.stats[region] = stats
	log.Printf("\t Updated stats for region %s: %+v", region, stats)
}

var rwRegionStats = NewMutexRegionStats()

type RotateRequest struct {
	Region  string `json:"region"`
	Success bool   `json:"success"`
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RotateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if req.Region == "" {
			http.Error(w, "region is required", http.StatusBadRequest)
			return
		}

		log.Printf("🧵 Recording proxy rotation for region %s", req.Region)

		// Old implementation
		// regionStatsMutex.Lock()
		// stats := regionStats[req.Region]
		// if req.Success {
		// 	stats.ActiveCount++
		// } else {
		// 	stats.ErrorCount++
		// }
		// regionStats[req.Region] = stats
		// regionStatsMutex.Unlock()

		rwRegionStats.RecordRotation(req.Region, req.Success)
		

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "recorded"}`))
	}
}
