#!/bin/bash

# SkaffaCity Cross-Platform Build Script

set -e

echo "🏗️  Building SkaffaCity Blockchain..."

# Build variables
VERSION=${VERSION:-$(git describe --tags 2>/dev/null || echo "v0.1.0")}
COMMIT=${COMMIT:-$(git rev-parse HEAD 2>/dev/null || echo "unknown")}
BUILD_TAGS="netgo"
BUILD_FLAGS="-tags ${BUILD_TAGS} -ldflags -w -s -X github.com/cosmos/cosmos-sdk/version.Name=skaffacity -X github.com/cosmos/cosmos-sdk/version.AppName=skaffacityd -X github.com/cosmos/cosmos-sdk/version.Version=${VERSION} -X github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT}"

# Create bin directory
mkdir -p bin

# Linux build (default for servers)
echo "🐧 Building for Linux..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_FLAGS} -o bin/skaffacityd ./cmd/skaffacityd

# Windows build (for development)
if [ "$1" == "all" ]; then
    echo "🪟 Building for Windows..."
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_FLAGS} -o bin/skaffacityd.exe ./cmd/skaffacityd
fi

# macOS build (for development)
if [ "$1" == "all" ]; then
    echo "🍎 Building for macOS..."
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_FLAGS} -o bin/skaffacityd-darwin ./cmd/skaffacityd
fi

echo "✅ Build complete!"
ls -la bin/
