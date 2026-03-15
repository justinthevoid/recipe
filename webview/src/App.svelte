<script lang="ts">
import { toast } from "svelte-sonner";
import CollapsibleSection from "$lib/components/CollapsibleSection.svelte";
import ColorBlender from "$lib/components/ColorBlender.svelte";
import ColorGrading from "$lib/components/ColorGrading.svelte";
import CorruptedFileView from "$lib/components/CorruptedFileView.svelte";
import ParameterSliderUnit from "$lib/components/ParameterSliderUnit.svelte";
import PhotoPreview from "$lib/components/PhotoPreview.svelte";
import SchemaDropdown from "$lib/components/SchemaDropdown.svelte";
import ToneCurveVisual from "$lib/components/ToneCurveVisual.svelte";
import { Toaster } from "$lib/components/ui/sonner";
import { getNested, np3AppStore } from "$lib/state/np3.svelte";
import type { IpcMessage, Np3Error, Np3OpenResponse, ParameterDefinition } from "$lib/types";
import { vscode } from "$lib/vscode";
import { getWasmStatus, initWasm, wasmGenerateLUT } from "$lib/wasm.svelte";

let status = $state<string>("Connecting...");
let originalMetadata = $state<Np3OpenResponse | null>(null);
let isDirty = $state(false);
let ipcPending = $state(false);
let pendingCount = $state(0);

// Preview mode state
let isPreviewMode = $state(false);
let previewImageData = $state<HTMLImageElement | null>(null);
let previewImagePath = $state<string | null>(null);
let previewDividerWidth = $state(65); // percentage for left pane
let isDraggingDivider = $state(false);

const MAX_IMAGE_SIZE = 10 * 1024 * 1024; // 10MB
const MAX_IMAGE_DIMENSION = 2048;

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

// biome-ignore lint/correctness/noUnusedVariables: used in template
function handleOpenImage() {
	vscode.postMessage({ type: "np3.open_image", payload: {} });
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
function handleTogglePreview() {
	isPreviewMode = !isPreviewMode;
	// Persist preview mode state
	_persistPreviewState();
}

function _persistPreviewState() {
	const saved = vscode.getState<Record<string, unknown>>() || {};
	vscode.setState({
		...saved,
		isPreviewMode,
		previewImagePath,
		previewDividerWidth,
	});
}

function _decodeAndLoadImage(base64Data: string, filename: string) {
	const byteString = atob(base64Data);
	const bytes = new Uint8Array(byteString.length);
	for (let i = 0; i < byteString.length; i++) {
		bytes[i] = byteString.charCodeAt(i);
	}
	const blob = new Blob([bytes]);
	const url = URL.createObjectURL(blob);

	const img = new Image();
	img.onload = () => {
		if (img.naturalWidth > MAX_IMAGE_DIMENSION || img.naturalHeight > MAX_IMAGE_DIMENSION) {
			_downscaleImage(img, (downscaled) => {
				URL.revokeObjectURL(url);
				previewImageData = downscaled;
				isPreviewMode = true;
				_persistPreviewState();
			});
		} else {
			// Revoke after WebGL upload will consume the image data on next tick
			previewImageData = img;
			isPreviewMode = true;
			_persistPreviewState();
			// Defer revoke to allow WebGL texImage2D to read the image first
			setTimeout(() => URL.revokeObjectURL(url), 1000);
		}
	};
	img.onerror = () => {
		URL.revokeObjectURL(url);
		toast.error(`Failed to load image: ${filename}`);
	};
	img.src = url;
}

function _downscaleImage(img: HTMLImageElement, callback: (result: HTMLImageElement) => void) {
	const canvas = document.createElement("canvas");
	const scale = MAX_IMAGE_DIMENSION / Math.max(img.naturalWidth, img.naturalHeight);
	canvas.width = Math.round(img.naturalWidth * scale);
	canvas.height = Math.round(img.naturalHeight * scale);
	const ctx = canvas.getContext("2d");
	if (ctx) {
		ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
	}
	const downscaled = new Image();
	downscaled.onload = () => callback(downscaled);
	downscaled.src = canvas.toDataURL("image/jpeg", 0.92);
}

function _handleDividerPointerDown(event: PointerEvent) {
	isDraggingDivider = true;
	(event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
}

function _handleDividerPointerMove(event: PointerEvent) {
	if (!isDraggingDivider) return;
	const container = (event.currentTarget as HTMLElement).parentElement;
	if (!container) return;
	const rect = container.getBoundingClientRect();
	const x = event.clientX - rect.left;
	const pct = (x / rect.width) * 100;
	// Min 30% for photo, max ~75% (leaving 320px min for controls at typical widths)
	previewDividerWidth = Math.max(30, Math.min(75, pct));
}

function _handleDividerPointerUp() {
	isDraggingDivider = false;
	_persistPreviewState();
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
		case "extension.wasm_uri": {
			const wasmPayload = message.payload as { uri: string };
			initWasm(wasmPayload.uri);
			break;
		}
		case "np3.image_data": {
			const imgPayload = message.payload as { data: string; filename: string };
			previewImagePath = imgPayload.filename;
			_decodeAndLoadImage(imgPayload.data, imgPayload.filename);
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
	const savedState = vscode.getState<{
		metadata: Np3OpenResponse;
		originalMetadata: Np3OpenResponse;
		isPreviewMode?: boolean;
		previewImagePath?: string;
		previewDividerWidth?: number;
	}>();
	if (savedState?.metadata) {
		np3AppStore.loadSuccess(savedState.metadata);
		originalMetadata = savedState.originalMetadata;
		np3AppStore.markSaved();
		isDirty = false;
		status = "Connected";
	}
	if (savedState?.isPreviewMode) {
		isPreviewMode = savedState.isPreviewMode;
	}
	if (savedState?.previewDividerWidth) {
		previewDividerWidth = savedState.previewDividerWidth;
	}
	if (savedState?.previewImagePath) {
		previewImagePath = savedState.previewImagePath;
		// Re-request image from extension host using stored path
		vscode.postMessage({ type: "np3.open_image", payload: { filePath: previewImagePath } });
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
					class="px-3 py-1.5 bg-(--vscode-button-secondaryBackground) text-(--vscode-button-secondaryForeground) text-sm font-medium rounded hover:opacity-90 transition-opacity"
					onclick={handleTogglePreview}
					disabled={!np3AppStore.isLoaded}
					title={isPreviewMode ? "Switch to Editor view" : "Switch to Preview view"}
					aria-label={isPreviewMode ? "Editor" : "Preview"}
				>
					{isPreviewMode ? "Editor" : "Preview"}
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

			{#if isPreviewMode}
				<!-- Preview Mode: Photo left, controls right -->
				<div class="flex h-full min-h-0 relative" style:--divider-width="{previewDividerWidth}%">
					<!-- Photo Preview Pane -->
					<div class="h-full min-w-0" style="width: var(--divider-width)">
						<PhotoPreview
							recipe={np3AppStore.metadata.recipe}
							imageData={previewImageData}
							wasmReady={getWasmStatus() === 'ready'}
							generateLUT={wasmGenerateLUT}
						/>
					</div>

					<!-- Resizable Divider -->
					<div
						class="w-1 cursor-col-resize flex-none hover:bg-(--vscode-focusBorder) transition-colors z-10"
						class:bg-(--vscode-focusBorder)={isDraggingDivider}
						role="separator"
						aria-label="Resize preview"
						onpointerdown={_handleDividerPointerDown}
						onpointermove={_handleDividerPointerMove}
						onpointerup={_handleDividerPointerUp}
					></div>

					<!-- Controls Pane -->
					<div class="flex-1 h-full overflow-y-auto pl-4 pr-2 flex flex-col gap-4 scrollbar-hide pb-20 min-w-[320px]">
						{#each grouped as group}
							<CollapsibleSection
								title={group.name}
								expanded={group.name === 'Tone' || group.name === 'Color Mixer'}
							>
								{#if group.name === 'Tone Curve'}
									<ToneCurveVisual
										points={np3AppStore.metadata.recipe?.pointCurve || []}
									/>
								{/if}

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
							</CollapsibleSection>
						{/each}
					</div>
				</div>

				<!-- Preview Footer -->
				<footer class="flex-none flex items-center gap-3 pt-2 border-t border-border">
					<button
						type="button"
						class="px-3 py-1.5 bg-(--vscode-button-background) text-(--vscode-button-foreground) text-sm font-medium rounded hover:opacity-90 transition-opacity"
						onclick={handleOpenImage}
					>
						Open Image
					</button>
					{#if previewImagePath}
						<span class="text-xs opacity-50 truncate">{previewImagePath}</span>
					{/if}
					<span class="ml-auto text-xs opacity-40">
						{#if getWasmStatus() === 'ready'}
							WASM Ready
						{:else if getWasmStatus() === 'loading'}
							Loading WASM...
						{:else if getWasmStatus() === 'error'}
							WASM Error
						{:else}
							WASM Idle
						{/if}
					</span>
				</footer>
			{:else}
				<!-- Editor Mode: Original two-column layout -->
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
			{/if}
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
