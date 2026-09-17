---
title: "Middleware"
weight: 4
---

# Middleware

Every Webknife command follows the same middleware pipeline.

## Pipeline

```
HTTP Listener
  │
  ▼
Request Observation (observe)
  │
  ▼
Basic Authentication (optional)
  │
  ▼
Request Header Transformation (optional)
  │
  ▼
Handler (static / proxy / echo / respond / redirect)
  │
  ▼
Response Header Transformation (optional)
  │
  ▼
Response Observation (events)
```

## Composition

Middleware composes via `http.Handler`. Each middleware wraps the next handler:

```go
handler = observe.Middleware(pub)(handler)
handler = auth.Basic(cfg, handler)
handler = headers.Middleware(ops, dir)(handler)
```

## Adding custom middleware

Webknife middleware follows the standard Go pattern:

```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // do something before
        next.ServeHTTP(w, r)
        // do something after
    })
}
```

This composes with any Webknife command without modification.
