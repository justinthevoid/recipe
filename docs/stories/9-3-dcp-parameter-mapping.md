# Story 9.3: DCP Parameter Mapping Documentation

Status: ready-for-dev

## Story

As a **developer**,
I want **comprehensive documentation of how UniversalRecipe parameters map to DNG Camera Profile (DCP) tone curves and color matrices**,
so that **I understand the conversion formulas, supported/unsupported features, and can verify correct DCP generation from any source format (NP3, XMP, lrtemplate, .costyle)**.

## Acceptance Criteria

**AC-1: Document UniversalRecipe → DCP Parameter Mappings**
- ✅ Add "DCP (DNG Camera Profile)" section to `docs/parameter-mapping.md`
- ✅ Document all supported UniversalRecipe parameters that map to DCP:
  - Exposure → Tone curve vertical shift (midpoint 128 → 128+shift)
  - Contrast → Tone curve slope (steepen/flatten)
  - Highlights → Tone curve top-end adjustment (192-255 range)
  - Shadows → Tone curve bottom-end adjustment (0-64 range)
- ✅ Specify exact conversion formulas with code examples
- ✅ Include visual diagrams showing linear → adjusted curve transformations
- ✅ Document precision considerations (float → int curve points, rounding)

**AC-2: Identify Unsupported DCP Features**
- ✅ List DCP features NOT implemented in Recipe v0.1.0:
  - Dual illuminant profiles (single illuminant only)
  - Full camera calibration (ForwardMatrix, CalibrationIlluminant)
  - HSV (Hue/Saturation/Value) lookup tables
  - Chromatic aberration correction
  - Lens correction profiles
  - Advanced color grading (3D LUTs)
- ✅ Explain why each feature is unsupported (identity matrices used instead)
- ✅ Note Recipe uses identity color matrices (no color calibration)
- ✅ Document fallback behavior (unsupported parameters ignored)

**AC-3: Define Tone Curve Conversion Formulas**
- ✅ Document 5-point tone curve generation algorithm:
  ```
  Linear base: (0,0) (64,64) (128,128) (192,192) (255,255)

  Step 1 - Apply Exposure:
    midpoint.Output = clamp(128 + exposure*64, 0, 255)

  Step 2 - Apply Contrast:
    for each point:
      deviation = point.Input - 128
      point.Output = clamp(128 + deviation*(1+contrast), 0, 255)

  Step 3 - Apply Highlights:
    highlightsShift = highlights * 32
    points[3].Output = clamp(points[3].Output + highlightsShift, points[2].Output, 255)
    points[4].Output = clamp(points[4].Output + highlightsShift, points[3].Output, 255)

  Step 4 - Apply Shadows:
    shadowsShift = shadows * 32
    points[0].Output = clamp(points[0].Output + shadowsShift, 0, points[1].Output)
    points[1].Output = clamp(points[1].Output + shadowsShift, points[0].Output, points[2].Output)
  ```
- ✅ Include worked examples (exposure +0.5, contrast +0.3, etc.)
- ✅ Document monotonic curve enforcement (output[i] >= output[i-1])
- ✅ Explain clamping to 0-255 range

**AC-4: Document Color Matrix Handling**
- ✅ Explain identity matrix usage (3x3 diagonal matrix):
  ```
  ColorMatrix1 = ColorMatrix2 = [
      1.0  0.0  0.0
      0.0  1.0  0.0
      0.0  0.0  1.0
  ]
  ```
- ✅ Clarify this means no color space transformation (camera RGB → XYZ passthrough)
- ✅ Explain ProfileCalibrationSignature = "com.adobe" (standard signature)
- ✅ Note Recipe focuses on tone/exposure only (not color calibration)
- ✅ Document when full calibration would be needed (out of scope for v0.1.0)

**AC-5: Provide Cross-Format Conversion Examples**
- ✅ Include minimum 3 example conversions in documentation:
  - **Example 1**: NP3 (Nikon) → DCP
    - NP3 exposure +0.5 → DCP midpoint 128 → 160
    - NP3 contrast +0.3 → DCP slope factor 1.3
  - **Example 2**: XMP (Adobe) → DCP
    - XMP Exposure2012 +0.5 → DCP tone curve
    - XMP Contrast2012 +30 → DCP slope adjustment
  - **Example 3**: lrtemplate (Lightroom) → DCP
    - lrtemplate crs:Exposure2012 → DCP tone curve
    - lrtemplate preserves visual appearance through DCP
- ✅ Show complete UniversalRecipe intermediate step
- ✅ Include visual before/after curve diagrams
- ✅ Document expected precision (±1 curve point due to rounding)

**AC-6: Test Mappings with Real Adobe DCP Samples**
- ✅ Load minimum 3 real DCP files from Adobe (Camera Raw, Lightroom)
- ✅ Parse tone curves from real samples
- ✅ Reverse-engineer curve → UniversalRecipe parameters
- ✅ Generate new DCP from derived UniversalRecipe
- ✅ Compare original vs. generated curves (visual + numerical)
- ✅ Document findings in parameter-mapping.md:
  - Accuracy of reverse-engineered parameters
  - Curve point differences (absolute max delta)
  - Visual similarity assessment
- ✅ Store real DCP samples in `testdata/dcp/adobe-samples/` with README

**AC-7: Create Reference Documentation**
- ✅ Add comprehensive "DCP Parameter Mapping" section to parameter-mapping.md:
  - Supported parameters table
  - Unsupported features list
  - Tone curve formulas (code + math)
  - Identity matrix explanation
  - Cross-format examples
  - Real sample analysis results
- ✅ Include references to Adobe DNG Specification 1.6
- ✅ Add glossary terms (IFD, Tone Curve, Color Matrix, Illuminant, HSV table)
- ✅ Link to implementation files (dcp/generate.go, dcp/parse.go)
- ✅ Update docs/index.md with DCP section link

## Tasks / Subtasks

### Task 1: Document DCP Parameter Mappings (AC-1)
- [ ] Add "DCP (DNG Camera Profile)" section to `docs/parameter-mapping.md`
- [ ] Create "Supported Parameters" table:
  | UniversalRecipe Parameter | DCP Representation          | Conversion Formula             |
  | ------------------------- | --------------------------- | ------------------------------ |
  | Exposure                  | Tone curve midpoint shift   | `midpoint = 128 + exposure*64` |
  | Contrast                  | Tone curve slope factor     | `slope = 1.0 + contrast`       |
  | Highlights                | Top-end curve adjustment    | `shift = highlights * 32`      |
  | Shadows                   | Bottom-end curve adjustment | `shift = shadows * 32`         |
- [ ] Document 5-point tone curve structure (0, 64, 128, 192, 255)
- [ ] Explain precision: float parameters → integer curve points (0-255)
- [ ] Note clamping/rounding behavior

### Task 2: List Unsupported DCP Features (AC-2)
- [ ] Create "Unsupported Features" subsection in parameter-mapping.md
- [ ] List features NOT implemented:
  - Dual illuminant profiles (D65/A or D50/A)
  - ForwardMatrix (camera RGB → XYZ transformation)
  - CalibrationIlluminant1/2 (white balance calibration)
  - HSV lookup tables (hue/saturation/value adjustments)
  - ChromaticAberration correction
  - LensCorrection profiles
  - 3D LUTs for color grading
- [ ] Explain Recipe uses identity matrices instead (no color transformation)
- [ ] Document fallback: Unsupported parameters ignored during conversion
- [ ] Note future work: Full calibration could be added in v0.2.0+

### Task 3: Define Tone Curve Formulas (AC-3)
- [ ] Create "Tone Curve Generation Algorithm" subsection
- [ ] Document step-by-step formula:
  ```
  Input:  UniversalRecipe (exposure, contrast, highlights, shadows)
  Output: 5-point tone curve [(0,Y0), (64,Y1), (128,Y2), (192,Y3), (255,Y4)]

  Step 1: Initialize linear curve
    points = [(0,0), (64,64), (128,128), (192,192), (255,255)]

  Step 2: Apply exposure (vertical shift of midpoint)
    exposureShift = int(exposure * 64.0)
    points[2].Output = clamp(128 + exposureShift, 0, 255)

  Step 3: Apply contrast (steepen/flatten slope)
    contrastFactor = 1.0 + contrast
    for each point:
      deviation = point.Input - 128
      point.Output = clamp(128 + int(deviation * contrastFactor), 0, 255)

  Step 4: Apply highlights (adjust top-end)
    highlightsShift = int(highlights * 32.0)
    points[3].Output = clamp(points[3].Output + highlightsShift, points[2].Output, 255)
    points[4].Output = clamp(points[4].Output + highlightsShift, points[3].Output, 255)

  Step 5: Apply shadows (adjust bottom-end)
    shadowsShift = int(shadows * 32.0)
    points[0].Output = clamp(points[0].Output + shadowsShift, 0, points[1].Output)
    points[1].Output = clamp(points[1].Output + shadowsShift, points[0].Output, points[2].Output)

  Step 6: Enforce monotonic curve
    for i = 1 to 4:
      if points[i].Output < points[i-1].Output:
        points[i].Output = points[i-1].Output
  ```
- [ ] Include worked example:
  - Input: exposure = +0.5, contrast = +0.3, highlights = -0.2, shadows = +0.1
  - Step-by-step calculation showing intermediate values
  - Final curve: [(0,3), (64,67), (128,160), (192,186), (255,249)]
- [ ] Add visual ASCII diagram showing curve transformation

### Task 4: Document Color Matrix Handling (AC-4)
- [ ] Create "Color Matrix Handling" subsection
- [ ] Explain identity matrix (3x3 diagonal):
  ```
  ColorMatrix1 = ColorMatrix2 = [
      1.0000  0.0000  0.0000
      0.0000  1.0000  0.0000
      0.0000  0.0000  1.0000
  ]
  ```
- [ ] Clarify meaning: No color space transformation (camera RGB passthrough)
- [ ] Explain ProfileCalibrationSignature = "com.adobe" (standard non-calibrated signature)
- [ ] Note Recipe v0.1.0 focuses on tone/exposure only (not color calibration)
- [ ] Document when full calibration needed:
  - Multi-illuminant support (D65 + tungsten)
  - Camera-specific color accuracy
  - Professional color grading workflows
  - Out of scope for initial release

### Task 5: Create Cross-Format Examples (AC-5)
- [ ] Add "Conversion Examples" subsection with 3 detailed examples
- [ ] **Example 1: NP3 → DCP**
  ```
  NP3 Input:
    Exposure: +0.5 (NP3 range: -1.0 to +1.0)
    Contrast: +0.3

  UniversalRecipe (intermediate):
    Exposure: 0.5
    Contrast: 0.3

  DCP Output (tone curve):
    Midpoint: 128 → 160 (128 + 0.5*64 = 160)
    Slope: 1.3x (1.0 + 0.3 = 1.3)
    Curve: [(0,0), (64,45), (128,160), (192,211), (255,255)]
  ```
- [ ] **Example 2: XMP → DCP**
  ```
  XMP Input (Adobe Camera Raw):
    crs:Exposure2012: +0.50 (stops)
    crs:Contrast2012: +30 (range -100 to +100)

  UniversalRecipe:
    Exposure: 0.5
    Contrast: 0.3 (30/100 = 0.3)

  DCP Output:
    Same as Example 1
  ```
- [ ] **Example 3: lrtemplate → DCP**
  ```
  lrtemplate Input (Lightroom template):
    s.Exposure2012 = 0.5
    s.Contrast2012 = 30

  Conversion path:
    lrtemplate → Parse → UniversalRecipe → Generate → DCP

  Expected precision: ±1 curve point due to float→int rounding
  ```
- [ ] Include visual before/after curve diagrams (ASCII or reference to images)

### Task 6: Test with Real Adobe DCP Samples (AC-6)
- [ ] Acquire 3 real DCP samples:
  - Download from Adobe Camera Raw (bundled profiles)
  - Download from Adobe Lightroom (Camera Profiles folder)
  - Use DCP files from Adobe DNG SDK samples
- [ ] Store in `testdata/dcp/adobe-samples/`:
  - `adobe-standard.dcp` (standard camera profile)
  - `adobe-landscape.dcp` (landscape preset)
  - `adobe-portrait.dcp` (portrait preset)
- [ ] Parse each sample using Story 9-1 Parse() function
- [ ] Analyze tone curves:
  - Extract 5-point curve values
  - Reverse-engineer exposure/contrast/highlights/shadows parameters
  - Document curve shapes (linear, S-curve, etc.)
- [ ] Generate new DCP from reverse-engineered UniversalRecipe
- [ ] Compare original vs. generated:
  - Curve point differences (calculate max absolute delta)
  - Visual similarity (plot curves in spreadsheet or tool)
  - XML structure comparison
- [ ] Document findings in parameter-mapping.md:
  ```
  ## Real Adobe DCP Analysis

  ### Sample 1: adobe-standard.dcp
  - Original curve: [(0,0), (64,64), (128,128), (192,192), (255,255)]
  - Derived UniversalRecipe: exposure=0.0, contrast=0.0 (neutral/linear)
  - Regenerated curve: [(0,0), (64,64), (128,128), (192,192), (255,255)]
  - Max delta: 0 points (exact match)

  ### Sample 2: adobe-landscape.dcp
  - Original curve: [(0,5), (64,70), (128,140), (192,200), (255,250)]
  - Derived UniversalRecipe: exposure=+0.38, contrast=+0.28, highlights=-0.16
  - Regenerated curve: [(0,5), (64,70), (128,140), (192,199), (255,249)]
  - Max delta: 1 point at (192) and (255) (acceptable precision)

  ### Sample 3: adobe-portrait.dcp
  - Similar analysis...
  ```
- [ ] Add README to testdata/dcp/adobe-samples/ explaining sample sources

### Task 7: Create Reference Documentation (AC-7)
- [ ] Consolidate all sections into parameter-mapping.md
- [ ] Add table of contents for DCP section:
  ```
  ### DCP (DNG Camera Profile)

  #### Overview
  #### Supported Parameters
  #### Unsupported Features
  #### Tone Curve Generation Algorithm
  #### Color Matrix Handling
  #### Conversion Examples
  #### Real Adobe DCP Analysis
  #### References
  ```
- [ ] Add glossary subsection:
  - **IFD (Image File Directory)**: TIFF metadata structure containing tags
  - **Tone Curve**: 5-point curve mapping input (0-255) to output (0-255)
  - **Color Matrix**: 3x3 matrix for RGB color space transformation
  - **Illuminant**: Light source for color calibration (D65, A, etc.)
  - **HSV Table**: Hue/Saturation/Value lookup table for color grading
- [ ] Add references subsection:
  - [Adobe DNG Specification 1.6](https://helpx.adobe.com/camera-raw/digital-negative.html)
  - [Source: internal/formats/dcp/generate.go] - Generation implementation
  - [Source: internal/formats/dcp/parse.go] - Parsing implementation
  - [Source: docs/tech-spec-epic-9.md] - DCP epic technical specification
- [ ] Link from docs/index.md:
  ```markdown
  ## Parameter Mapping Documentation

  - [DCP (DNG Camera Profile) Mapping](parameter-mapping.md#dcp-dng-camera-profile)
  ```
- [ ] Verify all cross-references work (internal links)

### Task 8: Manual Validation
- [ ] Review documentation for technical accuracy:
  - Verify formulas match implementation in dcp/generate.go
  - Check examples for mathematical correctness
  - Validate tone curve calculations manually (spreadsheet)
- [ ] Test documentation usability:
  - Can a new developer understand DCP conversion from docs alone?
  - Are examples clear and complete?
  - Is glossary helpful for unfamiliar terms?
- [ ] Peer review (optional):
  - Ask Justin to review for clarity
  - Check if format comparison makes sense
- [ ] Fix any issues found during validation

## Dev Notes

### Learnings from Previous Story

**From Story 9-2-dcp-generator (Status: drafted)**

Previous story not yet implemented - no file/pattern reuse available. This story creates the documentation foundation for DCP conversion.

[Source: docs/stories/9-2-dcp-generator.md]

### Architecture Alignment

**Tech Spec Epic 9 Alignment:**

Story 9-3 implements **AC-3 (DCP Parameter Mapping)** from tech-spec-epic-9.md.

**Documentation Flow:**
```
Real Adobe DCPs → Parse → Analyze → Reverse-engineer → Document formulas
         ↓
   parameter-mapping.md (DCP section)
         ↓
   Reference for developers, testers, users
```

**Supported Parameter Mapping:**
```
UniversalRecipe.Exposure   → DCP tone curve midpoint shift (128 → 128+shift)
UniversalRecipe.Contrast   → DCP tone curve slope (1.0 → 1.0+contrast)
UniversalRecipe.Highlights → DCP top-end adjustment (192-255 range)
UniversalRecipe.Shadows    → DCP bottom-end adjustment (0-64 range)
```

**Unsupported Features (Identity Matrix Approach):**
- No dual illuminant (D65/A)
- No ForwardMatrix (camera RGB → XYZ)
- No CalibrationIlluminant
- No HSV lookup tables
- No chromatic aberration correction

[Source: docs/tech-spec-epic-9.md#Detailed-Design]

### Documentation Strategy

**Parameter Mapping Document Structure:**

The `docs/parameter-mapping.md` file already contains mappings for:
- NP3 (Nikon Picture Control) - Story 1-2, 1-3
- XMP (Adobe Camera Raw) - Story 1-4, 1-5
- lrtemplate (Lightroom) - Story 1-6, 1-7
- .costyle (Capture One) - Story 8-1, 8-2

**Adding DCP Section:**

Follow established pattern from other formats:
1. Overview paragraph
2. Supported parameters table
3. Conversion formulas (code + math)
4. Unsupported features list
5. Cross-format examples
6. Real sample analysis (unique to DCP)

**Tone Curve Visualization:**

Include ASCII diagrams showing curve transformations:
```
Linear (neutral):
255 |                    •
    |                ·
128 |            •
    |        ·
  0 | ·——————————————————————
    0   64  128  192  255

Exposure +0.5:
255 |                    •
    |            ·
128 |        •          (midpoint shifts up)
    |    ·
  0 | ·——————————————————————
    0   64  128  192  255
```

[Source: docs/tech-spec-epic-9.md#Documentation-Requirements]

### Real DCP Sample Analysis Strategy

**Sample Acquisition:**

1. **Adobe Camera Raw** (Windows):
   - Location: `C:\Program Files\Adobe\Adobe Camera Raw\Camera Profiles\`
   - Files: `Adobe Standard.dcp`, `Camera Landscape.dcp`, etc.

2. **Lightroom Classic** (Windows):
   - Location: `C:\Program Files\Adobe\Adobe Lightroom Classic\Resources\Camera Profiles\`
   - Files: Multiple manufacturer-specific DCPs

3. **Adobe DNG SDK**:
   - Download: https://helpx.adobe.com/camera-raw/digital-negative.html
   - Sample DCPs in `dng_sdk/samples/` folder

**Analysis Process:**

1. Parse DCP using `dcp.Parse()` (Story 9-1)
2. Extract `ToneCurve.Points` (5 points)
3. Reverse-engineer parameters:
   - Exposure = (midpoint - 128) / 64.0
   - Contrast = (slope - 1.0)
   - Highlights/Shadows = calculated from curve ends
4. Generate new DCP using `dcp.Generate()` (Story 9-2)
5. Compare curves point-by-point
6. Calculate max delta (absolute difference)
7. Visual check (plot both curves)

**Expected Precision:**

- ±1 curve point due to float→int rounding
- ±2 points acceptable for complex curves
- Exact match for linear/neutral curves

[Source: docs/tech-spec-epic-9.md#Acceptance-Criteria]

### Cross-Format Conversion Paths

**Supported Conversion Paths through DCP:**

```
NP3 ──────────┐
              ├─→ UniversalRecipe ─→ DCP ─→ Lightroom/Camera Raw
XMP ──────────┤
              ├─→ (hub-and-spoke)
lrtemplate ───┤
              └─→ (tone adjustments only, no color calibration)
.costyle ─────┘
```

**Example Workflows:**

1. **Nikon → Adobe**:
   - NP3 file → Parse → UniversalRecipe → Generate → DCP
   - Load DCP in Lightroom → Apply to Canon RAW photos

2. **Capture One → Adobe**:
   - .costyle file → Parse → UniversalRecipe → Generate → DCP
   - Portable preset across ecosystems

3. **Lightroom → Camera Profile**:
   - lrtemplate → Parse → UniversalRecipe → Generate → DCP
   - Bake preset into camera profile (applied at import)

[Source: docs/tech-spec-epic-9.md#Use-Cases]

### Project Structure Notes

**New Files Created (Story 9-3):**
```
testdata/dcp/adobe-samples/
├── adobe-standard.dcp       # Real Adobe sample (NEW)
├── adobe-landscape.dcp      # Real Adobe sample (NEW)
├── adobe-portrait.dcp       # Real Adobe sample (NEW)
└── README.md                # Sample sources and licensing (NEW)
```

**Modified Files:**
- `docs/parameter-mapping.md` - Add DCP section (800+ lines)
- `docs/index.md` - Add link to DCP mapping documentation

**Files from Previous Stories (Referenced):**
- `internal/formats/dcp/generate.go` - Implementation reference (Story 9-2)
- `internal/formats/dcp/parse.go` - Implementation reference (Story 9-1)
- `testdata/dcp/` - Test samples (Story 9-1)

[Source: docs/tech-spec-epic-9.md#Components]

### Testing Strategy

**Documentation Testing:**

No unit tests required (documentation story), but perform:

1. **Technical Accuracy Verification:**
   - Manually validate formulas against dcp/generate.go implementation
   - Test examples with calculator/spreadsheet
   - Verify tone curve calculations step-by-step

2. **Real Sample Analysis:**
   - Load 3 Adobe DCP samples
   - Parse with dcp.Parse()
   - Reverse-engineer parameters
   - Regenerate with dcp.Generate()
   - Compare curves numerically and visually
   - Document max delta and precision

3. **Usability Testing:**
   - Can a developer understand DCP conversion from docs alone?
   - Are examples complete and runnable?
   - Is glossary helpful?

**Validation Checklist:**
- [ ] All formulas match implementation
- [ ] Examples are mathematically correct
- [ ] Real sample analysis complete (3 DCPs)
- [ ] Cross-references work (internal links)
- [ ] Glossary covers all technical terms

[Source: docs/tech-spec-epic-9.md#Test-Strategy-Summary]

### Known Risks

**RISK-19: Real Adobe DCPs unavailable**
- **Impact**: Cannot test with real samples for analysis
- **Mitigation**: Download from Adobe Camera Raw, Lightroom, or DNG SDK
- **Fallback**: Generate synthetic DCPs for analysis (less ideal)

**RISK-20: Reverse-engineering inaccurate**
- **Impact**: Derived parameters don't match original intent
- **Mitigation**: Visual validation, compare regenerated curves
- **Acceptable**: ±2 curve points due to discretization

**RISK-21: Documentation becomes outdated**
- **Impact**: Docs don't reflect implementation changes
- **Mitigation**: Link to source files, version documentation clearly
- **Process**: Update docs when generate.go changes (part of code review)

[Source: docs/tech-spec-epic-9.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-9.md#Acceptance-Criteria] - AC-3: DCP Parameter Mapping
- [Source: docs/tech-spec-epic-9.md#Data-Models-and-Contracts] - Tone curve formulas
- [Source: docs/parameter-mapping.md] - Existing format mappings (NP3, XMP, lrtemplate, .costyle)
- [Source: Adobe DNG Specification 1.6](https://helpx.adobe.com/camera-raw/digital-negative.html) - DCP format reference
- [Source: internal/formats/dcp/generate.go] - Generation implementation (Story 9-2)
- [Source: internal/formats/dcp/parse.go] - Parsing implementation (Story 9-1)

## Dev Agent Record

### Context Reference

- docs/stories/9-3-dcp-parameter-mapping.context.xml (generated 2025-11-09)

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
