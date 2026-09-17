---
title: "Building"
weight: 1
---

# Building

## Build the binary

```bash
make build
```

Output: `bin/webknife`

## Other targets

```bash
make test         # run all tests
make check        # fmt-check + vet + test
make fmt          # format code
make vet          # run go vet
make lint         # run golangci-lint
make cover        # open HTML coverage report
make clean        # remove build artifacts
make help         # show all targets
```

## Version injection

Build metadata is injected via linker flags:

```bash
make build VERSION=1.2.3
```

Or use the default (git-based):

```bash
make build
# built bin/webknife (v0.1.0-3-g17d4942)
```

## Static build

```bash
make release
```

Produces a stripped, statically linked binary.
