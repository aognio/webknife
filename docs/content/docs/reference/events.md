---
title: "Events"
weight: 3
---

# Events

Webknife uses an internal event bus for request lifecycle observability.

## Event types

| Event | Description |
|-------|-------------|
| `RequestReceived` | A request was received |
| `RequestCompleted` | A request was completed |
| `TLSConnectionAccepted` | A TLS connection was accepted |
| `ServerStarted` | The server started listening |
| `ServerStopped` | The server stopped |

## Publishing

Feature packages publish events through the `Publisher` port:

```go
pub.Publish(observe.RequestCompleted{
    Method:     r.Method,
    Path:       r.URL.Path,
    StatusCode: 200,
    Size:       1024,
    Duration:   1.234ms,
    RemoteAddr: r.RemoteAddr,
    Timestamp:  time.Now(),
})
```

## Subscribing

The application layer subscribes to events and routes them to log adapters:

```go
bus.Subscribe(func(event events.Event) {
    switch e := event.(type) {
    case observe.RequestCompleted:
        logHandler.Handle(logging.RequestCompleted{...})
    }
})
```

## Custom integrations

To add custom event handling, subscribe to the event bus in `internal/application/app.go`.
