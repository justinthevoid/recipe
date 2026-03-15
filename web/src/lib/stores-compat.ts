/**
 * Compatibility layer: re-exports Svelte writable stores for legacy components.
 * These components use $store syntax which requires actual Svelte stores.
 * Will be removed when components are migrated in Phase 3.
 */

import { writable } from "svelte/store";

// WASM state — legacy components still read this
export const wasmState = writable({
	status: "initializing" as "initializing" | "ready" | "error",
	error: null as string | null,
	version: "...",
});

// Uploaded files
export interface LegacyFile {
	id: number;
	file: File;
	name: string;
	size: number;
	format: string | null;
	status: "queued" | "processing" | "complete" | "error";
	progress: number;
	outputData: Uint8Array | null;
	outputFormat: string | null;
}

export const files = writable<LegacyFile[]>([]);

// Settings
export const settings = writable({
	targetFormat: "",
});

// Preview state
export const previewFile = writable<LegacyFile | null>(null);

// Editor recipe
export const currentRecipe = writable<Record<string, unknown> | null>(null);

// Helper functions
let fileIdCounter = 0;

export function addFile(file: File): number {
	const id = fileIdCounter++;
	files.update((list) => [
		...list,
		{
			id,
			file,
			name: file.name,
			size: file.size,
			status: "queued" as const,
			format: null,
			progress: 0,
			outputData: null,
			outputFormat: null,
		},
	]);
	return id;
}

export function updateFileStatus(
	id: number,
	updates: Partial<LegacyFile>,
): void {
	files.update((list) =>
		list.map((f) => (f.id === id ? { ...f, ...updates } : f)),
	);
}

export function removeFile(id: number): void {
	files.update((list) => list.filter((f) => f.id !== id));
}
