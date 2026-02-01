# eol-date - Task Runner

# Version from git
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`

# Default: show available commands
default:
    @just --list

# Setup project after checkout
setup:
    go mod download
    @command -v pre-commit >/dev/null && pre-commit install || echo "pre-commit not found, skipping hook installation"

# Build binary
build:
    go build -ldflags="-s -w -X 'main.version={{ version }}'" -trimpath -o build/eol-date ./cmd/eol-date

# Run all tests
# Requires: go install github.com/mfridman/tparse@latest
test:
    set -o pipefail && go test -json ./... | tparse -progress

# Run tests with coverage report
# Requires: go install github.com/mfridman/tparse@latest
cover:
    set -o pipefail && go test -json -coverprofile=coverage.out ./... | tparse -progress

# Open coverage report in browser
cover-report:
    go tool cover -html=coverage.out

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    golangci-lint run

# Run all checks (format, lint, test)
check:
    just fmt
    just lint
    just test

# Clean build artifacts
clean:
    rm -rf build/

# Install binary to $GOPATH/bin
install:
    go install -ldflags="-s -w -X 'main.version={{ version }}'" ./cmd/eol-date

# Tidy dependencies
tidy:
    go mod tidy

# Create and publish a release (requires git tag)
release *ARGS:
    goreleaser release --clean {{ARGS}}

# Generate demo GIF with vhs
demo:
    @command -v vhs >/dev/null || { echo "vhs not found, please install: https://github.com/charmbracelet/vhs"; exit 1; }
    vhs < docs/eol-date.tape
