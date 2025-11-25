<script>
    import { wasmState, addFile } from "../stores";

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
            // Basic validation (can be expanded)
            const ext = file.name.split(".").pop().toLowerCase();
            if (["np3", "xmp", "lrtemplate", "costyle", "dcp"].includes(ext)) {
                addFile(file);
            } else {
                // Handle invalid file (maybe toast notification?)
                console.warn("Skipping invalid file:", file.name);
            }
        });
        // Reset input
        if (fileInput) fileInput.value = "";
    }
</script>

<div
    class="glass-card"
    id="upload"
    class:drag-over={isDragging}
    on:dragover={handleDragOver}
    on:dragleave={handleDragLeave}
    on:drop={handleDrop}
    role="button"
    tabindex="0"
>
    <!-- Dropzone -->
    <div id="dropzone" class="dropzone">
        <svg class="upload-icon" fill="currentColor" viewBox="0 0 24 24">
            <path
                d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96zM14 13v4h-4v-4H7l5-5 5 5h-3z"
            />
        </svg>
        <h2>Drag & Drop Presets</h2>
        <p>Supports .np3, .xmp, .lrtemplate, .costyle, .dcp</p>

        <button
            id="browse-button"
            class="btn btn-primary"
            on:click={() => fileInput.click()}
        >
            Select Files
        </button>
        <input
            bind:this={fileInput}
            type="file"
            id="file-input"
            multiple
            hidden
            accept=".np3,.xmp,.lrtemplate,.costyle,.dcp"
            on:change={handleFileSelect}
        />

        <!-- Status Indicators -->
        <div
            id="status"
            class="status-indicator"
            class:ready={$wasmState.status === "ready"}
            class:error={$wasmState.status === "error"}
        >
            {#if $wasmState.status === "initializing"}
                Initializing Engine...
            {:else if $wasmState.status === "ready"}
                Ready to convert
            {:else if $wasmState.status === "error"}
                Error: {$wasmState.error}
            {/if}
        </div>
    </div>
</div>

<style>
    .glass-card.drag-over {
        border-color: var(--color-primary);
        background: rgba(255, 255, 255, 0.1);
    }

    .status-indicator.ready {
        color: var(--color-success);
    }

    .status-indicator.error {
        color: var(--color-error);
    }
</style>
