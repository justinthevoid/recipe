# NP3 Parser Enhancement Strategy - Phase 2

**Created**: 2025-11-07
**Status**: Phase 2 Day 1 - Design Complete
**Goal**: Replace heuristic parsing with exact offset-based extraction while maintaining backward compatibility

---

## Current Parser Analysis

### Existing Structure (parse.go:53-427)

The current parser uses a **4-stage pipeline**:

```
Parse() → validateFileStructure() → extractParameters() → validateParameters() → buildRecipe()
```

**Stage 1 - File Validation** (lines 80-98):
- ✅ Checks magic bytes "NCP"
- ✅ Validates minimum file size (300 bytes)
- **Status**: Keep unchanged

**Stage 2 - Parameter Extraction** (lines 137-251):
- ❌ Uses heuristic analysis of raw bytes (64-80)
- ❌ Analyzes RGB triplets (100-300) for saturation
- ❌ Analyzes tone curve pairs (150-500) for contrast
- ❌ Achieves ~95% accuracy but misses 44 parameters
- **Status**: Replace with exact offset extraction

**Stage 3 - Parameter Validation** (lines 359-386):
- ✅ Uses models.Validate* functions
- ⚠️ Only validates 5 parameters (sharpening, contrast, brightness, saturation, hue)
- **Status**: Expand to validate all 56 parameters

**Stage 4 - Recipe Building** (lines 388-427):
- ✅ Uses builder pattern
- ✅ Preserves rawData for round-trip
- ⚠️ Only sets 4 parameters (maps 5 extracted → 4 UniversalRecipe fields)
- **Status**: Expand to set all 56 parameters

---

## Enhancement Strategy: Dual-Mode Approach

### Core Principle: Gradual Migration

Rather than replacing the entire parser at once, we'll use a **dual-mode approach** that:
1. Keeps existing heuristic extraction as fallback
2. Adds exact offset extraction alongside
3. Prefers exact offsets when available
4. Maintains 100% backward compatibility

### Implementation Pattern: Additive, Not Replacement

```go
// OLD (heuristic)
params.contrast = estimateContrastFromToneCurve(toneCurve)

// NEW (exact + fallback)
if ValidateOffset(OffsetContrast) && len(data) > OffsetContrast {
    params.contrast = DecodeSigned8(data[OffsetContrast])  // Exact
} else {
    params.contrast = estimateContrastFromToneCurve(toneCurve)  // Fallback
}
```

**Benefits**:
- ✅ No breaking changes to existing behavior
- ✅ Improves accuracy from ~95% to ~100% for standard NP3 files
- ✅ Maintains support for variant/malformed files
- ✅ Easy to test incrementally (one parameter at a time)

---

## Detailed Implementation Plan

### Step 1: Extend np3Parameters Struct (parse.go:100-109)

**Current struct** (5 fields):
```go
type np3Parameters struct {
    name       string
    sharpening int
    contrast   int
    brightness float64
    saturation int
    hue        int
    rawData    []byte
}
```

**Enhanced struct** (48 fields):
```go
type np3Parameters struct {
    // Header/Metadata
    name string

    // Basic Adjustments (exact)
    sharpening float64  // Changed from int to float64 for Scaled4 encoding
    clarity    float64  // NEW

    // Advanced Adjustments (exact)
    midRangeSharpening float64  // NEW
    contrast           int
    highlights         int  // NEW
    shadows            int  // NEW
    whiteLevel         int  // NEW
    blackLevel         int  // NEW
    saturation         int

    // Color Blender (exact) - 8 colors × 3 values = 24 fields
    redHue        int  // NEW
    redChroma     int  // NEW
    redBrightness int  // NEW
    // ... (repeat for Orange, Yellow, Green, Cyan, Blue, Purple, Magenta)

    // Color Grading (exact) - 3 zones + 2 global = 11 fields
    highlightsZone ColorGradingZone  // NEW
    midtoneZone    ColorGradingZone  // NEW
    shadowsZone    ColorGradingZone  // NEW
    blending       int               // NEW
    balance        int               // NEW

    // Tone Curve (exact)
    toneCurvePointCount int             // NEW
    toneCurvePoints     []ToneCurvePoint // NEW
    toneCurveRaw        []uint16         // NEW

    // Legacy fields (deprecated, kept for fallback)
    brightness float64  // DEPRECATED: Use exposure calculation instead
    hue        int      // DEPRECATED: No UniversalRecipe equivalent

    // Metadata
    rawData []byte  // Keep for round-trip preservation
}
```

### Step 2: Add Exact Extraction Functions

Create new functions in parse.go for each parameter category:

#### 2.1 Basic Adjustments (parse.go:~450)

```go
// extractBasicAdjustments reads sharpening and clarity using exact offsets.
func extractBasicAdjustments(data []byte, params *np3Parameters) {
    // Sharpening (offset 82, Scaled4 encoding)
    if ValidateOffset(OffsetSharpening) && len(data) > OffsetSharpening {
        params.sharpening = DecodeScaled4(data[OffsetSharpening])
    }

    // Clarity (offset 92, Scaled4 encoding)
    if ValidateOffset(OffsetClarity) && len(data) > OffsetClarity {
        params.clarity = DecodeScaled4(data[OffsetClarity])
    }
}
```

**Lines**: ~15 lines
**Complexity**: Low (2 parameters, same encoding pattern)

#### 2.2 Advanced Adjustments (parse.go:~470)

```go
// extractAdvancedAdjustments reads all 7 advanced adjustment parameters using exact offsets.
func extractAdvancedAdjustments(data []byte, params *np3Parameters) {
    // Mid-Range Sharpening (offset 242, Scaled4)
    if ValidateOffset(OffsetMidRangeSharpening) && len(data) > OffsetMidRangeSharpening {
        params.midRangeSharpening = DecodeScaled4(data[OffsetMidRangeSharpening])
    }

    // Contrast (offset 272, Signed8)
    if ValidateOffset(OffsetContrast) && len(data) > OffsetContrast {
        params.contrast = DecodeSigned8(data[OffsetContrast])
    }

    // Highlights (offset 282, Signed8)
    if ValidateOffset(OffsetHighlights) && len(data) > OffsetHighlights {
        params.highlights = DecodeSigned8(data[OffsetHighlights])
    }

    // Shadows (offset 292, Signed8)
    if ValidateOffset(OffsetShadows) && len(data) > OffsetShadows {
        params.shadows = DecodeSigned8(data[OffsetShadows])
    }

    // White Level (offset 302, Signed8)
    if ValidateOffset(OffsetWhiteLevel) && len(data) > OffsetWhiteLevel {
        params.whiteLevel = DecodeSigned8(data[OffsetWhiteLevel])
    }

    // Black Level (offset 312, Signed8)
    if ValidateOffset(OffsetBlackLevel) && len(data) > OffsetBlackLevel {
        params.blackLevel = DecodeSigned8(data[OffsetBlackLevel])
    }

    // Saturation (offset 322, Signed8)
    if ValidateOffset(OffsetSaturation) && len(data) > OffsetSaturation {
        params.saturation = DecodeSigned8(data[OffsetSaturation])
    }
}
```

**Lines**: ~35 lines
**Complexity**: Medium (7 parameters, 2 encoding patterns: Scaled4 + Signed8)

#### 2.3 Color Blender (parse.go:~510)

```go
// extractColorBlender reads all 24 color blender parameters (8 colors × 3 values) using exact offsets.
func extractColorBlender(data []byte, params *np3Parameters) {
    // Red color
    if ValidateOffsetRange(OffsetRedHue, OffsetRedBrightness+1) && len(data) > OffsetRedBrightness {
        params.redHue = DecodeSigned8(data[OffsetRedHue])
        params.redChroma = DecodeSigned8(data[OffsetRedChroma])
        params.redBrightness = DecodeSigned8(data[OffsetRedBrightness])
    }

    // Orange color
    if ValidateOffsetRange(OffsetOrangeHue, OffsetOrangeBrightness+1) && len(data) > OffsetOrangeBrightness {
        params.orangeHue = DecodeSigned8(data[OffsetOrangeHue])
        params.orangeChroma = DecodeSigned8(data[OffsetOrangeChroma])
        params.orangeBrightness = DecodeSigned8(data[OffsetOrangeBrightness])
    }

    // ... (repeat for Yellow, Green, Cyan, Blue, Purple, Magenta)
}
```

**Lines**: ~50 lines
**Complexity**: Low (repetitive pattern, single encoding: Signed8)
**Optimization**: Could use a loop with offset array for DRY code

#### 2.4 Color Grading (parse.go:~565)

```go
// extractColorGrading reads all 11 color grading parameters (3 zones + 2 global) using exact offsets.
func extractColorGrading(data []byte, params *np3Parameters) {
    // Highlights zone (4 bytes: 2-byte hue + 1-byte chroma + 1-byte brightness)
    if ValidateOffsetRange(OffsetHighlightsHue, OffsetHighlightsBrightness+1) && len(data) > OffsetHighlightsBrightness {
        params.highlightsZone.Hue = DecodeHue12(data[OffsetHighlightsHue], data[OffsetHighlightsHue+1])
        params.highlightsZone.Chroma = DecodeSigned8(data[OffsetHighlightsChroma])
        params.highlightsZone.Brightness = DecodeSigned8(data[OffsetHighlightsBrightness])
    }

    // Midtone zone (4 bytes)
    if ValidateOffsetRange(OffsetMidtoneHue, OffsetMidtoneBrightness+1) && len(data) > OffsetMidtoneBrightness {
        params.midtoneZone.Hue = DecodeHue12(data[OffsetMidtoneHue], data[OffsetMidtoneHue+1])
        params.midtoneZone.Chroma = DecodeSigned8(data[OffsetMidtoneChroma])
        params.midtoneZone.Brightness = DecodeSigned8(data[OffsetMidtoneBrightness])
    }

    // Shadows zone (4 bytes)
    if ValidateOffsetRange(OffsetShadowsHue, OffsetShadowsBrightness+1) && len(data) > OffsetShadowsBrightness {
        params.shadowsZone.Hue = DecodeHue12(data[OffsetShadowsHue], data[OffsetShadowsHue+1])
        params.shadowsZone.Chroma = DecodeSigned8(data[OffsetShadowsChroma])
        params.shadowsZone.Brightness = DecodeSigned8(data[OffsetShadowsBrightness])
    }

    // Global parameters
    if ValidateOffset(OffsetColorGradingBlending) && len(data) > OffsetColorGradingBlending {
        params.blending = int(data[OffsetColorGradingBlending])  // No bias, direct value
    }

    if ValidateOffset(OffsetColorGradingBalance) && len(data) > OffsetColorGradingBalance {
        params.balance = DecodeSigned8(data[OffsetColorGradingBalance])
    }
}
```

**Lines**: ~30 lines
**Complexity**: Medium (3 zones with multi-byte hue encoding, 2 global params)

#### 2.5 Tone Curve (parse.go:~600)

```go
// extractToneCurve reads tone curve parameters using exact offsets.
func extractToneCurve(data []byte, params *np3Parameters) {
    // Point count (offset 404)
    if ValidateOffset(OffsetToneCurvePointCount) && len(data) > OffsetToneCurvePointCount {
        params.toneCurvePointCount = int(data[OffsetToneCurvePointCount])

        // Control points (offset 405, 2 bytes per point)
        if params.toneCurvePointCount > 0 && params.toneCurvePointCount <= 127 {
            pointsOffset := OffsetToneCurvePoints
            pointsEnd := pointsOffset + (params.toneCurvePointCount * 2)

            if ValidateOffsetRange(pointsOffset, pointsEnd) && len(data) >= pointsEnd {
                params.toneCurvePoints = make([]ToneCurvePoint, params.toneCurvePointCount)
                for i := 0; i < params.toneCurvePointCount; i++ {
                    offset := pointsOffset + (i * 2)
                    params.toneCurvePoints[i] = ToneCurvePoint{
                        X: data[offset],
                        Y: data[offset+1],
                    }
                }
            }
        }
    }

    // Raw curve data (offset 460, 257 × 16-bit big-endian)
    // Note: Standard 480-byte files cannot contain the full 514-byte raw curve
    if ValidateOffset(OffsetToneCurveRaw) && len(data) >= OffsetToneCurveRaw+514 {
        params.toneCurveRaw = make([]uint16, 257)
        for i := 0; i < 257; i++ {
            offset := OffsetToneCurveRaw + (i * 2)
            params.toneCurveRaw[i] = binary.BigEndian.Uint16(data[offset:offset+2])
        }
    }
}
```

**Lines**: ~35 lines
**Complexity**: High (variable-length arrays, 16-bit big-endian values, file size considerations)

### Step 3: Modify extractParameters() Function

Replace heuristic extraction with exact extraction:

```go
func extractParameters(data []byte) (*np3Parameters, error) {
    params := &np3Parameters{}

    // Store raw data for chunk preservation (KEEP)
    params.rawData = make([]byte, len(data))
    copy(params.rawData, data)

    // Extract preset name (KEEP - offset 20-60)
    if len(data) >= 60 {
        // ... (existing name extraction code)
    }

    // NEW: Exact offset-based extraction
    extractBasicAdjustments(data, params)
    extractAdvancedAdjustments(data, params)
    extractColorBlender(data, params)
    extractColorGrading(data, params)
    extractToneCurve(data, params)

    // FALLBACK: Heuristic extraction for missing/invalid offsets
    // Only runs if exact extraction failed (parameters still zero)
    if params.sharpening == 0 && params.contrast == 0 && params.saturation == 0 {
        // Extract using legacy heuristics (KEEP for backward compatibility)
        rawParams := extractRawParamBytes(data)
        colorData := extractColorData(data)
        toneCurve := extractToneCurveData(data)
        estimateParameters(params, rawParams, colorData, toneCurve)
    }

    return params, nil
}
```

**Changes**:
- Add exact extraction function calls
- Keep fallback to heuristics if exact extraction returns zeros
- Move legacy extraction code to separate functions for clarity

### Step 4: Extend buildRecipe() Function

Map all 48 extracted parameters to UniversalRecipe:

```go
func buildRecipe(params *np3Parameters) (*models.UniversalRecipe, error) {
    builder := models.NewRecipeBuilder()

    builder.WithSourceFormat("np3")

    if params.name != "" {
        builder.WithName(params.name)
    }

    // Basic Adjustments (NEW: exact values, no scaling needed)
    if params.sharpening != 0 {
        // Map -3.0 to +9.0 → 0 to 150
        builder.WithSharpness(int((params.sharpening + 3.0) * 12.5))
    }
    if params.clarity != 0 {
        // Map -5.0 to +5.0 → -100 to +100
        builder.WithClarity(int(params.clarity * 20))
    }

    // Advanced Adjustments (NEW)
    if params.midRangeSharpening != 0 {
        builder.WithMidRangeSharpening(params.midRangeSharpening)
    }
    if params.contrast != 0 {
        builder.WithContrast(params.contrast)
    }
    if params.highlights != 0 {
        builder.WithHighlights(params.highlights)
    }
    if params.shadows != 0 {
        builder.WithShadows(params.shadows)
    }
    if params.whiteLevel != 0 {
        builder.WithWhites(params.whiteLevel)
    }
    if params.blackLevel != 0 {
        builder.WithBlacks(params.blackLevel)
    }
    if params.saturation != 0 {
        builder.WithSaturation(params.saturation)
    }

    // Color Blender (NEW - 8 colors × 3 values)
    if params.redHue != 0 || params.redChroma != 0 || params.redBrightness != 0 {
        builder.WithColorAdjustment("red", params.redHue, params.redChroma, params.redBrightness)
    }
    // ... (repeat for other 7 colors)

    // Color Grading (NEW)
    if params.highlightsZone.Hue != 0 || params.midtoneZone.Hue != 0 || params.shadowsZone.Hue != 0 {
        builder.WithColorGrading(
            params.highlightsZone,
            params.midtoneZone,
            params.shadowsZone,
            params.blending,
            params.balance,
        )
    }

    // Tone Curve (NEW)
    if params.toneCurvePointCount > 0 {
        builder.WithToneCurve(params.toneCurvePoints, params.toneCurveRaw)
    }

    // Build and validate
    recipe, err := builder.Build()
    if err != nil {
        return nil, fmt.Errorf("build recipe: %w", err)
    }

    // Preserve raw binary data (KEEP)
    if params.rawData != nil && len(params.rawData) > 0 {
        if recipe.FormatSpecificBinary == nil {
            recipe.FormatSpecificBinary = make(map[string][]byte)
        }
        recipe.FormatSpecificBinary["np3_raw"] = params.rawData
    }

    return recipe, nil
}
```

---

## Testing Strategy

### Test 1: Backward Compatibility

Ensure existing tests still pass without modification:

```bash
go test ./internal/formats/np3/... -v
```

**Expected**: All existing round-trip tests pass with identical behavior.

### Test 2: Exact Offset Accuracy

Create new tests for exact offset extraction:

```go
func TestExactOffsetExtraction(t *testing.T) {
    // Load real NP3 file
    data, err := os.ReadFile("testdata/sample.np3")
    require.NoError(t, err)

    // Parse with exact offsets
    recipe, err := Parse(data)
    require.NoError(t, err)

    // Verify parameters match expected values from TypeScript implementation
    assert.Equal(t, expectedSharpening, recipe.Sharpness)
    assert.Equal(t, expectedClarity, recipe.Clarity)
    // ... (verify all 48 parameters)
}
```

### Test 3: Fallback Behavior

Test that heuristics still work when exact offsets are invalid:

```go
func TestFallbackToHeuristics(t *testing.T) {
    // Create malformed NP3 with invalid offsets
    data := make([]byte, 500)
    copy(data[:3], []byte{'N', 'C', 'P'})
    // Fill with invalid data at exact offsets

    // Parse should fall back to heuristics
    recipe, err := Parse(data)
    require.NoError(t, err)

    // Verify heuristic estimation still works
    assert.NotNil(t, recipe)
}
```

---

## Implementation Schedule

### Day 1 (Today): Design + Basic Adjustments
- ✅ Design strategy document (this file)
- ⏳ Extend np3Parameters struct
- ⏳ Implement extractBasicAdjustments()
- ⏳ Add tests for basic adjustments
- ⏳ Run backward compatibility tests

**Estimated**: 2-3 hours

### Day 2: Advanced Adjustments
- Implement extractAdvancedAdjustments()
- Update buildRecipe() for 7 new parameters
- Add validation for 7 new parameters
- Test with real NP3 files

**Estimated**: 3-4 hours

### Day 3: Color Blender
- Implement extractColorBlender()
- Update buildRecipe() for 24 color parameters
- Add validation for color blender
- Test color extraction accuracy

**Estimated**: 3-4 hours

### Day 4: Color Grading
- Implement extractColorGrading()
- Update buildRecipe() for 11 color grading parameters
- Test Hue12 encoding accuracy
- Verify zone-based color grading

**Estimated**: 2-3 hours

### Day 5: Tone Curve
- Implement extractToneCurve()
- Handle variable-length arrays
- Handle 16-bit big-endian values
- Test with extended NP3 formats

**Estimated**: 3-4 hours

### Day 6: Integration + Testing
- Full integration testing
- Test with 10+ real NP3 files
- Validate 100% parameter extraction
- Benchmark performance improvement

**Estimated**: 4-5 hours

### Day 7: Documentation + Review
- Update parse.go documentation
- Update format specification
- Create migration guide
- Code review and cleanup

**Estimated**: 2-3 hours

---

## Success Criteria

✅ **Accuracy**: All 56 parameters extracted with 100% accuracy from standard NP3 files
✅ **Compatibility**: All existing tests pass without modification
✅ **Fallback**: Heuristic extraction still works for malformed files
✅ **Performance**: No regression in parsing speed (<5% overhead acceptable)
✅ **Documentation**: All new code thoroughly documented
✅ **Testing**: 400+ test cases covering all extraction functions

---

## Risk Mitigation

**Risk 1: Breaking Existing Behavior**
- **Mitigation**: Dual-mode approach with fallback to heuristics
- **Validation**: Run full test suite after each parameter addition

**Risk 2: File Size Variants**
- **Mitigation**: ValidateOffset() checks before reading
- **Validation**: Test with 392-byte, 480-byte, and 978-byte files

**Risk 3: Encoding Errors**
- **Mitigation**: All encoding functions have 200+ round-trip tests
- **Validation**: Compare against TypeScript implementation results

**Risk 4: Performance Regression**
- **Mitigation**: Exact offset reads are faster than heuristic analysis
- **Validation**: Benchmark before/after implementation

---

**Document Status**: Design Complete - Ready for Implementation
**Last Updated**: 2025-11-07
**Next Review**: After Day 1 implementation (Basic Adjustments)
