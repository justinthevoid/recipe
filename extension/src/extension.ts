import * as vscode from "vscode";
import { Np3EditorPanel } from "./panels/Np3EditorPanel";

let providerInstance: Np3EditorPanel | null = null;
let statusBarItem: vscode.StatusBarItem | null = null;

export function activate(context: vscode.ExtensionContext) {
	const outputChannel = vscode.window.createOutputChannel("recipe");
	outputChannel.appendLine("Recipe NP3 Editor extension activated");

	// Create status bar item (P2-4c)
	statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
	statusBarItem.text = "$(circle-outline) np3tool";
	statusBarItem.tooltip = "Recipe NP3 Editor — no active session";
	context.subscriptions.push(statusBarItem);

	const provider = new Np3EditorPanel(context, outputChannel, statusBarItem);
	providerInstance = provider;

	const providerRegistration = vscode.window.registerCustomEditorProvider(
		"recipe.np3Editor",
		provider,
		{
			webviewOptions: { retainContextWhenHidden: true },
			supportsMultipleEditorsPerDocument: false,
		},
	);

	context.subscriptions.push(providerRegistration);

	// Register Save As command
	context.subscriptions.push(
		vscode.commands.registerCommand("recipe.saveAs", async () => {
			provider.triggerSaveAs();
		}),
	);

	// Register Reset All command
	context.subscriptions.push(
		vscode.commands.registerCommand("recipe.resetAll", () => {
			provider.triggerResetAll();
		}),
	);

	context.subscriptions.push(outputChannel);
}

export function deactivate() {
	providerInstance?.stopAll();
	providerInstance = null;
	statusBarItem?.dispose();
	statusBarItem = null;
}
