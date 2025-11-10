# Epic Technical Specification: Capture One Format Support

Date: 2025-11-08
Author: Justin
Epic ID: epic-8
Status: Draft

---

## Overview

Epic 8 adds Capture One .costyle format support to Recipe, enabling photographers to convert Nikon Picture Control recipes and other presets to/from Capture One's XML-based style format. This extends Recipe's format coverage to serve the 8% of professional photographers who use Capture One Pro as their primary editing platform.

The implementation follows Recipe's established hub-and-spoke architecture, with all conversions flowing through the UniversalRecipe intermediate representation. Capture One styles are XML-based (similar to Adobe XMP), making them straightforward to parse and generate using Go's standard library. Support includes both single .costyle files and .costylepack bundles (ZIP archives containing multiple styles).

## Objectives and Scope

**In Scope:**
- Parse Capture One .costyle XML files (exposure, contrast, saturation, temperature, tint, clarity, color balance)
- Generate valid .costyle XML files from UniversalRecipe representation
- Support .costylepack bundles (ZIP archives of multiple styles)
- Round-trip conversion testing with 95%+ parameter preservation
- Integration across all Recipe interfaces (CLI, TUI, Web)
- Parameter mapping documentation for Capture One-specific adjustments

**Out of Scope (Path A):**
- Capture One-specific features not representable in UniversalRecipe (local adjustments, layers)
- Advanced color grading beyond basic color balance (LUT integration)
- Capture One metadata beyond core adjustment parameters
- Capture One ICC profile generation
- Capture One catalog integration
- Support for Capture One versions prior to 2023

## System Architecture Alignment

**Components:**
- **New Package**: `internal/formats/costyle/` (following exact pattern of np3, xmp, lrtemplate)
  - `parse.go`: XML parsing using Go stdlib `encoding/xml`
  - `generate.go`: XML generation with formatted output
  - `types.go`: Go structs matching Capture One XML schema
  - `pack.go`: ZIP archive handling using `archive/zip`
  - `testdata/`: Real .costyle samples from Etsy/marketplaces

**Integration Points:**
- Extends `internal/converter/converter.go` to recognize .costyle format
- Adds Capture One format badge (purple #9C27B0) to web UI
- Updates `docs/parameter-mapping.md` with Capture One mappings
- Leverages existing format detection logic (extension + magic bytes)

**Constraints:**
- Must maintain <100ms conversion performance
- Must preserve hub-and-spoke architecture (no direct format-to-format conversion)
- Must use only Go stdlib + github.com/google/tiff (for TIFF in DCP epic only)
- Must maintain ≥85% test coverage
- Zero external dependencies for XML parsing (use encoding/xml)

## Detailed Design

### Services and Modules

| Module | Responsibility | Inputs | Outputs | Owner |
| ------ | -------------- | ------ | ------- | ----- |
| `costyle/parse.go` | Parse .costyle XML → UniversalRecipe | .costyle file bytes | `*universal.Recipe`, error | Dev (Epic 8) |
| `costyle/generate.go` | Generate .costyle XML ← UniversalRecipe | `*universal.Recipe` | .costyle file bytes, error | Dev (Epic 8) |
| `costyle/pack.go` | ZIP pack/unpack for .costylepack bundles | ZIP bytes / []Recipe | []Recipe / ZIP bytes, error | Dev (Epic 8) |
| `costyle/types.go` | Capture One XML schema types | - | Go struct definitions | Dev (Epic 8) |
| `converter/converter.go` | Format detection & routing (EXISTING) | File bytes, target format | Converted bytes, error | Dev (Epic 8 extension) |

### Data Models and Contracts

**Capture One .costyle XML Structure:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description>
      <!-- Core Adjustments -->
      <crs:Exposure>0.0</crs:Exposure>          <!-- -2.0 to +2.0 -->
      <crs:Contrast>0</crs:Contrast>            <!-- -100 to +100 -->
      <crs:Saturation>0</crs:Saturation>        <!-- -100 to +100 -->
      <crs:Temperature>0</crs:Temperature>      <!-- -100 to +100 -->
      <crs:Tint>0</crs:Tint>                    <!-- -100 to +100 -->
      <crs:Clarity>0</crs:Clarity>              <!-- -100 to +100 -->

      <!-- Color Balance -->
      <crs:ShadowsHue>0</crs:ShadowsHue>
      <crs:ShadowsSaturation>0</crs:ShadowsSaturation>
      <crs:MidtonesHue>0</crs:MidtonesHue>
      <crs:MidtonesSaturation>0</crs:MidtonesSaturation>
      <crs:HighlightsHue>0</crs:HighlightsHue>
      <crs:HighlightsSaturation>0</crs:HighlightsSaturation>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>
```

**Go Struct Mapping:**

```go
// types.go
package costyle

type CaptureOneStyle struct {
    XMLName xml.Name `xml:"xmpmeta"`
    RDF     RDF      `xml:"RDF"`
}

type RDF struct {
    Description Description `xml:"Description"`
}

type Description struct {
    // Core adjustments
    Exposure    float64 `xml:"Exposure,omitempty"`    // -2.0 to +2.0
    Contrast    int     `xml:"Contrast,omitempty"`    // -100 to +100
    Saturation  int     `xml:"Saturation,omitempty"`  // -100 to +100
    Temperature int     `xml:"Temperature,omitempty"` // -100 to +100
    Tint        int     `xml:"Tint,omitempty"`        // -100 to +100
    Clarity     int     `xml:"Clarity,omitempty"`     // -100 to +100

    // Color balance (shadows/midtones/highlights)
    ShadowsHue              int `xml:"ShadowsHue,omitempty"`
    ShadowsSaturation       int `xml:"ShadowsSaturation,omitempty"`
    MidtonesHue             int `xml:"MidtonesHue,omitempty"`
    MidtonesSaturation      int `xml:"MidtonesSaturation,omitempty"`
    HighlightsHue           int `xml:"HighlightsHue,omitempty"`
    HighlightsSaturation    int `xml:"HighlightsSaturation,omitempty"`
}
```

**UniversalRecipe Mapping:**

Capture One parameters map to UniversalRecipe fields:
- `Exposure` → `universal.Exposure` (direct mapping)
- `Contrast` → `universal.Contrast` (scale: C1 -100..100 = UR -1.0..1.0)
- `Saturation` → `universal.Saturation` (scale: C1 -100..100 = UR -1.0..1.0)
- `Temperature` → `universal.Temperature` (scale: C1 -100..100 = UR -100..100)
- `Tint` → `universal.Tint` (scale: C1 -100..100 = UR -100..100)
- `Clarity` → `universal.Clarity` (direct mapping)
- Color balance → `universal.ColorBalance{Shadows, Midtones, Highlights}` with hue/saturation per range

### APIs and Interfaces

**costyle.Parse Function:**

```go
// parse.go
package costyle

// Parse parses a Capture One .costyle file into UniversalRecipe
// Returns error if XML is invalid or structure doesn't match schema
func Parse(data []byte) (*universal.Recipe, error) {
    var style CaptureOneStyle
    if err := xml.Unmarshal(data, &style); err != nil {
        return nil, fmt.Errorf("failed to parse .costyle XML: %w", err)
    }

    return styleToUniversal(&style), nil
}

// styleToUniversal converts Capture One types to UniversalRecipe
func styleToUniversal(style *CaptureOneStyle) *universal.Recipe {
    // Implementation: map fields with range scaling
}
```

**costyle.Generate Function:**

```go
// generate.go
package costyle

// Generate creates a Capture One .costyle file from UniversalRecipe
// Returns formatted XML bytes suitable for Capture One import
func Generate(recipe *universal.Recipe) ([]byte, error) {
    style := universalToStyle(recipe)

    data, err := xml.MarshalIndent(style, "", "  ")
    if err != nil {
        return nil, fmt.Errorf("failed to generate .costyle XML: %w", err)
    }

    // Add XML declaration
    return append([]byte(xml.Header), data...), nil
}

// universalToStyle converts UniversalRecipe to Capture One types
func universalToStyle(recipe *universal.Recipe) *CaptureOneStyle {
    // Implementation: map fields with range scaling
}
```

**costyle.Pack Function (for .costylepack bundles):**

```go
// pack.go
package costyle

// ParsePack extracts multiple styles from a .costylepack ZIP archive
func ParsePack(data []byte) ([]*universal.Recipe, error) {
    r := bytes.NewReader(data)
    zipReader, err := zip.NewReader(r, int64(len(data)))
    if err != nil {
        return nil, fmt.Errorf("failed to read .costylepack ZIP: %w", err)
    }

    var recipes []*universal.Recipe
    for _, file := range zipReader.File {
        if filepath.Ext(file.Name) != ".costyle" {
            continue
        }

        rc, err := file.Open()
        if err != nil {
            return nil, fmt.Errorf("failed to open %s: %w", file.Name, err)
        }
        defer rc.Close()

        data, err := io.ReadAll(rc)
        if err != nil {
            return nil, fmt.Errorf("failed to read %s: %w", file.Name, err)
        }

        recipe, err := Parse(data)
        if err != nil {
            return nil, fmt.Errorf("failed to parse %s: %w", file.Name, err)
        }

        recipes = append(recipes, recipe)
    }

    return recipes, nil
}

// GeneratePack creates a .costylepack ZIP archive from multiple recipes
func GeneratePack(recipes []*universal.Recipe, names []string) ([]byte, error) {
    buf := new(bytes.Buffer)
    zipWriter := zip.NewWriter(buf)

    for i, recipe := range recipes {
        name := names[i]
        if !strings.HasSuffix(name, ".costyle") {
            name += ".costyle"
        }

        w, err := zipWriter.Create(name)
        if err != nil {
            return nil, fmt.Errorf("failed to create %s in archive: %w", name, err)
        }

        data, err := Generate(recipe)
        if err != nil {
            return nil, fmt.Errorf("failed to generate %s: %w", name, err)
        }

        if _, err := w.Write(data); err != nil {
            return nil, fmt.Errorf("failed to write %s: %w", name, err)
        }
    }

    if err := zipWriter.Close(); err != nil {
        return nil, fmt.Errorf("failed to close ZIP archive: %w", err)
    }

    return buf.Bytes(), nil
}
```

**converter.Convert Integration:**

```go
// internal/converter/converter.go (EXISTING - extend to add costyle)

func (c *Converter) Convert(data []byte, targetFormat string) ([]byte, error) {
    // Detect source format
    sourceFormat, err := detectFormat(data)
    if err != nil {
        return nil, fmt.Errorf("failed to detect format: %w", err)
    }

    // Parse source → UniversalRecipe
    var recipe *universal.Recipe
    switch sourceFormat {
    case "np3":
        recipe, err = np3.Parse(data)
    case "xmp":
        recipe, err = xmp.Parse(data)
    case "lrtemplate":
        recipe, err = lrtemplate.Parse(data)
    case "costyle":  // NEW
        recipe, err = costyle.Parse(data)
    case "dcp":  // Epic 9
        recipe, err = dcp.Parse(data)
    default:
        return nil, fmt.Errorf("unsupported source format: %s", sourceFormat)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to parse %s: %w", sourceFormat, err)
    }

    // Generate UniversalRecipe → target format
    var output []byte
    switch targetFormat {
    case "np3":
        output, err = np3.Generate(recipe)
    case "xmp":
        output, err = xmp.Generate(recipe)
    case "lrtemplate":
        output, err = lrtemplate.Generate(recipe)
    case "costyle":  // NEW
        output, err = costyle.Generate(recipe)
    case "dcp":  // Epic 9
        output, err = dcp.Generate(recipe)
    default:
        return nil, fmt.Errorf("unsupported target format: %s", targetFormat)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to generate %s: %w", targetFormat, err)
    }

    return output, nil
}

func detectFormat(data []byte) (string, error) {
    // Extension-based detection (primary)
    // Magic byte detection (fallback)
    // Add .costyle detection: check for XML header + xmpmeta namespace
    if bytes.Contains(data, []byte("<?xml")) &&
       bytes.Contains(data, []byte("x:xmpmeta")) &&
       bytes.Contains(data, []byte("crs:Exposure")) {
        return "costyle", nil
    }

    // ... existing format detection logic
}
```

### Workflows and Sequencing

**Parse Workflow (.costyle → UniversalRecipe):**

1. User uploads .costyle file via CLI/TUI/Web
2. `converter.detectFormat(data)` identifies file as Capture One style
3. `costyle.Parse(data)` called:
   a. `xml.Unmarshal()` parses XML structure into `CaptureOneStyle` struct
   b. Validate required fields exist (Exposure, Contrast, etc.)
   c. `styleToUniversal()` maps Capture One parameters to UniversalRecipe fields
   d. Scale ranges (Capture One -100..100 → UniversalRecipe -1.0..1.0 for contrast/saturation)
4. Return `*universal.Recipe` to converter
5. Converter routes to target format generator

**Generate Workflow (UniversalRecipe → .costyle):**

1. Converter calls `costyle.Generate(recipe)`
2. `universalToStyle()` maps UniversalRecipe → CaptureOneStyle:
   a. Direct copy for exposure, clarity, temperature, tint
   b. Scale ranges for contrast/saturation (UR -1.0..1.0 → C1 -100..100)
   c. Map color balance to shadows/midtones/highlights hue/saturation
3. `xml.MarshalIndent()` generates formatted XML
4. Prepend XML declaration: `<?xml version="1.0" encoding="UTF-8"?>`
5. Return XML bytes to converter
6. Converter returns to user (download/write file)

**Pack Workflow (.costylepack handling):**

1. User uploads .costylepack file (ZIP with multiple .costyle files)
2. Format detection identifies ZIP signature + .costyle extensions
3. `costyle.ParsePack(data)`:
   a. `zip.NewReader()` opens ZIP archive
   b. Iterate through ZIP entries, filter .costyle files
   c. Parse each .costyle file via `costyle.Parse()`
   d. Return slice of UniversalRecipes
4. User selects conversion target format
5. Convert each recipe individually
6. Optionally bundle outputs into target format pack (e.g., .costylepack → .zip of XMPs)

## Non-Functional Requirements

### Performance

**Targets (per PRD-Path-A.md, architecture-path-a.md):**
- Single .costyle conversion: <100ms (99th percentile)
- Batch .costylepack (50 files): <5 seconds total
- Memory usage: <10MB per conversion (can process large bundles)
- Web WASM build: Conversion time equivalent to native (WASM overhead <10%)

**Optimization Strategies:**
- Use streaming XML parsing for large files (encoding/xml is SAX-based, efficient)
- Preallocate slices for .costylepack bundles (avoid reallocations)
- Leverage Go's fast XML marshaling (no reflection overhead)
- Benchmark against existing formats (np3, xmp) to ensure parity

**Performance Monitoring:**
- Add benchmark tests: `BenchmarkParseCostyle`, `BenchmarkGenerateCostyle`
- CI/CD integration: Fail if regression >10% vs previous commit
- Document baseline: Expected <50ms for typical .costyle file

### Security

**Threats & Mitigations:**
- **XML Bomb (Billion Laughs)**: encoding/xml resistant by default (no entity expansion)
- **ZIP Bomb**: Limit extracted file sizes, check ZIP compression ratios before extraction
- **Path Traversal in .costylepack**: Validate ZIP entry names, reject paths with `..` or absolute paths
- **Malformed XML**: Graceful error handling, no panics on invalid input
- **File Upload**: Web UI validates file extensions client-side (UX), server validates magic bytes

**Privacy:**
- All processing client-side (WASM) or local (CLI/TUI)
- Zero external API calls for conversion
- No analytics, no tracking, no server uploads
- .costyle files may contain metadata (preset names) - preserve but don't expose beyond necessary

**Input Validation:**
- Check XML structure matches expected schema (required fields)
- Validate parameter ranges before conversion (exposure -2.0..2.0, etc.)
- Reject files >10MB (prevent DoS via large uploads)
- Sanitize filenames in .costylepack bundles (prevent path traversal)

### Reliability/Availability

**Error Handling:**
- Wrapped errors with context: `fmt.Errorf("failed to parse .costyle: %w", err)`
- Graceful degradation: Missing optional parameters use defaults (0 values)
- Round-trip validation: Test that Parse(Generate(recipe)) == recipe (95%+ accuracy)
- Atomic ZIP operations: Either entire .costylepack succeeds or fails (no partial state)

**Backwards Compatibility:**
- Support Capture One .costyle versions from 2023-2025
- Document known version differences in testdata README
- If new Capture One version changes schema, add version detection + fallback parsing

**Failure Modes:**
- Invalid XML: Return clear error message, don't crash
- Unsupported parameter: Log warning, skip parameter (don't fail entire conversion)
- ZIP extraction error: Return error with specific file that failed

### Observability

**Logging:**
- Use Go's `log/slog` structured logging (established in Epic 3)
- Log level DEBUG: Parsing steps, parameter mappings
- Log level INFO: Conversion success/failure, file counts for .costylepack
- Log level WARN: Skipped unsupported parameters, version mismatches
- Log level ERROR: Parse failures, ZIP errors

**Metrics:**
- Count conversions by source/target format pair (e.g., "costyle → xmp")
- Track parse/generate latency (p50, p95, p99)
- Monitor round-trip accuracy (% parameter preservation)

**Tracing:**
- Not required for personal project, but structure code for future instrumentation
- Function boundaries clear: Parse, Generate, styleToUniversal, universalToStyle

**Debugging:**
- Verbose mode (`-v` flag in CLI) dumps intermediate UniversalRecipe as JSON
- Test suite includes `testdata/sample.costyle` → expected `universal.Recipe` JSON

## Dependencies and Integrations

**Go Module Dependencies:**

Current (from go.mod):
```
module github.com/justin/recipe
go 1.25.1

require github.com/spf13/cobra v1.10.1
require (
    // Bubbletea v2 for TUI (existing)
    // Standard library dependencies
)
```

New Dependencies for Epic 8:
- **NONE** - All Capture One format support uses Go standard library only

Standard Library Usage:
- `encoding/xml`: XML parsing and generation (.costyle files)
- `archive/zip`: ZIP pack/unpack for .costylepack bundles
- `fmt`: Error wrapping with context
- `io`, `bytes`: Stream handling for ZIP extraction

**External Integrations:**
- **Capture One Pro** (validation): Generate .costyle files must load in Capture One software (manual testing)
- **Etsy/Marketplace .costyle samples** (test data): Acquire real preset files for test fixtures

**Version Constraints:**
- Go 1.25.1+ (required for go:wasmexport in existing codebase)
- Capture One 2023+ (supported .costyle format versions)
- No platform-specific dependencies (cross-platform: Windows, macOS, Linux, WASM)

## Acceptance Criteria (Authoritative)

**AC-1: Parse Capture One .costyle Files**
- ✅ Parse XML structure per Capture One style specification
- ✅ Extract exposure, contrast, saturation, temperature, tint, clarity adjustments
- ✅ Extract color balance (shadows, midtones, highlights) with hue/saturation
- ✅ Extract tone curve points if present (optional, may not map to UniversalRecipe)
- ✅ Handle missing or optional parameters gracefully (use zero values)
- ✅ Validate XML structure and report parsing errors with clear messages
- ✅ Support .costyle format versions currently in use (2023-2025)

**AC-2: Generate Capture One .costyle Files**
- ✅ Generate valid XML structure matching Capture One specification
- ✅ Map UniversalRecipe parameters to .costyle equivalents with correct scaling
- ✅ Include required XML elements (version, metadata, namespaces)
- ✅ Generate human-readable XML (formatted with indentation, not minified)
- ✅ Validate generated XML against schema (well-formed, valid structure)
- ✅ Generated files load successfully in Capture One software (manual validation)
- ✅ Handle edge cases (missing parameters use defaults, out-of-range values clamped)

**AC-3: Support .costylepack Bundles**
- ✅ Unzip .costylepack archives and extract individual .costyle files
- ✅ Parse each .costyle file within bundle (return slice of UniversalRecipes)
- ✅ Generate .costylepack by bundling multiple .costyle files into ZIP
- ✅ Maintain bundle metadata (name, description if present in ZIP comment)
- ✅ Handle large bundles (50+ styles) efficiently (<5 seconds total conversion)
- ✅ Validate ZIP structure and report extraction errors (corrupt ZIP, missing files)

**AC-4: Round-Trip Conversion Testing**
- ✅ Round-trip conversion preserves 95%+ of parameter values
- ✅ Key adjustments (exposure, contrast, saturation) preserved exactly (no precision loss)
- ✅ Document known limitations of lossy conversions (parameters not representable in UniversalRecipe)
- ✅ Test suite includes real-world .costyle samples from Etsy/marketplaces (minimum 5 files)
- ✅ Automated tests verify round-trip accuracy for all test files
- ✅ Visual validation in Capture One software confirms output quality (manual spot-check)

**AC-5: CLI/TUI/Web Integration**
- ✅ CLI: `recipe convert input.costyle --to xmp` works correctly
- ✅ TUI: Format menu includes "Capture One" option with purple badge
- ✅ Web: Upload .costyle files via drag-drop, convert to other formats
- ✅ Format detection automatically identifies .costyle files (extension + magic bytes)
- ✅ Converter.Convert() function extended to handle Capture One format
- ✅ Help text and documentation updated for new format (README, CLI help, web FAQ)

## Traceability Mapping

| AC ID | Spec Section(s) | Component(s)/API(s) | Test Idea |
| ----- | --------------- | ------------------- | --------- |
| AC-1 | Data Models, APIs (Parse) | `costyle/parse.go`, `costyle/types.go` | Unit test: Parse sample .costyle, verify UniversalRecipe fields match |
| AC-2 | Data Models, APIs (Generate) | `costyle/generate.go`, `costyle/types.go` | Unit test: Generate .costyle from UniversalRecipe, validate XML structure |
| AC-3 | APIs (Pack functions) | `costyle/pack.go` | Unit test: ParsePack() on .costylepack ZIP, verify count and content |
| AC-4 | Workflows (Round-trip) | All costyle components | Integration test: .costyle → UR → .costyle, compare input/output |
| AC-5 | Services (Converter integration) | `converter/converter.go`, web/js/converter.js | Integration test: CLI convert command, TUI format selection, web upload |

**Test Coverage Targets:**
- Unit tests: `costyle/parse_test.go`, `costyle/generate_test.go`, `costyle/pack_test.go`
- Integration tests: `converter/converter_test.go` (extend with costyle cases)
- Manual tests: Load generated .costyle in Capture One software, verify visual output
- Coverage target: ≥85% for costyle package (consistent with Recipe standards)

## Risks, Assumptions, Open Questions

### Risks

**RISK-1: Capture One .costyle format undocumented**
- **Severity**: Medium
- **Impact**: May not have complete XML schema, parameter definitions
- **Mitigation**: Acquire multiple .costyle samples from Etsy/marketplaces, reverse-engineer schema from examples
- **Owner**: Dev (Epic 8)

**RISK-2: .costyle format version differences**
- **Severity**: Low
- **Impact**: Capture One 2023 vs 2024 vs 2025 may have different XML structures
- **Mitigation**: Document version differences in testdata README, implement version detection if needed
- **Owner**: Dev (Epic 8)

**RISK-3: Round-trip accuracy below 95%**
- **Severity**: Medium
- **Impact**: Some parameters may not map cleanly (tone curves, advanced color grading)
- **Mitigation**: Document lossy mappings in parameter-mapping.md, focus on core adjustments (exposure, contrast, saturation)
- **Owner**: Dev (Epic 8)

**RISK-4: Capture One software unavailable for validation**
- **Severity**: Low
- **Impact**: Cannot manually validate generated .costyle files load correctly
- **Mitigation**: Use Capture One trial version (free for 30 days), acquire at implementation time
- **Owner**: Justin (manual testing)

### Assumptions

**ASSUMPTION-1**: Capture One .costyle is XML-based (similar to Adobe XMP)
- **Rationale**: Community documentation suggests XML format
- **Validation**: Inspect sample .costyle files from Etsy
- **Risk if false**: Would require different parsing approach (binary format)

**ASSUMPTION-2**: encoding/xml can handle .costyle XML structure
- **Rationale**: Go standard library XML parser is robust, handles XMP successfully
- **Validation**: Prototype parsing with sample .costyle file
- **Risk if false**: May need custom XML parser (unlikely)

**ASSUMPTION-3**: .costylepack is standard ZIP format with .costyle entries
- **Rationale**: Similar to Adobe XMP sidecar bundles
- **Validation**: Inspect .costylepack samples with ZIP utilities
- **Risk if false**: May be proprietary archive format (unlikely)

**ASSUMPTION-4**: 95% round-trip accuracy is achievable
- **Rationale**: Core adjustments (exposure, contrast, saturation) map 1:1 to UniversalRecipe
- **Validation**: Implement and test round-trip conversion
- **Risk if false**: May need to lower success criteria to 90% (acceptable for MVP)

### Open Questions

**Q-1**: What Capture One parameters are most critical to users?
- **Impact**: Prioritize parameter mapping development
- **Resolution**: Focus on core adjustments (exposure, contrast, saturation, temperature, tint) for MVP
- **Owner**: Justin (product decision)

**Q-2**: Should Recipe support .costyle export from non-Capture One sources?
- **Impact**: Allows Nikon/Adobe users to export to Capture One format
- **Resolution**: YES - this is the core value proposition (universal conversion)
- **Owner**: Justin (product decision, confirmed)

**Q-3**: How to handle Capture One-specific features not in UniversalRecipe?
- **Impact**: Local adjustments, layers, masking not representable
- **Resolution**: Document as "not supported" in parameter-mapping.md, skip during conversion
- **Owner**: Dev (Epic 8)

**Q-4**: Should .costylepack bundle conversion preserve filenames?
- **Impact**: User experience for batch conversions
- **Resolution**: YES - preserve original filenames with new extension (e.g., "Preset1.costyle" → "Preset1.xmp")
- **Owner**: Dev (Epic 8)

## Test Strategy Summary

### Test Levels

**Unit Tests (costyle package):**
- `parse_test.go`: Test XML parsing with valid/invalid inputs, missing fields, edge cases
- `generate_test.go`: Test XML generation, validate output structure, parameter scaling
- `pack_test.go`: Test ZIP pack/unpack, bundle handling, error cases
- Coverage target: ≥85% for costyle package

**Integration Tests (converter package):**
- Extend `converter_test.go` with costyle conversion paths:
  - `costyle → np3`, `costyle → xmp`, `costyle → lrtemplate`, `costyle → dcp`
  - `np3 → costyle`, `xmp → costyle`, `lrtemplate → costyle`, `dcp → costyle`
- Round-trip tests: `costyle → UR → costyle` (verify 95%+ accuracy)
- Batch conversion: `.costylepack → multiple outputs`

**Manual Validation:**
- Load generated .costyle files in Capture One Pro (trial version)
- Verify visual output matches expected adjustments (exposure, contrast, saturation visible)
- Test with 3-5 different preset styles (portrait, landscape, product)
- Document results in test report

### Test Frameworks

- Go testing: `go test ./internal/formats/costyle`
- Table-driven tests for parameter mapping variations
- Benchmark tests: `BenchmarkParseCostyle`, `BenchmarkGenerateCostyle`
- Golden file testing: Compare generated XML to expected output

### Coverage of ACs

| AC ID | Test Type | Test Location | Coverage |
| ----- | --------- | ------------- | -------- |
| AC-1 | Unit | `parse_test.go` | All parse paths, error cases |
| AC-2 | Unit | `generate_test.go` | All generate paths, validation |
| AC-3 | Unit | `pack_test.go` | ZIP handling, bundle operations |
| AC-4 | Integration | `converter_test.go` | Round-trip accuracy verification |
| AC-5 | Integration + Manual | CLI tests, web UI testing | Format detection, UI integration |

### Edge Cases

**Edge Case Testing:**
- Empty .costyle file (minimal XML)
- Maximum parameter values (exposure = +2.0, contrast = +100)
- Minimum parameter values (exposure = -2.0, contrast = -100)
- Missing optional parameters (tone curves, clarity)
- .costylepack with 0 files (empty ZIP)
- .costylepack with 100+ files (stress test)
- Corrupt ZIP archive (truncated, invalid CRC)
- Malformed XML (missing closing tags, invalid attributes)
- Non-costyle files in .costylepack (skip gracefully)
- .costyle with unknown Capture One version (future-proofing)

---

**Next Steps:**
1. Acquire 5-10 .costyle sample files from Etsy/marketplaces for test fixtures
2. Set up `internal/formats/costyle/` package structure
3. Implement `parse.go` and `types.go` based on sample XML structure
4. Implement `generate.go` with parameter mapping
5. Implement `pack.go` for .costylepack bundle handling
6. Write unit tests achieving ≥85% coverage
7. Integrate with converter and add format detection
8. Update documentation (parameter-mapping.md, README.md)
9. Manual validation in Capture One Pro trial version
10. Mark epic-8 as "contexted" in sprint-status.yaml

---

## Action Items from Code Reviews

### Story 8-3: .costylepack Bundle Support (2025-11-09)

**Priority: Medium | Target Story: 8-5 (costyle-integration)**

1. **Wrap errors in ConversionError type** (Architecture Pattern 5 compliance)
   - **Current**: `pack.go` returns standard `error` type from Unpack() and Pack()
   - **Required**: All conversion errors must be wrapped in `ConversionError` type per architecture.md Pattern 5
   - **Implementation**: Add error wrapping during converter.Convert() integration in Story 8-5
   - **Status**: Non-blocking (deferred to Story 8-5 when integration happens)
   - **Code reference**: internal/formats/costyle/pack.go:39-133 (Unpack), 161-239 (Pack)

2. **Optional: ZIP bomb protection** (Security enhancement)
   - **Current**: No decompression ratio validation or file count limits
   - **Enhancement**: Add configurable limits (e.g., max 1000 files, 100MB total uncompressed size)
   - **Implementation**: Add validation in Unpack() before extraction loop
   - **Status**: Optional (reasonable limits exist via sample file constraints)
   - **Code reference**: internal/formats/costyle/pack.go:58-67 (empty bundle check)
