// file-handler.js - File Upload Handling
// Epic 2, Stories 2-1, 2-2, 2-3: HTML Drag-Drop UI, File Upload, Format Detection
// Handles drag-drop, file picker, file reading, and data conversion

import { clearFormat } from './format-detector.js';

// Module-level state for current file (Story 2-2)
let currentFileData = null; // Uint8Array
let currentFileName = null;
let currentFileSize = 0;

/**
 * Initialize the drop zone with all event handlers
 * Sets up drag-drop, click-to-browse, keyboard accessibility, and file validation
 */
export function initializeDropZone() {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');

    if (!dropZone || !fileInput) {
        console.error('Drop zone or file input element not found');
        return;
    }

    console.log('Initializing drop zone...');

    // Click to open file picker
    dropZone.addEventListener('click', () => {
        fileInput.click();
    });

    // Keyboard accessibility - Enter or Space opens file picker
    dropZone.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            fileInput.click();
        }
    });

    // Prevent default drag behavior on drop zone and body
    // This prevents browser from opening file in new tab
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, preventDefaults, false);
        document.body.addEventListener(eventName, preventDefaults, false);
    });

    // Highlight drop zone when dragging over
    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => {
            dropZone.classList.add('drag-over');
        });
    });

    // Remove highlight when drag leaves or file is dropped
    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => {
            dropZone.classList.remove('drag-over');
        });
    });

    // Handle file drop
    dropZone.addEventListener('drop', handleDrop);

    // Handle file picker selection
    fileInput.addEventListener('change', handleFileSelect);

    console.log('Drop zone initialized successfully');
}

/**
 * Prevent default browser behavior for drag events
 * @param {Event} e - Drag event
 */
function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
}

/**
 * Handle file drop from drag-drop
 * @param {DragEvent} e - Drop event
 */
function handleDrop(e) {
    const dt = e.dataTransfer;
    const files = dt.files;

    if (files.length > 0) {
        handleFile(files[0]);
    } else {
        console.warn('No files in drop event');
    }
}

/**
 * Handle file selection from file picker
 * @param {Event} e - Change event from file input
 */
function handleFileSelect(e) {
    const files = e.target.files;
    if (files.length > 0) {
        handleFile(files[0]);
    }
}

/**
 * Process uploaded file
 * Story 2-2: Validates, reads file, converts to Uint8Array, displays metadata
 * @param {File} file - The file object to process
 */
async function handleFile(file) {
    console.log('File received:', file.name, file.size, 'bytes', file.type);

    // Clear previous file data (AC-6: Memory Management)
    clearFileData();

    // Validate file extension (Story 2-1)
    if (!isValidPresetFile(file.name)) {
        showError('Please upload a preset file (.np3, .xmp, .lrtemplate, .costyle, or .costylepack)');
        updateDropZoneState('error');
        return;
    }

    // AC-4: File Size Validation (10MB limit)
    const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB in bytes
    if (file.size > MAX_FILE_SIZE) {
        showError('File too large. Maximum size: 10MB');
        updateDropZoneState('error');
        return;
    }

    // AC-5: Check for empty file
    if (file.size === 0) {
        showError('File appears to be empty or corrupted');
        updateDropZoneState('error');
        return;
    }

    // AC-1: Show loading state for file reading
    showLoadingState('Reading file...');

    try {
        // AC-1: Read file as ArrayBuffer using FileReader API
        const arrayBuffer = await readFileAsArrayBuffer(file);

        // AC-2: Convert to Uint8Array (WASM-compatible format)
        currentFileData = new Uint8Array(arrayBuffer);
        currentFileName = file.name;
        currentFileSize = file.size;

        // AC-3: Display file metadata
        displayFileInfo(file.name, file.size);

        // Hide loading state
        hideLoadingState();

        // Update UI to success state
        hideError();
        updateDropZoneState('success');

        // Success feedback
        updateStatus('success', `File loaded: ${file.name}`);

        // Dispatch custom event for other modules (Story 2-3+)
        dispatchFileLoadedEvent();

        console.log('✅ File loaded successfully:', {
            name: currentFileName,
            size: currentFileSize,
            dataLength: currentFileData.length
        });

    } catch (error) {
        // AC-5: File read error handling
        console.error('File read error:', error);
        showError('Failed to read file. Please try again.');
        hideLoadingState();
        updateDropZoneState('error');
    }
}

/**
 * Read file as ArrayBuffer using FileReader API
 * AC-1: Promise wrapper around FileReader for async/await
 * @param {File} file - The file to read
 * @returns {Promise<ArrayBuffer>}
 */
function readFileAsArrayBuffer(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = (event) => {
            resolve(event.target.result);
        };

        reader.onerror = (event) => {
            reject(new Error('FileReader error: ' + event.target.error));
        };

        // Start reading
        reader.readAsArrayBuffer(file);
    });
}

/**
 * Display file metadata in UI
 * AC-3: Show file name, size, and extension
 * @param {string} fileName - Name of the uploaded file
 * @param {number} fileSize - Size in bytes
 */
function displayFileInfo(fileName, fileSize) {
    const fileInfoEl = document.getElementById('fileInfo');
    if (!fileInfoEl) return;

    const formattedSize = formatFileSize(fileSize);
    const extension = fileName.substring(fileName.lastIndexOf('.')).toLowerCase();

    // Use escapeHtml to prevent XSS attacks
    fileInfoEl.innerHTML = `
        <strong>File:</strong> ${escapeHtml(fileName)} (${formattedSize})
    `;
    fileInfoEl.style.display = 'block';
}

/**
 * Format bytes to human-readable size
 * AC-3: Format file size display
 * @param {number} bytes - File size in bytes
 * @returns {string} Formatted size string
 */
function formatFileSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
}

/**
 * Escape HTML to prevent XSS
 * AC-3: Security measure for displaying user-controlled text
 * @param {string} text - Text to escape
 * @returns {string} HTML-safe text
 */
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

/**
 * Clear file data from memory
 * AC-6: Prevent memory leaks when uploading new files
 * Story 2-3: Also clear detected format
 */
function clearFileData() {
    currentFileData = null;
    currentFileName = null;
    currentFileSize = 0;

    // Story 2-3: Clear detected format when new file uploaded
    clearFormat();

    // Also clear the UI
    const fileInfoEl = document.getElementById('fileInfo');
    if (fileInfoEl) {
        fileInfoEl.style.display = 'none';
        fileInfoEl.innerHTML = '';
    }
}

/**
 * Dispatch custom event when file is loaded
 * Allows other modules (Stories 2-3+) to react to file loading
 */
function dispatchFileLoadedEvent() {
    const event = new CustomEvent('fileLoaded', {
        detail: {
            fileName: currentFileName,
            fileSize: currentFileSize,
            fileData: currentFileData
        }
    });
    window.dispatchEvent(event);
}

/**
 * Validate file extension
 * @param {string} fileName - Name of the file
 * @returns {boolean} True if valid preset file
 */
function isValidPresetFile(fileName) {
    const validExtensions = ['.np3', '.xmp', '.lrtemplate', '.costyle', '.costylepack'];
    const lowerName = fileName.toLowerCase();
    return validExtensions.some(ext => lowerName.endsWith(ext));
}

/**
 * Update drop zone visual state
 * @param {string} state - 'error', 'success', or 'default'
 */
function updateDropZoneState(state) {
    const dropZone = document.getElementById('dropZone');
    if (!dropZone) return;

    dropZone.classList.remove('error', 'success');

    if (state === 'error') {
        dropZone.classList.add('error');
        // Remove error state after 3 seconds
        setTimeout(() => {
            dropZone.classList.remove('error');
            hideError();
        }, 3000);
    } else if (state === 'success') {
        dropZone.classList.add('success');
    }
}

/**
 * Display error message to user
 * @param {string} message - Error message to display
 */
function showError(message) {
    const errorEl = document.getElementById('errorMessage');
    if (errorEl) {
        errorEl.textContent = message;
        errorEl.style.display = 'block';
    }
}

/**
 * Hide error message
 */
function hideError() {
    const errorEl = document.getElementById('errorMessage');
    if (errorEl) {
        errorEl.style.display = 'none';
    }
}

/**
 * Show loading state
 * AC-1: Display loading indicator during file read
 * @param {string} message - Loading message to display
 */
function showLoadingState(message) {
    const statusEl = document.getElementById('status');
    if (statusEl) {
        statusEl.className = 'status loading';
        statusEl.textContent = message;
        statusEl.style.display = 'block';
    }
}

/**
 * Hide loading state
 */
function hideLoadingState() {
    const statusEl = document.getElementById('status');
    if (statusEl && statusEl.classList.contains('loading')) {
        statusEl.style.display = 'none';
    }
}

/**
 * Update status message
 * @param {string} type - Status type ('success', 'error', 'loading')
 * @param {string} message - Status message
 */
function updateStatus(type, message) {
    const statusEl = document.getElementById('status');
    if (statusEl) {
        statusEl.className = `status ${type}`;
        statusEl.textContent = message;
        statusEl.style.display = 'block';
    }
}

/**
 * Get current file data (exported for Story 2-3+)
 * AC-2: Provide access to file data for other modules
 * @returns {Uint8Array|null} Current file data
 */
export function getCurrentFileData() {
    return currentFileData;
}

/**
 * Get current file name (exported for Story 2-7)
 * @returns {string|null} Current file name
 */
export function getCurrentFileName() {
    return currentFileName;
}
