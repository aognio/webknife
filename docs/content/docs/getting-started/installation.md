---
title: "Installation"
weight: 2
---

# Installation

## Download a release binary

Precompiled binaries are available from [GitHub Releases](https://github.com/aognio/webknife/releases).

### Supported platforms

| Platform | Archive |
|----------|---------|
| Linux amd64 | `webknife_v0.1.0_linux_amd64.tar.gz` |
| Linux arm64 | `webknife_v0.1.0_linux_arm64.tar.gz` |
| macOS Intel | `webknife_v0.1.0_darwin_amd64.tar.gz` |
| macOS Apple Silicon | `webknife_v0.1.0_darwin_arm64.tar.gz` |
| Windows amd64 | `webknife_v0.1.0_windows_amd64.zip` |

### Extract and install

```bash
# Linux / macOS
tar xzf webknife_v0.1.0_linux_amd64.tar.gz
sudo mv webknife /usr/local/bin/

# Windows (PowerShell)
Expand-Archive webknife_v0.1.0_windows_amd64.zip -DestinationPath C:\webknife
```

### Verify checksum

Each release includes a `SHA256SUMS` file with SHA-256 hashes:

```bash
sha256sum -c SHA256SUMS --ignore-missing
```

Expected output:

```
webknife_v0.1.0_linux_amd64.tar.gz: OK
```

## Build from source

```bash
git clone https://github.com/aognio/webknife
cd webknife
make build
```

The binary is placed at `bin/webknife`.

## Using go install

```bash
go install github.com/aognio/webknife/cmd/webknife@latest
```

## Verify installation

```bash
webknife version
webknife help
```

## Version inspection

```bash
webknife version
# webknife v0.1.0 (commit=abc1234, built=2026-09-17T00:00:00Z, go=go1.24.13)
```

## Semantic versioning

Webknife follows [semantic versioning](https://semver.org/):

- **MAJOR** — incompatible API changes
- **MINOR** — new functionality in a backwards-compatible manner
- **PATCH** — backwards-compatible bug fixes

Release candidates use the `-rc.N` suffix (e.g., `v0.1.0-rc.1`).

## Release candidates

Pre-release tags are marked as prereleases on GitHub:

```bash
webknife version
# webknife v0.1.0-rc.1
```
