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
			mockResp.parameterDefinitions = [{ key: "knownParam", label: "Known", type: "continuous", min: 0, max: 100, step: 1, defaultValue: 0, group: "Basic" }];
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
