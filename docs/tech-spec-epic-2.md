# Technical Specification: Epic 2 - Web Interface

**Epic:** FR-2 - Web Interface
**Status:** Ready for Implementation
**Dependencies:** Epic 1 (Core Conversion Engine) - COMPLETED ✅
**Target Completion:** TBD
**Story Count:** 10 stories (2-1 through 2-10)

---

## Executive Summary

Epic 2 transforms Recipe from a Go library into a browser-based photo preset converter using WebAssembly. This epic delivers the MVP web interface enabling photographers to convert presets without installing software, with 100% client-side processing for privacy.

**Key Achievement Goal:** Deploy a functional web UI where users can drag-drop a preset file and download the converted result in <2 seconds total user experience time.

---

## Architecture Overview

### System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Browser (Client-Side Only)              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Static HTML/CSS/JS (web/)                           │  │
│  │  - Drag-drop UI                                      │  │
│  │  - File upload handling                              │  │
│  │  - Parameter preview                                 │  │
│  │  - Download trigger                                  │  │
│  └──────────────────────────────────────────────────────┘  │
│                         ↓ ↑                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  JavaScript Bridge (wasm_exec.js)                    │  │
│  │  - File API integration                              │  │
│  │  - Uint8Array ↔ Go []byte conversion               │  │
│  │  - Promise wrapping                                  │  │
│  └──────────────────────────────────────────────────────┘  │
│                         ↓ ↑                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  WASM Module (recipe.wasm - 1.03MB compressed)       │  │
│  │  ┌──────────────────────────────────────────────┐   │  │
│  │  │  Epic 1 Conversion Engine (unchanged)        │   │  │
│  │  │  - UniversalRecipe hub                       │   │  │
│  │  │  - NP3/XMP/lrtemplate parsers/generators    │   │  │
│  │  │  - ConversionError handling                  │   │  │
│  │  └──────────────────────────────────────────────┘   │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
         ↓ Download (Blob)
    User's Filesystem
```

**Key Principle:** Zero server communication - all processing happens in the browser via WebAssembly.

### Component Responsibilities

**1. HTML/CSS/JS Layer (Stories 2-1, 2-2, 2-4, 2-5, 2-7, 2-8, 2-9, 2-10)**
- User interface and interaction
- File I/O via File API
- Visual feedback and error display
- No business logic (conversion handled by WASM)

**2. JavaScript Bridge (Story 2-2, 2-3, 2-6)**
- Marshals data between JS and WASM
- Handles async Promise patterns
- Converts File → Uint8Array → Go []byte
- Converts Go []byte → Uint8Array → Blob → Download

**3. WASM Module (Story 2-6)**
- Epic 1 conversion engine (no changes required)
- Exposes: `convert()`, `detectFormat()`, `getVersion()`
- Thread-safe, stateless
- Binary size: 1.03MB compressed (66% under 3MB target)

---

## Technical Decisions

### Decision 1: Frontend Framework

**Options Evaluated:**
- Vanilla JavaScript (zero dependencies)
- React (familiar, mature ecosystem)
- Svelte (compile-time, small bundle)

**Decision:** **Vanilla JavaScript**

**Rationale:**
- Epic 2 UI is simple (10 stories, minimal state)
- Zero framework overhead (target: <100KB total JS)
- Faster initial load time
- Easier to understand for future contributors
- Can refactor to framework in Epic 2+ if complexity grows

**Trade-offs:**
- No component reusability (acceptable for small UI)
- Manual DOM manipulation (jQuery-style)
- State management via closures/objects

### Decision 2: WASM Loading Strategy

**Options:**
- Eager loading (load WASM on page load)
- Lazy loading (load WASM on first conversion)
- Service Worker caching

**Decision:** **Eager loading + Service Worker caching**

**Rationale:**
- Users come to site to convert (not to browse)
- Loading WASM upfront sets expectation ("Loading...")
- Service Worker caches WASM for instant subsequent visits
- Simpler error handling (fail fast if WASM doesn't load)

**Implementation:**
```javascript
// On page load
const go = new Go();
WebAssembly.instantiateStreaming(fetch('static/recipe.wasm'), go.importObject)
    .then(result => {
        go.run(result.instance);
        // Status: "WASM loaded, ready to convert"
    });
```

### Decision 3: File Format Auto-Detection

**Options:**
- User selects source format manually
- Auto-detect via file extension
- Auto-detect via content inspection

**Decision:** **Auto-detect via content inspection (WASM)**

**Rationale:**
- Better UX (fewer clicks)
- Handles renamed files correctly
- Epic 1's `DetectFormat()` already proven (100% accurate)
- WASM overhead negligible (<1ms)

**User flow:**
1. User uploads file → `detectFormat(bytes)`
2. Display: "Detected format: XMP"
3. Pre-select target format (inverse of source)
4. User clicks "Convert"

### Decision 4: Error Handling Strategy

**Options:**
- Silent failures with console logs
- JavaScript alerts
- Inline error messages in UI

**Decision:** **Inline error messages + console logging**

**Rationale:**
- Alerts are intrusive and block UI
- Console logs help debugging but users don't see them
- Inline messages are user-friendly and contextual

**Error Categories:**
1. **WASM Load Failure** → Red banner: "Failed to load converter"
2. **Invalid File** → Red banner: "Not a valid preset file"
3. **Conversion Error** → Red banner with details: "Conversion failed: [message]"

### Decision 5: Privacy Validation

**Requirement:** Verify zero network requests during conversion

**Decision:** **Manual DevTools testing + documentation**

**Rationale:**
- Automated E2E tests (Playwright) deferred to Epic 6
- Manual testing sufficient for MVP
- Document in Story 2-9 acceptance criteria
- CI/CD network monitoring in future epic

**Validation process:**
1. Open browser DevTools → Network tab
2. Upload file, convert, download
3. Verify: No XHR/fetch requests (only WASM load)
4. Document screenshot in Story 2-9

### Decision 6: Browser Compatibility Target

**Support Matrix:**

| Browser | Version | Status |
|---------|---------|--------|
| Chrome | Latest 2 | ✅ Full support |
| Firefox | Latest 2 | ✅ Full support |
| Safari | Latest 2 | ✅ Full support |
| Edge | Latest 2 | ✅ Full support (Chromium) |
| Safari iOS | Latest 2 | ⚠️ Limited (mobile secondary) |
| Chrome Android | Latest 2 | ⚠️ Limited (mobile secondary) |
| IE11 | Any | ❌ Not supported (no WASM) |

**Required Browser Features:**
- WebAssembly (2017+)
- File API (FileReader, Blob, URL.createObjectURL)
- Drag-and-drop events
- Promises / async-await

**Fallback for unsupported browsers:**
```html
<div class="unsupported-browser" style="display: none;">
    Your browser doesn't support WebAssembly. Please use Chrome, Firefox, or Safari.
</div>
```

---

## WASM Integration Details

### JavaScript API

The WASM module exposes three global functions:

#### `convert(inputBytes, fromFormat, toFormat)`

**Signature:**
```typescript
convert(
    inputBytes: Uint8Array,
    fromFormat: string,  // "np3" | "xmp" | "lrtemplate" | ""
    toFormat: string     // "np3" | "xmp" | "lrtemplate"
): Promise<Uint8Array>
```

**Behavior:**
- **Async:** Returns Promise (does not block browser UI)
- **Auto-detection:** `fromFormat = ""` triggers auto-detect
- **Error handling:** Promise rejects with error message string
- **Thread-safe:** Multiple conversions can run concurrently

**Example:**
```javascript
const fileData = new Uint8Array(await file.arrayBuffer());
try {
    const outputData = await convert(fileData, "xmp", "np3");
    downloadBlob(new Blob([outputData]), "preset.np3");
} catch (err) {
    showError(`Conversion failed: ${err}`);
}
```

#### `detectFormat(inputBytes)`

**Signature:**
```typescript
detectFormat(inputBytes: Uint8Array): Promise<string>
// Returns: "np3" | "xmp" | "lrtemplate"
// Rejects if format unknown
```

**Detection Logic (from Epic 1):**
- **NP3:** Magic bytes "NCP" + min 300 bytes
- **XMP:** XML structure + "crs:" or "x:xmpmeta"
- **lrtemplate:** Lua syntax "s = {"

**Performance:** <1ms

#### `getVersion()`

**Signature:**
```typescript
getVersion(): string
// Returns: "1.0.0-wasm"
```

**Purpose:** Display in UI footer, debugging

### Data Flow

**Upload → Convert → Download:**

```
1. User drops file on page
   ↓
2. JavaScript FileReader reads file as ArrayBuffer
   ↓
3. Convert ArrayBuffer → Uint8Array
   ↓
4. Call: detectFormat(uint8array) → "xmp"
   ↓
5. User selects target: "np3"
   ↓
6. Call: convert(uint8array, "xmp", "np3") → Promise<Uint8Array>
   ↓
7. Convert Uint8Array → Blob
   ↓
8. Create download URL: URL.createObjectURL(blob)
   ↓
9. Trigger download via <a download="...">
   ↓
10. User's filesystem has converted file
```

**Key constraint:** All data stays in browser memory (no persistence).

### Performance Expectations

**From WASM Preparation Report:**

| Operation | Native Go (Epic 1) | WASM Target | Expected WASM |
|-----------|-------------------|-------------|---------------|
| XMP parse | 0.045ms | <100ms | 5-50ms |
| XMP generate | 0.0085ms | <100ms | 1-20ms |
| lrtemplate parse | 0.067ms | <100ms | 5-50ms |
| NP3 parse | ~0.05ms | <100ms | 5-50ms |
| Total conversion | <1ms | <100ms | 10-100ms |

**Performance Buffer:** Epic 1 is 200-3500x faster than targets, providing cushion for WASM overhead (typically 10-100x slower than native).

**User perception:** <100ms feels instant. Target: conversion completes before user's eyes leave the "Convert" button.

---

## Security & Privacy

### Privacy Architecture

**Zero Server Communication:**
- Static site (no backend)
- WASM runs in browser sandbox
- Files never uploaded
- No analytics or tracking

**Validation:**
- Browser DevTools → Network tab → verify no XHR/fetch
- Document in Story 2-9
- Screenshot showing "0 requests" during conversion

### Content Security Policy (CSP)

**Recommended headers (Cloudflare Pages config):**

```http
Content-Security-Policy:
    default-src 'self';
    script-src 'self' 'wasm-unsafe-eval';
    style-src 'self' 'unsafe-inline';
    img-src 'self' data:;
    connect-src 'none';
```

**Key restrictions:**
- `connect-src 'none'` - No fetch/XHR allowed
- `wasm-unsafe-eval` - Required for WebAssembly.instantiate
- `unsafe-inline` - For inline styles (minimize use)

**Implementation:** Story 2-9 (Privacy Messaging)

### WASM Sandbox Security

**Browser guarantees:**
- WASM cannot access DOM directly
- WASM cannot make network requests
- WASM cannot access filesystem (only File API)
- Memory isolation (buffer overflows contained)

**Attack surface:**
- Malicious preset files (Epic 1 parsers validated with 1,479 files)
- XSS via user-controlled content (sanitize parameter display)

---

## Story Breakdown

### Story 2-1: HTML Drag-Drop UI
**Goal:** Basic page structure with drag-drop zone

**Deliverables:**
- `web/index.html` (replace test version)
- `web/static/style.css`
- Drag-drop visual feedback
- File picker fallback

**Acceptance Criteria:**
- Drop zone highlights on hover
- File picker opens on click
- Accepts .np3, .xmp, .lrtemplate files
- Error message for unsupported files

**Complexity:** Simple (1-2 days)

---

### Story 2-2: File Upload Handling
**Goal:** Convert File object to Uint8Array for WASM

**Deliverables:**
- JavaScript module: `file-handler.js`
- FileReader integration
- File metadata extraction (name, size)

**Acceptance Criteria:**
- Reads file as ArrayBuffer
- Converts to Uint8Array
- Handles files up to 10MB
- Displays file name and size

**Complexity:** Simple (1 day)

---

### Story 2-3: Format Detection
**Goal:** Auto-detect preset format using WASM

**Deliverables:**
- Call `detectFormat()` on file upload
- Display detected format badge
- Error handling for unknown formats

**Acceptance Criteria:**
- Correctly detects NP3, XMP, lrtemplate
- Shows format badge in UI
- <100ms detection time
- Clear error for invalid files

**Complexity:** Simple (0.5-1 day)
**Dependency:** Story 2-2 (need Uint8Array first)

---

### Story 2-4: Parameter Preview Display
**Goal:** Show extracted parameters before conversion

**Technical Approach:**
1. Call `convert(bytes, format, "xmp")` → get UniversalRecipe as XMP
2. Parse XMP XML in JavaScript (lightweight parser)
3. Display key parameters (Exposure, Contrast, etc.)

**Alternative (simpler):**
- Parse original file format directly in JavaScript
- For XMP: use DOMParser
- For lrtemplate: regex parsing
- For NP3: skip preview (binary, not user-readable)

**Deliverables:**
- Parameter display component
- XMP/lrtemplate parsers (JS)

**Acceptance Criteria:**
- Shows 10-15 key parameters
- Updates on file upload
- Warns for unmappable parameters

**Complexity:** Medium (2-3 days)
**Dependency:** Story 2-2, 2-3

---

### Story 2-5: Target Format Selection
**Goal:** User selects output format

**Deliverables:**
- Format selection dropdown/radio buttons
- Smart default (XMP→NP3, NP3→XMP)
- Validation (prevent same-format conversion)

**Acceptance Criteria:**
- Shows all 3 formats as options
- Pre-selects logical default
- Updates when source format changes
- Error if source == target

**Complexity:** Simple (0.5-1 day)
**Dependency:** Story 2-3 (need detected format)

---

### Story 2-6: WASM Conversion Execution ⚠️ CRITICAL PATH
**Goal:** Perform conversion using WASM module

**Deliverables:**
- Call `convert(bytes, fromFormat, toFormat)`
- Progress indicator during conversion
- Error handling and retry

**Acceptance Criteria:**
- Conversion completes successfully
- Time <100ms (95th percentile)
- Handles errors gracefully
- Logs performance metrics

**Complexity:** Medium (2-3 days)
**Dependencies:** All previous stories
**Risk:** WASM performance issues (mitigation: Epic 1's 200-3500x buffer)

**Performance Testing:**
- Test with 100+ sample files
- Measure P50, P95, P99 conversion times
- Document in acceptance criteria

---

### Story 2-7: File Download Trigger
**Goal:** Download converted file automatically

**Deliverables:**
- Blob creation from Uint8Array
- Download trigger via <a download>
- Filename generation (preserve name, change extension)

**Acceptance Criteria:**
- Download starts automatically on success
- Filename: `original_name.{target_ext}`
- Works across Chrome, Firefox, Safari
- User can cancel/retry

**Complexity:** Simple (0.5-1 day)
**Dependency:** Story 2-6

---

### Story 2-8: Error Handling UI
**Goal:** User-friendly error messages

**Error Categories:**
1. WASM load failure
2. Invalid file format
3. Conversion errors (from ConversionError)
4. Browser compatibility

**Deliverables:**
- Error display component
- Error message mapping (technical → user-friendly)
- Retry/reset actions

**Acceptance Criteria:**
- No technical jargon in errors
- Actionable suggestions ("Try re-exporting from Lightroom")
- Errors are dismissible
- State resets cleanly

**Complexity:** Simple (1-2 days)
**Dependency:** Story 2-6

---

### Story 2-9: Privacy Messaging
**Goal:** Communicate privacy-first architecture

**Deliverables:**
- Privacy badge/banner ("Your files never leave your device")
- FAQ section
- Network monitoring validation (manual)

**Acceptance Criteria:**
- Privacy message visible on landing page
- Explains WebAssembly processing
- DevTools screenshot shows 0 requests during conversion
- Legal disclaimer (if needed)

**Complexity:** Simple (1 day)
**Dependency:** None (can be done in parallel)

---

### Story 2-10: Responsive Design
**Goal:** Works on desktop and tablet

**Breakpoints:**
- Desktop: ≥1024px (primary)
- Tablet: 768-1023px (secondary)
- Mobile: <768px (minimal support, not primary target)

**Deliverables:**
- Responsive CSS
- Touch-friendly controls (drag-drop, buttons)
- Tested on actual devices

**Acceptance Criteria:**
- Usable on iPad and similar tablets
- Maintains functionality on small screens
- Layout doesn't break at any viewport size

**Complexity:** Medium (2-3 days)
**Dependency:** Stories 2-1 through 2-9 (UI complete)

---

## File Structure

**Proposed structure:**

```
web/
├── index.html                  # Main page (Story 2-1)
├── static/
│   ├── style.css              # Styles (Story 2-1, 2-10)
│   ├── main.js                # Main application logic
│   ├── file-handler.js        # File I/O (Story 2-2)
│   ├── format-detector.js     # Format detection (Story 2-3)
│   ├── parameter-preview.js   # Parameter display (Story 2-4)
│   ├── converter.js           # WASM integration (Story 2-6)
│   ├── downloader.js          # Download handling (Story 2-7)
│   ├── error-handler.js       # Error display (Story 2-8)
│   ├── recipe.wasm            # Compiled WASM (1.03MB compressed)
│   └── wasm_exec.js           # Go WASM runtime (~16KB)
├── serve.py                    # Dev server (Python)
├── serve.js                    # Dev server (Node.js)
└── README.md                   # Testing documentation
```

**Code organization:**
- ES6 modules (`import`/`export`)
- No bundler for MVP (native browser support)
- Can add bundler (Webpack/Rollup) in future if needed

---

## Testing Strategy

### Story-Level Testing

Each story includes acceptance criteria with specific test cases:

**Example (Story 2-3: Format Detection):**
- ✅ Upload `Classic Chrome.np3` → detects "np3"
- ✅ Upload `preset.xmp` → detects "xmp"
- ✅ Upload `preset.lrtemplate` → detects "lrtemplate"
- ✅ Upload `image.jpg` → error "Unknown format"
- ✅ Detection completes in <100ms

### Integration Testing

**Epic-level test (all stories working together):**

1. Open browser: `http://localhost:8080`
2. Drag-drop: `examples/np3/Denis Zeqiri/Classic Chrome.np3`
3. Verify: "Detected format: NP3"
4. Verify: Parameter preview shows Sharpening, Contrast, etc.
5. Select target: XMP
6. Click "Convert"
7. Verify: Conversion completes <100ms
8. Verify: Download starts automatically
9. Verify: File named `Classic Chrome.xmp`
10. Open in Lightroom → parameters match original

### Browser Compatibility Testing

**Manual testing matrix:**

| Test Case | Chrome | Firefox | Safari |
|-----------|--------|---------|--------|
| File upload | ✅ | ✅ | ✅ |
| Drag-drop | ✅ | ✅ | ✅ |
| WASM loading | ✅ | ✅ | ✅ |
| Conversion | ✅ | ✅ | ✅ |
| Download | ✅ | ✅ | ✅ |

### Performance Testing

**Benchmark scenarios:**

1. **Small file:** <50KB (typical XMP) → <50ms target
2. **Medium file:** 50-500KB (typical lrtemplate) → <100ms target
3. **Large file:** >500KB (edge case) → <200ms acceptable

**Measurement:**
```javascript
const startTime = performance.now();
const outputData = await convert(inputData, from, to);
const elapsedTime = performance.now() - startTime;
console.log(`Conversion: ${elapsedTime.toFixed(2)}ms`);
```

---

## Deployment (Cloudflare Pages)

### Build Configuration

**Repository:** `github.com/user/recipe` (private initially)
**Build command:** `scripts/build-wasm.sh`
**Output directory:** `web/`
**Environment variables:** None needed

### Custom Headers

**`web/_headers` file:**
```
/*
  Content-Security-Policy: default-src 'self'; script-src 'self' 'wasm-unsafe-eval'; style-src 'self' 'unsafe-inline'; connect-src 'none'
  X-Frame-Options: DENY
  X-Content-Type-Options: nosniff
  Referrer-Policy: no-referrer
```

### Caching Strategy

**Static assets:**
```
/static/*
  Cache-Control: public, max-age=31536000, immutable
```

**HTML:**
```
/*.html
  Cache-Control: public, max-age=3600
```

**WASM binary:**
```
/static/recipe.wasm
  Content-Type: application/wasm
  Cache-Control: public, max-age=31536000, immutable
```

---

## Success Criteria

Epic 2 is complete when:

1. ✅ All 10 stories marked "done" in sprint status
2. ✅ User can upload any preset file (NP3/XMP/lrtemplate)
3. ✅ Format auto-detection works 100% accurately
4. ✅ Conversion completes <100ms (P95)
5. ✅ Converted file downloads automatically
6. ✅ Privacy messaging visible and accurate
7. ✅ Works in Chrome, Firefox, Safari (latest 2 versions)
8. ✅ Deployed to Cloudflare Pages and accessible via public URL
9. ✅ Zero network requests during conversion (verified via DevTools)
10. ✅ Round-trip testing passes (convert A→B, open in software, verify parameters match)

---

## Open Questions

### 1. Should we add batch conversion to MVP?

**PRD says:** Out of scope for MVP (deferred to Phase 2)

**Recommendation:** Defer to Epic 2+. Single-file conversion is sufficient for MVP validation.

### 2. Should we show visual preset preview?

**PRD says:** Out of scope for MVP

**Alternative:** Just show parameter list (Story 2-4)

**Recommendation:** Defer visual preview to Epic 2+ or Epic 4 (TUI has live preview).

### 3. Do we need a landing page (marketing)?

**PRD says:** Minimal landing page with project description (FR-7.1)

**Recommendation:** Simple hero section above converter UI (Story 2-1):
- 3-sentence description
- Privacy promise
- Immediate access to converter

No separate marketing page needed for MVP.

---

## References

- **PRD Epic 2:** `docs/PRD.md` (FR-2: Web Interface, lines 243-332)
- **Epic 1 Retrospective:** `docs/epic-1-retrospective.md`
- **WASM Preparation Report:** `docs/wasm-preparation-report.md`
- **Test Interface:** `web/index.html` (functional prototype)

---

**Tech Spec Completed:** 2025-11-04
**Author:** Bob (Scrum Master), Justin (Developer)
**Status:** Ready for story drafting
