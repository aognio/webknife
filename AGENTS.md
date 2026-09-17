# AGENTS.md

## Build

```bash
make build     # build to bin/webknife
make test      # run all tests
make check     # fmt-check + vet + test
make help      # show all targets
```

## Run

```bash
./bin/webknife serve --listen :8080 --root .
./bin/webknife proxy --listen :8080 --upstream http://localhost:3000
./bin/webknife echo --listen :8080
./bin/webknife respond --listen :8080 --status 503
./bin/webknife redirect --listen :8080 --to https://example.com
```

## Project layout

- `cmd/webknife/` — composition root (wires CLI → application)
- `internal/adapters/cli/` — CLI parsing and validation (thin)
- `internal/adapters/logging/` — log adapters (text, JSON)
- `internal/application/` — orchestration layer (builds handlers, wires middleware)
- `internal/archtest/` — architectural dependency enforcement tests
- `internal/build/` — version metadata
- `internal/events/` — event types and event bus implementation
- `pkg/webknife/` — ports and shared types (Middleware, Publisher)
- `pkg/auth/` — Basic Auth middleware
- `pkg/echo/` — request inspection handler
- `pkg/headers/` — header manipulation middleware
- `pkg/observe/` — request observation middleware
- `pkg/proxy/` — reverse proxy handler
- `pkg/redact/` — secret detection and masking
- `pkg/redirect/` — redirect handler
- `pkg/respond/` — arbitrary response handler
- `pkg/server/` — HTTP server lifecycle
- `pkg/static/` — static file handler

## Conventions

- Independent features: each `pkg/*` package is standalone, no sibling imports
- Standard library first: no external dependencies unless they provide substantial value
- Composition: everything composes via `http.Handler`
- Ports in pkg/webknife: Publisher interface and Middleware type are shared contracts
- Architectural tests: run `go test ./internal/archtest/` to verify dependency rules
- Tests: use `httptest`, no external Internet, deterministic
