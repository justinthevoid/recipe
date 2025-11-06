# Story 2-7: File Download Trigger

**Epic:** Epic 2 - Web Interface (FR-2)
**Story ID:** 2-7
**Status:** drafted
**Created:** 2025-11-04
**Complexity:** Simple (1 day)

---

## User Story

**As a** photographer
**I want** to download my converted preset file
**So that** I can use it in my photography software

---

## Business Value

File download is the **final step** in the conversion flow - the moment users get their converted preset. This story completes the end-to-end user journey:

1. Upload file → 2. Detect format → 3. Preview parameters → 4. Select target format → 5. Convert → **6. Download**

**Success criteria:**
- Download works reliably (100% success rate)
- Correct filename and extension
- File opens in target software (Nikon ViewNX, Lightroom)

**This story delivers the "payoff" - users walk away with a usable preset file.**

---

## Acceptance Criteria

### AC-1: Display Download Button After Conversion

- [ ] Download button appears after conversion succeeds (Story 2-6)
- [ ] Button initially disabled (no converted data yet)
- [ ] Button enabled after `conversionComplete` event
- [ ] Button text: "Download [FileName].[ext]" (e.g., "Download Classic Chrome.xmp")
- [ ] Button style: Primary action (blue background, prominent)

**Test:**
1. Upload file → convert to XMP
2. Conversion succeeds
3. Verify: Download button appears below Convert button
4. Verify: Button text shows correct filename (e.g., "Download Classic Chrome.xmp")
5. Verify: Button enabled (not grayed out)

### AC-2: Generate Download Link from Converted Data

- [ ] Create Blob from converted Uint8Array
- [ ] Set correct MIME type:
  - NP3: `application/octet-stream` (binary)
  - XMP: `application/xml` or `text/xml`
  - lrtemplate: `text/plain` (Lua text)
- [ ] Generate object URL using `URL.createObjectURL()`
- [ ] Store URL for download trigger

**Test:**
1. Convert file to XMP
2. Verify: Blob created from Uint8Array
3. Verify: Blob size matches converted data size
4. Verify: Object URL generated (starts with "blob:http://...")
5. Verify: Console logs "Download link created: [URL]"

### AC-3: Trigger Browser Download on Button Click

- [ ] Create temporary `<a>` element with `download` attribute
- [ ] Set `href` to object URL
- [ ] Set `download` attribute to correct filename
- [ ] Programmatically click `<a>` to trigger download
- [ ] Browser saves file to default downloads folder

**Test:**
1. Convert file to XMP
2. Click "Download Classic Chrome.xmp" button
3. Verify: Browser download dialog appears (or file auto-saves)
4. Verify: File saved with correct name "Classic Chrome.xmp"
5. Verify: File saved to default downloads folder

### AC-4: Correct Filename and Extension

- [ ] Output filename = input filename with new extension
- [ ] Extension matches target format:
  - NP3: `.np3`
  - XMP: `.xmp`
  - lrtemplate: `.lrtemplate`
- [ ] Preserve original filename (don't add timestamps or random IDs)
- [ ] Handle special characters in filename (sanitize if needed)

**Examples:**
```
Input: Classic Chrome.np3 → Output (XMP): Classic Chrome.xmp
Input: My Preset.xmp → Output (NP3): My Preset.np3
Input: Auto Tone.lrtemplate → Output (XMP): Auto Tone.xmp
```

**Test:**
1. Upload "Classic Chrome.np3" → convert to XMP → download
2. Verify: Downloaded file named "Classic Chrome.xmp"
3. Upload "My Preset (v2).xmp" → convert to lrtemplate → download
4. Verify: Downloaded file named "My Preset (v2).lrtemplate"

### AC-5: Clean Up Object URLs (Memory Management)

- [ ] Revoke object URL after download completes
- [ ] Use `URL.revokeObjectURL()` to free memory
- [ ] Revoke previous URLs when new conversion starts
- [ ] No memory leak from accumulated object URLs

**Test:**
1. Convert file → download (object URL created)
2. Convert again → download (new object URL created, old one revoked)
3. Repeat 10 times
4. Verify: Memory usage stable (DevTools memory profiler)
5. Verify: Only one object URL active at a time

### AC-6: Handle Download Errors

- [ ] If download fails (browser blocks, no disk space):
  - Show error: "Download failed. Please check your browser settings and try again."
  - Log technical error to console
  - Keep download button enabled (allow retry)
- [ ] Handle browser popup blockers (if applicable)

**Test:**
1. Block downloads in browser settings
2. Click Download button
3. Verify: Error message displayed (user-friendly)
4. Verify: Console shows technical error
5. Unblock downloads → click Download → download succeeds

### AC-7: Visual Feedback During Download

- [ ] Button text changes to "Downloading..." (briefly)
- [ ] Button disabled during download (prevent double-click)
- [ ] Button returns to "Download [FileName]" after download completes
- [ ] Success message: "✓ Download complete!"

**Test:**
1. Click Download button
2. Verify: Button shows "Downloading..." (brief, <1s)
3. Download completes
4. Verify: Button returns to "Download Classic Chrome.xmp"
5. Verify: Success message displayed

### AC-8: Reset UI After Download (Optional)

- [ ] Optionally show "Convert Another File" button after download
- [ ] Allow user to upload new file without refreshing page
- [ ] Clear previous conversion data when new file uploaded

**Test:**
1. Complete full conversion flow → download file
2. Verify: "Convert Another File" button appears (optional)
3. Click "Convert Another File" → UI resets to default state
4. Upload new file → conversion flow works again

---

## Technical Approach

### Download Module

**File:** `web/static/downloader.js` (new file)

```javascript
// downloader.js - File download handling

let currentDownloadURL = null;

/**
 * Enable download button with converted file data
 * @param {Uint8Array} fileData - Converted file data
 * @param {string} fileName - Output filename with extension
 * @param {string} format - Target format ("np3" | "xmp" | "lrtemplate")
 */
export function enableDownload(fileData, fileName, format) {
    if (!fileData || !fileName) {
        throw new Error('Missing required parameters');
    }

    // Revoke previous download URL if exists
    revokeDownloadURL();

    // Create Blob with appropriate MIME type
    const mimeType = getMimeType(format);
    const blob = new Blob([fileData], { type: mimeType });

    // Create object URL
    currentDownloadURL = URL.createObjectURL(blob);

    console.log(`Download link created: ${fileName} (${blob.size} bytes)`);

    // Update download button
    updateDownloadButton(fileName);
}

/**
 * Get MIME type for format
 */
function getMimeType(format) {
    const mimeTypes = {
        np3: 'application/octet-stream', // Binary format
        xmp: 'application/xml',          // XML format
        lrtemplate: 'text/plain',        // Lua text format
    };

    return mimeTypes[format] || 'application/octet-stream';
}

/**
 * Update download button with filename
 */
function updateDownloadButton(fileName) {
    const downloadButton = document.getElementById('downloadButton');
    if (downloadButton) {
        downloadButton.disabled = false;
        downloadButton.textContent = `Download ${fileName}`;
        downloadButton.style.display = 'block';

        // Remove any previous event listeners (avoid duplicates)
        const newButton = downloadButton.cloneNode(true);
        downloadButton.parentNode.replaceChild(newButton, downloadButton);

        // Add new event listener
        newButton.addEventListener('click', () => handleDownload(fileName));
    }
}

/**
 * Handle download button click
 */
function handleDownload(fileName) {
    if (!currentDownloadURL) {
        showDownloadError('Download link not available. Please convert file again.');
        return;
    }

    console.log(`Downloading: ${fileName}`);

    // Show downloading state
    showDownloadingState();

    try {
        // Create temporary <a> element
        const link = document.createElement('a');
        link.href = currentDownloadURL;
        link.download = fileName;
        link.style.display = 'none';

        // Append to body (required for Firefox)
        document.body.appendChild(link);

        // Trigger download
        link.click();

        // Clean up
        document.body.removeChild(link);

        // Show success state (after brief delay)
        setTimeout(() => {
            showDownloadSuccess(fileName);
        }, 500);

    } catch (error) {
        console.error('Download error:', error);
        showDownloadError('Download failed. Please check your browser settings and try again.');
    }
}

/**
 * Show downloading state
 */
function showDownloadingState() {
    const downloadButton = document.getElementById('downloadButton');
    if (downloadButton) {
        downloadButton.disabled = true;
        downloadButton.textContent = 'Downloading...';
    }
}

/**
 * Show download success
 */
function showDownloadSuccess(fileName) {
    const downloadButton = document.getElementById('downloadButton');
    if (downloadButton) {
        downloadButton.disabled = false;
        downloadButton.textContent = `Download ${fileName}`;
    }

    // Show success message
    const statusEl = document.getElementById('downloadStatus');
    if (statusEl) {
        statusEl.className = 'status success';
        statusEl.textContent = '✓ Download complete!';
        statusEl.style.display = 'block';

        // Hide after 3 seconds
        setTimeout(() => {
            statusEl.style.display = 'none';
        }, 3000);
    }
}

/**
 * Show download error
 */
function showDownloadError(message) {
    const downloadButton = document.getElementById('downloadButton');
    if (downloadButton) {
        downloadButton.disabled = false; // Re-enable for retry
    }

    const errorEl = document.getElementById('downloadError');
    if (errorEl) {
        errorEl.textContent = message;
        errorEl.style.display = 'block';
    }
}

/**
 * Revoke current download URL (free memory)
 */
function revokeDownloadURL() {
    if (currentDownloadURL) {
        URL.revokeObjectURL(currentDownloadURL);
        currentDownloadURL = null;
        console.log('Previous download URL revoked');
    }
}

/**
 * Clear download state
 */
export function clearDownloadState() {
    revokeDownloadURL();

    const downloadButton = document.getElementById('downloadButton');
    if (downloadButton) {
        downloadButton.style.display = 'none';
        downloadButton.disabled = true;
    }

    const statusEl = document.getElementById('downloadStatus');
    if (statusEl) {
        statusEl.style.display = 'none';
    }

    const errorEl = document.getElementById('downloadError');
    if (errorEl) {
        errorEl.style.display = 'none';
    }
}
```

### Integration with Main Flow

**Update `main.js`:**

```javascript
// main.js - Integrate download

import { initializeDropZone, handleFile } from './file-handler.js';
import { detectFileFormat } from './format-detector.js';
import { displayParameters } from './parameter-display.js';
import { displayFormatSelector } from './format-selector.js';
import { convertFile, getConvertedFileData, getConvertedFileName, getConvertedFileFormat } from './converter.js';
import { enableDownload, clearDownloadState } from './downloader.js';
import { initializeWASM } from './wasm-loader.js';

// Initialize WASM
initializeWASM();

// Initialize UI
document.addEventListener('DOMContentLoaded', () => {
    initializeDropZone(handleFile);
});

// Listen for conversion complete event (Story 2-6)
window.addEventListener('conversionComplete', (event) => {
    const { convertedData, format } = event.detail;

    // Get converted file metadata
    const fileName = getConvertedFileName();

    // Enable download button
    enableDownload(convertedData, fileName, format);

    console.log('Download ready:', fileName);
});

// Listen for new file uploaded (clear previous download state)
window.addEventListener('fileLoaded', () => {
    clearDownloadState();
});
```

### CSS for Download Button

**Add to `web/static/style.css`:**

```css
/* Download button */
.download-button {
    width: 100%;
    margin-top: 1rem;
    padding: 0.875rem 1.5rem;
    background: #3182ce;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
}

.download-button:hover:not(:disabled) {
    background: #2c5aa0;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.download-button:active:not(:disabled) {
    transform: translateY(0);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.download-button:disabled {
    background: #cbd5e0;
    cursor: not-allowed;
    transform: none;
}

/* Download status */
#downloadStatus {
    margin-top: 0.5rem;
}

#downloadError {
    margin-top: 0.5rem;
}
```

### HTML Updates

**Add to `web/index.html`:**

```html
<!-- Download Button (appears after conversion) -->
<button id="downloadButton" class="download-button" style="display: none;" disabled>
    Download Converted Preset
</button>

<!-- Download Status -->
<div id="downloadStatus" class="status" style="display: none;" role="status" aria-live="polite"></div>

<!-- Download Error -->
<div id="downloadError" class="error-message" style="display: none;" role="alert"></div>
```

---

## Dependencies

### Required Before Starting

- ✅ Story 2-6 complete (converted file data available)

### No Blocking Dependencies

Story 2-7 is the final story in the core conversion flow. No other stories depend on it.

---

## Testing Plan

### Manual Testing

**Test Case 1: NP3 → XMP Download**
1. Upload `Classic Chrome.np3`
2. Convert to XMP
3. Verify: Download button appears with text "Download Classic Chrome.xmp"
4. Click Download button
5. Verify: Browser downloads file "Classic Chrome.xmp"
6. Open file in Lightroom → verify preset loads correctly

**Test Case 2: XMP → NP3 Download**
1. Upload `Classic Chrome.xmp`
2. Convert to NP3
3. Click Download button
4. Verify: Downloaded file "Classic Chrome.np3"
5. Open file in Nikon ViewNX → verify preset loads correctly

**Test Case 3: lrtemplate → XMP Download**
1. Upload `Auto Tone.lrtemplate`
2. Convert to XMP
3. Click Download button
4. Verify: Downloaded file "Auto Tone.xmp"
5. Verify: File size matches converted data size

**Test Case 4: Special Characters in Filename**
1. Upload file named "My Preset (v2) [Final].xmp"
2. Convert to NP3
3. Click Download button
4. Verify: Downloaded file "My Preset (v2) [Final].np3"
5. Verify: Filename preserves special characters

**Test Case 5: Multiple Downloads (Same File)**
1. Convert file → download
2. Click Download button again (without re-converting)
3. Verify: File downloads again (same data)
4. Repeat 5 times → verify all downloads succeed

**Test Case 6: Multiple Conversions (Different Files)**
1. Upload file A → convert to XMP → download
2. Upload file B → convert to NP3 → download
3. Verify: Each download has correct filename and data
4. Verify: Previous object URL revoked (memory stable)

**Test Case 7: Download Error Handling**
1. Block downloads in browser settings (Chrome: chrome://settings/content/pdfDocuments → Block)
2. Convert file → click Download button
3. Verify: Error message displayed (or browser shows blocked download notification)
4. Unblock downloads → click Download → download succeeds

**Test Case 8: Memory Management**
1. Convert and download 10 different files
2. Verify: Memory usage stable (DevTools memory profiler)
3. Verify: No memory leak from accumulated object URLs

**Test Case 9: Reset UI After Download**
1. Complete full flow: upload → convert → download
2. Upload new file (without refreshing page)
3. Verify: Download button disappears (old download cleared)
4. Convert new file → new download button appears

### File Validation Testing

**Validate downloaded files work in target software:**

1. **NP3 → XMP:**
   - Download "Classic Chrome.xmp"
   - Import into Lightroom CC
   - Verify: Preset appears in Presets panel
   - Apply to photo → verify parameters apply correctly

2. **XMP → NP3:**
   - Download "Classic Chrome.np3"
   - Copy to Nikon camera SD card
   - Load in camera Picture Control settings
   - Verify: Preset loads and parameters display correctly

3. **lrtemplate → XMP:**
   - Download "Auto Tone.xmp"
   - Import into Lightroom CC
   - Apply to photo → verify adjustments work

### Browser Compatibility

Test in:
- ✅ Chrome (latest) - Blob API, download attribute fully supported
- ✅ Firefox (latest) - Blob API, download attribute fully supported
- ✅ Safari (latest) - Blob API, download attribute fully supported

**Expected:** Identical behavior across browsers.

---

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Download works for all 3 formats (NP3, XMP, lrtemplate)
- [ ] Correct filenames and extensions
- [ ] Downloaded files open correctly in target software
- [ ] Object URLs properly revoked (no memory leak)
- [ ] Error handling tested (browser blocks, etc.)
- [ ] Manual testing completed in Chrome, Firefox, Safari
- [ ] File validation completed (files work in Lightroom, Nikon ViewNX)
- [ ] Code reviewed
- [ ] Integration with Story 2-6 verified
- [ ] Story marked "ready-for-dev" in sprint status

---

## Out of Scope

**Explicitly NOT in this story:**
- ❌ Download history (save last 5 downloads)
- ❌ Download location selection (use browser default)
- ❌ Download progress bar (files are small, instant download)
- ❌ Batch download (multiple files at once - Epic 3)

**This story only delivers:** Single file download - download converted preset with correct filename.

---

## Technical Notes

### Why Blob API?

**Alternative:** Data URLs (`data:application/octet-stream;base64,...`)

**Decision:** Use Blob API with object URLs

**Rationale:**
- **Performance:** Object URLs don't encode data as base64 (no overhead)
- **Memory:** Object URLs reference blob, not duplicate data
- **Size limit:** Data URLs limited to ~2MB in some browsers, Blobs have no limit
- **Clean up:** Object URLs can be revoked to free memory

### Download Attribute

The `download` attribute on `<a>` elements forces browser to download file instead of navigating:

```html
<a href="blob:http://..." download="Classic Chrome.xmp">Download</a>
```

**Browser support:** All modern browsers (Chrome 14+, Firefox 20+, Safari 10.1+)

### MIME Types

**Why different MIME types?**
- **NP3:** Binary format → `application/octet-stream` (generic binary)
- **XMP:** XML format → `application/xml` (browsers recognize as structured data)
- **lrtemplate:** Lua text → `text/plain` (plain text file)

**Impact:** Browsers may preview XMP/lrtemplate in tab (if user clicks link), but NP3 always downloads.

### Object URL Lifecycle

**Create:** `URL.createObjectURL(blob)` → returns `blob:http://localhost:8080/[uuid]`

**Use:** Set as `<a href="">` to trigger download

**Revoke:** `URL.revokeObjectURL(url)` → frees memory

**Best practice:** Revoke URLs as soon as download completes (or new conversion starts).

### Filename Sanitization

**Security:** Malicious filenames like `../../etc/passwd.xmp` could cause issues.

**Mitigation:** Browsers automatically sanitize download filenames (remove path separators, restrict characters).

**Recipe approach:** Trust browser sanitization (no custom logic needed).

---

## Follow-Up Stories

**After Story 2-7:**
- Story 2-8: Comprehensive error handling for all failure modes
- Story 2-9: Privacy messaging (reinforce "files never leave device")
- Story 2-10: Responsive design for mobile/tablet

**Future enhancements (not Epic 2):**
- Download history (store last 5 conversions for re-download)
- Download location picker (custom folder)
- Batch download (ZIP multiple conversions)
- Download analytics (track most popular conversions)

---

## References

- **Tech Spec:** `docs/tech-spec-epic-2.md` (Story 2-7 section)
- **PRD:** `docs/PRD.md` (FR-2.7: File Download)
- **Story 2-6:** `docs/stories/2-6-wasm-conversion-execution.md` (conversion data source)
- **Blob API Docs:** https://developer.mozilla.org/en-US/docs/Web/API/Blob
- **Download Attribute Docs:** https://developer.mozilla.org/en-US/docs/Web/HTML/Element/a#attr-download

---

**Story Created:** 2025-11-04
**Story Owner:** Justin (Developer)
**Reviewer:** Bob (Scrum Master)
**Estimated Effort:** 1 day
**Status:** Ready for SM approval → move to "ready-for-dev"
