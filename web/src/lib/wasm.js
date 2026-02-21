/* global Go, getVersion */
import { wasmState } from './stores';

let wasmInitialized = false;
let wasmInitPromise = null;

export async function initializeWASM() {
    if (wasmInitialized) return;
    if (wasmInitPromise) return wasmInitPromise;

    wasmState.update(s => ({ ...s, status: 'initializing' }));

    wasmInitPromise = (async () => {
        try {
            // Wait for Go to be available (loaded via script tag)
            if (typeof Go === 'undefined') {
                throw new Error('Go runtime not loaded');
            }

            const go = new Go();
            const result = await WebAssembly.instantiateStreaming(
                fetch('/recipe.wasm'),
                go.importObject
            );

            // Setup listener promise BEFORE running to catch the event
            // The Go program emits 'wasmReady' when it starts
            const readyPromise = new Promise(resolve => {
                window.addEventListener('wasmReady', resolve, { once: true });
            });

            // Run the Go program
            // This might be synchronous or asynchronous depending on the Go runtime wrapper,
            // but we must have the listener ready before this executes.
            go.run(result.instance);

            // Wait for wasmReady event
            await readyPromise;

            wasmInitialized = true;

            // Get version if available
            let version = 'unknown';
            if (typeof getVersion === 'function') {
                version = getVersion();
            }

            wasmState.set({
                status: 'ready',
                error: null,
                version
            });

        } catch (error) {
            console.error('WASM Init Failed:', error);
            wasmState.set({
                status: 'error',
                error: error.message,
                version: 'error'
            });
            throw error;
        }
    })();

    return wasmInitPromise;
}
