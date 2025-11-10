# Architecture - Recipe Path A Enhancements

## Executive Summary

Path A extends Recipe's proven hub-and-spoke conversion architecture to support two new formats (Capture One .costyle and DCP camera profiles) while enhancing the web interface with modern UX patterns and visual preset preview. The architecture maintains Recipe's core principles—privacy-first processing, zero external dependencies (except Google's TIFF library), and <100ms conversion performance—while adding professional polish through CSS-based preview and responsive design.

This architecture document defines the technical decisions and consistency rules for all AI agents implementing Path A's four epics.

## Decision Summary

| Category | Decision | Version | Affects Epics | Rationale |
| -------- | -------- | ------- | ------------- | --------- |
| Language | Go | 1.25.1+ | All | Established by Recipe core |
| XML Parsing | encoding/xml | Go stdlib | Epic 1, 2 | Zero dependencies, proven with XMP |
| TIFF Library | github.com/google/tiff | latest | Epic 2 | Complete TIFF support, Google-maintained |
| ZIP Handling | archive/zip | Go stdlib | Epic 1 | Standard library sufficient |
| Web Approach | Vanilla JS + WASM | ES6+ | Epic 3, 4 | No framework bloat, established |
| CSS Organization | Modular CSS files | - | Epic 3, 4 | main, components, layout, preview |
| Image Format | WebP | - | Epic 4 | Best compression for photos |
| Error Handling | fmt.Errorf with %w | Go stdlib | All | Wrapped context, clear messages |
| Test Organization | Co-located testdata/ | - | Epic 1, 2 | Per-format test samples |
| Deployment | Cloudflare Pages + GitHub Releases | - | All | Established by Recipe core |

## Project Structure

```
recipe/
├── cmd/
│   ├── recipe/              # CLI binary (existing)
│   └── recipe-tui/          # TUI binary (existing)
├── internal/
│   ├── converter/           # Hub-and-spoke engine (existing)
│   ├── formats/
│   │   ├── np3/             # Existing: Nikon Picture Control
│   │   ├── xmp/             # Existing: Adobe Lightroom
│   │   ├── lrtemplate/      # Existing: Lightroom Classic
│   │   ├── costyle/         # NEW: Capture One Format
│   │   │   ├── parse.go         # Parse .costyle XML → UniversalRecipe
│   │   │   ├── generate.go      # Generate .costyle XML ← UniversalRecipe
│   │   │   ├── types.go         # Capture One-specific types
│   │   │   ├── pack.go          # ZIP handling for .costylepack bundles
│   │   │   └── testdata/        # Sample .costyle files (Etsy/marketplaces)
│   │   │       ├── sample1.costyle
│   │   │       ├── sample2.costyle
│   │   │       ├── bundle.costylepack
│   │   │       └── README.md    # Source attribution
│   │   └── dcp/             # NEW: DNG Camera Profile
│   │       ├── parse.go         # Parse DCP TIFF → UniversalRecipe
│   │       ├── generate.go      # Generate DCP TIFF ← UniversalRecipe
│   │       ├── types.go         # DCP-specific types (matrices, curves)
│   │       ├── tiff.go          # TIFF tag reading/writing helpers
│   │       ├── profile.go       # XML camera profile parsing/generation
│   │       └── testdata/        # Sample DCP files (Adobe sources)
│   │           ├── camera1.dcp
│   │           ├── camera2.dcp
│   │           └── README.md    # Adobe source, license info
│   └── universal/           # UniversalRecipe intermediate (existing)
├── web/                     # ENHANCED: Web interface
│   ├── index.html           # Redesigned landing page
│   ├── css/
│   │   ├── main.css             # Global styles, CSS variables, color palette
│   │   ├── components.css       # Reusable components (badges, buttons, cards)
│   │   ├── layout.css           # Grid system, responsive breakpoints
│   │   └── preview.css          # Preview modal and before/after slider
│   ├── js/
│   │   ├── app.js               # Main application logic, initialization
│   │   ├── converter.js         # WASM conversion interface (existing)
│   │   ├── preview.js           # Preview modal and CSS filter mapping
│   │   ├── upload.js            # Batch upload, drag-drop, file handling
│   │   └── utils.js             # Shared utilities, format detection
│   ├── images/
│   │   ├── preview-portrait.webp    # Reference image: Portrait
│   │   ├── preview-landscape.webp   # Reference image: Landscape
│   │   └── preview-product.webp     # Reference image: Product/Still-life
│   ├── wasm_exec.js         # Go WASM runtime (existing)
│   └── recipe.wasm          # Compiled conversion engine (existing)
├── docs/
│   ├── PRD-Path-A.md            # This enhancement's requirements
│   ├── architecture-path-a.md   # This document
│   ├── parameter-mapping.md     # Updated with Capture One/DCP mappings
│   └── ...                      # Other existing docs
├── go.mod
├── go.sum
└── README.md
```

## Epic to Architecture Mapping

### Epic 1: Capture One Format Support

**Architecture Components:**
- `internal/formats/costyle/` package
  - `parse.go`: XML parsing using `encoding/xml`
  - `generate.go`: XML generation with formatted output
  - `pack.go`: ZIP archive handling using `archive/zip`
  - `types.go`: Go structs matching Capture One XML schema

**Integration Points:**
- Extends `converter.Convert()` to recognize .costyle format
- CLI: `recipe convert input.costyle --to xmp`
- TUI: Format menu includes "Capture One" option
- Web: Format detection via file extension + magic bytes
- Parameter mapping: New entries in `docs/parameter-mapping.md`

**Key Design Decisions:**
- Follow exact pattern of existing formats (np3, xmp, lrtemplate)
- Separate ZIP handling into dedicated `pack.go` file
- Test samples acquired from Etsy/marketplace .costyle presets

### Epic 2: DCP Camera Profile Support

**Architecture Components:**
- `internal/formats/dcp/` package
  - `parse.go`: Orchestrates TIFF reading → XML extraction
  - `generate.go`: Orchestrates XML generation → TIFF embedding
  - `tiff.go`: Low-level TIFF tag operations using `github.com/google/tiff`
  - `profile.go`: Adobe camera profile XML parsing/generation
  - `types.go`: Color matrices, tone curves, HSV tables

**Integration Points:**
- Extends `converter.Convert()` to recognize .dcp format
- TIFF tags: Read/write CameraProfile tag (50740)
- XML namespace: Adobe DNG Camera Profile namespace
- Parameter mapping: Color adjustments → DCP tone curves

**Key Design Decisions:**
- Split TIFF operations from XML profile logic
- Use identity matrices for non-calibration use cases
- Generate tone curves from exposure/contrast/highlights/shadows
- Document mapping limitations (dual illuminant not supported)

### Epic 3: Enhanced Web UI/UX

**Architecture Components:**
- `web/index.html`: Redesigned landing page
- `web/css/`:
  - `main.css`: CSS variables for color palette (format badge colors), typography, spacing
  - `components.css`: Reusable badge, button, card components
  - `layout.css`: Responsive grid (mobile <768px, tablet 768-1024px, desktop >1024px)
  - `preview.css`: Modal overlay, before/after slider styles
- `web/js/`:
  - `upload.js`: Drag-drop zone, batch file handling, progress tracking
  - `utils.js`: Format detection, file validation, error messaging

**Integration Points:**
- WASM interface via existing `converter.js`
- Batch operations: Queue files, convert sequentially or in parallel
- Format badges: Color-coded by format type (NP3=#FFC107, XMP=#0073E6, lrtemplate=#D81B60, Capture One=#9C27B0, DCP=#4CAF50)
- Responsive breakpoints:
  - Mobile (<768px): Single column, tap-to-browse
  - Tablet (768-1024px): Two-column grid
  - Desktop (>1024px): Three-column grid

**Key Design Decisions:**
- Vanilla JavaScript (no React/Vue/framework dependencies)
- Modern CSS (Grid, Flexbox, CSS Variables)
- Progressive enhancement (works without JS for basic conversion)
- BEM naming convention for CSS classes

### Epic 4: Image Preview System (Phase 1)

**Architecture Components:**
- `web/js/preview.js`: CSS filter mapping engine
  - `mapToCSSFilters(universalRecipe)`: Parameter → CSS filter string conversion
  - Modal management (open, close, tab switching)
  - Slider drag interaction handling
- `web/css/preview.css`: Modal styles, slider component
- `web/images/`: Three WebP reference images (portrait, landscape, product)

**Integration Points:**
- Reads UniversalRecipe parameters after parsing, before conversion
- Applies CSS filters to reference images in real-time
- "Convert now" button proceeds to actual conversion
- Clear labeling: "Approximate preview using CSS filters"

**CSS Filter Mapping Formula:**
```javascript
// Documented mapping formulas:
Exposure (-2.0 to +2.0) → brightness(0% to 200%)
  Formula: brightness = 100 + (exposure * 50)

Contrast (-100 to +100) → contrast(0% to 200%)
  Formula: contrast = 100 + contrast_value

Saturation (-100 to +100) → saturate(0% to 200%)
  Formula: saturation = 100 + saturation_value

Hue (degrees) → hue-rotate(Xdeg)
  Direct mapping

Temperature (warm tones) → sepia(%) approximation
  Formula: sepia = min(temperature / 100 * 30, 30)
```

**Limitations (documented in UI):**
- Tone curves not supported in Phase 1 (CSS filters can't replicate)
- HSL adjustments approximated
- Split toning not supported
- Vignette/grain not previewed

**Key Design Decisions:**
- Phase 1: CSS filters (instant, <100ms)
- Phase 2 (future): WebAssembly-based pixel-perfect preview
- Reference images: WebP format, 1200×800px, <150KB each
- Performance target: <100ms preview rendering, 60fps slider interaction

## Technology Stack Details

### Core Technologies

**Backend (Conversion Engine):**
- **Language:** Go 1.25.1+ with WASM compilation target (`GOOS=js GOARCH=wasm`)
- **Standard Libraries:**
  - `encoding/xml`: XML parsing for Capture One and DCP
  - `archive/zip`: .costylepack bundle handling
  - `fmt`: Error formatting with wrapping
  - `time`: Timestamp generation (RFC3339 format)
- **External Dependencies:**
  - `github.com/google/tiff@latest`: TIFF reading/writing for DCP format
  - (Minimal external dependencies maintained)

**Frontend (Web Interface):**
- **JavaScript:** ES6+ (no transpilation), vanilla JS (no framework)
- **CSS:** Modern CSS (Grid, Flexbox, CSS Variables, CSS Filters)
- **Image Format:** WebP for reference images
- **Browser Targets:** Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- **WASM Runtime:** Go's `wasm_exec.js` (bundled with Go distribution)

**Development Tools:**
- **Testing:** Go's built-in `testing` package, table-driven tests
- **Linting:** `golangci-lint` for Go code
- **Formatting:** `gofmt` for Go, consistent style for JS/CSS
- **Coverage:** `go test -cover`, target ≥85% per package

### Integration Points

**Format ↔ Converter:**
```go
// Every format implements this interface pattern:
type Format interface {
    Parse(data []byte) (UniversalRecipe, error)
    Generate(recipe UniversalRecipe) ([]byte, error)
}

// Converter orchestrates:
func Convert(input []byte, fromFormat, toFormat string) ([]byte, error) {
    // 1. Detect format if not specified
    // 2. Parse: inputFormat.Parse(input) → UniversalRecipe
    // 3. Generate: outputFormat.Generate(recipe) → output
    // 4. Return output bytes
}
```

**Converter ↔ Interfaces:**
- **CLI:** Direct function calls to `converter.Convert()`
- **TUI:** Bubbletea models wrap converter with progress updates
- **Web:** WASM exposes `convert()` function via `//go:wasmexport`

**Web UI ↔ WASM:**
```javascript
// converter.js wraps WASM calls
async function convertFile(inputBytes, fromFormat, toFormat) {
    // Call WASM convert() function
    // Returns output bytes or error
}

// upload.js manages batch operations
async function convertBatch(files, outputFormat) {
    for (const file of files) {
        const result = await convertFile(file.bytes, file.format, outputFormat);
        // Update progress, handle errors
    }
}

// preview.js reads UniversalRecipe before conversion
async function showPreview(inputBytes, format) {
    // Parse to UniversalRecipe (WASM call)
    // Map to CSS filters (JavaScript)
    // Apply to reference images (DOM manipulation)
}
```

## Implementation Patterns

These patterns ensure consistent implementation across all AI agents:

### Naming Conventions

**Go Code (Backend):**
- **Package names:** Lowercase, no underscores
  - ✅ `costyle`, `dcp`
  - ❌ `co_style`, `DCP`
- **File names:** Lowercase, underscores allowed
  - ✅ `parse.go`, `pack.go`, `parse_test.go`
  - ❌ `Parse.go`, `parseTest.go`
- **Type names:** PascalCase (exported), camelCase (internal)
  - ✅ `type CaptureOneStyle struct { ... }`
  - ✅ `type xmlElement struct { ... }` (internal)
- **Function names:** PascalCase (exported), camelCase (internal)
  - ✅ `func Parse(data []byte) error`
  - ✅ `func parseCostyle(data []byte) error` (internal helper)
- **Variable names:** camelCase
  - ✅ `universalRecipe`, `xmlData`, `tiffTags`

**Web Code (Frontend):**
- **File names:** kebab-case
  - ✅ `preview.css`, `upload.js`, `preview-portrait.webp`
  - ❌ `Preview.css`, `uploadJS.js`
- **CSS classes:** kebab-case with BEM
  - ✅ `.preview-modal`, `.preview-modal__slider`, `.preview-modal--active`
  - ❌ `.previewModal`, `.preview_modal`
- **JavaScript functions:** camelCase
  - ✅ `mapToCSSFilters()`, `handleFileUpload()`
- **JavaScript constants:** SCREAMING_SNAKE_CASE
  - ✅ `const MAX_FILE_SIZE = 10 * 1024 * 1024;`

### Code Organization

**Go Package Structure:**
Every format package MUST follow this structure:
```
internal/formats/{format}/
├── parse.go      # Parsing input format → UniversalRecipe
├── generate.go   # Generating output format ← UniversalRecipe
├── types.go      # Format-specific type definitions
├── {helper}.go   # Optional helpers (pack.go, tiff.go, profile.go)
└── testdata/     # Test sample files
    ├── sample1.{ext}
    ├── sample2.{ext}
    └── README.md # Attribution, sources, licenses
```

**Web Asset Organization:**
```
web/
├── index.html
├── css/          # Organized by purpose
│   ├── main.css      # First: Global styles, variables
│   ├── components.css # Second: Reusable components
│   ├── layout.css    # Third: Grid, responsive
│   └── {feature}.css # Fourth: Feature-specific (preview.css)
├── js/           # Organized by feature
│   ├── app.js        # Main entry point
│   ├── converter.js  # WASM interface
│   └── {feature}.js  # Feature modules (upload.js, preview.js)
└── images/       # All image assets
```

### Error Handling

**Pattern:** Always wrap errors with context using `fmt.Errorf()` and `%w`:
```go
// ✅ Correct
if err := parseXML(data); err != nil {
    return nil, fmt.Errorf("failed to parse .costyle XML: %w", err)
}

// ❌ Wrong
if err := parseXML(data); err != nil {
    return nil, err  // No context
}

// ❌ Wrong
if err := parseXML(data); err != nil {
    return nil, fmt.Errorf("failed to parse .costyle XML: %v", err)  // %v doesn't wrap
}
```

**User-facing errors:**
- Be specific and actionable
- ✅ "File format not recognized. Supported formats: .np3, .xmp, .lrtemplate, .costyle, .dcp"
- ❌ "Error 0x04F2"

**No panics in production code:**
- Handle errors gracefully, return error values
- Tests can use `t.Fatal()` but production code never panics

### Logging Strategy

**Minimal logging approach:**
- **CLI:** Progress to stdout, errors to stderr
- **TUI:** Bubbletea's built-in messaging system
- **Web:** Browser console for debugging only (no production logging)
- **No structured logging framework** (not a server application)

### Testing Patterns

**Test file naming:** Co-located with source
```
parse.go       → parse_test.go
generate.go    → generate_test.go
```

**Table-driven tests for conversions:**
```go
func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string  // Path to testdata file
        want    UniversalRecipe
        wantErr bool
    }{
        {name: "basic exposure", input: "testdata/exposure.costyle", want: ...},
        {name: "full adjustments", input: "testdata/complete.costyle", want: ...},
        {name: "invalid XML", input: "testdata/invalid.costyle", wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            data, _ := os.ReadFile(tt.input)
            got, err := Parse(data)
            // Assertions...
        })
    }
}
```

**Round-trip tests:**
```go
func TestRoundTrip(t *testing.T) {
    original, _ := os.ReadFile("testdata/sample.costyle")
    recipe, _ := Parse(original)
    generated, _ := Generate(recipe)
    recipe2, _ := Parse(generated)

    // Compare recipe and recipe2, allow <5% tolerance for float precision
}
```

**Coverage target:** ≥85% per package

## Consistency Rules

### Format Patterns

**XML Output:** Always formatted (human-readable), never minified
```go
// ✅ Correct: Formatted XML
encoder := xml.NewEncoder(buffer)
encoder.Indent("", "  ")

// ❌ Wrong: Minified XML
encoder := xml.NewEncoder(buffer) // No indentation
```

**TIFF Output:** Use `github.com/google/tiff` for all TIFF operations
- Do not implement custom TIFF tag writing
- Use library's tag structures and encoding

**ZIP Output:** Use `archive/zip` for all archive operations
- Create ZIP with standard compression (deflate)
- Maintain original file structure in .costylepack bundles

### Date/Time Format

**Timestamps in generated files:** RFC3339 format
```go
timestamp := time.Now().Format(time.RFC3339)
// Example: "2025-01-08T15:04:05Z"
```

### CSS Architecture

**Color Palette (CSS Variables in main.css):**
```css
:root {
    /* Format badge colors */
    --color-np3: #FFC107;          /* Nikon yellow */
    --color-xmp: #0073E6;          /* Adobe blue */
    --color-lrtemplate: #D81B60;   /* Magenta */
    --color-costyle: #9C27B0;      /* Purple */
    --color-dcp: #4CAF50;          /* Green */

    /* Neutral palette */
    --color-bg: #FFFFFF;
    --color-text: #212121;
    --color-border: #E0E0E0;
}
```

**Responsive Breakpoints:**
```css
/* Mobile first, then override */
.component { /* Mobile <768px styles */ }

@media (min-width: 768px) {
    .component { /* Tablet 768-1024px */ }
}

@media (min-width: 1024px) {
    .component { /* Desktop >1024px */ }
}
```

**BEM Naming:**
```css
.block { }               /* Component */
.block__element { }      /* Part of component */
.block--modifier { }     /* Variant of component */

/* Example: */
.preview-modal { }
.preview-modal__slider { }
.preview-modal__image { }
.preview-modal--active { }
```

## Data Architecture

### UniversalRecipe (Existing)

Path A extends the existing UniversalRecipe intermediate representation with parameters needed for Capture One and DCP:

**No changes to UniversalRecipe structure** - existing parameters are sufficient:
- Exposure, Contrast, Saturation, Hue (already present)
- Temperature, Tint (already present)
- Highlights, Shadows, Whites, Blacks (already present)
- Color balance (RGB curves, already present)

**Mapping considerations:**
- Capture One .costyle maps cleanly to existing UniversalRecipe params
- DCP tone curves mapped from exposure/contrast/highlights/shadows
- DCP color matrices use identity matrices (not full calibration)
- Unsupported DCP features (dual illuminant) documented in parameter-mapping.md

### Format-Specific Types

**Capture One (internal/formats/costyle/types.go):**
```go
type CaptureOneStyle struct {
    XMLName     xml.Name `xml:"style"`
    Version     string   `xml:"version,attr"`
    Exposure    float64  `xml:"exposure"`
    Contrast    int      `xml:"contrast"`
    Saturation  int      `xml:"saturation"`
    Temperature int      `xml:"temperature"`
    Tint        int      `xml:"tint"`
    // ... more fields
}
```

**DCP (internal/formats/dcp/types.go):**
```go
type DCPProfile struct {
    ColorMatrix1    [9]float64  // 3x3 color matrix
    ToneCurve       []TonePoint // Tone curve points
    HueSatMap       []HSVAdjust // HSV adjustments
    // ... more fields
}

type TonePoint struct {
    Input  float64
    Output float64
}
```

## Security Architecture

**Privacy-First Processing:**
- All conversion happens **client-side** (browser WASM or local CLI/TUI)
- **Zero data transmission** to external servers
- No analytics, tracking, or telemetry
- No external resource loading (fonts, scripts, images self-hosted)

**File Handling Security:**
- **File type validation:** Magic number checks, not just extension
  - `.costyle`: Verify XML header `<?xml`
  - `.dcp`: Verify TIFF magic bytes `49 49 2A 00` (little-endian) or `4D 4D 00 2A` (big-endian)
- **XML External Entity (XXE) prevention:** Disable external entity resolution in `encoding/xml`
- **ZIP bomb protection:** Limit uncompressed size ratio for .costylepack (max 100:1)
- **File size limits:** 10MB max per file (prevent DoS)
- **Graceful error handling:** No crashes, no sensitive error details exposed

**Web Application Security:**
- **Content Security Policy (CSP):** Strict headers via Cloudflare Pages
  - `default-src 'self'`
  - `script-src 'self' 'wasm-unsafe-eval'` (WASM requirement)
  - No `eval()` or `Function()` constructors in JavaScript
- **HTTPS-only:** Enforced by Cloudflare Pages deployment
- **No user-generated content:** No uploads to server (all client-side)

## Performance Considerations

**Conversion Performance:**
- **Target:** <100ms per conversion (existing Recipe performance)
- **WASM compilation:** Optimized with `GOOS=js GOARCH=wasm`
- **Concurrent batch processing:** Web Workers for multi-file conversions (future optimization)
- **Memory efficiency:** Stream large .costylepack bundles, don't load all in memory

**Web Interface Performance:**
- **Load time target:** <2 seconds on 3G connection
- **WASM caching:** Aggressive cache headers (1 year TTL)
- **Critical CSS:** Inline minimal CSS for first paint
- **Image optimization:** WebP format, <150KB each, lazy-load if needed
- **Preview rendering:** <100ms CSS filter application, 60fps slider interaction

**Lighthouse Performance Targets:**
- Performance score: ≥90
- Accessibility score: ≥95 (WCAG 2.1 AA)
- Best Practices score: ≥90
- SEO score: ≥90

## Deployment Architecture

**Web Interface:**
- **Platform:** Cloudflare Pages (existing)
- **Build command:** `make web` (builds WASM + copies assets)
- **Output directory:** `web/`
- **Deployment:** Automatic on `git push` to main branch
- **Rollback:** Cloudflare Pages deployment history (instant rollback)

**CLI/TUI Binaries:**
- **Platform:** GitHub Releases (existing)
- **Build:** GitHub Actions matrix build
  - `GOOS=windows GOARCH=amd64`
  - `GOOS=darwin GOARCH=amd64`
  - `GOOS=darwin GOARCH=arm64`
  - `GOOS=linux GOARCH=amd64`
  - `GOOS=linux GOARCH=arm64`
- **Versioning:** Semantic versioning (SemVer 2.0.0)
- **Release notes:** CHANGELOG.md with Path A enhancements

## Development Environment

### Prerequisites

**Required:**
- Go 1.25.1 or higher (`go version`)
- Git for version control
- Text editor with Go support (VS Code, GoLand, Vim, etc.)

**Optional:**
- `golangci-lint` for linting
- `make` for build automation (Makefile exists)
- Modern browser for web testing (Chrome, Firefox, Safari, or Edge)

### Setup Commands

**Clone and build:**
```bash
# Clone repository
git clone https://github.com/user/recipe.git
cd recipe

# Install dependencies (minimal - just google/tiff)
go mod download

# Build CLI binary
go build -o recipe ./cmd/recipe

# Build TUI binary
go build -o recipe-tui ./cmd/recipe-tui

# Build WASM for web
GOOS=js GOARCH=wasm go build -o web/recipe.wasm ./cmd/wasm

# Run tests
go test ./...

# Check coverage
go test -cover ./internal/...
```

**Development workflow:**
```bash
# Run CLI locally
./recipe convert input.np3 --to xmp

# Run TUI locally
./recipe-tui

# Serve web interface locally
# (Simple HTTP server, Python 3)
cd web
python3 -m http.server 8080
# Open http://localhost:8080 in browser

# Run tests with coverage
go test -coverprofile=coverage.out ./internal/formats/costyle/
go tool cover -html=coverage.out
```

**Adding test samples:**
```bash
# Acquire .costyle samples from Etsy/marketplaces
# Place in internal/formats/costyle/testdata/
# Document source in README.md

# Acquire DCP samples from Adobe
# Place in internal/formats/dcp/testdata/
# Document license in README.md
```

## Architecture Decision Records (ADRs)

### ADR-001: Maintain Zero External Dependencies

**Status:** Accepted

**Context:**
Recipe's core principle is privacy-first, local-only processing with minimal dependencies. Path A adds new formats requiring XML parsing, TIFF handling, and ZIP operations.

**Decision:**
Use Go standard library (`encoding/xml`, `archive/zip`) for all operations except TIFF, where `github.com/google/tiff` is justified as Google maintains both Go and this library.

**Consequences:**
- ✅ Privacy maintained (no external API calls)
- ✅ Minimal attack surface
- ✅ Fast compilation, smaller binaries
- ⚠️ Slight complexity for TIFF operations (but library handles it)

### ADR-002: CSS Filters for Preview (Phase 1)

**Status:** Accepted

**Context:**
Users want to preview preset effects before conversion. Options: CSS filters (instant but approximate) or WASM image processing (accurate but slower, larger bundle).

**Decision:**
Implement CSS filter-based preview for Phase 1, defer WASM-based accurate preview to Phase 2 (future).

**Consequences:**
- ✅ Instant preview (<100ms rendering)
- ✅ No additional bundle size
- ✅ Works across all browsers
- ⚠️ Approximate accuracy (tone curves not supported)
- ✅ Clear labeling manages user expectations

### ADR-003: Vanilla JavaScript (No Framework)

**Status:** Accepted

**Context:**
Path A enhances web UI with batch upload, preview modal, and responsive design. Could use React/Vue/Svelte or vanilla JS.

**Decision:**
Continue with vanilla JavaScript + modern CSS (no framework).

**Consequences:**
- ✅ No framework bloat (faster load times)
- ✅ Maintains Recipe's simplicity principle
- ✅ Easier WASM integration (no framework bridge layer)
- ⚠️ Slightly more verbose DOM manipulation code
- ✅ Future-proof (no framework version churn)

### ADR-004: Follow Existing Format Package Pattern

**Status:** Accepted

**Context:**
Recipe has established pattern for format packages (parse.go, generate.go, types.go, testdata/). New formats could diverge.

**Decision:**
Strictly follow existing pattern for Capture One and DCP formats, adding helper files (pack.go, tiff.go, profile.go) as needed but maintaining flat structure.

**Consequences:**
- ✅ Architectural consistency
- ✅ AI agents know exactly where code lives
- ✅ Easier navigation and maintenance
- ✅ Clear separation of concerns (parse vs generate vs helpers)

### ADR-005: Embedded Reference Images (Not Base64)

**Status:** Accepted

**Context:**
Preview feature needs reference images. Options: embed as Base64 in JS/CSS, serve as separate files, or external CDN.

**Decision:**
Serve optimized WebP images from `web/images/` directory, rely on browser caching.

**Consequences:**
- ✅ Clean separation of code and assets
- ✅ Leverages Cloudflare CDN + browser caching
- ✅ Easy to update/replace images
- ⚠️ Three HTTP requests (mitigated by caching)
- ✅ Smaller JS bundle size

---

_Generated by BMAD Decision Architecture Workflow v1.3.2_
_Date: 2025-01-08_
_For: Justin_
_Project: recipe Path A Enhancements_
