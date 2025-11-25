/**
 * Recipe UI - Creative Pro
 * Handles interactions for the new app-like interface.
 */

document.addEventListener('DOMContentLoaded', () => {
    initModals();
    initDragDrop();
    initWorkspace();
});

function initModals() {
    const modal = document.getElementById('legal-modal');
    const trigger = document.getElementById('legal-trigger');
    const close = document.getElementById('legal-close');

    if (trigger && modal && close) {
        trigger.addEventListener('click', (e) => {
            e.preventDefault();
            modal.classList.add('active');
        });

        close.addEventListener('click', () => {
            modal.classList.remove('active');
        });

        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.classList.remove('active');
            }
        });
    }
}

function initDragDrop() {
    const dropzone = document.getElementById('dropzone');
    const fileInput = document.getElementById('file-input');
    const browseBtn = document.getElementById('browse-button');

    // Browse Button
    if (browseBtn && fileInput) {
        browseBtn.addEventListener('click', (e) => {
            e.stopPropagation(); // Prevent bubbling to dropzone
            fileInput.click();
        });
    }

    // Drag Effects
    if (dropzone) {
        ['dragenter', 'dragover'].forEach(eventName => {
            dropzone.addEventListener(eventName, (e) => {
                e.preventDefault();
                e.stopPropagation();
                dropzone.classList.add('drag-over');
            }, false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            dropzone.addEventListener(eventName, (e) => {
                e.preventDefault();
                e.stopPropagation();
                dropzone.classList.remove('drag-over');
            }, false);
        });

        // Click to browse (if clicking on the box itself, not button)
        dropzone.addEventListener('click', () => {
            fileInput.click();
        });
    }
}

function initWorkspace() {
    // Logic to switch from "Empty State" to "Workspace State"
    // This will be triggered by main.js when files are added, 
    // but we can set up the listeners here if needed.

    // For now, we'll expose a global helper for main.js to call
    window.showWorkspace = () => {
        const dropzone = document.getElementById('dropzone');
        const workspace = document.getElementById('workspace');

        // Simplify dropzone to just be a top bar or smaller area?
        // For now, let's just reveal the workspace below.
        if (workspace) {
            workspace.classList.remove('hidden');
            workspace.style.display = 'block'; // Ensure it's visible
        }
    };
}

// Global Status Helper (used by main.js)
window.updateStatus = (msg, type = 'info') => {
    const el = document.getElementById('status');
    if (el) {
        el.textContent = msg;
        el.className = `status-indicator status-${type}`;
    }
};
