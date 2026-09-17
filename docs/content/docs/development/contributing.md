---
title: "Contributing"
weight: 3
---

# Contributing

## Getting started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Run `make check`
6. Submit a pull request

## Code style

- Standard library first — avoid external dependencies
- Independent features — no sibling imports between `pkg/*` packages
- Composition via `http.Handler`
- Tests use `httptest`, no external Internet

## Before submitting

```bash
make check   # fmt-check + vet + test
go test ./internal/archtest/  # dependency rules
```

## Adding a new feature

1. Create `pkg/newfeature/newfeature.go`
2. Add tests in `pkg/newfeature/newfeature_test.go`
3. Wire in `internal/application/app.go`
4. Add CLI flags in `internal/adapters/cli/cli.go`
5. Run `make check`
