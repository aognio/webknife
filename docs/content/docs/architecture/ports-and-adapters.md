---
title: "Ports and Adapters"
weight: 2
---

# Ports and Adapters

Webknife uses the Ports and Adapters (Hexagonal) architecture pattern.

## Ports

Ports are interfaces that define how feature packages communicate with the outside world. They live in `pkg/webknife/`:

```go
// Publisher is the port for publishing events.
type Publisher interface {
    Publish(event Event)
}

// Middleware wraps an http.Handler.
type Middleware func(http.Handler) http.Handler
```

## Adapters

Adapters implement ports for specific infrastructure:

| Port | Adapter | Location |
|------|---------|----------|
| `Publisher` | `eventPublisher` | `internal/application/` |
| Log handler | `TextHandler`, `JSONHandler` | `internal/adapters/logging/` |
| CLI parser | `cli.Parse` | `internal/adapters/cli/` |

## Why this matters

- Feature packages accept only ports, not concrete implementations
- The application layer wires adapters to ports
- You can swap adapters without changing feature code
- Testing is straightforward — inject a mock or stub adapter
