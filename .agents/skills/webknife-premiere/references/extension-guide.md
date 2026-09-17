# Extension Guide

## Adding a new capability to Webknife

Follow this checklist when adding a new handler, middleware, or command.

### 1. Classification

Ask these questions first:

- **Is this a handler, middleware, or adapter concern?**
  - Handler: returns a response (echo, respond, redirect, static, proxy)
  - Middleware: transforms requests or responses (auth, headers)
  - Adapter: CLI parsing, logging, event publishing

- **Does it belong in the reusable library?**
  - If other Go programs could use it: `pkg/webknife/`
  - If it's CLI-specific: `internal/adapters/cli/`

- **Does it produce meaningful domain events?**
  - Only if the operation is observable and useful to log
  - Not for internal implementation details

### 2. Implementation

1. Create `pkg/webknife/newcapability.go`
2. Define `NewCapabilityConfig` struct
3. Implement `NewCapabilityHandler(cfg) http.Handler`
4. Keep the handler focused — no auth, no logging, no header manipulation inside it

### 3. CLI wiring

1. Add config fields to `cli.Config`
2. Create `parseNewCapability(args, stdout)` in `cli.go`
3. Create `buildNewCapabilityHandler(cfg, bus)` that calls `wrapWithMiddleware`
4. Register in `ParseArgs` and `RunWithWriter`
5. Update `printUsage`

### 4. Composition

The new capability automatically gets:
- Basic Authentication (if `--auth` is provided)
- Header manipulation (if header flags are provided)
- Request observation (events are published)

Verify this works in tests.

### 5. Tests

Add to `pkg/webknife/webknife_test.go`:
- Basic behavior test
- Edge cases (empty input, invalid config)
- Composition with auth
- Composition with header middleware

Add to `internal/adapters/cli/cli_test.go`:
- CLI parsing test
- Flag defaults
- Required flag validation

### 6. Documentation

- Update `README.md` with usage examples
- Update the skill reference if architectural patterns changed
- Verify examples compile/run

### 7. Future pressure

Consider how this capability interacts with planned features:
- Virtual hosts: does it need per-host configuration?
- Routing: does it need path matching?
- Delays/failures: does it compose with fault injection?

Do not implement these now, but do not design in a way that blocks them.

## Anti-patterns to avoid

- **Do not create a handler that internally embeds auth.** Auth is middleware.
- **Do not create events for every function call.** Events are for observable facts.
- **Do not add flags that only work with one command.** Common flags stay common.
- **Do not put CLI parsing logic in `pkg/webknife`.** That belongs in `internal/adapters/cli`.
- **Do not create interfaces for things with one implementation.** Wait for a real second use case.
