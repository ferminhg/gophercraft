// Package handler is the HTTP driving adapter (Gin).
package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fermin/gophercraft/internal/domain/port"
)

// Server is the driving adapter for HTTP via Gin.
type Server struct {
	engine *gin.Engine
}

// NewServer returns an HTTP adapter with structured request logging and Recovery middleware.
func NewServer(logger port.Logger) (*Server, error) {
	engine := gin.New()
	if err := engine.SetTrustedProxies(nil); err != nil { // trust no proxies by default
		return nil, fmt.Errorf("set trusted proxies: %w", err)
	}
	engine.Use(LoggerMiddleware(logger), gin.Recovery())
	return &Server{engine: engine}, nil
}

// Engine exposes the Gin engine for integration tests only.
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// RegisterRoutes mounts application routes on the Gin engine.
func (s *Server) RegisterRoutes() {
	s.engine.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Run listens and serves HTTP on addr (e.g. ":3000").
func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}
