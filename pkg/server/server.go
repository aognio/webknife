// Package server provides HTTP server lifecycle management.
//
// It owns http.Server creation, timeouts, TLS startup, and graceful shutdown.
// This package has no dependencies on other Webknife feature packages.
package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Config holds configuration for the HTTP server.
type Config struct {
	Addr    string
	Handler http.Handler
}

// Server wraps http.Server with lifecycle management.
type Server struct {
	http *http.Server
}

// New creates a Server with the given config and sensible defaults.
func New(cfg Config) *Server {
	return &Server{
		http: &http.Server{
			Addr:         cfg.Addr,
			Handler:      cfg.Handler,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
			BaseContext: func(_ net.Listener) context.Context {
				return context.Background()
			},
		},
	}
}

// ListenAndServe starts the server. It returns nil on graceful shutdown.
func (s *Server) ListenAndServe() error {
	err := s.http.ListenAndServe()
	if err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// ListenAndServeTLS starts the server with TLS. It returns nil on graceful shutdown.
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	err := s.http.ListenAndServeTLS(certFile, keyFile)
	if err != http.ErrServerClosed {
		return fmt.Errorf("TLS server error: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

// DefaultTLSConfig returns a TLS configuration with sensible defaults.
func DefaultTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
