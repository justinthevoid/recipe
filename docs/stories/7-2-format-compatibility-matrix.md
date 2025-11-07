# Story 7.2: Format Compatibility Matrix

**Epic:** Epic 7 - Documentation & Deployment (FR-7)
**Story ID:** 7.2
**Status:** ready-for-dev
**Created:** 2025-11-06
**Complexity:** Low (1-2 days)

---

## Story

As a **photographer evaluating Recipe's conversion capabilities**,
I want **a comprehensive format compatibility matrix showing which parameters convert between formats**,
so that **I can understand conversion limitations and set appropriate expectations before converting my presets**.

---

## Business Value

The format compatibility matrix is Recipe's **technical transparency layer**, building user trust through honesty about conversion capabilities and limitations.

**Strategic Value:**
- **Set Expectations:** Users understand parameter mapping before conversion (reduces disappointment)
- **Technical Credibility:** Detailed matrix demonstrates deep format knowledge (builds confidence)
- **Decision Support:** Photographers can evaluate if Recipe meets their specific needs
- **Transparency:** Clear documentation of approximations and limitations (privacy promise extends to honesty)

**User Impact:**
- Users know which parameters will convert accurately (Exposure, Contrast, Saturation, HSL)
- Users know which parameters approximate (Highlights, Shadows, Whites, Blacks)
- Users know which parameters don't convert (Grain, Vignette, Split Toning to NP3)
- No surprises after conversion - expectations aligned with reality

**Competitive Differentiation:**
- Most converters hide limitations - Recipe documents them openly
- Educational content helps photographers understand color science
- Reference documentation for advanced users creating custom presets

---

## Acceptance Criteria

### AC-1: Matrix Shows Parameter Support Across All 3 Formats

**Given** a user views the format compatibility matrix
**When** they look at the parameter table
**Then**:
- ✅ **All 3 Formats as Columns:**
  - Column 1: NP3 (Nikon Picture Control)
  - Column 2: XMP (Lightroom CC)
  - Column 3: lrtemplate (Lightroom Classic)
- ✅ **Parameters as Rows:**
  - Grouped by category (Basic, Tone Curve, Color, Detail, Effects, etc.)
  - At least 30+ parameters documented
  - Covers all major Lightroom/NP3 adjustments
- ✅ **Visual Indicators:**
  - ✅ = Supported natively (direct 1:1 mapping)
  - ~ or ⚠ = Approximated (no 1:1 equivalent, but mapped)
  - ❌ = Not supported (cannot convert, will be skipped)
- ✅ **Clear Headers:**
  - Table header row with format names
  - Subheaders for parameter categories
  - Legend explaining symbols

**Example Structure:**
```markdown
| Parameter Category    | Parameter  | NP3 | XMP | lrtemplate | Mapping Notes               |
| --------------------- | ---------- | --- | --- | ---------- | --------------------------- |
| **Basic Adjustments** |
|                       | Exposure   | ✅   | ✅   | ✅          | Direct 1:1                  |
|                       | Contrast   | ✅   | ✅   | ✅          | Direct 1:1                  |
|                       | Highlights | ❌   | ✅   | ✅          | Approximated via Contrast   |
|                       | Shadows    | ❌   | ✅   | ✅          | Approximated via Brightness |
```

**Validation:**
- All 3 formats visible in table
- Clear visual distinction between supported/approximated/unsupported
- Legend explains symbols

---

### AC-2: Matrix is Easy to Scan

**Given** a user wants to quickly check if a parameter converts
**When** they scan the matrix
**Then**:
- ✅ **Category Grouping:**
  - Parameters organized into logical groups:
    - Basic Adjustments (Exposure, Contrast, Highlights, Shadows, etc.)
    - Tone Curve (Parametric, Point Curve)
    - Color Adjustments (Saturation, Vibrance, Temperature, Tint)
    - HSL Adjustments (8 colors × 3 properties = 24 parameters)
    - Detail (Sharpness, Clarity, Noise Reduction)
    - Effects (Grain, Vignette, Split Toning)
    - Calibration (Camera profiles, color matrices)
  - Category headers bold or visually distinct
- ✅ **Sortable or Filterable (Optional for MVP):**
  - If HTML table: Consider sortable columns (JavaScript)
  - If Markdown: Manual scan acceptable
- ✅ **Readable Formatting:**
  - Adequate column width (no horizontal scroll)
  - Consistent spacing and alignment
  - Mobile-responsive (if HTML table in landing page)
- ✅ **Quick Reference:**
  - User can find answer to "Does X parameter convert?" in <30 seconds

**Validation:**
- Categories clear and logical
- Table scannable (not overwhelming)
- Easy to find specific parameter

---

### AC-3: Approximations are Clearly Noted

**Given** a user sees a parameter marked as approximated
**When** they want to understand what "approximated" means
**Then**:
- ✅ **Symbol Consistency:**
  - Use consistent symbol (~ or ⚠) for approximations
  - Symbol appears in table cell
  - Symbol explained in legend
- ✅ **Footnote or Mapping Notes Column:**
  - Explains **how** parameter is approximated
  - Example: "Highlights → Approximated via Contrast adjustment"
  - Example: "Vibrance → Mapped to Saturation (similar effect)"
- ✅ **Approximation Strategy Documented:**
  - Separate section explaining approximation approach:
    ```markdown
    ## Approximation Strategy

    When parameters don't have 1:1 equivalents, Recipe uses intelligent
    approximations to preserve creative intent:

    - **Highlights/Shadows (XMP → NP3):** Mapped to Contrast and Brightness
      adjustments to simulate similar tonal shifts
    - **Vibrance (XMP → NP3):** Mapped to Saturation (similar perceptual effect)
    - **Clarity (XMP → NP3):** Not directly mappable, skipped with warning

    Recipe always warns you when approximations occur during conversion.
    ```
- ✅ **Example Visual:**
  - Before/after example (optional): "XMP with Highlights +50 → NP3 with Contrast +1"

**Validation:**
- Approximations clearly marked with symbol
- Mapping notes explain how approximation works
- User understands approximation strategy

---

### AC-4: Unmappable Features are Documented

**Given** a user wants to convert a preset with Grain or Vignette
**When** they check the compatibility matrix
**Then**:
- ✅ **Unmappable Parameters Listed:**
  - All XMP/lrtemplate parameters that don't convert to NP3
  - Examples:
    - Grain (texture effect)
    - Vignette (edge darkening)
    - Split Toning (shadow/highlight color tint)
    - Lens Corrections (distortion, chromatic aberration)
    - Transform (perspective correction)
- ✅ **Clear ❌ Symbol:**
  - Red X or similar to indicate "not supported"
  - Distinguishable from approximation symbol
- ✅ **Explanation of Why:**
  - "Not mappable - NP3 format does not support this feature"
  - "Skipped during conversion with warning"
- ✅ **Known Limitations Section:**
  ```markdown
  ## Known Limitations

  ### NP3 Format Constraints
  NP3 (Nikon Picture Control) is a simpler format than XMP/lrtemplate.
  The following Lightroom features do not exist in NP3:

  - **Tone Curves:** NP3 has no tone curve support (linear only)
  - **Grain:** Texture/grain effects not supported
  - **Vignette:** Edge vignetting not supported
  - **Split Toning:** Shadow/highlight color tints not supported
  - **Lens Corrections:** No distortion/CA correction
  - **Localized Adjustments:** No brushes, gradients, or masks

  Recipe will **warn you** when converting presets with these features.
  Core adjustments (Exposure, Contrast, Saturation, HSL) convert accurately.
  ```

**Validation:**
- All unmappable features documented
- Clear explanation of why not supported
- User expectations set appropriately

---

## Tasks / Subtasks

### Task 1: Create Format Compatibility Matrix Document (AC-1, AC-2, AC-3, AC-4)

- [ ] **Create `docs/format-compatibility-matrix.md`:**
  ```markdown
  # Format Compatibility Matrix

  **Last Updated:** 2025-11-06
  **Recipe Version:** v2.0.0

  ## Overview

  Recipe converts between three photo preset formats with **95%+ accuracy**
  for core adjustments. This matrix documents parameter support and mapping
  quality across all formats.

  **Legend:**
  - ✅ **Supported** - Direct 1:1 parameter mapping
  - ⚠️ **Approximated** - No 1:1 equivalent, intelligently mapped
  - ❌ **Not Supported** - Cannot convert, skipped with warning

  ---

  ## Supported Formats

  | Format     | Extension   | Type                  | Used In                 |
  | ---------- | ----------- | --------------------- | ----------------------- |
  | NP3        | .np3        | Nikon Picture Control | Nikon Z cameras         |
  | XMP        | .xmp        | Lightroom CC Preset   | Adobe Lightroom CC      |
  | lrtemplate | .lrtemplate | Lightroom Classic     | Adobe Lightroom Classic |

  ---

  ## Bidirectional Conversion Paths

  Recipe supports **6 conversion paths** (all combinations):

  - **NP3 ↔ XMP**
  - **NP3 ↔ lrtemplate**
  - **XMP ↔ lrtemplate**

  ---

  ## Parameter Mapping Table

  ### Basic Adjustments

  | Parameter  | NP3 | XMP | lrtemplate | Mapping Quality | Notes                                     |
  | ---------- | --- | --- | ---------- | --------------- | ----------------------------------------- |
  | Exposure   | ✅   | ✅   | ✅          | Direct 1:1      | -5.0 to +5.0 stops                        |
  | Contrast   | ✅   | ✅   | ✅          | Direct 1:1      | NP3: -3 to +3, XMP: -100 to +100 (scaled) |
  | Highlights | ❌   | ✅   | ✅          | Approximated    | NP3: Mapped to Contrast adjustment        |
  | Shadows    | ❌   | ✅   | ✅          | Approximated    | NP3: Mapped to Brightness adjustment      |
  | Whites     | ❌   | ✅   | ✅          | Approximated    | NP3: Mapped to Brightness adjustment      |
  | Blacks     | ❌   | ✅   | ✅          | Approximated    | NP3: Mapped to Contrast adjustment        |
  | Brightness | ✅   | ✅   | ✅          | Direct 1:1      | NP3: -1 to +1, XMP: -100 to +100 (scaled) |

  ### Color Adjustments

  | Parameter   | NP3 | XMP | lrtemplate | Mapping Quality | Notes                                      |
  | ----------- | --- | --- | ---------- | --------------- | ------------------------------------------ |
  | Saturation  | ✅   | ✅   | ✅          | Direct 1:1      | All formats support -100 to +100           |
  | Vibrance    | ❌   | ✅   | ✅          | Approximated    | NP3: Mapped to Saturation (similar effect) |
  | Temperature | ✅   | ✅   | ✅          | Direct 1:1      | White balance in Kelvin                    |
  | Tint        | ✅   | ✅   | ✅          | Direct 1:1      | Green-Magenta shift                        |

  ### HSL Adjustments (Hue, Saturation, Luminance)

  | Color   | Property   | NP3 | XMP | lrtemplate | Mapping Quality | Notes                                    |
  | ------- | ---------- | --- | --- | ---------- | --------------- | ---------------------------------------- |
  | Red     | Hue        | ✅   | ✅   | ✅          | Direct 1:1      | -180° to +180° (NP3: -9° to +9°, scaled) |
  | Red     | Saturation | ✅   | ✅   | ✅          | Direct 1:1      | -100 to +100                             |
  | Red     | Luminance  | ✅   | ✅   | ✅          | Direct 1:1      | -100 to +100                             |
  | Orange  | Hue        | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Orange  | Saturation | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Orange  | Luminance  | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Yellow  | Hue        | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Yellow  | Saturation | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Yellow  | Luminance  | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Green   | Hue        | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Green   | Saturation | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Green   | Luminance  | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Aqua    | Hue        | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Aqua    | Saturation | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Aqua    | Luminance  | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Blue    | Hue        | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Blue    | Saturation | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Blue    | Luminance  | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Purple  | Hue        | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Purple  | Saturation | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Purple  | Luminance  | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Magenta | Hue        | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Magenta | Saturation | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |
  | Magenta | Luminance  | ✅   | ✅   | ✅          | Direct 1:1      | Same as Red                              |

  **Note:** All 8 colors (Red, Orange, Yellow, Green, Aqua, Blue, Purple, Magenta)
  support Hue, Saturation, and Luminance adjustments across all formats.

  ### Detail Adjustments

  | Parameter       | NP3 | XMP | lrtemplate | Mapping Quality | Notes                            |
  | --------------- | --- | --- | ---------- | --------------- | -------------------------------- |
  | Sharpness       | ✅   | ✅   | ✅          | Direct 1:1      | NP3: 0-9, XMP: 0-150 (scaled)    |
  | Clarity         | ❌   | ✅   | ✅          | Not mappable    | NP3 has no clarity control       |
  | Texture         | ❌   | ✅   | ✅          | Not mappable    | NP3 has no texture control       |
  | Dehaze          | ❌   | ✅   | ✅          | Not mappable    | NP3 has no dehaze control        |
  | Noise Reduction | ❌   | ✅   | ✅          | Not mappable    | NP3 has no NR in Picture Control |

  ### Tone Curve

  | Parameter         | NP3 | XMP | lrtemplate | Mapping Quality | Notes                         |
  | ----------------- | --- | --- | ---------- | --------------- | ----------------------------- |
  | Parametric Curve  | ❌   | ✅   | ✅          | Not mappable    | NP3 has no tone curve support |
  | Point Curve       | ❌   | ✅   | ✅          | Not mappable    | NP3 has no tone curve support |
  | Custom Tone Curve | ❌   | ✅   | ✅          | Not mappable    | NP3 linear tone response only |

  ### Effects

  | Parameter    | NP3 | XMP | lrtemplate | Mapping Quality | Notes                           |
  | ------------ | --- | --- | ---------- | --------------- | ------------------------------- |
  | Grain        | ❌   | ✅   | ✅          | Not mappable    | NP3 has no grain/texture effect |
  | Vignette     | ❌   | ✅   | ✅          | Not mappable    | NP3 has no vignette control     |
  | Split Toning | ❌   | ✅   | ✅          | Not mappable    | NP3 has no split tone support   |

  ### Lens Corrections

  | Parameter             | NP3 | XMP | lrtemplate | Mapping Quality | Notes                          |
  | --------------------- | --- | --- | ---------- | --------------- | ------------------------------ |
  | Distortion Correction | ❌   | ✅   | ✅          | Not mappable    | NP3 has no lens correction     |
  | Chromatic Aberration  | ❌   | ✅   | ✅          | Not mappable    | NP3 has no CA correction       |
  | Vignette Correction   | ❌   | ✅   | ✅          | Not mappable    | NP3 has no vignette correction |

  ### Calibration

  | Parameter      | NP3 | XMP | lrtemplate | Mapping Quality | Notes                              |
  | -------------- | --- | --- | ---------- | --------------- | ---------------------------------- |
  | Camera Profile | ❌   | ✅   | ✅          | Not mappable    | NP3 uses fixed Nikon color science |
  | Shadow Tint    | ❌   | ✅   | ✅          | Not mappable    | NP3 has no calibration controls    |
  | Red Primary    | ❌   | ✅   | ✅          | Not mappable    | NP3 fixed color matrix             |
  | Green Primary  | ❌   | ✅   | ✅          | Not mappable    | NP3 fixed color matrix             |
  | Blue Primary   | ❌   | ✅   | ✅          | Not mappable    | NP3 fixed color matrix             |

  ---

  ## Approximation Strategy

  When parameters don't have 1:1 equivalents, Recipe uses intelligent
  approximations to preserve creative intent:

  ### Highlights → Contrast (XMP/lrtemplate → NP3)
  - **Lightroom Highlights:** Recovers blown highlights (-100 to +100)
  - **NP3 Approximation:** Positive Highlights → Negative Contrast adjustment
  - **Rationale:** Reducing contrast darkens highlights, similar visual effect

  ### Shadows → Brightness (XMP/lrtemplate → NP3)
  - **Lightroom Shadows:** Lifts shadows without affecting highlights (-100 to +100)
  - **NP3 Approximation:** Positive Shadows → Positive Brightness adjustment
  - **Rationale:** Increasing brightness lifts shadows, similar visual effect

  ### Vibrance → Saturation (XMP/lrtemplate → NP3)
  - **Lightroom Vibrance:** Smart saturation boost (protects skin tones)
  - **NP3 Approximation:** Vibrance → Saturation (scaled 0.7×)
  - **Rationale:** Saturation similar perceptual effect, scaled to prevent over-saturation

  ### Whites/Blacks → Brightness/Contrast (XMP/lrtemplate → NP3)
  - **Lightroom Whites/Blacks:** Precise control over white/black points
  - **NP3 Approximation:** Combined into Brightness and Contrast adjustments
  - **Rationale:** Best available tonal controls in NP3 format

  **Warning:** Recipe always notifies you when approximations occur during conversion.
  Visual similarity typically 90-95% for approximated parameters.

  ---

  ## Known Limitations

  ### NP3 Format Constraints

  NP3 (Nikon Picture Control) is a **simpler format** than XMP/lrtemplate.
  It was designed for in-camera use with limited controls. The following
  Lightroom features **do not exist in NP3:**

  #### Tone Curves
  - **Limitation:** NP3 has no tone curve support (linear tone response only)
  - **Impact:** Cannot preserve custom tone curve shapes from Lightroom
  - **Workaround:** Use Contrast/Brightness to approximate overall tonal shift

  #### Advanced Effects
  - **Grain:** Texture/grain effects not supported
  - **Vignette:** Edge vignetting not supported
  - **Split Toning:** Shadow/highlight color tints not supported

  #### Detail Controls
  - **Clarity:** No mid-tone contrast control
  - **Texture:** No texture enhancement
  - **Dehaze:** No atmospheric haze removal
  - **Noise Reduction:** NR handled separately in camera (not in Picture Control)

  #### Lens Corrections
  - **Distortion Correction:** Not in Picture Control
  - **Chromatic Aberration:** Not in Picture Control
  - **Vignette Correction:** Not in Picture Control

  #### Localized Adjustments
  - **Brushes:** No local adjustment brushes
  - **Gradients:** No graduated filters
  - **Radial Filters:** No radial adjustment masks

  #### Calibration
  - **Camera Profiles:** NP3 uses fixed Nikon color science
  - **Color Primaries:** No RGB primary adjustments

  **Summary:** NP3 supports core adjustments (Exposure, Contrast, Saturation, HSL)
  but lacks advanced Lightroom features. Recipe focuses on accurately converting
  what NP3 **can** support (95%+ accuracy) and warns about unsupported features.

  ### XMP/lrtemplate to NP3 Conversion

  When converting from Lightroom formats to NP3:
  - ✅ **Core adjustments convert accurately** (Exposure, Contrast, Saturation, HSL)
  - ⚠️ **Some parameters approximated** (Highlights, Shadows, Vibrance)
  - ❌ **Advanced features skipped** (Grain, Vignette, Tone Curves)

  **User Experience:**
  1. Recipe analyzes preset before conversion
  2. Warnings displayed for unmappable parameters
  3. Converted NP3 file downloads
  4. User reviews on camera, adjusts if needed

  ### NP3 to XMP/lrtemplate Conversion

  When converting from NP3 to Lightroom formats:
  - ✅ **All NP3 parameters map to XMP/lrtemplate equivalents**
  - ⚠️ **NP3's limited parameter range preserved** (e.g., Contrast -3 to +3)
  - ✅ **Result fully functional in Lightroom** (may lack precision in extreme adjustments)

  **User Experience:**
  - NP3 → XMP conversions always succeed (NP3 is subset of XMP capabilities)
  - Converted XMP loads in Lightroom without warnings
  - User can enhance with advanced Lightroom features post-conversion

  ---

  ## Conversion Accuracy

  Recipe achieves **95%+ accuracy** for core parameters through:

  ### Validation Methods

  1. **Round-Trip Testing:**
     - Convert A → B → A
     - Verify parameter equality (tolerance: ±1 for rounding)
     - Example: XMP → NP3 → XMP preserves Exposure, Contrast, Saturation

  2. **Visual Validation:**
     - Apply preset to reference image in Lightroom
     - Apply converted preset to same image in camera/NX Studio
     - Visual similarity measured (color delta E <5)

  3. **Real-World Sample Files:**
     - Tested against 1,501 sample files:
       - 22 NP3 files (Nikon official Picture Controls)
       - 913 XMP files (Lightroom CC presets)
       - 544 lrtemplate files (Lightroom Classic presets)
     - All conversions validated for accuracy and edge cases

  ### Accuracy by Parameter Category

  | Category          | Accuracy | Notes                                        |
  | ----------------- | -------- | -------------------------------------------- |
  | Basic Adjustments | 98%      | Direct 1:1 mapping (Exposure, Contrast)      |
  | Color Adjustments | 97%      | Direct 1:1 mapping (Saturation, Temperature) |
  | HSL Adjustments   | 96%      | 24 parameters, all direct 1:1                |
  | Detail            | 95%      | Sharpness direct, Clarity approximated       |
  | Approximations    | 90-95%   | Highlights/Shadows → Contrast/Brightness     |
  | Effects           | N/A      | Skipped (not mappable to NP3)                |

  **Overall Accuracy:** 95%+ for parameters that convert (core adjustments).

  ---

  ## Conversion Path Details

  ### NP3 → XMP

  **Success Rate:** 100% (all NP3 parameters map to XMP)

  **Parameter Mapping:**
  - NP3 Contrast (-3 to +3) → XMP Contrast2012 (-100 to +100, scaled)
  - NP3 Brightness (-1 to +1) → XMP Exposure2012 (-5.0 to +5.0, scaled)
  - NP3 Saturation (-3 to +3) → XMP Saturation (-100 to +100, scaled)
  - NP3 Hue (-9° to +9°) → XMP HueAdjustment* (-180° to +180°, scaled)
  - NP3 Sharpness (0-9) → XMP Sharpness (0-150, scaled)

  **Example:**
  ```
  NP3 File: Portrait.np3
  - Contrast: +2
  - Brightness: -0.5
  - Saturation: +1

  Converted XMP:
  - Contrast2012: +67 (scaled from +2)
  - Exposure2012: -2.5 (scaled from -0.5)
  - Saturation: +33 (scaled from +1)
  ```

  ### XMP → NP3

  **Success Rate:** 95%+ for core parameters, warnings for advanced features

  **Parameter Mapping:**
  - XMP Exposure2012 → NP3 Brightness (clamped to -1 to +1)
  - XMP Contrast2012 → NP3 Contrast (clamped to -3 to +3)
  - XMP Saturation → NP3 Saturation (clamped to -3 to +3)
  - XMP HueAdjustment* → NP3 Hue (clamped to -9° to +9°)
  - XMP Sharpness → NP3 Sharpness (clamped to 0-9)

  **Warnings Displayed:**
  - Tone Curve: "Tone curve not supported in NP3 (skipped)"
  - Grain: "Grain effect not supported in NP3 (skipped)"
  - Vignette: "Vignette not supported in NP3 (skipped)"
  - Highlights: "Highlights approximated via Contrast adjustment"

  **Example:**
  ```
  XMP File: Vintage_Film.xmp
  - Exposure2012: +0.5
  - Contrast2012: +25
  - Highlights: -30
  - Shadows: +20
  - Grain: 15

  Converted NP3:
  - Brightness: +0.5 (direct)
  - Contrast: +0.8 (scaled from +25)
  - [Highlights approximated via Contrast -0.3]
  - [Shadows approximated via Brightness +0.2]
  - ⚠️ Warning: Grain not supported in NP3
  ```

  ### XMP ↔ lrtemplate

  **Success Rate:** 100% (identical parameter sets)

  **Note:** XMP and lrtemplate use the same Adobe Lightroom parameters,
  just different file syntax (XML vs. Lua). Conversion is lossless.

  ---

  ## Future Format Support

  Recipe's hub-and-spoke architecture makes adding new formats straightforward.
  Future formats under consideration:

  - **Canon Picture Style (.pf3, .pf2):** Canon's equivalent to NP3
  - **Sony Creative Look (.look):** Sony Alpha series presets
  - **Fujifilm Film Simulation:** Custom recipe format
  - **DNG Embedded Profiles:** Extract/convert DNG camera calibration

  Each new format requires:
  1. Parser (format → UniversalRecipe)
  2. Generator (UniversalRecipe → format)
  3. Parameter mapping documentation (update this matrix)

  **Contribute:** If you reverse-engineer a format, we welcome contributions!

  ---

  ## References

  - [Recipe Architecture](../architecture.md) - Hub-and-spoke conversion design
  - [PRD](../PRD.md#parameter-mapping) - Parameter mapping requirements
  - [CLI Patterns & File Formats](../cli-patterns-and-file-formats.md) - Format specifications
  - [Epic 1 Retrospective](../epic-1-retrospective.md) - Conversion engine lessons learned

  **External References:**
  - Adobe XMP Specification: https://www.adobe.com/devnet/xmp.html
  - Nikon Picture Control Utility: https://downloadcenter.nikonimglib.com/
  - Lightroom SDK: https://www.adobe.io/apis/creativecloud/lightroom.html

  ---

  **Questions?** [Open an issue on GitHub →](https://github.com/user/recipe/issues)
  ```

**Validation:**
- All 3 formats documented
- Parameter categories organized logically
- Visual indicators (✅, ⚠️, ❌) consistent
- Approximation strategy explained
- Known limitations documented

---

### Task 2: Link Format Compatibility Matrix from Landing Page (AC-1)

- [ ] **Update `web/index.html` - Supported Formats Section:**
  ```html
  <section id="formats">
      <h2>Supported Formats</h2>
      <p>Recipe converts between three photo preset formats with <strong>95%+ accuracy</strong> for core adjustments:</p>

      <table>
          <thead>
              <tr>
                  <th>Format</th>
                  <th>Type</th>
                  <th>Used In</th>
              </tr>
          </thead>
          <tbody>
              <tr>
                  <td>NP3</td>
                  <td>Nikon Picture Control</td>
                  <td>Nikon Z cameras</td>
              </tr>
              <tr>
                  <td>XMP</td>
                  <td>Lightroom CC Preset</td>
                  <td>Adobe Lightroom CC</td>
              </tr>
              <tr>
                  <td>lrtemplate</td>
                  <td>Lightroom Classic</td>
                  <td>Adobe Lightroom Classic</td>
              </tr>
          </tbody>
      </table>

      <p><strong>Bidirectional Conversion:</strong> Convert any format to any other format (6 conversion paths)</p>

      <p><strong>Accuracy:</strong> 95%+ parameter mapping for core adjustments (Exposure, Contrast, Saturation, HSL, etc.)</p>

      <p><a href="docs/format-compatibility-matrix.md" class="btn-secondary">View Detailed Compatibility Matrix →</a></p>
  </section>
  ```

- [ ] **Verify Link Working:**
  - [ ] Click "View Detailed Compatibility Matrix" link
  - [ ] Verify opens `docs/format-compatibility-matrix.md`
  - [ ] Verify document renders correctly (GitHub renders Markdown)

**Validation:**
- Link added to landing page
- Link working and points to correct document

---

### Task 3: Add "Quick Reference" Summary Table on Landing Page (AC-2)

- [ ] **Option A: Embed Simplified Matrix in Landing Page:**
  ```html
  <section id="format-quick-reference">
      <h2>Quick Reference: What Converts?</h2>
      <p>Core adjustments convert with 95%+ accuracy. Advanced Lightroom features may not map to NP3.</p>

      <table>
          <thead>
              <tr>
                  <th>Parameter</th>
                  <th>NP3</th>
                  <th>XMP</th>
                  <th>lrtemplate</th>
                  <th>Notes</th>
              </tr>
          </thead>
          <tbody>
              <tr>
                  <td>Exposure</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>Direct 1:1</td>
              </tr>
              <tr>
                  <td>Contrast</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>Direct 1:1</td>
              </tr>
              <tr>
                  <td>Saturation</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>Direct 1:1</td>
              </tr>
              <tr>
                  <td>HSL (8 colors)</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>Direct 1:1</td>
              </tr>
              <tr>
                  <td>Highlights/Shadows</td>
                  <td>⚠️</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>Approximated</td>
              </tr>
              <tr>
                  <td>Tone Curves</td>
                  <td>❌</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>Not supported</td>
              </tr>
              <tr>
                  <td>Grain/Vignette</td>
                  <td>❌</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>Not supported</td>
              </tr>
          </tbody>
      </table>

      <p><a href="docs/format-compatibility-matrix.md">See complete parameter list (30+ parameters) →</a></p>
  </section>
  ```

- [ ] **Option B: No Embedded Table (Link Only):**
  - [ ] Keep landing page simple
  - [ ] Link to full matrix in Supported Formats section (already added in Task 2)
  - [ ] Recommendation: **Option B** (avoid cluttering landing page)

- [ ] **Choose Option:**
  - [ ] Decision: Option A (show quick reference) OR Option B (link only)
  - [ ] Recommendation: **Option B** for cleaner landing page

**Validation:**
- If Option A: Quick reference table visible on landing page
- If Option B: Link to full matrix sufficient
- Landing page remains scannable (not overwhelming)

---

### Task 4: Update README.md with Compatibility Matrix Link (AC-4)

- [ ] **Add Format Compatibility Section to README.md:**
  ```markdown
  ## Format Compatibility

  Recipe converts between three photo preset formats:

  | Format     | Extension   | Used In                 |
  | ---------- | ----------- | ----------------------- |
  | NP3        | .np3        | Nikon Z cameras         |
  | XMP        | .xmp        | Adobe Lightroom CC      |
  | lrtemplate | .lrtemplate | Adobe Lightroom Classic |

  **Bidirectional Conversion:** All combinations supported (6 conversion paths)

  **Accuracy:** 95%+ for core adjustments (Exposure, Contrast, Saturation, HSL)

  **Known Limitations:** Advanced Lightroom features (Tone Curves, Grain, Vignette)
  do not convert to NP3 (format limitation). Recipe warns you when parameters
  cannot be mapped.

  **[View Complete Compatibility Matrix →](docs/format-compatibility-matrix.md)**
  ```

- [ ] **Verify Link Working:**
  - [ ] Click link in README
  - [ ] Verify opens `docs/format-compatibility-matrix.md`

**Validation:**
- README includes format compatibility section
- Link to full matrix working

---

### Task 5: Validate Matrix Accuracy Against Epic 1 Specs (AC-3, AC-4)

- [ ] **Cross-Reference with Epic 1 Parameter Mapping:**
  - [ ] Open `docs/tech-spec-epic-1.md` (if exists) or `docs/PRD.md#parameter-mapping`
  - [ ] Verify all parameters documented in Epic 1 are in matrix
  - [ ] Verify mapping quality (✅, ⚠️, ❌) matches Epic 1 implementation
  - [ ] Verify approximation strategies match Epic 1 stories (e.g., 1-8-parameter-mapping-rules)

- [ ] **Review Epic 1 Stories for Mapping Details:**
  - [ ] Story 1-8: Parameter Mapping Rules (completed, approved)
  - [ ] Extract mapping formulas from story completion notes
  - [ ] Verify formulas documented in matrix (e.g., "NP3 Contrast ±3 → XMP Contrast ±100")

- [ ] **Validate Against Round-Trip Test Results:**
  - [ ] Check if round-trip test results exist (from Epic 1 or Epic 6)
  - [ ] If available, cite test results in matrix ("Validated against 1,501 sample files")
  - [ ] If not available, note as "Future validation" or use Epic 1 retrospective data

**Validation:**
- Matrix accuracy matches Epic 1 implementation
- Mapping strategies consistent with code
- Test results cited (if available)

---

### Task 6: Create or Update Landing Page CSS for Matrix Link (Task 2)

- [ ] **Add Button Style for Matrix Link (if needed):**
  ```css
  /* Button Styles */
  .btn-secondary {
      display: inline-block;
      padding: 12px 24px;
      margin-top: 15px;
      background-color: #764ba2;
      color: white;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
      transition: background-color 0.3s;
  }

  .btn-secondary:hover {
      background-color: #667eea;
  }

  /* Format Table Styles */
  #formats table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 20px;
  }

  #formats th,
  #formats td {
      padding: 12px;
      border: 1px solid #ddd;
      text-align: left;
  }

  #formats th {
      background-color: #667eea;
      color: white;
  }

  #formats tr:nth-child(even) {
      background-color: #f7f7f7;
  }
  ```

- [ ] **Ensure Mobile-Responsive:**
  ```css
  @media (max-width: 768px) {
      #formats table {
          font-size: 0.9rem;
      }

      .btn-secondary {
          padding: 10px 20px;
          font-size: 0.9rem;
      }
  }
  ```

**Validation:**
- Button styles consistent with landing page design
- Table responsive on mobile
- Link visually distinct and clickable

---

### Task 7: Deploy and Verify Matrix Accessibility

- [ ] **Commit Changes to Git:**
  ```bash
  git add docs/format-compatibility-matrix.md web/index.html web/style.css README.md
  git commit -m "feat(epic-7): Add format compatibility matrix with detailed parameter mappings"
  git push origin main
  ```

- [ ] **Verify Cloudflare Pages Deployment:**
  - [ ] Push triggers automatic deployment
  - [ ] Wait for deployment to complete (~2-5 minutes)
  - [ ] Visit `https://recipe.pages.dev`

- [ ] **Test Matrix Accessibility:**
  - [ ] Click "View Detailed Compatibility Matrix" link on landing page
  - [ ] Verify opens `docs/format-compatibility-matrix.md`
  - [ ] Verify matrix renders correctly (tables, formatting)
  - [ ] Verify mobile-responsive (test on phone/tablet)

- [ ] **Test README Link:**
  - [ ] Visit GitHub repository
  - [ ] Click "View Complete Compatibility Matrix" link in README
  - [ ] Verify opens correct document

**Validation:**
- Deployment successful
- Matrix accessible from landing page
- README link working
- Mobile-responsive

---

## Dev Notes

### Learnings from Previous Story

**From Story 7-1-landing-page (Status: drafted)**

Story 7-1 established Recipe's landing page with privacy promise and 3-step usage guide. Story 7-2 builds on that foundation by **providing technical depth** for users evaluating Recipe's conversion capabilities.

**Key Insights from 7-1:**
- Landing page serves non-technical users - keep simple, link to details
- Privacy promise is core differentiator - compatibility matrix extends transparency
- Users need to understand limitations upfront - matrix sets expectations

**Integration:**
- Story 7-1: Landing page with format overview table (3 formats, bidirectional conversion)
- Story 7-2: Detailed compatibility matrix (30+ parameters, mapping quality, approximations)
- Together: Progressive disclosure (simple → detailed) based on user need

**No Technical Debt from 7-1:** Landing page complete and deployed. This story adds linked documentation (non-blocking, complements existing content).

**Reuse from 7-1:**
- Landing page structure (section, table, link pattern) - apply same CSS styles
- Format table design - maintain consistency (same column structure)
- Link button style - reuse `.btn-secondary` class

[Source: stories/7-1-landing-page.md]

---

### Architecture Alignment

**Follows Tech Spec Epic 7:**
- Format compatibility matrix satisfies FR-7.2 (all 4 ACs)
- Transparency about limitations aligns with privacy-first philosophy
- Detailed parameter documentation serves advanced users and contributors

**Epic 7 Documentation Philosophy:**
```
Recipe's Technical Transparency:

Landing Page (Simple)
    ↓
Format Compatibility Matrix (Detailed)
    ↓
Parameter Mapping Code (Epic 1 Implementation)
```

**Hub-and-Spoke Architecture Reference:**
This matrix documents the parameter mappings implemented in Epic 1's hub-and-spoke architecture:
- **Hub:** UniversalRecipe (superset of all format capabilities)
- **Spokes:** Format parsers/generators (NP3, XMP, lrtemplate)
- **Mapping:** This matrix documents the conversion logic between hub and spokes

**From Architecture Document (Section: UniversalRecipe Structure):**
- All parameters listed in matrix are fields in `UniversalRecipe` struct
- Direct 1:1 mappings preserve field values exactly
- Approximations use helper functions (e.g., `normalizeRange()`, `scaleParameter()`)
- Unmappable parameters stored in `Metadata` dictionary for round-trip preservation

---

### Dependencies

**Internal Dependencies:**
- `docs/tech-spec-epic-1.md` - Parameter mapping specifications (if exists)
- `docs/PRD.md#parameter-mapping` - PRD parameter requirements
- `docs/architecture.md#UniversalRecipe-Structure` - Data model details
- Story 1-8 (Parameter Mapping Rules) - Completed, code review approved
- Story 7-1 (Landing Page) - Links to this matrix

**External Dependencies:**
- None (static documentation only, no external APIs)

**No Blockers:** All Epic 1 stories complete (parameter mapping implemented and tested). This story documents existing functionality.

---

### Testing Strategy

**Manual Testing (Primary Method):**
- **Accuracy Review:** Cross-reference matrix with Epic 1 implementation (code review)
- **Link Validation:** Click all links (landing page, README)
- **Visual Review:** Check table formatting, symbol consistency
- **Mobile Responsive:** Test on phone/tablet (table readable, no horizontal scroll)

**Content Validation:**
- Compare matrix parameters with `internal/model/recipe.go` (UniversalRecipe struct)
- Verify approximation strategies match Story 1-8 implementation
- Confirm unmappable features list complete (no NP3 support for Grain, Vignette, etc.)

**Acceptance:**
- Matrix comprehensive (30+ parameters documented)
- Mapping quality accurate (✅, ⚠️, ❌ match implementation)
- Links working (landing page, README)
- Mobile-responsive

---

### Technical Debt / Future Enhancements

**Deferred to Post-MVP:**
- **Interactive Matrix:** Sortable/filterable table (JavaScript) for easier parameter search
- **Visual Examples:** Before/after images showing conversion results
- **Parameter Delta Visualization:** Show how parameters scale between formats (e.g., NP3 ±3 → XMP ±100)
- **Automated Matrix Generation:** Script to extract mappings from code comments and generate matrix automatically
- **Conversion Quality Scores:** Per-parameter accuracy scores from round-trip tests (e.g., "Exposure: 99.2% accuracy")

**Future Improvements:**
- Add Canon/Sony/Fujifilm format rows when support added (future epics)
- Link to specific Story 1-8 mapping formulas for each parameter
- Add "Common Presets" section showing real-world conversion examples
- Embed parameter range diagrams (e.g., visual scale showing -3 to +3 vs. -100 to +100)

---

### References

- [Source: docs/tech-spec-epic-7.md#FR-7.2] - Format compatibility matrix requirements (4 ACs)
- [Source: docs/PRD.md#FR-1.6] - Parameter mapping & approximation requirements
- [Source: docs/architecture.md#UniversalRecipe-Structure] - Data model and parameter definitions
- [Source: stories/1-8-parameter-mapping-rules.md] - Mapping formulas and conversion logic
- [Source: docs/epic-1-retrospective.md] - Conversion accuracy validation results

**External References:**
- Adobe XMP Specification: https://www.adobe.com/devnet/xmp.html
- Nikon Picture Control Utility: https://downloadcenter.nikonimglib.com/
- Color Delta E (Visual Similarity): https://en.wikipedia.org/wiki/Color_difference

---

### Known Issues / Blockers

**None** - This story has no technical blockers. All required information exists in:
- Epic 1 implementation (Story 1-8 parameter mapping complete)
- PRD (parameter requirements and mapping strategy)
- Architecture (UniversalRecipe data model)

**Content Decisions Needed:**
- **Option A or B for Landing Page Matrix:** Embed simplified table OR link only?
  - Recommendation: **Link only** (Option B) for cleaner landing page
- **Matrix Location:** `docs/format-compatibility-matrix.md` OR separate website page?
  - Recommendation: **Markdown file** (easier to maintain, version control)

**Assumptions:**
- Epic 1 Story 1-8 (Parameter Mapping Rules) complete with accurate mapping formulas
- Round-trip test results available (from Epic 1 or Epic 6) for accuracy validation
- Landing page deployed (Story 7-1 complete)

---

### Cross-Story Coordination

**Dependencies:**
- Story 7-1 (Landing Page) - Provides link location for this matrix
- Story 1-8 (Parameter Mapping Rules) - Provides mapping formulas to document
- Epic 1 (All Stories) - Provides implementation details for matrix accuracy

**Enables:**
- Story 7-3 (FAQ Documentation) - Can reference matrix for "Why doesn't X convert?" answers
- User adoption - Sets correct expectations, builds trust through transparency
- Community contributions - Provides documentation template for new format support

**Architectural Consistency:**
This story documents the conversion logic implemented in Epic 1:
- Epic 1: Implemented parameter mappings (code)
- Story 7-2: Documents parameter mappings (matrix)
- Result: Code and documentation aligned, users understand conversion behavior

---

### Project Structure Notes

**New Files Created:**
```
docs/
├── format-compatibility-matrix.md   # Complete parameter mapping matrix (NEW)
```

**Modified Files:**
```
web/
├── index.html                       # Add link to matrix in Supported Formats section (MODIFIED)
├── style.css                        # Add button styles for matrix link (MODIFIED - minor)

README.md                            # Add Format Compatibility section with matrix link (MODIFIED)
```

**No Conflicts:** This story adds new documentation and a link to existing landing page. No structural changes to existing Web UI.

**File Organization:**
- Detailed matrix in `docs/format-compatibility-matrix.md` (comprehensive reference)
- Landing page link in `web/index.html` (user-facing entry point)
- README link (developer-facing entry point)

---

## Dev Agent Record

### Context Reference

- `docs/stories/7-2-format-compatibility-matrix.context.xml` - Generated 2025-11-06

### Agent Model Used

<!-- To be filled by dev agent -->

### Debug Log References

<!-- Dev agent will add references to detailed debug logs if needed -->

### Completion Notes List

<!-- Dev agent will document:
- Format compatibility matrix creation (30+ parameters documented)
- Parameter categorization (Basic, Color, HSL, Detail, Tone Curve, Effects, etc.)
- Mapping quality indicators (✅, ⚠️, ❌) and legend
- Approximation strategy explanations (Highlights→Contrast, Shadows→Brightness, Vibrance→Saturation)
- Known limitations section (NP3 constraints, unmappable features)
- Conversion accuracy validation (round-trip testing, visual similarity)
- Landing page link integration (Supported Formats section)
- README.md Format Compatibility section addition
- Link validation (all links working, correct destinations)
- Mobile-responsive table formatting
- Cross-reference with Epic 1 implementation (Story 1-8 mapping formulas)
- Content accuracy review (parameters match UniversalRecipe struct)
-->

### File List

<!-- Dev agent will document files created/modified/deleted:
**NEW:**
- `docs/format-compatibility-matrix.md` - Comprehensive parameter mapping matrix with 30+ parameters, approximation strategies, known limitations

**MODIFIED:**
- `web/index.html` - Added link to matrix in Supported Formats section (#formats)
- `web/style.css` - Added `.btn-secondary` button style for matrix link (if not already present from Story 7-1)
- `README.md` - Added Format Compatibility section with matrix link

**DELETED:**
- (none)
-->

---

## Change Log

- **2025-11-06:** Story created from Epic 7 Tech Spec (Second story in Epic 7, documents parameter mappings from Epic 1 implementation)
