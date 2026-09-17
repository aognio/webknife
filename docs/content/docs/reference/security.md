---
title: "Security"
weight: 2
---

# Security

Webknife is a diagnostic tool. Convenience-oriented defaults do not make it a hardened production server.

## Binding to public interfaces

By default, `--listen :8080` binds to all interfaces. Use `--listen 127.0.0.1:8080` to restrict to localhost.

## Privileged ports

Ports below 1024 require elevated permissions on most systems.

## Basic Authentication

The `--auth` flag passes credentials as a command-line argument. These are visible in:

- Process listings (`ps aux`)
- Shell history

For sensitive use, prefer TLS and environment variables in scripts.

## TLS

Use TLS for any traffic containing sensitive data. The server enforces TLS 1.2 minimum.

## Filesystem exposure

`webknife serve --root /` exposes the entire filesystem. Be deliberate about the root directory.

## Proxy exposure

`webknife proxy` makes the upstream accessible through the proxy. Anyone who can reach the proxy can reach the upstream.

## Sensitive headers

Webknife identifies `Authorization`, `Cookie`, `Set-Cookie`, and `Proxy-Authorization` as sensitive.

## Request body limits

The echo handler enforces a configurable maximum body size (default 1MB) to prevent memory exhaustion.

## Diagnostic defaults

Logging includes remote addresses, request paths, and status codes. Do not enable verbose logging in production environments with untrusted traffic.

## Temporary servers

Do not leave Webknife processes running unintentionally. They are diagnostic tools, not production services.
