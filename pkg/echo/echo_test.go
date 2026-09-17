package echo_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aognio/webknife/pkg/echo"
)

func TestEchoHandler_GET(t *testing.T) {
	handler := echo.New(echo.Config{})
	req := httptest.NewRequest("GET", "/test?foo=bar", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
	var inspection echo.Inspection
	if err := json.Unmarshal(w.Body.Bytes(), &inspection); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if inspection.Method != "GET" {
		t.Fatalf("expected GET, got %s", inspection.Method)
	}
	if inspection.Path != "/test" {
		t.Fatalf("expected /test, got %s", inspection.Path)
	}
	if inspection.Query["foo"][0] != "bar" {
		t.Fatalf("expected foo=bar, got %v", inspection.Query["foo"])
	}
}

func TestEchoHandler_POST(t *testing.T) {
	handler := echo.New(echo.Config{})
	body := `{"message":"hello"}`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection echo.Inspection
	if err := json.Unmarshal(w.Body.Bytes(), &inspection); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if inspection.Method != "POST" {
		t.Fatalf("expected POST, got %s", inspection.Method)
	}
	if inspection.Body != body {
		t.Fatalf("expected body %q, got %q", body, inspection.Body)
	}
	if inspection.ContentLength != int64(len(body)) {
		t.Fatalf("expected content_length %d, got %d", len(body), inspection.ContentLength)
	}
}

func TestEchoHandler_QueryParams(t *testing.T) {
	handler := echo.New(echo.Config{})
	req := httptest.NewRequest("GET", "/test?foo=bar&foo=baz&key=val", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection echo.Inspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if len(inspection.Query["foo"]) != 2 {
		t.Fatalf("expected 2 foo values, got %d", len(inspection.Query["foo"]))
	}
	if inspection.Query["foo"][0] != "bar" || inspection.Query["foo"][1] != "baz" {
		t.Fatalf("unexpected foo values: %v", inspection.Query["foo"])
	}
}

func TestEchoHandler_Headers(t *testing.T) {
	handler := echo.New(echo.Config{})
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Custom", "test-value")
	req.Header.Set("X-Debug", "true")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection echo.Inspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if inspection.Headers["X-Custom"][0] != "test-value" {
		t.Fatalf("expected X-Custom=test-value, got %v", inspection.Headers["X-Custom"])
	}
}

func TestEchoHandler_Cookies(t *testing.T) {
	handler := echo.New(echo.Config{})
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection echo.Inspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if len(inspection.Cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(inspection.Cookies))
	}
	if inspection.Cookies[0].Name != "session" {
		t.Fatalf("expected session cookie, got %s", inspection.Cookies[0].Name)
	}
}

func TestEchoHandler_BodyTruncation(t *testing.T) {
	handler := echo.New(echo.Config{MaxBodySize: 10})
	body := strings.Repeat("a", 20)
	req := httptest.NewRequest("POST", "/test", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection echo.Inspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if !inspection.BodyTruncated {
		t.Fatal("expected body_truncated=true")
	}
	if len(inspection.Body) != 10 {
		t.Fatalf("expected 10 body bytes, got %d", len(inspection.Body))
	}
}

func TestEchoHandler_BinaryBody(t *testing.T) {
	handler := echo.New(echo.Config{})
	body := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection echo.Inspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if !inspection.BodyBinary {
		t.Fatal("expected body_binary=true")
	}
}

func TestEchoHandler_RemoteAddr(t *testing.T) {
	handler := echo.New(echo.Config{})
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection echo.Inspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if inspection.RemoteAddr != "192.168.1.100:12345" {
		t.Fatalf("expected 192.168.1.100:12345, got %s", inspection.RemoteAddr)
	}
}
