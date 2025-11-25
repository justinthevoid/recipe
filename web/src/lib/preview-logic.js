/**
 * Preview module for CSS filter-based preset previews
 * Maps recipe parameters to CSS filters for instant preview
 */

/**
 * Convert UniversalRecipe parameters to CSS filter string
 * @param {Object} recipe - UniversalRecipe object with parameters
 * @returns {string} CSS filter string or "none"
 */
export function recipeToCSSFilters(recipe) {
    if (!recipe) return 'none';

    const filters = [];

    // Exposure → brightness (range: -2.0 to +2.0 → 0% to 200%)
    if (recipe.exposure != null && recipe.exposure !== 0 && typeof recipe.exposure === 'number') {
        const brightness = clamp((1.0 + recipe.exposure) * 100, 0, 200);
        filters.push(`brightness(${brightness}%)`);
    }

    // Contrast → contrast (range: -1.0 to +1.0 → 0% to 200%)
    if (recipe.contrast != null && recipe.contrast !== 0 && typeof recipe.contrast === 'number') {
        const contrast = clamp((1.0 + recipe.contrast) * 100, 0, 200);
        filters.push(`contrast(${contrast}%)`);
    }

    // Saturation → saturate (range: -1.0 to +1.0 → 0% to 200%)
    if (recipe.saturation != null && recipe.saturation !== 0 && typeof recipe.saturation === 'number') {
        const saturate = clamp((1.0 + recipe.saturation) * 100, 0, 200);
        filters.push(`saturate(${saturate}%)`);
    }

    // Hue → hue-rotate (range: -180 to +180 degrees)
    if (recipe.hue != null && recipe.hue !== 0 && typeof recipe.hue === 'number') {
        const hue = clamp(recipe.hue, -180, 180);
        filters.push(`hue-rotate(${hue}deg)`);
    }

    // Temperature → sepia + hue-rotate (approximation)
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
 * Clamp numeric value to min/max range
 */
export function clamp(value, min, max) {
    return Math.min(Math.max(value, min), max);
}
