// Package logging provides adapters for observing request lifecycle events.
//
// It subscribes to events published by the observe middleware and formats
// them as text or JSON logs. This package depends only on pkg/webknife
// (for the Middleware type) and the standard library.
package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

// Event is a minimal interface for events that carry loggable data.
// The observe package publishes concrete types that satisfy this.
type Event interface{}

// Handler processes log events.
type Handler interface {
	Handle(event Event)
}

// TextHandler logs events as human-readable text.
type TextHandler struct {
	output io.Writer
}

func NewTextHandler(w io.Writer) *TextHandler {
	return &TextHandler{output: w}
}

// RequestCompleted is the data shape we expect to log.
type RequestCompleted struct {
	Method     string
	Path       string
	StatusCode int
	Size       int64
	Duration   time.Duration
	RemoteAddr string
	Timestamp  time.Time
}

type ServerStarted struct {
	Addr string
	TLS  bool
}

type ServerStopped struct {
	Addr string
}

type AuthFailed struct {
	User string
	Path string
}

func (h *TextHandler) Handle(event Event) {
	switch e := event.(type) {
	case RequestCompleted:
		h.logRequest(e)
	case ServerStarted:
		h.logServerStarted(e)
	case ServerStopped:
		h.logServerStopped(e)
	case AuthFailed:
		h.logAuthFailed(e)
	}
}

func (h *TextHandler) logRequest(e RequestCompleted) {
	sizeStr := formatBytes(e.Size)
	ts := e.Timestamp.Format(time.RFC3339)
	log.SetOutput(h.output)
	log.Printf("%s %s %s %d %s %s %s",
		ts, e.Method, e.Path, e.StatusCode, sizeStr, e.Duration.Round(time.Microsecond), e.RemoteAddr)
}

func (h *TextHandler) logServerStarted(e ServerStarted) {
	ts := time.Now().Format(time.RFC3339)
	mode := "HTTP"
	if e.TLS {
		mode = "HTTPS"
	}
	log.SetOutput(h.output)
	log.Printf("%s %s server listening on %s", ts, mode, e.Addr)
}

func (h *TextHandler) logServerStopped(e ServerStopped) {
	ts := time.Now().Format(time.RFC3339)
	log.SetOutput(h.output)
	log.Printf("%s server stopped on %s", ts, e.Addr)
}

func (h *TextHandler) logAuthFailed(e AuthFailed) {
	ts := time.Now().Format(time.RFC3339)
	log.SetOutput(h.output)
	log.Printf("%s authentication failed user=%s path=%s", ts, e.User, e.Path)
}

// JSONHandler logs events as structured JSON.
type JSONHandler struct {
	output io.Writer
}

func NewJSONHandler(w io.Writer) *JSONHandler {
	return &JSONHandler{output: w}
}

func (h *JSONHandler) Handle(event Event) {
	switch e := event.(type) {
	case RequestCompleted:
		h.logRequest(e)
	case ServerStarted:
		h.logServerStarted(e)
	case ServerStopped:
		h.logServerStopped(e)
	case AuthFailed:
		h.logAuthFailed(e)
	}
}

type jsonEntry struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Size      string `json:"size"`
	Duration  string `json:"duration"`
	Remote    string `json:"remote"`
}

func (h *JSONHandler) logRequest(e RequestCompleted) {
	entry := jsonEntry{
		Timestamp: e.Timestamp.Format(time.RFC3339),
		Method:    e.Method,
		Path:      e.Path,
		Status:    e.StatusCode,
		Size:      formatBytes(e.Size),
		Duration:  e.Duration.Round(time.Microsecond).String(),
		Remote:    e.RemoteAddr,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		log.SetOutput(h.output)
		log.Printf("error marshaling log entry: %v", err)
		return
	}
	log.SetOutput(h.output)
	log.Println(string(data))
}

func (h *JSONHandler) logServerStarted(e ServerStarted) {
	ts := time.Now().Format(time.RFC3339)
	mode := "HTTP"
	if e.TLS {
		mode = "HTTPS"
	}
	log.SetOutput(h.output)
	log.Printf("%s %s server listening on %s", ts, mode, e.Addr)
}

func (h *JSONHandler) logServerStopped(e ServerStopped) {
	ts := time.Now().Format(time.RFC3339)
	log.SetOutput(h.output)
	log.Printf("%s server stopped on %s", ts, e.Addr)
}

func (h *JSONHandler) logAuthFailed(e AuthFailed) {
	ts := time.Now().Format(time.RFC3339)
	log.SetOutput(h.output)
	log.Printf("%s authentication failed user=%s path=%s", ts, e.User, e.Path)
}

func formatBytes(b int64) string {
	switch {
	case b >= 1<<20:
		return fmt.Sprintf("%.1fMB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1fKB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}
