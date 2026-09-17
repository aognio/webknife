package webknife_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/webknife/webknife/pkg/webknife"
)

func TestEchoHandler_GET(t *testing.T) {
	handler := webknife.NewEchoHandler(webknife.EchoConfig{})
	req := httptest.NewRequest("GET", "/test?foo=bar", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
	var inspection webknife.RequestInspection
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
	handler := webknife.NewEchoHandler(webknife.EchoConfig{})
	body := `{"message":"hello"}`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection webknife.RequestInspection
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
	handler := webknife.NewEchoHandler(webknife.EchoConfig{})
	req := httptest.NewRequest("GET", "/test?foo=bar&foo=baz&key=val", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection webknife.RequestInspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if len(inspection.Query["foo"]) != 2 {
		t.Fatalf("expected 2 foo values, got %d", len(inspection.Query["foo"]))
	}
	if inspection.Query["foo"][0] != "bar" || inspection.Query["foo"][1] != "baz" {
		t.Fatalf("unexpected foo values: %v", inspection.Query["foo"])
	}
}

func TestEchoHandler_Headers(t *testing.T) {
	handler := webknife.NewEchoHandler(webknife.EchoConfig{})
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Custom", "test-value")
	req.Header.Set("X-Debug", "true")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection webknife.RequestInspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if inspection.Headers["X-Custom"][0] != "test-value" {
		t.Fatalf("expected X-Custom=test-value, got %v", inspection.Headers["X-Custom"])
	}
}

func TestEchoHandler_Cookies(t *testing.T) {
	handler := webknife.NewEchoHandler(webknife.EchoConfig{})
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection webknife.RequestInspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if len(inspection.Cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(inspection.Cookies))
	}
	if inspection.Cookies[0].Name != "session" {
		t.Fatalf("expected session cookie, got %s", inspection.Cookies[0].Name)
	}
}

func TestEchoHandler_BodyTruncation(t *testing.T) {
	handler := webknife.NewEchoHandler(webknife.EchoConfig{MaxBodySize: 10})
	body := strings.Repeat("a", 20)
	req := httptest.NewRequest("POST", "/test", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection webknife.RequestInspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if !inspection.BodyTruncated {
		t.Fatal("expected body_truncated=true")
	}
	if len(inspection.Body) != 10 {
		t.Fatalf("expected 10 body bytes, got %d", len(inspection.Body))
	}
}

func TestEchoHandler_BinaryBody(t *testing.T) {
	handler := webknife.NewEchoHandler(webknife.EchoConfig{})
	body := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection webknife.RequestInspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if !inspection.BodyBinary {
		t.Fatal("expected body_binary=true")
	}
}

func TestEchoHandler_RemoteAddr(t *testing.T) {
	handler := webknife.NewEchoHandler(webknife.EchoConfig{})
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	var inspection webknife.RequestInspection
	json.Unmarshal(w.Body.Bytes(), &inspection)
	if inspection.RemoteAddr != "192.168.1.100:12345" {
		t.Fatalf("expected 192.168.1.100:12345, got %s", inspection.RemoteAddr)
	}
}

func TestRespondHandler_200(t *testing.T) {
	handler, err := webknife.NewRespondHandler(webknife.RespondConfig{
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
	handler, err := webknife.NewRespondHandler(webknife.RespondConfig{
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
	handler, err := webknife.NewRespondHandler(webknife.RespondConfig{
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
	handler, err := webknife.NewRespondHandler(webknife.RespondConfig{
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
	_, err := webknife.NewRespondHandler(webknife.RespondConfig{StatusCode: 99})
	if err == nil {
		t.Fatal("expected error for invalid status code")
	}
	_, err = webknife.NewRespondHandler(webknife.RespondConfig{StatusCode: 600})
	if err == nil {
		t.Fatal("expected error for invalid status code")
	}
}

func TestRedirectHandler_302(t *testing.T) {
	handler, err := webknife.NewRedirectHandler(webknife.RedirectConfig{
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
	handler, err := webknife.NewRedirectHandler(webknife.RedirectConfig{
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
	handler, err := webknife.NewRedirectHandler(webknife.RedirectConfig{
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
	_, err := webknife.NewRedirectHandler(webknife.RedirectConfig{
		StatusCode: 200,
		Target:     "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error for invalid redirect status")
	}
	_, err = webknife.NewRedirectHandler(webknife.RedirectConfig{
		StatusCode: 404,
		Target:     "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error for invalid redirect status")
	}
}

func TestHeaderMiddleware_SetRequest(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Received", r.Header.Get("X-Debug"))
		w.WriteHeader(http.StatusOK)
	})
	handler := webknife.HeaderMiddleware(webknife.HeaderMiddlewareConfig{
		Operations: []webknife.HeaderOperation{
			{Type: webknife.HeaderSet, Name: "X-Debug", Value: "true"},
		},
		Direction: webknife.RequestHeaders,
	})(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("X-Received") != "true" {
		t.Fatalf("expected X-Received=true, got %s", w.Header().Get("X-Received"))
	}
}

func TestHeaderMiddleware_AddRequest(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vals := r.Header.Values("X-Multi")
		w.Header().Set("X-Count", fmt.Sprintf("%d", len(vals)))
		w.WriteHeader(http.StatusOK)
	})
	handler := webknife.HeaderMiddleware(webknife.HeaderMiddlewareConfig{
		Operations: []webknife.HeaderOperation{
			{Type: webknife.HeaderAdd, Name: "X-Multi", Value: "a"},
			{Type: webknife.HeaderAdd, Name: "X-Multi", Value: "b"},
		},
		Direction: webknife.RequestHeaders,
	})(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("X-Count") != "2" {
		t.Fatalf("expected X-Count=2, got %s", w.Header().Get("X-Count"))
	}
}

func TestHeaderMiddleware_RemoveRequest(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("Authorization")
		w.Header().Set("X-Auth-Present", fmt.Sprintf("%t", val != ""))
		w.WriteHeader(http.StatusOK)
	})
	handler := webknife.HeaderMiddleware(webknife.HeaderMiddlewareConfig{
		Operations: []webknife.HeaderOperation{
			{Type: webknife.HeaderRemove, Name: "Authorization"},
		},
		Direction: webknife.RequestHeaders,
	})(inner)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("X-Auth-Present") != "false" {
		t.Fatalf("expected auth removed, got X-Auth-Present=%s", w.Header().Get("X-Auth-Present"))
	}
}

func TestHeaderMiddleware_SetResponse(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := webknife.HeaderMiddleware(webknife.HeaderMiddlewareConfig{
		Operations: []webknife.HeaderOperation{
			{Type: webknife.HeaderSet, Name: "Cache-Control", Value: "no-store"},
		},
		Direction: webknife.ResponseHeaders,
	})(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("expected no-store, got %s", w.Header().Get("Cache-Control"))
	}
}

func TestHeaderMiddleware_RemoveResponse(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "nginx/1.0")
		w.WriteHeader(http.StatusOK)
	})
	handler := webknife.HeaderMiddleware(webknife.HeaderMiddlewareConfig{
		Operations: []webknife.HeaderOperation{
			{Type: webknife.HeaderRemove, Name: "Server"},
		},
		Direction: webknife.ResponseHeaders,
	})(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("Server") != "" {
		t.Fatalf("expected Server removed, got %s", w.Header().Get("Server"))
	}
}

func TestHeaderMiddleware_CaseInsensitive(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Debug", r.Header.Get("X-Debug"))
		w.WriteHeader(http.StatusOK)
	})
	handler := webknife.HeaderMiddleware(webknife.HeaderMiddlewareConfig{
		Operations: []webknife.HeaderOperation{
			{Type: webknife.HeaderSet, Name: "x-debug", Value: "true"},
		},
		Direction: webknife.RequestHeaders,
	})(inner)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("X-Debug") != "true" {
		t.Fatalf("expected X-Debug=true, got %s", w.Header().Get("X-Debug"))
	}
}

func TestBuildHeaderMiddlewares(t *testing.T) {
	cfg := webknife.MultiHeaderConfig{
		SetRequestHeaders:     []string{"X-Debug: true"},
		RemoveResponseHeaders: []string{"Server"},
	}
	mws, err := webknife.BuildHeaderMiddlewares(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(mws) != 2 {
		t.Fatalf("expected 2 middlewares, got %d", len(mws))
	}
}

func TestBuildHeaderMiddlewares_Empty(t *testing.T) {
	cfg := webknife.MultiHeaderConfig{}
	mws, err := webknife.BuildHeaderMiddlewares(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(mws) != 0 {
		t.Fatalf("expected 0 middlewares, got %d", len(mws))
	}
}

func TestRedactValue(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"a", "****"},
		{"ab", "****"},
		{"abc", "****"},
		{"abcd", "****"},
		{"abcde", "ab*de"},
		{"abcdef", "ab**ef"},
		{"secret123", "se*****23"},
	}
	for _, tt := range tests {
		got := webknife.RedactValue(tt.input)
		if got != tt.want {
			t.Errorf("RedactValue(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRedactHeader(t *testing.T) {
	if !webknife.RedactHeader("Authorization") {
		t.Fatal("expected Authorization to be sensitive")
	}
	if !webknife.RedactHeader("authorization") {
		t.Fatal("expected authorization to be sensitive")
	}
	if !webknife.RedactHeader("Cookie") {
		t.Fatal("expected Cookie to be sensitive")
	}
	if webknife.RedactHeader("X-Custom") {
		t.Fatal("expected X-Custom to not be sensitive")
	}
}

func TestParseHeaders(t *testing.T) {
	headers, err := webknife.ParseHeaders([]string{"X-Foo: bar", "X-Baz: qux"})
	if err != nil {
		t.Fatal(err)
	}
	if headers["X-Foo"] != "bar" {
		t.Fatalf("expected X-Foo=bar, got %s", headers["X-Foo"])
	}
	if headers["X-Baz"] != "qux" {
		t.Fatalf("expected X-Baz=qux, got %s", headers["X-Baz"])
	}
}

func TestParseHeaders_Invalid(t *testing.T) {
	_, err := webknife.ParseHeaders([]string{"invalid"})
	if err == nil {
		t.Fatal("expected error for invalid header")
	}
	_, err = webknife.ParseHeaders([]string{": value"})
	if err == nil {
		t.Fatal("expected error for empty header name")
	}
}
