import { beforeEach, describe, expect, it, vi } from "vitest";
import type * as vscode from "vscode";

// Mock vscode module
vi.mock("vscode", () => ({
	Uri: {
		joinPath: (...args: unknown[]) => ({ fsPath: args.join("/") }),
	},
}));

// Mock child_process
const mockSpawn = vi.fn();
vi.mock("child_process", () => ({
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

			void manager.start();

			// Wait for the setTimeout in start()
			await new Promise((r) => setTimeout(r, 200));

			// The binary path should include platform-specific binary name
			expect(mockSpawn).toHaveBeenCalledTimes(1);
			const callArgs = mockSpawn.mock.calls[0];
			expect(callArgs[0]).toContain("np3tool");
		});

		it("should reject when spawn fails", async () => {
			mockSpawn.mockImplementation(() => {
				throw new Error("ENOENT");
			});

			await expect(manager.start()).rejects.toThrow("Failed to spawn Go binary");
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
