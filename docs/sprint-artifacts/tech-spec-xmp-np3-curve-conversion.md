# Tech-Spec: XMP/Lightroom to NP3 Curve Conversion Enhancement

**Created:** 2025-12-17  
**Status:** Ready for Development  
**Authors:** Winston (Architect), Mary (Analyst), Amelia (Dev)

---

## Overview

### Problem Statement

XMP/Lightroom presets use 50+ parameters for photo editing, while NP3's Picture Control format only supports ~48 parameters. Currently, 54% of XMP parameters are not being captured during conversion, resulting in significant visual differences in the converted NP3 output. The primary gap is:

1. **Wrong attribute names** in XMP parser (e.g., `HueRed` instead of `HueAdjustmentRed`)
2. **Missing parametric curve parsing** (7 parameters: Shadows, Darks, Lights, Highlights + 3 split points)
3. **No curve baking** - Lightroom effects cannot be approximated via NP3's tone curve

### Solution

A phased approach to improve XMP→NP3 conversion:

1. **Phase 0**: Generate test NP3 files with custom curves to validate NX Studio compatibility
2. **Phase 1**: Fix XMP parser attribute names to recover 33 missing parameters
3. **Phase 2**: Implement parametric curve → 257-entry LUT conversion
4. **Phase 3**: Add conversion warnings system for unsupported parameters

### Scope

| In Scope | Out of Scope |
|----------|--------------|
| Fix HSL attribute names (`HueAdjustmentRed`, etc.) | Clarity/Dehaze/Texture approximation (frequency-based effects) |
| Add `IncrementalTemperature`/`IncrementalTint` parsing | RGB channel curve conversion (NP3 only has master curve) |
| Parse parametric curve controls (7 params) | Lens correction parameters |
| Generate 257-entry LUT from parametric curves | Noise reduction parameters |
| Write LUT to NP3 offset 460 | Cross-format round-trip (costyle, lrtemplate) |
| NX Studio curve validation | Video LUT generation |
| Conversion warnings for unsupported params | |

---

## Context for Development

### Codebase Patterns

**NP3 Generator** (`internal/formats/np3/generate.go`):
- Uses TLV chunk structure for NX Studio compatibility
- Exact offset writing via helper functions (`writeBasicAdjustments`, `writeColorBlender`, etc.)
- `writeToneCurve()` at line 873 writes to offsets:
  - 404: Point count (1 byte)
  - 405: Control points (2 bytes per point: input, output)
  - 460: Raw curve LUT (257 × 16-bit big-endian = 514 bytes)
- Current `toneCurveRaw` field is `[]uint16` but rarely populated

**XMP Parser** (`internal/formats/xmp/parse.go`):
- Uses Go XML struct tags with `xml:"AttributeName,attr"` pattern
- `Description` struct (lines 76-225) defines all parsed attributes
- `extractParameters()` function converts string values to typed parameters
- Missing: `HueAdjustmentRed`, `IncrementalTemperature`, `ParametricShadows`, etc.

**Offset Constants** (`internal/formats/np3/offsets.go`):
- `OffsetToneCurvePointCount = 0x194` (404 decimal)
- `OffsetToneCurvePoints = 0x195` (405 decimal)
- `OffsetToneCurveRaw = 0x1CC` (460 decimal)

### Files to Reference

| File | Purpose |
|------|---------|
| `internal/formats/np3/generate.go` | NP3 binary generation, `writeToneCurve()` function |
| `internal/formats/np3/offsets.go` | Offset constants and encoding functions |
| `internal/formats/xmp/parse.go` | XMP parsing, `Description` struct |
| `internal/models/recipe.go` | `UniversalRecipe` struct |
| `docs/np3-format-specification.md` | NP3 format details |
| `docs/XMP_PARAMETER_COMPATIBILITY_ANALYSIS.md` | Gap analysis |

### Technical Decisions

1. **LUT Format**: 257-entry with 16-bit big-endian values (0-32767 range)
   - Entry 0 = black (input 0)
   - Entry 256 = white (input 255)
   - Middle entries interpolated

2. **Parametric Curve Zones**: Map Lightroom's 4-zone model
   - Shadows: 0-25% of tonal range
   - Darks: 25-50%
   - Lights: 50-75%
   - Highlights: 75-100%
   - Split points adjust zone boundaries

3. **Curve Baking Algorithm**: Apply parametric adjustments to linear base curve
   ```go
   for i := 0; i <= 256; i++ {
       normalized := float64(i) / 256.0
       output := applyParametricCurve(normalized, shadows, darks, lights, highlights, splits)
       lut[i] = uint16(output * 32767)
   }
   ```

---

## Implementation Plan

### Phase 0: NX Studio Curve Validation (Spike)

#### Task 0.1: Create Test Curve Generator

- [ ] **Add `GenerateTestCurve()` function** to `internal/formats/np3/generate.go`
  - Accepts curve type: "linear", "s-curve", "shadows-boost", "highlights-compress"
  - Returns `[257]uint16` suitable for writing to offset 460
  
- [ ] **Create test NP3 files** with known curves
  - `test_linear.np3` - Identity curve (y=x)
  - `test_scurve.np3` - Sigmoid contrast boost
  - `test_shadows_plus20.np3` - Lift shadows by 20%
  - `test_highlights_minus20.np3` - Compress highlights by 20%

#### Task 0.2: Manual NX Studio Testing

- [ ] **Justin validates in NX Studio**:
  1. Open each test NP3 file in Picture Control Utility 2 or NX Studio
  2. Verify curves display correctly in UI
  3. Apply to sample image, export, verify visual match
  4. Document any format quirks or rejections

---

### Phase 1: Fix XMP Parser

#### Task 1.1: Fix HSL Attribute Names

- [ ] **Update `Description` struct** in `parse.go` (lines 76-225)
  - Change `HueRed` → `HueAdjustmentRed` (and all 8 colors)
  - Change `SaturationRed` → `SaturationAdjustmentRed`  
  - Change `LuminanceRed` → `LuminanceAdjustmentRed`
  - Total: 24 attribute changes

#### Task 1.2: Add Incremental Temperature/Tint

- [ ] **Add new fields to `Description` struct**:
  ```go
  IncrementalTemperature string `xml:"IncrementalTemperature,attr"`
  IncrementalTint        string `xml:"IncrementalTint,attr"`
  ```

- [ ] **Update `extractParameters()`** to prefer Incremental over absolute values

#### Task 1.3: Add Parametric Curve Parsing

- [ ] **Add new fields to `Description` struct**:
  ```go
  ParametricShadows        string `xml:"ParametricShadows,attr"`
  ParametricDarks          string `xml:"ParametricDarks,attr"`
  ParametricLights         string `xml:"ParametricLights,attr"`
  ParametricHighlights     string `xml:"ParametricHighlights,attr"`
  ParametricShadowSplit    string `xml:"ParametricShadowSplit,attr"`
  ParametricMidtoneSplit   string `xml:"ParametricMidtoneSplit,attr"`
  ParametricHighlightSplit string `xml:"ParametricHighlightSplit,attr"`
  ```

- [ ] **Add fields to `xmpParameters` struct** and extraction logic
  
- [ ] **Add fields to `UniversalRecipe` struct** in `models/recipe.go`

---

### Phase 2: Parametric Curve → LUT Conversion

#### Task 2.1: Implement Curve Baking Algorithm

- [ ] **Create `curvebaker.go`** in `internal/formats/np3/`
  ```go
  // ParametricToLUT converts Lightroom parametric curve settings to a 257-entry LUT
  func ParametricToLUT(shadows, darks, lights, highlights int, 
                       shadowSplit, midtoneSplit, highlightSplit int) []uint16
  ```

- [ ] **Implement zone-based curve evaluation**:
  - Apply S-curve adjustments per zone based on slider values
  - Use split points to define zone boundaries
  - Ensure smooth transitions between zones (cubic interpolation)

#### Task 2.2: Integrate with NP3 Generator

- [ ] **Modify `convertToNP3Parameters()`** in `generate.go`
  - If recipe has parametric curve values, call `ParametricToLUT()`
  - Store result in `params.toneCurveRaw`

- [ ] **Verify `writeToneCurve()`** correctly writes LUT to offset 460

---

### Phase 3: Conversion Warnings

#### Task 3.1: Create Warning System

- [ ] **Define warning types** in `internal/models/warnings.go`
  ```go
  type ConversionWarning struct {
      Level       WarnLevel // Critical, Advisory, Info
      Parameter   string
      Message     string
      Alternative string // Suggested workaround
  }
  ```

- [ ] **Collect warnings during conversion** in NP3 generator

#### Task 3.2: Surface Warnings in CLI/Web

- [ ] **Add warnings to CLI output** with color-coded levels
- [ ] **Add warnings to web UI** with expandable details

---

## Acceptance Criteria

### Phase 0

- [ ] **AC-0.1**: Generated NP3 files with custom curves load in NX Studio without errors
- [ ] **AC-0.2**: Curve visualization in Picture Control Utility matches expected shape
- [ ] **AC-0.3**: Applied curves produce expected visual effect on sample images

### Phase 1

- [ ] **AC-1.1**: XMP files with `HueAdjustmentRed` etc. have values correctly extracted
- [ ] **AC-1.2**: Presets using `IncrementalTemperature` preserve white balance adjustments
- [ ] **AC-1.3**: All 7 parametric curve values are parsed from professional presets
- [ ] **AC-1.4**: Existing tests continue to pass (`go test ./internal/formats/xmp/...`)
- [ ] **AC-1.5**: Given the 6 Alex Ruskman test presets, when converted, then HSL values match source

### Phase 2

- [ ] **AC-2.1**: Parametric curve with Shadows=+50 produces visible shadow lift in LUT
- [ ] **AC-2.2**: Generated LUT values are within valid range (0-32767)
- [ ] **AC-2.3**: Round-trip test: Apply parametric curve in Lightroom → Convert to NP3 → Apply in NX Studio → Visual similarity ≥90% (SSIM ≥0.90)
- [ ] **AC-2.4**: LUT data correctly written to NP3 offset 460 (verified by hex inspection)

### Phase 3

- [ ] **AC-3.1**: Converting preset with Clarity=50 generates "Advisory" warning
- [ ] **AC-3.2**: Converting preset with unsupported HSL shows "Critical" warning
- [ ] **AC-3.3**: CLI displays color-coded warnings before conversion completes

---

## Additional Context

### Dependencies

- NP3 offset documentation in `docs/np3-format-specification.md`
- TypeScript reference: [ssssota/nikon-flexible-color-picture-control](https://github.com/ssssota/nikon-flexible-color-picture-control)
- Existing test fixtures in `testdata/xmp/` (914 files)

### Testing Strategy

**Automated Tests:**

1. **Unit tests for curve baking** (`internal/formats/np3/curvebaker_test.go`)
   ```bash
   go test -v ./internal/formats/np3/... -run Curve
   ```
   
2. **XMP parser attribute tests**
   ```bash
   go test -v ./internal/formats/xmp/... -run HSL
   go test -v ./internal/formats/xmp/... -run Parametric
   ```

3. **Integration tests for NP3 generation**
   ```bash
   go test -v ./internal/formats/np3/... -run TestGenerate
   ```

**Manual Verification (Justin):**

1. Open generated `test_*.np3` files in NX Studio → Verify curve display
2. Apply to sample RAW photo → Export JPEG → Compare to Lightroom export
3. Visual similarity assessment (target: ≥90%)

### Notes

- **File size**: Currently generating 1050-byte NP3 files, which accommodates the 514-byte raw curve (offset 460-973)
- **Saturation side effect**: NP3 curves affect saturation - lowering brightness increases saturation. Document for users.
- **Base curve**: NP3 curves are applied on top of a base Picture Control (Neutral, Standard, Vivid). Default to Neutral for most accurate conversion.

---

## Related Documents

- [xmp-to-np3-curve-research.md](../../.gemini/antigravity/brain/79032735-8ff4-401d-b456-ee0833440c7d/xmp-to-np3-curve-research.md) - Research findings
- [np3-format-specification.md](../np3-format-specification.md) - Binary format details
- [XMP_PARAMETER_COMPATIBILITY_ANALYSIS.md](../XMP_PARAMETER_COMPATIBILITY_ANALYSIS.md) - Gap analysis
- [parameter-mapping.md](../parameter-mapping.md) - Mapping rules
