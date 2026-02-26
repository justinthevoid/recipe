<script lang="ts">
import { vscode } from "$lib/vscode";
import CorruptedFileView from "$lib/components/CorruptedFileView.svelte";
import type { IpcMessage, Np3Error, Np3OpenResponse } from "$lib/types";

let status = $state<string>("Connecting...");
let toastMessage = $state<string | null>(null);
let toastType = $state<"error" | "success">("error");

// App state
let isCorrupted = $state<boolean>(false);
let isLoaded = $state<boolean>(false);
let currentError = $state<Np3Error | null>(null);
let metadata = $state<Np3OpenResponse | null>(null);
let currentFilename = $state<string>("Unknown File");

function showToast(message: string, type: "error" | "success" = "error") {
	toastMessage = message;
	toastType = type;
	setTimeout(() => {
		toastMessage = null;
	}, 5000);
}

function handleMessage(event: MessageEvent<IpcMessage>) {
	const message = event.data;

	switch (message.type) {
		case "np3.pong":
			status = `Connected — ${(message.payload as { status: string }).status}`;
			showToast("Binary connection established", "success");
			break;
		case "np3.metadata":
			metadata = message.payload as Np3OpenResponse;
			isCorrupted = false;
			isLoaded = true;
			currentError = null;
			showToast("File loaded successfully", "success");
			break;
		case "error":
		case "extension.error": {
			status = "Error";
			const errorPayload = message.payload as Np3Error;
			if (errorPayload.code === "ERR_INVALID_CHECKSUM" || errorPayload.code === "ERR_CORRUPTED_FILE") {
				isCorrupted = true;
				currentError = errorPayload;
			}
			showToast(errorPayload.message, "error");
			break;
		}
		// A mock trigger from extension to simulate opening a file
		case "extension.open": {
			const filePath = (message.payload as { filePath: string }).filePath;
			currentFilename = filePath.split(/[\\/]/).pop() || filePath;
			// Pass message to Go backend
			vscode.postMessage({ type: "np3.open", payload: { filePath } });
			break;
		}
		default:
			console.log("Unknown message type:", message.type);
	}
}

function handleRetry() {
    // Retry loading current file if possible, or clear state
	isCorrupted = false;
	currentError = null;
}

$effect(() => {
	window.addEventListener("message", handleMessage as EventListener);
	// Send initial ping
	vscode.postMessage({ type: "np3.ping", payload: {} });

	return () => {
		window.removeEventListener("message", handleMessage as EventListener);
	};
});
</script>

<main
	class="min-h-screen bg-(--vscode-editor-background) text-(--vscode-editor-foreground) p-4 flex flex-col"
>
	<div class="max-w-2xl mx-auto w-full flex-none">
		<h1 class="text-2xl font-bold mb-4">Recipe NP3 Editor</h1>
		<div class="flex items-center gap-2 mb-4">
			<span
				class="inline-block w-2 h-2 rounded-full"
				class:bg-green-500={status.includes("Connected")}
				class:bg-yellow-500={status === "Connecting..."}
				class:bg-red-500={status === "Error"}
			></span>
			<span class="text-sm opacity-80">{status}</span>
		</div>
	</div>

	<div class="flex-grow w-full relative pt-4">
		{#if isCorrupted}
			<CorruptedFileView error={currentError} filename={currentFilename} onRetry={handleRetry} />
		{:else if isLoaded && metadata}
			<div class="bg-zinc-900/50 p-6 rounded-lg border border-zinc-800">
				<h2 class="text-lg font-medium text-zinc-200 mb-2">Recipe Metadata</h2>
				<p class="text-sm text-zinc-400 mb-4 font-mono">Hash: {metadata.hash}</p>
				<pre class="text-xs bg-black/40 p-4 rounded overflow-auto max-h-[500px] text-zinc-300">{JSON.stringify(metadata.recipe, null, 2)}</pre>
			</div>
		{:else}
			<div class="h-64 flex flex-col items-center justify-center border-2 border-dashed border-zinc-800 rounded-lg text-zinc-500">
				<p>Waiting for NP3 file...</p>
			</div>
		{/if}
	</div>

	{#if toastMessage}
		<div
			class="fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg max-w-sm text-sm z-50"
			class:bg-red-900={toastType === "error"}
			class:text-red-100={toastType === "error"}
			class:bg-green-900={toastType === "success"}
			class:text-green-100={toastType === "success"}
		>
			{toastMessage}
		</div>
	{/if}
</main>
