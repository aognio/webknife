---
title: "Authentication"
weight: 6
---

# Basic Authentication

Protect any Webknife command with HTTP Basic Authentication.

## Basic usage

```bash
webknife serve --listen :8080 --root . --auth admin:secret
```

## How it works

The `--auth` flag accepts credentials in `user:password` format. It composes with any command — serve, proxy, echo, respond, or redirect.

## Behavior

- Requests without credentials receive `401 Unauthorized`
- The `WWW-Authenticate` header is set with the configured realm
- Credentials are compared using constant-time comparison (timing-attack safe)

## Examples

```bash
# Protect a static server
webknife serve --listen :8080 --root /var/www --auth admin:secret

# Protect a proxy
webknife proxy --listen :8080 --upstream http://localhost:3000 --auth admin:secret

# Protect echo (inspect only authenticated requests)
webknife echo --listen :8080 --auth admin:secret
```

## Security considerations

- Credentials are visible in process listings (`ps aux`) and shell history
- For sensitive use, prefer TLS and environment variables in scripts
- Basic Authentication sends credentials base64-encoded (not encrypted) — always use TLS in production

## Testing

```bash
scripts/test-webknife.sh 02   # static server with Basic Auth
```
