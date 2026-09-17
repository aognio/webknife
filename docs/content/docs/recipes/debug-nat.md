---
title: "Debug NAT"
weight: 2
---

# Is my NAT forwarding correct?

Verify that a router or cloud load balancer is forwarding traffic to the right internal server.

## Setup

On the internal server:

```bash
webknife echo --listen :8080
```

## Test from outside

```bash
curl http://your-public-ip:8080/
```

The echo response shows exactly what the request looks like after NAT:

- Source IP (the NAT'd address)
- Host header
- Any modified headers

This is useful when:

- Debugging port forwarding rules
- Testing cloud load balancer configuration
- Verifying Docker port mappings
- Diagnosing hairpin NAT issues
