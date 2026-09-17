# Architecture

## Dependency direction

```
CLI (internal/adapters/cli)
  ↓
Application (orchestration)
  ↓
Core (pkg/webknife)

HTTP / filesystem / logging / TLS
        ↑
      adapters
```

**Rule:** Never import from `internal/` into `pkg/`. Never import CLI concerns into the core library.

## Package boundaries

| Package | Responsibility | Importable from |
|---|---|---|
| `cmd/webknife` | Entry point only | Nowhere |
| `internal/adapters/cli` | CLI parsing, command routing | `cmd/webknife` |
| `internal/build` | Version metadata | `internal/adapters/cli`, `cmd/webknife` |
| `internal/events` | Event types, bus, middleware | `pkg/webknife`, `internal/adapters/cli` |
| `pkg/webknife` | All HTTP functionality | Anyone (public API) |

## Extension points

New capabilities are added as:
1. A handler in `pkg/webknife/` (e.g., `delay.go`)
2. A CLI command parser in `internal/adapters/cli/cli.go`
3. A builder function in `internal/adapters/cli/cli.go` that wires the handler into the pipeline

The composition in `wrapWithMiddleware` handles auth and header middleware automatically for any handler.

## Adding a new handler

1. Create `pkg/webknife/newhandler.go` with a config struct and constructor
2. The constructor returns `http.Handler` — nothing more
3. Add a CLI parser function in `cli.go`
4. Add a `buildNewHandler` function that wires it through `wrapWithMiddleware`
5. Register the command in `ParseArgs` and `RunWithWriter`

## Events are optional

Not every handler needs to publish events. Events are for:
- Lifecycle facts (server started/stopped)
- Observable operations worth logging (request completed, auth failed)
- Cross-cutting concerns (headers modified)

Do not create events for internal implementation details.
