---
title: "TLS"
weight: 7
---

# TLS

Serve over HTTPS with certificate and private key files.

## Basic usage

```bash
webknife serve --listen :443 --root . \
  --tls-cert cert.pem --tls-key key.pem
```

## How it works

When `--tls-cert` and `--tls-key` are both provided, Webknife starts an HTTPS server instead of HTTP.

## Options

| Flag | Description |
|------|-------------|
| `--tls-cert` | Path to TLS certificate file (PEM) |
| `--tls-key` | Path to TLS private key file (PEM) |

## Examples

```bash
# Serve static files over HTTPS
webknife serve --listen :443 --root ./public \
  --tls-cert /etc/ssl/certs/webknife.pem \
  --tls-key /etc/ssl/private/webknife.key

# Proxy over HTTPS
webknife proxy --listen :443 --upstream http://localhost:3000 \
  --tls-cert cert.pem --tls-key key.pem
```

## TLS defaults

- Minimum version: TLS 1.2
- Cipher suites: ECDHE with AES-GCM

## Security considerations

- Use TLS for any traffic containing sensitive data
- The server enforces TLS 1.2 minimum — older clients will not connect
- Certificate files should have restricted permissions
