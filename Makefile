SHELL := /bin/sh

GO ?= go
GOLANGCI_LINT ?= golangci-lint
APP_DIR ?= ./cmd/web
PKG ?= ./...
TIMEOUT ?= 60s
BINARY_NAME ?= url-fetcher
BUILD_DIR ?= ./build
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: help run test e2e-test lint lint-fix lint-setup build build-linux build-darwin build-windows build-all clean

help:
	@echo "Usage:"
	@echo "  make run                # Run the application"
	@echo "  make test               # Run unit tests"
	@echo "  make e2e-test           # Run end-to-end tests"
	@echo "  make lint               # Run golangci-lint"
	@echo "  make lint-fix           # Run linters with fixes where supported"
	@echo "  make lint-setup         # Install golangci-lint locally"
	@echo "  make build              # Build the application for current platform"
	@echo "  make build-linux        # Build for Linux (amd64)"
	@echo "  make build-darwin       # Build for macOS (amd64)"
	@echo "  make build-windows      # Build for Windows (amd64)"
	@echo "  make clean              # Clean build artifacts"

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

# Build targets
build: clean
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) $(APP_DIR)
	@echo "Binary created at $(BUILD_DIR)/$(BINARY_NAME)"

build-linux: clean
	@echo "Building $(BINARY_NAME) for Linux (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(APP_DIR)
	@echo "Binary created at $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

build-darwin: clean
	@echo "Building $(BINARY_NAME) for macOS (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(APP_DIR)
	@echo "Binary created at $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64"

build-windows: clean
	@echo "Building $(BINARY_NAME) for Windows (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(APP_DIR)
	@echo "Binary created at $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

# Build for all platforms
build-all: build-linux build-darwin build-windows
	@echo "Built binaries for all platforms in $(BUILD_DIR)/"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@echo "Clean complete"


