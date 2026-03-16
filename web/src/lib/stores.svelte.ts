/**
 * Web app state management using Svelte 5 runes.
 * Shared cross-island state lives in shared-stores.ts (nanostores).
 * This module handles island-local state (undo/redo history, etc.)
 * and provides convenience functions that write to both.
 */

import type { UniversalRecipe } from "@recipe/ui";
import {
	editorModeStore,
	currentRecipeStore,
	originalRecipeStore,
	currentFileNameStore,
	previewImageStore,
} from "./shared-stores";

const MAX_HISTORY = 50;

// Island-local undo/redo state (no need to share across islands)
class EditorHistory {
	recipeHistory = $state<UniversalRecipe[]>([]);
	historyIndex = $state(-1);

	canUndo = $derived(this.historyIndex > 0);
	canRedo = $derived(this.historyIndex < this.recipeHistory.length - 1);
	isDirty = $derived.by(() => {
		const current = currentRecipeStore.get();
		const original = originalRecipeStore.get();
		return current !== null && original !== null &&
			JSON.stringify(current) !== JSON.stringify(original);
	});
}

export const history = new EditorHistory();

// Recipe actions
export function openPreset(
	recipe: UniversalRecipe,
	fileName: string,
): void {
	currentRecipeStore.set(structuredClone(recipe));
	originalRecipeStore.set(structuredClone(recipe));
	currentFileNameStore.set(fileName);
	editorModeStore.set(true);
	history.recipeHistory = [structuredClone(recipe)];
	history.historyIndex = 0;
}

export function updateParameter(key: string, value: number): void {
	const current = currentRecipeStore.get();
	if (!current) return;

	const updated = structuredClone(current);

	// Discard future states when branching
	if (history.historyIndex < history.recipeHistory.length - 1) {
		history.recipeHistory = history.recipeHistory.slice(
			0,
			history.historyIndex + 1,
		);
	}

	const parts = key.split(".");
	if (parts.length === 1) {
		(updated as Record<string, unknown>)[key] = value;
	} else {
		let obj = updated as Record<string, unknown>;
		for (let i = 0; i < parts.length - 1; i++) {
			if (!obj[parts[i]] || typeof obj[parts[i]] !== "object") {
				obj[parts[i]] = {};
			}
			obj = obj[parts[i]] as Record<string, unknown>;
		}
		obj[parts[parts.length - 1]] = value;
	}

	currentRecipeStore.set(updated);

	history.recipeHistory.push(structuredClone(updated));
	if (history.recipeHistory.length > MAX_HISTORY) {
		history.recipeHistory.shift();
	} else {
		history.historyIndex++;
	}
}

export function undo(): void {
	if (!history.canUndo) return;
	history.historyIndex--;
	currentRecipeStore.set(structuredClone(
		history.recipeHistory[history.historyIndex],
	));
}

export function redo(): void {
	if (!history.canRedo) return;
	history.historyIndex++;
	currentRecipeStore.set(structuredClone(
		history.recipeHistory[history.historyIndex],
	));
}

export function resetRecipe(): void {
	const original = originalRecipeStore.get();
	if (!original) return;
	currentRecipeStore.set(structuredClone(original));
	history.recipeHistory = [structuredClone(original)];
	history.historyIndex = 0;
}

export function setPreviewImage(img: HTMLImageElement | null): void {
	previewImageStore.set(img);
}

export function closeEditor(): void {
	editorModeStore.set(false);
}
