---
title: "Configuration"
weight: 1
---

# Configuration

Webknife uses command-line flags for configuration. There are no config files.

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WEBKNIFE` | `./bin/webknife` | Binary path (for test harness) |
| `HOST` | `127.0.0.1` | Listen host (for test harness) |
| `PORT` | `8080` | Listen port (for test harness) |

## Defaults

| Flag | Default | Description |
|------|---------|-------------|
| `--listen` | `:8080` | Listen address |
| `--log-format` | `text` | Log format |
| `--max-body` | `1048576` | Max body size for echo (bytes) |
| `--status` | `200` | Status code for respond |
| `--status` | `302` | Status code for redirect |

## Override pattern

```bash
PORT=9090 scripts/test-webknife.sh 03
WEBKNIFE=/usr/local/bin/webknife scripts/test-webknife.sh 01
```
