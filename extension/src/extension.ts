import * as vscode from "vscode";
import { Np3EditorPanel } from "./panels/Np3EditorPanel";

export function activate(context: vscode.ExtensionContext) {
	const outputChannel = vscode.window.createOutputChannel("recipe");
	outputChannel.appendLine("Recipe NP3 Editor extension activated");

	const provider = new Np3EditorPanel(context, outputChannel);

	context.subscriptions.push(
		vscode.window.registerCustomEditorProvider("recipe.np3Editor", provider, {
			webviewOptions: { retainContextWhenHidden: true },
			supportsMultipleEditorsPerDocument: false,
		}),
	);

	context.subscriptions.push(outputChannel);
}

export function deactivate() {}
