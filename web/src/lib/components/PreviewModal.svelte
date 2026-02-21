<script>
    import { previewFile, currentRecipe } from "../stores";
    import { extractFullRecipe, generatePreset } from "../converter";
    import { detectFormatFromExtension } from "../format-detector";
    import { calculateColorMatrix, calculateTransferTable } from "../svg-logic";
    import { analyzeImage, calculateAutoExposure } from "../image-analysis";
    import Histogram from "./Histogram.svelte";
    import ColorBlender from "./editor/ColorBlender.svelte";
    import ColorGrading from "./editor/ColorGrading.svelte";
    import ToneCurve from "./editor/ToneCurve.svelte";

    let isOpen = false;
    let file = null;
    let loading = false;
    let error = null;
    let isSaving = false;

    // Editor Tabs
    let activeTab = "basic"; // basic, color, grading, curves

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

    // React to recipe changes
    currentRecipe.subscribe((recipe) => {
        if (recipe) {
            updateFilters(recipe);
        }
    });

    function close() {
        previewFile.set(null);
        currentRecipe.set(null);
    }

    function reset() {
        loading = false;
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
            if (f.isNew) {
                const defaultRecipe = {
                    name: "New Preset",
                    exposure: 0,
                    contrast: 0,
                    highlights: 0,
                    shadows: 0,
                    whites: 0,
                    blacks: 0,
                    clarity: 0,
                    dehaze: 0,
                    vibrance: 0,
                    saturation: 0,
                    temperature: 0,
                    tint: 0,
                    colorGrading: {
                        highlights: { hue: 0, chroma: 0, brightness: 0 },
                        midtone: { hue: 0, chroma: 0, brightness: 0 },
                        shadows: { hue: 0, chroma: 0, brightness: 0 },
                        blending: 50,
                        balance: 0,
                    },
                    toneCurveHighlights: 0,
                    toneCurveLights: 0,
                    toneCurveDarks: 0,
                    toneCurveShadows: 0,
                };
                currentRecipe.set(defaultRecipe);
                loading = false;
                return;
            }

            // Read file
            const buffer = await f.arrayBuffer();
            const bytes = new Uint8Array(buffer);
            const format = detectFormatFromExtension(f.name);

            // Extract full recipe
            const recipe = await extractFullRecipe(bytes, format);
            console.log("Loaded Recipe:", recipe);
            currentRecipe.set(recipe);
        } catch (e) {
            console.error(e);
            error = e.message;
        } finally {
            loading = false;
        }
    }

    function updateFilters(recipe) {
        if (!recipe) return;
        console.log("Updating Filters with:", recipe);

        // Extract values for SVG filters
        const temp = recipe.temperature || 0;
        const tint = recipe.tint || 0;
        const saturation = recipe.saturation || 0;
        const exposure = recipe.exposure || 0;
        const contrast = recipe.contrast || 0;
        const highlights = recipe.highlights || 0;
        const shadows = recipe.shadows || 0;
        const whites = recipe.whites || 0;
        const blacks = recipe.blacks || 0;

        const curveHighlights = recipe.toneCurveHighlights || 0;
        const curveLights = recipe.toneCurveLights || 0;
        const curveDarks = recipe.toneCurveDarks || 0;
        const curveShadows = recipe.toneCurveShadows || 0;

        colorMatrix = calculateColorMatrix(temp, tint, saturation);
        transferTable = calculateTransferTable(
            exposure,
            contrast,
            highlights,
            shadows,
            whites,
            blacks,
            curveHighlights,
            curveLights,
            curveDarks,
            curveShadows,
        );

        console.log("Transfer Table:", transferTable);
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
        currentRecipe.update((r) => {
            const currentExp = r.exposure || 0;
            const newExp = currentExp + evShift;
            return { ...r, exposure: parseFloat(newExp.toFixed(2)) };
        });
    }

    async function savePreset() {
        isSaving = true;
        try {
            let recipe;
            currentRecipe.subscribe((r) => (recipe = r))();

            if (!recipe) return;

            const np3Data = await generatePreset(recipe);

            // Create blob and download
            const blob = new Blob([np3Data], {
                type: "application/octet-stream",
            });
            const url = URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url;
            a.download = (recipe.name || "preset") + ".np3";
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
        } catch (e) {
            console.error("Save failed", e);
            error = "Failed to save preset: " + e.message;
        } finally {
            isSaving = false;
        }
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
                    <h3>Editor</h3>

                    {#if $currentRecipe}
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

                        <!-- Tabs -->
                        <div class="editor-tabs">
                            <button
                                class="tab-btn {activeTab === 'basic'
                                    ? 'active'
                                    : ''}"
                                on:click={() => (activeTab = "basic")}
                            >
                                Basic
                            </button>
                            <button
                                class="tab-btn {activeTab === 'curves'
                                    ? 'active'
                                    : ''}"
                                on:click={() => (activeTab = "curves")}
                            >
                                Curves
                            </button>
                            <button
                                class="tab-btn {activeTab === 'color'
                                    ? 'active'
                                    : ''}"
                                on:click={() => (activeTab = "color")}
                            >
                                Color
                            </button>
                            <button
                                class="tab-btn {activeTab === 'grading'
                                    ? 'active'
                                    : ''}"
                                on:click={() => (activeTab = "grading")}
                            >
                                Grading
                            </button>
                        </div>

                        <!-- Tab Content -->
                        <div class="tab-content">
                            {#if activeTab === "basic"}
                                <!-- Basic Controls -->
                                <div class="control-group">
                                    <h4>Light</h4>

                                    <label>
                                        <span>Exposure</span>
                                        <input
                                            type="range"
                                            min="-5"
                                            max="5"
                                            step="0.05"
                                            bind:value={$currentRecipe.exposure}
                                        />
                                        <span class="value"
                                            >{$currentRecipe.exposure}</span
                                        >
                                    </label>

                                    <label>
                                        <span>Contrast</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={$currentRecipe.contrast}
                                        />
                                        <span class="value"
                                            >{$currentRecipe.contrast}</span
                                        >
                                    </label>

                                    <label>
                                        <span>Highlights</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={
                                                $currentRecipe.highlights
                                            }
                                        />
                                        <span class="value"
                                            >{$currentRecipe.highlights}</span
                                        >
                                    </label>

                                    <label>
                                        <span>Shadows</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={$currentRecipe.shadows}
                                        />
                                        <span class="value"
                                            >{$currentRecipe.shadows}</span
                                        >
                                    </label>

                                    <label>
                                        <span>Whites</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={$currentRecipe.whites}
                                        />
                                        <span class="value"
                                            >{$currentRecipe.whites}</span
                                        >
                                    </label>

                                    <label>
                                        <span>Blacks</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={$currentRecipe.blacks}
                                        />
                                        <span class="value"
                                            >{$currentRecipe.blacks}</span
                                        >
                                    </label>
                                </div>

                                <div class="control-group">
                                    <h4>Presence</h4>

                                    <label>
                                        <span>Clarity</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={$currentRecipe.clarity}
                                        />
                                        <span class="value"
                                            >{$currentRecipe.clarity}</span
                                        >
                                    </label>

                                    <label>
                                        <span>Dehaze</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={$currentRecipe.dehaze}
                                        />
                                        <span class="value"
                                            >{$currentRecipe.dehaze}</span
                                        >
                                    </label>

                                    <label>
                                        <span>Vibrance</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={$currentRecipe.vibrance}
                                        />
                                        <span class="value"
                                            >{$currentRecipe.vibrance}</span
                                        >
                                    </label>

                                    <label>
                                        <span>Saturation</span>
                                        <input
                                            type="range"
                                            min="-100"
                                            max="100"
                                            step="1"
                                            bind:value={
                                                $currentRecipe.saturation
                                            }
                                        />
                                        <span class="value"
                                            >{$currentRecipe.saturation}</span
                                        >
                                    </label>
                                </div>
                            {:else if activeTab === "curves"}
                                <div class="control-group">
                                    <h4>Tone Curve</h4>
                                    <ToneCurve />
                                </div>
                            {:else if activeTab === "color"}
                                <div class="control-group">
                                    <h4>Color Mixer</h4>
                                    <ColorBlender />
                                </div>
                            {:else if activeTab === "grading"}
                                <div class="control-group">
                                    <h4>Color Grading</h4>
                                    <ColorGrading />
                                </div>
                            {/if}
                        </div>

                        <!-- Spacer to prevent content from being hidden behind sticky button -->
                        <div class="spacer"></div>
                    {:else if !loading}
                        <p>No parameters found.</p>
                    {/if}
                </div>

                <!-- Sticky Save Button Area -->
                <div class="save-section">
                    <button
                        class="btn-save"
                        on:click={savePreset}
                        disabled={isSaving}
                    >
                        {isSaving ? "Saving..." : "Save as NP3"}
                    </button>
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
        grid-template-columns: 1fr 350px;
        height: 100%;
        overflow: hidden;
        position: relative;
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
        appearance: none; /* Ensure standard property is set */
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
        padding-bottom: 100px; /* Space for sticky button */
        overflow-y: auto;
        background: rgba(255, 255, 255, 0.02);
        height: 100%;
        box-sizing: border-box;
    }

    .save-section {
        position: absolute;
        bottom: 0;
        right: 0;
        width: 350px;
        padding: 1.5rem;
        background: rgba(18, 18, 24, 0.95); /* Match modal bg */
        border-top: 1px solid rgba(255, 255, 255, 0.1);
        backdrop-filter: blur(10px);
        z-index: 20;
        box-sizing: border-box;
    }

    .spacer {
        height: 100px;
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

    h3 {
        margin-bottom: 1.5rem;
        font-size: 1.2rem;
        color: white;
    }

    h4 {
        font-size: 0.9rem;
        color: var(--color-primary);
        margin-bottom: 1rem;
        text-transform: uppercase;
        letter-spacing: 1px;
    }

    .control-group label {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        margin-bottom: 1rem;
    }

    .control-group span {
        font-size: 0.9rem;
        color: var(--text-secondary);
        display: flex;
        justify-content: space-between;
    }

    .control-group .value {
        color: var(--text-primary);
        font-family: monospace;
    }

    /* Custom Scrollbar for Params */
    .preview-params-section::-webkit-scrollbar {
        width: 6px; /* Thinner */
    }

    .preview-params-section::-webkit-scrollbar-track {
        background: rgba(0, 0, 0, 0.1);
    }

    .preview-params-section::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.2);
        border-radius: 4px;
    }

    .preview-params-section::-webkit-scrollbar-thumb:hover {
        background: rgba(255, 255, 255, 0.3);
    }

    /* Range Input Styling */
    input[type="range"] {
        width: 100%;
        background: transparent;
        -webkit-appearance: none;
        appearance: none;
        margin: 0.5rem 0;
    }

    input[type="range"]:focus {
        outline: none;
    }

    input[type="range"]::-webkit-slider-runnable-track {
        width: 100%;
        height: 6px; /* Thicker track */
        background: rgba(255, 255, 255, 0.15); /* More visible */
        border-radius: 3px;
        cursor: pointer;
        transition: background 0.2s;
    }

    input[type="range"]:hover::-webkit-slider-runnable-track {
        background: rgba(255, 255, 255, 0.25);
    }

    input[type="range"]::-webkit-slider-thumb {
        -webkit-appearance: none;
        appearance: none;
        height: 20px; /* Larger thumb */
        width: 20px;
        border-radius: 50%;
        background: var(--color-primary);
        border: 2px solid white; /* Add border for visibility */
        margin-top: -7px; /* Center thumb on track */
        cursor: pointer;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.4);
        transition:
            transform 0.1s,
            background 0.2s;
    }

    input[type="range"]::-webkit-slider-thumb:hover {
        transform: scale(1.1);
        background: #fff;
    }

    /* Firefox styles */
    input[type="range"]::-moz-range-track {
        width: 100%;
        height: 4px;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 2px;
        cursor: pointer;
    }

    input[type="range"]::-moz-range-thumb {
        height: 16px;
        width: 16px;
        border: none;
        border-radius: 50%;
        background: var(--color-primary);
        cursor: pointer;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
        transition:
            transform 0.1s,
            background 0.2s;
    }

    input[type="range"]::-moz-range-thumb:hover {
        transform: scale(1.1);
        background: #fff;
    }

    .btn-save {
        width: 100%;
        padding: 1rem;
        background: var(--color-primary);
        color: white;
        border: none;
        border-radius: 8px;
        font-weight: 700; /* Bolder */
        font-size: 1.1rem; /* Larger */
        cursor: pointer;
        transition: all 0.2s;
        margin-top: 0; /* Reset margin */
        box-shadow: 0 4px 12px rgba(var(--color-primary-rgb), 0.3); /* Glow */
    }

    .btn-save:hover {
        filter: brightness(1.1);
        transform: translateY(-1px);
    }

    .btn-save:disabled {
        opacity: 0.5;
        cursor: not-allowed;
        transform: none;
    }

    .editor-tabs {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 2rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        padding-bottom: 0.5rem;
    }

    .tab-btn {
        background: transparent;
        border: none;
        color: var(--text-secondary);
        padding: 0.5rem 1rem;
        cursor: pointer;
        font-size: 0.9rem;
        border-radius: 4px;
        transition: all 0.2s;
    }

    .tab-btn:hover {
        color: white;
        background: rgba(255, 255, 255, 0.05);
    }

    .tab-btn.active {
        color: var(--color-primary);
        background: rgba(255, 255, 255, 0.1);
        font-weight: 600;
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
