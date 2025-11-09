# Epic Technical Specification: Image Preview System (Phase 1)

Date: 2025-11-08
Author: Justin
Epic ID: epic-11
Status: Draft

---

## Overview

Epic 11 adds visual preset preview to Recipe's web interface, allowing users to see what a preset will do before converting. Phase 1 uses CSS filters for instant, approximate previews—trading pixel-perfect accuracy for speed and simplicity. This transforms Recipe from a "blind conversion" tool into a confident, visual experience where users can explore presets on reference images (Portrait, Landscape, Product) and use a before/after slider to compare effects.

The implementation leverages browser-native CSS filter functions (brightness, contrast, saturate, hue-rotate, sepia) to approximate UniversalRecipe parameters. Reference images are bundled as optimized WebP files (<200KB each). The preview modal integrates with Epic 10's upload interface, providing an optional "Preview" step before conversion.

**Phase 2 Vision (Future)**: WebAssembly-based accurate preview using Photon library for pixel-perfect rendering. Phase 1 CSS filters are the MVP commitment.

## Objectives and Scope

**In Scope (Phase 1):**
- CSS filter-based preview mapping UniversalRecipe → filter functions
- Three bundled reference images (Portrait, Landscape, Product) as WebP
- Preview modal with before/after comparison slider
- Tab navigation to switch between reference images
- Preset parameter display (e.g., "Exposure +0.7 • Contrast +15")
- Clear accuracy communication ("Approximate preview using CSS filters")
- Performance: <100ms preview rendering, 60fps slider interaction

**Out of Scope (Phase 1):**
- WebAssembly-based accurate preview (defer to Phase 2)
- Tone curve preview (CSS filters don't support curves)
- Custom image upload for preview (use bundled references only)
- Side-by-side before/after view (slider only)
- Preview export/download (preview is ephemeral, conversion is primary)
- Advanced color grading preview (HSL, parametric curves)

## System Architecture Alignment

**Components:**
- **New**: `web/js/preview.js` - Preview modal, CSS filter mapping, slider logic
- **New**: `web/css/preview.css` - Modal overlay, slider styles, tab navigation
- **New**: `web/images/preview-portrait.webp` - Portrait reference image (~180KB)
- **New**: `web/images/preview-landscape.webp` - Landscape reference image (~190KB)
- **New**: `web/images/preview-product.webp` - Product/Still-life reference image (~150KB)
- **Enhanced**: `web/js/app.js` - Integrate preview modal into upload workflow

**Integration Points:**
- Epic 10 upload interface: "Preview" button per uploaded file
- UniversalRecipe parameters: Map to CSS filter values
- Before/after slider: Integrated via CSS clip-path or two overlaid images

**Constraints:**
- Must maintain <100ms preview rendering (instant feel)
- Reference images total <600KB (fast download on 3G)
- CSS filters only (no JavaScript image processing in Phase 1)
- Works across all modern browsers (Chrome, Firefox, Safari, Edge)

## Detailed Design

### Services and Modules

| Module | Responsibility | Inputs | Outputs | Owner |
| ------ | -------------- | ------ | ------- | ----- |
| `js/preview.js` | Preview modal, CSS filter mapping, slider | `UniversalRecipe`, reference image | Modal DOM | Dev (Epic 11) |
| `css/preview.css` | Modal styles, slider, tabs | - | Styled modal | Dev (Epic 11) |
| `images/preview-*.webp` | Reference images | - | Visual previews | Dev (Epic 11) |

### Data Models and Contracts

**CSS Filter Mapping (preview.js):**

```javascript
// web/js/preview.js

/**
 * Map UniversalRecipe parameters to CSS filter string
 * @param {Object} recipe - UniversalRecipe parameters
 * @returns {string} CSS filter value (e.g., "brightness(1.2) contrast(1.15)")
 */
export function recipeToCSS Filters(recipe) {
    const filters = [];

    // Exposure (-2.0 to +2.0) → brightness(0% to 200%)
    if (recipe.exposure !== undefined && recipe.exposure !== 0) {
        const brightness = 100 + (recipe.exposure * 50); // -2.0 = 0%, 0 = 100%, +2.0 = 200%
        filters.push(`brightness(${brightness}%)`);
    }

    // Contrast (-1.0 to +1.0) → contrast(0% to 200%)
    if (recipe.contrast !== undefined && recipe.contrast !== 0) {
        const contrast = 100 + (recipe.contrast * 100); // -1.0 = 0%, 0 = 100%, +1.0 = 200%
        filters.push(`contrast(${contrast}%)`);
    }

    // Saturation (-1.0 to +1.0) → saturate(0% to 200%)
    if (recipe.saturation !== undefined && recipe.saturation !== 0) {
        const saturate = 100 + (recipe.saturation * 100); // -1.0 = 0%, 0 = 100%, +1.0 = 200%
        filters.push(`saturate(${saturate}%)`);
    }

    // Temperature (-100 to +100) → sepia() + hue-rotate() approximation
    // Warm = sepia + orange hue, Cool = blue hue
    if (recipe.temperature !== undefined && recipe.temperature !== 0) {
        if (recipe.temperature > 0) {
            // Warm: sepia for warmth, hue-rotate toward orange
            const sepia = recipe.temperature / 100 * 0.3; // Max 30% sepia
            const hue = recipe.temperature / 100 * 10; // Max 10deg toward orange
            filters.push(`sepia(${sepia})`);
            filters.push(`hue-rotate(${hue}deg)`);
        } else {
            // Cool: hue-rotate toward blue (negative temperature)
            const hue = recipe.temperature / 100 * 20; // Max -20deg toward blue
            filters.push(`hue-rotate(${hue}deg)`);
        }
    }

    // Tint (-100 to +100) → hue-rotate() (green/magenta shift)
    if (recipe.tint !== undefined && recipe.tint !== 0) {
        const hue = recipe.tint / 100 * 15; // Max ±15deg for tint
        filters.push(`hue-rotate(${hue}deg)`);
    }

    // Vibrance approximation (not perfect, use saturate as fallback)
    if (recipe.vibrance !== undefined && recipe.vibrance !== 0) {
        const vibrance = 100 + (recipe.vibrance * 30); // Subtle saturation boost
        filters.push(`saturate(${vibrance}%)`);
    }

    return filters.length > 0 ? filters.join(' ') : 'none';
}

/**
 * Generate human-readable parameter summary
 * @param {Object} recipe - UniversalRecipe parameters
 * @returns {string} Summary (e.g., "Exposure +0.7 • Contrast +15 • Warmth +10")
 */
export function formatRecipeSummary(recipe) {
    const parts = [];

    if (recipe.exposure) parts.push(`Exposure ${recipe.exposure > 0 ? '+' : ''}${recipe.exposure.toFixed(1)}`);
    if (recipe.contrast) parts.push(`Contrast ${recipe.contrast > 0 ? '+' : ''}${(recipe.contrast * 100).toFixed(0)}`);
    if (recipe.saturation) parts.push(`Saturation ${recipe.saturation > 0 ? '+' : ''}${(recipe.saturation * 100).toFixed(0)}`);
    if (recipe.temperature) parts.push(`${recipe.temperature > 0 ? 'Warmth' : 'Cool'} ${Math.abs(recipe.temperature).toFixed(0)}`);
    if (recipe.tint) parts.push(`Tint ${recipe.tint > 0 ? '+' : ''}${recipe.tint.toFixed(0)}`);

    return parts.length > 0 ? parts.join(' • ') : 'No adjustments';
}
```

**Preview Modal HTML Structure:**

```html
<!-- Modal created dynamically in preview.js -->
<div id="preview-modal" class="preview-modal" role="dialog" aria-labelledby="preview-title">
    <div class="preview-modal__overlay" aria-hidden="true"></div>
    <div class="preview-modal__content">
        <button class="preview-modal__close" aria-label="Close preview">&times;</button>

        <h2 id="preview-title">Preview: filename.np3</h2>
        <p class="preview-modal__disclaimer">Approximate preview using CSS filters</p>

        <!-- Reference image tabs -->
        <div class="preview-modal__tabs">
            <button class="preview-tab" data-image="portrait">Portrait</button>
            <button class="preview-tab" data-image="landscape">Landscape</button>
            <button class="preview-tab" data-image="product">Product</button>
        </div>

        <!-- Before/After Slider -->
        <div class="preview-modal__slider-container">
            <div class="preview-slider">
                <img src="images/preview-portrait.webp" alt="Before" class="preview-slider__before">
                <img src="images/preview-portrait.webp" alt="After" class="preview-slider__after" style="filter: brightness(120%) contrast(115%);">
                <input type="range" min="0" max="100" value="50" class="preview-slider__handle" aria-label="Adjust before/after comparison">
                <div class="preview-slider__labels">
                    <span>Before</span>
                    <span>After</span>
                </div>
            </div>
        </div>

        <!-- Parameter summary -->
        <p class="preview-modal__params">Exposure +0.7 • Contrast +15 • Warmth +10</p>

        <!-- Actions -->
        <div class="preview-modal__actions">
            <button class="btn btn--primary">Convert Now</button>
            <button class="btn btn--secondary">Cancel</button>
        </div>
    </div>
</div>
```

**Before/After Slider Implementation:**

```css
/* web/css/preview.css */
.preview-slider {
    position: relative;
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
    overflow: hidden;
}

.preview-slider__before,
.preview-slider__after {
    display: block;
    width: 100%;
    height: auto;
    user-select: none;
}

.preview-slider__before {
    position: relative;
    z-index: 1;
}

.preview-slider__after {
    position: absolute;
    top: 0;
    left: 0;
    z-index: 2;
    clip-path: inset(0 0 0 50%); /* Initially show right half (50%) */
    transition: clip-path 0.05s ease; /* Smooth slider drag */
}

.preview-slider__handle {
    position: absolute;
    top: 50%;
    left: 0;
    width: 100%;
    transform: translateY(-50%);
    z-index: 3;
    opacity: 0; /* Hidden range input, overlay on slider */
    cursor: ew-resize;
}

/* Visual slider handle (vertical line) */
.preview-slider::before {
    content: '';
    position: absolute;
    top: 0;
    left: 50%; /* Initially at 50% */
    width: 3px;
    height: 100%;
    background: white;
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.5);
    z-index: 4;
    pointer-events: none;
    transition: left 0.05s ease; /* Synchronized with clip-path */
}
```

```javascript
// web/js/preview.js - Slider interaction

function initSlider(sliderEl, afterImageEl) {
    const handle = sliderEl.querySelector('.preview-slider__handle');

    handle.addEventListener('input', (e) => {
        const value = e.target.value; // 0-100

        // Update clip-path (0% = all before, 100% = all after)
        const clipPercent = 100 - value; // Invert: slider 50 = clip 50% from right
        afterImageEl.style.clipPath = `inset(0 0 0 ${clipPercent}%)`;

        // Update visual handle position
        sliderEl.style.setProperty('--slider-position', `${value}%`);
    });

    // Keyboard accessibility
    handle.addEventListener('keydown', (e) => {
        if (e.key === 'ArrowLeft') {
            handle.value = Math.max(0, parseInt(handle.value) - 1);
            handle.dispatchEvent(new Event('input'));
        } else if (e.key === 'ArrowRight') {
            handle.value = Math.min(100, parseInt(handle.value) + 1);
            handle.dispatchEvent(new Event('input'));
        }
    });
}
```

### APIs and Interfaces

**preview.js Public API:**

```javascript
// web/js/preview.js

/**
 * Show preview modal for uploaded file
 * @param {UploadedFile} file - Uploaded file with UniversalRecipe
 * @param {Object} recipe - Parsed UniversalRecipe parameters
 */
export function showPreview(file, recipe) {
    // Create modal DOM
    const modal = createPreviewModal(file, recipe);
    document.body.appendChild(modal);

    // Initialize slider
    const sliderEl = modal.querySelector('.preview-slider');
    const afterImageEl = modal.querySelector('.preview-slider__after');
    initSlider(sliderEl, afterImageEl);

    // Apply CSS filters to "after" image
    const filterValue = recipeToCSS Filters(recipe);
    afterImageEl.style.filter = filterValue;

    // Tab switching
    const tabs = modal.querySelectorAll('.preview-tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const imageName = tab.dataset.image; // 'portrait', 'landscape', 'product'
            switchReferenceImage(modal, imageName);
        });
    });

    // Actions
    modal.querySelector('.btn--primary').addEventListener('click', () => {
        closePreview();
        // Trigger conversion (integrate with app.js)
        convertFile(file);
    });

    modal.querySelector('.btn--secondary').addEventListener('click', closePreview);
    modal.querySelector('.preview-modal__close').addEventListener('click', closePreview);

    // Keyboard: Esc to close
    document.addEventListener('keydown', handleEscapeKey);

    // Show modal (fade in animation)
    requestAnimationFrame(() => {
        modal.classList.add('preview-modal--visible');
    });
}

/**
 * Close preview modal
 */
export function closePreview() {
    const modal = document.getElementById('preview-modal');
    if (!modal) return;

    modal.classList.remove('preview-modal--visible');
    setTimeout(() => {
        modal.remove();
        document.removeEventListener('keydown', handleEscapeKey);
    }, 300); // Wait for fade-out animation
}

/**
 * Switch between reference images (Portrait, Landscape, Product)
 * @param {HTMLElement} modal - Modal element
 * @param {string} imageName - 'portrait', 'landscape', or 'product'
 */
function switchReferenceImage(modal, imageName) {
    const beforeImg = modal.querySelector('.preview-slider__before');
    const afterImg = modal.querySelector('.preview-slider__after');

    const imagePath = `images/preview-${imageName}.webp`;
    beforeImg.src = imagePath;
    afterImg.src = imagePath;

    // Update active tab
    modal.querySelectorAll('.preview-tab').forEach(tab => {
        tab.classList.toggle('preview-tab--active', tab.dataset.image === imageName);
    });
}

function handleEscapeKey(e) {
    if (e.key === 'Escape') {
        closePreview();
    }
}
```

**Integration with app.js (Epic 10):**

```javascript
// web/js/app.js - Add preview button to file cards

import { showPreview } from './preview.js';
import { convertFile } from './converter.js';

function createFileCard(file) {
    const card = document.createElement('div');
    card.className = 'file-card';

    card.innerHTML = `
        <div class="file-card__header">
            <span>${file.file.name}</span>
            <span class="format-badge format-badge--${file.format}">${file.format.toUpperCase()}</span>
        </div>
        <div class="file-card__actions">
            <button class="btn btn--preview" data-file-id="${file.id}">Preview</button>
            <button class="btn btn--convert" data-file-id="${file.id}">Convert</button>
        </div>
    `;

    // Preview button
    card.querySelector('.btn--preview').addEventListener('click', async () => {
        // Parse UniversalRecipe from file (WASM call)
        const recipe = await parseRecipe(file.data, file.format);
        showPreview(file, recipe);
    });

    return card;
}

async function parseRecipe(data, format) {
    // Call WASM parser to get UniversalRecipe JSON
    // (Existing converter has parse functionality)
    const recipeJSON = await parseFileToJSON(data, format);
    return JSON.parse(recipeJSON);
}
```

### Workflows and Sequencing

**Preview Workflow:**

1. User uploads preset file (NP3, XMP, lrtemplate, Capture One, DCP)
2. File card displayed with "Preview" button
3. User clicks "Preview" button
4. `app.js` calls WASM parser to extract UniversalRecipe parameters
5. `preview.js` `showPreview(file, recipe)` called:
   a. Create modal DOM with before/after slider
   b. Load default reference image (Portrait)
   c. Map UniversalRecipe → CSS filter string
   d. Apply CSS filter to "after" image (`style.filter = filterValue`)
   e. Initialize slider (range input updates clip-path)
   f. Display parameter summary ("Exposure +0.7 • Contrast +15")
6. User interacts:
   - Drag slider to compare before/after
   - Click tabs to switch reference images (Portrait, Landscape, Product)
   - Click "Convert Now" to proceed to conversion
   - Click "Cancel" or Esc to close preview
7. On "Convert Now": Close modal, trigger batch conversion with selected format

**Slider Interaction:**

1. User drags slider handle (range input 0-100)
2. Input event fires, value read (e.g., 75)
3. Calculate clip-path: `inset(0 0 0 ${100-75}%)`
4. Apply to "after" image: Shows 75% after, 25% before
5. Update visual handle position (white vertical line)
6. Smooth transition (CSS transition 0.05s)

**Tab Switching:**

1. User clicks "Landscape" tab
2. `switchReferenceImage('landscape')` called
3. Update `<img>` src attributes to `images/preview-landscape.webp`
4. CSS filter re-applied to "after" image (same filter value)
5. Active tab highlighted (CSS class `preview-tab--active`)

## Non-Functional Requirements

### Performance

**Targets:**
- Preview modal open: <100ms (instant feel)
- CSS filter application: <10ms (browser-native, GPU-accelerated)
- Slider interaction: 60fps (16ms per frame)
- Reference image load: <500ms on 3G (cached after first load)
- Tab switching: <50ms (image swap)

**Optimization Strategies:**
- Use CSS `will-change: filter` for GPU acceleration
- Preload all three reference images on page load (eager loading)
- Use WebP format for images (better compression than JPEG/PNG)
- Avoid JavaScript image processing (CSS filters only)
- Debounce slider input (optional, but CSS transitions handle this)

**Performance Monitoring:**
- Lighthouse audit: Interaction to Next Paint (INP) <200ms
- Manual testing: Drag slider smoothness on mobile devices
- Browser DevTools: Monitor filter application performance

### Security

**Threats & Mitigations:**
- **XSS via filename in modal title**: Sanitize filename with `textContent`
- **Malicious preset parameters**: Clamp CSS filter values (0-200% for brightness/contrast)
- **Resource exhaustion**: Limit reference image sizes (<200KB each)

**Privacy:**
- Reference images are static, bundled (no external requests)
- Preview is ephemeral (no preview data stored or uploaded)
- CSS filters applied client-side (no server processing)

**Input Validation:**
- Clamp CSS filter values to valid ranges (0-1000% max)
- Validate recipe parameters before mapping to CSS
- Handle missing/undefined recipe parameters gracefully (skip filter)

### Reliability/Availability

**Error Handling:**
- Missing reference image: Fallback to portrait (or show error message)
- Invalid recipe parameters: Show preview with no filters ("No adjustments")
- CSS filter not supported: Browser compatibility message (unlikely, modern browsers)

**Browser Compatibility:**
- Chrome 90+, Firefox 88+, Safari 14+, Edge 90+ (all support CSS filters)
- Graceful degradation: Older browsers show before image only (no filters)

**Failure Modes:**
- Modal doesn't open: Log error, allow conversion without preview
- Slider doesn't work: Show before/after images side-by-side (fallback)
- CSS filters don't apply: Show before image only

### Observability

**Logging:**
- Console log when preview modal opened (debug mode)
- Console warn if reference image fails to load
- Console error if CSS filter mapping fails

**Metrics:**
- Not required for personal project (no analytics)
- Could track preview usage (localStorage): "X previews opened"

**Debugging:**
- Verbose mode (`?debug=1`): Log CSS filter string, recipe parameters
- Browser DevTools: Inspect filter values, check GPU acceleration

## Dependencies and Integrations

**Dependencies:**

No external dependencies:
- Browser CSS filters (native support)
- HTML `<input type="range">` for slider
- CSS `clip-path` for before/after split

**Reference Images:**
- Create or acquire 3 public domain/CC0 images:
  - Portrait: Neutral skin tone, natural lighting
  - Landscape: Varied colors (sky, foliage, water)
  - Product: Still-life with textures, shadows
- Optimize with WebP: Target <200KB per image
- Embed in `web/images/` directory

**Integration Points:**
- Epic 10 upload interface: Add "Preview" button to file cards
- WASM converter: Parse preset → UniversalRecipe JSON
- Epic 10 modal styling: Consistent with overall UI design

**Version Constraints:**
- CSS Filters Level 1 (W3C Recommendation, all modern browsers)
- WebP image format (supported since Chrome 32, Firefox 65, Safari 14)

## Acceptance Criteria (Authoritative)

**AC-1: CSS Filter-Based Preview**
- ✅ Map UniversalRecipe parameters to CSS filter functions:
  - Exposure → brightness() (0% to 200%)
  - Contrast → contrast() (0% to 200%)
  - Saturation → saturate() (0% to 200%)
  - Hue adjustments → hue-rotate() (degrees)
  - Temperature/tint → sepia() + hue-rotate() approximation
- ✅ Preview renders in <100ms (instant, no processing delay)
- ✅ Preview updates in real-time as preset is selected (CSS filter applied instantly)
- ✅ Clear label: "Approximate preview using CSS filters"
- ✅ Works across all modern browsers (Chrome, Firefox, Safari, Edge)

**AC-2: Reference Image Bundle**
- ✅ Three reference images: Portrait, Landscape, Product/Still-life
- ✅ Images optimized for web (<200KB each, total <600KB)
- ✅ Images representative of common photography genres
- ✅ Images embedded in web bundle (no external requests, local `images/` directory)
- ✅ Licensing: Public domain or created specifically for Recipe (CC0 license)
- ✅ Images work well with typical preset adjustments (neutral starting point, no extreme colors)

**AC-3: Preview Modal Interface**
- ✅ Modal opens when "Preview" button clicked after upload
- ✅ Modal shows before/after slider with reference image
- ✅ Tabs to switch between reference images (Portrait, Landscape, Product)
- ✅ Preset parameters displayed: "Exposure +0.7 • Contrast +15 • Warmth +10"
- ✅ "Convert now" button in modal proceeds to conversion
- ✅ Close/cancel button returns to upload screen
- ✅ Modal keyboard accessible (Esc to close, Tab navigation, Arrow keys for slider)
- ✅ Mobile: Full-screen modal, touch-friendly controls (tap sides for snap, drag handle)

**AC-4: Accuracy Communication**
- ✅ Label: "Approximate preview using CSS filters" (prominent in modal)
- ✅ Tooltip/help text explains CSS filter limitations (hover/tap on label)
- ✅ No misleading claims about preview accuracy (avoid "realistic" or "accurate")
- ✅ Documentation explains preview vs. actual conversion differences (README section)
- ✅ Preview limitations listed (e.g., tone curves not supported in Phase 1)
- ✅ User expectations managed through transparent communication

**AC-5: Performance Optimization**
- ✅ Preview rendering completes in <100ms (CSS filter application instant)
- ✅ No blocking JavaScript during preview render (async/await, no long tasks)
- ✅ CSS filters applied via hardware acceleration (`will-change: filter`)
- ✅ Reference images cached after first load (browser cache, preload on page load)
- ✅ Smooth slider interaction (60fps minimum, CSS transitions)
- ✅ No performance degradation with multiple preview sessions (ephemeral modal, cleaned up on close)
- ✅ Works on mid-range mobile devices (tested on 3-year-old phones: iPhone 8, Android mid-range)

## Traceability Mapping

| AC ID | Spec Section(s) | Component(s)/API(s) | Test Idea |
| ----- | --------------- | ------------------- | --------- |
| AC-1 | Data Models (CSS filter mapping) | preview.js (recipeToCSS Filters) | Unit test: Recipe parameters → CSS filter string |
| AC-2 | Dependencies (Reference images) | images/preview-*.webp | Manual: Validate image size, licensing |
| AC-3 | APIs (showPreview, slider) | preview.js, preview.css | Manual: Open modal, interact with slider, tabs |
| AC-4 | Data Models (Disclaimer) | preview.html (modal template) | Manual: Check label visibility, documentation |
| AC-5 | NFRs (Performance) | preview.js, preview.css | Manual: Lighthouse audit, mobile device testing |

**Test Coverage Targets:**
- Unit tests: `preview.test.js` (CSS filter mapping, parameter formatting)
- Integration tests: Modal interaction (manual browser testing)
- Manual tests: Browser compatibility (Chrome, Firefox, Safari, Edge), Mobile devices (iPhone, Android)
- Coverage target: 75% for preview.js (lower due to DOM dependencies)

## Risks, Assumptions, Open Questions

### Risks

**RISK-1: CSS filter approximation too inaccurate**
- **Severity**: Medium
- **Impact**: Users may find preview misleading if accuracy is poor
- **Mitigation**: Clearly label as "Approximate", manage expectations via documentation
- **Owner**: Dev (Epic 11), Justin (user feedback)

**RISK-2: Reference images don't represent user's photos**
- **Severity**: Low
- **Impact**: Users may want to preview on their own images
- **Mitigation**: Document as Phase 2 feature (custom image upload), use diverse reference images
- **Owner**: Justin (product decision)

**RISK-3: Before/after slider janky on mobile**
- **Severity**: Low
- **Impact**: Poor user experience on touch devices
- **Mitigation**: Use CSS transitions, test on real devices, optimize clip-path performance
- **Owner**: Dev (Epic 11)

### Assumptions

**ASSUMPTION-1**: CSS filter approximation is "good enough" for MVP
- **Rationale**: Users want directional accuracy (warm/cool, bright/dark), not pixel-perfect
- **Validation**: User testing with real presets, gather feedback
- **Risk if false**: May need to prioritize Phase 2 (WASM-based preview) sooner

**ASSUMPTION-2**: Three reference images sufficient
- **Rationale**: Cover major photography genres (Portrait, Landscape, Product)
- **Validation**: User feedback on whether more genres needed
- **Risk if false**: May need to add more reference images (simple to add)

**ASSUMPTION-3**: Users understand "Approximate preview" disclaimer
- **Rationale**: Clear labeling and documentation set expectations
- **Validation**: Monitor user feedback for confusion
- **Risk if false**: May need more prominent warning or tooltip

**ASSUMPTION-4**: Before/after slider is best preview interaction
- **Rationale**: Slider is intuitive, widely used pattern
- **Validation**: User testing, compare to side-by-side view
- **Risk if false**: May need to add side-by-side option (easy to add)

### Open Questions

**Q-1**: Should users be able to upload their own images for preview?
- **Impact**: Significantly increases complexity (file upload, image decoding)
- **Resolution**: NO for Phase 1 - defer to Phase 2 (WebAssembly-based preview)
- **Owner**: Justin (product decision, confirmed)

**Q-2**: Should preview show parameter sliders (adjustable preview)?
- **Impact**: Would allow users to tweak parameters before conversion
- **Resolution**: NO for Phase 1 - preview is read-only, conversion happens after
- **Owner**: Justin (product decision, confirmed for MVP)

**Q-3**: Should preview support tone curves (not just basic adjustments)?
- **Impact**: CSS filters don't support arbitrary curves
- **Resolution**: NO for Phase 1 - document as limitation, requires Phase 2 WASM
- **Owner**: Dev (Epic 11), documented in parameter-mapping.md

**Q-4**: Should preview modal auto-open on file upload?
- **Impact**: User experience vs user control
- **Resolution**: NO - require explicit "Preview" button click (user-initiated)
- **Owner**: Dev (Epic 11)

## Test Strategy Summary

### Test Levels

**Unit Tests (preview.js):**
- `preview.test.js`: CSS filter mapping (recipeToCSS Filters), parameter formatting
- Coverage target: 75% for preview.js

**Integration Tests:**
- Modal interaction: Open preview, drag slider, switch tabs, close modal
- CSS filter application: Verify filter string applied to "after" image

**Manual Validation:**
- Browser compatibility: Chrome, Firefox, Safari, Edge (all latest versions)
- Mobile devices: iPhone (Safari), Android (Chrome), iPad (Safari)
- Performance: Lighthouse audit (INP <200ms), slider smoothness (60fps)
- Accuracy: Visual comparison of CSS filter preview vs actual conversion

### Test Frameworks

- Jest (or similar): JavaScript unit testing
- Lighthouse CI: Performance auditing
- BrowserStack: Cross-browser/device testing (manual)
- Manual testing: Real devices, visual accuracy assessment

### Coverage of ACs

| AC ID | Test Type | Test Location | Coverage |
| ----- | --------- | ------------- | -------- |
| AC-1 | Unit + Manual | preview.test.js, Browser testing | CSS filter mapping, browser support |
| AC-2 | Manual | Image validation | File sizes, licensing, visual quality |
| AC-3 | Manual | Browser testing | Modal interaction, keyboard nav, mobile |
| AC-4 | Manual | Visual inspection | Disclaimer visibility, documentation |
| AC-5 | Manual | Lighthouse, Device testing | Performance metrics, mobile smoothness |

### Edge Cases

**Edge Case Testing:**
- No adjustments (recipe all zeros): Preview shows no filter ("No adjustments")
- Maximum adjustments (exposure +2.0, contrast +1.0): Clamp CSS filters to valid range
- Negative adjustments (exposure -2.0): Test brightness 0% (completely black)
- Temperature/tint edge cases: Very warm (+100), very cool (-100)
- Reference image fails to load: Fallback to portrait or show error message
- Modal opened twice: Ensure only one modal instance (close previous)
- Slider at extremes (0% or 100%): Full before or full after
- Tab switching mid-slider drag: Handle gracefully (preserve slider position)
- Keyboard navigation: Tab through all interactive elements
- Mobile touch: Slider drag on touchscreen, tab taps

---

**Next Steps:**
1. Acquire or create 3 reference images (Portrait, Landscape, Product) - public domain/CC0
2. Optimize images with WebP: <200KB per image
3. Implement `web/js/preview.js` (CSS filter mapping, modal logic)
4. Implement `web/css/preview.css` (modal styles, slider, tabs)
5. Integrate preview button into file cards (Epic 10 app.js)
6. Write unit tests for CSS filter mapping
7. Manual testing: Browser compatibility, mobile devices, performance audit
8. Update documentation (README.md, FAQ) with preview limitations
9. Mark epic-11 as "contexted" in sprint-status.yaml
10. Celebrate! All 4 Path A epic tech specs complete 🎉
