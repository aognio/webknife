---
title: "Agent Development"
weight: 4
---

# Agent Development

Webknife is designed to be agent-friendly. Coding agents can discover and use the project through its structured conventions.

## Predictable boundaries

Each `pkg/*` package is independent. An agent can:

- Read `pkg/echo/` and understand request inspection without context from other packages
- Import `pkg/auth` without pulling in unrelated dependencies
- Add a new feature package without coupling to existing ones

## Progressive disclosure

Start with:

```bash
scripts/test-webknife.sh help
```

Then explore the source code. The package layout is self-documenting.

## Deterministic tests

```bash
go test ./...
```

Tests are fast, deterministic, and require no external services.

## Manual test scenarios

```bash
scripts/test-webknife.sh <scenario>
```

Numbered scenarios provide a stable reference for documentation, bug reports, and agent instructions.

## Small APIs

Each package exposes a minimal API:

- `static.New(cfg)` → `http.Handler`
- `echo.New(cfg)` → `http.Handler`
- `auth.Basic(cfg, next)` → `http.Handler`

No god objects, no complex initialization, no hidden coupling.

## Architecture enforcement

```bash
go test ./internal/archtest/
```

Automated tests verify dependency rules. An agent cannot accidentally introduce a cross-package import.
