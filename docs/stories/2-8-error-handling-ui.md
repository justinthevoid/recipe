# Story 2-8: Error Handling UI

**Epic:** Epic 2 - Web Interface (FR-2)
**Story ID:** 2-8
**Status:** drafted
**Created:** 2025-11-04
**Complexity:** Medium (1-2 days)

---

## User Story

**As a** photographer
**I want** clear error messages when something goes wrong
**So that** I know what happened and how to fix it

---

## Business Value

Error handling is often overlooked, but it's **critical for user trust**. When Recipe fails, users need:
- **Clear explanation** of what went wrong (not technical jargon)
- **Actionable next steps** (retry? upload different file? check browser settings?)
- **Confidence** that their data is safe (no server uploads, no data loss)

**Bad error:** "TypeError: Cannot read property 'length' of undefined at line 42"

**Good error:** "Unable to read file. Please try uploading again."

**This story ensures Recipe fails gracefully and guides users to success.**

---

## Acceptance Criteria

### AC-1: Centralized Error Display Component

- [ ] Create single error display component for all error types
- [ ] Location: Below page header, above main content (always visible)
- [ ] Style: Red background, white text, icon, dismissible
- [ ] Contains:
  - Error icon (⚠️ or ✗)
  - User-friendly error message
  - Technical details (collapsible, hidden by default)
  - Action buttons (Retry, Reset, Help)

**Visual design:**
```
┌─────────────────────────────────────────────────┐
│ ⚠️ Unable to read file. Please try again.       │
│                                                  │
│ [Show Details ▼] [Try Again] [Reset] [Help]     │
└─────────────────────────────────────────────────┘
```

**Test:**
1. Trigger any error (e.g., upload invalid file)
2. Verify: Error component appears at top of page
3. Verify: Red background, clear message, action buttons
4. Click "Show Details" → technical error expands
5. Click "×" (dismiss button) → error component hides

### AC-2: Comprehensive Error Message Library

- [ ] Define user-friendly messages for all error types:
  - **WASM loading failure** (Story 2-1)
  - **Invalid file type** (Story 2-1)
  - **File too large** (Story 2-2)
  - **File read error** (Story 2-2)
  - **Format detection failure** (Story 2-3)
  - **Parameter extraction failure** (Story 2-4)
  - **Conversion failure** (Story 2-6)
  - **Download failure** (Story 2-7)
  - **Browser compatibility issues**
  - **Network errors** (WASM fetch)

**Error message format:**
```
[User-Friendly Message]
What went wrong: [Simple explanation]
What to try: [Actionable next step]

[Technical Details] (collapsible)
Error: [Technical error message]
Stack: [Stack trace if available]
```

**Test:**
1. Review all error messages in library
2. Verify: No technical jargon in user-facing messages
3. Verify: All messages include actionable next step
4. Verify: Technical details available for debugging

### AC-3: Error Recovery Actions

- [ ] **Try Again:** Re-attempt last action (re-upload, re-convert, re-download)
- [ ] **Reset:** Clear all data and return to default state
- [ ] **Help:** Link to troubleshooting guide or FAQ
- [ ] **Dismiss:** Close error message (but keep data/state)

**Action behavior:**
```
Try Again → Re-run last operation
Reset     → Clear file data, hide all panels, return to drop zone
Help      → Open FAQ page in new tab
Dismiss   → Hide error message (allow user to continue)
```

**Test:**
1. Trigger conversion error → click "Try Again"
2. Verify: Conversion re-attempted with same data
3. Click "Reset" → verify: UI returns to default state (drop zone visible, no file data)
4. Click "Help" → verify: FAQ page opens in new tab
5. Click "Dismiss" → verify: Error message hides, user can continue

### AC-4: Error Logging (Console + Optional Telemetry)

- [ ] Log all errors to console with context:
  - Error type
  - Error message (technical + user-friendly)
  - Timestamp
  - User action that triggered error
  - File metadata (if applicable)
- [ ] Optional: Send anonymized error telemetry (MVP: just console logs)
- [ ] Never log sensitive data (file contents, user info)

**Console log format:**
```javascript
[2025-11-04 14:32:15] Recipe Error: ConversionError
Message: Unable to parse XMP file. File may be corrupted.
Technical: XMP parse error: unexpected token at line 42
Action: User clicked Convert button
File: Classic Chrome.xmp (15234 bytes)
```

**Test:**
1. Trigger any error
2. Open DevTools console
3. Verify: Error logged with full context
4. Verify: No sensitive data in logs (no file contents)
5. Verify: Timestamp and user action logged

### AC-5: Error Boundaries (Prevent Full UI Crash)

- [ ] Wrap critical components in error boundaries
- [ ] If component crashes, show fallback UI (not white screen)
- [ ] Fallback UI:
  - "Something went wrong"
  - "Please refresh the page"
  - "If problem persists, contact support"
- [ ] Log component crash to console

**Test:**
1. Simulate component crash (throw error in component render)
2. Verify: Fallback UI displays (not white screen or browser error)
3. Verify: Error logged to console
4. Refresh page → verify: App loads normally

### AC-6: Browser Compatibility Error Handling

- [ ] Detect unsupported browsers (IE11, old Chrome <90, etc.)
- [ ] Show message: "Recipe requires a modern browser. Please upgrade to Chrome, Firefox, or Safari."
- [ ] Detect missing APIs:
  - WebAssembly not supported
  - FileReader not supported
  - Blob API not supported
- [ ] Show specific error for missing API

**Test:**
1. Test in IE11 (or simulate with DevTools)
2. Verify: "Unsupported browser" message displayed
3. Test in browser with WASM disabled
4. Verify: "WebAssembly not supported" message displayed

### AC-7: Network Error Handling (WASM Loading)

- [ ] If WASM fails to load (network error, CDN down):
  - Show error: "Unable to load converter. Please check your internet connection."
  - Provide retry button
  - Log error to console
- [ ] Handle slow loading (>5s):
  - Show loading message: "Loading converter... (this may take a moment)"
  - Don't assume failure immediately

**Test:**
1. Simulate network error (DevTools → offline mode)
2. Load page
3. Verify: "Unable to load converter" error displayed
4. Click "Retry" → verify: WASM loading re-attempted
5. Re-enable network → retry → verify: WASM loads successfully

### AC-8: User Testing (Error Message Clarity)

- [ ] Test error messages with non-technical users (photographers)
- [ ] Verify messages are understandable without technical knowledge
- [ ] Verify action buttons are clear
- [ ] Iterate based on feedback

**Test:**
1. Show error messages to 3 non-technical users
2. Ask: "What do you think went wrong?" (should match error message)
3. Ask: "What would you do next?" (should match action buttons)
4. Iterate messages if confusion

---

## Technical Approach

### Centralized Error Handler

**File:** `web/static/error-handler.js` (new file)

```javascript
// error-handler.js - Centralized error handling

const ERROR_MESSAGES = {
    // WASM loading errors
    'wasm-load-failed': {
        title: 'Unable to Load Converter',
        message: 'Recipe couldn\'t load the conversion engine.',
        reason: 'Your internet connection may be unstable, or your browser doesn\'t support WebAssembly.',
        action: 'Check your internet connection and try refreshing the page. If the problem persists, try a different browser (Chrome, Firefox, or Safari).',
        recovery: ['retry', 'help'],
    },

    // File upload errors
    'invalid-file-type': {
        title: 'Invalid File Type',
        message: 'This file type isn\'t supported.',
        reason: 'Recipe only converts NP3, XMP, and lrtemplate preset files.',
        action: 'Please upload a valid preset file (.np3, .xmp, or .lrtemplate).',
        recovery: ['reset', 'help'],
    },

    'file-too-large': {
        title: 'File Too Large',
        message: 'This file exceeds the 10MB size limit.',
        reason: 'Preset files are typically <100KB. This file may be corrupted or not a preset.',
        action: 'Please check you\'ve uploaded the correct file.',
        recovery: ['reset', 'help'],
    },

    'file-read-error': {
        title: 'Unable to Read File',
        message: 'Recipe couldn\'t read your file.',
        reason: 'The file may be corrupted, or your browser blocked access.',
        action: 'Try uploading the file again. If the problem persists, try a different file.',
        recovery: ['retry', 'reset'],
    },

    // Format detection errors
    'format-detection-failed': {
        title: 'Unknown Format',
        message: 'Recipe couldn\'t identify this file\'s format.',
        reason: 'The file may be corrupted, or it may not be a valid preset.',
        action: 'Check you\'ve uploaded the correct file. Valid formats: NP3, XMP, lrtemplate.',
        recovery: ['reset', 'help'],
    },

    // Parameter extraction errors
    'parameter-extraction-failed': {
        title: 'Unable to Read Parameters',
        message: 'Recipe couldn\'t extract parameters from this file.',
        reason: 'The file may be corrupted or use an unsupported preset version.',
        action: 'You can still try converting the file - conversion may work even if parameter preview doesn\'t.',
        recovery: ['continue', 'reset'],
    },

    // Conversion errors
    'conversion-failed': {
        title: 'Conversion Failed',
        message: 'Recipe couldn\'t convert your preset.',
        reason: 'The file may be corrupted, or it may use unsupported features.',
        action: 'Try uploading a different preset, or check the file is valid.',
        recovery: ['retry', 'reset', 'help'],
    },

    // Download errors
    'download-failed': {
        title: 'Download Failed',
        message: 'Recipe couldn\'t download your converted preset.',
        reason: 'Your browser may have blocked the download, or there\'s not enough disk space.',
        action: 'Check your browser\'s download settings and try again.',
        recovery: ['retry', 'help'],
    },

    // Browser compatibility errors
    'browser-unsupported': {
        title: 'Unsupported Browser',
        message: 'Recipe requires a modern browser.',
        reason: 'Your browser doesn\'t support WebAssembly, which Recipe needs to convert presets.',
        action: 'Please upgrade to Chrome, Firefox, or Safari (latest version).',
        recovery: ['help'],
    },

    // Network errors
    'network-error': {
        title: 'Network Error',
        message: 'Recipe couldn\'t connect to the server.',
        reason: 'Your internet connection may be unstable.',
        action: 'Check your internet connection and try refreshing the page.',
        recovery: ['retry'],
    },

    // Generic fallback
    'unknown-error': {
        title: 'Something Went Wrong',
        message: 'Recipe encountered an unexpected error.',
        reason: 'This may be a bug, or your browser may not be supported.',
        action: 'Try refreshing the page. If the problem persists, please report this issue on GitHub.',
        recovery: ['retry', 'reset', 'help'],
    },
};

/**
 * Display error message
 * @param {string} errorType - Error type key from ERROR_MESSAGES
 * @param {Error} error - Original error object (for technical details)
 */
export function showError(errorType, error = null) {
    const errorData = ERROR_MESSAGES[errorType] || ERROR_MESSAGES['unknown-error'];

    // Log to console
    logError(errorType, errorData, error);

    // Display in UI
    renderErrorUI(errorData, error);
}

/**
 * Render error UI component
 */
function renderErrorUI(errorData, technicalError) {
    const container = document.getElementById('errorContainer');
    if (!container) {
        console.error('Error container not found');
        return;
    }

    let html = `
        <div class="error-panel" role="alert">
            <div class="error-header">
                <span class="error-icon">⚠️</span>
                <h3 class="error-title">${errorData.title}</h3>
                <button class="error-dismiss" aria-label="Dismiss error">×</button>
            </div>
            <div class="error-body">
                <p class="error-message"><strong>${errorData.message}</strong></p>
                <p class="error-reason">${errorData.reason}</p>
                <p class="error-action">
                    <strong>What to try:</strong> ${errorData.action}
                </p>
            </div>
    `;

    // Technical details (collapsible)
    if (technicalError) {
        html += `
            <div class="error-details">
                <button class="error-details-toggle" id="errorDetailsToggle">
                    Show Technical Details ▼
                </button>
                <div class="error-details-content" id="errorDetailsContent" style="display: none;">
                    <pre>${escapeHtml(technicalError.toString())}</pre>
                    ${technicalError.stack ? `<pre>${escapeHtml(technicalError.stack)}</pre>` : ''}
                </div>
            </div>
        `;
    }

    // Recovery actions
    html += `
            <div class="error-actions">
    `;

    for (const action of errorData.recovery) {
        const actionButtons = {
            'retry': '<button class="error-action-btn retry">Try Again</button>',
            'reset': '<button class="error-action-btn reset">Reset</button>',
            'continue': '<button class="error-action-btn continue">Continue Anyway</button>',
            'help': '<a href="https://github.com/justin/recipe#troubleshooting" target="_blank" class="error-action-btn help">Get Help</a>',
        };
        html += actionButtons[action] || '';
    }

    html += `
            </div>
        </div>
    `;

    container.innerHTML = html;
    container.style.display = 'block';

    // Attach event listeners
    attachErrorListeners();
}

/**
 * Attach event listeners to error UI
 */
function attachErrorListeners() {
    // Dismiss button
    const dismissBtn = document.querySelector('.error-dismiss');
    if (dismissBtn) {
        dismissBtn.addEventListener('click', hideError);
    }

    // Details toggle
    const detailsToggle = document.getElementById('errorDetailsToggle');
    if (detailsToggle) {
        detailsToggle.addEventListener('click', toggleErrorDetails);
    }

    // Action buttons
    const retryBtn = document.querySelector('.error-action-btn.retry');
    if (retryBtn) {
        retryBtn.addEventListener('click', handleRetry);
    }

    const resetBtn = document.querySelector('.error-action-btn.reset');
    if (resetBtn) {
        resetBtn.addEventListener('click', handleReset);
    }

    const continueBtn = document.querySelector('.error-action-btn.continue');
    if (continueBtn) {
        continueBtn.addEventListener('click', hideError);
    }
}

/**
 * Toggle technical details visibility
 */
function toggleErrorDetails() {
    const toggle = document.getElementById('errorDetailsToggle');
    const content = document.getElementById('errorDetailsContent');

    if (content.style.display === 'none') {
        content.style.display = 'block';
        toggle.textContent = 'Hide Technical Details ▲';
    } else {
        content.style.display = 'none';
        toggle.textContent = 'Show Technical Details ▼';
    }
}

/**
 * Handle retry action
 */
function handleRetry() {
    hideError();
    // Dispatch retry event for last action
    const event = new CustomEvent('errorRetry');
    window.dispatchEvent(event);
}

/**
 * Handle reset action
 */
function handleReset() {
    hideError();
    // Dispatch reset event
    const event = new CustomEvent('errorReset');
    window.dispatchEvent(event);
}

/**
 * Hide error UI
 */
export function hideError() {
    const container = document.getElementById('errorContainer');
    if (container) {
        container.style.display = 'none';
        container.innerHTML = '';
    }
}

/**
 * Log error to console
 */
function logError(errorType, errorData, technicalError) {
    const timestamp = new Date().toISOString();
    console.error(`[${timestamp}] Recipe Error: ${errorType}`);
    console.error('Message:', errorData.message);
    console.error('Reason:', errorData.reason);
    if (technicalError) {
        console.error('Technical:', technicalError);
    }

    // Optional: Send to telemetry service (not implemented in MVP)
    // sendErrorTelemetry(errorType, errorData, technicalError);
}

/**
 * Escape HTML to prevent XSS
 */
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

/**
 * Check browser compatibility
 */
export function checkBrowserCompatibility() {
    // Check WebAssembly support
    if (typeof WebAssembly === 'undefined') {
        showError('browser-unsupported', new Error('WebAssembly not supported'));
        return false;
    }

    // Check FileReader support
    if (typeof FileReader === 'undefined') {
        showError('browser-unsupported', new Error('FileReader API not supported'));
        return false;
    }

    // Check Blob support
    if (typeof Blob === 'undefined') {
        showError('browser-unsupported', new Error('Blob API not supported'));
        return false;
    }

    return true;
}
```

### Integration with Main Flow

**Update `main.js`:**

```javascript
// main.js - Integrate error handling

import { checkBrowserCompatibility, showError, hideError } from './error-handler.js';
import { initializeDropZone } from './file-handler.js';
import { initializeWASM } from './wasm-loader.js';

// Check browser compatibility on load
if (!checkBrowserCompatibility()) {
    // Error already displayed by checkBrowserCompatibility()
    throw new Error('Browser not supported');
}

// Initialize WASM with error handling
try {
    await initializeWASM();
} catch (error) {
    console.error('WASM initialization failed:', error);
    showError('wasm-load-failed', error);
}

// Initialize UI
document.addEventListener('DOMContentLoaded', () => {
    initializeDropZone();
});

// Listen for error reset event
window.addEventListener('errorReset', () => {
    // Clear all application state
    // Reset UI to default
    location.reload(); // Simple approach: just reload page
});

// Listen for error retry event
window.addEventListener('errorRetry', () => {
    // Re-attempt last action
    // Implementation depends on what failed
    console.log('Retry requested');
});

// Wrap critical operations in try-catch
window.addEventListener('formatDetected', async (event) => {
    try {
        await displayParameters(fileData, format);
        displayFormatSelector(format);
    } catch (error) {
        showError('parameter-extraction-failed', error);
    }
});

window.addEventListener('convertRequest', async (event) => {
    try {
        await convertFile(fileData, fromFormat, toFormat, fileName);
    } catch (error) {
        showError('conversion-failed', error);
    }
});
```

### CSS for Error UI

**Add to `web/static/style.css`:**

```css
/* Error panel */
.error-panel {
    position: fixed;
    top: 1rem;
    left: 50%;
    transform: translateX(-50%);
    z-index: 1000;
    max-width: 600px;
    width: 90%;
    background: #fed7d7;
    border: 2px solid #f56565;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.error-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1rem;
}

.error-icon {
    font-size: 1.5rem;
}

.error-title {
    flex: 1;
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: #742a2a;
}

.error-dismiss {
    background: none;
    border: none;
    font-size: 1.5rem;
    color: #742a2a;
    cursor: pointer;
    padding: 0.25rem;
    line-height: 1;
}

.error-dismiss:hover {
    color: #c53030;
}

.error-body {
    margin-bottom: 1rem;
}

.error-message {
    margin: 0 0 0.5rem 0;
    color: #742a2a;
    font-size: 1rem;
}

.error-reason {
    margin: 0 0 0.5rem 0;
    color: #9b2c2c;
    font-size: 0.875rem;
}

.error-action {
    margin: 0;
    color: #742a2a;
    font-size: 0.875rem;
}

/* Technical details */
.error-details {
    margin-bottom: 1rem;
    border-top: 1px solid #fc8181;
    padding-top: 1rem;
}

.error-details-toggle {
    background: none;
    border: none;
    color: #742a2a;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    padding: 0;
}

.error-details-toggle:hover {
    color: #c53030;
    text-decoration: underline;
}

.error-details-content {
    margin-top: 0.5rem;
    padding: 0.75rem;
    background: #fff5f5;
    border-radius: 4px;
    font-family: monospace;
    font-size: 0.75rem;
    color: #742a2a;
    max-height: 200px;
    overflow-y: auto;
}

.error-details-content pre {
    margin: 0;
    white-space: pre-wrap;
    word-break: break-word;
}

/* Action buttons */
.error-actions {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
}

.error-action-btn {
    padding: 0.5rem 1rem;
    border: 1px solid #742a2a;
    border-radius: 4px;
    background: #fff;
    color: #742a2a;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    text-decoration: none;
    display: inline-block;
    transition: all 0.2s ease;
}

.error-action-btn:hover {
    background: #fff5f5;
    border-color: #c53030;
    color: #c53030;
}
```

### HTML Updates

**Add to `web/index.html`:**

```html
<!-- Error Container (fixed position at top) -->
<div id="errorContainer" style="display: none;"></div>
```

---

## Dependencies

### Required Before Starting

- ✅ Stories 2-1 through 2-7 complete (all error scenarios defined)

### No Blocking Dependencies

Story 2-8 is a cross-cutting concern that enhances all previous stories.

---

## Testing Plan

### Manual Testing

**Test each error type:**

1. **WASM Load Failure:**
   - Simulate: DevTools → Network → Block `recipe.wasm`
   - Verify: "Unable to Load Converter" error displayed
   - Verify: Recovery actions: Retry, Help

2. **Invalid File Type:**
   - Upload `.jpg` file
   - Verify: "Invalid File Type" error
   - Verify: Recovery actions: Reset, Help

3. **File Too Large:**
   - Create 15MB file
   - Verify: "File Too Large" error

4. **File Read Error:**
   - Upload file, disconnect drive mid-read
   - Verify: "Unable to Read File" error

5. **Format Detection Failure:**
   - Upload corrupted file (random bytes)
   - Verify: "Unknown Format" error

6. **Parameter Extraction Failure:**
   - Upload corrupted XMP (invalid XML)
   - Verify: "Unable to Read Parameters" error
   - Verify: "Continue Anyway" button allows conversion

7. **Conversion Failure:**
   - Upload corrupted NP3 (invalid magic bytes)
   - Verify: "Conversion Failed" error

8. **Download Failure:**
   - Block downloads in browser
   - Verify: "Download Failed" error

9. **Browser Unsupported:**
   - Test in IE11 or old browser
   - Verify: "Unsupported Browser" error

10. **Network Error:**
    - Disconnect internet → load page
    - Verify: "Network Error" error

**Test error UI:**
1. Trigger any error
2. Verify: Error panel appears at top (fixed position)
3. Click "Show Technical Details" → details expand
4. Click "Hide Technical Details" → details collapse
5. Click "×" (dismiss) → error hides
6. Click recovery action → appropriate action taken

---

## Definition of Done

- [ ] All acceptance criteria met
- [ ] All 10 error types have user-friendly messages
- [ ] Error UI tested across all error scenarios
- [ ] Recovery actions work (retry, reset, help)
- [ ] Technical details logged to console
- [ ] Browser compatibility check implemented
- [ ] User testing completed (messages understandable)
- [ ] Manual testing in Chrome, Firefox, Safari
- [ ] Code reviewed
- [ ] Story marked "ready-for-dev" in sprint status

---

## Out of Scope

**Explicitly NOT in this story:**
- ❌ Error telemetry service (anonymized error reporting - future)
- ❌ A/B testing error messages (future optimization)
- ❌ Multilingual error messages (English only for MVP)

**This story only delivers:** Comprehensive error handling UI with user-friendly messages and recovery actions.

---

## References

- **Tech Spec:** `docs/tech-spec-epic-2.md` (Story 2-8 section)
- **PRD:** `docs/PRD.md` (FR-2.8: Error Handling)
- **Stories 2-1 through 2-7:** All error scenarios
- **UX Research:** Nielsen Norman Group - Error Message Guidelines

---

**Story Created:** 2025-11-04
**Story Owner:** Justin (Developer)
**Reviewer:** Bob (Scrum Master)
**Estimated Effort:** 1-2 days
**Status:** Ready for SM approval → move to "ready-for-dev"
