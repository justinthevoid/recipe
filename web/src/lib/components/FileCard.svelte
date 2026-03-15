<script>
import { detectFormatFromExtension } from "../format-detector";
import { removeFile } from "../stores-compat";
import { wasmExtractFullRecipe } from "../wasm.svelte";
import { openPreset } from "../stores.svelte";

export let file;

$: format = detectFormatFromExtension(file.name);

let isLoading = false;

function formatSize(bytes) {
	if (bytes === 0) return "0 B";
	const k = 1024;
	const sizes = ["B", "KB", "MB", "GB"];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return parseFloat((bytes / k ** i).toFixed(2)) + " " + sizes[i];
}

async function handleEdit() {
	if (isLoading) return;
	isLoading = true;
	try {
		const buffer = await file.file.arrayBuffer();
		const bytes = new Uint8Array(buffer);
		const fmt = detectFormatFromExtension(file.name);
		const recipe = await wasmExtractFullRecipe(bytes, fmt);
		openPreset(recipe, file.name);
	} catch (err) {
		console.error("Failed to open preset for editing:", err);
	} finally {
		isLoading = false;
	}
}
</script>

<div class="glass-regular rounded-xl p-4 transition-all hover:-translate-y-0.5 hover:shadow-lg">
	<div class="flex items-center justify-between">
		<div class="min-w-0 flex-1">
			<div class="text-sm font-medium text-foreground truncate">{file.name}</div>
			<div class="text-xs text-foreground-muted mt-0.5 flex items-center gap-2">
				{formatSize(file.size)}
				{#if file.status === "processing"}
					<span class="text-interactive">Converting...</span>
				{:else if file.status === "complete"}
					<span class="text-success">Converted</span>
				{:else if file.status === "error"}
					<span class="text-error">{file.error || "Error"}</span>
				{/if}
			</div>
		</div>
		<div class="flex items-center gap-2 ml-3">
			<span class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-interactive/10 text-interactive">
				{format.toUpperCase()}
			</span>

			{#if file.status === "complete"}
				<a
					href={file.outputUrl}
					download={file.outputName}
					class="p-1.5 rounded-lg text-success hover:bg-success/10 transition-colors"
					title="Download"
				>
					<svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
					</svg>
				</a>
			{/if}

			<!-- Edit button — opens WebGL2 editor -->
			<button
				class="px-2 py-1 rounded-lg text-xs font-medium text-interactive hover:bg-interactive/10 transition-colors disabled:opacity-50"
				on:click={handleEdit}
				disabled={isLoading}
				title="Edit in preview"
			>
				{isLoading ? "Loading..." : "Edit"}
			</button>

			<button
				class="p-1.5 rounded-lg text-foreground-muted hover:text-error hover:bg-error/10 transition-colors"
				on:click={() => removeFile(file.id)}
				aria-label="Remove file"
				title="Remove"
			>
				<svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>
		</div>
	</div>
</div>
