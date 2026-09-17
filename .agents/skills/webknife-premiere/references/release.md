# Release Engineering

## CI

CI runs automatically on push/PR to `main` and `wip`:

- Format check (`gofmt -l`)
- `go vet ./...`
- `go test -race -count=1 ./...`
- `go build ./...`
- Architecture dependency checks

## Local CI equivalent

```bash
make check
```

This reproduces the important CI checks locally.

## Release workflow

Triggered by:
- **Version tags** (`v*`) — creates a GitHub Release
- **Manual dispatch** — dry run, uploads artifacts without creating a Release

### Release dry run

From the `wip` branch, trigger the release workflow manually via GitHub Actions UI. This:

1. Runs validation (tests, vet, arch checks)
2. Cross-compiles for all platforms
3. Packages archives
4. Generates SHA256SUMS
5. Uploads as GitHub Actions artifacts (no Release created)

### Actual release

```bash
git checkout main
git merge wip
git push origin main

git tag -a v0.1.0 -m "Webknife v0.1.0"
git push origin v0.1.0
```

GitHub automatically:
1. Validates code
2. Cross-compiles for all platforms
3. Packages archives with LICENSE and README
4. Generates SHA256SUMS
5. Creates GitHub Release with generated notes

## Target platforms

| OS | Arch | Archive |
|----|------|---------|
| linux | amd64 | `.tar.gz` |
| linux | arm64 | `.tar.gz` |
| darwin | amd64 | `.tar.gz` |
| darwin | arm64 | `.tar.gz` |
| windows | amd64 | `.zip` |

## Version metadata

Injected via Go linker flags (`-ldflags`):

```
-X github.com/aognio/webknife/internal/build.Version=<version>
-X github.com/aognio/webknife/internal/build.Commit=<sha>
-X github.com/aognio/webknife/internal/build.Date=<utc-timestamp>
```

## Artifact naming

```
webknife_<version>_<os>_<arch>.tar.gz
webknife_<version>_windows_amd64.zip
SHA256SUMS
```

## Release tag discipline

- A published semantic release tag is immutable
- Do not delete and recreate published tags
- Use `-rc.N` suffix for release candidates
- Pre-release tags are marked as prereleases on GitHub

## Checksums

SHA-256 checksums are generated for every release archive:

```bash
sha256sum -c SHA256SUMS --ignore-missing
```

## Release checklist

- [ ] Working tree clean
- [ ] CI green (`make check`)
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
- [ ] v0.1.0 tag created from intended commit
