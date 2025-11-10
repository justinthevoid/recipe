// downloader.js - File download handling
// Epic 2, Stories 2-7 & 2-8: File Download Trigger & Error Handling
// Handles download button, Blob creation, and object URL management

import { showError as showErrorPanel } from './error-handler.js';

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

        // Story 2-8: Show centralized error panel
        showErrorPanel('download-failed', error);
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
