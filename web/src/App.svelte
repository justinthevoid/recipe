<script>
import { onMount } from "svelte";
import { initializeWASM } from "./lib/wasm";
import Hero from "./lib/components/Hero.svelte";
import FormatGrid from "./lib/components/FormatGrid.svelte";
import HowItWorks from "./lib/components/HowItWorks.svelte";
import FAQ from "./lib/components/FAQ.svelte";
import Footer from "./lib/components/Footer.svelte";
import Modal from "./lib/components/Modal.svelte";
import UploadZone from "./lib/components/UploadZone.svelte";
import FileList from "./lib/components/FileList.svelte";
import ActionPanel from "./lib/components/ActionPanel.svelte";
import PreviewModal from "./lib/components/PreviewModal.svelte";

let showLegalModal = false;

onMount(() => {
	initializeWASM();
});
</script>

<main class="app-shell">
  <Hero />

  <!-- Main Content -->
  <div class="container">
    <div class="glass-card-wrapper">
      <UploadZone />
      <FileList />
      <ActionPanel />
    </div>
  </div>

  <FormatGrid />
  <HowItWorks />
  <FAQ />
  <Footer onLegalClick={() => (showLegalModal = true)} />

  <Modal isOpen={showLegalModal} onClose={() => (showLegalModal = false)}>
    <h2>Legal Disclaimer</h2>
    <p><strong>Recipe is provided "AS IS" without warranty.</strong></p>
    <p>
      Files are processed locally via WebAssembly. No server uploads. No data
      stored.
    </p>
    <hr style="border-color: var(--glass-border); margin: 1rem 0;" />
    <h3>Reverse Engineering</h3>
    <p>
      Recipe uses reverse-engineered file formats for interoperability. Not
      affiliated with Nikon, Adobe, or Phase One.
    </p>
  </Modal>

  <PreviewModal />
</main>

<style>
  /* Component-specific styles if needed */
</style>
