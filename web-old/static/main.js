// main.js - Recipe Pro Entry Point
// Initializes WASM and Batch Upload Manager

import { initializeWASM } from './wasm-loader.js';
import { initializeUploadManager } from './upload.js';
import { startBatchConversion, cancelBatchConversion } from './batch-converter.js';

// Initialize App
document.addEventListener('DOMContentLoaded', async () => {
    console.log('Recipe Pro Initializing...');

    // 1. Initialize UI Managers
    initializeUploadManager();
    initializeBatchControls();

    // 2. Lazy Load WASM (or eager load if preferred for "Pro" feel)
    // Let's eager load it now since we have a status indicator in the sidebar
    try {
        updateStatus('Loading Engine...', 'loading');
        await initializeWASM();
        updateStatus('Ready', 'success');
    } catch (error) {
        console.error('WASM Init Failed:', error);
        updateStatus('Engine Failed', 'error');
    }
});

// Batch Control Logic
function initializeBatchControls() {
    const convertBtn = document.getElementById('convert-all-button');
    const formatSelect = document.getElementById('target-format');

    if (convertBtn && formatSelect) {
        // Enable button only when format selected
        formatSelect.addEventListener('change', () => {
            convertBtn.disabled = !formatSelect.value;
        });

        // Handle Conversion
        convertBtn.addEventListener('click', async () => {
            const format = formatSelect.value;
            if (!format) return;

            convertBtn.disabled = true;
            convertBtn.textContent = 'Converting...';
            updateStatus('Converting...', 'loading');

            try {
                await startBatchConversion(format);
                updateStatus('Conversion Complete', 'success');
            } catch (error) {
                console.error('Batch Error:', error);
                updateStatus('Conversion Failed', 'error');
            } finally {
                convertBtn.disabled = false;
                convertBtn.textContent = 'Convert All';
            }
        });
    }
}

// Global Status Helper (Polyfill if ui.js hasn't loaded yet, though it should have)
if (!window.updateStatus) {
    window.updateStatus = (msg, type) => {
        const el = document.getElementById('status');
        if (el) {
            el.textContent = msg;
            el.className = `status-indicator status-${type}`;
        }
    };
}
