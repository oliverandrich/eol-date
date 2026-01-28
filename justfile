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

# Run tests
test *ARGS:
    go test ./... {{ ARGS }}

# Run tests with coverage report
cover:
    go test -coverprofile=coverage.out ./...
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

# Cross-compile for all platforms using goreleaser
release:
    goreleaser build --clean --snapshot

# Create a full release with archives and checksums
release-dist:
    goreleaser release --clean --snapshot --skip=publish

# Generate demo GIF with vhs
demo:
    @command -v vhs >/dev/null || { echo "vhs not found, please install: https://github.com/charmbracelet/vhs"; exit 1; }
    vhs < docs/eol-date.tape
