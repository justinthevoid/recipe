# Story 11.3: Preview Modal with Before/After Slider

Status: ready-for-dev

## Story

As a **photographer evaluating preset conversions**,
I want **a modal interface with a before/after slider to preview presets on reference images**,
so that **I can visually compare the original and filtered versions before converting my files**.

## Acceptance Criteria

**AC-1: Modal Opens When Preview Button Clicked**
- ✅ "Preview" button visible after file upload and format detection
- ✅ Button state:
  - Enabled: When valid preset file uploaded
  - Disabled: When no file uploaded or invalid format
  - Label: "Preview Preset" or "See Preview"
- ✅ Click event triggers modal open:
  - Modal appears with fade-in animation (300ms)
  - Background dimmed with overlay (rgba(0,0,0,0.5))
  - Body scroll locked (prevent page scroll while modal open)
- ✅ Modal displays reference image:
  - Default reference image: Portrait (first tab selected)
  - CSS filters applied immediately on modal open (Story 11.1)
  - Before/after slider initialized at 50% position

**AC-2: Before/After Slider with Reference Image**
- ✅ Slider structure:
  - Left side: Original reference image (no filters)
  - Right side: Filtered reference image (CSS filters applied)
  - Vertical divider line: Separates before/after sections
  - Drag handle: Circle/pill on divider line (visual affordance)
- ✅ Slider interaction:
  - Mouse: Click and drag handle left/right
  - Touch: Tap and drag handle left/right
  - Keyboard: Arrow keys move handle (←/→ for 5% steps)
  - Snap behavior: Optionally snap to 0%, 50%, 100% positions
- ✅ Slider position indicator:
  - Visual cue: Handle position shows % filtered (0% = all original, 100% = all filtered)
  - Optional: Numeric label "50%" above/below handle
- ✅ Smooth rendering:
  - 60fps minimum (no janky updates)
  - CSS `clip-path` or `mask` for smooth reveal
  - GPU acceleration (`will-change: clip-path`)

**AC-3: Tabs to Switch Between Reference Images**
- ✅ Tab interface above slider:
  - Three tabs: "Portrait", "Landscape", "Product"
  - Active tab: Highlighted (border-bottom, color change)
  - Inactive tabs: Muted (gray, clickable)
- ✅ Tab click behavior:
  - Click tab → Reference image changes (Portrait/Landscape/Product)
  - Slider position preserved (e.g., 50% before → 50% after tab change)
  - CSS filters re-applied to new reference image
- ✅ Tab keyboard navigation:
  - Arrow keys (←/→) switch tabs
  - Tab key cycles through tabs + slider + buttons
  - Enter/Space activates selected tab
- ✅ Tab accessibility:
  - ARIA role="tablist", role="tab", aria-selected="true/false"
  - ARIA label: "Reference image: Portrait" (screen reader)

**AC-4: Preset Parameters Displayed**
- ✅ Parameter summary visible in modal:
  - Format: "Exposure +0.7 • Contrast +15 • Warmth +10"
  - Position: Below slider, above Convert button
  - Style: Small text (12-14px), muted color (gray)
- ✅ Parameters shown:
  - Exposure (if non-zero)
  - Contrast (if non-zero)
  - Saturation (if non-zero)
  - Hue (if non-zero)
  - Temperature/Tint (if non-zero)
  - Other parameters (if non-zero)
- ✅ Parameter formatting:
  - Positive values: "+0.7" (explicit plus sign)
  - Negative values: "-0.5" (explicit minus sign)
  - Zero values: Omitted (not shown)
  - Separator: " • " (bullet point)
  - Example: "Exposure +1.2 • Contrast +20 • Saturation -10"

**AC-5: Convert Now Button in Modal**
- ✅ "Convert Now" button visible in modal:
  - Position: Bottom right (or center on mobile)
  - Style: Primary button (blue, prominent)
  - Label: "Convert Now" or "Apply & Convert"
- ✅ Click behavior:
  - Close modal
  - Proceed to conversion (trigger conversion workflow)
  - Show progress indicator (Story 10.3)
  - Download converted file when complete (Story 2.7)
- ✅ Button state:
  - Enabled: Always (valid preset already uploaded)
  - Disabled: Never (preview only shown for valid presets)

**AC-6: Close/Cancel Button Returns to Upload Screen**
- ✅ Close button options:
  - "X" button (top right corner of modal)
  - "Cancel" button (bottom left of modal)
  - Click outside modal (background overlay)
  - Keyboard: Esc key
- ✅ Close behavior:
  - Modal fades out (300ms animation)
  - Background overlay removed
  - Body scroll unlocked
  - Return to upload screen (file still uploaded, ready to re-preview or convert)
- ✅ No conversion triggered:
  - Cancel does NOT convert file
  - Cancel does NOT reset upload state
  - User can re-open preview by clicking "Preview" button again

**AC-7: Modal Keyboard Accessible**
- ✅ Keyboard navigation:
  - Tab: Cycle through tabs, slider, Convert button, Cancel button, Close (X)
  - Shift+Tab: Reverse cycle
  - Enter/Space: Activate button or tab
  - Arrow keys (←/→): Switch tabs, move slider
  - Esc: Close modal
- ✅ Focus management:
  - Modal opened: Focus moves to first tab (Portrait)
  - Modal closed: Focus returns to "Preview" button
  - Focus trap: Tab navigation stays within modal (no focus on background elements)
- ✅ ARIA attributes:
  - `role="dialog"` on modal container
  - `aria-modal="true"` (prevents interaction with background)
  - `aria-labelledby="modal-title"` (links to modal heading)
  - `aria-describedby="modal-description"` (links to parameter summary)
- ✅ Screen reader support:
  - Modal title announced: "Preset Preview"
  - Tab changes announced: "Portrait selected" (ARIA live region)
  - Slider position announced: "50% filtered" (ARIA value text)

**AC-8: Mobile Full-Screen Modal with Touch-Friendly Controls**
- ✅ Mobile layout (<600px viewport):
  - Full-screen modal (no margins, 100vw x 100vh)
  - Tabs at top (horizontal scroll if needed)
  - Slider fills screen width (larger drag handle, easier to touch)
  - Convert button: Full-width, fixed at bottom
  - Cancel/Close: Top-left corner (X icon)
- ✅ Touch interactions:
  - Tap and drag slider handle (minimum 44x44px touch target)
  - Tap left/right sides of slider → Snap to 0% or 100%
  - Swipe tabs left/right to switch reference images
- ✅ Mobile optimizations:
  - Reference images: Load responsive 400px versions (Story 11.2)
  - Slider handle: Larger (60px diameter, easier to grab)
  - Buttons: Minimum 48px height (WCAG touch target)

## Tasks / Subtasks

### Task 1: Create Modal HTML Structure (AC-1, AC-2, AC-3)
- [ ] Create modal HTML template in `web/index.html`:
  ```html
  <!-- Preview Modal -->
  <div id="preview-modal" class="modal" role="dialog" aria-modal="true" aria-labelledby="modal-title" aria-describedby="modal-description" hidden>
    <div class="modal-overlay" aria-hidden="true"></div>
    <div class="modal-container">
      <div class="modal-header">
        <h2 id="modal-title">Preset Preview</h2>
        <button class="modal-close" aria-label="Close preview">
          <svg><!-- X icon --></svg>
        </button>
      </div>

      <div class="modal-body">
        <!-- Reference Image Tabs -->
        <div class="preview-tabs" role="tablist">
          <button class="preview-tab" role="tab" aria-selected="true" aria-controls="preview-image-portrait" id="tab-portrait">Portrait</button>
          <button class="preview-tab" role="tab" aria-selected="false" aria-controls="preview-image-landscape" id="tab-landscape">Landscape</button>
          <button class="preview-tab" role="tab" aria-selected="false" aria-controls="preview-image-product" id="tab-product">Product</button>
        </div>

        <!-- Before/After Slider -->
        <div class="preview-slider-container">
          <div class="preview-slider">
            <!-- Before image (left side) -->
            <div class="preview-before">
              <picture id="preview-image-portrait" role="tabpanel" aria-labelledby="tab-portrait">
                <source srcset="/images/preview-portrait-400w.webp" media="(max-width: 600px)" type="image/webp">
                <source srcset="/images/preview-portrait-800w.webp" media="(max-width: 1024px)" type="image/webp">
                <source srcset="/images/preview-portrait.webp" type="image/webp">
                <img src="/images/preview-portrait.jpg" alt="Portrait reference image (original)" loading="lazy" class="preview-image">
              </picture>
              <!-- Landscape and Product images (hidden by default) -->
            </div>

            <!-- After image (right side, with CSS filters applied) -->
            <div class="preview-after">
              <picture>
                <!-- Same images as above, but with CSS filters applied -->
                <img src="/images/preview-portrait.jpg" alt="Portrait reference image (filtered)" class="preview-image preview-image-filtered">
              </picture>
            </div>

            <!-- Slider handle -->
            <input type="range" min="0" max="100" value="50" class="preview-slider-handle" aria-label="Adjust preview comparison" aria-valuemin="0" aria-valuemax="100" aria-valuenow="50" aria-valuetext="50% filtered">
            <div class="preview-slider-divider"></div>
          </div>

          <!-- Parameter summary -->
          <div id="modal-description" class="preview-parameters" aria-live="polite">
            Exposure +0.7 • Contrast +15 • Warmth +10
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn btn-secondary modal-cancel">Cancel</button>
        <button class="btn btn-primary modal-convert">Convert Now</button>
      </div>
    </div>
  </div>
  ```
- [ ] Add "Preview" button to upload screen:
  ```html
  <!-- In file upload section -->
  <button id="preview-button" class="btn btn-primary" disabled>Preview Preset</button>
  ```
- [ ] Enable "Preview" button after file upload:
  ```javascript
  // In upload.js
  function handleFileUpload(file) {
    // ... existing upload logic ...

    // Enable preview button
    document.getElementById('preview-button').disabled = false;
  }
  ```

### Task 2: Implement Modal Open/Close Logic (AC-1, AC-6)
- [ ] Create `web/js/modal.js` with modal management:
  ```javascript
  // Open modal
  function openPreviewModal() {
    const modal = document.getElementById('preview-modal');
    modal.hidden = false;

    // Lock body scroll
    document.body.style.overflow = 'hidden';

    // Focus first tab
    document.getElementById('tab-portrait').focus();

    // Apply CSS filters to reference image
    applyPreviewFilters();

    // Trap focus within modal
    setupFocusTrap(modal);
  }

  // Close modal
  function closePreviewModal() {
    const modal = document.getElementById('preview-modal');
    modal.hidden = true;

    // Unlock body scroll
    document.body.style.overflow = '';

    // Return focus to preview button
    document.getElementById('preview-button').focus();
  }

  // Event listeners
  document.getElementById('preview-button').addEventListener('click', openPreviewModal);
  document.querySelector('.modal-close').addEventListener('click', closePreviewModal);
  document.querySelector('.modal-cancel').addEventListener('click', closePreviewModal);
  document.querySelector('.modal-overlay').addEventListener('click', closePreviewModal);

  // Keyboard: Esc to close
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && !document.getElementById('preview-modal').hidden) {
      closePreviewModal();
    }
  });

  // Focus trap (Tab cycles within modal)
  function setupFocusTrap(modal) {
    const focusableElements = modal.querySelectorAll('button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])');
    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];

    modal.addEventListener('keydown', (e) => {
      if (e.key === 'Tab') {
        if (e.shiftKey && document.activeElement === firstElement) {
          e.preventDefault();
          lastElement.focus();
        } else if (!e.shiftKey && document.activeElement === lastElement) {
          e.preventDefault();
          firstElement.focus();
        }
      }
    });
  }
  ```

### Task 3: Implement Before/After Slider (AC-2)
- [ ] Create slider CSS in `web/css/modal.css`:
  ```css
  .preview-slider {
    position: relative;
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    overflow: hidden;
  }

  .preview-before,
  .preview-after {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
  }

  .preview-after {
    clip-path: inset(0 0 0 var(--slider-position, 50%));
    will-change: clip-path; /* GPU acceleration */
  }

  .preview-image-filtered {
    filter: var(--preview-filter); /* CSS filters from Story 11.1 */
  }

  .preview-slider-divider {
    position: absolute;
    top: 0;
    left: var(--slider-position, 50%);
    width: 2px;
    height: 100%;
    background: white;
    box-shadow: 0 0 5px rgba(0,0,0,0.5);
    pointer-events: none;
  }

  .preview-slider-handle {
    position: absolute;
    top: 50%;
    left: var(--slider-position, 50%);
    transform: translate(-50%, -50%);
    width: 60px;
    height: 60px;
    background: white;
    border: 2px solid var(--color-primary);
    border-radius: 50%;
    cursor: ew-resize;
    box-shadow: 0 2px 8px rgba(0,0,0,0.2);
  }

  /* Mobile: Larger handle */
  @media (max-width: 600px) {
    .preview-slider-handle {
      width: 80px;
      height: 80px;
    }
  }
  ```
- [ ] Implement slider interaction in `web/js/slider.js`:
  ```javascript
  const slider = document.querySelector('.preview-slider-handle');
  const sliderContainer = document.querySelector('.preview-slider');
  let isDragging = false;

  // Mouse/touch drag
  slider.addEventListener('mousedown', startDrag);
  slider.addEventListener('touchstart', startDrag);

  function startDrag(e) {
    isDragging = true;
    updateSliderPosition(e);
  }

  document.addEventListener('mousemove', (e) => {
    if (isDragging) updateSliderPosition(e);
  });

  document.addEventListener('touchmove', (e) => {
    if (isDragging) updateSliderPosition(e);
  });

  document.addEventListener('mouseup', () => isDragging = false);
  document.addEventListener('touchend', () => isDragging = false);

  function updateSliderPosition(e) {
    const rect = sliderContainer.getBoundingClientRect();
    const x = (e.clientX || e.touches[0].clientX) - rect.left;
    const percentage = Math.max(0, Math.min(100, (x / rect.width) * 100));

    // Update CSS variable
    sliderContainer.style.setProperty('--slider-position', `${percentage}%`);

    // Update ARIA value
    slider.setAttribute('aria-valuenow', Math.round(percentage));
    slider.setAttribute('aria-valuetext', `${Math.round(percentage)}% filtered`);
  }

  // Keyboard: Arrow keys
  slider.addEventListener('keydown', (e) => {
    const currentValue = parseInt(slider.getAttribute('aria-valuenow'));
    let newValue = currentValue;

    if (e.key === 'ArrowLeft') {
      newValue = Math.max(0, currentValue - 5);
    } else if (e.key === 'ArrowRight') {
      newValue = Math.min(100, currentValue + 5);
    }

    if (newValue !== currentValue) {
      sliderContainer.style.setProperty('--slider-position', `${newValue}%`);
      slider.setAttribute('aria-valuenow', newValue);
      slider.setAttribute('aria-valuetext', `${newValue}% filtered`);
    }
  });
  ```

### Task 4: Implement Tab Switching (AC-3)
- [ ] Add tab switching logic to `web/js/tabs.js`:
  ```javascript
  const tabs = document.querySelectorAll('.preview-tab');
  const referenceImages = {
    'portrait': document.getElementById('preview-image-portrait'),
    'landscape': document.getElementById('preview-image-landscape'),
    'product': document.getElementById('preview-image-product')
  };

  tabs.forEach(tab => {
    tab.addEventListener('click', () => switchTab(tab));
  });

  function switchTab(selectedTab) {
    const imageType = selectedTab.id.replace('tab-', '');

    // Update tab states
    tabs.forEach(tab => {
      const isSelected = tab === selectedTab;
      tab.setAttribute('aria-selected', isSelected);
      tab.classList.toggle('active', isSelected);
    });

    // Switch reference image
    Object.keys(referenceImages).forEach(type => {
      const isVisible = type === imageType;
      referenceImages[type].hidden = !isVisible;

      // Apply CSS filters to newly visible image
      if (isVisible) {
        const filteredImage = referenceImages[type].querySelector('.preview-image-filtered');
        applyPreviewFilters(filteredImage);
      }
    });

    // Announce tab change to screen readers
    const announcement = document.createElement('div');
    announcement.setAttribute('role', 'status');
    announcement.setAttribute('aria-live', 'polite');
    announcement.textContent = `${imageType.charAt(0).toUpperCase() + imageType.slice(1)} selected`;
    document.body.appendChild(announcement);
    setTimeout(() => announcement.remove(), 1000);
  }

  // Keyboard: Arrow keys switch tabs
  document.querySelector('.preview-tabs').addEventListener('keydown', (e) => {
    const currentIndex = Array.from(tabs).indexOf(document.activeElement);
    let newIndex = currentIndex;

    if (e.key === 'ArrowLeft') {
      newIndex = Math.max(0, currentIndex - 1);
    } else if (e.key === 'ArrowRight') {
      newIndex = Math.min(tabs.length - 1, currentIndex + 1);
    }

    if (newIndex !== currentIndex) {
      tabs[newIndex].focus();
      switchTab(tabs[newIndex]);
    }
  });
  ```

### Task 5: Display Preset Parameters (AC-4)
- [ ] Create parameter display function in `web/js/preview.js`:
  ```javascript
  function formatPresetParameters(recipe) {
    const params = [];

    // Format each parameter (only if non-zero)
    if (recipe.exposure && recipe.exposure !== 0) {
      const sign = recipe.exposure > 0 ? '+' : '';
      params.push(`Exposure ${sign}${recipe.exposure.toFixed(1)}`);
    }

    if (recipe.contrast && recipe.contrast !== 0) {
      const sign = recipe.contrast > 0 ? '+' : '';
      params.push(`Contrast ${sign}${recipe.contrast}`);
    }

    if (recipe.saturation && recipe.saturation !== 0) {
      const sign = recipe.saturation > 0 ? '+' : '';
      params.push(`Saturation ${sign}${recipe.saturation}`);
    }

    if (recipe.hue && recipe.hue !== 0) {
      const sign = recipe.hue > 0 ? '+' : '';
      params.push(`Hue ${sign}${recipe.hue}\u00B0`); // degree symbol
    }

    if (recipe.temperature && recipe.temperature !== 0) {
      const sign = recipe.temperature > 0 ? '+' : '';
      params.push(`Warmth ${sign}${recipe.temperature}`);
    }

    // Join with bullet separator
    return params.join(' • ');
  }

  // Update parameter display in modal
  function updateParameterDisplay(recipe) {
    const paramElement = document.querySelector('.preview-parameters');
    paramElement.textContent = formatPresetParameters(recipe);
  }

  // Call when modal opens
  function openPreviewModal() {
    // ... existing modal open logic ...
    updateParameterDisplay(currentRecipe);
  }
  ```
- [ ] Add unit tests for parameter formatting:
  ```javascript
  // web/tests/preview.test.js
  describe('formatPresetParameters', () => {
    it('formats positive values with plus sign', () => {
      const recipe = { exposure: 0.7, contrast: 15 };
      expect(formatPresetParameters(recipe)).toBe('Exposure +0.7 • Contrast +15');
    });

    it('formats negative values with minus sign', () => {
      const recipe = { exposure: -0.5, saturation: -10 };
      expect(formatPresetParameters(recipe)).toBe('Exposure -0.5 • Saturation -10');
    });

    it('omits zero values', () => {
      const recipe = { exposure: 0.7, contrast: 0, saturation: -10 };
      expect(formatPresetParameters(recipe)).toBe('Exposure +0.7 • Saturation -10');
    });

    it('returns empty string for all zeros', () => {
      const recipe = { exposure: 0, contrast: 0 };
      expect(formatPresetParameters(recipe)).toBe('');
    });
  });
  ```

### Task 6: Implement Convert Now Button (AC-5)
- [ ] Add click handler for "Convert Now" button:
  ```javascript
  document.querySelector('.modal-convert').addEventListener('click', () => {
    // Close modal
    closePreviewModal();

    // Trigger conversion workflow (Story 2.6)
    triggerConversion(currentFile, currentRecipe);
  });

  function triggerConversion(file, recipe) {
    // Show progress indicator (Story 10.3)
    showProgressIndicator();

    // Start conversion (WASM, Story 2.6)
    convertFile(file, recipe)
      .then(convertedFile => {
        // Hide progress indicator
        hideProgressIndicator();

        // Trigger download (Story 2.7)
        downloadFile(convertedFile);
      })
      .catch(error => {
        // Show error (Story 2.8)
        showErrorMessage(error);
      });
  }
  ```

### Task 7: Modal Accessibility (AC-7)
- [ ] Add ARIA attributes to modal HTML (already in Task 1)
- [ ] Implement focus trap (already in Task 2)
- [ ] Test keyboard navigation:
  - Tab through all interactive elements (tabs, slider, buttons)
  - Arrow keys switch tabs and move slider
  - Esc closes modal
  - Focus returns to Preview button after close
- [ ] Test screen reader support:
  - NVDA (Windows): Announce modal title, tab changes, slider position
  - VoiceOver (macOS): Announce modal title, tab changes, slider position
  - TalkBack (Android): Announce modal title, tab changes, slider position
- [ ] Add ARIA live region for tab change announcements (already in Task 4)

### Task 8: Mobile Full-Screen Layout (AC-8)
- [ ] Add mobile-specific CSS:
  ```css
  @media (max-width: 600px) {
    .modal-container {
      width: 100vw;
      height: 100vh;
      max-width: none;
      margin: 0;
      border-radius: 0;
    }

    .preview-tabs {
      overflow-x: auto;
      white-space: nowrap;
      -webkit-overflow-scrolling: touch;
    }

    .preview-slider {
      height: calc(100vh - 200px); /* Full screen minus header/footer */
    }

    .modal-footer {
      position: fixed;
      bottom: 0;
      left: 0;
      width: 100%;
      padding: 16px;
      background: white;
      border-top: 1px solid #ddd;
    }

    .modal-convert {
      width: 100%;
      min-height: 48px; /* WCAG touch target */
    }
  }
  ```
- [ ] Add touch interactions:
  ```javascript
  // Tap left/right sides to snap slider to 0% or 100%
  sliderContainer.addEventListener('click', (e) => {
    const rect = sliderContainer.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const percentage = (x / rect.width) * 100;

    // Snap to 0% or 100% if tapped on sides
    if (percentage < 25) {
      sliderContainer.style.setProperty('--slider-position', '0%');
      slider.setAttribute('aria-valuenow', 0);
    } else if (percentage > 75) {
      sliderContainer.style.setProperty('--slider-position', '100%');
      slider.setAttribute('aria-valuenow', 100);
    }
  });

  // Swipe tabs left/right (optional)
  let touchStartX = 0;
  document.querySelector('.preview-tabs').addEventListener('touchstart', (e) => {
    touchStartX = e.touches[0].clientX;
  });

  document.querySelector('.preview-tabs').addEventListener('touchend', (e) => {
    const touchEndX = e.changedTouches[0].clientX;
    const deltaX = touchEndX - touchStartX;

    if (deltaX > 50) {
      // Swipe right → Previous tab
      const currentIndex = Array.from(tabs).indexOf(document.querySelector('.preview-tab[aria-selected="true"]'));
      if (currentIndex > 0) switchTab(tabs[currentIndex - 1]);
    } else if (deltaX < -50) {
      // Swipe left → Next tab
      const currentIndex = Array.from(tabs).indexOf(document.querySelector('.preview-tab[aria-selected="true"]'));
      if (currentIndex < tabs.length - 1) switchTab(tabs[currentIndex + 1]);
    }
  });
  ```

### Task 9: Integration with Story 11.1 (CSS Filters) and Story 11.2 (Reference Images)
- [ ] Apply CSS filters to filtered image:
  ```javascript
  function applyPreviewFilters(imageElement = null) {
    const targetImage = imageElement || document.querySelector('.preview-image-filtered');

    // Get CSS filter string from Story 11.1
    const filterString = recipeToCSSFilters(currentRecipe);

    // Apply filter to image
    targetImage.style.filter = filterString;

    // Update CSS variable for slider
    sliderContainer.style.setProperty('--preview-filter', filterString);
  }
  ```
- [ ] Verify reference images load correctly (Story 11.2):
  - Portrait, Landscape, Product images present in `web/images/`
  - Responsive images load (400w, 800w, 1200w)
  - WebP format loads on modern browsers, JPEG fallback on older browsers

### Task 10: Manual Testing and Documentation
- [ ] Test modal interaction:
  - Open modal: Click "Preview" button
  - Tabs: Click Portrait, Landscape, Product → Images switch
  - Slider: Drag handle left/right → Before/after reveals smoothly
  - Keyboard: Tab, Arrow keys, Esc → Navigation works
  - Convert: Click "Convert Now" → Modal closes, conversion starts
  - Cancel: Click "Cancel" or overlay → Modal closes, no conversion
- [ ] Test mobile layout:
  - iPhone (Safari): Full-screen modal, touch slider, tap sides to snap
  - Android (Chrome): Full-screen modal, touch slider, swipe tabs
  - iPad (Safari): Tablet layout, responsive images load correctly
- [ ] Performance testing:
  - Lighthouse audit: Modal opens in <300ms
  - Slider interaction: 60fps minimum (no janky rendering)
  - CSS filters: Apply in <100ms (GPU accelerated)
- [ ] Browser compatibility:
  - Chrome (latest): ✅ Modal, slider, CSS filters work
  - Firefox (latest): ✅ Modal, slider, CSS filters work
  - Safari (latest): ✅ Modal, slider, CSS filters work
  - Edge (latest): ✅ Modal, slider, CSS filters work
- [ ] Add README section:
  ```markdown
  ### Preview Modal (Epic 11)

  Recipe includes a preview modal to visualize presets before conversion:

  #### Features

  - **Before/After Slider**: Drag to compare original vs. filtered image
  - **Reference Images**: Switch between Portrait, Landscape, Product tabs
  - **Preset Parameters**: View exact adjustments (Exposure, Contrast, etc.)
  - **Keyboard Accessible**: Tab, Arrow keys, Esc navigation
  - **Mobile Optimized**: Full-screen layout, touch-friendly controls

  #### Usage

  1. Upload a preset file (.np3, .xmp, .lrtemplate)
  2. Click "Preview Preset" button
  3. Drag slider to compare before/after
  4. Switch reference image tabs (Portrait, Landscape, Product)
  5. Click "Convert Now" to proceed with conversion
  ```

## Dev Notes

### Learnings from Previous Story

**From Story 11-2-reference-image-bundle (Status: drafted)**

Previous story not yet implemented. Story 11.3 builds the modal interface that displays the reference images from Story 11.2.

**Integration with Story 11.2:**
- Reference images must be present at `web/images/preview-*.webp` and `web/images/preview-*.jpg`
- Responsive images: 400w (mobile), 800w (tablet), 1200px (desktop)
- Three reference images: Portrait (1200x1600), Landscape (1200x800), Product (1200x1200)
- WebP format with JPEG fallback (<200 KB each)

**Technical Requirements:**
- Modal loads reference images on demand (lazy loading with `loading="lazy"`)
- Tabs switch between reference images (Portrait/Landscape/Product)
- Slider applies CSS filters to reference images in real-time

[Source: docs/stories/11-2-reference-image-bundle.md]

**Integration with Story 11.1:**
- CSS filters from Story 11.1 applied to filtered image in modal
- Function `recipeToCSSFilters(recipe)` must be available
- Filter string format: `brightness(150%) contrast(130%) saturate(150%)`
- Filters must render in <100ms (GPU accelerated)

[Source: docs/stories/11-1-css-filter-mapping.md]

### Architecture Alignment

**Tech Spec Epic 11 Alignment:**

Story 11.3 implements **AC-3: Preview Modal Interface** from tech-spec-epic-11.md.

**Modal Interface Requirements:**

```
Component          Function                                    Performance Target
-----------------  ------------------------------------------  -------------------
Modal open         Fade in animation, lock body scroll         <300ms
Before/After       Slider with clip-path reveal                60fps rendering
Tabs               Switch reference images (P/L/P)             <100ms
Parameters         Display preset values (E+0.7, C+15, W+10)   N/A (static text)
Convert button     Trigger conversion workflow                 Immediate
Close button       Close modal, return focus to Preview btn    <300ms
```

**Keyboard Navigation Map:**

```
Key          Action
-----------  --------------------------------------------------------
Tab          Cycle: Tabs → Slider → Convert → Cancel → Close (X)
Shift+Tab    Reverse cycle
Arrow ←/→    Switch tabs (when tab focused) OR move slider (when slider focused)
Enter/Space  Activate button or tab
Esc          Close modal
```

[Source: docs/tech-spec-epic-11.md#AC-3]

### Before/After Slider Implementation Strategy

**CSS clip-path vs. Mask:**

| Technique | Browser Support | Performance | Complexity | Recommendation |
| --------- | --------------- | ----------- | ---------- | -------------- |
| clip-path | 97%+ (IE 11 no) | GPU accel.  | Simple     | ✅ Recommended  |
| mask      | 95%+ (IE 11 no) | GPU accel.  | Medium     | Alternative    |
| width     | 100%            | CPU only    | Simple     | Fallback       |

**Recipe Implementation: clip-path**

```css
.preview-after {
  clip-path: inset(0 0 0 var(--slider-position, 50%));
  will-change: clip-path; /* GPU acceleration */
}
```

**Why clip-path over mask?**
- Simpler syntax: `clip-path: inset(0 0 0 50%)` vs. `mask: linear-gradient(...)`
- Better browser support: 97% vs. 95% (Safari prefixed)
- Easier to animate: Single CSS variable (`--slider-position`)
- GPU acceleration: `will-change: clip-path` hint

**Slider Position Update:**

```javascript
// Update CSS variable on drag
function updateSliderPosition(percentage) {
  sliderContainer.style.setProperty('--slider-position', `${percentage}%`);
}
```

**Performance:**
- clip-path changes trigger composite layer (GPU accelerated)
- No reflow or repaint (only compositing)
- 60fps target achievable on mid-range devices

[Source: CSS clip-path - MDN Web Docs]

### Focus Trap Implementation

**What is a Focus Trap?**

A focus trap ensures keyboard focus stays within the modal while open, preventing users from accidentally tabbing to background elements.

**Implementation:**

```javascript
function setupFocusTrap(modal) {
  const focusableElements = modal.querySelectorAll('button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])');
  const firstElement = focusableElements[0];
  const lastElement = focusableElements[focusableElements.length - 1];

  modal.addEventListener('keydown', (e) => {
    if (e.key === 'Tab') {
      if (e.shiftKey && document.activeElement === firstElement) {
        e.preventDefault();
        lastElement.focus(); // Wrap to last element
      } else if (!e.shiftKey && document.activeElement === lastElement) {
        e.preventDefault();
        firstElement.focus(); // Wrap to first element
      }
    }
  });
}
```

**Behavior:**
- Tab on last element → Focus first element (wrap around)
- Shift+Tab on first element → Focus last element (wrap around)
- All tabs stay within modal (no background interaction)

**Accessibility Impact:**
- WCAG 2.1 Level AA Success Criterion 2.4.3 (Focus Order)
- Screen reader users cannot accidentally navigate outside modal
- Keyboard-only users cannot lose focus context

[Source: ARIA Authoring Practices - Modal Dialog]

### Mobile Touch Interactions

**Touch Target Sizes (WCAG 2.1 Level AA):**

| Element        | Minimum Size | Recipe Size             | Compliant? |
| -------------- | ------------ | ----------------------- | ---------- |
| Slider handle  | 44x44px      | 60x60px                 | ✅ Yes      |
| Tab button     | 44x44px      | 48x48px                 | ✅ Yes      |
| Convert button | 44x44px      | 100% width, 48px height | ✅ Yes      |
| Close (X)      | 44x44px      | 48x48px                 | ✅ Yes      |

**Mobile Slider Handle:**
- Desktop: 60px diameter (larger than standard 44px)
- Mobile: 80px diameter (even larger for easier touch)
- Drag handle visible with contrasting border (white + blue)

**Tap Behavior (Snap to Edges):**

```javascript
// Tap left side → Snap to 0% (all original)
// Tap right side → Snap to 100% (all filtered)
sliderContainer.addEventListener('click', (e) => {
  const percentage = (e.clientX - rect.left) / rect.width * 100;

  if (percentage < 25) {
    sliderContainer.style.setProperty('--slider-position', '0%');
  } else if (percentage > 75) {
    sliderContainer.style.setProperty('--slider-position', '100%');
  }
});
```

**Why Snap Behavior?**
- Users want to see "full before" (0%) or "full after" (100%) quickly
- Dragging to exact 0% or 100% is difficult on mobile
- Tapping sides provides instant comparison (one tap vs. drag)

[Source: Touch Target Sizes - Material Design]

### ARIA Live Regions for Tab Changes

**Problem:** Screen readers don't announce tab changes when user clicks/taps a tab.

**Solution:** Create temporary ARIA live region with announcement.

```javascript
function announceTabChange(tabName) {
  const announcement = document.createElement('div');
  announcement.setAttribute('role', 'status');
  announcement.setAttribute('aria-live', 'polite');
  announcement.textContent = `${tabName} selected`;
  document.body.appendChild(announcement);

  // Remove after 1 second (announcement read by screen reader)
  setTimeout(() => announcement.remove(), 1000);
}
```

**ARIA Live Region Attributes:**
- `role="status"`: Non-critical announcement (polite)
- `aria-live="polite"`: Announce after current speech completes
- `textContent`: Plain text announcement (no HTML)

**Screen Reader Behavior:**
- User clicks "Landscape" tab
- Screen reader announces: "Landscape selected"
- User understands reference image changed

**Alternative: aria-label on tab button**
- Less dynamic (static label)
- Doesn't announce change event
- Live region preferred for real-time feedback

[Source: ARIA Live Regions - MDN Web Docs]

### Modal Animation Performance

**CSS Transitions:**

```css
.modal {
  opacity: 0;
  transition: opacity 300ms ease-in-out;
}

.modal:not([hidden]) {
  opacity: 1;
}
```

**Why 300ms?**
- Material Design recommendation: 200-400ms
- Fast enough to feel instant (not sluggish)
- Slow enough to perceive motion (not jarring)
- 300ms = sweet spot for modal animations

**GPU Acceleration:**

```css
.modal-container {
  transform: translateZ(0); /* Force GPU layer */
  will-change: opacity; /* Hint to browser */
}
```

**Performance Metrics:**
- Modal fade-in: <300ms (target)
- Slider drag: 60fps minimum (target)
- CSS filter application: <100ms (target, from Story 11.1)

**Lighthouse Audit:**
- "Avoid large layout shifts" (modal doesn't shift layout)
- "Minimize main-thread work" (CSS transitions offloaded to compositor)

[Source: CSS Animation Performance - web.dev]

### Project Structure Notes

**New Files Created (Story 11.3):**
```
web/
├── css/
│   └── modal.css                (Modal styles, slider styles)
├── js/
│   ├── modal.js                 (Open/close logic, focus trap)
│   ├── slider.js                (Slider drag interaction)
│   ├── tabs.js                  (Tab switching logic)
│   └── preview.js               (Parameter display, filter application)
└── tests/
    └── preview.test.js          (Unit tests for parameter formatting)
```

**Modified Files:**
- `web/index.html` - Add modal HTML, "Preview" button
- `web/js/upload.js` - Enable "Preview" button after file upload
- `web/js/conversion.js` - Trigger conversion from "Convert Now" button

**Integration Points:**
- Story 11.1: CSS filter function `recipeToCSSFilters(recipe)` must be available
- Story 11.2: Reference images must be present in `web/images/`
- Story 2.6: Conversion workflow triggered by "Convert Now" button
- Story 2.7: File download triggered after conversion completes
- Story 10.3: Progress indicator shown during conversion

[Source: docs/tech-spec-epic-11.md#Services-and-Modules]

### Testing Strategy

**Unit Tests (web/tests/preview.test.js):**
- `formatPresetParameters()`: Parameter formatting (positive/negative/zero values)
- Coverage target: 100% (simple string formatting function)

**Manual Tests:**
- Modal open/close: Click "Preview", click "Cancel", click overlay, press Esc
- Slider interaction: Drag handle, arrow keys, tap sides (mobile)
- Tab switching: Click tabs, arrow keys, swipe (mobile)
- Keyboard navigation: Tab through all elements, Esc to close
- Screen reader: NVDA, VoiceOver, TalkBack (announcements, focus order)
- Mobile layout: Full-screen modal, touch targets (44x44px minimum)

**Browser Compatibility:**
- Chrome 18+: ✅ clip-path, CSS filters, WebP
- Firefox 35+: ✅ clip-path, CSS filters, WebP
- Safari 9.1+: ✅ clip-path, CSS filters, WebP (Safari 14+)
- Edge 12+: ✅ clip-path, CSS filters, WebP

**Performance Tests:**
- Lighthouse audit: "Time to Interactive" <3 seconds
- DevTools Performance: Slider drag 60fps minimum
- Mobile device testing: iPhone 8, Android mid-range (3-year-old phones)

[Source: docs/tech-spec-epic-11.md#Test-Strategy-Summary]

### Known Risks

**RISK-58: Slider interaction may be janky on low-end devices**
- **Impact**: Slider drag feels sluggish, <60fps rendering
- **Mitigation**: GPU acceleration (`will-change: clip-path`), test on 3-year-old phones
- **Test**: iPhone 8 (2017), Android mid-range (2021)

**RISK-59: Focus trap may conflict with browser extensions**
- **Impact**: Browser extensions (password managers, screen readers) may break focus trap
- **Mitigation**: Use standard ARIA dialog pattern, test with common extensions
- **Acceptable**: Extensions expected to handle ARIA dialogs correctly

**RISK-60: Mobile swipe gestures may conflict with browser back/forward**
- **Impact**: Swiping tabs left/right triggers browser back/forward navigation
- **Mitigation**: Implement swipe only within tabs container (prevent propagation)
- **Test**: iOS Safari (back gesture), Android Chrome (back gesture)

**RISK-61: CSS clip-path not supported on IE 11**
- **Impact**: Before/after slider doesn't work on IE 11
- **Mitigation**: Fallback to width-based slider (less performant but functional)
- **Acceptable**: IE 11 usage <1%, Recipe targets modern browsers

[Source: docs/tech-spec-epic-11.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-11.md#AC-3] - Preview modal interface requirements
- [Source: docs/stories/11-1-css-filter-mapping.md] - CSS filter function integration
- [Source: docs/stories/11-2-reference-image-bundle.md] - Reference image integration
- [CSS clip-path - MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/clip-path)
- [ARIA Authoring Practices - Modal Dialog](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/)
- [Touch Target Sizes - Material Design](https://m3.material.io/foundations/accessible-design/accessibility-basics)
- [ARIA Live Regions - MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/ARIA_Live_Regions)
- [CSS Animation Performance - web.dev](https://web.dev/animations-guide/)
- [Modal Design Patterns - Nielsen Norman Group](https://www.nngroup.com/articles/modal-nonmodal-dialog/)

## Dev Agent Record

### Context Reference

- docs/stories/11-3-preview-modal-interface.context.xml (Generated: 2025-11-09)

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
