package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/infrastructure/handler"
	"github.com/fermin/gophercraft/internal/infrastructure/logger"
	"github.com/fermin/gophercraft/internal/infrastructure/metrics"
)

func TestStatusEndpoint(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	s, err := handler.NewServer(logger.NewDiscardingZerologLogger(), metrics.NoopRecorder{}, nil)
	require.NoError(t, err, "new server")
	s.RegisterRoutes()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	s.Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Status string `json:"status"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err, "decode response")
	assert.Equal(t, "ok", body.Status)
}
