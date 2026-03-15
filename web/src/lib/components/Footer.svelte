<script>
export let onLegalClick;
</script>

<footer class="glass-regular mt-auto py-6 px-4">
	<div class="max-w-4xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
		<div class="text-center sm:text-left">
			<h4 class="text-sm font-bold text-foreground">Recipe</h4>
			<p class="text-xs text-foreground-muted">Pro Preset Converter</p>
		</div>

		<div class="glass-clear px-3 py-1.5 rounded-full">
			<span class="text-xs text-foreground-muted">
				Part of the <a href="https://shuttercoach.app" target="_blank" rel="noopener noreferrer" class="underline hover:text-interactive transition-colors">ShutterCoach</a> ecosystem
			</span>
		</div>

		<div class="flex items-center gap-4 text-xs">
			<a href="https://github.com/justin/recipe" target="_blank" rel="noopener noreferrer" class="text-foreground-muted hover:text-interactive transition-colors">GitHub</a>
			<a href="docs/privacy.md" class="text-foreground-muted hover:text-interactive transition-colors">Privacy</a>
			<!-- svelte-ignore a11y_invalid_attribute -->
			<a href="#" on:click|preventDefault={onLegalClick} class="text-foreground-muted hover:text-interactive transition-colors">Legal</a>
		</div>
	</div>
</footer>
