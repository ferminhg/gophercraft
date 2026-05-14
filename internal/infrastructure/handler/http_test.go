package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/fermin/gophercraft/internal/infrastructure/handler"
)

func TestStatusEndpoint(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	s, err := handler.NewServer()
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	s.RegisterRoutes()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	s.Engine().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("expected status ok, got %q", body.Status)
	}
}
