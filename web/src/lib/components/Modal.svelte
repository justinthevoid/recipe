<script>
export let isOpen = false;
export let onClose;
</script>

{#if isOpen}
	<!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events a11y_interactive_supports_focus -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
		on:click={onClose}
		on:keydown={(e) => e.key === 'Escape' && onClose()}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<!-- svelte-ignore a11y_no_noninteractive_element_interactions a11y_click_events_have_key_events -->
		<div
			class="glass-regular rounded-2xl p-6 max-w-lg mx-4 text-foreground relative"
			on:click|stopPropagation
			role="document"
		>
			<button
				class="absolute top-3 right-3 w-8 h-8 flex items-center justify-center rounded-lg text-foreground-muted hover:text-foreground hover:bg-surface-hover transition-colors text-xl"
				on:click={onClose}
			>
				&times;
			</button>
			<slot></slot>
		</div>
	</div>
{/if}
