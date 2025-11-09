# Story 10.5: Responsive Mobile Design and Touch Interactions

Status: ready-for-dev

## Story

As a **photographer using Recipe on mobile devices**,
I want **a fully responsive interface optimized for mobile, tablet, and desktop screens**,
so that **I can convert presets on any device without usability issues like tiny buttons or horizontal scrolling**.

## Acceptance Criteria

**AC-1: Mobile Layout (<768px)**
- ✅ Single-column stacked layout for all sections
- ✅ Hero section readable without zooming:
  - Title font size: 28px (down from 48px desktop)
  - Subtitle font size: 16px (minimum readable size)
  - CTA button full-width (100% width on mobile)
- ✅ Format badges stacked vertically (1 badge per row)
- ✅ Upload drop zone full-width with vertical padding
- ✅ File cards stacked (1 card per row, 100% width)
- ✅ Tap-to-browse upload (file picker triggered by tap)
- ✅ No drag-and-drop on mobile (iOS/Android limitations acknowledged)

**AC-2: Tablet Layout (768px-1023px)**
- ✅ Two-column grid for file cards
- ✅ Format badges in 2×3 grid (2 columns, 3 rows)
- ✅ Hero section optimized:
  - Title font size: 36px
  - Subtitle font size: 18px
- ✅ Upload drop zone spans full width (2 columns)
- ✅ Batch controls (format selector, "Convert All" button) side-by-side

**AC-3: Desktop Layout (>1024px)**
- ✅ Three-column grid for file cards
- ✅ Format badges in horizontal row (5 badges, 1 row)
- ✅ Hero section full-size:
  - Title font size: 48px
  - Subtitle font size: 20px
- ✅ Upload drop zone with horizontal padding
- ✅ Batch controls side-by-side with spacing

**AC-4: Touch-Friendly Tap Targets on Mobile**
- ✅ All interactive elements minimum 44px height (iOS guidelines)
- ✅ Buttons full-width on mobile or minimum 44px × 44px
- ✅ File card "Convert" button: 48px height
- ✅ File card "Remove" button: 44px height
- ✅ Format dropdown: 48px height
- ✅ Batch "Convert All" button: 56px height (primary action)
- ✅ Spacing between tap targets: minimum 8px gap

**AC-5: No Horizontal Scrolling on Any Device**
- ✅ Viewport meta tag set correctly: `<meta name="viewport" content="width=device-width, initial-scale=1">`
- ✅ All content fits within viewport width (no overflow-x)
- ✅ Long filenames wrap or truncate (no fixed-width containers)
- ✅ Images/SVGs scale to container width (max-width: 100%)
- ✅ Horizontal padding accounts for mobile screen edges (16px min)
- ✅ Test on smallest device: iPhone SE (320px width)

**AC-6: Readable Text Without Zooming**
- ✅ Body text minimum 16px font size (all devices)
- ✅ Small text (file size, error messages) minimum 14px
- ✅ Line height minimum 1.5 (WCAG AA compliance)
- ✅ Paragraph max-width: 65ch (optimal readability)
- ✅ Sufficient color contrast: 4.5:1 for body text, 3:1 for large text (WCAG AA)
- ✅ No text-rendering issues on iOS Safari (test with actual device)

**AC-7: Manual Testing on Real Devices**
- ✅ **iPhone Testing** (iOS Safari):
  - iPhone 13 Pro (390px width) - portrait and landscape
  - iPhone SE (320px width) - portrait (smallest common device)
  - Test tap targets, file upload, conversion flow
- ✅ **Android Testing** (Chrome):
  - Samsung Galaxy S21 (360px width) - portrait and landscape
  - Google Pixel 6 (393px width) - portrait
  - Test tap targets, file upload, conversion flow
- ✅ **iPad Testing** (Safari):
  - iPad Pro 11" (834px width) - portrait and landscape
  - iPad Air (820px width) - portrait
  - Test 2-column grid, touch targets
- ✅ **Desktop Browser Testing**:
  - Chrome (1920px width) - 3-column grid
  - Firefox (1920px width) - responsive design tools
  - Safari (1440px width) - macOS
  - Edge (1920px width) - Windows

## Tasks / Subtasks

### Task 1: Add Viewport Meta Tag and Mobile-First CSS Reset (AC-5, AC-6)
- [ ] Add viewport meta tag to `web/index.html`:
  ```html
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Recipe - Convert Photo Presets</title>
    <!-- ... -->
  </head>
  ```
- [ ] Add mobile-first CSS reset to `web/css/main.css`:
  ```css
  /* Mobile-first responsive design */
  * {
    box-sizing: border-box;
  }

  html {
    font-size: 16px; /* Minimum readable size */
    line-height: 1.5; /* WCAG AA compliance */
  }

  body {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    color: var(--color-text);
    background: var(--color-bg);
    overflow-x: hidden; /* Prevent horizontal scrolling */
  }

  img, svg {
    max-width: 100%; /* Scale images to container */
    height: auto;
  }

  p {
    max-width: 65ch; /* Optimal readability */
    margin: 0 0 1rem 0;
  }
  ```

### Task 2: Implement Mobile Hero Section (<768px) (AC-1)
- [ ] Add mobile-first hero styles to `web/css/layout.css`:
  ```css
  /* Mobile-first hero section */
  .hero {
    padding: 2rem 1rem; /* Vertical padding, horizontal edge padding */
    text-align: center;
  }

  .hero__title {
    font-size: 28px; /* Mobile size */
    line-height: 1.2;
    margin: 0 0 1rem 0;
  }

  .hero__subtitle {
    font-size: 16px; /* Minimum readable */
    line-height: 1.5;
    margin: 0 0 1.5rem 0;
  }

  .hero__cta {
    display: block;
    width: 100%; /* Full-width button on mobile */
    padding: 1rem 2rem;
    font-size: 18px;
    text-align: center;
  }

  /* Tablet (768px+): Increase font sizes */
  @media (min-width: 768px) {
    .hero {
      padding: 3rem 2rem;
    }

    .hero__title {
      font-size: 36px;
    }

    .hero__subtitle {
      font-size: 18px;
    }

    .hero__cta {
      width: auto; /* Auto-width button on tablet+ */
      display: inline-block;
    }
  }

  /* Desktop (1024px+): Full size */
  @media (min-width: 1024px) {
    .hero {
      padding: 4rem 2rem;
    }

    .hero__title {
      font-size: 48px;
    }

    .hero__subtitle {
      font-size: 20px;
    }
  }
  ```

### Task 3: Implement Responsive Format Badges (AC-1, AC-2, AC-3)
- [ ] Add responsive badge grid to `web/css/components.css`:
  ```css
  /* Mobile: Stacked badges (1 column) */
  .formats__badges {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    max-width: 400px;
    margin: 0 auto;
  }

  .badge {
    display: block;
    width: 100%;
    padding: 1rem;
    text-align: center;
    font-size: 18px;
    font-weight: var(--font-weight-bold);
  }

  /* Tablet (768px+): 2-column grid */
  @media (min-width: 768px) {
    .formats__badges {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      max-width: 600px;
    }

    .badge {
      font-size: 16px;
    }
  }

  /* Desktop (1024px+): Horizontal row */
  @media (min-width: 1024px) {
    .formats__badges {
      display: flex;
      flex-direction: row;
      justify-content: center;
      max-width: none;
      gap: 1.5rem;
    }

    .badge {
      width: auto;
      padding: 0.75rem 1.5rem;
      font-size: 16px;
    }
  }
  ```

### Task 4: Implement Responsive File Card Grid (AC-1, AC-2, AC-3)
- [ ] Update file card grid styles in `web/css/components.css`:
  ```css
  /* Mobile: 1 column (default) */
  .upload__files {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1rem;
    padding: 0 1rem; /* Edge padding */
  }

  /* Tablet (768px+): 2 columns */
  @media (min-width: 768px) {
    .upload__files {
      grid-template-columns: repeat(2, 1fr);
      gap: 1.5rem;
      padding: 0 2rem;
    }
  }

  /* Desktop (1024px+): 3 columns */
  @media (min-width: 1024px) {
    .upload__files {
      grid-template-columns: repeat(3, 1fr);
      gap: 2rem;
      max-width: 1400px;
      margin: 0 auto;
      padding: 0 2rem;
    }
  }
  ```

### Task 5: Implement Touch-Friendly Tap Targets (AC-4)
- [ ] Update button styles for mobile tap targets in `web/css/components.css`:
  ```css
  /* Mobile-first button styles */
  .button {
    min-height: 44px; /* iOS guideline */
    padding: 0.75rem 1.5rem;
    font-size: 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .button--primary {
    background: var(--color-primary);
    color: white;
  }

  .button--secondary {
    background: var(--color-secondary);
    color: var(--color-text);
  }

  /* File card buttons */
  .file-card__convert {
    min-height: 48px; /* Larger tap target for primary action */
    width: 100%; /* Full-width on mobile */
    margin-top: 0.5rem;
  }

  .file-card__remove {
    min-height: 44px;
    width: 100%;
    margin-top: 0.5rem;
  }

  .file-card__download {
    min-height: 48px;
    width: 100%;
  }

  /* Batch controls */
  #convert-all-btn {
    min-height: 56px; /* Extra large for primary batch action */
    font-size: 18px;
    width: 100%; /* Full-width on mobile */
  }

  /* Tablet (768px+): Buttons auto-width */
  @media (min-width: 768px) {
    .file-card__convert,
    .file-card__remove,
    .file-card__download,
    #convert-all-btn {
      width: auto;
    }

    .file-card__convert,
    .file-card__remove {
      display: inline-block;
      margin-right: 0.5rem;
    }
  }

  /* Format dropdown */
  .file-card__format-select {
    min-height: 48px;
    font-size: 16px;
    padding: 0.75rem;
  }

  /* Spacing between tap targets */
  .file-card__footer {
    display: flex;
    flex-direction: column;
    gap: 8px; /* Minimum spacing */
  }

  @media (min-width: 768px) {
    .file-card__footer {
      flex-direction: row;
      gap: 12px;
    }
  }
  ```

### Task 6: Implement Mobile Upload Drop Zone (AC-1)
- [ ] Update drop zone styles for mobile in `web/css/components.css`:
  ```css
  /* Mobile drop zone */
  .upload__dropzone {
    border: 2px dashed var(--color-border);
    border-radius: 8px;
    padding: 2rem 1rem; /* Reduced padding on mobile */
    text-align: center;
    background: var(--color-bg-light);
    cursor: pointer;
    min-height: 200px; /* Ensure touch-friendly size */
  }

  .upload__icon {
    width: 48px; /* Smaller icon on mobile */
    height: 48px;
    margin: 0 auto 1rem;
  }

  .upload__dropzone-title {
    font-size: 18px; /* Smaller on mobile */
    margin: 0 0 0.5rem 0;
  }

  .upload__dropzone-subtitle {
    font-size: 14px;
    margin: 0 0 1rem 0;
  }

  #browse-files {
    width: 100%; /* Full-width button on mobile */
    min-height: 48px;
  }

  /* Tablet (768px+): Increase size */
  @media (min-width: 768px) {
    .upload__dropzone {
      padding: 3rem 2rem;
      min-height: 250px;
    }

    .upload__icon {
      width: 56px;
      height: 56px;
    }

    .upload__dropzone-title {
      font-size: 20px;
    }

    .upload__dropzone-subtitle {
      font-size: 16px;
    }

    #browse-files {
      width: auto;
    }
  }

  /* Desktop (1024px+): Full size */
  @media (min-width: 1024px) {
    .upload__dropzone {
      padding: 4rem 2rem;
      min-height: 300px;
    }

    .upload__icon {
      width: 64px;
      height: 64px;
    }

    .upload__dropzone-title {
      font-size: 24px;
    }

    .upload__dropzone-subtitle {
      font-size: 18px;
    }
  }
  ```

### Task 7: Test on iPhone (iOS Safari) (AC-7)
- [ ] **iPhone 13 Pro Testing (390px width):**
  - Portrait mode:
    - Verify hero section readable (28px title, 16px subtitle)
    - Verify format badges stacked (1 per row)
    - Verify file cards stacked (1 per row)
    - Verify tap targets minimum 44px (test with finger, not stylus)
    - Verify "Convert" button 48px height (easy to tap)
    - Test file upload via tap (file picker opens)
    - Test conversion flow (tap "Convert", verify status updates)
  - Landscape mode:
    - Verify no horizontal scrolling
    - Verify buttons still touch-friendly
- [ ] **iPhone SE Testing (320px width - smallest device):**
  - Portrait mode:
    - Verify all content fits without horizontal scrolling
    - Verify text readable (16px minimum body text)
    - Verify long filenames truncate (no overflow)
    - Verify buttons full-width and 44px+ height
    - Test complete conversion flow (upload, convert, download)
  - Document any layout issues at 320px width

### Task 8: Test on Android (Chrome) (AC-7)
- [ ] **Samsung Galaxy S21 Testing (360px width):**
  - Portrait mode:
    - Verify hero section layout (28px title readable)
    - Verify format badges stacked
    - Verify file cards single column
    - Test tap targets (44px minimum)
    - Test file upload (Chrome file picker)
    - Test conversion flow (tap "Convert", download)
  - Landscape mode:
    - Verify responsive layout adapts
    - Verify no horizontal scrolling
- [ ] **Google Pixel 6 Testing (393px width):**
  - Portrait mode:
    - Verify layout similar to iPhone 13 Pro
    - Test tap targets and file upload
    - Test conversion flow

### Task 9: Test on iPad (Safari) (AC-7)
- [ ] **iPad Pro 11" Testing (834px width):**
  - Portrait mode:
    - Verify 2-column file card grid
    - Verify 2×3 format badge grid
    - Verify hero section 36px title
    - Test tap targets (44px minimum)
    - Test conversion flow
  - Landscape mode (1194px width):
    - Verify layout switches to 3-column grid (desktop breakpoint)
    - Verify format badges horizontal row
- [ ] **iPad Air Testing (820px width):**
  - Portrait mode:
    - Verify 2-column layout (tablet breakpoint)
    - Test tap targets and conversion flow

### Task 10: Test on Desktop Browsers (AC-7)
- [ ] **Chrome Desktop (1920px width):**
  - Verify 3-column file card grid
  - Verify format badges horizontal row (5 badges)
  - Verify hero section 48px title
  - Test responsive design tools:
    - Resize window 1920px → 768px → 320px
    - Verify breakpoints trigger correctly
    - Verify no layout jumps or janky transitions
- [ ] **Firefox Desktop (1920px width):**
  - Verify layout identical to Chrome
  - Test responsive design mode (built-in DevTools)
  - Test at 320px, 768px, 1024px, 1920px widths
- [ ] **Safari Desktop (1440px width - macOS):**
  - Verify layout matches Chrome/Firefox
  - Test at common macOS resolutions (1440px, 1680px, 1920px)
- [ ] **Edge Desktop (1920px width - Windows):**
  - Verify layout matches Chrome (Chromium-based)
  - Test Windows-specific rendering issues

### Task 11: Manual Testing Documentation
- [ ] Create testing checklist document: `docs/testing/responsive-design-testing.md`
- [ ] Document test results for each device:
  - Device name, screen size, OS version
  - Pass/Fail for each acceptance criteria
  - Screenshots of any layout issues
  - Notes on touch target size, readability
- [ ] Document known issues:
  - iOS Safari text rendering quirks
  - Android Chrome file picker differences
  - Tablet landscape mode edge cases
- [ ] Create GitHub issue for any failures found during testing

## Dev Notes

### Learnings from Previous Story

**From Story 10-4-individual-file-actions (Status: drafted)**

Previous story not yet implemented. Story 10.5 makes all functionality from Stories 10-2, 10-3, and 10-4 responsive across devices.

**Reuse from Previous Stories:**
- File card structure (Stories 10-2, 10-3, 10-4)
- Upload drop zone (Story 10-2)
- Format badges (Story 10-1)
- Conversion controls (Stories 10-3, 10-4)

**Responsive Adaptations Needed:**
- Grid layouts: 3-column → 2-column → 1-column
- Font sizes: 48px → 36px → 28px (hero title)
- Button widths: Auto-width → Full-width (mobile)
- Tap targets: Desktop click targets → 44px+ touch targets

**No New Functionality:**
- Story 10.5 is purely responsive design (no new features)
- All components exist from previous stories, just need responsive CSS

[Source: docs/stories/10-4-individual-file-actions.md]

### Architecture Alignment

**Tech Spec Epic 10 Alignment:**

Story 10.5 implements **AC-5 (Mobile-Responsive Design)** from tech-spec-epic-10.md.

**Responsive Breakpoints:**
```css
/* Mobile-first approach (default: <768px) */
.upload__files {
  grid-template-columns: 1fr; /* 1 column */
}

/* Tablet: 768px-1023px */
@media (min-width: 768px) {
  .upload__files {
    grid-template-columns: repeat(2, 1fr); /* 2 columns */
  }
}

/* Desktop: 1024px+ */
@media (min-width: 1024px) {
  .upload__files {
    grid-template-columns: repeat(3, 1fr); /* 3 columns */
  }
}
```

**Breakpoint Rationale:**
- **320px**: Smallest common device (iPhone SE)
- **768px**: Tablet portrait (iPad, Android tablets)
- **1024px**: Desktop / Tablet landscape (iPad Pro landscape)
- **1920px**: Desktop full HD (tested maximum)

[Source: docs/tech-spec-epic-10.md#Responsive-Design-Strategy]

### Mobile-First CSS Approach

**Why Mobile-First?**

Recipe uses mobile-first CSS (default styles for mobile, override for larger screens):

```css
/* Mobile default (applies to all screen sizes) */
.hero__title {
  font-size: 28px;
}

/* Tablet override (applies to 768px+) */
@media (min-width: 768px) {
  .hero__title {
    font-size: 36px;
  }
}

/* Desktop override (applies to 1024px+) */
@media (min-width: 1024px) {
  .hero__title {
    font-size: 48px;
  }
}
```

**Benefits:**
- Smaller CSS file size (mobile styles loaded first)
- Faster mobile performance (no overriding desktop styles)
- Progressive enhancement (mobile works, desktop enhances)
- Easier to maintain (add features for desktop, don't remove for mobile)

**Alternative (Desktop-First):**
Desktop-first uses `max-width` media queries, which is harder to maintain and results in larger CSS files for mobile.

[Source: CSS Best Practices - Mobile-First Responsive Design]

### Touch Target Guidelines

**iOS Human Interface Guidelines:**

Apple recommends minimum 44px × 44px tap targets:
- 44px: Minimum for all interactive elements
- 48px: Recommended for primary actions (better usability)
- 56px: Extra large for critical actions (e.g., "Convert All")

**Android Material Design Guidelines:**

Google recommends minimum 48dp tap targets:
- 48dp ≈ 48px (dp = density-independent pixels)
- 8dp spacing between tap targets

**Recipe Implementation:**
```css
.button {
  min-height: 44px; /* iOS minimum */
}

.file-card__convert {
  min-height: 48px; /* Recommended for primary action */
}

#convert-all-btn {
  min-height: 56px; /* Extra large for batch action */
}

.file-card__footer {
  gap: 8px; /* Minimum spacing between targets */
}
```

[Source: iOS Human Interface Guidelines - Touch Targets]
[Source: Material Design - Touch Targets]

### Viewport Meta Tag

**Critical Meta Tag:**

```html
<meta name="viewport" content="width=device-width, initial-scale=1">
```

**What it does:**
- `width=device-width`: Sets viewport width to device width (not desktop width)
- `initial-scale=1`: No zoom on page load (1:1 scale)

**Without this tag:**
- Mobile browsers render page at desktop width (980px)
- User sees tiny desktop layout, must zoom to read
- Responsive CSS doesn't work (viewport is 980px, not 320px)

**Common Mistakes:**
- `initial-scale=0.5` (page loads zoomed out, bad UX)
- `user-scalable=no` (prevents zooming, accessibility issue)
- Missing tag entirely (page renders as desktop on mobile)

[Source: MDN Web Docs - Viewport Meta Tag]

### Preventing Horizontal Scrolling

**Common Causes:**

1. **Fixed-width containers:**
   ```css
   /* Bad: Fixed width exceeds mobile viewport */
   .container {
     width: 1200px;
   }

   /* Good: Max-width with percentage fallback */
   .container {
     width: 100%;
     max-width: 1200px;
   }
   ```

2. **Oversized images:**
   ```css
   /* Bad: Image exceeds container */
   img {
     width: auto;
   }

   /* Good: Image scales to container */
   img {
     max-width: 100%;
     height: auto;
   }
   ```

3. **Negative margins:**
   ```css
   /* Bad: Negative margin extends beyond viewport */
   .element {
     margin-left: -50px;
   }

   /* Good: Use padding or transform */
   .element {
     transform: translateX(-50px);
   }
   ```

**Recipe Prevention:**
```css
body {
  overflow-x: hidden; /* Prevent horizontal scroll */
}

img, svg {
  max-width: 100%; /* Scale images */
}

* {
  box-sizing: border-box; /* Include padding in width calculations */
}
```

[Source: CSS Best Practices - Preventing Horizontal Scrolling]

### WCAG Accessibility Compliance

**WCAG AA Requirements:**

1. **Text Size:**
   - Body text: 16px minimum (1em)
   - Small text: 14px minimum (0.875em)

2. **Line Height:**
   - Minimum: 1.5 (150% of font size)
   - Optimal: 1.6-1.8 for body text

3. **Color Contrast:**
   - Body text (16px): 4.5:1 contrast ratio
   - Large text (24px+): 3:1 contrast ratio

4. **Paragraph Width:**
   - Maximum: 80 characters per line
   - Optimal: 45-75 characters (65ch)

**Recipe Implementation:**
```css
html {
  font-size: 16px; /* Minimum */
  line-height: 1.5; /* WCAG AA */
}

p {
  max-width: 65ch; /* Optimal readability */
}

.text-small {
  font-size: 14px; /* Minimum for small text */
}

/* Color contrast checked with WCAG tools */
:root {
  --color-text: #222; /* 4.5:1 on white background */
  --color-text-secondary: #666; /* 4.5:1 on white */
}
```

[Source: WCAG 2.1 AA Guidelines]

### Project Structure Notes

**Modified Files (Story 10.5):**
- `web/index.html` - Add viewport meta tag
- `web/css/main.css` - Mobile-first CSS reset, base styles
- `web/css/layout.css` - Responsive hero, container, grid layouts
- `web/css/components.css` - Responsive badges, buttons, file cards, drop zone

**No New Files Created:**
- Story 10.5 adds responsive CSS to existing files

**Files from Previous Stories (Made Responsive):**
- `web/index.html` - Hero, format badges, upload section (Story 10-1, 10-2)
- `web/css/components.css` - File cards, badges, buttons (Stories 10-2, 10-3, 10-4)

[Source: docs/tech-spec-epic-10.md#Services-and-Modules]

### Testing Strategy

**Manual Testing (Required):**

Recipe requires manual testing on real devices (not just browser DevTools):

**Why Real Device Testing?**
- Touch interactions differ from mouse clicks (tap vs click)
- iOS Safari has unique rendering quirks (text size, zooming)
- Android Chrome has different file picker UI
- DevTools responsive mode doesn't simulate touch accurately

**Testing Matrix:**
```
4 Device Categories × 2 Orientations = 8 Test Scenarios

Devices:
- iPhone (iOS Safari) - 320px-390px width
- Android (Chrome) - 360px-393px width
- iPad (Safari) - 820px-834px width
- Desktop (Chrome/Firefox/Safari/Edge) - 1024px-1920px width

Orientations:
- Portrait (primary)
- Landscape (secondary)
```

**Test Checklist Per Device:**
1. Visual layout (grid columns, font sizes, spacing)
2. Tap targets (44px minimum, easy to tap)
3. Text readability (16px minimum, no zooming needed)
4. No horizontal scrolling (all content fits viewport)
5. File upload (tap to browse, file picker opens)
6. Conversion flow (convert, download, remove)
7. Error handling (error messages readable)

**Documentation:**
All test results documented in `docs/testing/responsive-design-testing.md` with:
- Device name, screen size, OS version
- Pass/Fail for each AC
- Screenshots of layout issues
- Notes on touch usability

[Source: docs/tech-spec-epic-10.md#Test-Strategy-Summary]

### Known Risks

**RISK-38: iOS Safari text size adjustment**
- **Impact**: iOS Safari auto-adjusts text size in landscape mode (unreadable)
- **Mitigation**: Add `-webkit-text-size-adjust: 100%;` to prevent auto-sizing
- **Test**: Rotate iPhone to landscape, verify text doesn't shrink

**RISK-39: Android Chrome file picker differences**
- **Impact**: Android file picker UI differs from iOS (user confusion)
- **Mitigation**: Accept differences (native OS file picker, can't customize)
- **Documentation**: Note in testing docs that file picker varies by OS

**RISK-40: iPad landscape mode breakpoint edge case**
- **Impact**: iPad landscape (1194px) triggers desktop layout (1024px breakpoint), but touch targets still needed
- **Mitigation**: Ensure touch targets maintained at all breakpoints (44px minimum applies to all devices)
- **Test**: iPad landscape mode, verify tap targets still 44px+

**RISK-41: Long filenames overflow on mobile**
- **Impact**: 100-character filenames exceed mobile viewport width
- **Mitigation**: Truncate filenames at 30 characters (existing implementation from Story 10-2)
- **Test**: Upload file with 100-char name, verify truncation works

[Source: docs/tech-spec-epic-10.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-10.md#Acceptance-Criteria] - AC-5: Mobile-Responsive Design
- [Source: docs/tech-spec-epic-10.md#Responsive-Design-Strategy] - Breakpoints, mobile-first approach
- [Source: docs/stories/10-1-landing-page-redesign.md] - Hero section, format badges
- [Source: docs/stories/10-2-batch-file-upload.md] - Upload drop zone, file cards
- [Source: docs/stories/10-3-progress-indicators.md] - Status indicators, download buttons
- [Source: docs/stories/10-4-individual-file-actions.md] - Conversion controls, format dropdowns
- [iOS Human Interface Guidelines - Touch Targets](https://developer.apple.com/design/human-interface-guidelines/ios/visual-design/adaptivity-and-layout/)
- [Material Design - Touch Targets](https://material.io/design/usability/accessibility.html#layout-and-typography)
- [MDN: Viewport Meta Tag](https://developer.mozilla.org/en-US/docs/Web/HTML/Viewport_meta_tag)
- [WCAG 2.1 AA Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [CSS Best Practices - Mobile-First Design](https://developer.mozilla.org/en-US/docs/Web/CSS/Media_Queries/Using_media_queries)

## Dev Agent Record

### Context Reference

- `docs/stories/10-5-responsive-mobile-design.context.xml` (Generated: 2025-11-09)

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
