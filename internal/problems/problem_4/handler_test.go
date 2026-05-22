package problem_4

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
		"url":      "https://example.com",
		"selector": "invalid_selector_without_dot",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/extract/rule", bytes.NewReader(body))
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
		t.Errorf("Expected status 400 for bad selector, got %d", w.Code)
	}
}

func TestEdgeCase_ValidSelector(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"url":      "https://example.com",
		"selector": "div.product-card",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/extract/rule", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for valid selector, got %d", w.Code)
	}
}
