<script>
import { addFile, wasmState } from "../stores-compat";

let isDragging = false;
let fileInput;

function handleDragOver(e) {
	e.preventDefault();
	isDragging = true;
}

function handleDragLeave(e) {
	e.preventDefault();
	isDragging = false;
}

function handleDrop(e) {
	e.preventDefault();
	isDragging = false;
	if (e.dataTransfer.files) {
		processFiles(e.dataTransfer.files);
	}
}

function handleFileSelect(e) {
	if (e.target.files) {
		processFiles(e.target.files);
	}
}

function processFiles(fileList) {
	Array.from(fileList).forEach((file) => {
		const ext = file.name.split(".").pop().toLowerCase();
		if (["np3", "xmp"].includes(ext)) {
			addFile(file);
		} else {
			console.warn("Skipping invalid file:", file.name);
		}
	});
	if (fileInput) fileInput.value = "";
}
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	id="upload"
	class="glass-regular rounded-2xl p-8 text-center border-2 border-dashed transition-all {isDragging ? 'border-interactive bg-interactive/5' : 'border-border hover:border-interactive/50 hover:bg-surface-hover'}"
	on:dragover={handleDragOver}
	on:dragleave={handleDragLeave}
	on:drop={handleDrop}
	role="button"
	tabindex="0"
>
	<svg class="mx-auto mb-4 w-12 h-12 text-foreground-muted" fill="currentColor" viewBox="0 0 24 24">
		<path d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96zM14 13v4h-4v-4H7l5-5 5 5h-3z" />
	</svg>
	<h2 class="text-lg font-semibold text-foreground mb-1">Drag & Drop Presets</h2>
	<p class="text-sm text-foreground-muted mb-4">Supports .np3 and .xmp</p>

	<button
		class="px-5 py-2.5 rounded-xl font-medium bg-interactive text-interactive-foreground transition-all hover:-translate-y-0.5 hover:shadow-lg"
		on:click={() => fileInput.click()}
	>
		Select Files
	</button>
	<input
		bind:this={fileInput}
		type="file"
		multiple
		hidden
		accept=".np3,.xmp"
		on:change={handleFileSelect}
	/>

	<div class="mt-4 text-xs font-medium {$wasmState.status === 'ready' ? 'text-success' : $wasmState.status === 'error' ? 'text-error' : 'text-foreground-muted'}">
		{#if $wasmState.status === "initializing"}
			Initializing Engine...
		{:else if $wasmState.status === "ready"}
			Ready to convert
		{:else if $wasmState.status === "error"}
			Error: {$wasmState.error}
		{/if}
	</div>
</div>
