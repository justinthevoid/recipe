<script lang="ts">
import { toast } from "svelte-sonner";
import ColorBlender from "$lib/components/ColorBlender.svelte";
import ColorGrading from "$lib/components/ColorGrading.svelte";
import CorruptedFileView from "$lib/components/CorruptedFileView.svelte";
import ParameterSliderUnit from "$lib/components/ParameterSliderUnit.svelte";
import SchemaDropdown from "$lib/components/SchemaDropdown.svelte";
import ToneCurveVisual from "$lib/components/ToneCurveVisual.svelte";
import { Toaster } from "$lib/components/ui/sonner";
import { np3AppStore } from "$lib/state/np3.svelte";
import type { IpcMessage, Np3Error, Np3OpenResponse, ParameterDefinition } from "$lib/types";
import { vscode } from "$lib/vscode";

let status = $state<string>("Connecting...");
let originalMetadata = $state<Np3OpenResponse | null>(null);

// biome-ignore lint/correctness/noUnusedVariables: used in template
function handleParameterChange(key: string, value: number | string) {
	if (!np3AppStore.metadata) return;
	
	np3AppStore.patch((meta) => {
		// Optimization: Shallow copy the top-level object and the recipe
		// since patch() already handled saving the previous state to history.
		const newMeta = { ...meta };
		newMeta.recipe = {
			...meta.recipe,
			[key]: value
		};
		return newMeta as Np3OpenResponse;
	});

	vscode.postMessage({ type: "np3.patch", payload: { field: key, value } });
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
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
	const key = event.key.toLowerCase();

	if (isMod && key === 'z') {
		if (event.shiftKey) {
			handleRedo();
		} else {
			handleUndo();
		}
		event.preventDefault();
	} else if (isMod && key === 'c') {
		handleCopy();
		event.preventDefault();
	} else if (isMod && key === 'v') {
		handlePaste();
		event.preventDefault();
	}
}

function handleMessage(event: MessageEvent<IpcMessage>) {
	const message = event.data;

	switch (message.type) {
		case "np3.save_as_success":
			np3AppStore.saveAsSuccessful((message.payload as { filePath: string }).filePath);
			toast.success("File saved as new profile");
			break;
		case "np3.paste_response":
			if (message.payload) {
				const success = np3AppStore.handlePaste(message.payload as string);
				if (success) {
					toast.success("Parameters pasted from clipboard");
					// Sync with backend after paste
					vscode.postMessage({ type: "np3.sync_request", payload: np3AppStore.metadata });
				} else if (np3AppStore.currentError) {
					toast.error(np3AppStore.currentError.message);
				}
			}
			break;
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

// biome-ignore lint/correctness/noUnusedVariables: used in template
function handleSaveAs() {
	vscode.postMessage({ type: "np3.save_as", payload: {} });
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
function handleCopy() {
	np3AppStore.copyParameters();
	toast.info("Parameters copied to clipboard");
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
function handlePaste() {
	np3AppStore.pasteParameters();
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
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
						class="p-1.5 hover:bg-(--vscode-toolbar-hoverBackground) rounded text-sm font-medium transition-colors disabled:opacity-30"
						onclick={handleCopy}
						disabled={!np3AppStore.isLoaded}
						title="Copy Parameters"
						aria-label="Copy Parameters"
					>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><title>Copy</title><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg>
					</button>
					<button
						type="button"
						class="p-1.5 hover:bg-(--vscode-toolbar-hoverBackground) rounded text-sm font-medium transition-colors disabled:opacity-30"
						onclick={handlePaste}
						disabled={!np3AppStore.isLoaded}
						title="Paste Parameters"
						aria-label="Paste Parameters"
					>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><title>Paste</title><rect width="8" height="4" x="8" y="2" rx="1" ry="1"/><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/></svg>
					</button>
				</div>
				<div class="flex items-center gap-1 border-r border-border pr-2 mr-2">
					<button
						type="button"
						class="p-1.5 hover:bg-(--vscode-toolbar-hoverBackground) rounded disabled:opacity-30 disabled:hover:bg-transparent transition-colors"
						onclick={handleUndo}
						disabled={!np3AppStore.canUndo}
						title="Undo (Ctrl+Z)"
						aria-label="Undo"
					>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><title>Undo</title><path d="M3 7v6h6"/><path d="M21 17a9 9 0 0 0-9-9 9 9 0 0 0-6 2.3L3 13"/></svg>
					</button>
					<button
						type="button"
						class="p-1.5 hover:bg-(--vscode-toolbar-hoverBackground) rounded disabled:opacity-30 disabled:hover:bg-transparent transition-colors"
						onclick={handleRedo}
						disabled={!np3AppStore.canRedo}
						title="Redo (Ctrl+Shift+Z)"
						aria-label="Redo"
					>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><title>Redo</title><path d="M21 7v6h-6"/><path d="M3 17a9 9 0 0 1 9-9 9 9 0 0 1 6 2.3L21 13"/></svg>
					</button>
				</div>
				<button
					type="button"
					class="px-3 py-1.5 bg-(--vscode-button-secondaryBackground) text-(--vscode-button-secondaryForeground) text-sm font-medium rounded hover:opacity-90 disabled:opacity-30 transition-opacity mr-2"
					onclick={handleSaveAs}
					disabled={!np3AppStore.isLoaded}
					aria-label="Save As"
				>
					Save As...
				</button>
				<button
					type="button"
					class="px-3 py-1.5 bg-(--vscode-button-background) text-(--vscode-button-foreground) text-sm font-medium rounded hover:opacity-90 disabled:opacity-30 transition-opacity"
					onclick={handleResetAll}
					disabled={!np3AppStore.isLoaded}
					aria-label="Reset All"
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
			{@const parameterDefinitions = np3AppStore.metadata.parameterDefinitions || []}
			
            <div class="flex flex-col gap-2 mb-6 px-1">
				<input 
					type="text" 
					class="text-2xl font-semibold bg-transparent border-b border-transparent hover:border-border focus:border-(--vscode-focusBorder) focus:outline-none py-1 transition-colors w-full"
					value={np3AppStore.metadata.recipe?.name || ""}
					placeholder="Recipe Name"
					onchange={(e) => handleParameterChange('name', e.currentTarget.value)}
				/>
				<input 
					type="text" 
					class="text-sm opacity-80 bg-transparent border-b border-transparent hover:border-border focus:border-(--vscode-focusBorder) focus:outline-none py-1 transition-colors w-full"
					value={np3AppStore.metadata.recipe?.description || ""}
					placeholder="Description (Optional)"
					onchange={(e) => handleParameterChange('description', e.currentTarget.value)}
				/>
			</div>

			<div class="grid grid-cols-1 lg:grid-cols-2 gap-8 h-full">
				<!-- Left Lane: Tone & Curve -->
				<div class="h-full overflow-y-auto pr-4 flex flex-col gap-6 scrollbar-hide pb-20">
					<!-- Basic Toning -->
					<section class="flex flex-col gap-4">
						<header class="flex items-center justify-between border-l-2 border-(--vscode-button-background) pl-3 mb-2">
							<h2 class="text-xs font-bold uppercase tracking-widest opacity-70">Exposure & Tone</h2>
						</header>
						<div class="bg-muted/10 p-5 rounded-xl border border-border/50 backdrop-blur-sm flex flex-col gap-5">
							{#each parameterDefinitions.filter((p: ParameterDefinition) => p.group === 'Basic') as param}
								<ParameterSliderUnit 
									definition={param} 
									value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
									originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
									onchange={handleParameterChange}
								/>
							{/each}
						</div>
					</section>

					<!-- Tone Curve Visual -->
					<section class="flex flex-col gap-4">
						<header class="flex items-center justify-between border-l-2 border-(--vscode-button-background) pl-3 mb-2">
							<h2 class="text-xs font-bold uppercase tracking-widest opacity-70">Tone Curve</h2>
						</header>
						<div class="bg-muted/10 p-5 rounded-xl border border-border/50 backdrop-blur-sm">
							<ToneCurveVisual 
								points={np3AppStore.metadata.recipe?.pointCurve || []} 
							/>
							
							<!-- Parametric Controls -->
							<div class="mt-6 flex flex-col gap-4">
								{#each parameterDefinitions.filter((p: ParameterDefinition) => p.group === 'Tone Curve') as param}
									<ParameterSliderUnit 
										definition={param} 
										value={Number(np3AppStore.metadata?.recipe?.[param.key] ?? param.defaultValue)}
										originalValue={Number(originalMetadata?.recipe?.[param.key] ?? param.defaultValue)}
										onchange={handleParameterChange}
									/>
								{/each}
							</div>
						</div>
					</section>
				</div>

				<!-- Right Lane: Color & Detail -->
				<div class="h-full overflow-y-auto pr-4 flex flex-col gap-6 scrollbar-hide pb-20">
					<!-- Color Mixer -->
					<section class="flex flex-col gap-4">
						<header class="flex items-center justify-between border-l-2 border-(--vscode-button-background) pl-3 mb-2">
							<h2 class="text-xs font-bold uppercase tracking-widest opacity-70">Color Mixer</h2>
						</header>
						<div class="bg-muted/10 p-5 rounded-xl border border-border/50 backdrop-blur-sm">
							<ColorBlender 
								parameters={parameterDefinitions}
								recipe={np3AppStore.metadata.recipe}
								originalRecipe={originalMetadata?.recipe || {}}
								onchange={handleParameterChange}
							/>
						</div>
					</section>

					<!-- Color Grading -->
					<section class="flex flex-col gap-4">
						<header class="flex items-center justify-between border-l-2 border-(--vscode-button-background) pl-3 mb-2">
							<h2 class="text-xs font-bold uppercase tracking-widest opacity-70">Color Grading</h2>
						</header>
						<div class="bg-muted/10 p-5 rounded-xl border border-border/50 backdrop-blur-sm">
							<ColorGrading 
								parameters={parameterDefinitions}
								recipe={np3AppStore.metadata.recipe}
								originalRecipe={originalMetadata?.recipe || {}}
								onchange={handleParameterChange}
							/>
						</div>
					</section>
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
