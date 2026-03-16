# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) and similar AI agents when working with code in this repository.

## Project Overview

Recipe is a photo preset converter that converts between Nikon NP3 and Adobe Lightroom XMP formats. The project provides two interfaces (CLI, Web) that share a common Go conversion engine.

**Key Characteristics:**
- **Privacy-first**: All processing happens locally (no server uploads)
- **Performance**: Sub-millisecond conversions (<100ms WASM target)
- **Accuracy**: 98%+ conversion fidelity via exact offset mapping (48 NP3 parameters)
- **Architecture**: Hub-and-spoke pattern with UniversalRecipe intermediate representation

## Technology Stack

| Component | Technology | Notes |
|-----------|------------|-------|
| Core Engine | Go 1.25.1+ | `go:wasmexport` for WASM |
| CLI Framework | Cobra | github.com/spf13/cobra |
| **Web Frontend** | **Vite + Svelte 5** | `web/` directory |
| WASM | Go WebAssembly | Compiled from `cmd/wasm/` |
| Deployment | Cloudflare Pages | Auto-deploy on push to main |

## Essential Commands

### Building

```bash
# Build CLI for current platform
make cli
# or: go build -o recipe cmd/cli/*.go

# Build WASM for web interface (with size optimization)
make wasm
# or: GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/public/recipe.wasm cmd/wasm/main.go

# Build for all platforms
make cli-all
```

### Testing

```bash
# Run all tests (uses committed fixtures in package-level testdata/ dirs)
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for specific package
go test ./internal/formats/np3/
go test ./internal/converter/

# Run specific test
go test -run TestRoundTrip_NP3_XMP ./internal/converter/

# Run with coverage
make coverage

# Generate HTML coverage report
make coverage-html
```

### Web Development (Vite + Svelte 5)

```bash
# Navigate to web directory
cd web

# Install dependencies
npm install

# Start development server (hot reload)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
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

## Architecture

### Hub-and-Spoke Conversion Pattern

All conversions flow through a central `UniversalRecipe` intermediate representation:

```
NP3 ──Parse──→ UniversalRecipe ──Generate──→ XMP
XMP ──Parse──→ UniversalRecipe ──Generate──→ NP3
```

**Why this matters:**
- All conversions use the same API: `converter.Convert(input, from, to)`
- Parameter mapping logic is centralized in UniversalRecipe

### Core Conversion API

**Single entry point for all conversions:**

```go
// internal/converter/converter.go
func Convert(input []byte, from, to string) ([]byte, error)
```

**Critical Rules:**
1. All interfaces (CLI, WASM) MUST use `converter.Convert()` - never call format parsers directly
2. All conversion errors MUST be wrapped in `ConversionError` type
3. The API is thread-safe and stateless

### Package Structure

```
cmd/
├── cli/           # Cobra CLI application
└── wasm/          # WASM export entry point

internal/
├── converter/     # Core conversion engine (single source of truth)
├── formats/       # Format parsers/generators
│   ├── np3/       # Nikon binary format
│   └── xmp/       # Adobe Lightroom XML
├── models/        # UniversalRecipe data structures
├── inspect/       # Parameter inspection and diff tools
├── lut/           # LUT table handling
└── testutil/      # Test utilities

web/               # Vite + Svelte 5 frontend
├── src/
│   ├── App.svelte         # Main application component
│   ├── app.css            # Global styles
│   ├── lib/
│   │   ├── components/    # Svelte components (16 files)
│   │   ├── stores.js      # Svelte stores (state management)
│   │   ├── wasm.js        # WASM initialization
│   │   ├── converter.js   # WASM conversion wrapper
│   │   ├── format-detector.js
│   │   ├── parameter-extractor.js
│   │   ├── svg-logic.js   # Preview filter logic
│   │   └── preview-logic.js
│   └── main.js            # Entry point
├── public/                # Static assets (WASM binary)
├── vite.config.js         # Vite configuration
└── svelte.config.js       # Svelte configuration

docs/              # Core documentation
├── architecture.md
├── known-conversion-limitations.md
├── np3-format-specification.md
└── parameter-mapping.md

extension/         # VSCode extension (work in progress)
webview/           # Webview UI for VSCode extension
packages/          # Shared packages
```

### Format Package Pattern

**Every format package follows identical structure:**

```
internal/formats/{format}/
├── parse.go          # Parse([]byte) (*UniversalRecipe, error)
├── generate.go       # Generate(*UniversalRecipe) ([]byte, error)
├── {format}_test.go  # Table-driven tests with real samples
└── testdata/         # Committed test fixtures
```

**When adding a new format:**
1. Copy the structure from an existing format package
2. Implement Parse() and Generate() functions
3. Add test fixtures to `internal/formats/{format}/testdata/`
4. Update `converter.Convert()` switch statements

## Key Implementation Details

### NP3 Binary Format

The NP3 (Nikon Picture Control) format is a proprietary binary format analyzed through clean-room methods for interoperability.

**Critical implementation notes:**
- Magic bytes: "NCP" (0x4E, 0x43, 0x50)
- Fixed size: 1024 bytes
- Parameters stored at fixed byte offsets (documented in `docs/np3-format-specification.md`)
- Uses signed byte normalization (128 = zero point, ±127 = ±100%)
- Exact offset mapping for 48 parameters

**Key files:**
- `internal/formats/np3/parse.go` - Binary parsing logic
- `internal/formats/np3/generate.go` - Binary generation logic
- `internal/formats/np3/offsets.go` - Byte offset definitions
- `docs/np3-format-specification.md` - Complete format documentation

### XMP Format

**XMP (Adobe Lightroom CC):**
- XML format with `crs:` namespace for adjustments
- Uses ElementTree for parsing
- Full parameter support (50+ fields)

### Web Frontend (Vite + Svelte 5)

**Stack details:**
- **Framework**: Svelte 5.43.8 (latest stable with runes)
- **Build Tool**: Vite 7.2.4
- **Styling**: Vanilla CSS with CSS variables
- **State Management**: Svelte stores (`web/src/lib/stores.js`)

**Key Components:**
- `App.svelte` - Root component with WASM initialization
- `UploadZone.svelte` - Drag-and-drop file upload
- `FileList.svelte` / `FileCard.svelte` - File management
- `ActionPanel.svelte` - Conversion controls
- `PreviewModal.svelte` - Preset preview with adjustable parameters
- `SVGFilters.svelte` - CSS-based preset preview filters
- `Histogram.svelte` - Image histogram visualization

**WASM Integration:**
- `wasm.js` - Initializes Go WASM module
- `converter.js` - Wraps WASM conversion functions
- Exports: `convert()`, `generate()`, `extractFullRecipe()`
- Event: `wasmReady` fired when module is ready

**CSS Architecture:**
- Glassmorphism design system
- CSS custom properties for theming
- Mobile-responsive grid layouts
- `app.css` contains all global styles

### Error Handling Pattern

**All conversion errors use the ConversionError type:**

```go
type ConversionError struct {
    Operation string  // "parse", "generate", "validate", "detect"
    Format    string  // "np3", "xmp"
    Cause     error   // Underlying error
}
```

## Testing Strategy

**Committed test fixtures in package-level testdata/ directories:**
- 3-5 representative fixtures per format per package
- Existing `curve_tests/` synthetic fixtures in `internal/formats/np3/testdata/`

**Round-trip testing validates conversion fidelity:**
- Full fidelity path: NP3↔XMP
- Known limitations: XMP→NP3→XMP (some parameters unsupported by NP3)

**Test execution:**
- Tests complete in <2 seconds (parallel execution)
- All tests use table-driven pattern with real files
- No external data downloads needed — `go test ./...` works out of the box

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

## Important Constraints

### Format Limitations

**NP3 format has limited parameter support compared to XMP:**
- Not supported: Vibrance, Temperature/Tint, Grain Size/Roughness, Vignette, Custom Tone Curves (Point Curves and Parametric Curves)
- Well supported: Exposure, Contrast, Saturation, Sharpness, Highlights, Shadows, Whites, Blacks, Clarity, HSL Color, Color Grading

**IMPORTANT - XMP → NP3 Tone Adjustment Strategy:**

NP3 has a **critical limitation**: You can use EITHER tone curve OR basic tone parameters, but NOT BOTH simultaneously.

**Our conversion strategy: Direct Parameter Mapping (No Curve Generation)**

When converting XMP → NP3, we use direct parameter mapping instead of generating custom tone curves:

| XMP Parameter | NP3 Parameter | Byte Offset | Range |
|---------------|---------------|-------------|-------|
| `crs:Contrast2012` | Contrast | 0x110 | -100 to +100 |
| `crs:Highlights2012` | Highlights | 0x11A | -100 to +100 |
| `crs:Shadows2012` | Shadows | 0x124 | -100 to +100 |
| `crs:Whites2012` | White Level | 0x12E | -100 to +100 |
| `crs:Blacks2012` | Black Level | 0x138 | -100 to +100 |

**What Gets Lost in XMP → NP3 Conversion:**
- XMP Parametric Curve Sliders (`ToneCurveShadows`, `ToneCurveDarks`, `ToneCurveLights`, `ToneCurveHighlights`)
- XMP Custom Point Curves (`PointCurve`, `PointCurveRed`, `PointCurveGreen`, `PointCurveBlue`)

**Mitigation:**
- Curve data is preserved in `recipe.Metadata` for round-trip fidelity
- Conversion warnings inform users when curve data is lost
- This approach covers 95%+ of real-world XMP presets (most use basic adjustments, not custom curves)
- Users can create custom curves directly in NX Studio if needed

**Why We Don't Generate Curves:**
- NP3 cannot use both curves AND basic parameters simultaneously
- Direct parameter mapping is simpler, more accurate, and faster
- Previous curve generation attempts (257-entry LUT) failed to achieve acceptable visual fidelity

**NP3 Format Variants** (discovered via analysis of 160 samples):
- **392 bytes**: Minimal/compact format (chunk-based encoding) - 12 files
- **466 bytes**: Grain parameters variant - 6 files
- **480 bytes**: Standard format (direct offset mapping) - 12 files - PRIMARY IMPLEMENTATION
- **978-1,140 bytes**: Extended formats with metadata/descriptions (56+ files)
  - KOLORA format (1,140 bytes): Maximum parameters with full description text
  - Temperature/Tint/Vibrance likely present in extended variants (unconfirmed)

**Temperature/Tint/Vibrance Investigation Results**:
- NOT FOUND in 480-byte standard format after analyzing 160 samples
- Statistical analysis: Only 1 high-variance offset found (0xF2 = MidRangeSharpening)
- Hypothesis: These parameters may exist in 978-1,140 byte extended variants or use proprietary encoding

**Always test round-trip conversions when working with NP3:**
- XMP → NP3 → XMP may lose parameters (~85% fidelity)
- NP3 → XMP → NP3 preserves all (~98% fidelity)

### Performance Requirements

**All conversions must meet these targets:**
- WASM: <100ms (actual: 0.003-0.079ms)
- CLI: <20ms (actual: 0.003-0.079ms)
- Batch (100 files): <2s (actual: 37ms)
- Memory: <4096 B/op (actual: 8,890-29,026 B/op)

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
2. Builds Vite+Svelte production bundle
3. Deploys `web/dist/` to Cloudflare Pages
4. Live at https://recipe.pages.dev in 3-5 minutes

**Manual deployment:**
```bash
# Build WASM
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/public/recipe.wasm cmd/wasm/main.go

# Build web frontend
cd web && npm run build

# Deploy via Wrangler
wrangler pages deploy web/dist --project-name recipe
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

## Documentation

- `docs/architecture.md` - Architecture decisions
- `docs/np3-format-specification.md` - NP3 binary format details
- `docs/parameter-mapping.md` - Cross-format parameter mapping
- `docs/known-conversion-limitations.md` - Format-specific limitations

## Legal and Compliance

**Format analysis disclosure:**
- NP3 format analyzed through clean-room methods for interoperability
- Protected under DMCA Section 1201(f) for interoperability

**Privacy commitment:**
- Zero server uploads (all processing local/client-side)
- No analytics or tracking

## Quick Reference for Agents

### Common Tasks

| Task | Command |
|------|---------|
| Run tests | `go test ./...` |
| Build WASM | `make wasm` |
| Start web dev server | `cd web && npm run dev` |
| Check NP3 parsing | `go test ./internal/formats/np3/` |
| View coverage | `make coverage-html` |

### Important Files by Task

| Task | Files |
|------|-------|
| Add new parameter | `internal/models/recipe.go`, all format parsers |
| Fix NP3 parsing | `internal/formats/np3/parse.go`, `offsets.go` |
| Web UI changes | `web/src/lib/components/*.svelte`, `web/src/app.css` |
| WASM exports | `cmd/wasm/main.go`, `web/src/lib/wasm.js` |
| Preview filters | `web/src/lib/svg-logic.js`, `SVGFilters.svelte` |
