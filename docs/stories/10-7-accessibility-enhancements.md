# Story 10.7: Accessibility Enhancements (ARIA, Keyboard Navigation)

Status: ready-for-dev

## Story

As a **photographer with visual or motor impairments**,
I want **Recipe to be fully accessible via keyboard navigation and screen readers**,
so that **I can convert presets independently without requiring a mouse or visual interface**.

## Acceptance Criteria

**AC-1: Complete Keyboard Navigation**
- ✅ All interactive elements accessible via keyboard (Tab, Shift+Tab, Enter, Space, Arrow keys)
- ✅ Logical tab order: Upload button → File cards → Convert button → Download buttons
- ✅ Visible focus indicators on all interactive elements (2px solid outline, high contrast)
- ✅ No keyboard traps (user can navigate forward/backward from any element)
- ✅ Skip link at top of page: "Skip to main content" (hidden until focused)
- ✅ Escape key closes modals and cancels operations
- ✅ Enter/Space activates buttons and dropdowns (consistent behavior)

**AC-2: ARIA Labels and Roles**
- ✅ All interactive elements have descriptive ARIA labels:
  - Upload zone: `aria-label="Drag and drop files or click to browse"`
  - Format dropdown: `aria-label="Select target format for conversion"`
  - Convert button: `aria-label="Convert [filename] to [format]"`
  - Download button: `aria-label="Download converted [filename]"`
- ✅ ARIA roles assigned appropriately:
  - `role="button"` for clickable elements (not semantic `<button>`)
  - `role="status"` for progress indicators
  - `role="alert"` for error messages
  - `role="region"` for major sections (hero, upload area, file grid)
- ✅ ARIA live regions for dynamic content:
  - `aria-live="polite"` for file upload notifications
  - `aria-live="assertive"` for errors
  - `aria-atomic="true"` for complete status updates

**AC-3: Screen Reader Support**
- ✅ Meaningful page title: "Recipe - Convert Photo Presets"
- ✅ Landmark regions: `<header>`, `<main>`, `<section>`, `<footer>`
- ✅ Heading hierarchy (h1 → h2 → h3, no skipped levels):
  - h1: "Convert Photo Presets. Instantly. Privately."
  - h2: "Upload Your Presets" (upload section)
  - h2: "Batch Conversion" (file grid section)
- ✅ Alternative text for all visual content:
  - Format badge icons: `alt="NP3 format badge"`, `alt="XMP format badge"`
  - Status icons: `alt="Conversion complete"`, `alt="Error: Conversion failed"`
- ✅ Screen reader-only text for context:
  - "Converting 3 of 10 files" (batch progress)
  - "File 1: example.np3, queued" (per-file status)
- ✅ Form labels associated with inputs:
  - `<label for="file-input">Choose files</label>`
  - `<label for="target-format-[id]">Target format</label>`

**AC-4: Semantic HTML Structure**
- ✅ Use semantic HTML5 elements (not `<div>` soup):
  - `<button>` for all clickable actions (not `<div onclick>`)
  - `<select>` for format dropdown (native, keyboard accessible)
  - `<input type="file">` for file upload (native, accessible)
  - `<progress>` or `<meter>` for batch progress indicator
- ✅ Proper form structure:
  - `<form>` wraps file upload controls
  - `<fieldset>` groups related controls (batch conversion options)
  - `<legend>` labels fieldsets ("Conversion Options")
- ✅ Lists for grouped content:
  - `<ul>` for file grid (each `<li>` is a file card)
  - `<dl>` for file details (filename, size, format)
- ✅ No empty links or buttons (all have text or aria-label)

**AC-5: Focus Management**
- ✅ Focus moves to first file card after upload
- ✅ Focus moves to download button after successful conversion
- ✅ Focus moves to error message after failed conversion
- ✅ Focus returns to trigger element after modal closes
- ✅ Focus trapped inside modal (Tab cycles within modal)
- ✅ Focus restored to page when modal dismissed (Escape key)
- ✅ Autofocus disabled on page load (user controls focus)

**AC-6: Color Contrast and Visual Indicators**
- ✅ All text meets WCAG AA contrast ratios:
  - Body text (16px): 4.5:1 minimum
  - Large text (24px+): 3:1 minimum
  - UI components (buttons, borders): 3:1 minimum
- ✅ Focus indicators meet contrast requirements:
  - 3:1 contrast against background
  - 2px minimum width
- ✅ Error/success states use color + icon:
  - Error: Red color + ✕ icon + text ("Error: Invalid file")
  - Success: Green color + ✓ icon + text ("Conversion complete")
- ✅ Format badges use color + text (not color alone):
  - Badge contains format name: "NP3", "XMP", "lrtemplate"
- ✅ Hover states visible without color (underline, scale, border)

**AC-7: Reduced Motion Support**
- ✅ Respect user's `prefers-reduced-motion` preference:
  ```css
  @media (prefers-reduced-motion: reduce) {
    * {
      animation-duration: 0.01ms !important;
      transition-duration: 0.01ms !important;
    }
  }
  ```
- ✅ Disable non-essential animations:
  - Spinner animation → static icon
  - Hover scale effects → border highlight only
  - Progress bar transitions → instant fill
- ✅ Keep essential feedback:
  - Status changes (queued → complete) still visible
  - Error messages still displayed
- ✅ Test with OS setting enabled:
  - Windows: Settings → Ease of Access → Display → Show animations
  - macOS: System Preferences → Accessibility → Display → Reduce motion
  - Browsers: DevTools → Rendering → Emulate CSS media feature: prefers-reduced-motion

## Tasks / Subtasks

### Task 1: Implement Skip Link (AC-1, AC-2)
- [ ] Add skip link at top of `web/index.html`:
  ```html
  <a href="#main-content" class="skip-link">Skip to main content</a>
  ```
- [ ] Style skip link (hidden until focused):
  ```css
  .skip-link {
    position: absolute;
    top: -40px;
    left: 0;
    background: var(--color-primary);
    color: white;
    padding: 8px 16px;
    text-decoration: none;
    z-index: 1000;
  }

  .skip-link:focus {
    top: 0; /* Reveal on focus */
  }
  ```
- [ ] Add `id="main-content"` to `<main>` element:
  ```html
  <main id="main-content" role="main">
    <!-- Page content -->
  </main>
  ```
- [ ] Test: Tab on page load → Skip link appears → Enter skips to main content

### Task 2: Add ARIA Labels to Interactive Elements (AC-2)
- [ ] Upload zone (drag-drop area):
  ```html
  <div class="upload__dropzone"
       role="button"
       tabindex="0"
       aria-label="Drag and drop files or click to browse for preset files">
    <input type="file" id="file-input" multiple accept=".np3,.xmp,.lrtemplate,.costyle,.dcp" aria-label="Choose preset files">
  </div>
  ```
- [ ] Format dropdown (per-file):
  ```html
  <select id="target-format-{{fileId}}"
          class="file-card__format-select"
          aria-label="Select target format for {{filename}}">
    <option value="">Choose format</option>
    <option value="NP3">NP3 (Nikon)</option>
    <option value="XMP">XMP (Adobe)</option>
    <!-- ... -->
  </select>
  ```
- [ ] Convert button (per-file):
  ```html
  <button class="file-card__convert"
          aria-label="Convert {{filename}} to {{targetFormat}}">
    Convert
  </button>
  ```
- [ ] Batch convert button:
  ```html
  <button id="convert-all-btn"
          aria-label="Convert all {{fileCount}} files to {{targetFormat}}">
    Convert All ({{fileCount}})
  </button>
  ```
- [ ] Download button (per-file):
  ```html
  <a href="{{blobUrl}}"
     download="{{filename}}"
     class="file-card__download"
     aria-label="Download converted {{filename}}">
    Download
  </a>
  ```

### Task 3: Add ARIA Live Regions for Dynamic Content (AC-2)
- [ ] Create status announcement region in HTML:
  ```html
  <div id="status-announcements"
       class="sr-only"
       role="status"
       aria-live="polite"
       aria-atomic="true">
    <!-- JavaScript updates this text -->
  </div>
  ```
- [ ] Create error announcement region:
  ```html
  <div id="error-announcements"
       class="sr-only"
       role="alert"
       aria-live="assertive"
       aria-atomic="true">
    <!-- JavaScript updates this text -->
  </div>
  ```
- [ ] Add screen reader-only CSS utility:
  ```css
  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border-width: 0;
  }
  ```
- [ ] Update JavaScript to announce file uploads:
  ```javascript
  function announceStatus(message) {
    const statusEl = document.getElementById('status-announcements');
    statusEl.textContent = message;
  }

  // Usage:
  announceStatus('File uploaded: example.np3');
  announceStatus('Converting 3 of 10 files');
  announceStatus('Conversion complete');
  ```
- [ ] Update JavaScript to announce errors:
  ```javascript
  function announceError(message) {
    const errorEl = document.getElementById('error-announcements');
    errorEl.textContent = message;
  }

  // Usage:
  announceError('Error: Invalid file format');
  announceError('Error: Conversion failed for example.np3');
  ```

### Task 4: Add Semantic HTML and Landmarks (AC-3, AC-4)
- [ ] Update page structure with semantic HTML:
  ```html
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <title>Recipe - Convert Photo Presets</title>
  </head>
  <body>
    <a href="#main-content" class="skip-link">Skip to main content</a>

    <header role="banner">
      <!-- Logo, navigation (if any) -->
    </header>

    <main id="main-content" role="main">
      <section aria-labelledby="hero-heading">
        <h1 id="hero-heading">Convert Photo Presets. Instantly. Privately.</h1>
        <!-- Hero content -->
      </section>

      <section aria-labelledby="upload-heading">
        <h2 id="upload-heading">Upload Your Presets</h2>
        <form id="upload-form">
          <!-- File upload zone -->
        </form>
      </section>

      <section aria-labelledby="batch-heading" aria-live="polite">
        <h2 id="batch-heading">Batch Conversion</h2>
        <ul id="file-grid" class="upload__files" role="list">
          <!-- File cards (JavaScript-generated) -->
        </ul>
      </section>
    </main>

    <footer role="contentinfo">
      <!-- Privacy notice, links -->
    </footer>
  </body>
  </html>
  ```
- [ ] Replace non-semantic elements:
  - `<div class="button">` → `<button>`
  - `<div class="dropdown">` → `<select>`
  - `<div class="list">` → `<ul>` or `<ol>`
- [ ] Add `role="list"` to file grid (override CSS list-style reset):
  ```html
  <ul class="upload__files" role="list">
    <li class="file-card" role="listitem">
      <!-- File card content -->
    </li>
  </ul>
  ```

### Task 5: Implement Focus Management (AC-1, AC-5)
- [ ] Add visible focus indicators to all interactive elements:
  ```css
  button:focus,
  a:focus,
  input:focus,
  select:focus,
  [role="button"]:focus {
    outline: 2px solid var(--color-primary);
    outline-offset: 2px;
  }

  /* High contrast for visibility */
  button:focus-visible,
  a:focus-visible,
  select:focus-visible {
    outline: 2px solid #0056b3; /* Darker blue, 3:1 contrast */
    outline-offset: 2px;
  }
  ```
- [ ] Move focus to first file card after upload:
  ```javascript
  function handleFileUpload(files) {
    files.forEach(file => {
      const card = createFileCard(file);
      fileGrid.appendChild(card);
    });

    // Move focus to first uploaded file card
    const firstCard = fileGrid.querySelector('.file-card');
    if (firstCard) {
      const convertBtn = firstCard.querySelector('.file-card__convert');
      convertBtn.focus();
    }
  }
  ```
- [ ] Move focus to download button after conversion:
  ```javascript
  function handleConversionComplete(fileId) {
    updateFileStatus(fileId, 'complete');

    // Move focus to download button
    const card = document.getElementById(`file-card-${fileId}`);
    const downloadBtn = card.querySelector('.file-card__download');
    downloadBtn.focus();
  }
  ```
- [ ] Trap focus inside modal (if implemented):
  ```javascript
  function trapFocus(modalElement) {
    const focusableElements = modalElement.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );
    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];

    modalElement.addEventListener('keydown', (e) => {
      if (e.key === 'Tab') {
        if (e.shiftKey && document.activeElement === firstElement) {
          lastElement.focus();
          e.preventDefault();
        } else if (!e.shiftKey && document.activeElement === lastElement) {
          firstElement.focus();
          e.preventDefault();
        }
      }

      if (e.key === 'Escape') {
        closeModal();
      }
    });
  }
  ```

### Task 6: Add Keyboard Event Handlers (AC-1)
- [ ] Handle Enter/Space on upload zone:
  ```javascript
  const uploadZone = document.querySelector('.upload__dropzone');

  uploadZone.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      document.getElementById('file-input').click();
    }
  });
  ```
- [ ] Handle Escape to cancel operations:
  ```javascript
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      if (isModalOpen) {
        closeModal();
      }
      if (isBatchProcessing) {
        cancelBatchConversion();
      }
    }
  });
  ```
- [ ] Handle Arrow keys for dropdown navigation (native `<select>` handles this automatically)
- [ ] Test keyboard navigation:
  - Tab through all interactive elements
  - Verify tab order matches visual flow
  - Verify no keyboard traps
  - Verify Enter/Space activates buttons

### Task 7: Add Alternative Text for Visual Content (AC-3)
- [ ] Format badge icons (inline SVG or CSS background):
  ```html
  <span class="format-badge format-badge--np3" role="img" aria-label="NP3 format">
    NP3
  </span>
  ```
- [ ] Status icons (inline SVG with `<title>`):
  ```html
  <svg class="status-icon status-icon--complete" role="img" aria-labelledby="complete-icon-title">
    <title id="complete-icon-title">Conversion complete</title>
    <path d="M5 13l4 4L19 7" />
  </svg>
  ```
- [ ] Error icons:
  ```html
  <svg class="status-icon status-icon--error" role="img" aria-labelledby="error-icon-title">
    <title id="error-icon-title">Error: Conversion failed</title>
    <path d="M6 18L18 6M6 6l12 12" />
  </svg>
  ```
- [ ] Screen reader text for file status:
  ```html
  <div class="file-card__status">
    <span class="sr-only">Status: Queued for conversion</span>
    <svg class="status-icon status-icon--queued" aria-hidden="true">
      <!-- Icon SVG -->
    </svg>
  </div>
  ```

### Task 8: Implement Reduced Motion Support (AC-7)
- [ ] Add `prefers-reduced-motion` media query to CSS:
  ```css
  /* web/css/main.css */
  @media (prefers-reduced-motion: reduce) {
    *,
    *::before,
    *::after {
      animation-duration: 0.01ms !important;
      animation-iteration-count: 1 !important;
      transition-duration: 0.01ms !important;
      scroll-behavior: auto !important;
    }
  }
  ```
- [ ] Keep essential feedback (no animation, but still visible):
  ```css
  @media (prefers-reduced-motion: reduce) {
    .status-icon--processing {
      animation: none; /* Disable spinner rotation */
    }

    .file-card:hover {
      transform: none; /* Disable hover scale */
      border-color: var(--color-primary); /* Keep border highlight */
    }

    .upload__dropzone:hover {
      transform: none; /* Disable scale animation */
      border-color: var(--color-primary); /* Keep border highlight */
    }
  }
  ```
- [ ] Test with reduced motion enabled:
  - Windows: Settings → Ease of Access → Display → Show animations (OFF)
  - macOS: System Preferences → Accessibility → Display → Reduce motion (ON)
  - Chrome DevTools: Rendering → Emulate CSS media feature: prefers-reduced-motion: reduce

### Task 9: Color Contrast Audit (AC-6)
- [ ] Test all text with WebAIM Contrast Checker (https://webaim.org/resources/contrastchecker/):
  - Body text (#222 on #FFFFFF): ??? ratio (target: ≥4.5:1)
  - Hero title (#000 on #FFFFFF): ??? ratio (target: ≥4.5:1)
  - Button text (white on #2196F3): ??? ratio (target: ≥4.5:1)
  - Error text (#F44336 on #FFFFFF): ??? ratio (target: ≥4.5:1)
- [ ] Test UI components:
  - Focus indicator (#0056b3 on #FFFFFF): ??? ratio (target: ≥3:1)
  - Format badge borders: ??? ratio (target: ≥3:1)
  - File card borders: ??? ratio (target: ≥3:1)
- [ ] Fix any failing contrast ratios:
  - Darken light colors (increase saturation/luminance)
  - Lighten dark colors on dark backgrounds
  - Add border/outline for low-contrast elements
- [ ] Document contrast ratios in CSS comments:
  ```css
  :root {
    --color-text: #222222; /* 16.5:1 on white (WCAG AAA) */
    --color-primary: #0056b3; /* 8.2:1 on white (WCAG AAA) */
    --color-error: #d32f2f; /* 4.5:1 on white (WCAG AA) */
  }
  ```

### Task 10: Manual Accessibility Testing (AC-1, AC-2, AC-3)
- [ ] Test with keyboard only (no mouse):
  - Disconnect mouse or use `Tab` key exclusively
  - Navigate entire page from top to bottom
  - Verify all actions can be completed (upload, convert, download)
  - Verify logical tab order
  - Verify no keyboard traps
- [ ] Test with screen reader:
  - **Windows + NVDA (free)**:
    - Download: https://www.nvaccess.org/download/
    - Test: Navigate page, upload files, convert, download
  - **macOS + VoiceOver (built-in)**:
    - Enable: Cmd+F5
    - Test: Navigate page, upload files, convert, download
  - **Chrome + ChromeVox (extension)**:
    - Install extension
    - Test: Navigate page, upload files, convert, download
- [ ] Test with browser extensions:
  - **axe DevTools (Chrome/Firefox)**:
    - Install: https://www.deque.com/axe/devtools/
    - Run audit on https://recipe.justins.studio
    - Fix all Critical and Serious issues
    - Document Moderate/Minor issues (address if time permits)
  - **Lighthouse Accessibility Audit**:
    - Chrome DevTools → Lighthouse → Accessibility category
    - Target: ≥95 score (from Story 10.6)
    - Fix all failing audits
- [ ] Document any issues found:
  - Issue: "Download button not keyboard accessible"
  - Fix: Added `href` to link (was `<div onclick>`)
  - Verification: Tab to button, Enter downloads file

### Task 11: WCAG Compliance Validation (AC-1 through AC-7)
- [ ] Run WAVE accessibility evaluation (https://wave.webaim.org/):
  - Enter URL: https://recipe.justins.studio
  - Check for errors (0 errors target)
  - Review warnings (minimize warnings)
  - Verify ARIA attributes used correctly
- [ ] Run axe DevTools comprehensive audit:
  - Critical issues: 0 (must fix)
  - Serious issues: 0 (must fix)
  - Moderate issues: <5 (acceptable)
  - Minor issues: document only
- [ ] Verify WCAG 2.1 AA compliance manually:
  - ✅ 1.1.1 Non-text Content: All images have alt text
  - ✅ 1.3.1 Info and Relationships: Semantic HTML, ARIA labels
  - ✅ 1.4.3 Contrast (Minimum): 4.5:1 for text, 3:1 for UI
  - ✅ 2.1.1 Keyboard: All functionality available via keyboard
  - ✅ 2.1.2 No Keyboard Trap: Focus can move freely
  - ✅ 2.4.1 Bypass Blocks: Skip link implemented
  - ✅ 2.4.3 Focus Order: Logical tab order
  - ✅ 2.4.7 Focus Visible: Focus indicators on all elements
  - ✅ 3.2.1 On Focus: No unexpected context changes
  - ✅ 3.2.2 On Input: No unexpected context changes
  - ✅ 4.1.2 Name, Role, Value: ARIA attributes on all interactive elements
- [ ] Document WCAG compliance in README or docs:
  - "Recipe meets WCAG 2.1 AA standards for accessibility"
  - List key features: keyboard navigation, screen reader support, ARIA labels
  - Provide accessibility statement page (optional)

## Dev Notes

### Learnings from Previous Story

**From Story 10-6-performance-optimization (Status: drafted)**

Previous story not yet implemented. Story 10.7 adds accessibility enhancements to ensure Recipe's optimized interface is accessible to all users.

**Integration with Story 10.6:**
- Lighthouse audit (Story 10.6) includes Accessibility score ≥95 (AC-2)
- Story 10.7 implements accessibility features to meet Lighthouse targets
- Reduced motion (Story 10.7) complements performance optimization (Story 10.6)
- Focus indicators (Story 10.7) must not impact 60fps performance (Story 10.6)

**Performance Considerations:**
- Skip link adds minimal HTML (<50 bytes)
- ARIA labels add ~500 bytes to HTML (acceptable within 20KB budget)
- Focus indicators use CSS (no JavaScript overhead)
- Screen reader utilities (sr-only) add ~100 bytes CSS
- Total impact: <1KB (well within Story 10.6 performance budget)

[Source: docs/stories/10-6-performance-optimization.md]

### Architecture Alignment

**Tech Spec Epic 10 Alignment:**

Story 10.7 implements **accessibility requirements** mentioned in tech-spec-epic-10.md Overview: "accessibility (responsive design, keyboard navigation)".

**WCAG 2.1 AA Compliance:**
```
Principle 1: Perceivable
- 1.1.1 Non-text Content: Alt text for all images/icons
- 1.3.1 Info and Relationships: Semantic HTML, ARIA
- 1.4.3 Contrast (Minimum): 4.5:1 text, 3:1 UI

Principle 2: Operable
- 2.1.1 Keyboard: Full keyboard navigation
- 2.1.2 No Keyboard Trap: Focus moves freely
- 2.4.1 Bypass Blocks: Skip link
- 2.4.3 Focus Order: Logical tab order
- 2.4.7 Focus Visible: Focus indicators

Principle 3: Understandable
- 3.2.1 On Focus: No unexpected changes
- 3.2.2 On Input: No unexpected changes

Principle 4: Robust
- 4.1.2 Name, Role, Value: ARIA on all interactive elements
```

[Source: docs/tech-spec-epic-10.md#Overview]

### ARIA Patterns and Best Practices

**ARIA Landmark Roles:**

```html
<header role="banner">      <!-- Page header/logo -->
<nav role="navigation">     <!-- Navigation (if present) -->
<main role="main">          <!-- Main content -->
<section role="region">     <!-- Major sections -->
<footer role="contentinfo"> <!-- Page footer -->
```

**ARIA Live Regions:**

```html
<!-- Polite announcements (non-urgent) -->
<div role="status" aria-live="polite" aria-atomic="true">
  File uploaded: example.np3
</div>

<!-- Assertive announcements (urgent errors) -->
<div role="alert" aria-live="assertive" aria-atomic="true">
  Error: Conversion failed
</div>
```

**ARIA for Custom Controls:**

```html
<!-- Custom button (if not using <button>) -->
<div role="button" tabindex="0" aria-label="Convert file">
  Convert
</div>

<!-- Dropdown (prefer native <select>) -->
<div role="combobox" aria-expanded="false" aria-haspopup="listbox">
  <ul role="listbox">
    <li role="option">NP3</li>
    <li role="option">XMP</li>
  </ul>
</div>

<!-- Progress indicator -->
<div role="progressbar" aria-valuenow="30" aria-valuemin="0" aria-valuemax="100">
  30%
</div>
```

**When to Use `role` Attributes:**
- ✅ Use `role` when semantic HTML not available (e.g., `<div role="button">`)
- ❌ Don't use `role` when semantic HTML exists (e.g., use `<button>`, not `<div role="button">`)
- ✅ Use `role="img"` for decorative SVG with `aria-label`
- ✅ Use `role="status"` and `role="alert"` for dynamic announcements

[Source: ARIA Authoring Practices Guide (APG)]

### Keyboard Navigation Patterns

**Standard Keyboard Shortcuts:**

```
Tab             → Move focus to next interactive element
Shift+Tab       → Move focus to previous interactive element
Enter           → Activate button, submit form, follow link
Space           → Activate button, toggle checkbox, open dropdown
Escape          → Close modal, cancel operation, clear focus
Arrow Up/Down   → Navigate dropdown options (native <select>)
Arrow Left/Right → Navigate slider, tabs (if implemented)
Home/End        → Jump to first/last element in list
```

**Recipe-Specific Keyboard Shortcuts:**

```
Tab             → Navigate: Upload → File cards → Convert → Download
Shift+Tab       → Reverse navigation
Enter/Space     → Activate upload zone (click file picker)
Enter           → Convert file (from Convert button)
Enter           → Download file (from Download link)
Escape          → Cancel batch conversion (if in progress)
Escape          → Close modal (if open)
```

**Focus Order Logic:**

```
1. Skip link (top of page, hidden until focused)
2. Header/logo (if present)
3. Hero section (h1 title)
4. Upload zone (drag-drop area or file picker)
5. File grid (each file card)
   5a. Format dropdown (per file)
   5b. Convert button (per file)
   5c. Download button (per file, after conversion)
6. Batch controls (Convert All, Clear All)
7. Footer (privacy notice, links)
```

[Source: Keyboard Interaction Patterns - W3C WAI]

### Screen Reader Announcements

**Key Announcement Moments:**

```javascript
// File uploaded
announceStatus('File uploaded: example.np3, NP3 format, 25 KB');

// Batch conversion started
announceStatus('Starting batch conversion: 10 files');

// Per-file conversion progress
announceStatus('Converting file 3 of 10: photo.xmp');

// Conversion complete
announceStatus('Conversion complete: example.np3 converted to XMP format');

// Download ready
announceStatus('Download ready: example.xmp');

// Error occurred
announceError('Error: Invalid file format for example.txt');

// Batch complete
announceStatus('Batch conversion complete: 10 files converted successfully');
```

**Announcement Timing:**
- Use `aria-live="polite"` for non-urgent updates (file uploads, conversions)
- Use `aria-live="assertive"` for errors (interrupts screen reader)
- Use `aria-atomic="true"` for complete message replacement (not partial updates)
- Debounce rapid updates (avoid announcement spam)

**Screen Reader Testing Checklist:**
- [ ] Page title read on load: "Recipe - Convert Photo Presets"
- [ ] Heading hierarchy read correctly (h1 → h2 → h3)
- [ ] Landmark regions announced (header, main, footer)
- [ ] All buttons/links have descriptive labels
- [ ] File upload announced with filename, format, size
- [ ] Conversion progress announced (batch + per-file)
- [ ] Errors announced immediately (assertive)
- [ ] Download buttons announced with filename

[Source: Screen Reader Testing Guide - WebAIM]

### Focus Management Patterns

**Auto-Focus Guidelines:**

```javascript
// ✅ GOOD: Move focus after user action
function handleFileUpload(files) {
  // ... create file cards ...

  // Move focus to first card (user initiated upload)
  const firstCard = document.querySelector('.file-card');
  firstCard.querySelector('.file-card__convert').focus();
}

// ✅ GOOD: Move focus after conversion
function handleConversionComplete(fileId) {
  // ... update UI ...

  // Move focus to download button (logical next step)
  const downloadBtn = document.querySelector(`#download-${fileId}`);
  downloadBtn.focus();
}

// ❌ BAD: Auto-focus on page load
window.addEventListener('load', () => {
  document.getElementById('file-input').focus(); // Interrupts screen reader
});

// ❌ BAD: Move focus without user action
setInterval(() => {
  document.querySelector('.status').focus(); // Disorienting
}, 1000);
```

**Modal Focus Trap:**

```javascript
function trapFocus(modal) {
  const focusableElements = modal.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  );
  const firstElement = focusableElements[0];
  const lastElement = focusableElements[focusableElements.length - 1];

  // Focus first element on modal open
  firstElement.focus();

  modal.addEventListener('keydown', (e) => {
    if (e.key === 'Tab') {
      // Shift+Tab on first element → focus last element
      if (e.shiftKey && document.activeElement === firstElement) {
        lastElement.focus();
        e.preventDefault();
      }
      // Tab on last element → focus first element
      else if (!e.shiftKey && document.activeElement === lastElement) {
        firstElement.focus();
        e.preventDefault();
      }
    }

    // Escape closes modal and restores focus
    if (e.key === 'Escape') {
      closeModal();
      // Restore focus to trigger element
      modal.dataset.triggerElement?.focus();
    }
  });
}
```

[Source: Focus Management - ARIA Authoring Practices]

### Testing Tools and Resources

**Browser Extensions:**
- **axe DevTools** (Chrome/Firefox): https://www.deque.com/axe/devtools/
  - Automated WCAG 2.1 testing
  - Highlights specific issues with remediation guidance
  - Integrates with DevTools

- **WAVE** (Chrome/Firefox): https://wave.webaim.org/extension/
  - Visual feedback on accessibility issues
  - Sidebar shows errors, warnings, features
  - Color contrast checker

- **Lighthouse** (Chrome built-in):
  - DevTools → Lighthouse → Accessibility category
  - Automated audit with score 0-100
  - Provides actionable recommendations

**Screen Readers:**
- **NVDA (Windows, free)**: https://www.nvaccess.org/
  - Install and use with Firefox or Chrome
  - Keyboard: Insert+Down to read page

- **VoiceOver (macOS, built-in)**:
  - Enable: Cmd+F5
  - Keyboard: Control+Option+Arrow keys to navigate

- **JAWS (Windows, commercial)**: https://www.freedomscientific.com/products/software/jaws/
  - Industry standard for Windows
  - 40-minute demo mode available

**Testing Websites:**
- **WebAIM Contrast Checker**: https://webaim.org/resources/contrastchecker/
  - Test color contrast ratios
  - WCAG AA/AAA compliance

- **WAVE Web Accessibility Evaluation Tool**: https://wave.webaim.org/
  - Enter URL for full accessibility report
  - Visual overlay of issues

- **W3C Markup Validation**: https://validator.w3.org/
  - Validate HTML structure
  - Catch semantic errors

**References:**
- [WCAG 2.1 Quick Reference](https://www.w3.org/WAI/WCAG21/quickref/)
- [ARIA Authoring Practices Guide](https://www.w3.org/WAI/ARIA/apg/)
- [WebAIM Screen Reader Testing](https://webaim.org/articles/screenreader_testing/)
- [Inclusive Components](https://inclusive-components.design/)

### Project Structure Notes

**New Files Created (Story 10.7):**
```
(No new files - only modifications to existing web/ files)
```

**Modified Files:**
- `web/index.html` - Add ARIA labels, semantic HTML, skip link, landmark regions
- `web/css/main.css` - Add focus indicators, reduced motion support, sr-only utility
- `web/css/components.css` - Add focus styles for buttons, dropdowns, cards
- `web/js/app.js` - Add keyboard event handlers, focus management, ARIA live regions
- `web/js/upload.js` - Add screen reader announcements for file uploads

**Integration with Previous Stories:**
- Story 10.5 (Responsive Design): WCAG visual accessibility (contrast, text size)
- Story 10.6 (Performance): Lighthouse Accessibility score ≥95
- Story 10.7 (Accessibility): WCAG interactive accessibility (keyboard, ARIA, screen readers)

[Source: docs/tech-spec-epic-10.md#Services-and-Modules]

### Testing Strategy

**Accessibility Testing (Required):**

1. **Automated Testing:**
   - Lighthouse Accessibility audit: Target ≥95 (from Story 10.6)
   - axe DevTools: 0 Critical/Serious issues
   - WAVE: 0 Errors, minimize Warnings

2. **Manual Testing:**
   - Keyboard navigation: Tab through entire page, test all interactions
   - Screen reader: NVDA (Windows) or VoiceOver (macOS)
   - Focus management: Verify focus moves logically
   - Color contrast: WebAIM Contrast Checker for all text/UI

3. **Real Device Testing:**
   - Screen reader on mobile: TalkBack (Android), VoiceOver (iOS)
   - Keyboard navigation on desktop: Chrome, Firefox, Safari, Edge
   - Reduced motion: Test with OS setting enabled

**WCAG 2.1 AA Validation:**
- All 4 principles tested (Perceivable, Operable, Understandable, Robust)
- 11 key success criteria verified (see Task 11)
- Document compliance in accessibility statement

**User Testing (Optional):**
- Recruit users with disabilities (screen reader users, keyboard-only users)
- Observe real-world usage
- Gather feedback on accessibility improvements

[Source: docs/tech-spec-epic-10.md#Test-Strategy-Summary]

### Known Risks

**RISK-46: Screen reader compatibility across browsers**
- **Impact**: ARIA announcements may behave differently in NVDA vs VoiceOver vs JAWS
- **Mitigation**: Test with at least 2 screen readers (NVDA + VoiceOver)
- **Acceptable**: Minor differences acceptable, core functionality must work in all

**RISK-47: Focus indicators impact visual design**
- **Impact**: 2px outline may clash with existing design aesthetic
- **Mitigation**: Customize focus style to match brand (still maintain 3:1 contrast)
- **Example**: Use primary color outline instead of browser default blue

**RISK-48: ARIA overuse creates verbosity**
- **Impact**: Too many ARIA labels make screen reader experience noisy
- **Mitigation**: Use semantic HTML first, ARIA only when necessary
- **Test**: Screen reader users should confirm announcements are helpful, not overwhelming

**RISK-49: Keyboard shortcuts conflict with browser/OS shortcuts**
- **Impact**: Custom keyboard shortcuts (if added) may override browser defaults
- **Mitigation**: Avoid single-key shortcuts (use modifier keys: Ctrl+, Alt+)
- **Recipe**: Use native keyboard support only (Tab, Enter, Space, Escape)

[Source: docs/tech-spec-epic-10.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-10.md#Overview] - Mentions "keyboard navigation" as key feature
- [Source: docs/stories/10-5-responsive-mobile-design.md] - WCAG AA visual accessibility (contrast, text size)
- [Source: docs/stories/10-6-performance-optimization.md] - Lighthouse Accessibility score ≥95
- [WCAG 2.1 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/) - Web Content Accessibility Guidelines
- [ARIA Authoring Practices](https://www.w3.org/WAI/ARIA/apg/) - ARIA roles, states, properties
- [WebAIM Articles](https://webaim.org/articles/) - Screen reader testing, keyboard navigation
- [Inclusive Components](https://inclusive-components.design/) - Accessible component patterns
- [axe DevTools](https://www.deque.com/axe/devtools/) - Automated accessibility testing
- [WAVE Tool](https://wave.webaim.org/) - Web accessibility evaluation
- [Contrast Checker](https://webaim.org/resources/contrastchecker/) - WCAG contrast validation

## Dev Agent Record

### Context Reference

- `docs/stories/10-7-accessibility-enhancements.context.xml` - Generated 2025-11-09

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
