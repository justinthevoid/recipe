# Story 9.2: DNG Camera Profile (DCP) Generator

Status: ready-for-dev

## Story

As a **photographer**,
I want **Recipe to generate valid DNG Camera Profile (.dcp) files from UniversalRecipe representation**,
so that **I can convert my presets from other formats (NP3, XMP, lrtemplate, .costyle) to Adobe camera profiles and use them in Lightroom, Camera Raw, and with DNG files**.

## Acceptance Criteria

**AC-1: Generate Valid TIFF Structure**
- ✅ Create TIFF file with proper header (II or MM byte order marker)
- ✅ Generate Image File Directory (IFD) with required tags
- ✅ Include standard TIFF tags (ImageWidth, ImageLength, SamplesPerPixel, etc.)
- ✅ Write TIFF file using `github.com/google/tiff` library
- ✅ Validate TIFF structure (well-formed, valid byte order)
- ✅ Generated TIFF files open without errors in TIFF readers

**AC-2: Embed Camera Profile XML in Tag 50740**
- ✅ Generate Adobe Camera Profile XML from UniversalRecipe
- ✅ Embed XML data in TIFF tag 50740 (CameraProfile)
- ✅ Include required XML elements (ProfileName, ToneCurve, ColorMatrix)
- ✅ Use Adobe Camera Raw namespace (`http://ns.adobe.com/camera-raw-settings/1.0/`)
- ✅ Format XML with proper RDF structure (rdf:Seq, rdf:li elements)
- ✅ XML is human-readable (formatted with indentation)

**AC-3: Map UniversalRecipe to DCP Tone Curve**
- ✅ Generate 5-point tone curve from UniversalRecipe adjustments:
  - Start with linear curve (0,0 → 64,64 → 128,128 → 192,192 → 255,255)
  - Apply exposure: Vertical shift of midpoint (128 → 128+exposure*64)
  - Apply contrast: Steepen/flatten curve slope
  - Apply highlights: Adjust top-end points (192-255 range)
  - Apply shadows: Adjust bottom-end points (0-64 range)
- ✅ Clamp all curve points to valid 0-255 range
- ✅ Ensure monotonic curve (output[i] >= output[i-1])

**AC-4: Generate Identity Color Matrices**
- ✅ Create ColorMatrix1 with identity values (3x3 diagonal matrix):
  ```
  1.0 0.0 0.0
  0.0 1.0 0.0
  0.0 0.0 1.0
  ```
- ✅ Create ColorMatrix2 with same identity values (dual illuminant requirement)
- ✅ Format matrices as RDF lists (rdf:Seq with rdf:li elements)
- ✅ Include ProfileCalibrationSignature = "com.adobe" (standard signature)
- ✅ Skip full camera calibration (ForwardMatrix, CalibrationIlluminant optional)

**AC-5: Validate Generated DCP**
- ✅ Validate TIFF structure (magic bytes, IFD, tag 50740 exists)
- ✅ Validate XML structure (well-formed, Adobe namespace correct)
- ✅ Validate tone curve points (0-255 range, monotonic)
- ✅ Generated DCP loads in Adobe Camera Raw without errors (manual validation)
- ✅ Generated DCP loads in Lightroom Classic without errors (manual validation)
- ✅ Tone adjustments render correctly (visual spot-check)

**AC-6: Unit Test Coverage**
- ✅ Unit tests for Generate() function with various UniversalRecipe inputs
- ✅ Test edge cases (empty recipe, extreme values, minimal parameters)
- ✅ Test TIFF structure correctness (validate IFD, tags, byte order)
- ✅ Test XML generation (validate namespace, elements, formatting)
- ✅ Test coverage ≥85% for dcp/generate.go, dcp/tiff.go, dcp/profile.go
- ✅ All tests pass in CI

## Tasks / Subtasks

### Task 1: Implement Tone Curve Generation (AC-3)
- [ ] Implement `universalToToneCurve()` function in `profile.go`:
  ```go
  func universalToToneCurve(recipe *universal.Recipe) []ToneCurvePoint {
      // Start with linear 5-point curve
      points := []ToneCurvePoint{
          {Input: 0, Output: 0},
          {Input: 64, Output: 64},
          {Input: 128, Output: 128},
          {Input: 192, Output: 192},
          {Input: 255, Output: 255},
      }

      // Apply exposure (vertical shift of midpoint)
      exposureShift := int(recipe.Exposure * 64.0)
      points[2].Output = clamp(128 + exposureShift, 0, 255)

      // Apply contrast (steepen/flatten slope)
      contrastFactor := 1.0 + recipe.Contrast
      for i := range points {
          deviation := points[i].Input - 128
          points[i].Output = clamp(128 + int(float64(deviation)*contrastFactor), 0, 255)
      }

      // Apply highlights (adjust top-end points)
      highlightsShift := int(recipe.Highlights * 32.0)
      points[3].Output = clamp(points[3].Output + highlightsShift, points[2].Output, 255)
      points[4].Output = clamp(points[4].Output + highlightsShift, points[3].Output, 255)

      // Apply shadows (adjust bottom-end points)
      shadowsShift := int(recipe.Shadows * 32.0)
      points[0].Output = clamp(points[0].Output + shadowsShift, 0, points[1].Output)
      points[1].Output = clamp(points[1].Output + shadowsShift, points[0].Output, points[2].Output)

      return points
  }
  ```
- [ ] Ensure monotonic curve (output[i] >= output[i-1])
- [ ] Test with various parameter combinations

### Task 2: Generate Camera Profile XML (AC-2, AC-4)
- [ ] Implement `generateProfile()` function in `profile.go`:
  ```go
  func generateProfile(recipe *universal.Recipe) ([]byte, error) {
      // Create profile struct
      profile := &CameraProfile{
          Xmlns:       "http://ns.adobe.com/camera-raw-settings/1.0/",
          ProfileName: "Recipe Converted Profile",
          ToneCurve: &ToneCurve{
              Points: universalToToneCurve(recipe),
          },
          ColorMatrix1: identityMatrix(),
          ColorMatrix2: identityMatrix(),
      }

      // Marshal to XML with indentation
      xmlData, err := xml.MarshalIndent(profile, "", "  ")
      if err != nil {
          return nil, fmt.Errorf("failed to marshal camera profile XML: %w", err)
      }

      // Prepend XML declaration
      output := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
      output = append(output, xmlData...)
      return output, nil
  }
  ```
- [ ] Implement `identityMatrix()` helper (returns 3x3 identity matrix)
- [ ] Validate XML output (well-formed, proper namespace)

### Task 3: Create TIFF File with Embedded XML (AC-1, AC-2)
- [ ] Implement `createTIFF()` function in `tiff.go`:
  ```go
  func createTIFF(xmlData []byte) ([]byte, error) {
      // Create new TIFF structure
      tiffFile := tiff.New()

      // Add standard TIFF tags (minimal IFD)
      tiffFile.SetTag(tiff.TagImageWidth, 1)
      tiffFile.SetTag(tiff.TagImageLength, 1)
      tiffFile.SetTag(tiff.TagSamplesPerPixel, 3)
      tiffFile.SetTag(tiff.TagBitsPerSample, []int{8, 8, 8})
      tiffFile.SetTag(tiff.TagPhotometricInterpretation, tiff.PhotoRGB)

      // Embed camera profile XML in tag 50740
      const tagCameraProfile = 50740
      tiffFile.SetTag(tagCameraProfile, xmlData)

      // Write TIFF to buffer
      buf := new(bytes.Buffer)
      if err := tiff.Encode(buf, tiffFile); err != nil {
          return nil, fmt.Errorf("failed to encode TIFF: %w", err)
      }

      return buf.Bytes(), nil
  }
  ```
- [ ] Verify TIFF magic bytes (II or MM) in output
- [ ] Validate IFD structure

### Task 4: Implement Generate() Function (AC-1 to AC-4)
- [ ] Implement `Generate(recipe *universal.Recipe) ([]byte, error)` in `generate.go`:
  ```go
  func Generate(recipe *universal.Recipe) ([]byte, error) {
      // Step 1: Generate camera profile XML
      xmlData, err := generateProfile(recipe)
      if err != nil {
          return nil, err
      }

      // Step 2: Create TIFF with embedded XML
      tiffData, err := createTIFF(xmlData)
      if err != nil {
          return nil, err
      }

      return tiffData, nil
  }
  ```
- [ ] Add error handling for nil recipe input
- [ ] Validate recipe parameters before generation

### Task 5: Write Unit Tests (AC-6)
- [ ] Write `TestGenerate_ValidRecipe()` - Generate DCP from populated UniversalRecipe
- [ ] Write `TestGenerate_EmptyRecipe()` - Generate from neutral recipe (all zeros)
- [ ] Write `TestGenerate_ExtremeValues()` - Test with extreme parameters (exposure=+2.0, contrast=+1.0)
- [ ] Write `TestUniversalToToneCurve()` - Test tone curve generation formulas
- [ ] Write `TestIdentityMatrix()` - Verify identity matrix generation
- [ ] Write `TestTIFFStructure()` - Validate generated TIFF structure (magic bytes, IFD, tag 50740)
- [ ] Write `TestXMLValidation()` - Validate generated XML (namespace, elements)
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
  - Document input (UniversalRecipe), output (DCP TIFF bytes), error cases
  - Include example usage
- [ ] Update `docs/parameter-mapping.md` with DCP generation mappings:
  - Document tone curve generation formulas (exposure/contrast/highlights/shadows → 5-point curve)
  - Note precision considerations (float → int curve points)
  - Provide examples with visual curve diagrams (optional)
- [ ] Add README notes in `testdata/dcp/`:
  - Document generated DCP validation results
  - Note Adobe software compatibility (Camera Raw, Lightroom versions tested)
  - List known limitations (no dual illuminant, no HSV tables, identity matrices only)

## Dev Notes

### Learnings from Previous Story

**From Story 9-1-dcp-parser (Status: drafted)**

- **Package Structure**: `internal/formats/dcp/` with types.go, parse.go, tiff.go, profile.go
- **TIFF Library**: Using `github.com/google/tiff` for TIFF reading/writing
- **Tone Curve Analysis**: Extract exposure/contrast/highlights/shadows from curve shape
- **Identity Matrices**: 3x3 diagonal matrix for non-calibration use
- **XML Parsing**: Adobe Camera Raw namespace, RDF structure (rdf:Seq, rdf:li)
- **Test Coverage**: ≥85% target, test with real Adobe DCP samples

**Reuse from Story 9-1:**
- `types.go` - CameraProfile, ToneCurve, Matrix struct definitions (use for generation)
- `tiff.go` - TIFF tag operations (extend for writing)
- `profile.go` - XML camera profile helpers (add generation functions)
- Test samples in `testdata/dcp/` (use for round-trip testing)

[Source: docs/stories/9-1-dcp-parser.md#Dev-Notes]

### Architecture Alignment

**Tech Spec Epic 9 Alignment:**

Story 9-2 implements **AC-2 (Generate DCP Files)** from tech-spec-epic-9.md.

**Generation Flow:**
```
UniversalRecipe → universalToToneCurve() → generateProfile() → createTIFF() → .dcp bytes
```

**Tone Curve Generation (5-point curve):**
```
Linear base:   (0,0)   (64,64)   (128,128)   (192,192)   (255,255)
                ↓         ↓          ↓            ↓           ↓
Apply exposure:         shifts midpoint (128 → 128+shift)
Apply contrast:         steepens/flattens slope (multiply deviation from 128)
Apply highlights:       adjusts top-end points (192-255)
Apply shadows:          adjusts bottom-end points (0-64)
                ↓         ↓          ↓            ↓           ↓
Output curve:  (0,X)   (64,Y)    (128,Z)      (192,W)     (255,V)
```

**Identity Matrix (3x3):**
```
ColorMatrix1 = ColorMatrix2 = [
    1.0  0.0  0.0
    0.0  1.0  0.0
    0.0  0.0  1.0
]
```

[Source: docs/tech-spec-epic-9.md#Detailed-Design]

### TIFF Writing Pattern

**Using github.com/google/tiff for Writing:**
```go
import "github.com/google/tiff"

// Create TIFF
tiffFile := tiff.New()

// Set standard tags
tiffFile.SetTag(tiff.TagImageWidth, 1)
tiffFile.SetTag(tiff.TagImageLength, 1)

// Set custom tag (50740 = CameraProfile)
const tagCameraProfile = 50740
tiffFile.SetTag(tagCameraProfile, xmlData)

// Encode to bytes
buf := new(bytes.Buffer)
err := tiff.Encode(buf, tiffFile)
```

**Minimal TIFF IFD (Image File Directory):**
- ImageWidth: 1 pixel (minimal, not an actual image)
- ImageLength: 1 pixel
- SamplesPerPixel: 3 (RGB, standard)
- BitsPerSample: [8, 8, 8]
- PhotometricInterpretation: RGB

[Source: docs/tech-spec-epic-9.md#APIs-and-Interfaces]

### XML Generation Pattern

**Adobe Camera Profile XML Structure:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<crs:CameraProfile xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
  <crs:ProfileName>Recipe Converted Profile</crs:ProfileName>
  <crs:ToneCurve>
    <rdf:Seq>
      <rdf:li>0, 0</rdf:li>
      <rdf:li>64, 64</rdf:li>
      <rdf:li>128, 140</rdf:li>  <!-- Exposure shift -->
      <rdf:li>192, 200</rdf:li>
      <rdf:li>255, 255</rdf:li>
    </rdf:Seq>
  </crs:ToneCurve>
  <crs:ColorMatrix1>
    <rdf:Seq>
      <rdf:li>1.0 0.0 0.0</rdf:li>
      <rdf:li>0.0 1.0 0.0</rdf:li>
      <rdf:li>0.0 0.0 1.0</rdf:li>
    </rdf:Seq>
  </crs:ColorMatrix1>
  <crs:ColorMatrix2>
    <!-- Same as ColorMatrix1 -->
  </crs:ColorMatrix2>
</crs:CameraProfile>
```

**Use `xml.MarshalIndent()` for Formatting:**
```go
xmlData, err := xml.MarshalIndent(profile, "", "  ")
```

[Source: docs/tech-spec-epic-9.md#Data-Models-and-Contracts]

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
- `internal/formats/dcp/tiff.go` - Add TIFF writing functions (extend from 9-1)
- `internal/formats/dcp/profile.go` - Add XML generation functions (extend from 9-1)

**Files from Story 9-1 (Reused):**
- `types.go` - Struct definitions (used for both parse and generate)
- `testdata/dcp/` - Sample files (use for round-trip testing)

[Source: docs/tech-spec-epic-9.md#Components]

### Testing Strategy

**Unit Tests (Required for AC-6):**
- `TestGenerate_ValidRecipe()` - Generate from populated UniversalRecipe
- `TestGenerate_EmptyRecipe()` - Generate neutral preset
- `TestGenerate_ExtremeValues()` - Test with extreme parameters
- `TestUniversalToToneCurve()` - Verify tone curve formulas
- `TestIdentityMatrix()` - Verify matrix generation
- `TestTIFFStructure()` - Validate TIFF structure
- `TestXMLValidation()` - Validate XML output
- `TestRoundTrip_DCP()` - Generate → Parse → Compare
- Coverage target: ≥85% for generate.go

**Manual Validation (Required for AC-5):**
- Load generated DCP in Adobe Camera Raw (no errors)
- Load generated DCP in Lightroom Classic (no errors)
- Visual spot-check (tone adjustments render correctly)
- Document results in validation-report.md

[Source: docs/tech-spec-epic-9.md#Test-Strategy-Summary]

### Known Risks

**RISK-16: Generated DCPs rejected by Adobe software**
- **Impact**: Lightroom/Camera Raw refuse to load generated profiles
- **Mitigation**: Follow Adobe DNG Specification 1.6 exactly, validate with real Adobe samples
- **Fallback**: Adjust TIFF/XML structure based on validation errors

**RISK-17: Tone curve formula inaccuracy**
- **Impact**: Generated curve doesn't match visual expectations
- **Mitigation**: Visual validation in Lightroom, iterate on formulas
- **Target**: 90%+ visual similarity to original preset

**RISK-18: Identity matrices rejected**
- **Impact**: Adobe software expects full calibration matrices
- **Mitigation**: Test with identity matrices first, add calibration if required
- **Fallback**: Generate minimal valid matrices (close to identity)

[Source: docs/tech-spec-epic-9.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-9.md#Acceptance-Criteria] - AC-2: Generate DCP Files
- [Source: docs/tech-spec-epic-9.md#Data-Models-and-Contracts] - Tone curve generation formulas
- [Source: docs/tech-spec-epic-9.md#APIs-and-Interfaces] - Generate() function signature
- [Source: github.com/google/tiff] - TIFF library writing API
- [Source: Adobe DNG Specification 1.6] - DCP format requirements
- [Source: internal/formats/dcp/parse.go] - Parse() function (reverse operation)
- [Source: internal/formats/costyle/generate.go] - XML generation pattern (reference)

## Dev Agent Record

### Context Reference

- Story Context XML: `docs/stories/9-2-dcp-generator.context.xml` (Generated: 2025-11-09)

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
