package cli_test

import (
	"bytes"
	"testing"

	"github.com/webknife/webknife/internal/adapters/cli"
)

func TestParseArgs_NoCommand(t *testing.T) {
	_, _, err := cli.ParseArgs([]string{"webknife"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected error for no command")
	}
}

func TestParseArgs_Version(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	cmd, _, err := cli.ParseArgs([]string{"webknife", "version"}, stdout, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "" {
		t.Fatalf("expected empty cmd for version, got %q", cmd)
	}
	if !bytes.Contains(stdout.Bytes(), []byte("webknife")) {
		t.Fatal("expected version output")
	}
}

func TestParseArgs_Help(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	cmd, _, err := cli.ParseArgs([]string{"webknife", "help"}, stdout, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "" {
		t.Fatalf("expected empty cmd for help, got %q", cmd)
	}
	if !bytes.Contains(stdout.Bytes(), []byte("Usage")) {
		t.Fatal("expected usage output")
	}
}

func TestParseArgs_UnknownCommand(t *testing.T) {
	_, _, err := cli.ParseArgs([]string{"webknife", "unknown"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestParseArgs_Serve(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "serve", "--listen", ":9090", "--root", "/tmp"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "serve" {
		t.Fatalf("expected 'serve', got %q", cmd)
	}
	if cfg.ListenAddr != ":9090" {
		t.Fatalf("expected ':9090', got %q", cfg.ListenAddr)
	}
	if cfg.Root != "/tmp" {
		t.Fatalf("expected '/tmp', got %q", cfg.Root)
	}
}

func TestParseArgs_ServeDefaults(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "serve"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "serve" {
		t.Fatalf("expected 'serve', got %q", cmd)
	}
	if cfg.ListenAddr != ":8080" {
		t.Fatalf("expected ':8080', got %q", cfg.ListenAddr)
	}
	if cfg.Root != "." {
		t.Fatalf("expected '.', got %q", cfg.Root)
	}
}

func TestParseArgs_ServeAuth(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "serve", "--auth", "admin:secret"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "serve" {
		t.Fatalf("expected 'serve', got %q", cmd)
	}
	if cfg.Auth != "admin:secret" {
		t.Fatalf("expected 'admin:secret', got %q", cfg.Auth)
	}
}

func TestParseArgs_ServeTLS(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "serve", "--tls-cert", "cert.pem", "--tls-key", "key.pem"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "serve" {
		t.Fatalf("expected 'serve', got %q", cmd)
	}
	if cfg.TLSCert != "cert.pem" {
		t.Fatalf("expected 'cert.pem', got %q", cfg.TLSCert)
	}
	if cfg.TLSKey != "key.pem" {
		t.Fatalf("expected 'key.pem', got %q", cfg.TLSKey)
	}
}

func TestParseArgs_ServeLogFormat(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "serve", "--log-format", "json"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "serve" {
		t.Fatalf("expected 'serve', got %q", cmd)
	}
	if cfg.LogFormat != "json" {
		t.Fatalf("expected 'json', got %q", cfg.LogFormat)
	}
}

func TestParseArgs_Proxy(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "proxy", "--listen", ":9090", "--upstream", "http://localhost:3000"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "proxy" {
		t.Fatalf("expected 'proxy', got %q", cmd)
	}
	if cfg.ListenAddr != ":9090" {
		t.Fatalf("expected ':9090', got %q", cfg.ListenAddr)
	}
	if cfg.Upstream != "http://localhost:3000" {
		t.Fatalf("expected 'http://localhost:3000', got %q", cfg.Upstream)
	}
}

func TestParseArgs_ProxyNoUpstream(t *testing.T) {
	_, _, err := cli.ParseArgs([]string{"webknife", "proxy"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected error for missing upstream")
	}
}

func TestParseArgs_ProxyAuth(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "proxy", "--upstream", "http://localhost:3000", "--auth", "admin:secret"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "proxy" {
		t.Fatalf("expected 'proxy', got %q", cmd)
	}
	if cfg.Auth != "admin:secret" {
		t.Fatalf("expected 'admin:secret', got %q", cfg.Auth)
	}
}

func TestParseArgs_Echo(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "echo", "--listen", ":9090"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "echo" {
		t.Fatalf("expected 'echo', got %q", cmd)
	}
	if cfg.ListenAddr != ":9090" {
		t.Fatalf("expected ':9090', got %q", cfg.ListenAddr)
	}
}

func TestParseArgs_EchoDefaults(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "echo"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "echo" {
		t.Fatalf("expected 'echo', got %q", cmd)
	}
	if cfg.ListenAddr != ":8080" {
		t.Fatalf("expected ':8080', got %q", cfg.ListenAddr)
	}
	if cfg.EchoMaxBody != 1<<20 {
		t.Fatalf("expected max body 1MB, got %d", cfg.EchoMaxBody)
	}
}

func TestParseArgs_Respond(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "respond", "--status", "503", "--body", "Service unavailable"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "respond" {
		t.Fatalf("expected 'respond', got %q", cmd)
	}
	if cfg.RespondStatusCode != 503 {
		t.Fatalf("expected 503, got %d", cfg.RespondStatusCode)
	}
	if cfg.RespondBody != "Service unavailable" {
		t.Fatalf("expected 'Service unavailable', got %q", cfg.RespondBody)
	}
}

func TestParseArgs_RespondDefaults(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "respond"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "respond" {
		t.Fatalf("expected 'respond', got %q", cmd)
	}
	if cfg.RespondStatusCode != 200 {
		t.Fatalf("expected 200, got %d", cfg.RespondStatusCode)
	}
}

func TestParseArgs_RespondContentType(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "respond", "--status", "429", "--content-type", "application/json"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "respond" {
		t.Fatalf("expected 'respond', got %q", cmd)
	}
	if cfg.RespondStatusCode != 429 {
		t.Fatalf("expected 429, got %d", cfg.RespondStatusCode)
	}
	if cfg.RespondContentType != "application/json" {
		t.Fatalf("expected application/json, got %q", cfg.RespondContentType)
	}
}

func TestParseArgs_Redirect(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "redirect", "--to", "https://example.com"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "redirect" {
		t.Fatalf("expected 'redirect', got %q", cmd)
	}
	if cfg.RedirectTarget != "https://example.com" {
		t.Fatalf("expected https://example.com, got %q", cfg.RedirectTarget)
	}
	if cfg.RedirectStatusCode != 302 {
		t.Fatalf("expected 302, got %d", cfg.RedirectStatusCode)
	}
}

func TestParseArgs_RedirectStatus(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "redirect", "--to", "https://example.com", "--status", "301"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "redirect" {
		t.Fatalf("expected 'redirect', got %q", cmd)
	}
	if cfg.RedirectStatusCode != 301 {
		t.Fatalf("expected 301, got %d", cfg.RedirectStatusCode)
	}
}

func TestParseArgs_RedirectPreservePath(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{"webknife", "redirect", "--to", "https://example.com", "--preserve-path"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "redirect" {
		t.Fatalf("expected 'redirect', got %q", cmd)
	}
	if !cfg.RedirectPreservePath {
		t.Fatal("expected preserve-path=true")
	}
}

func TestParseArgs_RedirectNoTarget(t *testing.T) {
	_, _, err := cli.ParseArgs([]string{"webknife", "redirect"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected error for missing --to")
	}
}

func TestParseArgs_HeaderFlags(t *testing.T) {
	cmd, cfg, err := cli.ParseArgs([]string{
		"webknife", "proxy",
		"--upstream", "http://localhost:3000",
		"--set-request-header", "X-Debug: true",
		"--remove-response-header", "Server",
	}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "proxy" {
		t.Fatalf("expected 'proxy', got %q", cmd)
	}
	if len(cfg.SetRequestHeaders) != 1 || cfg.SetRequestHeaders[0] != "X-Debug: true" {
		t.Fatalf("unexpected set-request-headers: %v", cfg.SetRequestHeaders)
	}
	if len(cfg.RemoveResponseHeaders) != 1 || cfg.RemoveResponseHeaders[0] != "Server" {
		t.Fatalf("unexpected remove-response-headers: %v", cfg.RemoveResponseHeaders)
	}
}

func TestParseArgs_Verbose(t *testing.T) {
	_, cfg, err := cli.ParseArgs([]string{"webknife", "echo", "-v"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Verbose {
		t.Fatal("expected verbose=true")
	}
}

func TestParseArgs_EchoMaxBody(t *testing.T) {
	_, cfg, err := cli.ParseArgs([]string{"webknife", "echo", "--max-body", "1024"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.EchoMaxBody != 1024 {
		t.Fatalf("expected 1024, got %d", cfg.EchoMaxBody)
	}
}
