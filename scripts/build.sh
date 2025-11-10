#!/bin/bash

# Script to build executables for different platforms
# Usage: ./scripts/build.sh

APP_NAME="ecommerce-api"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
BIN_DIR="bin"

# Create bin directory if it doesn't exist
mkdir -p $BIN_DIR

echo "🔨 Building $APP_NAME..."
echo "Version: $VERSION"
echo "Build Time: $BUILD_TIME"
echo ""

# Build flags
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME"

# Build for Linux (Ubuntu)
echo "📦 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o ${BIN_DIR}/${APP_NAME}-linux-amd64 ./cmd/api
if [ $? -eq 0 ]; then
    echo "✓ Linux build successful: ${BIN_DIR}/${APP_NAME}-linux-amd64"
else
    echo "✗ Linux build failed"
    exit 1
fi

# Build for Windows
echo "📦 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o ${BIN_DIR}/${APP_NAME}-windows-amd64.exe ./cmd/api
if [ $? -eq 0 ]; then
    echo "✓ Windows build successful: ${BIN_DIR}/${APP_NAME}-windows-amd64.exe"
else
    echo "✗ Windows build failed"
    exit 1
fi

# Build for macOS (optional)
echo "📦 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o ${BIN_DIR}/${APP_NAME}-darwin-amd64 ./cmd/api
if [ $? -eq 0 ]; then
    echo "✓ macOS build successful: ${BIN_DIR}/${APP_NAME}-darwin-amd64"
else
    echo "✗ macOS build failed"
fi

# Build for macOS ARM (Apple Silicon)
echo "📦 Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o ${BIN_DIR}/${APP_NAME}-darwin-arm64 ./cmd/api
if [ $? -eq 0 ]; then
    echo "✓ macOS ARM build successful: ${BIN_DIR}/${APP_NAME}-darwin-arm64"
else
    echo "✗ macOS ARM build failed"
fi

echo ""
echo "✅ Build completed!"
echo "Executables are in the ${BIN_DIR}/ directory"
