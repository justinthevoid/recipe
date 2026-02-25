<script lang="ts">
import { vscode } from "$lib/vscode";

interface IpcMessage {
	type: string;
	payload: Record<string, unknown>;
}

let status = $state<string>("Connecting...");
let toastMessage = $state<string | null>(null);
let toastType = $state<"error" | "success">("error");

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
		case "error":
			status = "Error";
			showToast((message.payload as { message: string }).message, "error");
			break;
		default:
			console.log("Unknown message type:", message.type);
	}
}

$effect(() => {
	window.addEventListener("message", handleMessage);
	// Send initial ping
	vscode.postMessage({ type: "np3.ping", payload: {} });

	return () => {
		window.removeEventListener("message", handleMessage);
	};
});
</script>

<main class="min-h-screen bg-[var(--vscode-editor-background)] text-[var(--vscode-editor-foreground)] p-4">
	<div class="max-w-2xl mx-auto">
		<h1 class="text-2xl font-bold mb-4">Recipe NP3 Editor</h1>
		<div class="flex items-center gap-2 mb-4">
			<span class="inline-block w-2 h-2 rounded-full" class:bg-green-500={status.includes("Connected")} class:bg-yellow-500={status === "Connecting..."} class:bg-red-500={status === "Error"}></span>
			<span class="text-sm opacity-80">{status}</span>
		</div>
	</div>

	{#if toastMessage}
		<div
			class="fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg max-w-sm text-sm"
			class:bg-red-900={toastType === "error"}
			class:text-red-100={toastType === "error"}
			class:bg-green-900={toastType === "success"}
			class:text-green-100={toastType === "success"}
		>
			{toastMessage}
		</div>
	{/if}
</main>
