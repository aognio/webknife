package webknife

import (
	"fmt"
	"net/http"
	"strings"
)

type HeaderOperationType int

const (
	HeaderSet HeaderOperationType = iota
	HeaderAdd
	HeaderRemove
)

type HeaderOperation struct {
	Type  HeaderOperationType
	Name  string
	Value string
}

type HeaderDirection int

const (
	RequestHeaders HeaderDirection = iota
	ResponseHeaders
)

type HeaderMiddlewareConfig struct {
	Operations []HeaderOperation
	Direction  HeaderDirection
}

func HeaderMiddleware(cfg HeaderMiddlewareConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Direction == RequestHeaders {
				applyRequestHeaders(r, cfg.Operations)
			} else {
				rw := &responseHeaderWriter{
					ResponseWriter: w,
					operations:     cfg.Operations,
				}
				next.ServeHTTP(rw, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func applyRequestHeaders(r *http.Request, ops []HeaderOperation) {
	for _, op := range ops {
		name := http.CanonicalHeaderKey(op.Name)
		switch op.Type {
		case HeaderSet:
			r.Header.Set(name, op.Value)
		case HeaderAdd:
			r.Header.Add(name, op.Value)
		case HeaderRemove:
			r.Header.Del(name)
		}
	}
}

type responseHeaderWriter struct {
	http.ResponseWriter
	operations []HeaderOperation
	written    bool
}

func (rw *responseHeaderWriter) WriteHeader(code int) {
	if !rw.written {
		rw.written = true
		applyResponseHeaders(rw.ResponseWriter, rw.operations)
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseHeaderWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.written = true
		applyResponseHeaders(rw.ResponseWriter, rw.operations)
	}
	return rw.ResponseWriter.Write(b)
}

func applyResponseHeaders(w http.ResponseWriter, ops []HeaderOperation) {
	for _, op := range ops {
		name := http.CanonicalHeaderKey(op.Name)
		switch op.Type {
		case HeaderSet:
			w.Header().Set(name, op.Value)
		case HeaderAdd:
			w.Header().Add(name, op.Value)
		case HeaderRemove:
			w.Header().Del(name)
		}
	}
}

type MultiHeaderConfig struct {
	SetRequestHeaders     []string
	AddRequestHeaders     []string
	RemoveRequestHeaders  []string
	SetResponseHeaders    []string
	AddResponseHeaders    []string
	RemoveResponseHeaders []string
}

func BuildHeaderMiddlewares(cfg MultiHeaderConfig) ([]func(http.Handler) http.Handler, error) {
	var middlewares []func(http.Handler) http.Handler

	reqOps, err := buildOperations(
		cfg.SetRequestHeaders,
		cfg.AddRequestHeaders,
		cfg.RemoveRequestHeaders,
	)
	if err != nil {
		return nil, fmt.Errorf("request header error: %w", err)
	}
	if len(reqOps) > 0 {
		middlewares = append(middlewares, HeaderMiddleware(HeaderMiddlewareConfig{
			Operations: reqOps,
			Direction:  RequestHeaders,
		}))
	}

	respOps, err := buildOperations(
		cfg.SetResponseHeaders,
		cfg.AddResponseHeaders,
		cfg.RemoveResponseHeaders,
	)
	if err != nil {
		return nil, fmt.Errorf("response header error: %w", err)
	}
	if len(respOps) > 0 {
		middlewares = append(middlewares, HeaderMiddleware(HeaderMiddlewareConfig{
			Operations: respOps,
			Direction:  ResponseHeaders,
		}))
	}

	return middlewares, nil
}

func buildOperations(setFlags, addFlags, removeFlags []string) ([]HeaderOperation, error) {
	var ops []HeaderOperation

	for _, h := range setFlags {
		name, value, err := parseHeaderFlag(h)
		if err != nil {
			return nil, err
		}
		ops = append(ops, HeaderOperation{Type: HeaderSet, Name: name, Value: value})
	}

	for _, h := range addFlags {
		name, value, err := parseHeaderFlag(h)
		if err != nil {
			return nil, err
		}
		ops = append(ops, HeaderOperation{Type: HeaderAdd, Name: name, Value: value})
	}

	for _, h := range removeFlags {
		name := strings.TrimSpace(h)
		if name == "" {
			return nil, fmt.Errorf("empty header name in remove flag")
		}
		ops = append(ops, HeaderOperation{Type: HeaderRemove, Name: name})
	}

	return ops, nil
}

func parseHeaderFlag(h string) (string, string, error) {
	parts := strings.SplitN(h, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid header format: %s (expected 'Name: Value')", h)
	}
	name := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if name == "" {
		return "", "", fmt.Errorf("empty header name in: %s", h)
	}
	return name, value, nil
}
