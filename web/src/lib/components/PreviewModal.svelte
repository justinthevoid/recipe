<script>
    import { previewFile } from "../stores";
    import {
        extractPresetParameters,
        groupParameters,
        formatParameterValue,
    } from "../parameter-extractor";
    import { detectFormatFromExtension } from "../format-detector";
    import { calculateColorMatrix, calculateTransferTable } from "../svg-logic";
    import { analyzeImage, calculateAutoExposure } from "../image-analysis";
    import Histogram from "./Histogram.svelte";

    let isOpen = false;
    let file = null;
    let parameters = null;
    let groupedParams = null;
    let loading = false;
    let error = null;

    // Slider state
    let sliderPosition = 50;

    // Filter values
    let colorMatrix = "";
    let transferTable = "";

    // Image State
    const defaultImages = [
        "/images/portrait-original.jpg",
        "/images/landscape-original.jpg",
        "/images/product-original.jpg",
    ];
    let currentImageIndex = 0;
    let customImage = null;
    let fileInput;
    let imgElement; // Reference to the "Before" image for analysis

    // Analysis State
    let histogramData = { r: [], g: [], b: [] };
    let imageStats = null;
    let isAnalyzing = false;

    $: currentImageSrc = customImage || defaultImages[currentImageIndex];

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
        parameters = {};
        groupedParams = null;
        error = null;
        sliderPosition = 50;
        colorMatrix = "";
        transferTable = "";
        histogramData = { r: [], g: [], b: [] };
        // Don't reset custom image or index so user preference persists during session
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

            // Calculate filters
            updateFilters();
        } catch (e) {
            console.error(e);
            error = e.message;
        } finally {
            loading = false;
        }
    }

    function getVal(key, alias) {
        if (!parameters) return 0;
        if (parameters[key] != null && typeof parameters[key] === "number")
            return parameters[key];
        if (
            alias &&
            parameters[alias] != null &&
            typeof parameters[alias] === "number"
        )
            return parameters[alias];
        return 0;
    }

    function setVal(key, value) {
        if (!parameters) return;
        // Update parameter if it exists, or add it
        // We need to handle aliases (Exposure vs Exposure2012)
        if (parameters["Exposure2012"] !== undefined)
            parameters["Exposure2012"] = value;
        else if (parameters["Exposure"] !== undefined)
            parameters["Exposure"] = value;
        else parameters[key] = value; // Default to key if neither exists

        // Re-group and update
        groupedParams = groupParameters(parameters);
        updateFilters();
    }

    function updateFilters() {
        if (!parameters) return;

        const temp = getVal("Temperature");
        const tint = getVal("Tint");
        const saturation = getVal("Saturation");
        const exposure = getVal("Exposure", "Exposure2012");
        const contrast = getVal("Contrast", "Contrast2012");

        colorMatrix = calculateColorMatrix(temp, tint, saturation);
        transferTable = calculateTransferTable(exposure, contrast);
    }

    // Image Analysis Trigger
    async function handleImageLoad() {
        if (!imgElement) return;
        isAnalyzing = true;
        try {
            const result = await analyzeImage(imgElement);
            histogramData = result.histogram;
            imageStats = result;
        } catch (e) {
            console.error("Analysis failed", e);
        } finally {
            isAnalyzing = false;
        }
    }

    function autoTone() {
        if (!histogramData || !imageStats) return;

        const evShift = calculateAutoExposure(
            histogramData.l,
            imageStats.totalPixels,
        );

        // Apply shift to current exposure
        const currentExp = getVal("Exposure", "Exposure2012");
        const newExp = currentExp + evShift;

        setVal("Exposure", parseFloat(newExp.toFixed(2)));
    }

    // Image Navigation
    function nextImage() {
        customImage = null; // Clear custom image when navigating
        currentImageIndex = (currentImageIndex + 1) % defaultImages.length;
    }

    function prevImage() {
        customImage = null;
        currentImageIndex =
            (currentImageIndex - 1 + defaultImages.length) %
            defaultImages.length;
    }

    function triggerUpload() {
        fileInput.click();
    }

    function handleImageUpload(e) {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                customImage = e.target.result;
            };
            reader.readAsDataURL(file);
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

            <!-- Inline SVG Filters -->
            <svg
                style="position: absolute; width: 0; height: 0; overflow: hidden;"
                aria-hidden="true"
            >
                <defs>
                    <filter
                        id="preview-filter"
                        color-interpolation-filters="sRGB"
                    >
                        <feColorMatrix
                            type="matrix"
                            values={colorMatrix}
                            result="colored"
                        />
                        <feComponentTransfer in="colored" result="final">
                            <feFuncR type="table" tableValues={transferTable} />
                            <feFuncG type="table" tableValues={transferTable} />
                            <feFuncB type="table" tableValues={transferTable} />
                        </feComponentTransfer>
                    </filter>
                </defs>
            </svg>

            <div class="preview-layout">
                <!-- Left: Image Preview -->
                <div class="preview-image-section">
                    <h3>Instant Preview</h3>

                    <div class="image-container">
                        <!-- Before Image (Background) -->
                        <img
                            bind:this={imgElement}
                            src={currentImageSrc}
                            alt="Before"
                            class="preview-img before"
                            on:load={handleImageLoad}
                            crossorigin="anonymous"
                        />

                        <!-- After Image (Foreground, Clipped) -->
                        <img
                            src={currentImageSrc}
                            alt="After"
                            style="filter: url(#preview-filter); clip-path: inset(0 0 0 {sliderPosition}%); -webkit-filter: url(#preview-filter);"
                            class="preview-img after"
                        />

                        <!-- Slider Handle -->
                        <div
                            class="slider-handle"
                            style="left: {sliderPosition}%"
                        >
                            <div class="slider-line"></div>
                            <div class="slider-button">
                                <svg
                                    width="16"
                                    height="16"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path
                                        d="M18 8L22 12L18 16"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                    <path
                                        d="M6 8L2 12L6 16"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                </svg>
                            </div>
                        </div>

                        <!-- Range Input (Invisible Control) -->
                        <input
                            type="range"
                            min="0"
                            max="100"
                            bind:value={sliderPosition}
                            class="slider-input"
                            aria-label="Comparison slider"
                        />

                        <!-- Navigation Controls -->
                        <button
                            class="nav-btn prev"
                            on:click|stopPropagation={prevImage}
                            title="Previous Image"
                        >
                            <svg
                                width="24"
                                height="24"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <path
                                    d="M15 18l-6-6 6-6"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                            </svg>
                        </button>
                        <button
                            class="nav-btn next"
                            on:click|stopPropagation={nextImage}
                            title="Next Image"
                        >
                            <svg
                                width="24"
                                height="24"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <path
                                    d="M9 18l6-6-6-6"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                            </svg>
                        </button>

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

                        {#if loading || isAnalyzing}
                            <div class="loader-overlay">
                                {loading
                                    ? "Loading Preset..."
                                    : "Analyzing Image..."}
                            </div>
                        {/if}
                    </div>

                    <div class="image-toolbar">
                        <button class="btn-text" on:click={triggerUpload}>
                            <svg
                                width="16"
                                height="16"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <path
                                    d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                                <polyline
                                    points="17 8 12 3 7 8"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                                <line
                                    x1="12"
                                    y1="3"
                                    x2="12"
                                    y2="15"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                            </svg>
                            Upload Sample Image
                        </button>
                        <input
                            type="file"
                            accept="image/*"
                            bind:this={fileInput}
                            on:change={handleImageUpload}
                            style="display: none;"
                        />
                    </div>

                    <p class="disclaimer">
                        * Approximation using SVG filters. Actual conversion may
                        vary.
                    </p>
                </div>

                <!-- Right: Parameters -->
                <div class="preview-params-section">
                    <h3>Parameters</h3>

                    <!-- Histogram -->
                    <div class="histogram-section">
                        <Histogram data={histogramData} />
                        <button
                            class="btn-auto"
                            on:click={autoTone}
                            title="Auto-correct Exposure"
                        >
                            <svg
                                width="16"
                                height="16"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <path
                                    d="M12 2v20M2 12h20"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                                <circle cx="12" cy="12" r="4" />
                            </svg>
                            Auto Tone
                        </button>
                    </div>

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

    /* Navigation Buttons */
    .nav-btn {
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
        background: rgba(0, 0, 0, 0.5);
        border: 1px solid rgba(255, 255, 255, 0.2);
        color: white;
        width: 40px;
        height: 40px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        z-index: 25; /* Above slider input */
        transition: all 0.2s;
    }

    .nav-btn:hover {
        background: rgba(0, 0, 0, 0.8);
        transform: translateY(-50%) scale(1.1);
    }

    .nav-btn.prev {
        left: 1rem;
    }
    .nav-btn.next {
        right: 1rem;
    }

    /* Image Toolbar */
    .image-toolbar {
        margin-top: 1rem;
        display: flex;
        gap: 1rem;
    }

    .btn-text {
        background: none;
        border: none;
        color: var(--text-secondary);
        font-size: 0.9rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        transition: color 0.2s;
    }

    .btn-text:hover {
        color: var(--color-primary);
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

    .histogram-section {
        margin-bottom: 2rem;
        padding-bottom: 1rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }

    .btn-auto {
        width: 100%;
        padding: 0.5rem;
        background: rgba(255, 255, 255, 0.1);
        border: 1px solid rgba(255, 255, 255, 0.2);
        color: white;
        border-radius: 4px;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        font-size: 0.9rem;
        transition: all 0.2s;
    }

    .btn-auto:hover {
        background: rgba(255, 255, 255, 0.2);
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
        margin-top: 0.5rem;
        font-size: 0.8rem;
        color: var(--text-secondary);
        font-style: italic;
    }
</style>
