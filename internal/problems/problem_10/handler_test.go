package problem_10

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BUG HINT: E1

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	// Requesting the last page where `start + limit` aligns perfectly with the end of the slice,
	// but the `+ 1` bug pushes it out of bounds.
	req := httptest.NewRequest(http.MethodGet, "/scrape/results?page=3&limit=2", nil)
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HINT: check category E1")
			t.Fatalf("Handler panicked: %v", r)
		}
	}()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("HINT: check category E1")
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp PagedResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 item on the last page, got %d", len(resp.Data))
	}
}

func TestEdgeCase_PageOutOfBounds(t *testing.T) {
	handler := NewHandler()

	req := httptest.NewRequest(http.MethodGet, "/scrape/results?page=99&limit=2", nil)
	w := httptest.NewRecorder()

	// Should not panic, should return empty array
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for out of bounds page, got %d", w.Code)
	}
}
