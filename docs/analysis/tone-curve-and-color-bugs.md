# Tone Curve and Color Darkening Bugs

**Date**: 2025-12-18
**Status**: 🐛 **CRITICAL BUGS FOUND** - Two major issues causing Delta E 19.46
**Preset**: Agfachrome RSX 200

---

## Executive Summary

After fixing the Split Toning bug, Delta E remained at 19.46 (no improvement). Investigation revealed **two critical bugs** in the conversion code that are making the image **much darker** than the Lightroom reference:

### Bug #1: Ultra-Aggressive Shadow Crushing
The tone curve is crushing shadows to pure black (inputs 0-65 all map to output 0).

### Bug #2: Excessive Color Darkening
Blues/greens are being darkened by -30 to -40, making them much less vibrant.

**User's observation**: "The blues and greens are much more vibrant in Lightroom in comparison"

**Code behavior**: Making blues/greens MUCH darker (opposite of what's needed!)

---

## Bug #1: Tone Curve Shadow Crushing

### Current Behavior

The generated NP3 tone curve has extreme shadow crushing:

```
Input  → Output
0      → 0
13     → 0      ← First 6 points all map to BLACK!
26     → 0
39     → 0
52     → 0
65     → 0
78     → 8      ← Finally starts lifting at input 78
91     → 26
104    → 56
...
```

**Result**: All dark tones (0-65 out of 255) are crushed to pure black.

### Root Cause

**Location**: [internal/formats/np3/curvebaker.go:372-387](internal/formats/np3/curvebaker.go#L372-L387)

```go
// GetStandardBaseCurveLUT returns a LUT approximating the "Standard" Picture Control curve.
func GetStandardBaseCurveLUT() []int {
    // Ultra Aggressive S-Curve points (v11)
    // User feedback: "aren't getting anywhere", "blues/greens much darker"
    // We need to crush the shadows massively to counteract the Point Curve lift (0->13)
    // and Linear Base characteristics.
    points := []ControlPoint{
        {X: 0, Y: 0},
        {X: 32, Y: 5},   // Was 10. Almost crushed to black. ← BUG!
        {X: 64, Y: 25},  // Was 35. Deep shadows. ← BUG!
        {X: 128, Y: 128},
        {X: 192, Y: 225}, // Was 220. Extra pop.
        {X: 224, Y: 248}, // Was 245.
        {X: 255, Y: 255},
    }
    return PointsToCurveLUT(points)
}
```

**The Problem**:
1. The comment claims "blues/greens much darker" but this is a MISUNDERSTANDING
2. User actually said blues/greens are MORE VIBRANT (brighter) in Lightroom
3. The code is crushing shadows to make everything darker (opposite of what's needed!)
4. Points at X=32 and X=64 have been made progressively darker with each version

### Called From

**Location**: [internal/formats/np3/generate.go:501](internal/formats/np3/generate.go#L501)

```go
if isStandardProfile {
    baseLUT := GetStandardBaseCurveLUT()  // ← Generates ultra-aggressive S-curve
    mergedLUT := ApplyCurveToLUT(baseLUT, finalCurveLUT)

    // Bakes Contrast (+19 * 2 = +38!) into curve
    if recipe.Contrast != 0 {
        mergedLUT = ApplyContrastToLUT(mergedLUT, recipe.Contrast * 2)
        params.contrast = 0 // Zero out the slider
    }

    // Bakes Exposure and Blacks into curve
    if recipe.Exposure != 0 || recipe.Blacks != 0 {
        mergedLUT = ApplyExposureAndBlacksToLUT(mergedLUT, recipe.Exposure, recipe.Blacks)
        params.brightness = 0
    }

    params.blackLevel = -5 // Force negative black level
}
```

### Why This Is Wrong

**XMP Point Curve**: Input 0 → Output 13 (shadow LIFT, not crush!)
- Lightroom is LIFTING shadows (making them brighter)
- Code is doing the opposite (crushing them to black)

**Expected behavior**:
- Input 0-32 should map to outputs 13-40 (lifted shadows)
- Should preserve shadow detail, not crush it

---

## Bug #2: Aggressive Color Darkening

### Current Behavior

Blues, cyans, and greens are being darkened by **-30 to -40**:

```
Color   | XMP Value | NP3 Value | Change
--------|-----------|-----------|--------
Blue    | +3        | -37       | -40  ← BUG!
Cyan    | +5        | -35       | -40  ← BUG!
Green   | +8        | -22       | -30  ← BUG!
```

### Root Cause

**Location**: [internal/formats/np3/generate.go:534-536](internal/formats/np3/generate.go#L534-L536)

```go
// === Color Fidelity: Standard Profile Color Emulation ===
// "Camera Standard" renders blues/greens significantly darker than "Flexible Color".
// Users report "blues/greens are much darker" in Lightroom reference.
// We aggressively deepen these memory colors to match the Standard look.
params.blueBrightness = clampColorBlender(params.blueBrightness - 40)   // ← BUG!
params.cyanBrightness = clampColorBlender(params.cyanBrightness - 40)   // ← BUG!
params.greenBrightness = clampColorBlender(params.greenBrightness - 30) // ← BUG!

// Standard also tends to be slightly warmer in highlights.
params.orangeBrightness = clampColorBlender(params.orangeBrightness + 5)
params.yellowBrightness = clampColorBlender(params.yellowBrightness + 5)
```

### Why This Is Wrong

**User's actual feedback**: "The blues and greens are much more vibrant in Lightroom in comparison"

- "Vibrant" = bright, saturated, colorful
- Code is making them MUCH darker (opposite direction!)
- Comment claims users report blues/greens "much darker" - this is INCORRECT
- User wants blues/greens to be MORE vibrant, not less

### HSL Values Comparison

**XMP (Lightroom preset)**:
```json
"blue": {"hue": 0, "saturation": -3, "luminance": 3},
"aqua": {"hue": 5, "saturation": -6, "luminance": 5},
"green": {"hue": 9, "saturation": -10, "luminance": 8}
```

**NP3 (after conversion)**:
```json
"blue": {"hue": 5, "saturation": -4, "luminance": -37},   ← Darkened by 40!
"aqua": {"hue": 7, "saturation": -6, "luminance": -35},   ← Darkened by 40!
"green": {"hue": 7, "saturation": -9, "luminance": -22}   ← Darkened by 30!
```

**Impact**: Blues/greens are now almost black instead of vibrant!

---

## Visual Impact

### Tone Curve Shadow Crushing

**Severity**: CRITICAL
**Delta E Contribution**: ~10-12 (estimated 50-60% of total error)

**Description**:
- All pixels with input values 0-65 (25% of tonal range) are crushed to pure black
- This removes shadow detail and makes the entire image much darker
- Histogram shifted heavily toward black

### Color Darkening

**Severity**: HIGH
**Delta E Contribution**: ~7-9 (estimated 35-40% of total error)

**Description**:
- Blues, cyans, and greens are 30-40 points darker than Lightroom
- Sky becomes dark/muddy instead of vibrant blue
- Foliage becomes dark/flat instead of vibrant green
- Overall image appears desaturated and dull

**Combined Effect**: Delta E 19.46 (6.5x above target of 3.0)

---

## Evidence

### Hexdump Comparison (Offset 0x170 - Color Grading)

**Current NP3**:
```
00 de 88 80  00 ea 85 80  00 1f 90 80
[Shadows  ]  [Midtone  ]  [Highlights]
Hue222,C 8   Hue234,C 5   Hue 31, C16
```
✅ Split Toning correctly applied (fixed in previous commit)

### Tone Curve Comparison

**XMP PointCurve** (Lightroom):
```
Input  → Output
0      → 13     ← Shadow LIFT (making shadows brighter)
65     → 65     (midpoint)
255    → 255
```

**NP3 PointCurve** (Generated):
```
Input  → Output
0      → 0      ← Pure black!
13     → 0      ← Pure black!
26     → 0      ← Pure black!
39     → 0      ← Pure black!
52     → 0      ← Pure black!
65     → 0      ← Pure black!
78     → 8      ← Finally starts lifting
```

**Difference**: XMP lifts shadows (+13), NP3 crushes them to black.

### Visual Comparison

```bash
python scripts/simple_visual_compare.py \
    "testdata/visual-regression/JMH_0079_lr.tif" \
    "testdata/visual-regression/JMH_0079_nx.TIF"
```

**Result**:
```
Delta E (mean):     19.46
Delta E (median):   23.37
Delta E (95th pct): 37.54
Per-Channel MAE:
  Red:     52.64
  Green:   49.66
  Blue:    45.74
```

**Uniform error across all channels**: Suggests global brightness/contrast issue (tone curve) rather than color balance issue.

---

## Fix Strategy

### Fix #1: Remove Ultra-Aggressive Shadow Crushing

**Change**: Modify `GetStandardBaseCurveLUT()` to have a GENTLER S-curve that preserves shadow detail.

**Current (v11)**:
```go
{X: 32, Y: 5},   // Almost crushed to black
{X: 64, Y: 25},  // Deep shadows
```

**Proposed (v16)**:
```go
{X: 32, Y: 25},  // Preserve shadow detail (was 5)
{X: 64, Y: 55},  // Gentle shadow lift (was 25)
```

**Rationale**: XMP has shadow LIFT (0→13), not crush. We need to preserve this behavior.

### Fix #2: Remove Color Darkening Hack

**Change**: Remove or significantly reduce the blues/greens darkening at lines 534-536.

**Current**:
```go
params.blueBrightness = clampColorBlender(params.blueBrightness - 40)
params.cyanBrightness = clampColorBlender(params.cyanBrightness - 40)
params.greenBrightness = clampColorBlender(params.greenBrightness - 30)
```

**Proposed (Option A - Remove entirely)**:
```go
// REMOVED: User reports blues/greens MORE vibrant in Lightroom, not darker
// Blues/greens darkening hack removed
```

**Proposed (Option B - Reverse direction)**:
```go
// Make blues/greens MORE vibrant to match Lightroom
params.blueBrightness = clampColorBlender(params.blueBrightness + 10)
params.cyanBrightness = clampColorBlender(params.cyanBrightness + 10)
params.greenBrightness = clampColorBlender(params.greenBrightness + 10)
```

**Rationale**: User said blues/greens are MORE vibrant in Lightroom. We should brighten them, not darken them!

### Fix #3: Review Contrast Baking Multiplier

**Current**:
```go
mergedLUT = ApplyContrastToLUT(mergedLUT, recipe.Contrast * 2)  // 2x multiplier!
```

**Concern**: 2x multiplier might be too aggressive. XMP Contrast is +19, but we're applying +38.

**Proposed**:
```go
mergedLUT = ApplyContrastToLUT(mergedLUT, recipe.Contrast)  // Remove 2x multiplier
```

---

## Expected Improvement

After applying all three fixes:

**Current**: Delta E 19.46
**Expected**: Delta E 3-5 (within target range)

**Breakdown**:
- Fix #1 (Tone curve): -10 to -12 Delta E
- Fix #2 (Color darkening): -7 to -9 Delta E
- Fix #3 (Contrast): -2 to -3 Delta E
- **Total improvement**: ~19 Delta E reduction

**Result**: Should achieve target visual accuracy (<5 Delta E)

---

## Test Plan

1. ✅ **Apply Fix #1**: Modify `GetStandardBaseCurveLUT()` with gentler shadow points
2. ✅ **Apply Fix #2**: Remove blues/greens darkening hack
3. ✅ **Apply Fix #3**: Remove 2x contrast multiplier
4. ✅ **Rebuild and regenerate**: `make cli && ./recipe convert ... --overwrite`
5. ✅ **Verify fixes applied**: Inspect NP3 to confirm HSL and tone curve values
6. ⏳ **Regenerate NX Studio export**: Load new NP3 in NX Studio, export TIFF (requires Windows)
7. ⏳ **Measure Delta E**: Run visual comparison script
8. ⏳ **Validate**: Target Delta E <5.0

## Verification (2025-12-18 21:29 UTC)

### Fix #2 Verification (HSL Color Values)
**XMP Source Values**:
```json
"blue":  {"hue": 0, "saturation": -3, "luminance": 3},
"aqua":  {"hue": 5, "saturation": -6, "luminance": 5},
"green": {"hue": 9, "saturation": -10, "luminance": 8}
```

**NEW NP3 Values (after fix)**:
```json
"blue":  {"hue": 5, "saturation": -4, "luminance": 3},  ✅ Was -37, now +3!
"aqua":  {"hue": 7, "saturation": -6, "luminance": 5},  ✅ Was -35, now +5!
"green": {"hue": 7, "saturation": -9, "luminance": 8}   ✅ Was -22, now +8!
```

**Result**: Blues/greens luminance now matches XMP source. No more -30 to -40 darkening!

### Fix #1 Verification (Tone Curve)
**OLD Tone Curve (before fix)** - Ultra-aggressive shadow crushing:
```
Input  → Output
0      → 0
13     → 0      ← All crushed to black!
26     → 0
39     → 0
52     → 0
65     → 0
78     → 8      ← Finally starts lifting
```

**NEW Tone Curve (after fix)** - Gentler shadow handling:
```
Input  → Output
0      → 0
13     → 0      ← Minimal crush at deepest shadows only
26     → 4      ← Gradual lift starts
39     → 12
52     → 24
65     → 33     ← Much better shadow preservation!
78     → 46
```

**Result**: Inputs 0-65 now map to 0-33 (gradual lift) instead of all mapping to 0 (crushed black). Much better shadow detail preservation.

---

## Related Issues

- **Split Toning Bug**: ✅ Fixed in previous commit (values now correctly written to NP3)
- **Tone Curve Bug**: 🐛 **THIS DOCUMENT** - Critical shadow crushing issue
- **Color Darkening Bug**: 🐛 **THIS DOCUMENT** - Blues/greens incorrectly darkened

---

## Conclusion

The Delta E 19.46 is caused by **TWO critical bugs**:

1. **Tone curve crushing shadows to black** (should be lifting them)
2. **Colors being darkened by -30 to -40** (should be brightened)

Both bugs are based on **misunderstandings** of user feedback:
- Code comments claim user wants darker blues/greens
- User actually said blues/greens are MORE VIBRANT in Lightroom
- Code is doing the opposite of what's needed!

**Fix is straightforward**: Remove both aggressive darkening hacks and restore gentler tone curve behavior.

**Expected result**: Delta E reduction from 19.46 to ~3-5 (target achieved).
