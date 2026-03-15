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

	type Tab = "highlights" | "midtone" | "shadows" | "global";
	const tabs: { id: Tab; label: string }[] = [
		{ id: "highlights", label: "Highlights" },
		{ id: "midtone", label: "Midtone" },
		{ id: "shadows", label: "Shadows" },
		{ id: "global", label: "Global" }
	];

	let currentTab = $state<Tab>("highlights");

	function getParam(key: string) {
		return parameters.find(p => p.key === key);
	}

	function getVal(key: string) {
		return Number(getNested(recipe, key) ?? 0);
	}

	function getOrig(key: string) {
		return Number(getNested(originalRecipe, key) ?? 0);
	}

	function isZoneDirty(zone: Tab) {
		if (zone === "global") {
			return getVal(`colorGrading.blending`) !== getOrig(`colorGrading.blending`) ||
			       getVal(`colorGrading.balance`) !== getOrig(`colorGrading.balance`);
		}
		return getVal(`colorGrading.${zone}.hue`) !== getOrig(`colorGrading.${zone}.hue`) ||
		       getVal(`colorGrading.${zone}.chroma`) !== getOrig(`colorGrading.${zone}.chroma`) ||
		       getVal(`colorGrading.${zone}.brightness`) !== getOrig(`colorGrading.${zone}.brightness`);
	}

	function _getHueGradient() {
		return `linear-gradient(to right, #ff0000, #ffff00, #00ff00, #00ffff, #0000ff, #ff00ff, #ff0000)`;
	}

	function _getChromaGradient(h: number) {
		return `linear-gradient(to right, hsl(${h}, 0%, 50%), hsl(${h}, 50%, 50%), hsl(${h}, 100%, 50%))`;
	}

	function _getBrightnessGradient(h: number) {
		return `linear-gradient(to right, black, hsl(${h}, 100%, 50%), white)`;
	}

	let hueDef = $derived(getParam(`colorGrading.${currentTab}.hue`));
	let chromaDef = $derived(getParam(`colorGrading.${currentTab}.chroma`));
	let brightnessDef = $derived(getParam(`colorGrading.${currentTab}.brightness`));
	let blendDef = $derived(getParam("colorGrading.blending"));
	let balanceDef = $derived(getParam("colorGrading.balance"));
	let currentHue = $derived(getVal(`colorGrading.${currentTab}.hue`));
</script>

<div class="flex flex-col gap-4">
	<!-- Tab Navigation -->
	<div class="flex items-center w-full gap-1 p-1 bg-muted/50 rounded flex-wrap">
		{#each tabs as tab}
			<button
				type="button"
				class="relative flex-1 py-1.5 text-xs font-medium rounded-sm transition-all {currentTab === tab.id ? 'bg-background text-foreground shadow-xs' : 'text-muted-foreground hover:bg-muted hover:text-foreground'}"
				onclick={() => currentTab = tab.id}
			>
				{tab.label}
				{#if isZoneDirty(tab.id)}
					<div class="absolute top-1 right-2 h-1.5 w-1.5 rounded-full bg-modified"></div>
				{/if}
			</button>
		{/each}
	</div>

	<!-- Controls Section -->
	<div class="p-4 bg-muted/30 rounded border border-border flex flex-col gap-5">
		{#if currentTab !== "global"}
			<div class="flex flex-col gap-4">
				{#if hueDef}
					<ParameterSliderUnit
						definition={hueDef}
						value={currentHue}
						originalValue={getOrig(hueDef.key)}
						onchange={onchange}
						trackBackground={_getHueGradient()}
					/>
				{/if}

				{#if chromaDef}
					<ParameterSliderUnit
						definition={chromaDef}
						value={getVal(chromaDef.key)}
						originalValue={getOrig(chromaDef.key)}
						onchange={onchange}
						trackBackground={_getChromaGradient(currentHue)}
					/>
				{/if}

				{#if brightnessDef}
					<ParameterSliderUnit
						definition={brightnessDef}
						value={getVal(brightnessDef.key)}
						originalValue={getOrig(brightnessDef.key)}
						onchange={onchange}
						trackBackground={_getBrightnessGradient(currentHue)}
					/>
				{/if}
			</div>
		{:else}
			<div class="flex flex-col gap-4">
				{#if blendDef}
					<ParameterSliderUnit
						definition={blendDef}
						value={getVal(blendDef.key)}
						originalValue={getOrig(blendDef.key)}
						onchange={onchange}
					/>
				{/if}

				{#if balanceDef}
					<ParameterSliderUnit
						definition={balanceDef}
						value={getVal(balanceDef.key)}
						originalValue={getOrig(balanceDef.key)}
						onchange={onchange}
					/>
				{/if}
			</div>
		{/if}
	</div>
</div>
