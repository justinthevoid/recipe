// upload.js - Batch File Upload Manager
// Handles drag-drop, file validation, file cards, and empty state

import { convertFile } from './converter.js';
import { announceStatus, announceError, moveFocusToFirstFileCard } from './main.js';
import { applyPreviewFilter } from './preview.js';

/**
 * UploadManager Class
 * Manages batch file uploads with drag-and-drop support
 */
export class UploadManager {
    constructor() {
        // DOM references
        this.dropzone = document.getElementById('dropzone');
        this.fileInput = document.getElementById('file-input');
        this.browseButton = document.getElementById('browse-button');
        this.fileList = document.getElementById('file-list');

        // State
        this.uploadedFiles = new Map(); // fileId -> file object
        this.fileIdCounter = 0;

        // Supported formats
        this.supportedExtensions = ['.np3', '.xmp', '.lrtemplate', '.costyle', '.dcp'];

        // Initialize
        this.initEventListeners();
    }

    /**
     * Initialize all event listeners
     */
    initEventListeners() {
        // Browse button click -> trigger file input
        if (this.browseButton) {
            this.browseButton.addEventListener('click', (e) => {
                e.stopPropagation();
                this.fileInput.click();
            });
        }

        // Drop zone click -> trigger file input
        if (this.dropzone) {
            this.dropzone.addEventListener('click', () => {
                this.fileInput.click();
            });

            // Drag-and-drop events
            this.dropzone.addEventListener('dragover', (e) => {
                e.preventDefault();
                e.stopPropagation();
                this.dropzone.classList.add('drag-over');
            });

            this.dropzone.addEventListener('dragleave', (e) => {
                e.preventDefault();
                e.stopPropagation();
                if (e.target === this.dropzone) {
                    this.dropzone.classList.remove('drag-over');
                }
            });

            this.dropzone.addEventListener('drop', (e) => {
                e.preventDefault();
                e.stopPropagation();
                this.dropzone.classList.remove('drag-over');

                const files = e.dataTransfer.files;
                this.handleFiles(files);
            });
        }

        // File input change
        if (this.fileInput) {
            this.fileInput.addEventListener('change', (e) => {
                this.handleFiles(e.target.files);
                e.target.value = '';
            });
        }
    }

    /**
     * Handle uploaded files
     */
    handleFiles(fileList) {
        if (!fileList || fileList.length === 0) return;

        const files = Array.from(fileList);
        const validFiles = [];
        const rejectedFiles = [];

        files.forEach(file => {
            const validation = this.validateFile(file);
            if (validation.valid) {
                validFiles.push(file);
            } else {
                rejectedFiles.push({ name: file.name, reason: validation.message });
            }
        });

        // Add valid files
        validFiles.forEach(file => {
            this.addFileCard(file);
        });

        // Show workspace if we have files
        if (this.uploadedFiles.size > 0) {
            if (window.showWorkspace) window.showWorkspace();
        }

        // Handle errors
        if (rejectedFiles.length > 0) {
            const msg = `${rejectedFiles.length} file(s) rejected. ${rejectedFiles[0].reason}`;
            if (window.updateStatus) window.updateStatus(msg, 'error');
        } else if (validFiles.length > 0) {
            if (window.updateStatus) window.updateStatus(`${validFiles.length} file(s) ready`, 'success');
        }
    }

    /**
     * Validate file extension
     */
    validateFile(file) {
        const fileName = file.name.toLowerCase();
        const extension = this.supportedExtensions.find(ext => fileName.endsWith(ext));

        if (extension) {
            return { valid: true, extension: extension.replace('.', ''), message: '' };
        } else {
            return { valid: false, message: 'Unsupported format' };
        }
    }

    /**
     * Add file card to list
     */
    addFileCard(file) {
        const fileId = this.fileIdCounter++;
        const format = this.detectFormat(file.name);
        const fileSize = this.formatFileSize(file.size);
        const truncatedName = this.truncateFilename(file.name, 30);

        // Store file data
        this.uploadedFiles.set(fileId, {
            id: fileId,
            file: file,
            format: format,
            targetFormat: this.getDefaultTarget(format),
            status: 'queued'
        });

        // Create HTML (Glass Card Design)
        const cardHTML = `
            <div class="file-card" id="file-${fileId}">
                <div class="file-info">
                    <div class="file-details">
                        <div class="file-name">${truncatedName}</div>
                        <div class="file-meta">${fileSize}</div>
                    </div>
                </div>
                <div class="file-actions" style="display: flex; align-items: center; gap: 1rem;">
                    <span class="format-badge">${format.toUpperCase()}</span>
                    <button class="btn-icon remove-btn" data-id="${fileId}" aria-label="Remove file">
                        <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                    </button>
                </div>
            </div>
        `;

        this.fileList.insertAdjacentHTML('beforeend', cardHTML);

        // Add listeners
        const card = document.getElementById(`file-${fileId}`);
        card.querySelector('.remove-btn').addEventListener('click', () => this.removeFile(fileId));

        return fileId;
    }

    removeFile(fileId) {
        const card = document.getElementById(`file-${fileId}`);
        if (card) card.remove();
        this.uploadedFiles.delete(fileId);
    }

    detectFormat(filename) {
        const lower = filename.toLowerCase();
        if (lower.endsWith('.np3')) return 'np3';
        if (lower.endsWith('.xmp')) return 'xmp';
        if (lower.endsWith('.lrtemplate')) return 'lrtemplate';
        if (lower.endsWith('.costyle')) return 'costyle';
        if (lower.endsWith('.dcp')) return 'dcp';
        return 'unknown';
    }

    getDefaultTarget(source) {
        const targets = ['np3', 'xmp', 'lrtemplate', 'costyle', 'dcp'];
        return targets.find(t => t !== source) || 'xmp';
    }

    formatFileSize(bytes) {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
        return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
    }

    truncateFilename(name, max) {
        if (name.length <= max) return name;
        return name.substring(0, max - 3) + '...';
    }
}

// Global instance getter
let instance = null;
export function initializeUploadManager() {
    if (!instance) instance = new UploadManager();
    return instance;
}
export function getUploadManager() {
    return instance;
}
