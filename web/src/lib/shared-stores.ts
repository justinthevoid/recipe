/**
 * Cross-island state using nanostores.
 * Astro islands get separate module scopes, so Svelte $state singletons
 * are NOT shared between islands. Nanostores atoms are shared.
 *
 * Note: No $ prefix on exports — Svelte reserves $ for its own reactivity.
 */
import { atom } from "nanostores";
import type { UniversalRecipe } from "@recipe/ui";
import type { WasmStatus } from "./wasm.svelte";

// WASM state — shared across all islands
export const wasmStatusStore = atom<WasmStatus>("idle");
export const wasmErrorStore = atom<string | null>(null);
export const wasmVersionStore = atom<string>("...");

// Editor state — shared between ConversionCard and EditorView
export const editorModeStore = atom<boolean>(false);
export const currentRecipeStore = atom<UniversalRecipe | null>(null);
export const originalRecipeStore = atom<UniversalRecipe | null>(null);
export const currentFileNameStore = atom<string>("");
export const previewImageStore = atom<HTMLImageElement | null>(null);
