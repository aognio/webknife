# Conventions

## Go style

- Standard `gofmt` / `gofumpt` formatting
- No comments unless they explain why, not what
- Short variable names in small scopes, descriptive names at package boundary
- Error messages: lowercase, no punctuation, no `Error:` prefix
- Return errors, don't panic
- Use `fmt.Errorf` with `%w` for wrapping

## Naming

- Config structs: `XxxConfig` (e.g., `EchoConfig`, `ProxyConfig`)
- Constructors: `NewXxxHandler` returning `http.Handler`
- Middleware: `XxxMiddleware` or `BasicAuth` returning `func(http.Handler) http.Handler`
- Events: past tense nouns (`RequestCompleted`, not `OnRequest`)
- CLI parsers: `parseXxx(args, stdout)` returning `(string, *Config, error)`

## Files

- One concern per file: `auth.go`, `echo.go`, `proxy.go`
- Tests alongside source: `foo_test.go` in the same package
- No utility packages — keep helpers close to where they're used

## Dependencies

- Prefer Go standard library
- `net/http`, `net/http/httputil`, `net/http/httptest` are core
- `encoding/json` for echo output
- `crypto/subtle` for constant-time comparison
- Do not add dependencies for small amounts of code

## Error handling

- Errors propagate up, not logged at the point of failure
- CLI prints errors to stderr
- Handlers return HTTP error responses, not Go errors
- The event bus logs failures via the log adapter

## Documentation

- README documents user-facing behavior
- Skill documents agent-facing conventions
- Code is the source of truth when docs drift
