<script>
    import { removeFile } from "../stores";
    import { detectFormatFromExtension } from "../format-detector";

    export let file;

    // Format detection
    $: format = detectFormatFromExtension(file.name);

    // Size formatting
    function formatSize(bytes) {
        if (bytes === 0) return "0 B";
        const k = 1024;
        const sizes = ["B", "KB", "MB", "GB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    }
</script>

<div class="file-card">
    <div class="file-info">
        <div class="file-details">
            <div class="file-name">{file.name}</div>
            <div class="file-meta">{formatSize(file.size)}</div>
        </div>
    </div>
    <div
        class="file-actions"
        style="display: flex; align-items: center; gap: 1rem;"
    >
        <span class="format-badge">{format.toUpperCase()}</span>
        <button
            class="btn-icon remove-btn"
            on:click={() => removeFile(file.id)}
            aria-label="Remove file"
        >
            <svg
                width="20"
                height="20"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                ></path>
            </svg>
        </button>
    </div>
</div>
