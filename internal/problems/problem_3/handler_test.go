package problem_3

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	// Create a dummy upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer upstream.Close()

	handler := NewHandler()

	// Track goroutines to detect the leak
	initialGoroutines := runtime.NumGoroutine()

	// Hit it a bunch of times
	for i := 0; i < 50; i++ {
		payload := map[string]interface{}{
			"target_url": upstream.URL,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/scrape/health", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	}

	runtime.GC()
	finalGoroutines := runtime.NumGoroutine()

	// If there's a huge delta, bodies are leaking keep-alive connections
	if finalGoroutines > initialGoroutines+20 {
		t.Logf("HINT: check category D1")
		t.Errorf("Potential leak detected: Goroutines jumped from %d to %d", initialGoroutines, finalGoroutines)
	}
}

func TestEdgeCase_BadURL(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"target_url": "http://this-url-definitely-does-not-exist.local",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/scrape/health", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("Expected status 502 for bad URL, got %d", w.Code)
	}
}
