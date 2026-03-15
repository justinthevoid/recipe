/**
 * WASM loader for VSCode webview context.
 * Manages Go WASM initialization and exposes generateLUT for preview rendering.
 */

// Go class is defined globally by wasm_exec.js
declare class Go {
	importObject: WebAssembly.Imports;
	run(instance: WebAssembly.Instance): Promise<void>;
}

// Global function registered by WASM binary
declare function generateLUT(recipeJSON: string, size: number): Promise<Uint8Array>;

export type WasmStatus = "idle" | "loading" | "ready" | "error";

const WASM_READY_TIMEOUT_MS = 10_000;

let status = $state<WasmStatus>("idle");
let error = $state<string | null>(null);

export function getWasmStatus(): WasmStatus {
	return status;
}

export function getWasmError(): string | null {
	return error;
}

export async function initWasm(wasmUri: string): Promise<void> {
	if (status === "loading" || status === "ready") return;

	status = "loading";
	error = null;

	try {
		const go = new Go();

		const wasmReadyPromise = new Promise<void>((resolve, reject) => {
			const handler = () => {
				clearTimeout(timer);
				window.removeEventListener("wasmReady", handler);
				resolve();
			};
			const timer = setTimeout(() => {
				window.removeEventListener("wasmReady", handler);
				reject(new Error("WASM ready timeout — binary did not signal readiness"));
			}, WASM_READY_TIMEOUT_MS);
			window.addEventListener("wasmReady", handler);
		});

		const response = await fetch(wasmUri);
		const result = await WebAssembly.instantiateStreaming(response, go.importObject);

		// go.run() resolves when Go main() exits — catch unexpected exits
		go.run(result.instance).catch((err) => {
			console.error("WASM runtime exited unexpectedly:", err);
			status = "error";
			error = err instanceof Error ? err.message : "WASM runtime crashed";
		});

		await wasmReadyPromise;
		status = "ready";
	} catch (err) {
		status = "error";
		error = err instanceof Error ? err.message : "Failed to load WASM";
		console.error("WASM initialization failed:", err);
	}
}

export async function wasmGenerateLUT(recipeJSON: string, size: number): Promise<Float32Array> {
	if (status !== "ready") {
		throw new Error("WASM not ready");
	}

	const result = await generateLUT(recipeJSON, size);
	return new Float32Array(result.buffer, result.byteOffset, result.byteLength / 4);
}
