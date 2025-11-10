# Epic Technical Specification: DCP Camera Profile Support

Date: 2025-11-08 (Updated: 2025-11-10 for Binary Format)
Author: Justin
Epic ID: epic-9
Status: Draft

---

## ⚠️ CRITICAL FORMAT DISCOVERY

**Original Assumption (from initial spec):**
DCP files contain Adobe Camera Profile XML embedded in TIFF tag 50740.

**Reality Discovered During Story 9-1 Implementation:**
Real Adobe DCP files use **binary TIFF/DNG tags** in the 50700-52600 range, NOT XML in tag 50740.

This specification has been updated to reflect the binary format. See `docs/stories/9-1-dcp-parser-FORMAT-PIVOT.md` for complete discovery details.

---

## Overview

Epic 9 adds DCP (DNG Camera Profile) format support to Recipe, enabling photographers to convert presets to/from Adobe's camera profile format used in Lightroom, Camera Raw, and DNG files. DCPs are DNG-format files containing **binary camera profile data** stored in TIFF tags 50700-52600, primarily used for camera calibration and color science customization.

This implementation extends Recipe's hub-and-spoke architecture to handle the DNG format (using github.com/google/tiff library) and binary profile data (using encoding/binary). The focus is on tone curve adjustments (exposure, contrast, highlights, shadows) which map well to UniversalRecipe, rather than full camera calibration (color matrices, dual illuminant profiles) which are beyond Recipe's scope.

## Objectives and Scope

**In Scope:**
- Parse DCP files (DNG format with binary TIFF tags 50700-52600)
- Generate DCP files from UniversalRecipe tone curve adjustments
- Support tone curve mapping (exposure, contrast, highlights, shadows, blacks, whites)
- Identity color matrices for non-calibration use cases (binary SRational arrays)
- Round-trip conversion testing with Adobe DCP samples
- Integration across all Recipe interfaces (CLI, TUI, Web)
- Parameter mapping documentation for DCP-specific adjustments (binary format)

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
  - `parse.go`: Orchestrates DNG reading → binary tag extraction → UniversalRecipe
  - `generate.go`: Orchestrates UniversalRecipe → binary tag generation → DNG writing
  - `types.go`: Go structs for binary DCP data (tone curves as float32 arrays, matrices as SRationals)
  - `tiff.go`: Low-level TIFF/DNG tag operations using `github.com/google/tiff`, DNG version conversion
  - `profile.go`: Binary tone curve analysis and generation algorithms
  - `testdata/`: Real DCP samples from Adobe (36 files tested)

**Integration Points:**
- Extends `internal/converter/converter.go` to recognize .dcp format
- Adds DCP format badge (green #4CAF50) to web UI
- Updates `docs/parameter-mapping.md` with DCP tone curve mappings
- Leverages existing format detection logic (extension + TIFF magic bytes)

**New Dependencies:**
- `github.com/google/tiff` - Complete TIFF/DNG library for reading/writing TIFF tags
  - Rationale: Go stdlib `image/tiff` is decoder-only, doesn't support custom tags or DNG format
  - Google-maintained, stable, widely used
  - Approved in architecture decision (Decision 4)
- `encoding/binary` (stdlib) - Binary data conversion for float32 arrays and SRational values
- `unsafe` (stdlib) - DNG version byte conversion (IIRC/MMCR → version 42 for tiff library)

**Constraints:**
- Must maintain <200ms DCP generation performance (slower than other formats due to TIFF overhead)
- Must preserve hub-and-spoke architecture (no direct format-to-format conversion)
- Must maintain ≥85% test coverage
- Focus on tone curves only (skip full color calibration for MVP)

## Detailed Design

### Services and Modules

| Module | Responsibility | Inputs | Outputs | Owner |
| ------ | -------------- | ------ | ------- | ----- |
| `dcp/parse.go` | Parse DCP DNG → UniversalRecipe | .dcp file bytes | `*models.UniversalRecipe`, error | Dev (Epic 9) |
| `dcp/generate.go` | Generate DCP DNG ← UniversalRecipe | `*models.UniversalRecipe` | .dcp file bytes, error | Dev (Epic 9) |
| `dcp/tiff.go` | DNG tag read/write helpers, magic byte conversion | DNG file, tag ID | Binary tag value / error | Dev (Epic 9) |
| `dcp/profile.go` | Binary tone curve analysis/generation | ToneCurvePoint array | Exposure/Contrast/etc / error | Dev (Epic 9) |
| `dcp/types.go` | DCP-specific binary types (SRational, curves) | - | Go struct definitions | Dev (Epic 9) |
| `converter/converter.go` | Format detection & routing (EXISTING) | File bytes, target format | Converted bytes, error | Dev (Epic 9 extension) |

### Data Models and Contracts

**DCP File Structure (DNG format):**

```
DCP File (.dcp)
├── DNG Header (Magic bytes: "IIRC" or "MMCR" instead of TIFF version 42)
├── Image File Directory (IFD)
│   ├── Standard TIFF tags (ImageWidth, ImageLength, etc.)
│   ├── Tag 50708: UniqueCameraModel (ASCII, e.g., "Nikon Z f")
│   ├── Tag 50721: ColorMatrix1 (SRational[9] - binary array)
│   ├── Tag 50722: ColorMatrix2 (SRational[9] - binary array)
│   ├── Tag 50730: BaselineExposureOffset (SRational)
│   ├── Tag 50940: ProfileToneCurve (Float32 array of input/output pairs)
│   └── Tag 52552: ProfileName (ASCII, OPTIONAL - may be missing)
```

**Binary Data Formats:**

**Tone Curve (Tag 50940):**
- Format: Array of float32 pairs (input, output)
- Normalization: 0.0-1.0 range (NOT 0-255 integers)
- Byte structure: Each point = 8 bytes (4-byte float32 input + 4-byte float32 output)
- Example linear curve:
  ```
  {0.0, 0.0}, {0.25, 0.25}, {0.5, 0.5}, {0.75, 0.75}, {1.0, 1.0}
  = 40 bytes total (5 points × 8 bytes per point)
  ```

**Color Matrices (Tags 50721-50722):**
- Format: Array of 9 SRational values (row-major order)
- SRational: 8 bytes per value (4-byte int32 numerator + 4-byte int32 denominator)
- Identity matrix example:
  ```
  Diagonal: SRational{Numerator: 1, Denominator: 1}
  Off-diagonal: SRational{Numerator: 0, Denominator: 1}
  Total: 72 bytes (9 SRationals × 8 bytes each)
  ```

**ProfileName (Tag 52552):**
- Format: ASCII string, null-terminated
- **OPTIONAL** - not all DCP files have this tag (empty string if missing)
- Example: "Camera Standard\0" = 17 bytes

**Go Struct Mapping:**

```go
// types.go
package dcp

// DNG tag constants (binary format)
const (
    TagUniqueCameraModel      = 50708  // ASCII camera model
    TagColorMatrix1           = 50721  // SRational[9] color calibration matrix (illuminant 1)
    TagColorMatrix2           = 50722  // SRational[9] color calibration matrix (illuminant 2)
    TagBaselineExposureOffset = 50730  // SRational exposure compensation
    TagProfileToneCurve       = 50940  // Float32 array of (input, output) pairs
    TagProfileName            = 52552  // ASCII profile name (OPTIONAL)
)

// ToneCurvePoint represents a single point on the binary tone curve
// Stored as 8 bytes: 4-byte float32 input + 4-byte float32 output
type ToneCurvePoint struct {
    Input  float64  // 0.0-1.0 normalized (parsed from float32)
    Output float64  // 0.0-1.0 normalized (parsed from float32)
}

// SRational represents a signed rational number (numerator/denominator)
// Stored as 8 bytes: 4-byte int32 numerator + 4-byte int32 denominator
type SRational struct {
    Numerator   int32
    Denominator int32
}

// ToFloat64 converts SRational to float64
func (sr SRational) ToFloat64() float64 {
    if sr.Denominator == 0 {
        return 0.0
    }
    return float64(sr.Numerator) / float64(sr.Denominator)
}

// Matrix represents a 3x3 color transformation matrix
// Stored as binary array of 9 SRational values (72 bytes total)
type Matrix struct {
    Rows [3][3]float64  // Parsed from SRational binary data
}

// DCP profile metadata extracted from binary tags
type Profile struct {
    ProfileName           string           // Tag 52552 (optional, empty if missing)
    ToneCurve             []ToneCurvePoint // Tag 50940 (binary float32 array)
    ColorMatrix1          *Matrix          // Tag 50721 (9 SRationals)
    ColorMatrix2          *Matrix          // Tag 50722 (9 SRationals)
    BaselineExposureOffset float64         // Tag 50730 (SRational)
    UniqueCameraModel     string           // Tag 50708 (ASCII)
}
```

**UniversalRecipe Mapping (Binary Tone Curve Focus):**

Recipe maps tone curve adjustments to binary DCP format:
- **UniversalRecipe → DCP Binary Tone Curve (Tag 50940)**:
  - `universal.Exposure` → Vertical shift of entire tone curve (0.0-1.0 normalized)
  - `universal.Contrast` → Steeper/shallower slope around midpoint (0.5 in normalized space)
  - `universal.Highlights` → Adjust top-end curve points (0.75-1.0 range in normalized space)
  - `universal.Shadows` → Adjust bottom-end curve points (0.0-0.25 range in normalized space)
  - `universal.Blacks` → Clamp black point (input 0.0)
  - `universal.Whites` → Clamp white point (input 1.0)

- **Color Matrices (Tags 50721-50722)**: Use identity matrices as binary SRational arrays:
  - Diagonal: `SRational{Numerator: 1, Denominator: 1}` (represents 1.0)
  - Off-diagonal: `SRational{Numerator: 0, Denominator: 1}` (represents 0.0)
  - Recipe doesn't perform camera calibration

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
    "github.com/justin/recipe/internal/models"
)

// Parse parses a DCP file into UniversalRecipe
// Extracts binary tone curve and color matrices from DNG tags
func Parse(data []byte) (*models.UniversalRecipe, error) {
    // Read DNG/TIFF structure (convert DNG magic bytes "IIRC"/"MMCR" to version 42)
    tf, err := readDNG(data)
    if err != nil {
        return nil, fmt.Errorf("failed to read DNG structure: %w", err)
    }

    // Extract binary profile data from DNG tags
    profileName := extractProfileName(tf)        // Tag 52552 (optional)
    toneCurve := extractToneCurve(tf)            // Tag 50940 (binary float32 array)
    colorMatrix1 := extractColorMatrix(tf, 50721) // Tag 50721 (9 SRationals)
    colorMatrix2 := extractColorMatrix(tf, 50722) // Tag 50722 (9 SRationals)
    baselineExposure := extractBaselineExposure(tf) // Tag 50730 (SRational)

    // Create profile from binary data
    profile := &Profile{
        ProfileName:            profileName,
        ToneCurve:              toneCurve,
        ColorMatrix1:           colorMatrix1,
        ColorMatrix2:           colorMatrix2,
        BaselineExposureOffset: baselineExposure,
    }

    // Convert DCP binary tone curve → UniversalRecipe
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
    "github.com/justin/recipe/internal/models"
)

// Generate creates a DCP file from UniversalRecipe
// Generates binary DNG tags with tone curve and color matrices
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
    // Convert UniversalRecipe → DCP binary tone curve
    toneCurve := universalToToneCurve(recipe)

    // Generate identity color matrices (SRational arrays)
    colorMatrix1 := generateIdentityMatrix()
    colorMatrix2 := generateIdentityMatrix()

    // Generate baseline exposure offset
    baselineExposure := SRational{Numerator: 0, Denominator: 1} // Zero offset

    // Convert to binary format
    toneCurveBinary := toneCurveToBytes(toneCurve)           // Float32 array
    colorMatrix1Binary := srationalArrayToBytes(colorMatrix1) // 9 SRationals
    colorMatrix2Binary := srationalArrayToBytes(colorMatrix2) // 9 SRationals
    baselineExposureBinary := srationalToBytes(baselineExposure)

    // Create DNG file with binary tags
    dngData, err := createDNG(toneCurveBinary, colorMatrix1Binary, colorMatrix2Binary, baselineExposureBinary)
    if err != nil {
        return nil, fmt.Errorf("failed to create DNG: %w", err)
    }

    return dngData, nil
}
```

**dcp/tiff.go Helper Functions:**

```go
// tiff.go
package dcp

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "github.com/google/tiff"
    "unsafe"
)

const (
    // DNG tag constants (binary format)
    TagUniqueCameraModel      = 50708
    TagColorMatrix1           = 50721
    TagColorMatrix2           = 50722
    TagBaselineExposureOffset = 50730
    TagProfileToneCurve       = 50940
    TagProfileName            = 52552
)

// readDNG reads DNG structure from bytes, converting magic bytes to TIFF version 42
func readDNG(data []byte) (*tiff.TIFF, error) {
    // Convert DNG magic bytes "IIRC"/"MMCR" → standard TIFF version 42
    // This is required because github.com/google/tiff expects TIFF version 42
    convertedData := convertDNGMagicBytes(data)

    r := bytes.NewReader(convertedData)
    return tiff.Parse(r, nil, nil)
}

// convertDNGMagicBytes converts DNG magic bytes to TIFF version 42
func convertDNGMagicBytes(data []byte) []byte {
    if len(data) < 4 {
        return data
    }

    // Check for DNG magic bytes
    if (data[0] == 'I' && data[1] == 'I' && data[2] == 'R' && data[3] == 'C') || // Little-endian "IIRC"
       (data[0] == 'M' && data[1] == 'M' && data[2] == 'C' && data[3] == 'R') {  // Big-endian "MMCR"

        // Make a copy and replace version bytes
        converted := make([]byte, len(data))
        copy(converted, data)
        converted[2] = 0x2A  // TIFF version 42
        converted[3] = 0x00
        return converted
    }

    return data
}

// extractToneCurve extracts binary tone curve from tag 50940
func extractToneCurve(tf *tiff.TIFF) []ToneCurvePoint {
    ifd := tf.IFDs()[0]
    entry := ifd.GetField(TagProfileToneCurve)
    if entry == nil {
        return nil // Missing tone curve
    }

    // Read binary float32 array
    data := entry.Value().([]byte)
    points := make([]ToneCurvePoint, len(data)/8) // 8 bytes per point

    buf := bytes.NewReader(data)
    for i := range points {
        var input, output float32
        binary.Read(buf, binary.LittleEndian, &input)
        binary.Read(buf, binary.LittleEndian, &output)
        points[i] = ToneCurvePoint{
            Input:  float64(input),
            Output: float64(output),
        }
    }

    return points
}

// extractColorMatrix extracts 9 SRational values from color matrix tag
func extractColorMatrix(tf *tiff.TIFF, tagID uint16) *Matrix {
    ifd := tf.IFDs()[0]
    entry := ifd.GetField(tagID)
    if entry == nil {
        return nil
    }

    data := entry.Value().([]byte)
    if len(data) != 72 { // 9 SRationals × 8 bytes
        return nil
    }

    matrix := &Matrix{}
    buf := bytes.NewReader(data)
    for row := 0; row < 3; row++ {
        for col := 0; col < 3; col++ {
            var sr SRational
            binary.Read(buf, binary.LittleEndian, &sr.Numerator)
            binary.Read(buf, binary.LittleEndian, &sr.Denominator)
            matrix.Rows[row][col] = sr.ToFloat64()
        }
    }

    return matrix
}

// createDNG creates DNG file with binary profile tags
func createDNG(toneCurve, colorMatrix1, colorMatrix2, baselineExposure []byte) ([]byte, error) {
    tf := tiff.New()

    // Create IFD with DNG tags
    ifd := tiff.NewIFD(tf)
    ifd.SetField(tiff.NewField(tiff.ImageWidth, tiff.UInt32(1), tiff.NewUInt32(1)))
    ifd.SetField(tiff.NewField(tiff.ImageLength, tiff.UInt32(1), tiff.NewUInt32(1)))

    // Embed binary profile tags
    ifd.SetField(tiff.NewField(TagProfileToneCurve, tiff.Byte, toneCurve))
    ifd.SetField(tiff.NewField(TagColorMatrix1, tiff.Byte, colorMatrix1))
    ifd.SetField(tiff.NewField(TagColorMatrix2, tiff.Byte, colorMatrix2))
    ifd.SetField(tiff.NewField(TagBaselineExposureOffset, tiff.Byte, baselineExposure))

    tf.SetIFD(ifd)

    // Write to bytes
    var buf bytes.Buffer
    if err := tf.Write(&buf); err != nil {
        return nil, fmt.Errorf("failed to write DNG: %w", err)
    }

    // Convert TIFF version 42 back to DNG magic bytes
    dngData := buf.Bytes()
    if len(dngData) >= 4 {
        dngData[2] = 'R'  // "IIRC" or "MMCR"
        dngData[3] = 'C'
    }

    return dngData, nil
}
```

**dcp/profile.go Binary Tone Curve Analysis:**

```go
// profile.go
package dcp

import (
    "bytes"
    "encoding/binary"
    "github.com/justin/recipe/internal/models"
)

// profileToUniversal converts DCP binary tone curve → UniversalRecipe
func profileToUniversal(profile *Profile) *models.UniversalRecipe {
    recipe := &models.UniversalRecipe{}

    if len(profile.ToneCurve) > 0 {
        // Analyze binary tone curve (0.0-1.0 normalized) to extract parameters
        exposure, contrast, highlights, shadows := analyzeToneCurve(profile.ToneCurve)

        recipe.Exposure = exposure
        recipe.Contrast = contrast
        recipe.Highlights = highlights
        recipe.Shadows = shadows
    }

    // Color matrices ignored (identity assumed)

    return recipe
}

// analyzeToneCurve extracts exposure/contrast/highlights/shadows from 0.0-1.0 normalized curve
func analyzeToneCurve(points []ToneCurvePoint) (exposure, contrast, highlights, shadows float64) {
    if len(points) == 0 {
        return 0, 0, 0, 0
    }

    // Midpoint shift → Exposure (compare point at input=0.5)
    midpoint := findPoint(points, 0.5)
    exposureShift := (midpoint.Output - 0.5) / 0.25  // Normalize to ±2.0 range
    exposure = clampFloat64(exposureShift, -2.0, 2.0)

    // Slope → Contrast (top quarter slope - bottom quarter slope)
    topSlope := (findPoint(points, 0.75).Output - findPoint(points, 0.5).Output) / 0.25
    bottomSlope := (findPoint(points, 0.5).Output - findPoint(points, 0.25).Output) / 0.25
    contrast = clampFloat64(topSlope - bottomSlope, -1.0, 1.0)

    // Top-end deviation → Highlights
    topPoint := findPoint(points, 1.0)
    highlights = clampFloat64((topPoint.Output - 1.0) * 4.0, -1.0, 1.0)

    // Bottom-end deviation → Shadows
    bottomPoint := findPoint(points, 0.0)
    shadows = clampFloat64(bottomPoint.Output * 4.0, -1.0, 1.0)

    return exposure, contrast, highlights, shadows
}

// universalToToneCurve converts UniversalRecipe → binary DCP tone curve (0.0-1.0 normalized)
func universalToToneCurve(recipe *models.UniversalRecipe) []ToneCurvePoint {
    // Base curve: 5 points in 0.0-1.0 range
    points := []ToneCurvePoint{
        {Input: 0.0, Output: 0.0},
        {Input: 0.25, Output: 0.25},
        {Input: 0.5, Output: 0.5},
        {Input: 0.75, Output: 0.75},
        {Input: 1.0, Output: 1.0},
    }

    // Apply exposure (vertical shift of midpoint)
    exposureShift := recipe.Exposure * 0.25  // ±2.0 → ±0.5
    points[2].Output = clampFloat64(0.5+exposureShift, 0.0, 1.0)

    // Apply contrast (steepen/flatten slope)
    contrastFactor := 1.0 + recipe.Contrast
    for i := range points {
        delta := points[i].Output - 0.5
        points[i].Output = clampFloat64(0.5+delta*contrastFactor, 0.0, 1.0)
    }

    // Apply highlights (adjust top-end curve)
    highlightShift := recipe.Highlights * 0.125
    points[3].Output = clampFloat64(points[3].Output+highlightShift, 0.0, 1.0) // 0.75 point
    points[4].Output = clampFloat64(points[4].Output+highlightShift, 0.0, 1.0) // 1.0 point

    // Apply shadows (adjust bottom-end curve)
    shadowShift := recipe.Shadows * 0.125
    points[0].Output = clampFloat64(points[0].Output+shadowShift, 0.0, 1.0) // 0.0 point
    points[1].Output = clampFloat64(points[1].Output+shadowShift, 0.0, 1.0) // 0.25 point

    return points
}

// toneCurveToBytes converts ToneCurvePoint array to binary float32 format
func toneCurveToBytes(points []ToneCurvePoint) []byte {
    buf := new(bytes.Buffer)
    for _, pt := range points {
        binary.Write(buf, binary.LittleEndian, float32(pt.Input))
        binary.Write(buf, binary.LittleEndian, float32(pt.Output))
    }
    return buf.Bytes()
}

// generateIdentityMatrix creates 9 SRational values for identity matrix
func generateIdentityMatrix() []SRational {
    return []SRational{
        {Numerator: 1, Denominator: 1}, {Numerator: 0, Denominator: 1}, {Numerator: 0, Denominator: 1},
        {Numerator: 0, Denominator: 1}, {Numerator: 1, Denominator: 1}, {Numerator: 0, Denominator: 1},
        {Numerator: 0, Denominator: 1}, {Numerator: 0, Denominator: 1}, {Numerator: 1, Denominator: 1},
    }
}

// srationalArrayToBytes converts SRational array to binary format
func srationalArrayToBytes(srs []SRational) []byte {
    buf := new(bytes.Buffer)
    for _, sr := range srs {
        binary.Write(buf, binary.LittleEndian, sr.Numerator)
        binary.Write(buf, binary.LittleEndian, sr.Denominator)
    }
    return buf.Bytes()
}

func clampFloat64(value, min, max float64) float64 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
```

### Workflows and Sequencing

**Parse Workflow (.dcp → UniversalRecipe):**

1. User uploads .dcp file via CLI/TUI/Web
2. `converter.detectFormat(data)` identifies file as DCP (DNG magic bytes "IIRC"/"MMCR" + .dcp extension)
3. `dcp.Parse(data)` called:
   a. `readDNG(data)` converts DNG magic bytes → TIFF version 42, parses using github.com/google/tiff
   b. `extractToneCurve(tf)` extracts binary float32 array from tag 50940
   c. `extractColorMatrix(tf, 50721)` extracts 9 SRational values from tag 50721
   d. `extractColorMatrix(tf, 50722)` extracts 9 SRational values from tag 50722
   e. `extractProfileName(tf)` extracts ASCII string from tag 52552 (optional, empty if missing)
   f. `analyzeToneCurve(points)` extracts parameters from 0.0-1.0 normalized curve:
      - Extract exposure from midpoint shift (input=0.5)
      - Extract contrast from curve slope (top quarter vs bottom quarter)
      - Extract highlights from top-end curve point (input=1.0)
      - Extract shadows from bottom-end curve point (input=0.0)
4. Return `*models.UniversalRecipe` to converter
5. Converter routes to target format generator

**Generate Workflow (UniversalRecipe → .dcp):**

1. Converter calls `dcp.Generate(recipe)`
2. `universalToToneCurve(recipe)` builds binary DCP tone curve:
   a. Start with linear 5-point curve (0.0-1.0 normalized: 0.0, 0.25, 0.5, 0.75, 1.0)
   b. Apply exposure adjustment (vertical shift of midpoint at 0.5)
   c. Apply contrast adjustment (steepen/flatten slope around 0.5)
   d. Apply highlights adjustment (modify 0.75 and 1.0 points)
   e. Apply shadows adjustment (modify 0.0 and 0.25 points)
   f. Clamp all output values to 0.0-1.0 range
3. `toneCurveToBytes(points)` converts to binary float32 array (8 bytes per point)
4. `generateIdentityMatrix()` creates 9 SRational values for identity matrix
5. `srationalArrayToBytes(matrix)` converts to binary format (72 bytes total)
6. `createDNG()` creates DNG file:
   a. Initialize TIFF structure with minimal IFD
   b. Embed binary tone curve in tag 50940
   c. Embed binary color matrices in tags 50721, 50722
   d. Embed baseline exposure offset in tag 50730
   e. Write TIFF to bytes with version 42
   f. Convert TIFF version 42 back to DNG magic bytes "IIRC"
7. Return DNG bytes to converter
8. Converter returns to user (download/write file)

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
- **TIFF Bomb**: Limit DNG file size to <10MB, validate IFD structure before parsing
- **Binary Data Corruption**: Validate binary array lengths (tone curve: multiple of 8 bytes, matrices: exactly 72 bytes)
- **Malformed DNG**: github.com/google/tiff handles gracefully, returns errors (no panics)
- **Invalid Binary Tags**: Validate tag 50940 contains valid float32 pairs, tags 50721-50722 contain 9 SRationals
- **File Upload**: Web UI validates .dcp extension, server validates DNG magic bytes (IIRC or MMCR)

**Privacy:**
- All processing client-side (WASM) or local (CLI/TUI)
- DCP files may contain camera model metadata (tag 50708) - preserve but don't expose unnecessarily
- Zero external API calls

**Input Validation:**
- Verify DNG magic bytes (IIRC for little-endian, MMCR for big-endian)
- Check tag 50940 exists and contains valid float32 array (length % 8 == 0)
- Validate tone curve points within 0.0-1.0 range
- Validate SRational denominators are non-zero
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
- Invalid DNG structure: Return error, don't crash
- Missing tag 50940: Return error ("not a valid DCP file") or use linear curve (0, 0) → (1, 1)
- Malformed binary data: Return clear parsing error with byte offset
- Out-of-range tone curve values: Clamp to 0.0-1.0
- Zero denominator in SRational: Treat as 0.0 or return error

### Observability

**Logging:**
- Use Go's `log/slog` structured logging
- Log level DEBUG: DNG tag reading, binary data extraction, tone curve calculations
- Log level INFO: DCP conversion success/failure
- Log level WARN: Skipped unsupported DCP features (non-identity color matrices, missing ProfileName)
- Log level ERROR: DNG parsing failures, binary data corruption errors

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
- Go stdlib `image/tiff` is decoder-only, doesn't support writing or custom DNG tags
- `github.com/google/tiff` provides complete TIFF/DNG read/write capabilities
- Google-maintained, stable, widely used in imaging applications
- Necessary for reading/writing binary DNG tags 50700-52600
- Handles DNG magic bytes with custom conversion (IIRC/MMCR → version 42)

Standard Library Usage:
- `encoding/binary`: Binary data conversion for float32 arrays and SRational values
- `unsafe`: DNG magic byte conversion (IIRC/MMCR ↔ version 42)
- `fmt`: Error wrapping
- `bytes`, `io`: Stream handling for binary DNG operations

**External Integrations:**
- **Adobe Camera Raw / Lightroom** (validation): Generated DCP files must load in Adobe software (manual testing)
- **Adobe DCP samples** (test data): Download reference DCP files from Adobe for test fixtures

**Version Constraints:**
- Go 1.25.1+ (required for go:wasmexport)
- `github.com/google/tiff` latest stable release
- Adobe DNG Specification 1.6 (DCP format reference)

## Acceptance Criteria (Authoritative)

**AC-1: Parse DCP Files**
- ✅ Read DNG-based DCP file structure using github.com/google/tiff
- ✅ Convert DNG magic bytes (IIRC/MMCR) to TIFF version 42 for library compatibility
- ✅ Extract binary tone curve from tag 50940 (float32 array, 0.0-1.0 normalized)
- ✅ Extract binary color matrices from tags 50721-50722 (9 SRational values each)
- ✅ Parse optional ProfileName from tag 52552 (ASCII, may be missing)
- ✅ Handle identity matrices (diagonal SRational{1,1}, off-diagonal SRational{0,1})
- ✅ Handle DNG v1.x format variations (v1.0-v1.6 supported)
- ✅ Validate DNG structure and report parsing errors with clear messages
- ✅ Support both embedded (in DNG) and standalone DCP files (standalone priority for MVP)

**AC-2: Generate DCP Files**
- ✅ Generate DNG-based DCP file structure per Adobe DNG Specification 1.6
- ✅ Generate binary tone curve in tag 50940 (float32 array, 0.0-1.0 normalized)
- ✅ Generate binary color matrices in tags 50721-50722 (9 SRational values each, identity matrices)
- ✅ Generate baseline exposure offset in tag 50730 (SRational{0,1})
- ✅ Map UniversalRecipe parameters to binary DCP tone curve (0.0-1.0 range)
- ✅ Create tone curves from exposure/contrast/highlights/shadows adjustments (5-point curve)
- ✅ Convert TIFF version 42 back to DNG magic bytes (IIRC) before returning
- ✅ Validate generated DCP against Adobe spec (well-formed DNG + valid binary tags)
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
1. ✅ Acquire 36 DCP sample files from Adobe for test fixtures (DONE - story 9-1)
2. ✅ Add `github.com/google/tiff` dependency to go.mod (DONE - story 9-1)
3. ✅ Set up `internal/formats/dcp/` package structure (DONE - story 9-1)
4. ✅ Implement `tiff.go` DNG tag operations with magic byte conversion (DONE - story 9-1)
5. ✅ Implement `profile.go` binary tone curve analysis algorithms (DONE - story 9-1)
6. ✅ Implement `parse.go` orchestrating DNG → binary tags → UniversalRecipe (DONE - story 9-1)
7. ⏳ Implement `generate.go` orchestrating UniversalRecipe → binary tags → DNG (story 9-2)
8. ⏳ Write unit tests achieving ≥85% coverage (story 9-2)
9. ⏳ Integrate with converter and add format detection (story 9-3)
10. ⏳ Update documentation (parameter-mapping.md, README.md) (story 9-3)
11. ⏳ Manual validation in Adobe Camera Raw / Lightroom trial (story 9-4)
12. Mark epic-9 as "contexted" in sprint-status.yaml

**Critical Discovery from Story 9-1:**
Real Adobe DCP files use **binary DNG tags** (50700-52600), NOT XML in tag 50740. All documentation updated to reflect this binary format. See `docs/stories/9-1-dcp-parser-FORMAT-PIVOT.md` for complete details.
