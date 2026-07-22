package problem_4_2

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"domain":        "example.com",
		"authorization": "sometoken_without_bearer_prefix",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/extract/schedule", bytes.NewReader(body))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HINT: check category A1")
			t.Fatalf("Handler panicked: %v", r)
		}
	}()

	handler.ServeHTTP(w, req)

	// Since the input is bad, we'd actually expect a 400.
	// But right now it panics (500). The fix should handle the error and return a 400.
	if w.Code != http.StatusBadRequest {
		t.Logf("HINT: check category A1")
		t.Errorf("Expected status 400 for malformed authorization, got %d", w.Code)
	}
}

func TestEdgeCase_EmptyAuthorization(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"domain":        "example.com",
		"authorization": "",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/extract/schedule", bytes.NewReader(body))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HINT: check category A1")
			t.Fatalf("Handler panicked: %v", r)
		}
	}()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty authorization, got %d", w.Code)
	}
}

func TestEdgeCase_ValidAuthorization(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"domain":        "example.com",
		"authorization": "Bearer abc123token",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/extract/schedule", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for valid authorization, got %d", w.Code)
	}
}
