import { describe, expect, it } from "vitest";
import { Np3Store } from "./np3.svelte";

function makeMockResponse(hash = "abc123"): { hash: string; recipe: Record<string, unknown> } {
	return {
		hash,
		recipe: { name: "Test Recipe", version: 1 },
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
		});

		it("should set filename when provided", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse(), "test.np3");
			expect(store.currentFilename).toBe("test.np3");
		});

		it("should deep-copy response to prevent mutation leaks", () => {
			const store = new Np3Store();
			const response = makeMockResponse();

			store.loadSuccess(response);
			response.hash = "mutated";

			expect(store.metadata?.hash).toBe("abc123");
		});

		it("should clear corrupted state if previously set", () => {
			const store = new Np3Store();
			store.loadError({ message: "bad", code: "ERR_CORRUPTED_FILE" });
			expect(store.isCorrupted).toBe(true);

			store.loadSuccess(makeMockResponse());
			expect(store.isCorrupted).toBe(false);
			expect(store.currentError).toBeNull();
		});
	});

	describe("loadError", () => {
		it("should set corrupted state", () => {
			const store = new Np3Store();
			const error = { message: "checksum mismatch", code: "ERR_INVALID_CHECKSUM" };

			store.loadError(error);

			expect(store.isCorrupted).toBe(true);
			expect(store.currentError).toEqual(error);
			expect(store.metadata).toBeNull();
		});

		it("should set filename when provided", () => {
			const store = new Np3Store();
			store.loadError({ message: "bad", code: "ERR" }, "corrupt.np3");
			expect(store.currentFilename).toBe("corrupt.np3");
		});
	});

	describe("clearError", () => {
		it("should reset corruption state", () => {
			const store = new Np3Store();
			store.loadError({ message: "bad", code: "ERR" });

			store.clearError();

			expect(store.isCorrupted).toBe(false);
			expect(store.currentError).toBeNull();
		});
	});

	describe("patch", () => {
		it("should apply optimistic update and store history", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.patch((meta) => ({ ...meta, hash: "v2" }));

			expect(store.metadata?.hash).toBe("v2");
		});

		it("should do nothing if no metadata loaded", () => {
			const store = new Np3Store();
			// Should not throw
			store.patch((meta) => ({ ...meta, hash: "crash?" }));
			expect(store.metadata).toBeNull();
		});
	});

	describe("rollback", () => {
		it("should revert to previous state after patch", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.patch((meta) => ({ ...meta, hash: "v2" }));
			expect(store.metadata?.hash).toBe("v2");

			store.rollback();
			expect(store.metadata?.hash).toBe("v1");
		});

		it("should handle rollback with no history gracefully", () => {
			const store = new Np3Store();
			// Should not throw
			store.rollback();
			expect(store.metadata).toBeNull();
		});

		it("should support multiple rollbacks in sequence", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v1"));

			store.patch((meta) => ({ ...meta, hash: "v2" }));
			store.patch((meta) => ({ ...meta, hash: "v3" }));

			store.rollback();
			expect(store.metadata?.hash).toBe("v2");

			store.rollback();
			expect(store.metadata?.hash).toBe("v1");
		});
	});

	describe("history cap", () => {
		it("should cap history at MAX_HISTORY_SIZE (20)", () => {
			const store = new Np3Store();
			store.loadSuccess(makeMockResponse("v0"));

			// Apply 25 patches — history should only keep last 20
			for (let i = 1; i <= 25; i++) {
				store.patch((meta) => ({ ...meta, hash: `v${i}` }));
			}

			expect(store.metadata?.hash).toBe("v25");

			// Rollback all the way — should be capped at 20 successful rollbacks
			let rollbackCount = 0;
			for (let i = 0; i < 30; i++) {
				const before = store.metadata?.hash;
				store.rollback();
				const after = store.metadata?.hash;
				if (before === after) break; // no more history
				rollbackCount++;
			}

			expect(rollbackCount).toBeLessThanOrEqual(20);
			expect(rollbackCount).toBeGreaterThan(0);
		});
	});
});
