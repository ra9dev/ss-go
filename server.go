package ssgo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

// ServerOpt to enhance Server
type ServerOpt func(*Server)

// ServerWithRoute mounts routes at a root Server.Handler
func ServerWithRoute(route ServerRoute) ServerOpt {
	return func(srv *Server) {
		srv.routes = append(srv.routes, route)
	}
}

type (
	// Server abstraction for an HTTP protocol
	Server struct {
		port uint
		srv  *http.Server

		routes []ServerRoute
	}

	// ServerRoute that can be registered via ServerWithRoute
	ServerRoute struct {
		Pattern string
		Handler http.Handler
	}
)

// NewServer constructor
func NewServer(port uint, opts ...ServerOpt) Server {
	srv := Server{
		port: port,
	}

	for _, opt := range opts {
		opt(&srv)
	}

	return srv
}

// Handler mounts multiple ServerRoute and handles HTTP requests
func (s Server) Handler() http.Handler {
	router := chi.NewRouter()

	for _, route := range s.routes {
		router.Mount(route.Pattern, route.Handler)
	}

	return router
}

// Run HTTP Server
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)

	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.Handler(),
	}

	ServerLogger.Printf("Application server LISTENS HTTP on %s...", s.srv.Addr)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server failed to serve: %w", err)
	}

	return nil
}

// Shutdown HTTP Server
func (s *Server) Shutdown() error {
	if err := s.srv.Shutdown(context.TODO()); err != nil {
		return fmt.Errorf("failed to shutdown: %w", err)
	}

	ServerLogger.Printf("Application server STOPPED listening HTTP on %s", s.srv.Addr)

	return nil
}
