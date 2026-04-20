# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) and similar AI agents when working with code in this repository.

## Project Overview

Recipe is a photo preset converter for Nikon NP3 (Picture Control), Adobe Lightroom XMP, Fujifilm CoStyle, and Adobe DCP formats. A single Go conversion engine is shared across a CLI, a web app, and a VSCode extension — all in one monorepo.

**Key Characteristics:**
- **Privacy-first**: All processing happens locally (no server uploads)
- **Performance**: Sub-millisecond conversions (<100ms WASM target)
- **Accuracy**: ~98% NP3↔XMP fidelity via exact offset mapping (56 NP3 parameters via dual-mode extraction)
- **Architecture**: Hub-and-spoke pattern with `UniversalRecipe` intermediate representation

## Repository Layout

This is a **monorepo** managed with **bun workspaces** (see root `package.json`).

| Workspace           | Purpose                                                                          |
| ------------------- | -------------------------------------------------------------------------------- |
| `cmd/`, `internal/` | Go module `github.com/justin/recipe` — conversion engine, CLI, WASM, debug tools |
| `web/`              | Public website + converter at recipe.shuttercoach.app (Astro + Svelte 5)         |
| `extension/`        | VSCode extension `@recipe/extension` — custom editor for `.np3` files            |
| `webview/`          | `@recipe/webview` — Svelte 5 UI bundled into the VSCode extension                |
| `packages/ui`       | `@recipe/ui` — shared Svelte 5 + Tailwind 4 component library                    |

## Technology Stack

| Component        | Technology                          | Notes                                                   |
| ---------------- | ----------------------------------- | ------------------------------------------------------- |
| Core Engine      | Go 1.25.1                           | Module `github.com/justin/recipe`; uses `go:wasmexport` |
| CLI Framework    | Cobra                               | `github.com/spf13/cobra`                                |
| Web Frontend     | Astro 6 + Svelte 5 + Tailwind 4     | Islands architecture, SSG                               |
| VSCode Extension | TypeScript + Svelte 5 webview       | Bundles WASM + `np3tool` binary                         |
| Shared UI        | `@recipe/ui` (Svelte 5, Tailwind 4) | Consumed by `web/` and `webview/`                       |
| WASM             | Go WebAssembly                      | Compiled from `cmd/wasm/`                               |
| Package Manager  | bun (workspaces)                    | `bun.lock` at repo root                                 |
| Linter/Formatter | Biome                               | `biome.json` at repo root                               |
| Deployment       | Cloudflare Pages                    | Auto-deploy on push to main                             |

## Essential Commands

### Building

```bash
# ---- Go (CLI / WASM) ----
make cli               # Build CLI with version injection → ./recipe
make cli-all           # Cross-compile for linux/darwin/windows (amd64/arm64) → ./bin/
make wasm              # Build stripped WASM for web → web/public/recipe.wasm
make wasm-dev          # WASM without -s -w (readable stack traces for debugging)

# ---- Monorepo (bun / workspaces) ----
bun install            # Install all workspace dependencies
bun run build          # Full build: np3tool + WASM + webview + extension
bun run dev            # Run `dev` script across every workspace in parallel
bun run lint           # Biome check across extension/, webview/, web/
bun run lint:fix       # Biome autofix

# ---- Targeted workspace builds ----
bun run build:go          # Build np3tool binary → extension/bin/np3tool
bun run build:wasm-ext    # Build WASM for VSCode extension webview
bun run build:wasm-web    # Build WASM for the public website
bun run build:extension   # Package the VSCode extension
bun run build:webview     # Build the Svelte webview bundle
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

### Web Development (Astro + Svelte 5)

```bash
# Navigate to web directory
cd web

# Install dependencies
npm install

# Start dev server (hot reload) — runs `astro dev`
npm run dev

# Production build — runs `astro build`, output in web/dist/
npm run build

# Preview production build
npm run preview

# Type-check (astro check + svelte-check)
npm run check

# Run Vitest component tests
npm test
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
├── cli/                 # Cobra CLI (primary user-facing binary)
├── wasm/                # WASM export entry point for web + VSCode webview
├── np3tool/             # Helper binary bundled into the VSCode extension
├── debug_curve/         # Dev tool for tone-curve investigation
├── test_curve/          # Ad-hoc curve test harness
└── test_compensation/   # Ad-hoc exposure-compensation test harness

internal/
├── converter/     # Core conversion engine — Convert(), round-trip tests, benchmarks
├── formats/       # Format parsers/generators (one package per format)
│   ├── np3/       # Nikon Picture Control (binary) — primary target
│   ├── xmp/       # Adobe Lightroom (XML, crs: namespace)
├── models/        # UniversalRecipe + parameter definitions, builders, validation, warnings
├── inspect/       # Parameter inspection and diff tools (used by `recipe inspect`)
├── lut/           # LUT table handling (tone curves, color grading)
├── verify/        # Validation helpers for conversion correctness
├── apperr/        # Application error types (wraps conversion errors)
├── utils/         # Shared low-level helpers
└── testutil/      # Test helpers for fixture loading

web/               # Astro + Svelte 5 frontend (Tailwind 4)
├── src/
│   ├── pages/
│   │   └── index.astro           # Main page (SSG)
│   ├── layouts/
│   │   └── Layout.astro          # Root layout with SEO meta
│   ├── components/               # Mixed .astro + .svelte
│   │   ├── NavBar.astro
│   │   ├── Footer.astro
│   │   ├── Explainer.astro
│   │   ├── FAQ.astro
│   │   ├── AuroraBackground.svelte
│   │   ├── ConversionCard.svelte # Main upload/convert flow
│   │   ├── CurveEditor.svelte    # Tone curves editor
│   │   └── EditorView.svelte     # Metadata editing modal
│   ├── lib/
│   │   ├── wasm.svelte.ts        # WASM initialization (runes)
│   │   ├── converter.svelte.ts   # WASM conversion wrapper
│   │   ├── stores.svelte.ts      # State (nanostores + runes)
│   │   ├── shared-stores.ts
│   │   ├── format-detector.ts
│   │   ├── parameter-counter.ts
│   │   ├── image-analysis.js
│   │   ├── image-loader.ts
│   │   ├── preview-logic.js
│   │   └── utils.ts              # clsx / tailwind-merge helpers
│   ├── styles/                   # Tailwind 4 + theme tokens
│   └── __tests__/                # Vitest component tests
├── public/                       # Static assets (WASM binary)
├── remotion/                     # Remotion project for OG images/video
├── astro.config.mjs              # Astro config (sitemap, Svelte, Tailwind)
├── svelte.config.js
├── vitest.config.ts
└── wrangler.jsonc                # Cloudflare Pages config

docs/              # Core documentation
├── architecture.md
├── known-conversion-limitations.md
├── np3-format-specification.md
└── parameter-mapping.md

extension/         # VSCode extension (@recipe/extension) — custom .np3 editor
├── src/           # TypeScript extension host code
└── bin/np3tool    # Bundled Go helper (built via `bun run build:go`)

webview/           # @recipe/webview — Svelte 5 UI for the VSCode custom editor
packages/
└── ui/            # @recipe/ui — shared Svelte 5 + Tailwind 4 component library

testdata/          # Top-level fixtures (nksc/, np3/, visual-regression/)
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
- Magic bytes: "NCP" (0x4E, 0x43, 0x50) at offset 0
- **Variable size**: 392 / 466 / 480 / 978–1,140 bytes (see variants below). Minimum accepted: 300 bytes.
- Parameters stored at fixed byte offsets (see `internal/formats/np3/offsets.go`, documented in `docs/np3-format-specification.md`)
- Uses signed byte normalization (0x80 = zero point, ±0x7F ≈ ±100%)
- **Dual-mode extraction** (`internal/formats/np3/parse.go`):
  1. Exact offset extraction — 56 parameters at known byte positions (primary, ~100% on 480-byte files)
  2. Heuristic fallback — scans raw bytes when exact extraction fails (~95% on variants)

**Key files:**
- `internal/formats/np3/parse.go` - Binary parsing logic
- `internal/formats/np3/generate.go` - Binary generation logic
- `internal/formats/np3/offsets.go` - Byte offset definitions
- `docs/np3-format-specification.md` - Complete format documentation

### XMP Format

**XMP (Adobe Lightroom / Lightroom Classic):**
- XML sidecar format with `crs:` namespace for camera-raw adjustments
- Parsed with Go's `encoding/xml` (stdlib) — no external dependencies
- Full parameter support (50+ fields including tone curves, HSL, color grading, split-toning)

### Web Frontend (Astro + Svelte 5)

**Stack details:**
- **Meta-framework**: Astro 6 (SSG) — `src/pages/index.astro` is the single page
- **Islands**: Svelte 5.53+ components with runes (`$state`, `$derived`, `$effect`)
- **Styling**: Tailwind CSS 4 via `@tailwindcss/vite`, plus `bits-ui` primitives
- **State**: `nanostores` + Svelte 5 runes (`*.svelte.ts` files in `src/lib/`)
- **Icons**: `@lucide/svelte`
- **Testing**: Vitest + `@testing-library/svelte` (jsdom)

**Key Components:**
- `ConversionCard.svelte` - Upload, detect, convert, download flow
- `EditorView.svelte` - Glassmorphism modal for editing NP3 name/description + tone curves
- `CurveEditor.svelte` - Visual tone curves editor
- `AuroraBackground.svelte` - Animated background
- `NavBar.astro` / `Footer.astro` / `Explainer.astro` / `FAQ.astro` - Static sections

**WASM Integration:**
- `lib/wasm.svelte.ts` - Loads `/recipe.wasm` and initializes the Go runtime
- `lib/converter.svelte.ts` - Wraps exported WASM functions
- Exports: `convert()`, `generate()`, `extractFullRecipe()`
- WASM binary served from `web/public/recipe.wasm`

**Tailwind / design system:**
- Tailwind 4 theme tokens (e.g. `text-foreground`, `text-interactive`, `glass-regular`) defined in `src/styles/`
- Use `clsx` + `tailwind-merge` (`lib/utils.ts`) for conditional classes
- Glassmorphism utilities for cards/modals
- Shared primitives live in `@recipe/ui` (`packages/ui/src/`) — prefer importing from there before inlining new components

### VSCode Extension

**Goal:** open `.np3` files inside VSCode with a visual editor (custom editor with `viewType: "recipe.np3Editor"`).

**Architecture:**
- **Extension host** (`extension/src/`, TypeScript) — registers the custom editor, shells out to `np3tool` for any native-only operations, and hosts the webview.
- **Webview** (`webview/src/`, Svelte 5) — the UI the user sees inside the editor tab. Built with Vite into a static bundle the extension loads.
- **Native helper** (`extension/bin/np3tool`, built from `cmd/np3tool/`) — used when WASM isn't sufficient.
- **WASM** (`extension/dist/webview/recipe.wasm`) — same engine as the website, wired into the webview for fast in-process conversion.

**Typical dev loop:**
```bash
bun install
bun run build        # Builds np3tool + WASM + webview + extension once
# Then "Run Extension" from VSCode's Run & Debug panel
```

The webview consumes `@recipe/ui` for shared components — changes to `packages/ui` affect both the website and the extension.

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

**Go 1.25+ with `go:wasmexport` directive** (entry point: `cmd/wasm/main.go`):

```go
//go:wasmexport convertPreset
func convertPreset(inputPtr, inputLen uint32, srcFormat, dstFormat string) (uint32, uint32, string)
```

**Key details:**
- Direct memory access (zero reflection overhead)
- Returns `(outputPtr, outputLen, errorMsg)` tuple
- Binary size: ~4.0 MB stripped, ~1.13 MB gzipped (built with `-ldflags="-s -w"`)
- Target: <100ms conversions (actual: 0.003–0.079 ms)
- Consumed by **both** the public website (`web/public/recipe.wasm`) and the VSCode extension webview (`extension/dist/webview/recipe.wasm`) — same binary, different `wasm_exec.js` copies

## Important Constraints

### Format Limitations

**NP3 format has limited parameter support compared to XMP:**
- Not supported: Vibrance, Temperature/Tint, Grain Size/Roughness, Vignette, Custom Tone Curves (Point Curves and Parametric Curves)
- Well supported: Exposure, Contrast, Saturation, Sharpness, Highlights, Shadows, Whites, Blacks, Clarity, HSL Color, Color Grading

**IMPORTANT - XMP → NP3 Tone Adjustment Strategy:**

NP3 has a **critical limitation**: You can use EITHER tone curve OR basic tone parameters, but NOT BOTH simultaneously.

**Our conversion strategy: Direct Parameter Mapping (No Curve Generation)**

When converting XMP → NP3, we use direct parameter mapping instead of generating custom tone curves:

| XMP Parameter        | NP3 Parameter | Byte Offset | Range        |
| -------------------- | ------------- | ----------- | ------------ |
| `crs:Contrast2012`   | Contrast      | 0x110       | -100 to +100 |
| `crs:Highlights2012` | Highlights    | 0x11A       | -100 to +100 |
| `crs:Shadows2012`    | Shadows       | 0x124       | -100 to +100 |
| `crs:Whites2012`     | White Level   | 0x12E       | -100 to +100 |
| `crs:Blacks2012`     | Black Level   | 0x138       | -100 to +100 |

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
2. Builds Astro production bundle (`astro build`)
3. Deploys `web/dist/` to Cloudflare Pages (config in `web/wrangler.jsonc`)
4. Live at https://recipe.shuttercoach.app in 3-5 minutes

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
- `README.md`, `CONTRIBUTING.md`, `CHANGELOG.md`, `SECURITY.md` at repo root
- Per-workspace READMEs: `web/README.md`, `extension/README.md`, `webview/README.md`

## Legal and Compliance

**Format analysis disclosure:**
- NP3 format analyzed through clean-room methods for interoperability
- Protected under DMCA Section 1201(f) for interoperability

**Privacy commitment:**
- Zero server uploads (all processing local/client-side)
- No analytics or tracking

## Quick Reference for Agents

### Common Tasks

| Task                    | Command                                 |
| ----------------------- | --------------------------------------- |
| Install JS deps         | `bun install` (from repo root)          |
| Run Go tests            | `go test ./...`                         |
| Run all JS/Svelte tests | `bun run --filter '*' test`             |
| Lint JS/TS/Svelte       | `bun run lint`                          |
| Build WASM (web)        | `make wasm` or `bun run build:wasm-web` |
| Build WASM (extension)  | `bun run build:wasm-ext`                |
| Build VSCode extension  | `bun run build:extension`               |
| Start web dev server    | `cd web && npm run dev`                 |
| Check NP3 parsing       | `go test ./internal/formats/np3/`       |
| View Go coverage        | `make coverage-html`                    |

### Important Files by Task

| Task                | Files                                                                                                                       |
| ------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| Add a new parameter | `internal/models/recipe.go`, `internal/models/parameter_definitions.go`, each `internal/formats/*/parse.go` + `generate.go` |
| Fix NP3 parsing     | `internal/formats/np3/parse.go`, `internal/formats/np3/offsets.go`                                                          |
| Add a new format    | New package under `internal/formats/<name>/`, then update `internal/converter/converter.go`                                 |
| CLI command changes | `cmd/cli/` (`root.go`, `convert.go`, `batch.go`, `inspect.go`, `diff.go`, `format.go`)                                      |
| Web UI changes      | `web/src/components/*.{svelte,astro}`, `web/src/styles/`, `packages/ui/src/`                                                |
| VSCode extension    | `extension/src/`, `webview/src/`                                                                                            |
| Shared components   | `packages/ui/src/` (consumed by web + webview)                                                                              |
| WASM exports        | `cmd/wasm/main.go`, `web/src/lib/wasm.svelte.ts`                                                                            |
| Preview filters     | `web/src/lib/preview-logic.js`, `web/src/lib/image-analysis.js`                                                             |
| Error types         | `internal/converter/error.go`, `internal/apperr/`                                                                           |
