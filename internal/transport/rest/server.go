// Package rest implements setup, teardown and handlers for the REST API
package rest

import (
	"context"
	"net/http"
	"time"
)

// Server represents API server of application
type Server struct {
	httpServer *http.Server
}

// Run creates configuration for the server and starts it
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 Mb
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown stops the Server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
