package problem_8

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BUG HINT: D3

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	// Legitimate DELETE request to terminate a session
	req := httptest.NewRequest(http.MethodDelete, "/session/active?session_id=sess_123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Since the method check is backwards, it currently returns 405 Method Not Allowed
	if w.Code != http.StatusOK {
		t.Logf("HINT: check category D3")
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Code == http.StatusOK {
		var resp SessionResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.Status != "terminated" {
			t.Errorf("Expected status 'terminated', got '%s'", resp.Status)
		}
	}
}

func TestEdgeCase_MissingSessionID(t *testing.T) {
	handler := NewHandler()

	req := httptest.NewRequest(http.MethodDelete, "/session/active", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing session ID, got %d", w.Code)
	}
}
