package events

import (
	"net/http"
	"time"
)

type Event interface{}

type ServerStarted struct {
	Addr string
	TLS  bool
}

type ServerStopped struct {
	Addr string
}

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

type RequestFailed struct {
	Method     string
	Host       string
	Path       string
	Error      error
	Duration   time.Duration
	RemoteAddr string
	Timestamp  time.Time
}

type ProxyRequestStarted struct {
	Method   string
	Path     string
	Upstream string
}

type ProxyRequestCompleted struct {
	Method     string
	Path       string
	Upstream   string
	StatusCode int
	Duration   time.Duration
}

type AuthenticationSucceeded struct {
	Method string
	Path   string
	User   string
}

type AuthenticationFailed struct {
	Method string
	Path   string
	User   string
}

type TLSConnectionAccepted struct {
	RemoteAddr string
	ServerName string
}

type RequestInspected struct {
	Method string
	Path   string
}

type ResponseGenerated struct {
	Method     string
	Path       string
	StatusCode int
}

type RedirectGenerated struct {
	Method     string
	Path       string
	Target     string
	StatusCode int
}

type RequestHeadersModified struct {
	Method string
	Path   string
	Count  int
}

type ResponseHeadersModified struct {
	Method string
	Path   string
	Count  int
}

type Publisher interface {
	Publish(event Event)
}

type Middleware func(http.Handler) http.Handler

func RequestMiddleware(pub Publisher) Middleware {
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
