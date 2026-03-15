import { type ChildProcess, spawn } from "node:child_process";
import * as crypto from "node:crypto";
import * as fs from "node:fs/promises";
import * as path from "node:path";
import { createInterface, type Interface } from "node:readline";
import type * as vscode from "vscode";

export interface IpcMessage {
	type: string;
	payload: Record<string, unknown>;
	requestId?: string;
}

export class BinaryManager {
	private process: ChildProcess | null = null;
	private readline: Interface | null = null;
	private pendingRequests: Map<
		string,
		{
			resolve: (value: IpcMessage) => void;
			reject: (reason: Error) => void;
			timeout: NodeJS.Timeout;
		}
	> = new Map();

	// Crash recovery state (P2-8a)
	private restartCount = 0;
	private maxRestarts = 3;
	private lastFilePath: string | null = null;
	private intentionalStop = false;
	onCrash?: (restarted: boolean) => void;

	constructor(
		private readonly context: vscode.ExtensionContext,
		private readonly outputChannel: vscode.OutputChannel,
	) {}

	private getBinaryPath(): string {
		const platform = process.platform;
		const binaryName = platform === "win32" ? "np3tool.exe" : "np3tool";
		return path.join(this.context.extensionPath, "bin", binaryName);
	}

	private async createBackup(filePath: string): Promise<void> {
		const parsedPath = path.parse(filePath);
		const backupPath = path.join(parsedPath.dir, `${parsedPath.name}${parsedPath.ext}.bak`);
		await fs.copyFile(filePath, backupPath);
	}

	async start(filePath: string): Promise<void> {
		this.lastFilePath = filePath;
		this.intentionalStop = false;
		const binaryPath = this.getBinaryPath();

		// Verify binary exists before spawn (P1-6a)
		try {
			await fs.access(binaryPath, fs.constants.X_OK);
		} catch {
			throw new Error(`np3tool binary not found at ${binaryPath}. Please reinstall the extension.`);
		}

		// Create backup before spawning
		try {
			await this.createBackup(filePath);
		} catch (err) {
			this.outputChannel.appendLine(
				`Warning: Failed to create backup for ${filePath}: ${err instanceof Error ? err.message : String(err)}`,
			);
		}

		return new Promise<void>((resolve, reject) => {
			try {
				this.process = spawn(binaryPath, [], {
					stdio: ["pipe", "pipe", "pipe"],
				});
			} catch (err) {
				reject(
					new Error(
						`Failed to spawn Go binary at ${binaryPath}: ${err instanceof Error ? err.message : String(err)}`,
					),
				);
				return;
			}

			this.process.on("error", (err) => {
				this.outputChannel.appendLine(`Go binary error: ${err.message}`);
				reject(new Error(`Failed to start Go binary: ${err.message}`));
			});

			this.process.on("exit", (code, signal) => {
				this.outputChannel.appendLine(`Go binary exited (code=${code}, signal=${signal})`);
				this.cleanup();

				// Crash detection and auto-restart (P2-8a)
				if (!this.intentionalStop && (code !== 0 || signal)) {
					this.outputChannel.appendLine(
						`Unexpected exit detected (restart ${this.restartCount}/${this.maxRestarts})`,
					);
					if (this.restartCount < this.maxRestarts && this.lastFilePath) {
						this.restartCount++;
						this.start(this.lastFilePath)
							.then(() => {
								this.outputChannel.appendLine("Auto-restart successful");
								this.onCrash?.(true);
							})
							.catch((err) => {
								this.outputChannel.appendLine(
									`Auto-restart failed: ${err instanceof Error ? err.message : String(err)}`,
								);
								this.onCrash?.(false);
							});
					} else {
						this.onCrash?.(false);
					}
				}
			});

			// Pipe stderr to output channel (Go debug logging)
			if (this.process.stderr) {
				this.process.stderr.on("data", (data: Buffer) => {
					this.outputChannel.appendLine(`[np3tool] ${data.toString().trim()}`);
				});
			}

			// Read JSONL responses from stdout
			if (this.process.stdout) {
				this.readline = createInterface({
					input: this.process.stdout,
					crlfDelay: Infinity,
				});

				this.readline.on("line", (line: string) => {
					try {
						const message: IpcMessage = JSON.parse(line);
						this.handleResponse(message);
					} catch (_err) {
						this.outputChannel.appendLine(`Failed to parse Go response: ${line}`);
					}
				});
			}

			// Send ping and wait for pong as readiness check (P1-6b)
			const pingId = crypto.randomUUID();
			const pingTimeout = setTimeout(() => {
				this.pendingRequests.delete(pingId);
				reject(new Error("np3tool binary did not respond to ping within 5 seconds"));
			}, 5000);

			this.pendingRequests.set(pingId, {
				resolve: (response) => {
					clearTimeout(pingTimeout);
					this.restartCount = 0; // Reset on successful start (F8)
					const version = (response.payload as { version?: string })?.version;
					if (version) {
						this.outputChannel.appendLine(`np3tool version: ${version}`);
					}
					resolve();
				},
				reject: (err) => {
					clearTimeout(pingTimeout);
					reject(err);
				},
				timeout: pingTimeout,
			});

			const pingMessage = JSON.stringify({ type: "np3.ping", payload: {}, requestId: pingId });
			this.process.stdin?.write(`${pingMessage}\n`);
		});
	}

	async send(message: IpcMessage): Promise<IpcMessage> {
		if (!this.process || !this.process.stdin) {
			throw new Error("Go binary is not running");
		}

		return new Promise<IpcMessage>((resolve, reject) => {
			const requestId = crypto.randomUUID();

			const timeout = setTimeout(() => {
				this.pendingRequests.delete(requestId);
				reject(new Error(`IPC request timed out: ${message.type}`));
			}, 10000);

			this.pendingRequests.set(requestId, { resolve, reject, timeout });

			const outgoing = { ...message, requestId };
			const jsonLine = `${JSON.stringify(outgoing)}\n`;
			this.process?.stdin?.write(jsonLine, (err) => {
				if (err) {
					clearTimeout(timeout);
					this.pendingRequests.delete(requestId);
					reject(new Error(`Failed to write to Go binary: ${err.message}`));
				}
			});
		});
	}

	async open(filePath: string): Promise<Record<string, unknown>> {
		const response = await this.send({
			type: "np3.open",
			payload: { filePath },
		});
		return response.payload;
	}

	private handleResponse(message: IpcMessage): void {
		// Match by requestId
		if (message.requestId) {
			const pending = this.pendingRequests.get(message.requestId);
			if (pending) {
				clearTimeout(pending.timeout);
				this.pendingRequests.delete(message.requestId);
				pending.resolve(message);
				return;
			}
		}

		// Broadcast unsolicited messages (e.g., messages without requestId)
		this.outputChannel.appendLine(`Unsolicited Go message: ${JSON.stringify(message)}`);
	}

	stop(): void {
		this.intentionalStop = true;
		if (this.process && !this.process.killed) {
			this.process.stdin?.end();
			this.process.kill();
		}
		this.cleanup();
	}

	private cleanup(): void {
		if (this.readline) {
			this.readline.close();
			this.readline = null;
		}

		// Reject all pending requests
		for (const [_id, pending] of this.pendingRequests) {
			clearTimeout(pending.timeout);
			pending.reject(new Error("Go binary terminated"));
		}
		this.pendingRequests.clear();

		this.process = null;
	}
}
