import type { CopyPastePayload, Np3Error, Np3OpenResponse } from "../types";
import { NP3_SCHEMA_VERSION } from "../types";

const MAX_HISTORY_SIZE = 100;

/**
 * Deep-clone helper. Uses structuredClone and $state.snapshot()
 * to handle Svelte 5's $state reactive proxies efficiently.
 */
function deepClone<T>(value: T): T {
	if (!value || typeof value !== "object") return value;
	return structuredClone($state.snapshot(value)) as T;
}

function getVsCode() {
	return (window as Window & { acquireVsCodeApi?: () => { postMessage: (msg: unknown) => void } }).acquireVsCodeApi?.() || { postMessage: () => { } };
}

function validateRecipe(recipe: unknown, allowedKeys?: string[]): boolean {
	if (!recipe || typeof recipe !== "object") return false;
	const r = recipe as Record<string, unknown>;

	// If we have allowedKeys, ensure no extra junk is being pasted
	if (allowedKeys && allowedKeys.length > 0) {
		const keys = Object.keys(r);
		for (const key of keys) {
			// Basic metadata fields are always allowed
			if (key === 'name' || key === 'description' || key === 'version' || key === 'pointCurve') continue;
			if (!allowedKeys.includes(key)) return false;
		}
	}

	// Basic check for common fields or pointCurve structure
	if (r.pointCurve && !Array.isArray(r.pointCurve)) return false;
	return true;
}

export class Np3Store {
	isCorrupted = $state<boolean>(false);
	isLoaded = $state<boolean>(false);
	currentError = $state<Np3Error | null>(null);
	metadata = $state<Np3OpenResponse | null>(null);
	currentFilename = $state<string>("Unknown File");

	// Undo/Redo stacks (capped at MAX_HISTORY_SIZE entries)
	#undoStack = $state<Np3OpenResponse[]>([]);
	#redoStack = $state<Np3OpenResponse[]>([]);

	// Derived properties for UI state
	canUndo = $derived(this.#undoStack.length > 0);
	canRedo = $derived(this.#redoStack.length > 0);

	loadSuccess(response: Np3OpenResponse, filename?: string) {
		this.metadata = deepClone(response);
		if (filename) this.currentFilename = filename;
		this.isLoaded = true;
		this.isCorrupted = false;
		this.currentError = null;
		this.#undoStack = [];
		this.#redoStack = [];
	}

	loadError(error: Np3Error, filename?: string) {
		this.currentError = error;
		if (filename) this.currentFilename = filename;
		this.isCorrupted = true;
		this.metadata = null;
	}

	clearError() {
		this.isCorrupted = false;
		this.currentError = null;
	}

	patch(updater: (meta: Np3OpenResponse) => Np3OpenResponse) {
		if (!this.metadata) return;

		// Clear transient errors on new user action
		if (this.currentError?.code.startsWith("PASTE_")) {
			this.currentError = null;
		}

		// Save current state to undo stack
		this.#undoStack.push(deepClone(this.metadata));
		if (this.#undoStack.length > MAX_HISTORY_SIZE) {
			this.#undoStack = this.#undoStack.slice(-MAX_HISTORY_SIZE);
		}

		// Clear redo stack on new change
		this.#redoStack = [];

		// Apply update optimistically
		this.metadata = updater(this.metadata);
	}

	rollback() {
		if (this.#undoStack.length > 0) {
			const previous = this.#undoStack.pop();
			if (previous) {
				this.metadata = deepClone(previous);
			}
		}
	}

	undo() {
		if (!this.metadata || this.#undoStack.length === 0) return;

		// Push current state to redo stack
		this.#redoStack.push(deepClone(this.metadata));

		// Revert to last undo state
		const previous = this.#undoStack.pop();
		if (previous) {
			this.metadata = deepClone(previous);
		}
	}

	redo() {
		if (!this.metadata || this.#redoStack.length === 0) return;

		// Push current state to undo stack
		this.#undoStack.push(deepClone(this.metadata));

		// Apply last redo state
		const next = this.#redoStack.pop();
		if (next) {
			this.metadata = deepClone(next);
		}
	}
	saveAsSuccessful(filePath: string) {
		this.currentFilename = filePath.split(/[/\\]/).pop() || filePath;
		this.currentError = null;
	}

	async copyParameters() {
		if (!this.metadata) return;

		const payload: CopyPastePayload = {
			version: NP3_SCHEMA_VERSION,
			recipe: $state.snapshot(this.metadata.recipe),
		};

		getVsCode().postMessage({
			type: "np3.copy",
			payload: JSON.stringify(payload),
		});
	}

	async pasteParameters() {
		getVsCode().postMessage({ type: "np3.paste_request" });
	}

	handlePaste(clipboardData: string): boolean {
		try {
			const data = JSON.parse(clipboardData) as CopyPastePayload;

			if (!data || typeof data !== "object" || data.version === undefined) {
				throw new Error("Invalid format");
			}

			if (data.version !== NP3_SCHEMA_VERSION) {
				this.currentError = {
					message: "Copy/Paste version mismatch. Please ensure both editors are up to date.",
					code: "PASTE_VERSION_MISMATCH",
				};
				return false;
			}

			const allowedKeys = this.metadata?.parameterDefinitions?.map(p => p.key) || [];
			if (!validateRecipe(data.recipe, allowedKeys)) {
				throw new Error("Invalid recipe data or unknown parameters");
			}

			// Apply the recipe data
			this.patch((meta) => {
				// Optimization: We don't need a full deepClone here since patch() handles history
				// and structuredClone/snapshot is expensive.
				const updatedRecipe = {
					...meta.recipe,
					...data.recipe,
				};
				return { ...meta, recipe: updatedRecipe };
			});

			this.currentError = null;
			return true;
		} catch (e) {
			this.currentError = {
				message: (e instanceof Error) ? e.message : "Invalid clipboard data. Please copy from the NP3 Editor.",
				code: (e instanceof Error && e.message.includes("version mismatch")) ? "PASTE_VERSION_MISMATCH" : "PASTE_INVALID_DATA",
			};
			return false;
		}
	}
}

export const np3AppStore = new Np3Store();
