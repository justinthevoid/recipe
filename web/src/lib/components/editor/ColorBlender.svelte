<script>
    import { currentRecipe } from "../../stores";

    const colors = [
        { id: "Red", label: "Red", color: "#ff4d4d" },
        { id: "Orange", label: "Orange", color: "#ff9f4d" },
        { id: "Yellow", label: "Yellow", color: "#ffd700" },
        { id: "Green", label: "Green", color: "#4dff4d" },
        { id: "Aqua", label: "Aqua", color: "#4dffff" },
        { id: "Blue", label: "Blue", color: "#4d4dff" },
        { id: "Purple", label: "Purple", color: "#9f4dff" },
        { id: "Magenta", label: "Magenta", color: "#ff4dff" },
    ];

    let activeColor = "Red";

    function ensureColorStruct(colorId) {
        if (!$currentRecipe[colorId]) {
            $currentRecipe[colorId] = { Hue: 0, Saturation: 0, Luminance: 0 };
        }
    }
</script>

<div class="color-blender">
    <div class="color-tabs">
        {#each colors as color}
            <button
                class="color-tab {activeColor === color.id ? 'active' : ''}"
                style="--tab-color: {color.color}"
                on:click={() => (activeColor = color.id)}
                title={color.label}
            >
                <div class="color-dot"></div>
            </button>
        {/each}
    </div>

    <div class="sliders">
        {#if $currentRecipe}
            {#each colors as color}
                {#if activeColor === color.id}
                    <!-- Ensure struct exists on render -->
                    {(ensureColorStruct(color.id), "")}

                    <div class="slider-group">
                        <label>
                            <span>Hue</span>
                            <input
                                type="range"
                                min="-180"
                                max="180"
                                step="1"
                                bind:value={$currentRecipe[color.id].Hue}
                            />
                            <span class="value"
                                >{$currentRecipe[color.id].Hue}</span
                            >
                        </label>

                        <label>
                            <span>Saturation</span>
                            <input
                                type="range"
                                min="-100"
                                max="100"
                                step="1"
                                bind:value={$currentRecipe[color.id].Saturation}
                            />
                            <span class="value"
                                >{$currentRecipe[color.id].Saturation}</span
                            >
                        </label>

                        <label>
                            <span>Luminance</span>
                            <input
                                type="range"
                                min="-100"
                                max="100"
                                step="1"
                                bind:value={$currentRecipe[color.id].Luminance}
                            />
                            <span class="value"
                                >{$currentRecipe[color.id].Luminance}</span
                            >
                        </label>
                    </div>
                {/if}
            {/each}
        {/if}
    </div>
</div>

<style>
    .color-blender {
        background: rgba(0, 0, 0, 0.2);
        padding: 1rem;
        border-radius: 8px;
        margin-bottom: 2rem;
    }

    .color-tabs {
        display: flex;
        justify-content: space-between;
        margin-bottom: 1.5rem;
        background: rgba(0, 0, 0, 0.3);
        padding: 0.5rem;
        border-radius: 20px;
    }

    .color-tab {
        width: 24px;
        height: 24px;
        border-radius: 50%;
        border: 2px solid transparent;
        background: transparent;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s;
        padding: 0;
    }

    .color-dot {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        background-color: var(--tab-color);
        transition: all 0.2s;
    }

    .color-tab:hover .color-dot {
        transform: scale(1.2);
    }

    .color-tab.active {
        border-color: var(--tab-color);
        background: rgba(255, 255, 255, 0.1);
    }

    .color-tab.active .color-dot {
        transform: scale(0.8);
    }

    .slider-group {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    label {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    span {
        font-size: 0.8rem;
        color: var(--text-secondary);
        display: flex;
        justify-content: space-between;
    }

    .value {
        color: var(--text-primary);
        font-family: monospace;
    }

    input[type="range"] {
        width: 100%;
        background: transparent;
        -webkit-appearance: none;
        appearance: none;
    }

    input[type="range"]::-webkit-slider-runnable-track {
        width: 100%;
        height: 4px;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 2px;
    }

    input[type="range"]::-webkit-slider-thumb {
        -webkit-appearance: none;
        height: 16px;
        width: 16px;
        border-radius: 50%;
        background: var(--color-primary);
        margin-top: -6px;
        cursor: pointer;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    }
</style>
