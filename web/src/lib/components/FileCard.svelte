<script>
import { removeFile, previewFile } from "../stores";
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

<div class="file-card {file.status}">
    <div class="file-info">
        <div class="file-details">
            <div class="file-name">{file.name}</div>
            <div class="file-meta">
                {formatSize(file.size)}
                {#if file.status === "processing"}
                    <span class="status-text processing">Converting...</span>
                {:else if file.status === "complete"}
                    <span class="status-text success">Converted</span>
                {:else if file.status === "error"}
                    <span class="status-text error"
                        >{file.error || "Error"}</span
                    >
                {/if}
            </div>
        </div>
    </div>
    <div
        class="file-actions"
        style="display: flex; align-items: center; gap: 1rem;"
    >
        <span class="format-badge">{format.toUpperCase()}</span>

        {#if file.status === "complete"}
            <a
                href={file.outputUrl}
                download={file.outputName}
                class="btn-icon download-btn"
                title="Download"
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
                        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                    ></path>
                </svg>
            </a>
        {/if}

        <!-- Preview Button -->
        <button
            class="btn-icon"
            on:click={() => previewFile.set(file.file)}
            aria-label="Preview file"
            title="Preview"
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
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                />
            </svg>
        </button>

        <button
            class="btn-icon remove-btn"
            on:click={() => removeFile(file.id)}
            aria-label="Remove file"
            title="Remove"
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

<style>
    .status-text {
        margin-left: 0.5rem;
        font-size: 0.8rem;
        font-weight: 500;
    }
    .status-text.processing {
        color: var(--color-secondary);
    }
    .status-text.success {
        color: #4cd964;
    }
    .status-text.error {
        color: #ff3b30;
    }

    .download-btn {
        color: #4cd964;
    }
    .download-btn:hover {
        background: rgba(76, 217, 100, 0.1);
    }
</style>
