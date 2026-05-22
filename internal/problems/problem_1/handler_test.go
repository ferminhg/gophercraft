package problem_1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BUG HINT: B2

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"target_url":   "https://example.com",
		"proxy_region": "eu-west-1",
		"js_rendering": true,
		// Notice: custom_headers is intentionally omitted here
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/proxy/session", bytes.NewReader(body))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HINT: check category B2")
			t.Fatalf("Handler panicked: %v", r)
		}
	}()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("HINT: check category B2")
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp SessionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.SessionID == "" {
		t.Error("Expected a session_id in response")
	}
}

func TestEdgeCase_EmptyBody(t *testing.T) {
	handler := NewHandler()

	req := httptest.NewRequest(http.MethodPost, "/proxy/session", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty body, got %d", w.Code)
	}
}
