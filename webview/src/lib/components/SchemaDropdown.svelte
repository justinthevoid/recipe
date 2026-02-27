<script lang="ts">
	import { Select } from 'bits-ui';
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

	// Convert value to string for bits-ui select
	let localValue = $state("");

	// Sync localValue with external value changes
	$effect(() => {
		localValue = value.toString();
	});

	let isDirty = $derived(value !== originalValue);

	// Pre-flight schema validation
	let isValid = $derived(
		(definition.options ?? []).some((opt) => opt.value === value)
	);

	function _handleValueChange(v: string) {
		if (v === undefined) return;
		const numVal = parseInt(v, 10);
		if (Number.isNaN(numVal)) return;
		onchange(definition.key, numVal);
	}

	function _handleDoubleClick(e: MouseEvent) {
		e.preventDefault();
		onchange(definition.key, definition.defaultValue);
	}

	function _handleLabelKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			onchange(definition.key, definition.defaultValue);
		}
	}

	let selectedLabel = $derived(
		definition.options?.find(opt => opt.value.toString() === localValue)?.label ?? "Invalid"
	);
</script>

<div class="flex flex-col gap-1 w-full">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-1.5 cursor-pointer" role="button" tabindex="0" ondblclick={_handleDoubleClick} onkeydown={_handleLabelKeydown} title="Double-click to reset to default">
			<label
				for={definition.key}
				class="text-xs font-medium text-(--vscode-editor-foreground) select-none cursor-pointer"
			>
				{definition.label}
			</label>
			{#if isDirty}
				<div
					class="w-1 h-1 rounded-full bg-(--vscode-editorOverviewRuler-modifiedForeground)"
					aria-label="Modified"
					title="Modified"
				></div>
			{/if}
		</div>

		<Select.Root
			type="single"
			bind:value={localValue}
			onValueChange={_handleValueChange}
			disabled={!isValid}
		>
			<Select.Trigger
				id={definition.key}
				data-slot="select-trigger"
				class="w-[180px] h-6 px-2 text-xs font-mono text-left focus:ring-1 focus:ring-(--vscode-focusBorder) transition-colors duration-100 {!isValid ? 'warning-hash-pattern text-(--vscode-editorWarning-foreground)' : 'bg-(--vscode-input-background) text-(--vscode-editor-foreground) border-(--vscode-input-border)'}"
				title={!isValid ? `Invalid value: ${value}` : ''}
			>
				{selectedLabel}
			</Select.Trigger>
			<Select.Content class="bg-(--vscode-dropdown-background) border-(--vscode-dropdown-border) text-(--vscode-dropdown-foreground)">
				{#each definition.options || [] as option}
					<Select.Item
						value={option.value.toString()}
						class="text-xs font-mono hover:bg-(--vscode-list-hoverBackground) focus:bg-(--vscode-list-focusBackground)"
					>
						{option.label}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>
</div>
