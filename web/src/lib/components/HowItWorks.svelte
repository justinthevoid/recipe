<section id="how-it-works" class="py-20 px-4">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-2xl font-bold text-foreground">How It Works</h2>
			<p class="mt-2 text-foreground-muted">Simple, fast, and secure conversion workflow.</p>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			{#each [
				{ step: "1", title: "Upload", desc: "Drag & drop your preset files or click to browse. We support single files or bulk upload." },
				{ step: "2", title: "Convert", desc: "Select your target format. Our WASM engine processes everything locally on your device." },
				{ step: "3", title: "Download", desc: "Get your converted presets instantly. No server uploads, no waiting, no data usage." },
			] as item}
				<div class="glass-regular rounded-xl p-6 text-center transition-all hover:-translate-y-1">
					<div class="inline-flex items-center justify-center w-10 h-10 rounded-full mb-4 font-bold text-lg bg-gradient-to-r from-aurora-cyan to-aurora-violet bg-clip-text text-transparent border border-border">
						{item.step}
					</div>
					<h3 class="text-lg font-semibold text-foreground mb-2">{item.title}</h3>
					<p class="text-sm text-foreground-muted">{item.desc}</p>
				</div>
			{/each}
		</div>
	</div>
</section>
