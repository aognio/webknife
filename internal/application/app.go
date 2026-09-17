// Package application provides the composition layer that wires handlers,
// middleware, logging, and server lifecycle together.
//
// It imports feature packages and adapters, building complete handler
// pipelines from parsed CLI configuration. This is where the dependency
// arrows converge — feature packages never import each other.
package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aognio/webknife/internal/adapters/logging"
	"github.com/aognio/webknife/internal/build"
	"github.com/aognio/webknife/internal/events"
	"github.com/aognio/webknife/pkg/auth"
	"github.com/aognio/webknife/pkg/echo"
	"github.com/aognio/webknife/pkg/headers"
	"github.com/aognio/webknife/pkg/observe"
	"github.com/aognio/webknife/pkg/proxy"
	"github.com/aognio/webknife/pkg/redirect"
	"github.com/aognio/webknife/pkg/respond"
	"github.com/aognio/webknife/pkg/server"
	"github.com/aognio/webknife/pkg/static"
	"github.com/aognio/webknife/pkg/webknife"
)

// Config is the fully-parsed application configuration.
// The CLI adapter produces this; the application layer consumes it.
type Config struct {
	Command string

	ListenAddr string
	Root       string
	Upstream   string
	Auth       string
	TLSCert    string
	TLSKey     string
	LogFormat  string
	Verbose    bool

	EchoMaxBody int64

	RespondStatusCode  int
	RespondBody        string
	RespondBodyFile    string
	RespondContentType string

	RedirectTarget       string
	RedirectStatusCode   int
	RedirectPreservePath bool

	SetRequestHeaders     []string
	AddRequestHeaders     []string
	RemoveRequestHeaders  []string
	SetResponseHeaders    []string
	AddResponseHeaders    []string
	RemoveResponseHeaders []string
}

// Run builds the handler for the given command, wires middleware, and starts the server.
func Run(cfg *Config) error {
	bus := events.NewEventBus()
	var logHandler logging.Handler
	if cfg.LogFormat == "json" {
		logHandler = logging.NewJSONHandler(os.Stdout)
	} else {
		logHandler = logging.NewTextHandler(os.Stdout)
	}

	// Subscribe to observe events and translate them for the log adapter.
	bus.Subscribe(func(event events.Event) {
		switch e := event.(type) {
		case observe.RequestCompleted:
			logHandler.Handle(logging.RequestCompleted{
				Method:     e.Method,
				Path:       e.Path,
				StatusCode: e.StatusCode,
				Size:       e.Size,
				Duration:   e.Duration,
				RemoteAddr: e.RemoteAddr,
				Timestamp:  e.Timestamp,
			})
		case observe.TLSConnectionAccepted:
			// TLS events are logged at server start, not per-request.
		}
	})

	pub := &eventPublisher{bus: bus}

	handler, err := buildHandler(cfg, pub)
	if err != nil {
		return err
	}

	addr := cfg.ListenAddr
	if addr == "" {
		addr = ":8080"
	}

	fmt.Fprintf(os.Stderr, "webknife %s\n", build.Version)

	tlsMode := cfg.TLSCert != "" && cfg.TLSKey != ""
	if tlsMode {
		logHandler.Handle(logging.ServerStarted{Addr: addr, TLS: true})
	} else {
		logHandler.Handle(logging.ServerStarted{Addr: addr, TLS: false})
	}

	srv := server.New(server.Config{
		Addr:    addr,
		Handler: handler,
	})

	// Graceful shutdown on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if tlsMode {
			errCh <- srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
		} else {
			errCh <- srv.ListenAndServe()
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		logHandler.Handle(logging.ServerStopped{Addr: addr})
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		logHandler.Handle(logging.ServerStopped{Addr: addr})
		return err
	}
}

func buildHandler(cfg *Config, pub webknife.Publisher) (http.Handler, error) {
	var handler http.Handler

	switch cfg.Command {
	case "serve":
		h, err := static.New(static.Config{Root: cfg.Root})
		if err != nil {
			return nil, fmt.Errorf("invalid root directory: %w", err)
		}
		handler = h
	case "proxy":
		h, err := proxy.New(proxy.Config{Upstream: cfg.Upstream})
		if err != nil {
			return nil, fmt.Errorf("invalid upstream URL: %w", err)
		}
		handler = h
	case "echo":
		handler = echo.New(echo.Config{MaxBodySize: cfg.EchoMaxBody})
	case "respond":
		h, err := respond.New(respond.Config{
			StatusCode:  cfg.RespondStatusCode,
			Body:        cfg.RespondBody,
			BodyFile:    cfg.RespondBodyFile,
			ContentType: cfg.RespondContentType,
		})
		if err != nil {
			return nil, err
		}
		handler = h
	case "redirect":
		h, err := redirect.New(redirect.Config{
			StatusCode:   cfg.RedirectStatusCode,
			Target:       cfg.RedirectTarget,
			PreservePath: cfg.RedirectPreservePath,
		})
		if err != nil {
			return nil, err
		}
		handler = h
	default:
		return nil, fmt.Errorf("unknown command: %s", cfg.Command)
	}

	// Header manipulation middleware.
	headerMws, err := headers.BuildMiddlewares(headers.MultiConfig{
		SetRequest:     cfg.SetRequestHeaders,
		AddRequest:     cfg.AddRequestHeaders,
		RemoveRequest:  cfg.RemoveRequestHeaders,
		SetResponse:    cfg.SetResponseHeaders,
		AddResponse:    cfg.AddResponseHeaders,
		RemoveResponse: cfg.RemoveResponseHeaders,
	})
	if err != nil {
		return nil, err
	}
	for i := len(headerMws) - 1; i >= 0; i-- {
		handler = headerMws[i](handler)
	}

	// Basic auth middleware.
	if cfg.Auth != "" {
		user, pass := parseAuth(cfg.Auth)
		handler = auth.Basic(auth.Config{
			Username: user,
			Password: pass,
		}, handler)
	}

	// Observation middleware (outermost).
	handler = observe.Middleware(pub)(handler)

	return handler, nil
}

func parseAuth(auth string) (string, string) {
	for i := 0; i < len(auth); i++ {
		if auth[i] == ':' {
			return auth[:i], auth[i+1:]
		}
	}
	return auth, ""
}

// eventPublisher adapts events.EventBus to the webknife.Publisher port.
type eventPublisher struct {
	bus *events.EventBus
}

func (p *eventPublisher) Publish(event webknife.Event) {
	// Translate observe events to internal events and publish.
	switch e := event.(type) {
	case observe.RequestReceived:
		p.bus.Publish(events.RequestReceived{
			Method:     e.Method,
			Host:       e.Host,
			Path:       e.Path,
			RemoteAddr: e.RemoteAddr,
			Timestamp:  e.Timestamp,
		})
	case observe.RequestCompleted:
		p.bus.Publish(events.RequestCompleted{
			Method:     e.Method,
			Host:       e.Host,
			Path:       e.Path,
			StatusCode: e.StatusCode,
			Size:       e.Size,
			Duration:   e.Duration,
			RemoteAddr: e.RemoteAddr,
			Timestamp:  e.Timestamp,
		})
	case observe.TLSConnectionAccepted:
		p.bus.Publish(events.TLSConnectionAccepted{
			RemoteAddr: e.RemoteAddr,
			ServerName: e.ServerName,
		})
	}
}
