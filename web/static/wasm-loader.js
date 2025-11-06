// wasm-loader.js - WASM Module Initialization
// Epic 2, Story 2-1: HTML Drag-Drop UI
// Loads recipe.wasm and exposes getVersion() to populate footer

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

        // Show user-friendly error message
        showError(`Unable to initialize converter. Please refresh the page and try again.`);
    }
}

/**
 * Display error message to user
 * @param {string} message - Error message to display
 */
function showError(message) {
    const errorEl = document.getElementById('errorMessage');
    if (errorEl) {
        errorEl.textContent = message;
        errorEl.style.display = 'block';
    }
}
