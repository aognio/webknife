---
title: "Testing"
weight: 2
---

# Testing

Webknife has two types of tests: automated Go tests and a manual test harness.

## Automated tests

```bash
make test
```

Tests use `httptest` extensively and require no external Internet access. They cover:

- Static file serving, missing files, HEAD requests, path traversal
- Basic Auth success, failure, missing credentials
- Proxy forwarding, response propagation, header forwarding
- Echo request inspection (GET, POST, query, headers, cookies, body truncation)
- Arbitrary response generation, custom headers, body from file
- Redirect status codes, destination, path preservation
- Header manipulation (set, add, remove) for request and response
- Event publishing and middleware composition
- CLI argument parsing for all commands
- Secret redaction
- Architectural dependency rules

### Architectural tests

```bash
go test ./internal/archtest/
```

Enforces dependency rules:
- Feature packages don't import siblings
- Feature packages don't import `internal/*`
- `pkg/webknife` doesn't import `internal/*`
- CLI doesn't import feature packages directly

## Manual test harness

For end-to-end HTTP experimentation:

```bash
scripts/test-webknife.sh help     # list all scenarios
scripts/test-webknife.sh 03       # run scenario 03 (echo)
```

Each scenario starts a Webknife server in the foreground and prints `curl` commands for manual inspection. Press `Ctrl+C` to stop.

### When to use each

| Tool | Purpose |
|------|---------|
| `go test ./...` | Deterministic automated correctness |
| `scripts/test-webknife.sh` | Manual end-to-end HTTP experimentation |

### Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WEBKNIFE` | `./bin/webknife` | Path to webknife binary |
| `HOST` | `127.0.0.1` | Listen host |
| `PORT` | `8080` | Listen port |

### Examples

```bash
PORT=9090 scripts/test-webknife.sh 03
WEBKNIFE=/usr/local/bin/webknife scripts/test-webknife.sh 01
```

Or via Make:

```bash
make manual-test TEST=01
```

### Available scenarios

| # | Description |
|---|-------------|
| 01 | Static file server |
| 02 | Static server with Basic Authentication |
| 03 | Echo / request inspection |
| 04 | Custom HTTP response |
| 05 | Custom error response / HTTP 503 |
| 06 | HTTP redirect |
| 07 | Request header manipulation |
| 08 | Response header manipulation |
| 09 | Reverse proxy |
| 10 | Structured JSON logging |
