// main.js - Recipe Web Interface Entry Point
// Epic 2, Stories 2-1, 2-2, 2-3, 2-4, 2-5, 2-8, 2-9, 2-10: HTML Drag-Drop UI, File Upload, Format Detection, Parameter Display, Format Selection, Error Handling, Privacy Messaging, Responsive Design
// Epic 10, Story 10-2: Batch File Upload with Drag-and-Drop
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
import { initializeUploadManager, getUploadManager } from './upload.js';
import { startBatchConversion, cancelBatchConversion } from './batch-converter.js';

/**
 * Story 10-7: ARIA Live Region Announcement Functions (AC-2)
 * Provides screen reader announcements for file uploads, conversions, and errors
 */

/**
 * Announce status updates to screen readers (polite - non-urgent)
 * @param {string} message - The message to announce
 */
export function announceStatus(message) {
    const statusEl = document.getElementById('status-announcements');
    if (statusEl) {
        statusEl.textContent = message;
        console.log('[ARIA Status]', message);
    }
}

/**
 * Announce errors to screen readers (assertive - urgent)
 * @param {string} message - The error message to announce
 */
export function announceError(message) {
    const errorEl = document.getElementById('error-announcements');
    if (errorEl) {
        errorEl.textContent = message;
        console.log('[ARIA Error]', message);
    }
}

/**
 * Clear status announcements
 */
export function clearStatusAnnouncement() {
    const statusEl = document.getElementById('status-announcements');
    if (statusEl) {
        statusEl.textContent = '';
    }
}

/**
 * Clear error announcements
 */
export function clearErrorAnnouncement() {
    const errorEl = document.getElementById('error-announcements');
    if (errorEl) {
        errorEl.textContent = '';
    }
}

/**
 * Story 10-7: Keyboard Event Handlers (AC-1)
 * Adds keyboard accessibility to drag-drop zone
 */
function initializeKeyboardAccessibility() {
    const dropzone = document.getElementById('dropzone');
    if (dropzone) {
        // Handle Enter/Space on dropzone to trigger file picker (AC-1)
        dropzone.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                document.getElementById('file-input').click();
                announceStatus('File picker opened. Select files to upload.');
            }
        });
    }

    // Handle Escape key to cancel operations (AC-1)
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            // Cancel batch conversion if in progress
            const batchProgress = document.getElementById('batch-progress');
            if (batchProgress && !batchProgress.hidden) {
                cancelBatchConversion();
                announceStatus('Batch conversion cancelled.');
            }
        }
    });
}

/**
 * Story 10-7: Focus Management (AC-5)
 * Moves focus to relevant elements after user actions
 */

/**
 * Move focus to first file card after upload
 * @param {string} fileId - The ID of the first uploaded file
 */
export function moveFocusToFirstFileCard(fileId) {
    setTimeout(() => {
        const firstCard = document.querySelector(`[data-file-id="${fileId}"]`);
        if (firstCard) {
            const formatSelect = firstCard.querySelector('select');
            if (formatSelect) {
                formatSelect.focus();
                console.log('[Focus] Moved to first file card format selector');
            }
        }
    }, 100); // Small delay to ensure DOM is updated
}

/**
 * Move focus to download button after successful conversion
 * @param {string} fileId - The ID of the converted file
 */
export function moveFocusToDownloadButton(fileId) {
    setTimeout(() => {
        const card = document.querySelector(`[data-file-id="${fileId}"]`);
        if (card) {
            const downloadBtn = card.querySelector('.file-card__download');
            if (downloadBtn) {
                downloadBtn.focus();
                console.log('[Focus] Moved to download button for file', fileId);
            }
        }
    }, 100);
}

/**
 * Move focus to error message after conversion failure
 * @param {string} fileId - The ID of the file that failed
 */
export function moveFocusToErrorMessage(fileId) {
    setTimeout(() => {
        const card = document.querySelector(`[data-file-id="${fileId}"]`);
        if (card) {
            const errorMsg = card.querySelector('.file-card__error');
            if (errorMsg) {
                errorMsg.setAttribute('tabindex', '-1'); // Make focusable
                errorMsg.focus();
                console.log('[Focus] Moved to error message for file', fileId);
            }
        }
    }, 100);
}

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

    // Story 10-7: Initialize keyboard accessibility (Enter/Space on dropzone, Escape to cancel)
    initializeKeyboardAccessibility();

    // Story 10-2: Initialize batch upload manager (drag-drop, file validation, file cards)
    initializeUploadManager();

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

    // Story 10-6: REMOVED eager WASM loading for performance
    // WASM now lazy-loads on first file upload (see handleFileLoaded function)
    // This saves ~4MB (~500KB gzipped) on initial page load
    // Performance improvement: TTI reduced by ~1-1.5 seconds on 3G
    console.log('WASM will load on first file upload (lazy-loading enabled)');

    // Story 10-2: Drag-drop is now handled by UploadManager (batch upload)
    // Old initializeDropZone() removed to prevent duplicate file processing

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

    // Story 10-3: Initialize batch conversion controls
    initializeBatchControls();

    // Update status banner to show ready state (WASM loads on first upload)
    const statusEl = document.getElementById('status');
    if (statusEl) {
        statusEl.className = 'status ready';
        statusEl.textContent = 'Ready to convert';
    }

    console.log('Initialization complete');
}

/**
 * Initialize batch conversion controls (Story 10-3)
 * Sets up format selector and convert all button
 */
function initializeBatchControls() {
    const targetFormatSelect = document.getElementById('target-format');
    const convertAllButton = document.getElementById('convert-all-button');
    const cancelBatchButton = document.getElementById('cancel-batch');

    if (!targetFormatSelect || !convertAllButton || !cancelBatchButton) {
        console.error('Batch control elements not found');
        return;
    }

    // Enable/disable convert button based on format selection
    targetFormatSelect.addEventListener('change', () => {
        const selectedFormat = targetFormatSelect.value;
        convertAllButton.disabled = !selectedFormat;
    });

    // Convert all button click handler
    convertAllButton.addEventListener('click', async () => {
        const selectedFormat = targetFormatSelect.value;
        if (!selectedFormat) {
            console.warn('No target format selected');
            return;
        }

        // Disable controls during conversion
        convertAllButton.disabled = true;
        targetFormatSelect.disabled = true;

        try {
            await startBatchConversion(selectedFormat);
        } catch (error) {
            console.error('Batch conversion failed:', error);
            showErrorPanel('batch-conversion-failed', error);
        } finally {
            // Re-enable controls after conversion
            convertAllButton.disabled = false;
            targetFormatSelect.disabled = false;
        }
    });

    // Cancel batch button click handler
    cancelBatchButton.addEventListener('click', () => {
        cancelBatchConversion();
    });
}

/**
 * Handle file loaded event
 * Stories 2-2, 2-3, & 2-9: Called when file has been read and converted to Uint8Array
 * Story 10-6: Lazy-load WASM on first file upload
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

    // Story 10-6: Lazy-load WASM on first file upload (performance optimization)
    // Only load WASM when user actually needs it (first file uploaded)
    // This defers ~4MB (~500KB gzipped) from initial page load
    try {
        await initializeWASM();
        console.log('WASM initialized on first file upload');
    } catch (error) {
        console.error('Failed to initialize WASM on file upload:', error);
        // Error already shown by wasm-loader.js via showErrorPanel
        return; // Cannot proceed without WASM
    }

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
