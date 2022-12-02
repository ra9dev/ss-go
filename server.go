package ssgo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const defaultServerPort = 8080

type ServerOpt func(*Server)

func ServerWithRoute(route ServerRoute) ServerOpt {
	return func(srv *Server) {
		srv.routes = append(srv.routes, route)
	}
}

type (
	Server struct {
		port uint
		srv  *http.Server

		routes []ServerRoute
	}

	ServerRoute struct {
		Pattern string
		Handler http.Handler
	}
)

func NewServer(port uint, opts ...ServerOpt) Server {
	srv := Server{
		port: port,
	}

	for _, opt := range opts {
		opt(&srv)
	}

	return srv
}

func (s Server) Handler() http.Handler {
	router := chi.NewRouter()

	for _, route := range s.routes {
		router.Mount(route.Pattern, route.Handler)
	}

	return router
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)

	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.Handler(),
	}

	log.Printf("Application server LISTENS HTTP on %s...", s.srv.Addr)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server failed to serve: %w", err)
	}

	return nil
}

func (s *Server) Shutdown() error {
	if err := s.srv.Shutdown(context.TODO()); err != nil {
		return err
	}

	log.Printf("Application server STOPPED listening HTTP on %s", s.srv.Addr)

	return nil
}
