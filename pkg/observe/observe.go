// Package observe provides HTTP request observation middleware that publishes
// lifecycle events through a Publisher port.
//
// This package depends only on the webknife ports package for the Publisher interface.
package observe

import (
	"net/http"
	"time"

	"github.com/aognio/webknife/pkg/webknife"
)

// Event types published by this middleware.

type RequestReceived struct {
	Method     string
	Host       string
	Path       string
	RemoteAddr string
	Timestamp  time.Time
}

type RequestCompleted struct {
	Method     string
	Host       string
	Path       string
	StatusCode int
	Size       int64
	Duration   time.Duration
	RemoteAddr string
	Timestamp  time.Time
}

type TLSConnectionAccepted struct {
	RemoteAddr string
	ServerName string
}

// Middleware returns middleware that publishes request lifecycle events.
func Middleware(pub webknife.Publisher) webknife.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			addr := r.RemoteAddr

			if r.TLS != nil {
				pub.Publish(TLSConnectionAccepted{
					RemoteAddr: r.RemoteAddr,
					ServerName: r.TLS.ServerName,
				})
			}
			pub.Publish(RequestReceived{
				Method:     r.Method,
				Host:       r.Host,
				Path:       r.URL.Path,
				RemoteAddr: r.RemoteAddr,
				Timestamp:  start,
			})

			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, r)

			pub.Publish(RequestCompleted{
				Method:     r.Method,
				Host:       r.Host,
				Path:       r.URL.Path,
				StatusCode: sw.status,
				Size:       sw.bytes,
				Duration:   time.Since(start),
				RemoteAddr: addr,
				Timestamp:  time.Now(),
			})
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}
