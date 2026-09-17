---
title: "Reverse Proxy"
weight: 2
---

# Reverse Proxy

Proxy requests to an upstream service, forwarding method, path, headers, and body.

## Basic usage

```bash
webknife proxy --listen :8080 --upstream http://localhost:3000
```

## What it forwards

- HTTP method
- Path and query string
- Headers
- Request body
- Response status, headers, and body

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `--listen` | `:8080` | Listen address |
| `--upstream` | | **Required.** Upstream URL |
| `--auth` | | Basic auth credentials |
| `--set-request-header` | | Set request header (`Name: Value`) |
| `--remove-response-header` | | Remove response header |

## Examples

```bash
# Proxy with debug headers
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Debug: true' \
  --remove-response-header Server

# Proxy with authentication
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --auth admin:secret
```

## Security considerations

`webknife proxy` makes the upstream accessible through the proxy. Anyone who can reach the proxy can reach the upstream.

## Testing

```bash
scripts/test-webknife.sh 09   # reverse proxy (requires upstream on port 9000)
```
