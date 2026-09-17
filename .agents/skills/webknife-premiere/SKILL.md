---
name: webknife-premiere
description: "Comprehensive guidance for working on Webknife — a single-binary, CLI-first HTTP diagnostic laboratory. Use when modifying, extending, or debugging Webknife; when adding new commands, handlers, or middleware; when reviewing architecture or test coverage; or when onboarding to the project."
user-invocable: true
license: MIT
metadata:
  author: webknife
  version: "0.3.0"
  openclaw:
    emoji: "🔪"
    requires:
      bins: [go, git, make]
---

**Persona:** You are a contributor to Webknife. You preserve the project's architectural invariants, compose via `http.Handler`, prefer standard library, and avoid speculative abstractions.

# Webknife

A single-binary, CLI-first HTTP laboratory for diagnostics, testing, proxying, serving, inspection, and experimentation.

## When to use this skill

- Modifying or extending any Webknife code
- Adding new CLI commands, handlers, or middleware
- Reviewing architecture or test coverage
- Debugging issues in existing functionality
- Onboarding to the project for the first time

## Core rules

1. **Independent features.** Each `pkg/*` package is standalone. Never add imports between sibling feature packages.
2. **Standard library first.** Go's HTTP stack is excellent. Do not add external dependencies for small amounts of straightforward code.
3. **Composition over inheritance.** Everything composes via `http.Handler`. No inheritance, no god objects.
4. **Ports live in pkg/webknife.** The `Publisher` interface and `Middleware` type are the only shared contracts.
5. **Preserve dependency direction.** cmd/ → internal/application/ → pkg/*. Never reverse this. Never import internal/ from pkg/.
6. **No speculative abstractions.** Do not create interfaces, repositories, factories, or service layers unless they solve a current problem.
7. **Manual tests are separate.** Use `go test ./...` for automated correctness. Use `scripts/test-webknife.sh <scenario>` for manual HTTP experimentation.

## Project layout

```
cmd/webknife/                  Composition root (wires CLI → application)
internal/
  adapters/
    cli/                       CLI parsing and validation (thin)
    logging/                   Log adapters (text, JSON)
  application/                 Orchestration layer (builds handlers, wires middleware)
  archtest/                    Architectural dependency enforcement tests
  build/                       Version and build metadata
  events/                      Event types and event bus implementation
pkg/
  webknife/                    Ports and shared types (Middleware, Publisher)
  auth/                        Basic Auth middleware
  echo/                        Request inspection handler
  headers/                     Header manipulation middleware
  observe/                     Request observation middleware
  proxy/                       Reverse proxy handler
  redact/                      Secret detection and masking
  redirect/                    Redirect handler
  respond/                     Arbitrary response handler
  server/                      HTTP server lifecycle
  static/                      Static file handler
```

## HTTP pipeline

Every command follows this composition model:

```
HTTP Listener
  ↓
Request Observation (pkg/observe)
  ↓
Basic Authentication (optional, pkg/auth)
  ↓
Request Header Transformation (optional, pkg/headers)
  ↓
Handler (echo / respond / redirect / static / proxy)
  ↓
Response Header Transformation (optional, pkg/headers)
  ↓
Response Observation (events)
```

This pipeline is assembled in `internal/application/app.go` → `buildHandler()`.

## Adding a new feature

1. Create `pkg/newfeature/newfeature.go` — handler or middleware, no sibling imports
2. Add event types to `internal/events/events.go` if needed
3. Wire in `internal/application/app.go` → `buildHandler()` switch
4. Add CLI flags in `internal/adapters/cli/cli.go`
5. Add tests in `pkg/newfeature/newfeature_test.go`
6. Run `go test ./...` and `go test ./internal/archtest/`

## References

- [architecture.md](references/architecture.md) — dependency direction, package boundaries, extension points
- [cli.md](references/cli.md) — command naming, flag conventions, adding new commands
- [events.md](references/events.md) — event types, publishing, consumption, what deserves an event
- [http-pipeline.md](references/http-pipeline.md) — middleware composition, handler patterns, request lifecycle
- [testing.md](references/testing.md) — test patterns, httptest, deterministic tests
- [security.md](references/security.md) — credential redaction, path traversal, proxy safety
- [conventions.md](references/conventions.md) — code style, naming, Go idioms
- [extension-guide.md](references/extension-guide.md) — how to add a new capability
