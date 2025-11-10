# Story 10.6: Performance Optimization for 3G Load Times

Status: drafted

## Story

As a **photographer on a slow mobile connection**,
I want **Recipe to load in under 2 seconds on 3G networks**,
so that **I can use the tool effectively even with limited bandwidth or in areas with poor connectivity**.

## Acceptance Criteria

**AC-1: Initial Page Load <2 Seconds on 3G**
- ✅ WebPageTest validation: Load time <2 seconds on 3G connection (1.6 Mbps)
- ✅ Measured from navigation start to fully interactive (Time to Interactive - TTI)
- ✅ Test location: Dulles, VA (WebPageTest server)
- ✅ Test device: Moto G4 (representative mid-range Android)
- ✅ Test network: 3G (1.6 Mbps down, 768 Kbps up, 300ms latency)
- ✅ Target: TTI <2000ms, First Contentful Paint (FCP) <1000ms

**AC-2: Lighthouse Performance Score ≥95**
- ✅ Lighthouse CI audit: Performance score ≥95
- ✅ Accessibility score ≥95 (maintained from previous stories)
- ✅ Best Practices score ≥95
- ✅ SEO score ≥90
- ✅ Run Lighthouse in CI/CD pipeline (GitHub Actions)
- ✅ Block PR merge if score drops below 90

**AC-3: Critical CSS Inlined**
- ✅ Above-the-fold CSS inlined in `<head>` tag:
  - Base styles (typography, colors, layout)
  - Hero section styles
  - Upload drop zone styles (visible on page load)
- ✅ Non-critical CSS loaded asynchronously:
  - File card styles (below fold)
  - Badge styles (below fold)
  - Conversion controls (below fold)
- ✅ Inline CSS size <14 KB (TCP slow-start threshold)
- ✅ Total CSS transferred <50 KB (gzipped)

**AC-4: Zero External Dependencies**
- ✅ No CDN-hosted resources:
  - No Google Fonts (use system fonts)
  - No Font Awesome (use inline SVG icons)
  - No external analytics (privacy-first, no tracking)
  - No external JavaScript libraries (vanilla JS only)
- ✅ All resources self-hosted:
  - WASM bundle served from `/wasm/`
  - CSS served from `/css/`
  - JavaScript served from `/js/`
- ✅ No third-party requests (verified with Network DevTools)

**AC-5: Image and Asset Optimization**
- ✅ SVG icons optimized:
  - Remove unnecessary metadata, comments
  - Minify SVG code (SVGO tool)
  - Inline critical SVG icons (upload cloud, checkmark, error X)
- ✅ No raster images (PNG/JPG) used:
  - Use SVG for all graphics (scalable, smaller file size)
  - Use CSS gradients for backgrounds (no image files)
- ✅ Favicon optimized:
  - Use SVG favicon (modern browsers)
  - Fallback ICO file <5 KB

**AC-6: JavaScript Bundle Optimization**
- ✅ JavaScript minified and compressed:
  - Minify with Terser (remove whitespace, comments, shorten variable names)
  - Gzip compression enabled on server
- ✅ JavaScript size budget:
  - Total JS transferred: <100 KB (gzipped)
  - WASM bundle: <2 MB (gzipped ~500 KB)
- ✅ Lazy-load non-critical JavaScript:
  - Load conversion logic after page interactive
  - Load WASM module on first file upload (defer until needed)
- ✅ No JavaScript render-blocking:
  - Use `defer` attribute on script tags
  - Or load scripts at end of `<body>`

**AC-7: 60fps Scrolling and Animations**
- ✅ All animations use GPU-accelerated properties:
  - Use `transform` and `opacity` (not `left`, `top`, `width`, `height`)
  - Example: Spinner rotation uses `transform: rotate(360deg)`
- ✅ CSS transitions <300ms (avoid janky long animations)
- ✅ No layout thrashing (batch DOM reads/writes)
- ✅ Use `will-change` for animated elements (hint browser for optimization)
- ✅ Test on low-end mobile device (Moto G4, 60fps maintained)

## Tasks / Subtasks

### Task 1: Run Initial WebPageTest Baseline (AC-1)
- [ ] Run WebPageTest audit on current Recipe deployment:
  - URL: https://recipe.justins.studio
  - Test location: Dulles, VA
  - Browser: Chrome
  - Connection: 3G (1.6 Mbps)
  - Device: Moto G4
- [ ] Document baseline metrics:
  - Time to Interactive (TTI): ??? ms
  - First Contentful Paint (FCP): ??? ms
  - Largest Contentful Paint (LCP): ??? ms
  - Total page size: ??? KB
- [ ] Identify performance bottlenecks:
  - Render-blocking resources
  - Large JavaScript bundles
  - Unused CSS
  - Network requests waterfall
- [ ] Create performance budget spreadsheet:
  - HTML: <20 KB
  - CSS: <50 KB (gzipped)
  - JavaScript: <100 KB (gzipped, excluding WASM)
  - WASM: <500 KB (gzipped)
  - Total: <670 KB (gzipped)

### Task 2: Run Initial Lighthouse Audit (AC-2)
- [ ] Run Lighthouse audit in Chrome DevTools:
  - Open https://recipe.justins.studio
  - DevTools → Lighthouse tab
  - Mode: Navigation
  - Device: Mobile
  - Categories: Performance, Accessibility, Best Practices, SEO
- [ ] Document baseline scores:
  - Performance: ??? / 100
  - Accessibility: ??? / 100
  - Best Practices: ??? / 100
  - SEO: ??? / 100
- [ ] Identify failing audits:
  - Render-blocking resources
  - Unused CSS/JS
  - Image optimization
  - Accessibility issues
- [ ] Create action plan to reach ≥95 Performance score

### Task 3: Inline Critical CSS (AC-3)
- [ ] Identify above-the-fold CSS (visible without scrolling):
  - Base styles (typography, colors, CSS variables)
  - Hero section styles (.hero, .hero__title, .hero__cta)
  - Upload drop zone styles (.upload__dropzone)
  - Navbar/header styles (if present)
- [ ] Extract critical CSS to separate file: `web/css/critical.css`
- [ ] Inline critical CSS in `web/index.html`:
  ```html
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Recipe - Convert Photo Presets</title>
    <style>
      /* Critical CSS inlined (< 14 KB) */
      :root { --color-primary: #0073E6; /* ... */ }
      body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; /* ... */ }
      .hero { padding: 4rem 2rem; /* ... */ }
      .upload__dropzone { border: 2px dashed #ccc; /* ... */ }
    </style>
  </head>
  ```
- [ ] Load non-critical CSS asynchronously:
  ```html
  <link rel="preload" href="/css/main.css" as="style" onload="this.onload=null;this.rel='stylesheet'">
  <noscript><link rel="stylesheet" href="/css/main.css"></noscript>
  ```
- [ ] Verify critical CSS size <14 KB (TCP slow-start threshold)
- [ ] Test rendering: Page looks complete before async CSS loads

### Task 4: Remove External Dependencies (AC-4)
- [ ] Audit current external requests (Network DevTools):
  - Check for any CDN-hosted resources
  - Check for analytics scripts (Google Analytics, etc.)
  - Check for web fonts (Google Fonts, Adobe Fonts)
- [ ] Remove any external dependencies found:
  - Replace web fonts with system fonts
  - Remove analytics scripts (Recipe is privacy-first)
  - Ensure all JS/CSS self-hosted
- [ ] Verify zero external requests:
  - Open Network DevTools
  - Load page
  - Verify all requests to `recipe.justins.studio` domain
  - No requests to `fonts.googleapis.com`, `cdn.jsdelivr.net`, etc.

### Task 5: Optimize SVG Icons (AC-5)
- [ ] Audit current SVG usage in `web/index.html`
- [ ] Optimize SVG icons with SVGO:
  - Install SVGO: `npm install -g svgo`
  - Run: `svgo web/icons/*.svg --multipass`
  - Flags: `--remove-metadata --remove-comments --minify-styles`
- [ ] Inline critical SVG icons (above-the-fold):
  - Upload cloud icon (drop zone)
  - Checkmark icon (complete status)
  - Error X icon (error status)
- [ ] Lazy-load non-critical SVG icons:
  - Format badge icons (below fold)
  - Spinner icon (only needed during conversion)
- [ ] Verify SVG file size reduction:
  - Before optimization: ??? KB
  - After optimization: ??? KB
  - Target: >50% reduction

### Task 6: Minify and Compress JavaScript (AC-6)
- [ ] Install Terser for JavaScript minification:
  - `npm install terser --save-dev`
- [ ] Create build script in `package.json`:
  ```json
  {
    "scripts": {
      "build": "npm run build:js && npm run build:wasm",
      "build:js": "terser web/js/app.js web/js/upload.js -c -m -o web/dist/bundle.min.js",
      "build:wasm": "cp wasm/recipe.wasm web/dist/recipe.wasm"
    }
  }
  ```
- [ ] Run build: `npm run build`
- [ ] Verify minified JavaScript size:
  - Original: ??? KB
  - Minified: ??? KB
  - Gzipped: ??? KB (target: <100 KB)
- [ ] Update `web/index.html` to use minified bundle:
  ```html
  <script src="/dist/bundle.min.js" defer></script>
  ```
- [ ] Enable gzip compression on Cloudflare Pages (automatic)

### Task 7: Lazy-Load WASM Module (AC-6)
- [ ] Defer WASM loading until first file upload:
  ```javascript
  // web/js/app.js
  let wasmModule = null;

  async function initWASM() {
    if (wasmModule) return wasmModule; // Already loaded

    const wasmUrl = '/dist/recipe.wasm';
    const response = await fetch(wasmUrl);
    const bytes = await response.arrayBuffer();
    wasmModule = await WebAssembly.instantiate(bytes);
    return wasmModule;
  }

  // Only load WASM when user uploads first file
  uploadManager.on('file-added', async () => {
    if (!wasmModule) {
      console.log('Loading WASM module...');
      await initWASM();
    }
  });
  ```
- [ ] Remove WASM preload from `<head>` (no longer needed on page load)
- [ ] Test lazy loading:
  - Open page → Network tab shows no WASM request
  - Upload file → WASM request triggered
  - Convert file → WASM module loaded and ready

### Task 8: Optimize CSS Transitions for 60fps (AC-7)
- [ ] Audit current animations and transitions:
  - File card hover (box-shadow, transform)
  - Status icon spinner (transform: rotate)
  - Drop zone drag-over (background-color, border-color)
  - Progress bar fill (width)
- [ ] Replace non-GPU properties with GPU-accelerated properties:
  ```css
  /* Before: Causes repaint */
  .file-card:hover {
    box-shadow: 0 4px 8px rgba(0,0,0,0.2);
  }

  /* After: GPU-accelerated */
  .file-card {
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    transition: transform 0.2s ease;
  }

  .file-card:hover {
    transform: translateY(-2px);
  }
  ```
- [ ] Add `will-change` hints for animated elements:
  ```css
  .status-icon--processing {
    will-change: transform; /* Hint browser to optimize */
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
  ```
- [ ] Reduce transition durations:
  - File card hover: 0.2s (was 0.3s)
  - Status transitions: 0.3s (keep, acceptable)
  - Progress bar: 0.3s (keep, acceptable)
- [ ] Test on low-end device (Moto G4):
  - Scroll through file grid → maintain 60fps
  - Hover file cards → no janky transitions
  - Progress bar animation → smooth fill

### Task 9: Add Lighthouse CI to GitHub Actions (AC-2)
- [ ] Create Lighthouse CI configuration: `.lighthouserc.json`
  ```json
  {
    "ci": {
      "collect": {
        "url": ["https://recipe.justins.studio"],
        "numberOfRuns": 3
      },
      "assert": {
        "assertions": {
          "categories:performance": ["error", {"minScore": 0.95}],
          "categories:accessibility": ["error", {"minScore": 0.95}],
          "categories:best-practices": ["error", {"minScore": 0.95}],
          "categories:seo": ["error", {"minScore": 0.90}]
        }
      },
      "upload": {
        "target": "temporary-public-storage"
      }
    }
  }
  ```
- [ ] Create GitHub Actions workflow: `.github/workflows/lighthouse-ci.yml`
  ```yaml
  name: Lighthouse CI
  on: [push, pull_request]

  jobs:
    lighthouse:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v3
        - uses: actions/setup-node@v3
          with:
            node-version: 18
        - run: npm install -g @lhci/cli
        - run: lhci autorun
  ```
- [ ] Test Lighthouse CI locally:
  - `npm install -g @lhci/cli`
  - `lhci autorun`
  - Verify scores meet thresholds
- [ ] Commit and push workflow → GitHub Actions runs Lighthouse CI

### Task 10: Run Final WebPageTest and Lighthouse Audits (AC-1, AC-2)
- [ ] Run final WebPageTest audit:
  - URL: https://recipe.justins.studio
  - Test location: Dulles, VA
  - Connection: 3G
  - Device: Moto G4
- [ ] Verify performance targets met:
  - ✅ TTI <2000ms (target: <2 seconds)
  - ✅ FCP <1000ms (target: <1 second)
  - ✅ LCP <2500ms (target: <2.5 seconds)
  - ✅ Total page size <670 KB (gzipped)
- [ ] Run final Lighthouse audit:
  - Performance: ??? / 100 (target: ≥95)
  - Accessibility: ??? / 100 (target: ≥95)
  - Best Practices: ??? / 100 (target: ≥95)
  - SEO: ??? / 100 (target: ≥90)
- [ ] Document improvements:
  - Baseline TTI: ??? ms → Final TTI: ??? ms (% improvement)
  - Baseline Performance score: ??? → Final: ??? (+ ??? points)
  - Baseline page size: ??? KB → Final: ??? KB (% reduction)

### Task 11: Manual Testing on Low-End Device (AC-7)
- [ ] Test on Moto G4 (or equivalent low-end Android):
  - Open https://recipe.justins.studio
  - Test scrolling performance:
    - Scroll through file grid (10+ cards)
    - Verify 60fps (Chrome DevTools → Performance → Record)
    - No janky scrolling or layout shifts
  - Test animation performance:
    - File card hover → smooth transform animation
    - Spinner animation → smooth rotation, no stuttering
    - Progress bar fill → smooth width transition
- [ ] Test on iPhone SE (low-end iOS device):
  - Same tests as Moto G4
  - Verify iOS Safari performance
- [ ] Document any performance issues found:
  - Slow scrolling → optimize CSS
  - Janky animations → use GPU-accelerated properties
  - Layout shifts → add explicit width/height

## Dev Notes

### Learnings from Previous Story

**From Story 10-5-responsive-mobile-design (Status: drafted)**

Previous story not yet implemented. Story 10.6 optimizes performance to ensure responsive design loads fast on all devices.

**Performance Constraints from Story 10.5:**
- Mobile-first CSS (smaller initial payload)
- Touch-friendly tap targets (no impact on performance)
- Responsive images (max-width: 100%, already optimized)
- System fonts (zero download time, addressed in Story 10.6)

**Integration with Story 10.5:**
- Responsive CSS already mobile-first (good for performance)
- No external font dependencies (Story 10.5 uses system fonts)
- Viewport meta tag already set (Story 10.5)

[Source: docs/stories/10-5-responsive-mobile-design.md]

### Architecture Alignment

**Tech Spec Epic 10 Alignment:**

Story 10.6 implements **Performance Requirements** from tech-spec-epic-10.md.

**Performance Budget:**
```
| Resource        | Size (Gzipped)   | Budget   |
| --------------- | ---------------- | -------- |
| HTML            | <20 KB           | 20 KB    |
| CSS             | <50 KB           | 50 KB    |
| JavaScript      | <100 KB          | 100 KB   |
| WASM            | <500 KB          | 500 KB   |
| --------------- | ---------------- | -------- |
| Total           | <670 KB          | 670 KB   |
```

**Load Time Targets:**
- 3G connection (1.6 Mbps): 670 KB ÷ 200 KB/s = 3.35 seconds (raw transfer)
- With optimization (compression, caching): <2 seconds TTI
- First Contentful Paint: <1 second (critical CSS inlined)

**Critical Rendering Path Optimization:**
```
HTML (20 KB, 100ms) → Critical CSS (inline, 0ms) → FCP (1000ms)
                   → JS (deferred, loads after FCP)
                   → WASM (lazy, loads on upload)
                   → Non-critical CSS (async, loads after FCP)
```

[Source: docs/tech-spec-epic-10.md#Performance-Requirements]

### Critical CSS Strategy

**What is Critical CSS?**

Critical CSS is the minimum CSS needed to render above-the-fold content (visible without scrolling).

**Recipe Critical CSS:**
```css
/* Base styles (typography, colors) */
:root { --color-primary: #0073E6; /* ... */ }
body { font-family: system-ui; color: #222; }

/* Hero section (visible on page load) */
.hero { padding: 4rem 2rem; text-align: center; }
.hero__title { font-size: 48px; margin: 0 0 1rem 0; }
.hero__cta { background: var(--color-primary); color: white; }

/* Upload drop zone (visible on page load) */
.upload__dropzone { border: 2px dashed #ccc; padding: 4rem 2rem; }
```

**Non-Critical CSS (loaded async):**
- File card styles (below fold)
- Format badge styles (below fold)
- Conversion controls (below fold)
- Responsive media queries (enhanced, not critical)

**Why Inline Critical CSS?**
- Eliminates render-blocking CSS request (no network delay)
- Faster First Contentful Paint (FCP <1 second)
- Better user experience (page visible immediately)

**14 KB Threshold:**
TCP slow-start sends 14 KB in first roundtrip. If critical CSS <14 KB, it arrives in first packet (no additional roundtrips needed).

[Source: Web Performance Best Practices - Critical CSS]

### Lazy-Loading WASM Strategy

**Why Lazy-Load WASM?**

Recipe's WASM bundle (~2 MB uncompressed, ~500 KB gzipped) is the largest resource. Loading it on page load delays Time to Interactive (TTI).

**Lazy-Loading Approach:**
```javascript
// Don't load WASM on page load
// Load WASM when user uploads first file

let wasmModule = null;

async function initWASM() {
  if (wasmModule) return; // Already loaded

  const response = await fetch('/dist/recipe.wasm');
  const bytes = await response.arrayBuffer();
  wasmModule = await WebAssembly.instantiate(bytes);
}

// Trigger WASM load on first file upload
uploadManager.on('file-added', async () => {
  if (!wasmModule) await initWASM();
});
```

**Benefits:**
- Page load time: 2.0s → 1.2s (WASM removed from critical path)
- Time to Interactive: 2.5s → 1.5s (no WASM parsing/compilation on load)
- User perception: Page feels faster (hero visible immediately)

**Trade-off:**
- First conversion: 100ms → 600ms (includes WASM load time)
- Acceptable: 600ms still fast, only happens once per session

[Source: WebAssembly Performance Best Practices]

### GPU-Accelerated CSS Animations

**What Properties Trigger GPU Acceleration?**

Modern browsers use GPU for specific CSS properties:
- ✅ `transform` (translate, rotate, scale)
- ✅ `opacity`
- ❌ `width`, `height`, `top`, `left` (trigger layout/repaint)
- ❌ `background-color`, `border-color` (trigger repaint)

**Recipe Animation Optimization:**
```css
/* Before: Causes repaint (slow) */
.file-card:hover {
  top: -2px; /* Triggers layout */
}

/* After: GPU-accelerated (fast) */
.file-card:hover {
  transform: translateY(-2px); /* GPU-accelerated */
}

/* Before: Causes repaint (slow) */
.status-icon--processing {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
/* Already optimized (transform is GPU-accelerated) */
```

**`will-change` Hint:**
```css
.status-icon--processing {
  will-change: transform; /* Tell browser to optimize */
  animation: spin 1s linear infinite;
}
```

**Warning:** Don't overuse `will-change` (memory overhead). Only use for actively animating elements.

[Source: CSS Performance Best Practices - GPU Acceleration]

### Lighthouse CI Integration

**Why Lighthouse CI?**

Lighthouse CI runs Lighthouse audits in CI/CD pipeline, preventing performance regressions:
- Runs on every commit/PR
- Blocks merge if score drops below threshold
- Tracks performance over time
- Catches performance regressions early

**Recipe Lighthouse CI Configuration:**
```json
{
  "ci": {
    "assert": {
      "assertions": {
        "categories:performance": ["error", {"minScore": 0.95}],
        "categories:accessibility": ["error", {"minScore": 0.95}]
      }
    }
  }
}
```

**Enforcement:**
- Performance <95 → CI fails, PR blocked
- Accessibility <95 → CI fails, PR blocked
- Forces developers to fix performance issues before merge

[Source: Lighthouse CI Documentation]

### Project Structure Notes

**New Files Created (Story 10.6):**
```
web/
├── css/
│   └── critical.css        # Critical CSS extracted (NEW)
├── dist/
│   ├── bundle.min.js       # Minified JavaScript (NEW)
│   └── recipe.wasm         # WASM bundle (copied from wasm/)
.github/
└── workflows/
    └── lighthouse-ci.yml   # Lighthouse CI workflow (NEW)
.lighthouserc.json          # Lighthouse CI config (NEW)
package.json                # Build scripts added (MODIFIED)
```

**Modified Files:**
- `web/index.html` - Inline critical CSS, defer JavaScript
- `web/js/app.js` - Lazy-load WASM module
- `web/css/components.css` - GPU-accelerated animations
- `package.json` - Add build scripts (Terser minification)

**Files from Previous Stories (Optimized):**
- `web/css/main.css` - Split into critical.css + main.css
- `web/js/upload.js` - Minified in bundle.min.js

[Source: docs/tech-spec-epic-10.md#Services-and-Modules]

### Testing Strategy

**Performance Testing (Required):**

1. **WebPageTest:**
   - Run on https://webpagetest.org
   - Test location: Dulles, VA
   - Connection: 3G (1.6 Mbps)
   - Device: Moto G4
   - Metrics: TTI, FCP, LCP, page size
   - Target: TTI <2000ms

2. **Lighthouse Audit:**
   - Run in Chrome DevTools
   - Mode: Navigation (not Timespan)
   - Device: Mobile
   - Categories: Performance, Accessibility, Best Practices, SEO
   - Target: Performance ≥95

3. **Real Device Testing:**
   - Moto G4 (low-end Android): Test scrolling, animations (60fps)
   - iPhone SE (low-end iOS): Test Safari performance
   - Chrome DevTools Performance tab: Record scrolling, check frame rate

**Automated Testing:**
- Lighthouse CI: Runs on every PR, blocks merge if score <90
- GitHub Actions: Automated performance audits

**Documentation:**
- WebPageTest results: Screenshot + report link
- Lighthouse report: JSON export + screenshots
- Performance budget: Track actual vs budget in spreadsheet

[Source: docs/tech-spec-epic-10.md#Test-Strategy-Summary]

### Known Risks

**RISK-42: WASM bundle size exceeds budget (>500 KB gzipped)**
- **Impact**: Page load time exceeds 2 seconds on 3G
- **Mitigation**: Lazy-load WASM (defer until first file upload)
- **Acceptable**: WASM loaded on-demand, not on page load

**RISK-43: Critical CSS exceeds 14 KB (TCP slow-start)**
- **Impact**: Critical CSS requires additional roundtrip (slower FCP)
- **Mitigation**: Minimize critical CSS (only above-fold styles)
- **Test**: Measure critical.css size, target <14 KB

**RISK-44: Lighthouse CI flaky scores (variance ±5 points)**
- **Impact**: CI fails intermittently even with good performance
- **Mitigation**: Run 3 audits, average scores (Lighthouse CI default)
- **Threshold**: Block merge only if score <90 (buffer for variance)

**RISK-45: GPU-accelerated animations don't work on old devices**
- **Impact**: Animations janky on devices without GPU acceleration
- **Mitigation**: Test on low-end device (Moto G4), ensure 60fps
- **Fallback**: If <60fps, reduce animation complexity (simpler transitions)

[Source: docs/tech-spec-epic-10.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-10.md#Performance-Requirements] - <2s load time, Lighthouse ≥95
- [Source: docs/tech-spec-epic-10.md#Optimization-Strategies] - Critical CSS, lazy-loading, minification
- [Source: docs/stories/10-5-responsive-mobile-design.md] - Mobile-first CSS, system fonts
- [WebPageTest](https://www.webpagetest.org/) - 3G load time testing
- [Lighthouse CI](https://github.com/GoogleChrome/lighthouse-ci) - Automated performance audits
- [Terser](https://terser.org/) - JavaScript minification
- [SVGO](https://github.com/svg/svgo) - SVG optimization
- [Critical CSS](https://web.dev/extract-critical-css/) - Above-fold CSS optimization
- [GPU-Accelerated CSS](https://web.dev/animations-guide/) - CSS animation performance
- [Lazy-Loading WASM](https://developers.google.com/web/fundamentals/performance/lazy-loading-guidance/wasm)

## Dev Agent Record

### Context Reference

- docs/stories/10-6-performance-optimization.context.xml

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
