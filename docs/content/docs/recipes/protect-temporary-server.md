---
title: "Protect Temp Server"
weight: 4
---

# Can I quickly protect this temporary directory?

Serve a directory over HTTP with Basic Authentication in one command.

```bash
webknife serve --listen :8080 --root /tmp/shared-files --auth admin:secret
```

Now only authenticated users can access the files:

```bash
# Without auth — 401
curl http://your-server:8080/

# With auth — 200
curl -u admin:secret http://your-server:8080/
```

## When to use this

- Sharing files temporarily with a colleague
- Exposing build artifacts during a deployment
- Serving internal documentation on a dev server
- Quick file transfer over a network

This is not a replacement for proper access control. It is a quick convenience barrier for temporary use.
