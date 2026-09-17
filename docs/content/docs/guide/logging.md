---
title: "Logging"
weight: 9
---

# Logging and Observability

Webknife logs every request with concise human-readable or structured JSON output.

## Log formats

### Text (default)

```bash
webknife echo --listen :8080
```

Output:

```
2026-09-17T00:00:36-05:00 GET /test 200 1.2KB 1.234ms 127.0.0.1:54321
```

### JSON

```bash
webknife echo --listen :8080 --log-format json
```

Output:

```json
{"timestamp":"2026-09-17T00:00:36-05:00","method":"GET","path":"/test","status":200,"size":"1.2KB","duration":"1.234ms","remote":"127.0.0.1:54321"}
```

## What gets logged

Every request produces a log entry with:

- Timestamp
- HTTP method
- Request path
- Response status code
- Response size (human-readable)
- Request duration
- Remote client address

## Event-driven observability

Webknife uses an internal event bus. Request lifecycle events are published through a `Publisher` interface. This enables:

- Custom log adapters
- Metrics collection
- Integration with observability platforms

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `--log-format` | `text` | `text` for human-readable, `json` for structured |
