PROJECT_NAME := "lazyfile"

# Default command
_default:
    @just --list --unsorted

# Sync Go modules
tidy:
    go mod tidy
    go work sync
    @echo "All modules synced, Go workspace ready!"

# CLI local run wrapper
cli *args:
    @go run . {{ args }}

# Execute unit tests
test:
    @echo "Running unit tests!"
    go clean -testcache
    go test -cover ./...

# Run coverage and open a report
view-coverage:
    go clean -testcache
    go test -coverpkg="./..." -coverprofile="coverage.out" -covermode="count" ./...
    go tool cover -html="coverage.out" -o coverage.html
    xdg-open coverage.html

# Build CLI binary
build version="0.0.0":
    #!/usr/bin/env bash
    echo "Building {{ PROJECT_NAME }} binary..."
    go mod download all
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version={{ version }}" -o ./{{ PROJECT_NAME }} .
    echo "Built binary for {{ PROJECT_NAME }} {{ version }} successfully!"
