// error-handler.js - Centralized error handling for Recipe
// Story 2-8: Error Handling UI

/**
 * Comprehensive error message library
 * Maps error type keys to user-friendly error objects
 */
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
 * Display error message in UI
 * @param {string} errorType - Error type key from ERROR_MESSAGES
 * @param {Error} error - Original error object (for technical details)
 */
export function showError(errorType, error = null) {
    const errorData = ERROR_MESSAGES[errorType] || ERROR_MESSAGES['unknown-error'];

    // Log to console with full context
    logError(errorType, errorData, error);

    // Display in UI
    renderErrorUI(errorData, error);
}

/**
 * Hide error UI component
 */
export function hideError() {
    const container = document.getElementById('errorContainer');
    if (container) {
        container.style.display = 'none';
        container.innerHTML = '';
    }
}

/**
 * Check browser compatibility on initialization
 * Story 6-4 (AC-7): Enhanced browser detection with dedicated unsupported browser screen
 * @returns {boolean} True if browser is supported, false otherwise
 */
export function checkBrowserCompatibility() {
    // Check WebAssembly support
    const hasWasm = typeof WebAssembly !== 'undefined';

    // Check FileReader API support
    const hasFileReader = typeof FileReader !== 'undefined';

    // Check Blob support
    const hasBlob = typeof Blob !== 'undefined';

    // If any required feature is missing, show unsupported browser screen
    if (!hasWasm || !hasFileReader || !hasBlob) {
        showUnsupportedBrowserMessage();
        return false;
    }

    return true;
}

/**
 * Show dedicated unsupported browser message screen
 * Story 6-4 (AC-7): Replace entire app content with clear unsupported browser message
 * @private
 */
function showUnsupportedBrowserMessage() {
    const appContainer = document.getElementById('app') || document.body;
    appContainer.innerHTML = `
        <div class="unsupported-browser">
            <h1>Unsupported Browser</h1>
            <p class="message">Recipe requires a modern browser with WebAssembly support.</p>
            <p class="info">Please use one of the following browsers:</p>
            <ul class="browser-list">
                <li>Chrome (version 131 or newer)</li>
                <li>Firefox (version 132 or newer)</li>
                <li>Safari (version 18.0 or newer)</li>
                <li>Edge (version 131 or newer)</li>
            </ul>
            <p class="download-links">
                <a href="https://www.google.com/chrome/" target="_blank" rel="noopener noreferrer">Download Chrome</a> |
                <a href="https://www.mozilla.org/firefox/" target="_blank" rel="noopener noreferrer">Download Firefox</a> |
                <a href="https://www.microsoft.com/edge" target="_blank" rel="noopener noreferrer">Download Edge</a>
            </p>
            <p class="technical-info">
                <strong>Technical Details:</strong><br>
                Your browser is missing required features for Recipe:
                ${typeof WebAssembly === 'undefined' ? '✗ WebAssembly support' : '✓ WebAssembly support'}<br>
                ${typeof FileReader === 'undefined' ? '✗ FileReader API support' : '✓ FileReader API support'}<br>
                ${typeof Blob === 'undefined' ? '✗ Blob API support' : '✓ Blob API support'}
            </p>
        </div>
    `;

    // Log to console
    console.error('Browser not supported:', {
        hasWebAssembly: typeof WebAssembly !== 'undefined',
        hasFileReader: typeof FileReader !== 'undefined',
        hasBlob: typeof Blob !== 'undefined',
        userAgent: navigator.userAgent
    });
}

/**
 * Render error UI component
 * @private
 */
function renderErrorUI(errorData, technicalError) {
    const container = document.getElementById('errorContainer');
    if (!container) {
        console.error('Error container not found in DOM');
        return;
    }

    let html = `
        <div class="error-panel" role="alert" aria-live="assertive">
            <div class="error-header">
                <span class="error-icon" aria-hidden="true">⚠️</span>
                <h3 class="error-title">${escapeHtml(errorData.title)}</h3>
                <button class="error-dismiss" aria-label="Dismiss error" type="button">×</button>
            </div>
            <div class="error-body">
                <p class="error-message"><strong>${escapeHtml(errorData.message)}</strong></p>
                <p class="error-reason">${escapeHtml(errorData.reason)}</p>
                <p class="error-action">
                    <strong>What to try:</strong> ${escapeHtml(errorData.action)}
                </p>
            </div>
    `;

    // Technical details (collapsible)
    if (technicalError) {
        const errorStack = technicalError.stack || '';
        html += `
            <div class="error-details">
                <button class="error-details-toggle" id="errorDetailsToggle" type="button">
                    Show Technical Details ▼
                </button>
                <div class="error-details-content" id="errorDetailsContent" style="display: none;">
                    <pre>${escapeHtml(technicalError.toString())}</pre>
                    ${errorStack ? `<pre>${escapeHtml(errorStack)}</pre>` : ''}
                </div>
            </div>
        `;
    }

    // Recovery actions
    html += `
            <div class="error-actions">
    `;

    const actionButtons = {
        'retry': '<button class="error-action-btn retry" type="button">Try Again</button>',
        'reset': '<button class="error-action-btn reset" type="button">Reset</button>',
        'continue': '<button class="error-action-btn continue" type="button">Continue Anyway</button>',
        'help': '<a href="https://github.com/justin/recipe#troubleshooting" target="_blank" rel="noopener noreferrer" class="error-action-btn help">Get Help</a>',
    };

    for (const action of errorData.recovery) {
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
 * Attach event listeners to error UI elements
 * @private
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
 * @private
 */
function toggleErrorDetails() {
    const toggle = document.getElementById('errorDetailsToggle');
    const content = document.getElementById('errorDetailsContent');

    if (!toggle || !content) return;

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
 * @private
 */
function handleRetry() {
    hideError();
    // Dispatch retry event for last action
    const event = new CustomEvent('errorRetry');
    window.dispatchEvent(event);
}

/**
 * Handle reset action
 * @private
 */
function handleReset() {
    hideError();
    // Dispatch reset event
    const event = new CustomEvent('errorReset');
    window.dispatchEvent(event);
}

/**
 * Log error to console with full context
 * @private
 */
function logError(errorType, errorData, technicalError) {
    const timestamp = new Date().toISOString();
    console.error(`[${timestamp}] Recipe Error: ${errorType}`);
    console.error('User-Friendly Message:', errorData.message);
    console.error('Reason:', errorData.reason);
    console.error('Action:', errorData.action);

    if (technicalError) {
        console.error('Technical Error:', technicalError);
        if (technicalError.stack) {
            console.error('Stack Trace:', technicalError.stack);
        }
    }

    // Optional: Send to telemetry service (not implemented in MVP)
    // sendErrorTelemetry(errorType, errorData, technicalError);
}

/**
 * Escape HTML to prevent XSS attacks
 * @private
 * @param {string} text - Text to escape
 * @returns {string} Escaped HTML-safe text
 */
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
