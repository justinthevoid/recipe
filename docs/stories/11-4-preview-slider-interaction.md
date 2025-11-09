# Story 11.4: Interactive Slider Enhancements (Keyboard, Touch, Accessibility)

Status: ready-for-dev

## Story

As a **photographer using keyboard-only navigation or touch devices**,
I want **enhanced slider interactions with keyboard shortcuts, touch gestures, and accessibility features**,
so that **I can efficiently compare before/after previews using my preferred input method**.

## Acceptance Criteria

**AC-1: Keyboard Shortcuts for Slider Control**
- ✅ Arrow keys move slider:
  - Left Arrow (←): Move slider left by 5% (show more "before")
  - Right Arrow (→): Move slider right by 5% (show more "after")
  - Repeat key hold: Continuous movement (smooth, not discrete jumps)
- ✅ Home/End keys for instant snap:
  - Home key: Snap slider to 0% (full "before")
  - End key: Snap slider to 100% (full "after")
- ✅ Number keys for quick positions:
  - 0 key: Snap to 0% (full before)
  - 5 key: Snap to 50% (center)
  - 1 key: Snap to 100% (full after)
- ✅ Spacebar toggle:
  - Press spacebar: Toggle between 0% and 100% (quick comparison)
  - Release spacebar: Return to previous position
- ✅ Page Up/Down for larger steps:
  - Page Up: Move slider left by 10%
  - Page Down: Move slider right by 10%

**AC-2: Touch Gestures for Mobile**
- ✅ Swipe gestures:
  - Swipe right: Move slider right (show more "after")
  - Swipe left: Move slider left (show more "before")
  - Swipe velocity: Faster swipe = larger movement
- ✅ Double-tap to snap:
  - Double-tap left side: Snap to 0% (full before)
  - Double-tap right side: Snap to 100% (full after)
  - Double-tap center: Snap to 50% (center)
- ✅ Pinch gesture (optional):
  - Pinch in: Zoom into image (for detail inspection)
  - Pinch out: Zoom out to full view
  - Note: Requires careful implementation to avoid conflicts with browser zoom
- ✅ Long press:
  - Long press slider handle: Show position tooltip ("50% filtered")
  - Release: Hide tooltip

**AC-3: Visual Feedback and Indicators**
- ✅ Slider position percentage display:
  - Numeric label: "50%" displayed above/below slider handle
  - Updates in real-time as slider moves
  - Fade in when dragging, fade out after 1 second of inactivity
- ✅ Before/After labels:
  - "Before" label on left edge of image
  - "After" label on right edge of image
  - Persistent (always visible) or toggle with "L" key
- ✅ Slider track indicators:
  - Subtle gradient: Left side (original) → Right side (filtered)
  - Optional: Tick marks at 0%, 25%, 50%, 75%, 100%
- ✅ Drag state visual:
  - Handle grows slightly when dragging (scale 1.1x)
  - Cursor changes to `ew-resize` when hovering handle
  - Handle shadow increases when dragging (visual depth)

**AC-4: Smooth Animations and Transitions**
- ✅ Slider movement animation:
  - Keyboard arrow keys: Smooth transition (150ms ease-out)
  - Snap to position (Home/End/Numbers): Smooth transition (300ms ease-in-out)
  - Drag: Instant (no animation, follow cursor)
- ✅ Handle scale animation:
  - Hover: Scale to 1.05x (100ms transition)
  - Active/Dragging: Scale to 1.1x (100ms transition)
  - Release: Return to 1.0x (150ms ease-out)
- ✅ Percentage label animation:
  - Fade in when dragging: 200ms fade-in
  - Fade out after 1s inactivity: 300ms fade-out
- ✅ Performance:
  - All animations: 60fps minimum (GPU accelerated)
  - No janky movement during drag or keyboard input

**AC-5: Accessibility Enhancements for Screen Readers**
- ✅ ARIA slider role:
  - `role="slider"` on slider handle
  - `aria-valuemin="0"`, `aria-valuemax="100"`, `aria-valuenow="50"`
  - `aria-valuetext="50% filtered"` (human-readable value)
  - `aria-label="Preview comparison slider"`
- ✅ Live region for position announcements:
  - ARIA live region announces position changes
  - Throttled announcements (every 10% change, not every 1%)
  - Example: "60% filtered" (announced when crossing 60% threshold)
- ✅ Instructions for screen reader users:
  - `aria-describedby` points to instructions element
  - Instructions text: "Use left and right arrow keys to adjust preview comparison. Press Home for full before, End for full after."
- ✅ Keyboard focus visible:
  - Clear focus ring when slider handle focused (2px solid blue)
  - Focus ring remains visible during keyboard interaction
  - Focus ring removed when mouse dragging (`:focus-visible` polyfill)

**AC-6: Multi-Touch Support for Tablets**
- ✅ Two-finger drag:
  - Two-finger drag left/right: Move slider (alternative to single-finger)
  - Prevents accidental page scroll on tablets
- ✅ Touch target size:
  - Slider handle: Minimum 60px diameter (desktop), 80px (mobile)
  - Exceeds WCAG 2.1 Level AAA (44x44px minimum)
- ✅ Touch feedback:
  - Haptic feedback on iOS/Android (if supported)
  - Visual feedback: Handle pulses when tapped
- ✅ Palm rejection:
  - Ignore palm touches (detect large touch area vs. finger)
  - Only respond to intentional finger touches

**AC-7: Performance Optimization for Slider Rendering**
- ✅ GPU acceleration:
  - `will-change: clip-path` on `.preview-after` element
  - `transform: translateZ(0)` to force GPU layer
- ✅ Debounced position updates:
  - Slider position updates: Max 60fps (16.67ms interval)
  - ARIA announcements: Throttled to 10% intervals
  - Percentage label updates: Real-time (no throttle)
- ✅ Lazy rendering:
  - CSS filters applied only to visible image (current tab)
  - Hidden tabs: Filters not applied until tab switch
- ✅ Performance targets:
  - Drag latency: <16ms (60fps)
  - Keyboard input latency: <50ms
  - Animation smoothness: 60fps minimum
  - Mobile (iPhone 8, 2017): 60fps minimum

## Tasks / Subtasks

### Task 1: Implement Keyboard Shortcuts (AC-1)
- [ ] Add keyboard event listeners in `web/js/slider.js`:
  ```javascript
  const sliderHandle = document.querySelector('.preview-slider-handle');

  sliderHandle.addEventListener('keydown', (e) => {
    const currentPos = parseInt(sliderHandle.getAttribute('aria-valuenow'));
    let newPos = currentPos;

    switch (e.key) {
      case 'ArrowLeft':
        newPos = Math.max(0, currentPos - 5);
        animateSliderTo(newPos, 150); // Smooth transition
        break;
      case 'ArrowRight':
        newPos = Math.min(100, currentPos + 5);
        animateSliderTo(newPos, 150);
        break;
      case 'Home':
        animateSliderTo(0, 300);
        break;
      case 'End':
        animateSliderTo(100, 300);
        break;
      case '0':
        animateSliderTo(0, 300);
        break;
      case '5':
        animateSliderTo(50, 300);
        break;
      case '1':
        animateSliderTo(100, 300);
        break;
      case 'PageUp':
        newPos = Math.max(0, currentPos - 10);
        animateSliderTo(newPos, 200);
        break;
      case 'PageDown':
        newPos = Math.min(100, currentPos + 10);
        animateSliderTo(newPos, 200);
        break;
      case ' ': // Spacebar
        e.preventDefault(); // Prevent page scroll
        if (!e.repeat) {
          // Store current position and toggle to opposite
          const togglePos = currentPos < 50 ? 100 : 0;
          sliderHandle.dataset.toggleReturn = currentPos;
          animateSliderTo(togglePos, 200);
        }
        break;
    }
  });

  sliderHandle.addEventListener('keyup', (e) => {
    if (e.key === ' ') {
      // Return to previous position
      const returnPos = parseInt(sliderHandle.dataset.toggleReturn || 50);
      animateSliderTo(returnPos, 200);
    }
  });

  // Smooth animation helper
  function animateSliderTo(targetPos, duration) {
    const startPos = parseInt(sliderHandle.getAttribute('aria-valuenow'));
    const startTime = performance.now();

    function animate(currentTime) {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const easeProgress = easeInOutCubic(progress);
      const newPos = startPos + (targetPos - startPos) * easeProgress;

      updateSliderPosition(Math.round(newPos));

      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    }

    requestAnimationFrame(animate);
  }

  function easeInOutCubic(t) {
    return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
  }
  ```

### Task 2: Implement Touch Gestures (AC-2)
- [ ] Add touch gesture detection in `web/js/gestures.js`:
  ```javascript
  const sliderContainer = document.querySelector('.preview-slider');
  let touchStartX = 0;
  let touchStartTime = 0;
  let lastTapTime = 0;

  // Swipe gesture
  sliderContainer.addEventListener('touchstart', (e) => {
    touchStartX = e.touches[0].clientX;
    touchStartTime = Date.now();
  });

  sliderContainer.addEventListener('touchmove', (e) => {
    const touchX = e.touches[0].clientX;
    const deltaX = touchX - touchStartX;
    const rect = sliderContainer.getBoundingClientRect();
    const percentage = (deltaX / rect.width) * 100;

    // Update slider based on swipe
    const currentPos = parseInt(sliderHandle.getAttribute('aria-valuenow'));
    const newPos = Math.max(0, Math.min(100, currentPos + percentage));
    updateSliderPosition(newPos);

    touchStartX = touchX; // Reset for next move
  });

  // Double-tap gesture
  sliderContainer.addEventListener('touchend', (e) => {
    const now = Date.now();
    const timeSinceLastTap = now - lastTapTime;

    if (timeSinceLastTap < 300) {
      // Double-tap detected
      const rect = sliderContainer.getBoundingClientRect();
      const tapX = e.changedTouches[0].clientX - rect.left;
      const tapPercent = (tapX / rect.width) * 100;

      // Snap to nearest position
      let snapPos = 50;
      if (tapPercent < 33) {
        snapPos = 0; // Left third → Snap to 0%
      } else if (tapPercent > 66) {
        snapPos = 100; // Right third → Snap to 100%
      }

      animateSliderTo(snapPos, 300);
    }

    lastTapTime = now;
  });

  // Long press gesture
  let longPressTimer = null;

  sliderHandle.addEventListener('touchstart', (e) => {
    longPressTimer = setTimeout(() => {
      // Show position tooltip
      showPositionTooltip();
    }, 500); // 500ms long press threshold
  });

  sliderHandle.addEventListener('touchend', () => {
    clearTimeout(longPressTimer);
    hidePositionTooltip();
  });

  function showPositionTooltip() {
    const tooltip = document.getElementById('slider-tooltip');
    const currentPos = sliderHandle.getAttribute('aria-valuenow');
    tooltip.textContent = `${currentPos}%`;
    tooltip.hidden = false;
  }

  function hidePositionTooltip() {
    const tooltip = document.getElementById('slider-tooltip');
    tooltip.hidden = true;
  }
  ```

### Task 3: Add Visual Feedback and Indicators (AC-3)
- [ ] Create percentage label HTML:
  ```html
  <!-- Add to modal HTML -->
  <div class="slider-percentage-label" id="slider-percentage" aria-hidden="true">50%</div>
  <div class="slider-tooltip" id="slider-tooltip" hidden>50%</div>
  ```
- [ ] Add percentage label CSS:
  ```css
  .slider-percentage-label {
    position: absolute;
    top: -30px;
    left: var(--slider-position, 50%);
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.8);
    color: white;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
    opacity: 0;
    transition: opacity 200ms ease-in-out;
    pointer-events: none;
  }

  .slider-percentage-label.visible {
    opacity: 1;
  }

  .preview-slider-before::before,
  .preview-slider-after::after {
    position: absolute;
    font-size: 14px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
    pointer-events: none;
  }

  .preview-slider-before::before {
    content: 'Before';
    top: 10px;
    left: 10px;
  }

  .preview-slider-after::after {
    content: 'After';
    top: 10px;
    right: 10px;
  }
  ```
- [ ] Update percentage label on slider move:
  ```javascript
  function updateSliderPosition(percentage) {
    // ... existing slider update logic ...

    // Update percentage label
    const percentLabel = document.getElementById('slider-percentage');
    percentLabel.textContent = `${Math.round(percentage)}%`;
    percentLabel.style.left = `${percentage}%`;

    // Show label when dragging
    percentLabel.classList.add('visible');

    // Hide label after 1 second of inactivity
    clearTimeout(window.percentLabelTimeout);
    window.percentLabelTimeout = setTimeout(() => {
      percentLabel.classList.remove('visible');
    }, 1000);
  }
  ```

### Task 4: Implement Smooth Animations (AC-4)
- [ ] Add CSS transitions for handle:
  ```css
  .preview-slider-handle {
    transition: transform 100ms ease-out, box-shadow 100ms ease-out;
  }

  .preview-slider-handle:hover {
    transform: translate(-50%, -50%) scale(1.05);
  }

  .preview-slider-handle:active,
  .preview-slider-handle.dragging {
    transform: translate(-50%, -50%) scale(1.1);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
  }
  ```
- [ ] Add dragging state class:
  ```javascript
  sliderHandle.addEventListener('mousedown', () => {
    sliderHandle.classList.add('dragging');
  });

  document.addEventListener('mouseup', () => {
    sliderHandle.classList.remove('dragging');
  });
  ```

### Task 5: Accessibility Enhancements (AC-5)
- [ ] Add ARIA attributes to slider HTML:
  ```html
  <div class="preview-slider-handle"
       role="slider"
       aria-valuemin="0"
       aria-valuemax="100"
       aria-valuenow="50"
       aria-valuetext="50% filtered"
       aria-label="Preview comparison slider"
       aria-describedby="slider-instructions"
       tabindex="0">
  </div>

  <div id="slider-instructions" class="sr-only">
    Use left and right arrow keys to adjust preview comparison. Press Home for full before, End for full after.
  </div>

  <div role="status" aria-live="polite" aria-atomic="true" class="sr-only" id="slider-live-region"></div>
  ```
- [ ] Implement throttled ARIA announcements:
  ```javascript
  let lastAnnouncedPosition = 50;

  function updateSliderPosition(percentage) {
    // ... existing update logic ...

    // Update ARIA attributes
    sliderHandle.setAttribute('aria-valuenow', Math.round(percentage));
    sliderHandle.setAttribute('aria-valuetext', `${Math.round(percentage)}% filtered`);

    // Announce position changes every 10%
    const roundedPos = Math.round(percentage / 10) * 10;
    if (roundedPos !== lastAnnouncedPosition) {
      const liveRegion = document.getElementById('slider-live-region');
      liveRegion.textContent = `${roundedPos}% filtered`;
      lastAnnouncedPosition = roundedPos;
    }
  }
  ```
- [ ] Add visible focus ring CSS:
  ```css
  .preview-slider-handle:focus-visible {
    outline: 2px solid var(--color-primary);
    outline-offset: 4px;
  }

  /* Hide focus ring when dragging with mouse */
  .preview-slider-handle:focus:not(:focus-visible) {
    outline: none;
  }
  ```

### Task 6: Multi-Touch Support (AC-6)
- [ ] Detect multi-touch gestures:
  ```javascript
  sliderContainer.addEventListener('touchmove', (e) => {
    if (e.touches.length === 2) {
      // Two-finger drag
      const touch1X = e.touches[0].clientX;
      const touch2X = e.touches[1].clientX;
      const averageX = (touch1X + touch2X) / 2;

      const rect = sliderContainer.getBoundingClientRect();
      const percentage = ((averageX - rect.left) / rect.width) * 100;
      updateSliderPosition(Math.max(0, Math.min(100, percentage)));

      e.preventDefault(); // Prevent page scroll
    }
  });
  ```
- [ ] Add haptic feedback (iOS/Android):
  ```javascript
  function triggerHapticFeedback() {
    if ('vibrate' in navigator) {
      navigator.vibrate(10); // 10ms vibration
    }

    // iOS Taptic Engine (requires iOS 13+)
    if (window.Taptic && typeof window.Taptic.impact === 'function') {
      window.Taptic.impact({ style: 'light' });
    }
  }

  // Trigger haptic when slider snaps to position
  function snapToPosition(position) {
    animateSliderTo(position, 300);
    triggerHapticFeedback();
  }
  ```

### Task 7: Performance Optimization (AC-7)
- [ ] Add GPU acceleration CSS:
  ```css
  .preview-after {
    will-change: clip-path;
    transform: translateZ(0); /* Force GPU layer */
  }

  .preview-slider-handle {
    will-change: transform;
  }
  ```
- [ ] Debounce slider position updates:
  ```javascript
  let rafId = null;

  function requestSliderUpdate(percentage) {
    if (rafId) {
      cancelAnimationFrame(rafId);
    }

    rafId = requestAnimationFrame(() => {
      updateSliderPosition(percentage);
      rafId = null;
    });
  }

  // Use requestSliderUpdate instead of direct updateSliderPosition
  sliderContainer.addEventListener('mousemove', (e) => {
    if (isDragging) {
      const percentage = calculatePercentage(e);
      requestSliderUpdate(percentage); // Debounced to 60fps
    }
  });
  ```
- [ ] Performance monitoring:
  ```javascript
  // Track drag performance
  let dragFrames = 0;
  let dragStartTime = 0;

  sliderHandle.addEventListener('mousedown', () => {
    dragFrames = 0;
    dragStartTime = performance.now();
  });

  document.addEventListener('mousemove', () => {
    if (isDragging) {
      dragFrames++;
    }
  });

  document.addEventListener('mouseup', () => {
    if (isDragging) {
      const dragDuration = performance.now() - dragStartTime;
      const fps = (dragFrames / dragDuration) * 1000;
      console.log(`Slider drag performance: ${fps.toFixed(1)} fps`);

      // Warn if performance below target
      if (fps < 60) {
        console.warn(`Slider performance below target (${fps.toFixed(1)} fps < 60 fps)`);
      }
    }
  });
  ```

### Task 8: Unit Tests for Slider Interactions
- [ ] Create `web/tests/slider.test.js`:
  ```javascript
  describe('Slider Keyboard Shortcuts', () => {
    it('moves slider left with arrow key', () => {
      sliderHandle.setAttribute('aria-valuenow', 50);
      simulateKeypress('ArrowLeft');
      expect(sliderHandle.getAttribute('aria-valuenow')).toBe('45');
    });

    it('snaps to 0% with Home key', () => {
      sliderHandle.setAttribute('aria-valuenow', 50);
      simulateKeypress('Home');
      setTimeout(() => {
        expect(sliderHandle.getAttribute('aria-valuenow')).toBe('0');
      }, 350); // Wait for animation
    });

    it('toggles to opposite with Spacebar', () => {
      sliderHandle.setAttribute('aria-valuenow', 30);
      simulateKeypress(' ', 'keydown');
      setTimeout(() => {
        expect(sliderHandle.getAttribute('aria-valuenow')).toBe('100');
      }, 250);
    });
  });

  describe('ARIA Announcements', () => {
    it('announces position every 10%', () => {
      updateSliderPosition(50);
      expect(liveRegion.textContent).toBe('50% filtered');

      updateSliderPosition(55); // No announcement (within 10% threshold)
      expect(liveRegion.textContent).toBe('50% filtered');

      updateSliderPosition(60); // Announcement (crossed 10% threshold)
      expect(liveRegion.textContent).toBe('60% filtered');
    });
  });
  ```

### Task 9: Manual Testing
- [ ] Keyboard testing:
  - Arrow keys: Smooth movement in 5% increments
  - Home/End: Instant snap to 0%/100%
  - Number keys (0, 5, 1): Snap to 0%, 50%, 100%
  - Spacebar: Toggle between 0% and 100%
  - Page Up/Down: Movement in 10% increments
- [ ] Touch testing (mobile):
  - Swipe left/right: Slider moves smoothly
  - Double-tap sides: Snaps to 0% or 100%
  - Long press handle: Tooltip appears
  - Two-finger drag: Alternative slider control
- [ ] Accessibility testing:
  - Screen reader (NVDA): Announces position changes every 10%
  - Keyboard focus: Visible focus ring on slider handle
  - ARIA attributes: Correct aria-valuenow, aria-valuetext updates
- [ ] Performance testing:
  - Drag performance: Monitor FPS (DevTools Performance panel)
  - Target: 60fps minimum on desktop, 60fps on mobile (iPhone 8)
  - Animation smoothness: No janky transitions

### Task 10: Documentation
- [ ] Add keyboard shortcuts to README:
  ```markdown
  ### Slider Keyboard Shortcuts

  | Key          | Action                             |
  | ------------ | ---------------------------------- |
  | ← →          | Move slider 5% (smooth transition) |
  | Home         | Snap to 0% (full before)           |
  | End          | Snap to 100% (full after)          |
  | 0, 5, 1      | Snap to 0%, 50%, 100%              |
  | Page Up/Down | Move slider 10%                    |
  | Spacebar     | Toggle between 0% and 100% (hold)  |
  ```
- [ ] Add touch gestures to README:
  ```markdown
  ### Touch Gestures

  | Gesture          | Action                                  |
  | ---------------- | --------------------------------------- |
  | Swipe left/right | Move slider                             |
  | Double-tap       | Snap to 0%/50%/100% (based on position) |
  | Long press       | Show position tooltip                   |
  | Two-finger drag  | Alternative slider control              |
  ```

## Dev Notes

### Learnings from Previous Story

**From Story 11-3-preview-modal-interface (Status: drafted)**

Previous story not yet implemented. Story 11.4 enhances the slider interaction features introduced in Story 11.3.

**Story 11.3 Coverage:**
- Basic slider structure (before/after with clip-path)
- Mouse drag interaction
- Touch drag interaction
- Basic keyboard support (Arrow keys)
- Slider position at 50% default

**Story 11.4 Enhancements:**
- Advanced keyboard shortcuts (Home/End, Number keys, Spacebar toggle, Page Up/Down)
- Touch gestures (Swipe, Double-tap, Long press, Multi-touch)
- Visual feedback (Percentage label, Before/After labels, Drag state)
- Smooth animations (Handle scale, Snap transitions, Fade in/out)
- Accessibility enhancements (ARIA live regions, Throttled announcements, Focus visible)
- Performance optimizations (GPU acceleration, Debounced updates, RAF)

[Source: docs/stories/11-3-preview-modal-interface.md]

### Architecture Alignment

**Tech Spec Epic 11 Alignment:**

Story 11.4 provides **enhancements** to AC-3 (Preview Modal Interface) focusing on advanced slider interactions.

**Slider Interaction Matrix:**

```
Input Method    Basic (Story 11.3)         Enhanced (Story 11.4)
--------------  -------------------------  -----------------------------------
Mouse           Drag handle                + Hover scale, Active scale
Touch           Tap and drag               + Swipe, Double-tap, Long press, Multi-touch
Keyboard        Arrow keys (←/→)           + Home/End, 0/5/1, Page Up/Down, Spacebar
Accessibility   ARIA valuenow              + Live region, Throttled announcements
Visual          Handle, Divider line       + Percentage label, Before/After labels, Animations
Performance     Basic rendering            + GPU acceleration, Debounced RAF updates
```

[Source: docs/tech-spec-epic-11.md#AC-3]

### Keyboard Shortcuts Design Rationale

**Why These Specific Shortcuts?**

| Shortcut     | Rationale                                               | Precedent                    |
| ------------ | ------------------------------------------------------- | ---------------------------- |
| Arrow keys   | Universal slider control (5% increments = fine control) | Media players, image editors |
| Home/End     | Standard "go to start/end" convention                   | Text editors, browsers       |
| 0, 5, 1      | Quick numeric positions (0%/50%/100%)                   | Lightroom (0-9 for ratings)  |
| Spacebar     | Toggle preview (common in photo editors)                | Photoshop (before/after)     |
| Page Up/Down | Larger steps (10% = coarse control)                     | Document viewers, scrolling  |

**Rationale for 5% vs. 1% Increments:**
- 1% too granular: Requires 100 key presses to traverse slider (tedious)
- 10% too coarse: Only 10 positions, limited precision
- 5% sweet spot: 20 key presses to traverse, sufficient precision

[Source: Keyboard Interaction Patterns - Nielsen Norman Group]

### Touch Gesture Design

**Swipe vs. Drag:**

| Interaction | Description                         | Use Case                         |
| ----------- | ----------------------------------- | -------------------------------- |
| Drag        | Continuous finger contact on handle | Precise control (existing, 11.3) |
| Swipe       | Flick gesture anywhere on image     | Quick adjustments (new, 11.4)    |

**Why Swipe in Addition to Drag?**
- Faster: Flick gesture quicker than dragging handle
- Larger touch area: Entire image vs. small handle (60-80px)
- Familiar: Swipe gestures common in mobile apps (galleries, carousels)

**Double-Tap Snap Zones:**

```
|←  0%  →|← 50% →|← 100% →|
 ←33.3%→  33-66%  ←66.7%→

Tap left third → Snap to 0%
Tap center third → Snap to 50%
Tap right third → Snap to 100%
```

[Source: Mobile Touch Gestures - Material Design]

### ARIA Live Region Throttling

**Problem:** Announcing every 1% change overwhelms screen readers.

**Solution:** Throttle announcements to 10% intervals.

```javascript
// Bad: Announces every 1% (50 announcements from 0% to 50%)
function updateSliderPosition(percentage) {
  liveRegion.textContent = `${percentage}% filtered`; // Too frequent
}

// Good: Announces every 10% (5 announcements from 0% to 50%)
function updateSliderPosition(percentage) {
  const roundedPos = Math.round(percentage / 10) * 10;
  if (roundedPos !== lastAnnouncedPosition) {
    liveRegion.textContent = `${roundedPos}% filtered`;
    lastAnnouncedPosition = roundedPos;
  }
}
```

**Benefits:**
- Reduces cognitive load: Users hear "0%, 10%, 20%, 30%, 40%, 50%" vs. "0%, 1%, 2%, 3%... 50%"
- Prevents announcement queue overflow
- Screen reader remains responsive

[Source: WCAG 2.1 SC 4.1.3 - Status Messages]

### GPU Acceleration for Smooth Rendering

**CSS Properties and GPU Acceleration:**

| Property    | GPU Accelerated? | Use Case                         |
| ----------- | ---------------- | -------------------------------- |
| `clip-path` | ✅ Yes            | Before/after reveal (Story 11.3) |
| `transform` | ✅ Yes            | Handle scale animation           |
| `opacity`   | ✅ Yes            | Percentage label fade            |
| `width`     | ❌ No             | Avoid for slider positioning     |
| `left`      | ❌ No             | Avoid for slider positioning     |

**Why `will-change`?**

```css
.preview-after {
  will-change: clip-path; /* Browser hint: "This will change frequently" */
}
```

**Effect:**
- Browser promotes element to GPU layer (compositing)
- Changes handled by compositor thread (not main thread)
- Result: 60fps rendering even during JavaScript execution

**Warning:** Overuse of `will-change` causes memory overhead. Use sparingly.

[Source: CSS will-change - MDN Web Docs]

### RequestAnimationFrame for Smooth Animations

**Problem:** `setTimeout` not synced with browser repaint (janky animations).

**Solution:** `requestAnimationFrame` syncs with 60fps repaint cycle.

```javascript
// Bad: setTimeout (not synced with repaint)
function animateSlider(targetPos, duration) {
  const fps = 60;
  const interval = 1000 / fps;
  const steps = duration / interval;

  let currentStep = 0;
  const timer = setInterval(() => {
    currentStep++;
    const progress = currentStep / steps;
    updateSliderPosition(startPos + (targetPos - startPos) * progress);

    if (currentStep >= steps) clearInterval(timer);
  }, interval); // May not align with repaint
}

// Good: requestAnimationFrame (synced with repaint)
function animateSlider(targetPos, duration) {
  const startTime = performance.now();

  function animate(currentTime) {
    const elapsed = currentTime - startTime;
    const progress = Math.min(elapsed / duration, 1);
    updateSliderPosition(startPos + (targetPos - startPos) * progress);

    if (progress < 1) {
      requestAnimationFrame(animate); // Next frame
    }
  }

  requestAnimationFrame(animate);
}
```

**Benefits:**
- Synced with browser repaint (60fps)
- Pauses when tab inactive (saves CPU)
- More efficient than `setTimeout`/`setInterval`

[Source: requestAnimationFrame - web.dev]

### Haptic Feedback on Mobile

**iOS Taptic Engine:**

```javascript
// Check if available (iOS 13+)
if (window.Taptic && typeof window.Taptic.impact === 'function') {
  window.Taptic.impact({ style: 'light' }); // Light vibration
}
```

**Android Vibration API:**

```javascript
if ('vibrate' in navigator) {
  navigator.vibrate(10); // 10ms vibration
}
```

**Use Cases:**
- Slider snaps to position (0%, 50%, 100%)
- Handle reaches edge (0% or 100%)
- Double-tap gesture detected

**Benefits:**
- Physical feedback improves perceived responsiveness
- Confirms action without visual confirmation
- Enhances mobile user experience

**Caution:** Overuse of haptics can annoy users. Use sparingly for significant events only.

[Source: Haptic Feedback - Apple Human Interface Guidelines]

### Performance Monitoring in Production

**FPS Tracking During Drag:**

```javascript
let dragFrames = 0;
let dragStartTime = 0;

sliderHandle.addEventListener('mousedown', () => {
  dragFrames = 0;
  dragStartTime = performance.now();
});

document.addEventListener('mousemove', () => {
  if (isDragging) dragFrames++;
});

document.addEventListener('mouseup', () => {
  if (isDragging) {
    const fps = (dragFrames / (performance.now() - dragStartTime)) * 1000;
    console.log(`Drag FPS: ${fps.toFixed(1)}`);

    // Analytics (optional)
    if (fps < 60) {
      trackPerformanceIssue('slider-drag', { fps });
    }
  }
});
```

**Benefits:**
- Identify performance regressions in production
- Track device-specific performance (iPhone 8 vs. iPhone 15)
- Inform optimization priorities

[Source: Performance Monitoring - web.dev]

### Project Structure Notes

**New Files Created (Story 11.4):**
```
web/
├── js/
│   ├── gestures.js              (Touch gestures: swipe, double-tap, long press)
│   └── slider-animations.js     (Animation helpers: easeInOutCubic, animateSliderTo)
└── tests/
    └── slider.test.js           (Unit tests for keyboard shortcuts, ARIA)
```

**Modified Files:**
- `web/js/slider.js` - Add keyboard shortcuts, performance optimizations
- `web/css/modal.css` - Add percentage label, Before/After labels, animations
- `web/index.html` - Add percentage label, tooltip, ARIA live region

**Integration Points:**
- Story 11.3: Slider structure (builds on existing implementation)
- Story 11.1: CSS filters (applied to slider images)
- Story 11.2: Reference images (used in slider)

[Source: docs/tech-spec-epic-11.md#Services-and-Modules]

### Testing Strategy

**Unit Tests (web/tests/slider.test.js):**
- Keyboard shortcuts: Arrow keys, Home/End, Number keys, Spacebar, Page Up/Down
- ARIA announcements: Throttled to 10% intervals
- Coverage target: 90% (complex interaction logic)

**Manual Tests:**
- Keyboard: All shortcuts work, animations smooth
- Touch: Swipe, double-tap, long press, multi-touch
- Accessibility: Screen reader announces position, focus visible
- Performance: 60fps during drag (DevTools Performance panel)
- Mobile: Test on iPhone 8 (2017), Android mid-range (2021)

**Browser Compatibility:**
- Chrome 18+: ✅ All features supported
- Firefox 35+: ✅ All features supported
- Safari 9.1+: ✅ All features supported (Taptic Engine iOS 13+)
- Edge 12+: ✅ All features supported

**Performance Targets:**
- Drag latency: <16ms (60fps)
- Keyboard latency: <50ms
- Animation smoothness: 60fps minimum
- Mobile (iPhone 8): 60fps minimum

[Source: docs/tech-spec-epic-11.md#Test-Strategy-Summary]

### Known Risks

**RISK-62: Spacebar toggle may conflict with page scroll**
- **Impact**: Pressing spacebar scrolls page instead of toggling slider
- **Mitigation**: `e.preventDefault()` when slider focused
- **Test**: Verify spacebar doesn't scroll page when slider has focus

**RISK-63: Haptic feedback not supported on all devices**
- **Impact**: No haptic feedback on older devices or browsers
- **Mitigation**: Feature detection, graceful degradation
- **Acceptable**: Haptic feedback is enhancement, not required

**RISK-64: Multi-touch gestures may conflict with browser zoom**
- **Impact**: Two-finger drag triggers browser zoom instead of slider control
- **Mitigation**: Careful touch event handling, test on real devices
- **Test**: iPhone (Safari), Android (Chrome)

**RISK-65: Performance degradation on low-end devices**
- **Impact**: Slider drag <60fps on older phones
- **Mitigation**: GPU acceleration, debounced updates, test on iPhone 8 (2017)
- **Acceptable**: Target 60fps, acceptable 45fps on very old devices (<5 years)

[Source: docs/tech-spec-epic-11.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-11.md#AC-3] - Preview modal interface (slider requirements)
- [Source: docs/stories/11-3-preview-modal-interface.md] - Basic slider implementation
- [Keyboard Interaction Patterns - Nielsen Norman Group](https://www.nngroup.com/articles/keyboard-accessibility/)
- [Mobile Touch Gestures - Material Design](https://m3.material.io/foundations/interaction/gestures)
- [WCAG 2.1 SC 4.1.3 - Status Messages](https://www.w3.org/WAI/WCAG21/Understanding/status-messages.html)
- [CSS will-change - MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/will-change)
- [requestAnimationFrame - web.dev](https://web.dev/requestanimationframe/)
- [Haptic Feedback - Apple Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/playing-haptics)
- [Performance Monitoring - web.dev](https://web.dev/performance-budgets-101/)

## Dev Agent Record

### Context Reference

- docs/stories/11-4-preview-slider-interaction.context.xml

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
