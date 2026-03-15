<script>
import { files } from "../stores-compat";
import FileCard from "./FileCard.svelte";
</script>

{#if $files.length > 0}
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
		{#each $files as file (file.id)}
			<FileCard {file} />
		{/each}
	</div>
{/if}
