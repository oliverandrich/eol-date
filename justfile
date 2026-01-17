# __ProjectName__ - Task Runner

# Version from git
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`

# Default: show available commands
default:
    @just --list

# Build binary
build:
    go build -ldflags="-s -w -X 'main.version={{version}}'" -trimpath -o build/__ProjectName__ ./cmd/__ProjectName__

# Run tests
test *ARGS:
    go test ./... {{ARGS}}

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
    go install -ldflags="-s -w -X 'main.version={{version}}'" ./cmd/__ProjectName__

# Release mit goreleaser erstellen
release:
    goreleaser release --clean

# Lokaler Test-Build ohne Release (Snapshot)
release-snapshot:
    goreleaser release --snapshot --clean
