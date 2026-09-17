---
title: "Composition Root"
weight: 3
---

# Composition Root

The top-level binary is the composition root. It wires independent components into a working application.

## How it works

```
cmd/webknife/main.go
    │
    ├── cli.Parse(args)        → Config
    │
    └── application.Run(cfg)
            │
            ├── build handler   (from feature packages)
            ├── wrap middleware (auth, headers, observe)
            └── start server
```

## The principle

> Feature packages provide behavior. The top-level binary provides composition.

This means:

- `pkg/static` knows how to serve files, but not how to start a server
- `pkg/auth` knows how to check credentials, but not which routes to protect
- `internal/application` decides which handler to build and which middleware to wrap
- `cmd/webknife` decides when to call the application layer

## Benefits

- **Testing** — each layer can be tested independently
- **Library reuse** — import `pkg/echo` without the CLI or server
- **Agent development** — predictable package boundaries for coding agents
- **Extension** — adding a new feature never couples to existing ones
