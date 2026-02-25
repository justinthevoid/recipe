<script lang="ts">
	import { AlertTriangle, FileWarning, RefreshCcw } from "lucide-svelte";
	import { Button } from "$lib/components/ui/button";

	export let error: { message: string, code: string } | null = null;
	export let filename: string = "Unknown File";
	export let onRetry: () => void = () => {};
</script>

<div class="h-full flex flex-col items-center justify-center p-8 bg-stripes-danger bg-opacity-10 dark:bg-opacity-5 rounded-lg border-2 border-dashed border-red-500/30">
	<div class="max-w-md w-full text-center space-y-6 bg-(--vscode-editor-background) p-8 rounded-xl shadow-xl border border-red-900/20">
		<div class="flex justify-center">
			<div class="relative">
				<FileWarning size={64} class="text-red-500 opacity-80" />
				<AlertTriangle size={32} class="absolute -bottom-2 -right-2 text-yellow-500 fill-zinc-900" />
			</div>
		</div>

		<div class="space-y-2">
			<h2 class="text-xl font-semibold text-red-500">File Corrupted or Unsupported</h2>
			<p class="text-sm text-zinc-400 font-mono bg-zinc-900/50 py-1 px-3 rounded inline-block truncate max-w-full">
				{filename}
			</p>
		</div>

		{#if error}
			<div class="bg-red-500/10 border border-red-500/20 rounded-md p-4 text-left">
				<p class="text-xs font-mono text-red-400 mb-1 border-b border-red-500/20 pb-1">Code: {error.code}</p>
				<p class="text-sm text-red-300">{error.message}</p>
			</div>
		{/if}

		<div class="text-sm text-zinc-400 leading-relaxed text-left space-y-2">
			<p>This .np3 file appears to be damaged or uses an unsupported format version. The internal checksum validation failed.</p>
			<p>For your safety, the recipe editor operates in read-only mode to prevent further data corruption.</p>
		</div>

		<div class="pt-4 flex justify-center">
			<Button variant="outline" class="gap-2 text-(--vscode-button-foreground) bg-(--vscode-button-background) hover:bg-(--vscode-button-hoverBackground) border-none" onclick={onRetry}>
				<RefreshCcw size={16} />
				Try Reloading
			</Button>
		</div>
	</div>
</div>

<style>
	/* Warning hash pattern */
	.bg-stripes-danger {
		background-image: linear-gradient(
			45deg,
			rgba(239, 68, 68, 0.1) 25%,
			transparent 25%,
			transparent 50%,
			rgba(239, 68, 68, 0.1) 50%,
			rgba(239, 68, 68, 0.1) 75%,
			transparent 75%,
			transparent
		);
		background-size: 32px 32px;
	}
</style>
