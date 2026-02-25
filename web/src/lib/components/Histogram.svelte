<script>
import { onMount, afterUpdate } from "svelte";

export let data = { r: [], g: [], b: [] };
export let height = 100;

let canvas;
let ctx;

$: if (canvas && data) {
	drawHistogram();
}

function drawHistogram() {
	if (!canvas || !data.r.length) return;

	const ctx = canvas.getContext("2d");
	const w = canvas.width;
	const h = canvas.height;

	ctx.clearRect(0, 0, w, h);

	// Find max value to normalize height
	const maxVal = Math.max(...data.r, ...data.g, ...data.b);

	if (maxVal === 0) return;

	// Draw Channel helper
	const drawChannel = (hist, color) => {
		ctx.fillStyle = color;
		ctx.beginPath();
		ctx.moveTo(0, h);

		for (let i = 0; i < 256; i++) {
			const x = (i / 255) * w;
			const val = hist[i];
			const barHeight = (val / maxVal) * h;
			const y = h - barHeight;
			ctx.lineTo(x, y);
		}

		ctx.lineTo(w, h);
		ctx.closePath();
		ctx.fill();
	};

	// Use additive blending for RGB overlap
	ctx.globalCompositeOperation = "screen";

	drawChannel(data.r, "rgba(255, 50, 50, 0.8)");
	drawChannel(data.g, "rgba(50, 255, 50, 0.8)");
	drawChannel(data.b, "rgba(50, 50, 255, 0.8)");

	ctx.globalCompositeOperation = "source-over";
}

onMount(() => {
	drawHistogram();
});
</script>

<div class="histogram-container" style="height: {height}px">
    <canvas bind:this={canvas} width="300" {height}></canvas>
</div>

<style>
    .histogram-container {
        width: 100%;
        background: rgba(0, 0, 0, 0.3);
        border-radius: 4px;
        overflow: hidden;
        margin-bottom: 1rem;
    }

    canvas {
        width: 100%;
        height: 100%;
        display: block;
    }
</style>
