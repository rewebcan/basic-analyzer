SHELL := /bin/sh

GO ?= go
APP_DIR ?= ./cmd/web
PKG ?= ./...
TIMEOUT ?= 60s

.PHONY: help run test integration-test itest

help:
	@echo "Usage:"
	@echo "  make run                # Run the application"
	@echo "  make test               # Run unit tests"
	@echo "  make integration-test   # Run integration tests (requires -tags=integration)"
	@echo "  make itest              # Alias for integration-test"

run:
	$(GO) run $(APP_DIR)

test:
	$(GO) test $(PKG) -count=1 -race -timeout $(TIMEOUT) -coverprofile=coverage.out

integration-test itest:
	$(GO) test -tags=integration $(PKG) -count=1 -race -timeout $(TIMEOUT)


