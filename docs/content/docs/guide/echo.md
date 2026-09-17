---
title: "Echo / Inspection"
weight: 3
---

# Echo / Request Inspection

Reflect incoming requests as structured JSON with full diagnostic detail.

## Basic usage

```bash
webknife echo --listen :8080
```

## What it returns

The echo handler returns a JSON object with:

- HTTP method
- Scheme (http/https)
- Host
- Path
- Query parameters
- Protocol version
- Content length
- Remote address
- All request headers
- Full URI
- TLS info (when applicable)
- Request body (text or binary)
- Body truncation indicator
- Cookies

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `--listen` | `:8080` | Listen address |
| `--max-body` | `1048576` (1MB) | Maximum body size for inspection |

## Examples

```bash
# Basic inspection
curl http://localhost:8080/

# POST with body
curl -X POST http://localhost:8080/api -d '{"key":"value"}'

# With query params
curl "http://localhost:8080/search?q=hello&page=1"
```

## Body handling

- Text bodies are included as-is
- Binary bodies are hex-encoded
- Bodies exceeding `--max-body` are truncated (indicator set to true)

## Testing

```bash
scripts/test-webknife.sh 03   # echo / request inspection
```
