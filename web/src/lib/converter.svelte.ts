/**
 * Conversion wrapper for web context.
 * Wraps WASM functions with validation, error handling, and download support.
 */

import type { UniversalRecipe } from "@recipe/ui";
import { isWasmReady, wasmConvert, wasmGenerate } from "./wasm.svelte";

class ConversionError extends Error {
	userMessage: string;

	constructor(message: string) {
		super(message);
		this.name = "ConversionError";
		this.userMessage = getUserFriendlyErrorMessage(message);
	}
}

function getUserFriendlyErrorMessage(technical: string): string {
	const mappings: Record<string, string> = {
		"NP3 magic bytes": "File appears corrupted or not a valid NP3 preset.",
		"XMP parse error": "Unable to parse XMP file. File may be corrupted.",
		"Invalid XMP format": "Unable to parse XMP file. File may be corrupted.",
		"WASM not ready": "Converter not ready. Please refresh the page.",
	};

	for (const [key, friendly] of Object.entries(mappings)) {
		if (technical.includes(key)) return friendly;
	}

	return "Conversion failed. File may be corrupted or unsupported.";
}

function validateOutput(data: Uint8Array, format: string): void {
	if (!data || data.length === 0) {
		throw new ConversionError("Converted data is empty");
	}

	switch (format) {
		case "np3":
			if (data[0] !== 0x4e || data[1] !== 0x43 || data[2] !== 0x50) {
				throw new ConversionError("Invalid NP3 magic bytes");
			}
			break;
		case "xmp": {
			const header = new TextDecoder().decode(data.slice(0, 5));
			if (!header.startsWith("<?xml")) {
				throw new ConversionError("Invalid XMP format");
			}
			break;
		}
	}
}

function generateFileName(originalName: string, targetFormat: string): string {
	const baseName = originalName.replace(/\.(np3|xmp)$/i, "");
	const extensions: Record<string, string> = {
		np3: ".np3",
		xmp: ".xmp",
	};
	return `${baseName}${extensions[targetFormat] ?? `.${targetFormat}`}`;
}

export async function convertFile(
	fileData: Uint8Array,
	sourceFormat: string,
	targetFormat: string,
	originalFileName: string,
): Promise<{ data: Uint8Array; fileName: string }> {
	if (!isWasmReady()) {
		throw new ConversionError("WASM not ready");
	}

	if (sourceFormat === targetFormat) {
		throw new ConversionError("Cannot convert to same format");
	}

	try {
		const outputData = await wasmConvert(fileData, sourceFormat, targetFormat);
		validateOutput(outputData, targetFormat);
		return {
			data: outputData,
			fileName: generateFileName(originalFileName, targetFormat),
		};
	} catch (error) {
		if (error instanceof ConversionError) throw error;
		throw new ConversionError(
			error instanceof Error ? error.message : "Unknown conversion error",
		);
	}
}

export async function generatePreset(
	recipe: UniversalRecipe,
): Promise<Uint8Array> {
	if (!isWasmReady()) {
		throw new ConversionError("WASM not ready");
	}

	try {
		return await wasmGenerate(JSON.stringify(recipe));
	} catch (error) {
		throw new ConversionError(
			error instanceof Error ? error.message : "Generation failed",
		);
	}
}

export async function convertAndDownload(
	recipe: UniversalRecipe,
	targetFormat: string,
	presetName: string,
): Promise<void> {
	let data: Uint8Array;

	if (targetFormat === "np3") {
		// Direct generation
		data = await generatePreset(recipe);
	} else {
		// Two-step: recipe → NP3 → target
		const np3Data = await generatePreset(recipe);
		const result = await convertFile(np3Data, "np3", targetFormat, presetName);
		data = result.data;
	}

	// Trigger download
	const blob = new Blob([data]);
	const url = URL.createObjectURL(blob);
	const a = document.createElement("a");
	a.href = url;
	a.download = generateFileName(presetName, targetFormat);
	a.click();
	URL.revokeObjectURL(url);
}

export { ConversionError };
