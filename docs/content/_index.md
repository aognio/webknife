---
title: "Documentation"
type: docs
---

<p align="center">
  <img src="/images/webknife-logo.png" alt="Webknife" width="300">
</p>

# Webknife

> A single-binary, CLI-first HTTP laboratory for diagnostics, testing, proxying, serving, inspection, and experimentation.

Webknife is a disposable HTTP tool you keep in your terminal. It starts in seconds, requires no configuration files, and does one thing well.

```bash
# Serve files
webknife serve --listen :8080 --root .

# Inspect requests
webknife echo --listen :8080

# Reverse proxy
webknife proxy --listen :8080 --upstream http://localhost:3000
```

---

## [Quick Start](getting-started/)

Get from zero to useful in a few minutes.

## [Feature Guide](guide/)

Deep dive into every capability.

## [Recipes](recipes/)

Problem-oriented examples for common tasks.

## [CLI Reference](cli/)

Complete command and flag documentation.

## [Architecture](architecture/)

How Webknife is designed and built.

## [Development](development/)

Building, testing, contributing, and agent development.
