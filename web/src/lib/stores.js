import { writable } from 'svelte/store';

// WASM Initialization State
export const wasmState = writable({
    status: 'initializing', // initializing, ready, error
    error: null,
    version: '...'
});

// Uploaded Files
// Each item: { id, file, name, size, format, status (queued, processing, complete, error), outputData, outputFormat }
export const files = writable([]);

// User Settings
export const settings = writable({
    targetFormat: '' // np3, xmp, etc.
});

// Helper to add a file
let fileIdCounter = 0;
export function addFile(file) {
    const id = fileIdCounter++;
    files.update(list => [...list, {
        id,
        file,
        name: file.name,
        size: file.size,
        status: 'queued',
        format: null, // Will be detected
        progress: 0
    }]);
    return id;
}

// Helper to update file status
export function updateFileStatus(id, updates) {
    files.update(list => list.map(f => f.id === id ? { ...f, ...updates } : f));
}

// Helper to remove file
export function removeFile(id) {
    files.update(list => list.filter(f => f.id !== id));
}
