<script lang="ts">
	import { onMount } from "svelte";
	import { initWasm, isWasmReady, wasmExtractFullRecipe } from "$lib/wasm.svelte";
	import { wasmStatusStore } from "$lib/shared-stores";
	import { convertFile } from "$lib/converter.svelte";
	import { detectFormatFromExtension, getOppositeFormat, getFormatLabel } from "$lib/format-detector";
	import { countParameters } from "$lib/parameter-counter";
	import { openPreset } from "$lib/stores.svelte";
	import type { UniversalRecipe } from "@recipe/ui";

	type CardState = "idle" | "file-ready" | "converting" | "complete" | "error";

	let cardState = $state<CardState>("idle");
	let isDragging = $state(false);
	let isDemo = $state(false);
	let wasmStatus = $state(wasmStatusStore.get());

	// File state
	let fileName = $state("");
	let fileSize = $state(0);
	let fileData = $state<Uint8Array | null>(null);
	let sourceFormat = $state("");
	let targetFormat = $state("");

	// Result state
	let resultData = $state<Uint8Array | null>(null);
	let resultFileName = $state("");
	let mappedCount = $state(0);
	let skippedCount = $state(0);
	let paramExpanded = $state(false);
	let convertedRecipe = $state<Record<string, unknown> | null>(null);

	// Error state
	let errorMessage = $state("");

	let fileInput = $state<HTMLInputElement>(undefined!);

	onMount(() => {
		initWasm();
		const unsub = wasmStatusStore.subscribe(v => { wasmStatus = v; });
		return unsub;
	});

	function formatBytes(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		return `${(bytes / 1024).toFixed(1)} KB`;
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave() {
		isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		const file = e.dataTransfer?.files[0];
		if (file) processFile(file);
	}

	function handleFileSelect(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (file) processFile(file);
		input.value = "";
	}

	async function processFile(file: File) {
		const format = detectFormatFromExtension(file.name);
		if (format === "unknown") return; // Silently reject invalid files

		fileName = file.name;
		fileSize = file.size;
		sourceFormat = format;
		targetFormat = getOppositeFormat(format);

		const buffer = await file.arrayBuffer();
		fileData = new Uint8Array(buffer);

		cardState = "file-ready";
		isDemo = false;
	}

	async function processData(data: Uint8Array, name: string, format: string) {
		fileName = name;
		fileSize = data.byteLength;
		fileData = data;
		sourceFormat = format;
		targetFormat = getOppositeFormat(format);
		cardState = "file-ready";
	}

	async function handleConvert() {
		if (!fileData) return;
		cardState = "converting";

		// Brief animation delay for perceived quality
		await new Promise(r => setTimeout(r, 350));

		try {
			const result = await convertFile(fileData, sourceFormat, targetFormat, fileName);
			resultData = result.data;
			resultFileName = result.fileName;

			// Extract recipe to count parameters
			const recipe = await wasmExtractFullRecipe(result.data, targetFormat);
			convertedRecipe = recipe;
			const counts = countParameters(recipe, sourceFormat);
			mappedCount = counts.mapped;
			skippedCount = counts.skipped;

			paramExpanded = isDemo;
			cardState = "complete";
		} catch (err: unknown) {
			const error = err as { userMessage?: string; message?: string };
			errorMessage = error.userMessage ?? error.message ?? "Conversion failed.";
			cardState = "error";
		}
	}

	function handleDownload() {
		if (!resultData) return;
		const blob = new Blob([resultData]);
		const url = URL.createObjectURL(blob);
		const a = document.createElement("a");
		a.href = url;
		a.download = resultFileName;
		a.click();
		URL.revokeObjectURL(url);
	}

	function handleEditInPreview() {
		if (!convertedRecipe) return;
		openPreset(convertedRecipe as UniversalRecipe, resultFileName);
	}

	function reset() {
		cardState = "idle";
		fileData = null;
		resultData = null;
		convertedRecipe = null;
		isDemo = false;
		paramExpanded = false;
	}

	async function handleDemo() {
		isDemo = true;
		try {
			const response = await fetch("/demo-preset.np3");
			const buffer = await response.arrayBuffer();
			const data = new Uint8Array(buffer);
			await processData(data, "demo-preset.np3", "np3");
			// Auto-convert
			await handleConvert();
		} catch {
			errorMessage = "Failed to load demo preset.";
			cardState = "error";
		}
	}
</script>

<div
	class="glass-regular rounded-2xl transition-all duration-300 relative overflow-hidden
		{cardState === 'idle' && isDragging ? 'border-interactive border-2' : ''}
		{cardState === 'idle' && !isDragging ? 'border-dashed border-2 border-foreground-muted/20' : ''}
		{cardState === 'file-ready' ? 'border-solid border border-foreground-muted/30 shadow-lg shadow-aurora-violet/5' : ''}
		{cardState === 'converting' ? 'border-solid border border-foreground-muted/30' : ''}
		{cardState === 'complete' ? 'border-solid border border-success/20 shadow-lg shadow-success/5' : ''}
		{cardState === 'error' ? 'border-solid border-2 border-error/40 shadow-lg shadow-error/10' : ''}"
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
	role="region"
	aria-label="Preset converter"
>
	<!-- Converting shimmer overlay -->
	{#if cardState === "converting"}
		<div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/5 to-transparent animate-shimmer pointer-events-none"></div>
	{/if}

	<div class="p-6 md:p-8">
		{#if cardState === "idle"}
			<!-- Idle state -->
			<div class="flex flex-col items-center gap-4 py-4">
				<div class="text-4xl opacity-50">☁️</div>
				<p class="text-sm text-foreground-muted text-center">
					Drop your <span class="font-medium text-foreground">.np3</span> or <span class="font-medium text-foreground">.xmp</span> file here
				</p>
				<button
					type="button"
					class="px-5 py-2.5 text-sm font-medium bg-interactive text-interactive-foreground rounded-lg hover:opacity-90 transition-opacity"
					onclick={() => fileInput.click()}
				>
					Select File
				</button>
				<button
					type="button"
					class="text-xs text-foreground-muted hover:text-interactive transition-colors"
					onclick={handleDemo}
					disabled={wasmStatus !== "ready"}
				>
					or try a demo
				</button>
				<input
					bind:this={fileInput}
					type="file"
					accept=".np3,.xmp"
					class="hidden"
					onchange={handleFileSelect}
				/>
			</div>
			<!-- WASM status -->
			<div class="text-center mt-2">
				{#if wasmStatus === "loading"}
					<span class="text-xs text-foreground-muted/50">Initializing engine…</span>
				{:else if wasmStatus === "error"}
					<span class="text-xs text-error">Engine failed to load</span>
				{:else if wasmStatus === "ready"}
					<span class="text-xs text-success/50">Ready</span>
				{/if}
			</div>

		{:else if cardState === "file-ready"}
			<!-- File ready state -->
			<div class="flex flex-col gap-4">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-3">
						<span class="text-xs font-bold uppercase px-2 py-0.5 rounded
							{sourceFormat === 'np3' ? 'bg-aurora-violet/20 text-aurora-violet' : 'bg-aurora-cyan/20 text-aurora-cyan'}">
							{getFormatLabel(sourceFormat)}
						</span>
						<div>
							<p class="text-sm font-medium text-foreground truncate max-w-[200px]">{fileName}</p>
							<p class="text-xs text-foreground-muted">{formatBytes(fileSize)}</p>
						</div>
					</div>
					<button
						type="button"
						class="text-foreground-muted hover:text-foreground text-lg leading-none transition-colors"
						onclick={reset}
						aria-label="Clear file"
					>
						×
					</button>
				</div>
				<button
					type="button"
					class="w-full py-3 text-sm font-medium bg-interactive text-interactive-foreground rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50"
					disabled={wasmStatus !== "ready"}
					onclick={handleConvert}
				>
					{wasmStatus !== "ready" ? "Initializing engine…" : `Convert to ${getFormatLabel(targetFormat)}`}
				</button>
			</div>

		{:else if cardState === "converting"}
			<!-- Converting state -->
			<div class="flex flex-col items-center gap-4 py-4">
				<div class="text-sm text-foreground-muted animate-pulse">Converting…</div>
			</div>

		{:else if cardState === "complete"}
			<!-- Results state -->
			<div class="flex flex-col gap-4">
				<div class="flex items-center gap-3">
					<span class="text-xs font-bold uppercase px-2 py-0.5 rounded
						{targetFormat === 'np3' ? 'bg-aurora-violet/20 text-aurora-violet' : 'bg-aurora-cyan/20 text-aurora-cyan'}">
						{getFormatLabel(targetFormat)}
					</span>
					<p class="text-sm font-medium text-foreground">
						Converted to {getFormatLabel(targetFormat)}
					</p>
				</div>

				<!-- Parameter count -->
				<details class="text-xs text-foreground-muted" bind:open={paramExpanded}>
					<summary class="cursor-pointer hover:text-foreground transition-colors select-none">
						{mappedCount} mapped · {skippedCount} skipped
					</summary>
					{#if convertedRecipe}
						<div class="mt-2 p-3 bg-surface-elevated/50 rounded-lg max-h-48 overflow-y-auto">
							<pre class="text-xs text-foreground-muted whitespace-pre-wrap">{JSON.stringify(convertedRecipe, null, 2)}</pre>
						</div>
					{/if}
				</details>

				<!-- Download -->
				<button
					type="button"
					class="w-full py-3 text-sm font-medium bg-interactive text-interactive-foreground rounded-lg hover:opacity-90 transition-opacity"
					onclick={handleDownload}
				>
					Download {resultFileName}
				</button>

				<!-- Secondary actions -->
				<div class="flex items-center gap-4 justify-center">
					<button
						type="button"
						class="text-xs text-foreground-muted hover:text-interactive transition-colors"
						onclick={handleEditInPreview}
					>
						Edit in Preview
					</button>
					<button
						type="button"
						class="text-xs text-foreground-muted hover:text-interactive transition-colors"
						onclick={reset}
					>
						Convert another
					</button>
				</div>

				<!-- Install hints -->
				<details class="text-xs text-foreground-muted">
					<summary class="cursor-pointer hover:text-foreground transition-colors select-none">
						Where does this go?
					</summary>
					<div class="mt-2 space-y-1.5 text-xs">
						{#if targetFormat === "xmp"}
							<p><strong class="text-foreground">macOS:</strong> <code class="text-foreground-muted/80">~/Library/Application Support/Adobe/Lightroom/Develop Presets/</code></p>
							<p><strong class="text-foreground">Windows:</strong> <code class="text-foreground-muted/80">%APPDATA%\Adobe\Lightroom\Develop Presets\</code></p>
							<p><strong class="text-foreground">Or:</strong> Lightroom → Develop → Presets → Import</p>
						{:else}
							<p><strong class="text-foreground">Nikon:</strong> Import via NX Studio → Picture Control → Import</p>
						{/if}
					</div>
				</details>
			</div>

		{:else if cardState === "error"}
			<!-- Error state -->
			<div class="flex flex-col items-center gap-4 py-4">
				<p class="text-sm text-error text-center">{errorMessage}</p>
				<button
					type="button"
					class="px-5 py-2.5 text-sm font-medium bg-interactive text-interactive-foreground rounded-lg hover:opacity-90 transition-opacity"
					onclick={reset}
				>
					Try again
				</button>
			</div>
		{/if}
	</div>
</div>

<style>
	@keyframes shimmer {
		0% { transform: translateX(-100%); }
		100% { transform: translateX(100%); }
	}
	.animate-shimmer {
		animation: shimmer 1.5s infinite;
	}
</style>
