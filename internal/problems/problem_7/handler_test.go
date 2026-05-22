package problem_7

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	// A valid request that intentionally omits the optional limit_config
	payload := map[string]interface{}{
		"account_id": "acc_12345",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/bandwidth/limit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("HINT: check category B1")
			t.Fatalf("Handler panicked: %v", r)
		}
	}()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("HINT: check category B1")
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEdgeCase_WithLimits(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"account_id": "acc_12345",
		"limit_config": map[string]interface{}{
			"max_bytes": 1024000,
			"max_reqs":  100,
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/bandwidth/limit", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 when limits are provided, got %d", w.Code)
	}
}
