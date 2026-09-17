# webknife

A single-binary, CLI-first HTTP and network diagnostic Swiss Army knife.

> A portable HTTP laboratory that can be started in seconds from a shell.

## What it is

webknife is a temporary diagnostics tool for developers, sysadmins, DevOps/SRE engineers, and testers troubleshooting HTTP infrastructure from a terminal.

It is **not** intended to replace Nginx, Apache, Caddy, HAProxy, mitmproxy, or full production web servers. It excels at temporary diagnostics, experimentation, inspection, proxying, simulation, testing, and debugging.

## Features

- **Static file serving** — serve files from any directory with correct MIME types
- **Reverse proxy** — proxy requests to an upstream service
- **Echo/inspect** — reflect incoming requests as structured JSON
- **Arbitrary responses** — return any HTTP status code with custom body and headers
- **Redirects** — redirect to any URL with configurable status codes
- **Header manipulation** — set, add, or remove request and response headers
- **HTTP Basic Authentication** — optional, composable with all handlers
- **TLS** — support for certificate and private key files
- **Request logging** — concise human-readable or structured JSON output
- **Secret redaction** — centralized sensitive header/parameter masking
- **Event-driven architecture** — observable lifecycle events

## Installation

```bash
go install github.com/webknife/webknife/cmd/webknife@latest
```

Or build from source:

```bash
git clone https://github.com/webknife/webknife
cd webknife
go build -o webknife ./cmd/webknife
```

## Quick Start

```bash
# Inspect exactly what a client sends
webknife echo --listen :8080

# Serve the current directory
webknife serve --listen :8080 --root .

# Simulate an unavailable service
webknife respond --listen :8080 --status 503

# Simulate rate limiting
webknife respond \
  --listen :8080 \
  --status 429 \
  --header 'Retry-After: 60'

# Temporary redirect
webknife redirect \
  --listen :8080 \
  --to https://example.com

# Reverse proxy another service
webknife proxy --listen :8080 --upstream http://localhost:3000

# Modify traffic while proxying
webknife proxy \
  --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Debug: true' \
  --remove-response-header Server
```

## Commands

### `webknife echo`

Inspect incoming HTTP requests. Returns structured JSON with method, path, query parameters, headers, body, TLS info, and more.

```
webknife echo [options]

Options:
  --listen     Listen address (default ":8080")
  --max-body   Max body size for inspection in bytes (default 1048576)
  --auth       Basic auth credentials as user:password
  --tls-cert   TLS certificate file
  --tls-key    TLS private key file
  --log-format Log format: text or json (default "text")
```

Example output:

```json
{
  "method": "POST",
  "scheme": "http",
  "host": "localhost:8080",
  "path": "/test",
  "query": {"foo": ["bar", "baz"]},
  "proto": "HTTP/1.1",
  "content_length": 19,
  "remote_addr": "127.0.0.1:12345",
  "headers": {"Content-Type": ["application/json"]},
  "request_uri": "/test?foo=bar&foo=baz",
  "body": "{\"message\":\"hello\"}"
}
```

### `webknife respond`

Return an arbitrary HTTP response. Useful for simulating error conditions, rate limiting, or mock APIs.

```
webknife respond [options]

Options:
  --listen       Listen address (default ":8080")
  --status       HTTP status code (default 200)
  --body         Response body
  --body-file    Response body from file
  --content-type Content-Type header
  --header       Response header (Name: Value, repeatable)
  --auth         Basic auth credentials as user:password
  --tls-cert     TLS certificate file
  --tls-key      TLS private key file
  --log-format   Log format: text or json (default "text")
```

Examples:

```bash
# Return 503
webknife respond --listen :8080 --status 503 --body 'Service unavailable'

# Rate limit with retry header
webknife respond \
  --listen :8080 \
  --status 429 \
  --header 'Retry-After: 60' \
  --body '{"error":"rate limited"}' \
  --content-type application/json

# Serve a static JSON response from file
webknife respond \
  --status 200 \
  --body-file ./response.json \
  --content-type application/json
```

### `webknife redirect`

Redirect all requests to a target URL with configurable status codes.

```
webknife redirect [options]

Options:
  --listen        Listen address (default ":8080")
  --to            Redirect target URL (required)
  --status        Redirect status code: 301, 302, 303, 307, 308 (default 302)
  --preserve-path Preserve original path and query string
  --auth          Basic auth credentials as user:password
  --tls-cert      TLS certificate file
  --tls-key       TLS private key file
  --log-format    Log format: text or json (default "text")
```

Examples:

```bash
# Simple redirect
webknife redirect --listen :8080 --to https://example.com

# Permanent redirect preserving path
webknife redirect \
  --listen :8080 \
  --to https://example.com \
  --status 301 \
  --preserve-path
```

With `--preserve-path`, a request to `http://localhost:8080/foo?a=1` becomes `https://example.com/foo?a=1`.

### `webknife serve`

Serve static files from a directory.

```
webknife serve [options]

Options:
  --listen     Listen address (default ":8080")
  --root       Document root directory (default ".")
  --auth       Basic auth credentials as user:password
  --tls-cert   TLS certificate file
  --tls-key    TLS private key file
  --log-format Log format: text or json (default "text")
```

### `webknife proxy`

Reverse proxy to an upstream service.

```
webknife proxy [options]

Options:
  --listen     Listen address (default ":8080")
  --upstream   Upstream URL (required)
  --auth       Basic auth credentials as user:password
  --tls-cert   TLS certificate file
  --tls-key    TLS private key file
  --log-format Log format: text or json (default "text")
```

### `webknife version`

Print the version.

### `webknife help`

Show usage information.

## Header Manipulation

All server commands support header manipulation flags:

```bash
# Set/replace a request header
--set-request-header 'X-Debug: true'

# Add a request header (appends, does not replace)
--add-request-header 'X-Multi: value'

# Remove a request header
--remove-request-header Authorization

# Set/replace a response header
--set-response-header 'Cache-Control: no-store'

# Add a response header
--add-response-header 'X-Multi: value'

# Remove a response header
--remove-response-header Server
```

Header names are case-insensitive per HTTP semantics. These flags compose with authentication and all handlers:

```bash
webknife proxy \
  --listen :8080 \
  --upstream http://localhost:3000 \
  --auth admin:secret \
  --set-request-header 'X-Debug: true' \
  --remove-response-header Server
```

## Composition Pipeline

Request processing follows this pipeline:

```
HTTP Listener
  ↓
Request Observation (events)
  ↓
Basic Authentication (optional)
  ↓
Request Header Transformation (optional)
  ↓
Handler (echo / respond / redirect / static / proxy)
  ↓
Response Header Transformation (optional)
  ↓
Response Observation (events)
```

## Logging

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

## Secret Redaction

Sensitive headers (`Authorization`, `Cookie`, `Set-Cookie`, `Proxy-Authorization`) and query parameters (`token`, `key`, `secret`, `password`, `api_key`, `apikey`) are identified for redaction. The `RedactValue` function masks sensitive values while preserving first and last 2 characters for debugging.

## Architecture

```
cmd/webknife/          CLI entry point
internal/
  adapters/cli/        Command-line parsing and command routing
  events/              Event types and event bus
pkg/webknife/          Core library (reusable without CLI)
  echo.go              Request inspection handler
  respond.go           Arbitrary response handler
  redirect.go          Redirect handler
  headers.go           Header manipulation middleware
  auth.go              Basic authentication middleware
  proxy.go             Reverse proxy handler
  static.go            Static file handler
  logging.go           Log adapter
  redact.go            Secret redaction
  server.go            HTTP server and TLS config
  publisher.go         Event publisher adapter
```

### Design principles

- **CLI-first** — the binary is the primary interface
- **Library-first** — core functionality is usable without the CLI
- **Event-first** — important observable events are represented explicitly
- **Ports and adapters** — clean separation between domain logic and infrastructure

## Testing

```bash
go test ./...
```

## Current Limitations

- No automatic self-signed TLS certificates
- No ACME/Let's Encrypt support
- No mTLS
- No virtual hosts or path-based routing
- No WebSockets or HTTP/2
- No configuration files (YAML/TOML)
- No environment variable configuration

## Planned Direction

Future iterations may add:

- Delays, throttling, failure injection
- Webhook catching, request persistence, replay
- WebSockets, HTTP/2, HTTP/3
- Self-signed TLS, ACME, mTLS, TLS inspection
- Configurable middleware pipelines
- YAML/TOML configuration

## License

MIT
