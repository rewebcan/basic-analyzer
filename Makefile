SHELL := /bin/sh

GO ?= go
GOLANGCI_LINT ?= golangci-lint
APP_DIR ?= ./cmd/web
PKG ?= ./...
TIMEOUT ?= 60s

.PHONY: help run test integration-test itest lint lint-fix lint-setup

help:
	@echo "Usage:"
	@echo "  make run                # Run the application"
	@echo "  make test               # Run unit tests"
	@echo "  make integration-test   # Run integration tests (requires -tags=integration)"
	@echo "  make itest              # Alias for integration-test"
	@echo "  make lint               # Run golangci-lint"
	@echo "  make lint-fix           # Run linters with fixes where supported"
	@echo "  make lint-setup         # Install golangci-lint locally"

run:
	$(GO) run $(APP_DIR)

test:
	$(GO) test $(PKG) -count=1 -race -timeout $(TIMEOUT) -coverprofile=coverage.out

integration-test itest:
	$(GO) test -tags=integration $(PKG) -count=1 -race -timeout $(TIMEOUT)

lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint not found. Run 'make lint-setup' first."; exit 1; }
	$(GOLANGCI_LINT) run ./...

lint-fix:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint not found. Run 'make lint-setup' first."; exit 1; }
	$(GOLANGCI_LINT) run --fix ./...

lint-setup:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Ensure $$GOPATH/bin (or Go install bin dir) is on your $$PATH."


