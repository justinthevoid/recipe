import type { Np3Error, Np3OpenResponse } from "../types";

const MAX_HISTORY_SIZE = 20;

/**
 * Deep-clone helper. Uses JSON serialization because structuredClone
 * cannot handle Svelte 5's $state reactive proxies.
 */
function deepClone<T>(value: T): T {
	return JSON.parse(JSON.stringify(value));
}

export class Np3Store {
	isCorrupted = $state<boolean>(false);
	isLoaded = $state<boolean>(false);
	currentError = $state<Np3Error | null>(null);
	metadata = $state<Np3OpenResponse | null>(null);
	currentFilename = $state<string>("Unknown File");

	// Rollback buffer for optimistic UI (capped at MAX_HISTORY_SIZE entries)
	#history = $state<Np3OpenResponse[]>([]);

	loadSuccess(response: Np3OpenResponse, filename?: string) {
		this.metadata = deepClone(response);
		if (filename) this.currentFilename = filename;
		this.isLoaded = true;
		this.isCorrupted = false;
		this.currentError = null;
		this.#history = [deepClone(response)];
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

		// Save current state to history, cap at MAX_HISTORY_SIZE
		this.#history.push(deepClone(this.metadata));
		if (this.#history.length > MAX_HISTORY_SIZE) {
			this.#history = this.#history.slice(-MAX_HISTORY_SIZE);
		}

		// Apply update optimistically
		this.metadata = updater(this.metadata);
	}

	rollback() {
		if (this.#history.length > 0) {
			const previous = this.#history.pop();
			if (previous) {
				this.metadata = deepClone(previous);
			}
		}
	}
}

export const np3AppStore = new Np3Store();
