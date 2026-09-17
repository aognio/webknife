# Security

## Credential redaction

- `Authorization`, `Cookie`, `Set-Cookie`, `Proxy-Authorization` are identified as sensitive headers
- `RedactValue` masks values: first 2 chars + asterisks + last 2 chars
- Query params `token`, `key`, `secret`, `password`, `api_key`, `apikey` are flagged
- Never log Basic Auth credentials

## Path traversal

- `NewStaticHandler` uses `filepath.Abs` and `http.Dir` which prevent escaping the root
- Test: `GET /../etc/passwd` must not return 200

## Proxy safety

- `httputil.NewSingleHostReverseProxy` is used correctly
- `Director` sets scheme, host, and `req.Host`
- Do not add `X-Forwarded-For` or `X-Real-IP` by default — let the upstream decide

## TLS

- Minimum version: TLS 1.2
- Cipher suites: ECDHE with AES-GCM only
- Certificate and key files are read from disk, not embedded

## Request body limits

- Echo handler has `MaxBodySize` (default 1MB)
- Bodies exceeding the limit are truncated, not rejected
- `BodyTruncated` flag is set in the response

## Defaults that matter

- `--listen :8080` binds to all interfaces — warn users
- No request timeout by default beyond `http.Server` settings
- Log format defaults to text — structured JSON requires opt-in

## What not to do

- Do not store secrets in environment variables that appear in process listings
- Do not log `Authorization` headers even in debug mode
- Do not expose internal Go types in echo output
- Do not assume TLS is enabled — document when it's not
