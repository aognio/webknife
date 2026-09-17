# Events

## Event types

All events implement the `Event` interface (empty interface). Current events:

| Event | When |
|---|---|
| `ServerStarted` | Server begins listening |
| `ServerStopped` | Server shuts down gracefully |
| `RequestReceived` | Request arrives (before handler) |
| `RequestCompleted` | Request finishes (after handler) |
| `RequestFailed` | Request errors |
| `ProxyRequestStarted` | Proxy forwards a request |
| `ProxyRequestCompleted` | Proxy receives upstream response |
| `AuthenticationSucceeded` | Auth passes |
| `AuthenticationFailed` | Auth fails |
| `TLSConnectionAccepted` | TLS handshake completes |
| `RequestInspected` | Echo handler processes request |
| `ResponseGenerated` | Respond handler generates response |
| `RedirectGenerated` | Redirect handler generates redirect |
| `RequestHeadersModified` | Request headers transformed |
| `ResponseHeadersModified` | Response headers transformed |

## What deserves an event

**Yes:** Lifecycle facts (server start/stop), authentication outcomes, operations that change observable state.

**No:** Internal implementation details, high-frequency operations without logging value, events that duplicate what `RequestCompleted` already provides.

## Publishing

Events are published via the `Publisher` interface:

```go
type Publisher interface {
    Publish(event Event)
}
```

The `EventBus` implementation fans out to subscribers. The `EventPublisher` in `pkg/webknife/publisher.go` wraps the bus for use in handlers.

## Consumption

The `LogAdapter` subscribes to events and logs them. Other consumers can subscribe via `bus.Subscribe(fn)`.

## Event naming

Events use past tense: `RequestCompleted`, not `RequestComplete` or `OnRequest`. They describe what happened, not what should happen.

## Rule

> Events describe facts that already happened, not commands disguised as events.
