---
title: "Quick Start"
weight: 3
---

# Quick Start

Get from zero to useful in five minutes.

## Build

```bash
git clone https://github.com/aognio/webknife
cd webknife
make build
```

## Serve files

```bash
webknife serve --listen :8080 --root .
```

Open `http://localhost:8080` in your browser.

## Inspect a request

In another terminal:

```bash
webknife echo --listen :9090
curl http://localhost:9090/test?foo=bar
```

You get structured JSON showing the full request details.

## Reverse proxy

```bash
webknife proxy --listen :8080 --upstream http://localhost:3000
```

All requests to `:8080` are forwarded to your upstream service.

## Add Basic Authentication

```bash
webknife serve --listen :8080 --root . --auth admin:secret
```

```bash
curl http://localhost:8080/          # 401
curl -u admin:secret http://localhost:8080/  # 200
```

## Simulate errors

```bash
webknife respond --listen :8080 --status 503 --body "Service unavailable"
```

## Check logs

Webknife logs every request to stdout. Use `--log-format json` for structured output:

```bash
webknife echo --listen :8080 --log-format json
```
