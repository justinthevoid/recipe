<section class="py-20 px-4">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-2xl font-bold text-foreground">Supported Formats</h2>
			<p class="mt-2 text-foreground-muted">Bidirectional conversion between Nikon NP3 and Adobe Lightroom XMP.</p>
		</div>

		<div class="grid grid-cols-2 gap-4 max-w-md mx-auto">
			{#each [
				{ name: "NP3", desc: "Nikon Picture Control", color: "var(--color-aurora-violet)" },
				{ name: "XMP", desc: "Adobe Lightroom CC", color: "var(--color-aurora-cyan)" },
			] as format}
				<div class="glass-regular rounded-xl p-5 text-center transition-all hover:-translate-y-1 hover:shadow-lg">
					<div class="w-3 h-3 rounded-full mx-auto mb-3" style="background: {format.color}; box-shadow: 0 0 15px {format.color};"></div>
					<div class="text-sm font-bold text-foreground">{format.name}</div>
					<div class="text-xs text-foreground-muted mt-1">{format.desc}</div>
				</div>
			{/each}
		</div>
	</div>
</section>
