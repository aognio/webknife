---
title: "Test Port 80"
weight: 1
---

# Is port 80 actually reaching this VPS?

Verify that a remote server's port 80 is reachable from the outside.

## On the server

```bash
sudo webknife respond --listen :80 --status 200 --body "reachable"
```

> Ports below 1024 require elevated permissions on most systems.

## From the client

```bash
curl http://your-server-ip/
# Expected: reachable
```

## What this tells you

- Port forwarding from the internet reaches the server
- No firewall is blocking port 80
- The server is listening and responding

If the curl times out, the issue is in your firewall, NAT, or cloud security group — not in your application.
