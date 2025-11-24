// wasm-loader.js - WASM Module Initialization
// Loads recipe.wasm and exposes getVersion() to populate footer

import { showError as showErrorPanel } from './error-handler.js';

// Track WASM initialization state
let wasmInitialized = false;
let wasmInitPromise = null;

/**
 * Initialize the WASM module
 * Loads recipe.wasm, starts Go runtime, and dispatches wasmReady event
 * Can be called multiple times safely - will only initialize once
 * @returns {Promise<void>}
 */
export async function initializeWASM() {
    // If already initialized, return immediately
    if (wasmInitialized) {
        console.log('WASM already initialized');
        return;
    }

    // If initialization in progress, wait for it
    if (wasmInitPromise) {
        console.log('WASM initialization already in progress, waiting...');
        return wasmInitPromise;
    }

    console.log('Initializing WASM module...');

    const statusEl = document.getElementById('status');
    const versionEl = document.getElementById('version');

    // Store promise to prevent duplicate initialization
    wasmInitPromise = (async () => {
        try {
            // Register wasmReady listener BEFORE loading WASM to avoid race condition
            // (event is dispatched from cmd/wasm/main.go after initialization)
            window.addEventListener('wasmReady', () => {
                console.log('WASM module ready');
                wasmInitialized = true;

                // Update version in footer using getVersion() function
                // getVersion() is exposed by cmd/wasm/main.go
                if (typeof getVersion === 'function') {
                    const version = getVersion();
                    if (versionEl) {
                        versionEl.textContent = version;
                    }
                    console.log('Recipe version:', version);
                } else {
                    if (versionEl) {
                        versionEl.textContent = 'unknown';
                    }
                    console.warn('getVersion() function not available');
                }

                // Update status banner
                if (statusEl) {
                    statusEl.className = 'status ready';
                    statusEl.textContent = 'Ready to convert';
                }

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

            // Reset promise so retry is possible
            wasmInitPromise = null;

            // Update status banner with error
            if (statusEl) {
                statusEl.className = 'status error';
                statusEl.textContent = `Failed to load converter: ${error.message}`;
            }

            // Update version to show error
            if (versionEl) {
                versionEl.textContent = 'error';
            }

            // Show centralized error panel
            showErrorPanel('wasm-load-failed', error);

            // Re-throw to propagate error
            throw error;
        }
    })();

    return wasmInitPromise;
}

/**
 * Check if WASM is initialized
 * @returns {boolean} True if WASM is ready
 */
export function isWASMReady() {
    return wasmInitialized;
}
