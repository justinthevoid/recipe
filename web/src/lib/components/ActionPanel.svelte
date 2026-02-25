<script>
import { convertFile } from "../converter";
import { detectFormatFromExtension } from "../format-detector";
import { currentRecipe, files, previewFile, settings, updateFileStatus } from "../stores";

let converting = false;
let error = null;

// Supported target formats
const formats = [
	{ value: "xmp", label: "XMP (Adobe/Lightroom)" },
	{ value: "lrtemplate", label: "LRTEMPLATE (Legacy Lightroom)" },
	{ value: "np3", label: "NP3 (Nikon)" },
	{ value: "costyle", label: "COSTYLE (Capture One)" },
	{ value: "dcp", label: "DCP (DNG Profile)" },
];

// Default target format
if (!$settings.targetFormat) {
	settings.update((s) => ({ ...s, targetFormat: "xmp" }));
}

async function handleConvert() {
	if (converting || $files.length === 0) return;

	converting = true;
	error = null;
	const targetFormat = $settings.targetFormat;

	try {
		// Process files sequentially
		for (const fileData of $files) {
			if (fileData.status === "complete") continue; // Skip already converted

			updateFileStatus(fileData.id, { status: "processing" });

			try {
				const buffer = await fileData.file.arrayBuffer();
				const bytes = new Uint8Array(buffer);
				const sourceFormat = detectFormatFromExtension(fileData.name);

				const outputData = await convertFile(bytes, sourceFormat, targetFormat, fileData.name);

				// Create download URL
				const blob = new Blob([outputData], {
					type: "application/octet-stream",
				});
				const url = URL.createObjectURL(blob);

				// Generate new filename
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
	// Simple implementation: trigger click on all download links
	// A better approach would be to zip them, but for now this works for small batches
	$files.forEach((f) => {
		if (f.status === "complete" && f.outputUrl) {
			const a = document.createElement("a");
			a.href = f.outputUrl;
			a.download = f.outputName;
			a.click();
		}
	});
}

function createNewPreset() {
	const defaultRecipe = {
		Name: "New Preset",
		Exposure: 0,
		Contrast: 0,
		Highlights: 0,
		Shadows: 0,
		Whites: 0,
		Blacks: 0,
		Clarity: 0,
		Dehaze: 0,
		Vibrance: 0,
		Saturation: 0,
		Temperature: 0,
		Tint: 0,
		ColorGrading: {
			Highlights: { Hue: 0, Chroma: 0, Brightness: 0 },
			Midtone: { Hue: 0, Chroma: 0, Brightness: 0 },
			Shadows: { Hue: 0, Chroma: 0, Brightness: 0 },
			Blending: 50,
			Balance: 0,
		},
		ToneCurveHighlights: 0,
		ToneCurveLights: 0,
		ToneCurveDarks: 0,
		ToneCurveShadows: 0,
	};

	// Create a mock file object
	const newFile = {
		name: "New Preset.np3",
		isNew: true,
		arrayBuffer: async () => new ArrayBuffer(0),
	};

	previewFile.set(newFile);
	// currentRecipe will be set by PreviewModal based on isNew flag
}
</script>

<div class="action-panel glass-card">
    <div class="controls">
        <div class="control-group">
            <label for="format-select">Convert to:</label>
            <div class="select-wrapper">
                <select
                    id="format-select"
                    bind:value={$settings.targetFormat}
                    disabled={converting}
                >
                    {#each formats as format}
                        <option value={format.value}>{format.label}</option>
                    {/each}
                </select>
            </div>
        </div>

        <div class="actions">
            <button
                class="btn-primary"
                on:click={handleConvert}
                disabled={converting || $files.length === 0}
            >
                {#if converting}
                    Converting...
                {:else}
                    Convert All
                {/if}
            </button>

            {#if $files.some((f) => f.status === "complete")}
                <button class="btn-secondary" on:click={downloadAll}>
                    Download All
                </button>
            {/if}
        </div>

        <div class="create-action">
            <button class="btn-create" on:click={createNewPreset}>
                <svg
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                >
                    <path
                        d="M12 5v14M5 12h14"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
                Create New Preset
            </button>
        </div>
    </div>

    {#if error}
        <div class="error-banner">
            {error}
        </div>
    {/if}
</div>

<style>
    .action-panel {
        margin-top: 2rem;
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .controls {
        display: flex;
        align-items: center;
        justify-content: space-between;
        flex-wrap: wrap;
        gap: 1.5rem;
    }

    .control-group {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    label {
        color: var(--text-secondary);
        font-weight: 500;
    }

    .select-wrapper {
        position: relative;
    }

    select {
        appearance: none;
        background: rgba(0, 0, 0, 0.2);
        border: 1px solid var(--glass-border);
        border-radius: 8px;
        padding: 0.75rem 2.5rem 0.75rem 1rem;
        color: var(--text-primary);
        font-size: 1rem;
        cursor: pointer;
        min-width: 200px;
        transition: all 0.2s ease;
    }

    select:hover {
        background: rgba(0, 0, 0, 0.3);
        border-color: rgba(255, 255, 255, 0.2);
    }

    select:focus {
        outline: none;
        border-color: var(--color-primary);
        box-shadow: 0 0 0 2px rgba(100, 108, 255, 0.2);
    }

    /* Custom arrow for select */
    .select-wrapper::after {
        content: "";
        position: absolute;
        right: 1rem;
        top: 50%;
        transform: translateY(-50%);
        width: 0;
        height: 0;
        border-left: 5px solid transparent;
        border-right: 5px solid transparent;
        border-top: 5px solid var(--text-secondary);
        pointer-events: none;
    }

    .actions {
        display: flex;
        gap: 1rem;
    }

    .btn-primary {
        background: linear-gradient(
            135deg,
            var(--color-primary) 0%,
            var(--color-secondary) 100%
        );
        border: none;
        border-radius: 8px;
        padding: 0.75rem 2rem;
        color: white;
        font-weight: 600;
        cursor: pointer;
        transition:
            transform 0.2s,
            opacity 0.2s;
    }

    .btn-primary:hover:not(:disabled) {
        transform: translateY(-2px);
        opacity: 0.9;
    }

    .btn-primary:disabled {
        opacity: 0.5;
        cursor: not-allowed;
        transform: none;
    }

    .btn-secondary {
        background: rgba(255, 255, 255, 0.1);
        border: 1px solid var(--glass-border);
        border-radius: 8px;
        padding: 0.75rem 1.5rem;
        color: var(--text-primary);
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-secondary:hover {
        background: rgba(255, 255, 255, 0.15);
    }

    .error-banner {
        background: rgba(255, 59, 48, 0.1);
        border: 1px solid rgba(255, 59, 48, 0.3);
        color: #ff3b30;
        padding: 1rem;
        border-radius: 8px;
        font-size: 0.9rem;
    }

    @media (max-width: 600px) {
        .controls {
            flex-direction: column;
            align-items: stretch;
        }

        .control-group {
            flex-direction: column;
            align-items: stretch;
        }

        .actions {
            flex-direction: column;
        }

        .btn-primary,
        .btn-secondary,
        .btn-create {
            width: 100%;
        }
    }

    .create-action {
        width: 100%;
        border-top: 1px solid rgba(255, 255, 255, 0.1);
        padding-top: 1.5rem;
    }

    .btn-create {
        width: 100%;
        background: rgba(255, 255, 255, 0.05);
        border: 1px dashed rgba(255, 255, 255, 0.3);
        border-radius: 8px;
        padding: 1rem;
        color: var(--text-primary);
        font-weight: 500;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        transition: all 0.2s;
    }

    .btn-create:hover {
        background: rgba(255, 255, 255, 0.1);
        border-color: var(--color-primary);
        color: var(--color-primary);
    }
</style>
