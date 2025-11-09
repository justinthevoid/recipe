# Story 11.5: Clear Disclaimer About CSS Filter Approximation

Status: ready-for-dev

## Story

As a **photographer evaluating preset previews**,
I want **clear communication about the limitations of CSS filter-based previews**,
so that **I understand the preview is approximate and set appropriate expectations for the actual conversion**.

## Acceptance Criteria

**AC-1: Prominent Disclaimer Label in Modal**
- ✅ Disclaimer label visible in preview modal:
  - Text: "⚠️ Approximate preview using CSS filters"
  - Position: Top of modal (above reference image tabs) OR below slider
  - Style: Warning/info badge (yellow/blue background, icon, readable text)
  - Font size: 14px minimum (readable at all screen sizes)
- ✅ Label always visible:
  - Persistent (not hidden or dismissible)
  - Visible on modal open (no user action required)
  - Visible on all tabs (Portrait, Landscape, Product)
- ✅ Visual prominence:
  - Contrasts with background (WCAG AA 4.5:1 minimum)
  - Icon (⚠️ or ℹ️) draws attention
  - Stands out from other modal elements (not buried in UI)

**AC-2: Tooltip/Help Text Explains Limitations**
- ✅ Tooltip triggered by:
  - Hover: Mouse hover over disclaimer label
  - Tap: Tap/click disclaimer label on touch devices
  - Keyboard: Focus label + Enter/Space key
- ✅ Tooltip content explains CSS filter limitations:
  - "This preview uses CSS filters to approximate preset adjustments."
  - "Actual conversions may differ, especially for:"
  - "• Tone curves (not supported in preview)"
  - "• Advanced color grading (simplified in preview)"
  - "• White balance (approximated with sepia + hue rotation)"
  - "The preview is useful for quick comparison, but not pixel-perfect."
- ✅ Tooltip placement:
  - Desktop: Below disclaimer label (popover style)
  - Mobile: Full-width banner at bottom of modal
  - Auto-dismisses after 5 seconds OR on click outside
- ✅ Tooltip accessibility:
  - `aria-describedby` links label to tooltip
  - Screen reader announces tooltip when label focused
  - Keyboard: Esc key dismisses tooltip

**AC-3: No Misleading Claims About Preview Accuracy**
- ✅ Avoid misleading terms in UI:
  - ❌ "Realistic preview" (too strong)
  - ❌ "Accurate preview" (implies pixel-perfect)
  - ❌ "True preview" (misleading)
  - ✅ "Approximate preview" (honest)
  - ✅ "CSS filter-based preview" (transparent about method)
  - ✅ "Quick comparison" (sets expectations)
- ✅ UI copy review:
  - Check all modal text (buttons, labels, help text)
  - Ensure no overpromising language
  - Use neutral, transparent phrasing
- ✅ User expectations managed:
  - Preview presented as "quick comparison tool"
  - Actual conversion outcome emphasized as final result
  - Users encouraged to convert sample files for testing

**AC-4: Documentation Explains Preview vs. Actual Conversion Differences**
- ✅ README section: "Preview System Limitations"
  - Explains CSS filter-based approach
  - Lists supported parameters (exposure, contrast, saturation, hue)
  - Lists unsupported parameters (tone curves, advanced color grading)
  - Provides examples of preview vs. actual conversion differences
- ✅ FAQ entry: "Why doesn't my preview match the converted file exactly?"
  - Answer: "The preview uses CSS filters for speed (<100ms). Actual conversions use full parameter processing, which may produce slightly different results, especially for tone curves and advanced color grading."
- ✅ Comparison table in docs:
  ```
  | Feature          | Preview (CSS Filters) | Actual Conversion |
  | ---------------- | --------------------- | ----------------- |
  | Exposure         | ✅ Supported           | ✅ Supported       |
  | Contrast         | ✅ Supported           | ✅ Supported       |
  | Saturation       | ✅ Supported           | ✅ Supported       |
  | Hue              | ✅ Supported           | ✅ Supported       |
  | Temperature/Tint | ⚠️ Approximated        | ✅ Supported       |
  | Tone Curves      | ❌ Not supported       | ✅ Supported       |
  | Advanced Color   | ❌ Not supported       | ✅ Supported       |
  | Speed            | <100ms                | 100-500ms         |
  ```

**AC-5: Preview Limitations Listed in Documentation**
- ✅ Explicit limitations section in README:
  - **Limitation 1**: Tone curves not supported in preview (Phase 1)
    - Impact: Presets with custom tone curves show simplified preview
    - Workaround: Convert sample file to see actual result
  - **Limitation 2**: Advanced color grading approximated
    - Impact: Split toning, HSL adjustments simplified
    - Workaround: Preview useful for basic adjustments, test conversion for accuracy
  - **Limitation 3**: White balance approximated
    - Impact: Temperature/tint mapped to sepia + hue rotation (not exact)
    - Workaround: Preview shows direction of change, not exact color
  - **Limitation 4**: Browser CSS filter support varies
    - Impact: Older browsers (Safari <9.1) may not show filters
    - Workaround: Use modern browser (Chrome, Firefox, Safari 14+)
- ✅ Phase roadmap in docs:
  - **Phase 1 (Current)**: CSS filter-based preview (basic parameters)
  - **Phase 2 (Future)**: Canvas-based preview (tone curves, advanced color)
  - **Phase 3 (Future)**: WebGL-based preview (pixel-perfect, all parameters)

**AC-6: User Expectations Managed Through Transparent Communication**
- ✅ Modal disclaimer sets expectations:
  - Preview = "approximate" (not exact)
  - Actual conversion = "final result" (authoritative)
  - Preview useful for "quick comparison" (speed vs. accuracy trade-off)
- ✅ First-time user guidance (optional):
  - On first modal open: Show tooltip automatically
  - Tooltip explains preview limitations
  - User can dismiss with "Got it" button
  - Preference saved in localStorage (don't show again)
- ✅ "Try a sample conversion" CTA:
  - Encourage users to convert sample files
  - Test actual conversion quality vs. preview
  - Build confidence in Recipe's accuracy
- ✅ Community transparency:
  - GitHub README: Honest about preview limitations
  - Issue template: Ask users to provide sample files for debugging
  - Documentation: Clear about Phase 1 limitations and future roadmap

## Tasks / Subtasks

### Task 1: Add Disclaimer Label to Modal (AC-1)
- [ ] Update modal HTML in `web/index.html`:
  ```html
  <!-- Preview Modal -->
  <div id="preview-modal" class="modal" role="dialog" aria-modal="true">
    <div class="modal-container">
      <div class="modal-header">
        <h2 id="modal-title">Preset Preview</h2>
        <button class="modal-close" aria-label="Close preview">×</button>
      </div>

      <!-- Disclaimer Label -->
      <div class="preview-disclaimer" role="status" aria-live="polite" tabindex="0" aria-describedby="disclaimer-tooltip">
        <svg class="disclaimer-icon" aria-hidden="true">
          <use href="#icon-info"></use>
        </svg>
        <span class="disclaimer-text">Approximate preview using CSS filters</span>
        <svg class="disclaimer-help-icon" aria-hidden="true">
          <use href="#icon-help"></use>
        </svg>
      </div>

      <!-- Tooltip (hidden by default) -->
      <div id="disclaimer-tooltip" class="disclaimer-tooltip" hidden>
        <p>This preview uses CSS filters to approximate preset adjustments.</p>
        <p>Actual conversions may differ, especially for:</p>
        <ul>
          <li>Tone curves (not supported in preview)</li>
          <li>Advanced color grading (simplified in preview)</li>
          <li>White balance (approximated with sepia + hue rotation)</li>
        </ul>
        <p>The preview is useful for quick comparison, but not pixel-perfect.</p>
        <button class="tooltip-close" aria-label="Close tooltip">Got it</button>
      </div>

      <!-- Rest of modal content (tabs, slider, etc.) -->
      ...
    </div>
  </div>
  ```
- [ ] Add disclaimer CSS in `web/css/modal.css`:
  ```css
  .preview-disclaimer {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 16px;
    margin-bottom: 16px;
    background: #FEF3C7; /* Yellow-100 */
    border: 1px solid #F59E0B; /* Yellow-600 */
    border-radius: 8px;
    font-size: 14px;
    color: #78350F; /* Yellow-900 */
    cursor: pointer;
    transition: background-color 200ms ease;
  }

  .preview-disclaimer:hover,
  .preview-disclaimer:focus {
    background: #FDE68A; /* Yellow-200 */
    outline: 2px solid #F59E0B;
    outline-offset: 2px;
  }

  .disclaimer-icon {
    width: 20px;
    height: 20px;
    fill: #F59E0B; /* Yellow-600 */
  }

  .disclaimer-text {
    flex: 1;
    font-weight: 500;
  }

  .disclaimer-help-icon {
    width: 16px;
    height: 16px;
    fill: #F59E0B;
  }
  ```

### Task 2: Implement Tooltip Functionality (AC-2)
- [ ] Add tooltip JavaScript in `web/js/disclaimer.js`:
  ```javascript
  const disclaimerLabel = document.querySelector('.preview-disclaimer');
  const disclaimerTooltip = document.getElementById('disclaimer-tooltip');
  const tooltipClose = document.querySelector('.tooltip-close');

  // Show tooltip on hover/click
  disclaimerLabel.addEventListener('mouseenter', showTooltip);
  disclaimerLabel.addEventListener('click', showTooltip);
  disclaimerLabel.addEventListener('focus', showTooltip);

  // Hide tooltip on click outside or Esc
  document.addEventListener('click', (e) => {
    if (!disclaimerLabel.contains(e.target) && !disclaimerTooltip.contains(e.target)) {
      hideTooltip();
    }
  });

  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      hideTooltip();
    }
  });

  tooltipClose.addEventListener('click', hideTooltip);

  function showTooltip() {
    disclaimerTooltip.hidden = false;

    // Auto-dismiss after 10 seconds
    clearTimeout(window.tooltipTimeout);
    window.tooltipTimeout = setTimeout(hideTooltip, 10000);
  }

  function hideTooltip() {
    disclaimerTooltip.hidden = true;
    clearTimeout(window.tooltipTimeout);
  }

  // Keyboard: Enter/Space to show tooltip
  disclaimerLabel.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      showTooltip();
    }
  });
  ```
- [ ] Add tooltip CSS:
  ```css
  .disclaimer-tooltip {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin-top: 8px;
    padding: 16px;
    background: white;
    border: 1px solid #D1D5DB; /* Gray-300 */
    border-radius: 8px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
    z-index: 1000;
    font-size: 14px;
    color: #374151; /* Gray-700 */
  }

  .disclaimer-tooltip p {
    margin: 0 0 8px 0;
  }

  .disclaimer-tooltip ul {
    margin: 8px 0;
    padding-left: 20px;
  }

  .disclaimer-tooltip li {
    margin-bottom: 4px;
  }

  .tooltip-close {
    margin-top: 12px;
    padding: 8px 16px;
    background: var(--color-primary);
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  /* Mobile: Full-width tooltip */
  @media (max-width: 600px) {
    .disclaimer-tooltip {
      position: fixed;
      bottom: 0;
      left: 0;
      right: 0;
      top: auto;
      margin: 0;
      border-radius: 16px 16px 0 0;
    }
  }
  ```

### Task 3: Review UI Copy for Misleading Claims (AC-3)
- [ ] Audit all modal text:
  - Modal title: "Preset Preview" ✅ (neutral)
  - Disclaimer label: "Approximate preview using CSS filters" ✅ (honest)
  - Tabs: "Portrait", "Landscape", "Product" ✅ (neutral)
  - Buttons: "Convert Now", "Cancel" ✅ (neutral)
  - Parameter display: "Exposure +0.7 • Contrast +15" ✅ (factual)
- [ ] Replace any misleading terms:
  - Search codebase for: "realistic", "accurate", "true", "exact", "perfect"
  - Replace with: "approximate", "preview", "quick comparison"
- [ ] Add neutral phrasing guide to docs:
  ```markdown
  ### UI Copy Guidelines: Preview Modal

  **Approved Terms:**
  - "Approximate preview"
  - "CSS filter-based preview"
  - "Quick comparison"
  - "Preview vs. actual conversion"

  **Avoid:**
  - "Realistic preview" (too strong)
  - "Accurate preview" (implies pixel-perfect)
  - "True preview" (misleading)
  - "Exact preview" (overpromising)
  ```

### Task 4: Add Documentation Section (AC-4)
- [ ] Create `docs/preview-limitations.md`:
  ```markdown
  # Preview System Limitations

  Recipe's preview system uses CSS filters for instant comparison (<100ms). This approach prioritizes speed over pixel-perfect accuracy.

  ## How It Works

  1. Upload a preset file (.np3, .xmp, .lrtemplate)
  2. Recipe parses preset parameters (exposure, contrast, etc.)
  3. Parameters mapped to CSS filter functions:
     - Exposure → `brightness()`
     - Contrast → `contrast()`
     - Saturation → `saturate()`
     - Hue → `hue-rotate()`
     - Temperature/Tint → `sepia()` + `hue-rotate()` (approximated)
  4. Filters applied to reference image in real-time

  ## Supported vs. Unsupported Parameters

  | Feature          | Preview (CSS Filters) | Actual Conversion | Notes                              |
  | ---------------- | --------------------- | ----------------- | ---------------------------------- |
  | Exposure         | ✅ Supported           | ✅ Supported       | Accurate                           |
  | Contrast         | ✅ Supported           | ✅ Supported       | Accurate                           |
  | Saturation       | ✅ Supported           | ✅ Supported       | Accurate                           |
  | Hue              | ✅ Supported           | ✅ Supported       | Accurate                           |
  | Temperature/Tint | ⚠️ Approximated        | ✅ Supported       | Simplified in preview              |
  | Tone Curves      | ❌ Not supported       | ✅ Supported       | Phase 1 limitation                 |
  | Advanced Color   | ❌ Not supported       | ✅ Supported       | HSL, split toning simplified       |
  | Speed            | <100ms                | 100-500ms         | Preview instant, conversion slower |

  ## Known Limitations

  ### 1. Tone Curves Not Supported (Phase 1)

  **Impact:** Presets with custom tone curves (S-curves, highlights/shadows) show simplified preview.

  **Example:**
  - Preset has S-curve: Boosted shadows, crushed highlights
  - Preview shows: Linear exposure adjustment (no curve)
  - Actual conversion: Full S-curve applied

  **Workaround:** Convert a sample file to see actual tone curve result.

  ### 2. Advanced Color Grading Approximated

  **Impact:** HSL adjustments, split toning, color calibration simplified in preview.

  **Example:**
  - Preset has split toning: Blue shadows, orange highlights
  - Preview shows: Slight hue shift (approximation)
  - Actual conversion: Full split toning applied

  **Workaround:** Preview useful for basic color direction, convert sample for accuracy.

  ### 3. White Balance Approximated

  **Impact:** Temperature/tint mapped to `sepia()` + `hue-rotate()` (not exact color science).

  **Example:**
  - Preset has +20 temperature (warmer)
  - Preview shows: Sepia tone + yellow hue shift (approximation)
  - Actual conversion: Accurate white balance adjustment

  **Workaround:** Preview shows direction (warmer/cooler), convert sample for exact color.

  ### 4. Browser Support Varies

  **Impact:** Older browsers may not support CSS filters.

  **Minimum Browser Versions:**
  - Chrome 18+ ✅
  - Firefox 35+ ✅
  - Safari 9.1+ ✅
  - Edge 12+ ✅
  - IE 11 ❌ (not supported)

  **Workaround:** Use modern browser for preview. Conversion works on all browsers (server-side).

  ## Future Roadmap

  ### Phase 1 (Current): CSS Filters
  - Speed: <100ms
  - Parameters: Exposure, Contrast, Saturation, Hue (basic)
  - Limitations: No tone curves, approximated white balance

  ### Phase 2 (Planned): Canvas-based Preview
  - Speed: <500ms
  - Parameters: All except advanced color grading
  - Improvements: Tone curves supported, accurate white balance

  ### Phase 3 (Future): WebGL-based Preview
  - Speed: <200ms
  - Parameters: All parameters (pixel-perfect)
  - Improvements: GPU-accelerated, full color science

  ## FAQ

  **Q: Why doesn't my preview match the converted file exactly?**

  A: The preview uses CSS filters for speed (<100ms). Actual conversions use full parameter processing, which may produce slightly different results, especially for tone curves and advanced color grading.

  **Q: Can I trust the preview for basic adjustments?**

  A: Yes! Exposure, contrast, saturation, and hue are accurate. Temperature/tint are approximated but show the correct direction.

  **Q: Should I skip the preview and just convert?**

  A: No! The preview is useful for quick comparison. Use it to evaluate presets, then convert a sample file to verify accuracy.

  **Q: Will Phase 2 be more accurate?**

  A: Yes. Phase 2 will add tone curve support and accurate white balance using Canvas API.
  ```
- [ ] Add FAQ entry to `docs/FAQ.md`:
  ```markdown
  ## Q: Why doesn't my preview match the converted file exactly?

  The preview uses CSS filters for instant comparison (<100ms). This approach prioritizes speed over pixel-perfect accuracy.

  **What's accurate in preview:**
  - Exposure, Contrast, Saturation, Hue ✅

  **What's approximated in preview:**
  - Temperature/Tint ⚠️ (simplified)
  - Tone Curves ❌ (not supported in Phase 1)
  - Advanced Color Grading ❌ (not supported in Phase 1)

  **Recommendation:** Use the preview for quick comparison, then convert a sample file to verify accuracy.

  See [Preview Limitations](preview-limitations.md) for details.
  ```

### Task 5: List Preview Limitations in README (AC-5)
- [ ] Add section to `README.md`:
  ```markdown
  ### Preview System (Epic 11)

  Recipe includes a CSS filter-based preview system for instant preset comparison.

  #### Features

  - **Instant Preview**: <100ms rendering (no processing delay)
  - **Before/After Slider**: Drag to compare original vs. filtered
  - **Reference Images**: Portrait, Landscape, Product tabs
  - **Keyboard/Touch Support**: Arrow keys, swipe gestures, double-tap

  #### Limitations (Phase 1)

  ⚠️ **The preview is approximate, not pixel-perfect.**

  | Parameter        | Preview Accuracy |
  | ---------------- | ---------------- |
  | Exposure         | ✅ Accurate       |
  | Contrast         | ✅ Accurate       |
  | Saturation       | ✅ Accurate       |
  | Hue              | ✅ Accurate       |
  | Temperature/Tint | ⚠️ Approximated   |
  | Tone Curves      | ❌ Not supported  |
  | Advanced Color   | ❌ Not supported  |

  **Recommendation:** Use preview for quick comparison. Convert a sample file to verify accuracy.

  See [Preview Limitations](docs/preview-limitations.md) for technical details.

  #### Roadmap

  - **Phase 1 (Current)**: CSS filters (basic parameters)
  - **Phase 2 (Planned)**: Canvas API (tone curves, accurate white balance)
  - **Phase 3 (Future)**: WebGL (pixel-perfect, all parameters)
  ```

### Task 6: First-Time User Guidance (AC-6, Optional)
- [ ] Add first-time tooltip logic:
  ```javascript
  // Check if user has seen tooltip before
  const hasSeenTooltip = localStorage.getItem('preview-disclaimer-seen');

  if (!hasSeenTooltip) {
    // Show tooltip automatically on first modal open
    showTooltip();

    // Add "Don't show again" checkbox
    const tooltipContent = document.getElementById('disclaimer-tooltip');
    const dontShowAgain = document.createElement('label');
    dontShowAgain.innerHTML = `
      <input type="checkbox" id="dont-show-again"> Don't show this again
    `;
    tooltipContent.appendChild(dontShowAgain);

    // Save preference when "Got it" clicked
    tooltipClose.addEventListener('click', () => {
      if (document.getElementById('dont-show-again').checked) {
        localStorage.setItem('preview-disclaimer-seen', 'true');
      }
      hideTooltip();
    });
  }
  ```

### Task 7: Manual Testing
- [ ] Visual testing:
  - Disclaimer label visible in modal ✅
  - Warning icon and text readable ✅
  - Contrasts with background (WCAG AA 4.5:1) ✅
  - Label stands out (not buried in UI) ✅
- [ ] Tooltip testing:
  - Hover: Tooltip appears below label ✅
  - Click: Tooltip appears on click ✅
  - Keyboard: Focus label + Enter shows tooltip ✅
  - Mobile: Tooltip appears as bottom banner ✅
  - Auto-dismiss: Tooltip disappears after 10 seconds ✅
  - Esc key: Tooltip dismisses ✅
- [ ] UI copy audit:
  - Search for misleading terms ("realistic", "accurate") ✅
  - Replace with neutral terms ("approximate", "preview") ✅
  - All modal text reviewed ✅
- [ ] Documentation review:
  - README section: Preview limitations listed ✅
  - FAQ entry: "Why doesn't preview match?" ✅
  - Docs: Preview vs. actual conversion table ✅

### Task 8: Accessibility Testing
- [ ] Screen reader testing (NVDA, VoiceOver):
  - Disclaimer label announced: "Warning: Approximate preview using CSS filters" ✅
  - Tooltip announced when label focused ✅
  - `aria-describedby` links label to tooltip ✅
- [ ] Keyboard navigation:
  - Tab to disclaimer label ✅
  - Enter/Space: Show tooltip ✅
  - Esc: Dismiss tooltip ✅
  - Focus visible: Clear outline on label ✅
- [ ] Color contrast:
  - Disclaimer label text: 4.5:1 minimum (WCAG AA) ✅
  - Tooltip text: 4.5:1 minimum ✅
  - Warning icon: Sufficient contrast ✅

### Task 9: Documentation Updates
- [ ] Update `docs/index.md` with link to preview-limitations.md
- [ ] Update landing page FAQ section (Epic 7) with preview limitation question
- [ ] Add "Preview System" section to docs navigation
- [ ] Cross-reference from README to detailed docs

### Task 10: Community Transparency
- [ ] GitHub README: Honest about limitations ✅
- [ ] Issue template: Request sample files for debugging
  ```markdown
  ### Bug Report: Preview Mismatch

  **Describe the issue:**
  The preview doesn't match my converted file.

  **Sample files (required):**
  - [ ] Preset file (.np3, .xmp, .lrtemplate)
  - [ ] Original image
  - [ ] Converted image (actual output)
  - [ ] Screenshot of preview (expected vs. actual)

  **Expected behavior:**
  Preview should closely match converted file.

  **Actual behavior:**
  Preview shows [describe difference].

  **Note:** Phase 1 preview has known limitations (tone curves, advanced color). See [Preview Limitations](docs/preview-limitations.md).
  ```
- [ ] CONTRIBUTING.md: Note about preview limitations in Phase 1

## Dev Notes

### Learnings from Previous Story

**From Story 11-4-preview-slider-interaction (Status: drafted)**

Previous story not yet implemented. Story 11.5 adds disclaimer and transparency features to the preview modal from Stories 11.3 and 11.4.

**Story 11.4 Coverage:**
- Advanced slider interactions (keyboard shortcuts, touch gestures)
- Visual feedback (percentage labels, Before/After labels)
- Accessibility (ARIA live regions, focus visible)
- Performance optimizations (GPU acceleration, RAF)

**Story 11.5 Addition:**
- Disclaimer label explaining CSS filter limitations
- Tooltip with detailed explanation of preview vs. conversion differences
- Documentation of limitations (tone curves, advanced color grading)
- User expectation management (preview = approximate, not pixel-perfect)

[Source: docs/stories/11-4-preview-slider-interaction.md]

### Architecture Alignment

**Tech Spec Epic 11 Alignment:**

Story 11.5 implements **AC-4: Accuracy Communication** from tech-spec-epic-11.md.

**Disclaimer Requirements:**

```
Component          Content                                     Placement
-----------------  ------------------------------------------  ---------------------------
Label              "Approximate preview using CSS filters"     Top of modal OR below slider
Tooltip            Detailed limitations explanation            Below label (desktop), bottom banner (mobile)
Documentation      Preview vs. conversion comparison table     README, docs/preview-limitations.md
FAQ                "Why doesn't preview match?"                docs/FAQ.md
UI Copy            Neutral terms (avoid "realistic")           All modal text
```

**User Expectation Management Strategy:**

```
Communication Layer   Message
--------------------  --------------------------------------------------------
Modal Disclaimer      "This is approximate (not exact)"
Tooltip               "Tone curves not supported, white balance approximated"
Documentation         "Phase 1 = basic, Phase 2 = tone curves, Phase 3 = pixel-perfect"
FAQ                   "Use preview for quick comparison, convert sample to verify"
Community             "Honest about limitations, roadmap for improvements"
```

[Source: docs/tech-spec-epic-11.md#AC-4]

### Disclaimer Label Design

**Why Warning/Info Badge Style?**

| Style Option  | Pros                             | Cons                          | Decision    |
| ------------- | -------------------------------- | ----------------------------- | ----------- |
| Plain text    | Minimal, unobtrusive             | Easy to miss, low prominence  | ❌ No        |
| Warning badge | High prominence, draws attention | May alarm users unnecessarily | ✅ Yes       |
| Info badge    | Friendly, informative            | Lower urgency than warning    | Alternative |
| Modal header  | Always visible                   | Takes valuable header space   | ❌ No        |

**Recipe Choice: Warning badge (yellow background, ⚠️ icon)**

**Rationale:**
- High prominence: Users can't miss it
- Honest communication: Sets expectations upfront
- Not alarming: Yellow (caution) vs. red (error)
- Persistent: Always visible (not dismissible)

**Alternative: Info badge (blue background, ℹ️ icon)**
- Less alarming than warning
- Still prominent and informative
- May be preferred if users find warning too strong

[Source: UI Patterns - Warning vs. Info Messages]

### Tooltip Content Strategy

**Structure:**

1. **What**: "This preview uses CSS filters..."
2. **Why**: "...to approximate preset adjustments."
3. **Limitations**: "Actual conversions may differ for: [list]"
4. **Reassurance**: "Useful for quick comparison, not pixel-perfect."

**Why This Structure?**

- **What + Why**: Context first (user understands the approach)
- **Limitations**: Specifics second (user knows what to expect)
- **Reassurance**: Positive note last (user still finds preview valuable)

**Tone:**
- Honest: "not pixel-perfect" (transparent)
- Neutral: "may differ" (not alarming)
- Helpful: "useful for quick comparison" (emphasizes value)

[Source: UX Writing Principles - Honest, Helpful, Human]

### CSS Filter Limitations (Technical)

**Why Can't CSS Filters Do Tone Curves?**

CSS filters are **global transformations** (apply to entire image). Tone curves are **per-pixel transformations** (different adjustment for each brightness level).

**Example: S-Curve (Contrast Enhancement)**

```
Brightness Input → Brightness Output
0 (black)        → 10 (slightly lifted shadows)
128 (midtone)    → 128 (unchanged)
255 (white)      → 245 (slightly crushed highlights)
```

**CSS Filter Equivalent:**

```css
filter: brightness(1.1) contrast(1.2);
/* Global: Brightens everything, increases contrast everywhere */
/* Cannot target shadows separately from highlights */
```

**Result:** CSS filters can't replicate S-curve behavior (different adjustments per tone).

**Solution (Phase 2):** Canvas API with pixel manipulation.

```javascript
// Canvas API: Per-pixel tone curve
const imageData = ctx.getImageData(0, 0, width, height);
for (let i = 0; i < imageData.data.length; i += 4) {
  const r = imageData.data[i];
  const g = imageData.data[i + 1];
  const b = imageData.data[i + 2];

  // Apply tone curve lookup table
  imageData.data[i] = toneCurveLUT[r];
  imageData.data[i + 1] = toneCurveLUT[g];
  imageData.data[i + 2] = toneCurveLUT[b];
}
ctx.putImageData(imageData, 0, 0);
```

[Source: CSS Filters vs. Canvas API - MDN Web Docs]

### Phase Roadmap Justification

**Why Not Canvas/WebGL in Phase 1?**

| Approach    | Speed  | Accuracy     | Complexity | Phase |
| ----------- | ------ | ------------ | ---------- | ----- |
| CSS Filters | <100ms | 80% (basic)  | Low        | 1     |
| Canvas API  | <500ms | 95% (curves) | Medium     | 2     |
| WebGL       | <200ms | 100% (all)   | High       | 3     |

**Phase 1 Trade-off:**
- Priority: **Speed** (instant preview, <100ms)
- Acceptable: **80% accuracy** (basic parameters only)
- Rationale: Users want quick comparison, not pixel-perfect

**Phase 2 Upgrade:**
- Priority: **Accuracy** (tone curves, white balance)
- Acceptable: **500ms delay** (still fast, noticeable lag)
- Rationale: Users willing to wait 500ms for accurate preview

**Phase 3 Vision:**
- Goal: **Best of both** (speed + accuracy)
- Technology: **WebGL** (GPU-accelerated pixel processing)
- Timeline: Long-term (complex implementation)

[Source: Recipe Product Roadmap - Phase A Enhancements]

### User Expectation Management Psychology

**Transparency Builds Trust:**

Research shows users prefer **honest limitations** over **hidden surprises**.

**Example: Misleading Preview**
- User sees preview: "Looks great!"
- Converts file: "Doesn't match preview. Recipe is broken."
- Result: **Loss of trust**, negative review

**Example: Transparent Preview**
- User sees disclaimer: "Preview approximate, not exact"
- User converts sample: "Close enough, as expected"
- Result: **Trust maintained**, expectations met

**Recipe Strategy:**
- **Honest**: Upfront about limitations
- **Helpful**: Explains why limitations exist
- **Hopeful**: Roadmap shows future improvements

[Source: The Psychology of Transparency - Nielsen Norman Group]

### WCAG Color Contrast Requirements

**Disclaimer Label Contrast:**

```
Yellow-100 background: #FEF3C7
Yellow-900 text: #78350F

Contrast ratio: 8.2:1 ✅ (exceeds WCAG AAA 7:1)
```

**Why Yellow (Warning) Color?**

- **Yellow**: Universal "caution" color (not error, not success)
- **High contrast**: Dark text on light background (8.2:1)
- **Friendly**: Less alarming than red

**Alternative: Blue (Info) Color**

```
Blue-100 background: #DBEAFE
Blue-900 text: #1E3A8A

Contrast ratio: 9.1:1 ✅ (exceeds WCAG AAA 7:1)
```

**Decision:** Yellow (warning) preferred for prominence, blue (info) acceptable alternative.

[Source: WCAG 2.1 SC 1.4.3 - Contrast (Minimum)]

### Project Structure Notes

**New Files Created (Story 11.5):**
```
web/
├── js/
│   └── disclaimer.js            (Tooltip show/hide, first-time guidance)
└── css/
    └── disclaimer.css           (Disclaimer label, tooltip styles) OR add to modal.css

docs/
├── preview-limitations.md       (Detailed technical limitations)
└── FAQ.md                       (Updated with preview limitation question)
```

**Modified Files:**
- `web/index.html` - Add disclaimer label, tooltip HTML
- `web/css/modal.css` - Add disclaimer styles (alternative to separate file)
- `README.md` - Add preview limitations section
- `docs/index.md` - Add link to preview-limitations.md
- `.github/ISSUE_TEMPLATE/bug_report.md` - Add preview mismatch template

**Integration Points:**
- Story 11.3: Disclaimer appears in modal created in 11.3
- Story 11.1: Tooltip references CSS filter implementation from 11.1
- Epic 7: FAQ entry cross-references landing page FAQ

[Source: docs/tech-spec-epic-11.md#Services-and-Modules]

### Testing Strategy

**Manual Tests:**
- Visual: Disclaimer label prominent, readable, contrasts well ✅
- Tooltip: Appears on hover/click/keyboard, auto-dismisses ✅
- UI copy: No misleading terms ("realistic", "accurate") ✅
- Documentation: README, FAQ, preview-limitations.md complete ✅

**Accessibility Tests:**
- Screen reader (NVDA): Disclaimer announced, tooltip read ✅
- Keyboard: Tab to label, Enter shows tooltip, Esc dismisses ✅
- Color contrast: Label text 4.5:1 minimum (WCAG AA) ✅

**User Testing (Optional):**
- Show preview to 5 users
- Ask: "What do you expect from the converted file?"
- Validate: Users understand preview is approximate

[Source: docs/tech-spec-epic-11.md#Test-Strategy-Summary]

### Known Risks

**RISK-66: Users may ignore disclaimer and expect pixel-perfect preview**
- **Impact**: Users disappointed when conversion differs from preview
- **Mitigation**: Prominent disclaimer, tooltip auto-shows first time, documentation
- **Acceptable**: Some users will still be surprised, but transparency minimizes complaints

**RISK-67: Disclaimer may alarm users unnecessarily**
- **Impact**: Users avoid using preview feature
- **Mitigation**: Use friendly warning (yellow) not error (red), emphasize "useful for quick comparison"
- **Test**: User testing to validate tone is appropriate

**RISK-68: Tooltip may be too long for mobile**
- **Impact**: Tooltip text truncated or hard to read on small screens
- **Mitigation**: Mobile: Full-width bottom banner with larger text, scrollable if needed
- **Test**: iPhone SE (smallest modern screen), Android budget phones

**RISK-69: Documentation may be too technical for non-technical users**
- **Impact**: Users don't understand limitations (tone curves, color science)
- **Mitigation**: FAQ uses simple language, README has visual table
- **Acceptable**: Advanced users read docs, casual users rely on disclaimer

[Source: docs/tech-spec-epic-11.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-11.md#AC-4] - Accuracy communication requirements
- [Source: docs/stories/11-3-preview-modal-interface.md] - Modal structure
- [Source: docs/stories/11-1-css-filter-mapping.md] - CSS filter implementation
- [UI Patterns - Warning vs. Info Messages](https://ui-patterns.com/patterns/FeedbackMessages)
- [UX Writing Principles - Honest, Helpful, Human](https://www.nngroup.com/articles/tone-of-voice-dimensions/)
- [CSS Filters vs. Canvas API - MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Tutorial/Pixel_manipulation_with_canvas)
- [The Psychology of Transparency - Nielsen Norman Group](https://www.nngroup.com/articles/transparency/)
- [WCAG 2.1 SC 1.4.3 - Contrast (Minimum)](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html)

## Dev Agent Record

### Context Reference

- docs/stories/11-5-preview-disclaimer.context.xml (Generated: 2025-11-09)

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
