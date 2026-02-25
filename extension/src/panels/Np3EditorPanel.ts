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
		return { uri, dispose: () => {} };
	}

	async resolveCustomEditor(
		_document: vscode.CustomDocument,
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
			await binaryManager.start();
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

		// Forward webview messages to Go binary
		webviewPanel.webview.onDidReceiveMessage(
			async (message: IpcMessage) => {
				this.outputChannel.appendLine(`Webview → Go: ${JSON.stringify(message)}`);
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
