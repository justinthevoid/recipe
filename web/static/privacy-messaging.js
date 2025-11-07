// privacy-messaging.js - Privacy messaging and education
// Story 2-9: Privacy Messaging
// Handles privacy reminders, FAQ, and education

/**
 * Show privacy reminder when file is uploaded
 * AC-2: Display privacy message after file upload, auto-fade after 5 seconds
 */
export function showPrivacyReminder() {
    const statusEl = document.getElementById('privacyStatus');
    if (statusEl) {
        statusEl.className = 'status privacy';
        statusEl.textContent = '✓ File loaded. Processing locally in your browser - no server upload.';
        statusEl.style.display = 'block';
        statusEl.style.opacity = '1';

        // Fade out after 5 seconds (AC-2: 3-5 seconds, using 5 for visibility)
        setTimeout(() => {
            statusEl.style.opacity = '0';
            // Wait for fade transition (0.5s) before hiding
            setTimeout(() => {
                statusEl.style.display = 'none';
                statusEl.style.opacity = '1'; // Reset for next time
            }, 500);
        }, 5000);
    }
}

/**
 * Show privacy reminder after conversion
 * AC-3: Include privacy message in conversion success
 */
export function showConversionPrivacyMessage() {
    const statusEl = document.getElementById('conversionStatus');
    if (statusEl) {
        // Append privacy message to existing conversion success message
        const currentText = statusEl.textContent;
        statusEl.textContent = `${currentText} Your preset was converted entirely in your browser.`;
    }
}

/**
 * Initialize privacy FAQ
 * AC-4: Initialize FAQ toggle and privacy badge click handler
 */
export function initializePrivacyFAQ() {
    // FAQ toggle button
    const faqToggle = document.getElementById('privacyFAQToggle');
    if (faqToggle) {
        faqToggle.addEventListener('click', togglePrivacyFAQ);
    }

    // Privacy badge click → jump to FAQ and auto-expand
    const privacyBadge = document.getElementById('privacyBadge');
    if (privacyBadge) {
        privacyBadge.addEventListener('click', (e) => {
            e.preventDefault();
            const faqSection = document.getElementById('privacyFAQ');
            const faqContent = document.getElementById('privacyFAQContent');

            if (faqSection) {
                // Smooth scroll to FAQ section
                faqSection.scrollIntoView({ behavior: 'smooth' });

                // Auto-expand FAQ if collapsed
                if (faqContent && faqContent.style.display === 'none') {
                    togglePrivacyFAQ();
                }
            }
        });
    }
}

/**
 * Toggle privacy FAQ visibility
 * AC-4: Expand/collapse FAQ section
 */
function togglePrivacyFAQ() {
    const faqContent = document.getElementById('privacyFAQContent');
    const faqToggle = document.getElementById('privacyFAQToggle');

    if (faqContent && faqToggle) {
        const isHidden = faqContent.style.display === 'none' || !faqContent.style.display;

        if (isHidden) {
            // Expand FAQ
            faqContent.style.display = 'block';
            faqToggle.textContent = 'Hide Privacy FAQ ▲';
            faqToggle.setAttribute('aria-expanded', 'true');
        } else {
            // Collapse FAQ
            faqContent.style.display = 'none';
            faqToggle.textContent = 'Privacy FAQ ▼';
            faqToggle.setAttribute('aria-expanded', 'false');
        }
    }
}
