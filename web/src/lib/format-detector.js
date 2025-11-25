/* global detectFormat */
// format-detector.js - WASM format detection wrapper
// Epic 2, Story 2-3: Format Detection
// Provides automatic format detection for preset files (NP3, XMP, lrtemplate, Costyle, Costylepack)

let currentFormat = null;

/**
 * Detect preset file format using WASM
 * @param {Uint8Array} fileData - Raw file bytes
 * @returns {Promise<string>} Format: "np3" | "xmp" | "lrtemplate" | "costyle" | "costylepack"
 */
export async function detectFileFormat(fileData) {
    if (!fileData || fileData.length === 0) {
        throw new Error('No file data provided');
    }

    // Check if WASM is ready
    if (typeof detectFormat !== 'function') {
        throw new Error('WASM module not loaded');
    }

    console.log(`Detecting format for ${fileData.length} bytes...`);
    const startTime = performance.now();

    try {
        // Call WASM function (returns Promise<string>)
        const format = await detectFormat(fileData);

        const elapsedTime = performance.now() - startTime;
        console.log(`Format detected: ${format} (${elapsedTime.toFixed(2)}ms)`);

        // Store for later use
        currentFormat = format;

        return format;

    } catch (error) {
        console.error('Format detection failed:', error);
        throw new Error(`Unable to detect format: ${error.message || error}`);
    }
}

/**
 * Get currently detected format
 * @returns {string|null} "np3" | "xmp" | "lrtemplate" | "costyle" | "costylepack" | null
 */
export function getCurrentFormat() {
    return currentFormat;
}

/**
 * Clear detected format (when new file uploaded)
 */
export function clearFormat() {
    currentFormat = null;
}

/**
 * Get display name for format
 * @param {string} format - "np3" | "xmp" | "lrtemplate" | "costyle" | "costylepack"
 * @returns {string} Human-readable format name
 */
export function getFormatDisplayName(format) {
    const displayNames = {
        'np3': 'NP3 (Nikon Picture Control)',
        'xmp': 'XMP (Lightroom CC)',
        'lrtemplate': 'lrtemplate (Lightroom Classic)',
        'costyle': 'Costyle (Capture One)',
        'costylepack': 'Costylepack (Capture One Bundle)'
    };
    return displayNames[format] || format.toUpperCase();
}

/**
 * Get format badge color
 * @param {string} format
 * @returns {string} CSS class for badge color
 */
export function getFormatBadgeClass(format) {
    const badgeClasses = {
        'np3': 'badge-blue',         // Nikon blue
        'xmp': 'badge-purple',       // Adobe purple
        'lrtemplate': 'badge-teal',  // Lightroom teal
        'costyle': 'badge-purple',   // Capture One purple (brand color)
        'costylepack': 'badge-purple' // Capture One bundle purple
    };
    return badgeClasses[format] || 'badge-gray';
}
/**
 * Simple synchronous format detection from filename extension
 * @param {string} filename
 * @returns {string} Format: "np3" | "xmp" | "lrtemplate" | "costyle" | "dcp" | "unknown"
 */
export function detectFormatFromExtension(filename) {
    const ext = filename.split('.').pop().toLowerCase();
    const map = {
        'np3': 'np3',
        'xmp': 'xmp',
        'lrtemplate': 'lrtemplate',
        'costyle': 'costyle',
        'costylepack': 'costylepack',
        'dcp': 'dcp'
    };
    return map[ext] || 'unknown';
}
