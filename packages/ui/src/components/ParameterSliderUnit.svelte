<script lang="ts">
	import type { ParameterDefinition } from '../types';

	let {
		definition,
		value,
		originalValue,
		onchange,
		trackBackground
	}: {
		definition: ParameterDefinition;
		value: number;
		originalValue: number;
		onchange: (key: string, value: number) => void;
		trackBackground?: string;
	} = $props();

	let localValue = $state(0);
	let isInputError = $state(false);

	$effect(() => {
		localValue = value;
	});

	let isDirty = $derived(value !== originalValue);

	let progressPercent = $derived(
		definition.max !== definition.min
			? ((localValue - definition.min) / (definition.max - definition.min)) * 100
			: 0
	);

	let isCentered = $derived(definition.min < 0 && definition.max > 0);

	let zeroPercent = $derived(
		isCentered ? ((0 - definition.min) / (definition.max - definition.min)) * 100 : 0
	);

	let trackStyle = $derived.by(() => {
		if (trackBackground) {
			return trackBackground;
		}
		if (!isCentered) {
			return `linear-gradient(to right, var(--color-interactive) ${progressPercent}%, var(--color-surface-elevated) ${progressPercent}%)`;
		}

		if (localValue > 0) {
			return `linear-gradient(to right, var(--color-surface-elevated) ${zeroPercent}%, var(--color-interactive) ${zeroPercent}%, var(--color-interactive) ${progressPercent}%, var(--color-surface-elevated) ${progressPercent}%)`;
		} else {
			return `linear-gradient(to right, var(--color-surface-elevated) ${progressPercent}%, var(--color-interactive) ${progressPercent}%, var(--color-interactive) ${zeroPercent}%, var(--color-surface-elevated) ${zeroPercent}%)`;
		}
	});

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

	// biome-ignore lint/correctness/noUnusedVariables: used in Svelte template
	function handleDoubleClick() {
		localValue = definition.defaultValue;
		onchange(definition.key, localValue);
	}

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
	<div class="flex items-center justify-between pointer-events-none">
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="flex items-center gap-1.5 cursor-pointer pointer-events-auto" ondblclick={handleDoubleClick}>
			<label for={definition.key} class="text-xs font-medium text-foreground select-none cursor-pointer">
				{definition.label}
			</label>
			{#if isDirty}
				<div class="h-1 w-1 rounded-full bg-modified" aria-label="Modified" title="Modified"></div>
			{/if}
		</div>
		<input
			type="number"
			id="{definition.key}-number"
			class="w-16 h-6 px-1 text-right text-xs font-mono origin-right pointer-events-auto text-foreground bg-surface-elevated border focus:outline-none transition-colors duration-100 {isInputError ? 'border-error' : 'border-border focus:border-focus focus:ring-1 focus:ring-focus'}"
			value={localValue}
			min={definition.min}
			max={definition.max}
			step={definition.step}
			onblur={handleNumberDone}
			onkeydown={handleNumberDone}
			oninput={handleNumberInput}
		/>
	</div>

	<!-- Bottom row: Slider -->
	<input
		type="range"
		id={definition.key}
		class="parameter-slider w-full h-1 appearance-none cursor-pointer rounded-sm"
		style="background: {trackStyle}"
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
input[type=range] {
	-webkit-appearance: none;
	appearance: none;
}
input[type=range]::-webkit-slider-thumb {
	-webkit-appearance: none;
	width: 12px;
	height: 12px;
	border-radius: 50%;
	background: var(--color-foreground);
	cursor: pointer;
}
input[type=number]::-webkit-inner-spin-button,
input[type=number]::-webkit-outer-spin-button {
	-webkit-appearance: none;
	margin: 0;
}
input[type=number] {
	-moz-appearance: textfield;
	appearance: textfield;
}
</style>
