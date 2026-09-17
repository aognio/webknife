---
title: "Simulate Errors"
weight: 6
---

# How does my client behave when the server returns 503?

Simulate error responses to test client retry logic, error handling, and fallback behavior.

## Service unavailable

```bash
webknife respond --listen :8080 --status 503 --body "Service temporarily unavailable"
```

## Rate limiting

```bash
webknife respond --listen :8080 --status 429 \
  --body '{"error":"rate limited"}' \
  --content-type application/json \
  --set-response-header 'Retry-After: 60'
```

## Not found

```bash
webknife respond --listen :8080 --status 404 --body "Resource not found"
```

## Forbidden

```bash
webknife respond --listen :8080 --status 403 --body "Access denied"
```

## Testing client behavior

These simulations are useful for:

- Verifying retry logic in HTTP clients
- Testing circuit breaker implementations
- Validating error display in frontend applications
- Checking monitoring and alerting for error responses
- Drilling incident response procedures
