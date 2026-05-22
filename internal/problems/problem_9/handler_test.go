package problem_9

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// BUG HINT: A2

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	req := httptest.NewRequest(http.MethodGet, "/proxy/regions", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Since the backend returns ErrNoRegionsAvailable, the API should return a 404 (or similar),
	// but currently it masks everything as a 500.
	if w.Code != http.StatusNotFound {
		t.Logf("HINT: check category A2")
		t.Errorf("Expected status 404 when no regions are available, got %d", w.Code)
	}
}

func TestEdgeCase_MethodNotAllowed(t *testing.T) {
	handler := NewHandler()

	req := httptest.NewRequest(http.MethodPost, "/proxy/regions", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for POST request, got %d", w.Code)
	}
}
