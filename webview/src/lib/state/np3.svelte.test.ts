import { describe, expect, it } from "vitest";
import { Np3Store } from "./np3.svelte";

function makeMockResponse(hash = "abc123"): { hash: string; recipe: Record<string, unknown>; parameters: never[] } {
	return {
		hash,
		recipe: { name: "Test Recipe", version: 1 },
		parameters: [],
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
			store.patch(m => ({ ...m, hash: "v2" }));
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
			store.patch(m => ({ ...m, hash: "v2" }));
			store.undo();
			expect(store.canRedo).toBe(true);

			store.patch(m => ({ ...m, hash: "v3" }));
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
			store.patch(m => ({ ...m, hash: "v2" }));
			store.patch(m => ({ ...m, hash: "v3" }));

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
			store.patch(m => ({ ...m, hash: "v2" }));

			store.rollback();
			expect(store.metadata?.hash).toBe("v1");
		});
	});
});
