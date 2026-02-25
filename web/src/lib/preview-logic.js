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
	if (!recipe) return "none";

	const filters = [];

	// Helper to get value from multiple possible keys
	const getVal = (...keys) => {
		for (const key of keys) {
			if (recipe[key] != null && typeof recipe[key] === "number") {
				return recipe[key];
			}
		}
		return 0;
	};

	// Exposure → brightness (range: -5.0 to +5.0, typically -2 to +2)
	// Map: 0 -> 100%, +1 -> 150%, -1 -> 50%
	const exposure = getVal("Exposure", "Exposure2012");
	if (exposure !== 0) {
		const brightness = clamp((1.0 + exposure) * 100, 0, 200);
		filters.push(`brightness(${brightness}%)`);
	}

	// Contrast → contrast (range: -100 to +100)
	// Map: 0 -> 100%, +100 -> 200%, -100 -> 0%
	// Note: Lightroom Contrast is often -100 to 100, legacy might be -1.0 to 1.0
	// We'll assume -100 to 100 for 2012 process, normalize to -1.0 to 1.0
	let contrast = getVal("Contrast", "Contrast2012");
	if (Math.abs(contrast) > 1.0) contrast /= 100; // Normalize if > 1

	if (contrast !== 0) {
		const contrastVal = clamp((1.0 + contrast) * 100, 0, 200);
		filters.push(`contrast(${contrastVal}%)`);
	}

	// Saturation → saturate (range: -100 to +100)
	// Map: 0 -> 100%, +100 -> 200%, -100 -> 0%
	let saturation = getVal("Saturation");
	if (Math.abs(saturation) > 1.0) saturation /= 100; // Normalize if > 1

	if (saturation !== 0) {
		const saturate = clamp((1.0 + saturation) * 100, 0, 200);
		filters.push(`saturate(${saturate}%)`);
	}

	// Hue (rarely used in basic panel, but supported)
	const hue = getVal("Hue");
	if (hue !== 0) {
		const hueRotate = clamp(hue, -180, 180);
		filters.push(`hue-rotate(${hueRotate}deg)`);
	}

	// Temperature → sepia + hue-rotate (approximation)
	// Range: -100 to +100
	const temperature = getVal("Temperature");
	if (temperature !== 0) {
		const temp = clamp(temperature, -100, 100);
		const sepia = Math.abs(temp) * 0.005; // Scale down effect
		const hueShift = temp * 0.2;
		filters.push(`sepia(${sepia * 100}%)`);
		filters.push(`hue-rotate(${hueShift}deg)`);
	}

	return filters.length > 0 ? filters.join(" ") : "none";
}

/**
 * Clamp numeric value to min/max range
 */
export function clamp(value, min, max) {
	return Math.min(Math.max(value, min), max);
}
