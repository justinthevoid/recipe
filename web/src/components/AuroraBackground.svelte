<div class="aurora-container" aria-hidden="true">
	<!-- Static canvas layer -->
	<div class="aurora-layer aurora-static"></div>

	<!-- Slow breathing (30s) — all viewports -->
	<div class="aurora-layer aurora-slow"></div>

	<!-- Medium breathing (15s) — tablet+ -->
	<div class="aurora-layer aurora-medium"></div>

	<!-- Fast breathing (8s) — desktop+ -->
	<div class="aurora-layer aurora-fast"></div>

	<!-- Film grain overlay -->
	<svg class="aurora-grain" xmlns="http://www.w3.org/2000/svg">
		<filter id="grain">
			<feTurbulence type="fractalNoise" baseFrequency="0.65" numOctaves="3" stitchTiles="stitch" />
		</filter>
		<rect width="100%" height="100%" filter="url(#grain)" opacity="0.04" />
	</svg>
</div>

<style>
	.aurora-container {
		position: fixed;
		inset: -10%;
		z-index: -1;
		pointer-events: none;
		overflow: hidden;
	}

	.aurora-layer {
		position: absolute;
		inset: 0;
	}

	.aurora-static {
		background:
			radial-gradient(ellipse 80% 50% at 20% 40%, var(--color-aurora-deep-blue) 0%, transparent 70%),
			radial-gradient(ellipse 60% 40% at 80% 20%, var(--color-aurora-violet) 0%, transparent 60%),
			radial-gradient(ellipse 50% 60% at 50% 80%, var(--color-aurora-cyan) 0%, transparent 65%),
			radial-gradient(ellipse 70% 45% at 70% 60%, var(--color-aurora-purple) 0%, transparent 55%);
		opacity: 0.28;
	}

	.aurora-slow {
		background: radial-gradient(ellipse 90% 60% at 40% 50%, var(--color-aurora-soft-blue) 0%, transparent 70%);
		opacity: 0.15;
		animation: aurora-breathe-slow 30s ease-in-out infinite;
	}

	.aurora-medium {
		background: radial-gradient(ellipse 70% 50% at 60% 30%, var(--color-aurora-soft-pink) 0%, transparent 65%);
		opacity: 0.18;
		animation: aurora-breathe-medium 15s ease-in-out infinite;
		display: none;
	}

	.aurora-fast {
		background: radial-gradient(ellipse 50% 40% at 30% 70%, var(--color-aurora-cyan) 0%, transparent 60%);
		opacity: 0.14;
		animation: aurora-breathe-fast 8s ease-in-out infinite;
		display: none;
	}

	.aurora-grain {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
	}

	@keyframes aurora-breathe-slow {
		0%, 100% { transform: scale(1) translateY(0); opacity: 0.15; }
		50% { transform: scale(1.05) translateY(-2%); opacity: 0.2; }
	}

	@keyframes aurora-breathe-medium {
		0%, 100% { transform: scale(1) translateX(0); opacity: 0.18; }
		50% { transform: scale(1.08) translateX(3%); opacity: 0.22; }
	}

	@keyframes aurora-breathe-fast {
		0%, 100% { transform: scale(1) rotate(0deg); opacity: 0.14; }
		50% { transform: scale(1.1) rotate(1deg); opacity: 0.18; }
	}

	@media (min-width: 768px) {
		.aurora-medium { display: block; }
	}

	@media (min-width: 1024px) {
		.aurora-fast { display: block; }
	}

	@media (prefers-reduced-motion: reduce) {
		.aurora-slow,
		.aurora-medium,
		.aurora-fast {
			animation: none;
		}
	}

	@media (prefers-reduced-transparency: reduce) {
		.aurora-container {
			background: var(--color-canvas-base);
		}
		.aurora-layer,
		.aurora-grain {
			display: none;
		}
	}
</style>
