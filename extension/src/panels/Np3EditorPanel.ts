import * as vscode from "vscode";
import { BinaryManager, type IpcMessage } from "../backend/BinaryManager";

export class Np3EditorPanel implements vscode.CustomReadonlyEditorProvider {
	private readonly extensionUri: vscode.Uri;
	private readonly outputChannel: vscode.OutputChannel;

	constructor(
		private readonly context: vscode.ExtensionContext,
		outputChannel: vscode.OutputChannel,
	) {
		this.extensionUri = context.extensionUri;
		this.outputChannel = outputChannel;
	}

	openCustomDocument(
		uri: vscode.Uri,
		_openContext: vscode.CustomDocumentOpenContext,
		_token: vscode.CancellationToken,
	): vscode.CustomDocument {
		return { uri, dispose: () => { } };
	}

	async resolveCustomEditor(
		document: vscode.CustomDocument,
		webviewPanel: vscode.WebviewPanel,
		_token: vscode.CancellationToken,
	): Promise<void> {
		webviewPanel.webview.options = {
			enableScripts: true,
			localResourceRoots: [vscode.Uri.joinPath(this.extensionUri, "dist", "webview")],
		};

		webviewPanel.webview.html = this.getHtmlForWebview(webviewPanel.webview);

		// Spawn Go binary
		const binaryManager = new BinaryManager(this.context, this.outputChannel);

		try {
			await binaryManager.start(document.uri.fsPath);
			this.outputChannel.appendLine("Go binary started successfully");
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : "Unknown error starting Go binary";
			this.outputChannel.appendLine(`Failed to start Go binary: ${errorMessage}`);
			webviewPanel.webview.postMessage({
				type: "error",
				payload: { message: errorMessage, code: "BINARY_SPAWN_FAILED" },
			});
			return;
		}

		// Handle webview messages
		webviewPanel.webview.onDidReceiveMessage(
			async (message: IpcMessage) => {
				this.outputChannel.appendLine(`Webview → Extension: ${JSON.stringify(message)}`);

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
							webviewPanel.webview.postMessage(response);
						} catch (err) {
							webviewPanel.webview.postMessage({
								type: "error",
								payload: { message: (err as Error).message, code: "SAVE_AS_FAILED" },
							});
						}
					}
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

				try {
					const response = await binaryManager.send(message);
					this.outputChannel.appendLine(`Go → Webview: ${JSON.stringify(response)}`);
					webviewPanel.webview.postMessage(response);
				} catch (err) {
					const errorMessage = err instanceof Error ? err.message : "IPC communication error";
					this.outputChannel.appendLine(`IPC Error: ${errorMessage}`);
					webviewPanel.webview.postMessage({
						type: "error",
						payload: { message: errorMessage, code: "IPC_ERROR" },
					});
				}
			},
			undefined,
			this.context.subscriptions,
		);

		// Clean up binary on panel close
		webviewPanel.onDidDispose(() => {
			binaryManager.stop();
		});
	}

	/**
	 * Public method to trigger Save As from command palette.
	 * In a real custom editor, we'd track active panels.
	 */
	public triggerSaveAs() {
		// This is a simplified implementation. 
		// In a production app, we would broadcast to the active webview.
		// For now, we'll rely on the webview UI button to trigger the IPC "np3.save_as".
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
