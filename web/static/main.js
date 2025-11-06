// main.js - Recipe Web Interface Entry Point
// Epic 2, Stories 2-1, 2-2, 2-3, 2-4, 2-5: HTML Drag-Drop UI, File Upload, Format Detection, Parameter Display, Format Selection
// Initializes WASM module, drag-drop, and file loading

import { initializeWASM } from './wasm-loader.js';
import { initializeDropZone, getCurrentFileData, getCurrentFileName } from './file-handler.js';
import { detectFileFormat, getFormatDisplayName, getFormatBadgeClass, clearFormat } from './format-detector.js';
import { displayParameters, clearParameterPanel } from './parameter-display.js';
import { displayFormatSelector, clearFormatSelector } from './format-selector.js';

/**
 * Main initialization function
 * Called when DOM is fully loaded
 */
function init() {
    console.log('Recipe Web Interface - Stories 2-1 & 2-2');
    console.log('Initializing...');

    // Initialize WASM module (loads recipe.wasm and updates version)
    initializeWASM();

    // Initialize drag-drop zone (sets up event handlers)
    initializeDropZone();

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

    console.log('Initialization complete');
}

/**
 * Handle file loaded event
 * Stories 2-2 & 2-3: Called when file has been read and converted to Uint8Array
 * Triggers format detection
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

    // Story 2-3: Format detection
    await handleFormatDetection(fileData);
}

/**
 * Handle format detection for uploaded file
 * Story 2-3: Detect format using WASM and display badge
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
        // Detection failed
        console.error('Format detection error:', error);
        hideDetectionLoading();
        showError('Unable to detect format. Please upload a valid preset file (.np3, .xmp, or .lrtemplate)');

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
 * @param {string} format - Detected format ("np3" | "xmp" | "lrtemplate")
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
        // Parameter extraction failed
        console.error('Parameter extraction error:', error);
        hideParameterLoading();
        showParameterError('Unable to extract parameters. The file may be corrupted or in an unsupported format.');

        // Clear any partial display
        clearParameterPanel();
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
 * Story 2-6: Will implement actual conversion logic
 * @param {CustomEvent} event - Convert request event with fromFormat and toFormat
 */
function handleConvertRequest(event) {
    const { fromFormat, toFormat } = event.detail;
    console.log('Convert request received:', fromFormat, '→', toFormat);
    console.log('Story 2-6 will implement conversion logic here');
    // Story 2-6: Implement conversion using WASM here
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

// Wait for DOM to be ready before initializing
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    // DOM already loaded (script loaded after DOMContentLoaded)
    init();
}
