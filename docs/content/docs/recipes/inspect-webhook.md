---
title: "Inspect Webhooks"
weight: 3
---

# What exactly is this webhook sending?

Set up a temporary echo server to capture incoming webhook payloads.

## Start the inspector

```bash
webknife echo --listen :8080
```

## Point your webhook at it

Update your webhook URL to `http://your-server:8080/`.

## Inspect the payload

The echo output shows:

- Full request body (JSON, form data, etc.)
- All headers the webhook sender set
- Content-Type
- Signature headers (if present)
- Query parameters

This is invaluable when:

- Integrating with third-party services (GitHub, Stripe, Slack)
- Debugging webhook format assumptions
- Verifying signature validation
- Understanding content types being sent
