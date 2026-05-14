package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server is the driving adapter for HTTP via Gin.
type Server struct {
	engine *gin.Engine
}

// NewServer returns an HTTP adapter with Logger and Recovery middleware.
func NewServer() *Server {
	engine := gin.New()
	engine.SetTrustedProxies(nil) // trust no proxies by default
	engine.Use(gin.Logger(), gin.Recovery())
	return &Server{engine: engine}
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
