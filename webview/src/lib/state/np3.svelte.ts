import type { Np3Error, Np3OpenResponse } from "../types";

const MAX_HISTORY_SIZE = 100;

/**
 * Deep-clone helper. Uses structuredClone and $state.snapshot() 
 * to handle Svelte 5's $state reactive proxies efficiently.
 */
function deepClone<T>(value: T): T {
	return structuredClone($state.snapshot(value)) as T;
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
}

export const np3AppStore = new Np3Store();
