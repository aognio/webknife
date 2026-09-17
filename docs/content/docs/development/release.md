---
title: "Release"
weight: 5
---

# Release

How Webknife releases are built and published.

## Platforms

| OS | Arch | Archive |
|----|------|---------|
| Linux | amd64 | `webknife_<version>_linux_amd64.tar.gz` |
| Linux | arm64 | `webknife_<version>_linux_arm64.tar.gz` |
| macOS | amd64 (Intel) | `webknife_<version>_darwin_amd64.tar.gz` |
| macOS | arm64 (Apple Silicon) | `webknife_<version>_darwin_arm64.tar.gz` |
| Windows | amd64 | `webknife_<version>_windows_amd64.zip` |

## CI pipeline

Runs on every push and pull request:

- Format check
- `go vet`
- `go test -race`
- Build verification
- Architecture dependency checks

## Release pipeline

Triggered by version tags (`v*`) or manual dispatch.

### Manual dry run

Trigger the release workflow from GitHub Actions to test the build pipeline without creating a Release. Artifacts are uploaded for inspection.

### Tagged release

```bash
git tag -a v0.1.0 -m "Webknife v0.1.0"
git push origin v0.1.0
```

This creates a GitHub Release with:
- Cross-compiled binaries for all platforms
- SHA256SUMS checksum file
- Auto-generated release notes

## Checksum verification

```bash
sha256sum -c SHA256SUMS --ignore-missing
# webknife_v0.1.0_linux_amd64.tar.gz: OK
```

## Version metadata

```bash
webknife version
# webknife v0.1.0 (commit=abc1234, built=2026-09-17T00:00:00Z, go=go1.24.13)
```

## Release checklist

- [ ] `make check` passes
- [ ] Architecture checks pass
- [ ] Hugo docs build
- [ ] Manual test scenarios exercised
- [ ] Release workflow dry-run succeeds
- [ ] Artifact tested on at least one platform
- [ ] SHA256SUMS verified
- [ ] Version metadata correct
- [ ] README accurate
- [ ] wip merged into main
- [ ] Final commit reviewed
- [ ] Tag created from intended commit

## Release tag discipline

A published semantic release tag is immutable. Use release candidates (`-rc.N`) for pre-release testing. Do not delete and recreate published tags.
