/**
 * Web app state management using Svelte 5 runes.
 */

import type { UniversalRecipe } from "@recipe/ui";

// File management
export interface UploadedFile {
	id: number;
	file: File;
	name: string;
	size: number;
	format: string | null;
	status: "queued" | "processing" | "complete" | "error";
	outputData: Uint8Array | null;
	outputFormat: string | null;
	errorMessage?: string;
}

let fileIdCounter = 0;

// Reactive store object — all state in one place for reliable reactivity
class AppStore {
	files = $state<UploadedFile[]>([]);
	currentRecipe = $state<UniversalRecipe | null>(null);
	originalRecipe = $state<UniversalRecipe | null>(null);
	previewImage = $state<HTMLImageElement | null>(null);
	editorMode = $state(false);
	targetFormat = $state("xmp");
	currentFileName = $state("");

	// Undo/redo
	recipeHistory = $state<UniversalRecipe[]>([]);
	historyIndex = $state(-1);

	// Derived
	isDirty = $derived(
		this.currentRecipe !== null &&
			this.originalRecipe !== null &&
			JSON.stringify(this.currentRecipe) !==
				JSON.stringify(this.originalRecipe),
	);
	canUndo = $derived(this.historyIndex > 0);
	canRedo = $derived(this.historyIndex < this.recipeHistory.length - 1);
	canConvert = $derived(
		this.currentRecipe !== null && this.files.length > 0,
	);
}

export const store = new AppStore();

const MAX_HISTORY = 50;

// File actions
export function addFile(file: File): number {
	const id = fileIdCounter++;
	store.files.push({
		id,
		file,
		name: file.name,
		size: file.size,
		format: null,
		status: "queued",
		outputData: null,
		outputFormat: null,
	});
	return id;
}

export function updateFileStatus(
	id: number,
	updates: Partial<UploadedFile>,
): void {
	const idx = store.files.findIndex((f) => f.id === id);
	if (idx !== -1) {
		store.files[idx] = { ...store.files[idx], ...updates };
	}
}

export function removeFile(id: number): void {
	const idx = store.files.findIndex((f) => f.id === id);
	if (idx !== -1) {
		store.files.splice(idx, 1);
	}
}

// Recipe actions
export function openPreset(
	recipe: UniversalRecipe,
	fileName: string,
): void {
	store.currentRecipe = structuredClone(recipe);
	store.originalRecipe = structuredClone(recipe);
	store.currentFileName = fileName;
	store.editorMode = true;
	store.recipeHistory = [structuredClone(recipe)];
	store.historyIndex = 0;
}

export function updateParameter(key: string, value: number): void {
	if (!store.currentRecipe) return;

	// Discard future states when branching
	if (store.historyIndex < store.recipeHistory.length - 1) {
		store.recipeHistory = store.recipeHistory.slice(
			0,
			store.historyIndex + 1,
		);
	}

	const parts = key.split(".");
	if (parts.length === 1) {
		(store.currentRecipe as Record<string, unknown>)[key] = value;
	} else {
		let obj = store.currentRecipe as Record<string, unknown>;
		for (let i = 0; i < parts.length - 1; i++) {
			if (!obj[parts[i]] || typeof obj[parts[i]] !== "object") {
				obj[parts[i]] = {};
			}
			obj = obj[parts[i]] as Record<string, unknown>;
		}
		obj[parts[parts.length - 1]] = value;
	}

	store.recipeHistory.push(structuredClone(store.currentRecipe));
	if (store.recipeHistory.length > MAX_HISTORY) {
		store.recipeHistory.shift();
	} else {
		store.historyIndex++;
	}
}

export function undo(): void {
	if (!store.canUndo) return;
	store.historyIndex--;
	store.currentRecipe = structuredClone(
		store.recipeHistory[store.historyIndex],
	);
}

export function redo(): void {
	if (!store.canRedo) return;
	store.historyIndex++;
	store.currentRecipe = structuredClone(
		store.recipeHistory[store.historyIndex],
	);
}

export function resetRecipe(): void {
	if (!store.originalRecipe) return;
	store.currentRecipe = structuredClone(store.originalRecipe);
	store.recipeHistory = [structuredClone(store.originalRecipe)];
	store.historyIndex = 0;
}

export function setPreviewImage(img: HTMLImageElement | null): void {
	store.previewImage = img;
}

export function closeEditor(): void {
	store.editorMode = false;
}

export function setTargetFormat(format: string): void {
	store.targetFormat = format;
}
