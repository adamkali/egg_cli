# egg_cli Makefile

# Variables
BINARY_NAME=egg_cli
VERSION?=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

.PHONY: help build clean test tidy tag push-tag all

# Default target
all: tidy test build

help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "✅ Binary built: $(BINARY_NAME)"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "✅ Clean complete"

test: ## Run tests
	@echo "Running tests..."
	@go test ./...
	@echo "✅ Tests complete"

tidy: ## Tidy and verify dependencies
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod verify
	@echo "✅ Dependencies tidied"

##@ Release

tag: ## Create and push a new git tag (usage: make tag VERSION=v1.0.0)
ifndef VERSION
	$(error VERSION is required. Usage: make tag VERSION=v1.0.0)
endif
	@echo "Creating tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "✅ Tag $(VERSION) created and pushed"

push-tag: ## Push existing tags to remote
	@echo "Pushing tags to remote..."
	@git push --tags
	@echo "✅ Tags pushed"

##@ Information

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"