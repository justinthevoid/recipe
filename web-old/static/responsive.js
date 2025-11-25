// responsive.js - Responsive adaptations for touch devices
// Epic 2, Story 2-10: Responsive Design
// Handles touch detection, UI adaptations, and device orientation

/**
 * Detect if device is touch-enabled
 * Uses multiple methods for cross-browser compatibility
 * @returns {boolean} True if touch device detected
 */
export function isTouchDevice() {
    // Check multiple touch detection methods
    return (
        ('ontouchstart' in window) ||              // Standard touch events
        (navigator.maxTouchPoints > 0) ||          // Modern browsers (Pointer Events)
        (navigator.msMaxTouchPoints > 0)           // IE/Edge legacy
    );
}

/**
 * Adapt UI for touch devices
 * Modifies drop zone text, collapses parameter panel on mobile, and logs detection
 */
export function adaptForTouch() {
    if (!isTouchDevice()) {
        console.log('Non-touch device detected - keeping default UI');
        return;
    }

    console.log('Touch device detected - adapting UI...');

    // Update drop zone text for touch devices (AC-4, AC-5)
    const dropZonePrimaryText = document.querySelector('.drop-zone .primary-text');
    if (dropZonePrimaryText) {
        dropZonePrimaryText.textContent = 'Tap to select your preset file';
        console.log('Drop zone text updated for touch');
    }

    // Collapse parameter panel by default on mobile (<768px) (AC-5)
    if (window.innerWidth < 768) {
        const parameterPanel = document.querySelector('.parameter-panel');
        if (parameterPanel) {
            parameterPanel.classList.add('collapsed');
            console.log('Parameter panel collapsed for mobile');
        }
    }

    // Add touch-device class to body for CSS targeting
    document.body.classList.add('touch-device');

    console.log('Touch adaptations complete');
}

/**
 * Handle screen orientation change
 * Logs orientation changes (most layout adapts via CSS automatically)
 */
export function handleOrientationChange() {
    // Modern approach: screen.orientation API
    if (screen.orientation) {
        screen.orientation.addEventListener('change', () => {
            console.log('Orientation changed:', screen.orientation.type);
            console.log('New dimensions:', window.innerWidth, 'x', window.innerHeight);

            // Re-check parameter panel collapse state on orientation change
            if (isTouchDevice() && window.innerWidth < 768) {
                const parameterPanel = document.querySelector('.parameter-panel');
                if (parameterPanel && !parameterPanel.classList.contains('collapsed')) {
                    parameterPanel.classList.add('collapsed');
                    console.log('Parameter panel collapsed after rotation to portrait');
                }
            }
        });
    } else {
        // Fallback: orientationchange event (deprecated but still works)
        window.addEventListener('orientationchange', () => {
            console.log('Orientation changed (legacy event)');
            console.log('New dimensions:', window.innerWidth, 'x', window.innerHeight);
        });
    }
}

/**
 * Initialize responsive adaptations
 * Called once on page load
 */
export function initializeResponsive() {
    console.log('Initializing responsive adaptations...');

    // Log initial viewport info
    console.log('Viewport dimensions:', window.innerWidth, 'x', window.innerHeight);
    console.log('Device pixel ratio:', window.devicePixelRatio);
    console.log('Screen dimensions:', screen.width, 'x', screen.height);

    // Detect and adapt for touch
    adaptForTouch();

    // Set up orientation change handling
    handleOrientationChange();

    console.log('Responsive initialization complete');
}
