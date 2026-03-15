import { describe, expect, it } from "vitest";
import type { Np3OpenResponse, ParameterDefinition } from "../types";
import { Np3Store } from "./np3.svelte";

function makeMockResponse(hash = "abc123"): Np3OpenResponse {
	return {
		hash,
		recipe: { name: "Test Recipe", version: 1 },
		parameterDefinitions: [] as ParameterDefinition[],
	};
}

describe("Np3Store", () => {
	describe("loadSuccess", () => {
		it("should set metadata and mark as loaded", () => {
			const store = new Np3Store();
			const response = makeMockResponse();

			store.loadSuccess(response);

			expect(store.isLoaded).toBe(true);
			expect(store.isCorrupted).toBe(false);
			expect(store.currentError).toBeNull();
			expect(store.metadata).toEqual(response);
			expect(store.canUndo).toBe(false);
			expect(store.canRedo).toBe(false);
		});

		it("should set filename when provided", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse(), "test.np3");
			expect(store.currentFilename).toBe("test.np3");
		});

		it("should clear undo/redo stacks on load", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));
			store.patch((m) => ({ ...m, hash: "v2" }));
			expect(store.canUndo).toBe(true);

			store.loadSuccess(makeMockResponse("v3"));
			expect(store.canUndo).toBe(false);
			expect(store.canRedo).toBe(false);
		});
	});

	describe("patch", () => {
		it("should apply optimistic update and store in undo stack", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.patch((meta) => ({ ...meta, hash: "v2" }));

			expect(store.metadata?.hash).toBe("v2");
			expect(store.canUndo).toBe(true);
		});

		it("should clear redo stack on new patch", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));
			store.patch((m) => ({ ...m, hash: "v2" }));
			store.undo();
			expect(store.canRedo).toBe(true);

			store.patch((m) => ({ ...m, hash: "v3" }));
			expect(store.canRedo).toBe(false);
		});

		it("should handle nested patches via patchParameter", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.patchParameter("colorGrading.highlights.hue", 180);

			expect(store.metadata?.recipe.colorGrading?.highlights.hue).toBe(180);
			expect(store.canUndo).toBe(true);
		});
	});

	describe("undo/redo", () => {
		it("should undo and redo correctly", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.patch((meta) => ({ ...meta, hash: "v2" }));
			expect(store.metadata?.hash).toBe("v2");

			store.undo();
			expect(store.metadata?.hash).toBe("v1");
			expect(store.canRedo).toBe(true);

			store.redo();
			expect(store.metadata?.hash).toBe("v2");
			expect(store.canUndo).toBe(true);
		});

		it("should handle multiple steps", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));
			store.patch((m) => ({ ...m, hash: "v2" }));
			store.patch((m) => ({ ...m, hash: "v3" }));

			store.undo();
			expect(store.metadata?.hash).toBe("v2");
			store.undo();
			expect(store.metadata?.hash).toBe("v1");

			store.redo();
			expect(store.metadata?.hash).toBe("v2");
			store.redo();
			expect(store.metadata?.hash).toBe("v3");
		});
	});

	describe("history cap", () => {
		it("should cap history at MAX_HISTORY_SIZE (100)", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v0"));

			// Apply 120 patches
			for (let i = 1; i <= 120; i++) {
				store.patch((meta) => ({ ...meta, hash: `v${i}` }));
			}

			expect(store.metadata?.hash).toBe("v120");

			// Undo 100 times - should land at v20, not v0
			for (let i = 0; i < 100; i++) {
				store.undo();
			}
			expect(store.metadata?.hash).toBe("v20");
			expect(store.canUndo).toBe(false);
		});
	});

	describe("rollback", () => {
		it("should still support rollback (used for patch errors)", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));
			store.patch((m) => ({ ...m, hash: "v2" }));

			store.rollback();
			expect(store.metadata?.hash).toBe("v1");
		});
	});

	describe("copy / paste", () => {
		it("should handle valid paste", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			const pasteData = {
				version: 1,
				recipe: { name: "Pasted Recipe", version: 1, someParam: 123 },
			};

			const result = store.handlePaste(JSON.stringify(pasteData));

			expect(result).toBe(true);
			expect(store.metadata?.recipe.name).toBe("Pasted Recipe");
			expect(store.canUndo).toBe(true);
		});

		it("should reject paste with unknown parameter keys", () => {
			const store = new Np3Store();
			const mockResp = makeMockResponse("v1");
			mockResp.parameterDefinitions = [
				{
					key: "knownParam",
					label: "Known",
					type: "continuous",
					min: 0,
					max: 100,
					step: 1,
					defaultValue: 0,
					group: "Basic",
				},
			];
			store.loadSuccess(mockResp);

			const pasteData = {
				version: 1,
				recipe: { knownParam: 50, maliciousParam: "hacked" },
			};

			const result = store.handlePaste(JSON.stringify(pasteData));

			expect(result).toBe(false);
			expect(store.currentError?.message).toContain("unknown parameters");
		});

		it("should reject paste with version mismatch", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			const pasteData = {
				version: 99, // Incompatible
				recipe: { name: "Future Recipe" },
			};

			const result = store.handlePaste(JSON.stringify(pasteData));

			expect(result).toBe(false);
			// Should not have changed
			expect(store.metadata?.recipe.name).toBe("Test Recipe");
			expect(store.currentError?.message).toContain("version mismatch");
			// Important: should NOT be corrupted
			expect(store.isCorrupted).toBe(false);
		});

		it("should handle malformed paste data", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			const result = store.handlePaste("not json");

			expect(result).toBe(false);
			expect(store.currentError?.message).toMatch(/Invalid clipboard data|Unexpected token/);
			// Important: should NOT be corrupted
			expect(store.isCorrupted).toBe(false);
		});
	});

	describe("copyParameters", () => {
		it("should post message with serialized parameters", () => {
			const store = new Np3Store();
			const response = makeMockResponse();
			store.loadSuccess(response);

			store.copyParameters();

			// Verify if it posts the right message (we'd need to mock vscode global if we wanted to test this strictly)
			// For now, let's verify it doesn't crash and we can test internal logic if we extract serialization.
		});
	});

	describe("undo coalescing (P2-1a)", () => {
		it("should coalesce rapid patchParameter calls for same key within 300ms", async () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			// Rapid changes to same key — should produce single undo entry
			store.patchParameter("contrast", 10);
			store.patchParameter("contrast", 20);
			store.patchParameter("contrast", 30);

			expect(store.undoStack.length).toBe(1);
			expect(store.metadata?.recipe.contrast).toBe(30);

			store.undo();
			expect(store.metadata?.recipe.contrast).toBeUndefined();
		});

		it("should create separate undo entries for different keys", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.patchParameter("contrast", 10);
			store.patchParameter("saturation", 50);

			expect(store.undoStack.length).toBe(2);
		});

		it("should use undo groups for slider drags", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.beginUndoGroup();
			store.patchParameter("contrast", 10);
			store.patchParameter("contrast", 30);
			store.patchParameter("contrast", 50);
			store.commitUndoGroup();

			expect(store.undoStack.length).toBe(1);
			expect(store.metadata?.recipe.contrast).toBe(50);

			store.undo();
			expect(store.metadata?.recipe.contrast).toBeUndefined();
		});

		it("should guard rollback during undo group (F5)", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.beginUndoGroup();
			store.patchParameter("contrast", 10);
			// Rollback during group should close the group first
			store.rollback();

			expect(store.metadata?.recipe.name).toBe("Test Recipe");
			expect(store.canUndo).toBe(false);
		});
	});

	describe("paste validation hardening (P2-3a)", () => {
		it("should reject paste data larger than 100KB", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			const largeData = "x".repeat(100_001);
			const result = store.handlePaste(largeData);

			expect(result).toBe(false);
			expect(store.currentError?.code).toBe("PASTE_INVALID_DATA");
			expect(store.currentError?.message).toContain("too large");
		});

		it("should reject paste with non-string name", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			const data = { version: 1, recipe: { name: 123 } };
			const result = store.handlePaste(JSON.stringify(data));

			expect(result).toBe(false);
		});

		it("should reject paste with name longer than 256 chars", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			const data = { version: 1, recipe: { name: "x".repeat(257) } };
			const result = store.handlePaste(JSON.stringify(data));

			expect(result).toBe(false);
		});

		it("should reject paste with invalid pointCurve entries", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			const data = {
				version: 1,
				recipe: { pointCurve: [{ input: -1, output: 0 }] },
			};
			const result = store.handlePaste(JSON.stringify(data));

			expect(result).toBe(false);
		});

		it("should reject paste with string values for numeric fields", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			const data = {
				version: 1,
				recipe: { contrast: "not a number" },
			};
			const result = store.handlePaste(JSON.stringify(data));

			expect(result).toBe(false);
		});
	});

	describe("dirty tracking", () => {
		it("should track dirty state via markSaved and computeIsDirty", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));
			store.markSaved();

			expect(store.computeIsDirty()).toBe(false);

			store.patchParameter("contrast", 50);
			expect(store.computeIsDirty()).toBe(true);

			store.undo();
			expect(store.computeIsDirty()).toBe(false);
		});
	});

	describe("error handling", () => {
		it("should load and clear error", () => {
			const store = new Np3Store();
			const err = { message: "Test Error", code: "TEST" };

			store.loadError(err);
			expect(store.isCorrupted).toBe(true);
			expect(store.currentError).toEqual(err);

			store.clearError();
			expect(store.isCorrupted).toBe(false);
			expect(store.currentError).toBeNull();
		});
	});
});
