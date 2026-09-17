package webknife

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/webknife/webknife/internal/events"
)

type ServerConfig struct {
	ListenAddr string
	TLSCert    string
	TLSKey     string
	Handler    http.Handler
}

func NewServer(cfg ServerConfig, pub events.Publisher) *http.Server {
	pub.Publish(events.ServerStarted{
		Addr: cfg.ListenAddr,
		TLS:  cfg.TLSCert != "",
	})

	return &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      cfg.Handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return context.Background()
		},
	}
}

func ListenAndServe(srv *http.Server, cfg ServerConfig, pub events.Publisher) error {
	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		return listenAndServeTLS(srv, cfg, pub)
	}
	return listenAndServe(srv, pub)
}

func listenAndServe(srv *http.Server, pub events.Publisher) error {
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	pub.Publish(events.ServerStopped{Addr: srv.Addr})
	return nil
}

func listenAndServeTLS(srv *http.Server, cfg ServerConfig, pub events.Publisher) error {
	err := srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
	if err != http.ErrServerClosed {
		return fmt.Errorf("TLS server error: %w", err)
	}
	pub.Publish(events.ServerStopped{Addr: srv.Addr})
	return nil
}

func NewTLSConfig() *tls.Config {
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
