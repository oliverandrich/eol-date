# __ProjectName__ - Task Runner
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

# Release mit goreleaser erstellen
release:
    goreleaser release --clean

# Lokaler Test-Build ohne Release (Snapshot)
release-snapshot:
    goreleaser release --snapshot --clean

# Generate demo GIF with vhs
demo:
    vhs < docs/eol-date.tape
