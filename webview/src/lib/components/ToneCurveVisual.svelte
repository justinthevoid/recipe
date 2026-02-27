<script lang="ts">
	let {
		points = [],
		onchange
	}: {
		points: { input: number; output: number }[];
		onchange?: (points: { input: number; output: number }[]) => void;
	} = $props();

	// Ensure points are sorted by input
	let sortedPoints = $derived([...points].sort((a, b) => a.input - b.input));

	// Fallback to linear if no points
	let displayPoints = $derived(
		sortedPoints.length > 0 
			? sortedPoints 
			: [{ input: 0, output: 0 }, { input: 255, output: 255 }]
	);

	// Generate SVG path data (inverted Y)
	let pathData = $derived(
		displayPoints.length > 0
			? `M ${displayPoints[0].input} ${255 - displayPoints[0].output} ` +
			  displayPoints.slice(1).map(p => `L ${p.input} ${255 - p.output}`).join(' ')
			: ""
	);
</script>

<div class="relative w-full aspect-square bg-(--vscode-editor-background) border border-border rounded overflow-hidden select-none max-w-[400px] mx-auto">
	<!-- Grid Background -->
	<div class="absolute inset-0 pointer-events-none opacity-20">
		<svg width="100%" height="100%" viewBox="0 0 255 255" preserveAspectRatio="none">
			<title>Tone Curve Grid</title>
			<defs>
				<pattern id="grid" width="63.75" height="63.75" patternUnits="userSpaceOnUse">
					<path d="M 63.75 0 L 0 0 0 63.75" fill="none" stroke="currentColor" stroke-width="0.5"/>
				</pattern>
			</defs>
			<rect width="100%" height="100%" fill="url(#grid)" />
			<line x1="0" y1="0" x2="255" y2="255" stroke="currentColor" stroke-width="0.5" stroke-dasharray="2,2" />
		</svg>
	</div>

	<!-- Curve SVG -->
	<svg 
		width="100%" 
		height="100%" 
		viewBox="0 0 255 255" 
		preserveAspectRatio="none"
		class="relative"
	>
		<title>Tone Curve</title>
		<!-- Active Path -->
		<path 
			d={pathData} 
			fill="none" 
			stroke="var(--vscode-button-background)" 
			stroke-width="2" 
			vector-effect="non-scaling-stroke"
		/>

		<!-- Control Points -->
		{#each displayPoints as point}
			<circle 
				cx={point.input} 
				cy={255 - point.output} 
				r="4" 
				class="fill-(--vscode-button-background)"
			/>
		{/each}
	</svg>

	<!-- Axes Labels -->
	<div class="absolute bottom-1 left-1 text-[8px] opacity-50 uppercase tracking-tighter">Shadows</div>
	<div class="absolute top-1 right-1 text-[8px] opacity-50 uppercase tracking-tighter">Highlights</div>
</div>
