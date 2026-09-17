package headers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/webknife/webknife/pkg/headers"
)

func TestMiddleware_SetRequest(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Received", r.Header.Get("X-Debug"))
		w.WriteHeader(http.StatusOK)
	})
	handler := headers.Middleware([]headers.Operation{
		{Type: headers.Set, Name: "X-Debug", Value: "true"},
	}, headers.RequestHeaders)(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("X-Received") != "true" {
		t.Fatalf("expected X-Received=true, got %s", w.Header().Get("X-Received"))
	}
}

func TestMiddleware_AddRequest(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vals := r.Header.Values("X-Multi")
		w.Header().Set("X-Count", fmt.Sprintf("%d", len(vals)))
		w.WriteHeader(http.StatusOK)
	})
	handler := headers.Middleware([]headers.Operation{
		{Type: headers.Add, Name: "X-Multi", Value: "a"},
		{Type: headers.Add, Name: "X-Multi", Value: "b"},
	}, headers.RequestHeaders)(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("X-Count") != "2" {
		t.Fatalf("expected X-Count=2, got %s", w.Header().Get("X-Count"))
	}
}

func TestMiddleware_RemoveRequest(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("Authorization")
		w.Header().Set("X-Auth-Present", fmt.Sprintf("%t", val != ""))
		w.WriteHeader(http.StatusOK)
	})
	handler := headers.Middleware([]headers.Operation{
		{Type: headers.Remove, Name: "Authorization"},
	}, headers.RequestHeaders)(inner)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("X-Auth-Present") != "false" {
		t.Fatalf("expected auth removed, got X-Auth-Present=%s", w.Header().Get("X-Auth-Present"))
	}
}

func TestMiddleware_SetResponse(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := headers.Middleware([]headers.Operation{
		{Type: headers.Set, Name: "Cache-Control", Value: "no-store"},
	}, headers.ResponseHeaders)(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("expected no-store, got %s", w.Header().Get("Cache-Control"))
	}
}

func TestMiddleware_RemoveResponse(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "nginx/1.0")
		w.WriteHeader(http.StatusOK)
	})
	handler := headers.Middleware([]headers.Operation{
		{Type: headers.Remove, Name: "Server"},
	}, headers.ResponseHeaders)(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("Server") != "" {
		t.Fatalf("expected Server removed, got %s", w.Header().Get("Server"))
	}
}

func TestMiddleware_CaseInsensitive(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Debug", r.Header.Get("X-Debug"))
		w.WriteHeader(http.StatusOK)
	})
	handler := headers.Middleware([]headers.Operation{
		{Type: headers.Set, Name: "x-debug", Value: "true"},
	}, headers.RequestHeaders)(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("X-Debug") != "true" {
		t.Fatalf("expected X-Debug=true, got %s", w.Header().Get("X-Debug"))
	}
}

func TestBuildMiddlewares(t *testing.T) {
	cfg := headers.MultiConfig{
		SetRequest:     []string{"X-Debug: true"},
		RemoveResponse: []string{"Server"},
	}
	mws, err := headers.BuildMiddlewares(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(mws) != 2 {
		t.Fatalf("expected 2 middlewares, got %d", len(mws))
	}
}

func TestBuildMiddlewares_Empty(t *testing.T) {
	cfg := headers.MultiConfig{}
	mws, err := headers.BuildMiddlewares(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(mws) != 0 {
		t.Fatalf("expected 0 middlewares, got %d", len(mws))
	}
}
