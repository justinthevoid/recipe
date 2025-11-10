// wasm-loader.js - WASM Module Initialization
// Epic 2, Stories 2-1 & 2-8: HTML Drag-Drop UI & Error Handling
// Loads recipe.wasm and exposes getVersion() to populate footer

import { showError as showErrorPanel } from './error-handler.js';

/**
 * Initialize the WASM module
 * Loads recipe.wasm, starts Go runtime, and dispatches wasmReady event
 * @returns {Promise<void>}
 */
export async function initializeWASM() {
    console.log('Initializing WASM module...');

    const statusEl = document.getElementById('status');
    const versionEl = document.getElementById('version');

    try {
        // Register wasmReady listener BEFORE loading WASM to avoid race condition
        // (event is dispatched from cmd/wasm/main.go after initialization)
        window.addEventListener('wasmReady', () => {
            console.log('WASM module ready');

            // Update version in footer using getVersion() function
            // getVersion() is exposed by cmd/wasm/main.go
            if (typeof getVersion === 'function') {
                const version = getVersion();
                versionEl.textContent = version;
                console.log('Recipe version:', version);
            } else {
                versionEl.textContent = 'unknown';
                console.warn('getVersion() function not available');
            }

            // Update status banner
            statusEl.className = 'status ready';
            statusEl.textContent = 'Ready to convert';

            // Log available WASM functions for debugging
            console.log('WASM functions available:', {
                convert: typeof convert === 'function',
                detectFormat: typeof detectFormat === 'function',
                getVersion: typeof getVersion === 'function'
            });
        });

        // Create Go runtime instance (global Go class from wasm_exec.js)
        const go = new Go();

        // Load WASM binary
        const result = await WebAssembly.instantiateStreaming(
            fetch('recipe.wasm'),
            go.importObject
        );

        // Start Go program (runs in background, exposes functions)
        // This will trigger the 'wasmReady' event
        go.run(result.instance);

    } catch (error) {
        console.error('Failed to load WASM module:', error);

        // Update status banner with error
        statusEl.className = 'status error';
        statusEl.textContent = `Failed to load converter: ${error.message}`;

        // Update version to show error
        versionEl.textContent = 'error';

        // Story 2-8: Show centralized error panel
        showErrorPanel('wasm-load-failed', error);
    }
}
