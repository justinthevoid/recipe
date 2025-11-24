// upload.js - Batch File Upload Manager
// Story 10-2: Batch File Upload with Drag-and-Drop
// Story 10-4: Individual File Conversion Controls
// Story 10-7: Accessibility Enhancements (ARIA announcements)
// Handles drag-drop, file validation, file cards, and empty state

import { convertFile } from './converter.js';
import { announceStatus, announceError, moveFocusToFirstFileCard } from './main.js';
import { applyPreviewFilter } from './preview.js';

/**
 * UploadManager Class
 * Manages batch file uploads with drag-and-drop support
 * Story 10-4: Added individual file conversion support
 */
export class UploadManager {
    constructor() {
        // DOM references
        this.dropzone = document.getElementById('dropzone');
        this.fileInput = document.getElementById('file-input');
        this.browseButton = document.getElementById('browse-button');
        this.errorDisplay = document.getElementById('upload-error');
        this.fileGrid = document.getElementById('file-grid');

        // State
        this.uploadedFiles = new Map(); // fileId -> file object
        this.fileIdCounter = 0;

        // Supported formats (AC-2, AC-3)
        this.supportedExtensions = ['.np3', '.xmp', '.lrtemplate', '.costyle', '.dcp'];

        // Initialize
        this.initEventListeners();

        // Ensure error display is hidden on init
        this.hideError();
    }

    /**
     * Initialize all event listeners (AC-1, AC-2)
     */
    initEventListeners() {
        // Browse button click -> trigger file input
        this.browseButton.addEventListener('click', (e) => {
            e.stopPropagation();
            this.fileInput.click();
        });

        // Drop zone click -> trigger file input
        this.dropzone.addEventListener('click', () => {
            this.fileInput.click();
        });

        // File input change -> handle selected files
        this.fileInput.addEventListener('change', (e) => {
            this.handleFiles(e.target.files);
            // Reset input so same file can be selected again
            e.target.value = '';
        });

        // Drag-and-drop events (AC-1)
        this.dropzone.addEventListener('dragover', (e) => {
            e.preventDefault();
            e.stopPropagation();
            this.dropzone.classList.add('upload__dropzone--drag-over');
        });

        this.dropzone.addEventListener('dragleave', (e) => {
            e.preventDefault();
            e.stopPropagation();
            // Only remove highlight if leaving the dropzone itself (not child elements)
            if (e.target === this.dropzone) {
                this.dropzone.classList.remove('upload__dropzone--drag-over');
            }
        });

        this.dropzone.addEventListener('drop', (e) => {
            e.preventDefault();
            e.stopPropagation();
            this.dropzone.classList.remove('upload__dropzone--drag-over');

            // Handle dropped files
            const files = e.dataTransfer.files;
            this.handleFiles(files);

            // Show success animation (AC-1)
            this.showSuccessAnimation();
        });
    }

    /**
     * Handle uploaded files (AC-2, AC-3)
     * Validates files and creates file cards
     * @param {FileList} fileList - Files from input or drag-drop
     */
    handleFiles(fileList) {
        if (!fileList || fileList.length === 0) {
            return;
        }

        const files = Array.from(fileList);
        const validFiles = [];
        const rejectedFiles = [];

        // Validate each file (AC-3)
        files.forEach(file => {
            const validation = this.validateFile(file);
            if (validation.valid) {
                validFiles.push(file);
            } else {
                rejectedFiles.push({
                    name: file.name,
                    reason: validation.message
                });
            }
        });

        // Add valid files to grid
        let firstFileId = null;
        validFiles.forEach((file, index) => {
            const fileId = this.addFileCard(file);
            if (index === 0) {
                firstFileId = fileId;
            }
        });

        // Story 10-7: Announce file upload to screen readers (AC-2)
        if (validFiles.length > 0) {
            const fileNames = validFiles.map(f => f.name).join(', ');
            if (validFiles.length === 1) {
                announceStatus(`File uploaded: ${fileNames}`);
            } else {
                announceStatus(`${validFiles.length} files uploaded: ${fileNames}`);
            }

            // Story 10-7: Move focus to first file card (AC-5)
            if (firstFileId) {
                moveFocusToFirstFileCard(firstFileId);
            }
        }

        // Show error for rejected files (AC-3)
        if (rejectedFiles.length > 0) {
            this.showError(validFiles.length, rejectedFiles);
            // Story 10-7: Announce errors to screen readers (AC-2)
            const errorSummary = rejectedFiles.map(f => `${f.name}: ${f.reason}`).join('. ');
            announceError(`Error: ${rejectedFiles.length} file(s) rejected. ${errorSummary}`);
        } else if (validFiles.length > 0) {
            // Hide error if all files valid
            this.hideError();
        }

        // Show batch controls and hide empty state if files uploaded (AC-8, Story 10-3)
        if (this.uploadedFiles.size > 0) {
            this.hideEmptyState();
            this.showBatchControls();
        }
    }

    /**
     * Validate file extension (AC-3)
     * @param {File} file - File to validate
     * @returns {Object} Validation result { valid: boolean, extension: string, message: string }
     */
    validateFile(file) {
        const fileName = file.name.toLowerCase();
        const extension = this.supportedExtensions.find(ext => fileName.endsWith(ext));

        if (extension) {
            return {
                valid: true,
                extension: extension.replace('.', ''),
                message: ''
            };
        } else {
            return {
                valid: false,
                extension: '',
                message: `Unsupported file type: ${file.name}. Supported formats: NP3, XMP, lrtemplate, .costyle, DCP`
            };
        }
    }

    /**
     * Add file card to grid (AC-4, AC-5, AC-6, AC-7)
     * @param {File} file - File object
     */
    addFileCard(file) {
        const fileId = this.fileIdCounter++;
        const format = this.detectFormat(file.name);
        const fileSize = this.formatFileSize(file.size);
        const truncatedName = this.truncateFilename(file.name, 30);

        // Get default target format (Story 10-4, AC-1)
        const validTargets = this.getValidTargetFormats(format);
        const defaultTarget = validTargets[0];

        // Store file data (Story 10-3: Extended with output data fields)
        // Story 10-4: Added targetFormat and abortController fields
        this.uploadedFiles.set(fileId, {
            id: fileId,
            file: file,
            format: format,
            targetFormat: defaultTarget, // Story 10-4: Selected target format
            status: 'queued',
            outputData: null,
            outputFormat: null,
            error: null,
            abortController: null // Story 10-4, AC-7: For cancellation support
        });

        // Create file card HTML
        const cardHTML = this.createFileCard(fileId, truncatedName, file.name, format, fileSize);

        // Append to grid
        this.fileGrid.insertAdjacentHTML('beforeend', cardHTML);

        // Add event listeners
        const card = document.getElementById(`file-${fileId}`);

        // Remove button listener
        const removeButton = card.querySelector('.file-card__remove');
        removeButton.addEventListener('click', () => {
            this.removeFile(fileId);
        });

        // Download button listener (Story 10-3, AC-5)
        const downloadButton = card.querySelector('.file-card__download');
        downloadButton.addEventListener('click', () => {
            this.downloadFile(fileId);
        });

        // Format dropdown listener (Story 10-4, AC-1)
        const formatSelect = card.querySelector('.file-card__format-select');
        formatSelect.addEventListener('change', (e) => {
            const fileData = this.uploadedFiles.get(fileId);
            if (fileData) {
                fileData.targetFormat = e.target.value;
            }
        });

        // Preview button listener (Story 11-1, AC-3)
        const previewButton = card.querySelector('.file-card__preview');
        previewButton.addEventListener('click', () => {
            this.showPreviewForFile(fileId);
        });

        // Convert button listener (Story 10-4, AC-2)
        const convertButton = card.querySelector('.file-card__convert');
        convertButton.addEventListener('click', () => {
            this.convertIndividualFile(fileId);
        });

        // Retry button listener (Story 10-4, AC-5)
        const retryButton = card.querySelector('.file-card__retry');
        retryButton.addEventListener('click', () => {
            this.convertIndividualFile(fileId);
        });

        // Cancel button listener (Story 10-4, AC-7)
        const cancelButton = card.querySelector('.file-card__cancel');
        cancelButton.addEventListener('click', () => {
            this.cancelIndividualConversion(fileId);
        });

        // Story 10-7: Return fileId for focus management (AC-5)
        return fileId;
    }

    /**
     * Create file card HTML (AC-4, AC-5, AC-6, AC-7)
     * Story 10-3: Added status indicators and download button
     * Story 10-4: Added format dropdown and individual convert button (AC-1, AC-2)
     * @param {number} fileId - Unique file ID
     * @param {string} displayName - Truncated filename for display
     * @param {string} fullName - Full filename for title attribute
     * @param {string} format - Detected format
     * @param {string} fileSize - Formatted file size
     * @returns {string} HTML string for file card
     */
    createFileCard(fileId, displayName, fullName, format, fileSize) {
        // Get valid target formats for this source format (Story 10-4, AC-1)
        const validTargets = this.getValidTargetFormats(format);
        const defaultTarget = validTargets[0]; // First format alphabetically

        // Build format dropdown options
        const formatOptions = validTargets.map(targetFormat =>
            `<option value="${targetFormat}" ${targetFormat === defaultTarget ? 'selected' : ''}>
                ${targetFormat.toUpperCase()}
            </option>`
        ).join('');

        return `
            <li class="file-card" id="file-${fileId}" data-status="queued" data-file-id="${fileId}" role="listitem">
                <div class="file-card__header">
                    <span class="file-card__filename" title="${fullName}">${displayName}</span>
                    <span class="badge badge--${format}" role="img" aria-label="${format.toUpperCase()} format">${format.toUpperCase()}</span>
                </div>
                <div class="file-card__body">
                    <div class="file-card__size">${fileSize}</div>
                    <div class="file-card__status" role="status" aria-live="polite">
                        <span class="file-card__status-icon" aria-hidden="true">⏱️</span>
                        <span class="file-card__status-text">Queued</span>
                    </div>
                </div>
                <div class="file-card__conversion">
                    <label class="file-card__format-label" for="format-${fileId}">Convert to:</label>
                    <select class="file-card__format-select" id="format-${fileId}" data-file-id="${fileId}" aria-label="Select target format for ${fullName}">
                        ${formatOptions}
                    </select>
                    <button class="file-card__preview" data-file-id="${fileId}" aria-label="Preview ${fullName} adjustments">Preview</button>
                    <button class="file-card__convert" data-file-id="${fileId}" aria-label="Convert ${fullName} to ${defaultTarget.toUpperCase()}">Convert</button>
                    <button class="file-card__retry" data-file-id="${fileId}" hidden aria-label="Retry converting ${fullName}">Retry</button>
                    <button class="file-card__cancel" data-file-id="${fileId}" hidden aria-label="Cancel converting ${fullName}">Cancel</button>
                </div>
                <div class="file-card__footer">
                    <button class="file-card__download" data-file-id="${fileId}" hidden aria-label="Download converted ${fullName}">Download ${format.toUpperCase()}</button>
                    <button class="file-card__remove" data-file-id="${fileId}" aria-label="Remove ${fullName}">Remove</button>
                </div>
            </li>
        `;
    }

    /**
     * Remove file from grid and memory (AC-7)
     * @param {number} fileId - File ID to remove
     */
    removeFile(fileId) {
        const card = document.getElementById(`file-${fileId}`);
        if (!card) return;

        // Add removing animation
        card.classList.add('file-card--removing');

        // Remove from DOM after animation
        setTimeout(() => {
            card.remove();

            // Remove from state
            this.uploadedFiles.delete(fileId);

            // Show empty state and hide batch controls if no files left (AC-8, Story 10-3)
            if (this.uploadedFiles.size === 0) {
                this.showEmptyState();
                this.hideBatchControls();
            }
        }, 300); // Match animation duration
    }

    /**
     * Detect format from filename (AC-5)
     * @param {string} filename - File name
     * @returns {string} Format identifier (np3, xmp, lrtemplate, costyle, dcp)
     */
    detectFormat(filename) {
        const lowerName = filename.toLowerCase();
        if (lowerName.endsWith('.np3')) return 'np3';
        if (lowerName.endsWith('.xmp')) return 'xmp';
        if (lowerName.endsWith('.lrtemplate')) return 'lrtemplate';
        if (lowerName.endsWith('.costyle')) return 'costyle';
        if (lowerName.endsWith('.dcp')) return 'dcp';
        return 'unknown';
    }

    /**
     * Get valid target formats for a source format (Story 10-4, AC-1)
     * Returns all supported formats except the source format
     * @param {string} sourceFormat - Source format identifier
     * @returns {string[]} Array of valid target formats (sorted alphabetically)
     */
    getValidTargetFormats(sourceFormat) {
        const allFormats = ['costyle', 'dcp', 'lrtemplate', 'np3', 'xmp'];
        return allFormats.filter(format => format !== sourceFormat).sort();
    }

    /**
     * Format file size for display (AC-6)
     * @param {number} bytes - File size in bytes
     * @returns {string} Human-readable file size
     */
    formatFileSize(bytes) {
        if (bytes < 1024) {
            return `${bytes} bytes`;
        } else if (bytes < 1024 * 1024) {
            return `${(bytes / 1024).toFixed(1)} KB`;
        } else {
            return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
        }
    }

    /**
     * Truncate long filenames (AC-4)
     * @param {string} name - File name
     * @param {number} max - Maximum length
     * @returns {string} Truncated filename
     */
    truncateFilename(name, max) {
        if (name.length <= max) {
            return name;
        }
        const extension = name.substring(name.lastIndexOf('.'));
        const baseName = name.substring(0, name.lastIndexOf('.'));
        const truncatedBase = baseName.substring(0, max - extension.length - 3);
        return `${truncatedBase}...${extension}`;
    }

    /**
     * Show error message for rejected files (AC-3)
     * @param {number} accepted - Number of accepted files
     * @param {Array} rejected - Array of rejected file objects
     */
    showError(accepted, rejected) {
        if (!this.errorDisplay) {
            console.error('Error display element not found');
            return;
        }

        const rejectedCount = rejected.length;
        let message = '';

        if (accepted > 0) {
            // Mixed batch: some accepted, some rejected
            message = `${accepted} file(s) uploaded, ${rejectedCount} file(s) rejected`;
        } else {
            // All rejected
            message = `${rejectedCount} file(s) rejected`;
        }

        // Show first rejection reason as example
        if (rejected.length > 0) {
            message += `. ${rejected[0].reason}`;
        }

        this.errorDisplay.textContent = message;
        this.errorDisplay.hidden = false;
    }

    /**
     * Hide error message
     */
    hideError() {
        if (!this.errorDisplay) {
            return;
        }
        this.errorDisplay.hidden = true;
        this.errorDisplay.textContent = '';
    }

    /**
     * Show success animation on drop (AC-1)
     */
    showSuccessAnimation() {
        this.dropzone.classList.add('upload__dropzone--success');
        setTimeout(() => {
            this.dropzone.classList.remove('upload__dropzone--success');
        }, 500);
    }

    /**
     * Show empty state (AC-8)
     */
    showEmptyState() {
        this.dropzone.style.display = 'flex';
    }

    /**
     * Hide empty state (AC-8)
     */
    hideEmptyState() {
        // Don't hide dropzone completely, just make it smaller
        // (Future enhancement: could collapse to compact mode)
    }

    /**
     * Show batch controls (Story 10-3)
     */
    showBatchControls() {
        const batchControls = document.getElementById('batch-controls');
        if (batchControls) {
            batchControls.removeAttribute('hidden');
        }
    }

    /**
     * Hide batch controls (Story 10-3)
     */
    hideBatchControls() {
        const batchControls = document.getElementById('batch-controls');
        if (batchControls) {
            batchControls.setAttribute('hidden', '');
        }
    }

    /**
     * Update file status (Story 10-3, AC-2, AC-3)
     * Updates file card UI based on conversion state
     * @param {number} fileId - File ID to update
     * @param {string} status - Status ('queued', 'processing', 'complete', 'error')
     * @param {string} errorMessage - Error message (optional, for error state)
     */
    updateFileStatus(fileId, status, errorMessage = null) {
        const fileData = this.uploadedFiles.get(fileId);
        if (!fileData) {
            console.error('File not found:', fileId);
            return;
        }

        // Update file data
        fileData.status = status;
        if (errorMessage) {
            fileData.error = errorMessage;
        }

        // Get card elements
        const card = document.getElementById(`file-${fileId}`);
        if (!card) {
            console.error('Card element not found:', fileId);
            return;
        }

        const statusIcon = card.querySelector('.file-card__status-icon');
        const statusText = card.querySelector('.file-card__status-text');
        const downloadBtn = card.querySelector('.file-card__download');
        const convertBtn = card.querySelector('.file-card__convert');
        const retryBtn = card.querySelector('.file-card__retry');
        const cancelBtn = card.querySelector('.file-card__cancel');

        // Update card data-status attribute (triggers CSS transition)
        card.setAttribute('data-status', status);

        // Update icon and text based on status
        switch (status) {
            case 'queued':
                statusIcon.textContent = '⏱️';
                statusText.textContent = 'Queued';
                downloadBtn.setAttribute('hidden', '');
                convertBtn.removeAttribute('hidden');
                retryBtn.setAttribute('hidden', '');
                cancelBtn.setAttribute('hidden', '');
                break;
            case 'processing':
                statusIcon.textContent = '⏳';
                statusText.textContent = 'Converting...';
                downloadBtn.setAttribute('hidden', '');
                convertBtn.setAttribute('hidden', '');
                retryBtn.setAttribute('hidden', '');
                cancelBtn.removeAttribute('hidden'); // Show cancel button (Story 10-4, AC-7)
                break;
            case 'complete':
                statusIcon.textContent = '✓';
                statusText.textContent = 'Complete';
                downloadBtn.removeAttribute('hidden');
                convertBtn.setAttribute('hidden', '');
                retryBtn.setAttribute('hidden', '');
                cancelBtn.setAttribute('hidden', '');
                break;
            case 'error':
                statusIcon.textContent = '✕';
                const truncatedError = this.truncateError(errorMessage || 'Conversion failed', 50);
                statusText.textContent = `Error: ${truncatedError}`;
                statusText.title = errorMessage || 'Conversion failed'; // Full message on hover
                downloadBtn.setAttribute('hidden', '');
                convertBtn.setAttribute('hidden', '');
                retryBtn.removeAttribute('hidden'); // Show retry button (Story 10-4, AC-5)
                cancelBtn.setAttribute('hidden', '');
                break;
        }
    }

    /**
     * Truncate error message (Story 10-3, AC-6)
     * @param {string} message - Error message
     * @param {number} maxLength - Maximum length
     * @returns {string} Truncated message
     */
    truncateError(message, maxLength) {
        if (message.length <= maxLength) {
            return message;
        }
        return message.substring(0, maxLength - 3) + '...';
    }

    /**
     * Update batch progress (Story 10-3, AC-1)
     * Updates progress bar and title text
     * @param {number} completed - Number of completed files
     * @param {number} total - Total number of files
     */
    updateBatchProgress(completed, total) {
        const batchProgress = document.getElementById('batch-progress');
        const title = document.getElementById('batch-progress-title');
        const fill = document.getElementById('batch-progress-fill');
        const cancelBtn = document.getElementById('cancel-batch');

        if (total === 0) {
            batchProgress.setAttribute('hidden', '');
            return;
        }

        batchProgress.removeAttribute('hidden');
        title.textContent = `Converting ${completed} of ${total}...`;
        cancelBtn.removeAttribute('hidden');

        const percentage = (completed / total) * 100;
        fill.style.width = `${percentage}%`;
    }

    /**
     * Show batch completion message (Story 10-3, AC-1)
     * @param {number} successCount - Number of successfully converted files
     * @param {number} errorCount - Number of failed files (optional)
     */
    showBatchComplete(successCount, errorCount = 0) {
        const title = document.getElementById('batch-progress-title');
        const cancelBtn = document.getElementById('cancel-batch');

        if (errorCount === 0) {
            title.textContent = `All ${successCount} files converted successfully`;
        } else {
            const total = successCount + errorCount;
            title.textContent = `${successCount} of ${total} converted, ${errorCount} failed`;
        }

        cancelBtn.setAttribute('hidden', '');
    }

    /**
     * Show batch cancelled message (Story 10-3, AC-7)
     * @param {number} completed - Number of completed files before cancellation
     * @param {number} total - Total number of files
     */
    showBatchCancelled(completed, total) {
        const title = document.getElementById('batch-progress-title');
        const cancelBtn = document.getElementById('cancel-batch');

        title.textContent = `Conversion cancelled (${completed} of ${total} complete)`;
        cancelBtn.setAttribute('hidden', '');
    }

    /**
     * Show preview for uploaded file (Story 11-1, AC-3, H1, H2)
     * Extracts parameters from preset file and applies CSS filter preview
     * @param {number} fileId - File ID to preview
     */
    async showPreviewForFile(fileId) {
        const fileData = this.uploadedFiles.get(fileId);
        if (!fileData) {
            console.error('File not found:', fileId);
            return;
        }

        console.log(`Showing preview for file ${fileId}: ${fileData.file.name}`);

        try {
            // Check if WASM is ready
            if (typeof extractParameters !== 'function') {
                throw new Error('WASM module not loaded. Please refresh the page.');
            }

            // Read file as Uint8Array
            const arrayBuffer = await fileData.file.arrayBuffer();
            const uint8Array = new Uint8Array(arrayBuffer);

            // Extract parameters using WASM
            const jsonString = await extractParameters(uint8Array, fileData.format);
            const recipe = JSON.parse(jsonString);

            console.log('Extracted parameters:', recipe);

            // Apply CSS filter preview
            applyPreviewFilter(recipe);

            // Show preview modal (Story 11-3 will implement the modal UI)
            const modal = document.getElementById('preview-modal');
            if (modal) {
                modal.removeAttribute('hidden');
                modal.style.display = 'flex';
            }

            // Announce to screen readers (Story 10-7, AC-2)
            announceStatus(`Preview loaded for ${fileData.file.name}`);

        } catch (error) {
            console.error(`Preview failed for file ${fileId}:`, error);

            // Show user-friendly error message
            const errorMessage = error.message || 'Unable to generate preview';
            announceError(`Preview error: ${errorMessage}`);

            // Optionally show error in UI (could use a toast notification)
            alert(`Preview Error: ${errorMessage}\n\nPlease ensure the file is a valid preset file.`);
        }
    }

    /**
     * Convert individual file (Story 10-4, AC-2, AC-3)
     * Converts a single file to its selected target format
     * @param {number} fileId - File ID to convert
     */
    async convertIndividualFile(fileId) {
        const fileData = this.uploadedFiles.get(fileId);
        if (!fileData) {
            console.error('File not found:', fileId);
            return;
        }

        // Check if file is already converted or processing
        if (fileData.status === 'complete') {
            console.warn('File already converted:', fileId);
            return;
        }

        if (fileData.status === 'processing') {
            console.warn('File conversion already in progress:', fileId);
            return;
        }

        // Clear previous error (Story 10-4, AC-5 - retry functionality)
        if (fileData.status === 'error') {
            fileData.error = null;
        }

        console.log(`Converting file ${fileId}: ${fileData.format} → ${fileData.targetFormat}`);

        // Create AbortController for cancellation support (Story 10-4, AC-7)
        fileData.abortController = new AbortController();

        // Update status to processing
        this.updateFileStatus(fileId, 'processing');

        try {
            // Read file as Uint8Array
            const arrayBuffer = await fileData.file.arrayBuffer();
            const uint8Array = new Uint8Array(arrayBuffer);

            // Convert file using WASM
            const outputData = await convertFile(
                uint8Array,
                fileData.format,
                fileData.targetFormat,
                fileData.file.name
            );

            // Store output data
            fileData.outputData = outputData;
            fileData.outputFormat = fileData.targetFormat;

            // Update status to complete
            this.updateFileStatus(fileId, 'complete');

            console.log(`File ${fileId} converted successfully`);

        } catch (error) {
            console.error(`Conversion failed for file ${fileId}:`, error);

            // Map error to user-friendly message
            const errorMessage = this.mapConversionError(error);

            // Update status to error
            this.updateFileStatus(fileId, 'error', errorMessage);
        }

        // Update batch progress (Story 10-4, AC-6)
        // Count files by status to show overall progress
        this.updateOverallProgress();
    }

    /**
     * Cancel individual file conversion (Story 10-4, AC-7)
     * Aborts conversion for a single file and reverts status to queued
     * @param {number} fileId - File ID to cancel
     */
    cancelIndividualConversion(fileId) {
        const fileData = this.uploadedFiles.get(fileId);
        if (!fileData) {
            console.error('File not found:', fileId);
            return;
        }

        // Only cancel if file is currently processing
        if (fileData.status !== 'processing') {
            console.warn('File is not processing, cannot cancel:', fileId);
            return;
        }

        console.log(`Cancelling conversion for file ${fileId}`);

        // Abort the conversion if AbortController exists
        if (fileData.abortController) {
            fileData.abortController.abort();
            fileData.abortController = null;
        }

        // Revert status to queued (Story 10-4, AC-7)
        this.updateFileStatus(fileId, 'queued');

        console.log(`File ${fileId} conversion cancelled, status reverted to queued`);
    }

    /**
     * Update overall batch progress (Story 10-4, AC-6)
     * Counts all converted files (both batch and individual) to show progress
     */
    updateOverallProgress() {
        const files = Array.from(this.uploadedFiles.values());
        const total = files.length;
        const completed = files.filter(f => f.status === 'complete').length;
        const errors = files.filter(f => f.status === 'error').length;

        // If all files are converted (complete or error), show completion message
        if (completed + errors === total && total > 0) {
            this.showBatchComplete(completed, errors);
        } else if (completed > 0 || errors > 0) {
            // Show progress for partially converted files
            this.updateBatchProgress(completed + errors, total);
        }
    }

    /**
     * Map conversion error to user-friendly message (Story 10-4, AC-5)
     * @param {Error} error - Error object from conversion
     * @returns {string} User-friendly error message
     */
    mapConversionError(error) {
        // Check if error has userMessage property (from ConversionError)
        if (error.userMessage) {
            return error.userMessage;
        }

        // Map common error patterns
        const errorMappings = {
            'magic bytes': 'File appears corrupted or invalid',
            'parse error': 'Unable to parse file. File may be corrupted',
            'syntax error': 'Invalid file format',
            'not loaded': 'Converter not ready. Please refresh the page',
            'already in progress': 'Please wait for current conversion to complete',
            'corrupted': 'File appears corrupted',
            'unsupported': 'File version not supported'
        };

        const errorText = error.message || error.toString();
        for (const [pattern, message] of Object.entries(errorMappings)) {
            if (errorText.toLowerCase().includes(pattern)) {
                return message;
            }
        }

        // Default fallback
        return 'Conversion failed. File may be corrupted or unsupported.';
    }

    /**
     * Download file (Story 10-3, AC-5)
     * Triggers browser download of converted file
     * @param {number} fileId - File ID to download
     */
    downloadFile(fileId) {
        const fileData = this.uploadedFiles.get(fileId);
        if (!fileData || !fileData.outputData) {
            console.error('No output data available for file:', fileId);
            return;
        }

        // Generate filename: original_converted.ext
        const originalName = fileData.file.name.split('.').slice(0, -1).join('.');
        const extension = this.getFormatExtension(fileData.outputFormat);
        const filename = `${originalName}_converted.${extension}`;

        // Create Blob and trigger download
        const blob = new Blob([fileData.outputData], { type: 'application/octet-stream' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        a.click();
        URL.revokeObjectURL(url);

        console.log('File downloaded:', filename);
    }

    /**
     * Get format file extension (Story 10-3, AC-5)
     * @param {string} format - Format name
     * @returns {string} File extension
     */
    getFormatExtension(format) {
        const extensions = {
            'np3': 'np3',
            'xmp': 'xmp',
            'lrtemplate': 'lrtemplate',
            'costyle': 'costyle',
            'dcp': 'dcp'
        };
        return extensions[format] || 'bin';
    }
}

// Initialize upload manager on page load
let uploadManager;

export function initializeUploadManager() {
    uploadManager = new UploadManager();
    console.log('UploadManager initialized');
    return uploadManager;
}

export function getUploadManager() {
    return uploadManager;
}

// Auto-initialize if DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeUploadManager);
} else {
    initializeUploadManager();
}
