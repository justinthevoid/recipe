<script>
import { files } from "../stores";
import FileCard from "./FileCard.svelte";
</script>

{#if $files.length > 0}
    <div class="file-list">
        {#each $files as file (file.id)}
            <FileCard {file} />
        {/each}
    </div>
{/if}
