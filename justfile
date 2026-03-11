# Default recipe: list available commands
default:
    @just --list

# Run all tests
test *args:
    go test -json {{args}} ./... | tparse

# Run linter
lint:
    golangci-lint run ./...

# Format all Go files
fmt:
    gofmt -w .
    goimports -w .

# Run tests with coverage report
coverage:
    #!/usr/bin/env bash
    set -euo pipefail
    go test -json -coverprofile=coverage.out ./... | tparse
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report: coverage.html"

# Tidy module dependencies
tidy:
    go mod tidy

# Install binary to $GOPATH/bin
install:
    go install -ldflags="-s -w -X 'main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)'" ./cmd/eol-date

# Check that all required dev tools are installed
setup:
    #!/usr/bin/env bash
    set -euo pipefail
    ok=true
    check() {
        if command -v "$1" &>/dev/null; then
            printf "  %-20s %s\n" "$1" "$(command -v "$1")"
        else
            printf "  %-20s MISSING — %s\n" "$1" "$2"
            ok=false
        fi
    }
    echo "Checking dev tools:"
    check go              "https://go.dev/dl/"
    check golangci-lint   "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    check tparse          "go install github.com/mfridman/tparse@latest"
    check goimports       "go install golang.org/x/tools/cmd/goimports@latest"
    check pre-commit      "https://pre-commit.com/#install"
    echo ""
    if $ok; then
        echo "All tools installed."
        echo "Run 'pre-commit install' to set up git hooks."
    else
        echo "Some tools are missing. Install them and re-run 'just setup'."
        exit 1
    fi

# Generate demo GIF with vhs
demo:
    @command -v vhs >/dev/null || { echo "vhs not found, please install: https://github.com/charmbracelet/vhs"; exit 1; }
    vhs < docs/eol-date.tape
