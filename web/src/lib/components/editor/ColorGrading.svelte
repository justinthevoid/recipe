<script>
    import { currentRecipe } from "../../stores";

    // Ensure ColorGrading struct exists
    $: if ($currentRecipe && !$currentRecipe.ColorGrading) {
        $currentRecipe.ColorGrading = {
            Highlights: { Hue: 0, Chroma: 0, Brightness: 0 },
            Midtone: { Hue: 0, Chroma: 0, Brightness: 0 },
            Shadows: { Hue: 0, Chroma: 0, Brightness: 0 },
            Blending: 0,
            Balance: 0,
        };
    }

    const zones = [
        { id: "Shadows", label: "Shadows" },
        { id: "Midtone", label: "Midtones" },
        { id: "Highlights", label: "Highlights" },
    ];

    let activeZone = "Midtone";
</script>

<div class="color-grading">
    <div class="zone-tabs">
        {#each zones as zone}
            <button
                class="zone-tab {activeZone === zone.id ? 'active' : ''}"
                on:click={() => (activeZone = zone.id)}
            >
                {zone.label}
            </button>
        {/each}
    </div>

    {#if $currentRecipe && $currentRecipe.ColorGrading}
        <div class="wheel-controls">
            <!-- We'll use sliders for now, implementing a visual wheel is complex and can be done later if needed -->
            <div class="slider-group">
                <label>
                    <span>Hue</span>
                    <input
                        type="range"
                        min="0"
                        max="360"
                        step="1"
                        bind:value={$currentRecipe.ColorGrading[activeZone].Hue}
                        style="background: linear-gradient(to right, red, yellow, lime, cyan, blue, magenta, red);"
                    />
                    <span class="value"
                        >{$currentRecipe.ColorGrading[activeZone].Hue}°</span
                    >
                </label>

                <label>
                    <span>Saturation</span>
                    <input
                        type="range"
                        min="0"
                        max="100"
                        step="1"
                        bind:value={
                            $currentRecipe.ColorGrading[activeZone].Chroma
                        }
                    />
                    <span class="value"
                        >{$currentRecipe.ColorGrading[activeZone].Chroma}</span
                    >
                </label>

                <label>
                    <span>Luminance</span>
                    <input
                        type="range"
                        min="-100"
                        max="100"
                        step="1"
                        bind:value={
                            $currentRecipe.ColorGrading[activeZone].Brightness
                        }
                    />
                    <span class="value"
                        >{$currentRecipe.ColorGrading[activeZone]
                            .Brightness}</span
                    >
                </label>
            </div>

            <div class="global-controls">
                <h4>Global</h4>
                <label>
                    <span>Blending</span>
                    <input
                        type="range"
                        min="0"
                        max="100"
                        step="1"
                        bind:value={$currentRecipe.ColorGrading.Blending}
                    />
                    <span class="value"
                        >{$currentRecipe.ColorGrading.Blending}</span
                    >
                </label>

                <label>
                    <span>Balance</span>
                    <input
                        type="range"
                        min="-100"
                        max="100"
                        step="1"
                        bind:value={$currentRecipe.ColorGrading.Balance}
                    />
                    <span class="value"
                        >{$currentRecipe.ColorGrading.Balance}</span
                    >
                </label>
            </div>
        </div>
    {/if}
</div>

<style>
    .color-grading {
        background: rgba(0, 0, 0, 0.2);
        padding: 1rem;
        border-radius: 8px;
        margin-bottom: 2rem;
    }

    .zone-tabs {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 1.5rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }

    .zone-tab {
        background: transparent;
        border: none;
        color: var(--text-secondary);
        padding: 0.5rem 1rem;
        cursor: pointer;
        font-size: 0.85rem;
        border-bottom: 2px solid transparent;
        transition: all 0.2s;
    }

    .zone-tab.active {
        color: var(--color-primary);
        border-bottom-color: var(--color-primary);
    }

    .slider-group {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        margin-bottom: 2rem;
    }

    .global-controls {
        border-top: 1px solid rgba(255, 255, 255, 0.1);
        padding-top: 1rem;
    }

    .global-controls h4 {
        font-size: 0.8rem;
        color: var(--text-secondary);
        margin-bottom: 1rem;
        text-transform: uppercase;
    }

    label {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        margin-bottom: 1rem;
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
