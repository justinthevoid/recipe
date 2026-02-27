<script lang="ts">
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
		return Number(recipe?.[key] ?? 0);
	}

	function getOrig(key: string) {
		return Number(originalRecipe?.[key] ?? 0);
	}

	function _getNormalizedKey(key: string, type: string) {
		// handle legacy or specific key mapping if needed. But currently it's "red.hue", "red.saturation", "red.luminance".
		// Note: types are capitalized later. The struct uses lower case for subfields: "red.hue", etc.
		return `${key}.${type.toLowerCase()}`;
	}

	function isColorDirty(prefix: string) {
		return getVal(`${prefix}.hue`) !== getOrig(`${prefix}.hue`) ||
		       getVal(`${prefix}.saturation`) !== getOrig(`${prefix}.saturation`) ||
		       getVal(`${prefix}.luminance`) !== getOrig(`${prefix}.luminance`);
	}

	// Generate colorful track gradients based on NX Studio
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
	<!-- Color Swatch Selection Row (Radio button style) -->
	<div class="flex items-center justify-between px-1">
		{#each colors as color, i}
			<button
				type="button"
				class="relative flex items-center justify-center rounded-full transition-all {selectedColorIndex === i ? 'w-6 h-6 ring-2 ring-(--vscode-focusBorder) ring-offset-1 ring-offset-(--vscode-sideBar-background)' : 'w-5 h-5 opacity-80 hover:opacity-100 hover:scale-110'}"
				style="background-color: {color.hex}"
				onclick={() => selectedColorIndex = i}
				title={color.name}
				aria-label="Select {color.name}"
			>
				{#if isColorDirty(color.key)}
					<div class="absolute -top-1 -right-1 h-2 w-2 rounded-full border border-(--vscode-sideBar-background) bg-(--vscode-editorOverviewRuler-modifiedForeground)"></div>
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

