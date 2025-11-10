# Story 9.2: DNG Camera Profile (DCP) Generator

Status: ready-for-dev

> ⚠️ **CRITICAL FORMAT DISCOVERY**: Real Adobe DCP files use **binary TIFF tags** (50700-52600), NOT XML in tag 50740. This story was updated after implementing story 9-1 revealed the actual DCP format. See `docs/stories/9-1-dcp-parser-FORMAT-PIVOT.md` for details.

## Story

As a **photographer**,
I want **Recipe to generate valid DNG Camera Profile (.dcp) files from UniversalRecipe representation**,
so that **I can convert my presets from other formats (NP3, XMP, lrtemplate, .costyle) to Adobe camera profiles and use them in Lightroom, Camera Raw, and with DNG files**.

## Acceptance Criteria

**AC-1: Generate Valid DNG/TIFF Structure**
- ✅ Create TIFF/DNG file with proper header (II or MM byte order marker)
- ✅ Use DNG version magic bytes ("IIRC" for little-endian or "MMCR" for big-endian)
- ✅ Generate Image File Directory (IFD) with required binary tags
- ✅ Include standard TIFF tags (ImageWidth, ImageLength, SamplesPerPixel, etc.)
- ✅ Write TIFF/DNG file using `github.com/google/tiff` library or manual binary writing
- ✅ Validate DNG structure (well-formed, valid byte order, correct version bytes)
- ✅ Generated DCP files open without errors in TIFF/DNG readers

**AC-2: Generate Binary Profile Data Tags**
- ✅ Write tag 50708 (UniqueCameraModel) as ASCII string
- ✅ Write tag 50940 (ProfileToneCurve) as binary float32 array (input/output pairs)
- ✅ Write tags 50721-50722 (ColorMatrix1/2) as 9 SRational values each
- ✅ Write tag 50730 (BaselineExposureOffset) as SRational
- ✅ Write tag 52552 (ProfileName) as ASCII string (OPTIONAL)
- ✅ Use correct binary data types (Float32, SRational, ASCII)
- ✅ All tag data follows DNG 1.0-1.6 binary format

**AC-3: Map UniversalRecipe to Binary Tone Curve**
- ✅ Generate tone curve as array of float32 (input, output) pairs:
  - Start with linear curve: {0.0, 0.0}, {0.25, 0.25}, {0.5, 0.5}, {0.75, 0.75}, {1.0, 1.0}
  - Apply exposure: Vertical shift of midpoint (0.5 → 0.5 + exposure*0.25)
  - Apply contrast: Steepen/flatten curve slope
  - Apply highlights: Adjust top-end points (0.75-1.0 range)
  - Apply shadows: Adjust bottom-end points (0.0-0.25 range)
- ✅ Clamp all curve points to valid 0.0-1.0 range
- ✅ Ensure monotonic curve (output[i] >= output[i-1])
- ✅ Write as binary float32 pairs (8 bytes per point: 4 bytes input + 4 bytes output)

**AC-4: Generate Binary Color Matrices**
- ✅ Create ColorMatrix1 as 9 SRational values (identity matrix):
  ```
  SRational{1, 1}, SRational{0, 1}, SRational{0, 1},
  SRational{0, 1}, SRational{1, 1}, SRational{0, 1},
  SRational{0, 1}, SRational{0, 1}, SRational{1, 1}
  ```
- ✅ Create ColorMatrix2 with same identity values (dual illuminant requirement)
- ✅ Write as binary SRational (8 bytes each: 4 bytes numerator + 4 bytes denominator)
- ✅ Skip full camera calibration (ForwardMatrix, CalibrationIlluminant optional)

**AC-5: Validate Generated DCP**
- ✅ Validate DNG structure (magic bytes "IIRC"/"MMCR", IFD, required tags exist)
- ✅ Validate binary tag data (correct lengths, valid SRational denominators)
- ✅ Validate tone curve points (0.0-1.0 range, monotonic, float32 format)
- ✅ Generated DCP loads in Adobe Camera Raw without errors (manual validation)
- ✅ Generated DCP loads in Lightroom Classic without errors (manual validation)
- ✅ Tone adjustments render correctly (visual spot-check)

**AC-6: Unit Test Coverage**
- ✅ Unit tests for Generate() function with various UniversalRecipe inputs
- ✅ Test edge cases (empty recipe, extreme values, minimal parameters)
- ✅ Test DNG structure correctness (validate IFD, tags, byte order, DNG version)
- ✅ Test binary tag generation (validate float32 arrays, SRational values)
- ✅ Test coverage ≥85% for dcp/generate.go, dcp/tiff.go, dcp/profile.go
- ✅ All tests pass in CI

## Tasks / Subtasks

### Task 1: Implement Tone Curve Generation (AC-3)
- [ ] Implement `universalToToneCurve()` function in `profile.go`:
  ```go
  func universalToToneCurve(recipe *universal.Recipe) []ToneCurvePoint {
      // Start with linear 5-point curve (0.0-1.0 normalized)
      points := []ToneCurvePoint{
          {Input: 0.0, Output: 0.0},
          {Input: 0.25, Output: 0.25},
          {Input: 0.5, Output: 0.5},
          {Input: 0.75, Output: 0.75},
          {Input: 1.0, Output: 1.0},
      }

      // Apply exposure (vertical shift of midpoint)
      exposureShift := recipe.Exposure * 0.25
      points[2].Output = clampFloat64(0.5 + exposureShift, 0.0, 1.0)

      // Apply contrast (steepen/flatten slope)
      contrastFactor := 1.0 + recipe.Contrast
      for i := range points {
          deviation := points[i].Input - 0.5
          points[i].Output = clampFloat64(0.5 + deviation*contrastFactor, 0.0, 1.0)
      }

      // Apply highlights (adjust top-end points)
      highlightsShift := recipe.Highlights * 0.125
      points[3].Output = clampFloat64(points[3].Output + highlightsShift, points[2].Output, 1.0)
      points[4].Output = clampFloat64(points[4].Output + highlightsShift, points[3].Output, 1.0)

      // Apply shadows (adjust bottom-end points)
      shadowsShift := recipe.Shadows * 0.125
      points[0].Output = clampFloat64(points[0].Output + shadowsShift, 0.0, points[1].Output)
      points[1].Output = clampFloat64(points[1].Output + shadowsShift, points[0].Output, points[2].Output)

      return points
  }
  ```
- [ ] Ensure monotonic curve (output[i] >= output[i-1])
- [ ] Test with various parameter combinations

### Task 2: Generate Binary Tag Data (AC-2, AC-4)
- [ ] Implement `generateBinaryToneCurve()` function in `profile.go`:
  ```go
  func generateBinaryToneCurve(points []ToneCurvePoint) []byte {
      // Each point is 2 float32 values (8 bytes total)
      buf := new(bytes.Buffer)
      for _, pt := range points {
          binary.Write(buf, binary.LittleEndian, float32(pt.Input))
          binary.Write(buf, binary.LittleEndian, float32(pt.Output))
      }
      return buf.Bytes()
  }
  ```
- [ ] Implement `generateColorMatrix()` function (returns 9 SRational values):
  ```go
  func generateColorMatrix() []SRational {
      // Identity matrix: diagonal 1.0, off-diagonal 0.0
      return []SRational{
          {Numerator: 1, Denominator: 1}, {Numerator: 0, Denominator: 1}, {Numerator: 0, Denominator: 1},
          {Numerator: 0, Denominator: 1}, {Numerator: 1, Denominator: 1}, {Numerator: 0, Denominator: 1},
          {Numerator: 0, Denominator: 1}, {Numerator: 0, Denominator: 1}, {Numerator: 1, Denominator: 1},
      }
  }
  ```
- [ ] Implement `srationalToBytes()` helper to convert SRational array to binary
- [ ] Validate binary data lengths (tone curve: N*8 bytes, matrix: 72 bytes)

### Task 3: Create DNG File with Binary Tags (AC-1, AC-2)
- [ ] Implement `createDNG()` function in `tiff.go`:
  ```go
  func createDNG(toneCurve []byte, colorMatrix1 []byte, colorMatrix2 []byte, profileName string) ([]byte, error) {
      // Create buffer with DNG magic bytes
      buf := new(bytes.Buffer)

      // Write DNG header (little-endian)
      buf.Write([]byte{0x49, 0x49})       // "II" (little-endian)
      buf.Write([]byte{0x52, 0x43})       // "RC" (DNG version instead of 42)

      // Build IFD with binary tags
      ifd := buildIFD()

      // Add DNG profile tags
      ifd.SetTag(TagUniqueCameraModel, "Recipe Converted Camera")
      ifd.SetTag(TagProfileToneCurve, toneCurve)        // Float32 array
      ifd.SetTag(TagColorMatrix1, colorMatrix1)         // SRational array
      ifd.SetTag(TagColorMatrix2, colorMatrix2)         // SRational array
      ifd.SetTag(TagBaselineExposureOffset, []byte{0, 0, 0, 0, 1, 0, 0, 0}) // SRational{0, 1}
      if profileName != "" {
          ifd.SetTag(TagProfileName, profileName)       // ASCII string (OPTIONAL)
      }

      // Write IFD to buffer
      if err := ifd.Write(buf); err != nil {
          return nil, fmt.Errorf("failed to write IFD: %w", err)
      }

      return buf.Bytes(), nil
  }
  ```
- [ ] Implement `buildIFD()` helper to create minimal IFD with standard TIFF tags
- [ ] Verify DNG magic bytes ("IIRC" or "MMCR") in output
- [ ] Validate IFD structure with binary tag data

### Task 4: Implement Generate() Function (AC-1 to AC-4)
- [ ] Implement `Generate(recipe *universal.Recipe) ([]byte, error)` in `generate.go`:
  ```go
  func Generate(recipe *universal.Recipe) ([]byte, error) {
      // Step 1: Generate tone curve points
      points := universalToToneCurve(recipe)

      // Step 2: Convert to binary float32 array
      toneCurve := generateBinaryToneCurve(points)

      // Step 3: Generate identity color matrices
      matrix1 := generateColorMatrix()
      matrix2 := generateColorMatrix()
      colorMatrix1Bytes := srationalArrayToBytes(matrix1)
      colorMatrix2Bytes := srationalArrayToBytes(matrix2)

      // Step 4: Get profile name from metadata (if exists)
      profileName := ""
      if name, ok := recipe.Metadata["profile_name"]; ok {
          profileName = name
      }

      // Step 5: Create DNG with binary tags
      dngData, err := createDNG(toneCurve, colorMatrix1Bytes, colorMatrix2Bytes, profileName)
      if err != nil {
          return nil, err
      }

      return dngData, nil
  }
  ```
- [ ] Add error handling for nil recipe input
- [ ] Validate recipe parameters before generation

### Task 5: Write Unit Tests (AC-6)
- [ ] Write `TestGenerate_ValidRecipe()` - Generate DCP from populated UniversalRecipe
- [ ] Write `TestGenerate_EmptyRecipe()` - Generate from neutral recipe (all zeros)
- [ ] Write `TestGenerate_ExtremeValues()` - Test with extreme parameters (exposure=+2.0, contrast=+1.0)
- [ ] Write `TestUniversalToToneCurve()` - Test tone curve generation formulas with 0.0-1.0 values
- [ ] Write `TestGenerateColorMatrix()` - Verify identity matrix generation as SRational
- [ ] Write `TestDNGStructure()` - Validate generated DNG structure (magic bytes "IIRC", IFD, binary tags)
- [ ] Write `TestBinaryTagData()` - Validate binary tag data (float32 arrays, SRational values)
- [ ] Write `TestRoundTrip_DCP()` - Generate → Parse → Compare (verify 95%+ accuracy)
- [ ] Run tests: `go test ./internal/formats/dcp/`
- [ ] Verify coverage: `go test -cover ./internal/formats/dcp/` (target ≥85%)

### Task 6: Manual Validation in Adobe Software (AC-5)
- [ ] Generate 3 test DCP files from UniversalRecipe:
  - Neutral preset (all parameters zero)
  - Portrait preset (exposure +0.5, contrast +0.3, highlights -0.2)
  - Landscape preset (exposure +0.3, saturation +0.4, shadows +0.2)
- [ ] Test in Adobe Camera Raw:
  - Open Camera Raw (Photoshop or standalone)
  - Import generated DCP profiles
  - Apply to test image
  - Verify no errors during import/application
- [ ] Test in Lightroom Classic:
  - Import generated DCP profiles into Lightroom
  - Apply to test image in Develop module
  - Verify tone adjustments visible (exposure, contrast, highlights, shadows)
- [ ] Document validation results in `testdata/dcp/validation-report.md`

### Task 7: Documentation (AC-1 to AC-5)
- [ ] Add function comment for `Generate()`:
  - Document input (UniversalRecipe), output (DCP TIFF/DNG bytes), error cases
  - Include example usage
- [ ] Update `docs/parameter-mapping.md` with DCP generation mappings:
  - Document tone curve generation formulas (exposure/contrast/highlights/shadows → float32 pairs)
  - Note precision considerations (float64 → float32 curve points)
  - Provide examples with visual curve diagrams (optional)
- [ ] Add README notes in `testdata/dcp/`:
  - Document generated DCP validation results
  - Note Adobe software compatibility (Camera Raw, Lightroom versions tested)
  - List known limitations (no dual illuminant, no HSV tables, identity matrices only)

## Dev Notes

### Critical Format Discovery

**From Story 9-1-dcp-parser (Status: ✅ COMPLETED)**

Real Adobe DCP files use **binary TIFF/DNG tags** (50700-52600 range), NOT XML in tag 50740. This required a complete rewrite of story 9-1 and impacts all generation logic in story 9-2.

**Binary DNG Tag Format:**

| Tag ID | Name | Type | Binary Format |
|--------|------|------|---------------|
| **50708** | UniqueCameraModel | ASCII | Null-terminated string |
| **50721** | ColorMatrix1 | SRational[9] | 72 bytes (9 × 8-byte SRational) |
| **50722** | ColorMatrix2 | SRational[9] | 72 bytes (9 × 8-byte SRational) |
| **50730** | BaselineExposureOffset | SRational | 8 bytes (numerator + denominator) |
| **50940** | ProfileToneCurve | Float32[] | N × 8 bytes (N input/output pairs) |
| **52552** | ProfileName | ASCII | Null-terminated string (**OPTIONAL**) |

**Data Type Details:**
- **SRational**: 8 bytes = 4-byte signed int32 numerator + 4-byte signed int32 denominator
- **Float32**: 4 bytes = IEEE 754 single-precision float
- **Tone Curve Point**: 8 bytes = 4-byte float32 input + 4-byte float32 output
- **Normalization**: All tone curve values are 0.0-1.0 normalized (NOT 0-255 integers)

**DNG Version Bytes:**
- Standard TIFF uses version 42 (0x002A)
- DNG uses "IIRC" (0x49494352) for little-endian or "MMCR" (0x4D4D4352) for big-endian
- When using `github.com/google/tiff` library, may need to convert DNG version to TIFF version 42

[Source: docs/stories/9-1-dcp-parser-FORMAT-PIVOT.md]

### Architecture Alignment

**Tech Spec Epic 9 Alignment:**

Story 9-2 implements **AC-2 (Generate DCP Files)** from tech-spec-epic-9.md.

**Generation Flow (Updated for Binary Format):**
```
UniversalRecipe → universalToToneCurve() → generateBinaryToneCurve() → createDNG() → .dcp bytes
```

**Binary Tone Curve Generation (5-point curve, 0.0-1.0 normalized):**
```
Linear base:   {0.0, 0.0}  {0.25, 0.25}  {0.5, 0.5}  {0.75, 0.75}  {1.0, 1.0}
                   ↓            ↓            ↓            ↓            ↓
Apply exposure:              shifts midpoint (0.5 → 0.5+shift)
Apply contrast:              steepens/flattens slope (multiply deviation from 0.5)
Apply highlights:            adjusts top-end points (0.75-1.0)
Apply shadows:               adjusts bottom-end points (0.0-0.25)
                   ↓            ↓            ↓            ↓            ↓
Output curve:   {0.0, X}    {0.25, Y}    {0.5, Z}    {0.75, W}    {1.0, V}
                   ↓            ↓            ↓            ↓            ↓
Binary format:  [float32 0.0, float32 X] [float32 0.25, float32 Y] ... (40 bytes total)
```

**Identity Matrix Binary Format (SRational[9]):**
```
ColorMatrix1 = ColorMatrix2 = [
    SRational{1, 1}  SRational{0, 1}  SRational{0, 1}
    SRational{0, 1}  SRational{1, 1}  SRational{0, 1}
    SRational{0, 1}  SRational{0, 1}  SRational{1, 1}
]

Binary: 72 bytes = 9 SRationals × 8 bytes each
  [1,0,0,0, 1,0,0,0]  [0,0,0,0, 1,0,0,0]  [0,0,0,0, 1,0,0,0]
  [0,0,0,0, 1,0,0,0]  [1,0,0,0, 1,0,0,0]  [0,0,0,0, 1,0,0,0]
  [0,0,0,0, 1,0,0,0]  [0,0,0,0, 1,0,0,0]  [1,0,0,0, 1,0,0,0]
```

[Source: docs/tech-spec-epic-9.md#Detailed-Design - requires update]

### DNG/TIFF Writing Pattern

**Binary Tag Writing Approach:**

Since `github.com/google/tiff` may not support writing custom DNG tags easily, consider manual binary writing:

```go
import (
    "bytes"
    "encoding/binary"
)

func createDNG(toneCurve, colorMatrix1, colorMatrix2 []byte, profileName string) ([]byte, error) {
    buf := new(bytes.Buffer)

    // Write DNG header (little-endian)
    buf.Write([]byte{0x49, 0x49})       // "II" (little-endian byte order)
    buf.Write([]byte{0x52, 0x43})       // "RC" (DNG version instead of TIFF 42)

    // Write IFD offset (4 bytes, offset to first IFD)
    binary.Write(buf, binary.LittleEndian, uint32(8))

    // Build IFD with all tags
    ifd := buildIFD()
    ifd.AddTag(TagProfileToneCurve, TypeFloat, uint32(len(toneCurve)/4), toneCurve)
    ifd.AddTag(TagColorMatrix1, TypeSRational, 9, colorMatrix1)
    ifd.AddTag(TagColorMatrix2, TypeSRational, 9, colorMatrix2)
    // ... add more tags

    // Write IFD to buffer
    ifd.Write(buf)

    return buf.Bytes(), nil
}
```

**TIFF Tag Structure (12 bytes per tag):**
```
Bytes 0-1:   Tag ID (uint16)
Bytes 2-3:   Data Type (uint16) - 1=Byte, 2=ASCII, 3=Short, 5=Rational, 10=SRational, 11=Float
Bytes 4-7:   Count (uint32) - number of values
Bytes 8-11:  Value/Offset (uint32) - if ≤4 bytes: value, else: offset to data
```

**Alternative: Use github.com/google/tiff for Reading, Manual Writing:**
```go
// Option 1: Extend google/tiff library
tiffFile := tiff.New()
tiffFile.SetField(TagProfileToneCurve, toneCurve)

// Option 2: Manual binary writing (more control)
// See above approach
```

[Source: Adobe DNG Specification 1.6, TIFF 6.0 Specification]

### Binary Data Conversion Helpers

**Float32 to Bytes:**
```go
func float32ToBytes(f float32) []byte {
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, f)
    return buf.Bytes()
}
```

**SRational to Bytes:**
```go
type SRational struct {
    Numerator   int32
    Denominator int32
}

func srationalToBytes(sr SRational) []byte {
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, sr.Numerator)
    binary.Write(buf, binary.LittleEndian, sr.Denominator)
    return buf.Bytes()
}

func srationalArrayToBytes(srs []SRational) []byte {
    buf := new(bytes.Buffer)
    for _, sr := range srs {
        binary.Write(buf, binary.LittleEndian, sr.Numerator)
        binary.Write(buf, binary.LittleEndian, sr.Denominator)
    }
    return buf.Bytes()
}
```

[Source: internal/formats/dcp/tiff.go - parser uses similar conversion]

### Project Structure Notes

**New Files Created (Story 9-2):**
```
internal/formats/dcp/
├── generate.go           # Generate() function (NEW)
├── generate_test.go      # Unit tests for generation (NEW)
└── testdata/dcp/
    └── validation-report.md  # Adobe validation results (NEW)
```

**Modified Files:**
- `internal/formats/dcp/tiff.go` - Add DNG writing functions (binary tag writing)
- `internal/formats/dcp/profile.go` - Add binary tone curve generation

**Files from Story 9-1 (Reused):**
- `types.go` - Struct definitions (ToneCurvePoint, Matrix, SRational - updated for binary format)
- `testdata/dcp/` - Sample files (36 real Adobe DCP files for round-trip testing)

[Source: docs/tech-spec-epic-9.md#Components - requires update]

### Testing Strategy

**Unit Tests (Required for AC-6):**
- `TestGenerate_ValidRecipe()` - Generate from populated UniversalRecipe
- `TestGenerate_EmptyRecipe()` - Generate neutral preset
- `TestGenerate_ExtremeValues()` - Test with extreme parameters
- `TestUniversalToToneCurve()` - Verify tone curve formulas (0.0-1.0 range)
- `TestGenerateColorMatrix()` - Verify SRational identity matrix
- `TestDNGStructure()` - Validate DNG magic bytes, IFD, binary tags
- `TestBinaryTagData()` - Validate float32 and SRational binary formats
- `TestRoundTrip_DCP()` - Generate → Parse → Compare
- Coverage target: ≥85% for generate.go

**Manual Validation (Required for AC-5):**
- Load generated DCP in Adobe Camera Raw (no errors)
- Load generated DCP in Lightroom Classic (no errors)
- Visual spot-check (tone adjustments render correctly)
- Document results in validation-report.md

[Source: docs/tech-spec-epic-9.md#Test-Strategy-Summary - requires update]

### Known Risks

**RISK-16: Generated DCPs rejected by Adobe software**
- **Impact**: Lightroom/Camera Raw refuse to load generated profiles
- **Mitigation**: Follow Adobe DNG Specification 1.6 exactly, use binary format (not XML), validate with real Adobe samples
- **Fallback**: Adjust DNG structure based on validation errors

**RISK-17: Tone curve formula inaccuracy**
- **Impact**: Generated curve doesn't match visual expectations
- **Mitigation**: Use 0.0-1.0 normalized values (not 0-255 integers), visual validation in Lightroom, iterate on formulas
- **Target**: 90%+ visual similarity to original preset

**RISK-18: Binary format compatibility**
- **Impact**: Adobe software expects specific binary formats (float32, SRational)
- **Mitigation**: Test binary conversions carefully, validate byte order (little-endian), use exact data types from DNG spec
- **Fallback**: Inspect real Adobe DCP files to understand exact binary format

**RISK-19: Optional ProfileName handling**
- **Impact**: Some DCP files don't have tag 52552 (ProfileName)
- **Mitigation**: Make ProfileName optional (don't write tag 52552 if empty string)
- **Testing**: Test both with and without ProfileName

[Source: docs/tech-spec-epic-9.md#Risks-Assumptions-Open-Questions - requires update]

### References

- [Source: docs/stories/9-1-dcp-parser-FORMAT-PIVOT.md] - Binary format discovery
- [Source: docs/tech-spec-epic-9.md#Acceptance-Criteria] - AC-2: Generate DCP Files (requires update)
- [Source: docs/tech-spec-epic-9.md#Data-Models-and-Contracts] - Tone curve formulas (requires update)
- [Source: docs/tech-spec-epic-9.md#APIs-and-Interfaces] - Generate() function signature
- [Source: Adobe DNG Specification 1.6] - Binary tag definitions, DNG version bytes
- [Source: TIFF 6.0 Specification] - TIFF IFD structure, tag format
- [Source: internal/formats/dcp/parse.go] - Parse() function (reverse operation, binary format)
- [Source: internal/formats/dcp/tiff.go] - Binary tag extraction (reference for generation)

## Dev Agent Record

### Context Reference

- Story Context XML: `docs/stories/9-2-dcp-generator.context.xml` (Generated: 2025-11-09, **requires update for binary format**)

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
