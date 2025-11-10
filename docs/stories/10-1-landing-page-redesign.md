# Story 10.1: Landing Page Redesign with Hero Section

Status: ready-for-dev

## Story

As a **photographer visiting Recipe's website**,
I want **a modern, visually appealing landing page with clear value proposition and format badges**,
so that **I immediately understand Recipe's purpose, see supported formats at a glance, and feel confident using the tool**.

## Acceptance Criteria

**AC-1: Hero Section with Clear Value Proposition**
- ✅ Hero section at top of landing page with headline: "Convert Photo Presets. Instantly. Privately."
- ✅ Subheadline explaining Recipe's purpose (1-2 sentences max):
  - "Transform your photo editing presets between Nikon, Adobe, Lightroom, Capture One, and DNG Camera Profiles."
  - "All processing happens in your browser—your files never leave your device."
- ✅ Hero uses clean typography (system fonts: -apple-system, BlinkMacSystemFont, "Segoe UI")
- ✅ Visually distinct from main content (background color, padding, centered text)
- ✅ Call-to-action: "Get Started" button scrolls to upload section

**AC-2: Visual Format Badges**
- ✅ Format badges displayed prominently below hero section
- ✅ 5 badges total with defined brand colors:
  - NP3 (Nikon): #FFC107 (yellow)
  - XMP (Adobe): #0073E6 (blue)
  - lrtemplate (Lightroom): #D81B60 (magenta)
  - Capture One (.costyle): #9C27B0 (purple)
  - DCP (Camera Profile): #4CAF50 (green)
- ✅ Each badge shows:
  - Format abbreviation (e.g., "NP3", "XMP", "DCP")
  - Full name on hover (e.g., "Nikon Picture Control", "Adobe Camera Raw")
- ✅ Badges use consistent styling (rounded corners, bold text, white text on colored background)
- ✅ Accessible: Color not sole indicator (includes format name text)

**AC-3: Single-Page Layout**
- ✅ All core functionality on single page (no navigation to separate pages)
- ✅ Sections organized vertically:
  1. Hero section (value proposition)
  2. Format badges (supported formats)
  3. Upload section (drag-drop or file picker)
  4. Conversion section (format selection, download)
- ✅ Smooth scroll behavior between sections (CSS scroll-behavior: smooth)
- ✅ No page reloads during conversion workflow (SPA behavior with vanilla JS)

**AC-4: Responsive Design**
- ✅ Layout adapts to three breakpoints:
  - Mobile: 320px-767px (single column, stacked badges)
  - Tablet: 768px-1023px (2-column grid, wrapped badges)
  - Desktop: 1024px+ (centered max-width 1200px, badges in row)
- ✅ Hero section responsive:
  - Mobile: Headline font-size 28px, centered text
  - Tablet: Headline font-size 36px
  - Desktop: Headline font-size 48px
- ✅ Format badges responsive:
  - Mobile: Stack vertically (5 rows)
  - Tablet: 2-3 per row
  - Desktop: All 5 in single row
- ✅ Touch-friendly on mobile (button minimum 44x44px tap targets)

**AC-5: Performance - Fast Load Time**
- ✅ Page loads in <2 seconds on 3G connection (WebPageTest validation)
- ✅ Zero external dependencies:
  - No CDN fonts (use system fonts)
  - No analytics trackers (no Google Analytics, no Cloudflare Web Analytics)
  - No external CSS frameworks (vanilla CSS only)
- ✅ Optimized assets:
  - CSS minified and inline (or single stylesheet <10KB)
  - JavaScript deferred or async (non-blocking)
  - WASM loaded on-demand (only when user uploads file)
- ✅ Progressive enhancement: Basic HTML/CSS works without JavaScript

**AC-6: Clean Typography**
- ✅ System font stack for performance:
  ```css
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  ```
- ✅ Font sizes defined using CSS variables:
  - `--font-size-base: 16px` (body text)
  - `--font-size-large: 20px` (subheadings)
  - `--font-size-hero: 48px` (hero headline, responsive)
- ✅ Line height optimized for readability:
  - Body text: 1.6
  - Headings: 1.2
- ✅ Font weights:
  - Normal: 400
  - Bold: 600 (headings, badges)

**AC-7: No External Dependencies**
- ✅ Zero third-party scripts or stylesheets:
  - No Google Fonts (use system fonts)
  - No Bootstrap, Tailwind, or CSS frameworks
  - No jQuery or utility libraries
  - No analytics (Google Analytics, Plausible, Cloudflare Web Analytics)
- ✅ All assets self-hosted in `web/` directory:
  - `web/css/main.css` (global styles)
  - `web/css/components.css` (badges, buttons)
  - `web/css/layout.css` (responsive grid)
  - `web/js/app.js` (initialization)
- ✅ Privacy-focused: Zero tracking, zero external requests

## Tasks / Subtasks

### Task 1: Create Hero Section (AC-1)
- [ ] Update `web/index.html` to add hero section:
  ```html
  <section class="hero">
    <div class="hero__container">
      <h1 class="hero__title">Convert Photo Presets. Instantly. Privately.</h1>
      <p class="hero__subtitle">
        Transform your photo editing presets between Nikon, Adobe, Lightroom, Capture One, and DNG Camera Profiles.<br>
        All processing happens in your browser—your files never leave your device.
      </p>
      <a href="#upload" class="hero__cta button button--primary">Get Started</a>
    </div>
  </section>
  ```
- [ ] Add hero styles to `web/css/main.css`:
  ```css
  .hero {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 4rem 2rem;
    text-align: center;
  }

  .hero__title {
    font-size: var(--font-size-hero); /* 48px desktop */
    font-weight: 600;
    margin: 0 0 1rem 0;
    line-height: 1.2;
  }

  .hero__subtitle {
    font-size: var(--font-size-large); /* 20px */
    line-height: 1.6;
    max-width: 700px;
    margin: 0 auto 2rem auto;
    opacity: 0.95;
  }

  .hero__cta {
    display: inline-block;
    padding: 1rem 2rem;
    font-size: var(--font-size-large);
    text-decoration: none;
    border-radius: 8px;
  }
  ```
- [ ] Make hero responsive (mobile: 28px, tablet: 36px, desktop: 48px)
- [ ] Test "Get Started" button smooth scroll to #upload section

### Task 2: Create Format Badge System (AC-2)
- [ ] Define CSS variables for format colors in `web/css/main.css`:
  ```css
  :root {
    --color-np3: #FFC107;        /* Nikon yellow */
    --color-xmp: #0073E6;        /* Adobe blue */
    --color-lrtemplate: #D81B60; /* Magenta */
    --color-costyle: #9C27B0;    /* Capture One purple */
    --color-dcp: #4CAF50;        /* DCP green */
  }
  ```
- [ ] Create badge component styles in `web/css/components.css`:
  ```css
  .badge {
    display: inline-block;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-weight: 600;
    font-size: var(--font-size-small);
    color: white;
    text-align: center;
  }

  .badge--np3 { background-color: var(--color-np3); }
  .badge--xmp { background-color: var(--color-xmp); }
  .badge--lrtemplate { background-color: var(--color-lrtemplate); }
  .badge--costyle { background-color: var(--color-costyle); }
  .badge--dcp { background-color: var(--color-dcp); }

  .badge:hover {
    opacity: 0.9;
    transform: translateY(-2px);
    transition: all 0.2s ease;
  }
  ```
- [ ] Add format badges section to `web/index.html`:
  ```html
  <section class="formats">
    <div class="formats__container">
      <h2 class="formats__title">Supported Formats</h2>
      <div class="formats__badges">
        <span class="badge badge--np3" title="Nikon Picture Control">NP3</span>
        <span class="badge badge--xmp" title="Adobe Camera Raw">XMP</span>
        <span class="badge badge--lrtemplate" title="Lightroom Template">lrtemplate</span>
        <span class="badge badge--costyle" title="Capture One Style">.costyle</span>
        <span class="badge badge--dcp" title="DNG Camera Profile">DCP</span>
      </div>
    </div>
  </section>
  ```
- [ ] Make badges responsive (stack on mobile, wrap on tablet, row on desktop)

### Task 3: Implement Single-Page Layout (AC-3)
- [ ] Update `web/index.html` structure with sections:
  ```html
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Recipe - Convert Photo Presets</title>
    <link rel="stylesheet" href="css/main.css">
    <link rel="stylesheet" href="css/components.css">
    <link rel="stylesheet" href="css/layout.css">
  </head>
  <body>
    <section class="hero" id="hero"><!-- Hero content --></section>
    <section class="formats" id="formats"><!-- Format badges --></section>
    <section class="upload" id="upload"><!-- Upload UI --></section>
    <section class="conversion" id="conversion"><!-- Conversion controls --></section>

    <script src="js/app.js" defer></script>
  </body>
  </html>
  ```
- [ ] Add smooth scroll CSS to `web/css/main.css`:
  ```css
  html {
    scroll-behavior: smooth;
  }
  ```
- [ ] Test smooth scroll from hero CTA button to #upload section
- [ ] Verify no page reloads during workflow (all interactions via JavaScript)

### Task 4: Implement Responsive Design (AC-4)
- [ ] Create `web/css/layout.css` with breakpoints:
  ```css
  /* Mobile-first approach */
  .hero__container,
  .formats__container,
  .upload__container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1rem;
  }

  /* Mobile: 320px-767px */
  @media (max-width: 767px) {
    .hero__title {
      font-size: 28px;
    }

    .formats__badges {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .badge {
      width: 100%; /* Full width on mobile */
    }
  }

  /* Tablet: 768px-1023px */
  @media (min-width: 768px) and (max-width: 1023px) {
    .hero__title {
      font-size: 36px;
    }

    .formats__badges {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }
  }

  /* Desktop: 1024px+ */
  @media (min-width: 1024px) {
    .hero__title {
      font-size: 48px;
    }

    .formats__badges {
      display: flex;
      justify-content: center;
      gap: 1rem;
    }
  }
  ```
- [ ] Test responsive breakpoints in browser DevTools (320px, 768px, 1024px, 1920px)
- [ ] Verify touch targets ≥44x44px on mobile (buttons, badges)

### Task 5: Optimize for Performance (AC-5)
- [ ] Measure baseline performance with WebPageTest:
  - Test URL: https://recipe.justins.studio (or local build)
  - Connection: 3G (1.6 Mbps, 300ms RTT)
  - Metric: Load time (target: <2 seconds)
- [ ] Remove external dependencies:
  - Check index.html for external <link> or <script> tags
  - Remove any CDN references (Google Fonts, analytics)
  - Verify all assets are self-hosted in web/ directory
- [ ] Optimize CSS:
  - Minify main.css, components.css, layout.css (optional for dev)
  - Consider inlining critical CSS in <head> (for <10KB total)
  - Use CSS variables for performance (browser-native)
- [ ] Optimize JavaScript:
  - Add `defer` attribute to <script> tags (non-blocking)
  - Load WASM only when user uploads file (lazy loading)
- [ ] Progressive enhancement:
  - Verify basic HTML/CSS renders without JavaScript enabled
  - Core message (hero, format badges) visible without JS
- [ ] Document performance results in story completion notes

### Task 6: Implement Clean Typography (AC-6)
- [ ] Define system font stack in `web/css/main.css`:
  ```css
  :root {
    --font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    --font-size-base: 16px;
    --font-size-small: 14px;
    --font-size-large: 20px;
    --font-size-hero: 48px; /* Responsive via media queries */
    --font-weight-normal: 400;
    --font-weight-bold: 600;
    --line-height-body: 1.6;
    --line-height-heading: 1.2;
  }

  body {
    font-family: var(--font-family);
    font-size: var(--font-size-base);
    font-weight: var(--font-weight-normal);
    line-height: var(--line-height-body);
  }

  h1, h2, h3 {
    font-weight: var(--font-weight-bold);
    line-height: var(--line-height-heading);
  }
  ```
- [ ] Apply font sizes consistently across components
- [ ] Test font rendering on Windows, macOS, Linux (system fonts should adapt)

### Task 7: Remove All External Dependencies (AC-7)
- [ ] Audit `web/index.html` for external resources:
  - Check <link> tags (stylesheets, fonts)
  - Check <script> tags (analytics, libraries)
  - Check <img> tags (external images)
- [ ] Remove any third-party scripts:
  - No Google Analytics
  - No Cloudflare Web Analytics
  - No Plausible, Fathom, or other trackers
- [ ] Verify all assets are self-hosted:
  - CSS files in web/css/
  - JavaScript files in web/js/
  - WASM file in web/
- [ ] Add Content Security Policy (CSP) header (optional):
  ```html
  <meta http-equiv="Content-Security-Policy" content="default-src 'self'; script-src 'self' 'wasm-unsafe-eval';">
  ```
- [ ] Test landing page works completely offline (after initial load)

### Task 8: Manual Testing and Validation
- [ ] Test landing page on multiple browsers:
  - Chrome 120+ (desktop, mobile)
  - Firefox 120+ (desktop, mobile)
  - Safari 17+ (desktop, iOS)
  - Edge 120+ (desktop)
- [ ] Test responsive design at breakpoints:
  - 320px (iPhone SE)
  - 375px (iPhone 12)
  - 768px (iPad portrait)
  - 1024px (iPad landscape)
  - 1920px (desktop)
- [ ] Test touch interactions on mobile:
  - Tap hero CTA button (smooth scroll to upload)
  - Tap format badges (hover effect visible)
  - All tap targets ≥44x44px
- [ ] Test performance on 3G:
  - Use Chrome DevTools throttling (Slow 3G preset)
  - Verify load time <2 seconds
  - Check Network waterfall for external requests (should be zero)
- [ ] Test accessibility:
  - Keyboard navigation (Tab through all interactive elements)
  - Screen reader (VoiceOver on macOS, NVDA on Windows)
  - Color contrast (badges pass WCAG AA minimum 4.5:1)
- [ ] Test without JavaScript:
  - Disable JavaScript in browser
  - Verify hero and format badges visible
  - Upload UI should gracefully degrade (show message)
- [ ] Screenshot results for validation (desktop, tablet, mobile views)

## Dev Notes

### Learnings from Previous Story

**From Story 9-4-dcp-compatibility-validation (Status: drafted)**

Previous story (Epic 9) not yet implemented. Epic 10 represents a shift from backend format support (DCP) to frontend UI/UX enhancements.

**Key Differences:**
- Epic 9: Backend Go code (parsers, generators, TIFF/XML)
- Epic 10: Frontend HTML/CSS/JS (landing page, responsive design)
- Epic 9: Manual testing in Adobe software
- Epic 10: Browser testing (Chrome, Firefox, Safari, responsive)

[Source: docs/stories/9-4-dcp-compatibility-validation.md]

### Architecture Alignment

**Tech Spec Epic 10 Alignment:**

Story 10.1 implements **AC-1 (Redesigned Landing Page)** and **AC-2 (Visual Format Badges)** from tech-spec-epic-10.md.

**Component Structure:**
```
web/
├── index.html           # Landing page structure (MODIFIED)
├── css/
│   ├── main.css        # Global styles, CSS variables (NEW)
│   ├── components.css  # Badge, button components (NEW)
│   └── layout.css      # Responsive grid, breakpoints (NEW)
├── js/
│   └── app.js          # Initialization (MODIFIED)
└── recipe.wasm         # Existing WASM engine (UNCHANGED)
```

**Hero Section Design:**
```
┌────────────────────────────────────┐
│     HERO SECTION (gradient bg)     │
│                                     │
│  Convert Photo Presets.            │
│  Instantly. Privately.             │
│                                     │
│  Transform your photo editing...   │
│  All processing happens...         │
│                                     │
│     [Get Started Button]           │
└────────────────────────────────────┘
```

**Format Badges (Desktop):**
```
┌─────┐ ┌─────┐ ┌──────────┐ ┌────────┐ ┌─────┐
│ NP3 │ │ XMP │ │lrtemplate│ │.costyle│ │ DCP │
└─────┘ └─────┘ └──────────┘ └────────┘ └─────┘
 Yellow   Blue     Magenta      Purple    Green
```

[Source: docs/tech-spec-epic-10.md#Detailed-Design]

### CSS Architecture

**File Organization:**

1. **main.css** - Global styles, CSS variables, reset
   - CSS variables (colors, fonts, spacing)
   - Typography (body, headings)
   - Base element styles (html, body, a, button)

2. **components.css** - Reusable UI components
   - Badge system (.badge, .badge--np3, .badge--xmp, etc.)
   - Button styles (.button, .button--primary, .button--secondary)
   - Card components (for future upload cards)

3. **layout.css** - Responsive grid system
   - Container widths (max-width: 1200px)
   - Breakpoints (mobile: <768px, tablet: 768-1023px, desktop: 1024px+)
   - Section spacing (hero, formats, upload, conversion)

**CSS Methodology:**

- **BEM naming**: `.block__element--modifier` (e.g., `.hero__title`, `.badge--np3`)
- **Mobile-first**: Base styles for mobile, media queries for larger screens
- **CSS variables**: Use `var(--color-np3)` instead of hardcoded colors

[Source: docs/tech-spec-epic-10.md#Data-Models-and-Contracts]

### Responsive Design Strategy

**Breakpoints:**

1. **Mobile (320px-767px):**
   - Single column layout
   - Hero title: 28px
   - Format badges: Stack vertically (5 rows, full width)
   - Touch targets: ≥44x44px

2. **Tablet (768px-1023px):**
   - 2-column grid
   - Hero title: 36px
   - Format badges: 2-3 per row (wrapped)

3. **Desktop (1024px+):**
   - Centered max-width 1200px
   - Hero title: 48px
   - Format badges: All 5 in single row

**Mobile-First Approach:**

Start with mobile styles (no media query), then progressively enhance:
```css
/* Mobile (default) */
.formats__badges {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

/* Tablet */
@media (min-width: 768px) {
  .formats__badges {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
  }
}

/* Desktop */
@media (min-width: 1024px) {
  .formats__badges {
    display: flex;
    flex-direction: row;
  }
}
```

[Source: docs/tech-spec-epic-10.md#Acceptance-Criteria]

### Performance Optimization

**3G Load Time Target: <2 seconds**

**Optimization Strategies:**

1. **Zero External Requests:**
   - No CDN fonts (system fonts only)
   - No analytics trackers
   - No external CSS/JS libraries

2. **Asset Optimization:**
   - CSS: 3 files, total <20KB (unminified), <10KB (minified)
   - JS: Defer loading, async where possible
   - WASM: Load on-demand (only when user uploads file)

3. **Critical CSS Inlining:**
   - Inline hero + badges CSS in <head> (optional)
   - Defer non-critical CSS (upload, conversion sections)

4. **Progressive Enhancement:**
   - HTML/CSS works without JavaScript (hero, badges visible)
   - JavaScript enhances UX (smooth scroll, drag-drop, conversion)

**WebPageTest Validation:**

Test URL: https://recipe.justins.studio
- Connection: 3G (1.6 Mbps, 300ms RTT)
- Location: Dulles, Virginia
- Metrics:
  - First Contentful Paint (FCP): <1.5s
  - Largest Contentful Paint (LCP): <2.0s
  - Total Blocking Time (TBT): <300ms

[Source: docs/tech-spec-epic-10.md#System-Architecture-Alignment]

### Format Badge Colors (Brand Alignment)

**Rationale for Color Choices:**

| Format     | Color   | Hex     | Brand Alignment                         |
| ---------- | ------- | ------- | --------------------------------------- |
| NP3        | Yellow  | #FFC107 | Nikon's brand color (yellow accents)    |
| XMP        | Blue    | #0073E6 | Adobe's brand color (Adobe blue)        |
| lrtemplate | Magenta | #D81B60 | Lightroom's accent color (vibrant pink) |
| .costyle   | Purple  | #9C27B0 | Capture One's brand color (purple)      |
| DCP        | Green   | #4CAF50 | Generic/universal (no specific brand)   |

**Accessibility:**

All badge colors pass WCAG AA contrast ratio (4.5:1 minimum) when paired with white text:
- Yellow #FFC107 + white: 1.9:1 (FAIL - use dark text instead)
- Blue #0073E6 + white: 4.6:1 (PASS)
- Magenta #D81B60 + white: 5.2:1 (PASS)
- Purple #9C27B0 + white: 6.7:1 (PASS)
- Green #4CAF50 + white: 3.4:1 (FAIL - use dark text instead)

**Fix for Yellow/Green Badges:**
```css
.badge--np3,
.badge--dcp {
  color: #212121; /* Dark text for yellow/green backgrounds */
}
```

[Source: docs/tech-spec-epic-10.md#Data-Models-and-Contracts]

### Project Structure Notes

**New Files Created (Story 10.1):**
```
web/
├── css/
│   ├── main.css        # Global styles, CSS variables (NEW)
│   ├── components.css  # Badge, button components (NEW)
│   └── layout.css      # Responsive grid, breakpoints (NEW)
```

**Modified Files:**
- `web/index.html` - Add hero section, format badges section
- `web/js/app.js` - Initialize smooth scroll, WASM lazy loading

**Files from Previous Stories (Reused):**
- `web/recipe.wasm` - Existing WASM engine (no changes)
- `web/js/converter.js` - WASM conversion interface (no changes)

[Source: docs/tech-spec-epic-10.md#System-Architecture-Alignment]

### Testing Strategy

**Browser Compatibility Testing:**

| Browser      | Desktop | Mobile  | Notes                             |
| ------------ | ------- | ------- | --------------------------------- |
| Chrome 120+  | ✅       | ✅       | Primary browser, best DevTools    |
| Firefox 120+ | ✅       | ✅       | Test flexbox/grid compatibility   |
| Safari 17+   | ✅       | ✅ (iOS) | Test system font rendering        |
| Edge 120+    | ✅       | ❌       | Chromium-based, similar to Chrome |

**Responsive Testing:**

| Device    | Width  | Hero Font | Badge Layout   |
| --------- | ------ | --------- | -------------- |
| iPhone SE | 320px  | 28px      | Vertical stack |
| iPhone 12 | 375px  | 28px      | Vertical stack |
| iPad Mini | 768px  | 36px      | 2-3 per row    |
| iPad Pro  | 1024px | 48px      | Single row     |
| Desktop   | 1920px | 48px      | Single row     |

**Performance Testing:**

1. **WebPageTest** (3G connection):
   - First Contentful Paint: <1.5s
   - Largest Contentful Paint: <2.0s
   - Total Blocking Time: <300ms

2. **Chrome DevTools** (Slow 3G):
   - Network waterfall: Zero external requests
   - Coverage: Unused CSS/JS <10%

3. **Lighthouse** (Mobile):
   - Performance: ≥90
   - Accessibility: ≥95
   - Best Practices: 100
   - SEO: ≥90

**Accessibility Testing:**

- Keyboard navigation (Tab through all elements)
- Screen reader (VoiceOver, NVDA)
- Color contrast (WCAG AA 4.5:1 minimum)
- Touch targets (44x44px minimum on mobile)

[Source: docs/tech-spec-epic-10.md#Test-Strategy-Summary]

### Known Risks

**RISK-26: System fonts render inconsistently across platforms**
- **Impact**: Typography looks different on Windows vs. macOS vs. Linux
- **Mitigation**: Test on all platforms, use font-weight 400/600 only (widely supported)
- **Acceptable**: Minor rendering differences expected (system fonts are platform-specific)

**RISK-27: Performance exceeds 2-second target on 3G**
- **Impact**: Landing page feels slow on mobile networks
- **Mitigation**: Inline critical CSS, defer JavaScript, lazy-load WASM
- **Fallback**: Accept 2.5s load time if all optimizations applied

**RISK-28: Badge colors fail WCAG contrast**
- **Impact**: Accessibility issues for low-vision users
- **Mitigation**: Use dark text on yellow/green badges (instead of white)
- **Target**: All badges pass WCAG AA (4.5:1 minimum contrast)

[Source: docs/tech-spec-epic-10.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-10.md#Acceptance-Criteria] - AC-1: Redesigned Landing Page, AC-2: Visual Format Badges
- [Source: docs/tech-spec-epic-10.md#Data-Models-and-Contracts] - CSS variables, badge colors
- [Source: docs/tech-spec-epic-10.md#System-Architecture-Alignment] - Component structure, constraints
- [Source: web/index.html] - Current landing page (to be modified)
- [Source: web/js/app.js] - Current app initialization (to be enhanced)
- [Source: Story 2-1-html-drag-drop-ui.md] - Original web UI implementation (reference for existing patterns)
- [WebPageTest](https://www.webpagetest.org/) - Performance testing tool
- [WCAG 2.1 Contrast Checker](https://webaim.org/resources/contrastchecker/) - Accessibility validation

## Dev Agent Record

### Context Reference

- Context File: docs/stories/10-1-landing-page-redesign.context.xml (Generated: 2025-11-09)

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
