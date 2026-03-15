const MAX_DIMENSION = 2048;
const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
const ACCEPTED_TYPES = ["image/jpeg", "image/png", "image/webp"];

export async function loadImage(file: File): Promise<HTMLImageElement> {
	if (!ACCEPTED_TYPES.includes(file.type)) {
		throw new Error(`Unsupported image type: ${file.type}. Use JPEG, PNG, or WebP.`);
	}

	if (file.size > MAX_FILE_SIZE) {
		console.warn(`Large image (${(file.size / 1024 / 1024).toFixed(1)}MB) — loading may be slow`);
	}

	const url = URL.createObjectURL(file);

	const img = await new Promise<HTMLImageElement>((resolve, reject) => {
		const el = new Image();
		el.onload = () => resolve(el);
		el.onerror = () => reject(new Error("Failed to load image"));
		el.src = url;
	});

	// Downscale if needed
	if (img.naturalWidth > MAX_DIMENSION || img.naturalHeight > MAX_DIMENSION) {
		// Revoke original blob URL — downscale creates a new data URL
		URL.revokeObjectURL(url);
		return await downscale(img);
	}

	// Don't revoke the blob URL — WebGL needs it alive for texImage2D.
	// The URL will be garbage collected when the image element is released.
	return img;
}

async function downscale(img: HTMLImageElement): Promise<HTMLImageElement> {
	const canvas = document.createElement("canvas");
	const scale = Math.min(
		MAX_DIMENSION / img.naturalWidth,
		MAX_DIMENSION / img.naturalHeight,
	);

	canvas.width = Math.round(img.naturalWidth * scale);
	canvas.height = Math.round(img.naturalHeight * scale);

	const ctx = canvas.getContext("2d");
	if (!ctx) throw new Error("Canvas 2D context unavailable");

	ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

	const dataUrl = canvas.toDataURL("image/jpeg", 0.92);

	// Wait for the new image to fully load before returning
	return new Promise<HTMLImageElement>((resolve, reject) => {
		const result = new Image();
		result.onload = () => resolve(result);
		result.onerror = () => reject(new Error("Failed to create downscaled image"));
		result.src = dataUrl;
	});
}

export async function loadSampleImage(path: string): Promise<HTMLImageElement> {
	return new Promise<HTMLImageElement>((resolve, reject) => {
		const img = new Image();
		img.onload = () => resolve(img);
		img.onerror = () => reject(new Error(`Failed to load sample image: ${path}`));
		img.src = path;
	});
}
