// batch-converter.js - Batch conversion orchestration
// Story 10-3: Progress Indicators for Batch Conversions
// Handles batch file conversion with progress updates and cancellation

import { convertFile } from './converter.js';
import { getUploadManager } from './upload.js';

let conversionAborted = false;

/**
 * Start batch conversion (Story 10-3, AC-1, AC-4, AC-6, AC-7)
 * Converts all uploaded files sequentially with progress updates
 * @param {string} targetFormat - Target format for all files
 */
export async function startBatchConversion(targetFormat) {
    const uploadManager = getUploadManager();
    if (!uploadManager) {
        console.error('UploadManager not initialized');
        return;
    }

    // Reset abortion flag
    conversionAborted = false;

    // Get all uploaded files
    const files = Array.from(uploadManager.uploadedFiles.values());
    const total = files.length;

    if (total === 0) {
        console.warn('No files to convert');
        return;
    }

    console.log(`Starting batch conversion: ${total} files to ${targetFormat}`);

    let completed = 0;
    let successCount = 0;
    let errorCount = 0;

    // Initialize batch progress
    uploadManager.updateBatchProgress(completed, total);

    // Convert files sequentially
    for (const fileData of files) {
        // Check if conversion was aborted
        if (conversionAborted) {
            console.log('Batch conversion aborted by user');
            uploadManager.showBatchCancelled(completed, total);
            return;
        }

        // Update status to processing
        uploadManager.updateFileStatus(fileData.id, 'processing');

        try {
            // Read file as Uint8Array
            const arrayBuffer = await fileData.file.arrayBuffer();
            const uint8Array = new Uint8Array(arrayBuffer);

            // Convert file
            const outputData = await convertFile(
                uint8Array,
                fileData.format,
                targetFormat,
                fileData.file.name
            );

            // Store output data in file data
            fileData.outputData = outputData;
            fileData.outputFormat = targetFormat;

            // Update status to complete
            uploadManager.updateFileStatus(fileData.id, 'complete');

            successCount++;

        } catch (error) {
            console.error(`Conversion failed for file ${fileData.file.name}:`, error);

            // Map error to user-friendly message
            const errorMessage = mapConversionError(error);

            // Update status to error
            uploadManager.updateFileStatus(fileData.id, 'error', errorMessage);

            errorCount++;
        }

        // Update progress
        completed++;
        uploadManager.updateBatchProgress(completed, total);
    }

    // Show completion message
    uploadManager.showBatchComplete(successCount, errorCount);

    console.log(`Batch conversion complete: ${successCount} succeeded, ${errorCount} failed`);
}

/**
 * Cancel batch conversion (Story 10-3, AC-7)
 * Sets abortion flag to stop processing remaining files
 */
export function cancelBatchConversion() {
    if (conversionAborted) {
        console.warn('Batch conversion already cancelled');
        return;
    }

    // Show confirmation dialog
    const confirmed = confirm('Cancel all in-progress conversions?\n\nFiles already converted will remain available.');

    if (confirmed) {
        conversionAborted = true;
        console.log('Batch conversion cancellation requested');
    }
}

/**
 * Map conversion error to user-friendly message (Story 10-3, AC-6)
 * @param {Error} error - Conversion error
 * @returns {string} User-friendly error message
 */
function mapConversionError(error) {
    // If error has userMessage (from ConversionError class), use it
    if (error.userMessage) {
        return error.userMessage;
    }

    // Otherwise map common errors
    const errorMessage = error.message || '';

    if (errorMessage.includes('Invalid file format') || errorMessage.includes('magic bytes')) {
        return 'Invalid file format';
    }

    if (errorMessage.includes('parse') || errorMessage.includes('Parse')) {
        return `Parsing failed: ${errorMessage}`;
    }

    if (errorMessage.includes('Conversion failed')) {
        return `Conversion failed: ${errorMessage}`;
    }

    if (errorMessage.includes('Unsupported')) {
        return 'Unsupported format';
    }

    // Default fallback
    return 'Conversion failed';
}

/**
 * Check if batch conversion is aborted
 * @returns {boolean} True if aborted
 */
export function isBatchConversionAborted() {
    return conversionAborted;
}
