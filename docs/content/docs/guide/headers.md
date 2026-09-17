---
title: "Headers"
weight: 8
---

# Header Manipulation

Set, add, or remove request and response headers on any command.

## Basic usage

```bash
webknife echo --listen :8080 --set-request-header 'X-Custom: value'
```

## Operations

| Operation | Description |
|-----------|-------------|
| `--set-request-header` | Set or replace a request header |
| `--add-request-header` | Add a value to a request header |
| `--remove-request-header` | Remove a request header |
| `--set-response-header` | Set or replace a response header |
| `--add-response-header` | Add a value to a response header |
| `--remove-response-header` | Remove a response header |

## Format

Headers use `Name: Value` format:

```bash
--set-request-header 'X-Debug: true'
--remove-response-header Server
```

## Examples

```bash
# Inject a debug header into proxied requests
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Debug: true'

# Remove Server header from responses
webknife echo --listen :8080 \
  --set-response-header 'X-Debug: true' \
  --remove-response-header Server

# Multiple headers
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Request-Source: webknife' \
  --set-request-header 'X-Custom: value' \
  --remove-response-header Server
```

## Composability

Header manipulation composes with every command. Apply it to:

- `serve` — modify headers on static file responses
- `proxy` — inject or strip headers on proxied traffic
- `echo` — see how header injection looks in inspection output
- `respond` — add headers to custom responses
- `redirect` — add headers to redirect responses

## Testing

```bash
scripts/test-webknife.sh 07   # request header manipulation
scripts/test-webknife.sh 08   # response header manipulation
```
