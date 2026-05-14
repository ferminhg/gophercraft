package handler

import "net/http"

// Server is a placeholder driving adapter for HTTP.
type Server struct{}

// NewServer returns a new HTTP adapter shell.
func NewServer() *Server {
	return &Server{}
}

// RegisterRoutes mounts routes on the given mux (placeholder).
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	_ = s
	_ = mux
}
