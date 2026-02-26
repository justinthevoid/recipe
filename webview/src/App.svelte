<script lang="ts">
import { toast } from "svelte-sonner";
import CorruptedFileView from "$lib/components/CorruptedFileView.svelte";
import ParameterSliderUnit from "$lib/components/ParameterSliderUnit.svelte";
import { Toaster } from "$lib/components/ui/sonner";
import { np3AppStore } from "$lib/state/np3.svelte";
import type { IpcMessage, Np3Error, Np3OpenResponse, ParameterDefinition } from "$lib/types";
import { vscode } from "$lib/vscode";

let status = $state<string>("Connecting...");
let originalMetadata = $state<Np3OpenResponse | null>(null);

// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
function handleParameterChange(key: string, value: number) {
	if (!np3AppStore.metadata) return;
	
	np3AppStore.patch((meta) => {
		const newMeta = JSON.parse(JSON.stringify(meta));
		if (!newMeta.recipe) newMeta.recipe = {};
		newMeta.recipe[key] = value;
		return newMeta;
	});

	vscode.postMessage({ type: "np3.patch", payload: { field: key, value } });
}

// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
function handleResetAll() {
	vscode.postMessage({ type: "np3.reset", payload: {} });
}

function handleMessage(event: MessageEvent<IpcMessage>) {
	const message = event.data;

	switch (message.type) {
		case "np3.pong":
			status = `Connected — ${(message.payload as { status: string }).status}`;
			toast.success("Binary connection established");
			break;
		case "np3.metadata":
			np3AppStore.loadSuccess(message.payload as Np3OpenResponse);
			originalMetadata = JSON.parse(JSON.stringify(message.payload));
			toast.success("File loaded successfully");
			break;
		case "np3.sync":
			if (message.payload) {
				np3AppStore.metadata = JSON.parse(JSON.stringify(message.payload));
			}
			break;
		case "np3.patch_error":
			np3AppStore.rollback();
			toast.error((message.payload as Np3Error).message || "Failed to edit parameter");
			break;
		case "error":
		case "extension.error": {
			status = "Error";
			const errorPayload = message.payload as Np3Error;
			if (errorPayload.code === "ERR_INVALID_CHECKSUM" || errorPayload.code === "ERR_CORRUPTED_FILE") {
				np3AppStore.loadError(errorPayload);
			}
			toast.error(errorPayload.message);
			break;
		}
		// A mock trigger from extension to simulate opening a file
		case "extension.open": {
			const filePath = (message.payload as { filePath: string }).filePath;
			np3AppStore.currentFilename = filePath.split(/[\\/]/).pop() || filePath;
			// Pass message to Go backend
			vscode.postMessage({ type: "np3.open", payload: { filePath } });
			break;
		}
		default:
			console.log("Unknown message type:", message.type);
	}
}

// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
function handleRetry() {
    // Retry loading current file if possible, or clear state
	np3AppStore.clearError();
}

$effect(() => {
	window.addEventListener("message", handleMessage as EventListener);
	// Send initial ping
	vscode.postMessage({ type: "np3.ping", payload: {} });
	// Tell extension host we are ready to receive the file to open
	vscode.postMessage({ type: "webview.ready", payload: {} });

	return () => {
		window.removeEventListener("message", handleMessage as EventListener);
	};
});
</script>

<main
	class="h-screen w-full overflow-hidden bg-background text-foreground p-4 flex flex-col gap-4"
>
	<header class="flex-none max-w-7xl mx-auto w-full flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold mb-4">Recipe NP3 Editor</h1>
			<div class="flex items-center gap-2">
				<span
					class="inline-block w-2 h-2 rounded-full"
					class:bg-(--vscode-testing-iconPassed)={status.includes("Connected")}
					class:bg-(--vscode-editorWarning-foreground)={status === "Connecting..."}
					class:bg-(--vscode-errorForeground)={status === "Error"}
				></span>
				<span class="text-sm opacity-80">{status}</span>
			</div>
		</div>
		{#if np3AppStore.isLoaded && !np3AppStore.isCorrupted}
			<button
				type="button"
				class="px-3 py-1.5 bg-(--vscode-button-background) text-(--vscode-button-foreground) text-sm font-medium rounded hover:opacity-90 transition-opacity"
				onclick={handleResetAll}
			>
				Reset All
			</button>
		{/if}
	</header>

	<div class="flex-1 w-full max-w-7xl mx-auto overflow-hidden min-h-0">
		{#if np3AppStore.isCorrupted}
			<CorruptedFileView error={np3AppStore.currentError} filename={np3AppStore.currentFilename} onRetry={handleRetry} />
		{:else if np3AppStore.isLoaded && np3AppStore.metadata}
			{@const parameters = np3AppStore.metadata.parameters || []}
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 h-full">
				<!-- Left Lane: Core Edits -->
				<div class="h-full overflow-y-auto pr-2 flex flex-col gap-4">
					<div class="bg-card text-card-foreground p-4 rounded-lg border border-border">
						<h2 class="text-lg font-medium mb-2">Basic</h2>
						<div class="flex flex-col gap-4">
							{#each parameters.filter(p => p.group === 'Basic') as param}
								<ParameterSliderUnit 
									definition={param} 
									value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
									originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
									onchange={handleParameterChange}
								/>
							{/each}
						</div>
					</div>
					<div class="bg-card text-card-foreground p-4 rounded-lg border border-border">
						<h2 class="text-lg font-medium mb-2">Tone Curve</h2>
						<div class="flex flex-col gap-4">
							{#each parameters.filter(p => p.group === 'Tone Curve') as param}
								<ParameterSliderUnit 
									definition={param} 
									value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
									originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
									onchange={handleParameterChange}
								/>
							{/each}
						</div>
					</div>
					<div class="bg-card text-card-foreground p-4 rounded-lg border border-border">
						<h2 class="text-lg font-medium mb-2">Detail</h2>
						<div class="flex flex-col gap-4">
							{#each parameters.filter(p => p.group === 'Detail') as param}
								<ParameterSliderUnit 
									definition={param} 
									value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
									originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
									onchange={handleParameterChange}
								/>
							{/each}
						</div>
					</div>
				</div>

				<!-- Right Lane: Grading/Effects -->
				<div class="h-full overflow-y-auto pr-2 flex flex-col gap-4">
					<div class="bg-card text-card-foreground p-4 rounded-lg border border-border">
						<h2 class="text-lg font-medium mb-2">Color Mixer</h2>
						<div class="flex flex-col gap-4">
							{#each parameters.filter(p => p.group === 'Color Mixer') as param}
								<ParameterSliderUnit 
									definition={param} 
									value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
									originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
									onchange={handleParameterChange}
								/>
							{/each}
						</div>
					</div>
					<div class="bg-card text-card-foreground p-4 rounded-lg border border-border">
						<h2 class="text-lg font-medium mb-2">Geometry</h2>
						<div class="flex flex-col gap-4">
							{#each parameters.filter(p => p.group === 'Geometry') as param}
								<ParameterSliderUnit 
									definition={param} 
									value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
									originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
									onchange={handleParameterChange}
								/>
							{/each}
						</div>
					</div>

				</div>
			</div>
		{:else}
			<div class="h-full flex flex-col items-center justify-center border-2 border-dashed border-border rounded-lg text-muted-foreground">
				<p>Waiting for NP3 file...</p>
			</div>
		{/if}
	</div>

	<Toaster 
		theme="system" 
		position="bottom-right" 
		toastOptions={{
			classes: {
				toast: "bg-[var(--vscode-notifications-background)] text-[var(--vscode-notifications-foreground)] border-[var(--vscode-widget-border)]",
				title: "text-[var(--vscode-notifications-foreground)]",
				description: "text-[var(--vscode-notifications-foreground)] opacity-90",
				error: "!text-[var(--vscode-notificationsErrorIcon-foreground)]",
				actionButton: "bg-[var(--vscode-button-background)] text-[var(--vscode-button-foreground)]",
				cancelButton: "bg-[var(--vscode-button-secondaryBackground)] text-[var(--vscode-button-secondaryForeground)]",
			}
		}}
	/>
</main>
