---
title: "Static Files"
weight: 1
---

# Static File Serving

Serve files from any directory with correct MIME types and directory handling.

## Basic usage

```bash
webknife serve --listen :8080 --root ./public
```

The simplest form serves the current directory:

```bash
webknife serve --listen :8080 --root .
```

## How it works

Webknife uses Go's standard `http.FileServer`. This means:

- Correct MIME type detection
- Directory listing support
- HEAD request support
- Path traversal protection (cannot escape the document root)

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `--listen` | `:8080` | Listen address |
| `--root` | `.` | Document root directory |
| `--auth` | | Basic auth credentials (`user:password`) |
| `--tls-cert` | | TLS certificate file |
| `--tls-key` | | TLS private key file |
| `--log-format` | `text` | Log format: `text` or `json` |

## Examples

```bash
# Serve with authentication
webknife serve --listen :8080 --root /var/www --auth admin:secret

# Serve with TLS
webknife serve --listen :443 --root ./public --tls-cert cert.pem --tls-key key.pem
```

## Testing

```bash
scripts/test-webknife.sh 01   # basic static serving
scripts/test-webknife.sh 02   # static with auth
```
