import { beforeEach, describe, expect, it, vi } from "vitest";
import type * as vscode from "vscode";

// Mock vscode module
vi.mock("vscode", () => ({
	Uri: {
		joinPath: (...args: unknown[]) => ({ fsPath: args.join("/") }),
	},
}));

// Mock node:fs/promises
import * as fs from "node:fs/promises";

vi.mock("node:fs/promises", () => ({
	copyFile: vi.fn(),
	access: vi.fn(),
	constants: { X_OK: 1 },
}));

// Mock node:crypto
vi.mock("node:crypto", () => ({
	randomUUID: () => "test-uuid-1234",
}));

// Mock child_process
const mockSpawn = vi.fn();
vi.mock("node:child_process", () => ({
	spawn: (...args: unknown[]) => mockSpawn(...args),
}));

// Mock readline
vi.mock("readline", () => ({
	createInterface: vi.fn(() => ({
		on: vi.fn(),
		close: vi.fn(),
	})),
}));

import { BinaryManager, type IpcMessage } from "./BinaryManager";

function createMockContext(): vscode.ExtensionContext {
	return {
		extensionPath: "/mock/extension",
		subscriptions: [],
	} as unknown as vscode.ExtensionContext;
}

function createMockOutputChannel(): vscode.OutputChannel {
	return {
		appendLine: vi.fn(),
		append: vi.fn(),
		clear: vi.fn(),
		show: vi.fn(),
		hide: vi.fn(),
		dispose: vi.fn(),
		replace: vi.fn(),
		name: "recipe",
	} as unknown as vscode.OutputChannel;
}

describe("BinaryManager", () => {
	let manager: BinaryManager;
	let mockContext: vscode.ExtensionContext;
	let mockOutputChannel: vscode.OutputChannel;

	beforeEach(() => {
		vi.clearAllMocks();
		mockContext = createMockContext();
		mockOutputChannel = createMockOutputChannel();
		manager = new BinaryManager(mockContext, mockOutputChannel);
		// By default, binary exists
		vi.mocked(fs.access).mockResolvedValue(undefined);
	});

	describe("start", () => {
		const testFilePath = "/mock/path/test.np3";
		const expectedBackupPath = "/mock/path/test.np3.bak";

		it("should spawn the Go binary", async () => {
			const lineHandler: ((line: string) => void)[] = [];
			const mockProcess = {
				on: vi.fn(),
				stdin: { write: vi.fn(), end: vi.fn() },
				stdout: { on: vi.fn() },
				stderr: { on: vi.fn() },
				killed: false,
				kill: vi.fn(),
			};
			mockSpawn.mockReturnValue(mockProcess);

			// Capture readline line handler to simulate pong response
			const mockCreateInterface = vi.fn(() => ({
				on: vi.fn((event: string, handler: (line: string) => void) => {
					if (event === "line") lineHandler.push(handler);
				}),
				close: vi.fn(),
			}));
			vi.doMock("node:readline", () => ({ createInterface: mockCreateInterface }));

			const startPromise = manager.start(testFilePath);

			// Wait a tick for the ping to be written
			await new Promise((r) => setTimeout(r, 50));

			// Simulate pong response
			if (lineHandler.length > 0) {
				lineHandler[0](
					JSON.stringify({
						type: "np3.pong",
						payload: { status: "ok", version: "1.0.0" },
						requestId: "test-uuid-1234",
					}),
				);
			}

			// Wait for start to resolve via pong
			await new Promise((r) => setTimeout(r, 100));

			expect(mockSpawn).toHaveBeenCalledTimes(1);
			const callArgs = mockSpawn.mock.calls[0];
			expect(callArgs[0]).toContain("np3tool");
		});

		it("should create a .bak file before spawning the Go binary", async () => {
			const mockProcess = {
				on: vi.fn(),
				stdin: { write: vi.fn(), end: vi.fn() },
				stdout: { on: vi.fn() },
				stderr: { on: vi.fn() },
				killed: false,
				kill: vi.fn(),
			};
			mockSpawn.mockReturnValue(mockProcess);

			void manager.start(testFilePath);

			await new Promise((r) => setTimeout(r, 200));

			expect(fs.copyFile).toHaveBeenCalledWith(testFilePath, expectedBackupPath);
			const callOrder = vi.mocked(fs.copyFile).mock.invocationCallOrder[0];
			const spawnOrder = mockSpawn.mock.invocationCallOrder[0];
			expect(callOrder).toBeLessThan(spawnOrder);
		});

		it("should gracefully handle backup creation failure without crashing", async () => {
			const mockProcess = {
				on: vi.fn(),
				stdin: { write: vi.fn(), end: vi.fn() },
				stdout: { on: vi.fn() },
				stderr: { on: vi.fn() },
				killed: false,
				kill: vi.fn(),
			};
			mockSpawn.mockReturnValue(mockProcess);

			vi.mocked(fs.copyFile).mockRejectedValueOnce(new Error("EACCES: permission denied"));

			void manager.start(testFilePath);

			await new Promise((r) => setTimeout(r, 200));

			expect(mockSpawn).toHaveBeenCalledTimes(1);
			expect(mockOutputChannel.appendLine).toHaveBeenCalledWith(
				expect.stringContaining(
					"Warning: Failed to create backup for /mock/path/test.np3: EACCES: permission denied",
				),
			);
		});

		it("should reject when spawn fails", async () => {
			mockSpawn.mockImplementation(() => {
				throw new Error("ENOENT");
			});

			await expect(manager.start(testFilePath)).rejects.toThrow("Failed to spawn Go binary");
		});

		it("should reject when binary does not exist", async () => {
			vi.mocked(fs.access).mockRejectedValueOnce(new Error("ENOENT"));

			await expect(manager.start(testFilePath)).rejects.toThrow("np3tool binary not found");
		});
	});

	describe("stop", () => {
		it("should clean up resources when stopped", () => {
			expect(() => manager.stop()).not.toThrow();
		});
	});

	describe("send", () => {
		it("should reject when binary is not running", async () => {
			const message: IpcMessage = { type: "np3.ping", payload: {} };
			await expect(manager.send(message)).rejects.toThrow("Go binary is not running");
		});
	});
});
