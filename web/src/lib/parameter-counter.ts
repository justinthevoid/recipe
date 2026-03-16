/**
 * Count mapped and skipped parameters from a converted recipe.
 */

const KNOWN_PARAMS: Record<string, number> = {
	np3: 48,
	xmp: 50,
};

/**
 * Count non-default numeric values in a recipe object.
 * Recursively traverses nested objects (HSL, color grading).
 */
function countNonDefault(
	obj: Record<string, unknown>,
	defaults: Record<string, number> = {},
): number {
	let count = 0;
	for (const [key, value] of Object.entries(obj)) {
		if (typeof value === "number") {
			const defaultVal = defaults[key] ?? 0;
			if (value !== defaultVal) count++;
		} else if (
			value !== null &&
			typeof value === "object" &&
			!Array.isArray(value)
		) {
			count += countNonDefault(value as Record<string, unknown>);
		}
	}
	return count;
}

export function countParameters(
	recipe: Record<string, unknown>,
	sourceFormat: string,
): { mapped: number; skipped: number } {
	const mapped = countNonDefault(recipe);
	const totalSource = KNOWN_PARAMS[sourceFormat] ?? 48;
	const skipped = Math.max(0, totalSource - mapped);
	return { mapped, skipped };
}
