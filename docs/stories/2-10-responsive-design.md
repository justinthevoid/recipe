# Story 2-10: Responsive Design

**Epic:** Epic 2 - Web Interface (FR-2)
**Story ID:** 2-10
**Status:** drafted
**Created:** 2025-11-04
**Complexity:** Medium (1-2 days)

---

## User Story

**As a** photographer using a mobile device or tablet
**I want** Recipe to work well on my screen size
**So that** I can convert presets on-the-go without needing my desktop computer

---

## Business Value

While Epic 2 is **desktop-first** (per PRD), mobile/tablet responsiveness expands Recipe's usability:

**Use cases:**
- **On-location:** Photographer at shoot wants to quickly convert preset
- **Tablet users:** iPad/Android tablet users (large screen, but touch input)
- **Mobile browsing:** User discovers Recipe on mobile, tries immediately

**Market data:**
- ~40% of web traffic is mobile (general web stats)
- ~20% of photography workflow tools used on mobile (estimate)

**This story ensures Recipe works everywhere, even if desktop remains primary.**

---

## Acceptance Criteria

### AC-1: Mobile Breakpoint (<768px) - Single Column Layout

- [ ] All content stacks vertically (single column)
- [ ] Drop zone:
  - Full width (100% minus padding)
  - Smaller text (reduce font sizes by ~20%)
  - Larger touch target (min 48×48px for buttons)
- [ ] Parameter display:
  - Single column (parameter name/value stack vertically)
  - Smaller font sizes
  - Collapsible by default (save vertical space)
- [ ] Format selector:
  - Full width options
  - Larger radio buttons (easier to tap)
  - Format descriptions remain visible
- [ ] Convert/Download buttons:
  - Full width
  - Larger tap targets (min 48px height)

**Test:**
1. Resize browser to 375px width (iPhone size)
2. Verify: All elements stack vertically
3. Verify: No horizontal scrolling
4. Verify: Buttons easy to tap (not too small)

### AC-2: Tablet Breakpoint (768-1023px) - Hybrid Layout

- [ ] Drop zone:
  - Centered, 90% width (max 600px)
  - Full-size text
- [ ] Parameter display:
  - Two-column grid (parameter name | value)
  - OR single column if clearer on tablet
- [ ] Format selector:
  - Centered, 90% width
  - Radio buttons larger than desktop (touch-friendly)
- [ ] Convert/Download buttons:
  - Centered, full width or fixed width (e.g., 400px)

**Test:**
1. Resize browser to 800px width (iPad size)
2. Verify: Layout adapts (not cramped like mobile, not spread like desktop)
3. Verify: Comfortable touch targets
4. Verify: Content readable without pinch-zoom

### AC-3: Desktop Breakpoint (≥1024px) - Full Layout (Already Implemented)

- [ ] Drop zone: 600px width, centered
- [ ] Parameter display: Two-column grid
- [ ] Format selector: 600px width, centered
- [ ] All hover states functional

**Test:**
1. Resize browser to 1920px width (desktop)
2. Verify: Layout matches Stories 2-1 through 2-9 specifications
3. Verify: No changes needed (already responsive from previous stories)

### AC-4: Touch-Friendly Interactions

- [ ] **No hover-dependent UX:** Essential actions don't require hover
  - Format selection works via tap (not just hover)
  - Buttons show active state on tap
  - Error details expand on tap (not hover)
- [ ] **Larger tap targets:**
  - Minimum 48×48px for all interactive elements (WCAG 2.1 AAA)
  - Spacing between tap targets (min 8px)
- [ ] **No drag-drop on mobile:**
  - Drag-drop disabled on touch devices (not reliable)
  - File picker opens on tap of drop zone
  - Show message: "Tap to select file" (not "Drag or click")

**Test:**
1. Test on touch device (iPad, Android tablet)
2. Verify: All buttons tappable (not too small)
3. Verify: No accidental taps (adequate spacing)
4. Tap drop zone → file picker opens (drag-drop doesn't interfere)

### AC-5: Mobile-Specific UI Adaptations

- [ ] **Drop zone text:**
  - Desktop: "Drop your preset file here or click to browse"
  - Mobile: "Tap to select your preset file"
- [ ] **Parameter panel:**
  - Desktop: Expanded by default
  - Mobile: Collapsed by default (save vertical space)
- [ ] **Error messages:**
  - Desktop: Fixed position at top
  - Mobile: Full-width at top, smaller dismiss button
- [ ] **Footer:**
  - Desktop: Horizontal links
  - Mobile: Vertical links (stacked)

**Test:**
1. Load on mobile device
2. Verify: Drop zone text says "Tap to select"
3. Verify: Parameter panel collapsed by default
4. Upload file → trigger error → verify error message readable
5. Scroll to footer → verify links stacked vertically

### AC-6: Performance on Mobile Devices

- [ ] WASM loads <5s on 3G connection
- [ ] Conversion completes <500ms on mobile (acceptable overhead)
- [ ] Page renders <3s on 3G (First Contentful Paint)
- [ ] No janky scrolling (60fps smooth scroll)

**Test:**
1. Test on real mobile device (or DevTools → throttle to "Slow 3G")
2. Measure WASM load time (console logs)
3. Verify: WASM loads <5s
4. Convert file → verify conversion <500ms
5. Scroll page → verify smooth (no jank)

### AC-7: Mobile Testing (iOS Safari, Android Chrome)

- [ ] **iOS Safari:**
  - File picker works (no iOS-specific bugs)
  - WASM loads and runs
  - Download works (saves to Files app)
- [ ] **Android Chrome:**
  - File picker works
  - WASM loads and runs
  - Download works (saves to Downloads folder)
- [ ] **Mobile-specific bugs:**
  - No viewport zoom issues
  - No keyboard overlap (input fields visible when keyboard open)
  - No broken layout on screen rotation

**Test:**
1. Test full conversion flow on iPhone (iOS Safari)
2. Test full conversion flow on Android device (Chrome)
3. Verify: File picker, conversion, download all work
4. Rotate device → verify layout adapts (no broken elements)
5. Verify: No pinch-zoom required to read text

### AC-8: Responsive Images and Icons

- [ ] All icons scale properly (use SVG or CSS, not fixed-size PNG)
- [ ] Format badges readable at all sizes
- [ ] Privacy badge text readable on mobile (may need smaller font)
- [ ] Error icons scale appropriately

**Test:**
1. Resize browser from 320px to 1920px
2. Verify: All icons/badges scale smoothly (no pixelation)
3. Verify: Text remains readable at all sizes

---

## Technical Approach

### Responsive CSS Strategy

**File:** `web/static/style.css` (update with responsive breakpoints)

**Mobile-first approach:**
- Base styles for mobile (<768px)
- Tablet styles in `@media (min-width: 768px)`
- Desktop styles in `@media (min-width: 1024px)`

```css
/* ========================================
   BASE STYLES (Mobile <768px)
   ======================================== */

body {
    font-size: 14px; /* Smaller base font for mobile */
    padding: 1rem; /* Tighter padding */
}

header h1 {
    font-size: 2rem; /* Smaller on mobile */
}

/* Drop zone */
.drop-zone {
    width: 100%;
    padding: 2rem 1rem; /* Less vertical padding */
}

.drop-zone .primary-text {
    font-size: 1rem;
}

.drop-zone .secondary-text {
    display: none; /* Hide "or click to browse" on mobile */
}

/* Parameter panel */
.parameter-panel {
    padding: 1rem;
}

.parameter-grid {
    grid-template-columns: 1fr; /* Single column */
}

.parameter-row {
    flex-direction: column; /* Stack name/value */
    align-items: flex-start;
}

/* Format selector */
.format-options {
    grid-template-columns: 1fr; /* Single column */
}

.format-option {
    padding: 1rem;
}

/* Buttons */
.convert-button,
.download-button {
    width: 100%;
    padding: 1rem; /* Larger tap target */
    font-size: 1rem;
}

/* Footer */
footer {
    font-size: 0.875rem;
}

.footer-links a {
    display: block; /* Stack vertically */
    margin-bottom: 0.5rem;
}

/* ========================================
   TABLET STYLES (768px - 1023px)
   ======================================== */

@media (min-width: 768px) {
    body {
        font-size: 16px; /* Standard size */
        padding: 1.5rem;
    }

    header h1 {
        font-size: 2.5rem;
    }

    /* Drop zone */
    .drop-zone {
        width: 90%;
        max-width: 600px;
        margin: 0 auto;
        padding: 2.5rem;
    }

    .drop-zone .secondary-text {
        display: block; /* Show on tablet */
    }

    /* Parameter panel */
    .parameter-grid {
        grid-template-columns: 1fr 1fr; /* Two columns */
    }

    .parameter-row {
        flex-direction: row; /* Horizontal */
        align-items: center;
    }

    /* Format selector */
    .format-options {
        grid-template-columns: 1fr; /* Still single column (easier to read) */
    }

    /* Buttons */
    .convert-button,
    .download-button {
        width: 100%;
        max-width: 600px;
        margin: 0 auto;
    }

    /* Footer */
    .footer-links a {
        display: inline-block; /* Horizontal */
        margin-right: 1rem;
        margin-bottom: 0;
    }
}

/* ========================================
   DESKTOP STYLES (≥1024px)
   ======================================== */

@media (min-width: 1024px) {
    body {
        padding: 2rem;
    }

    header h1 {
        font-size: 3rem;
    }

    /* Drop zone */
    .drop-zone {
        width: 600px;
        padding: 3rem;
    }

    /* Parameter panel */
    .parameter-panel {
        padding: 1.5rem;
    }

    /* Format selector */
    .format-options {
        grid-template-columns: 1fr; /* Single column (full width options clearer) */
    }

    /* Buttons */
    .convert-button,
    .download-button {
        width: 600px;
    }
}

/* ========================================
   TOUCH-SPECIFIC STYLES
   ======================================== */

/* Larger tap targets on touch devices */
@media (hover: none) and (pointer: coarse) {
    /* Touch device detected */

    .drop-zone {
        padding: 2.5rem 1.5rem;
    }

    .convert-button,
    .download-button,
    .error-action-btn,
    .format-option {
        min-height: 48px; /* WCAG AAA touch target */
    }

    /* Disable hover effects (not applicable on touch) */
    .drop-zone:hover {
        background: #f7fafc; /* Same as default */
        border-color: #cbd5e0;
    }
}
```

### Touch Detection and Adaptation

**File:** `web/static/responsive.js` (new file)

```javascript
// responsive.js - Responsive adaptations

/**
 * Detect if device is touch-enabled
 */
export function isTouchDevice() {
    return ('ontouchstart' in window) ||
           (navigator.maxTouchPoints > 0) ||
           (navigator.msMaxTouchPoints > 0);
}

/**
 * Adapt UI for touch devices
 */
export function adaptForTouch() {
    if (isTouchDevice()) {
        // Update drop zone text
        const dropZonePrimaryText = document.querySelector('.drop-zone .primary-text');
        if (dropZonePrimaryText) {
            dropZonePrimaryText.textContent = 'Tap to select your preset file';
        }

        // Disable drag-drop events (not reliable on touch)
        const dropZone = document.getElementById('dropZone');
        if (dropZone) {
            // Remove drag event listeners (if any)
            // Keep click event (works on touch)
        }

        // Collapse parameter panel by default on mobile
        if (window.innerWidth < 768) {
            const parameterPanel = document.querySelector('.parameter-panel');
            if (parameterPanel) {
                parameterPanel.classList.add('collapsed');
            }
        }

        console.log('Touch device detected - UI adapted');
    }
}

/**
 * Handle screen orientation change
 */
export function handleOrientationChange() {
    window.addEventListener('orientationchange', () => {
        console.log('Orientation changed - re-rendering UI');

        // Optional: Re-render components that need adjustment
        // Most CSS should handle this automatically
    });
}

/**
 * Initialize responsive adaptations
 */
export function initializeResponsive() {
    adaptForTouch();
    handleOrientationChange();
}
```

### Integration with Main Flow

**Update `main.js`:**

```javascript
// main.js - Initialize responsive adaptations

import { initializeResponsive } from './responsive.js';

// Initialize responsive adaptations
document.addEventListener('DOMContentLoaded', () => {
    initializeResponsive();
    // ... rest of initialization
});
```

---

## Dependencies

### Required Before Starting

- ✅ Stories 2-1 through 2-9 complete (all desktop UI implemented)

### No Blocking Dependencies

Story 2-10 is the final story in Epic 2.

---

## Testing Plan

### Manual Testing (Desktop Browser Resize)

**Test Case 1: Mobile Breakpoint (375px - iPhone)**
1. Resize browser to 375×667px (iPhone 6/7/8 size)
2. Verify: Drop zone text says "Tap to select"
3. Verify: All content single column (no horizontal scroll)
4. Verify: Buttons full width, easy to tap
5. Upload file → convert → download → verify full flow works

**Test Case 2: Tablet Breakpoint (768px - iPad)**
1. Resize browser to 768×1024px (iPad size)
2. Verify: Drop zone centered, ~90% width
3. Verify: Parameter panel shows two columns (or single if clearer)
4. Verify: Format selector comfortable to read
5. Complete full conversion flow

**Test Case 3: Desktop Breakpoint (1920px)**
1. Resize browser to 1920×1080px
2. Verify: Layout matches Stories 2-1 through 2-9 (no changes)

**Test Case 4: Extreme Sizes**
1. Resize to 320px (iPhone SE, smallest modern phone)
2. Verify: UI still usable (no broken layout)
3. Resize to 2560px (large desktop)
4. Verify: Content doesn't stretch awkwardly (max-width enforced)

### Real Device Testing

**iOS Safari (iPhone):**
1. Open Recipe on iPhone (iOS Safari)
2. Tap drop zone → file picker opens
3. Select `.xmp` file from Files app
4. Verify: File uploads, format detected, parameters displayed
5. Select target format → tap Convert
6. Tap Download → verify file saves to Files app
7. Open downloaded file in Lightroom → verify preset works

**Android Chrome:**
1. Open Recipe on Android device (Chrome)
2. Complete full conversion flow
3. Verify: Download saves to Downloads folder
4. Verify: File is valid

**iPad Safari:**
1. Test on iPad (tablet size)
2. Verify: Layout comfortable (not cramped like phone, not sparse like desktop)
3. Complete conversion flow

### Performance Testing (Mobile)

**3G Network Simulation:**
1. DevTools → Network → Throttle to "Slow 3G"
2. Load page → measure WASM load time
3. Verify: WASM loads <5s
4. Convert file → measure conversion time
5. Verify: Conversion <500ms

**Low-End Device:**
1. Test on older device (e.g., iPhone 6, Android phone from 2018)
2. Verify: Page loads and conversion works
3. Verify: No crashes or freezes

---

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Responsive breakpoints tested (320px, 768px, 1024px, 1920px)
- [ ] Touch-friendly interactions verified
- [ ] Mobile-specific UI adaptations implemented
- [ ] Real device testing completed (iOS Safari, Android Chrome, iPad)
- [ ] Performance tested on 3G and low-end device
- [ ] No horizontal scrolling at any breakpoint
- [ ] Code reviewed
- [ ] Story marked "ready-for-dev" in sprint status

---

## Out of Scope

**Explicitly NOT in this story:**
- ❌ Mobile app (native iOS/Android - Recipe is web-only)
- ❌ Progressive Web App (PWA) - install to home screen (future enhancement)
- ❌ Offline mode (Service Worker, IndexedDB - future enhancement)
- ❌ Mobile-specific features (camera integration, etc.)

**This story only delivers:** Responsive web design - Recipe works on mobile/tablet browsers.

---

## Technical Notes

### Why Mobile-First CSS?

**Approach:** Start with mobile styles, add desktop styles via `@media (min-width: ...)`

**Alternative:** Desktop-first (start with desktop, override for mobile)

**Decision:** Mobile-first (best practice)

**Rationale:**
- Easier to enhance (add features) than remove (simplify)
- Mobile styles tend to be simpler (fewer columns, smaller fonts)
- Performance: Mobile devices load fewer styles (don't parse desktop styles)

### Touch Detection

**Challenge:** Detect touch devices reliably

**Approach:** Multiple checks:
```javascript
('ontouchstart' in window) ||  // Standard
(navigator.maxTouchPoints > 0) || // Modern browsers
(navigator.msMaxTouchPoints > 0)  // IE/Edge legacy
```

**Alternative:** Assume touch if screen width <768px

**Decision:** Use touch detection (more accurate - tablets can be 1024px but touch)

### Drag-Drop on Touch Devices

**Problem:** Drag-drop events (`dragover`, `drop`) don't work reliably on touch devices

**Solution:** Disable drag-drop on touch, rely on file picker

**Implementation:** Keep click event on drop zone (opens file picker), ignore drag events on touch

**User impact:** Touch users just tap → file picker opens (simpler UX anyway)

### Mobile Performance

**WASM considerations:**
- WASM binary is 1.03MB (compressed) - loads fast even on 3G
- WASM compilation may be slower on low-end devices (~1-2s)
- Conversion performance acceptable (<500ms on mobile)

**Optimization opportunities (future):**
- Service Worker caching (WASM loads instantly on repeat visits)
- Progressive loading (show UI before WASM fully loaded)
- Lazy-load non-critical features

---

## Follow-Up Stories

**After Story 2-10:**
- Epic 2 complete! All 10 stories drafted and ready for implementation.

**Future mobile enhancements (not Epic 2):**
- Progressive Web App (PWA) - install to home screen
- Offline mode - convert without internet
- Service Worker caching - instant load on repeat visits
- Mobile-specific optimizations (reduce WASM size, faster conversion)

---

## References

- **Tech Spec:** `docs/tech-spec-epic-2.md` (Story 2-10 section)
- **PRD:** `docs/PRD.md` (FR-2.10: Responsive Design)
- **WCAG 2.1 Touch Targets:** https://www.w3.org/WAI/WCAG21/Understanding/target-size.html
- **Mobile-First CSS:** https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/Responsive/Mobile_first

---

**Story Created:** 2025-11-04
**Story Owner:** Justin (Developer)
**Reviewer:** Bob (Scrum Master)
**Estimated Effort:** 1-2 days
**Status:** Ready for SM approval → move to "ready-for-dev"
