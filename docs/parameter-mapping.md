# Parameter Mapping Rules

## Overview

This document provides comprehensive parameter mapping rules for converting photo editing presets between NP3, XMP, and lrtemplate formats. The Recipe conversion engine uses a **hub-and-spoke pattern** where `UniversalRecipe` serves as the central hub, and each format implements bidirectional conversions (Parse and Generate).

**Key Concepts:**
- **Direct Mapping**: Parameters with identical ranges and semantics (XMP ↔ lrtemplate)
- **Approximation Mapping**: Parameters requiring formula-based range conversion (NP3 ↔ XMP/lrtemplate)
- **Unmappable Parameters**: Advanced features present in XMP/lrtemplate but absent in NP3

**Visual Similarity Goal**: ≥95% visual similarity for common adjustments (Exposure, Contrast, Saturation)

## Table of Contents

1. [Direct Parameter Mappings (XMP ↔ lrtemplate)](#direct-parameter-mappings-xmp--lrtemplate)
2. [Approximation Mappings (NP3 ↔ XMP/lrtemplate)](#approximation-mappings-np3--xmplrtemplate)
3. [Unmappable Parameters](#unmappable-parameters)
4. [Bidirectional Conversion Paths](#bidirectional-conversion-paths)
5. [Visual Similarity Validation](#visual-similarity-validation)
6. [Metadata Dictionary Usage](#metadata-dictionary-usage)
7. [Error Reporting and Warnings](#error-reporting-and-warnings)
8. [Implementation Guidance](#implementation-guidance)

---

## Direct Parameter Mappings (XMP ↔ lrtemplate)

XMP and lrtemplate formats use **identical parameter names and ranges** with only syntax differences (XML vs Lua). This enables 1:1 mapping without any transformation.

### Complete Mapping Table

| Category | XMP Field Name | lrtemplate Field Name | UniversalRecipe Field | Range | Notes |
|----------|----------------|----------------------|----------------------|-------|-------|
| **Basic Adjustments** |
| Exposure | `Exposure2012` | `Exposure2012` | `Exposure` | -5.0 to +5.0 | Float, Process Version 2012 |
| Contrast | `Contrast2012` | `Contrast2012` | `Contrast` | -100 to +100 | Integer |
| Highlights | `Highlights2012` | `Highlights2012` | `Highlights` | -100 to +100 | Integer |
| Shadows | `Shadows2012` | `Shadows2012` | `Shadows` | -100 to +100 | Integer |
| Whites | `Whites2012` | `Whites2012` | `Whites` | -100 to +100 | Integer |
| Blacks | `Blacks2012` | `Blacks2012` | `Blacks` | -100 to +100 | Integer |
| **Color Adjustments** |
| Saturation | `Saturation` | `Saturation` | `Saturation` | -100 to +100 | Integer |
| Vibrance | `Vibrance` | `Vibrance` | `Vibrance` | -100 to +100 | Integer |
| **Presence** |
| Texture | `Texture` | `Texture` | `Texture` | -100 to +100 | Integer |
| Clarity | `Clarity2012` | `Clarity2012` | `Clarity` | -100 to +100 | Integer |
| Dehaze | `Dehaze` | `Dehaze` | `Dehaze` | -100 to +100 | Integer |
| **Sharpening** |
| Sharpness | `Sharpness` | `Sharpness` | `Sharpness` | 0 to 150 | Integer |
| Sharpness Radius | `SharpenRadius` | `SharpenRadius` | `SharpnessRadius` | 0.5 to 3.0 | Float |
| Sharpness Detail | `SharpenDetail` | `SharpenDetail` | `SharpnessDetail` | 0 to 100 | Integer |
| Sharpness Masking | `SharpenEdgeMasking` | `SharpenEdgeMasking` | `SharpnessMasking` | 0 to 100 | Integer |
| **White Balance** |
| Temperature | `Temperature` | `Temperature` | `Temperature` | 2000 to 50000 | Integer (Kelvin), nullable |
| Tint | `Tint` | `Tint` | `Tint` | -150 to +150 | Integer |

### HSL Adjustments (8 Colors × 3 Properties = 24 Parameters)

Each of 8 colors (Red, Orange, Yellow, Green, Aqua, Blue, Purple, Magenta) has 3 adjustment properties:

| Property | XMP Pattern | lrtemplate Pattern | UniversalRecipe Field | Range |
|----------|-------------|-------------------|----------------------|-------|
| Hue | `HueAdjustment{Color}` | `HueAdjustment{Color}` | `{Color}.Hue` | -100 to +100 |
| Saturation | `SaturationAdjustment{Color}` | `SaturationAdjustment{Color}` | `{Color}.Saturation` | -100 to +100 |
| Luminance | `LuminanceAdjustment{Color}` | `LuminanceAdjustment{Color}` | `{Color}.Luminance` | -100 to +100 |

**Example**: Red color adjustments
- XMP: `HueAdjustmentRed`, `SaturationAdjustmentRed`, `LuminanceAdjustmentRed`
- lrtemplate: Same field names
- UniversalRecipe: `Red.Hue`, `Red.Saturation`, `Red.Luminance`

### Tone Curves

| Type | XMP Field | lrtemplate Field | UniversalRecipe Field | Format |
|------|-----------|------------------|----------------------|--------|
| Master Curve | `ToneCurvePV2012` | `ToneCurvePV2012` | `PointCurve` | Array of `{Input, Output}` pairs (0-255) |
| Red Channel | `ToneCurvePV2012Red` | `ToneCurvePV2012Red` | `PointCurveRed` | Array of `{Input, Output}` pairs |
| Green Channel | `ToneCurvePV2012Green` | `ToneCurvePV2012Green` | `PointCurveGreen` | Array of `{Input, Output}` pairs |
| Blue Channel | `ToneCurvePV2012Blue` | `ToneCurvePV2012Blue` | `PointCurveBlue` | Array of `{Input, Output}` pairs |

### Split Toning

| Parameter | XMP Field | lrtemplate Field | UniversalRecipe Field | Range |
|-----------|-----------|------------------|----------------------|-------|
| Shadow Hue | `SplitToningShadowHue` | `SplitToningShadowHue` | `SplitShadowHue` | 0 to 360 |
| Shadow Saturation | `SplitToningShadowSaturation` | `SplitToningShadowSaturation` | `SplitShadowSaturation` | 0 to 100 |
| Highlight Hue | `SplitToningHighlightHue` | `SplitToningHighlightHue` | `SplitHighlightHue` | 0 to 360 |
| Highlight Saturation | `SplitToningHighlightSaturation` | `SplitToningHighlightSaturation` | `SplitHighlightSaturation` | 0 to 100 |
| Balance | `SplitToningBalance` | `SplitToningBalance` | `SplitBalance` | -100 to +100 |

### Effects

| Effect | XMP Field | lrtemplate Field | UniversalRecipe Field | Range |
|--------|-----------|------------------|----------------------|-------|
| Grain Amount | `GrainAmount` | `GrainAmount` | `GrainAmount` | 0 to 100 |
| Grain Size | `GrainSize` | `GrainSize` | `GrainSize` | 0 to 100 |
| Grain Roughness | `GrainFrequency` | `GrainFrequency` | `GrainRoughness` | 0 to 100 |

### Vignette

| Parameter | XMP Field | lrtemplate Field | UniversalRecipe Field | Range |
|-----------|-----------|------------------|----------------------|-------|
| Amount | `VignetteAmount` | `VignetteAmount` | `VignetteAmount` | -100 to +100 |
| Midpoint | `VignetteMidpoint` | `VignetteMidpoint` | `VignetteMidpoint` | 0 to 100 |
| Roundness | `VignetteRoundness` | `VignetteRoundness` | `VignetteRoundness` | -100 to +100 |
| Feather | `VignetteFeather` | `VignetteFeather` | `VignetteFeather` | 0 to 100 |

### Camera Calibration

| Parameter | XMP Field | lrtemplate Field | UniversalRecipe Field | Range |
|-----------|-----------|------------------|----------------------|-------|
| Red Hue | `RedHue` | `RedHue` | `CameraProfile.RedHue` | -100 to +100 |
| Red Saturation | `RedSaturation` | `RedSaturation` | `CameraProfile.RedSaturation` | -100 to +100 |
| Green Hue | `GreenHue` | `GreenHue` | `CameraProfile.GreenHue` | -100 to +100 |
| Green Saturation | `GreenSaturation` | `GreenSaturation` | `CameraProfile.GreenSaturation` | -100 to +100 |
| Blue Hue | `BlueHue` | `BlueHue` | `CameraProfile.BlueHue` | -100 to +100 |
| Blue Saturation | `BlueSaturation` | `BlueSaturation` | `CameraProfile.BlueSaturation` | -100 to +100 |

### Field Naming Conventions

1. **Process Version 2012**: Many fields use `2012` suffix (e.g., `Exposure2012`, `Contrast2012`) indicating Process Version 2012 - Lightroom's modern processing engine.

2. **No ProcessVersion Suffix**: Some fields don't have version suffix (e.g., `Saturation`, `Temperature`) because they've remained consistent across versions.

3. **HSL Pattern**: HSL fields use pattern `PropertyAdjustmentColor` (e.g., `HueAdjustmentRed`)

4. **Tone Curve Pattern**: Tone curves use `ToneCurvePV2012` prefix with optional channel suffix

### Code Example: XMP ↔ lrtemplate Mapping

```go
// XMP → UniversalRecipe (Parse)
recipe.Exposure = parseFloat(xmpDoc, "Exposure2012")
recipe.Contrast = parseInt(xmpDoc, "Contrast2012")
recipe.Saturation = parseInt(xmpDoc, "Saturation")

// UniversalRecipe → lrtemplate (Generate)
buf.WriteString(fmt.Sprintf("\t\t\tExposure2012 = %.2f,\n", recipe.Exposure))
buf.WriteString(fmt.Sprintf("\t\t\tContrast2012 = %d,\n", recipe.Contrast))
buf.WriteString(fmt.Sprintf("\t\t\tSaturation = %d,\n", recipe.Saturation))
```

**Key Insight**: XMP and lrtemplate use **identical field names**, so conversion between these two formats is simply format translation (XML ↔ Lua) with no parameter transformation needed.

---

## Approximation Mappings (NP3 ↔ XMP/lrtemplate)

NP3 format supports only **5 core parameters** in a fixed 1024-byte binary structure. Converting to/from XMP/lrtemplate requires approximation formulas to map between different ranges.

### NP3 Format Limitations

| Capability | NP3 | XMP/lrtemplate |
|------------|-----|----------------|
| Parameters | 5 (Sharpening, Contrast, Brightness, Saturation, Hue) | 50+ |
| Adjustment Precision | Low (3-9 steps) | High (±100 or ±5.0) |
| HSL Support | Global hue only (-9° to +9°) | 8 colors × 3 properties |
| Tone Curves | None | Parametric + RGB point curves |
| Split Toning | None | Full support |
| Effects | None | Grain, vignette, etc. |

### Approximation Formulas

#### 1. Sharpening / Sharpness

**NP3 Range**: 0-9 (10 discrete levels)
**UniversalRecipe Range**: 0-150

**Formula**:
- **Parse (NP3 → Universal)**: `UniversalSharpness = NP3Sharpening × 10`
- **Generate (Universal → NP3)**: `NP3Sharpening = UniversalSharpness / 10`

**Rationale**: NP3's 10 levels (0-9) map cleanly to 0-90 range in Universal. Values 91-150 are clamped to 9.

**Code Example**:
```go
const NP3_SHARPNESS_SCALE = 10

// Parser: internal/formats/np3/parse.go:398
builder.WithSharpness(params.sharpening * NP3_SHARPNESS_SCALE)  // 0-9 → 0-90

// Generator: internal/formats/np3/generate.go:69
params.sharpening = recipe.Sharpness / NP3_SHARPNESS_SCALE  // 0-150 → 0-9 (clamped)
```

**Visual Impact**: Sharpening affects edge definition. 10 levels provide adequate granularity for most uses. Loss of precision above level 9 (90/150) may result in slightly less sharp images when converting from XMP/lrtemplate to NP3.

---

#### 2. Contrast

**NP3 Range**: -3 to +3 (7 discrete levels)
**UniversalRecipe Range**: -100 to +100

**Formula**:
- **Parse (NP3 → Universal)**: `UniversalContrast = NP3Contrast × 33`
- **Generate (Universal → NP3)**: `NP3Contrast = UniversalContrast / 33`

**Rationale**: NP3's 7 levels map to ±99 range (3 × 33 = 99). Scale factor 33 chosen to maximize range utilization while staying within ±100 bounds.

**Mapping Table**:

| NP3 Contrast | Universal Contrast | Visual Effect |
|--------------|-------------------|---------------|
| -3 | -99 | Maximum contrast reduction |
| -2 | -66 | Strong contrast reduction |
| -1 | -33 | Moderate contrast reduction |
| 0 | 0 | No adjustment |
| +1 | +33 | Moderate contrast increase |
| +2 | +66 | Strong contrast increase |
| +3 | +99 | Maximum contrast increase |

**Code Example**:
```go
const NP3_CONTRAST_SCALE = 33

// Parser: internal/formats/np3/parse.go:399
builder.WithContrast(params.contrast * NP3_CONTRAST_SCALE)  // -3 to +3 → -99 to +99

// Generator: internal/formats/np3/generate.go:77
params.contrast = recipe.Contrast / NP3_CONTRAST_SCALE  // -100 to +100 → -3 to +3
```

**Visual Impact**: Contrast is a critical parameter for visual similarity (weight: 1.0). The 33-unit steps are perceptually significant, so rounding may be noticeable when converting from XMP/lrtemplate to NP3.

---

#### 3. Brightness / Exposure

**NP3 Range**: -1.0 to +1.0 (floating point)
**UniversalRecipe Range**: -5.0 to +5.0 (Exposure field)

**Formula**:
- **Parse (NP3 → Universal)**: `UniversalExposure = NP3Brightness` (direct copy, within bounds)
- **Generate (Universal → NP3)**: `NP3Brightness = clamp(UniversalExposure, -1.0, 1.0)`

**Rationale**: NP3 brightness maps to UniversalRecipe Exposure field. NP3's ±1.0 range is subset of Universal's ±5.0 range. Values outside ±1.0 are clamped when generating NP3.

**Code Example**:
```go
// Parser: internal/formats/np3/parse.go:400
builder.WithExposure(params.brightness)  // -1.0 to +1.0 → Exposure field

// Generator: internal/formats/np3/generate.go:96-101
params.brightness = recipe.Exposure
if params.brightness > 1.0 {
    params.brightness = 1.0
} else if params.brightness < -1.0 {
    params.brightness = -1.0
}
```

**Visual Impact**: Exposure is the most critical parameter for visual similarity (weight: 1.0). Clamping to ±1.0 range means extreme exposure adjustments (+2.0 to +5.0 or -2.0 to -5.0) cannot be represented in NP3, resulting in visible brightness differences.

---

#### 4. Saturation

**NP3 Range**: -3 to +3 (7 discrete levels)
**UniversalRecipe Range**: -100 to +100

**Formula**:
- **Parse (NP3 → Universal)**: `UniversalSaturation = NP3Saturation × 33`
- **Generate (Universal → NP3)**: `NP3Saturation = UniversalSaturation / 33`

**Rationale**: Identical to Contrast mapping. Scale factor 33 maximizes range utilization (3 × 33 = 99 ≈ 100).

**Mapping Table**: Same as Contrast (see above)

**Code Example**:
```go
const NP3_SATURATION_SCALE = 33

// Parser: internal/formats/np3/parse.go:401
builder.WithSaturation(params.saturation * NP3_SATURATION_SCALE)  // -3 to +3 → -99 to +99

// Generator: internal/formats/np3/generate.go:87
params.saturation = recipe.Saturation / NP3_SATURATION_SCALE  // -100 to +100 → -3 to +3
```

**Visual Impact**: Saturation is very important for visual similarity (weight: 0.8). The 33-unit steps are perceptually noticeable, especially in highly saturated images (landscapes, sunsets).

---

#### 5. Hue

**NP3 Range**: -9° to +9° (global hue shift)
**UniversalRecipe Range**: 8 colors × HSL properties (no direct equivalent)

**Status**: **No Direct Mapping**

**Challenge**: NP3 provides a global hue shift (-9° to +9°) affecting all colors uniformly. UniversalRecipe only supports per-color hue adjustments (Red, Orange, Yellow, etc.). These are fundamentally different operations:
- NP3 Hue: Shifts entire color wheel (affects all colors equally)
- Universal HSL: Adjusts individual color ranges independently

**Approximation Strategy**: Best-effort mapping to Temperature/Tint fields

**Current Implementation**:
- **Parse (NP3 → Universal)**: NP3 Hue is **ignored** (no target field)
- **Generate (Universal → NP3)**: NP3 Hue set to **0** (neutral default)

**Future Enhancement** (optional):
```go
// Approximate global hue shift using Temperature/Tint
// Hue shift toward warm (red/yellow): increase Temperature
// Hue shift toward cool (blue/cyan): decrease Temperature
if np3Hue > 0 {
    recipe.Temperature = 5500 + (int(np3Hue) * 200)  // Approximate: +200K per degree
} else if np3Hue < 0 {
    recipe.Temperature = 5500 - (int(np3Hue) * 200)
}
```

**Warning Message**: "NP3 global hue adjustment cannot be precisely mapped to per-color HSL adjustments. Hue shift will not be preserved in conversion."

**Visual Impact**: Global hue shifts in NP3 presets will be lost during conversion to XMP/lrtemplate, potentially resulting in significant color cast differences.

---

### Summary: Approximation Accuracy

| Parameter | Scale Factor | Precision Loss | Visual Impact | Weight |
|-----------|-------------|----------------|---------------|--------|
| Sharpness | ×10 | Low (10 levels) | Moderate | 0.5 |
| Contrast | ×33 | Moderate (7 levels) | **High** | **1.0** |
| Brightness/Exposure | 1:1 (clamped) | High (±1.0 vs ±5.0) | **Critical** | **1.0** |
| Saturation | ×33 | Moderate (7 levels) | **High** | **0.8** |
| Hue | **Unmappable** | **Complete loss** | Variable | 0.5 |

**Overall**: NP3 ↔ XMP/lrtemplate conversions can achieve ≥95% visual similarity for **basic adjustments** (exposure, contrast, saturation within ±1 stop). Advanced features and extreme values cannot be represented in NP3.

---

## Unmappable Parameters

XMP and lrtemplate formats support 50+ advanced parameters that **cannot be mapped to NP3** due to format limitations. These parameters must be handled gracefully to avoid silent data loss.

### Complete List of Unmappable Parameters (NP3 Generation Only)

| Category | Parameters | Reason Unmappable | Strategy |
|----------|-----------|-------------------|----------|
| **Tone Controls** | Highlights, Shadows, Whites, Blacks | NP3 has only single Contrast parameter | Store in Metadata, warn user |
| **Presence** | Texture, Clarity, Dehaze, Vibrance | NP3 has no equivalent | Store in Metadata, warn user |
| **HSL Adjustments** | 8 colors × 3 properties (24 total) | NP3 has only global Hue | Store in Metadata, **critical warning** |
| **Tone Curves** | ToneCurvePV2012, ToneCurvePV2012Red, ToneCurvePV2012Green, ToneCurvePV2012Blue | NP3 has no tone curve support | Store in Metadata, warn user |
| **Split Toning** | SplitToningShadowHue, SplitToningShadowSaturation, SplitToningHighlightHue, SplitToningHighlightSaturation, SplitToningBalance | NP3 has no split toning | Store in Metadata, warn user |
| **Effects - Grain** | GrainAmount, GrainSize, GrainRoughness | NP3 has no grain effect | Store in Metadata, warn user |
| **Vignette** | VignetteAmount, VignetteMidpoint, VignetteRoundness, VignetteFeather | NP3 has no vignette effect | Store in Metadata, warn user |
| **Sharpening Detail** | SharpnessRadius, SharpnessDetail, SharpnessMasking | NP3 has only basic Sharpening 0-9 | Store in Metadata, warn user |
| **White Balance Detail** | Tint | NP3 has no separate Tint control | Store in Metadata, warn user |
| **Camera Calibration** | RedHue, RedSaturation, GreenHue, GreenSaturation, BlueHue, BlueSaturation | NP3 has no camera calibration | Store in Metadata, warn user |

### Decision Matrix: How to Handle Each Unmappable Parameter

| Parameter Type | Action | Warning Level | Preserve in Metadata? |
|----------------|--------|---------------|----------------------|
| Highlights, Shadows, Whites, Blacks | Warn, continue | Advisory | Yes |
| Texture, Clarity, Dehaze | Warn, continue | Advisory | Yes |
| Vibrance | Warn, continue | Advisory | Yes (close to Saturation) |
| **HSL Adjustments (24 params)** | **Warn, continue** | **Critical** | **Yes** (high impact) |
| **Tone Curves** | **Warn, continue** | **Critical** | **Yes** (complex data) |
| Split Toning | Warn, continue | Advisory | Yes |
| Grain | Warn, continue | Informational | Yes |
| Vignette | Warn, continue | Informational | Yes |
| Sharpening Detail | Warn, continue | Informational | No (Basic Sharpness preserved) |
| Tint | Warn, continue | Informational | No |
| Camera Calibration | Warn, continue | Informational | Yes |

### User Warning Messages

#### Critical Warnings (High Visual Impact)

```
⚠️  CRITICAL: NP3 does not support HSL adjustments
Your preset contains per-color hue, saturation, and luminance adjustments for
{affected_colors} which will be lost in NP3 conversion. This may significantly
alter the color rendering of your images.

Recommendation: Use XMP or lrtemplate format to preserve full color control.
```

```
⚠️  CRITICAL: NP3 does not support tone curves
Your preset uses custom tone curves (ToneCurvePV2012) which will be lost in
NP3 conversion. This will change the tonal distribution of your images.

Recommendation: Use XMP or lrtemplate format to preserve tone curves.
```

#### Advisory Warnings (Moderate Visual Impact)

```
⚠️  ADVISORY: Advanced tone controls lost in NP3 conversion
Your preset uses Highlights, Shadows, Whites, and/or Blacks adjustments which
cannot be represented in NP3 (only Contrast is available). Image tonality may differ.

Preserved in Universal→NP3: Contrast approximation
Lost: Separate highlight/shadow control
```

```
⚠️  ADVISORY: Presence adjustments not supported in NP3
Your preset includes Texture, Clarity, and/or Dehaze adjustments which will be
lost in NP3 conversion. Image clarity and local contrast may differ.
```

#### Informational Warnings (Low Visual Impact)

```
ℹ️  INFO: NP3 does not support grain effects
GrainAmount, GrainSize, and GrainRoughness will be lost in conversion. Original
values preserved in metadata for round-trip conversion.
```

```
ℹ️  INFO: NP3 does not support vignette effects
Vignette parameters will be lost in NP3 conversion. Original values preserved in
metadata for round-trip conversion.
```

### Interface-Specific Warning Format

#### CLI/TUI (Detailed)
```bash
$ recipe convert preset.xmp preset.np3

⚠️  WARNINGS: Generating NP3 from XMP preset "Cinematic Look"
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CRITICAL ISSUES (will significantly impact visual appearance):
  • HSL adjustments: Red hue +15, saturation -20 (LOST)
  • Tone curves: Custom RGB curves with 8 points (LOST)

ADVISORY ISSUES (may impact appearance):
  • Highlights: -45 (approximated by Contrast)
  • Shadows: +30 (approximated by Contrast)
  • Texture: +15 (LOST)
  • Clarity: +25 (LOST)

INFORMATIONAL:
  • Grain: Amount 35, Size 25 (preserved in metadata)
  • Vignette: Amount -40 (preserved in metadata)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ NP3 file generated: preset.np3
⚠️  Visual similarity: ~75% (due to HSL and tone curve loss)

Continue? [Y/n]
```

#### Web Interface (Concise)
```
⚠️  4 warnings when converting to NP3

Critical: HSL adjustments and tone curves will be lost
Advisory: Some tone controls approximated
Info: Grain and vignette effects not supported

[View Details] [Continue Anyway] [Use XMP Instead]
```

### When to Fail vs Proceed with Warnings

**Default Behavior**: **Proceed with warnings** (best for user flexibility)

```go
// Generate NP3 with warnings
func Generate(recipe *models.UniversalRecipe) ([]byte, []ConversionWarning, error) {
    warnings := collectUnmappableWarnings(recipe)

    // Always proceed - user can decide based on warnings
    data, err := generateBinary(recipe)
    return data, warnings, err
}
```

**Strict Mode** (optional CLI flag `--strict`): **Fail on unmappable parameters**

```bash
$ recipe convert preset.xmp preset.np3 --strict

❌ ERROR: Conversion failed in strict mode

NP3 format cannot represent the following parameters in your preset:
  • HSL adjustments (Red, Blue)
  • Tone curves (RGB channels)

Use --no-strict to generate NP3 with warnings, or convert to XMP/lrtemplate instead.
```

---

## Bidirectional Conversion Paths

The Recipe engine supports **6 bidirectional conversion paths** between the 3 formats:

### Conversion Path Matrix

| Source Format | Target Format | Mapping Type | Data Loss | Typical Accuracy |
|---------------|---------------|--------------|-----------|------------------|
| NP3 | XMP | Approximation | None (upconvert) | 100% |
| NP3 | lrtemplate | Approximation | None (upconvert) | 100% |
| XMP | NP3 | Approximation | **High** (40+ params lost) | 60-80% |
| XMP | lrtemplate | Direct (1:1) | None | 100% |
| lrtemplate | NP3 | Approximation | **High** (40+ params lost) | 60-80% |
| lrtemplate | XMP | Direct (1:1) | None | 100% |

### Path Details

#### 1. NP3 → XMP (Approximation, No Loss)

**Process**: Parse NP3 → UniversalRecipe → Generate XMP

**Mapping**:
- Sharpening (0-9) → Sharpness (0-90) via ×10
- Contrast (-3 to +3) → Contrast (-99 to +99) via ×33
- Brightness (-1.0 to +1.0) → Exposure (-1.0 to +1.0)
- Saturation (-3 to +3) → Saturation (-99 to +99) via ×33
- Hue (-9° to +9°) → **No equivalent** (ignored)

**Zero-Value Handling**: NP3 parser sets UniversalRecipe fields; XMP generator omits zero values → Clean XMP output

**Edge Cases**:
- NP3 file with all-zero parameters → Minimal XMP preset (only ProcessVersion field)
- Corrupt NP3 data → Parser validation catches invalid ranges before conversion

**Example**:
```go
// NP3 values
sharpening=5, contrast=2, brightness=0.5, saturation=1

// Converted to UniversalRecipe
Sharpness=50, Contrast=66, Exposure=0.5, Saturation=33

// Generated XMP
<crs:Sharpness>50</crs:Sharpness>
<crs:Contrast2012>66</crs:Contrast2012>
<crs:Exposure2012>0.5</crs:Exposure2012>
<crs:Saturation>33</crs:Saturation>
```

---

#### 2. NP3 → lrtemplate (Approximation, No Loss)

**Process**: Parse NP3 → UniversalRecipe → Generate lrtemplate

**Mapping**: Identical to NP3 → XMP (same parameters, different syntax)

**Output Format**: Lua table instead of XML

**Example**:
```lua
s = {
    id = "00000000-0000-0000-0000-000000000000",
    internalName = "NP3 Preset",
    title = "NP3 Preset",
    type = "Develop",
    value = {
        settings = {
            ProcessVersion = "10.0",
            Sharpness = 50,
            Contrast2012 = 66,
            Exposure2012 = 0.5,
            Saturation = 33,
        },
        uuid = "00000000-0000-0000-0000-000000000001",
    },
    version = 0,
}
```

---

#### 3. XMP → NP3 (Approximation, HIGH Data Loss)

**Process**: Parse XMP → UniversalRecipe → Generate NP3

**Mapping**: Reverse formulas from Path #1, with clamping

**Data Loss**:
- Only 5 parameters preserved (Exposure, Contrast, Saturation, Sharpness, Hue)
- 40+ advanced parameters lost (HSL, tone curves, split toning, grain, vignette, etc.)
- Values outside NP3 ranges clamped

**Rounding Rules**:
- Sharpness: `sharpening = int(sharpness / 10)` → ±5 units precision loss
- Contrast: `contrast = int(contrast / 33)` → ±16 units precision loss
- Saturation: `saturation = int(saturation / 33)` → ±16 units precision loss
- Exposure: Clamped to ±1.0 (values outside lost)

**Edge Cases**:
- XMP preset with only unmappable parameters (HSL only) → Empty NP3 (all zeros)
- XMP with Exposure = +3.0 → Clamped to +1.0 in NP3 (image 2 stops darker)

**Example with Data Loss**:
```xml
<!-- Original XMP -->
<crs:Exposure2012>2.5</crs:Exposure2012>
<crs:Contrast2012>75</crs:Contrast2012>
<crs:Highlights2012>-50</crs:Highlights2012>
<crs:Texture>25</crs:Texture>
<crs:HueAdjustmentRed>+15</crs:HueAdjustmentRed>

<!-- Converted to NP3 (internal) -->
sharpening=0, contrast=2, brightness=1.0 (clamped!), saturation=0

<!-- Lost in conversion -->
- Exposure clamped: 2.5 → 1.0 (1.5 stops lost)
- Highlights: -50 (no NP3 equivalent)
- Texture: 25 (no NP3 equivalent)
- HSL Red: +15 (no NP3 equivalent)

⚠️  Visual similarity: ~70% (significant brightness difference)
```

---

#### 4. XMP → lrtemplate (Direct 1:1, No Loss)

**Process**: Parse XMP → UniversalRecipe → Generate lrtemplate

**Mapping**: Direct field-name mapping (no transformation)

**Accuracy**: 100% - Byte-for-byte accuracy achievable (only formatting differs)

**Zero-Value Handling**: Both generators omit zero values → Equivalent output

**Round-Trip Validation**:
```go
// Round-trip test pattern (from lrtemplate_test.go:1243-1260)
originalXMP := parseXMP("preset.xmp")
lrtemplateData := generateLrtemplate(originalXMP)
recoveredRecipe := parseLrtemplate(lrtemplateData)

// Assert: All 50+ parameters match exactly
assert.Equal(originalXMP.Exposure, recoveredRecipe.Exposure)
assert.Equal(originalXMP.Contrast, recoveredRecipe.Contrast)
// ... all other parameters
```

**Edge Cases**:
- Empty XMP (no parameters) → Empty lrtemplate (only structure)
- XMP with ProcessVersion mismatch → Normalized to PV2012 in UniversalRecipe

---

#### 5. lrtemplate → NP3 (Approximation, HIGH Data Loss)

**Process**: Parse lrtemplate → UniversalRecipe → Generate NP3

**Mapping**: Identical to Path #3 (XMP → NP3) - same data loss

**Example**: See Path #3 example (same behavior)

---

#### 6. lrtemplate → XMP (Direct 1:1, No Loss)

**Process**: Parse lrtemplate → UniversalRecipe → Generate XMP

**Mapping**: Direct field-name mapping (no transformation)

**Accuracy**: 100% - Perfect round-trip accuracy

**Round-Trip Validation**:
```go
// Round-trip test achieves 100% accuracy
originalLrtemplate := parseLrtemplate("preset.lrtemplate")
xmpData := generateXMP(originalLrtemplate)
recoveredRecipe := parseXMP(xmpData)

// Assert: Exact match for all parameters
assert.Equal(originalLrtemplate, recoveredRecipe)
```

---

### Conversion Path Selection Guide

**Use Case**: "I have an NP3 preset, which format should I convert to?"
- **Answer**: XMP or lrtemplate (equivalent output)
- **Reason**: No data loss, all NP3 parameters map cleanly

**Use Case**: "I have an XMP preset, which format should I convert to?"
- **To lrtemplate**: ✅ Perfect conversion (1:1 mapping)
- **To NP3**: ⚠️  Only if you accept losing 40+ advanced parameters

**Use Case**: "I have a lrtemplate preset, which format should I convert to?"
- **To XMP**: ✅ Perfect conversion (1:1 mapping)
- **To NP3**: ⚠️  Only if you accept losing 40+ advanced parameters

**Use Case**: "I need to edit an NP3 preset with advanced features"
1. Convert NP3 → XMP (no loss)
2. Edit in Lightroom (access all 50+ parameters)
3. Save as XMP or lrtemplate
4. **Do NOT convert back to NP3** (will lose advanced edits)

---

## Visual Similarity Validation

To ensure conversions maintain ≥95% visual similarity, the Recipe engine uses weighted parameter comparison with tolerance thresholds.

### Similarity Metrics

#### Parameter Delta Calculation

For each parameter, calculate the delta between original and round-trip values:

```go
// Integer parameter delta
func intDelta(original, roundtrip int) int {
    return abs(original - roundtrip)
}

// Float parameter delta
func floatDelta(original, roundtrip float64) float64 {
    return abs(original - roundtrip)
}
```

#### Weighted Importance

Not all parameters contribute equally to visual similarity. Critical parameters (Exposure, Contrast) have higher weights.

| Parameter | Weight | Rationale |
|-----------|--------|-----------|
| Exposure | 1.0 | **Most critical** - Affects overall brightness |
| Contrast | 1.0 | **Most critical** - Affects tonal distribution |
| Saturation | 0.8 | Very important - Affects color intensity |
| Highlights | 0.7 | Important - Affects bright tones |
| Shadows | 0.7 | Important - Affects dark tones |
| Whites | 0.6 | Important - Affects brightest areas |
| Blacks | 0.6 | Important - Affects darkest areas |
| Vibrance | 0.6 | Important - Affects color selectively |
| Clarity | 0.6 | Important - Affects local contrast |
| Sharpness | 0.5 | Moderately important - Affects edge definition |
| Texture | 0.5 | Moderately important - Affects fine detail |
| Dehaze | 0.5 | Moderately important - Affects atmospheric haze |
| HSL (per color) | 0.4 | Moderately important - Affects specific colors |
| Tone Curves | 0.7 | Important - Affects tonal mapping |
| Split Toning | 0.3 | Less important - Affects color grading |
| Grain | 0.2 | Less visible - Adds texture |
| Vignette | 0.3 | Less visible - Affects corners |
| Temperature | 0.6 | Important - Affects color cast |
| Tint | 0.4 | Moderately important - Affects green/magenta shift |

#### Tolerance Thresholds

Define acceptable deltas for integer and float parameters:

| Parameter Type | Tolerance | Reasoning |
|----------------|-----------|-----------|
| Integer parameters | ±1 | Rounding errors acceptable |
| Float parameters | ±0.01 | Sub-pixel precision acceptable |
| Approximated NP3 parameters | ±33 (Contrast, Saturation) | Scale factor precision |
| Approximated NP3 parameters | ±10 (Sharpness) | Scale factor precision |

### Similarity Score Calculation

#### Formula

```
Similarity = 100% - (Σ (weight_i × delta_i / tolerance_i) / Σ weight_i) × 100%
```

Where:
- `weight_i` = Importance weight for parameter i
- `delta_i` = Absolute difference between original and round-trip value
- `tolerance_i` = Acceptable threshold for parameter i

#### Code Example

```go
func calculateSimilarity(original, roundtrip *models.UniversalRecipe) float64 {
    totalWeightedDelta := 0.0
    totalWeight := 0.0

    // Exposure (weight: 1.0, tolerance: 0.01)
    if original.Exposure != 0 || roundtrip.Exposure != 0 {
        delta := math.Abs(original.Exposure - roundtrip.Exposure)
        normalizedDelta := delta / 0.01
        totalWeightedDelta += 1.0 * normalizedDelta
        totalWeight += 1.0
    }

    // Contrast (weight: 1.0, tolerance: 1)
    if original.Contrast != 0 || roundtrip.Contrast != 0 {
        delta := float64(abs(original.Contrast - roundtrip.Contrast))
        normalizedDelta := delta / 1.0
        totalWeightedDelta += 1.0 * normalizedDelta
        totalWeight += 1.0
    }

    // Saturation (weight: 0.8, tolerance: 1)
    if original.Saturation != 0 || roundtrip.Saturation != 0 {
        delta := float64(abs(original.Saturation - roundtrip.Saturation))
        normalizedDelta := delta / 1.0
        totalWeightedDelta += 0.8 * normalizedDelta
        totalWeight += 0.8
    }

    // ... repeat for all parameters

    // Calculate similarity percentage
    avgNormalizedDelta := totalWeightedDelta / totalWeight
    similarity := 100.0 - (avgNormalizedDelta * 100.0)

    // Clamp to 0-100% range
    if similarity < 0 {
        similarity = 0
    }
    if similarity > 100 {
        similarity = 100
    }

    return similarity
}
```

### Test Scenarios for ≥95% Similarity

#### Scenario 1: XMP ↔ lrtemplate Round-Trip (Direct Mapping)

**Expected Similarity**: 100%

```go
func TestXMPToLrtemplateRoundTrip(t *testing.T) {
    original := &models.UniversalRecipe{
        Exposure:   1.5,
        Contrast:   25,
        Saturation: -15,
        Sharpness:  80,
        // ... all other parameters
    }

    // Round-trip: UniversalRecipe → lrtemplate → UniversalRecipe
    lrtemplateData, _ := lrtemplate.Generate(original)
    recovered, _ := lrtemplate.Parse(lrtemplateData)

    similarity := calculateSimilarity(original, recovered)
    assert.GreaterOrEqual(t, similarity, 100.0, "XMP ↔ lrtemplate should be 100% similar")
}
```

**Result**: 100% similarity (no precision loss, direct field mapping)

---

#### Scenario 2: NP3 Round-Trip (Approximation Mapping)

**Expected Similarity**: ≥95% (within tolerance for NP3 approximation)

```go
func TestNP3RoundTrip(t *testing.T) {
    original := &models.UniversalRecipe{
        Exposure:   0.5,   // Within NP3 range (-1.0 to +1.0)
        Contrast:   66,    // 2 × 33 (NP3 level 2)
        Saturation: -33,   // -1 × 33 (NP3 level -1)
        Sharpness:  50,    // 5 × 10 (NP3 level 5)
    }

    // Round-trip: UniversalRecipe → NP3 → UniversalRecipe
    np3Data, _ := np3.Generate(original)
    recovered, _ := np3.Parse(np3Data)

    similarity := calculateSimilarity(original, recovered)
    assert.GreaterOrEqual(t, similarity, 95.0, "NP3 round-trip should be ≥95% similar")

    // Verify specific tolerances
    assert.InDelta(t, original.Exposure, recovered.Exposure, 0.01, "Exposure delta ≤0.01")
    assert.InDelta(t, original.Contrast, recovered.Contrast, 1, "Contrast delta ≤1")
    assert.InDelta(t, original.Saturation, recovered.Saturation, 1, "Saturation delta ≤1")
    assert.InDelta(t, original.Sharpness, recovered.Sharpness, 1, "Sharpness delta ≤1")
}
```

**Result**: 97-100% similarity (perfect round-trip for values aligned with NP3 scale factors)

---

#### Scenario 3: XMP → NP3 → XMP (Clamping and Unmappable Parameters)

**Expected Similarity**: 60-80% (significant data loss)

```go
func TestXMPToNP3WithDataLoss(t *testing.T) {
    original := &models.UniversalRecipe{
        Exposure:   2.5,    // Will be clamped to 1.0 in NP3
        Contrast:   75,     // Will round to 66 (NP3 level 2)
        Highlights: -50,    // Unmappable to NP3 (LOST)
        Texture:    25,     // Unmappable to NP3 (LOST)
        Red: models.ColorAdjustment{
            Hue: 15,        // Unmappable to NP3 (LOST)
        },
    }

    // Conversion: UniversalRecipe → NP3 → UniversalRecipe
    np3Data, warnings := np3.Generate(original)
    recovered, _ := np3.Parse(np3Data)

    // Verify warnings generated
    assert.NotEmpty(t, warnings, "Should warn about unmappable parameters")
    assert.Contains(t, warnings[0].Message, "HSL adjustments")

    // Calculate similarity
    similarity := calculateSimilarity(original, recovered)

    // Similarity will be reduced due to:
    // - Exposure clamped: 2.5 → 1.0 (major delta, weight 1.0)
    // - Contrast rounded: 75 → 66 (minor delta, weight 1.0)
    // - Highlights lost: -50 → 0 (moderate delta, weight 0.7)
    // - Texture lost: 25 → 0 (moderate delta, weight 0.5)
    // - HSL Red lost: 15 → 0 (minor delta, weight 0.4)

    assert.Less(t, similarity, 95.0, "Should be <95% due to clamping and data loss")
    assert.GreaterOrEqual(t, similarity, 60.0, "Should be ≥60% (basic params preserved)")
}
```

**Result**: 70-75% similarity (Exposure clamping has major impact, unmappable parameters reduce score)

---

#### Scenario 4: Measuring Similarity with Unmappable Parameters Present

When unmappable parameters exist in the original but are lost in conversion, they should contribute to the similarity score reduction.

**Strategy**: Include unmappable parameters in similarity calculation with weight 0 in round-trip value:

```go
func calculateSimilarityWithUnmappable(original, roundtrip *models.UniversalRecipe, warnings []ConversionWarning) float64 {
    similarity := calculateSimilarity(original, roundtrip)

    // Reduce similarity based on unmappable parameter warnings
    unmappableDeduction := 0.0
    for _, warning := range warnings {
        switch warning.Severity {
        case "critical":
            unmappableDeduction += 5.0  // HSL, tone curves
        case "advisory":
            unmappableDeduction += 3.0  // Highlights, Texture, etc.
        case "informational":
            unmappableDeduction += 1.0  // Grain, vignette
        }
    }

    // Cap deduction at 40% (to avoid negative scores)
    if unmappableDeduction > 40.0 {
        unmappableDeduction = 40.0
    }

    finalSimilarity := similarity - unmappableDeduction
    if finalSimilarity < 0 {
        finalSimilarity = 0
    }

    return finalSimilarity
}
```

---

### Continuous Validation

#### Test Data

Use real-world sample files for validation:

- `testdata/np3/*.np3` (22 samples) - Test NP3 round-trip accuracy
- `testdata/xmp/*.xmp` (913 samples) - Test XMP ↔ lrtemplate accuracy
- `testdata/lrtemplate/*.lrtemplate` (544 samples) - Test lrtemplate ↔ XMP accuracy

#### Benchmark Tests

```go
func BenchmarkSimilarityCalculation(b *testing.B) {
    original := loadTestRecipe("testdata/xmp/sample.xmp")
    roundtrip := loadTestRecipe("testdata/xmp/sample_roundtrip.xmp")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        calculateSimilarity(original, roundtrip)
    }
}
```

**Performance Target**: <1ms per similarity calculation

---

## Metadata Dictionary Usage

The `UniversalRecipe.Metadata` map preserves unmappable format-specific data, enabling round-trip conversions without total data loss.

### When to Use Metadata Map

| Scenario | Use Metadata? | Rationale |
|----------|---------------|-----------|
| Unmappable parameters (XMP → NP3) | ✅ Yes | Preserve HSL, tone curves, etc. for potential round-trip |
| Format-specific raw data (NP3 color data) | ✅ Yes | Preserve binary structures for exact reconstruction |
| Temporary conversion artifacts | ❌ No | Use local variables instead |
| Debugging information | ❌ No | Use logging, not persistent metadata |

### Key Naming Convention

**Format**: `{format_prefix}_{field_name}`

**Examples**:
- `xmp_tone_curve_pv2012` - XMP tone curve data
- `xmp_hsl_red_hue` - XMP HSL red hue adjustment
- `lrtemplate_split_toning_balance` - lrtemplate split toning balance
- `np3_unknown_bytes_64_80` - NP3 raw parameter bytes (debug)

**Rationale**: Format prefix prevents key collisions when round-tripping through multiple formats.

### Serialization Format

#### Simple Values (Integers, Floats, Strings)

```go
// Store simple values directly
recipe.Metadata["xmp_highlights"] = -50
recipe.Metadata["xmp_texture"] = 25
recipe.Metadata["lrtemplate_preset_name"] = "Cinematic Look"
```

#### Complex Structures (Arrays, Nested Objects)

Use JSON serialization for complex structures:

```go
import "encoding/json"

// Store tone curve (array of points)
curveJSON, _ := json.Marshal(recipe.PointCurve)
recipe.Metadata["xmp_tone_curve_pv2012"] = string(curveJSON)

// Store HSL adjustments (struct)
hslJSON, _ := json.Marshal(recipe.Red)
recipe.Metadata["xmp_hsl_red"] = string(hslJSON)

// Store split toning (multiple fields)
splitToningJSON, _ := json.Marshal(map[string]int{
    "shadowHue":        recipe.SplitShadowHue,
    "shadowSaturation": recipe.SplitShadowSaturation,
    "highlightHue":     recipe.SplitHighlightHue,
    "highlightSaturation": recipe.SplitHighlightSaturation,
    "balance":          recipe.SplitBalance,
})
recipe.Metadata["xmp_split_toning"] = string(splitToningJSON)
```

### Metadata Lifecycle

#### Phase 1: Add During Parse

When parsing XMP/lrtemplate that will be converted to NP3, store unmappable parameters:

```go
// internal/formats/xmp/parse.go
func Parse(data []byte) (*models.UniversalRecipe, error) {
    recipe := &models.UniversalRecipe{}

    // Parse basic parameters (always mapped)
    recipe.Exposure = parseFloat(doc, "Exposure2012")
    recipe.Contrast = parseInt(doc, "Contrast2012")

    // Parse HSL (may be unmappable to NP3)
    recipe.Red.Hue = parseInt(doc, "HueAdjustmentRed")
    // Also store in metadata for NP3 round-trip preservation
    recipe.Metadata["xmp_hsl_red_hue"] = parseInt(doc, "HueAdjustmentRed")

    // Parse tone curves (unmappable to NP3)
    recipe.PointCurve = parseToneCurve(doc, "ToneCurvePV2012")
    curveJSON, _ := json.Marshal(recipe.PointCurve)
    recipe.Metadata["xmp_tone_curve_pv2012"] = string(curveJSON)

    return recipe, nil
}
```

#### Phase 2: Preserve During Generate

When generating NP3, preserve metadata even though parameters can't be represented:

```go
// internal/formats/np3/generate.go
func Generate(recipe *models.UniversalRecipe) ([]byte, []ConversionWarning, error) {
    // Generate NP3 binary (only 5 basic parameters)
    data := encodeBinary(recipe)

    // Metadata is preserved in UniversalRecipe but not written to NP3 binary
    // (NP3 format has no metadata section)

    // Generate warnings for unmappable parameters
    warnings := []ConversionWarning{}
    if recipe.Metadata["xmp_tone_curve_pv2012"] != nil {
        warnings = append(warnings, ConversionWarning{
            Severity: "critical",
            Message:  "NP3 does not support tone curves - adjustment will be lost",
        })
    }

    return data, warnings, nil
}
```

#### Phase 3: Retrieve During Round-Trip

When converting back from NP3 to XMP, retrieve metadata:

```go
// internal/formats/xmp/generate.go
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
    // Generate XMP from basic parameters
    writeField(buf, "Exposure2012", recipe.Exposure)
    writeField(buf, "Contrast2012", recipe.Contrast)

    // Retrieve tone curve from metadata if present
    if curveData, ok := recipe.Metadata["xmp_tone_curve_pv2012"]; ok {
        var curve []models.ToneCurvePoint
        json.Unmarshal([]byte(curveData.(string)), &curve)
        writeToneCurve(buf, "ToneCurvePV2012", curve)
    }

    // Retrieve HSL from metadata if present
    if redHue, ok := recipe.Metadata["xmp_hsl_red_hue"]; ok {
        writeField(buf, "HueAdjustmentRed", redHue.(int))
    }

    return buf.Bytes(), nil
}
```

#### Phase 4: Warn User If Present

When metadata exists, inform user that data was preserved for round-trip:

```go
func generateWithMetadataWarning(recipe *models.UniversalRecipe) ([]byte, []string, error) {
    data, err := Generate(recipe)

    messages := []string{}
    if len(recipe.Metadata) > 0 {
        messages = append(messages, fmt.Sprintf(
            "ℹ️  %d unmappable parameters preserved in metadata for round-trip conversion",
            len(recipe.Metadata),
        ))
    }

    return data, messages, err
}
```

### Code Examples for Each Unmappable Parameter Type

#### Example 1: Tone Curves (Complex Array)

```go
// Store during XMP parse
curveJSON, _ := json.Marshal(recipe.PointCurve)
recipe.Metadata["xmp_tone_curve_pv2012"] = string(curveJSON)

// Retrieve during XMP generate
if curveData, ok := recipe.Metadata["xmp_tone_curve_pv2012"]; ok {
    var curve []models.ToneCurvePoint
    json.Unmarshal([]byte(curveData.(string)), &curve)
    writeToneCurveToXML(buf, curve)
}
```

#### Example 2: HSL Adjustments (Struct Array)

```go
// Store all 8 colors during XMP parse
for _, colorName := range []string{"Red", "Orange", "Yellow", "Green", "Aqua", "Blue", "Purple", "Magenta"} {
    adjustment := getColorAdjustment(recipe, colorName)
    adjJSON, _ := json.Marshal(adjustment)
    recipe.Metadata[fmt.Sprintf("xmp_hsl_%s", strings.ToLower(colorName))] = string(adjJSON)
}

// Retrieve during XMP generate
for _, colorName := range []string{"Red", "Orange", "Yellow", "Green", "Aqua", "Blue", "Purple", "Magenta"} {
    key := fmt.Sprintf("xmp_hsl_%s", strings.ToLower(colorName))
    if adjData, ok := recipe.Metadata[key]; ok {
        var adj models.ColorAdjustment
        json.Unmarshal([]byte(adjData.(string)), &adj)
        writeHSLToXML(buf, colorName, adj)
    }
}
```

#### Example 3: Simple Unmappable Parameters (Integer/Float)

```go
// Store during XMP parse
recipe.Metadata["xmp_highlights"] = recipe.Highlights
recipe.Metadata["xmp_shadows"] = recipe.Shadows
recipe.Metadata["xmp_texture"] = recipe.Texture

// Retrieve during XMP generate
if highlights, ok := recipe.Metadata["xmp_highlights"]; ok {
    writeField(buf, "Highlights2012", highlights.(int))
}
if shadows, ok := recipe.Metadata["xmp_shadows"]; ok {
    writeField(buf, "Shadows2012", shadows.(int))
}
if texture, ok := recipe.Metadata["xmp_texture"]; ok {
    writeField(buf, "Texture", texture.(int))
}
```

#### Example 4: Split Toning (Multiple Related Fields)

```go
// Store as grouped object during lrtemplate parse
splitToningJSON, _ := json.Marshal(map[string]int{
    "shadowHue":        recipe.SplitShadowHue,
    "shadowSaturation": recipe.SplitShadowSaturation,
    "highlightHue":     recipe.SplitHighlightHue,
    "highlightSaturation": recipe.SplitHighlightSaturation,
    "balance":          recipe.SplitBalance,
})
recipe.Metadata["lrtemplate_split_toning"] = string(splitToningJSON)

// Retrieve during lrtemplate generate
if splitData, ok := recipe.Metadata["lrtemplate_split_toning"]; ok {
    var split map[string]int
    json.Unmarshal([]byte(splitData.(string)), &split)
    writeLuaField(buf, "SplitToningShadowHue", split["shadowHue"])
    writeLuaField(buf, "SplitToningShadowSaturation", split["shadowSaturation"])
    writeLuaField(buf, "SplitToningHighlightHue", split["highlightHue"])
    writeLuaField(buf, "SplitToningHighlightSaturation", split["highlightSaturation"])
    writeLuaField(buf, "SplitToningBalance", split["balance"])
}
```

---

## Error Reporting and Warnings

Consistent error reporting across all conversion operations using the `ConversionError` type.

### ConversionError Structure

```go
// Defined in internal/formats/*/errors.go (each format package)
type ConversionError struct {
    Operation string  // "parse" or "generate"
    Format    string  // "np3", "xmp", "lrtemplate"
    Cause     error   // Underlying error
}

func (e *ConversionError) Error() string {
    return fmt.Sprintf("%s %s: %v", e.Operation, e.Format, e.Cause)
}
```

### Error Categories

| Category | Operation | Format | Cause | Example |
|----------|-----------|--------|-------|---------|
| **Validation** | parse | any | Invalid file structure | "parse np3: file too small: got 100 bytes, minimum 300 bytes required" |
| **Validation** | parse | any | Invalid magic bytes | "parse np3: invalid magic bytes: expected \"NCP\", got \"XMP\"" |
| **Validation** | parse | any | Parameter out of range | "parse xmp: validate contrast: value 150 out of range [-100, 100]" |
| **Validation** | generate | any | Nil recipe | "generate lrtemplate: recipe is nil" |
| **Encoding** | generate | any | Encoding failure | "generate xmp: xml marshal failed: invalid UTF-8" |
| **Unmappable** | generate | np3 | Advanced parameters | "generate np3: HSL adjustments cannot be represented (critical data loss)" |

### Warning Priority Levels

| Level | Description | Visual Severity | When to Use |
|-------|-------------|-----------------|-------------|
| **Critical** | High visual impact, significant data loss | 🔴 Red | HSL adjustments, tone curves, extreme value clamping |
| **Advisory** | Moderate visual impact, some data loss | 🟡 Yellow | Highlights/Shadows, Texture, Clarity, Dehaze |
| **Informational** | Low visual impact, minor data loss | 🔵 Blue | Grain, vignette, sharpening detail |

### ConversionWarning Structure

```go
type ConversionWarning struct {
    Severity string  // "critical", "advisory", "informational"
    Message  string  // User-friendly explanation
    Category string  // "unmappable", "clamping", "approximation"
    Details  map[string]interface{}  // Additional context
}
```

### User-Friendly Explanations

#### Why Parameters Can't Be Mapped

**HSL Adjustments**:
> "NP3 format was designed for simple global adjustments and only supports a single hue shift (-9° to +9°) affecting all colors uniformly. XMP and lrtemplate allow independent hue, saturation, and luminance adjustments for 8 individual colors (Red, Orange, Yellow, Green, Aqua, Blue, Purple, Magenta). This per-color control cannot be simplified to a global adjustment without losing the specific color grading intent."

**Tone Curves**:
> "NP3 uses a fixed lookup table for tonal adjustments with no support for custom tone curves. XMP and lrtemplate allow defining custom curves with multiple control points to precisely map input tones to output tones. These curves enable advanced creative looks (crushed blacks, S-curves, film emulation) that cannot be represented in NP3's simple contrast parameter."

**Highlights/Shadows**:
> "NP3 has only a single Contrast parameter that affects the entire tonal range uniformly. XMP and lrtemplate provide separate Highlights and Shadows controls for independent adjustment of bright and dark areas. When converting to NP3, these are approximated by the Contrast value, but the independent control is lost."

#### When to Fail Conversion vs Proceed with Warnings

**Default Behavior**: **Proceed with warnings**

```go
func Generate(recipe *models.UniversalRecipe) ([]byte, []ConversionWarning, error) {
    data, err := encodeBinary(recipe)
    if err != nil {
        return nil, nil, &ConversionError{
            Operation: "generate",
            Format:    "np3",
            Cause:     err,
        }
    }

    // Collect warnings (non-fatal)
    warnings := collectUnmappableWarnings(recipe)

    // Return data + warnings, allowing user to decide
    return data, warnings, nil
}
```

**Strict Mode**: **Fail on unmappable parameters**

```go
func GenerateStrict(recipe *models.UniversalRecipe) ([]byte, error) {
    // Check for unmappable parameters BEFORE generating
    if hasUnmappableParameters(recipe) {
        return nil, &ConversionError{
            Operation: "generate",
            Format:    "np3",
            Cause:     fmt.Errorf("preset contains unmappable parameters (use --no-strict to ignore)"),
        }
    }

    data, _, err := Generate(recipe)
    return data, err
}
```

---

## Implementation Guidance

### Where to Implement Mapping Logic

**Rule**: Mapping logic belongs in `generate.go` files, NOT in a centralized converter.

**Rationale**: Each format generator is responsible for converting from `UniversalRecipe` (hub) to its target format, including:
- Range conversions (approximation formulas)
- Format-specific encoding (XML, Lua, binary)
- Unmappable parameter handling (warnings, metadata)

**Pattern**:
```
internal/formats/
├── np3/
│   ├── parse.go       ← Maps NP3 → UniversalRecipe
│   └── generate.go    ← Maps UniversalRecipe → NP3 (approximation formulas HERE)
├── xmp/
│   ├── parse.go       ← Maps XMP → UniversalRecipe
│   └── generate.go    ← Maps UniversalRecipe → XMP (1:1 mapping HERE)
└── lrtemplate/
    ├── parse.go       ← Maps lrtemplate → UniversalRecipe
    └── generate.go    ← Maps UniversalRecipe → lrtemplate (1:1 mapping HERE)
```

**Anti-Pattern** (DO NOT DO):
```go
// ❌ BAD: Centralized converter with format-specific logic
package converter

func Convert(recipe *UniversalRecipe, targetFormat string) ([]byte, error) {
    switch targetFormat {
    case "np3":
        // NP3-specific logic here...
    case "xmp":
        // XMP-specific logic here...
    }
}
```

**Correct Pattern**:
```go
// ✅ GOOD: Format-specific generator
package np3

func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
    // NP3-specific approximation formulas
    params := convertToNP3Parameters(recipe)
    return encodeBinary(params)
}
```

### Constants to Define

Define scaling constants at package level in `generate.go`:

```go
// internal/formats/np3/generate.go
const (
    // Approximation scale factors for NP3 ↔ UniversalRecipe conversion
    NP3_CONTRAST_SCALE    = 33  // Maps NP3 ±3 to Universal ±99
    NP3_SATURATION_SCALE  = 33  // Maps NP3 ±3 to Universal ±99
    NP3_SHARPNESS_SCALE   = 10  // Maps NP3 0-9 to Universal 0-90
    NP3_BRIGHTNESS_SCALE  = 1.0 // Maps NP3 ±1.0 to Universal Exposure ±1.0
)
```

**Rationale**: Named constants make formulas self-documenting and easier to adjust if scale factors change.

### Pseudo-Code for Approximation Algorithms

#### Algorithm 1: Integer Range Scaling (Contrast, Saturation)

```
Input: universalValue (integer, range -100 to +100)
Output: np3Value (integer, range -3 to +3)

ALGORITHM ScaleToNP3Integer(universalValue, scale=33):
    np3Value ← universalValue / scale
    np3Value ← CLAMP(np3Value, -3, +3)
    RETURN np3Value

EXAMPLE:
    ScaleToNP3Integer(75, 33) → 75/33=2.27 → 2 (rounded down)
    ScaleToNP3Integer(-99, 33) → -99/33=-3.0 → -3
    ScaleToNP3Integer(120, 33) → 120/33=3.64 → 3 (clamped)
```

#### Algorithm 2: Integer Range Expansion (Contrast, Saturation)

```
Input: np3Value (integer, range -3 to +3)
Output: universalValue (integer, range -99 to +99)

ALGORITHM ScaleFromNP3Integer(np3Value, scale=33):
    universalValue ← np3Value * scale
    RETURN universalValue

EXAMPLE:
    ScaleFromNP3Integer(2, 33) → 2*33=66
    ScaleFromNP3Integer(-3, 33) → -3*33=-99
```

#### Algorithm 3: Float Clamping (Brightness/Exposure)

```
Input: universalValue (float, range -5.0 to +5.0)
Output: np3Value (float, range -1.0 to +1.0)

ALGORITHM ClampToNP3Float(universalValue, min=-1.0, max=1.0):
    IF universalValue < min THEN
        RETURN min
    ELSE IF universalValue > max THEN
        RETURN max
    ELSE
        RETURN universalValue

EXAMPLE:
    ClampToNP3Float(0.5, -1.0, 1.0) → 0.5
    ClampToNP3Float(2.5, -1.0, 1.0) → 1.0 (clamped)
    ClampToNP3Float(-3.0, -1.0, 1.0) → -1.0 (clamped)
```

### Reference Existing Implementations

| Format | File | Lines | Purpose |
|--------|------|-------|---------|
| NP3 Parse | `internal/formats/np3/parse.go` | 398-401 | Shows approximation expansion (NP3 → Universal) |
| NP3 Generate | `internal/formats/np3/generate.go` | 66-108 | Shows approximation compression (Universal → NP3) |
| XMP Generate | `internal/formats/xmp/generate.go` | entire | Shows direct 1:1 mapping (Universal → XMP) |
| lrtemplate Generate | `internal/formats/lrtemplate/generate.go` | 114-250 | Shows complete parameter list and zero-value omission |
| Round-Trip Test | `internal/formats/lrtemplate/lrtemplate_test.go` | 1243-1260 | Shows similarity validation pattern |

### Common Mistakes to Avoid

#### Mistake 1: Using Magic Numbers Instead of Constants

```go
// ❌ BAD: Magic number 33 is not self-documenting
params.contrast = recipe.Contrast / 33

// ✅ GOOD: Named constant explains the purpose
params.contrast = recipe.Contrast / NP3_CONTRAST_SCALE
```

#### Mistake 2: Off-By-One Errors in Range Mapping

```go
// ❌ BAD: Incorrect range (0-90 instead of 0-99)
params.contrast = recipe.Contrast / 10  // 10 levels, but range is -99 to +99

// ✅ GOOD: Correct scale factor for range
params.contrast = recipe.Contrast / 33  // Maps ±99 to ±3
```

#### Mistake 3: Forgetting to Clamp Values

```go
// ❌ BAD: No clamping, out-of-range values cause invalid NP3 files
params.contrast = recipe.Contrast / 33

// ✅ GOOD: Clamp to valid NP3 range
params.contrast = recipe.Contrast / 33
if params.contrast > 3 {
    params.contrast = 3
} else if params.contrast < -3 {
    params.contrast = -3
}
```

#### Mistake 4: Inconsistent Rounding (Float → Int)

```go
// ❌ BAD: Inconsistent rounding behavior
params.contrast = int(recipe.Contrast / 33.0)  // Truncates toward zero

// ✅ GOOD: Consistent rounding behavior
params.contrast = recipe.Contrast / 33  // Integer division (Go default)
```

#### Mistake 5: Not Preserving Unmappable Parameters in Metadata

```go
// ❌ BAD: Silent data loss
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
    // HSL adjustments are ignored without warning or metadata preservation
    return encodeBinary(recipe)
}

// ✅ GOOD: Preserve in metadata and warn user
func Generate(recipe *models.UniversalRecipe) ([]byte, []ConversionWarning, error) {
    warnings := []ConversionWarning{}

    // Preserve HSL in metadata
    if recipe.Red.Hue != 0 {
        recipe.Metadata["xmp_hsl_red_hue"] = recipe.Red.Hue
        warnings = append(warnings, ConversionWarning{
            Severity: "critical",
            Message:  "HSL Red hue adjustment (+15) will be lost in NP3",
        })
    }

    return encodeBinary(recipe), warnings, nil
}
```

#### Mistake 6: Not Testing Round-Trip Accuracy

```go
// ❌ BAD: No validation of mapping accuracy
func TestNP3Generate(t *testing.T) {
    recipe := &models.UniversalRecipe{Contrast: 66}
    data, _ := Generate(recipe)
    assert.NotNil(t, data)  // Only checks that generation succeeded
}

// ✅ GOOD: Validate round-trip accuracy
func TestNP3RoundTrip(t *testing.T) {
    original := &models.UniversalRecipe{Contrast: 66}

    // Round-trip
    np3Data, _ := Generate(original)
    recovered, _ := Parse(np3Data)

    // Validate parameters match (within tolerance)
    assert.InDelta(t, original.Contrast, recovered.Contrast, 1)
}
```

---

## Conclusion

This parameter mapping documentation provides:

1. ✅ Complete mapping tables for all 50+ parameters
2. ✅ Approximation formulas with rationale for NP3 conversions
3. ✅ Unmappable parameter identification and handling strategies
4. ✅ Bidirectional conversion path specifications
5. ✅ Visual similarity validation strategy with weighted metrics
6. ✅ Metadata dictionary usage patterns and lifecycle
7. ✅ Error reporting and user warning guidelines
8. ✅ Implementation guidance with code examples and common mistakes

**Next Steps**: Use this documentation to:
- Implement new format parsers/generators (follow same patterns)
- Validate existing implementations (cross-reference with this spec)
- Guide future format additions (Epic 7: Documentation & Deployment)
- Support CLI/TUI/Web warning messages (Epics 2, 3, 4)

**References**:
- Architecture: `docs/architecture.md` (UniversalRecipe structure, hub-and-spoke pattern)
- Tech Spec: `docs/tech-spec-epic-1.md` (format specifications, requirements)
- PRD: `docs/PRD.md` (FR-1.6: Parameter Mapping & Approximation, ≥95% similarity goal)
- Code: `internal/formats/*/parse.go` and `generate.go` (actual implementations)
