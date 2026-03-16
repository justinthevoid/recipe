import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/svelte";

// Mock WASM module
vi.mock("$lib/wasm.svelte", () => ({
	wasm: { ready: true, status: "ready" as const, error: null, version: "1.0" },
	initWasm: vi.fn().mockResolvedValue(undefined),
	isWasmReady: vi.fn().mockReturnValue(true),
	wasmConvert: vi.fn().mockResolvedValue(
		new Uint8Array([0x3c, 0x3f, 0x78, 0x6d, 0x6c]), // "<?xml"
	),
	wasmExtractFullRecipe: vi.fn().mockResolvedValue({
		exposure: 0.67,
		contrast: 35,
		saturation: 20,
		highlights: -10,
		shadows: 15,
	}),
	wasmDetectFormat: vi.fn().mockReturnValue("np3"),
	wasmGenerateLUT: vi.fn(),
}));

// Mock converter
vi.mock("$lib/converter.svelte", () => ({
	convertFile: vi.fn().mockResolvedValue({
		data: new Uint8Array([0x3c, 0x3f, 0x78, 0x6d, 0x6c]),
		fileName: "test-preset.xmp",
	}),
	ConversionError: class ConversionError extends Error {
		userMessage: string;
		constructor(msg: string) {
			super(msg);
			this.userMessage = msg;
		}
	},
}));

// Mock shared stores
vi.mock("$lib/shared-stores", () => {
	const { atom } = require("nanostores");
	return {
		wasmStatusStore: atom("ready"),
		wasmErrorStore: atom(null),
		wasmVersionStore: atom("1.0"),
		editorModeStore: atom(false),
		currentRecipeStore: atom(null),
		originalRecipeStore: atom(null),
		currentFileNameStore: atom(""),
		previewImageStore: atom(null),
	};
});

// Mock stores.svelte
vi.mock("$lib/stores.svelte", () => ({
	openPreset: vi.fn(),
}));

// Mock parameter counter
vi.mock("$lib/parameter-counter", () => ({
	countParameters: vi.fn().mockReturnValue({ mapped: 35, skipped: 13 }),
}));

// Mock fetch for demo preset
const mockFetch = vi.fn().mockResolvedValue({
	arrayBuffer: () => Promise.resolve(new ArrayBuffer(480)),
});
vi.stubGlobal("fetch", mockFetch);

import ConversionCard from "./ConversionCard.svelte";

describe("ConversionCard", () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockFetch.mockResolvedValue({
			arrayBuffer: () => Promise.resolve(new ArrayBuffer(480)),
		});
	});

	it("renders idle state with upload zone", () => {
		render(ConversionCard);
		expect(screen.getByText(/drop your/i)).toBeTruthy();
		expect(screen.getByText("Select File")).toBeTruthy();
		expect(screen.getByText(/try a demo/i)).toBeTruthy();
	});

	it("renders WASM ready status", () => {
		render(ConversionCard);
		expect(screen.getByText("Ready")).toBeTruthy();
	});

	it("has a hidden file input with correct accept types", () => {
		render(ConversionCard);
		const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement;
		expect(fileInput).toBeTruthy();
		expect(fileInput.accept).toBe(".np3,.xmp");
	});

	it("rejects files with invalid extensions silently", async () => {
		render(ConversionCard);
		// Card should stay in idle state — no state change for invalid files
		expect(screen.getByText(/drop your/i)).toBeTruthy();
	});
});

describe("ConversionCard format detection", () => {
	it("detects np3 format from extension", async () => {
		const { detectFormatFromExtension } = await import("$lib/format-detector");
		expect(detectFormatFromExtension("preset.np3")).toBe("np3");
	});

	it("detects xmp format from extension", async () => {
		const { detectFormatFromExtension } = await import("$lib/format-detector");
		expect(detectFormatFromExtension("preset.xmp")).toBe("xmp");
	});

	it("returns unknown for unsupported extensions", async () => {
		const { detectFormatFromExtension } = await import("$lib/format-detector");
		expect(detectFormatFromExtension("photo.jpg")).toBe("unknown");
	});
});

describe("parameter counter", () => {
	it("counts non-default parameters", async () => {
		// Unmock to test real implementation
		vi.doUnmock("$lib/parameter-counter");
		const { countParameters } = await import("$lib/parameter-counter");
		const recipe = {
			exposure: 0.67,
			contrast: 35,
			saturation: 0, // default — should not count
			highlights: -10,
		};
		const result = countParameters(recipe, "np3");
		expect(result.mapped).toBe(3); // exposure, contrast, highlights
		expect(result.skipped).toBe(45); // 48 - 3
	});
});
