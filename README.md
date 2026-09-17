<p align="center">
  <img src="assets/images/webknife-logo.png" alt="Webknife" width="400">
</p>

# Webknife

A single-binary, CLI-first HTTP laboratory for diagnostics, testing, proxying, serving, inspection, and experimentation.

## Why Webknife?

When troubleshooting HTTP infrastructure, you often need a quick, disposable server that does one thing well — serve files, proxy a port, return a specific status code, or show you exactly what a client is sending. Webknife is that tool. It starts in seconds from a shell, requires no configuration files, and does not pretend to be a production web server.

It is **not** Nginx, Apache, Caddy, HAProxy, or mitmproxy. It is closer to a Swiss Army knife you keep in your terminal for moments like:

- Verifying whether NAT or firewall forwarding reaches a server
- Temporarily serving a directory over HTTP
- Putting Basic Authentication in front of something
- Quickly reverse proxying a local service to test CORS or headers
- Inspecting incoming webhook requests without writing code
- Simulating HTTP 503 or 429 error responses
- Testing redirects before implementing them
- Adding diagnostic headers to proxied traffic
- Debugging an application from an SSH session
- Experimenting with TLS certificates

## Features

| Capability | Description |
|---|---|
| **Static file serving** | Serve files from any directory with correct MIME types and directory handling |
| **Reverse proxy** | Proxy requests to an upstream service, forwarding method, path, headers, and body |
| **Echo / inspect** | Reflect incoming requests as structured JSON with full diagnostic detail |
| **Arbitrary responses** | Return any HTTP status code with custom body, Content-Type, and headers |
| **Redirects** | Redirect to any URL with configurable status codes (301, 302, 303, 307, 308) |
| **Header manipulation** | Set, add, or remove request and response headers on any command |
| **Basic Authentication** | Optional auth composable with all handlers |
| **TLS** | Support for certificate and private key files |
| **Request logging** | Concise human-readable or structured JSON output |
| **Event-driven** | Observable lifecycle events for integration and observability |
| **Library-first** | Core functionality usable as a Go package without the CLI |

## Quick Start

```bash
# Inspect exactly what a client sends
webknife echo --listen :8080

# Serve the current directory
webknife serve --listen :8080 --root .

# Protect it with Basic Authentication
webknife serve --listen :8080 --root . --auth admin:secret

# Reverse proxy a local application
webknife proxy --listen :8080 --upstream http://localhost:3000

# Simulate an unavailable service
webknife respond --listen :8080 --status 503 --body 'Service temporarily unavailable'

# Simulate rate limiting
webknife respond --listen :8080 --status 429 --header 'Retry-After: 60'

# Create a redirect
webknife redirect --listen :8080 --to https://example.com --preserve-path

# Proxy with header manipulation
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Debug: true' \
  --remove-response-header Server
```

## Installation

### From source

```bash
git clone https://github.com/aognio/webknife
cd webknife
make build
```

The binary is placed at `bin/webknife`.

### Using `go install`

```bash
go install github.com/aognio/webknife/cmd/webknife@latest
```

### With Make

```bash
make build        # build to bin/webknife
make test         # run test suite
make check        # fmt-check + vet + test
make help         # show all available targets
```

## Static File Serving

```bash
webknife serve --listen :8080 --root ./public
```

Serves files from the specified directory. Directory listings use Go's standard `http.FileServer`. HEAD requests are supported. MIME types are handled automatically. Path traversal outside the document root is prevented.

The simplest invocation:

```bash
webknife serve --listen :8080 --root .
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `--listen` | `:8080` | Listen address |
| `--root` | `.` | Document root directory |
| `--auth` | | Basic auth credentials (`user:password`) |
| `--tls-cert` | | TLS certificate file |
| `--tls-key` | | TLS private key file |
| `--log-format` | `text` | Log format: `text` or `json` |

## Reverse Proxy

```bash
webknife proxy --listen :8080 --upstream http://localhost:3000
```

Proxies all requests to the upstream service. Correctly forwards:

- HTTP method
- Path and query string
- Headers
- Request body
- Response status, headers, and body

Authentication composes with proxying:

```bash
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --auth admin:secret
```

## Request Inspection / Echo

```bash
webknife echo --listen :8080
```

Returns structured JSON describing the incoming request:

```json
{
  "method": "POST",
  "scheme": "http",
  "host": "localhost:8080",
  "path": "/webhook",
  "query": {"token": ["abc123"]},
  "proto": "HTTP/1.1",
  "content_length": 42,
  "remote_addr": "192.168.1.5:54321",
  "headers": {
    "Content-Type": ["application/json"],
    "X-Webhook-Secret": ["verify-me"]
  },
  "request_uri": "/webhook?token=abc123",
  "body": "{\"event\":\"push\",\"ref\":\"refs/heads/main\"}"
}
```

Includes TLS info when available. Binary bodies are hex-encoded. Large bodies are truncated at the configurable limit (default 1MB).

```bash
webknife echo --listen :8080 --max-body 4096
```

## Basic Authentication

```bash
webknife serve --listen :8080 --root . --auth admin:secret
```

Returns `401 Unauthorized` with `WWW-Authenticate: Basic realm="..."` when credentials are missing or wrong. Uses constant-time comparison to prevent timing attacks. Credentials never appear in logs.

Authentication composes with every handler: `serve`, `proxy`, `echo`, `respond`, and `redirect`.

## TLS

```bash
webknife serve --listen :443 \
  --root ./public \
  --tls-cert ./cert.pem \
  --tls-key ./key.pem
```

Works with all commands. TLS configuration uses Go's recommended minimum version (TLS 1.2) and cipher suites.

## Custom Responses

```bash
# Return a 503 error
webknife respond --listen :8080 --status 503 --body 'Service unavailable'

# Return a JSON error with headers
webknife respond --listen :8080 \
  --status 429 \
  --header 'Retry-After: 60' \
  --body '{"error":"rate limited"}' \
  --content-type application/json

# Serve a response body from a file
webknife respond --listen :8080 \
  --status 200 \
  --body-file ./response.json \
  --content-type application/json
```

## Redirects

```bash
# Simple redirect
webknife redirect --listen :8080 --to https://example.com

# Permanent redirect
webknife redirect --listen :8080 --to https://example.com --status 301

# Preserve the original path and query string
webknife redirect --listen :8080 \
  --to https://example.com \
  --preserve-path
```

With `--preserve-path`, a request to `http://localhost:8080/foo?a=1` becomes `https://example.com/foo?a=1`.

Supported status codes: 301, 302, 303, 307, 308.

## Header Manipulation

All commands support request and response header flags:

```bash
# Set (replace) a request header
--set-request-header 'X-Debug: true'

# Add a request header (appends, does not replace)
--add-request-header 'X-Custom: value'

# Remove a request header
--remove-request-header Authorization

# Set a response header
--set-response-header 'Cache-Control: no-store'

# Add a response header
--add-response-header 'X-Multi: value'

# Remove a response header
--remove-response-header Server
```

These compose with authentication and all handlers:

```bash
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --auth admin:secret \
  --set-request-header 'X-Debug: true' \
  --remove-response-header Server
```

## Logging and Observability

### Text format (default)

```
2026-09-16T21:42:10-05:00 GET /index.html 200 4.2KB 1.7ms 192.168.1.1:54321
```

### JSON format

```bash
webknife serve --listen :8080 --root . --log-format json
```

```json
{"timestamp":"2026-09-16T21:42:10-05:00","method":"GET","path":"/index.html","status":200,"size":"4.2KB","duration":"1.7ms","remote":"192.168.1.1:54321"}
```

### Events

Webknife publishes lifecycle events that can be consumed programmatically when using the library:

- `ServerStarted`, `ServerStopped`
- `RequestReceived`, `RequestCompleted`, `RequestFailed`
- `ProxyRequestStarted`, `ProxyRequestCompleted`
- `AuthenticationSucceeded`, `AuthenticationFailed`
- `TLSConnectionAccepted`
- `RequestInspected`, `ResponseGenerated`, `RedirectGenerated`
- `RequestHeadersModified`, `ResponseHeadersModified`

Events describe facts that already happened, not commands.

## CLI Reference

```
webknife <command> [options]

Commands:
  serve     Serve static files
  proxy     Reverse proxy to an upstream
  echo      Inspect incoming requests
  respond   Return an arbitrary HTTP response
  redirect  Redirect to a URL
  version   Show version
  help      Show this help
```

### Common flags (all commands)

| Flag | Default | Description |
|---|---|---|
| `--listen` | `:8080` | Listen address |
| `--auth` | | Basic auth credentials (`user:password`) |
| `--tls-cert` | | TLS certificate file |
| `--tls-key` | | TLS private key file |
| `--log-format` | `text` | Log format: `text` or `json` |
| `-v`, `--verbose` | `false` | Verbose output |

### Header flags (all commands)

| Flag | Description |
|---|---|
| `--set-request-header 'Name: Value'` | Set/replace a request header |
| `--add-request-header 'Name: Value'` | Add a request header |
| `--remove-request-header Name` | Remove a request header |
| `--set-response-header 'Name: Value'` | Set/replace a response header |
| `--add-response-header 'Name: Value'` | Add a response header |
| `--remove-response-header Name` | Remove a response header |

## Using Webknife as a Go Library

The core functionality lives in `pkg/webknife` and is usable independently of the CLI:

```go
package main

import (
    "log"
    "net/http"

    "github.com/aognio/webknife/pkg/webknife"
)

func main() {
    // Echo handler
    echo := webknife.NewEchoHandler(webknife.EchoConfig{
        MaxBodySize: 1 << 20,
    })

    // Compose with Basic Auth
    handler := webknife.BasicAuth(webknife.BasicAuthConfig{
        Username: "admin",
        Password: "secret",
    }, echo)

    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

The CLI is an adapter. The library is the product.

## Architecture

```
cmd/webknife/                  Composition root (wires CLI → application)
internal/
  adapters/
    cli/                       CLI parsing and validation (thin)
    logging/                   Log adapters (text, JSON)
  application/                 Orchestration layer (builds handlers, wires middleware)
  archtest/                    Architectural dependency enforcement tests
  build/                       Version and build metadata
  events/                      Event types and event bus implementation
pkg/
  webknife/                    Ports and shared types (Middleware, Publisher)
  auth/                        Basic Auth middleware
  echo/                        Request inspection handler
  headers/                     Header manipulation middleware
  observe/                     Request observation middleware
  proxy/                       Reverse proxy handler
  redact/                      Secret detection and masking
  redirect/                    Redirect handler
  respond/                     Arbitrary response handler
  server/                      HTTP server lifecycle
  static/                      Static file handler
```

### Dependency direction

```
cmd/webknife/ (composition root)
    │
    ▼
internal/application/ (wires everything)
    │
    ├──► internal/adapters/cli/    (CLI parsing)
    ├──► internal/adapters/logging/ (log output)
    │
    └──► pkg/ (feature packages)
         ├──► pkg/webknife/        (ports: Publisher, Middleware)
         ├──► pkg/observe/         (publishes events via Publisher port)
         ├──► pkg/auth/
         ├──► pkg/headers/
         ├──► pkg/echo/
         ├──► pkg/respond/
         ├──► pkg/redirect/
         ├──► pkg/proxy/
         ├──► pkg/static/
         └──► pkg/server/
```

**Architectural rules (enforced by tests):**
- Feature packages (`pkg/*`) never import each other (except `pkg/webknife` for ports)
- Feature packages never import `internal/*`
- `pkg/webknife` (ports) never imports `internal/*`
- CLI adapter never imports feature packages directly (goes through application layer)

### Design principles

- **CLI-first** — the binary is the primary interface for operators
- **Library-first** — every `pkg/*` feature is independently reusable without the CLI
- **Ports and Adapters** — feature packages define no shared interfaces; the application layer wires them together
- **Composition** — middleware and handlers compose via `http.Handler`, not inheritance
- **Standard library first** — Go's stdlib provides the HTTP server, reverse proxy, file server, and TLS; external dependencies are avoided unless they provide substantial concrete value
- **Independent features** — each `pkg/*` package has zero sibling dependencies; adding a new feature never couples to existing ones

### Request processing pipeline

```
HTTP Listener
  │
  ▼
Request Observation (events.RequestMiddleware)
  │
  ▼
Basic Authentication (optional)
  │
  ▼
Request Header Transformation (optional)
  │
  ▼
Handler (echo / respond / redirect / static / proxy)
  │
  ▼
Response Header Transformation (optional)
  │
  ▼
Response Observation (events)
```

## Security Considerations

Webknife is a diagnostic tool. Convenience-oriented defaults do not make it a hardened production server.

**Credentials on the command line.** The `--auth` flag passes credentials as a command-line argument. These are visible in process listings (`ps aux`) and shell history. For sensitive use, prefer TLS and environment variables in scripts.

**Binding to public interfaces.** By default, `--listen :8080` binds to all interfaces. Use `--listen 127.0.0.1:8080` to restrict to localhost.

**Privileged ports.** Ports below 1024 require elevated permissions on most systems.

**Filesystem exposure.** `webknife serve --root /` exposes the entire filesystem. Be deliberate about the root directory.

**Proxy exposure.** `webknife proxy` makes the upstream accessible through the proxy. Anyone who can reach the proxy can reach the upstream.

**Sensitive headers.** Webknife identifies `Authorization`, `Cookie`, `Set-Cookie`, and `Proxy-Authorization` as sensitive. The `RedactValue` utility masks these in output.

**Request body limits.** The echo handler enforces a configurable maximum body size (default 1MB) to prevent memory exhaustion from large uploads.

**TLS.** Use TLS for any traffic containing sensitive data. The server enforces TLS 1.2 minimum.

**Diagnostic defaults.** Logging includes remote addresses, request paths, and status codes. Do not enable verbose logging in production environments with untrusted traffic.

## Recipes

### Verify port forwarding works

```bash
# On the server:
webknife respond --listen :8080 --status 200 --body 'reachable'

# From the client:
curl http://server:8080/
```

### Serve a directory temporarily

```bash
webknife serve --listen :8080 --root ./dist
```

### Put authentication in front of a static site

```bash
webknife serve --listen :8080 --root ./site --auth deploy:secret-token
```

### Inspect an incoming webhook

```bash
webknife echo --listen :9090 --log-format json
```

Point your webhook at `http://your-server:9090/` and see the full request in JSON.

### Proxy traffic to a development server

```bash
webknife proxy --listen :8080 --upstream http://localhost:3000
```

### Add a diagnostic header while proxying

```bash
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Forwarded-By: webknife'
```

### Return a deliberate HTTP 503

```bash
webknife respond --listen :8080 --status 503 --body 'Under maintenance'
```

### Simulate rate limiting

```bash
webknife respond --listen :8080 \
  --status 429 \
  --header 'Retry-After: 60' \
  --body '{"error":"rate limited"}' \
  --content-type application/json
```

### Test HTTPS with an existing certificate

```bash
webknife serve --listen :443 \
  --root ./public \
  --tls-cert ./server.crt \
  --tls-key ./server.key
```

### Inspect traffic passing through a proxy

```bash
# Upstream server
webknife serve --listen :3000 --root ./api

# Proxy with debug headers
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Debug: true' \
  --set-response-header 'X-Proxied-By: webknife'
```

## Development

### Prerequisites

- Go 1.22+
- GNU Make

### Build

```bash
make build     # binary at bin/webknife
make release   # stripped, static binary
```

### Test

```bash
make test      # run all tests
make test-race # with race detector
make check     # fmt-check + vet + test
```

### Format

```bash
make fmt       # format with gofumpt or gofmt
make fmt-check # fail if unformatted
```

### Version injection

The build injects version, commit, and date via `ldflags`:

```bash
make build VERSION=v1.0.0
```

Or let it default to `git describe`.

## Testing

```bash
make test
```

Tests use `httptest` extensively and require no external Internet access. Test coverage includes:

- Static file serving, missing files, HEAD requests, path traversal
- Basic Auth success, failure, missing credentials
- Proxy forwarding, response propagation, header forwarding
- Echo request inspection (GET, POST, query, headers, cookies, body truncation, binary bodies)
- Arbitrary response generation, custom headers, body from file
- Redirect status codes, destination, path preservation
- Header manipulation (set, add, remove) for request and response
- Event publishing and middleware composition
- CLI argument parsing for all commands
- Secret redaction

### Manual Test Harness

A Bash-based manual test harness is available for interactive exploration:

```bash
scripts/test-webknife.sh help
```

This displays all available scenarios:

```bash
scripts/test-webknife.sh 01   # Static file server
scripts/test-webknife.sh 03   # Echo / request inspection
scripts/test-webknife.sh 09   # Reverse proxy
```

Each scenario starts a Webknife server in the foreground and prints `curl` commands to run from another terminal. Press `Ctrl+C` to stop.

Environment variables for overriding defaults:

```bash
PORT=9090 scripts/test-webknife.sh 03
WEBKNIFE=/usr/local/bin/webknife scripts/test-webknife.sh 01
```

Or via Make:

```bash
make manual-test TEST=01
```

## Project Philosophy

Webknife is built on a few principles:

1. **Diagnose, don't replace.** It fills the gap between `curl` and a production web server.
2. **Start in seconds.** No config files, no Docker, no dependencies. Just run the binary.
3. **Compose, don't accumulate.** Every capability composes with every other through `http.Handler`.
4. **Observe, don't guess.** Events and logging make behavior visible.
5. **Library first.** The CLI is a convenience adapter over a reusable Go library.
6. **Standard library first.** Go's HTTP stack is excellent. Webknife builds on it, not around it.

## Roadmap

Future iterations may add:

- Virtual hosts and host/path/method routing
- Request matching and configurable middleware pipelines
- Delays, bandwidth throttling, connection resets
- Probabilistic failures and fault injection
- Webhook capture, request persistence, and replay
- WebSockets, HTTP/2, HTTP/3 experimentation
- Self-signed certificates, ACME, mTLS, TLS diagnostics
- TCP utilities
- YAML/TOML and environment variable configuration

These are architectural direction, not current capabilities. The architecture is designed to accommodate them naturally.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes following the project conventions
4. Run `make check` to verify
5. Submit a pull request

Follow existing code style. Prefer standard library. Write tests for new functionality. Keep the dependency graph clean.

## License

MIT License - Copyright (c) 2026 Antonio Ognio

Made with ❤️ from 🇵🇪. El Perú es clave 🔑.
