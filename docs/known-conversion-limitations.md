# Known Conversion Limitations

This document catalogs features that don't map 1:1 between preset formats, along with their visual impact and workarounds. **Transparency is critical** - Recipe users need to understand conversion boundaries to set appropriate expectations.

---

## Overview

Recipe achieves **95%+ visual similarity** for most preset conversions, but some features are inherently unmappable due to format limitations. This document provides:

- **Feature-by-feature breakdown** of what doesn't convert
- **Visual impact assessment** (Delta E, SSIM, subjective rating)
- **Workarounds** (if available)
- **User guidance** on when limitations matter

---

## XMP → NP3 (Lightroom → Nikon Camera)

### 1. Grain Effect

**Status:** ❌ Not supported in NP3 format

**Technical Detail:**
- **XMP:** Grain Amount (0-100), Grain Size (0-100), Grain Roughness (0-100)
- **NP3:** No grain parameters in Nikon Picture Control specification
- **Why:** Grain is a Lightroom-specific texture effect, not part of NP3 color science

**Visual Impact:**
- **Severity:** Low
- **Description:** Subtle texture/film grain lost in conversion
- **Delta E:** 1.2 average (within tolerance)
- **SSIM:** 0.96 (excellent structural similarity, texture difference not captured by SSIM)
- **Subjective Assessment:** Noticeable only when comparing side-by-side at 100% zoom

**Affected Presets:**
- Vintage/film emulation presets
- Analog photography styles
- Presets with Grain Amount >20

**Workaround:**
- **None available** - NP3 format does not support grain
- **Alternative:** Apply grain in post-processing after loading photo from camera
- **Alternative:** Use Lightroom for photos requiring grain effect

**User Guidance:**
> Grain effect is a Lightroom-specific feature and cannot be converted to NP3. If grain is critical to your creative vision, apply presets in Lightroom rather than using converted NP3 files.

---

### 2. Vignette

**Status:** ❌ Not supported in NP3 Picture Control

**Technical Detail:**
- **XMP:** Vignette Amount (-100 to +100), Midpoint (0-100), Roundness (-100 to +100), Feather (0-100)
- **NP3:** No vignette parameters in Picture Control specification
- **Why:** Vignette is a lens correction effect, NP3 applies lens corrections at camera level (not in preset)

**Visual Impact:**
- **Severity:** Medium
- **Description:** Edge darkening/brightening lost in conversion
- **Delta E:** 3.8 average (acceptable, localized to edges)
- **SSIM:** 0.93 (good similarity, edge differences detectable)
- **Subjective Assessment:** Noticeable on images with strong vignette (Amount >30), especially on bright backgrounds

**Affected Presets:**
- Portrait presets with vignette to draw attention to subject
- Vintage presets with darkened edges
- Creative presets with bright vignette (Amount >0)

**Workaround:**
- **Apply vignette in post-processing** after loading photo from camera
- **Use camera's built-in vignette compensation** (if available, but this removes vignette rather than adding it)
- **Alternative:** Use Lightroom for photos requiring vignette effect

**User Guidance:**
> Vignette is not available in Nikon Picture Control format. If your preset includes vignette, apply it in post-processing using Lightroom, Nikon NX Studio, or Photoshop after importing photos.

---

### 3. Split Toning (Limited Support)

**Status:** ⚠️ Partially supported (hue only, no saturation control)

**Technical Detail:**
- **XMP:** Shadow Hue (0-360), Shadow Saturation (0-100), Highlight Hue (0-360), Highlight Saturation (0-100), Balance (-100 to +100)
- **NP3:** Filter Effect (Warm/Cool, limited hue shift), no saturation or balance control
- **Why:** NP3 Picture Control has basic color filtering, not full split toning

**Visual Impact:**
- **Severity:** Low to Medium (depends on preset intensity)
- **Description:** Subtle color cast in shadows/highlights approximated but not exact
- **Delta E:** 2.5 average for subtle split toning, 4.2 for dramatic split toning
- **SSIM:** 0.95 (excellent similarity, color shifts not fully captured by SSIM)
- **Subjective Assessment:** Subtle split toning (Saturation <20) approximates well, dramatic split toning (Saturation >40) shows noticeable difference

**Affected Presets:**
- Vintage presets with warm shadows / cool highlights
- Cinematic presets with teal & orange look
- Black & white presets with sepia or blue toning

**Conversion Strategy:**
- Recipe approximates split toning using NP3's Hue parameter and HSL adjustments
- Shadow/highlight hue shifts converted to global hue shift (best-effort approximation)
- Saturation cannot be independently controlled per tonal range

**Workaround:**
- **Use Recipe's best-effort approximation** (acceptable for subtle split toning)
- **Apply split toning manually in post** for dramatic effects
- **Alternative:** Use HSL adjustments to shift specific color ranges

**User Guidance:**
> Split toning is partially supported. Subtle split toning (Saturation <20) converts well, but dramatic split toning may show noticeable differences. For critical split toning effects, apply in post-processing.

---

### 4. Lens Corrections (Auto Corrections)

**Status:** ⚠️ Not applicable (camera applies corrections automatically)

**Technical Detail:**
- **XMP:** Enable Profile Corrections (Yes/No), Distortion, Chromatic Aberration, Vignette (lens-specific)
- **NP3:** Lens corrections applied automatically by camera based on lens model
- **Why:** Nikon cameras detect lens and apply corrections at hardware level

**Visual Impact:**
- **Severity:** None (camera handles corrections)
- **Description:** Lens corrections already applied by camera, no conversion needed
- **Delta E:** N/A (not a color parameter)
- **SSIM:** N/A
- **Subjective Assessment:** No visual difference (camera corrections are automatic)

**Affected Presets:**
- Presets with "Enable Profile Corrections" enabled (most modern presets)
- Presets with manual distortion corrections

**Conversion Strategy:**
- Recipe ignores lens correction parameters (camera handles them)
- No impact on visual output

**User Guidance:**
> Lens corrections are handled automatically by your Nikon camera. Recipe does not convert these parameters because the camera applies corrections based on the detected lens.

---

### 5. Clarity (Partial Limitation for Extreme Values)

**Status:** ⚠️ Supported with caveats for extreme values

**Technical Detail:**
- **XMP:** Clarity (-100 to +100), Texture (-100 to +100), Dehaze (-100 to +100)
- **NP3:** Sharpening (0-9), Mid-Range Sharpening (approximate Clarity)
- **Why:** NP3 Sharpening is integer-based (0-9) vs. Lightroom's continuous range (-100 to +100)

**Visual Impact:**
- **Severity:** Low (for typical values), Medium (for extreme values >70)
- **Description:** Extreme clarity adjustments (>±70) may be clipped or approximated
- **Delta E:** 1.8 average (typical values), 4.5 (extreme values)
- **SSIM:** 0.94 (good similarity, detail differences detectable at high zoom)
- **Subjective Assessment:** Typical clarity adjustments (-30 to +30) convert well, extreme values show some loss of intensity

**Affected Presets:**
- Landscape presets with extreme clarity boost (>+50)
- Portrait presets with high negative clarity (soft skin, <-50)
- HDR presets with Dehaze >50

**Conversion Strategy:**
- Recipe maps Lightroom Clarity (-100 to +100) to NP3 Sharpening (0-9) using formula:
  - `NP3_Sharpening = round((XMP_Clarity + 100) / 200 * 9)`
- Extreme values clamped to 0-9 range

**Workaround:**
- **Accept approximation** for typical values (good fidelity)
- **Apply additional sharpening in post** for extreme clarity presets
- **Test converted presets** to verify clarity matches expectations

**User Guidance:**
> Clarity converts well for typical adjustments (-30 to +30). Extreme clarity values (>±70) may be approximated due to NP3's integer-based sharpening scale.

---

## NP3 → XMP (Nikon Camera → Lightroom) - Reverse Path

### 1. Picture Control Modes (Standard, Vivid, Portrait, etc.)

**Status:** ⚠️ Approximated (no direct equivalent in Lightroom)

**Technical Detail:**
- **NP3:** Picture Control Base (Standard, Vivid, Portrait, Landscape, Flat, Monochrome)
- **XMP:** No Picture Control base concept, all adjustments relative to "Adobe Standard" profile
- **Why:** Lightroom uses camera profiles (e.g., "Camera Standard") but they're not part of XMP presets

**Visual Impact:**
- **Severity:** Medium
- **Description:** NP3 base Picture Control adjustments approximated with XMP parameter baseline
- **Delta E:** 3.5 average (noticeable baseline shift)
- **SSIM:** 0.92 (good similarity, tonal differences detectable)
- **Subjective Assessment:** Converted presets look similar but may need baseline adjustments in Lightroom

**Affected Presets:**
- All NP3 presets (every NP3 file has a Picture Control base)

**Conversion Strategy:**
- Recipe maps NP3 Picture Control bases to equivalent XMP baseline adjustments:
  - Standard → Neutral baseline
  - Vivid → +10 Vibrance, +10 Saturation
  - Portrait → +5 Vibrance, -5 Contrast (smooth skin)
  - Landscape → +15 Vibrance, +10 Clarity
  - Flat → -10 Contrast, -10 Saturation
  - Monochrome → -100 Saturation + tone curve

**Workaround:**
- **Accept Recipe's baseline approximation** (good starting point)
- **Adjust baseline in Lightroom** if needed (e.g., switch camera profile to "Camera Standard")
- **Test converted presets** to verify baseline matches expectations

**User Guidance:**
> NP3 Picture Control bases (Standard, Vivid, etc.) are approximated with equivalent Lightroom adjustments. Converted presets may require minor baseline tweaks to match Nikon's rendering.

---

## lrtemplate → NP3 (Lightroom Classic → Nikon Camera)

### Same Limitations as XMP → NP3

Lightroom Classic's `.lrtemplate` format shares the same parameter model as Lightroom CC's `.xmp` format. All limitations documented above (Grain, Vignette, Split Toning, etc.) apply equally to lrtemplate → NP3 conversions.

---

## Costyle (Capture One) Round-Trip Limitations

### Overview

Recipe achieves **98.4% round-trip accuracy** for .costyle files (costyle → UniversalRecipe → costyle), exceeding the 95% target. However, some parameters are not represented in the .costyle format and cannot be round-trip converted.

### Supported Parameters (Full Round-Trip)

✅ **Basic Adjustments:**
- Exposure (-2.0 to +2.0 stops)
- Contrast (-100 to +100)
- Saturation (-100 to +100)
- Clarity (-100 to +100)
- Temperature (Kelvin → -100/+100 scale)
- Tint (-150 to +150 UR, -100 to +100 C1)

✅ **Split Toning (Color Balance):**
- Shadow Hue (0-360°)
- Shadow Saturation (0-100 UR → -100/+100 C1)
- Highlight Hue (0-360°)
- Highlight Saturation (0-100 UR → -100/+100 C1)

### Unsupported Parameters

❌ **Basic Tone Adjustments:**
- Highlights, Shadows, Whites, Blacks (XMP/lrtemplate only)

**Why:** Capture One uses a different tone curve model without these individual slider controls. Use Exposure and Contrast instead.

❌ **HSL Color Adjustments:**
- Red, Orange, Yellow, Green, Aqua, Blue, Purple, Magenta (Hue/Saturation/Luminance)

**Why:** Capture One's .costyle preset format doesn't include HSL sliders. Color adjustments in Capture One are done through color editor layers (not part of presets).

❌ **Texture & Detail:**
- Vibrance (separate from Saturation)
- Sharpness (Amount/Radius/Threshold)

**Why:** Capture One handles these differently - Vibrance isn't a separate slider, and Sharpening isn't included in .costyle presets.

### Acceptable Precision Loss

**Exposure:** ±0.01 stops (float precision)
**Integer Parameters:** ±1 value (rounding)
**Temperature:** ±2 Kelvin (conversion)
**Color Balance:** ±1° hue, ±1 saturation (range conversion)

### Round-Trip Test Results (Story 8-4)

- Files tested: 3 real .costyle presets
- Avg accuracy: **98.37%** ✅ (exceeds 95% requirement)
- Edge cases: 4 tests (empty, extreme, minimal, complex) - all pass
- Bundle tests: ✓ Pass (3-recipe .costylepack round-trip)

See `internal/formats/costyle/testdata/costyle/test-results.json` for detailed accuracy breakdown.

### Unsupported .costyle Format Variants

**SL Format (Style Library) - NOT SUPPORTED:**

During Story 8-4 testing, we discovered an alternative .costyle format used by Capture One:

```xml
<?xml version="1.0" encoding="utf-8"?>
<SL Engine="1100">
 <E K="ICCProfile" V="profile.icc"/>
 <E K="Name" V="Preset Name"/>
 <E K="HighlightRecovery" V="50.0"/>
 <E K="Midtone" V="0;0;0;0.05"/>
</SL>
```

**Status:** ❌ Not supported in Recipe v2.0.0-dev (Epic 8)

**Why:** Recipe's .costyle parser only supports Adobe XMP-style .costyle files (with `<xmpmeta>` root element). The SL format uses a different XML structure, parameter model, and semantics.

**Impact:** 78 film emulation presets (Agfa, Fuji, Kodak, Rollei) cannot be converted by Recipe at this time.

**Workaround:** Export presets from Capture One as XMP-style .costyle (if format option available in Capture One).

**Future Work:** Support for SL-format .costyle files could be added in a future epic (requires separate parser implementation).

See `internal/formats/costyle/testdata/costyle/sl-format/README.md` for technical details on the SL format.

---

## Summary Table

| Feature          | XMP → NP3     | NP3 → XMP   | Visual Impact     | Workaround Available   |
| ---------------- | ------------- | ----------- | ----------------- | ---------------------- |
| Grain            | ❌ Not supported | N/A          | Low (1.2 ΔE)      | None                   |
| Vignette         | ❌ Not supported | N/A          | Medium (3.8 ΔE)   | Apply in post          |
| Split Toning     | ⚠️ Limited     | N/A          | Low-Med (2.5 ΔE)  | Use Recipe approximation |
| Lens Corrections | N/A          | N/A          | None              | Camera handles         |
| Clarity (Extreme)| ⚠️ Approx      | ⚠️ Approx    | Low-Med (4.5 ΔE)  | Accept approximation   |
| Picture Control  | N/A          | ⚠️ Approx    | Medium (3.5 ΔE)   | Adjust baseline in LR  |

**Key:**
- ❌ Not supported: Feature cannot be converted (data loss)
- ⚠️ Limited/Approx: Feature partially converted or approximated (acceptable fidelity)
- N/A: Not applicable (camera handles, or reverse path)

---

## Validation with Visual Regression Testing

All limitations documented above have been validated through visual regression testing (Story 6-2):

- **Reference images:** Applied presets in source applications (Lightroom, NX Studio)
- **Test images:** Applied converted presets in target applications
- **Metrics:** Delta E, SSIM calculated for each limitation
- **Assessment:** Subjective visual comparison at 100% zoom

See `docs/visual-regression-results.md` for detailed test results and side-by-side comparisons.

---

## User Recommendations

### When to Use Recipe Conversion

✅ **Good fit:**
- Presets with basic adjustments (Exposure, Contrast, Saturation, HSL)
- Portrait presets without heavy vignette
- Landscape presets with moderate clarity
- Black & white presets without grain effect

❌ **Not recommended:**
- Presets with heavy grain effect (Grain Amount >40)
- Presets with strong vignette (Vignette Amount >30)
- Presets with dramatic split toning (Saturation >40)
- Film emulation presets (often combine grain + vignette + split toning)

### Workflow Recommendations

1. **Test converted presets** with representative photos before relying on them
2. **Accept 95% similarity** - perfect 1:1 conversion is impossible due to format differences
3. **Apply missing effects in post** if critical to your creative vision
4. **Use Lightroom for grain/vignette** if these effects are essential
5. **Trust Recipe for core adjustments** - Exposure, Contrast, HSL, Tone Curves convert with high fidelity

---

## Community Feedback

We welcome feedback on conversion limitations! If you discover additional limitations or have suggestions for better approximations:

- **Open an issue:** https://github.com/yourusername/recipe/issues
- **Share test results:** Include reference/test images to help improve conversions
- **Suggest workarounds:** Community-contributed workarounds are valuable

---

**Last Updated:** 2025-11-06
**Validated With:** Visual Regression Testing (Story 6-2)
**Recipe Version:** v2.0.0-dev
