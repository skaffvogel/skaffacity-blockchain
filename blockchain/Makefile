PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BUILD_TAGS := netgo
BUILD_FLAGS := -tags "$(BUILD_TAGS)" -ldflags '-X github.com/cosmos/cosmos-sdk/version.Name=skaffacity \
  -X github.com/cosmos/cosmos-sdk/version.AppName=skaffacityd \
  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) -w -s'

# Default target
all: build

# Install for development
install: go.sum
	@echo "Installing skaffacityd..."
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/skaffacityd

# Build binary for production/Pterodactyl
build: clean-build
	@echo "Building skaffacityd binary..."
	@mkdir -p bin
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -a -installsuffix cgo -o bin/skaffacityd ./cmd/skaffacityd
	@echo "Binary built: bin/skaffacityd"
	@ls -la bin/skaffacityd

# Pterodactyl-specific build (minimal, static)
build-pterodactyl: clean-build
	@echo "Building for Pterodactyl deployment..."
	@echo "Setting up Go cache directories..."
	@mkdir -p bin /tmp/go-cache /tmp/go-mod
	GOCACHE=/tmp/go-cache GOMODCACHE=/tmp/go-mod \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags='-w -s -extldflags "-static"' \
		-a -installsuffix cgo \
		-o bin/skaffacityd ./cmd/skaffacityd
	@echo "Pterodactyl binary built: bin/skaffacityd"
	@echo "Binary size: $(shell du -h bin/skaffacityd | cut -f1)"

# Ultra-minimal build for space-constrained environments
build-minimal: clean-build
	@echo "Building minimal binary for space-constrained deployment..."
	@echo "Cleaning all caches first..."
	@go clean -cache -modcache 2>/dev/null || true
	@mkdir -p bin
	GOCACHE=off GOTMPDIR=/tmp \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-mod=readonly \
		-ldflags='-w -s -extldflags "-static"' \
		-trimpath \
		-o bin/skaffacityd ./cmd/skaffacityd
	@echo "Cleaning up after build..."
	@go clean -cache 2>/dev/null || true
	@echo "Minimal binary built: bin/skaffacityd"
	@echo "Binary size: $(shell du -h bin/skaffacityd | cut -f1)"

# Copy essential files for deployment
prepare-deployment: build-pterodactyl
	@echo "Preparing deployment files..."
	@mkdir -p deployment
	@cp bin/skaffacityd deployment/
	@cp -r cmd deployment/ 2>/dev/null || true
	@cp go.mod deployment/
	@cp go.sum deployment/
	@echo "Deployment ready in ./deployment/"

# Initialize blockchain
init:
	@echo "Initializing blockchain..."
	@mkdir -p ~/.skaffacity
	./bin/skaffacityd init local --chain-id skaffacity-1 --home ~/.skaffacity

# Start blockchain
start:
	@echo "Starting blockchain..."
	./bin/skaffacityd start --home ~/.skaffacity

# Test
test:
	@go test -mod=readonly $(PACKAGES)

# Clean build artifacts
clean-build:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf deployment/

# Full clean
clean: clean-build
	@echo "Full cleanup..."
	@go clean -cache -modcache 2>/dev/null || true
	@rm -rf build/
	@rm -rf temp-*

# Development helpers
dev-setup: build init
	@echo "Development environment ready"

# Show build info
info:
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(COMMIT)"
	@echo "Version: $(VERSION)"
	@echo "Build flags: $(BUILD_FLAGS)"

.PHONY: all install build build-pterodactyl build-minimal prepare-deployment init start test clean-build clean dev-setup info
