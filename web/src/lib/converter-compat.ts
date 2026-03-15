/**
 * Compatibility wrapper for legacy components that import from "../converter".
 * Delegates to wasm.svelte.ts functions.
 * Will be removed when components are migrated in Phase 3.
 */

import {
	isWasmReady,
	wasmConvert,
	wasmGenerate,
	wasmExtractFullRecipe,
} from "./wasm.svelte";

let convertedFileData: Uint8Array | null = null;
let convertedFileName: string | null = null;
let convertedFileFormat: string | null = null;
let isConverting = false;

function generateConvertedFileName(
	originalFileName: string,
	targetFormat: string,
): string {
	const baseName = originalFileName.replace(
		/\.(np3|xmp)$/i,
		"",
	);
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
): Promise<Uint8Array> {
	if (!isWasmReady()) throw new Error("WASM module not loaded");
	if (isConverting) throw new Error("Conversion already in progress");

	isConverting = true;
	try {
		const outputData = await wasmConvert(fileData, sourceFormat, targetFormat);
		convertedFileData = outputData;
		convertedFileFormat = targetFormat;
		convertedFileName = generateConvertedFileName(
			originalFileName,
			targetFormat,
		);
		return outputData;
	} finally {
		isConverting = false;
	}
}

export async function generatePreset(
	recipe: Record<string, unknown>,
): Promise<Uint8Array> {
	if (!isWasmReady()) throw new Error("WASM module not loaded");
	return await wasmGenerate(JSON.stringify(recipe));
}

export async function extractFullRecipe(
	fileData: Uint8Array,
	format: string,
): Promise<Record<string, unknown>> {
	if (!isWasmReady()) throw new Error("WASM module not loaded");
	return (await wasmExtractFullRecipe(fileData, format)) as Record<
		string,
		unknown
	>;
}

export function getConvertedFileData(): Uint8Array | null {
	return convertedFileData;
}
export function getConvertedFileName(): string | null {
	return convertedFileName;
}
export function getConvertedFileFormat(): string | null {
	return convertedFileFormat;
}
export function clearConvertedData(): void {
	convertedFileData = null;
	convertedFileName = null;
	convertedFileFormat = null;
}
export function isConversionInProgress(): boolean {
	return isConverting;
}
