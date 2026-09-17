# HTTP Pipeline

## Composition model

Every request flows through this pipeline:

```
HTTP Listener
  ↓
events.RequestMiddleware (observation)
  ↓
BasicAuth (optional)
  ↓
Request Header Middleware (optional)
  ↓
Handler
  ↓
Response Header Middleware (optional)
  ↓
(events complete)
```

## How it's assembled

In `cli.go`, `wrapWithMiddleware` takes a handler and layers middleware on top:

```go
func wrapWithMiddleware(handler http.Handler, cfg *Config, bus *events.EventBus) (http.Handler, error) {
    // 1. Apply response header middleware (outermost, wraps next)
    // 2. Apply request header middleware
    // 3. Apply Basic Auth if configured
    // 4. Apply request observation middleware
    return handler, nil
}
```

The order matters:
- **Request observation** runs first (outermost) to capture the full lifecycle
- **Auth** runs before the handler to reject unauthorized requests
- **Header middleware** runs between auth and handler
- **Response headers** are applied just before writing

## Handler pattern

Every handler in `pkg/webknife` follows the same shape:

```go
type SomeConfig struct { /* options */ }

func NewSomeHandler(cfg SomeConfig) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // implementation
    })
}
```

Handlers return `http.Handler`. They do not publish events directly — the observation middleware handles that.

## Middleware pattern

Middleware in `pkg/webknife` is a function:

```go
func SomeMiddleware(cfg SomeConfig) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // before
            next.ServeHTTP(w, r)
            // after
        })
    }
}
```

Or as a simple function:

```go
func BasicAuth(cfg BasicAuthConfig, next http.Handler) http.Handler { ... }
```

## Header manipulation

The `HeaderMiddleware` wraps the `ResponseWriter` to intercept `WriteHeader` and `Write` calls, applying response header operations at the right moment.

Request header operations modify `r.Header` directly before passing to the next handler.

## Adding a new handler

1. Create `pkg/webknife/newhandler.go`
2. Define config struct and constructor returning `http.Handler`
3. Add to `wrapWithMiddleware` if it needs different pipeline treatment
4. Or just wire it through the existing pipeline in the CLI builder
