<script lang="ts">
	let {
		title,
		expanded = false,
		children,
	}: {
		title: string;
		expanded?: boolean;
		children: import("svelte").Snippet;
	} = $props();

	let isExpanded = $state(expanded);

	function _handleToggle() {
		isExpanded = !isExpanded;
	}
</script>

<section class="flex flex-col">
	<button
		type="button"
		class="flex items-center gap-2 border-l-2 border-(--vscode-button-background) pl-3 py-2 hover:bg-(--vscode-list-hoverBackground) rounded-r transition-colors cursor-pointer w-full text-left"
		onclick={_handleToggle}
		aria-expanded={isExpanded}
	>
		<span class="text-xs opacity-60 w-3">{isExpanded ? "▼" : "▶"}</span>
		<h2 class="text-xs font-bold uppercase tracking-widest opacity-70">{title}</h2>
	</button>

	{#if isExpanded}
		<div class="bg-muted/10 p-5 rounded-xl border border-border/50 backdrop-blur-sm mt-2">
			{@render children()}
		</div>
	{/if}
</section>
