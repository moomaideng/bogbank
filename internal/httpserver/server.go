// Package httpserver is the chi and Huma stack with liveness and readiness.
package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Pinger is a dependency readyz checks. *bun.DB implements it.
type Pinger interface {
	PingContext(ctx context.Context) error
}

type Server struct {
	api  huma.API
	http *http.Server
}

// New mounts health routes and Huma on one chi router. Health stays off the
// OpenAPI document so probes are not part of the versioned API.
func New(addr, serviceName string, readiness []Pinger) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	registerHealthRoutes(router, readiness)

	api := humachi.New(router, huma.DefaultConfig(serviceName, "1.0.0"))
	return &Server{
		api: api,
		http: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) HumaAPI() huma.API {
	return s.api
}

func (s *Server) Handler() http.Handler {
	return s.http.Handler
}

// Start listens until ctx is cancelled, then shuts down.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "address", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	slog.Info("shutting down")
	return s.http.Shutdown(shutdownCtx)
}
