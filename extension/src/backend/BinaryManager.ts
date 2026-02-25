import * as vscode from "vscode";
import * as path from "path";
import { ChildProcess, spawn } from "child_process";
import { createInterface, Interface } from "readline";

export interface IpcMessage {
	type: string;
	payload: Record<string, unknown>;
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
	private requestCounter = 0;

	constructor(
		private readonly context: vscode.ExtensionContext,
		private readonly outputChannel: vscode.OutputChannel,
	) {}

	private getBinaryPath(): string {
		const platform = process.platform;
		const binaryName = platform === "win32" ? "np3tool.exe" : "np3tool";
		return path.join(this.context.extensionPath, "bin", binaryName);
	}

	async start(): Promise<void> {
		const binaryPath = this.getBinaryPath();

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
					} catch (err) {
						this.outputChannel.appendLine(`Failed to parse Go response: ${line}`);
					}
				});
			}

			// Consider started once the process is spawned without immediate error
			// Use a short delay to catch immediate spawn failures
			setTimeout(() => {
				if (this.process && !this.process.killed) {
					resolve();
				}
			}, 100);
		});
	}

	async send(message: IpcMessage): Promise<IpcMessage> {
		if (!this.process || !this.process.stdin) {
			throw new Error("Go binary is not running");
		}

		return new Promise<IpcMessage>((resolve, reject) => {
			const id = `req_${++this.requestCounter}`;

			const timeout = setTimeout(() => {
				this.pendingRequests.delete(message.type);
				reject(new Error(`IPC request timed out: ${message.type}`));
			}, 10000);

			this.pendingRequests.set(message.type, { resolve, reject, timeout });

			const jsonLine = JSON.stringify(message) + "\n";
			this.process!.stdin!.write(jsonLine, (err) => {
				if (err) {
					clearTimeout(timeout);
					this.pendingRequests.delete(message.type);
					reject(new Error(`Failed to write to Go binary: ${err.message}`));
				}
			});
		});
	}

	private handleResponse(message: IpcMessage): void {
		// Map response types to request types (e.g., np3.pong → np3.ping)
		const requestType = this.getRequestTypeForResponse(message.type);
		const pending = this.pendingRequests.get(requestType);

		if (pending) {
			clearTimeout(pending.timeout);
			this.pendingRequests.delete(requestType);
			pending.resolve(message);
		} else {
			// Broadcast unsolicited messages (e.g., errors)
			this.outputChannel.appendLine(`Unsolicited Go message: ${JSON.stringify(message)}`);
		}
	}

	private getRequestTypeForResponse(responseType: string): string {
		// Convention: response is the "answer" to a request
		// np3.pong → np3.ping, error can match any pending request
		const mappings: Record<string, string> = {
			"np3.pong": "np3.ping",
		};
		return mappings[responseType] || responseType;
	}

	stop(): void {
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
		for (const [type, pending] of this.pendingRequests) {
			clearTimeout(pending.timeout);
			pending.reject(new Error("Go binary terminated"));
		}
		this.pendingRequests.clear();

		this.process = null;
	}
}
