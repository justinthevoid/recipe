# Story 11.2: Reference Image Bundle (Portrait, Landscape, Product)

Status: ready-for-dev

## Story

As a **photographer evaluating preset conversions**,
I want **three representative reference images (Portrait, Landscape, Product) to preview presets on**,
so that **I can see how a preset affects different photography genres before converting my own files**.

## Acceptance Criteria

**AC-1: Three Reference Images Provided**
- ✅ Portrait reference image:
  - Subject: Person (headshot or upper body)
  - Lighting: Natural, soft (not harsh shadows)
  - Colors: Skin tones, neutral background
  - Use case: Test preset impact on portraits, skin tone adjustments
- ✅ Landscape reference image:
  - Subject: Outdoor scene (mountains, sky, trees, water)
  - Lighting: Daylight (blue sky, clouds, greenery)
  - Colors: Blues (sky), greens (foliage), warm tones (earth)
  - Use case: Test preset impact on nature, landscape photography
- ✅ Product/Still-life reference image:
  - Subject: Inanimate object (product, food, still-life setup)
  - Lighting: Studio or controlled (even, neutral)
  - Colors: Neutral background, varied product colors
  - Use case: Test preset impact on commercial, product photography
- ✅ Images representative of common photography genres:
  - Cover 80%+ of typical preset use cases
  - Neutral starting point (not extreme colors or lighting)
  - Work well with typical preset adjustments (exposure, contrast, saturation)

**AC-2: Image Optimization for Web**
- ✅ File format: WebP (modern, efficient compression)
  - Fallback: JPEG for older browsers (via `<picture>` element)
- ✅ Image size: <200 KB each (gzipped)
  - Portrait: Target 150 KB
  - Landscape: Target 180 KB (more detail)
  - Product: Target 120 KB (simpler composition)
  - Total bundle: <600 KB (all three images)
- ✅ Image dimensions:
  - Width: 1200px (desktop preview)
  - Height: Variable (maintain aspect ratio)
  - Aspect ratios:
    - Portrait: 3:4 or 2:3 (vertical orientation)
    - Landscape: 16:9 or 3:2 (horizontal orientation)
    - Product: 4:3 or 1:1 (square or near-square)
- ✅ Compression:
  - WebP quality: 80-85% (balance size vs. quality)
  - JPEG quality: 85-90% (fallback)
  - Lossless optimization (remove EXIF metadata, color profiles)

**AC-3: Licensing and Attribution**
- ✅ Licensing: Public domain (CC0) or created specifically for Recipe
  - No copyright restrictions
  - No attribution required (but appreciated)
  - Commercial use allowed
- ✅ Sources:
  - Option 1: Use public domain images from Unsplash, Pexels, Pixabay (verify CC0 license)
  - Option 2: Create images specifically for Recipe (Justin's photos)
  - Option 3: Commission photographer for CC0 images
- ✅ Attribution (optional but recommended):
  - Include `CREDITS.md` file in `web/images/` directory
  - List image sources, photographers, licenses
  - Example: "Portrait by John Doe (CC0), Landscape by Jane Smith (CC0)"
- ✅ No licensing conflicts:
  - Verify no restrictive licenses (e.g., Creative Commons BY-SA, ND)
  - Ensure images can be distributed with Recipe under MIT license

**AC-4: Images Embedded in Web Bundle**
- ✅ Images stored in `web/images/` directory:
  - `web/images/preview-portrait.webp`
  - `web/images/preview-landscape.webp`
  - `web/images/preview-product.webp`
  - Fallback JPEGs:
    - `web/images/preview-portrait.jpg`
    - `web/images/preview-landscape.jpg`
    - `web/images/preview-product.jpg`
- ✅ No external requests (images self-hosted):
  - No CDN dependencies (Cloudflare, AWS S3, etc.)
  - All images served from Recipe domain
  - Faster load times (no DNS lookup, no external latency)
- ✅ Images preloaded on page load (optional, for performance):
  ```html
  <link rel="preload" as="image" href="/images/preview-portrait.webp" type="image/webp">
  <link rel="preload" as="image" href="/images/preview-landscape.webp" type="image/webp">
  <link rel="preload" as="image" href="/images/preview-product.webp" type="image/webp">
  ```

**AC-5: Images Work Well with Typical Presets**
- ✅ Neutral starting point:
  - Images not over-saturated (avoid already vibrant colors)
  - Images not underexposed or overexposed (neutral exposure)
  - Images have good tonal range (shadows, midtones, highlights)
- ✅ Test with typical preset adjustments:
  - Exposure ±1.0: Image should not clip highlights or crush blacks
  - Contrast ±0.5: Image should show improved separation (not muddy)
  - Saturation ±0.5: Image should look more/less colorful (not oversaturated)
  - Hue ±30°: Image should show subtle color shift (not unnatural)
  - Temperature ±20: Image should look warmer/cooler (not extreme)
- ✅ Images representative of real-world usage:
  - Portrait: Skin tones render naturally with typical adjustments
  - Landscape: Sky, foliage, water render well with typical adjustments
  - Product: Neutral background, varied colors respond predictably

**AC-6: Responsive Image Handling**
- ✅ Responsive image sizes for mobile:
  - Desktop (1200px width): Full-size WebP
  - Tablet (800px width): Medium-size WebP (800px)
  - Mobile (400px width): Small-size WebP (400px)
- ✅ Use `<picture>` element for responsive images:
  ```html
  <picture>
    <source srcset="/images/preview-portrait-400w.webp" media="(max-width: 600px)" type="image/webp">
    <source srcset="/images/preview-portrait-800w.webp" media="(max-width: 1024px)" type="image/webp">
    <source srcset="/images/preview-portrait.webp" type="image/webp">
    <img src="/images/preview-portrait.jpg" alt="Portrait reference image" loading="lazy">
  </picture>
  ```
- ✅ Lazy loading for images (defer until needed):
  - `loading="lazy"` attribute on `<img>` elements
  - Images load only when preview modal opened (not on page load)
- ✅ Total image bundle:
  - Desktop: 600 KB (3 full-size images)
  - Tablet: 400 KB (3 medium-size images)
  - Mobile: 200 KB (3 small-size images)

**AC-7: Image Selection and Validation**
- ✅ Image selection criteria documented:
  - Portrait: Representative skin tones (light, medium, dark)
  - Landscape: Typical outdoor colors (blue sky, green foliage)
  - Product: Neutral background (white, gray, black)
- ✅ Visual validation:
  - Images look good at full size (no pixelation, no compression artifacts)
  - Images look good with typical CSS filters applied
  - Images work well in before/after slider (clear comparison)
- ✅ Technical validation:
  - Image dimensions correct (1200px width)
  - Image file sizes within budget (<200 KB each)
  - Image format correct (WebP primary, JPEG fallback)
  - Images load successfully in all browsers

## Tasks / Subtasks

### Task 1: Select Reference Images (AC-1, AC-3, AC-5)
- [ ] Portrait image selection:
  - Search public domain sources: Unsplash, Pexels, Pixabay
  - Keywords: "portrait", "headshot", "person", "natural light"
  - Criteria:
    - Skin tones visible (light, medium, or dark)
    - Neutral background (not distracting)
    - Soft, natural lighting (not harsh shadows)
    - Neutral expression (friendly, approachable)
    - High resolution (2000px+ width, will be downscaled)
  - License: Verify CC0 (public domain)
  - Download high-resolution version
- [ ] Landscape image selection:
  - Search public domain sources: Unsplash, Pexels, Pixabay
  - Keywords: "landscape", "mountains", "nature", "outdoor", "sky"
  - Criteria:
    - Blue sky visible (test saturation, hue adjustments)
    - Green foliage visible (test color shifts)
    - Good tonal range (shadows in foreground, highlights in sky)
    - Neutral exposure (not overexposed, not underexposed)
    - High resolution (2000px+ width, will be downscaled)
  - License: Verify CC0 (public domain)
  - Download high-resolution version
- [ ] Product/Still-life image selection:
  - Search public domain sources: Unsplash, Pexels, Pixabay
  - Keywords: "product", "still life", "food", "studio", "white background"
  - Criteria:
    - Neutral background (white, gray, or black)
    - Varied product colors (test saturation, hue)
    - Even lighting (no harsh shadows)
    - Simple composition (not cluttered)
    - High resolution (2000px+ width, will be downscaled)
  - License: Verify CC0 (public domain)
  - Download high-resolution version

### Task 2: Optimize Images for Web (AC-2)
- [ ] Install image optimization tools:
  - WebP converter: `cwebp` (Google WebP tools)
    - macOS: `brew install webp`
    - Windows: Download from https://developers.google.com/speed/webp/download
    - Linux: `sudo apt-get install webp`
  - ImageMagick (for resizing, JPEG conversion):
    - macOS: `brew install imagemagick`
    - Windows: Download from https://imagemagick.org/
    - Linux: `sudo apt-get install imagemagick`
- [ ] Resize images to target dimensions:
  ```bash
  # Desktop size (1200px width)
  convert portrait-original.jpg -resize 1200x portrait-1200w.jpg
  convert landscape-original.jpg -resize 1200x landscape-1200w.jpg
  convert product-original.jpg -resize 1200x product-1200w.jpg

  # Tablet size (800px width)
  convert portrait-original.jpg -resize 800x portrait-800w.jpg
  convert landscape-original.jpg -resize 800x landscape-800w.jpg
  convert product-original.jpg -resize 800x product-800w.jpg

  # Mobile size (400px width)
  convert portrait-original.jpg -resize 400x portrait-400w.jpg
  convert landscape-original.jpg -resize 400x landscape-400w.jpg
  convert product-original.jpg -resize 400x product-400w.jpg
  ```
- [ ] Convert to WebP format:
  ```bash
  # Desktop WebP (quality 80)
  cwebp -q 80 portrait-1200w.jpg -o preview-portrait.webp
  cwebp -q 80 landscape-1200w.jpg -o preview-landscape.webp
  cwebp -q 80 product-1200w.jpg -o preview-product.webp

  # Tablet WebP (quality 80)
  cwebp -q 80 portrait-800w.jpg -o preview-portrait-800w.webp
  cwebp -q 80 landscape-800w.jpg -o preview-landscape-800w.webp
  cwebp -q 80 product-800w.jpg -o preview-product-800w.webp

  # Mobile WebP (quality 75, smaller file)
  cwebp -q 75 portrait-400w.jpg -o preview-portrait-400w.webp
  cwebp -q 75 landscape-400w.jpg -o preview-landscape-400w.webp
  cwebp -q 75 product-400w.jpg -o preview-product-400w.webp
  ```
- [ ] Optimize JPEG fallbacks:
  ```bash
  # Desktop JPEG (quality 85)
  convert portrait-1200w.jpg -quality 85 -strip preview-portrait.jpg
  convert landscape-1200w.jpg -quality 85 -strip preview-landscape.jpg
  convert product-1200w.jpg -quality 85 -strip preview-product.jpg
  ```
- [ ] Verify file sizes:
  ```bash
  ls -lh web/images/preview-*.webp
  ls -lh web/images/preview-*.jpg

  # Target:
  # preview-portrait.webp: <150 KB
  # preview-landscape.webp: <180 KB
  # preview-product.webp: <120 KB
  # Total WebP: <600 KB
  ```
- [ ] Adjust quality if needed:
  - If file size exceeds target: Reduce quality (-q 70 or -q 75)
  - If visual quality poor: Increase quality (-q 85 or -q 90)
  - Balance: Target 80-85% quality for good size/quality trade-off

### Task 3: Create Image Directory Structure (AC-4)
- [ ] Create `web/images/` directory:
  ```bash
  mkdir -p web/images
  ```
- [ ] Move optimized images to directory:
  ```bash
  mv preview-portrait*.webp web/images/
  mv preview-landscape*.webp web/images/
  mv preview-product*.webp web/images/
  mv preview-portrait.jpg web/images/
  mv preview-landscape.jpg web/images/
  mv preview-product.jpg web/images/
  ```
- [ ] Verify directory structure:
  ```
  web/
  └── images/
      ├── preview-portrait.webp         (Desktop, ~150 KB)
      ├── preview-portrait-800w.webp    (Tablet, ~80 KB)
      ├── preview-portrait-400w.webp    (Mobile, ~30 KB)
      ├── preview-portrait.jpg          (Fallback, ~180 KB)
      ├── preview-landscape.webp        (Desktop, ~180 KB)
      ├── preview-landscape-800w.webp   (Tablet, ~100 KB)
      ├── preview-landscape-400w.webp   (Mobile, ~40 KB)
      ├── preview-landscape.jpg         (Fallback, ~220 KB)
      ├── preview-product.webp          (Desktop, ~120 KB)
      ├── preview-product-800w.webp     (Tablet, ~60 KB)
      ├── preview-product-400w.webp     (Mobile, ~25 KB)
      └── preview-product.jpg           (Fallback, ~140 KB)
  ```

### Task 4: Create Attribution File (AC-3)
- [ ] Create `web/images/CREDITS.md` file:
  ```markdown
  # Reference Image Credits

  Recipe uses the following reference images for preset preview:

  ## Portrait

  - **Image**: Preview portrait reference
  - **Photographer**: [Name or "Unknown"]
  - **Source**: [Unsplash/Pexels/Pixabay URL or "Created for Recipe"]
  - **License**: CC0 (Public Domain)
  - **Original URL**: [URL or "N/A"]

  ## Landscape

  - **Image**: Preview landscape reference
  - **Photographer**: [Name or "Unknown"]
  - **Source**: [Unsplash/Pexels/Pixabay URL or "Created for Recipe"]
  - **License**: CC0 (Public Domain)
  - **Original URL**: [URL or "N/A"]

  ## Product

  - **Image**: Preview product reference
  - **Photographer**: [Name or "Unknown"]
  - **Source**: [Unsplash/Pexels/Pixabay URL or "Created for Recipe"]
  - **License**: CC0 (Public Domain)
  - **Original URL**: [URL or "N/A"]

  ---

  All images are licensed under CC0 (Public Domain) and can be used without attribution.
  Recipe is licensed under the MIT License.
  ```
- [ ] Document image sources for each reference image
- [ ] Verify licenses are correctly attributed (CC0 only)

### Task 5: Implement Responsive Image HTML (AC-6)
- [ ] Add `<picture>` elements to preview modal HTML:
  ```html
  <!-- web/index.html (or preview modal template) -->
  <div id="preview-image-container">
    <picture id="preview-image-portrait">
      <source srcset="/images/preview-portrait-400w.webp" media="(max-width: 600px)" type="image/webp">
      <source srcset="/images/preview-portrait-800w.webp" media="(max-width: 1024px)" type="image/webp">
      <source srcset="/images/preview-portrait.webp" type="image/webp">
      <img src="/images/preview-portrait.jpg"
           alt="Portrait reference image for preset preview"
           loading="lazy"
           class="preview-image">
    </picture>

    <picture id="preview-image-landscape" hidden>
      <source srcset="/images/preview-landscape-400w.webp" media="(max-width: 600px)" type="image/webp">
      <source srcset="/images/preview-landscape-800w.webp" media="(max-width: 1024px)" type="image/webp">
      <source srcset="/images/preview-landscape.webp" type="image/webp">
      <img src="/images/preview-landscape.jpg"
           alt="Landscape reference image for preset preview"
           loading="lazy"
           class="preview-image">
    </picture>

    <picture id="preview-image-product" hidden>
      <source srcset="/images/preview-product-400w.webp" media="(max-width: 600px)" type="image/webp">
      <source srcset="/images/preview-product-800w.webp" media="(max-width: 1024px)" type="image/webp">
      <source srcset="/images/preview-product.webp" type="image/webp">
      <img src="/images/preview-product.jpg"
           alt="Product reference image for preset preview"
           loading="lazy"
           class="preview-image">
    </picture>
  </div>
  ```
- [ ] Add JavaScript to switch between reference images (Story 11.3):
  ```javascript
  // web/js/preview.js
  function showReferenceImage(imageType) {
    // Hide all images
    document.getElementById('preview-image-portrait').hidden = true;
    document.getElementById('preview-image-landscape').hidden = true;
    document.getElementById('preview-image-product').hidden = true;

    // Show selected image
    const imageMap = {
      'portrait': 'preview-image-portrait',
      'landscape': 'preview-image-landscape',
      'product': 'preview-image-product'
    };

    const imageId = imageMap[imageType] || 'preview-image-portrait';
    document.getElementById(imageId).hidden = false;
  }

  // Default to portrait on load
  showReferenceImage('portrait');
  ```

### Task 6: Optional Preload Optimization (AC-4)
- [ ] Add `<link rel="preload">` for critical images (optional):
  ```html
  <!-- web/index.html <head> -->
  <!-- Preload portrait image (default preview) -->
  <link rel="preload" as="image" href="/images/preview-portrait.webp" type="image/webp">

  <!-- Note: Only preload portrait (default), lazy-load others -->
  <!-- Landscape and Product loaded on-demand when tabs clicked -->
  ```
- [ ] Test impact on page load time:
  - Without preload: Measure page load time (Lighthouse)
  - With preload: Measure page load time (Lighthouse)
  - Decision: Keep preload if improvement >200ms, remove if negligible

### Task 7: Visual Quality Validation (AC-7)
- [ ] Test images at full size:
  - Open each image in browser at 100% zoom
  - Check for pixelation, compression artifacts
  - Verify colors look natural (not oversaturated, not washed out)
  - Verify sharpness (not blurry, not over-sharpened)
- [ ] Test images with CSS filters applied (Story 11.1):
  - Apply `brightness(150%)` → Image should not blow out highlights
  - Apply `contrast(130%)` → Image should show improved separation
  - Apply `saturate(150%)` → Image should look more colorful (not garish)
  - Apply `hue-rotate(30deg)` → Image should show subtle color shift
  - Apply combination filters → Image should remain usable
- [ ] Test images in before/after slider (Story 11.4):
  - Open slider at 0% (before): Original image visible
  - Open slider at 50%: Half original, half filtered
  - Open slider at 100% (after): Fully filtered image visible
  - Verify slider interaction smooth (no janky rendering)

### Task 8: Browser Compatibility Testing (AC-7)
- [ ] Test WebP support on all browsers:
  - Chrome (latest): ✅ WebP supported
  - Firefox (latest): ✅ WebP supported
  - Safari (latest): ✅ WebP supported (since Safari 14, 2020)
  - Edge (latest): ✅ WebP supported
  - Safari iOS (latest): ✅ WebP supported (since iOS 14)
- [ ] Test JPEG fallback on older browsers:
  - Safari 13 (macOS 10.15 Catalina): Should load JPEG fallback
  - IE 11 (if testing): Should load JPEG fallback
- [ ] Test responsive images on mobile:
  - iPhone (Safari): Verify 400w WebP loads on mobile
  - Android (Chrome): Verify 400w WebP loads on mobile
  - Tablet (iPad): Verify 800w WebP loads on tablet

### Task 9: Performance Testing (AC-2, AC-6)
- [ ] Measure image load times:
  - Desktop (1200px): Measure time to load all 3 full-size images
  - Tablet (800px): Measure time to load all 3 medium-size images
  - Mobile (400px): Measure time to load all 3 small-size images
  - Target: <1 second for all images on 3G connection
- [ ] Lighthouse audit:
  - Run Lighthouse on https://recipe.justins.studio
  - Check "Properly size images" audit (should pass)
  - Check "Serve images in next-gen formats" audit (should pass, WebP)
  - Check "Defer offscreen images" audit (should pass, lazy loading)
- [ ] Network waterfall analysis (DevTools):
  - Open DevTools → Network tab
  - Load page → Open preview modal
  - Verify images load only when modal opened (lazy loading)
  - Verify correct image size loaded (responsive srcset)

### Task 10: Documentation (AC-1, AC-3)
- [ ] Add reference images section to README:
  ```markdown
  ## Reference Images (Epic 11)

  Recipe includes three reference images for preset preview:

  ### Portrait
  - **Use case**: Test presets on portraits, skin tones
  - **Lighting**: Natural, soft
  - **Colors**: Skin tones, neutral background

  ### Landscape
  - **Use case**: Test presets on nature, landscapes
  - **Lighting**: Daylight (blue sky, greenery)
  - **Colors**: Blues, greens, warm tones

  ### Product
  - **Use case**: Test presets on commercial, product photography
  - **Lighting**: Studio, even
  - **Colors**: Neutral background, varied products

  ### Licensing

  All reference images are licensed under CC0 (Public Domain).
  See `web/images/CREDITS.md` for attribution details.
  ```
- [ ] Add image optimization documentation (for contributors):
  ```markdown
  ### Adding Reference Images

  To add or replace reference images:

  1. Select high-resolution images (2000px+ width)
  2. Verify CC0 license (public domain)
  3. Resize to target dimensions (1200px, 800px, 400px)
  4. Convert to WebP (quality 80):
     ```bash
     cwebp -q 80 image-1200w.jpg -o preview-name.webp
     ```
  5. Create JPEG fallback (quality 85):
     ```bash
     convert image-1200w.jpg -quality 85 -strip preview-name.jpg
     ```
  6. Update `CREDITS.md` with attribution
  7. Test file sizes (<200 KB each)
  ```

## Dev Notes

### Learnings from Previous Story

**From Story 11-1-css-filter-mapping (Status: drafted)**

Previous story not yet implemented. Story 11.2 provides reference images that Story 11.1 will apply CSS filters to.

**Integration with Story 11.1:**
- CSS filters (Story 11.1) applied to reference images (Story 11.2)
- Images must work well with typical filter adjustments (exposure, contrast, saturation, hue)
- Images should have neutral starting point (not extreme colors or exposure)
- Images must render smoothly with GPU-accelerated filters (<100ms)

**Technical Requirements:**
- Images must support `filter` CSS property (all modern browsers)
- Images must be high enough quality to show filter effects clearly
- Images must have good tonal range (shadows, midtones, highlights)

[Source: docs/stories/11-1-css-filter-mapping.md]

### Architecture Alignment

**Tech Spec Epic 11 Alignment:**

Story 11.2 implements **AC-2: Reference Image Bundle** from tech-spec-epic-11.md.

**Reference Image Requirements:**

```
Image Type   Dimensions  File Size   Aspect Ratio   Lighting        Use Case
-----------  ----------  ----------  -------------  --------------  ---------------------------
Portrait     1200x1600   <150 KB     3:4            Natural, soft   Portraits, skin tones
Landscape    1200x800    <180 KB     16:9           Daylight        Nature, landscapes, outdoor
Product      1200x1200   <120 KB     1:1            Studio, even    Commercial, products, food
-----------  ----------  ----------  -------------  --------------  ---------------------------
Total                    <600 KB
```

**Responsive Image Sizes:**

```
Breakpoint   Width   Portrait   Landscape   Product   Total
-----------  ------  ---------  ----------  --------  ------
Desktop      1200px  150 KB     180 KB      120 KB    450 KB
Tablet       800px   80 KB      100 KB      60 KB     240 KB
Mobile       400px   30 KB      40 KB       25 KB     95 KB
```

[Source: docs/tech-spec-epic-11.md#AC-2]

### WebP Format Benefits

**Why WebP Over JPEG?**

| Feature                | WebP           | JPEG   | Improvement     |
| ---------------------- | -------------- | ------ | --------------- |
| File size (quality 80) | 150 KB         | 220 KB | 32% smaller     |
| Compression            | Lossy/Lossless | Lossy  | More flexible   |
| Transparency           | ✅ Yes          | ❌ No   | Better support  |
| Animation              | ✅ Yes          | ❌ No   | Bonus feature   |
| Browser support        | 95%+           | 100%   | Fallback needed |

**WebP Savings:**
- Portrait: 220 KB (JPEG) → 150 KB (WebP) = **32% reduction**
- Landscape: 280 KB (JPEG) → 180 KB (WebP) = **36% reduction**
- Product: 180 KB (JPEG) → 120 KB (WebP) = **33% reduction**
- Total: 680 KB (JPEG) → 450 KB (WebP) = **34% reduction**

**Performance Impact:**
- 3G connection (1.6 Mbps = 200 KB/s):
  - JPEG bundle (680 KB): 3.4 seconds
  - WebP bundle (450 KB): 2.25 seconds
  - **Savings: 1.15 seconds** (34% faster)

[Source: WebP Compression Study - Google Developers]

### Image Licensing Best Practices

**CC0 (Public Domain) Requirements:**

✅ **Allowed:**
- Use images without attribution
- Modify images (resize, crop, filter)
- Distribute images with Recipe
- Commercial use

❌ **Not allowed (other licenses):**
- CC BY: Requires attribution
- CC BY-SA: Requires share-alike (same license)
- CC BY-ND: No derivatives (can't modify)
- CC BY-NC: Non-commercial only (Recipe is commercial-use)

**Verification Checklist:**
1. Check license on source website (Unsplash, Pexels, Pixabay)
2. Verify "CC0" or "Public Domain" label
3. Read license terms (some sites have restrictions)
4. Download full-resolution version (not thumbnail)
5. Document source in `CREDITS.md` (optional but recommended)

**Recommended Sources:**
- Unsplash: All images CC0 (verified)
- Pexels: Most images CC0 (check individual license)
- Pixabay: All images CC0 (verified)
- Wikimedia Commons: Mixed licenses (check each image)

[Source: Creative Commons CC0 License]

### Image Optimization Trade-offs

**Quality vs. Size Trade-off:**

| WebP Quality    | File Size  | Visual Quality | Use Case                         |
| --------------- | ---------- | -------------- | -------------------------------- |
| 100 (lossless)  | 400 KB     | Perfect        | Archival (too large for web)     |
| 90              | 220 KB     | Excellent      | Print quality (overkill for web) |
| 85              | 180 KB     | Very good      | High-quality web images          |
| **80 (target)** | **150 KB** | **Good**       | **Web standard (recommended)**   |
| 75              | 120 KB     | Acceptable     | Mobile-only, fast connections    |
| 70              | 100 KB     | Fair           | Budget phones, slow connections  |
| 60              | 80 KB      | Poor           | Visible compression artifacts    |

**Recipe Target: Quality 80**
- Good balance between file size and visual quality
- Compression artifacts minimal (not visible at normal viewing distance)
- File size within budget (<200 KB per image)
- Works well on 3G connections (<2 second load time)

**Optimization Strategy:**
1. Start with quality 80
2. Check file size (target <200 KB)
3. If too large: Reduce quality to 75 or 70
4. If quality poor: Increase quality to 85
5. Balance: Prefer smaller file size unless quality significantly degrades

[Source: Image Optimization Best Practices - web.dev]

### Responsive Image Strategy

**Why Multiple Image Sizes?**

| Device  | Viewport | Image Width | File Size | Network | Use Case       |
| ------- | -------- | ----------- | --------- | ------- | -------------- |
| Desktop | 1920px   | 1200px      | 150 KB    | WiFi    | Full detail    |
| Tablet  | 1024px   | 800px       | 80 KB     | WiFi/4G | Reduced detail |
| Mobile  | 375px    | 400px       | 30 KB     | 3G/4G   | Minimal detail |

**Bandwidth Savings:**
- Desktop user (1920px viewport): Loads 1200px image (150 KB)
- Mobile user (375px viewport): Loads 400px image (30 KB)
- **Savings: 120 KB** (80% reduction for mobile users)

**Implementation:**
```html
<picture>
  <source srcset="/images/preview-portrait-400w.webp" media="(max-width: 600px)">
  <source srcset="/images/preview-portrait-800w.webp" media="(max-width: 1024px)">
  <source srcset="/images/preview-portrait.webp">
  <img src="/images/preview-portrait.jpg" alt="Portrait" loading="lazy">
</picture>
```

**Browser Behavior:**
- Mobile (375px viewport): Loads 400w image (first `<source>` matches)
- Tablet (768px viewport): Loads 800w image (second `<source>` matches)
- Desktop (1920px viewport): Loads 1200px image (default `<source>`)
- Fallback: Loads JPEG if WebP not supported

[Source: Responsive Images - MDN Web Docs]

### Project Structure Notes

**New Files Created (Story 11.2):**
```
web/
└── images/
    ├── preview-portrait.webp         (Desktop, 1200px, ~150 KB)
    ├── preview-portrait-800w.webp    (Tablet, 800px, ~80 KB)
    ├── preview-portrait-400w.webp    (Mobile, 400px, ~30 KB)
    ├── preview-portrait.jpg          (Fallback, 1200px, ~180 KB)
    ├── preview-landscape.webp        (Desktop, 1200px, ~180 KB)
    ├── preview-landscape-800w.webp   (Tablet, 800px, ~100 KB)
    ├── preview-landscape-400w.webp   (Mobile, 400px, ~40 KB)
    ├── preview-landscape.jpg         (Fallback, 1200px, ~220 KB)
    ├── preview-product.webp          (Desktop, 1200px, ~120 KB)
    ├── preview-product-800w.webp     (Tablet, 800px, ~60 KB)
    ├── preview-product-400w.webp     (Mobile, 400px, ~25 KB)
    ├── preview-product.jpg           (Fallback, 1200px, ~140 KB)
    └── CREDITS.md                    (Attribution, licensing)
```

**Modified Files:**
- `web/index.html` - Add `<picture>` elements for responsive images
- `web/js/preview.js` - Add `showReferenceImage()` function (image switching)

**Integration Points:**
- Story 11.1: CSS filters applied to these reference images
- Story 11.3: Preview modal displays reference images
- Story 11.4: Slider shows before/after on reference images

[Source: docs/tech-spec-epic-11.md#Services-and-Modules]

### Testing Strategy

**Visual Quality Testing:**
- Manual inspection: Open each image at 100% zoom, check for artifacts
- Filter testing: Apply CSS filters, verify images still look good
- Slider testing: Test before/after slider, verify smooth comparison

**Performance Testing:**
- File size validation: Verify each image <200 KB (total <600 KB)
- Load time testing: Measure load time on 3G connection (<2 seconds)
- Lighthouse audit: Verify "Serve images in next-gen formats" passes

**Browser Compatibility Testing:**
- WebP support: Test on Chrome, Firefox, Safari, Edge
- JPEG fallback: Test on Safari 13, IE 11 (if applicable)
- Responsive images: Test on mobile, tablet, desktop viewports

**Licensing Verification:**
- Verify CC0 license on source websites
- Document sources in `CREDITS.md`
- Ensure no restrictive licenses (BY-SA, ND, NC)

[Source: docs/tech-spec-epic-11.md#Test-Strategy-Summary]

### Known Risks

**RISK-54: Reference images may not represent all photography genres**
- **Impact**: Users with niche photography styles (macro, astrophotography) may not find representative preview
- **Mitigation**: Choose diverse images (portrait, landscape, product cover 80%+ use cases)
- **Acceptable**: Phase 1 focuses on common genres, Phase 2 can add more categories

**RISK-55: WebP not supported on Safari 13 or older**
- **Impact**: Users on macOS Catalina (Safari 13) will load JPEG fallbacks (larger files)
- **Mitigation**: Provide JPEG fallbacks, acceptable 5% of users
- **Acceptable**: Safari 14+ (2020) supports WebP, most users upgraded

**RISK-56: Image file sizes exceed budget (>200 KB each)**
- **Impact**: Page load time exceeds 2 seconds on 3G connection
- **Mitigation**: Reduce WebP quality to 75 or 70, verify visual quality acceptable
- **Test**: Measure load times on WebPageTest (3G connection)

**RISK-57: Images don't work well with extreme preset adjustments**
- **Impact**: Preview looks unnatural with high exposure/contrast/saturation
- **Mitigation**: Choose neutral starting images, test with typical adjustments (±1.0 exposure, ±0.5 contrast)
- **Acceptable**: Extreme adjustments (exposure +2.0) expected to look unusual

[Source: docs/tech-spec-epic-11.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-11.md#AC-2] - Reference image bundle requirements
- [Source: docs/stories/11-1-css-filter-mapping.md] - CSS filter integration
- [WebP Image Format - Google Developers](https://developers.google.com/speed/webp)
- [Responsive Images - MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn/HTML/Multimedia_and_embedding/Responsive_images)
- [Image Optimization - web.dev](https://web.dev/fast/#optimize-your-images)
- [CC0 License - Creative Commons](https://creativecommons.org/publicdomain/zero/1.0/)
- [Unsplash License](https://unsplash.com/license) - CC0 (Public Domain)
- [Pexels License](https://www.pexels.com/license/) - CC0 (Public Domain)
- [Pixabay License](https://pixabay.com/service/license/) - CC0 (Public Domain)

## Dev Agent Record

### Context Reference

- `docs/stories/11-2-reference-image-bundle.context.xml` - Generated 2025-11-09

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
