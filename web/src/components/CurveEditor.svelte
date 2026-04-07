<script lang="ts">
	import type { ToneCurvePoint } from "@recipe/ui";

	interface Props {
		points: ToneCurvePoint[];
		color?: string;
		disabled?: boolean;
		onchange?: (points: ToneCurvePoint[]) => void;
	}

	let {
		points,
		color = "rgba(255,255,255,0.9)",
		disabled = false,
		onchange,
	}: Props = $props();

	const PAD = 16;
	const SZ = 220; // logical curve area (0–255 maps to 0–SZ pixels)
	const TOTAL = SZ + 2 * PAD;
	const R = 6; // handle radius

	let canvas: HTMLCanvasElement;
	let dragging = -1;
	// Pending drag state: drawn immediately but not committed to history until pointerUp
	let pendingPoints: ToneCurvePoint[] | null = null;

	// Ensure (0,0) and (255,255) anchors exist; return sorted copy
	function norm(pts: ToneCurvePoint[]): ToneCurvePoint[] {
		const r = pts.filter((p) => p.input >= 0 && p.input <= 255);
		if (!r.some((p) => p.input === 0)) r.push({ input: 0, output: 0 });
		if (!r.some((p) => p.input === 255)) r.push({ input: 255, output: 255 });
		return r.sort((a, b) => a.input - b.input);
	}

	// Strip identity anchors before emitting to parent — keeps store data minimal
	function emit(pts: ToneCurvePoint[]): ToneCurvePoint[] {
		return pts.filter(
			(p) => !(p.input === 0 && p.output === 0) && !(p.input === 255 && p.output === 255),
		);
	}

	// Curve-space ↔ canvas-pixel conversions
	function px(input: number): number { return PAD + (input / 255) * SZ; }
	function py(output: number): number { return PAD + ((255 - output) / 255) * SZ; }
	function ix(x: number): number { return Math.round(Math.max(0, Math.min(255, ((x - PAD) / SZ) * 255))); }
	function iy(y: number): number { return Math.round(Math.max(0, Math.min(255, (1 - (y - PAD) / SZ) * 255))); }

	// Fritsch–Carlson monotone cubic Hermite tangents
	function mkTangents(pts: ToneCurvePoint[]): number[] {
		const n = pts.length;
		if (n < 2) return new Array(n).fill(0);
		const d: number[] = [];
		for (let i = 0; i < n - 1; i++) {
			const dx = pts[i + 1].input - pts[i].input;
			d.push(dx === 0 ? 0 : (pts[i + 1].output - pts[i].output) / dx);
		}
		const m = new Array<number>(n).fill(0);
		m[0] = d[0];
		m[n - 1] = d[n - 2];
		for (let i = 1; i < n - 1; i++) {
			m[i] = d[i - 1] * d[i] <= 0 ? 0 : (d[i - 1] + d[i]) / 2;
		}
		for (let i = 0; i < n - 1; i++) {
			if (d[i] === 0) { m[i] = m[i + 1] = 0; continue; }
			const a = m[i] / d[i], b = m[i + 1] / d[i], sq = a * a + b * b;
			if (sq > 9) { const t = 3 / Math.sqrt(sq); m[i] = t * a * d[i]; m[i + 1] = t * b * d[i]; }
		}
		return m;
	}

	function evalSpline(pts: ToneCurvePoint[], m: number[], x: number): number {
		const n = pts.length;
		if (n === 0) return x;
		if (x <= pts[0].input) return pts[0].output;
		if (x >= pts[n - 1].input) return pts[n - 1].output;
		let lo = 0, hi = n - 1;
		while (hi - lo > 1) { const mid = (lo + hi) >> 1; if (pts[mid].input <= x) lo = mid; else hi = mid; }
		const h = pts[hi].input - pts[lo].input;
		if (h === 0) return pts[lo].output;
		const t = (x - pts[lo].input) / h, t2 = t * t, t3 = t2 * t;
		return (2 * t3 - 3 * t2 + 1) * pts[lo].output + (t3 - 2 * t2 + t) * h * m[lo]
			+ (-2 * t3 + 3 * t2) * pts[hi].output + (t3 - t2) * h * m[hi];
	}

	function drawWith(drawPoints: ToneCurvePoint[], strokeColor: string) {
		if (!canvas) return;
		const ctx = canvas.getContext("2d");
		if (!ctx) return;
		const dpr = window.devicePixelRatio || 1;
		// Only reallocate the GPU framebuffer when dimensions actually change
		const targetW = TOTAL * dpr;
		const targetH = TOTAL * dpr;
		if (canvas.width !== targetW || canvas.height !== targetH) {
			canvas.width = targetW;
			canvas.height = targetH;
			canvas.style.width = `${TOTAL}px`;
			canvas.style.height = `${TOTAL}px`;
		}
		ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
		ctx.clearRect(0, 0, TOTAL, TOTAL);

		// Background
		ctx.fillStyle = "rgba(0,0,0,0.4)";
		ctx.fillRect(0, 0, TOTAL, TOTAL);

		// Chart area
		ctx.fillStyle = "rgba(255,255,255,0.02)";
		ctx.fillRect(PAD, PAD, SZ, SZ);
		ctx.strokeStyle = "rgba(255,255,255,0.1)";
		ctx.lineWidth = 0.5;
		ctx.strokeRect(PAD, PAD, SZ, SZ);

		// Grid lines at 25 / 50 / 75 %
		ctx.strokeStyle = "rgba(255,255,255,0.06)";
		ctx.lineWidth = 0.5;
		for (const t of [0.25, 0.5, 0.75]) {
			ctx.beginPath(); ctx.moveTo(PAD + t * SZ, PAD); ctx.lineTo(PAD + t * SZ, PAD + SZ); ctx.stroke();
			ctx.beginPath(); ctx.moveTo(PAD, PAD + t * SZ); ctx.lineTo(PAD + SZ, PAD + t * SZ); ctx.stroke();
		}

		// Identity diagonal (dashed)
		ctx.strokeStyle = "rgba(255,255,255,0.12)";
		ctx.lineWidth = 1;
		ctx.setLineDash([4, 4]);
		ctx.beginPath(); ctx.moveTo(PAD, PAD + SZ); ctx.lineTo(PAD + SZ, PAD); ctx.stroke();
		ctx.setLineDash([]);

		const pts = norm(drawPoints);
		const m = mkTangents(pts);

		// Spline
		ctx.strokeStyle = strokeColor;
		ctx.lineWidth = 1.5;
		ctx.beginPath();
		for (let i = 0; i <= SZ; i++) {
			const inputVal = (i / SZ) * 255;
			const outputVal = Math.max(0, Math.min(255, evalSpline(pts, m, inputVal)));
			if (i === 0) ctx.moveTo(px(inputVal), py(outputVal));
			else ctx.lineTo(px(inputVal), py(outputVal));
		}
		ctx.stroke();

		// Handles
		for (const pt of pts) {
			const isAnchor = pt.input === 0 || pt.input === 255;
			ctx.beginPath();
			ctx.arc(px(pt.input), py(pt.output), R, 0, Math.PI * 2);
			ctx.fillStyle = isAnchor ? "rgba(255,255,255,0.15)" : strokeColor;
			ctx.fill();
			ctx.strokeStyle = "rgba(255,255,255,0.65)";
			ctx.lineWidth = 1.5;
			ctx.stroke();
		}
	}

	$effect(() => {
		// When not dragging, redraw from prop; during drag, drawWith is called directly
		if (dragging === -1) drawWith(points, color);
	});

	function hit(cx: number, cy: number): number {
		const pts = norm(points);
		for (let i = 0; i < pts.length; i++) {
			if (Math.hypot(cx - px(pts[i].input), cy - py(pts[i].output)) <= R + 4) return i;
		}
		return -1;
	}

	function pointerDown(e: PointerEvent) {
		if (disabled) return;
		e.preventDefault();
		const rect = canvas.getBoundingClientRect();
		const cx = e.clientX - rect.left, cy = e.clientY - rect.top;
		const idx = hit(cx, cy);
		if (idx !== -1) {
			dragging = idx;
			canvas.setPointerCapture(e.pointerId);
		} else {
			// Add new point — instantaneous, commit immediately (not a drag)
			const input = ix(cx), output = iy(cy);
			if (input === 0 || input === 255) return;
			const updated = norm([...points, { input, output }]);
			drawWith(updated, color);
			onchange?.(emit(updated));
		}
	}

	function pointerMove(e: PointerEvent) {
		if (disabled || dragging === -1) return;
		const rect = canvas.getBoundingClientRect();
		const cx = e.clientX - rect.left, cy = e.clientY - rect.top;
		const pts = norm(points);
		const pt = pts[dragging];
		if (!pt) return;
		const output = iy(cy);
		let updated: ToneCurvePoint[];
		if (pt.input === 0 || pt.input === 255) {
			// Anchors: vertical movement only
			pts[dragging] = { input: pt.input, output };
			updated = pts;
		} else {
			const minI = dragging > 0 ? pts[dragging - 1].input + 1 : 1;
			const maxI = dragging < pts.length - 1 ? pts[dragging + 1].input - 1 : 254;
			pts[dragging] = { input: Math.max(minI, Math.min(maxI, ix(cx))), output };
			updated = pts;
		}
		// Draw immediately for smooth feedback; defer store commit to pointerUp
		drawWith(updated, color);
		pendingPoints = updated;
	}

	function pointerUp() {
		// Commit a single undo history entry for the entire drag gesture
		if (pendingPoints !== null) {
			onchange?.(emit(pendingPoints));
			pendingPoints = null;
		}
		dragging = -1;
	}

	function pointerCancel() {
		// Revert visual to last committed state
		pendingPoints = null;
		dragging = -1;
		drawWith(points, color);
	}

	function contextMenu(e: MouseEvent) {
		if (disabled) return;
		e.preventDefault();
		const rect = canvas.getBoundingClientRect();
		const idx = hit(e.clientX - rect.left, e.clientY - rect.top);
		if (idx === -1) return;
		const pts = norm(points);
		const pt = pts[idx];
		if (pt.input === 0 || pt.input === 255) return;
		pts.splice(idx, 1);
		onchange?.(emit(pts));
	}
</script>

<canvas
	bind:this={canvas}
	style="display:block; cursor:{disabled ? 'default' : 'crosshair'};"
	class="rounded select-none"
	onpointerdown={pointerDown}
	onpointermove={pointerMove}
	onpointerup={pointerUp}
	onpointercancel={pointerCancel}
	oncontextmenu={contextMenu}
	role="img"
	aria-label="Tone curve editor — drag to move points, click to add, right-click to remove"
></canvas>
