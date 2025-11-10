# Epic Technical Specification: Enhanced Web UI/UX

Date: 2025-11-08
Author: Justin
Epic ID: epic-10
Status: Draft

---

## Overview

Epic 10 redesigns Recipe's web interface to provide a modern, polished user experience that matches the quality of the conversion engine. The enhancement focuses on visual appeal (format badges, clean typography), usability (batch drag-drop uploads, progress indicators), and accessibility (responsive design, keyboard navigation). This epic transforms Recipe from a functional tool into a delightful, professional-grade application.

The implementation uses Recipe's existing vanilla JavaScript + WebAssembly approach (no frameworks), extending the current `web/` directory with organized CSS modules (`main.css`, `components.css`, `layout.css`), refactored JavaScript modules (`upload.js`, `preview.js`), and responsive design patterns. All conversions remain client-side using the existing WASM engine.

## Objectives and Scope

**In Scope:**
- Redesigned landing page with hero section and clear value proposition
- Visual format badge system (NP3, XMP, lrtemplate, Capture One, DCP with brand colors)
- Batch file upload via drag-drop and file picker
- Progress indicators for multi-file conversions (overall progress, per-file status)
- Mobile-responsive design (320px+, 768px+, 1024px+)
- Before/after comparison slider (integrated with Epic 11 preview)
- Improved conversion flow (instant format detection, streamlined UI)

**Out of Scope (Path A):**
- JavaScript framework adoption (React, Vue, Svelte) - remain vanilla
- Backend API (all processing stays client-side WASM)
- User accounts, preferences, or persistent state (localStorage only)
- Advanced animations or transitions (keep performant on low-end devices)
- Internationalization (English only for MVP)
- Dark mode (nice-to-have, defer to future)

## System Architecture Alignment

**Components:**
- **Enhanced**: `web/index.html` - Redesigned landing page with hero section
- **New**: `web/css/main.css` - Global styles, CSS variables, color palette
- **New**: `web/css/components.css` - Reusable badge, button, card components
- **New**: `web/css/layout.css` - Responsive grid system, breakpoints
- **New**: `web/css/preview.css` - Preview modal and before/after slider (Epic 11 integration)
- **Enhanced**: `web/js/app.js` - Main application logic, initialization
- **New**: `web/js/upload.js` - Batch upload, drag-drop, file handling
- **Existing**: `web/js/converter.js` - WASM conversion interface (no changes needed)
- **Enhanced**: `web/js/utils.js` - Shared utilities, format detection

**Integration Points:**
- WASM engine via existing `recipe.wasm` and `converter.js`
- Format badges for all formats (NP3, XMP, lrtemplate, Capture One, DCP)
- Preview system (Epic 11) integrated via modal overlay
- Cloudflare Pages deployment (existing infrastructure)

**Constraints:**
- Must maintain <2 seconds load time on 3G connection
- Zero external dependencies (no CDN fonts, no analytics)
- Must work without JavaScript for basic file download (progressive enhancement)
- Must maintain existing WASM conversion performance (<100ms per file)

## Detailed Design

### Services and Modules

| Module | Responsibility | Inputs | Outputs | Owner |
| ------ | -------------- | ------ | ------- | ----- |
| `index.html` | Landing page structure, hero section | - | HTML DOM | Dev (Epic 10) |
| `css/main.css` | Global styles, CSS variables, typography | - | Styled document | Dev (Epic 10) |
| `css/components.css` | Badge, button, card components | - | Reusable components | Dev (Epic 10) |
| `css/layout.css` | Responsive grid, breakpoints | - | Adaptive layout | Dev (Epic 10) |
| `js/upload.js` | Batch upload, drag-drop handling | File[] | Uploaded file cards | Dev (Epic 10) |
| `js/app.js` | Application initialization, routing | - | Initialized app | Dev (Epic 10) |
| `js/utils.js` | Format detection, utilities | File bytes | Format string | Dev (Epic 10) |

### Data Models and Contracts

**CSS Variables (main.css):**

```css
/* web/css/main.css */
:root {
    /* Format Badge Colors */
    --color-np3: #FFC107;        /* Nikon yellow */
    --color-xmp: #0073E6;        /* Adobe blue */
    --color-lrtemplate: #D81B60; /* Magenta */
    --color-costyle: #9C27B0;    /* Capture One purple */
    --color-dcp: #4CAF50;        /* DCP green */

    /* UI Colors */
    --color-primary: #2196F3;    /* Primary blue */
    --color-success: #4CAF50;    /* Success green */
    --color-error: #F44336;      /* Error red */
    --color-warning: #FF9800;    /* Warning orange */

    /* Typography */
    --font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
    --font-size-base: 16px;
    --font-size-large: 20px;
    --font-size-small: 14px;
    --font-weight-normal: 400;
    --font-weight-bold: 600;

    /* Spacing */
    --spacing-xs: 4px;
    --spacing-sm: 8px;
    --spacing-md: 16px;
    --spacing-lg: 24px;
    --spacing-xl: 32px;

    /* Breakpoints (for reference, actual use in media queries) */
    --breakpoint-mobile: 768px;
    --breakpoint-tablet: 1024px;
}
```

**File Upload Data Model:**

```javascript
// web/js/upload.js

/**
 * UploadedFile represents a file uploaded by the user
 * @typedef {Object} UploadedFile
 * @property {File} file - Browser File object
 * @property {string} id - Unique ID (timestamp + random)
 * @property {string} format - Detected format (np3, xmp, lrtemplate, costyle, dcp)
 * @property {string} status - Status: 'queued', 'processing', 'complete', 'error'
 * @property {ArrayBuffer} data - File contents as ArrayBuffer
 * @property {Uint8Array} convertedData - Converted file bytes (after conversion)
 * @property {string} targetFormat - User-selected target format
 * @property {string} error - Error message if status === 'error'
 */

/**
 * Batch conversion progress
 * @typedef {Object} BatchProgress
 * @property {number} total - Total files to convert
 * @property {number} completed - Files converted successfully
 * @property {number} failed - Files that failed conversion
 * @property {number} queued - Files waiting to be converted
 */
```

**Format Badge Component (BEM naming):**

```html
<!-- components.css defines .format-badge and modifiers -->
<span class="format-badge format-badge--np3">NP3</span>
<span class="format-badge format-badge--xmp">XMP</span>
<span class="format-badge format-badge--lrtemplate">lrtemplate</span>
<span class="format-badge format-badge--costyle">Capture One</span>
<span class="format-badge format-badge--dcp">DCP</span>
```

```css
/* web/css/components.css */
.format-badge {
    display: inline-block;
    padding: var(--spacing-xs) var(--spacing-sm);
    border-radius: 4px;
    font-size: var(--font-size-small);
    font-weight: var(--font-weight-bold);
    text-transform: uppercase;
    color: white;
}

.format-badge--np3 { background-color: var(--color-np3); }
.format-badge--xmp { background-color: var(--color-xmp); }
.format-badge--lrtemplate { background-color: var(--color-lrtemplate); }
.format-badge--costyle { background-color: var(--color-costyle); }
.format-badge--dcp { background-color: var(--color-dcp); }
```

### APIs and Interfaces

**upload.js Public API:**

```javascript
// web/js/upload.js

/**
 * Initialize drag-drop upload zone
 * @param {HTMLElement} dropZone - Drop zone element
 * @param {Function} onFilesAdded - Callback (files: UploadedFile[]) => void
 */
export function initUpload(dropZone, onFilesAdded) {
    dropZone.addEventListener('dragover', handleDragOver);
    dropZone.addEventListener('drop', handleDrop);

    // File picker
    const fileInput = dropZone.querySelector('input[type="file"]');
    fileInput.addEventListener('change', handleFileSelect);
}

/**
 * Detect format from file bytes
 * @param {Uint8Array} data - File bytes
 * @param {string} filename - Original filename
 * @returns {string} Format: 'np3', 'xmp', 'lrtemplate', 'costyle', 'dcp', or 'unknown'
 */
export async function detectFormat(data, filename) {
    // Check file extension first
    const ext = filename.split('.').pop().toLowerCase();
    if (['np3', 'xmp', 'lrtemplate', 'costyle', 'dcp'].includes(ext)) {
        return ext;
    }

    // Magic byte detection (fallback)
    if (data[0] === 0x4E && data[1] === 0x49 && data[2] === 0x43 && data[3] === 0x4F) {
        return 'np3'; // "NICO" magic bytes
    }
    if (data[0] === 0x3C && data[1] === 0x3F && data[2] === 0x78 && data[3] === 0x6D) {
        // "<?xm" XML header - could be XMP or Capture One
        // Check for specific namespace
        const text = new TextDecoder().decode(data.slice(0, 1000));
        if (text.includes('x:xmpmeta')) {
            if (text.includes('crs:Exposure')) {
                return 'costyle'; // Capture One
            }
            return 'xmp'; // Adobe XMP
        }
    }
    if (data[0] === 0x49 && data[1] === 0x49 || data[0] === 0x4D && data[1] === 0x4D) {
        return 'dcp'; // TIFF magic bytes (II or MM)
    }

    return 'unknown';
}

/**
 * Create file card DOM element
 * @param {UploadedFile} file - Uploaded file
 * @returns {HTMLElement} Card element
 */
export function createFileCard(file) {
    const card = document.createElement('div');
    card.className = 'file-card';
    card.dataset.fileId = file.id;

    card.innerHTML = `
        <div class="file-card__header">
            <span class="file-card__filename">${file.file.name}</span>
            <span class="format-badge format-badge--${file.format}">${file.format.toUpperCase()}</span>
        </div>
        <div class="file-card__size">${formatFileSize(file.file.size)}</div>
        <div class="file-card__status" data-status="${file.status}">
            ${getStatusIcon(file.status)} ${getStatusText(file.status)}
        </div>
        <div class="file-card__actions">
            ${file.status === 'complete' ? '<button class="btn btn--download">Download</button>' : ''}
        </div>
    `;

    return card;
}

/**
 * Update file card status
 * @param {string} fileId - File ID
 * @param {string} status - New status: 'queued', 'processing', 'complete', 'error'
 * @param {string} [error] - Error message if status === 'error'
 */
export function updateFileStatus(fileId, status, error = null) {
    const card = document.querySelector(`[data-file-id="${fileId}"]`);
    if (!card) return;

    const statusEl = card.querySelector('.file-card__status');
    statusEl.dataset.status = status;
    statusEl.innerHTML = `${getStatusIcon(status)} ${error || getStatusText(status)}`;

    // Add download button if complete
    if (status === 'complete') {
        const actionsEl = card.querySelector('.file-card__actions');
        actionsEl.innerHTML = '<button class="btn btn--download">Download</button>';
    }
}
```

**app.js Main Application:**

```javascript
// web/js/app.js

import { initUpload, detectFormat, createFileCard, updateFileStatus } from './upload.js';
import { convertFile } from './converter.js'; // Existing WASM interface
import { showPreview } from './preview.js'; // Epic 11

// State
let uploadedFiles = [];
let batchProgress = { total: 0, completed: 0, failed: 0, queued: 0 };

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', () => {
    const dropZone = document.getElementById('drop-zone');
    const fileGrid = document.getElementById('file-grid');
    const batchConvertBtn = document.getElementById('batch-convert');
    const startFreshBtn = document.getElementById('start-fresh');

    // Initialize upload
    initUpload(dropZone, async (files) => {
        for (const file of files) {
            const data = await file.arrayBuffer();
            const format = await detectFormat(new Uint8Array(data), file.name);

            const uploadedFile = {
                file,
                id: Date.now() + Math.random(),
                format,
                status: 'queued',
                data,
                convertedData: null,
                targetFormat: null,
                error: null
            };

            uploadedFiles.push(uploadedFile);
            fileGrid.appendChild(createFileCard(uploadedFile));
        }

        updateBatchProgress();
    });

    // Batch convert
    batchConvertBtn.addEventListener('click', () => {
        const targetFormat = document.getElementById('target-format-select').value;
        batchConvert(targetFormat);
    });

    // Start fresh
    startFreshBtn.addEventListener('click', () => {
        uploadedFiles = [];
        fileGrid.innerHTML = '';
        batchProgress = { total: 0, completed: 0, failed: 0, queued: 0 };
        updateBatchProgress();
    });
});

async function batchConvert(targetFormat) {
    batchProgress.total = uploadedFiles.length;
    batchProgress.queued = uploadedFiles.length;

    for (const file of uploadedFiles) {
        if (file.status !== 'queued') continue;

        updateFileStatus(file.id, 'processing');

        try {
            // Convert using WASM
            const converted = await convertFile(new Uint8Array(file.data), targetFormat);
            file.convertedData = converted;
            file.targetFormat = targetFormat;
            file.status = 'complete';

            updateFileStatus(file.id, 'complete');
            batchProgress.completed++;
        } catch (err) {
            file.status = 'error';
            file.error = err.message;

            updateFileStatus(file.id, 'error', err.message);
            batchProgress.failed++;
        }

        batchProgress.queued--;
        updateBatchProgress();
    }
}

function updateBatchProgress() {
    const progressEl = document.getElementById('batch-progress');
    const { total, completed, failed, queued } = batchProgress;

    if (total === 0) {
        progressEl.textContent = 'Upload files to begin';
        return;
    }

    if (queued === 0) {
        progressEl.textContent = `Complete: ${completed} files, ${failed} errors`;
    } else {
        progressEl.textContent = `Converting ${completed + failed + 1} of ${total}...`;
    }
}
```

### Workflows and Sequencing

**File Upload Workflow:**

1. User drags files into drop zone or clicks to browse
2. Browser FileList → Array of File objects
3. For each file:
   a. Read file as ArrayBuffer
   b. Detect format via extension + magic bytes
   c. Create UploadedFile object (status: 'queued')
   d. Append file card to grid (shows filename, format badge, file size)
4. Enable batch convert button

**Batch Conversion Workflow:**

1. User selects target format from dropdown (e.g., "XMP")
2. User clicks "Convert All" button
3. For each uploaded file (status === 'queued'):
   a. Update card status to 'processing' (show spinner)
   b. Call `convertFile(data, targetFormat)` (WASM conversion)
   c. On success:
      - Store converted bytes in file.convertedData
      - Update card status to 'complete' (show checkmark + download button)
      - Increment batchProgress.completed
   d. On error:
      - Store error message in file.error
      - Update card status to 'error' (show error icon + message)
      - Increment batchProgress.failed
4. Update batch progress text: "Converting 3 of 10..." → "Complete: 8 files, 2 errors"

**Individual File Conversion (Optional):**

1. User clicks individual file card
2. Modal opens with target format dropdown (per-file selection)
3. User selects format, clicks "Convert"
4. Same conversion logic as batch, but for single file

## Non-Functional Requirements

### Performance

**Targets:**
- Initial page load: <2 seconds on 3G connection (1.6 Mbps)
- File upload feedback: <100ms (drag-over visual feedback)
- Format detection: <50ms per file (client-side, no network)
- Conversion: <100ms per file (existing WASM performance)
- Responsive layout: 60fps scrolling/animations on mobile

**Optimization Strategies:**
- Use system fonts (no web font download)
- Inline critical CSS (above-the-fold styles)
- Lazy load preview images (Epic 11) until modal opened
- Use CSS transforms for animations (GPU-accelerated)
- Minimize DOM manipulation (batch updates)

**Performance Monitoring:**
- Lighthouse CI: Target score ≥95 for Performance, Accessibility, Best Practices
- WebPageTest: Load time <2s on 3G
- Manual testing: Real devices (iPhone SE, Android budget phone)

### Security

**Threats & Mitigations:**
- **XSS via filename display**: Sanitize filenames before innerHTML (use textContent)
- **File bomb (huge uploads)**: Warn user if file >10MB (but allow, client-side processing)
- **Malicious files**: All processing in WASM sandbox (no native code execution)
- **CORS issues**: Not applicable (no external API calls)

**Privacy:**
- Zero analytics, zero tracking (as per Recipe core principle)
- All processing client-side (WASM)
- No localStorage usage for uploaded files (only ephemeral state)
- Clear "Privacy: Your files never leave your device" messaging

**Input Validation:**
- Accept only supported file extensions via file picker
- Validate format via magic bytes (reject unknown formats gracefully)
- No maximum file count (batch upload)

### Reliability/Availability

**Error Handling:**
- Unsupported file format: Clear error message "Unsupported file type"
- Conversion failure: Show per-file error message (don't halt batch)
- WASM load failure: Fallback message "Please refresh the page"
- Drag-drop not supported: Graceful degradation to file picker

**Browser Compatibility:**
- Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- Fallback for older browsers: Basic file picker (no drag-drop)
- Progressive enhancement: Works without JavaScript for file download link

**Failure Modes:**
- JavaScript disabled: Show static conversion info, download links
- WASM not supported: Show browser upgrade message
- Network offline: All conversion works (client-side only)

### Observability

**Logging:**
- Console warnings for unsupported file formats
- Console errors for conversion failures (include file name, format)
- No server-side logging (client-side app)

**Metrics (client-side only, no analytics):**
- Not required for personal project
- Could add optional localStorage metrics (opt-in): conversions per format

**Debugging:**
- Verbose mode (add `?debug=1` to URL): Log all conversion steps to console
- File format detection debugging (show detected magic bytes)

## Dependencies and Integrations

**Dependencies:**

Current (no changes):
```
- None (vanilla JavaScript)
- WASM: recipe.wasm (existing)
- Go runtime: wasm_exec.js (existing)
```

New (zero external dependencies):
- Standard Web APIs: FileReader, Blob, URL.createObjectURL
- CSS Grid, Flexbox (native browser support)
- ES6 modules (import/export)

**External Integrations:**
- Cloudflare Pages (existing deployment)
- Browser compatibility testing: BrowserStack (manual testing)

**Version Constraints:**
- ES6+ JavaScript (target: Chrome 90+, Firefox 88+, Safari 14+)
- CSS Grid Level 1 (all modern browsers)
- WebAssembly 1.0 (all modern browsers)

## Acceptance Criteria (Authoritative)

**AC-1: Redesigned Landing Page**
- ✅ Hero section with clear value proposition: "Convert Photo Presets. Instantly. Privately."
- ✅ Visual format badges displayed prominently with brand colors
- ✅ Single-page layout (no navigation to other pages for core conversion)
- ✅ Clean typography using system fonts for performance
- ✅ Responsive design: works on mobile (320px+), tablet (768px+), desktop (1024px+)
- ✅ Fast load time: <2 seconds on 3G connection (WebPageTest validation)
- ✅ No external dependencies (no CDN fonts, no analytics trackers)

**AC-2: Visual Format Badges**
- ✅ Badge system implemented with defined colors (NP3=#FFC107, XMP=#0073E6, lrtemplate=#D81B60, Capture One=#9C27B0, DCP=#4CAF50)
- ✅ Badges shown on landing page, in upload cards, in conversion dropdowns
- ✅ Accessible: Color not sole indicator (includes format name text)
- ✅ Responsive: Badges scale/stack appropriately on mobile
- ✅ Consistent styling across all interface elements (BEM naming)

**AC-3: Batch File Upload with Drag-and-Drop**
- ✅ Large drop zone on landing page invites drag-and-drop
- ✅ Visual feedback on drag-over (highlight, scale animation)
- ✅ Support multiple file selection via file picker
- ✅ Accept all supported formats (.np3, .xmp, .lrtemplate, .costyle, .dcp)
- ✅ Reject unsupported files with clear error message ("Unsupported file type")
- ✅ No file size limit (client-side processing handles large files)
- ✅ Display uploaded files as individual cards in grid layout
- ✅ Each card shows: filename, detected format badge, file size, conversion status

**AC-4: Progress Indicators**
- ✅ Batch conversion shows overall progress: "Converting 3 of 10..."
- ✅ Per-file status indicators: queued (clock icon), processing (spinner), complete (checkmark), error (X icon)
- ✅ Smooth transitions between states (CSS transitions, no janky updates)
- ✅ Visual feedback during processing (spinner, progress bar optional)
- ✅ Completion state with download buttons per file
- ✅ Error state shows specific error message per file
- ✅ Users can cancel in-progress batch conversions (abort button)

**AC-5: Mobile-Responsive Design**
- ✅ Mobile (<768px): Single column, stacked layout, tap-to-browse upload
- ✅ Tablet (768-1024px): Two-column grid for batch files
- ✅ Desktop (>1024px): Three-column grid, full features visible
- ✅ Touch-friendly targets (44px minimum) on mobile
- ✅ No horizontal scrolling on any device size
- ✅ Readable text without zooming (16px minimum body text)
- ✅ Test on real devices: iPhone, Android, iPad, Desktop browsers (manual testing)

**AC-6: Before/After Comparison Slider** (Epic 11 integration)
- ✅ Slider implemented with draggable handle
- ✅ Smooth drag interaction with visual feedback
- ✅ Keyboard accessible (arrow keys move slider)
- ✅ Mobile: Tap-and-hold to compare, or tap sides to snap 50/50
- ✅ Slider position persists during modal session
- ✅ Clear visual indicators (before/after labels)
- ✅ Works across all browsers (Chrome, Firefox, Safari, Edge)

**AC-7: Improved Conversion Flow**
- ✅ Format detection happens instantly on upload (client-side, <50ms)
- ✅ "Convert to..." dropdown shows only valid target formats
- ✅ Batch convert: Single action converts all files to same format
- ✅ Individual convert: Per-file format selection available (optional)
- ✅ Conversion happens instantly (<100ms per file)
- ✅ Download buttons appear immediately after conversion
- ✅ Option to "Convert More" without page reload (persistent state)
- ✅ "Start Fresh" clears all files and resets interface

## Traceability Mapping

| AC ID | Spec Section(s) | Component(s)/API(s) | Test Idea |
| ----- | --------------- | ------------------- | --------- |
| AC-1 | Data Models (HTML/CSS) | index.html, main.css | Manual: Lighthouse audit, 3G load test |
| AC-2 | Data Models (CSS) | components.css (format-badge) | Manual: Visual regression test |
| AC-3 | APIs (upload.js) | initUpload, detectFormat, createFileCard | Unit test: Drag-drop events, format detection |
| AC-4 | APIs (app.js) | updateFileStatus, updateBatchProgress | Manual: Multi-file conversion, observe status |
| AC-5 | Data Models (CSS) | layout.css (responsive grid) | Manual: BrowserStack device testing |
| AC-6 | APIs (preview.js) | showPreview, slider interaction | Manual: Browser compatibility testing |
| AC-7 | Workflows (Conversion) | app.js (batchConvert) | Manual: End-to-end conversion flow |

**Test Coverage Targets:**
- Unit tests: `upload.test.js` (format detection, file validation)
- Integration tests: End-to-end conversion flow (CLI automated testing not applicable for web UI)
- Manual tests: Browser compatibility (Chrome, Firefox, Safari, Edge), Device testing (iPhone, Android, iPad)
- Coverage target: 80% for JavaScript modules (lower than backend due to DOM dependencies)

## Risks, Assumptions, Open Questions

### Risks

**RISK-1: 3G load time target (<2s) difficult to achieve**
- **Severity**: Medium
- **Impact**: WASM bundle (~2MB) may exceed budget on slow connections
- **Mitigation**: Inline critical CSS, lazy load preview images, optimize WASM size
- **Owner**: Dev (Epic 10)

**RISK-2: Mobile drag-drop inconsistent across browsers**
- **Severity**: Low
- **Impact**: Some mobile browsers may not support drag-drop well
- **Mitigation**: Prioritize file picker for mobile, drag-drop as enhancement
- **Owner**: Dev (Epic 10)

**RISK-3: CSS Grid not supported on very old browsers**
- **Severity**: Very Low
- **Impact**: Users on IE11 or Safari <10 won't see grid layout
- **Mitigation**: Acceptable - Recipe targets modern browsers (2020+)
- **Owner**: Justin (product decision)

### Assumptions

**ASSUMPTION-1**: System fonts provide sufficient visual quality
- **Rationale**: System fonts load instantly, widely available, professional appearance
- **Validation**: Visual design review
- **Risk if false**: May need to add web fonts (increase load time)

**ASSUMPTION-2**: Vanilla JavaScript sufficient (no framework needed)
- **Rationale**: Recipe UI is simple (single-page, no complex state)
- **Validation**: Implement MVP and assess maintainability
- **Risk if false**: May need to refactor to framework later (defer to Path B/C)

**ASSUMPTION-3**: Batch conversion can be sequential (no parallelism needed)
- **Rationale**: Conversions are fast (<100ms), sequential provides clear progress
- **Validation**: Test with 50+ file batch
- **Risk if false**: May need Web Workers for parallel conversion (nice-to-have)

### Open Questions

**Q-1**: Should Recipe support saving conversion preferences (localStorage)?
- **Impact**: User convenience vs privacy principle (no tracking)
- **Resolution**: DEFER to future - start with no persistent state, add if users request
- **Owner**: Justin (product decision)

**Q-2**: Should individual file conversion be supported, or only batch?
- **Impact**: Additional UI complexity vs flexibility
- **Resolution**: YES - provide both batch (primary) and individual (secondary) conversion
- **Owner**: Dev (Epic 10)

**Q-3**: Should Recipe show tooltips/help text for first-time users?
- **Impact**: User onboarding vs clean UI
- **Resolution**: YES - subtle help text in hero section, tooltips for non-obvious controls
- **Owner**: Dev (Epic 10), UX decision

## Test Strategy Summary

### Test Levels

**Unit Tests (JavaScript modules):**
- `upload.test.js`: Format detection, file validation, createFileCard
- `utils.test.js`: Utility functions, formatFileSize
- Coverage target: 80% for JavaScript modules

**Integration Tests:**
- End-to-end conversion flow: Upload → Detect → Convert → Download
- Cross-browser compatibility (manual): Chrome, Firefox, Safari, Edge

**Manual Validation:**
- Lighthouse audit: Performance ≥95, Accessibility ≥95, Best Practices ≥95
- WebPageTest: Load time <2s on 3G
- Device testing: iPhone (Safari), Android (Chrome), iPad, Desktop browsers
- Responsive design: Test all breakpoints (320px, 768px, 1024px)

### Test Frameworks

- Jest (or similar): JavaScript unit testing
- Lighthouse CI: Performance auditing
- BrowserStack: Cross-browser/device testing (manual)

### Coverage of ACs

| AC ID | Test Type | Test Location | Coverage |
| ----- | --------- | ------------- | -------- |
| AC-1 | Manual | Lighthouse, WebPageTest | Load time, performance score |
| AC-2 | Manual | Visual regression | Badge colors, consistency |
| AC-3 | Unit + Manual | upload.test.js, Browser testing | Drag-drop, file validation |
| AC-4 | Manual | Browser testing | Status indicators, progress text |
| AC-5 | Manual | BrowserStack devices | Responsive layout, touch targets |
| AC-6 | Manual | Browser testing | Slider interaction, keyboard nav |
| AC-7 | Manual | End-to-end testing | Full conversion flow |

### Edge Cases

**Edge Case Testing:**
- Upload 100+ files (stress test batch conversion)
- Upload single file >10MB (warn but allow)
- Upload unsupported file type (show clear error)
- Upload files with special characters in filename (sanitize display)
- Drag-drop on mobile Safari (may not work, fallback to file picker)
- Offline mode (ensure conversion still works, WASM is local)
- Browser back button (state lost - acceptable for MVP)
- Multiple target format selections (batch convert, individual convert)
- Cancel batch conversion mid-process (abort button)

---

**Next Steps:**
1. Create `web/css/` directory structure (main.css, components.css, layout.css, preview.css)
2. Refactor `web/index.html` with hero section and drop zone
3. Implement `web/js/upload.js` (drag-drop, format detection)
4. Implement `web/js/app.js` (batch conversion logic)
5. Implement responsive grid in `layout.css` (mobile, tablet, desktop breakpoints)
6. Implement format badge components in `components.css` (BEM naming)
7. Write unit tests for upload.js and utils.js
8. Manual testing: Lighthouse audit, WebPageTest, BrowserStack devices
9. Update README.md with new web UI features
10. Mark epic-10 as "contexted" in sprint-status.yaml
