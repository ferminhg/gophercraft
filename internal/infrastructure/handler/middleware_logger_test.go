package handler_test

import (
	"bytes"
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

func TestLoggerMiddleware_LogsRequestFields(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	log := logger.NewZerologLoggerWithWriter(&buf, "info", false, nil)

	s, err := handler.NewServer(log, metrics.NoopRecorder{}, nil)
	require.NoError(t, err)
	s.RegisterRoutes()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	s.Engine().ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var line map[string]any
	dec := json.NewDecoder(bytes.NewReader(bytes.TrimSpace(buf.Bytes())))
	require.NoError(t, dec.Decode(&line))

	assert.Equal(t, "request completed", line["message"])
	assert.Equal(t, "GET", line["method"])
	assert.Equal(t, "/status", line["path"])
	assert.Equal(t, float64(200), line["status"]) // JSON numbers decode as float64
	assert.Contains(t, line, "latency_ms")
}
