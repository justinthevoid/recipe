import * as vscode from "vscode";
import { BinaryManager, type IpcMessage } from "../backend/BinaryManager";

interface PanelEntry {
	webviewPanel: vscode.WebviewPanel;
	binaryManager: BinaryManager;
	document: vscode.CustomDocument;
}

export class Np3EditorPanel implements vscode.CustomReadonlyEditorProvider {
	private readonly extensionUri: vscode.Uri;
	private readonly outputChannel: vscode.OutputChannel;
	private readonly statusBarItem: vscode.StatusBarItem | null;

	// Multi-panel support (P1-4a)
	private panels = new Map<string, PanelEntry>();
	private activePanelUri: string | null = null;

	// Dirty document tracking (P1-3a)
	private dirtyDocuments = new Set<string>();

	constructor(
		private readonly context: vscode.ExtensionContext,
		outputChannel: vscode.OutputChannel,
		statusBarItem?: vscode.StatusBarItem,
	) {
		this.extensionUri = context.extensionUri;
		this.outputChannel = outputChannel;
		this.statusBarItem = statusBarItem ?? null;
	}

	openCustomDocument(
		uri: vscode.Uri,
		_openContext: vscode.CustomDocumentOpenContext,
		_token: vscode.CancellationToken,
	): vscode.CustomDocument {
		return { uri, dispose: () => {} };
	}

	async resolveCustomEditor(
		document: vscode.CustomDocument,
		webviewPanel: vscode.WebviewPanel,
		_token: vscode.CancellationToken,
	): Promise<void> {
		const documentUri = document.uri.toString();

		// Enforce max concurrent panels limit (P3-2b)
		const maxPanels = vscode.workspace.getConfiguration("recipe").get<number>("maxConcurrentPanels", 5);
		if (this.panels.size >= maxPanels) {
			vscode.window.showInformationMessage(
				`Maximum of ${maxPanels} NP3 files can be open simultaneously. Close an existing file first.`,
			);
			webviewPanel.dispose();
			return;
		}

		webviewPanel.webview.options = {
			enableScripts: true,
			localResourceRoots: [vscode.Uri.joinPath(this.extensionUri, "dist", "webview")],
		};

		webviewPanel.webview.html = this.getHtmlForWebview(webviewPanel.webview);

		// Spawn Go binary
		const binaryManager = new BinaryManager(this.context, this.outputChannel);

		// Wire crash recovery notifications (P2-8a)
		binaryManager.onCrash = (restarted: boolean) => {
			if (restarted) {
				this.outputChannel.appendLine("Binary restarted after crash — reopening file");
				this.updateStatusBar("connected");
				// Re-open the file in the new binary and notify webview to resync (F2)
				webviewPanel.webview.postMessage({
					type: "extension.open",
					payload: { filePath: document.uri.fsPath },
				});
				if (this.dirtyDocuments.has(documentUri)) {
					vscode.window.showWarningMessage(
						"np3tool binary crashed and was restarted. You had unsaved changes — please verify your edits.",
					);
				} else {
					vscode.window.showWarningMessage(
						"np3tool binary crashed and was restarted. The file has been reloaded.",
					);
				}
			} else {
				this.updateStatusBar("error");
				vscode.window.showErrorMessage(
					"np3tool binary crashed and could not be restarted. Please close and reopen the file.",
				);
			}
		};

		// Track panel in map
		this.panels.set(documentUri, { webviewPanel, binaryManager, document });

		// Track active panel
		const updateActivePanel = () => {
			if (webviewPanel.active) {
				this.activePanelUri = documentUri;
			} else if (this.activePanelUri === documentUri) {
				this.activePanelUri = null;
			}
		};
		updateActivePanel();
		webviewPanel.onDidChangeViewState(updateActivePanel);

		try {
			await binaryManager.start(document.uri.fsPath);
			this.outputChannel.appendLine("Go binary started successfully");
			this.updateStatusBar("connected");
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : "Unknown error starting Go binary";
			this.outputChannel.appendLine(`Failed to start Go binary: ${errorMessage}`);
			this.updateStatusBar("error");

			// Offer retry on spawn failure (P2-8b)
			const action = await vscode.window.showErrorMessage(
				`Failed to start np3tool: ${errorMessage}`,
				"Retry",
			);
			if (action === "Retry") {
				try {
					await binaryManager.start(document.uri.fsPath);
					this.outputChannel.appendLine("Go binary started on retry");
					this.updateStatusBar("connected");
				} catch (retryErr) {
					const retryMsg = retryErr instanceof Error ? retryErr.message : "Unknown error";
					this.outputChannel.appendLine(`Retry failed: ${retryMsg}`);
					webviewPanel.webview.postMessage({
						type: "error",
						payload: { message: retryMsg, code: "BINARY_SPAWN_FAILED" },
					});
					return;
				}
			} else {
				webviewPanel.webview.postMessage({
					type: "error",
					payload: { message: errorMessage, code: "BINARY_SPAWN_FAILED" },
				});
				return;
			}
		}

		// Handle webview messages
		webviewPanel.webview.onDidReceiveMessage(
			async (message: IpcMessage) => {
				const reqId = (message as { requestId?: string }).requestId;
				this.outputChannel.appendLine(`[${reqId ?? "no-id"}] Webview → Extension: ${message.type}`);

				// Intercept webview.ready to trigger file load
				if (message.type === "webview.ready") {
					webviewPanel.webview.postMessage({
						type: "extension.open",
						payload: { filePath: document.uri.fsPath },
					});
					return;
				}

				// Handle save_as from webview (if triggered via UI button)
				if (message.type === "np3.save_as") {
					this.executeSaveAs(webviewPanel, binaryManager, document);
					return;
				}

				// Handle clipboard copy
				if (message.type === "np3.copy") {
					const payload = message.payload;
					if (typeof payload === "string") {
						await vscode.env.clipboard.writeText(payload);
					} else {
						this.outputChannel.appendLine("Error: np3.copy payload is not a string");
					}
					return;
				}

				// Handle clipboard paste request
				if (message.type === "np3.paste_request") {
					const text = await vscode.env.clipboard.readText();
					webviewPanel.webview.postMessage({
						type: "np3.paste_response",
						payload: text,
					});
					return;
				}

				// Handle open backup request (P1-5a)
				if (message.type === "np3.open_backup") {
					const bakPath = `${document.uri.fsPath}.bak`;
					try {
						const bakUri = vscode.Uri.file(bakPath);
						await vscode.commands.executeCommand("vscode.open", bakUri);
					} catch {
						webviewPanel.webview.postMessage({
							type: "error",
							payload: {
								message: "Backup file not found or could not be opened.",
								code: "IO_ERROR",
							},
						});
					}
					return;
				}

				// Handle reveal in finder request (P1-5a)
				if (message.type === "np3.reveal_in_finder") {
					await vscode.commands.executeCommand("revealFileInOS", document.uri);
					return;
				}

				try {
					const response = await binaryManager.send(message);
					const respId = (response as { requestId?: string }).requestId;
					this.outputChannel.appendLine(`[${respId ?? "no-id"}] Go → Webview: ${response.type}`);

					// Intercept responses for dirty tracking (P1-3a)
					this.updateDirtyState(documentUri, response);

					webviewPanel.webview.postMessage(response);
				} catch (err) {
					const errorMessage = err instanceof Error ? err.message : "IPC communication error";
					this.outputChannel.appendLine(`[${reqId ?? "no-id"}] IPC Error: ${errorMessage}`);
					webviewPanel.webview.postMessage({
						type: "error",
						payload: { message: errorMessage, code: "IPC_ERROR" },
					});
				}
			},
			undefined,
			this.context.subscriptions,
		);

		// Check if backup exists and notify webview (P1-5a)
		try {
			const bakPath = `${document.uri.fsPath}.bak`;
			await vscode.workspace.fs.stat(vscode.Uri.file(bakPath));
			// backup exists — will be communicated via error payload if needed
		} catch {
			// no backup — that's fine
		}

		// Clean up binary on panel close
		webviewPanel.onDidDispose(() => {
			// Warn if dirty (P1-3d)
			if (this.dirtyDocuments.has(documentUri)) {
				const filename = document.uri.fsPath.split(/[\\/]/).pop() || document.uri.fsPath;
				vscode.window.showWarningMessage(
					`You have unsaved changes to ${filename}. They will be lost.`,
				);
				this.dirtyDocuments.delete(documentUri);
			}

			this.panels.delete(documentUri);
			if (this.activePanelUri === documentUri) {
				this.activePanelUri = null;
			}
			binaryManager.stop();
			this.updateStatusBar(this.panels.size > 0 ? "connected" : "disconnected");
		});
	}

	/**
	 * Update dirty state based on Go response type.
	 */
	private updateDirtyState(documentUri: string, response: IpcMessage): void {
		const payload = response.payload as { dirty?: boolean } | undefined;

		switch (response.type) {
			case "np3.patch_success":
			case "np3.sync":
				if (payload?.dirty) {
					this.dirtyDocuments.add(documentUri);
				}
				break;
			case "np3.save_success":
			case "np3.metadata": // from reset or open
				this.dirtyDocuments.delete(documentUri);
				break;
			case "np3.save_as_success":
				this.dirtyDocuments.delete(documentUri);
				break;
		}
	}

	/**
	 * Public method to trigger Save As from command palette.
	 */
	public triggerSaveAs() {
		if (this.activePanelUri) {
			const entry = this.panels.get(this.activePanelUri);
			if (entry) {
				this.executeSaveAs(entry.webviewPanel, entry.binaryManager, entry.document);
			}
		}
	}

	/**
	 * Public method to trigger Reset All from command palette.
	 */
	public triggerResetAll() {
		if (this.activePanelUri) {
			const entry = this.panels.get(this.activePanelUri);
			if (entry) {
				entry.webviewPanel.webview.postMessage({
					type: "extension.triggerReset",
					payload: {},
				});
			}
		}
	}

	/**
	 * Stop all binary processes (called from deactivate).
	 */
	public stopAll(): void {
		for (const [_uri, entry] of this.panels) {
			entry.binaryManager.stop();
		}
		this.panels.clear();
		this.dirtyDocuments.clear();
		this.updateStatusBar("disconnected");
	}

	/**
	 * Update status bar item based on binary connection state (P2-4c).
	 */
	private updateStatusBar(state: "connected" | "error" | "disconnected"): void {
		if (!this.statusBarItem) return;

		switch (state) {
			case "connected":
				this.statusBarItem.text = "$(check) np3tool";
				this.statusBarItem.tooltip = `Recipe NP3 Editor — ${this.panels.size} active session(s)`;
				this.statusBarItem.show();
				break;
			case "error":
				this.statusBarItem.text = "$(error) np3tool";
				this.statusBarItem.tooltip = "Recipe NP3 Editor — connection error";
				this.statusBarItem.show();
				break;
			case "disconnected":
				if (this.panels.size === 0) {
					this.statusBarItem.hide();
				} else {
					this.statusBarItem.text = "$(circle-outline) np3tool";
					this.statusBarItem.tooltip = `Recipe NP3 Editor — ${this.panels.size} active session(s)`;
				}
				break;
		}
	}

	/**
	 * Common Save As logic shared between UI button and Command Palette.
	 */
	private async executeSaveAs(
		webviewPanel: vscode.WebviewPanel,
		binaryManager: BinaryManager,
		document: vscode.CustomDocument,
	) {
		const uri = await vscode.window.showSaveDialog({
			filters: { "Nikon NP3": ["np3"] },
			defaultUri: document.uri,
		});

		if (uri) {
			try {
				const response = await binaryManager.send({
					type: "np3.save_as",
					payload: { filePath: uri.fsPath },
				});

				// Update dirty state on successful save_as
				this.updateDirtyState(document.uri.toString(), response);

				webviewPanel.webview.postMessage(response);
			} catch (err) {
				webviewPanel.webview.postMessage({
					type: "error",
					payload: { message: (err as Error).message, code: "SAVE_AS_FAILED" },
				});
			}
		}
	}

	private getHtmlForWebview(webview: vscode.Webview): string {
		const distPath = vscode.Uri.joinPath(this.extensionUri, "dist", "webview");

		const scriptUri = webview.asWebviewUri(vscode.Uri.joinPath(distPath, "webview.js"));
		const styleUri = webview.asWebviewUri(vscode.Uri.joinPath(distPath, "webview.css"));

		const nonce = getNonce();

		return `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src ${webview.cspSource}; script-src 'nonce-${nonce}';">
	<link rel="stylesheet" href="${styleUri}">
	<title>Recipe NP3 Editor</title>
</head>
<body>
	<div id="app"></div>
	<script nonce="${nonce}" src="${scriptUri}"></script>
</body>
</html>`;
	}
}

function getNonce(): string {
	let text = "";
	const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
	for (let i = 0; i < 32; i++) {
		text += possible.charAt(Math.floor(Math.random() * possible.length));
	}
	return text;
}
