<script lang="ts">
	import { getNested } from "../utils";
	import type { ParameterDefinition, UniversalRecipe } from "../types";
	import ParameterSliderUnit from "./ParameterSliderUnit.svelte";

	let {
		parameters = [],
		recipe = {},
		originalRecipe = {},
		onchange
	}: {
		parameters: ParameterDefinition[];
		recipe: UniversalRecipe;
		originalRecipe: UniversalRecipe;
		onchange: (field: string, value: number) => void;
	} = $props();

	let selectedColorIndex = $state(0);

	const colors = [
		{ name: "Red", key: "red", hex: "#e5484d", h: 358 },
		{ name: "Orange", key: "orange", hex: "#f76b15", h: 25 },
		{ name: "Yellow", key: "yellow", hex: "#ffe629", h: 53 },
		{ name: "Green", key: "green", hex: "#46a758", h: 131 },
		{ name: "Cyan", key: "aqua", hex: "#00c2d7", h: 186 },
		{ name: "Blue", key: "blue", hex: "#0090ff", h: 206 },
		{ name: "Purple", key: "purple", hex: "#8e4ec6", h: 272 },
		{ name: "Magenta", key: "magenta", hex: "#e93d82", h: 336 }
	];

	const selectedColor = $derived(colors[selectedColorIndex]);

	function _getParam(key: string) {
		return parameters.find(p => p.key === key);
	}

	function getVal(key: string) {
		return Number(getNested(recipe, key) ?? 0);
	}

	function getOrig(key: string) {
		return Number(getNested(originalRecipe, key) ?? 0);
	}

	function _getNormalizedKey(key: string, type: string) {
		return `${key}.${type.toLowerCase()}`;
	}

	function isColorDirty(prefix: string) {
		return getVal(`${prefix}.hue`) !== getOrig(`${prefix}.hue`) ||
		       getVal(`${prefix}.saturation`) !== getOrig(`${prefix}.saturation`) ||
		       getVal(`${prefix}.luminance`) !== getOrig(`${prefix}.luminance`);
	}

	function _getHueGradient(h: number) {
		return `linear-gradient(to right, hsl(${h - 50}, 80%, 50%), hsl(${h}, 80%, 50%), hsl(${h + 50}, 80%, 50%))`;
	}

	function _getChromaGradient(h: number) {
		return `linear-gradient(to right, hsl(${h}, 0%, 50%), hsl(${h}, 80%, 50%), hsl(${h}, 100%, 50%))`;
	}

	function _getBrightnessGradient(h: number) {
		return `linear-gradient(to right, black, hsl(${h}, 80%, 50%), white)`;
	}
</script>

<div class="flex flex-col gap-4">
	<!-- Color Swatch Selection Row -->
	<div class="flex items-center gap-2 px-1">
		{#each colors as color, i}
			<button
				type="button"
				class="relative flex-shrink-0 rounded-full transition-all"
				style="background-color: {color.hex}; width: {selectedColorIndex === i ? '24px' : '20px'}; height: {selectedColorIndex === i ? '24px' : '20px'}; {selectedColorIndex === i ? 'box-shadow: 0 0 0 2px var(--color-interactive);' : 'opacity: 0.8;'}"
				onclick={() => selectedColorIndex = i}
				title={color.name}
				aria-label="Select {color.name}"
			>
				{#if isColorDirty(color.key)}
					<div class="absolute -top-1 -right-1 rounded-full" style="width: 8px; height: 8px; background: var(--color-modified); border: 1px solid var(--color-canvas-base);"></div>
				{/if}
			</button>
		{/each}
	</div>

	<!-- Controls Section -->
	<div class="p-4 bg-muted/30 rounded border border-border flex flex-col gap-5">
		<div class="flex flex-col gap-4">
			{#each ["Hue", "Saturation", "Luminance"] as type}
				{@const fieldKey = _getNormalizedKey(selectedColor.key, type)}
				{@const def = _getParam(fieldKey)}
				{@const gradient = type === "Hue" ? _getHueGradient(selectedColor.h) : type === "Saturation" ? _getChromaGradient(selectedColor.h) : _getBrightnessGradient(selectedColor.h)}

				{#if def}
					<ParameterSliderUnit
						definition={def}
						value={getVal(fieldKey)}
						originalValue={getOrig(fieldKey)}
						onchange={onchange}
						trackBackground={gradient}
					/>
				{/if}
			{/each}
		</div>
	</div>
</div>
