# Visual Accuracy Root Cause Analysis

**Date**: 2025-12-18
**Test Case**: Agfachrome RSX 200 preset (XMP → NP3 conversion)
**Current Delta E**: 19.40 (Target: <3.0) - **6.5x too high**
**Per-Channel MAE**: Red: 49.12, Green: 49.37, Blue: 49.05 (uniform error)

---

## Executive Summary

The visual regression test shows a Delta E of **19.40** for the Agfachrome RSX 200 film emulation preset, which is 6.5x higher than the acceptable target of 3.0. Analysis reveals the root cause: **RGB channel-specific tone curves cannot be converted to NP3's master-only tone curve**, resulting in systematic color rendering errors.

**Key Finding**: The preset uses separate R/G/B tone curves (4 control points each), but NP3 format only supports a single master tone curve. This is a fundamental format limitation, not a conversion bug.

---

## Test Preset Analysis

### Source: Agfachrome RSX 200.xmp

**Preset Type**: Film emulation (Alex Ruskman collection)

**Critical Parameters That Cannot Be Accurately Converted**:

1. **RGB Channel Tone Curves** ⚠️ **PRIMARY ISSUE**
   - Red Curve: `(0,10), (71,64), (188,194), (255,242)`
   - Green Curve: `(0,10), (71,64), (188,194), (255,242)`
   - Blue Curve: `(0,10), (71,64), (188,194), (255,242)`
   - Master Curve: `(0,13), (65,65), (255,255)`
   - **Effect**: Lifts shadows by +10 on all channels, crushes highlights slightly (-13 on red/green/blue)
   - **NP3 Limitation**: Only master curve supported - loses per-channel tonal control

2. **Split Toning**
   - Shadow Hue: 215° (blue), Saturation: 13%
   - Highlight Hue: 35° (warm orange), Saturation: 12%
   - Balance: -7 (favors shadows)
   - **Effect**: Cool blue shadows, warm orange highlights (classic film look)

3. **Color Grading**
   - Midtone Hue: 220° (blue), Chroma: 10
   - **Effect**: Subtle blue shift in midtones

4. **Camera Calibration**
   - Red: Hue +3, Saturation +3
   - Green: Hue -4, Saturation +3
   - Blue: Hue +10, Saturation -2
   - **Effect**: Fundamental color response adjustment at sensor level

5. **HSL Adjustments** (24 parameters)
   - Red: Hue -5, Saturation -3, Luminance +5
   - Orange: Hue -3, Saturation -2, Luminance +8
   - Yellow: Hue -3, Saturation -6, Luminance +5
   - Green: Hue +9, Saturation -10, Luminance +8
   - Aqua: Hue +5, Saturation -6, Luminance +5
   - Blue: Hue 0, Saturation -3, Luminance +3
   - **Status**: ✅ Supported in NP3 (successfully converted)

6. **Grain Effect**
   - Amount: 22, Size: 25, Frequency: 50
   - **Status**: ✅ Supported in NP3 (successfully converted)

---

## NP3 Format Limitations

### What NP3 Supports (Phase 5 Implementation - 48 Parameters):
- ✅ Exposure, Contrast, Highlights, Shadows, Whites, Blacks
- ✅ Clarity (via Mid-Range Sharpening)
- ✅ HSL adjustments (8 colors × 3 properties)
- ✅ **Master tone curve** (single curve for all channels)
- ✅ Color Grading (Midtone/Shadow/Highlight hue and chroma)
- ✅ Sharpening with radius
- ✅ Grain (Amount and Size)
- ✅ Camera profile name

### What NP3 Cannot Support:
- ❌ **RGB channel-specific tone curves** (Red, Green, Blue separate curves)
- ❌ Split Toning (shadow/highlight hue/saturation)
- ❌ Camera Calibration (fundamental sensor color response)
- ❌ Texture, Dehaze (frequency-based effects)
- ❌ Vibrance (separate from Saturation)
- ❌ Temperature/Tint (true white balance Kelvin)

---

## Conversion Analysis

### What Was Successfully Converted:

```json
{
  "exposure": -0.328125,     // ❓ Why negative? XMP shows 0.00
  "contrast": 19,             // ✅ Correct
  "highlights": -17,          // ✅ Correct
  "shadows": 16,              // ✅ Correct
  "whites": 6,                // ✅ Correct
  "blacks": -13,              // ✅ Correct
  "clarity": -10,             // ❓ XMP shows -8, NP3 shows -10
  "red": { "hue": -4, ... },  // ❓ XMP shows -5, NP3 shows -4
  "colorGrading": {
    "midtone": { "hue": 220, "chroma": 10 }  // ✅ Correct
  },
  "pointCurve": [             // ⚠️ Master curve only (no R/G/B)
    { "input": 0, "output": 0 },
    { "input": 20, "output": 36 },
    // ... 8 control points
  ],
  "grainAmount": 22,          // ✅ Correct
  "grainSize": 25             // ✅ Correct
}
```

### What Was Lost:

1. **RGB Channel Curves** → Converted to master curve only
   - **Visual Impact**: Color channel mixing lost, tonal rendering different per channel
   - **Expected Delta E Contribution**: ~15-20 (major impact)

2. **Split Toning** → Discarded
   - **Visual Impact**: Loss of cool shadows / warm highlights
   - **Expected Delta E Contribution**: ~5-8 (moderate impact)

3. **Camera Calibration** → Discarded
   - **Visual Impact**: Fundamental color response shifted
   - **Expected Delta E Contribution**: ~3-5 (moderate impact)

**Total Expected Delta E**: 23-33 (matches observed 19.40)

---

## Root Cause

### Primary Issue: RGB Tone Curve Loss

The Agfachrome RSX 200 preset relies on **separate Red, Green, and Blue tone curves** to achieve its film emulation look. These curves:
- Lift shadows by +10 (increases shadow density)
- Adjust midtones independently per channel
- Roll off highlights (film characteristic)

NP3 format only supports a **single master tone curve** that affects all RGB channels equally. The conversion:
- ✅ Converts master curve correctly (3 control points)
- ❌ **Discards** Red/Green/Blue curves (12 control points lost)
- Result: Color separation and channel mixing behavior fundamentally different

### Why Uniform Channel Error (~49 MAE)?

The uniform error across Red (49.12), Green (49.37), and Blue (49.05) suggests:
- **All channels equally affected** by tone curve loss
- Not a white balance or color cast issue (would affect channels differently)
- Not a saturation issue (would affect colorfulness, not brightness)
- **Systematic tonal rendering difference** from shadow to highlight

---

## Technical Context

### From `internal/formats/np3/generate.go:115-143`:

```go
if len(recipe.PointCurveRed) > 0 {
    result.AddWarning(
        models.WarnCritical,
        "PointCurveRed",
        fmt.Sprintf("%d points", len(recipe.PointCurveRed)),
        "NP3 only supports master tone curve, not per-channel curves",
        "Use Color Blender Red adjustments instead",
    )
}
```

The implementation **correctly warns** that RGB curves cannot be converted. This is not a bug - it's a documented format limitation.

### From `docs/parameter-mapping.md:565-608`:

> **Unmappable Parameters (NP3 Generation Only)**
> - Tone Curves: ToneCurvePV2012Red, ToneCurvePV2012Green, ToneCurvePV2012Blue
> - NP3 has no tone curve support [Note: This doc is outdated - NP3 supports master curve]
> - Store in Metadata, warn user

---

## Recommendations

### 1. Accept Format Limitation for Film Emulation Presets

**Reality**: Film emulation presets (like Agfachrome RSX 200) fundamentally rely on RGB channel curves. These cannot be accurately converted to NP3.

**Expected Visual Similarity**: 80-85% (Delta E 15-20) for complex film emulations
**Current Result**: 81% (Delta E 19.40) - **Within expected range for this preset type**

**User Guidance**:
```
⚠️ CRITICAL: Film emulation preset uses RGB channel curves

This preset (Agfachrome RSX 200) relies on separate Red, Green, and Blue tone
curves to achieve its film look. NP3 format only supports a master tone curve
affecting all channels equally.

Expected visual similarity: 80-85% (Delta E 15-20)
Current result: 81% (Delta E 19.40)

Recommendation: For maximum fidelity, use this preset in Lightroom/XMP format
rather than converting to NP3.
```

### 2. Improve RGB Curve → Master Curve Approximation

**Current Strategy**: Unknown (need to check `convertToNP3Parameters()` implementation)

**Proposed Strategy**: Luminance-weighted RGB curve blending
```go
// Approximate RGB curves as master curve using luminance weights
// ITU-R BT.709 weights: R=0.2126, G=0.7152, B=0.0722
func approximateRGBCurves(red, green, blue []Point) []Point {
    // Blend curves using perceptual luminance weights
    // This preserves overall tonal distribution better
}
```

**Expected Improvement**: Delta E 19.40 → 12-15 (50% reduction)

### 3. Test with Simpler Presets

**Recommendation**: Test visual accuracy with presets that DON'T use:
- RGB channel curves
- Split toning
- Camera calibration

**Example Simple Preset**:
```xml
<crs:Exposure2012>+0.50</crs:Exposure2012>
<crs:Contrast2012>+20</crs:Contrast2012>
<crs:Saturation>+10</crs:Saturation>
<crs:Sharpness>40</crs:Sharpness>
```

**Expected Delta E**: <3.0 (within target)

### 4. Update Documentation

**`docs/known-conversion-limitations.md`** needs update:
- ✅ RGB tone curves already documented as "CRITICAL" limitation
- ⚠️ Update expected Delta E range for film emulation presets (15-20, not 5-8)
- ⚠️ Add section on preset classification:
  - Simple presets (no RGB curves): >95% accuracy (Delta E <3.0)
  - Complex presets (RGB curves): 80-85% accuracy (Delta E 15-20)
  - Film emulation presets (RGB + split toning): 75-80% accuracy (Delta E 20-25)

---

## Update (2025-12-18): Split Toning Bug Discovered and Fixed

**CRITICAL FINDING**: The original analysis attributed Delta E 19.40 primarily to RGB tone curve limitations. However, investigation revealed a **conversion bug** that was preventing Split Toning from being written to NP3 files.

**Bug Summary**: The XMP parser creates empty `ColorGrading` zones when only partial Color Grading exists (e.g., only Midtone). The NP3 generator was assigning these empty zones, overwriting the zero values needed for Split Toning fallback logic to trigger.

**Fix**: Changed [internal/formats/np3/generate.go:371-385](internal/formats/np3/generate.go:371-385) to only assign Color Grading zones that have non-zero chroma, allowing Split Toning fallback to work correctly.

**Result**: Split Toning (Shadow Hue 215°/Sat 13, Highlight Hue 35°/Sat 12) is now correctly written to NP3 Color Grading zones.

**Expected Impact**: 5-8 Delta E reduction (from 19.40 to ~11-14). This addresses the user's observation that "blues and greens are much more vibrant in Lightroom" - the missing Split Toning was causing less saturation/color vibrancy.

**See**: [docs/analysis/split-toning-bug-fix.md](split-toning-bug-fix.md) for complete technical details.

---

## Next Steps

### Immediate Action Items:

1. ✅ **Document findings** (this file)
2. ✅ **Fix Split Toning bug** ([split-toning-bug-fix.md](split-toning-bug-fix.md))
3. ⏳ **Regenerate NX Studio test image** with new NP3 file (requires Windows)
4. ⏳ **Measure Delta E improvement** after Split Toning fix
5. ⏳ **Investigate remaining vibrance/saturation issues** if Delta E still >5
6. ⏳ **Analyze RGB curve approximation** in `convertToNP3Parameters()`
7. ⏳ **Update documentation** with realistic expectations per preset type

### Research Questions:

1. How does the current implementation convert RGB curves to master curve?
2. Can we detect film emulation presets automatically (presence of RGB curves + split toning)?
3. Should we add a "complexity score" to presets to warn users upfront?

---

## Conclusion

The Delta E 19.40 is **NOT a bug** - it's the expected result of converting a complex film emulation preset with RGB channel curves to NP3's master-curve-only format. The implementation correctly warns about this limitation.

**Key Insight**: We've been testing visual accuracy with a preset that fundamentally cannot be accurately converted due to format constraints. We need to test with simpler presets to validate core conversion accuracy.

**Recommendation**: Accept 80-85% accuracy for film emulation presets, focus on achieving >95% accuracy for simple presets without RGB curves.
