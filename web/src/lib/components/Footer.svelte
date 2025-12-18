<script>
    export let onLegalClick;
    export let version = "1.0.0";
</script>

<footer class="app-footer">
    <div class="container footer-content">
        <div class="footer-brand">
            <h4>Recipe</h4>
            <p style="color: var(--text-muted);">Pro Preset Converter</p>
        </div>
        <div class="footer-links">
            <a href="https://github.com/justin/recipe" target="_blank">GitHub</a
            >
            <a href="docs/privacy.md">Privacy</a>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a href="#" on:click|preventDefault={onLegalClick}
                >Legal Disclaimer</a
            >
        </div>
        <div style="color: var(--text-muted);">
            v<span>{version}</span>
        </div>
    </div>
</footer>
