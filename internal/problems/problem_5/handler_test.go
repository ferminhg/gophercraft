package problem_5

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Returns200_WhenFixed(t *testing.T) {
	handler := NewHandler()

	// A target URL with spaces and special characters that must be encoded
	req := httptest.NewRequest(http.MethodGet, "/proxy/forward?target_url=https://example.com/search?q=hello world & stuff", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("HINT: check category E3")
		t.Errorf("Expected status 200, got %d. The URL was likely malformed.", w.Code)
	}

	var resp ForwardResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if !strings.Contains(resp.UpstreamURL, "hello+world+%26+stuff") && !strings.Contains(resp.UpstreamURL, "hello%20world%20%26%20stuff") {
		t.Errorf("Upstream URL does not appear to be properly escaped: %s", resp.UpstreamURL)
	}
}

func TestEdgeCase_MissingTarget(t *testing.T) {
	handler := NewHandler()

	req := httptest.NewRequest(http.MethodGet, "/proxy/forward", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing target_url, got %d", w.Code)
	}
}
