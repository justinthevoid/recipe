import * as vscode from "vscode";
import { Np3EditorPanel } from "./panels/Np3EditorPanel";

export function activate(context: vscode.ExtensionContext) {
	const outputChannel = vscode.window.createOutputChannel("recipe");
	outputChannel.appendLine("Recipe NP3 Editor extension activated");

	const provider = new Np3EditorPanel(context, outputChannel);
	const providerRegistration = vscode.window.registerCustomEditorProvider("recipe.np3Editor", provider, {
		webviewOptions: { retainContextWhenHidden: true },
		supportsMultipleEditorsPerDocument: false,
	});

	context.subscriptions.push(providerRegistration);

	// Register Save As command
	context.subscriptions.push(
		vscode.commands.registerCommand("recipe.saveAs", async () => {
			const activeEditor = vscode.window.activeTextEditor || vscode.window.tabGroups.activeTabGroup.activeTab;
			// Note: Custom editors don't always set activeTextEditor. 
			// We need to find the active Np3EditorPanel instance.
			// However, for Story 3.1, we'll assume the panel will trigger this via IPC 
			// OR we use the active tab.
			provider.triggerSaveAs();
		})
	);

	context.subscriptions.push(outputChannel);
}

export function deactivate() { }
