# Makefile for webknife — a single-binary HTTP diagnostic Swiss Army knife.
#
# Delete what you don't use. An unused target is worse than a missing one.

# --- Preamble ----------------------------------------------------------------
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.DEFAULT_GOAL := help

V ?= 0
ifeq ($(V),1)
  Q :=
else
  Q := @
endif

# --- Project -----------------------------------------------------------------
BINARY  := webknife
MODULE  := github.com/webknife/webknife
CMD     := ./cmd/$(BINARY)
BUILDPKG := $(MODULE)/internal/build

OUTDIR  := bin
TARGET  := $(OUTDIR)/$(BINARY)

PREFIX  ?= $(HOME)/.local
BINDIR  ?= $(PREFIX)/bin
DESTDIR ?=

GO      ?= go
GOFLAGS ?=

ifeq ($(origin VERSION),undefined)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
endif
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE   := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X $(BUILDPKG).Version=$(VERSION) \
           -X $(BUILDPKG).Commit=$(COMMIT) \
           -X $(BUILDPKG).Date=$(DATE)

SOURCES := $(shell find . -name '*.go' -not -path './vendor/*' -not -path './.development/*' 2>/dev/null)

# Extra arguments: make run ARGS='echo --listen :9090'
ARGS ?=

##@ Build

.PHONY: build
build: $(TARGET) ## Build the binary into bin/

$(TARGET): $(SOURCES) | $(OUTDIR)
	$(Q)$(GO) build $(GOFLAGS) -trimpath -ldflags '$(LDFLAGS)' -o $@ $(CMD)
	$(Q)echo "built $@ ($(VERSION))"

$(OUTDIR):
	$(Q)mkdir -p $@

.PHONY: release
release: ## Build a stripped, statically linked binary
	$(Q)mkdir -p $(OUTDIR)
	$(Q)CGO_ENABLED=0 $(GO) build $(GOFLAGS) -trimpath \
		-ldflags '-s -w $(LDFLAGS)' -o $(TARGET) $(CMD)
	$(Q)echo "built $(TARGET) ($(VERSION), stripped)"

.PHONY: run
run: $(TARGET) ## Build and run (make run ARGS='echo --listen :9090')
	$(Q)$(TARGET) $(ARGS)

##@ Quality

.PHONY: test
test: ## Run tests
	$(Q)$(GO) test $(ARGS) ./...

.PHONY: test-race
test-race: ## Run tests with the race detector, uncached
	$(Q)$(GO) test -race -count=1 ./...

.PHONY: cover
cover: ## Run tests and open an HTML coverage report
	$(Q)$(GO) test -coverprofile=coverage.out ./...
	$(Q)$(GO) tool cover -html=coverage.out

.PHONY: vet
vet: ## Run go vet
	$(Q)$(GO) vet ./...

.PHONY: fmt
fmt: ## Format code (gofumpt if available, else gofmt)
	@if command -v gofumpt >/dev/null 2>&1; then \
		gofumpt -w .; \
	else \
		gofmt -s -w .; \
	fi

.PHONY: fmt-check
fmt-check: ## Fail if any file is unformatted
	$(Q)out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "unformatted files:" >&2; \
		echo "$$out" >&2; \
		exit 1; \
	fi

.PHONY: lint
lint: require-golangci-lint ## Run golangci-lint
	$(Q)golangci-lint run ./...

.PHONY: check
check: fmt-check vet test ## Run everything CI runs

##@ Modules

.PHONY: tidy
tidy: ## Run go mod tidy
	$(Q)$(GO) mod tidy

.PHONY: tidy-check
tidy-check: ## Fail if go.mod/go.sum are not tidy
	$(Q)$(GO) mod tidy -diff

##@ Install

.PHONY: install
install: $(TARGET) ## Install the binary into PREFIX/bin
	$(Q)install -D -m 0755 $(TARGET) $(DESTDIR)$(BINDIR)/$(BINARY)
	$(Q)echo "installed $(DESTDIR)$(BINDIR)/$(BINARY)"

.PHONY: uninstall
uninstall: ## Remove the installed binary
	$(Q)rm -f $(DESTDIR)$(BINDIR)/$(BINARY)

.PHONY: clean
clean: ## Remove build artifacts
	$(Q)rm -rf $(OUTDIR) coverage.out

##@ Help

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z0-9_%\/-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
	@echo ""

.PHONY: require-%
require-%:
	@command -v $* >/dev/null 2>&1 || { \
		echo "error: '$*' is required but not installed" >&2; exit 1; }

.PHONY: print-%
print-%:
	@echo '$*=$($*)'
