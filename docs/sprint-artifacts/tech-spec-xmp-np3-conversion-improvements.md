# Tech-Spec: XMP to NP3 Conversion Improvements

**Created:** 2025-12-18
**Status:** Ready for Development

## Overview

### Problem Statement

Many XMP (Lightroom) parameters have no direct equivalent in the NP3 (Nikon Picture Control) format. When converting professional XMP presets to NP3, these incompatible parameters are lost, resulting in presets that don't match the original Lightroom look. This forces preset creators to manually recreate adjustments in Nikon NX Studio.

### Solution

Implement **intelligent parameter conversion strategies** that approximate XMP-only parameters using available NP3 equivalents. Instead of silently dropping unsupported parameters, the converter will:

1. **Map parameters to closest NP3 equivalents** (e.g., Vibrance → Saturation + HSL adjustments)
2. **Convert parametric curves to point curves** (already implemented, needs verification)
3. **Convert RGB channel curves to Color Blender adjustments** (new)
4. **Provide clear warnings** when approximations are made

### Scope

**In Scope:**
- Vibrance → Saturation + selective HSL saturation boost
- RGB channel curves → Color Blender Hue/Chroma/Brightness adjustments
- Temperature/Tint → informational warning (no NP3 equivalent)
- Dehaze → Contrast + Blacks approximation
- Camera Calibration → Color Blender adjustments
- Parametric curve → Point curve conversion (verify existing implementation)

**Out of Scope:**
- Adding new NP3 offset discoveries (format already fully mapped)
- XMP parameters that affect non-color aspects (lens corrections, perspective, etc.)
- Noise reduction (not in NP3 format)

## Context for Development

### Codebase Patterns

- **UniversalRecipe model:** `internal/models/recipe.go` - intermediate representation
- **XMP Parser:** `internal/formats/xmp/parse.go` - extracts 60+ parameters with fallback support
- **NP3 Generator:** `internal/formats/np3/generate.go` - writes 48 parameters to exact offsets
- **NP3 Offsets:** `internal/formats/np3/offsets.go` - all byte offsets defined
- **Encoding helpers:** `EncodeSigned8()`, `EncodeScaled4()`, `EncodeHue12()`

### Files to Reference

| File | Purpose |
|------|---------|
| [parse.go](file:///home/justin/void/recipe/internal/formats/xmp/parse.go) | XMP parser with modern + legacy HSL attribute names |
| [generate.go](file:///home/justin/void/recipe/internal/formats/np3/generate.go) | NP3 generator with exact offset writers |
| [offsets.go](file:///home/justin/void/recipe/internal/formats/np3/offsets.go) | All NP3 byte offset definitions |
| [recipe.go](file:///home/justin/void/recipe/internal/models/recipe.go) | UniversalRecipe model definition |
| [XMP_PARAMETER_COMPATIBILITY_ANALYSIS.md](file:///home/justin/void/recipe/docs/XMP_PARAMETER_COMPATIBILITY_ANALYSIS.md) | Research on XMP/NP3 compatibility |

### Technical Decisions

1. **Approximation over loss:** Better to have an approximate conversion than dropping parameters entirely
2. **Warnings are mandatory:** Users must be informed when approximations are made
3. **Preserve existing behavior:** Don't break existing round-trip fidelity for parameters that already work

## Implementation Plan

### Tasks

#### Task 1: Vibrance → Saturation + HSL Approximation

- [ ] Add `convertVibrance()` function in `generate.go`
- [ ] Map Vibrance value to: 
  - 30% applied to global Saturation
  - 70% distributed across non-skin HSL channels (Orange, Yellow, Green, Aqua, Blue)
- [ ] Add warning when Vibrance is approximated

#### Task 2: RGB Channel Curves → Color Blender

- [ ] Add `convertRGBCurvesToColorBlender()` function in `generate.go`
- [ ] Analyze RGB curve shapes to determine dominant color shifts
- [ ] Map curve deviations from linear to Color Blender Hue/Chroma adjustments:
  - Red curve lift → Red/Orange Brightness increase
  - Blue curve lift → Blue/Purple Brightness increase
  - Inverse relationship for curve dips
- [ ] Add critical warning (WarnCritical) since this is a rough approximation

#### Task 3: Dehaze → Contrast + Blacks Approximation

- [ ] Modify `collectConversionWarnings()` to apply Dehaze approximation
- [ ] Apply Dehaze value as:
  - 50% added to Contrast
  - -50% added to Blacks (making shadows deeper)
- [ ] Add advisory warning

#### Task 4: Camera Calibration → Color Blender

- [ ] Add `convertCameraCalibration()` function
- [ ] Map Camera Calibration RGB adjustments to Color Blender:
  - RedHue/RedSaturation → Red Color Blender Hue/Chroma
  - GreenHue/GreenSaturation → Green Color Blender
  - BlueHue/BlueSaturation → Blue/Cyan Color Blender
- [ ] Add warning for ShadowTint (no equivalent)

#### Task 5: Verify Parametric Curve Conversion

- [ ] Review `ParametricToControlPoints()` in `curvegen.go`
- [ ] Add test cases for various parametric curve configurations
- [ ] Verify split point handling (ShadowSplit, MidtoneSplit, HighlightSplit)

#### Task 6: Temperature/Tint Warning

- [ ] Add informational warning when Temperature or Tint has non-zero values
- [ ] Suggest user manually adjust in NX Studio

### Acceptance Criteria

- [ ] AC 1: Given an XMP file with Vibrance +50, when converted to NP3, then Saturation and non-skin HSL channels are boosted proportionally
- [ ] AC 2: Given an XMP file with a non-linear Red channel curve, when converted to NP3, then Red/Orange Color Blender brightness is adjusted
- [ ] AC 3: Given an XMP file with Dehaze +30, when converted to NP3, then Contrast is increased and Blacks are deepened
- [ ] AC 4: Given an XMP file with Camera Calibration adjustments, when converted to NP3, then Color Blender values approximate the color shifts
- [ ] AC 5: All approximations generate appropriate warnings in ConversionResult
- [ ] AC 6: Existing test suite passes (no regressions)

## Additional Context

### Dependencies

- None (pure Go, no external libraries)

### Testing Strategy

**Existing Tests:**
```bash
# Run all NP3 tests
go test ./internal/formats/np3/... -v

# Run with coverage
go test ./internal/formats/np3/... -coverprofile=coverage.out
```

**New Tests to Add:**
1. `TestVibranceApproximation` - verify Vibrance → Saturation + HSL mapping
2. `TestRGBCurveConversion` - verify RGB curves → Color Blender mapping
3. `TestDehazeApproximation` - verify Dehaze → Contrast + Blacks
4. `TestCameraCalibrationConversion` - verify Camera Calibration → Color Blender

**Manual Verification:**
1. Convert a professional XMP preset (e.g., from testdata) with Vibrance
2. Load resulting NP3 in Nikon NX Studio
3. Visually compare color saturation to Lightroom

### Notes

The XMP parser has already been updated to correctly parse:
- Modern HSL attribute names (`HueAdjustmentRed`, `SaturationAdjustmentRed`, etc.)
- Incremental Temperature/Tint (relative white balance)
- Camera Calibration RGB adjustments
- Parametric curve parameters
- All Color Grading parameters

The focus of this tech-spec is on the **conversion logic** in the NP3 generator, not parsing.

---

## Reference Preset Analysis: Agfachrome RSX 200

This professional film emulation preset serves as the primary validation target.

### Parameters Used

| Category | Parameters | NP3 Support |
|----------|-----------|-------------|
| **Basic Tone** | Contrast +19, Highlights -17, Shadows +16, Whites +6, Blacks -13 | ✅ Full |
| **Presence** | Texture -3, Clarity -8 | ✅ Full |
| **Master Curve** | 3 points: (0,13), (65,65), (255,255) - lifted blacks | ✅ Full |
| **RGB Curves** | R/G/B identical: (0,10), (71,64), (188,194), (255,242) | ❌ Need approx |
| **Parametric Curve** | Darks -10, Lights +5, Highlights +10, splits 25/50/75 | ✅ Converted to points |
| **HSL Hue** | Red -5, Orange -3, Yellow -3, Green +9, Aqua +5 | ✅ Full |
| **HSL Saturation** | Red -3, Orange -2, Yellow -6, Green -10, Aqua -6, Blue -3 | ✅ Full |
| **HSL Luminance** | Red +5, Orange +8, Yellow +5, Green +8, Aqua +5, Blue +3 | ✅ Full |
| **Split Toning** | Shadow: 215°/13, Highlight: 35°/12, Balance -7 | ✅ Via Color Grading |
| **Color Grading** | Midtone: 220°/10, Blending 50 | ✅ Full |
| **Camera Calibration** | Red +3/+3, Green -4/+3, Blue +10/-2 | ❌ Need approx |
| **White Balance** | Incremental Temp +9, Tint +3 | ❌ No NP3 equivalent |
| **Grain** | Amount 22, Size 25, Frequency 50 | ✅ Full |

### Critical Conversion Challenges

1. **RGB Channel Curves:** All three channels use identical S-curves with lifted blacks (0→10) and crushed whites (255→242). Since NP3 only has a master curve, these need approximation via Color Blender.

2. **Camera Calibration:** RedHue +3, BlueHue +10 significantly shift color science. Need to map to Color Blender.

3. **White Balance:** IncrementalTemperature +9 adds warmth - no NP3 equivalent but could be approximated via Color Grading warm tones.

---

## Visual Comparison Testing

### Available Test Resources

**NEF Test Files:** 11 high-resolution RAW files in `testdata/visual-regression/images/`:
- DSC_1631.nef through DSC_2226.nef (Nikon Z camera samples)
- PNK_0716.NEF through PNK_0739.NEF (Beach portraits)

**Tools Available:**
- ImageMagick (`convert`, `magick`) - for image comparison and Delta E calculation
- exiftool - for metadata extraction

### Automated Visual Comparison Pipeline

```bash
# Step 1: Convert XMP preset to NP3
./recipe convert "testdata/xmp/Agfachrome RSX 200.xmp" -o output.np3

# Step 2: Export reference image from Lightroom with XMP preset applied
# (Manual step - export as 16-bit TIFF to testdata/visual-regression/reference/)

# Step 3: Apply NP3 in NX Studio and export test image
# (Manual step - export as 16-bit TIFF to testdata/visual-regression/test/)

# Step 4: Calculate Delta E using ImageMagick
magick compare -metric RMSE reference.tif test.tif diff.png

# Step 5: Calculate SSIM (if ImageMagick supports it)
# Alternative: Use Python script with scikit-image

# Step 6: Generate visual diff overlay
magick composite reference.tif test.tif -compose difference diff.tif
```

### Python Visual Comparison Script (New)

Create `scripts/visual_compare.py`:

```python
#!/usr/bin/env python3
"""Visual comparison between Lightroom and NX Studio exports."""
from PIL import Image
import numpy as np
from skimage.metrics import structural_similarity as ssim
from colormath.color_objects import LabColor, sRGBColor
from colormath.color_conversions import convert_color
from colormath.color_diff import delta_e_cie2000

def calculate_delta_e(img1_path, img2_path, sample_points=1000):
    """Calculate average Delta E between two images."""
    img1 = np.array(Image.open(img1_path).convert('RGB'))
    img2 = np.array(Image.open(img2_path).convert('RGB'))
    
    # Sample random points
    h, w = img1.shape[:2]
    indices = np.random.choice(h * w, sample_points, replace=False)
    
    delta_es = []
    for idx in indices:
        y, x = divmod(idx, w)
        
        rgb1 = sRGBColor(img1[y,x,0]/255, img1[y,x,1]/255, img1[y,x,2]/255)
        rgb2 = sRGBColor(img2[y,x,0]/255, img2[y,x,1]/255, img2[y,x,2]/255)
        
        lab1 = convert_color(rgb1, LabColor)
        lab2 = convert_color(rgb2, LabColor)
        
        delta_es.append(delta_e_cie2000(lab1, lab2))
    
    return np.mean(delta_es), np.max(delta_es), np.percentile(delta_es, 95)

def calculate_ssim(img1_path, img2_path):
    """Calculate SSIM between two images."""
    img1 = np.array(Image.open(img1_path).convert('L'))
    img2 = np.array(Image.open(img2_path).convert('L'))
    return ssim(img1, img2)
```

### Acceptance Threshold

| Metric | Target | Acceptable |
|--------|--------|------------|
| Delta E (avg) | < 3.0 | < 5.0 |
| Delta E (95th pctl) | < 5.0 | < 8.0 |
| SSIM | > 0.95 | > 0.90 |

### Manual Verification Steps

1. **Export from Lightroom:**
   - Open any NEF from `testdata/visual-regression/images/`
   - Apply "Agfachrome RSX 200" preset
   - Export as 16-bit TIFF (Adobe RGB, no sharpening)
   - Save to `testdata/visual-regression/reference/agfachrome-rsx-200-{filename}.tif`

2. **Apply NP3 in NX Studio:**
   - Run: `./recipe convert "testdata/xmp/Agfachrome RSX 200.xmp" -o Agfachrome_RSX_200.np3`
   - Open same NEF in NX Studio
   - Load the converted NP3 Picture Control
   - Export with identical settings
   - Save to `testdata/visual-regression/test/agfachrome-rsx-200-{filename}.tif`

3. **Visual Side-by-Side:**
   - Open both TIFFs in image viewer
   - Compare overall color cast, contrast, shadow detail
   - Check specific areas: blue sky, green foliage, skin tones (if present)
