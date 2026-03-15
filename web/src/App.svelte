<script lang="ts">
import { onMount } from "svelte";
import ActionPanel from "./lib/components/ActionPanel.svelte";
import AuroraBackground from "./lib/components/AuroraBackground.svelte";
import EditorView from "./lib/components/EditorView.svelte";
import FAQ from "./lib/components/FAQ.svelte";
import FileList from "./lib/components/FileList.svelte";
import Footer from "./lib/components/Footer.svelte";
import FormatGrid from "./lib/components/FormatGrid.svelte";
import Hero from "./lib/components/Hero.svelte";
import HowItWorks from "./lib/components/HowItWorks.svelte";
import Modal from "./lib/components/Modal.svelte";
import UploadZone from "./lib/components/UploadZone.svelte";
import { initWasm, wasm, wasmGenerateLUT } from "./lib/wasm.svelte";
import { store } from "./lib/stores.svelte";

let showLegalModal = $state(false);

onMount(() => {
	initWasm();
});
</script>

<AuroraBackground />

{#if store.editorMode}
	<EditorView
		wasmReady={wasm.ready}
		generateLUT={wasmGenerateLUT}
	/>
{:else}
	<main class="min-h-screen flex flex-col text-foreground font-sans relative">
		<Hero />

		<div class="max-w-4xl mx-auto w-full px-4 py-8">
			<UploadZone />
			<FileList />
			<ActionPanel />
		</div>

		<FormatGrid />
		<HowItWorks />
		<FAQ />
		<Footer onLegalClick={() => (showLegalModal = true)} />

		<Modal isOpen={showLegalModal} onClose={() => (showLegalModal = false)}>
			<h2 class="text-lg font-bold text-foreground mb-2">Legal Disclaimer</h2>
			<p class="text-sm text-foreground-muted"><strong class="text-foreground">Recipe is provided "AS IS" without warranty.</strong></p>
			<p class="text-sm text-foreground-muted mt-2">
				Files are processed locally via WebAssembly. No server uploads. No data stored.
			</p>
			<hr class="border-border my-4" />
			<h3 class="text-sm font-bold text-foreground mb-1">Reverse Engineering</h3>
			<p class="text-sm text-foreground-muted">
				Recipe uses reverse-engineered file formats for interoperability. Not
				affiliated with Nikon, Adobe, or Phase One.
			</p>
		</Modal>
	</main>
{/if}
