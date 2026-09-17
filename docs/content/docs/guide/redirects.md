---
title: "Redirects"
weight: 5
---

# Redirects

Redirect to any URL with configurable status codes.

## Basic usage

```bash
webknife redirect --listen :8080 --to https://example.com
```

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `--listen` | `:8080` | Listen address |
| `--to` | | **Required.** Target URL |
| `--status` | `302` | Status code (301, 302, 303, 307, 308) |
| `--preserve-path` | `false` | Preserve original path in redirect |

## Status codes

| Code | Meaning |
|------|---------|
| 301 | Moved Permanently |
| 302 | Found (temporary) |
| 303 | See Other |
| 307 | Temporary Redirect |
| 308 | Permanent Redirect |

## Examples

```bash
# Permanent redirect
webknife redirect --listen :8080 --to https://example.com --status 301

# Preserve path
webknife redirect --listen :8080 --to https://example.com --preserve-path
# /foo?a=1 -> https://example.com/foo?a=1
```

## Testing

```bash
scripts/test-webknife.sh 06   # HTTP redirect
```
