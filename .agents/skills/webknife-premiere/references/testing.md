# Testing

## Principles

- Tests are part of the iteration, not afterthoughts
- Use `httptest` for all HTTP testing
- No external Internet access required
- Deterministic: no time-dependent tests, no random ports
- Fast: unit tests should complete in milliseconds

## Patterns

### Handler test

```go
func TestHandler_Behavior(t *testing.T) {
    handler := webknife.NewSomeHandler(webknife.SomeConfig{...})
    req := httptest.NewRequest("GET", "/path", nil)
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    if w.Code != expected {
        t.Fatalf("expected %d, got %d", expected, w.Code)
    }
}
```

### Middleware test

```go
func TestMiddleware_ModifiesRequest(t *testing.T) {
    inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // assert on modified request
    })
    handler := webknife.SomeMiddleware(cfg)(inner)
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
}
```

### Proxy test

Start a test upstream server, then proxy to it:

```go
func TestProxy_ForwardsRequest(t *testing.T) {
    upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("upstream"))
    }))
    defer upstream.Close()

    handler, _ := webknife.NewProxyHandler(webknife.ProxyConfig{Upstream: upstream.URL})
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    // assert on w
}
```

### CLI test

```go
func TestParseArgs_Command(t *testing.T) {
    cmd, cfg, err := cli.ParseArgs([]string{"webknife", "command", "--flag", "value"}, stdout, stderr)
    if err != nil { t.Fatal(err) }
    if cmd != "command" { t.Fatalf(...) }
}
```

## Test file placement

- `pkg/webknife/webknife_test.go` — all handler/middleware tests
- `internal/adapters/cli/cli_test.go` — CLI parsing tests
- `internal/events/events_test.go` — event bus tests

## What to test

- Static file serving, missing files, HEAD requests, path traversal
- Basic Auth success, failure, missing credentials
- Proxy forwarding, response propagation, header forwarding
- Echo: GET, POST, query params, headers, cookies, body truncation, binary bodies
- Respond: status codes, custom headers, body from file
- Redirect: status codes, destination, path preservation
- Headers: set, add, remove for request and response
- Events: publishing, subscribing
- CLI: all command parsing, flag defaults, required flags

## Coverage

```bash
make test
make test-race
make cover   # opens HTML coverage report
```
