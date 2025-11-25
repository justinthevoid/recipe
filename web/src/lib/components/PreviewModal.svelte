<script>
    import { previewFile } from "../stores";
                        />

                        <!-- Labels -->
                        <div
                            class="label label-before"
                            style="opacity: {sliderPosition < 10 ? 0 : 1}"
                        >
                            Before
                        </div>
                        <div
                            class="label label-after"
                            style="opacity: {sliderPosition > 90 ? 0 : 1}"
                        >
                            After
                        </div>

                        {#if loading}
                            <div class="loader-overlay">Loading...</div>
                        {/if}
                    </div>

                    <p class="disclaimer">
                        * Approximation using SVG filters. Actual conversion may
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
        position: relative;
    }

    .image-container {
        position: relative;
        max-width: 100%;
        border-radius: 8px;
        overflow: hidden;
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
        user-select: none;
    }

    .preview-img {
        display: block;
        max-width: 100%;
        height: auto;
        pointer-events: none;
    }

    .preview-img.after {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
    }

    /* Slider Controls */
    .slider-input {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        opacity: 0;
        cursor: col-resize;
        z-index: 20;
        margin: 0;
    }

    .slider-handle {
        position: absolute;
        top: 0;
        bottom: 0;
        width: 2px;
        background: rgba(255, 255, 255, 0.8);
        pointer-events: none; /* Let clicks pass to input */
        z-index: 10;
        box-shadow: 0 0 5px rgba(0, 0, 0, 0.5);
    }

    .slider-button {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 32px;
        height: 32px;
        background: white;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #333;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
    }

    .label {
        position: absolute;
        top: 1rem;
        background: rgba(0, 0, 0, 0.6);
        color: white;
        padding: 0.25rem 0.75rem;
        border-radius: 4px;
        font-size: 0.8rem;
        pointer-events: none;
        transition: opacity 0.2s;
        z-index: 5;
    }

    .label-before {
        left: 1rem;
    }

    .label-after {
        right: 1rem;
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
    }

    .param-name {
        color: var(--text-secondary);
        font-size: 0.9rem;
    }

    .param-value {
        color: var(--text-primary);
        font-family: monospace;
        font-size: 0.9rem;
    }

    .error-msg {
        color: #ff3b30;
        padding: 1rem;
        background: rgba(255, 59, 48, 0.1);
        border-radius: 8px;
    }

    .loader-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        z-index: 30;
    }

    .disclaimer {
        margin-top: 1rem;
        font-size: 0.8rem;
        color: var(--text-secondary);
        font-style: italic;
    }
</style>
