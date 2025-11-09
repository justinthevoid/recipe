# Story 9.1: DNG Camera Profile (DCP) Parser

Status: ready-for-dev

## Story

As a **photographer**,
I want **Recipe to parse DNG Camera Profile (.dcp) files and extract tone curve adjustments**,
so that **I can convert Adobe camera profiles to other preset formats (NP3, XMP, lrtemplate, .costyle) and use them across different editing software**.

## Acceptance Criteria

**AC-1: Parse TIFF-based DCP Structure**
- ✅ Read TIFF file using `github.com/google/tiff` library
- ✅ Validate TIFF magic bytes (II for little-endian, MM for big-endian)
- ✅ Parse TIFF Image File Directory (IFD) structure
- ✅ Extract standard TIFF tags (ImageWidth, ImageLength, etc.)
- ✅ Handle both standalone .dcp files and embedded DCP in DNG (standalone priority)
- ✅ Report parsing errors with clear messages (corrupt TIFF, missing IFD, etc.)

**AC-2: Extract Camera Profile XML**
- ✅ Locate TIFF tag 50740 (CameraProfile tag)
- ✅ Extract XML data from tag 50740
- ✅ Validate XML structure (check for Adobe Camera Raw namespace)
- ✅ Parse ProfileName element (profile identity)
- ✅ Handle missing tag 50740 gracefully (return error: "not a valid DCP file")
- ✅ Support DCP v1.0-v1.6 format variations

**AC-3: Parse Tone Curve Adjustments**
- ✅ Parse `<crs:ToneCurve>` XML element (piecewise linear curve points)
- ✅ Extract tone curve points as (input, output) pairs (0-255 range)
- ✅ Analyze tone curve shape to extract parameters:
  - Exposure: Midpoint shift (point 128 vertical offset)
  - Contrast: Curve slope (difference between top/bottom points)
  - Highlights: Top-end curve shape (points 192-255)
  - Shadows: Bottom-end curve shape (points 0-64)
- ✅ Handle missing tone curve (use linear curve as fallback)
- ✅ Clamp extracted values to valid UniversalRecipe ranges

**AC-4: Parse Color Matrices (Identity Matrices)**
- ✅ Parse `<crs:ColorMatrix1>` and `<crs:ColorMatrix2>` (if present)
- ✅ Recognize identity matrices (diagonal 1.0, off-diagonal 0.0)
- ✅ Skip full camera calibration (ForwardMatrix, CalibrationIlluminant) for MVP
- ✅ Log warning if non-identity matrices detected ("color calibration not supported")
- ✅ Store matrix data in Metadata map for future use (optional)

**AC-5: Return UniversalRecipe Representation**
- ✅ Map parsed tone curve to UniversalRecipe fields:
  - recipe.Exposure (from midpoint shift)
  - recipe.Contrast (from curve slope)
  - recipe.Highlights (from top-end curve points)
  - recipe.Shadows (from bottom-end curve points)
- ✅ Preserve profile metadata (name, description if present)
- ✅ Store unsupported DCP data in Metadata map (matrices, HSV tables)
- ✅ Return populated `*universal.Recipe` and nil error on success

**AC-6: Handle Parsing Errors**
- ✅ Validate TIFF magic bytes before parsing (fail fast if not TIFF)
- ✅ Report specific errors (line number or tag ID if possible)
- ✅ Handle corrupt TIFF files gracefully (no panics)
- ✅ Handle malformed XML in tag 50740 (return descriptive error)
- ✅ Handle out-of-range tone curve values (clamp to 0-255)

**AC-7: Unit Test Coverage**
- ✅ Unit tests for Parse() function with real Adobe DCP samples
- ✅ Test edge cases (missing tone curve, identity matrices, corrupt TIFF)
- ✅ Test with minimum 3 real-world DCP files from Adobe
- ✅ Test coverage ≥85% for dcp/parse.go, dcp/tiff.go, dcp/profile.go
- ✅ All tests pass in CI

## Tasks / Subtasks

### Task 1: Set Up DCP Package Structure (AC: All)
- [ ] Create `internal/formats/dcp/` directory
- [ ] Create `types.go` - Define Go structs matching DCP XML schema:
  ```go
  type CameraProfile struct {
      XMLName     xml.Name `xml:"CameraProfile"`
      Xmlns       string   `xml:"xmlns:crs,attr"`
      ProfileName string   `xml:"crs:ProfileName"`
      ToneCurve   *ToneCurve `xml:"crs:ToneCurve"`
      ColorMatrix1 *Matrix   `xml:"crs:ColorMatrix1,omitempty"`
      ColorMatrix2 *Matrix   `xml:"crs:ColorMatrix2,omitempty"`
  }

  type ToneCurve struct {
      Points []ToneCurvePoint `xml:"rdf:li"`
  }

  type ToneCurvePoint struct {
      Input  int `xml:",chardata"` // 0-255
      Output int // Parsed from chardata "input, output"
  }

  type Matrix struct {
      Rows []MatrixRow `xml:"rdf:li"`
  }

  type MatrixRow struct {
      Values [3]float64 `xml:",chardata"` // Parsed from "v1 v2 v3"
  }
  ```
- [ ] Create `parse.go` - Implement Parse(data []byte) function
- [ ] Create `tiff.go` - TIFF tag reading helpers
- [ ] Create `profile.go` - XML camera profile parsing
- [ ] Create `parse_test.go` - Unit tests
- [ ] Create `testdata/dcp/` - Sample .dcp files

### Task 2: Add github.com/google/tiff Dependency (AC-1)
- [ ] Add dependency to go.mod:
  ```bash
  go get github.com/google/tiff
  ```
- [ ] Verify dependency downloads successfully
- [ ] Import in tiff.go: `import "github.com/google/tiff"`
- [ ] Document dependency in docs/dependencies.md (if exists)

### Task 3: Implement TIFF Reading (AC-1)
- [ ] Implement `readTIFF()` helper in `tiff.go`:
  ```go
  func readTIFF(data []byte) (*tiff.TIFF, error) {
      // Validate TIFF magic bytes
      if len(data) < 4 {
          return nil, fmt.Errorf("file too small to be a TIFF")
      }

      // Check for little-endian (II) or big-endian (MM)
      if !bytes.Equal(data[:2], []byte{'I', 'I'}) && !bytes.Equal(data[:2], []byte{'M', 'M'}) {
          return nil, fmt.Errorf("invalid TIFF magic bytes")
      }

      // Parse TIFF using google/tiff library
      tiffFile, err := tiff.Decode(bytes.NewReader(data))
      if err != nil {
          return nil, fmt.Errorf("failed to parse TIFF structure: %w", err)
      }

      return tiffFile, nil
  }
  ```
- [ ] Implement `extractCameraProfileTag()` helper:
  ```go
  func extractCameraProfileTag(tiffFile *tiff.TIFF) ([]byte, error) {
      // TIFF tag 50740 = CameraProfile
      const tagCameraProfile = 50740

      // Get tag value (XML bytes)
      tag := tiffFile.GetTag(tagCameraProfile)
      if tag == nil {
          return nil, fmt.Errorf("tag 50740 (CameraProfile) not found - not a valid DCP file")
      }

      xmlData, ok := tag.Value.([]byte)
      if !ok {
          return nil, fmt.Errorf("tag 50740 value is not bytes")
      }

      return xmlData, nil
  }
  ```
- [ ] Handle TIFF parsing errors (corrupt file, missing IFD)

### Task 4: Implement Camera Profile XML Parsing (AC-2, AC-3, AC-4)
- [ ] Implement `parseProfile()` function in `profile.go`:
  ```go
  func parseProfile(xmlData []byte) (*CameraProfile, error) {
      var profile CameraProfile

      // Unmarshal XML
      if err := xml.Unmarshal(xmlData, &profile); err != nil {
          return nil, fmt.Errorf("failed to parse camera profile XML: %w", err)
      }

      // Validate namespace
      if !strings.Contains(profile.Xmlns, "adobe.com/camera-raw-settings") {
          return nil, fmt.Errorf("invalid camera profile XML namespace")
      }

      return &profile, nil
  }
  ```
- [ ] Implement tone curve parsing:
  - Parse `<rdf:li>` elements in `<crs:ToneCurve>`
  - Each point format: "input, output" (e.g., "0, 0", "255, 255")
  - Split string, parse integers:
    ```go
    func parseToneCurvePoint(s string) (ToneCurvePoint, error) {
        parts := strings.Split(s, ",")
        if len(parts) != 2 {
            return ToneCurvePoint{}, fmt.Errorf("invalid tone curve point format")
        }
        input, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
        output, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
        return ToneCurvePoint{Input: input, Output: output}, nil
    }
    ```
- [ ] Implement color matrix parsing:
  - Parse `<rdf:li>` elements in `<crs:ColorMatrix1>`
  - Each row format: "v1 v2 v3" (e.g., "1.0 0.0 0.0")
  - Split string, parse floats
- [ ] Detect identity matrices (diagonal 1.0, off-diagonal 0.0)

### Task 5: Analyze Tone Curve Shape (AC-3)
- [ ] Implement `analyzeToneCurve()` function in `profile.go`:
  ```go
  func analyzeToneCurve(points []ToneCurvePoint) (exposure, contrast, highlights, shadows float64) {
      if len(points) == 0 {
          return 0, 0, 0, 0 // Linear curve
      }

      // Find midpoint (input=128)
      midpoint := findPoint(points, 128)
      // Exposure = vertical shift from linear (128 → X)
      exposure = float64(midpoint.Output - 128) / 128.0 // Normalize to -1.0/+1.0

      // Contrast = slope difference (top - bottom)
      topPoint := findPoint(points, 192)
      bottomPoint := findPoint(points, 64)
      slopeDiff := float64(topPoint.Output - bottomPoint.Output) / 128.0
      contrast = (slopeDiff - 1.0) // Normalize (linear slope = 1.0)

      // Highlights = top-end deviation (192-255 range)
      highlightsPoint := findPoint(points, 255)
      highlights = float64(highlightsPoint.Output - 255) / 64.0 // Normalize

      // Shadows = bottom-end deviation (0-64 range)
      shadowsPoint := findPoint(points, 0)
      shadows = float64(shadowsPoint.Output - 0) / 64.0 // Normalize

      return exposure, contrast, highlights, shadows
  }

  func findPoint(points []ToneCurvePoint, input int) ToneCurvePoint {
      // Find exact match or interpolate
      for _, p := range points {
          if p.Input == input {
              return p
          }
      }
      // Interpolate if exact match not found
      // ... linear interpolation logic ...
  }
  ```
- [ ] Clamp extracted values to valid UniversalRecipe ranges:
  - Exposure: -2.0 to +2.0
  - Contrast: -1.0 to +1.0
  - Highlights: -1.0 to +1.0
  - Shadows: -1.0 to +1.0

### Task 6: Implement Parse() Function (AC-5)
- [ ] Implement `Parse(data []byte) (*universal.Recipe, error)` in `parse.go`:
  ```go
  func Parse(data []byte) (*universal.Recipe, error) {
      // Step 1: Read TIFF structure
      tiffFile, err := readTIFF(data)
      if err != nil {
          return nil, fmt.Errorf("failed to read DCP TIFF: %w", err)
      }

      // Step 2: Extract camera profile XML from tag 50740
      xmlData, err := extractCameraProfileTag(tiffFile)
      if err != nil {
          return nil, err
      }

      // Step 3: Parse camera profile XML
      profile, err := parseProfile(xmlData)
      if err != nil {
          return nil, err
      }

      // Step 4: Analyze tone curve to extract parameters
      exposure, contrast, highlights, shadows := analyzeToneCurve(profile.ToneCurve.Points)

      // Step 5: Create UniversalRecipe
      recipe := &universal.Recipe{
          Exposure:   exposure,
          Contrast:   contrast,
          Highlights: highlights,
          Shadows:    shadows,
          Metadata:   make(map[string]interface{}),
      }

      // Step 6: Store profile name in metadata
      recipe.Metadata["profile_name"] = profile.ProfileName

      // Step 7: Check for non-identity matrices (log warning)
      if profile.ColorMatrix1 != nil && !isIdentityMatrix(profile.ColorMatrix1) {
          log.Warn("DCP contains color calibration matrices (not supported)")
          recipe.Metadata["color_matrix_1"] = profile.ColorMatrix1
      }

      return recipe, nil
  }
  ```
- [ ] Implement `isIdentityMatrix()` helper:
  ```go
  func isIdentityMatrix(matrix *Matrix) bool {
      expected := [][]float64{
          {1.0, 0.0, 0.0},
          {0.0, 1.0, 0.0},
          {0.0, 0.0, 1.0},
      }
      for i := 0; i < 3; i++ {
          for j := 0; j < 3; j++ {
              if math.Abs(matrix.Rows[i].Values[j] - expected[i][j]) > 0.001 {
                  return false
              }
          }
      }
      return true
  }
  ```

### Task 7: Error Handling (AC-6)
- [ ] Validate TIFF magic bytes in `readTIFF()` (fail fast)
- [ ] Wrap all errors with descriptive messages using `fmt.Errorf` with `%w` verb
- [ ] Handle corrupt TIFF gracefully (catch panics from google/tiff library if needed)
- [ ] Handle malformed XML in tag 50740 (xml.Unmarshal errors)
- [ ] Clamp out-of-range tone curve values before returning:
  ```go
  func clampFloat64(value, min, max float64) float64 {
      if value < min {
          return min
      }
      if value > max {
          return max
      }
      return value
  }

  // Apply before returning
  recipe.Exposure = clampFloat64(exposure, -2.0, 2.0)
  recipe.Contrast = clampFloat64(contrast, -1.0, 1.0)
  recipe.Highlights = clampFloat64(highlights, -1.0, 1.0)
  recipe.Shadows = clampFloat64(shadows, -1.0, 1.0)
  ```

### Task 8: Acquire Real DCP Sample Files (AC-7)
- [ ] Download Adobe DCP samples:
  - Adobe Camera Raw SDK includes sample DCP files
  - URL: https://www.adobe.com/go/dng_profiles (if available)
  - Or create synthetic DCP files using Adobe DNG Profile Editor
- [ ] Acquire 3+ different DCP profiles:
  - Standard Adobe camera profile (e.g., Canon 5D Mark IV)
  - Nikon camera profile (e.g., Nikon D850)
  - Custom DCP profile with tone curves (portrait, landscape)
- [ ] Add samples to `testdata/dcp/` directory:
  ```
  testdata/dcp/
  ├── README.md           # Document sample sources
  ├── canon-standard.dcp
  ├── nikon-standard.dcp
  └── custom-portrait.dcp
  ```
- [ ] Verify samples are valid DCP files (open in Lightroom or Adobe DNG Converter)

### Task 9: Write Unit Tests (AC-7)
- [ ] Write `TestParse_ValidDCP()` - Parse real Adobe DCP sample:
  ```go
  func TestParse_ValidDCP(t *testing.T) {
      data, err := os.ReadFile("testdata/dcp/canon-standard.dcp")
      require.NoError(t, err)

      recipe, err := dcp.Parse(data)
      require.NoError(t, err)
      assert.NotNil(t, recipe)

      // Verify tone curve parameters extracted
      assert.InDelta(t, 0.0, recipe.Exposure, 0.1) // Near neutral
      assert.InDelta(t, 0.0, recipe.Contrast, 0.1)

      // Verify profile name in metadata
      assert.Contains(t, recipe.Metadata, "profile_name")
  }
  ```
- [ ] Write `TestParse_MissingToneCurve()` - DCP without tone curve:
  - Should return neutral recipe (all parameters zero)
- [ ] Write `TestParse_NonIdentityMatrix()` - DCP with color calibration:
  - Should log warning, store matrix in metadata
- [ ] Write `TestParse_CorruptTIFF()` - Malformed TIFF file:
  - Should return error: "invalid TIFF magic bytes" or "failed to parse TIFF"
- [ ] Write `TestParse_MissingTag50740()` - TIFF without camera profile tag:
  - Should return error: "not a valid DCP file"
- [ ] Write `TestParse_MalformedXML()` - Tag 50740 contains invalid XML:
  - Should return error: "failed to parse camera profile XML"
- [ ] Write `TestAnalyzeToneCurve()` - Unit test for tone curve analysis:
  - Test with linear curve (0,0 → 255,255)
  - Test with exposure shift (+0.5, +1.0)
  - Test with contrast adjustment (steeper slope)
  - Test with highlights/shadows adjustment
- [ ] Run tests: `go test ./internal/formats/dcp/`
- [ ] Verify coverage: `go test -cover ./internal/formats/dcp/` (target ≥85%)

### Task 10: Documentation (AC-1 to AC-6)
- [ ] Add package comment in `parse.go`:
  ```go
  // Package dcp provides parsing and generation of DNG Camera Profile (.dcp) files.
  // DCPs are TIFF-based files containing Adobe Camera Profile XML data embedded
  // in TIFF tag 50740. Recipe supports tone curve adjustments only (exposure,
  // contrast, highlights, shadows), not full camera calibration.
  ```
- [ ] Add function comment for `Parse()`:
  - Document input (DCP TIFF bytes), output (UniversalRecipe), error cases
  - Include example usage
- [ ] Update `docs/parameter-mapping.md` with DCP mappings:
  - Document tone curve → UniversalRecipe mapping formulas
  - Note unsupported DCP features (color matrices, HSV tables, dual illuminant)
  - Provide examples (midpoint shift = exposure, slope = contrast)
- [ ] Add README in `testdata/dcp/`:
  - Document sample DCP file sources
  - Note DCP version support (v1.0-v1.6)
  - List known limitations

## Dev Notes

### Learnings from Previous Story

**From Story 8-5-costyle-integration (Status: drafted)**

- **Integration Pattern**: Extend Converter.Convert() to route format to Parse/Generate
- **Format Detection**: Extension check first, then magic bytes validation
- **Error Handling**: Consistent error messages across CLI/TUI/Web interfaces
- **Documentation**: Update README, help text, FAQ for new format

**Reuse from Story 8-5:**
- Format detection pattern (extension + magic bytes)
- Converter integration approach (will be used in Story 9-4)
- Error message consistency (apply to DCP parsing errors)
- Test coverage standard (≥85%)

[Source: docs/stories/8-5-costyle-integration.md#Dev-Notes]

### Architecture Alignment

**Tech Spec Epic 9 Alignment:**

Story 9-1 implements **AC-1 (Parse DCP Files)** from tech-spec-epic-9.md.

**Parse Flow:**
```
.dcp file bytes → readTIFF() → extractCameraProfileTag() → parseProfile() → analyzeToneCurve() → UniversalRecipe
```

**TIFF Structure:**
```
DCP File (.dcp)
├── TIFF Header (II or MM)
├── Image File Directory (IFD)
│   ├── Tag 50740: CameraProfile (XML data)
│   └── Other standard TIFF tags
```

**Tone Curve Analysis:**
- Exposure: Midpoint (128) vertical shift
- Contrast: Slope difference (top - bottom points)
- Highlights: Top-end curve shape (192-255 points)
- Shadows: Bottom-end curve shape (0-64 points)

[Source: docs/tech-spec-epic-9.md#Detailed-Design]

### TIFF Library (github.com/google/tiff)

**Why google/tiff?**
- Go stdlib `image/tiff` is decoder-only, doesn't support writing
- google/tiff supports custom TIFF tag reading/writing (tag 50740)
- Google-maintained, stable, widely used in production
- Approved in architecture decision (Decision 4)

**Usage Pattern:**
```go
import "github.com/google/tiff"

// Read TIFF
tiffFile, err := tiff.Decode(reader)

// Get custom tag
tag := tiffFile.GetTag(50740) // CameraProfile
xmlData := tag.Value.([]byte)
```

**Error Handling:**
- `tiff.Decode()` returns errors for corrupt files (no panics)
- Validate magic bytes before calling Decode (fail fast)

[Source: docs/tech-spec-epic-9.md#System-Architecture-Alignment]

### DCP XML Format (Adobe Camera Profile)

**XML Namespace:**
```xml
<crs:CameraProfile xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
```

**Tone Curve Format:**
```xml
<crs:ToneCurve>
  <rdf:Seq>
    <rdf:li>0, 0</rdf:li>
    <rdf:li>64, 50</rdf:li>
    <rdf:li>128, 128</rdf:li>
    <rdf:li>192, 200</rdf:li>
    <rdf:li>255, 255</rdf:li>
  </rdf:Seq>
</crs:ToneCurve>
```

**Color Matrix Format:**
```xml
<crs:ColorMatrix1>
  <rdf:Seq>
    <rdf:li>1.0 0.0 0.0</rdf:li>
    <rdf:li>0.0 1.0 0.0</rdf:li>
    <rdf:li>0.0 0.0 1.0</rdf:li>
  </rdf:Seq>
</crs:ColorMatrix1>
```

[Source: docs/tech-spec-epic-9.md#Data-Models-and-Contracts]

### Project Structure Notes

**New Files Created (Story 9-1):**
```
internal/formats/dcp/
├── types.go           # DCP-specific structs (NEW)
├── parse.go           # Parse() function (NEW)
├── tiff.go            # TIFF reading helpers (NEW)
├── profile.go         # XML camera profile parsing (NEW)
├── parse_test.go      # Unit tests (NEW)
└── testdata/dcp/      # Sample DCP files (NEW)
    ├── README.md
    ├── canon-standard.dcp
    ├── nikon-standard.dcp
    └── custom-portrait.dcp
```

**Modified Files:**
- `go.mod` - Add github.com/google/tiff dependency

**Files from Epic 8 (Pattern Reference):**
- `internal/formats/costyle/parse.go` - XML parsing pattern
- `internal/formats/costyle/types.go` - Struct definition pattern

[Source: docs/tech-spec-epic-9.md#Components]

### Testing Strategy

**Unit Tests (Required for AC-7):**
- `TestParse_ValidDCP()` - Parse real Adobe DCP sample
- `TestParse_MissingToneCurve()` - DCP without tone curve
- `TestParse_NonIdentityMatrix()` - DCP with color calibration
- `TestParse_CorruptTIFF()` - Malformed TIFF file
- `TestParse_MissingTag50740()` - TIFF without camera profile tag
- `TestParse_MalformedXML()` - Invalid XML in tag 50740
- `TestAnalyzeToneCurve()` - Tone curve analysis logic
- Coverage target: ≥85% for parse.go, tiff.go, profile.go

**Manual Validation (Story 9-4):**
- Load generated DCP in Adobe Camera Raw (verify no errors)
- Load generated DCP in Lightroom Classic (verify no errors)
- Visual spot-check (verify tone adjustments render correctly)

[Source: docs/tech-spec-epic-9.md#Test-Strategy-Summary]

### Known Risks

**RISK-13: DCP format variations across versions**
- **Impact**: v1.0-v1.6 may have structural differences
- **Mitigation**: Test with multiple DCP versions, document supported versions
- **Fallback**: Support v1.4+ only (most common)

**RISK-14: Tone curve analysis inaccuracy**
- **Impact**: Extracted exposure/contrast may not match visual output
- **Mitigation**: Visual validation in Lightroom (Story 9-4), adjust formulas if needed
- **Target**: 90%+ visual similarity (accept some precision loss)

**RISK-15: Adobe DCP samples unavailable**
- **Impact**: Cannot test with real-world files
- **Mitigation**: Use Adobe DNG Profile Editor to create synthetic DCPs
- **Fallback**: Generate minimal valid DCP for testing (identity matrices + linear curve)

[Source: docs/tech-spec-epic-9.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-9.md#Acceptance-Criteria] - AC-1: Parse DCP Files
- [Source: docs/tech-spec-epic-9.md#Data-Models-and-Contracts] - DCP XML structure
- [Source: docs/tech-spec-epic-9.md#APIs-and-Interfaces] - Parse() function signature
- [Source: github.com/google/tiff] - TIFF library documentation
- [Source: Adobe DNG Specification 1.6] - DCP format specification (external)
- [Source: internal/formats/costyle/parse.go] - XML parsing pattern (reference)

## Dev Agent Record

### Context Reference

- docs/stories/9-1-dcp-parser.context.xml

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
