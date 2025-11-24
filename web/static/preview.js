/**
 * Preview module for CSS filter-based preset previews
 * Provides instant approximate preview using browser-native CSS filters
 *
 * Story: 11-1-css-filter-mapping
 * Epic: 11 (Image Preview System)
 */

/**
 * Convert UniversalRecipe parameters to CSS filter string
 *
 * Maps recipe parameters to CSS filters for instant preview:
 * - Exposure (-2.0 to +2.0) → brightness (0% to 200%)
 * - Contrast (-1.0 to +1.0) → contrast (0% to 200%)
 * - Saturation (-1.0 to +1.0) → saturate (0% to 200%)
 * - Hue (-180 to +180) → hue-rotate (-180deg to +180deg)
 * - Temperature (-100 to +100) → sepia + hue-rotate (approximation)
 *
 * @param {Object} recipe - UniversalRecipe object with parameters
 * @returns {string} CSS filter string (e.g., "brightness(150%) contrast(130%)") or "none"
 *
 * @example
 * const recipe = { exposure: 0.5, contrast: 0.3, saturation: -0.5 };
 * recipeToCSSFilters(recipe); // "brightness(150%) contrast(130%) saturate(50%)"
 */
export function recipeToCSSFilters(recipe) {
    if (!recipe) return 'none';

    const filters = [];

    // Exposure → brightness (range: -2.0 to +2.0 → 0% to 200%)
    // Formula: brightness = (1.0 + exposure) * 100%
    // Example: Exposure +0.5 → brightness(150%)
    if (recipe.exposure != null && recipe.exposure !== 0 && typeof recipe.exposure === 'number') {
        const brightness = clamp((1.0 + recipe.exposure) * 100, 0, 200);
        filters.push(`brightness(${brightness}%)`);
    }

    // Contrast → contrast (range: -1.0 to +1.0 → 0% to 200%)
    // Formula: contrast = (1.0 + contrast) * 100%
    // Example: Contrast +0.3 → contrast(130%)
    if (recipe.contrast != null && recipe.contrast !== 0 && typeof recipe.contrast === 'number') {
        const contrast = clamp((1.0 + recipe.contrast) * 100, 0, 200);
        filters.push(`contrast(${contrast}%)`);
    }

    // Saturation → saturate (range: -1.0 to +1.0 → 0% to 200%)
    // Formula: saturate = (1.0 + saturation) * 100%
    // Example: Saturation -0.5 → saturate(50%)
    if (recipe.saturation != null && recipe.saturation !== 0 && typeof recipe.saturation === 'number') {
        const saturate = clamp((1.0 + recipe.saturation) * 100, 0, 200);
        filters.push(`saturate(${saturate}%)`);
    }

    // Hue → hue-rotate (range: -180 to +180 degrees)
    // Formula: hue-rotate = hue (direct mapping)
    // Example: Hue +30 → hue-rotate(30deg)
    if (recipe.hue != null && recipe.hue !== 0 && typeof recipe.hue === 'number') {
        const hue = clamp(recipe.hue, -180, 180);
        filters.push(`hue-rotate(${hue}deg)`);
    }

    // Temperature → sepia + hue-rotate (approximation only)
    // Warm temperatures (+) → sepia() + positive hue-rotate
    // Cool temperatures (-) → sepia() + negative hue-rotate
    // Formula: sepia(|temperature| * 0.3) hue-rotate(temperature * 0.5deg)
    // Note: This is an approximation - CSS lacks true color temperature adjustment
    if (recipe.temperature != null && recipe.temperature !== 0 && typeof recipe.temperature === 'number') {
        const temp = clamp(recipe.temperature, -100, 100);
        const sepia = Math.abs(temp) * 0.3;
        const hueShift = temp * 0.5;
        filters.push(`sepia(${sepia})`);
        filters.push(`hue-rotate(${hueShift}deg)`);
    }

    return filters.length > 0 ? filters.join(' ') : 'none';
}

/**
 * Apply CSS filter to preview image element
 *
 * @param {Object} recipe - UniversalRecipe object
 * @returns {void}
 */
export function applyPreviewFilter(recipe) {
    const previewImage = document.getElementById('preview-image');
    if (!previewImage) {
        console.warn('Preview image element not found (id="preview-image")');
        return;
    }

    const filterString = recipeToCSSFilters(recipe);
    previewImage.style.filter = filterString;

    console.log(`Preview filter applied: ${filterString}`);
}

/**
 * Clamp numeric value to min/max range
 *
 * @param {number} value - Value to clamp
 * @param {number} min - Minimum allowed value
 * @param {number} max - Maximum allowed value
 * @returns {number} Clamped value
 *
 * @example
 * clamp(250, 0, 200); // 200
 * clamp(-50, 0, 200); // 0
 * clamp(150, 0, 200); // 150
 */
export function clamp(value, min, max) {
    return Math.min(Math.max(value, min), max);
}

/**
 * Detect CSS filter support using CSS.supports() API
 *
 * @returns {boolean} True if browser supports CSS filters
 */
export function isCSSFilterSupported() {
    // Check if CSS.supports API exists
    if (!window.CSS || !window.CSS.supports) {
        return false;
    }

    // Test CSS filter support
    return CSS.supports('filter', 'brightness(100%)');
}

/**
 * Check browser compatibility and show fallback if needed
 *
 * @returns {boolean} True if browser is compatible
 */
export function checkBrowserCompatibility() {
    if (!isCSSFilterSupported()) {
        console.warn('CSS filters not supported in this browser');

        // Hide preview feature gracefully
        const previewModal = document.getElementById('preview-modal');
        if (previewModal) {
            previewModal.style.display = 'none';
        }

        // Show message to user
        const message = 'Preview not available in this browser. Please use Chrome 18+, Firefox 35+, Safari 9.1+, or Edge 12+.';
        alert(message);

        return false;
    }

    return true;
}

/**
 * Show disclaimer help text explaining CSS filter limitations
 *
 * @returns {void}
 */
export function showDisclaimerHelp() {
    const helpText = `This preview uses CSS filters for instant feedback.

Limitations:
• Approximates exposure, contrast, saturation, hue
• Temperature/tint is simplified (not accurate color science)
• Tone curves not supported in Phase 1
• Actual conversion results may differ

For accurate results, convert the preset and view in your photo editor.`;

    alert(helpText); // Simple alert for Phase 1, modal in future
}

// Initialize event listeners on page load
if (typeof document !== 'undefined') {
    document.addEventListener('DOMContentLoaded', () => {
        // Check browser compatibility
        checkBrowserCompatibility();

        // Attach disclaimer help button event listener
        const helpButton = document.querySelector('.preview-disclaimer__help');
        if (helpButton) {
            helpButton.addEventListener('click', showDisclaimerHelp);
        }
    });
}
