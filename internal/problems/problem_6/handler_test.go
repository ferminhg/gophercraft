package problem_6

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	initialGoroutines := runtime.NumGoroutine()

	// Simulate clients that disconnect early
	for i := 0; i < 50; i++ {
		payload := map[string]interface{}{
			"url": "https://example.com/heavy",
		}
		body, _ := json.Marshal(payload)

		// Create a request with a context that cancels immediately
		ctx, cancel := context.WithCancel(context.Background())
		req := httptest.NewRequest(http.MethodPost, "/scrape/async", bytes.NewReader(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		// Cancel immediately to trigger the r.Context().Done() path
		cancel()

		handler.ServeHTTP(w, req)
	}

	// Wait a moment for all heavy scrapes to attempt to write to their channels
	time.Sleep(200 * time.Millisecond)

	finalGoroutines := runtime.NumGoroutine()

	// If the goroutines are leaked, final will be much higher than initial
	if finalGoroutines > initialGoroutines+20 {
		t.Logf("HINT: check category C2")
		t.Errorf("Goroutine leak detected: jumped from %d to %d", initialGoroutines, finalGoroutines)
	}
}

func TestEdgeCase_SuccessfulCompletion(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"url": "https://example.com/fast",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/scrape/async", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for successful completion, got %d", w.Code)
	}
}
