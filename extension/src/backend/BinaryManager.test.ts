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
	});

	describe("start", () => {
		const testFilePath = "/mock/path/test.np3";
		const expectedBackupPath = "/mock/path/test.np3.bak";

		it("should spawn the Go binary", async () => {
			const mockProcess = {
				on: vi.fn(),
				stdin: { write: vi.fn(), end: vi.fn() },
				stdout: { on: vi.fn() },
				stderr: { on: vi.fn() },
				killed: false,
				kill: vi.fn(),
			};
			mockSpawn.mockReturnValue(mockProcess);

			// Simulate the process spawning without immediate error
			mockProcess.on.mockImplementation(
				(_event: string, _handler: (...args: unknown[]) => void) => {
					// Don't call error handler - successful spawn
				},
			);

			void manager.start(testFilePath);

			// Wait for the setTimeout in start()
			await new Promise((r) => setTimeout(r, 200));

			// The binary path should include platform-specific binary name
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
			// Check if copyFile was called before spawn
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

			// Should not throw
			await expect(manager.start(testFilePath)).resolves.toBeUndefined();

			// Wait for setTimeout
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
	});

	describe("stop", () => {
		it("should clean up resources when stopped", () => {
			// Manager should not throw even when no process is running
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
