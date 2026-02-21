# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) and similar AI agents when working with code in this repository.

## Project Overview

Recipe is a universal photo preset converter that converts between Nikon NP3, Adobe Lightroom XMP, and lrtemplate formats. The project provides two interfaces (CLI, Web) that share a common Go conversion engine.

**Key Characteristics:**
- **Privacy-first**: All processing happens locally (no server uploads)
- **Performance**: Sub-millisecond conversions (<100ms WASM target)
- **Accuracy**: 98%+ conversion fidelity via exact offset mapping (48 NP3 parameters)
- **Architecture**: Hub-and-spoke pattern with UniversalRecipe intermediate representation

**Note**: Capture One Costyle format support is disabled. TUI interface is archived in `.archive/tui/`.

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
XMP ──Parse──→ UniversalRecipe ──Generate──→ lrtemplate
lrtemplate ──Parse──→ UniversalRecipe ──Generate──→ NP3
DCP ──Parse──→ UniversalRecipe ──Generate──→ DCP  (hub integrated)
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
│   ├── np3/       # Nikon binary format (35 files)
│   ├── xmp/       # Adobe Lightroom XML
│   ├── lrtemplate/# Lightroom Classic Lua
│   ├── dcp/       # DNG Camera Profiles (fully hub integrated)
│   └── costyle/   # Capture One presets (DISABLED - 96 files)
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

testdata/          # 1,531 real sample files (73 NP3, 914 XMP, 544 lrtemplate)
docs/              # Core documentation (user guides, format specs, architecture)
.reverse-engineering/      # Reverse engineering artifacts
├── docs/          # NX Studio analysis, DCP research, findings
└── scripts/       # Python RE scripts (88 files)
.archive/          # Archived implementation artifacts
├── tui/           # Bubbletea TUI (archived)
└── docs/          # Epic/story/retrospective docs
scripts/           # Build utility scripts (benchmark, WASM build)
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
- `internal/formats/np3/parse.go` - Binary parsing logic (44KB)
- `internal/formats/np3/generate.go` - Binary generation logic (37KB)
- `internal/formats/np3/offsets.go` - Byte offset definitions
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

### Capture One Costyle Format

**New format support (Epic 8):**
- `.costyle` - Individual preset XML files
- `.costylepack` - ZIP bundles containing multiple presets
- Round-trip accuracy: 98.4%
- 96 implementation files in `internal/formats/costyle/`

### DNG Camera Profile (DCP) Format

> **STATUS: Fully integrated into converter hub ✅**

**Implementation Complete (Epic 9):**
- `internal/formats/dcp/parse.go` (3KB) - Parses DCP files to UniversalRecipe
- `internal/formats/dcp/generate.go` (5KB) - Generates DCP from UniversalRecipe
- `internal/formats/dcp/profile.go` (16KB) - Tone curve analysis, color matrix generation
- `internal/formats/dcp/tiff.go` (15KB) - Low-level TIFF/DNG binary handling
- Round-trip tests pass (`TestRoundTrip_DCP` in generate_test.go)

**What DCP Generate Creates:**
- Valid DNG Camera Profile with `IIRC` magic bytes
- 5-point tone curve from Exposure/Contrast/Highlights/Shadows
- Nikon Z f calibrated color matrices (ColorMatrix1 for 2856K, ColorMatrix2 for D65)
- Forward matrices (XYZ → camera RGB)
- 3D identity LUT (90×16×16 = 23,040 HSV entries)
- Baseline exposure offset (-0.15 EV for Nikon Z f)
- Profile name and camera model from metadata

**What DCP Parse Extracts:**
- Tone curve → Exposure, Contrast, Highlights, Shadows (via curve analysis)
- Color matrices stored in recipe.Metadata (with warning if non-identity)
- Baseline exposure offset
- Profile name

**Hub Integration Complete:**

The DCP format is now fully integrated into `converter.go`. You can convert:
- `converter.Convert(data, "np3", "dcp")` - NP3 → DCP
- `converter.Convert(data, "xmp", "dcp")` - XMP → DCP
- `converter.Convert(data, "dcp", "xmp")` - DCP → XMP
- `converter.Convert(data, "dcp", "np3")` - DCP → NP3
- Auto-detection via `IIRC` magic bytes (0x49, 0x49, 0x52, 0x43)

**Use Cases:**
- **NP3 → DCP**: Convert Nikon Picture Controls to Adobe Camera Profiles
- **XMP → DCP**: Create camera profiles from Lightroom presets
- **DCP → XMP**: Extract tone adjustments from existing camera profiles

**Warm Color Matrix Variant:**

A custom warm DCP generator is available to address Nikon vs Adobe color rendering differences:

- `internal/formats/dcp/profile_warm.go` - Custom warm Color Matrix 2
- `cmd/cli/generate_warm_dcp.go` - CLI utility for generation
- Triggered via `recipe.Metadata["use_warm_matrix"] = true`

**Warm Matrix Coefficients** (vs Adobe Standard):
```
Adobe Standard:              Warm Custom:
[1.1607  -0.4491  -0.0977]  →  [1.25   -0.35   0.08]
[-0.4522  1.2460   0.2304]      [-0.40   1.20   0.15]
[-0.0458  0.1519   0.7616]      [-0.02   0.10   0.85]
```

Key change: Blue→Red coefficient -0.0977 → +0.08 (+0.1777 warmth boost)

**Expected results**: 10-15% warmer rendering, closer to Nikon NX Studio output
**Target accuracy**: 92-96% Delta E (vs 85% with Adobe Standard)

**Generate warm DCP**:
```bash
go run cmd/cli/generate_warm_dcp.go
# Output: output/Nikon_Zf_Warm_Custom.dcp (277KB)
```

**Install in Lightroom**:
```bash
# Windows
copy output\Nikon_Zf_Warm_Custom.dcp %APPDATA%\Adobe\CameraRaw\CameraProfiles\
# Restart Lightroom, select "Nikon Z f Warm Custom" profile
```

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
- `PreviewModal.svelte` - Preset preview with adjustable parameters (40KB)
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
- `app.css` contains all global styles (17KB)

### Error Handling Pattern

**All conversion errors use the ConversionError type:**

```go
type ConversionError struct {
    Operation string  // "parse", "generate", "validate", "detect"
    Format    string  // "np3", "xmp", "lrtemplate"
    Cause     error   // Underlying error
}
```

## Nikon NX Studio Reverse Engineering

> **IMPORTANT**: Extensive reverse engineering documentation exists in `docs/` from Codex deep-dive sessions.

### Key Documents

| File | Description |
|------|-------------|
| `docs/NXSTUDIO_FINDINGS.md` | NX Studio directory analysis, Picture Control system |
| `docs/REVERSE_ENGINEERING_SUMMARY.md` | Comprehensive binary analysis summary |
| `docs/PICCON21_ANALYSIS.md` | PicCon21.bin calibration data structure |
| `docs/ADOBE_DCP_ANALYSIS.md` | Adobe DCP profile color matrices extracted |
| `docs/FINAL_CONCLUSIONS.md` | Root cause analysis and recommendations |
| `docs/reverse_engineering/` | 12 detailed analysis files |

### Critical Findings

**Nikon's Color Processing Architecture:**
1. **Polaris.dll** (3.5 MB) - Main color processing engine, ICC/ICM profile management
2. **Rome2.dll** (9.7 MB) - Rendering, color balance, LCH color space editing
3. **picture_control.n5m** (266 KB) - Picture Control service, custom curves

**PicCon21.bin Discovery:**
- Contains 5600x3728 14-bit RAW calibration image (23 MB compressed)
- Compression type 34713 (Nikon proprietary NEF)
- Used for color calibration chart reference

**Picture Control Types:**
- Standard: STANDARD, NEUTRAL, VIVID, PORTRAIT, LANDSCAPE, FLAT
- Monochrome: BW, FLAT_MONOCHROME, DEEP_TONE_MONOCHROME
- Creative: FLEXIBLE_COLOR, RICH_TONE_PORTRAIT (and 20+ numbered presets)

**Adobe DCP Color Matrices (Nikon Z f):**
```
Color Matrix 2 (D65 Daylight):
[ 1.1607  -0.4491  -0.0977 ]   ← Blue in red = NEGATIVE (cool shift)
[-0.4522   1.2460   0.2304 ]
[-0.0458   0.1519   0.7616 ]   ← Blue diagonal too low
```

**Root Cause of Color Mismatch:**
- Adobe uses negative blue→red coefficient (-0.0977) causing cooler reds
- Nikon likely uses positive value (+0.05) for warmer rendering
- Temperature compensation (+1000K) failed because it shifts white point, not matrix coefficients
- Solution: Create custom DCP with modified Color Matrix 2

### Analysis Scripts

Key Python scripts in `scripts/` for research:
- `reverse_engineer_nx.py` - Comprehensive DLL analysis
- `find_nikon_matrices.py` - Matrix pattern search
- `extract_dcp_lut.py` - DCP LUT extraction
- `process_piccon21_calibration.py` - Calibration data processing
- `frida_*.js` / `run_frida_*.bat` - Runtime hooking scripts

**NX Studio Binaries (copied to testdata/nxstudio/):**
```
testdata/nxstudio/
├── PicCon.bin          (legacy Picture Control DB, 9.2 MB)
├── PicCon21.bin        (modern Picture Control DB, 26 MB)
├── Polaris.dll         (color engine, 3.5 MB)
├── Rome2.dll           (render engine, 9.7 MB)
├── prm.bin             (camera calibration parameters, 17 MB)
└── Services/picture_control.n5m (Picture Control service, 266 KB)
```

## Testing Strategy

**1,531 real sample files across all formats:**
- 73 NP3 files
- 914 XMP files
- 544 lrtemplate files

**Round-trip testing validates conversion fidelity:**
- Full fidelity paths: NP3↔XMP, NP3↔lrtemplate, XMP↔lrtemplate
- Known limitations: XMP→NP3→XMP and lrtemplate→NP3→lrtemplate (some parameters unsupported by NP3)
- Costyle round-trip: 98.4% accuracy
- DCP round-trip: Tests pass (`TestRoundTrip_DCP`)

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

## Important Constraints

### Format Limitations

**NP3 format has limited parameter support compared to XMP/lrtemplate:**
- ❌ Not supported: Vibrance, Temperature/Tint, Grain Size/Roughness, Vignette, Custom Tone Curves (Point Curves and Parametric Curves)
- ✅ Well supported: Exposure, Contrast, Saturation, Sharpness, Highlights, Shadows, Whites, Blacks, Clarity, HSL Color, Color Grading

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
- See `docs/implementation-artifacts/sprint-change-proposal-2025-12-24.md` for full technical analysis

**NP3 Format Variants** (discovered via analysis of 160 samples):
- **392 bytes**: Minimal/compact format (chunk-based encoding) - 12 files
- **466 bytes**: Grain parameters variant - 6 files
- **480 bytes**: Standard format (direct offset mapping) - 12 files - PRIMARY IMPLEMENTATION
- **978-1,140 bytes**: Extended formats with metadata/descriptions (56+ files)
  - KOLORA format (1,140 bytes): Maximum parameters with full description text
  - Temperature/Tint/Vibrance likely present in extended variants (unconfirmed)

**Temperature/Tint/Vibrance Investigation Results**:
- ❌ NOT FOUND in 480-byte standard format after analyzing 160 samples
- Statistical analysis: Only 1 high-variance offset found (0xF2 = MidRangeSharpening)
- Hypothesis: These parameters may exist in 978-1,140 byte extended variants or use proprietary encoding

**Always test round-trip conversions when working with NP3:**
- XMP → NP3 → XMP may lose parameters (~85% fidelity)
- NP3 → XMP → NP3 preserves all (~98% fidelity)

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

## Useful Documentation

**Must-read for new contributors:**
- `docs/architecture.md` - Complete architecture with ADRs
- `docs/np3-format-specification.md` - NP3 binary format details
- `docs/parameter-mapping.md` - Cross-format parameter mapping
- `docs/PRD.md` - Product requirements and user stories

**Reverse Engineering:**
- `docs/FINAL_CONCLUSIONS.md` - Root cause analysis for color issues
- `docs/REVERSE_ENGINEERING_SUMMARY.md` - Binary analysis details
- `docs/ADOBE_DCP_ANALYSIS.md` - Adobe profile color matrices
- `docs/reverse_engineering/REVERSE_ENGINEERING_REPORT.md` - Full report

**For specific tasks:**
- `docs/browser-compatibility.md` - Web interface browser support
- `docs/performance-benchmarks.md` - Benchmark methodology
- `docs/format-compatibility-matrix.md` - Conversion accuracy matrix
- `docs/known-conversion-limitations.md` - Format-specific limitations

## Legal and Compliance

**Reverse engineering disclosure:**
- NP3 format reverse-engineered through clean-room analysis
- Protected under DMCA Section 1201(f) for interoperability
- Recommended for private/personal use until full legal assessment

**Privacy commitment:**
- Zero server uploads (all processing local/client-side)
- No analytics or tracking
- See `docs/faq.md` for complete privacy FAQ

## Quick Reference for Agents

### Common Tasks

| Task | Command |
|------|---------|
| Run tests | `go test ./...` |
| Build WASM | `make wasm` |
| Start web dev server | `cd web && npm run dev` |
| Check NP3 parsing | `go test ./internal/formats/np3/` |
| Check DCP parsing | `go test ./internal/formats/dcp/` |
| View coverage | `make coverage-html` |

### Important Files by Task

| Task | Files |
|------|-------|
| Add new parameter | `internal/models/recipe.go`, all format parsers |
| Fix NP3 parsing | `internal/formats/np3/parse.go`, `offsets.go` |
| Integrate DCP | `internal/converter/converter.go`, `internal/formats/dcp/` |
| Web UI changes | `web/src/lib/components/*.svelte`, `web/src/app.css` |
| WASM exports | `cmd/wasm/main.go`, `web/src/lib/wasm.js` |
| Preview filters | `web/src/lib/svg-logic.js`, `SVGFilters.svelte` |

### Files to Never Touch

- `testdata/` samples (reference files for tests)
- `testdata/nxstudio/` binaries (copyrighted Nikon files)
- Binary analysis output in `docs/reverse_engineering/*.json`
