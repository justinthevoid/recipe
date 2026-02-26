<script lang="ts">
	import { AlertTriangle, FileWarning, RefreshCcw } from "lucide-svelte";
	import { Button } from "$lib/components/ui/button";

	let { error = null, filename = "Unknown File", onRetry = () => {} }: {
		error: { message: string, code: string, rawData?: string } | null;
		filename: string;
		onRetry: () => void;
	} = $props();

	let hexDump = $derived.by(() => {
		if (!error?.rawData) return "";
		try {
			const binary = atob(error.rawData);
			let hex = "";
			for (let i = 0; i < Math.min(binary.length, 1024); i += 16) {
				const chunk = binary.slice(i, i + 16);
				const hexChunk = Array.from(chunk).map(c => c.charCodeAt(0).toString(16).padStart(2, '0')).join(' ');
				const asciiChunk = Array.from(chunk).map(c => {
					const code = c.charCodeAt(0);
					return (code >= 32 && code <= 126) ? c : '.';
				}).join('');
				hex += `${i.toString(16).padStart(8, '0')}  ${hexChunk.padEnd(48, ' ')}  |${asciiChunk}|\n`;
			}
			if (binary.length > 1024) {
				hex += `... (${binary.length - 1024} more bytes truncated)`;
			}
			return hex;
		} catch {
			return "Failed to decode raw binary data.";
		}
	});
</script>

<div class="h-full flex flex-col p-8 bg-stripes-danger bg-opacity-10 dark:bg-opacity-5 rounded-lg border-2 border-dashed border-(--vscode-errorForeground)/30 overflow-y-auto">
	<div class="max-w-4xl mx-auto w-full space-y-6 bg-(--vscode-editor-background) p-8 rounded-xl shadow-xl border border-(--vscode-errorForeground)/20">
		<div class="flex items-start gap-6">
			<div class="relative shrink-0 mt-2">
				<FileWarning size={64} class="text-(--vscode-errorForeground) opacity-80" />
				<AlertTriangle size={32} class="absolute -bottom-2 -right-2 text-(--vscode-editorWarning-foreground) fill-background" />
			</div>
			
			<div class="flex-1 space-y-4">
				<div class="space-y-2">
					<h2 class="text-xl font-semibold text-(--vscode-errorForeground)">File Corrupted or Unsupported</h2>
					<p class="text-sm text-muted-foreground font-mono bg-muted py-1 px-3 rounded inline-block truncate max-w-full">
						{filename}
					</p>
				</div>

				<div class="text-sm text-muted-foreground leading-relaxed space-y-2">
					<p>This .np3 file appears to be damaged or uses an unsupported format version. The internal checksum validation failed.</p>
					<p>For your safety, the recipe editor operates in read-only mode to prevent further data corruption.</p>
				</div>
				
				<div class="pt-2">
					<Button variant="outline" class="gap-2 text-(--vscode-button-foreground) bg-(--vscode-button-background) hover:bg-(--vscode-button-hoverBackground) border-none" onclick={onRetry}>
						<RefreshCcw size={16} />
						Try Reloading
					</Button>
				</div>
			</div>
		</div>

		{#if error}
			<div class="bg-(--vscode-errorForeground)/10 border border-(--vscode-errorForeground)/20 rounded-md p-4 mt-6">
				<p class="text-xs font-mono text-(--vscode-errorForeground) mb-1 border-b border-(--vscode-errorForeground)/20 pb-1">Code: {error.code}</p>
				<p class="text-sm text-(--vscode-errorForeground) opacity-80">{error.message}</p>
			</div>

			{#if hexDump}
				<div class="mt-6">
					<h3 class="text-sm font-semibold text-foreground mb-2">Raw Binary Payload (Preview)</h3>
					<pre class="bg-muted p-4 rounded-md text-xs font-mono text-muted-foreground overflow-x-auto border border-border selection:bg-(--vscode-errorForeground)/20">{hexDump}</pre>
				</div>
			{/if}
		{/if}
	</div>
</div>

<style>
	/* Warning hash pattern — uses VS Code error color with low opacity */
	.bg-stripes-danger {
		background-image: linear-gradient(
			45deg,
			color-mix(in srgb, var(--vscode-errorForeground) 10%, transparent) 25%,
			transparent 25%,
			transparent 50%,
			color-mix(in srgb, var(--vscode-errorForeground) 10%, transparent) 50%,
			color-mix(in srgb, var(--vscode-errorForeground) 10%, transparent) 75%,
			transparent 75%,
			transparent
		);
		background-size: 32px 32px;
	}
</style>
