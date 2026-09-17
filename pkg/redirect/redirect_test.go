package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/webknife/webknife/pkg/redirect"
)

func TestRedirectHandler_302(t *testing.T) {
	handler, err := redirect.New(redirect.Config{
		StatusCode: 302,
		Target:     "https://example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 302 {
		t.Fatalf("expected 302, got %d", w.Code)
	}
	if w.Header().Get("Location") != "https://example.com" {
		t.Fatalf("expected Location=https://example.com, got %s", w.Header().Get("Location"))
	}
}

func TestRedirectHandler_301(t *testing.T) {
	handler, err := redirect.New(redirect.Config{
		StatusCode: 301,
		Target:     "https://example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/old", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 301 {
		t.Fatalf("expected 301, got %d", w.Code)
	}
}

func TestRedirectHandler_PreservePath(t *testing.T) {
	handler, err := redirect.New(redirect.Config{
		StatusCode:   302,
		Target:       "https://example.com",
		PreservePath: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/foo?a=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 302 {
		t.Fatalf("expected 302, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if loc != "https://example.com/foo?a=1" {
		t.Fatalf("expected https://example.com/foo?a=1, got %s", loc)
	}
}

func TestRedirectHandler_InvalidStatus(t *testing.T) {
	_, err := redirect.New(redirect.Config{
		StatusCode: 200,
		Target:     "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error for invalid redirect status")
	}
	_, err = redirect.New(redirect.Config{
		StatusCode: 404,
		Target:     "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error for invalid redirect status")
	}
}
