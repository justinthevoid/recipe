import type { CopyPastePayload, Np3Error, Np3OpenResponse, ParameterDefinition } from "../types";
import { NP3_SCHEMA_VERSION } from "../types";

export const GROUP_ORDER = [
	"Basic",
	"Tone Curve",
	"Color Mixer",
	"Color Grading",
	"Detail",
	"Geometry",
	"System",
];

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
	return (
		(
			window as Window & { acquireVsCodeApi?: () => { postMessage: (msg: unknown) => void } }
		).acquireVsCodeApi?.() || { postMessage: () => {} }
	);
}

/**
 * Gets a nested property value using dot notation.
 */
export function getNested(obj: unknown, path: string): any {
	if (!obj || !path) return undefined;
	const parts = path.split(".");
	let curr: any = obj;
	for (const part of parts) {
		if (curr === null || curr === undefined) return undefined;
		curr = curr[part];
	}
	return curr;
}

/**
 * Sets a nested property value using dot notation, returning a new object.
 */
export function setNested(obj: Record<string, any>, path: string, value: any): Record<string, any> {
	if (!obj || !path) return obj;
	const parts = path.split(".");
	const newObj = { ...obj };
	let curr = newObj;

	for (let i = 0; i < parts.length - 1; i++) {
		const part = parts[i];
		// Create nested object if it doesn't exist
		curr[part] = curr[part] ? { ...curr[part] } : {};
		curr = curr[part];
	}

	curr[parts[parts.length - 1]] = value;
	return newObj;
}

function validateRecipe(recipe: unknown, allowedKeys?: string[]): boolean {
	if (!recipe || typeof recipe !== "object") return false;
	const r = recipe as Record<string, unknown>;

	// If we have allowedKeys, ensure no extra junk is being pasted
	if (allowedKeys && allowedKeys.length > 0) {
		const keys = Object.keys(r);
		for (const key of keys) {
			// Basic metadata fields are always allowed
			if (key === "name" || key === "description" || key === "version" || key === "pointCurve")
				continue;
			if (!allowedKeys.includes(key)) return false;
		}
	}

	// Validate metadata field types (F4/P2-3a)
	if (r.name !== undefined && (typeof r.name !== "string" || (r.name as string).length > 256))
		return false;
	if (
		r.description !== undefined &&
		(typeof r.description !== "string" || (r.description as string).length > 256)
	)
		return false;
	if (r.version !== undefined && typeof r.version !== "number") return false;

	// Validate pointCurve structure
	if (r.pointCurve !== undefined) {
		if (!Array.isArray(r.pointCurve)) return false;
		for (const pt of r.pointCurve as unknown[]) {
			if (!pt || typeof pt !== "object") return false;
			const point = pt as Record<string, unknown>;
			if (typeof point.input !== "number" || typeof point.output !== "number") return false;
			if (point.input < 0 || point.input > 255 || point.output < 0 || point.output > 255)
				return false;
		}
	}

	// Validate numeric fields are actually numbers (not strings, objects, arrays)
	for (const [key, val] of Object.entries(r)) {
		if (key === "name" || key === "description" || key === "version" || key === "pointCurve")
			continue;
		// Color adjustment objects and colorGrading are allowed as objects
		if (val !== null && val !== undefined && typeof val === "object") continue;
		// Scalar values should be numbers
		if (val !== null && val !== undefined && typeof val !== "number") return false;
	}

	return true;
}

export class Np3Store {
	metadata = $state<Np3OpenResponse | null>(null);
	isLoaded = $state(false);
	isCorrupted = $state(false);
	currentError = $state<Np3Error | null>(null);
	currentFilename = $state<string | null>(null);

	// Saved state for dirty tracking (F9/F33/F39)
	private savedRecipeSnapshot: string | null = null;

	// History
	undoStack = $state<Np3OpenResponse[]>([]);
	redoStack = $state<Np3OpenResponse[]>([]);

	// Undo coalescing state (P2-1a)
	#inUndoGroup = false;
	#undoGroupSnapshot: Np3OpenResponse | null = null;
	#debounceKey: string | null = null;
	#debounceTimer: ReturnType<typeof setTimeout> | null = null;

	groupedParameters = $derived.by(() => {
		if (!this.metadata?.parameterDefinitions) return [];

		const groups: Record<string, ParameterDefinition[]> = {};
		for (const p of this.metadata.parameterDefinitions) {
			if (!groups[p.group]) groups[p.group] = [];
			groups[p.group].push(p);
		}

		return Object.keys(groups)
			.sort((a, b) => {
				const idxA = GROUP_ORDER.indexOf(a);
				const idxB = GROUP_ORDER.indexOf(b);
				if (idxA === -1 && idxB === -1) return a.localeCompare(b);
				if (idxA === -1) return 1;
				if (idxB === -1) return -1;
				return idxA - idxB;
			})
			.map((name) => ({
				name,
				parameters: groups[name],
			}));
	});

	// Derived properties for UI state
	canUndo = $derived(this.undoStack.length > 0);
	canRedo = $derived(this.redoStack.length > 0);

	/**
	 * Begin an undo group — captures a snapshot before the group starts.
	 * All patchParameter calls within the group skip pushToUndo.
	 * Call commitUndoGroup() when the group is complete (P2-1a).
	 */
	beginUndoGroup() {
		if (!this.metadata || this.#inUndoGroup) return;
		this.#inUndoGroup = true;
		this.#undoGroupSnapshot = deepClone(this.metadata);
	}

	/**
	 * Commit an undo group — pushes the pre-group snapshot to the undo stack.
	 */
	commitUndoGroup() {
		if (!this.#inUndoGroup || !this.#undoGroupSnapshot) return;
		this.#inUndoGroup = false;

		this.undoStack.push(this.#undoGroupSnapshot);
		if (this.undoStack.length > MAX_HISTORY_SIZE) {
			this.undoStack.shift();
		}
		this.redoStack = [];
		this.#undoGroupSnapshot = null;
	}

	/**
	 * Debounced undo push for rapid changes to the same parameter key (P2-1a/F18).
	 * Pushes to undo immediately on first change, then suppresses subsequent
	 * same-key changes within 300ms — coalescing rapid edits into a single undo entry.
	 */
	#debouncedPushToUndo(paramKey: string) {
		if (this.#inUndoGroup) return;

		if (this.#debounceKey === paramKey && this.#debounceTimer) {
			// Same key within debounce window — extend timer, skip push
			clearTimeout(this.#debounceTimer);
		} else {
			// Different key or no active debounce — push immediately
			this.#pushToUndo();
			this.#debounceKey = paramKey;
		}

		// Reset the debounce window
		this.#debounceTimer = setTimeout(() => {
			this.#debounceTimer = null;
			this.#debounceKey = null;
		}, 300);
	}

	loadSuccess(response: Np3OpenResponse, filename?: string) {
		this.metadata = deepClone(response);
		this.isLoaded = true;
		this.isCorrupted = false;
		this.currentError = null;
		if (filename) {
			this.currentFilename = filename;
		}

		// Reset history on new file load
		this.undoStack = [];
		this.redoStack = [];
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

	/**
	 * Save current state to undo stack before making a change
	 */
	#pushToUndo() {
		if (!this.metadata) return;

		this.undoStack.push(deepClone(this.metadata));
		if (this.undoStack.length > MAX_HISTORY_SIZE) {
			this.undoStack.shift();
		}
		// Any new change clears the redo stack
		this.redoStack = [];
	}

	patch(updater: (meta: Np3OpenResponse) => Np3OpenResponse) {
		if (!this.metadata) return;

		// Clear transient errors on new user action
		if (this.currentError?.code?.startsWith("PASTE_")) {
			this.currentError = null;
		}

		this.#pushToUndo();
		this.metadata = updater(this.metadata);
	}

	/**
	 * Patches a specific parameter in the recipe, supporting nested paths.
	 * Uses debounced undo push for rapid edits to the same parameter (P2-1a).
	 */
	patchParameter(path: string, value: any) {
		if (!this.metadata) return;

		this.#debouncedPushToUndo(path);
		this.metadata.recipe = setNested(this.metadata.recipe, path, value);
	}

	rollback() {
		// Guard rollback during undo groups (F5)
		if (this.#inUndoGroup) {
			this.commitUndoGroup();
		}
		// Clear any active debounce
		if (this.#debounceTimer) {
			clearTimeout(this.#debounceTimer);
			this.#debounceTimer = null;
			this.#debounceKey = null;
		}
		if (this.undoStack.length > 0) {
			this.metadata = this.undoStack.pop()!;
		}
	}

	undo() {
		if (!this.metadata || this.undoStack.length === 0) return;

		const current = deepClone(this.metadata);
		this.redoStack.push(current);

		const previous = this.undoStack.pop()!;
		this.metadata = previous;
	}

	redo() {
		if (!this.metadata || this.redoStack.length === 0) return;

		const current = deepClone(this.metadata);
		this.undoStack.push(current);

		const next = this.redoStack.pop()!;
		this.metadata = next;
	}
	saveAsSuccessful(filePath: string) {
		this.currentFilename = filePath.split(/[/\\]/).pop() || filePath;
		this.currentError = null;
	}

	/**
	 * Mark the current recipe state as saved (F39).
	 * Called on loadSuccess, np3.save_success, and np3.save_as_success.
	 */
	markSaved() {
		if (this.metadata?.recipe) {
			this.savedRecipeSnapshot = JSON.stringify($state.snapshot(this.metadata.recipe));
		}
	}

	/**
	 * Compute whether the current recipe differs from the last saved state (F9).
	 */
	computeIsDirty(): boolean {
		if (!this.metadata?.recipe || !this.savedRecipeSnapshot) return false;
		return JSON.stringify($state.snapshot(this.metadata.recipe)) !== this.savedRecipeSnapshot;
	}

	async copyParameters() {
		if (!this.metadata) return;

		const payload: CopyPastePayload = {
			version: NP3_SCHEMA_VERSION,
			recipe: deepClone(this.metadata.recipe),
		};

		const text = JSON.stringify(payload, null, 2);
		// In a real VS Code environment, this would use vscode.env.clipboard.writeText
		// but since we're in a webview, we post a message to the extension
		getVsCode().postMessage({ type: "np3.copy", payload: text });
	}

	async pasteParameters() {
		getVsCode().postMessage({ type: "np3.paste_request" });
	}

	handlePaste(clipboardData: string): boolean {
		try {
			// Size limit: reject payloads larger than 100KB (P2-3a)
			if (clipboardData.length > 100_000) {
				this.currentError = {
					message: "Paste data is too large (max 100KB).",
					code: "PASTE_INVALID_DATA",
				};
				return false;
			}

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

			const allowedKeys = this.metadata?.parameterDefinitions?.map((p) => p.key) || [];
			if (!validateRecipe(data.recipe, allowedKeys)) {
				throw new Error("Invalid recipe data or unknown parameters");
			}

			this.#pushToUndo();
			this.metadata = {
				...this.metadata!,
				recipe: {
					...this.metadata?.recipe,
					...data.recipe,
				},
			};

			this.currentError = null;
			return true;
		} catch (e) {
			this.currentError = {
				message:
					e instanceof Error
						? e.message
						: "Invalid clipboard data. Please copy from the NP3 Editor.",
				code:
					e instanceof Error && e.message.includes("version mismatch")
						? "PASTE_VERSION_MISMATCH"
						: "PASTE_INVALID_DATA",
			};
			return false;
		}
	}
}

export const np3AppStore = new Np3Store();
