<script lang="ts">
import { toast } from "svelte-sonner";
import CorruptedFileView from "$lib/components/CorruptedFileView.svelte";
import ParameterSliderUnit from "$lib/components/ParameterSliderUnit.svelte";
import SchemaDropdown from "$lib/components/SchemaDropdown.svelte";
import { Toaster } from "$lib/components/ui/sonner";
import { np3AppStore } from "$lib/state/np3.svelte";
import type { IpcMessage, Np3Error, Np3OpenResponse, ParameterDefinition } from "$lib/types";
import { vscode } from "$lib/vscode";

let status = $state<string>("Connecting...");
let originalMetadata = $state<Np3OpenResponse | null>(null);

function handleParameterChange(key: string, value: number) {
	if (!np3AppStore.metadata) return;
	
	np3AppStore.patch((meta) => {
		const newMeta = structuredClone($state.snapshot(meta)) as Np3OpenResponse;
		if (!newMeta.recipe) newMeta.recipe = {} as Record<string, unknown>;
		(newMeta.recipe as Record<string, unknown>)[key] = value;
		return newMeta;
	});

	vscode.postMessage({ type: "np3.patch", payload: { field: key, value } });
}

function handleResetAll() {
	vscode.postMessage({ type: "np3.reset", payload: {} });
}

function handleUndo() {
	if (!np3AppStore.canUndo) return;
	np3AppStore.undo();
	// Sync with backend after undo
	vscode.postMessage({ type: "np3.sync_request", payload: np3AppStore.metadata });
}

function handleRedo() {
	if (!np3AppStore.canRedo) return;
	np3AppStore.redo();
	// Sync with backend after redo
	vscode.postMessage({ type: "np3.sync_request", payload: np3AppStore.metadata });
}

function handleKeyDown(event: KeyboardEvent) {
	const isMod = event.ctrlKey || event.metaKey;
	if (isMod && event.key.toLowerCase() === 'z') {
		if (event.shiftKey) {
			handleRedo();
		} else {
			handleUndo();
		}
		event.preventDefault();
	}
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
			originalMetadata = structuredClone(message.payload as Np3OpenResponse);
			toast.success("File loaded successfully");
			break;
		case "np3.sync":
			if (message.payload) {
				np3AppStore.metadata = structuredClone(message.payload as Np3OpenResponse);
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
		case "extension.open": {
			const filePath = (message.payload as { filePath: string }).filePath;
			np3AppStore.currentFilename = filePath.split(/[\\/]/).pop() || filePath;
			vscode.postMessage({ type: "np3.open", payload: { filePath } });
			break;
		}
		default:
			console.log("Unknown message type:", message.type);
	}
}

function handleRetry() {
	np3AppStore.clearError();
	vscode.postMessage({ type: "webview.ready", payload: {} });
}

$effect(() => {
	window.addEventListener("message", handleMessage as EventListener);
	window.addEventListener("keydown", handleKeyDown);
	vscode.postMessage({ type: "np3.ping", payload: {} });
	vscode.postMessage({ type: "webview.ready", payload: {} });

	return () => {
		window.removeEventListener("message", handleMessage as EventListener);
		window.removeEventListener("keydown", handleKeyDown);
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
			<div class="flex items-center gap-2">
				<div class="flex items-center gap-1 border-r border-border pr-2 mr-2">
					<button
						type="button"
						class="p-1.5 hover:bg-(--vscode-toolbar-hoverBackground) rounded disabled:opacity-30 disabled:hover:bg-transparent transition-colors"
						onclick={handleUndo}
						disabled={!np3AppStore.canUndo}
						title="Undo (Ctrl+Z)"
					>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><title>Undo</title><path d="M3 7v6h6"/><path d="M21 17a9 9 0 0 0-9-9 9 9 0 0 0-6 2.3L3 13"/></svg>
					</button>
					<button
						type="button"
						class="p-1.5 hover:bg-(--vscode-toolbar-hoverBackground) rounded disabled:opacity-30 disabled:hover:bg-transparent transition-colors"
						onclick={handleRedo}
						disabled={!np3AppStore.canRedo}
						title="Redo (Ctrl+Shift+Z)"
					>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><title>Redo</title><path d="M21 7v6h-6"/><path d="M3 17a9 9 0 0 1 9-9 9 9 0 0 1 6 2.3L21 13"/></svg>
					</button>
				</div>
				<button
					type="button"
					class="px-3 py-1.5 bg-(--vscode-button-background) text-(--vscode-button-foreground) text-sm font-medium rounded hover:opacity-90 transition-opacity"
					onclick={handleResetAll}
				>
					Reset All
				</button>
			</div>
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
								{#if param.type === 'continuous'}
									<ParameterSliderUnit 
										definition={param} 
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{:else if param.type === 'discrete'}
									<SchemaDropdown
										definition={param}
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{/if}
							{/each}
						</div>
					</div>
					<div class="bg-card text-card-foreground p-4 rounded-lg border border-border">
						<h2 class="text-lg font-medium mb-2">Tone Curve</h2>
						<div class="flex flex-col gap-4">
							{#each parameters.filter(p => p.group === 'Tone Curve') as param}
								{#if param.type === 'continuous'}
									<ParameterSliderUnit 
										definition={param} 
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{:else if param.type === 'discrete'}
									<SchemaDropdown
										definition={param}
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{/if}
							{/each}
						</div>
					</div>
					<div class="bg-card text-card-foreground p-4 rounded-lg border border-border">
						<h2 class="text-lg font-medium mb-2">Detail</h2>
						<div class="flex flex-col gap-4">
							{#each parameters.filter(p => p.group === 'Detail') as param}
								{#if param.type === 'continuous'}
									<ParameterSliderUnit 
										definition={param} 
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{:else if param.type === 'discrete'}
									<SchemaDropdown
										definition={param}
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{/if}
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
								{#if param.type === 'continuous'}
									<ParameterSliderUnit 
										definition={param} 
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{:else if param.type === 'discrete'}
									<SchemaDropdown
										definition={param}
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{/if}
							{/each}
						</div>
					</div>
					<div class="bg-card text-card-foreground p-4 rounded-lg border border-border">
						<h2 class="text-lg font-medium mb-2">Geometry</h2>
						<div class="flex flex-col gap-4">
							{#each parameters.filter(p => p.group === 'Geometry') as param}
								{#if param.type === 'continuous'}
									<ParameterSliderUnit 
										definition={param} 
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{:else if param.type === 'discrete'}
									<SchemaDropdown
										definition={param}
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{/if}
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
