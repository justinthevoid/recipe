/* global extractParameters */
/**
 * Extract and display parameters from preset file
 */

/**
 * Extract parameters from preset file using WASM
 * @param {Uint8Array} fileData - Raw file bytes
 * @param {string} format - Detected format ("np3" | "xmp" | "lrtemplate")
 * @returns {Promise<Object>} Parameters object
 */
export async function extractPresetParameters(fileData, format) {
	if (!fileData || !format) {
		throw new Error("File data and format required");
	}

	// Check if WASM is ready
	if (typeof extractParameters !== "function") {
		throw new Error("WASM module not loaded");
	}

	try {
		// Call WASM function (returns Promise<string> containing JSON)
		const jsonString = await extractParameters(fileData, format);
		return JSON.parse(jsonString);
	} catch (error) {
		console.error("Parameter extraction failed:", error);
		throw new Error(`Unable to extract parameters: ${error.message || error}`);
	}
}

/**
 * Group parameters by category
 */
export function groupParameters(parameters) {
	const groups = {
		"Basic Adjustments": {},
		"Color Adjustments": {},
		"Detail Adjustments": {},
	};

	// Basic adjustments
	const basicParams = [
		"Exposure",
		"Exposure2012",
		"Contrast",
		"Contrast2012",
		"Highlights",
		"Highlights2012",
		"Shadows",
		"Shadows2012",
		"Whites",
		"Whites2012",
		"Blacks",
		"Blacks2012",
	];

	// Color adjustments
	const colorParams = ["Vibrance", "Saturation", "Temperature", "Tint"];

	// Detail adjustments
	const detailParams = [
		"Clarity",
		"Clarity2012",
		"Sharpness",
		"Dehaze",
		"Texture",
		"GrainAmount",
		"GrainSize",
	];

	for (const [key, value] of Object.entries(parameters)) {
		if (value === null || value === undefined) continue;
		if (key === "Name") continue;

		if (basicParams.includes(key)) {
			groups["Basic Adjustments"][key] = value;
		} else if (colorParams.includes(key)) {
			groups["Color Adjustments"][key] = value;
		} else if (detailParams.includes(key)) {
			groups["Detail Adjustments"][key] = value;
		}
	}

	// Remove empty groups
	for (const [category, params] of Object.entries(groups)) {
		if (Object.keys(params).length === 0) {
			delete groups[category];
		}
	}

	return groups;
}

/**
 * Format parameter value for display
 */
export function formatParameterValue(value) {
	if (value === null || value === undefined) return "—";
	if (typeof value === "number") {
		if (value > 0) return `+${value.toFixed(2)}`;
		if (value < 0) return value.toFixed(2);
		return "0";
	}
	return String(value);
}
