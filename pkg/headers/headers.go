// Package headers provides HTTP header manipulation middleware.
//
// This package has no dependencies on other Webknife feature packages.
package headers

import (
	"fmt"
	"net/http"
	"strings"
)

// OperationType describes how a header operation modifies headers.
type OperationType int

const (
	Set    OperationType = iota // Set/replace a header value
	Add                         // Add a header value (appends)
	Remove                      // Remove a header entirely
)

// Operation describes a single header transformation.
type Operation struct {
	Type  OperationType
	Name  string
	Value string
}

// Direction specifies whether operations apply to request or response headers.
type Direction int

const (
	RequestHeaders Direction = iota
	ResponseHeaders
)

// Middleware returns a middleware that applies header operations.
func Middleware(ops []Operation, dir Direction) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if dir == RequestHeaders {
				applyRequest(r, ops)
				next.ServeHTTP(w, r)
				return
			}
			rw := &responseWriter{ResponseWriter: w, operations: ops}
			next.ServeHTTP(rw, r)
		})
	}
}

func applyRequest(r *http.Request, ops []Operation) {
	for _, op := range ops {
		name := http.CanonicalHeaderKey(op.Name)
		switch op.Type {
		case Set:
			r.Header.Set(name, op.Value)
		case Add:
			r.Header.Add(name, op.Value)
		case Remove:
			r.Header.Del(name)
		}
	}
}

type responseWriter struct {
	http.ResponseWriter
	operations []Operation
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.written = true
		applyResponse(rw.ResponseWriter, rw.operations)
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.written = true
		applyResponse(rw.ResponseWriter, rw.operations)
	}
	return rw.ResponseWriter.Write(b)
}

func applyResponse(w http.ResponseWriter, ops []Operation) {
	for _, op := range ops {
		name := http.CanonicalHeaderKey(op.Name)
		switch op.Type {
		case Set:
			w.Header().Set(name, op.Value)
		case Add:
			w.Header().Add(name, op.Value)
		case Remove:
			w.Header().Del(name)
		}
	}
}

// MultiConfig holds header operations for both request and response directions.
type MultiConfig struct {
	SetRequest     []string
	AddRequest     []string
	RemoveRequest  []string
	SetResponse    []string
	AddResponse    []string
	RemoveResponse []string
}

// BuildMiddlewares creates middleware functions from a multi-config.
// Returns nil (no error) if no operations are configured.
func BuildMiddlewares(cfg MultiConfig) ([]func(http.Handler) http.Handler, error) {
	var mws []func(http.Handler) http.Handler

	reqOps, err := buildOps(cfg.SetRequest, cfg.AddRequest, cfg.RemoveRequest)
	if err != nil {
		return nil, fmt.Errorf("request header error: %w", err)
	}
	if len(reqOps) > 0 {
		mws = append(mws, Middleware(reqOps, RequestHeaders))
	}

	respOps, err := buildOps(cfg.SetResponse, cfg.AddResponse, cfg.RemoveResponse)
	if err != nil {
		return nil, fmt.Errorf("response header error: %w", err)
	}
	if len(respOps) > 0 {
		mws = append(mws, Middleware(respOps, ResponseHeaders))
	}

	return mws, nil
}

func buildOps(setFlags, addFlags, removeFlags []string) ([]Operation, error) {
	var ops []Operation
	for _, h := range setFlags {
		name, value, err := parseFlag(h)
		if err != nil {
			return nil, err
		}
		ops = append(ops, Operation{Type: Set, Name: name, Value: value})
	}
	for _, h := range addFlags {
		name, value, err := parseFlag(h)
		if err != nil {
			return nil, err
		}
		ops = append(ops, Operation{Type: Add, Name: name, Value: value})
	}
	for _, h := range removeFlags {
		name := strings.TrimSpace(h)
		if name == "" {
			return nil, fmt.Errorf("empty header name in remove flag")
		}
		ops = append(ops, Operation{Type: Remove, Name: name})
	}
	return ops, nil
}

func parseFlag(h string) (string, string, error) {
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
