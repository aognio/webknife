---
title: "Proxy Development Service"
weight: 5
---

# Can I proxy this development service?

Expose a local development server with added headers, authentication, or logging.

## Basic proxy

```bash
webknife proxy --listen :8080 --upstream http://localhost:3000
```

## Add debug headers

```bash
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --set-request-header 'X-Forwarded-For: debug'
```

## Add authentication

```bash
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --auth admin:secret
```

## Remove upstream headers

```bash
webknife proxy --listen :8080 \
  --upstream http://localhost:3000 \
  --remove-response-header Server \
  --remove-response-header X-Powered-By
```

## Use cases

- Testing CORS behavior through a proxy
- Adding authentication to a development server
- Stripping headers from upstream responses
- Logging all proxied requests for debugging
