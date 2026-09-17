---
title: "Introduction"
weight: 1
---

# Introduction

Webknife is a single-binary, CLI-first HTTP laboratory. It is a diagnostic tool for developers, sysadmins, and SREs who need a quick, disposable HTTP server that does one thing well.

## What Webknife is

- A **static file server** you can start in one command
- A **reverse proxy** for inspecting or modifying traffic
- A **request inspector** that shows exactly what clients send
- A **response simulator** for testing error paths
- A **redirect tester** for verifying URL redirects
- A **header manipulation** tool for debugging HTTP headers
- A **Basic Auth** protector for temporary servers

## What Webknife is not

Webknife is not Nginx, Apache, Caddy, HAProxy, or mitmproxy. It is not a production web server. It is a Swiss Army knife you keep in your terminal for moments when you need to quickly diagnose, test, or prototype HTTP interactions.

## Who uses Webknife

- **Developers** inspecting webhook payloads or API requests
- **Sysadmins** verifying port forwarding or NAT configuration
- **SRE/DevOps** simulating error responses during incident drills
- **Go developers** embedding HTTP handlers in their own tools
- **Coding agents** automating HTTP testing workflows
