<script lang="ts">
import { toast } from "svelte-sonner";
import ColorBlender from "$lib/components/ColorBlender.svelte";
import ColorGrading from "$lib/components/ColorGrading.svelte";
import CorruptedFileView from "$lib/components/CorruptedFileView.svelte";
import ParameterSliderUnit from "$lib/components/ParameterSliderUnit.svelte";
import SchemaDropdown from "$lib/components/SchemaDropdown.svelte";
import ToneCurveVisual from "$lib/components/ToneCurveVisual.svelte";
import { Toaster } from "$lib/components/ui/sonner";
import { getNested, np3AppStore } from "$lib/state/np3.svelte";
import type { IpcMessage, Np3Error, Np3OpenResponse, ParameterDefinition } from "$lib/types";
import { vscode } from "$lib/vscode";

let status = $state<string>("Connecting...");
let originalMetadata = $state<Np3OpenResponse | null>(null);
let isDirty = $state(false);
let ipcPending = $state(false);
let pendingCount = $state(0);

function trackIpcSend() {
	pendingCount++;
	ipcPending = true;
}

function trackIpcReceive() {
	pendingCount = Math.max(0, pendingCount - 1);
	ipcPending = pendingCount > 0;
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
function handleParameterChange(key: string, value: number | string) {
	if (!np3AppStore.metadata) return;

	np3AppStore.patchParameter(key, value);
	trackIpcSend();
	vscode.postMessage({ type: "np3.patch", payload: { field: key, value } });
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
function handleResetAll() {
	trackIpcSend();
	vscode.postMessage({ type: "np3.reset", payload: {} });
}

function handleUndo() {
	if (!np3AppStore.canUndo) return;
	np3AppStore.undo();
	// Sync with backend after undo — send only recipe (F36)
	trackIpcSend();
	vscode.postMessage({ type: "np3.sync_request", payload: { recipe: np3AppStore.metadata?.recipe } });
	isDirty = np3AppStore.computeIsDirty();
}

function handleRedo() {
	if (!np3AppStore.canRedo) return;
	np3AppStore.redo();
	// Sync with backend after redo — send only recipe (F36)
	trackIpcSend();
	vscode.postMessage({ type: "np3.sync_request", payload: { recipe: np3AppStore.metadata?.recipe } });
	isDirty = np3AppStore.computeIsDirty();
}

function handleSave() {
	if (!isDirty) return;
	trackIpcSend();
	vscode.postMessage({ type: "np3.save", payload: {} });
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
	} else if (isMod && key === 's') {
		if (event.shiftKey) {
			handleSaveAs();
		} else {
			handleSave();
		}
		event.preventDefault();
	}
}

function handleMessage(event: MessageEvent<IpcMessage>) {
	const message = event.data;

	switch (message.type) {
		case "np3.save_as_success": {
			trackIpcReceive();
			const payload = message.payload as { filePath: string; dirty?: boolean };
			np3AppStore.saveAsSuccessful(payload.filePath);
			np3AppStore.markSaved();
			isDirty = false;
			toast.success("File saved as new profile");
			break;
		}
		case "np3.paste_response":
			if (message.payload) {
				const success = np3AppStore.handlePaste(message.payload as string);
				if (success) {
					toast.success("Parameters pasted from clipboard");
					// Sync with backend after paste — send only recipe (F36)
					vscode.postMessage({ type: "np3.sync_request", payload: { recipe: np3AppStore.metadata?.recipe } });
					isDirty = true;
				} else if (np3AppStore.currentError) {
					toast.error(np3AppStore.currentError.message);
				}
			}
			break;
		case "np3.metadata": {
			trackIpcReceive();
			np3AppStore.loadSuccess(message.payload as Np3OpenResponse);
			originalMetadata = structuredClone(message.payload as Np3OpenResponse);
			np3AppStore.markSaved();
			isDirty = false;
			status = "Connected";
			// Persist state for webview recovery (P2-9a)
			vscode.setState({ metadata: message.payload, originalMetadata });
			toast.success("File loaded successfully");
			break;
		}
		case "np3.patch_success": {
			trackIpcReceive();
			isDirty = true;
			break;
		}
		case "np3.save_success": {
			trackIpcReceive();
			np3AppStore.markSaved();
			isDirty = false;
			toast.success("File saved");
			break;
		}
		case "np3.save_error": {
			trackIpcReceive();
			const errPayload = message.payload as Np3Error;
			toast.error(errPayload.message || "Failed to save file");
			break;
		}
		case "np3.sync": {
			trackIpcReceive();
			if (message.payload) {
				const syncPayload = message.payload as { recipe?: unknown; dirty?: boolean };
				if (syncPayload.recipe && np3AppStore.metadata) {
					np3AppStore.metadata = {
						...np3AppStore.metadata,
						recipe: structuredClone(syncPayload.recipe) as Np3OpenResponse["recipe"],
					};
				}
			}
			break;
		}
		case "np3.sync_error": {
			trackIpcReceive();
			const syncErrPayload = message.payload as Np3Error;
			toast.error(syncErrPayload.message || "Sync failed");
			break;
		}
		case "np3.patch_error":
			trackIpcReceive();
			np3AppStore.rollback();
			isDirty = np3AppStore.computeIsDirty();
			toast.error((message.payload as Np3Error).message || "Failed to edit parameter");
			break;
		case "error":
		case "extension.error": {
			trackIpcReceive();
			status = "Error";
			const errorPayload = message.payload as Np3Error;
			if (errorPayload.code === "ERR_INVALID_CHECKSUM" || errorPayload.code === "ERR_CORRUPTED_FILE") {
				np3AppStore.loadError(errorPayload);
			}
			toast.error(errorPayload.message);
			break;
		}
		case "extension.open": {
			trackIpcSend();
			const filePath = (message.payload as { filePath: string }).filePath;
			np3AppStore.currentFilename = filePath.split(/[\\/]/).pop() || filePath;
			vscode.postMessage({ type: "np3.open", payload: { filePath } });
			break;
		}
		case "extension.triggerReset": {
			handleResetAll();
			break;
		}
		default:
			console.log("Unknown message type:", message.type);
	}
}

function handleSaveAs() {
	vscode.postMessage({ type: "np3.save_as", payload: {} });
}

function handleCopy() {
	np3AppStore.copyParameters();
	toast.info("Parameters copied to clipboard");
}

function handlePaste() {
	np3AppStore.pasteParameters();
}

function handleOpenBackup() {
	vscode.postMessage({ type: "np3.open_backup", payload: {} });
}

function handleRevealInFinder() {
	vscode.postMessage({ type: "np3.reveal_in_finder", payload: {} });
}

$effect(() => {
	window.addEventListener("message", handleMessage as EventListener);
	window.addEventListener("keydown", handleKeyDown);

	// Restore state if available (P2-9a)
	const savedState = vscode.getState<{ metadata: Np3OpenResponse; originalMetadata: Np3OpenResponse }>();
	if (savedState?.metadata) {
		np3AppStore.loadSuccess(savedState.metadata);
		originalMetadata = savedState.originalMetadata;
		np3AppStore.markSaved();
		isDirty = false;
		status = "Connected";
	}

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
	{#if ipcPending}
		<div class="fixed top-0 left-0 right-0 h-0.5 bg-(--vscode-progressBar-background) z-50 animate-pulse"></div>
	{/if}

	<header class="flex-none max-w-7xl mx-auto w-full flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold mb-4">
				Recipe NP3 Editor
				{#if isDirty}
					<span class="text-sm font-normal opacity-60 ml-2">[Modified]</span>
				{/if}
			</h1>
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
						title="Copy Parameters (Ctrl+C)"
						aria-label="Copy Parameters"
					>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><title>Copy</title><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg>
					</button>
					<button
						type="button"
						class="p-1.5 hover:bg-(--vscode-toolbar-hoverBackground) rounded text-sm font-medium transition-colors disabled:opacity-30"
						onclick={handlePaste}
						disabled={!np3AppStore.isLoaded}
						title="Paste Parameters (Ctrl+V)"
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
					class="px-3 py-1.5 bg-(--vscode-button-background) text-(--vscode-button-foreground) text-sm font-medium rounded hover:opacity-90 disabled:opacity-30 transition-opacity"
					onclick={handleSave}
					disabled={!isDirty}
					title="Save (Ctrl+S)"
					aria-label="Save"
				>
					Save
				</button>
				<button
					type="button"
					class="px-3 py-1.5 bg-(--vscode-button-secondaryBackground) text-(--vscode-button-secondaryForeground) text-sm font-medium rounded hover:opacity-90 disabled:opacity-30 transition-opacity mr-2"
					onclick={handleSaveAs}
					disabled={!np3AppStore.isLoaded}
					title="Save As... (Ctrl+Shift+S)"
					aria-label="Save As"
				>
					Save As...
				</button>
				<button
					type="button"
					class="px-3 py-1.5 bg-(--vscode-button-secondaryBackground) text-(--vscode-button-secondaryForeground) text-sm font-medium rounded hover:opacity-90 disabled:opacity-30 transition-opacity"
					onclick={handleResetAll}
					disabled={!np3AppStore.isLoaded}
					aria-label="Reset All"
					title="Reset All — revert to last saved state"
				>
					Reset All
				</button>
				<button
					type="button"
					class="p-1.5 hover:bg-(--vscode-toolbar-hoverBackground) rounded text-sm transition-colors opacity-50 hover:opacity-100"
					title="Keyboard Shortcuts&#10;&#10;Ctrl+S — Save&#10;Ctrl+Shift+S — Save As&#10;Ctrl+Z — Undo&#10;Ctrl+Shift+Z — Redo&#10;Ctrl+C — Copy Parameters&#10;Ctrl+V — Paste Parameters&#10;Double-click label — Reset to default"
					aria-label="Keyboard shortcuts help"
				>
					?
				</button>
			</div>
		{/if}
	</header>

	<div class="flex-1 w-full max-w-7xl mx-auto overflow-hidden min-h-0">
		{#if np3AppStore.isCorrupted}
			<CorruptedFileView
				error={np3AppStore.currentError}
				filename={np3AppStore.currentFilename}
				onOpenBackup={handleOpenBackup}
				onRevealInFinder={handleRevealInFinder}
			/>
		{:else if np3AppStore.isLoaded && np3AppStore.metadata}
			{@const grouped = np3AppStore.groupedParameters}
			{@const leftLane = grouped.filter(g => g.parameters[0]?.lane === 'left')}
			{@const rightLane = grouped.filter(g => g.parameters[0]?.lane !== 'left')}

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
					{#each leftLane as group}
						<section class="flex flex-col gap-4">
							<header class="flex items-center justify-between border-l-2 border-(--vscode-button-background) pl-3 mb-2">
								<h2 class="text-xs font-bold uppercase tracking-widest opacity-70">{group.name}</h2>
							</header>
							<div class="bg-muted/10 p-5 rounded-xl border border-border/50 backdrop-blur-sm flex flex-col gap-5">
								{#if group.name === 'Tone Curve'}
									<ToneCurveVisual
										points={np3AppStore.metadata.recipe?.pointCurve || []}
									/>
								{/if}

								{#each group.parameters as param}
									{#if param.type === 'discrete'}
										<SchemaDropdown
											definition={param}
											value={Number(getNested(np3AppStore.metadata?.recipe, param.key) ?? param.defaultValue)}
											originalValue={Number(getNested(originalMetadata?.recipe, param.key) ?? param.defaultValue)}
											onchange={handleParameterChange}
										/>
									{:else}
										<ParameterSliderUnit
											definition={param}
											value={Number(getNested(np3AppStore.metadata?.recipe, param.key) ?? param.defaultValue)}
											originalValue={Number(getNested(originalMetadata?.recipe, param.key) ?? param.defaultValue)}
											onchange={handleParameterChange}
										/>
									{/if}
								{/each}
							</div>
						</section>
					{/each}
				</div>

				<!-- Right Lane: Color, Detail, System -->
				<div class="h-full overflow-y-auto pr-4 flex flex-col gap-6 scrollbar-hide pb-20">
					{#each rightLane as group}
						<section class="flex flex-col gap-4">
							<header class="flex items-center justify-between border-l-2 border-(--vscode-button-background) pl-3 mb-2">
								<h2 class="text-xs font-bold uppercase tracking-widest opacity-70">{group.name}</h2>
							</header>
							<div class="bg-muted/10 p-5 rounded-xl border border-border/50 backdrop-blur-sm">
								{#if group.name === 'Color Mixer'}
									<ColorBlender
										parameters={group.parameters}
										recipe={np3AppStore.metadata.recipe}
										originalRecipe={originalMetadata?.recipe || {}}
										onchange={handleParameterChange}
									/>
								{:else if group.name === 'Color Grading'}
									<ColorGrading
										parameters={group.parameters}
										recipe={np3AppStore.metadata.recipe}
										originalRecipe={originalMetadata?.recipe || {}}
										onchange={handleParameterChange}
									/>
								{:else}
									<div class="flex flex-col gap-5">
										{#each group.parameters as param}
											{#if param.type === 'discrete'}
												<SchemaDropdown
													definition={param}
													value={Number(getNested(np3AppStore.metadata?.recipe, param.key) ?? param.defaultValue)}
													originalValue={Number(getNested(originalMetadata?.recipe, param.key) ?? param.defaultValue)}
													onchange={handleParameterChange}
												/>
											{:else}
												<ParameterSliderUnit
													definition={param}
													value={Number(getNested(np3AppStore.metadata?.recipe, param.key) ?? param.defaultValue)}
													originalValue={Number(getNested(originalMetadata?.recipe, param.key) ?? param.defaultValue)}
													onchange={handleParameterChange}
												/>
											{/if}
										{/each}
									</div>
								{/if}
							</div>
						</section>
					{/each}
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
