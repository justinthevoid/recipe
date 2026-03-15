<script lang="ts">
	import {
		CollapsibleSection,
		ColorBlender,
		ColorGrading,
		getNested,
		ParameterSliderUnit,
		PhotoPreview,
		ToneCurveVisual,
	} from "@recipe/ui";
	import type { ParameterDefinition, UniversalRecipe } from "@recipe/ui";
	import { convertAndDownload } from "$lib/converter.svelte";
	import {
		store,
		closeEditor,
		resetRecipe,
		setPreviewImage,
		undo,
		redo,
		updateParameter,
	} from "$lib/stores.svelte";
	import { loadImage } from "$lib/image-loader";

	let {
		wasmReady,
		generateLUT,
	}: {
		wasmReady: boolean;
		generateLUT: (recipeJSON: string, size: number) => Promise<Float32Array>;
	} = $props();

	let isConverting = $state(false);
	let convertError = $state<string | null>(null);
	let showFormatPicker = $state(false);

	// Static parameter definitions (same as Go backend)
	const basicParams: ParameterDefinition[] = [
		{ key: "exposure", label: "Exposure", type: "continuous", min: -5, max: 5, step: 0.01, defaultValue: 0, group: "Basic" },
		{ key: "contrast", label: "Contrast", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Basic" },
		{ key: "highlights", label: "Highlights", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Basic" },
		{ key: "shadows", label: "Shadows", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Basic" },
		{ key: "whites", label: "Whites", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Basic" },
		{ key: "blacks", label: "Blacks", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Basic" },
		{ key: "clarity", label: "Clarity", type: "continuous", min: -5, max: 5, step: 0.25, defaultValue: 0, group: "Basic" },
		{ key: "saturation", label: "Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Basic" },
		{ key: "sharpness", label: "Sharpness", type: "continuous", min: 0, max: 150, step: 1, defaultValue: 0, group: "Basic" },
	];

	const colorMixerParams: ParameterDefinition[] = [
		{ key: "red.hue", label: "Red Hue", type: "continuous", min: -180, max: 180, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "red.saturation", label: "Red Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "red.luminance", label: "Red Luminance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "orange.hue", label: "Orange Hue", type: "continuous", min: -180, max: 180, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "orange.saturation", label: "Orange Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "orange.luminance", label: "Orange Luminance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "yellow.hue", label: "Yellow Hue", type: "continuous", min: -180, max: 180, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "yellow.saturation", label: "Yellow Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "yellow.luminance", label: "Yellow Luminance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "green.hue", label: "Green Hue", type: "continuous", min: -180, max: 180, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "green.saturation", label: "Green Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "green.luminance", label: "Green Luminance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "aqua.hue", label: "Aqua Hue", type: "continuous", min: -180, max: 180, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "aqua.saturation", label: "Aqua Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "aqua.luminance", label: "Aqua Luminance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "blue.hue", label: "Blue Hue", type: "continuous", min: -180, max: 180, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "blue.saturation", label: "Blue Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "blue.luminance", label: "Blue Luminance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "purple.hue", label: "Purple Hue", type: "continuous", min: -180, max: 180, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "purple.saturation", label: "Purple Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "purple.luminance", label: "Purple Luminance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "magenta.hue", label: "Magenta Hue", type: "continuous", min: -180, max: 180, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "magenta.saturation", label: "Magenta Saturation", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
		{ key: "magenta.luminance", label: "Magenta Luminance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Mixer" },
	];

	const colorGradingParams: ParameterDefinition[] = [
		{ key: "colorGrading.highlights.hue", label: "Highlights Hue", type: "continuous", min: 0, max: 360, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.highlights.chroma", label: "Highlights Chroma", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.highlights.brightness", label: "Highlights Brightness", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.midtone.hue", label: "Midtone Hue", type: "continuous", min: 0, max: 360, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.midtone.chroma", label: "Midtone Chroma", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.midtone.brightness", label: "Midtone Brightness", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.shadows.hue", label: "Shadows Hue", type: "continuous", min: 0, max: 360, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.shadows.chroma", label: "Shadows Chroma", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.shadows.brightness", label: "Shadows Brightness", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Grading" },
		{ key: "colorGrading.blending", label: "Blending", type: "continuous", min: 0, max: 100, step: 1, defaultValue: 50, group: "Color Grading" },
		{ key: "colorGrading.balance", label: "Balance", type: "continuous", min: -100, max: 100, step: 1, defaultValue: 0, group: "Color Grading" },
	];

	const formats = ["np3", "xmp", "lrtemplate"];

	function handleParameterChange(key: string, value: number) {
		updateParameter(key, value);
	}

	async function handlePhotoUpload() {
		const input = document.createElement("input");
		input.type = "file";
		input.accept = "image/jpeg,image/png,image/webp";
		input.onchange = async () => {
			const file = input.files?.[0];
			if (!file) return;
			try {
				const img = await loadImage(file);
				setPreviewImage(img);
			} catch (err) {
				console.error("Image load failed:", err);
			}
		};
		input.click();
	}

	async function handleConvert(format: string) {
		showFormatPicker = false;
		const recipe = store.currentRecipe;
		if (!recipe) return;

		isConverting = true;
		convertError = null;

		try {
			await convertAndDownload(recipe, format, store.currentFileName || "preset");
		} catch (err) {
			convertError = err instanceof Error ? err.message : "Conversion failed";
		} finally {
			isConverting = false;
		}
	}

	function getRecipeValue(key: string): number {
		if (!store.currentRecipe) return 0;
		return Number(getNested(store.currentRecipe, key) ?? 0);
	}

	function getOriginalValue(key: string): number {
		if (!store.originalRecipe) return 0;
		return Number(getNested(store.originalRecipe, key) ?? 0);
	}
</script>

<div class="fixed inset-0 bg-canvas-base z-50 flex flex-col">
	<!-- Toolbar -->
	<div class="flex items-center justify-between px-4 h-12 border-b border-border bg-surface-elevated/80 backdrop-blur-sm shrink-0">
		<div class="flex items-center gap-3">
			<button
				type="button"
				class="text-sm text-foreground-muted hover:text-foreground transition-colors"
				onclick={closeEditor}
			>
				← Back
			</button>
			<span class="text-sm text-foreground font-medium truncate max-w-[200px]">
				{store.currentFileName || "Untitled"}
			</span>
			{#if store.isDirty}
				<span class="text-xs text-modified">Modified</span>
			{/if}
		</div>

		<div class="flex items-center gap-2">
			<button
				type="button"
				class="px-2 py-1 text-xs text-foreground-muted hover:text-foreground disabled:opacity-30 transition-colors"
				disabled={!store.canUndo}
				onclick={undo}
				title="Undo"
			>
				Undo
			</button>
			<button
				type="button"
				class="px-2 py-1 text-xs text-foreground-muted hover:text-foreground disabled:opacity-30 transition-colors"
				disabled={!store.canRedo}
				onclick={redo}
				title="Redo"
			>
				Redo
			</button>

			<button
				type="button"
				class="px-2 py-1 text-xs text-foreground-muted hover:text-foreground transition-colors"
				onclick={resetRecipe}
				title="Reset All"
			>
				Reset
			</button>

			<button
				type="button"
				class="px-2 py-1 text-xs text-foreground-muted hover:text-foreground transition-colors"
				onclick={handlePhotoUpload}
			>
				Change Photo
			</button>

			<div class="relative">
				<button
					type="button"
					class="px-3 py-1.5 text-xs font-medium bg-interactive text-interactive-foreground rounded transition-colors hover:opacity-90 disabled:opacity-50"
					disabled={isConverting}
					onclick={() => showFormatPicker = !showFormatPicker}
				>
					{isConverting ? "Converting..." : "Convert & Download"}
				</button>

				{#if showFormatPicker}
					<!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events -->
					<div
						class="absolute right-0 top-full mt-1 bg-surface-elevated border border-border rounded shadow-lg z-10 min-w-[140px]"
						onclick={(e: MouseEvent) => e.stopPropagation()}
					>
						{#each formats as fmt}
							<button
								type="button"
								class="w-full text-left px-3 py-2 text-xs text-foreground hover:bg-surface-hover transition-colors"
								onclick={() => handleConvert(fmt)}
							>
								{fmt.toUpperCase()}
								{#if fmt !== "np3"}
									<span class="text-foreground-muted ml-1">(via NP3)</span>
								{/if}
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	</div>

	{#if convertError}
		<div class="px-4 py-2 bg-error/10 text-error text-xs border-b border-error/20">
			{convertError}
		</div>
	{/if}

	<!-- Main editor area -->
	<div class="flex flex-1 overflow-hidden min-h-0" style="height: calc(100vh - 3rem);">
		<!-- Canvas area -->
		<div class="flex-1 min-w-0 flex flex-col p-4">
			{#if !store.previewImage}
				<div class="flex flex-col items-center justify-center flex-1 bg-surface-elevated/30 rounded-xl gap-4 border border-dashed border-border">
					<span class="text-4xl">📷</span>
					<p class="text-sm text-foreground-muted text-center max-w-xs">
						Upload a photo to see a live preview of your preset adjustments
					</p>
					<button
						type="button"
						class="px-4 py-2 text-sm font-medium bg-interactive text-interactive-foreground rounded-lg hover:opacity-90 transition-opacity"
						onclick={handlePhotoUpload}
					>
						Upload Photo
					</button>
				</div>
			{:else}
				<div class="flex-1 min-h-0" style="position: relative;">
					<PhotoPreview
						recipe={store.currentRecipe}
						imageData={store.previewImage}
						{wasmReady}
						{generateLUT}
					/>
				</div>
			{/if}
		</div>

		<!-- Sidebar -->
		<div class="w-80 border-l border-border overflow-y-auto bg-surface-elevated/50 p-4 space-y-2 shrink-0 hidden lg:block">
			<CollapsibleSection title="Basic" expanded={true}>
				<div class="flex flex-col gap-3">
					{#each basicParams as def}
						<ParameterSliderUnit
							definition={def}
							value={getRecipeValue(def.key)}
							originalValue={getOriginalValue(def.key)}
							onchange={handleParameterChange}
						/>
					{/each}
				</div>
			</CollapsibleSection>

			<CollapsibleSection title="Color Mixer">
				<ColorBlender
					parameters={colorMixerParams}
					recipe={store.currentRecipe ?? {}}
					originalRecipe={store.originalRecipe ?? {}}
					onchange={handleParameterChange}
				/>
			</CollapsibleSection>

			<CollapsibleSection title="Tone Curve">
				<ToneCurveVisual
					points={store.currentRecipe?.pointCurve ?? []}
				/>
			</CollapsibleSection>

			<CollapsibleSection title="Color Grading">
				<ColorGrading
					parameters={colorGradingParams}
					recipe={store.currentRecipe ?? {}}
					originalRecipe={store.originalRecipe ?? {}}
					onchange={handleParameterChange}
				/>
			</CollapsibleSection>
		</div>
	</div>
</div>
