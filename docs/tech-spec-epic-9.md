# Epic Technical Specification: DCP Camera Profile Support

Date: 2025-11-08
Author: Justin
Epic ID: epic-9
Status: Draft

---

## Overview

Epic 9 adds DCP (DNG Camera Profile) format support to Recipe, enabling photographers to convert presets to/from Adobe's camera profile format used in Lightroom, Camera Raw, and DNG files. DCPs are TIFF-based files containing XML camera profile data, primarily used for camera calibration and color science customization.

This implementation extends Recipe's hub-and-spoke architecture to handle the TIFF container format (using github.com/google/tiff library) and XML camera profile data (using encoding/xml). The focus is on tone curve adjustments (exposure, contrast, highlights, shadows) which map well to UniversalRecipe, rather than full camera calibration (color matrices, dual illuminant profiles) which are beyond Recipe's scope.

## Objectives and Scope

**In Scope:**
- Parse DCP files (TIFF container + embedded XML camera profile)
- Generate DCP files from UniversalRecipe tone curve adjustments
- Support tone curve mapping (exposure, contrast, highlights, shadows, blacks, whites)
- Identity color matrices for non-calibration use cases
- Round-trip conversion testing with Adobe DCP samples
- Integration across all Recipe interfaces (CLI, TUI, Web)
- Parameter mapping documentation for DCP-specific adjustments

**Out of Scope (Path A):**
- Full camera calibration (ForwardMatrix, ColorMatrix, CalibrationIlluminant)
- Dual illuminant profiles (D65 + Tungsten light sources)
- Hue/Saturation/Value (HSV) look tables (complex color grading)
- Embedded DCP extraction from DNG image files (focus on standalone .dcp)
- DCP profile creation from scratch (color checker measurements)
- DCP metadata beyond core tone curves (copyright, camera model matching)

## System Architecture Alignment

**Components:**
- **New Package**: `internal/formats/dcp/` (following exact pattern of np3, xmp, lrtemplate, costyle)
  - `parse.go`: Orchestrates TIFF reading → XML extraction → UniversalRecipe
  - `generate.go`: Orchestrates UniversalRecipe → XML generation → TIFF embedding
  - `types.go`: Go structs matching DCP XML schema (tone curves, matrices)
  - `tiff.go`: Low-level TIFF tag operations using `github.com/google/tiff`
  - `profile.go`: Adobe camera profile XML parsing/generation
  - `testdata/`: Real DCP samples from Adobe

**Integration Points:**
- Extends `internal/converter/converter.go` to recognize .dcp format
- Adds DCP format badge (green #4CAF50) to web UI
- Updates `docs/parameter-mapping.md` with DCP tone curve mappings
- Leverages existing format detection logic (extension + TIFF magic bytes)

**New Dependency:**
- `github.com/google/tiff` - Complete TIFF library for reading/writing TIFF tags
  - Rationale: Go stdlib `image/tiff` is decoder-only, doesn't support custom tags
  - Google-maintained, stable, widely used
  - Approved in architecture decision (Decision 4)

**Constraints:**
- Must maintain <200ms DCP generation performance (slower than other formats due to TIFF overhead)
- Must preserve hub-and-spoke architecture (no direct format-to-format conversion)
- Must maintain ≥85% test coverage
- Focus on tone curves only (skip full color calibration for MVP)

## Detailed Design

### Services and Modules

| Module | Responsibility | Inputs | Outputs | Owner |
| ------ | -------------- | ------ | ------- | ----- |
| `dcp/parse.go` | Parse DCP TIFF → UniversalRecipe | .dcp file bytes | `*universal.Recipe`, error | Dev (Epic 9) |
| `dcp/generate.go` | Generate DCP TIFF ← UniversalRecipe | `*universal.Recipe` | .dcp file bytes, error | Dev (Epic 9) |
| `dcp/tiff.go` | TIFF tag read/write helpers | TIFF file, tag ID | Tag value / error | Dev (Epic 9) |
| `dcp/profile.go` | Adobe XML camera profile parsing | XML bytes | `Profile` struct / error | Dev (Epic 9) |
| `dcp/types.go` | DCP-specific types (matrices, curves) | - | Go struct definitions | Dev (Epic 9) |
| `converter/converter.go` | Format detection & routing (EXISTING) | File bytes, target format | Converted bytes, error | Dev (Epic 9 extension) |

### Data Models and Contracts

**DCP File Structure (TIFF container):**

```
DCP File (.dcp)
├── TIFF Header (Little/Big Endian marker)
├── Image File Directory (IFD)
│   ├── Standard TIFF tags (ImageWidth, ImageLength, etc.)
│   ├── Tag 50740: CameraProfile (XML data embedded)
│   │   └── Adobe Camera Profile XML (profile data)
│   └── Other DNG tags (optional)
```

**Adobe Camera Profile XML Structure (Embedded in Tag 50740):**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<crs:CameraProfile xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
  <crs:ProfileName>Recipe Converted Profile</crs:ProfileName>

  <!-- Tone Curve (Primary Focus for Recipe) -->
  <crs:ToneCurve>
    <rdf:Seq>
      <rdf:li>0, 0</rdf:li>      <!-- Input, Output (0-255 range) -->
      <rdf:li>64, 70</rdf:li>    <!-- Lifted shadows example -->
      <rdf:li>128, 128</rdf:li>  <!-- Midpoint -->
      <rdf:li>192, 190</rdf:li>  <!-- Crushed highlights example -->
      <rdf:li>255, 255</rdf:li>
    </rdf:Seq>
  </crs:ToneCurve>

  <!-- Color Matrices (Identity for non-calibration) -->
  <crs:ColorMatrix1>
    <rdf:Seq>
      <rdf:li>1.0 0.0 0.0</rdf:li>
      <rdf:li>0.0 1.0 0.0</rdf:li>
      <rdf:li>0.0 0.0 1.0</rdf:li>
    </rdf:Seq>
  </crs:ColorMatrix1>

  <!-- Optional: Hue/Saturation/Value tables (skip for MVP) -->
</crs:CameraProfile>
```

**Go Struct Mapping:**

```go
// types.go
package dcp

import "encoding/xml"

// CameraProfile represents Adobe DNG Camera Profile XML
type CameraProfile struct {
    XMLName      xml.Name   `xml:"CameraProfile"`
    Xmlns        string     `xml:"xmlns,attr"`
    ProfileName  string     `xml:"ProfileName"`
    ToneCurve    *ToneCurve `xml:"ToneCurve,omitempty"`
    ColorMatrix1 *Matrix    `xml:"ColorMatrix1,omitempty"`
    ColorMatrix2 *Matrix    `xml:"ColorMatrix2,omitempty"`
}

// ToneCurve represents a DCP tone curve (piecewise linear)
type ToneCurve struct {
    Points []ToneCurvePoint `xml:"Seq>li"`
}

// ToneCurvePoint is a single (input, output) point on the tone curve
type ToneCurvePoint struct {
    Input  int // 0-255
    Output int // 0-255
}

// Matrix represents a 3x3 color transformation matrix
type Matrix struct {
    Rows []MatrixRow `xml:"Seq>li"`
}

// MatrixRow is one row of the matrix (3 values)
type MatrixRow struct {
    Values [3]float64
}

// UnmarshalXML custom unmarshaler for ToneCurvePoint "0, 0" format
func (p *ToneCurvePoint) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var s string
    if err := d.DecodeElement(&s, &start); err != nil {
        return err
    }
    // Parse "input, output" format
    _, err := fmt.Sscanf(s, "%d, %d", &p.Input, &p.Output)
    return err
}

// MarshalXML custom marshaler for ToneCurvePoint "0, 0" format
func (p ToneCurvePoint) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    s := fmt.Sprintf("%d, %d", p.Input, p.Output)
    return e.EncodeElement(s, start)
}
```

**UniversalRecipe Mapping (Tone Curve Focus):**

Recipe maps tone curve adjustments to DCP:
- **UniversalRecipe → DCP Tone Curve**:
  - `universal.Exposure` → Vertical shift of entire tone curve
  - `universal.Contrast` → Steeper/shallower slope around midpoint
  - `universal.Highlights` → Adjust top-end curve points (192-255 range)
  - `universal.Shadows` → Adjust bottom-end curve points (0-64 range)
  - `universal.Blacks` → Clamp black point (input 0)
  - `universal.Whites` → Clamp white point (input 255)

- **Color Matrices**: Use identity matrices (1.0 on diagonal, 0.0 elsewhere) - Recipe doesn't perform camera calibration

- **Unsupported DCP Features** (documented in parameter-mapping.md):
  - Dual illuminant profiles (D65 + Tungsten)
  - Hue/Saturation/Value look tables
  - Forward/Reverse color matrices (camera calibration)
  - Camera model metadata matching

### APIs and Interfaces

**dcp.Parse Function:**

```go
// parse.go
package dcp

import (
    "fmt"
    "github.com/google/tiff"
    "github.com/justin/recipe/internal/universal"
)

// Parse parses a DCP file into UniversalRecipe
// Extracts tone curve from embedded Adobe Camera Profile XML
func Parse(data []byte) (*universal.Recipe, error) {
    // Read TIFF structure
    tf, err := readTIFF(data)
    if err != nil {
        return nil, fmt.Errorf("failed to read TIFF structure: %w", err)
    }

    // Extract CameraProfile tag (50740)
    profileXML, err := extractCameraProfileTag(tf)
    if err != nil {
        return nil, fmt.Errorf("failed to extract camera profile: %w", err)
    }

    // Parse XML camera profile
    profile, err := parseProfile(profileXML)
    if err != nil {
        return nil, fmt.Errorf("failed to parse camera profile XML: %w", err)
    }

    // Convert DCP tone curve → UniversalRecipe
    return profileToUniversal(profile), nil
}
```

**dcp.Generate Function:**

```go
// generate.go
package dcp

import (
    "fmt"
    "github.com/google/tiff"
    "github.com/justin/recipe/internal/universal"
)

// Generate creates a DCP file from UniversalRecipe
// Embeds tone curve as Adobe Camera Profile XML in TIFF tag 50740
func Generate(recipe *universal.Recipe) ([]byte, error) {
    // Convert UniversalRecipe → DCP tone curve
    profile := universalToProfile(recipe)

    // Generate XML camera profile
    profileXML, err := generateProfile(profile)
    if err != nil {
        return nil, fmt.Errorf("failed to generate camera profile XML: %w", err)
    }

    // Create TIFF with embedded profile
    tiffData, err := createTIFF(profileXML)
    if err != nil {
        return nil, fmt.Errorf("failed to create TIFF: %w", err)
    }

    return tiffData, nil
}
```

**dcp/tiff.go Helper Functions:**

```go
// tiff.go
package dcp

import (
    "bytes"
    "github.com/google/tiff"
)

const (
    // DNG/DCP TIFF tags
    TagCameraProfile = 50740 // Adobe Camera Profile XML
)

// readTIFF reads TIFF structure from bytes
func readTIFF(data []byte) (*tiff.TIFF, error) {
    r := bytes.NewReader(data)
    return tiff.Parse(r, nil, nil)
}

// extractCameraProfileTag extracts XML from tag 50740
func extractCameraProfileTag(tf *tiff.TIFF) ([]byte, error) {
    // Navigate IFD to find tag 50740
    ifd := tf.IFDs()[0] // First IFD
    entry := ifd.GetField(TagCameraProfile)
    if entry == nil {
        return nil, fmt.Errorf("camera profile tag not found (tag 50740)")
    }

    // Extract byte data
    data, ok := entry.Value().([]byte)
    if !ok {
        return nil, fmt.Errorf("camera profile tag has invalid type")
    }

    return data, nil
}

// createTIFF creates TIFF file with embedded camera profile XML
func createTIFF(profileXML []byte) ([]byte, error) {
    tf := tiff.New()

    // Create IFD with required TIFF tags
    ifd := tiff.NewIFD(tf)
    ifd.SetField(tiff.NewField(tiff.ImageWidth, tiff.UInt32(1), tiff.NewUInt32(1)))
    ifd.SetField(tiff.NewField(tiff.ImageLength, tiff.UInt32(1), tiff.NewUInt32(1)))

    // Embed camera profile XML in tag 50740
    ifd.SetField(tiff.NewField(TagCameraProfile, tiff.Byte, profileXML))

    tf.SetIFD(ifd)

    // Write TIFF to bytes
    var buf bytes.Buffer
    if err := tf.Write(&buf); err != nil {
        return nil, fmt.Errorf("failed to write TIFF: %w", err)
    }

    return buf.Bytes(), nil
}
```

**dcp/profile.go XML Parsing:**

```go
// profile.go
package dcp

import (
    "encoding/xml"
    "fmt"
)

// parseProfile parses Adobe Camera Profile XML
func parseProfile(data []byte) (*CameraProfile, error) {
    var profile CameraProfile
    if err := xml.Unmarshal(data, &profile); err != nil {
        return nil, fmt.Errorf("failed to parse camera profile XML: %w", err)
    }

    return &profile, nil
}

// generateProfile generates Adobe Camera Profile XML
func generateProfile(profile *CameraProfile) ([]byte, error) {
    // Set namespace
    profile.Xmlns = "http://ns.adobe.com/camera-raw-settings/1.0/"

    data, err := xml.MarshalIndent(profile, "", "  ")
    if err != nil {
        return nil, fmt.Errorf("failed to generate camera profile XML: %w", err)
    }

    // Prepend XML declaration
    return append([]byte(xml.Header), data...), nil
}

// profileToUniversal converts DCP tone curve → UniversalRecipe
func profileToUniversal(profile *CameraProfile) *universal.Recipe {
    recipe := &universal.Recipe{}

    if profile.ToneCurve != nil {
        // Analyze tone curve shape to extract exposure, contrast, highlights, shadows
        curve := profile.ToneCurve

        // Midpoint shift → Exposure (compare point at input=128)
        midpoint := findPointOutput(curve, 128)
        recipe.Exposure = float64(midpoint-128) / 128.0 * 2.0 // Scale to -2.0..+2.0

        // Slope → Contrast (compare endpoints vs midpoint)
        slope := calculateSlope(curve)
        recipe.Contrast = (slope - 1.0) // Deviation from linear (1.0 = neutral)

        // Top-end points → Highlights (average output at 192-255)
        recipe.Highlights = calculateHighlights(curve)

        // Bottom-end points → Shadows (average output at 0-64)
        recipe.Shadows = calculateShadows(curve)
    }

    // Color matrices ignored (identity assumed)

    return recipe
}

// universalToProfile converts UniversalRecipe → DCP tone curve
func universalToProfile(recipe *universal.Recipe) *CameraProfile {
    profile := &CameraProfile{
        ProfileName: "Recipe Converted Profile",
    }

    // Build tone curve from UniversalRecipe adjustments
    curve := &ToneCurve{
        Points: buildToneCurve(recipe.Exposure, recipe.Contrast, recipe.Highlights, recipe.Shadows),
    }
    profile.ToneCurve = curve

    // Identity color matrix (no camera calibration)
    profile.ColorMatrix1 = identityMatrix()

    return profile
}

// buildToneCurve creates piecewise linear tone curve from adjustments
func buildToneCurve(exposure, contrast, highlights, shadows float64) []ToneCurvePoint {
    // Base curve: 5 points (0, 64, 128, 192, 255)
    points := []ToneCurvePoint{
        {Input: 0, Output: 0},
        {Input: 64, Output: 64},
        {Input: 128, Output: 128},
        {Input: 192, Output: 192},
        {Input: 255, Output: 255},
    }

    // Apply exposure (vertical shift)
    shift := int(exposure * 64) // -2.0..+2.0 → -128..+128
    for i := range points {
        points[i].Output = clamp(points[i].Output+shift, 0, 255)
    }

    // Apply contrast (steepen/flatten slope around midpoint)
    contrastFactor := 1.0 + contrast
    for i := range points {
        delta := float64(points[i].Output - 128)
        points[i].Output = clamp(128+int(delta*contrastFactor), 0, 255)
    }

    // Apply highlights (adjust top-end curve)
    highlightShift := int(highlights * 32)
    points[3].Output = clamp(points[3].Output+highlightShift, 0, 255) // 192 point
    points[4].Output = clamp(points[4].Output+highlightShift, 0, 255) // 255 point

    // Apply shadows (adjust bottom-end curve)
    shadowShift := int(shadows * 32)
    points[0].Output = clamp(points[0].Output+shadowShift, 0, 255) // 0 point
    points[1].Output = clamp(points[1].Output+shadowShift, 0, 255) // 64 point

    return points
}

func clamp(value, min, max int) int {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

func identityMatrix() *Matrix {
    return &Matrix{
        Rows: []MatrixRow{
            {Values: [3]float64{1.0, 0.0, 0.0}},
            {Values: [3]float64{0.0, 1.0, 0.0}},
            {Values: [3]float64{0.0, 0.0, 1.0}},
        },
    }
}
```

### Workflows and Sequencing

**Parse Workflow (.dcp → UniversalRecipe):**

1. User uploads .dcp file via CLI/TUI/Web
2. `converter.detectFormat(data)` identifies file as DCP (TIFF magic bytes + .dcp extension)
3. `dcp.Parse(data)` called:
   a. `readTIFF(data)` parses TIFF structure using github.com/google/tiff
   b. `extractCameraProfileTag()` extracts XML from TIFF tag 50740
   c. `parseProfile(xml)` unmarshals Adobe Camera Profile XML
   d. `profileToUniversal()` analyzes tone curve shape:
      - Extract exposure from midpoint shift
      - Extract contrast from curve slope
      - Extract highlights from top-end curve points (192-255)
      - Extract shadows from bottom-end curve points (0-64)
4. Return `*universal.Recipe` to converter
5. Converter routes to target format generator

**Generate Workflow (UniversalRecipe → .dcp):**

1. Converter calls `dcp.Generate(recipe)`
2. `universalToProfile()` builds DCP tone curve:
   a. Start with linear 5-point curve (0, 64, 128, 192, 255)
   b. Apply exposure adjustment (vertical shift)
   c. Apply contrast adjustment (steepen/flatten slope)
   d. Apply highlights adjustment (modify 192-255 points)
   e. Apply shadows adjustment (modify 0-64 points)
   f. Clamp all output values to 0-255 range
   g. Create identity color matrix (no calibration)
3. `generateProfile()` marshals CameraProfile to XML
4. `createTIFF()` creates TIFF file:
   a. Initialize TIFF structure with minimal IFD
   b. Embed camera profile XML in tag 50740
   c. Write TIFF to bytes
5. Return TIFF bytes to converter
6. Converter returns to user (download/write file)

## Non-Functional Requirements

### Performance

**Targets:**
- DCP parsing: <200ms (99th percentile) - slower than other formats due to TIFF overhead
- DCP generation: <200ms (99th percentile)
- Memory usage: <15MB per conversion (TIFF library may allocate more than XML-only formats)
- Web WASM build: DCP conversion functional (WASM overhead <20% due to TIFF complexity)

**Optimization Strategies:**
- Stream TIFF tag reading (don't load entire TIFF into memory)
- Cache parsed profile structure (avoid re-parsing XML)
- Use efficient tone curve algorithms (5-point curves, no complex interpolation)

**Performance Monitoring:**
- Benchmark tests: `BenchmarkParseDCP`, `BenchmarkGenerateDCP`
- Compare to existing formats (expect 2x slower due to TIFF operations)
- Document baseline: <150ms for typical DCP file

### Security

**Threats & Mitigations:**
- **TIFF Bomb**: Limit TIFF file size to <10MB, validate IFD structure before parsing
- **XML Bomb**: encoding/xml resistant by default (no entity expansion)
- **Malformed TIFF**: github.com/google/tiff handles gracefully, returns errors (no panics)
- **Invalid Camera Profile XML**: Validate against Adobe namespace, reject unknown tags
- **File Upload**: Web UI validates .dcp extension, server validates TIFF magic bytes (II or MM)

**Privacy:**
- All processing client-side (WASM) or local (CLI/TUI)
- DCP files may contain camera model metadata - preserve but don't expose unnecessarily
- Zero external API calls

**Input Validation:**
- Verify TIFF magic bytes (II for little-endian, MM for big-endian)
- Check tag 50740 exists and contains valid XML
- Validate tone curve points within 0-255 range
- Reject files >10MB

### Reliability/Availability

**Error Handling:**
- Wrapped errors with context: `fmt.Errorf("failed to read TIFF: %w", err)`
- Graceful degradation: Missing tone curve uses linear (passthrough)
- TIFF library errors handled explicitly (no panics)
- Round-trip validation: Parse(Generate(recipe)) preserves tone curve shape

**Backwards Compatibility:**
- Support DCP v1.0-v1.6 formats (Adobe DNG Specification)
- Document known version differences in testdata README

**Failure Modes:**
- Invalid TIFF structure: Return error, don't crash
- Missing tag 50740: Return error ("not a valid DCP file")
- Malformed XML: Return clear parsing error
- Out-of-range tone curve values: Clamp to 0-255

### Observability

**Logging:**
- Use Go's `log/slog` structured logging
- Log level DEBUG: TIFF tag reading, XML parsing steps, tone curve calculations
- Log level INFO: DCP conversion success/failure
- Log level WARN: Skipped unsupported DCP features (color matrices, HSV tables)
- Log level ERROR: TIFF parsing failures, XML errors

**Metrics:**
- Count DCP conversions by source/target format pair
- Track parse/generate latency (p50, p95, p99)
- Monitor TIFF library memory allocations

**Debugging:**
- Verbose mode dumps intermediate tone curve points as JSON
- Test suite includes visual comparison (generated DCP applied in Lightroom)

## Dependencies and Integrations

**Go Module Dependencies:**

New Dependency for Epic 9:
```
require github.com/google/tiff latest
```

Rationale:
- Go stdlib `image/tiff` is decoder-only, doesn't support writing or custom tags
- `github.com/google/tiff` provides complete TIFF read/write capabilities
- Google-maintained, stable, widely used in imaging applications
- Necessary for embedding camera profile XML in TIFF tag 50740

Standard Library Usage:
- `encoding/xml`: Adobe Camera Profile XML parsing/generation
- `fmt`: Error wrapping
- `bytes`, `io`: Stream handling for TIFF operations

**External Integrations:**
- **Adobe Camera Raw / Lightroom** (validation): Generated DCP files must load in Adobe software (manual testing)
- **Adobe DCP samples** (test data): Download reference DCP files from Adobe for test fixtures

**Version Constraints:**
- Go 1.25.1+ (required for go:wasmexport)
- `github.com/google/tiff` latest stable release
- Adobe DNG Specification 1.6 (DCP format reference)

## Acceptance Criteria (Authoritative)

**AC-1: Parse DCP Files**
- ✅ Read TIFF-based DCP file structure using github.com/google/tiff
- ✅ Extract XML camera profile data from TIFF tag 50740
- ✅ Parse color matrices (identity matrices recognized, full calibration optional)
- ✅ Parse tone curve adjustments (piecewise linear curve points)
- ✅ Parse hue/saturation/value tables if present (skip for MVP, log warning)
- ✅ Handle DCP v1.x format variations (v1.0-v1.6 supported)
- ✅ Validate DCP structure and report parsing errors with clear messages
- ✅ Support both embedded (in DNG) and standalone DCP files (standalone priority for MVP)

**AC-2: Generate DCP Files**
- ✅ Generate TIFF-based DCP file structure per Adobe DNG Specification 1.6
- ✅ Embed XML camera profile data in TIFF tag 50740
- ✅ Map UniversalRecipe color parameters to DCP tone curve equivalents
- ✅ Generate required matrices (identity matrices for non-calibration use)
- ✅ Create tone curves from exposure/contrast/highlights/shadows adjustments
- ✅ Validate generated DCP against Adobe spec (well-formed TIFF + valid XML)
- ✅ Generated DCPs load in Adobe Camera Raw and Lightroom without errors
- ✅ Document mapping limitations and best practices (parameter-mapping.md)

**AC-3: DCP Parameter Mapping**
- ✅ Document which UniversalRecipe parameters map to DCP tone curves
- ✅ Identify unsupported DCP features (dual illuminant, full color calibration, HSV tables)
- ✅ Define conversion formulas for tone curves (exposure/contrast/highlights/shadows → 5-point curve)
- ✅ Handle color space transformations correctly (identity matrices, no calibration)
- ✅ Create reference documentation for DCP mapping (parameter-mapping.md section)
- ✅ Include examples of common conversions (NP3 → DCP, XMP → DCP)
- ✅ Test mapping with real DCP samples from Adobe (minimum 3 sample files)

**AC-4: Compatibility Validation**
- ✅ Generated DCPs load without errors in Adobe Camera Raw
- ✅ Generated DCPs load without errors in Lightroom Classic
- ✅ Preset adjustments render visually similar to original (tone curve visible)
- ✅ Performance: DCP generation completes in <200ms (accept slower than other formats)
- ✅ Test with multiple camera models (Nikon, Canon samples - identity matrices work universally)
- ✅ Document known compatibility issues or edge cases (testdata README)

## Traceability Mapping

| AC ID | Spec Section(s) | Component(s)/API(s) | Test Idea |
| ----- | --------------- | ------------------- | --------- |
| AC-1 | Data Models, APIs (Parse) | `dcp/parse.go`, `dcp/tiff.go`, `dcp/profile.go` | Unit test: Parse Adobe sample DCP, verify tone curve extraction |
| AC-2 | Data Models, APIs (Generate) | `dcp/generate.go`, `dcp/tiff.go`, `dcp/profile.go` | Unit test: Generate DCP from UniversalRecipe, validate TIFF structure |
| AC-3 | APIs (Tone curve mapping) | `dcp/profile.go` (profileToUniversal, universalToProfile) | Unit test: Round-trip tone curve, verify curve shape preservation |
| AC-4 | Workflows (Adobe validation) | All dcp components | Manual test: Load generated DCP in Lightroom, verify visual output |

**Test Coverage Targets:**
- Unit tests: `dcp/parse_test.go`, `dcp/generate_test.go`, `dcp/tiff_test.go`, `dcp/profile_test.go`
- Integration tests: `converter/converter_test.go` (extend with DCP cases)
- Manual tests: Load generated DCP in Adobe Camera Raw/Lightroom, compare to original
- Coverage target: ≥85% for dcp package

## Risks, Assumptions, Open Questions

### Risks

**RISK-1: TIFF library complexity**
- **Severity**: Medium
- **Impact**: github.com/google/tiff may have learning curve, unexpected behaviors
- **Mitigation**: Start with simple TIFF operations (read tag 50740), expand as needed
- **Owner**: Dev (Epic 9)

**RISK-2: Adobe DCP format underdocumented**
- **Severity**: Medium
- **Impact**: DNG Specification may not cover all edge cases
- **Mitigation**: Test with multiple real Adobe DCP samples, reverse-engineer if needed
- **Owner**: Dev (Epic 9)

**RISK-3: Tone curve mapping imprecise**
- **Severity**: Low
- **Impact**: DCP tone curve is piecewise linear, may not perfectly match UniversalRecipe adjustments
- **Mitigation**: Use 5-point curve (sufficient for most presets), document limitations
- **Owner**: Dev (Epic 9)

**RISK-4: Adobe software unavailable for validation**
- **Severity**: Low
- **Impact**: Cannot manually validate DCPs load correctly
- **Mitigation**: Use Lightroom trial (free for 7 days), Camera Raw (free with Photoshop trial)
- **Owner**: Justin (manual testing)

### Assumptions

**ASSUMPTION-1**: DCP tone curves sufficient for Recipe MVP
- **Rationale**: Full camera calibration (color matrices) is complex, tone curves are 80% of value
- **Validation**: Test with typical preset conversions (NP3 → DCP)
- **Risk if false**: Users may need full color calibration (defer to future enhancement)

**ASSUMPTION-2**: Identity color matrices acceptable
- **Rationale**: Most presets adjust tone curves, not camera-specific color calibration
- **Validation**: Generated DCPs load in Adobe software without errors
- **Risk if false**: May need to implement ForwardMatrix/ColorMatrix (complex)

**ASSUMPTION-3**: Standalone .dcp files are primary use case
- **Rationale**: Embedded DCP extraction from DNG images is complex
- **Validation**: User feedback on whether embedded DCP support needed
- **Risk if false**: May need to add DNG image parsing (defer to future)

**ASSUMPTION-4**: 5-point tone curve sufficient
- **Rationale**: More points add complexity without significant quality gain
- **Validation**: Visual comparison of generated DCP vs original preset
- **Risk if false**: May need 9-point or 17-point curves (simple to add more points)

### Open Questions

**Q-1**: Should Recipe support embedded DCP extraction from DNG images?
- **Impact**: Significantly increases complexity (full DNG image parser)
- **Resolution**: NO for MVP - focus on standalone .dcp files only
- **Owner**: Justin (product decision)

**Q-2**: Should Recipe generate dual illuminant profiles (D65 + Tungsten)?
- **Impact**: Doubles complexity, most users don't need it
- **Resolution**: NO for MVP - single illuminant (D65 assumed) with identity matrices
- **Owner**: Justin (product decision, confirmed)

**Q-3**: How to handle DCP HSV (Hue/Saturation/Value) look tables?
- **Impact**: Advanced color grading not representable in UniversalRecipe
- **Resolution**: SKIP for MVP - log warning if HSV tables present, don't convert
- **Owner**: Dev (Epic 9), document in parameter-mapping.md

**Q-4**: Should DCP generation include camera model metadata?
- **Impact**: Camera model matching for Adobe software compatibility
- **Resolution**: NO for MVP - use generic profile name "Recipe Converted Profile"
- **Owner**: Dev (Epic 9)

## Test Strategy Summary

### Test Levels

**Unit Tests (dcp package):**
- `parse_test.go`: TIFF reading, camera profile XML extraction, tone curve parsing
- `generate_test.go`: TIFF creation, camera profile XML generation, tone curve building
- `tiff_test.go`: TIFF tag operations (read tag 50740, write tag 50740)
- `profile_test.go`: XML parsing/generation, tone curve algorithms
- Coverage target: ≥85% for dcp package

**Integration Tests (converter package):**
- Extend `converter_test.go` with DCP conversion paths:
  - `dcp → np3`, `dcp → xmp`, `dcp → lrtemplate`, `dcp → costyle`
  - `np3 → dcp`, `xmp → dcp`, `lrtemplate → dcp`, `costyle → dcp`
- Round-trip tests: `dcp → UR → dcp` (verify tone curve shape preservation)

**Manual Validation:**
- Load generated DCP files in Adobe Camera Raw (Photoshop trial)
- Load generated DCP files in Lightroom Classic (trial version)
- Apply DCP to sample RAW images, verify tone curve adjustments visible
- Test with 3-5 different preset styles (high contrast, lifted shadows, crushed highlights)
- Document results in test report

### Test Frameworks

- Go testing: `go test ./internal/formats/dcp`
- Table-driven tests for tone curve calculations
- Benchmark tests: `BenchmarkParseDCP`, `BenchmarkGenerateDCP`
- Visual regression: Screenshot Lightroom before/after DCP application

### Coverage of ACs

| AC ID | Test Type | Test Location | Coverage |
| ----- | --------- | ------------- | -------- |
| AC-1 | Unit | `parse_test.go`, `tiff_test.go` | TIFF parsing, XML extraction |
| AC-2 | Unit | `generate_test.go`, `tiff_test.go` | TIFF generation, XML embedding |
| AC-3 | Unit | `profile_test.go` | Tone curve mapping formulas |
| AC-4 | Manual | Adobe Camera Raw / Lightroom | Visual validation, load errors |

### Edge Cases

**Edge Case Testing:**
- Empty DCP (minimal TIFF with no tone curve)
- Maximum tone curve adjustments (exposure +2.0, contrast +1.0)
- Minimum tone curve adjustments (exposure -2.0, contrast -1.0)
- Missing color matrices (should still generate identity)
- DCP with HSV tables (skip gracefully, log warning)
- Corrupt TIFF structure (truncated file, invalid IFD)
- Malformed camera profile XML (missing namespace, invalid points)
- DCP v1.0 vs v1.6 format differences
- Large DCP files (>5MB - stress test)
- DCP with unknown camera model metadata

---

**Next Steps:**
1. Acquire 3-5 DCP sample files from Adobe for test fixtures
2. Add `github.com/google/tiff` dependency to go.mod
3. Set up `internal/formats/dcp/` package structure
4. Implement `tiff.go` TIFF tag operations
5. Implement `profile.go` Adobe Camera Profile XML parsing
6. Implement `parse.go` orchestrating TIFF → XML → UniversalRecipe
7. Implement `generate.go` orchestrating UniversalRecipe → XML → TIFF
8. Write unit tests achieving ≥85% coverage
9. Integrate with converter and add format detection
10. Update documentation (parameter-mapping.md, README.md)
11. Manual validation in Adobe Camera Raw / Lightroom trial
12. Mark epic-9 as "contexted" in sprint-status.yaml
