# Story 2-6: WASM Conversion Execution

**Epic:** Epic 2 - Web Interface (FR-2)
**Story ID:** 2-6
**Status:** drafted
**Created:** 2025-11-04
**Complexity:** Medium (2-3 days)
**Priority:** CRITICAL PATH - Core conversion functionality

---

## User Story

**As a** photographer
**I want** to convert my preset file to a different format
**So that** I can use my presets in different photography software

---

## Business Value

This is the **core deliverable** of Recipe - the actual format conversion. Everything else in Epic 2 is UI scaffolding leading to this moment.

**Success criteria:**
- Conversion works reliably (100% success rate on valid files)
- Fast performance (<100ms for typical presets)
- Clear error messages (if conversion fails)
- Validated output (converted file is usable)

**This story delivers the "magic" - upload file, click button, get converted preset.**

---

## Acceptance Criteria

### AC-1: Call WASM Convert Function

- [ ] Listen for `convertRequest` event (dispatched by Story 2-5)
- [ ] Extract source file data (from Story 2-2)
- [ ] Extract source format (from Story 2-3)
- [ ] Extract target format (from Story 2-5)
- [ ] Call WASM `convert(fileData, sourceFormat, targetFormat)`
- [ ] Handle async Promise (conversion may take 1-100ms)

**Test:**
1. Upload `Classic Chrome.np3`
2. Select XMP target format
3. Click "Convert" button
4. Verify: `convert()` called with correct arguments
5. Verify: Console shows "Converting np3 → xmp..."

### AC-2: Display Conversion Status (Loading → Success/Error)

- [ ] **Before conversion:** Button shows "Convert to [Format]" (enabled)
- [ ] **During conversion:** Button shows "Converting..." (disabled), spinner icon
- [ ] **Success:** Button shows "Converted!" (disabled), checkmark icon
- [ ] **Error:** Button shows "Conversion Failed" (enabled for retry), error icon

**Visual states:**
```
Before:   [Convert to XMP]        (blue, enabled)
During:   [⟳ Converting...]       (gray, disabled)
Success:  [✓ Converted!]          (green, disabled)
Error:    [✗ Conversion Failed]   (red, enabled)
```

**Test:**
1. Click Convert → button changes to "Converting..." with spinner
2. Conversion succeeds → button changes to "Converted!" with checkmark
3. Upload new file → button resets to "Convert to [Format]"

### AC-3: Handle Conversion Success

- [ ] Store converted file data in memory (Uint8Array)
- [ ] Store converted file metadata (size, format)
- [ ] Display success message: "✓ Conversion complete! Your [format] preset is ready."
- [ ] Enable download button (Story 2-7 dependency)
- [ ] Log conversion stats to console (size, time, format)

**Test:**
1. Convert NP3 → XMP
2. Verify: Success message displayed
3. Verify: Console shows "Conversion complete: np3 → xmp ([time]ms, [size] bytes)"
4. Verify: Download button appears (Story 2-7)
5. Verify: `getConvertedFileData()` returns Uint8Array

### AC-4: Handle Conversion Errors (User-Friendly)

- [ ] If `convert()` rejects (parse error, unsupported feature):
  - Display error: "Conversion failed: [user-friendly message]"
  - Show "Try Again" button (re-enable convert button)
  - Log technical error to console (for debugging)
  - Don't crash UI (error is recoverable)

**Error message mapping:**
```
Technical error → User-friendly message
---
"NP3 magic bytes invalid" → "File appears corrupted or not a valid NP3 preset."
"XMP parse error" → "Unable to parse XMP file. File may be corrupted."
"lrtemplate syntax error" → "Invalid Lightroom preset format."
"Unsupported NP3 version" → "NP3 preset version not supported."
```

**Test:**
1. Upload corrupted NP3 file (invalid magic bytes)
2. Click Convert
3. Verify: Error message displayed (user-friendly, not technical)
4. Verify: Console shows technical error (for debugging)
5. Click "Try Again" → Convert button re-enabled
6. Upload valid file → conversion works

### AC-5: Performance Target (<100ms P95)

- [ ] Measure conversion time (performance.now())
- [ ] 95% of conversions complete <100ms
- [ ] Display conversion time in console log
- [ ] No browser UI blocking during conversion

**Performance benchmark:**
- Small files (<50KB): <50ms
- Medium files (50-100KB): <100ms
- Large files (>100KB): <200ms

**Test:**
1. Convert 20 different files (mix of NP3, XMP, lrtemplate)
2. Record conversion times (console logs)
3. Verify: P95 < 100ms
4. Verify: Average time < 50ms
5. Verify: Browser remains responsive during conversion

### AC-6: Memory Management (Clear Old Conversions)

- [ ] Clear previous converted data when new conversion starts
- [ ] Release Uint8Array references (avoid memory leak)
- [ ] Handle rapid conversions (debounce if needed)

**Test:**
1. Upload file → convert to XMP
2. Upload same file → convert to lrtemplate
3. Verify: Only most recent conversion stored in memory
4. Verify: Memory usage stable (DevTools memory profiler)
5. Convert 10 files in succession → memory doesn't grow unbounded

### AC-7: Conversion Validation (Verify Output)

- [ ] After conversion, verify output is valid format:
  - NP3: Starts with "NCP" magic bytes
  - XMP: Valid XML structure (<?xml, crs: namespace)
  - lrtemplate: Valid Lua syntax (s = {)
- [ ] If validation fails, treat as conversion error
- [ ] Log validation result to console

**Test:**
1. Convert NP3 → XMP
2. Verify: Output starts with "<?xml" (valid XMP)
3. Verify: Console shows "Validation: XMP output valid"
4. Convert XMP → NP3
5. Verify: Output starts with "NCP" (valid NP3)

### AC-8: Disable Multiple Concurrent Conversions

- [ ] If conversion is in progress, disable Convert button
- [ ] Ignore additional button clicks during conversion
- [ ] Only one conversion can run at a time

**Test:**
1. Click Convert button
2. Immediately click Convert button again (while first conversion running)
3. Verify: Second click ignored
4. Verify: Only one conversion runs
5. First conversion completes → button re-enabled

---

## Technical Approach

### Conversion Module

**File:** `web/static/converter.js` (new file)

```javascript
// converter.js - WASM conversion wrapper

let convertedFileData = null;
let convertedFileName = null;
let convertedFileFormat = null;
let isConverting = false;

/**
 * Convert preset file using WASM
 * @param {Uint8Array} fileData - Source file data
 * @param {string} sourceFormat - "np3" | "xmp" | "lrtemplate"
 * @param {string} targetFormat - "np3" | "xmp" | "lrtemplate"
 * @param {string} originalFileName - Original file name (for output naming)
 * @returns {Promise<Uint8Array>} Converted file data
 */
export async function convertFile(fileData, sourceFormat, targetFormat, originalFileName) {
    if (!fileData || !sourceFormat || !targetFormat) {
        throw new Error('Missing required parameters');
    }

    if (sourceFormat === targetFormat) {
        throw new Error('Cannot convert to same format');
    }

    if (isConverting) {
        throw new Error('Conversion already in progress');
    }

    // Check if WASM is ready
    if (typeof convert !== 'function') {
        throw new Error('WASM module not loaded');
    }

    console.log(`Converting ${sourceFormat} → ${targetFormat}...`);
    const startTime = performance.now();

    isConverting = true;

    try {
        // Call WASM function (returns Promise<Uint8Array>)
        const outputData = await convert(fileData, sourceFormat, targetFormat);

        const elapsedTime = performance.now() - startTime;
        console.log(`Conversion complete: ${sourceFormat} → ${targetFormat} (${elapsedTime.toFixed(2)}ms, ${outputData.length} bytes)`);

        // Validate output
        validateConvertedData(outputData, targetFormat);

        // Store converted data
        convertedFileData = outputData;
        convertedFileFormat = targetFormat;
        convertedFileName = generateConvertedFileName(originalFileName, targetFormat);

        isConverting = false;

        return outputData;

    } catch (error) {
        isConverting = false;
        console.error('Conversion failed:', error);
        throw new ConversionError(error.message || error);
    }
}

/**
 * Validate converted data matches expected format
 */
function validateConvertedData(data, format) {
    if (!data || data.length === 0) {
        throw new Error('Converted data is empty');
    }

    try {
        switch (format) {
            case 'np3':
                // Check NP3 magic bytes: "NCP" (0x4E 0x43 0x50)
                if (data[0] !== 0x4E || data[1] !== 0x43 || data[2] !== 0x50) {
                    throw new Error('Invalid NP3 magic bytes');
                }
                console.log('Validation: NP3 output valid (magic bytes correct)');
                break;

            case 'xmp':
                // Check XMP starts with XML declaration
                const xmpHeader = new TextDecoder().decode(data.slice(0, 5));
                if (!xmpHeader.startsWith('<?xml')) {
                    throw new Error('Invalid XMP format (missing XML declaration)');
                }
                console.log('Validation: XMP output valid (XML structure correct)');
                break;

            case 'lrtemplate':
                // Check lrtemplate starts with Lua syntax "s = {"
                const lrtemplateHeader = new TextDecoder().decode(data.slice(0, 10)).trim();
                if (!lrtemplateHeader.startsWith('s = {')) {
                    throw new Error('Invalid lrtemplate format (missing Lua syntax)');
                }
                console.log('Validation: lrtemplate output valid (Lua syntax correct)');
                break;

            default:
                throw new Error(`Unknown format: ${format}`);
        }
    } catch (validationError) {
        throw new Error(`Validation failed: ${validationError.message}`);
    }
}

/**
 * Generate output file name based on input file name and target format
 */
function generateConvertedFileName(originalFileName, targetFormat) {
    // Remove original extension
    const baseName = originalFileName.replace(/\.(np3|xmp|lrtemplate)$/i, '');

    // Add new extension
    const extensions = {
        np3: '.np3',
        xmp: '.xmp',
        lrtemplate: '.lrtemplate',
    };

    return `${baseName}${extensions[targetFormat]}`;
}

/**
 * Custom error class for conversion errors
 */
class ConversionError extends Error {
    constructor(message) {
        super(message);
        this.name = 'ConversionError';
        this.userMessage = getUserFriendlyErrorMessage(message);
    }
}

/**
 * Map technical error messages to user-friendly messages
 */
function getUserFriendlyErrorMessage(technicalError) {
    const errorMappings = {
        'NP3 magic bytes invalid': 'File appears corrupted or not a valid NP3 preset.',
        'Invalid NP3 magic bytes': 'File appears corrupted or not a valid NP3 preset.',
        'XMP parse error': 'Unable to parse XMP file. File may be corrupted.',
        'Invalid XMP format': 'Unable to parse XMP file. File may be corrupted.',
        'lrtemplate syntax error': 'Invalid Lightroom preset format.',
        'Invalid lrtemplate format': 'Invalid Lightroom preset format.',
        'Unsupported NP3 version': 'NP3 preset version not supported.',
        'Unsupported XMP version': 'XMP preset version not supported.',
    };

    // Check for exact matches
    for (const [technical, friendly] of Object.entries(errorMappings)) {
        if (technicalError.includes(technical)) {
            return friendly;
        }
    }

    // Default fallback
    return 'Conversion failed. File may be corrupted or unsupported.';
}

/**
 * Get converted file data
 */
export function getConvertedFileData() {
    return convertedFileData;
}

/**
 * Get converted file name
 */
export function getConvertedFileName() {
    return convertedFileName;
}

/**
 * Get converted file format
 */
export function getConvertedFileFormat() {
    return convertedFileFormat;
}

/**
 * Clear converted data from memory
 */
export function clearConvertedData() {
    convertedFileData = null;
    convertedFileName = null;
    convertedFileFormat = null;
    isConverting = false;
}

/**
 * Check if conversion is in progress
 */
export function isConversionInProgress() {
    return isConverting;
}
```

### Integration with Main Flow

**Update `main.js`:**

```javascript
// main.js - Integrate conversion

import { initializeDropZone, handleFile, getCurrentFileData, getCurrentFileName } from './file-handler.js';
import { detectFileFormat } from './format-detector.js';
import { displayParameters } from './parameter-display.js';
import { displayFormatSelector, getSourceFormat, getTargetFormat } from './format-selector.js';
import { convertFile, getConvertedFileData, getConvertedFileName } from './converter.js';
import { initializeWASM } from './wasm-loader.js';

// Initialize WASM
initializeWASM();

// Initialize UI
document.addEventListener('DOMContentLoaded', () => {
    initializeDropZone(handleFile);
});

// Listen for format detected event
window.addEventListener('formatDetected', async (event) => {
    const { format } = event.detail;
    const fileData = getCurrentFileData();

    try {
        // Display parameters (Story 2-4)
        await displayParameters(fileData, format);

        // Display format selector (Story 2-5)
        displayFormatSelector(format);

    } catch (error) {
        console.error('Error:', error);
    }
});

// Listen for convert request event (Story 2-6)
window.addEventListener('convertRequest', async (event) => {
    const { fromFormat, toFormat } = event.detail;

    // Show converting state
    showConvertingState();

    try {
        // Get source file data and name
        const fileData = getCurrentFileData();
        const fileName = getCurrentFileName();

        // Perform conversion
        const convertedData = await convertFile(fileData, fromFormat, toFormat, fileName);

        // Show success state
        showConversionSuccess(toFormat);

        // Enable download button (Story 2-7 will implement download)
        enableDownloadButton();

        // Dispatch conversion complete event
        dispatchConversionCompleteEvent(convertedData, toFormat);

    } catch (error) {
        console.error('Conversion error:', error);
        showConversionError(error);
    }
});

function showConvertingState() {
    const convertButton = document.getElementById('convertButton');
    if (convertButton) {
        convertButton.disabled = true;
        convertButton.innerHTML = '⟳ Converting...';
        convertButton.classList.add('converting');
    }

    // Show status message
    const statusEl = document.getElementById('conversionStatus');
    if (statusEl) {
        statusEl.className = 'status loading';
        statusEl.textContent = 'Converting preset...';
        statusEl.style.display = 'block';
    }
}

function showConversionSuccess(targetFormat) {
    const convertButton = document.getElementById('convertButton');
    if (convertButton) {
        convertButton.disabled = true;
        convertButton.innerHTML = '✓ Converted!';
        convertButton.classList.remove('converting');
        convertButton.classList.add('success');
    }

    // Show success message
    const statusEl = document.getElementById('conversionStatus');
    if (statusEl) {
        statusEl.className = 'status success';
        statusEl.textContent = `✓ Conversion complete! Your ${targetFormat.toUpperCase()} preset is ready.`;
    }
}

function showConversionError(error) {
    const convertButton = document.getElementById('convertButton');
    if (convertButton) {
        convertButton.disabled = false; // Re-enable for retry
        convertButton.innerHTML = '✗ Conversion Failed';
        convertButton.classList.remove('converting');
        convertButton.classList.add('error');
    }

    // Show error message (user-friendly)
    const errorEl = document.getElementById('conversionError');
    if (errorEl) {
        const userMessage = error.userMessage || error.message || 'Conversion failed';
        errorEl.textContent = userMessage;
        errorEl.style.display = 'block';
    }

    // Hide status
    const statusEl = document.getElementById('conversionStatus');
    if (statusEl) {
        statusEl.style.display = 'none';
    }
}

function enableDownloadButton() {
    const downloadButton = document.getElementById('downloadButton');
    if (downloadButton) {
        downloadButton.disabled = false;
        downloadButton.style.display = 'block';
    }
}

function dispatchConversionCompleteEvent(convertedData, format) {
    const event = new CustomEvent('conversionComplete', {
        detail: { convertedData, format }
    });
    window.dispatchEvent(event);
}
```

### CSS Updates

**Add to `web/static/style.css`:**

```css
/* Convert button states */
.convert-button.converting {
    background: #cbd5e0;
    cursor: wait;
}

.convert-button.success {
    background: #48bb78;
    cursor: default;
}

.convert-button.error {
    background: #f56565;
}

.convert-button.error:hover {
    background: #e53e3e;
}

/* Conversion status */
.status {
    padding: 1rem;
    border-radius: 6px;
    margin-top: 1rem;
    font-size: 0.875rem;
    font-weight: 500;
}

.status.loading {
    background: #edf2f7;
    color: #4a5568;
}

.status.success {
    background: #c6f6d5;
    color: #22543d;
}

.status.error {
    background: #fed7d7;
    color: #742a2a;
}

/* Spinner icon animation */
@keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
}

.converting::before {
    content: '⟳';
    display: inline-block;
    animation: spin 1s linear infinite;
}
```

### HTML Updates

**Add to `web/index.html`:**

```html
<!-- Conversion Status -->
<div id="conversionStatus" class="status" style="display: none;" role="status" aria-live="polite"></div>

<!-- Conversion Error -->
<div id="conversionError" class="error-message" style="display: none;" role="alert"></div>

<!-- Download Button (Story 2-7 will implement) -->
<button id="downloadButton" class="download-button" style="display: none;" disabled>
    Download Converted Preset
</button>
```

---

## Dependencies

### Required Before Starting

- ✅ Story 2-2 complete (file data available)
- ✅ Story 2-3 complete (format detected)
- ✅ Story 2-5 complete (target format selected)
- ✅ WASM convert() function implemented (cmd/wasm/main.go)

### Blocks These Stories

- Story 2-7 (File Download) - needs converted file data

---

## Testing Plan

### Manual Testing

**Test Case 1: NP3 → XMP Conversion**
1. Upload `Classic Chrome.np3`
2. Select XMP target format
3. Click "Convert to XMP"
4. Verify: Button changes to "Converting..." (disabled)
5. Verify: Conversion completes <100ms
6. Verify: Button changes to "Converted!" (green, disabled)
7. Verify: Success message displayed
8. Verify: Console shows conversion stats (time, size)
9. Verify: Download button appears

**Test Case 2: XMP → NP3 Conversion**
1. Upload `Classic Chrome.xmp`
2. Select NP3 target format
3. Click "Convert to NP3"
4. Verify: Conversion succeeds
5. Verify: Output validated (starts with "NCP" magic bytes)
6. Verify: Console shows "Validation: NP3 output valid"

**Test Case 3: lrtemplate → XMP Conversion**
1. Upload `auto tone.lrtemplate`
2. Select XMP target format
3. Click "Convert to XMP"
4. Verify: Conversion succeeds
5. Verify: Output validated (starts with "<?xml")
6. Verify: Console shows "Validation: XMP output valid"

**Test Case 4: Conversion Error Handling**
1. Upload corrupted NP3 file (modify first 3 bytes to invalidate magic bytes)
2. Select XMP target format
3. Click "Convert to XMP"
4. Verify: Conversion fails
5. Verify: Error message displayed (user-friendly, not technical)
6. Verify: Console shows technical error (for debugging)
7. Verify: Button shows "✗ Conversion Failed" (red, enabled for retry)
8. Upload valid file → click "Convert to XMP" → conversion succeeds

**Test Case 5: Multiple Conversions (Same File)**
1. Upload NP3 file → convert to XMP (success)
2. Without uploading new file, convert to lrtemplate (should fail - need to re-upload)
3. Upload same NP3 file again → convert to lrtemplate (success)
4. Verify: Each conversion clears previous converted data
5. Verify: Memory usage stable

**Test Case 6: Rapid Conversions**
1. Upload file → click Convert button
2. Immediately click Convert button again (while first conversion running)
3. Verify: Second click ignored
4. Verify: Only one conversion runs
5. First conversion completes → button re-enabled

**Test Case 7: Performance Benchmark**
1. Convert 20 different files (mix of NP3, XMP, lrtemplate)
2. Record conversion times (console logs)
3. Verify: P95 < 100ms
4. Verify: Average time < 50ms
5. Verify: All files convert successfully

**Test Case 8: Memory Management**
1. Upload file → convert (store converted data in memory)
2. Upload new file → convert (should clear previous converted data)
3. Repeat 10 times
4. Verify: Memory usage stable (DevTools memory profiler)
5. Verify: Only most recent conversion stored

### Automated Testing (Optional for MVP)

```javascript
// Unit test for conversion

import { convertFile } from './converter.js';

// Mock WASM function
global.convert = async (data, fromFormat, toFormat) => {
    // Simulate conversion delay
    await new Promise(resolve => setTimeout(resolve, 10));

    if (fromFormat === 'np3' && toFormat === 'xmp') {
        // Return mock XMP data
        const xmpData = '<?xml version="1.0"?><x:xmpmeta>...</x:xmpmeta>';
        return new TextEncoder().encode(xmpData);
    }

    throw new Error('Unsupported conversion');
};

// Test successful conversion
const np3Data = new Uint8Array([0x4E, 0x43, 0x50, /* ... */]);
try {
    const xmpData = await convertFile(np3Data, 'np3', 'xmp', 'test.np3');
    console.assert(xmpData.length > 0, 'Conversion output empty');
    console.log('✓ Conversion test passed');
} catch (error) {
    console.error('✗ Conversion test failed:', error);
}

// Test error handling
try {
    await convertFile(np3Data, 'np3', 'np3', 'test.np3'); // Same format
    console.error('✗ Should have thrown error for same-format conversion');
} catch (error) {
    console.log('✓ Error handling test passed');
}
```

### Browser Compatibility

Test in:
- ✅ Chrome (latest) - WASM, Promises fully supported
- ✅ Firefox (latest) - WASM, Promises fully supported
- ✅ Safari (latest) - WASM, Promises fully supported

**Expected:** Identical behavior across browsers.

---

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Conversion works for all format combinations (9 total: 3×3 minus 3 same-format)
- [ ] Performance target met (<100ms P95)
- [ ] Error handling tested with corrupted files
- [ ] Output validation works for all formats
- [ ] Memory management verified (no leaks)
- [ ] Manual testing completed in Chrome, Firefox, Safari
- [ ] Code reviewed
- [ ] Integration with Stories 2-2, 2-3, 2-5 verified
- [ ] Story marked "ready-for-dev" in sprint status

---

## Out of Scope

**Explicitly NOT in this story:**
- ❌ File download (Story 2-7)
- ❌ Batch conversion (multiple files at once - Epic 3)
- ❌ Conversion options (quality settings, metadata preservation - future)

**This story only delivers:** Format conversion - transform file from one format to another using WASM.

---

## Technical Notes

### Why Validate Output?

**Alternative:** Trust WASM converter always produces valid output

**Decision:** Validate output after conversion

**Rationale:**
- **Catch WASM bugs:** If Epic 1 converter has bugs, validation detects them
- **User confidence:** "Your preset is ready" is more credible if we checked it
- **Debug aid:** If validation fails, we know issue is in converter (not parsing)
- **Performance:** Validation is fast (~1ms) - check first 10 bytes

### Error Message Philosophy

**Principle:** Show user-friendly errors to users, technical errors to console

**Example:**
```javascript
// Technical error (for developers):
console.error('NP3 parse error: invalid magic bytes at offset 0x00')

// User-friendly error (for photographers):
showError('File appears corrupted or not a valid NP3 preset.')
```

**Why?** Users don't care about "magic bytes" or "offset 0x00" - they just want to know if they uploaded the wrong file.

### Performance Expectations

**Native Go (Epic 1):** 0.002-0.067ms per conversion

**WASM target (Epic 2):** <100ms per conversion

**Overhead:** 10-100x slower than native (acceptable for browser use)

**Why so much slower?**
- WASM→Go interop overhead
- JavaScript Promise overhead
- Browser memory management
- No optimization flags in WASM build (yet)

**Still fast enough:** 100ms is imperceptible to users (<200ms is "instant" in UX research)

### Memory Management Strategy

**Problem:** Each conversion creates large Uint8Arrays (50-100KB). Holding multiple conversions risks memory exhaustion.

**Solution:** Clear previous conversion when starting new one

**Implementation:**
```javascript
// Before new conversion
clearConvertedData();

// After conversion
convertedFileData = outputData; // Store new conversion
```

**Result:** Only one conversion in memory at a time.

### Same-Format Conversion

**Why prevent?** Converting XMP→XMP is a no-op (output = input). Wastes computation and confuses users ("I converted but nothing changed!").

**Implementation:** Check `sourceFormat === targetFormat` before conversion, throw error.

---

## Follow-Up Stories

**After Story 2-6:**
- Story 2-7: Download converted file with correct extension and filename
- Story 2-8: Comprehensive error handling (network errors, WASM crashes)

**Future enhancements (not Epic 2):**
- Conversion preview (show converted parameters before download)
- Batch conversion (convert multiple files at once)
- Conversion options (metadata preservation, quality settings)
- Conversion history (store last 5 conversions for re-download)

---

## References

- **Tech Spec:** `docs/tech-spec-epic-2.md` (Story 2-6 section)
- **PRD:** `docs/PRD.md` (FR-2.6: WASM Conversion)
- **WASM Entry Point:** `cmd/wasm/main.go` (convert() function)
- **Epic 1 Converter:** `internal/converter/converter.go` (conversion logic)
- **Story 2-2:** `docs/stories/2-2-file-upload-handling.md` (file data source)
- **Story 2-3:** `docs/stories/2-3-format-detection.md` (format detection)
- **Story 2-5:** `docs/stories/2-5-target-format-selection.md` (target format)

---

**Story Created:** 2025-11-04
**Story Owner:** Justin (Developer)
**Reviewer:** Bob (Scrum Master)
**Estimated Effort:** 2-3 days
**Priority:** CRITICAL PATH
**Status:** Ready for SM approval → move to "ready-for-dev"
