// Package handler is the HTTP driving adapter (Gin).
package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fermin/gophercraft/internal/domain/port"
)

// Server is the driving adapter for HTTP via Gin.
type Server struct {
	engine          *gin.Engine
	metricsGatherer prometheus.Gatherer
}

// NewServer returns an HTTP adapter with structured request logging and Recovery middleware.
// When metricsGatherer is nil, GET /metrics is not registered.
func NewServer(logger port.Logger, metrics port.MetricsRecorder, metricsGatherer prometheus.Gatherer) (*Server, error) {
	engine := gin.New()
	if err := engine.SetTrustedProxies(nil); err != nil { // trust no proxies by default
		return nil, fmt.Errorf("set trusted proxies: %w", err)
	}
	engine.Use(metricsMiddleware(metrics), LoggerMiddleware(logger), gin.Recovery())
	return &Server{engine: engine, metricsGatherer: metricsGatherer}, nil
}

// Engine exposes the Gin engine for integration tests only.
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// RegisterRoutes mounts application routes on the Gin engine.
func (s *Server) RegisterRoutes() {
	if s.metricsGatherer != nil {
		s.engine.GET("/metrics", gin.WrapH(promhttp.HandlerFor(s.metricsGatherer, promhttp.HandlerOpts{})))
	}
	s.engine.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Run listens and serves HTTP on addr (e.g. ":3000").
func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func metricsMiddleware(rec port.MetricsRecorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		if route == "/metrics" {
			return
		}

		rec.RecordHTTPRequest(c.Request.Method, route, c.Writer.Status(), time.Since(start).Seconds())
	}
}
