<script>
let items = [
	{
		question: "Is my data private?",
		answer:
			"Yes, 100%. Recipe processes all files locally in your browser using WebAssembly. Your files never leave your device, and we don't collect any analytics or data.",
		isOpen: false,
	},
	{
		question: "Is Recipe legal?",
		answer:
			"Recipe uses reverse-engineered file formats for interoperability purposes, which is generally protected under fair use. We recommend private use only. Not affiliated with Nikon or Adobe.",
		isOpen: false,
	},
	{
		question: "Why don't some parameters convert?",
		answer:
			"Different formats support different features. For example, Lightroom's Grain effect doesn't exist in Nikon's NP3 format. Recipe approximates where possible.",
		isOpen: false,
	},
];

function toggle(index) {
	items[index].isOpen = !items[index].isOpen;
}
</script>

<section class="section">
    <div class="container">
        <div class="section-header">
            <h2 class="section-title">FAQ & Privacy</h2>
        </div>

        <div class="faq-grid">
            {#each items as item, i}
                <div class="faq-item" class:active={item.isOpen}>
                    <div
                        class="faq-question"
                        on:click={() => toggle(i)}
                        role="button"
                        tabindex="0"
                        on:keydown={(e) => e.key === "Enter" && toggle(i)}
                    >
                        {item.question}
                        <span>{item.isOpen ? "−" : "+"}</span>
                    </div>
                    {#if item.isOpen}
                        <div class="faq-answer" style="display: block;">
                            {@html item.answer}
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
</section>
