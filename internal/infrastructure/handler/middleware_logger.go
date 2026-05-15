package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fermin/gophercraft/internal/domain/port"
)

// LoggerMiddleware emits a structured log line after each request.
func LoggerMiddleware(logger port.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)

		logger.Info("request completed",
			"method", method,
			"path", path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
		)
	}
}
