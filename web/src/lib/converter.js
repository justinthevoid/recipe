// converter.js - WASM conversion wrapper
/* global convert, generate, extractFullRecipe */
// Epic 2, Story 2-6: WASM Conversion Execution
// Handles file conversion using WASM module with error handling and validation

let convertedFileData = null;
let convertedFileName = null;
let convertedFileFormat = null;
let isConverting = false;

/**
 * Convert preset file using WASM
 * @param {Uint8Array} fileData - Source file data
 * @param {string} sourceFormat - "np3" | "xmp" | "lrtemplate"
 * @param {string} targetFormat - "np3" | "xmp" | "lrtemplate"
 * @param {string} originalFileName - Original file name (for output naming)
 * @returns {Promise<Uint8Array>} Converted file data
 */
export async function convertFile(fileData, sourceFormat, targetFormat, originalFileName) {
	if (!fileData || !sourceFormat || !targetFormat) {
		throw new Error("Missing required parameters");
	}

	if (sourceFormat === targetFormat) {
		throw new Error("Cannot convert to same format");
	}

	if (isConverting) {
		throw new Error("Conversion already in progress");
	}

	// Check if WASM is ready
	if (typeof convert !== "function") {
		throw new Error("WASM module not loaded");
	}

	console.log(`Converting ${sourceFormat} → ${targetFormat}...`);
	const startTime = performance.now();

	isConverting = true;

	try {
		// Call WASM function (returns Promise<Uint8Array>)
		const outputData = await convert(fileData, sourceFormat, targetFormat);

		const elapsedTime = performance.now() - startTime;
		console.log(
			`Conversion complete: ${sourceFormat} → ${targetFormat} (${elapsedTime.toFixed(2)}ms, ${outputData.length} bytes)`,
		);

		// Validate output
		validateConvertedData(outputData, targetFormat);

		// Store converted data
		convertedFileData = outputData;
		convertedFileFormat = targetFormat;
		convertedFileName = generateConvertedFileName(originalFileName, targetFormat);

		isConverting = false;

		return outputData;
	} catch (error) {
		isConverting = false;
		console.error("Conversion failed:", error);
		throw new ConversionError(error.message || error);
	}
}

/**
 * Generate NP3 preset from recipe JSON
 * @param {Object} recipe - UniversalRecipe object
 * @returns {Promise<Uint8Array>} Generated NP3 file data
 */
export async function generatePreset(recipe) {
	if (!recipe) {
		throw new Error("Missing recipe data");
	}

	if (typeof generate !== "function") {
		throw new Error("WASM module not loaded");
	}

	try {
		const jsonString = JSON.stringify(recipe);
		const outputData = await generate(jsonString);
		return outputData;
	} catch (error) {
		console.error("Generation failed:", error);
		throw new ConversionError(error.message || error);
	}
}

/**
 * Extract full recipe from file data
 * @param {Uint8Array} fileData - Source file data
 * @param {string} format - Source format
 * @returns {Promise<Object>} Full UniversalRecipe object
 */
export async function extractFullRecipe(fileData, format) {
	if (!fileData || !format) {
		throw new Error("Missing required parameters");
	}

	if (typeof extractFullRecipe !== "function") {
		throw new Error("WASM module not loaded");
	}

	try {
		const jsonString = await window.extractFullRecipe(fileData, format);
		return JSON.parse(jsonString);
	} catch (error) {
		console.error("Extraction failed:", error);
		throw new ConversionError(error.message || error);
	}
}

/**
 * Validate converted data matches expected format
 */
function validateConvertedData(data, format) {
	if (!data || data.length === 0) {
		throw new Error("Converted data is empty");
	}

	try {
		switch (format) {
			case "np3":
				// Check NP3 magic bytes: "NCP" (0x4E 0x43 0x50)
				if (data[0] !== 0x4e || data[1] !== 0x43 || data[2] !== 0x50) {
					throw new Error("Invalid NP3 magic bytes");
				}
				console.log("Validation: NP3 output valid (magic bytes correct)");
				break;

			case "xmp": {
				// Check XMP starts with XML declaration
				const xmpHeader = new TextDecoder().decode(data.slice(0, 5));
				if (!xmpHeader.startsWith("<?xml")) {
					throw new Error("Invalid XMP format (missing XML declaration)");
				}
				console.log("Validation: XMP output valid (XML structure correct)");
				break;
			}

			case "lrtemplate": {
				// Check lrtemplate starts with Lua syntax "s = {"
				const lrtemplateHeader = new TextDecoder().decode(data.slice(0, 10)).trim();
				if (!lrtemplateHeader.startsWith("s = {")) {
					throw new Error("Invalid lrtemplate format (missing Lua syntax)");
				}
				console.log("Validation: lrtemplate output valid (Lua syntax correct)");
				break;
			}

			default:
				throw new Error(`Unknown format: ${format}`);
		}
	} catch (validationError) {
		throw new Error(`Validation failed: ${validationError.message}`);
	}
}

/**
 * Generate output file name based on input file name and target format
 */
function generateConvertedFileName(originalFileName, targetFormat) {
	// Remove original extension
	const baseName = originalFileName.replace(/\.(np3|xmp|lrtemplate)$/i, "");

	// Add new extension
	const extensions = {
		np3: ".np3",
		xmp: ".xmp",
		lrtemplate: ".lrtemplate",
	};

	return `${baseName}${extensions[targetFormat]}`;
}

/**
 * Custom error class for conversion errors
 */
class ConversionError extends Error {
	constructor(message) {
		super(message);
		this.name = "ConversionError";
		this.userMessage = getUserFriendlyErrorMessage(message);
	}
}

/**
 * Map technical error messages to user-friendly messages
 */
function getUserFriendlyErrorMessage(technicalError) {
	const errorMappings = {
		"NP3 magic bytes invalid": "File appears corrupted or not a valid NP3 preset.",
		"Invalid NP3 magic bytes": "File appears corrupted or not a valid NP3 preset.",
		"XMP parse error": "Unable to parse XMP file. File may be corrupted.",
		"Invalid XMP format": "Unable to parse XMP file. File may be corrupted.",
		"lrtemplate syntax error": "Invalid Lightroom preset format.",
		"Invalid lrtemplate format": "Invalid Lightroom preset format.",
		"Unsupported NP3 version": "NP3 preset version not supported.",
		"Unsupported XMP version": "XMP preset version not supported.",
		"WASM module not loaded": "Converter not ready. Please refresh the page.",
		"Conversion already in progress": "Please wait for current conversion to complete.",
	};

	// Check for exact matches
	for (const [technical, friendly] of Object.entries(errorMappings)) {
		if (technicalError.includes(technical)) {
			return friendly;
		}
	}

	// Default fallback
	return "Conversion failed. File may be corrupted or unsupported.";
}

/**
 * Get converted file data
 */
export function getConvertedFileData() {
	return convertedFileData;
}

/**
 * Get converted file name
 */
export function getConvertedFileName() {
	return convertedFileName;
}

/**
 * Get converted file format
 */
export function getConvertedFileFormat() {
	return convertedFileFormat;
}

/**
 * Clear converted data from memory
 */
export function clearConvertedData() {
	convertedFileData = null;
	convertedFileName = null;
	convertedFileFormat = null;
	isConverting = false;
}

/**
 * Check if conversion is in progress
 */
export function isConversionInProgress() {
	return isConverting;
}
