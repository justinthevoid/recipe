# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Recipe is a universal photo preset converter that converts between Nikon NP3, Adobe Lightroom XMP, and lrtemplate formats. The project provides three interfaces (CLI, TUI, Web) that all share a common conversion engine.

**Key Characteristics:**
- **Privacy-first**: All processing happens locally (no server uploads)
- **Performance**: Sub-millisecond conversions (<100ms WASM target)
- **Accuracy**: 98%+ conversion fidelity via exact offset mapping (48 NP3 parameters)
- **Architecture**: Hub-and-spoke pattern with UniversalRecipe intermediate representation

## Essential Commands

### Building

```bash
# Build CLI for current platform
make cli
# or: go build -o recipe cmd/cli/*.go

# Build TUI (interactive file browser)
make tui
# or: go build -o recipe-tui cmd/tui/*.go

# Build WASM for web interface (with size optimization)
make wasm
# or: GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/recipe.wasm cmd/wasm/main.go

# Build for all platforms
make cli-all
make tui-all
```

### Testing

```bash
# Run all tests (1,531 sample files)
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for specific package
go test ./internal/formats/np3/
go test ./internal/converter/

# Run specific test
go test -run TestRoundTrip_NP3_XMP ./internal/converter/

# Run with coverage (current: 89.5%)
make coverage

# Generate HTML coverage report
make coverage-html
```

### Performance Benchmarks

```bash
# Run conversion benchmarks
make benchmark

# Run all benchmarks including detection
make benchmark-all

# CPU profiling
make profile-cpu
# View: go tool pprof -http=:8080 cpu.prof

# Memory profiling
make profile-mem
# View: go tool pprof -http=:8080 mem.prof
```

### Local Development

```bash
# Serve web interface locally for WASM testing
cd web
python3 -m http.server 8080
# Open http://localhost:8080
```

## Architecture

### Hub-and-Spoke Conversion Pattern

All conversions flow through a central `UniversalRecipe` intermediate representation:

```
NP3 ──Parse──→ UniversalRecipe ──Generate──→ XMP
XMP ──Parse──→ UniversalRecipe ──Generate──→ lrtemplate
lrtemplate ──Parse──→ UniversalRecipe ──Generate──→ NP3
```

**Why this matters:**
- Adding a new format requires only 2 functions (Parse + Generate), not N converters
- All conversions use the same API: `converter.Convert(input, from, to)`
- Parameter mapping logic is centralized in UniversalRecipe

### Core Conversion API

**Single entry point for all conversions:**

```go
// internal/converter/converter.go
func Convert(input []byte, from, to string) ([]byte, error)
```

**Critical Rules:**
1. All interfaces (CLI, TUI, WASM) MUST use `converter.Convert()` - never call format parsers directly
2. All conversion errors MUST be wrapped in `ConversionError` type
3. The API is thread-safe and stateless

### Package Structure

```
cmd/
├── cli/           # Cobra CLI application
├── tui/           # Bubbletea TUI application
└── wasm/          # WASM export entry point

internal/
├── converter/     # Core conversion engine (single source of truth)
├── formats/       # Format parsers/generators
│   ├── np3/       # Nikon binary format
│   ├── xmp/       # Adobe Lightroom XML
│   ├── lrtemplate/# Lightroom Classic Lua
│   └── dcp/       # DCP color profiles (support)
├── models/        # UniversalRecipe data structures
├── inspect/       # Parameter inspection and diff tools
└── lut/           # LUT table handling

web/               # Static web interface (vanilla JS + WASM)
testdata/          # 1,531 real sample files (73 NP3, 914 XMP, 544 lrtemplate)
```

### Format Package Pattern

**Every format package follows identical structure:**

```
internal/formats/{format}/
├── parse.go          # Parse([]byte) (*UniversalRecipe, error)
├── generate.go       # Generate(*UniversalRecipe) ([]byte, error)
└── {format}_test.go  # Table-driven tests with real samples
```

**When adding a new format:**
1. Copy the structure from an existing format package
2. Implement Parse() and Generate() functions
3. Add test files to `testdata/{format}/`
4. Update `converter.Convert()` switch statements

## Key Implementation Details

### NP3 Binary Format

The NP3 (Nikon Picture Control) format is a proprietary binary format that was reverse-engineered through clean-room analysis.

**Critical implementation notes:**
- Magic bytes: "NCP" (0x4E, 0x43, 0x50)
- Fixed size: 1024 bytes
- Parameters stored at fixed byte offsets (documented in `docs/np3-format-specification.md`)
- Uses signed byte normalization (128 = zero point, ±127 = ±100%)
- Phase 5 implementation includes exact offset mapping for 48 parameters

**Key files:**
- `internal/formats/np3/parse.go` - Binary parsing logic
- `internal/formats/np3/generate.go` - Binary generation logic
- `docs/np3-format-specification.md` - Complete format documentation

### XMP/lrtemplate XML Formats

Both use XML/Lua text formats with standard parsing:

**XMP (Adobe Lightroom CC):**
- XML format with `crs:` namespace for adjustments
- Uses ElementTree for parsing
- Full parameter support (50+ fields)

**lrtemplate (Lightroom Classic):**
- Lua table syntax wrapped in XML
- Direct string parsing (not proper Lua interpreter)
- Identical parameter set to XMP

### Error Handling Pattern

**All conversion errors use the ConversionError type:**

```go
type ConversionError struct {
    Operation string  // "parse", "generate", "validate", "detect"
    Format    string  // "np3", "xmp", "lrtemplate"
    Cause     error   // Underlying error
}
```

**Usage:**
```go
if err != nil {
    return &ConversionError{
        Operation: "parse",
        Format:    "np3",
        Cause:     err,
    }
}
```

### Testing Strategy

**1,531 real sample files across all formats:**
- 73 NP3 files
- 914 XMP files
- 544 lrtemplate files

**Round-trip testing validates conversion fidelity:**
- Full fidelity paths: NP3↔XMP, NP3↔lrtemplate, XMP↔lrtemplate (when staying in Adobe ecosystem)
- Known limitations: XMP→NP3→XMP and lrtemplate→NP3→lrtemplate (some parameters unsupported by NP3)

**Test execution:**
- Tests complete in <2 seconds (parallel execution)
- Coverage: 89.5% across internal packages
- All tests use table-driven pattern with real files

### WASM Implementation

**Go 1.24+ with `go:wasmexport` directive:**

```go
//go:wasmexport convertPreset
func convertPreset(inputPtr, inputLen uint32, srcFormat, dstFormat string) (uint32, uint32, string)
```

**Key details:**
- Direct memory access (zero reflection overhead)
- Returns (outputPtr, outputLen, errorMsg) tuple
- Binary size: 4.0 MB stripped, 1.13 MB gzipped
- Target: <100ms conversions (actual: 0.003-0.079ms)

### Logging with slog

**Structured logging using Go 1.21+ stdlib slog:**

```go
logger.Info("conversion complete",
    slog.String("file", path),
    slog.String("from", sourceFormat),
    slog.String("to", targetFormat),
    slog.Duration("elapsed", elapsed),
)
```

**Log levels:**
- DEBUG: Internal state, parsed values (verbose mode only)
- INFO: Conversion start/complete, file counts
- WARN: Missing optional fields, approximated values
- ERROR: Conversion failures, invalid input

## Development Workflow

### Adding a New Feature

1. **Check architecture**: Does it fit hub-and-spoke pattern?
2. **Add tests first**: Create test files in `testdata/`
3. **Implement**: Follow existing package patterns
4. **Benchmark**: Ensure performance meets targets
5. **Document**: Update relevant docs in `docs/`

### Modifying Conversion Logic

1. **Update UniversalRecipe**: If adding new parameters (`internal/models/recipe.go`)
2. **Update all format packages**: Parse + Generate for np3, xmp, lrtemplate
3. **Add round-trip tests**: Validate conversion fidelity
4. **Update documentation**: Parameter mapping matrix in `docs/`

### Working with Binary Formats

**For NP3 format debugging:**
```bash
# Inspect binary structure with hex dump
./recipe inspect portrait.np3 --binary

# Compare parameters before/after conversion
./recipe diff original.np3 converted.xmp
```

## Important Constraints

### Format Limitations

**NP3 format has limited parameter support compared to XMP/lrtemplate:**
- ❌ Not supported: Vibrance, Temperature/Tint (partial), Grain, Vignette, Parametric Tone Curves
- ✅ Well supported: Exposure, Contrast, Saturation, Sharpness, Highlights, Shadows, Whites, Blacks, Clarity, HSL Color, Color Grading, Tone Curve control points

**Always test round-trip conversions when working with NP3:**
- XMP → NP3 → XMP may lose parameters
- NP3 → XMP → NP3 preserves all (full fidelity)

### Performance Requirements

**All conversions must meet these targets:**
- WASM: <100ms (actual: 0.003-0.079ms) ✅
- CLI: <20ms (actual: 0.003-0.079ms) ✅
- Batch (100 files): <2s (actual: 37ms) ✅
- Memory: <4096 B/op (actual: 8,890-29,026 B/op) ✅

### Privacy Guarantee

**Web interface must maintain zero network requests:**
- No analytics, tracking, or telemetry
- No file uploads to servers
- All processing via WebAssembly in browser
- Validate with browser DevTools Network tab

## Deployment

### Web Interface (Cloudflare Pages)

**Automatic deployment on push to main:**
1. GitHub Actions builds WASM binary (`.github/workflows/deploy-pages.yml`)
2. Deploys `web/static/` to Cloudflare Pages
3. Live at https://recipe.pages.dev in 3-5 minutes

**Manual deployment:**
```bash
# Build WASM
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/static/recipe.wasm cmd/wasm/main.go

# Deploy via Wrangler
wrangler pages deploy web/static --project-name recipe
```

### CLI Binaries (GitHub Releases)

**Release artifacts built for all platforms:**
- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64, arm64

**Create release:**
```bash
# Tag version
git tag v2.x.x
git push origin v2.x.x

# GitHub Actions automatically builds all binaries
```

## Useful Documentation

**Must-read for new contributors:**
- `docs/architecture.md` - Complete architecture with ADRs
- `docs/np3-format-specification.md` - NP3 binary format details
- `docs/parameter-mapping.md` - Cross-format parameter mapping
- `docs/PRD.md` - Product requirements and user stories

**For specific tasks:**
- `docs/browser-compatibility.md` - Web interface browser support
- `docs/performance-benchmarks.md` - Benchmark methodology
- `docs/format-compatibility-matrix.md` - Conversion accuracy matrix
- `docs/known-conversion-limitations.md` - Format-specific limitations

## Technology Stack

- **Language**: Go 1.25.1+ (leverages `go:wasmexport` for WASM)
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **TUI Framework**: Bubbletea v2 (charm.land/bubbletea/v2)
- **Web**: Vanilla JavaScript (ES6+) + WebAssembly
- **Deployment**: Cloudflare Pages (web), GitHub Releases (CLI/TUI)
- **Dependencies**: Minimal - only Cobra and Bubbletea (zero dependencies for core library)

## Legal and Compliance

**Reverse engineering disclosure:**
- NP3 format reverse-engineered through clean-room analysis
- Protected under DMCA Section 1201(f) for interoperability
- Recommended for private/personal use until full legal assessment

**Privacy commitment:**
- Zero server uploads (all processing local/client-side)
- No analytics or tracking
- See `docs/faq.md` for complete privacy FAQ
