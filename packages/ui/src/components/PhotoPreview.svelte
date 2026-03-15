<script lang="ts">
	import type { UniversalRecipe } from "../types";

	const VERTEX_SHADER = `#version 300 es
		layout(location = 0) in vec2 a_position;
		layout(location = 1) in vec2 a_texCoord;
		out vec2 v_texCoord;
		void main() {
			gl_Position = vec4(a_position, 0.0, 1.0);
			v_texCoord = a_texCoord;
		}
	`;

	const LUT_FRAGMENT_SHADER = `#version 300 es
		precision highp float;
		uniform sampler2D u_image;
		uniform highp sampler3D u_lut;
		uniform float u_lutSize;
		in vec2 v_texCoord;
		out vec4 fragColor;
		void main() {
			vec4 color = texture(u_image, v_texCoord);
			// Half-texel offset for correct trilinear interpolation
			vec3 lutCoord = (color.rgb * (u_lutSize - 1.0) + 0.5) / u_lutSize;
			vec3 mapped = texture(u_lut, lutCoord).rgb;
			fragColor = vec4(mapped, color.a);
		}
	`;

	const SHARPEN_FRAGMENT_SHADER = `#version 300 es
		precision highp float;
		uniform sampler2D u_texture;
		uniform float u_sharpness;
		uniform vec2 u_texelSize;
		in vec2 v_texCoord;
		out vec4 fragColor;
		void main() {
			vec3 center = texture(u_texture, v_texCoord).rgb;

			vec3 blurred =
				texture(u_texture, v_texCoord + vec2(-1.0, -1.0) * u_texelSize).rgb * (1.0/16.0) +
				texture(u_texture, v_texCoord + vec2( 0.0, -1.0) * u_texelSize).rgb * (2.0/16.0) +
				texture(u_texture, v_texCoord + vec2( 1.0, -1.0) * u_texelSize).rgb * (1.0/16.0) +
				texture(u_texture, v_texCoord + vec2(-1.0,  0.0) * u_texelSize).rgb * (2.0/16.0) +
				center * (4.0/16.0) +
				texture(u_texture, v_texCoord + vec2( 1.0,  0.0) * u_texelSize).rgb * (2.0/16.0) +
				texture(u_texture, v_texCoord + vec2(-1.0,  1.0) * u_texelSize).rgb * (1.0/16.0) +
				texture(u_texture, v_texCoord + vec2( 0.0,  1.0) * u_texelSize).rgb * (2.0/16.0) +
				texture(u_texture, v_texCoord + vec2( 1.0,  1.0) * u_texelSize).rgb * (1.0/16.0);

			vec3 result = center + u_sharpness * (center - blurred);
			fragColor = vec4(clamp(result, 0.0, 1.0), 1.0);
		}
	`;

	const LUT_SIZE = 17;

	let {
		recipe = null,
		imageData = null,
		wasmReady = false,
		generateLUT,
	}: {
		recipe: UniversalRecipe | null;
		imageData: HTMLImageElement | null;
		wasmReady: boolean;
		generateLUT: (recipeJSON: string, size: number) => Promise<Float32Array>;
	} = $props();

	let canvas: HTMLCanvasElement | undefined = $state();
	let gl: WebGL2RenderingContext | null = $state(null);
	let contextLost = $state(false);
	let webgl2Unsupported = $state(false);

	// WebGL resources
	let lutProgram: WebGLProgram | null = null;
	let sharpenProgram: WebGLProgram | null = null;
	let imageTexture: WebGLTexture | null = null;
	let lutTexture: WebGLTexture | null = null;
	let fbo: WebGLFramebuffer | null = null;
	let fboTexture: WebGLTexture | null = null;
	let vao: WebGLVertexArrayObject | null = null;
	let hasFloatLinear = false;

	// Track last uploaded image for ResizeObserver
	let lastImage: HTMLImageElement | null = null;

	// Debounce state (non-reactive)
	let generationCounter = 0;
	let pendingGeneration = false;
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;

	function compileShader(glCtx: WebGL2RenderingContext, type: number, source: string): WebGLShader | null {
		const shader = glCtx.createShader(type);
		if (!shader) return null;
		glCtx.shaderSource(shader, source);
		glCtx.compileShader(shader);
		if (!glCtx.getShaderParameter(shader, glCtx.COMPILE_STATUS)) {
			console.error("Shader compile error:", glCtx.getShaderInfoLog(shader));
			glCtx.deleteShader(shader);
			return null;
		}
		return shader;
	}

	function createProgram(glCtx: WebGL2RenderingContext, vertSrc: string, fragSrc: string): WebGLProgram | null {
		const vert = compileShader(glCtx, glCtx.VERTEX_SHADER, vertSrc);
		const frag = compileShader(glCtx, glCtx.FRAGMENT_SHADER, fragSrc);
		if (!vert || !frag) return null;

		const prog = glCtx.createProgram();
		if (!prog) return null;
		glCtx.attachShader(prog, vert);
		glCtx.attachShader(prog, frag);
		glCtx.linkProgram(prog);

		if (!glCtx.getProgramParameter(prog, glCtx.LINK_STATUS)) {
			console.error("Program link error:", glCtx.getProgramInfoLog(prog));
			glCtx.deleteProgram(prog);
			return null;
		}

		glCtx.deleteShader(vert);
		glCtx.deleteShader(frag);
		return prog;
	}

	function initWebGL(glCtx: WebGL2RenderingContext) {
		hasFloatLinear = !!glCtx.getExtension("OES_texture_float_linear");
		if (!hasFloatLinear) {
			console.warn("OES_texture_float_linear not available — falling back to NEAREST filtering for LUT");
		}

		lutProgram = createProgram(glCtx, VERTEX_SHADER, LUT_FRAGMENT_SHADER);
		sharpenProgram = createProgram(glCtx, VERTEX_SHADER, SHARPEN_FRAGMENT_SHADER);

		if (!lutProgram || !sharpenProgram) {
			console.error("Failed to create WebGL programs");
			return;
		}

		const positions = new Float32Array([
			-1, -1,  0, 1,
			 1, -1,  1, 1,
			-1,  1,  0, 0,
			 1,  1,  1, 0,
		]);

		vao = glCtx.createVertexArray();
		glCtx.bindVertexArray(vao);

		const buf = glCtx.createBuffer();
		glCtx.bindBuffer(glCtx.ARRAY_BUFFER, buf);
		glCtx.bufferData(glCtx.ARRAY_BUFFER, positions, glCtx.STATIC_DRAW);

		glCtx.enableVertexAttribArray(0);
		glCtx.vertexAttribPointer(0, 2, glCtx.FLOAT, false, 16, 0);

		glCtx.enableVertexAttribArray(1);
		glCtx.vertexAttribPointer(1, 2, glCtx.FLOAT, false, 16, 8);

		glCtx.bindVertexArray(null);

		lutTexture = glCtx.createTexture();
		glCtx.bindTexture(glCtx.TEXTURE_3D, lutTexture);
		glCtx.texParameteri(glCtx.TEXTURE_3D, glCtx.TEXTURE_WRAP_S, glCtx.CLAMP_TO_EDGE);
		glCtx.texParameteri(glCtx.TEXTURE_3D, glCtx.TEXTURE_WRAP_T, glCtx.CLAMP_TO_EDGE);
		glCtx.texParameteri(glCtx.TEXTURE_3D, glCtx.TEXTURE_WRAP_R, glCtx.CLAMP_TO_EDGE);
		const filterMode = hasFloatLinear ? glCtx.LINEAR : glCtx.NEAREST;
		glCtx.texParameteri(glCtx.TEXTURE_3D, glCtx.TEXTURE_MIN_FILTER, filterMode);
		glCtx.texParameteri(glCtx.TEXTURE_3D, glCtx.TEXTURE_MAG_FILTER, filterMode);

		imageTexture = glCtx.createTexture();

		fbo = glCtx.createFramebuffer();
		fboTexture = glCtx.createTexture();

		// Upload identity LUT so image renders immediately without WASM
		uploadIdentityLUT(glCtx);
	}

	function uploadIdentityLUT(glCtx: WebGL2RenderingContext) {
		const size = LUT_SIZE;
		const data = new Float32Array(size * size * size * 4);
		for (let b = 0; b < size; b++) {
			for (let g = 0; g < size; g++) {
				for (let r = 0; r < size; r++) {
					const idx = (b * size * size + g * size + r) * 4;
					data[idx] = r / (size - 1);
					data[idx + 1] = g / (size - 1);
					data[idx + 2] = b / (size - 1);
					data[idx + 3] = 1.0;
				}
			}
		}
		uploadLUT(glCtx, data);
	}

	function uploadImage(glCtx: WebGL2RenderingContext, img: HTMLImageElement) {
		lastImage = img;
		glCtx.bindTexture(glCtx.TEXTURE_2D, imageTexture);
		glCtx.texParameteri(glCtx.TEXTURE_2D, glCtx.TEXTURE_WRAP_S, glCtx.CLAMP_TO_EDGE);
		glCtx.texParameteri(glCtx.TEXTURE_2D, glCtx.TEXTURE_WRAP_T, glCtx.CLAMP_TO_EDGE);
		glCtx.texParameteri(glCtx.TEXTURE_2D, glCtx.TEXTURE_MIN_FILTER, glCtx.LINEAR);
		glCtx.texParameteri(glCtx.TEXTURE_2D, glCtx.TEXTURE_MAG_FILTER, glCtx.LINEAR);
		glCtx.texImage2D(glCtx.TEXTURE_2D, 0, glCtx.RGBA, glCtx.RGBA, glCtx.UNSIGNED_BYTE, img);

		updateCanvasSize(glCtx, img);
		setupFboTexture(glCtx);

		// Render immediately with current LUT (identity if WASM not ready)
		render(glCtx, 0);
	}

	function updateCanvasSize(glCtx: WebGL2RenderingContext, img: HTMLImageElement) {
		if (!canvas) return;
		const container = canvas.parentElement;
		if (!container) return;

		let containerWidth = container.clientWidth;
		let containerHeight = container.clientHeight;

		// Fallback: if container hasn't laid out yet, use parent of parent
		if (containerWidth < 10 || containerHeight < 10) {
			const outer = container.parentElement;
			if (outer) {
				containerWidth = outer.clientWidth;
				containerHeight = outer.clientHeight;
			}
		}
		if (containerWidth < 10 || containerHeight < 10) return;
		const imgAspect = img.naturalWidth / img.naturalHeight;
		const containerAspect = containerWidth / containerHeight;

		let drawWidth: number;
		let drawHeight: number;

		if (imgAspect > containerAspect) {
			drawWidth = containerWidth;
			drawHeight = containerWidth / imgAspect;
		} else {
			drawHeight = containerHeight;
			drawWidth = containerHeight * imgAspect;
		}

		canvas.width = Math.round(drawWidth * devicePixelRatio);
		canvas.height = Math.round(drawHeight * devicePixelRatio);
		canvas.style.width = `${Math.round(drawWidth)}px`;
		canvas.style.height = `${Math.round(drawHeight)}px`;

		glCtx.viewport(0, 0, canvas.width, canvas.height);
	}

	function setupFboTexture(glCtx: WebGL2RenderingContext) {
		if (!canvas) return;
		glCtx.bindTexture(glCtx.TEXTURE_2D, fboTexture);
		glCtx.texImage2D(
			glCtx.TEXTURE_2D, 0, glCtx.RGBA, canvas.width, canvas.height, 0,
			glCtx.RGBA, glCtx.UNSIGNED_BYTE, null,
		);
		glCtx.texParameteri(glCtx.TEXTURE_2D, glCtx.TEXTURE_MIN_FILTER, glCtx.LINEAR);
		glCtx.texParameteri(glCtx.TEXTURE_2D, glCtx.TEXTURE_MAG_FILTER, glCtx.LINEAR);
		glCtx.texParameteri(glCtx.TEXTURE_2D, glCtx.TEXTURE_WRAP_S, glCtx.CLAMP_TO_EDGE);
		glCtx.texParameteri(glCtx.TEXTURE_2D, glCtx.TEXTURE_WRAP_T, glCtx.CLAMP_TO_EDGE);
	}

	function uploadLUT(glCtx: WebGL2RenderingContext, data: Float32Array) {
		glCtx.bindTexture(glCtx.TEXTURE_3D, lutTexture);
		glCtx.texImage3D(
			glCtx.TEXTURE_3D, 0, glCtx.RGBA32F,
			LUT_SIZE, LUT_SIZE, LUT_SIZE, 0,
			glCtx.RGBA, glCtx.FLOAT, data,
		);
	}

	function render(glCtx: WebGL2RenderingContext, sharpness: number) {
		if (!lutProgram || !vao || !imageTexture || !lutTexture) return;

		const needsSharpen = sharpness > 0 && sharpenProgram && canvas;

		if (needsSharpen && fbo && fboTexture) {
			glCtx.bindFramebuffer(glCtx.FRAMEBUFFER, fbo);
			glCtx.framebufferTexture2D(glCtx.FRAMEBUFFER, glCtx.COLOR_ATTACHMENT0, glCtx.TEXTURE_2D, fboTexture, 0);
		} else {
			glCtx.bindFramebuffer(glCtx.FRAMEBUFFER, null);
		}

		glCtx.useProgram(lutProgram);
		glCtx.bindVertexArray(vao);

		glCtx.activeTexture(glCtx.TEXTURE0);
		glCtx.bindTexture(glCtx.TEXTURE_2D, imageTexture);
		glCtx.uniform1i(glCtx.getUniformLocation(lutProgram, "u_image"), 0);

		glCtx.activeTexture(glCtx.TEXTURE1);
		glCtx.bindTexture(glCtx.TEXTURE_3D, lutTexture);
		glCtx.uniform1i(glCtx.getUniformLocation(lutProgram, "u_lut"), 1);
		glCtx.uniform1f(glCtx.getUniformLocation(lutProgram, "u_lutSize"), LUT_SIZE);

		glCtx.drawArrays(glCtx.TRIANGLE_STRIP, 0, 4);

		if (needsSharpen && sharpenProgram && canvas) {
			glCtx.bindFramebuffer(glCtx.FRAMEBUFFER, null);
			glCtx.useProgram(sharpenProgram);

			glCtx.activeTexture(glCtx.TEXTURE0);
			glCtx.bindTexture(glCtx.TEXTURE_2D, fboTexture);
			glCtx.uniform1i(glCtx.getUniformLocation(sharpenProgram, "u_texture"), 0);
			// Scale sharpness: 150 is NP3 max, map to a subtle 0-0.5 range for preview
			glCtx.uniform1f(glCtx.getUniformLocation(sharpenProgram, "u_sharpness"), Math.min(sharpness / 300.0, 0.5));
			glCtx.uniform2f(
				glCtx.getUniformLocation(sharpenProgram, "u_texelSize"),
				1.0 / canvas.width,
				1.0 / canvas.height,
			);

			glCtx.drawArrays(glCtx.TRIANGLE_STRIP, 0, 4);
		}

		glCtx.bindVertexArray(null);
	}

	async function regenerateLUT() {
		if (!gl || !wasmReady || !recipe || !generateLUT) return;

		generationCounter++;

		if (pendingGeneration) return;
		pendingGeneration = true;

		while (true) {
			const myGeneration = generationCounter;

			await new Promise<void>((resolve) => {
				if (debounceTimer) clearTimeout(debounceTimer);
				debounceTimer = setTimeout(resolve, 16);
			});

			if (myGeneration !== generationCounter) continue;

			try {
				const recipeJSON = JSON.stringify(recipe);
				const lutData = await generateLUT(recipeJSON, LUT_SIZE);

				if (myGeneration !== generationCounter) continue;

				if (gl) {
					uploadLUT(gl, lutData);
					render(gl, recipe?.sharpness ?? 0);
				}
				break;
			} catch (err) {
				console.error("LUT generation failed:", err);
				break;
			}
		}

		pendingGeneration = false;
	}

	function _handleContextLost(event: Event) {
		event.preventDefault();
		contextLost = true;
		console.warn("WebGL context lost");
	}

	function _handleContextRestored() {
		contextLost = false;
		console.info("WebGL context restored — reinitializing");
		if (!canvas) return;
		const glCtx = canvas.getContext("webgl2");
		if (!glCtx) return;
		gl = glCtx;
		initWebGL(glCtx);
		if (imageData) uploadImage(glCtx, imageData);
		regenerateLUT();
	}

	// Initialize WebGL context
	$effect(() => {
		if (!canvas) return;

		canvas.addEventListener("webglcontextlost", _handleContextLost);
		canvas.addEventListener("webglcontextrestored", _handleContextRestored);

		const glCtx = canvas.getContext("webgl2");
		if (!glCtx) {
			webgl2Unsupported = true;
			return;
		}
		gl = glCtx;
		initWebGL(glCtx);

		// ResizeObserver to handle container size changes
		const container = canvas.parentElement;
		let resizeObserver: ResizeObserver | null = null;
		if (container) {
			resizeObserver = new ResizeObserver(() => {
				if (glCtx && lastImage) {
					updateCanvasSize(glCtx, lastImage);
					setupFboTexture(glCtx);
					render(glCtx, recipe?.sharpness ?? 0);
				}
			});
			resizeObserver.observe(container);
		}

		return () => {
			resizeObserver?.disconnect();
			canvas?.removeEventListener("webglcontextlost", _handleContextLost);
			canvas?.removeEventListener("webglcontextrestored", _handleContextRestored);
			if (glCtx) {
				if (lutProgram) glCtx.deleteProgram(lutProgram);
				if (sharpenProgram) glCtx.deleteProgram(sharpenProgram);
				if (imageTexture) glCtx.deleteTexture(imageTexture);
				if (lutTexture) glCtx.deleteTexture(lutTexture);
				if (fboTexture) glCtx.deleteTexture(fboTexture);
				if (fbo) glCtx.deleteFramebuffer(fbo);
				if (vao) glCtx.deleteVertexArray(vao);
			}
			if (debounceTimer) clearTimeout(debounceTimer);
		};
	});

	// Upload image when it changes
	$effect(() => {
		if (!gl || !imageData) return;
		uploadImage(gl, imageData);
		regenerateLUT();
	});

	// Regenerate LUT when any recipe property changes (including nested color/grading objects)
	$effect(() => {
		if (!recipe || !gl || !imageData) return;
		// Deep-read all properties to subscribe to nested changes
		void JSON.stringify(recipe);
		regenerateLUT();
	});
</script>

{#if webgl2Unsupported}
	<div class="flex items-center justify-center h-full bg-canvas-base text-foreground rounded-lg">
		<p class="text-sm opacity-70">WebGL2 not available — preview requires GPU acceleration</p>
	</div>
{:else if contextLost}
	<div class="flex items-center justify-center h-full bg-canvas-base text-foreground rounded-lg">
		<p class="text-sm opacity-70">GPU context lost — restoring...</p>
	</div>
{:else if !imageData}
	<div class="flex flex-col items-center justify-center h-full bg-canvas-base text-foreground rounded-lg gap-3">
		<span class="text-3xl">📷</span>
		<p class="text-sm opacity-70 text-center max-w-xs">
			Drop an image or click Open Image to see a live preview of your adjustments
		</p>
	</div>
{:else}
	<div
		class="flex items-center justify-center bg-canvas-base rounded-lg overflow-hidden"
		style="position: absolute; inset: 0;"
	>
		<canvas bind:this={canvas}></canvas>
	</div>
{/if}
