---
title: "Overview"
weight: 1
---

# Architecture Overview

Webknife follows a layered architecture with clear dependency direction and independent feature packages.

## Package layout

```
cmd/webknife/                  Composition root
internal/
  adapters/
    cli/                       CLI parsing and validation
    logging/                   Log adapters (text, JSON)
  application/                 Orchestration layer
  archtest/                    Dependency enforcement tests
  build/                       Version metadata
  events/                      Event bus implementation
pkg/
  webknife/                    Ports (Publisher, Middleware)
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

## Dependency direction

```
cmd/webknife/ (composition root)
    │
    ▼
internal/application/ (wires everything)
    │
    ├──► internal/adapters/cli/    (CLI parsing)
    ├──► internal/adapters/logging/ (log output)
    │
    └──► pkg/ (feature packages)
         ├──► pkg/webknife/        (ports)
         ├──► pkg/observe/
         ├──► pkg/auth/
         ├──► pkg/headers/
         ├──► pkg/echo/
         ├──► pkg/respond/
         ├──► pkg/redirect/
         ├──► pkg/proxy/
         ├──► pkg/static/
         └──► pkg/server/
```

## Architectural rules

1. Feature packages never import each other (except `pkg/webknife` for ports)
2. Feature packages never import `internal/*`
3. `pkg/webknife` never imports `internal/*`
4. CLI adapter never imports feature packages directly
