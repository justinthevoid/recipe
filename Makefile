.PHONY: cli cli-all tui tui-all clean wasm test branding

# Version injection (git-based or override with VERSION=x.y.z)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build CLI for current platform
cli:
	go build -ldflags="-X main.version=$(VERSION)" -o recipe cmd/cli/*.go

# Build CLI for all platforms
cli-all:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/recipe-linux-amd64 cmd/cli/*.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/recipe-darwin-amd64 cmd/cli/*.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o bin/recipe-darwin-arm64 cmd/cli/*.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/recipe-windows-amd64.exe cmd/cli/*.go

# Build TUI for current platform
tui:
	go build -o recipe-tui cmd/tui/*.go

# Build TUI for all platforms
tui-all:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/recipe-tui-linux-amd64 cmd/tui/*.go
	GOOS=darwin GOARCH=amd64 go build -o bin/recipe-tui-darwin-amd64 cmd/tui/*.go
	GOOS=darwin GOARCH=arm64 go build -o bin/recipe-tui-darwin-arm64 cmd/tui/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/recipe-tui-windows-amd64.exe cmd/tui/*.go

# Build WASM module with size optimization (production)
wasm:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/public/recipe.wasm cmd/wasm/main.go
	@echo "WASM binary size:"
	@ls -lh web/public/recipe.wasm 2>/dev/null || dir web\public\recipe.wasm

# Build WASM module without optimization (development)
wasm-dev:
	GOOS=js GOARCH=wasm go build -o web/public/recipe.wasm cmd/wasm/main.go

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total

# Generate HTML coverage report
coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run performance benchmarks
benchmark:
	@echo "Running performance benchmarks..."
	go test -bench="BenchmarkConvert_(NP3|XMP)" -benchmem -run=^$$ ./internal/converter/ | tee benchmarks.txt
	@echo ""
	@echo "Results saved to benchmarks.txt"

# Run all benchmarks (including detection and overhead)
benchmark-all:
	go test -bench=. -benchmem ./internal/converter/ | tee benchmarks.txt

# CPU profiling
profile-cpu:
	go test -bench=BenchmarkConvert_NP3_to_XMP -cpuprofile=cpu.prof ./internal/converter/
	@echo "CPU profile generated: cpu.prof"
	@echo "View with: go tool pprof -http=:8080 cpu.prof"

# Memory profiling
profile-mem:
	go test -bench=BenchmarkConvert_NP3_to_XMP -memprofile=mem.prof ./internal/converter/
	@echo "Memory profile generated: mem.prof"
	@echo "View with: go tool pprof -http=:8080 mem.prof"

# Build web interface (Astro + WASM)
web:
	@echo "Building optimized web interface..."
	@$(MAKE) wasm
	cd web && npm run build
	@echo "Web build complete! Output: web/dist/"

# Start web dev server (Astro)
web-dev:
	cd web && npm run dev

# Clean build artifacts
clean:
	rm -f recipe recipe.exe recipe-tui recipe-tui.exe
	rm -f coverage.out coverage.html
	rm -f benchmarks.txt *.prof
	rm -rf bin/
	rm -f web/static/bundle.min.js web/static/bundle.min.js.map


branding: ## Render branding assets (logo, OG image) via Remotion
	cd web/remotion && bun run render:all

# Check import graph constraints
check-imports:
	@echo "Checking import graph violations..."
	@go list -f '{{.ImportPath}}: {{join .Imports ", "}}' ./... | \
		grep -E 'internal/formats.*(cmd/|internal/batch)' && \
		echo "ERROR: Forbidden import detected" && exit 1 || echo "Import graph clean"

# Download and verify test fixtures
setup-fixtures:
	go run scripts/download_fixtures.go

# Check fixture documentation
check-fixture-docs:
	@[ -f testdata/nx-fixtures/README.md ] || (echo "ERROR: testdata/nx-fixtures/README.md missing" && exit 1)
	@echo "All fixture directories documented."
