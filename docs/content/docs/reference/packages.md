---
title: "Packages"
weight: 4
---

# Packages

Webknife's Go packages and their responsibilities.

## Ports

| Package | Description |
|---------|-------------|
| `pkg/webknife` | Shared types: `Publisher`, `Middleware`, `Event`, `Chain` |

## Feature packages

| Package | Description |
|---------|-------------|
| `pkg/static` | Static file server |
| `pkg/proxy` | Reverse proxy |
| `pkg/echo` | Request inspection (JSON) |
| `pkg/respond` | Arbitrary HTTP responses |
| `pkg/redirect` | HTTP redirects |
| `pkg/redact` | Secret detection and masking |

## Middleware packages

| Package | Description |
|---------|-------------|
| `pkg/auth` | Basic Authentication |
| `pkg/headers` | Header manipulation |
| `pkg/observe` | Request observation (publishes events) |

## Infrastructure packages

| Package | Description |
|---------|-------------|
| `pkg/server` | HTTP server lifecycle |
| `internal/application` | Orchestration layer |
| `internal/adapters/cli` | CLI parsing |
| `internal/adapters/logging` | Log adapters (text, JSON) |
| `internal/events` | Event bus implementation |
| `internal/build` | Version metadata |
| `internal/archtest` | Dependency enforcement tests |

## Import rules

- Feature packages never import siblings
- Feature packages never import `internal/*`
- `pkg/webknife` never imports `internal/*`
- CLI never imports feature packages directly
