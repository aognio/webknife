---
title: "Installation"
weight: 2
---

# Installation

## From source (recommended)

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

## Pre-built binaries

Check the [GitHub releases](https://github.com/aognio/webknife/releases) page.

## Verify installation

```bash
./bin/webknife version
./bin/webknife help
```
