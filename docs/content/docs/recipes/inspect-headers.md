---
title: "Inspect Headers"
weight: 7
---

# What headers is my client actually sending?

Use echo to see the exact headers your HTTP client sends.

## Start echo

```bash
webknife echo --listen :8080
```

## Test from your client

```bash
curl -v http://localhost:8080/
```

The JSON output includes every header:

```json
{
  "headers": {
    "Accept": ["*/*"],
    "User-Agent": ["curl/7.68.0"],
    "Host": ["localhost:8080"]
  }
}
```

## Common debugging scenarios

### Debug proxy headers

```bash
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Debug: true'
curl -v http://localhost:8080/
# X-Debug appears in echo output
```

### Debug authentication headers

```bash
curl -u user:pass http://localhost:8080/
# Authorization header appears in echo output
```

### Debug cookie handling

```bash
curl -b "session=abc123" http://localhost:8080/
# Cookie header appears in echo output
```
