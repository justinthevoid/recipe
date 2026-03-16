# Contributing to Recipe

Thanks for your interest in contributing to Recipe! This guide will help you get started.

## Reporting Issues

- **Bugs**: Open a [GitHub Issue](https://github.com/justinthevoid/recipe/issues) with steps to reproduce, expected vs actual behavior, and your OS/browser version.
- **Feature requests**: Open a GitHub Issue describing your use case and the desired behavior.

## Development Setup

**Prerequisites:**
- Go 1.25.1+
- Bun (for workspace dependency management)
- Node.js 18+ (for web frontend)

```bash
# Clone the repo
git clone https://github.com/justinthevoid/recipe.git
cd recipe

# Install dependencies
bun install

# Run tests (uses committed fixtures, no external data needed)
go test ./...

# Build the CLI
make cli

# Build WASM for the web interface
make wasm

# Start the web dev server
cd web && npm run dev
```

## Running Tests

All tests use committed fixture files in package-level `testdata/` directories. No external downloads needed.

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/formats/np3/
go test ./internal/converter/

# Run with coverage
make coverage
```

## Building

```bash
# CLI binary
make cli

# CLI for all platforms
make cli-all

# WASM module
make wasm

# Web frontend
cd web && npm run build
```

## Pull Request Process

1. Fork the repo and create a feature branch from `main`
2. Make your changes, following existing code patterns
3. Ensure all tests pass: `go test ./...`
4. Ensure the build succeeds: `go build ./...`
5. Open a PR against `main` with a clear description of the change

## Code Style

- Follow existing patterns in the codebase
- Go code follows standard `gofmt` formatting
- Web code uses Biome for linting/formatting

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
