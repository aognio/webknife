package respond_test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/webknife/webknife/pkg/respond"
)

func TestRespondHandler_200(t *testing.T) {
	handler, err := respond.New(respond.Config{
		StatusCode: 200,
		Body:       "ok",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Fatalf("expected 'ok', got %q", w.Body.String())
	}
}

func TestRespondHandler_503(t *testing.T) {
	handler, err := respond.New(respond.Config{
		StatusCode: 503,
		Body:       "Service unavailable",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 503 {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestRespondHandler_CustomHeaders(t *testing.T) {
	handler, err := respond.New(respond.Config{
		StatusCode:  429,
		Headers:     map[string]string{"Retry-After": "60"},
		Body:        `{"error":"rate limited"}`,
		ContentType: "application/json",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 429 {
		t.Fatalf("expected 429, got %d", w.Code)
	}
	if w.Header().Get("Retry-After") != "60" {
		t.Fatalf("expected Retry-After=60, got %s", w.Header().Get("Retry-After"))
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected application/json, got %s", w.Header().Get("Content-Type"))
	}
}

func TestRespondHandler_BodyFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "response.json")
	if err := os.WriteFile(tmpFile, []byte(`{"status":"ok"}`), 0644); err != nil {
		t.Fatal(err)
	}
	handler, err := respond.New(respond.Config{
		StatusCode:  200,
		BodyFile:    tmpFile,
		ContentType: "application/json",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Body.String() != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestRespondHandler_InvalidStatusCode(t *testing.T) {
	_, err := respond.New(respond.Config{StatusCode: 99})
	if err == nil {
		t.Fatal("expected error for invalid status code")
	}
	_, err = respond.New(respond.Config{StatusCode: 600})
	if err == nil {
		t.Fatal("expected error for invalid status code")
	}
}
