package problem_2

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	// Simulate concurrent requests that will trigger a race condition panic
	var wg sync.WaitGroup
	numRequests := 100
	wg.Add(numRequests)

	errChan := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					t.Logf("HINT: check category C1")
					t.Errorf("Handler panicked: %v", r)
				}
			}()

			payload := map[string]interface{}{
				"region":  "us-east-1",
				"success": true,
			}
			body, _ := json.Marshal(payload)

			req := httptest.NewRequest(http.MethodPost, "/proxy/rotate", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		}()
	}

	wg.Wait()
	close(errChan)
}

func TestEdgeCase_MissingRegion(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"success": true,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/proxy/rotate", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing region, got %d", w.Code)
	}
}
