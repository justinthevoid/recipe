/**
 * Format detection for preset files.
 */

const EXTENSION_MAP: Record<string, string> = {
	np3: "np3",
	xmp: "xmp",
};

/**
 * Detect format from filename extension.
 */
export function detectFormatFromExtension(filename: string): string {
	const ext = filename.split(".").pop()?.toLowerCase() ?? "";
	return EXTENSION_MAP[ext] ?? "unknown";
}

/**
 * Get the opposite format for conversion.
 */
export function getOppositeFormat(format: string): string {
	return format === "np3" ? "xmp" : "np3";
}

/**
 * Get human-readable format label.
 */
export function getFormatLabel(format: string): string {
	return format === "np3" ? "NP3" : "XMP";
}
