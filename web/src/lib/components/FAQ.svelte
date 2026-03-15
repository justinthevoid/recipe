<script>
let items = [
	{
		question: "Is my data private?",
		answer: "Yes, 100%. Recipe processes all files locally in your browser using WebAssembly. Your files never leave your device, and we don't collect any analytics or data.",
		isOpen: false,
	},
	{
		question: "Is Recipe legal?",
		answer: "Recipe uses reverse-engineered file formats for interoperability purposes, which is generally protected under fair use. We recommend private use only. Not affiliated with Nikon or Adobe.",
		isOpen: false,
	},
	{
		question: "Why don't some parameters convert?",
		answer: "NP3 and XMP support different features. For example, Lightroom's Grain effect doesn't exist in Nikon's NP3 format. Recipe preserves what it can and approximates the rest.",
		isOpen: false,
	},
];

function toggle(index) {
	items[index].isOpen = !items[index].isOpen;
}
</script>

<section class="py-20 px-4">
	<div class="max-w-2xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-2xl font-bold text-foreground">FAQ & Privacy</h2>
		</div>

		<div class="space-y-3">
			{#each items as item, i}
				<div class="glass-regular rounded-xl overflow-hidden transition-all">
					<button
						class="w-full text-left px-5 py-4 flex items-center justify-between cursor-pointer hover:bg-surface-hover transition-colors"
						on:click={() => toggle(i)}
					>
						<span class="text-sm font-medium text-foreground">{item.question}</span>
						<span class="text-foreground-muted text-lg ml-4 transition-transform {item.isOpen ? 'rotate-45' : ''}">{item.isOpen ? "−" : "+"}</span>
					</button>
					{#if item.isOpen}
						<div class="px-5 pb-4 text-sm text-foreground-muted leading-relaxed">
							{item.answer}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</section>
