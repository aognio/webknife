---
title: "Overview"
weight: 1
---

# CLI Reference

```
webknife <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| `serve` | Serve static files |
| `proxy` | Reverse proxy to an upstream |
| `echo` | Inspect incoming requests |
| `respond` | Return an arbitrary HTTP response |
| `redirect` | Redirect to a URL |
| `version` | Show version |
| `help` | Show help |

## Common options

These flags apply to all commands:

| Flag | Default | Description |
|------|---------|-------------|
| `--listen` | `:8080` | Listen address |
| `--auth` | | Basic auth credentials (`user:password`) |
| `--tls-cert` | | TLS certificate file |
| `--tls-key` | | TLS private key file |
| `--log-format` | `text` | Log format: `text` or `json` |
| `-v, --verbose` | `false` | Verbose output |

## Header manipulation

These flags apply to all commands:

| Flag | Description |
|------|-------------|
| `--set-request-header 'Name: Value'` | Set/replace a request header |
| `--add-request-header 'Name: Value'` | Add a value to a request header |
| `--remove-request-header Name` | Remove a request header |
| `--set-response-header 'Name: Value'` | Set/replace a response header |
| `--add-response-header 'Name: Value'` | Add a value to a response header |
| `--remove-response-header Name` | Remove a response header |

## Examples

```bash
webknife serve --listen :8080 --root ./public
webknife proxy --listen :8080 --upstream http://localhost:3000
webknife echo --listen :8080
webknife respond --listen :8080 --status 503 --body 'Service unavailable'
webknife redirect --listen :8080 --to https://example.com
```
