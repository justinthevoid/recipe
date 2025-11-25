<script>
    import { previewFile } from "../stores";
    import {
        extractPresetParameters,
        groupParameters,
        formatParameterValue,
    } from "../parameter-extractor";
    import { recipeToCSSFilters } from "../preview-logic";
    import { detectFormatFromExtension } from "../format-detector";

    let isOpen = false;
    let file = null;
    let parameters = null;
    let groupedParams = null;
    let cssFilter = "none";
    let loading = false;
    let error = null;

    // Subscribe to store
    previewFile.subscribe(async (f) => {
        file = f;
        isOpen = !!f;
        if (f) {
            await loadPreview(f);
        } else {
            reset();
        }
    });

    function close() {
        previewFile.set(null);
    }

    function reset() {
        parameters = null;
        groupedParams = null;
        cssFilter = "none";
        error = null;
    }

    async function loadPreview(f) {
        loading = true;
        error = null;
        try {
            // Read file
            const buffer = await f.arrayBuffer();
            const bytes = new Uint8Array(buffer);
            const format = detectFormatFromExtension(f.name);

            // Extract parameters
            parameters = await extractPresetParameters(bytes, format);
            groupedParams = groupParameters(parameters);

            // Generate CSS filter
            cssFilter = recipeToCSSFilters(parameters);
        } catch (e) {
            console.error(e);
            error = e.message;
        } finally {
            loading = false;
        }
    }
</script>

{#if isOpen}
    <div
        class="modal-overlay active"
        on:click={close}
        role="button"
        tabindex="0"
        on:keydown={(e) => e.key === "Escape" && close()}
    >
        <div
            class="modal-content glass-card preview-modal"
            on:click|stopPropagation
            role="document"
            tabindex="0"
        >
            <button class="modal-close" on:click={close}>&times;</button>

            <div class="preview-layout">
                <!-- Left: Image Preview -->
                <div class="preview-image-section">
                    <h3>Instant Preview</h3>
                    <div class="image-container">
                        <img
                            src="/images/portrait-original.jpg"
                            alt="Preview"
                            style="filter: {cssFilter};"
                            class="preview-img"
                        />
                        {#if loading}
                            <div class="loader-overlay">Loading...</div>
                        {/if}
                    </div>
                    <p class="disclaimer">
                        * Approximation using CSS filters. Actual conversion may
                        vary.
                    </p>
                </div>

                <!-- Right: Parameters -->
                <div class="preview-params-section">
                    <h3>Parameters</h3>
                    {#if error}
                        <div class="error-msg">{error}</div>
                    {:else if groupedParams}
                        <div class="params-list">
                            {#each Object.entries(groupedParams) as [category, params]}
                                <div class="param-group">
                                    <h4>{category}</h4>
                                    {#each Object.entries(params) as [name, value]}
                                        <div class="param-row">
                                            <span class="param-name"
                                                >{name}</span
                                            >
                                            <span class="param-value"
                                                >{formatParameterValue(
                                                    value,
                                                )}</span
                                            >
                                        </div>
                                    {/each}
                                </div>
                            {/each}
                        </div>
                    {:else if !loading}
                        <p>No parameters found.</p>
                    {/if}
                </div>
            </div>
        </div>
    </div>
{/if}

<style>
    .preview-modal {
        max-width: 900px;
        width: 90%;
        max-height: 90vh;
        display: flex;
        flex-direction: column;
        padding: 0;
        overflow: hidden;
    }

    .preview-layout {
        display: grid;
        grid-template-columns: 1fr 300px;
        height: 100%;
        overflow: hidden;
    }

    .preview-image-section {
        padding: 2rem;
        background: rgba(0, 0, 0, 0.2);
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-right: 1px solid var(--glass-border);
    }

    .image-container {
        position: relative;
        max-width: 100%;
        border-radius: 8px;
        overflow: hidden;
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
    }

    .preview-img {
        display: block;
        max-width: 100%;
        height: auto;
        transition: filter 0.3s ease;
    }

    .preview-params-section {
        padding: 2rem;
        overflow-y: auto;
        background: rgba(255, 255, 255, 0.02);
    }

    .param-group {
        margin-bottom: 1.5rem;
    }

    .param-group h4 {
        color: var(--color-primary);
        margin-bottom: 0.5rem;
        font-size: 0.9rem;
        text-transform: uppercase;
        letter-spacing: 1px;
    }

    .param-row {
        display: flex;
        justify-content: space-between;
        padding: 0.25rem 0;
        border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        font-size: 0.9rem;
    }

    .param-name {
        color: var(--text-secondary);
    }

    .param-value {
        font-family: monospace;
        color: var(--text-primary);
    }

    .disclaimer {
        margin-top: 1rem;
        font-size: 0.8rem;
        color: var(--text-muted);
        font-style: italic;
    }

    .modal-close {
        top: 1rem;
        right: 1rem;
        z-index: 10;
    }

    @media (max-width: 768px) {
        .preview-layout {
            grid-template-columns: 1fr;
            overflow-y: auto;
        }
        .preview-image-section {
            border-right: none;
            border-bottom: 1px solid var(--glass-border);
        }
    }
</style>
