<script lang="ts">
	import type { ToneCurvePoint, UniversalRecipe } from "@recipe/ui";
	import { convertAndDownload } from "$lib/converter.svelte";
	import {
		history,
		closeEditor,
		resetRecipe,
		undo,
		redo,
		updateParameter,
		updateMetadata,
		updateCurvePoints,
		type CurveChannel,
	} from "$lib/stores.svelte";
	import CurveEditor from "./CurveEditor.svelte";
	import {
		editorModeStore,
		currentRecipeStore,
		originalRecipeStore,
		currentFileNameStore,
	} from "$lib/shared-stores";

	// Read shared nanostores reactively via subscriptions
	let editorMode = $state(editorModeStore.get());
	let currentRecipe = $state<UniversalRecipe | null>(currentRecipeStore.get());
	let originalRecipe = $state<UniversalRecipe | null>(originalRecipeStore.get());
	let currentFileName = $state(currentFileNameStore.get());

	$effect(() => {
		const unsubs = [
			editorModeStore.subscribe(v => { editorMode = v; }),
			currentRecipeStore.subscribe(v => { currentRecipe = v; }),
			originalRecipeStore.subscribe(v => { originalRecipe = v; }),
			currentFileNameStore.subscribe(v => { currentFileName = v; }),
		];
		return () => unsubs.forEach(u => u());
	});

	let isConverting = $state(false);
	let convertError = $state<string | null>(null);
	let showFormatPicker = $state(false);

	const formats = ["np3", "xmp"];

	// Inline getNested — avoids @recipe/ui runtime dependency
	function getNested(obj: unknown, path: string): unknown {
		const parts = path.split(".");
		let curr: unknown = obj;
		for (const p of parts) curr = (curr as Record<string, unknown>)?.[p];
		return curr;
	}

	function getVal(key: string): number {
		return Number(getNested(currentRecipe, key) ?? 0);
	}

	function handleChange(key: string, value: number) {
		updateParameter(key, value);
	}

	async function handleConvert(format: string) {
		showFormatPicker = false;
		const recipe = currentRecipe;
		if (!recipe) return;
		isConverting = true;
		convertError = null;
		try {
			await convertAndDownload(recipe, format, currentFileName || "preset");
		} catch (err) {
			convertError = err instanceof Error ? err.message : "Conversion failed";
		} finally {
			isConverting = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === "Escape") closeEditor();
	}

	// Parameter definitions
	const basicParams = [
		{ key: "exposure",   label: "Exposure",   min: -5,   max: 5,   step: 0.01 },
		{ key: "contrast",   label: "Contrast",   min: -100, max: 100, step: 1 },
		{ key: "highlights", label: "Highlights", min: -100, max: 100, step: 1 },
		{ key: "shadows",    label: "Shadows",    min: -100, max: 100, step: 1 },
		{ key: "whites",     label: "Whites",     min: -100, max: 100, step: 1 },
		{ key: "blacks",     label: "Blacks",     min: -100, max: 100, step: 1 },
		{ key: "clarity",    label: "Clarity",    min: -100, max: 100, step: 1 },
		{ key: "saturation", label: "Saturation", min: -100, max: 100, step: 1 },
		{ key: "sharpness",  label: "Sharpness",  min: 0,    max: 150, step: 1 },
	] as const;

	const colorChannels = ["red", "orange", "yellow", "green", "aqua", "blue", "purple", "magenta"] as const;

	const colorGradingParams = [
		{ key: "colorGrading.highlights.hue",    label: "Highlights Hue",        min: 0,    max: 360, step: 1 },
		{ key: "colorGrading.highlights.chroma", label: "Highlights Chroma",     min: -100, max: 100, step: 1 },
		{ key: "colorGrading.highlights.brightness", label: "Highlights Brightness", min: -100, max: 100, step: 1 },
		{ key: "colorGrading.midtone.hue",       label: "Midtone Hue",           min: 0,    max: 360, step: 1 },
		{ key: "colorGrading.midtone.chroma",    label: "Midtone Chroma",        min: -100, max: 100, step: 1 },
		{ key: "colorGrading.midtone.brightness","label": "Midtone Brightness",  min: -100, max: 100, step: 1 },
		{ key: "colorGrading.shadows.hue",       label: "Shadows Hue",           min: 0,    max: 360, step: 1 },
		{ key: "colorGrading.shadows.chroma",    label: "Shadows Chroma",        min: -100, max: 100, step: 1 },
		{ key: "colorGrading.shadows.brightness","label": "Shadows Brightness",  min: -100, max: 100, step: 1 },
		{ key: "colorGrading.blending",          label: "Blending",              min: 0,    max: 100, step: 1 },
		{ key: "colorGrading.balance",           label: "Balance",               min: -100, max: 100, step: 1 },
	] as const;

	const parametricCurveParams = [
		{ key: "toneCurveShadows",    label: "Shadows",    min: -100, max: 100, step: 1 },
		{ key: "toneCurveDarks",      label: "Darks",      min: -100, max: 100, step: 1 },
		{ key: "toneCurveLights",     label: "Lights",     min: -100, max: 100, step: 1 },
		{ key: "toneCurveHighlights", label: "Highlights", min: -100, max: 100, step: 1 },
	] as const;

	const curveChannelDefs = [
		{ id: "master" as const,  label: "Master", key: "pointCurve"      as CurveChannel, color: "rgba(255,255,255,0.9)" },
		{ id: "red"    as const,  label: "Red",    key: "pointCurveRed"   as CurveChannel, color: "#ff6b6b" },
		{ id: "green"  as const,  label: "Green",  key: "pointCurveGreen" as CurveChannel, color: "#6bff8a" },
		{ id: "blue"   as const,  label: "Blue",   key: "pointCurveBlue"  as CurveChannel, color: "#6b9fff" },
	];

	function hasCurveData(r: UniversalRecipe | null): boolean {
		if (!r) return false;
		return curveChannelDefs.some((ch) => {
			const v = r[ch.key];
			return Array.isArray(v) && v.length > 0;
		}) || [r.toneCurveShadows, r.toneCurveDarks, r.toneCurveLights, r.toneCurveHighlights]
			.some((v) => v != null && v !== 0);
	}

	function hasParametricCurves(r: UniversalRecipe | null): boolean {
		if (!r) return false;
		return [r.toneCurveShadows, r.toneCurveDarks, r.toneCurveLights, r.toneCurveHighlights]
			.some((v) => v != null && v !== 0);
	}

	let activeCurveTab = $state<"master" | "red" | "green" | "blue">("master");

	const availableCurveTabs = $derived(
		curveChannelDefs.filter((ch) => {
			const v = currentRecipe?.[ch.key];
			return Array.isArray(v) && v.length > 0;
		}),
	);

	// Keep active tab valid when recipe changes
	$effect(() => {
		const tabs = availableCurveTabs;
		if (tabs.length > 0 && !tabs.some((t) => t.id === activeCurveTab)) {
			activeCurveTab = tabs[0].id;
		}
	});

	const activeTabDef = $derived(
		availableCurveTabs.find((t) => t.id === activeCurveTab) ?? availableCurveTabs[0],
	);

	const activePoints = $derived.by(() => {
		if (!activeTabDef || !currentRecipe) return [] as ToneCurvePoint[];
		const val = currentRecipe[activeTabDef.key];
		return Array.isArray(val) ? (val as ToneCurvePoint[]) : ([] as ToneCurvePoint[]);
	});

	const _enc = new TextEncoder();

	// Clamp a string to at most maxBytes UTF-8 bytes, cutting at code-point boundaries.
	function limitToBytes(value: string, maxBytes: number): string {
		if (_enc.encode(value).length <= maxBytes) return value;
		let result = "";
		let count = 0;
		for (const ch of value) {
			const bytes = _enc.encode(ch).length;
			if (count + bytes > maxBytes) break;
			result += ch;
			count += bytes;
		}
		return result;
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if editorMode}
<!-- Backdrop -->
<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
<div
	class="fixed inset-0 z-50 flex items-center justify-center p-4"
	role="presentation"
	onclick={(e) => { if (e.target === e.currentTarget) closeEditor(); }}
>
	<!-- Dim overlay — pointer-events-none so clicks pass through to the outer wrapper for close-on-backdrop -->
	<div class="absolute inset-0 bg-black/60 backdrop-blur-sm pointer-events-none"></div>

	<!-- Modal card -->
	<div
		class="glass-regular rounded-2xl relative z-10 w-full max-w-lg flex flex-col"
		style="max-height: 85dvh;"
	>
		<!-- Header -->
		<div class="flex items-center justify-between px-5 py-3.5 border-b border-white/5 shrink-0">
			<div class="flex items-center gap-2 min-w-0">
				<span class="text-sm font-medium text-foreground truncate">
					{currentFileName || "Untitled"}
				</span>
				{#if history.isDirty}
					<span class="text-xs text-modified shrink-0">Modified</span>
				{/if}
			</div>
			<button
				type="button"
				class="text-foreground-muted hover:text-foreground text-xl leading-none transition-colors ml-4 shrink-0"
				onclick={closeEditor}
				aria-label="Close"
			>
				×
			</button>
		</div>

		<!-- Scrollable body -->
		<div class="flex-1 overflow-y-auto px-5 py-4 space-y-4 min-h-0">

			<!-- Metadata section -->
			<details open>
				<summary class="text-xs font-semibold uppercase tracking-wider text-foreground-muted cursor-pointer select-none hover:text-foreground transition-colors py-1">
					Metadata
				</summary>
				<div class="mt-3 space-y-3">
					<label class="flex items-center gap-3 text-xs">
						<span class="w-24 shrink-0 text-foreground-muted">Name</span>
						<input
							type="text"
							value={currentRecipe?.name ?? ""}
							oninput={(e) => {
								if (e.isComposing) return;
								const clamped = limitToBytes((e.currentTarget as HTMLInputElement).value, 20);
								(e.currentTarget as HTMLInputElement).value = clamped;
								updateMetadata("name", clamped);
							}}
							class="flex-1 bg-transparent border border-white/10 rounded px-2 py-1 text-foreground focus:outline-none focus:border-white/30"
						/>
					</label>
					<label class="flex items-start gap-3 text-xs">
						<span class="w-24 shrink-0 text-foreground-muted pt-1">Description</span>
						<textarea
							rows="2"
							value={currentRecipe?.description ?? ""}
							oninput={(e) => {
								if (e.isComposing) return;
								const clamped = limitToBytes((e.currentTarget as HTMLTextAreaElement).value, 256);
								(e.currentTarget as HTMLTextAreaElement).value = clamped;
								updateMetadata("description", clamped);
							}}
							class="flex-1 bg-transparent border border-white/10 rounded px-2 py-1 text-foreground focus:outline-none focus:border-white/30 resize-none"
						></textarea>
					</label>
				</div>
			</details>

			<!-- Basic section (open by default) -->
			<details open>
				<summary class="text-xs font-semibold uppercase tracking-wider text-foreground-muted cursor-pointer select-none hover:text-foreground transition-colors py-1">
					Basic
				</summary>
				<div class="mt-3 space-y-3">
					{#each basicParams as p}
						{@const val = getVal(p.key)}
						<label class="flex items-center gap-3 text-xs">
							<span class="w-24 shrink-0 text-foreground-muted">{p.label}</span>
							<input
								type="range"
								min={p.min}
								max={p.max}
								step={p.step}
								value={val}
								oninput={(e) => handleChange(p.key, +(e.currentTarget as HTMLInputElement).value)}
								class="flex-1 accent-interactive h-1 cursor-pointer"
							/>
							<span class="w-10 text-right tabular-nums text-foreground shrink-0">
								{val > 0 ? `+${val}` : val}
							</span>
						</label>
					{/each}
				</div>
			</details>

			<!-- Color Mixer section -->
			<details>
				<summary class="text-xs font-semibold uppercase tracking-wider text-foreground-muted cursor-pointer select-none hover:text-foreground transition-colors py-1">
					Color Mixer
				</summary>
				<div class="mt-3 space-y-4">
					{#each colorChannels as channel}
						<div>
							<p class="text-xs text-foreground-muted/60 mb-2 capitalize">{channel}</p>
							<div class="space-y-2">
								{#each ["hue", "saturation", "luminance"] as aspect}
									{@const key = `${channel}.${aspect}`}
									{@const val = getVal(key)}
									{@const min = aspect === "hue" ? -180 : -100}
									{@const max = aspect === "hue" ? 180 : 100}
									<label class="flex items-center gap-3 text-xs">
										<span class="w-24 shrink-0 text-foreground-muted capitalize">{aspect}</span>
										<input
											type="range"
											{min}
											{max}
											step={1}
											value={val}
											oninput={(e) => handleChange(key, +(e.currentTarget as HTMLInputElement).value)}
											class="flex-1 accent-interactive h-1 cursor-pointer"
										/>
										<span class="w-10 text-right tabular-nums text-foreground shrink-0">
											{val > 0 ? `+${val}` : val}
										</span>
									</label>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			</details>

			<!-- Color Grading section -->
			<details>
				<summary class="text-xs font-semibold uppercase tracking-wider text-foreground-muted cursor-pointer select-none hover:text-foreground transition-colors py-1">
					Color Grading
				</summary>
				<div class="mt-3 space-y-3">
					{#each colorGradingParams as p}
						{@const val = getVal(p.key)}
						<label class="flex items-center gap-3 text-xs">
							<span class="w-32 shrink-0 text-foreground-muted">{p.label}</span>
							<input
								type="range"
								min={p.min}
								max={p.max}
								step={p.step}
								value={val}
								oninput={(e) => handleChange(p.key, +(e.currentTarget as HTMLInputElement).value)}
								class="flex-1 accent-interactive h-1 cursor-pointer"
							/>
							<span class="w-10 text-right tabular-nums text-foreground shrink-0">
								{val > 0 ? `+${val}` : val}
							</span>
						</label>
					{/each}
				</div>
			</details>

			<!-- Tone Curves section — only shown when recipe has curve data -->
			{#if hasCurveData(currentRecipe)}
			<details>
				<summary class="text-xs font-semibold uppercase tracking-wider text-foreground-muted cursor-pointer select-none hover:text-foreground transition-colors py-1">
					Tone Curves
				</summary>
				<div class="mt-3 space-y-3">

					<!-- Point curve channel tabs + canvas editor -->
					{#if availableCurveTabs.length > 0}
					<div>
						<!-- Channel tabs -->
						<div class="flex gap-1 mb-3">
							{#each availableCurveTabs as tab}
								{@const isActive = activeCurveTab === tab.id}
								<button
									type="button"
									class="px-2.5 py-1 text-xs rounded transition-colors {isActive ? 'text-black font-medium' : 'text-foreground-muted hover:text-foreground'}"
									style={isActive
										? (tab.id !== 'master' ? `background:${tab.color}` : 'background:rgba(255,255,255,0.9)')
										: 'background:rgba(255,255,255,0.08)'}
									onclick={() => { activeCurveTab = tab.id; }}
								>{tab.label}</button>
							{/each}
						</div>

						<!-- Canvas curve editor -->
						{#if activeTabDef}
						<div class="flex justify-center">
							<CurveEditor
								points={activePoints}
								color={activeTabDef.color}
								onchange={(pts) => updateCurvePoints(activeTabDef.key, pts)}
							/>
						</div>
						<p class="text-xs text-foreground-muted/50 text-center mt-1.5">
							Click to add point · Right-click to remove · Drag to adjust
						</p>
						{/if}
					</div>
					{/if}

					<!-- Parametric curve sliders -->
					{#if hasParametricCurves(currentRecipe)}
					<div class="space-y-3 pt-1">
						<p class="text-xs text-foreground-muted/60 uppercase tracking-wider">Parametric</p>
						{#each parametricCurveParams as p}
							{@const val = getVal(p.key)}
							<label class="flex items-center gap-3 text-xs">
								<span class="w-24 shrink-0 text-foreground-muted">{p.label}</span>
								<input
									type="range"
									min={p.min}
									max={p.max}
									step={p.step}
									value={val}
									oninput={(e) => handleChange(p.key, +(e.currentTarget as HTMLInputElement).value)}
									class="flex-1 accent-interactive h-1 cursor-pointer"
								/>
								<span class="w-10 text-right tabular-nums text-foreground shrink-0">
									{val > 0 ? `+${val}` : val}
								</span>
							</label>
						{/each}
					</div>
					{/if}
				</div>
			</details>
			{/if}
		</div>

		<!-- Footer -->
		<div class="px-5 py-3.5 border-t border-white/5 shrink-0 space-y-2">
			{#if convertError}
				<p class="text-xs text-error">{convertError}</p>
			{/if}
			<div class="flex items-center justify-between gap-2">
				<!-- Left: undo/redo/reset -->
				<div class="flex items-center gap-1">
					<button
						type="button"
						class="px-2 py-1 text-xs text-foreground-muted hover:text-foreground disabled:opacity-30 transition-colors"
						disabled={!history.canUndo}
						onclick={undo}
					>Undo</button>
					<button
						type="button"
						class="px-2 py-1 text-xs text-foreground-muted hover:text-foreground disabled:opacity-30 transition-colors"
						disabled={!history.canRedo}
						onclick={redo}
					>Redo</button>
					<button
						type="button"
						class="px-2 py-1 text-xs text-foreground-muted hover:text-foreground transition-colors"
						onclick={resetRecipe}
					>Reset</button>
				</div>

				<!-- Right: Convert & Download -->
				<div class="relative">
					<button
						type="button"
						class="px-3 py-1.5 text-xs font-medium bg-interactive text-interactive-foreground rounded-lg transition-opacity hover:opacity-90 disabled:opacity-50"
						disabled={isConverting}
						onclick={() => { showFormatPicker = !showFormatPicker; convertError = null; }}
					>
						{isConverting ? "Converting…" : "Convert & Download"}
					</button>

					{#if showFormatPicker}
						<!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events -->
						<div
							class="absolute right-0 bottom-full mb-1 glass-regular rounded-lg shadow-lg z-10 min-w-[140px] overflow-hidden"
							onclick={(e) => e.stopPropagation()}
						>
							{#each formats as fmt}
								<button
									type="button"
									class="w-full text-left px-3 py-2 text-xs text-foreground hover:bg-white/5 transition-colors"
									onclick={() => handleConvert(fmt)}
								>
									{fmt.toUpperCase()}
								</button>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
</div>
{/if}

<style>
	details > summary {
		list-style: none;
	}
	details > summary::before {
		content: "▶ ";
		font-size: 0.6rem;
		opacity: 0.5;
		transition: transform 0.15s;
		display: inline-block;
	}
	details[open] > summary::before {
		transform: rotate(90deg);
	}
	details > summary::-webkit-details-marker {
		display: none;
	}
</style>
