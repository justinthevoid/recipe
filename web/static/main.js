// main.js - Recipe Web Interface Entry Point
// Epic 2, Stories 2-1, 2-2, 2-3, 2-4, 2-5, 2-8, 2-9, 2-10: HTML Drag-Drop UI, File Upload, Format Detection, Parameter Display, Format Selection, Error Handling, Privacy Messaging, Responsive Design
// Initializes WASM module, drag-drop, and file loading

import { initializeWASM } from './wasm-loader.js';
import { initializeDropZone, getCurrentFileData, getCurrentFileName } from './file-handler.js';
import { detectFileFormat, getFormatDisplayName, getFormatBadgeClass, clearFormat } from './format-detector.js';
import { displayParameters, clearParameterPanel } from './parameter-display.js';
import { displayFormatSelector, clearFormatSelector } from './format-selector.js';
import { convertFile, getConvertedFileData, getConvertedFileName, clearConvertedData } from './converter.js';
import { enableDownload, clearDownloadState } from './downloader.js';
import { checkBrowserCompatibility, showError as showErrorPanel, hideError } from './error-handler.js';
import { showPrivacyReminder, showConversionPrivacyMessage, initializePrivacyFAQ } from './privacy-messaging.js';
import { initializeResponsive } from './responsive.js';

/**
 * Main initialization function
 * Called when DOM is fully loaded
 */
function init() {
    console.log('Recipe Web Interface - Stories 2-1, 2-2, 2-8, 2-9, & 2-10');
    console.log('Initializing...');

    // Story 2-8: Check browser compatibility before initialization
    if (!checkBrowserCompatibility()) {
        console.error('Browser compatibility check failed');
        return; // Stop initialization if browser is unsupported
    }

    // Story 2-10: Initialize responsive adaptations (touch detection, orientation handling)
    initializeResponsive();

    // Story 2-8: Global error handlers (AC-5: Error Boundaries)
    // Catch uncaught exceptions and unhandled promise rejections
    window.addEventListener('error', (event) => {
        console.error('Unhandled error:', event.error);
        showErrorPanel('unknown-error', event.error);
        event.preventDefault(); // Prevent default browser error handling
    });

    window.addEventListener('unhandledrejection', (event) => {
        console.error('Unhandled promise rejection:', event.reason);
        showErrorPanel('unknown-error', new Error(event.reason));
        event.preventDefault(); // Prevent default browser rejection handling
    });

    // Initialize WASM module (loads recipe.wasm and updates version)
    initializeWASM();

    // Initialize drag-drop zone (sets up event handlers)
    initializeDropZone();

    // Story 2-9: Initialize privacy FAQ (AC-4: FAQ toggle and privacy badge click)
    initializePrivacyFAQ();

    // Story 2-2: Listen for file loaded event
    // This allows other modules (Stories 2-3+) to react when a file is ready
    window.addEventListener('fileLoaded', handleFileLoaded);

    // Story 2-4: Listen for format detected event
    // Triggers parameter extraction after format detection
    window.addEventListener('formatDetected', handleFormatDetected);

    // Story 2-5: Listen for format selected event (optional - for analytics, logging)
    window.addEventListener('formatSelected', handleFormatSelected);

    // Story 2-6: Listen for convert request event (will implement conversion in Story 2-6)
    window.addEventListener('convertRequest', handleConvertRequest);

    // Story 2-7: Listen for conversion complete event (enable download)
    window.addEventListener('conversionComplete', handleConversionComplete);

    // Story 2-8: Listen for error recovery events
    window.addEventListener('errorRetry', handleErrorRetry);
    window.addEventListener('errorReset', handleErrorReset);

    console.log('Initialization complete');
}

/**
 * Handle file loaded event
 * Stories 2-2, 2-3, & 2-9: Called when file has been read and converted to Uint8Array
 * Triggers format detection and privacy reminder
 * @param {CustomEvent} event - File loaded event with file details
 */
async function handleFileLoaded(event) {
    const { fileName, fileSize, fileData } = event.detail;

    console.log('File loaded event received:', {
        fileName,
        fileSize,
        dataType: fileData instanceof Uint8Array ? 'Uint8Array' : typeof fileData,
        dataLength: fileData ? fileData.length : 0
    });

    // Story 2-7: Clear previous download state when new file uploaded
    clearDownloadState();

    // Story 2-9: Show privacy reminder (AC-2: Privacy message after file upload)
    showPrivacyReminder();

    // Story 2-3: Format detection
    await handleFormatDetection(fileData);
}

/**
 * Handle format detection for uploaded file
 * Story 2-3: Detect format using WASM and display badge
 * Story 2-8: Enhanced error handling with centralized error panel
 * @param {Uint8Array} fileData - Raw file bytes
 */
async function handleFormatDetection(fileData) {
    // Show loading state
    showDetectionLoading();

    try {
        // Detect format using WASM
        const format = await detectFileFormat(fileData);

        // Display format badge
        displayFormatBadge(format);

        // Hide loading state
        hideDetectionLoading();

        // Notify other components (Story 2-5 will listen)
        dispatchFormatDetectedEvent(format);

    } catch (error) {
        // Detection failed - Story 2-8: Use centralized error handler
        console.error('Format detection error:', error);
        hideDetectionLoading();

        // Show centralized error panel
        showErrorPanel('format-detection-failed', error);

        // Clear format state
        clearFormat();

        // Reset UI (allow retry)
        resetAfterError();
    }
}

/**
 * Show loading indicator during format detection
 */
function showDetectionLoading() {
    const statusEl = document.getElementById('status');
    statusEl.className = 'status loading';
    statusEl.textContent = 'Detecting format...';
    statusEl.style.display = 'block';
}

/**
 * Hide loading indicator
 */
function hideDetectionLoading() {
    const statusEl = document.getElementById('status');
    if (statusEl.classList.contains('loading')) {
        statusEl.style.display = 'none';
    }
}

/**
 * Display format badge in file info section
 * @param {string} format - Detected format ("np3" | "xmp" | "lrtemplate" | "costyle" | "costylepack")
 */
function displayFormatBadge(format) {
    const fileInfoEl = document.getElementById('fileInfo');
    const displayName = getFormatDisplayName(format);
    const badgeClass = getFormatBadgeClass(format);

    // Remove existing badge if any
    const existingBadge = fileInfoEl.querySelector('.format-badge');
    if (existingBadge) {
        existingBadge.remove();
    }

    // Add format badge to file info
    const badge = document.createElement('span');
    badge.className = `format-badge ${badgeClass}`;
    badge.textContent = displayName;
    fileInfoEl.appendChild(badge);
}

/**
 * Dispatch formatDetected event for other modules
 * @param {string} format - Detected format
 */
function dispatchFormatDetectedEvent(format) {
    const event = new CustomEvent('formatDetected', {
        detail: { format }
    });
    window.dispatchEvent(event);
    console.log(`formatDetected event dispatched: ${format}`);
}

/**
 * Handle format detected event
 * Story 2-4: Extract and display parameters after format detection
 * Story 2-5: Display format selector after parameter display
 * Story 2-8: Enhanced error handling with centralized error panel (non-blocking)
 * @param {CustomEvent} event - Format detected event with format string
 */
async function handleFormatDetected(event) {
    const { format } = event.detail;

    console.log('Format detected event received:', format);

    // Get current file data
    const fileData = getCurrentFileData();
    if (!fileData) {
        console.error('No file data available for parameter extraction');
        return;
    }

    // Show parameter extraction loading state
    showParameterLoading();

    try {
        // Story 2-4: Extract and display parameters
        const parameters = await displayParameters(fileData, format);

        // Hide loading state
        hideParameterLoading();

        console.log('Parameters displayed successfully:', Object.keys(parameters).length, 'parameters');

        // Story 2-5: Display format selector after parameters
        displayFormatSelector(format);

    } catch (error) {
        // Parameter extraction failed - Story 2-8: Non-blocking error (can continue)
        console.error('Parameter extraction error:', error);
        hideParameterLoading();

        // Show centralized error panel (non-blocking - includes "Continue Anyway" option)
        showErrorPanel('parameter-extraction-failed', error);

        // Clear any partial display
        clearParameterPanel();

        // Still display format selector (conversion may work even if parameter display failed)
        displayFormatSelector(format);
    }
}

/**
 * Show loading indicator during parameter extraction
 */
function showParameterLoading() {
    const statusEl = document.getElementById('parameterStatus');
    if (statusEl) {
        statusEl.className = 'status loading';
        statusEl.textContent = 'Extracting parameters...';
        statusEl.style.display = 'block';
    }

    // Hide any previous errors
    const errorEl = document.getElementById('parameterError');
    if (errorEl) {
        errorEl.style.display = 'none';
    }
}

/**
 * Hide parameter loading indicator
 */
function hideParameterLoading() {
    const statusEl = document.getElementById('parameterStatus');
    if (statusEl) {
        statusEl.style.display = 'none';
    }
}

/**
 * Show parameter extraction error
 * @param {string} message - Error message
 */
function showParameterError(message) {
    const errorEl = document.getElementById('parameterError');
    if (errorEl) {
        errorEl.textContent = message;
        errorEl.style.display = 'block';
    }
}

/**
 * Show error message to user
 * @param {string} message - Error message
 */
function showError(message) {
    const errorEl = document.getElementById('errorMessage');
    errorEl.textContent = message;
    errorEl.style.display = 'block';
}

/**
 * Handle format selected event
 * Story 2-5: Optional analytics/logging when user changes format selection
 * @param {CustomEvent} event - Format selected event with format string
 */
function handleFormatSelected(event) {
    const { format } = event.detail;
    console.log('User selected target format:', format);
    // Future: Add analytics tracking here if needed
}

/**
 * Handle convert request event
 * Story 2-6: Perform WASM conversion with UI state management
 * Story 2-8: Enhanced error handling with centralized error panel
 * @param {CustomEvent} event - Convert request event with fromFormat and toFormat
 */
async function handleConvertRequest(event) {
    const { fromFormat, toFormat } = event.detail;
    console.log('Convert request received:', fromFormat, '→', toFormat);

    // Show converting state
    showConvertingState();

    try {
        // Clear previous conversion data
        clearConvertedData();

        // Get source file data and name
        const fileData = getCurrentFileData();
        const fileName = getCurrentFileName();

        if (!fileData || !fileName) {
            throw new Error('No file data available for conversion');
        }

        // Perform conversion
        const convertedData = await convertFile(fileData, fromFormat, toFormat, fileName);

        // Show success state
        showConversionSuccess(toFormat);

        // Enable download button (Story 2-7 will implement download)
        enableDownloadButton();

        // Dispatch conversion complete event
        dispatchConversionCompleteEvent(convertedData, toFormat);

    } catch (error) {
        // Conversion failed - Story 2-8: Use centralized error handler
        console.error('Conversion error:', error);
        showConversionError(error);

        // Show centralized error panel
        showErrorPanel('conversion-failed', error);
    }
}

/**
 * Show converting state (Story 2-6)
 */
function showConvertingState() {
    const convertButton = document.getElementById('convertButton');
    if (convertButton) {
        convertButton.disabled = true;
        convertButton.innerHTML = '⟳ Converting...';
        convertButton.classList.add('converting');
        convertButton.classList.remove('success', 'error');
    }

    // Show status message
    const statusEl = document.getElementById('conversionStatus');
    if (statusEl) {
        statusEl.className = 'status loading';
        statusEl.textContent = 'Converting preset...';
        statusEl.style.display = 'block';
    }

    // Hide any previous errors
    const errorEl = document.getElementById('conversionError');
    if (errorEl) {
        errorEl.style.display = 'none';
    }
}

/**
 * Show conversion success state (Stories 2-6 & 2-9)
 * @param {string} targetFormat - Target format name
 */
function showConversionSuccess(targetFormat) {
    const convertButton = document.getElementById('convertButton');
    if (convertButton) {
        convertButton.disabled = true;
        convertButton.innerHTML = '✓ Converted!';
        convertButton.classList.remove('converting', 'error');
        convertButton.classList.add('success');
    }

    // Show success message
    const statusEl = document.getElementById('conversionStatus');
    if (statusEl) {
        statusEl.className = 'status success';
        statusEl.textContent = `✓ Conversion complete! Your ${targetFormat.toUpperCase()} preset is ready.`;
    }

    // Story 2-9: Add privacy reminder to success message (AC-3)
    showConversionPrivacyMessage();
}

/**
 * Show conversion error state (Story 2-6)
 * @param {Error} error - Conversion error
 */
function showConversionError(error) {
    const convertButton = document.getElementById('convertButton');
    if (convertButton) {
        convertButton.disabled = false; // Re-enable for retry
        convertButton.innerHTML = '✗ Conversion Failed';
        convertButton.classList.remove('converting', 'success');
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

/**
 * Handle conversion complete event (Story 2-7)
 * Enable download button with converted file data
 * @param {CustomEvent} event - Conversion complete event
 */
function handleConversionComplete(event) {
    const { convertedData, format } = event.detail;

    // Get converted file metadata
    const fileName = getConvertedFileName();

    if (!fileName) {
        console.error('Cannot enable download: filename not available');
        return;
    }

    // Enable download button
    enableDownload(convertedData, fileName, format);

    console.log('Download ready:', fileName);
}

/**
 * Dispatch conversion complete event (Story 2-7 will listen)
 * @param {Uint8Array} convertedData - Converted file data
 * @param {string} format - Target format
 */
function dispatchConversionCompleteEvent(convertedData, format) {
    const event = new CustomEvent('conversionComplete', {
        detail: { convertedData, format }
    });
    window.dispatchEvent(event);
    console.log(`conversionComplete event dispatched: ${format}`);
}

/**
 * Reset UI after detection error (allow retry)
 * Clears file data but keeps drop zone ready
 */
function resetAfterError() {
    // Hide file info
    const fileInfoEl = document.getElementById('fileInfo');
    fileInfoEl.style.display = 'none';
    fileInfoEl.innerHTML = '';

    // Clear format selector
    clearFormatSelector();

    // Note: File data will be cleared when user uploads next file
    // (file-handler.js clears on new upload)
}

/**
 * Handle error retry event (Story 2-8)
 * Re-attempt the last operation based on current state
 */
function handleErrorRetry() {
    console.log('Error retry requested');
    hideError();

    // Determine what to retry based on current UI state
    const fileData = getCurrentFileData();

    if (!fileData) {
        // No file loaded - prompt user to upload file
        console.log('No file to retry with');
        return;
    }

    // Check if we have a convert button (conversion phase)
    const convertButton = document.getElementById('convertButton');
    if (convertButton && !convertButton.disabled) {
        // Retry conversion
        convertButton.click();
        return;
    }

    // Otherwise, retry format detection
    handleFormatDetection(fileData);
}

/**
 * Handle error reset event (Story 2-8)
 * Clear all state and return to initial UI
 */
function handleErrorReset() {
    console.log('Error reset requested');
    hideError();

    // Simple approach: reload the page to reset all state
    location.reload();
}

/**
 * Enable download button (Story 2-7)
 * Shows download button after successful conversion
 */
function enableDownloadButton() {
    const downloadButton = document.getElementById('downloadButton');
    if (downloadButton) {
        downloadButton.disabled = false;
        downloadButton.style.display = 'block';
    }
}

// Wait for DOM to be ready before initializing
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    // DOM already loaded (script loaded after DOMContentLoaded)
    init();
}
