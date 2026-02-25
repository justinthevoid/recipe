<script>
import { currentRecipe } from "../../stores";

// Parametric Curve sliders
// ToneCurveShadows, ToneCurveDarks, ToneCurveLights, ToneCurveHighlights
// Range: -100 to +100
</script>

<div class="tone-curve">
    {#if $currentRecipe}
        <div class="slider-group">
            <h4>Parametric Curve</h4>

            <label>
                <span>Highlights</span>
                <input
                    type="range"
                    min="-100"
                    max="100"
                    step="1"
                    bind:value={$currentRecipe.ToneCurveHighlights}
                />
                <span class="value">{$currentRecipe.ToneCurveHighlights}</span>
            </label>

            <label>
                <span>Lights</span>
                <input
                    type="range"
                    min="-100"
                    max="100"
                    step="1"
                    bind:value={$currentRecipe.ToneCurveLights}
                />
                <span class="value">{$currentRecipe.ToneCurveLights}</span>
            </label>

            <label>
                <span>Darks</span>
                <input
                    type="range"
                    min="-100"
                    max="100"
                    step="1"
                    bind:value={$currentRecipe.ToneCurveDarks}
                />
                <span class="value">{$currentRecipe.ToneCurveDarks}</span>
            </label>

            <label>
                <span>Shadows</span>
                <input
                    type="range"
                    min="-100"
                    max="100"
                    step="1"
                    bind:value={$currentRecipe.ToneCurveShadows}
                />
                <span class="value">{$currentRecipe.ToneCurveShadows}</span>
            </label>
        </div>
    {/if}
</div>

<style>
    .tone-curve {
        background: rgba(0, 0, 0, 0.2);
        padding: 1rem;
        border-radius: 8px;
        margin-bottom: 2rem;
    }

    .slider-group {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    h4 {
        font-size: 0.9rem;
        color: var(--color-primary);
        margin-bottom: 1rem;
        text-transform: uppercase;
        letter-spacing: 1px;
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
