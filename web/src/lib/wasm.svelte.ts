/**
 * WASM loader for web context.
 * Manages Go WASM initialization and exposes conversion + LUT generation functions.
 *
 * Uses nanostores for cross-island state sharing (Astro islands are isolated).
 */

import { wasmStatusStore, wasmErrorStore, wasmVersionStore } from "./shared-stores";

// Go class is defined globally by wasm_exec.js
declare class Go {
	importObject: WebAssembly.Imports;
	run(instance: WebAssembly.Instance): Promise<void>;
}

// Global functions registered by WASM binary
declare function generateLUT(
	recipeJSON: string,
	size: number,
): Promise<Uint8Array>;
declare function convert(
	data: Uint8Array,
	from: string,
	to: string,
): Promise<Uint8Array>;
declare function generate(recipeJSON: string): Promise<Uint8Array>;
declare function detectFormat(data: Uint8Array): string;
declare function getVersion(): string;

export type WasmStatus = "idle" | "loading" | "ready" | "error";

const WASM_READY_TIMEOUT_MS = 10_000;

// Island-local reactive wrappers that read from shared nanostores
class WasmState {
	get status() { return wasmStatusStore.get(); }
	get error() { return wasmErrorStore.get(); }
	get version() { return wasmVersionStore.get(); }
	get ready() { return wasmStatusStore.get() === "ready"; }
}

export const wasm = new WasmState();

export async function initWasm(): Promise<void> {
	const status = wasmStatusStore.get();
	if (status === "loading" || status === "ready") return;

	wasmStatusStore.set("loading");
	wasmErrorStore.set(null);

	try {
		if (typeof Go === "undefined") {
			throw new Error("Go runtime not loaded — ensure wasm_exec.js is included");
		}

		const go = new Go();

		const wasmReadyPromise = new Promise<void>((resolve, reject) => {
			const handler = () => {
				clearTimeout(timer);
				window.removeEventListener("wasmReady", handler);
				resolve();
			};
			const timer = setTimeout(() => {
				window.removeEventListener("wasmReady", handler);
				reject(
					new Error(
						"WASM ready timeout — binary did not signal readiness",
					),
				);
			}, WASM_READY_TIMEOUT_MS);
			window.addEventListener("wasmReady", handler);
		});

		const response = await fetch(`/recipe.wasm?v=${Date.now()}`);
		const result = await WebAssembly.instantiateStreaming(
			response,
			go.importObject,
		);

		go.run(result.instance).catch((err) => {
			console.error("WASM runtime exited unexpectedly:", err);
			wasmStatusStore.set("error");
			wasmErrorStore.set(err instanceof Error ? err.message : "WASM runtime crashed");
		});

		await wasmReadyPromise;

		if (typeof getVersion === "function") {
			wasmVersionStore.set(getVersion());
		}

		wasmStatusStore.set("ready");
	} catch (err) {
		wasmStatusStore.set("error");
		wasmErrorStore.set(err instanceof Error ? err.message : "Failed to load WASM");
		console.error("WASM initialization failed:", err);
	}
}

export function isWasmReady(): boolean {
	return wasmStatusStore.get() === "ready";
}

export async function wasmConvert(
	data: Uint8Array,
	from: string,
	to: string,
): Promise<Uint8Array> {
	if (!wasm.ready) throw new Error("WASM not ready");
	return await convert(data, from, to);
}

export async function wasmGenerate(
	recipeJSON: string,
): Promise<Uint8Array> {
	if (!wasm.ready) throw new Error("WASM not ready");
	return await generate(recipeJSON);
}

export async function wasmGenerateLUT(
	recipeJSON: string,
	size: number,
): Promise<Float32Array> {
	if (!wasm.ready) throw new Error("WASM not ready");
	const result = await generateLUT(recipeJSON, size);
	return new Float32Array(
		result.buffer,
		result.byteOffset,
		result.byteLength / 4,
	);
}

export function wasmDetectFormat(data: Uint8Array): string {
	if (!wasm.ready) throw new Error("WASM not ready");
	return detectFormat(data);
}

export async function wasmExtractFullRecipe(
	data: Uint8Array,
	format: string,
): Promise<Record<string, unknown>> {
	if (!wasm.ready) throw new Error("WASM not ready");
	// Use window.extractFullRecipe to avoid shadowing
	const jsonString = await (
		window as unknown as {
			extractFullRecipe: (
				data: Uint8Array,
				format: string,
			) => Promise<string>;
		}
	).extractFullRecipe(data, format);
	return JSON.parse(jsonString);
}
