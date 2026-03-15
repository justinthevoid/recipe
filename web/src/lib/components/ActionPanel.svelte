<script>
import { convertFile } from "../converter-compat";
import { detectFormatFromExtension } from "../format-detector";
import { files, settings, updateFileStatus } from "../stores-compat";

let converting = false;
let error = null;

const formats = [
	{ value: "xmp", label: "XMP (Adobe/Lightroom)" },
	{ value: "np3", label: "NP3 (Nikon)" },
];

if (!$settings.targetFormat) {
	settings.update((s) => ({ ...s, targetFormat: "xmp" }));
}

async function handleConvert() {
	if (converting || $files.length === 0) return;

	converting = true;
	error = null;
	const targetFormat = $settings.targetFormat;

	try {
		for (const fileData of $files) {
			if (fileData.status === "complete") continue;

			updateFileStatus(fileData.id, { status: "processing" });

			try {
				const buffer = await fileData.file.arrayBuffer();
				const bytes = new Uint8Array(buffer);
				const sourceFormat = detectFormatFromExtension(fileData.name);

				const outputData = await convertFile(bytes, sourceFormat, targetFormat, fileData.name);

				const blob = new Blob([outputData], { type: "application/octet-stream" });
				const url = URL.createObjectURL(blob);
				const newName = fileData.name.replace(/\.[^/.]+$/, "") + "." + targetFormat;

				updateFileStatus(fileData.id, {
					status: "complete",
					outputUrl: url,
					outputName: newName,
				});
			} catch (e) {
				console.error(e);
				updateFileStatus(fileData.id, {
					status: "error",
					error: e.message,
				});
			}
		}
	} catch (e) {
		error = e.message;
	} finally {
		converting = false;
	}
}

function downloadAll() {
	$files.forEach((f) => {
		if (f.status === "complete" && f.outputUrl) {
			const a = document.createElement("a");
			a.href = f.outputUrl;
			a.download = f.outputName;
			a.click();
		}
	});
}
</script>

<div class="glass-regular rounded-xl mt-4 p-5">
	<div class="flex flex-wrap items-center justify-between gap-4">
		<div class="flex items-center gap-3">
			<label for="format-select" class="text-sm font-medium text-foreground-muted">Convert to:</label>
			<div class="relative">
				<select
					id="format-select"
					bind:value={$settings.targetFormat}
					disabled={converting}
					class="appearance-none bg-canvas-base border border-border rounded-lg px-4 py-2 pr-8 text-sm text-foreground cursor-pointer min-w-[200px] transition-colors hover:border-interactive/50 focus:outline-none focus:border-focus focus:ring-1 focus:ring-focus"
				>
					{#each formats as format}
						<option value={format.value}>{format.label}</option>
					{/each}
				</select>
				<div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none border-l-[5px] border-r-[5px] border-t-[5px] border-l-transparent border-r-transparent border-t-foreground-muted"></div>
			</div>
		</div>

		<div class="flex gap-3">
			<button
				class="px-5 py-2.5 rounded-xl font-semibold bg-interactive text-interactive-foreground transition-all hover:-translate-y-0.5 hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
				on:click={handleConvert}
				disabled={converting || $files.length === 0}
			>
				{converting ? "Converting..." : "Convert All"}
			</button>

			{#if $files.some((f) => f.status === "complete")}
				<button
					class="px-4 py-2.5 rounded-xl font-medium glass-regular text-foreground transition-all hover:bg-surface-hover"
					on:click={downloadAll}
				>
					Download All
				</button>
			{/if}
		</div>
	</div>

	{#if error}
		<div class="mt-4 p-3 rounded-lg bg-error/10 border border-error/30 text-error text-sm">
			{error}
		</div>
	{/if}
</div>
