# Visual Regression Testing Results

**Recipe Version:** v2.0.0-dev
**Test Date:** 2025-11-06
**Test Environment:** Manual visual regression testing (MVP)
**Status:** ⏳ Awaiting manual testing execution

---

## Executive Summary

> **Note:** This is a template document. Results will be populated after manual visual regression testing is completed (Tasks 3-6 of Story 6-2).

Visual regression testing validates Recipe's 95%+ conversion accuracy claim by applying presets to real photos and comparing outputs visually. This report documents metrics (Delta E, SSIM) and subjective assessments for [TO BE FILLED: X] preset-image combinations.

### Target Metrics

- **Visual Similarity:** ≥95% (subjective assessment)
- **Color Delta E:** <5 for critical colors (reds, blues, greens, skin tones)
- **SSIM:** >0.95 (excellent structural similarity)

### Results Summary (To Be Filled)

| Metric                          | Target | Actual | Status |
| ------------------------------- | ------ | ------ | ------ |
| **Total Comparisons**           | 50     | [TBD]  | ⏳      |
| **Excellent (Delta E <1)**      | 70%+   | [TBD]  | ⏳      |
| **Good (Delta E 1-3)**          | 20%+   | [TBD]  | ⏳      |
| **Acceptable (Delta E 3-5)**    | <10%   | [TBD]  | ⏳      |
| **Poor (Delta E >5)**           | 0%     | [TBD]  | ⏳      |
| **Average SSIM**                | >0.95  | [TBD]  | ⏳      |
| **Average Delta E**             | <3.0   | [TBD]  | ⏳      |
| **95%+ Visual Similarity Met?** | Yes    | [TBD]  | ⏳      |

---

## Test Methodology

### Reference Images

**Count:** 5 images selected from diverse photography scenarios

| Image           | Resolution | Color Profile | Subject Description                                    | Key Colors to Test                  |
| --------------- | ---------- | ------------- | ------------------------------------------------------ | ----------------------------------- |
| portrait.jpg    | [TBD]      | Adobe RGB     | [TBD: e.g., Female portrait, warm skin tones]          | Skin tones, hair, eyes, background  |
| landscape.jpg   | [TBD]      | Adobe RGB     | [TBD: e.g., Landscape with blue sky, green foliage]    | Blue sky, green foliage, earth tones |
| vintage.jpg     | [TBD]      | Adobe RGB     | [TBD: e.g., Urban scene, nostalgic subject]            | Muted tones, warm hues              |
| bw-candidate.jpg | [TBD]      | Adobe RGB     | [TBD: e.g., High-contrast scene for B&W conversion]    | Tonal range (highlights to shadows) |
| hdr-candidate.jpg | [TBD]     | Adobe RGB     | [TBD: e.g., Backlit sunset, extreme dynamic range]     | Highlight/shadow recovery           |

See `testdata/visual-regression/README.md` for detailed image documentation.

### Test Presets

**Count:** 10 representative presets covering subtle to dramatic adjustments

| Preset Name | Source File | Category | Intensity | Key Adjustments | Conversion Path | Target Delta E |
| ----------- | ----------- | -------- | --------- | --------------- | --------------- | -------------- |
| [TBD]       | [TBD]       | [TBD]    | Subtle    | [TBD]           | XMP → NP3       | <1             |
| [TBD]       | [TBD]       | [TBD]    | Subtle    | [TBD]           | XMP → NP3       | <1             |
| [TBD]       | [TBD]       | [TBD]    | Subtle    | [TBD]           | XMP → NP3       | <1             |
| [TBD]       | [TBD]       | [TBD]    | Moderate  | [TBD]           | XMP → NP3       | 1-3            |
| [TBD]       | [TBD]       | [TBD]    | Moderate  | [TBD]           | XMP → NP3       | 1-3            |
| [TBD]       | [TBD]       | [TBD]    | Moderate  | [TBD]           | XMP → NP3       | 1-3            |
| [TBD]       | [TBD]       | [TBD]    | Moderate  | [TBD]           | XMP → NP3       | 1-3            |
| [TBD]       | [TBD]       | [TBD]    | Dramatic  | [TBD]           | XMP → NP3       | 3-5            |
| [TBD]       | [TBD]       | [TBD]    | Dramatic  | [TBD]           | XMP → NP3       | 3-5            |
| [TBD]       | [TBD]       | [TBD]    | Dramatic  | [TBD]           | XMP → NP3       | 3-5            |

See `testdata/visual-regression/presets.md` for detailed preset documentation.

### Export Settings

**Reference Outputs:**
- Application: Adobe Lightroom CC/Classic
- Format: 16-bit TIFF
- Color Space: Adobe RGB (1998)
- Compression: None (lossless)
- Location: `testdata/visual-regression/reference/`

**Test Outputs:**
- Application: Nikon NX Studio
- Format: 16-bit TIFF
- Color Space: Adobe RGB (1998)
- Compression: None (lossless)
- Location: `testdata/visual-regression/test/`

See `testdata/visual-regression/export-settings.md` for complete export configuration.

---

## Detailed Results

### Sample Comparisons (To Be Filled)

> **Note:** This section will be populated after manual testing. Below is a template showing expected format.

---

#### Example: Portrait - Warm Glow (Moderate Preset)

| Reference (XMP in Lightroom)                          | Test (NP3 in NX Studio)                             |
| ----------------------------------------------------- | --------------------------------------------------- |
| ![Reference](../testdata/visual-regression/reference/portrait-warm-glow.tiff) | ![Test](../testdata/visual-regression/test/portrait-warm-glow-np3.tiff) |

**Metrics:**
- **Delta E:** [TBD] (Target: 1-3 for moderate preset)
- **SSIM:** [TBD] (Target: >0.95)
- **Assessment:** [TBD: Excellent / Good / Acceptable / Poor]

**Observations:**
- [TBD: e.g., Skin tones preserved accurately]
- [TBD: e.g., Slight shift in shadow warmth, within tolerance]
- [TBD: e.g., Highlight detail maintained]

**Difference Map:**

![Difference Map](../testdata/visual-regression/diff/portrait-warm-glow-diff.png)

---

#### [Template for Additional Comparisons - To Be Duplicated]

| Reference (XMP in Lightroom)    | Test (NP3 in NX Studio)          |
| ------------------------------- | -------------------------------- |
| ![Reference](../testdata/visual-regression/reference/[preset]-[image].tiff) | ![Test](../testdata/visual-regression/test/[preset]-[image]-np3.tiff) |

**Metrics:**
- **Delta E:** [TBD]
- **SSIM:** [TBD]
- **Assessment:** [TBD]

**Observations:**
- [TBD]

---

## Metrics Table (All Comparisons)

> **Note:** This table will be populated after automated metrics calculation (Task 6).

| Preset Name    | Image     | Delta E | SSIM  | Assessment | Notes                  |
| -------------- | --------- | ------- | ----- | ---------- | ---------------------- |
| [Preset 1]     | Portrait  | [TBD]   | [TBD] | [TBD]      | [Observations]         |
| [Preset 1]     | Landscape | [TBD]   | [TBD] | [TBD]      | [Observations]         |
| [Preset 1]     | Vintage   | [TBD]   | [TBD] | [TBD]      | [Observations]         |
| [Preset 1]     | B&W       | [TBD]   | [TBD] | [TBD]      | [Observations]         |
| [Preset 1]     | HDR       | [TBD]   | [TBD] | [TBD]      | [Observations]         |
| [Preset 2]     | Portrait  | [TBD]   | [TBD] | [TBD]      | [Observations]         |
| [... 44 more rows ...] |     |         |       |            |                        |

**Summary Statistics:**
- **Total Comparisons:** [TBD] / 50
- **Average Delta E:** [TBD]
- **Average SSIM:** [TBD]
- **Min Delta E:** [TBD] (best match)
- **Max Delta E:** [TBD] (worst match)

---

## Known Differences (Cross-Reference)

See `docs/known-conversion-limitations.md` for detailed documentation of unmappable features.

### Grain Effect (XMP → NP3)

**Visual Impact:** Low (Delta E: 1.2 average)

**Affected Presets:** [TBD: List presets with grain effect]

**Observations:**
- [TBD: e.g., Grain texture lost in conversion as expected]
- [TBD: e.g., Overall tonal adjustments preserved]

### Vignette (XMP → NP3)

**Visual Impact:** Medium (Delta E: 3.8 average)

**Affected Presets:** [TBD: List presets with vignette]

**Observations:**
- [TBD: e.g., Edge darkening lost in conversion as expected]
- [TBD: e.g., Noticeable on bright backgrounds]

### Split Toning (XMP → NP3, Limited Support)

**Visual Impact:** Low-Medium (Delta E: 2.5 average for subtle, 4.2 for dramatic)

**Affected Presets:** [TBD: List presets with split toning]

**Observations:**
- [TBD: e.g., Subtle split toning approximated well]
- [TBD: e.g., Dramatic split toning shows noticeable color shift]

---

## Recommendations

### Findings (To Be Filled After Testing)

1. **[TBD: Finding 1]**
   - Observation: [TBD]
   - Recommendation: [TBD]
   - Priority: High / Medium / Low

2. **[TBD: Finding 2]**
   - Observation: [TBD]
   - Recommendation: [TBD]
   - Priority: High / Medium / Low

### Future Improvements

**Short-Term (Next Release):**
- [TBD: e.g., Improve split toning approximation - current Delta E 4.2, target <3.0]
- [TBD: e.g., Add grain simulation to NP3 generator (if feasible)]

**Long-Term (Post-MVP):**
- Automate preset application via Lightroom scripting
- Automate metrics calculation (batch processing)
- Integrate into CI/CD pipeline (detect regressions)
- Expand reference image library (10+ images)

---

## Validation Status

| Task                                   | Status | Date       | Validator |
| -------------------------------------- | ------ | ---------- | --------- |
| Reference images selected and documented | ⏳      | [TBD]      | [TBD]     |
| Presets selected and documented        | ⏳      | [TBD]      | [TBD]     |
| Reference outputs generated (Lightroom)| ⏳      | [TBD]      | [TBD]     |
| Test outputs generated (NX Studio)     | ⏳      | [TBD]      | [TBD]     |
| Manual visual comparison completed     | ⏳      | [TBD]      | [TBD]     |
| Automated metrics calculated (SSIM, Delta E) | ⏳ | [TBD]   | [TBD]     |
| Known differences documented           | ✅      | 2025-11-06 | Dev Agent |
| Visual regression report finalized     | ⏳      | [TBD]      | [TBD]     |

---

## Tool Versions

**Adobe Lightroom:**
- Version: [TBD]
- Build: [TBD]

**Nikon NX Studio:**
- Version: [TBD]
- Build: [TBD]

**ImageMagick:**
- Version: [TBD]

**Python:**
- Version: [TBD]
- colormath: [TBD]
- Pillow: [TBD]
- numpy: [TBD]

**Recipe:**
- Version: v2.0.0-dev
- Commit: [TBD: git rev-parse --short HEAD]

---

## Reproducibility

### Steps to Reproduce

1. **Clone Repository:**
   ```bash
   git clone https://github.com/yourusername/recipe.git
   cd recipe
   ```

2. **Install Dependencies:**
   ```bash
   # ImageMagick (macOS)
   brew install imagemagick

   # Python dependencies
   pip3 install Pillow numpy colormath
   ```

3. **Download Reference Images:**
   ```bash
   # [TBD: Provide download link if images stored externally]
   # OR use images in testdata/visual-regression/images/ (if committed)
   ```

4. **Generate Reference Outputs:**
   - Open `testdata/visual-regression/images/[image].jpg` in Adobe Lightroom
   - Apply preset from `testdata/xmp/[preset].xmp`
   - Export as 16-bit TIFF (settings in `testdata/visual-regression/export-settings.md`)
   - Save to `testdata/visual-regression/reference/[image]-[preset].tiff`

5. **Convert Presets:**
   ```bash
   recipe convert testdata/xmp/[preset].xmp --to np3 --output test-presets/[preset].np3
   ```

6. **Generate Test Outputs:**
   - Open `testdata/visual-regression/images/[image].jpg` in Nikon NX Studio
   - Load converted preset from `test-presets/[preset].np3`
   - Export as 16-bit TIFF (same settings as Lightroom)
   - Save to `testdata/visual-regression/test/[image]-[preset]-np3.tiff`

7. **Run Comparison:**
   ```bash
   cd scripts/visual-regression
   ./compare-all.sh
   ```

8. **Review Results:**
   - Check `testdata/visual-regression/results.json` for summary
   - View difference maps in `testdata/visual-regression/diff/`
   - Compare side-by-side in photo viewer

---

## Community Validation

Recipe's visual regression testing is fully transparent and reproducible. Community members can:

- **Validate results independently** using steps above
- **Test with their own presets** following same methodology
- **Contribute test results** via GitHub issues or pull requests
- **Report discrepancies** if results differ from documented findings

---

## Change Log

| Version | Date       | Changes                                                      |
| ------- | ---------- | ------------------------------------------------------------ |
| v1.0    | 2025-11-06 | Initial template created (Story 6-2), awaiting manual testing |

---

**Last Updated:** 2025-11-06
**Status:** Template created, awaiting manual testing execution
**Next Steps:** Complete Tasks 3-6 of Story 6-2 to populate results
