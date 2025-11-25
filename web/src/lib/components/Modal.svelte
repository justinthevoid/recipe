<script>
    export let isOpen = false;
    export let onClose;
</script>

{#if isOpen}
    <div
        class="modal-overlay active"
        on:click={onClose}
        role="dialog"
        aria-modal="true"
    >
        <div
            class="modal-content glass-card"
            on:click|stopPropagation
            role="document"
        >
            <button class="modal-close" on:click={onClose}>&times;</button>
            <slot></slot>
        </div>
    </div>
{/if}
