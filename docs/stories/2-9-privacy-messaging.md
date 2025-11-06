# Story 2-9: Privacy Messaging

**Epic:** Epic 2 - Web Interface (FR-2)
**Story ID:** 2-9
**Status:** drafted
**Created:** 2025-11-04
**Complexity:** Simple (0.5-1 day)

---

## User Story

**As a** photographer
**I want** clear assurance that my presets are private
**So that** I feel confident Recipe won't leak my proprietary editing styles to competitors or third parties

---

## Business Value

Privacy is Recipe's **core differentiator**. Photographers often treat presets as trade secrets - they represent years of creative development and can be worth thousands of dollars.

**Competitive advantage:**
- **Recipe:** Files never leave your device (WASM in browser)
- **Competitors:** Upload files to server (privacy risk, GDPR concerns)

**Trust factors:**
- Clear, prominent privacy messaging
- Technical transparency (explain WASM)
- Open source auditability (GitHub link)

**This story builds user confidence to actually use Recipe with valuable presets.**

---

## Acceptance Criteria

### AC-1: Privacy Badge in Header

- [ ] Display privacy badge below main title
- [ ] Badge text: "🔒 100% Privacy - Your files never leave your device"
- [ ] Badge style: Subtle background, small font, always visible
- [ ] Optional: Clickable badge → expands privacy explanation

**Visual design:**
```
┌─────────────────────────────────────┐
│         🍳 Recipe                   │
│   Convert photo presets between     │
│         formats                      │
│                                      │
│ 🔒 100% Privacy - Your files never  │
│    leave your device                 │
└─────────────────────────────────────┘
```

**Test:**
1. Load page
2. Verify: Privacy badge visible in header
3. Verify: Badge text clear and readable
4. Verify: Badge doesn't interfere with main content

### AC-2: Privacy Messaging on File Upload

- [ ] Show privacy reminder when file is uploaded:
  - "✓ File loaded. Processing locally in your browser - no server upload."
- [ ] Display briefly (3-5 seconds), then fade out
- [ ] Subtle, non-intrusive (don't interrupt workflow)

**Test:**
1. Upload file
2. Verify: Privacy message appears after file loads
3. Verify: Message fades out after ~3 seconds
4. Verify: Message doesn't block other UI elements

### AC-3: Privacy Reminder After Conversion

- [ ] Include privacy message in conversion success:
  - "✓ Conversion complete! Your preset was converted entirely in your browser."
- [ ] Emphasize no data was uploaded or stored
- [ ] Link to privacy explanation (FAQ)

**Test:**
1. Complete conversion flow
2. Verify: Success message includes privacy reminder
3. Verify: Link to privacy FAQ (if clicked, opens FAQ section)

### AC-4: Privacy FAQ Section

- [ ] Create expandable FAQ section in footer
- [ ] Questions to answer:
  - **Q: Are my files uploaded to a server?**
    A: No. All processing happens in your browser using WebAssembly. Your files never leave your device.

  - **Q: Does Recipe store or track my files?**
    A: No. Recipe doesn't store any files, and we don't track what you convert. There are no servers, databases, or accounts.

  - **Q: How does Recipe convert files without a server?**
    A: Recipe uses WebAssembly (WASM) to run the conversion engine directly in your browser. Think of it like a mini-application running entirely on your computer.

  - **Q: Can I verify Recipe is private?**
    A: Yes! Recipe is open source. You can review the entire codebase on GitHub to confirm no data is sent anywhere.

  - **Q: What about analytics or tracking?**
    A: Recipe doesn't use tracking or analytics. No cookies, no telemetry, no fingerprinting.

- [ ] FAQ collapsible (click to expand)
- [ ] Link from privacy badge to FAQ

**Test:**
1. Scroll to footer
2. Click "Privacy FAQ" → section expands
3. Read all Q&A → verify accuracy
4. Click privacy badge in header → jumps to FAQ section

### AC-5: Technical Transparency (Explain WASM)

- [ ] Include brief explanation of WebAssembly:
  - "Recipe uses WebAssembly (WASM) to run the conversion engine in your browser. WASM is a safe, sandboxed technology supported by all modern browsers. Your files are processed locally and never transmitted over the network."
- [ ] Link to WebAssembly explainer (external resource)
- [ ] Avoid technical jargon (explain in plain language)

**Test:**
1. Read WASM explanation
2. Verify: Non-technical users can understand
3. Verify: Link to external WASM explainer works

### AC-6: Open Source Transparency

- [ ] Display GitHub link in footer:
  - "📂 Open Source - View the code on GitHub"
- [ ] Emphasize auditability:
  - "All Recipe code is public. You can verify exactly what it does."
- [ ] Link to GitHub repository

**Test:**
1. Click GitHub link in footer
2. Verify: Opens Recipe repository in new tab
3. Verify: Repository is public and accessible

### AC-7: No Tracking or Analytics (Or Minimal Anonymous)

- [ ] **Option 1 (Preferred):** No analytics at all
- [ ] **Option 2:** Minimal anonymous analytics:
  - Only aggregate stats (e.g., "X conversions today")
  - No personally identifiable information
  - No file metadata (names, sizes, contents)
  - Clear disclosure in FAQ
- [ ] Never use third-party trackers (Google Analytics, Facebook Pixel, etc.)

**Test:**
1. Open DevTools → Network tab
2. Load page, complete conversion flow
3. Verify: No requests to analytics services (google-analytics.com, facebook.com, etc.)
4. Verify: No tracking cookies set
5. If analytics used: Verify only anonymous aggregate data sent

### AC-8: Privacy Policy (Optional)

- [ ] Optional: Add minimal privacy policy page
- [ ] Content:
  - "Recipe doesn't collect, store, or transmit your files."
  - "Recipe doesn't use cookies or tracking."
  - "Recipe is a static web application - no backend servers."
- [ ] Link from footer

**Test:**
1. Click "Privacy Policy" link in footer
2. Verify: Policy page loads
3. Verify: Policy is accurate and matches actual behavior

---

## Technical Approach

### Privacy Messaging Component

**File:** `web/static/privacy-messaging.js` (new file)

```javascript
// privacy-messaging.js - Privacy messaging and education

/**
 * Show privacy reminder when file is uploaded
 */
export function showPrivacyReminder() {
    const statusEl = document.getElementById('privacyStatus');
    if (statusEl) {
        statusEl.className = 'status privacy';
        statusEl.textContent = '✓ File loaded. Processing locally in your browser - no server upload.';
        statusEl.style.display = 'block';

        // Fade out after 5 seconds
        setTimeout(() => {
            statusEl.style.opacity = '0';
            setTimeout(() => {
                statusEl.style.display = 'none';
                statusEl.style.opacity = '1';
            }, 500); // Wait for fade transition
        }, 5000);
    }
}

/**
 * Show privacy reminder after conversion
 */
export function showConversionPrivacyMessage() {
    const statusEl = document.getElementById('conversionStatus');
    if (statusEl) {
        const currentText = statusEl.textContent;
        statusEl.textContent = `${currentText} Your preset was converted entirely in your browser.`;
    }
}

/**
 * Initialize privacy FAQ
 */
export function initializePrivacyFAQ() {
    const faqToggle = document.getElementById('privacyFAQToggle');
    if (faqToggle) {
        faqToggle.addEventListener('click', togglePrivacyFAQ);
    }

    // Privacy badge click → jump to FAQ
    const privacyBadge = document.getElementById('privacyBadge');
    if (privacyBadge) {
        privacyBadge.addEventListener('click', (e) => {
            e.preventDefault();
            const faqSection = document.getElementById('privacyFAQ');
            if (faqSection) {
                faqSection.scrollIntoView({ behavior: 'smooth' });
                // Expand FAQ if collapsed
                if (faqSection.style.display === 'none') {
                    togglePrivacyFAQ();
                }
            }
        });
    }
}

/**
 * Toggle privacy FAQ visibility
 */
function togglePrivacyFAQ() {
    const faqContent = document.getElementById('privacyFAQContent');
    const faqToggle = document.getElementById('privacyFAQToggle');

    if (faqContent && faqToggle) {
        if (faqContent.style.display === 'none') {
            faqContent.style.display = 'block';
            faqToggle.textContent = 'Hide Privacy FAQ ▲';
        } else {
            faqContent.style.display = 'none';
            faqToggle.textContent = 'Privacy FAQ ▼';
        }
    }
}
```

### Integration with Main Flow

**Update `main.js`:**

```javascript
// main.js - Integrate privacy messaging

import { showPrivacyReminder, showConversionPrivacyMessage, initializePrivacyFAQ } from './privacy-messaging.js';

// Initialize privacy FAQ on load
document.addEventListener('DOMContentLoaded', () => {
    initializePrivacyFAQ();
});

// Show privacy reminder when file is loaded
window.addEventListener('fileLoaded', () => {
    showPrivacyReminder();
});

// Show privacy message after conversion
window.addEventListener('conversionComplete', () => {
    showConversionPrivacyMessage();
});
```

### HTML Updates

**Update `web/index.html`:**

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🍳 Recipe - Photo Preset Converter</title>
    <link rel="stylesheet" href="static/style.css">
</head>
<body>
    <!-- Hero Section with Privacy Badge -->
    <header>
        <h1>🍳 Recipe</h1>
        <p class="tagline">Convert photo presets between formats</p>
        <p class="formats">Nikon NP3 ↔ Lightroom XMP ↔ lrtemplate</p>

        <!-- Privacy Badge -->
        <a href="#privacyFAQ" id="privacyBadge" class="privacy-badge">
            🔒 100% Privacy - Your files never leave your device
        </a>
    </header>

    <!-- Privacy Status (shown after file upload) -->
    <div id="privacyStatus" class="status privacy" style="display: none;" role="status" aria-live="polite"></div>

    <!-- Main content (drop zone, etc.) -->
    <main>
        <!-- ... existing content ... -->
    </main>

    <!-- Footer with Privacy FAQ -->
    <footer>
        <div class="footer-links">
            <p>
                Recipe v<span id="version">Loading...</span> |
                <a href="https://github.com/justin/recipe" target="_blank" rel="noopener">📂 Open Source - View on GitHub</a>
            </p>
        </div>

        <!-- Privacy FAQ Section -->
        <div class="privacy-faq" id="privacyFAQ">
            <button id="privacyFAQToggle" class="faq-toggle">Privacy FAQ ▼</button>

            <div id="privacyFAQContent" class="faq-content" style="display: none;">
                <h3>Frequently Asked Questions About Privacy</h3>

                <div class="faq-item">
                    <h4>Are my files uploaded to a server?</h4>
                    <p><strong>No.</strong> All processing happens in your browser using WebAssembly. Your files never leave your device.</p>
                </div>

                <div class="faq-item">
                    <h4>Does Recipe store or track my files?</h4>
                    <p><strong>No.</strong> Recipe doesn't store any files, and we don't track what you convert. There are no servers, databases, or accounts.</p>
                </div>

                <div class="faq-item">
                    <h4>How does Recipe convert files without a server?</h4>
                    <p>Recipe uses <strong>WebAssembly (WASM)</strong> to run the conversion engine directly in your browser. Think of it like a mini-application running entirely on your computer. <a href="https://webassembly.org/" target="_blank" rel="noopener">Learn more about WebAssembly</a>.</p>
                </div>

                <div class="faq-item">
                    <h4>Can I verify Recipe is private?</h4>
                    <p><strong>Yes!</strong> Recipe is open source. You can review the entire codebase on <a href="https://github.com/justin/recipe" target="_blank" rel="noopener">GitHub</a> to confirm no data is sent anywhere.</p>
                </div>

                <div class="faq-item">
                    <h4>What about analytics or tracking?</h4>
                    <p>Recipe doesn't use tracking or analytics. No cookies, no telemetry, no fingerprinting. You can verify this by inspecting network requests in your browser's developer tools.</p>
                </div>

                <div class="faq-item">
                    <h4>Is my browser safe running WebAssembly?</h4>
                    <p><strong>Yes.</strong> WebAssembly runs in a sandboxed environment with the same security as JavaScript. It cannot access your file system, network, or other browser tabs without explicit permission.</p>
                </div>
            </div>
        </div>

        <p class="disclaimer">
            Files processed locally via WebAssembly. No server uploads. No data stored.
        </p>
    </footer>

    <!-- Load WASM and initialize -->
    <script src="static/wasm_exec.js"></script>
    <script src="static/main.js" type="module"></script>
</body>
</html>
```

### CSS for Privacy Messaging

**Add to `web/static/style.css`:**

```css
/* Privacy badge */
.privacy-badge {
    display: inline-block;
    margin-top: 0.5rem;
    padding: 0.5rem 1rem;
    background: #c6f6d5;
    color: #22543d;
    border: 1px solid #9ae6b4;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    text-decoration: none;
    transition: all 0.2s ease;
}

.privacy-badge:hover {
    background: #9ae6b4;
    border-color: #68d391;
}

/* Privacy status message */
.status.privacy {
    background: #c6f6d5;
    color: #22543d;
    border: 1px solid #9ae6b4;
    transition: opacity 0.5s ease;
}

/* Privacy FAQ */
.privacy-faq {
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid #e2e8f0;
}

.faq-toggle {
    background: none;
    border: none;
    color: #4a5568;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    padding: 0;
    margin-bottom: 1rem;
}

.faq-toggle:hover {
    color: #2d3748;
    text-decoration: underline;
}

.faq-content {
    margin-top: 1rem;
}

.faq-content h3 {
    margin: 0 0 1rem 0;
    font-size: 1.125rem;
    font-weight: 600;
    color: #2d3748;
}

.faq-item {
    margin-bottom: 1.5rem;
}

.faq-item h4 {
    margin: 0 0 0.5rem 0;
    font-size: 1rem;
    font-weight: 600;
    color: #2d3748;
}

.faq-item p {
    margin: 0;
    font-size: 0.875rem;
    color: #4a5568;
    line-height: 1.6;
}

.faq-item a {
    color: #3182ce;
    text-decoration: underline;
}

.faq-item a:hover {
    color: #2c5aa0;
}

/* Footer disclaimer */
.disclaimer {
    margin-top: 1rem;
    font-size: 0.75rem;
    color: #718096;
    text-align: center;
}
```

---

## Dependencies

### Required Before Starting

- ✅ Stories 2-1 through 2-7 complete (core conversion flow)

### No Blocking Dependencies

Story 2-9 enhances existing stories with privacy messaging.

---

## Testing Plan

### Manual Testing

**Test Case 1: Privacy Badge Visibility**
1. Load page
2. Verify: Privacy badge visible below tagline
3. Verify: Badge text readable and clear
4. Click privacy badge → verify: jumps to Privacy FAQ section

**Test Case 2: Privacy Reminder on File Upload**
1. Upload file
2. Verify: Privacy message appears: "✓ File loaded. Processing locally..."
3. Wait 5 seconds
4. Verify: Message fades out and disappears

**Test Case 3: Privacy Reminder After Conversion**
1. Complete conversion flow
2. Verify: Success message includes: "Your preset was converted entirely in your browser."
3. Verify: Message reinforces privacy

**Test Case 4: Privacy FAQ Functionality**
1. Scroll to footer
2. Click "Privacy FAQ ▼" → FAQ expands
3. Read all 6 Q&A items
4. Verify: All answers accurate and non-technical
5. Click "Hide Privacy FAQ ▲" → FAQ collapses
6. Click privacy badge in header → FAQ expands (if collapsed)

**Test Case 5: External Links**
1. Click "Learn more about WebAssembly" link in FAQ
2. Verify: Opens https://webassembly.org/ in new tab
3. Click "View on GitHub" link in footer
4. Verify: Opens Recipe repository in new tab

**Test Case 6: No Tracking Verification**
1. Open DevTools → Network tab
2. Load page, complete full conversion flow
3. Verify: No requests to third-party tracking services
4. Verify: No cookies set (check DevTools → Application → Cookies)
5. Verify: No local storage or session storage used (unless explicitly documented)

**Test Case 7: User Understanding**
1. Show FAQ to 3 non-technical users
2. Ask: "Does Recipe upload your files to a server?" (Answer should be "No")
3. Ask: "Where does conversion happen?" (Answer should be "In my browser")
4. Ask: "Can Recipe access my files?" (Answer should be "Only files I upload")
5. Verify: All users understand privacy model

---

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Privacy badge visible in header
- [ ] Privacy reminders show after upload and conversion
- [ ] Privacy FAQ comprehensive and accurate
- [ ] External links work (GitHub, WebAssembly explainer)
- [ ] No tracking or analytics (or minimal anonymous, disclosed)
- [ ] User testing completed (non-technical users understand privacy)
- [ ] Manual testing in Chrome, Firefox, Safari
- [ ] Code reviewed
- [ ] Story marked "ready-for-dev" in sprint status

---

## Out of Scope

**Explicitly NOT in this story:**
- ❌ GDPR compliance documentation (Recipe doesn't collect data, so GDPR doesn't apply)
- ❌ Cookie consent banner (no cookies used)
- ❌ Terms of Service (not needed for MVP)
- ❌ Multilingual privacy policy (English only for MVP)

**This story only delivers:** Privacy messaging, education, and transparency to build user trust.

---

## Technical Notes

### Why Privacy Matters for Recipe

**Problem:** Photographers treat presets as trade secrets. Uploading to a server creates:
- **Privacy risk:** Server breach exposes proprietary editing styles
- **Legal risk:** GDPR, CCPA compliance overhead
- **Trust risk:** Users don't know what happens to their files

**Solution:** WASM-based local processing eliminates all these risks.

### WASM as Privacy Feature

**Technical advantage:** WASM runs in browser sandbox:
- No network access (unless explicitly granted)
- No file system access (unless explicitly granted)
- Same security model as JavaScript
- Code is auditable (open source on GitHub)

**Marketing angle:** "Recipe is physically incapable of uploading your files - the technology doesn't allow it."

### Minimal Analytics (If Implemented)

**If analytics are needed** (e.g., for growth metrics):

**Option 1 (Preferred):** No analytics

**Option 2:** Minimal anonymous aggregate:
```javascript
// Example: Only count conversions (no file metadata)
fetch('https://recipe-analytics.example.com/event', {
    method: 'POST',
    body: JSON.stringify({
        event: 'conversion',
        fromFormat: 'np3',  // Generic format name only
        toFormat: 'xmp',     // Generic format name only
        // NO file names, sizes, or contents
    }),
});
```

**Option 3:** Self-hosted Plausible Analytics (privacy-focused, open source)

**Best practice:** Disclose all analytics in FAQ, make opt-out easy.

---

## Follow-Up Stories

**After Story 2-9:**
- Story 2-10: Responsive design for mobile/tablet

**Future enhancements (not Epic 2):**
- Privacy certification (e.g., audit by security firm)
- Privacy-focused analytics (Plausible, Fathom)
- GDPR/CCPA documentation (if Recipe expands to collect any data)
- Multilingual privacy policy

---

## References

- **Tech Spec:** `docs/tech-spec-epic-2.md` (Story 2-9 section)
- **PRD:** `docs/PRD.md` (FR-2.9: Privacy Messaging)
- **WebAssembly Security:** https://webassembly.org/docs/security/
- **Privacy by Design:** GDPR Article 25

---

**Story Created:** 2025-11-04
**Story Owner:** Justin (Developer)
**Reviewer:** Bob (Scrum Master)
**Estimated Effort:** 0.5-1 day
**Status:** Ready for SM approval → move to "ready-for-dev"
