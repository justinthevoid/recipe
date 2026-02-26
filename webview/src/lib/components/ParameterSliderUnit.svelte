<script lang="ts">
	import type { ParameterDefinition } from '$lib/types';

	let {
		definition,
		value,
		originalValue,
		onchange
	}: {
		definition: ParameterDefinition;
		value: number;
		originalValue: number;
		onchange: (key: string, value: number) => void;
	} = $props();

	// H2 fix: Use $state(0) and sync entirely through $effect to avoid
	// the state_referenced_locally compiler warning.
	let localValue = $state(0);
	let isInputError = $state(false);

	// Sync localValue with external value changes (e.g. from Reset All or initial prop)
	$effect(() => {
		localValue = value;
	});

	// 2.9 Dirty-state indicator
	let isDirty = $derived(value !== originalValue);

	// L1 fix: Compute slider progress percentage for active-portion track coloring
	let progressPercent = $derived(
		definition.max !== definition.min
			? ((localValue - definition.min) / (definition.max - definition.min)) * 100
			: 0
	);

	// 2.7 Emit onchange callback ONLY on pointerup (H4 fix: removed duplicate onchange)
	// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
	function handleSliderCommit(e: Event) {
		const target = e.target as HTMLInputElement;
		const val = Number(target.value);
		localValue = val;
		onchange(definition.key, val);
	}

	// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
	function handleInput(e: Event) {
		const target = e.target as HTMLInputElement;
		localValue = Number(target.value);
	}

	// 2.8 Implement double-click reset
	// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
	function handleDoubleClick() {
		localValue = definition.defaultValue;
		onchange(definition.key, localValue);
	}

	// 2.10 Direct numeric entry
	// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
	function handleNumberDone(e: KeyboardEvent | FocusEvent) {
		if (e.type === 'keydown' && (e as KeyboardEvent).key !== 'Enter') return;
		
		let val = Number((e.target as HTMLInputElement).value);
		let hasError = false;

		if (val < definition.min) {
			val = definition.min;
			hasError = true;
		}
		if (val > definition.max) {
			val = definition.max;
			hasError = true;
		}

		if (hasError) {
			isInputError = true;
			setTimeout(() => {
				isInputError = false;
			}, 300);
		}

		localValue = val;
		onchange(definition.key, localValue);
	}

	// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
	function handleNumberInput(e: Event) {
		localValue = Number((e.target as HTMLInputElement).value);
	}

	// 2.11 Keyboard accessibility (H3 fix: added Arrow ±1 step, Home/End)
	// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
	function handleKeyDown(e: KeyboardEvent) {
		let val = localValue;
		let handled = false;

		if (e.key === 'Home') {
			e.preventDefault();
			val = definition.min;
			handled = true;
		} else if (e.key === 'End') {
			e.preventDefault();
			val = definition.max;
			handled = true;
		} else if (e.shiftKey && (e.key === 'ArrowLeft' || e.key === 'ArrowRight')) {
			e.preventDefault();
			const direction = e.key === 'ArrowRight' ? 1 : -1;
			val = localValue + direction * definition.step * 10;
			handled = true;
		}

		if (handled) {
			val = Math.max(definition.min, Math.min(definition.max, val));
			localValue = val;
			onchange(definition.key, localValue);
		}
	}
</script>

<div class="flex flex-col gap-1 w-full">
	<!-- Top row layout -->
	<div class="flex items-center justify-between">
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="flex items-center gap-1.5 cursor-pointer" ondblclick={handleDoubleClick}>
			<label for={definition.key} class="text-xs font-medium text-(--vscode-editor-foreground) select-none cursor-pointer">
				{definition.label}
			</label>
			{#if isDirty}
				<div class="h-1 w-1 rounded-full bg-(--vscode-editorOverviewRuler-modifiedForeground)" aria-label="Modified" title="Modified"></div>
			{/if}
		</div>
		<input
			type="number"
			id="{definition.key}-number"
			class="w-16 h-6 px-1 text-right text-xs font-mono text-(--vscode-editor-foreground) bg-(--vscode-input-background) border focus:outline-none transition-colors duration-100 {isInputError ? 'border-(--vscode-inputValidation-errorBorder)' : 'border-(--vscode-input-border) focus:border-(--vscode-focusBorder) focus:ring-1 focus:ring-(--vscode-focusBorder)'}"
			value={localValue}
			min={definition.min}
			max={definition.max}
			step={definition.step}
			onblur={handleNumberDone}
			onkeydown={handleNumberDone}
			oninput={handleNumberInput}
		/>
	</div>

	<!-- Bottom row: Slider (L1 fix: active portion coloring via gradient) -->
	<input
		type="range"
		id={definition.key}
		class="parameter-slider w-full h-1 appearance-none cursor-pointer rounded-sm"
		style="--progress: {progressPercent}%"
		value={localValue}
		min={definition.min}
		max={definition.max}
		step={definition.step}
		oninput={handleInput}
		onpointerup={handleSliderCommit}
		ondblclick={handleDoubleClick}
		onkeydown={handleKeyDown}
		aria-label={definition.label}
		aria-valuemin={definition.min}
		aria-valuemax={definition.max}
		aria-valuenow={localValue}
	/>
</div>

<style>
/* Remove default webkit styles to ensure custom styling works */
input[type=range] {
	-webkit-appearance: none; 
}
/* L1 fix: Active portion coloring — left of thumb uses button bg, right uses muted track */
input[type=range].parameter-slider {
	background: linear-gradient(
		to right,
		var(--vscode-button-background) var(--progress, 0%),
		var(--vscode-input-background) var(--progress, 0%)
	);
}
input[type=range]::-webkit-slider-thumb {
	-webkit-appearance: none;
	width: 12px;
	height: 12px;
	border-radius: 50%;
	background: var(--vscode-editor-foreground);
	cursor: pointer;
}
/* Hide number input spinners */
input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
	-webkit-appearance: none; 
	margin: 0; 
}
input[type=number] {
	-moz-appearance: textfield;
}
</style>
