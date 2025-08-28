SHELL := /bin/sh

GO ?= go
GOLANGCI_LINT ?= golangci-lint
APP_DIR ?= ./cmd/web
PKG ?= ./...
TIMEOUT ?= 60s

.PHONY: help run test e2e-test lint lint-fix lint-setup

help:
	@echo "Usage:"
	@echo "  make run                # Run the application"
	@echo "  make test               # Run unit tests"
	@echo "  make e2e-test           # Run end-to-end tests"
	@echo "  make lint               # Run golangci-lint"
	@echo "  make lint-fix           # Run linters with fixes where supported"
	@echo "  make lint-setup         # Install golangci-lint locally"

run:
	$(GO) run $(APP_DIR)

test:
	$(GO) test $(PKG) -count=1 -race -timeout $(TIMEOUT) -coverprofile=coverage.out

e2e-test:
	$(GO) test ./cmd/web/tests/ -count=1 -race -timeout $(TIMEOUT) -v

lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint not found. Run 'make lint-setup' first."; exit 1; }
	$(GOLANGCI_LINT) run ./...

lint-fix:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint not found. Run 'make lint-setup' first."; exit 1; }
	$(GOLANGCI_LINT) run --fix ./...

lint-setup:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Ensure $$GOPATH/bin (or Go install bin dir) is on your $$PATH."


