---
title: "Custom Responses"
weight: 4
---

# Custom Responses

Return an arbitrary HTTP status code with custom body, Content-Type, and headers.

## Basic usage

```bash
webknife respond --listen :8080 --status 200 --body "All systems go"
```

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `--listen` | `:8080` | Listen address |
| `--status` | `200` | HTTP status code |
| `--body` | | Response body |
| `--body-file` | | Response body from file |
| `--content-type` | | Content-Type header |
| `--set-response-header` | | Set response header (`Name: Value`) |

## Examples

```bash
# Simulate service unavailable
webknife respond --listen :8080 --status 503 --body "Service unavailable"

# Simulate rate limiting
webknife respond --listen :8080 --status 429 --body '{"error":"rate limited"}' \
  --content-type application/json \
  --set-response-header 'Retry-After: 60'

# Serve a static response file
webknife respond --listen :8080 --status 200 --body-file response.json \
  --content-type application/json
```

## Testing

```bash
scripts/test-webknife.sh 04   # custom response (200)
scripts/test-webknife.sh 05   # error response (503)
```
