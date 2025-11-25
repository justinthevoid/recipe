<script>
    import { files, settings, updateFileStatus } from "../stores";
    import { convertFile } from "../converter";
    import { detectFormatFromExtension } from "../format-detector";

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
                    const sourceFormat = detectFormatFromExtension(
                        fileData.name,
                    );

                    const outputData = await convertFile(
                        bytes,
                        sourceFormat,
                        targetFormat,
                        fileData.name,
                    );

                    // Create download URL
                    const blob = new Blob([outputData], {
                        type: "application/octet-stream",
                    });
                    const url = URL.createObjectURL(blob);

                    // Generate new filename
                    const newName =
                        fileData.name.replace(/\.[^/.]+$/, "") +
                        "." +
                        targetFormat;

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
        .btn-secondary {
            width: 100%;
        }
    }
</style>
