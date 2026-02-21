# P2P Network Makefile

# Build variables
BINARY_NAME=p2p-node
BUILD_DIR=bin
GO=go
GOFLAGS=-v
LDFLAGS=

# Colors
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: all build clean test test-integration test-unit lint benchmark install dev

all: build

# Build the binary
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/node
	@echo "$(GREEN)Build complete!$(NC)"

# Build for different platforms
build-darwin:
	@echo "$(GREEN)Building for darwin...$(NC)"
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/node
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/node
	@echo "$(GREEN)Darwin build complete!$(NC)"

build-linux:
	@echo "$(GREEN)Building for linux...$(NC)"
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/node
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/node
	@echo "$(GREEN)Linux build complete!$(NC)"

build-windows:
	@echo "$(GREEN)Building for windows...$(NC)"
	GOOS=windows $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe ./cmd/node
	@echo "$(GREEN)Windows build complete!$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	$(GO) clean
	@echo "$(GREEN)Clean complete!$(NC)"

# Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)Tests complete!$(NC)"

# Run unit tests only
test-unit:
	@echo "$(GREEN)Running unit tests...$(NC)"
	$(GO) test -v -race ./test/unit/...
	@echo "$(GREEN)Unit tests complete!$(NC)"

# Run integration tests
test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GO) test -v -race ./test/integration/...
	@echo "$(GREEN)Integration tests complete!$(NC)"

# Run benchmarks
benchmark:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	$(GO) test -bench=. -benchmem ./benchmarks/...
	@echo "$(GREEN)Benchmarks complete!$(NC)"

# Run linter
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@which golangci-lint > /dev/null || $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...
	@echo "$(GREEN)Linting complete!$(NC)"

# Install dependencies
install:
	@echo "$(GREEN)Installing dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)Installation complete!$(NC)"

# Development setup
dev: install build

# Show help
help:
	@echo "P2P Network Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build           - Build the binary"
	@echo "  build-darwin   - Build for Darwin (macOS)"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-windows  - Build for Windows"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests"
	@echo "  benchmark      - Run benchmarks"
	@echo "  lint           - Run linter"
	@echo "  install        - Install dependencies"
	@echo "  dev            - Development setup (install + build)"
	@echo "  help           - Show this help message"
