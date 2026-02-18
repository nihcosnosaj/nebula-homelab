#!/bin/bash
VERSION="1.0.0"
PLATFORMS=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")

for platform in "${PLATFORMS[@]}"; do
    OS=${platform%/*}
    ARCH=${platform#*/}
    OUTPUT="dist/nebula-${VERSION}-${OS}-${ARCH}"

    echo "Building for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT ./cmd/nebula/main.go
done