# Reference Image Preparation Guide (Story 11.2)

This guide provides step-by-step instructions for selecting, downloading, optimizing, and integrating the three reference images for Recipe's preset preview system.

## Prerequisites

Install required image processing tools:

### macOS
```bash
brew install webp imagemagick
```

### Windows
1. **WebP Tools**: Download from https://developers.google.com/speed/webp/download
   - Extract to `C:\Program Files\WebP\bin\`
   - Add to PATH: System Properties → Environment Variables → Path → Add `C:\Program Files\WebP\bin`

2. **ImageMagick**: Download from https://imagemagick.org/script/download.php#windows
   - Run installer, check "Add application directory to system path"

### Linux (Debian/Ubuntu)
```bash
sudo apt-get install webp imagemagick
```

---

## Step 1: Select Reference Images

### Portrait Image Selection

**Search Keywords**: "portrait headshot person natural light"

**Recommended Sources**:
- Unsplash: https://unsplash.com/s/photos/portrait
- Pexels: https://www.pexels.com/search/portrait/
- Pixabay: https://pixabay.com/images/search/portrait/

**Selection Criteria**:
- ✅ Natural, soft lighting (not harsh shadows)
- ✅ Visible skin tones (light, medium, or dark)
- ✅ Neutral background (not distracting)
- ✅ Friendly, neutral expression
- ✅ High resolution (2000px+ width minimum)
- ✅ Good tonal range (shadows, midtones, highlights)
- ✅ **MUST be CC0 license** (verify on source page)

**Example Search Results**:
- Unsplash: Look for portraits with "natural light" tag
- Pexels: Filter by "Free to use" license
- Pixabay: All images are CC0 by default

**Download**:
- Download **full resolution** version (not thumbnail)
- Save as `portrait-original.jpg` in `web/images/` directory

---

### Landscape Image Selection

**Search Keywords**: "landscape mountains sky nature outdoor"

**Recommended Sources**:
- Unsplash: https://unsplash.com/s/photos/landscape
- Pexels: https://www.pexels.com/search/landscape/
- Pixabay: https://pixabay.com/images/search/landscape/

**Selection Criteria**:
- ✅ Blue sky visible (tests saturation, hue adjustments)
- ✅ Green foliage visible (tests color shifts)
- ✅ Good tonal range (shadows in foreground, highlights in sky)
- ✅ Neutral exposure (not overexposed, not underexposed)
- ✅ Daylight scene (not sunset/sunrise with extreme colors)
- ✅ High resolution (2000px+ width minimum)
- ✅ **MUST be CC0 license**

**Example Search Results**:
- Unsplash: Look for "mountain landscape" or "nature valley"
- Pexels: Filter for "outdoor nature" with blue sky
- Pixabay: Search "landscape mountains"

**Download**:
- Download **full resolution** version
- Save as `landscape-original.jpg` in `web/images/` directory

---

### Product Image Selection

**Search Keywords**: "product still life studio white background"

**Recommended Sources**:
- Unsplash: https://unsplash.com/s/photos/product-photography
- Pexels: https://www.pexels.com/search/product/
- Pixabay: https://pixabay.com/images/search/product/

**Selection Criteria**:
- ✅ Neutral background (white, gray, or black)
- ✅ Varied product colors (tests saturation, hue)
- ✅ Even, studio-quality lighting (no harsh shadows)
- ✅ Simple composition (not cluttered)
- ✅ High resolution (2000px+ width minimum)
- ✅ **MUST be CC0 license**

**Example Search Results**:
- Unsplash: Look for "product photography white background"
- Pexels: Filter for "still life" or "product shoot"
- Pixabay: Search "product white background"

**Download**:
- Download **full resolution** version
- Save as `product-original.jpg` in `web/images/` directory

---

## Step 2: Optimize Images

Run these commands from the `web/images/` directory:

### Resize to Multiple Sizes

```bash
# Portrait (3:4 aspect ratio recommended)
convert portrait-original.jpg -resize 1200x -quality 100 portrait-1200w.jpg
convert portrait-original.jpg -resize 800x -quality 100 portrait-800w.jpg
convert portrait-original.jpg -resize 400x -quality 100 portrait-400w.jpg

# Landscape (16:9 or 3:2 aspect ratio recommended)
convert landscape-original.jpg -resize 1200x -quality 100 landscape-1200w.jpg
convert landscape-original.jpg -resize 800x -quality 100 landscape-800w.jpg
convert landscape-original.jpg -resize 400x -quality 100 landscape-400w.jpg

# Product (1:1 or 4:3 aspect ratio recommended)
convert product-original.jpg -resize 1200x -quality 100 product-1200w.jpg
convert product-original.jpg -resize 800x -quality 100 product-800w.jpg
convert product-original.jpg -resize 400x -quality 100 product-400w.jpg
```

### Convert to WebP Format

```bash
# Desktop WebP (quality 80 - target <200KB each)
cwebp -q 80 portrait-1200w.jpg -o preview-portrait.webp
cwebp -q 80 landscape-1200w.jpg -o preview-landscape.webp
cwebp -q 80 product-1200w.jpg -o preview-product.webp

# Tablet WebP (quality 80 - target <100KB each)
cwebp -q 80 portrait-800w.jpg -o preview-portrait-800w.webp
cwebp -q 80 landscape-800w.jpg -o preview-landscape-800w.webp
cwebp -q 80 product-800w.jpg -o preview-product-800w.webp

# Mobile WebP (quality 75 - target <40KB each)
cwebp -q 75 portrait-400w.jpg -o preview-portrait-400w.webp
cwebp -q 75 landscape-400w.jpg -o preview-landscape-400w.webp
cwebp -q 75 product-400w.jpg -o preview-product-400w.webp
```

### Create JPEG Fallbacks

```bash
# Desktop JPEG fallbacks (quality 85, strip EXIF)
convert portrait-1200w.jpg -quality 85 -strip preview-portrait.jpg
convert landscape-1200w.jpg -quality 85 -strip preview-landscape.jpg
convert product-1200w.jpg -quality 85 -strip preview-product.jpg
```

---

## Step 3: Verify File Sizes

```bash
# Check WebP file sizes
ls -lh preview-*.webp

# Expected output (approximate):
# preview-portrait.webp        ~150 KB (target <150KB)
# preview-portrait-800w.webp   ~80 KB  (target <80KB)
# preview-portrait-400w.webp   ~30 KB  (target <30KB)
# preview-landscape.webp       ~180 KB (target <180KB)
# preview-landscape-800w.webp  ~100 KB (target <100KB)
# preview-landscape-400w.webp  ~40 KB  (target <40KB)
# preview-product.webp         ~120 KB (target <120KB)
# preview-product-800w.webp    ~60 KB  (target <60KB)
# preview-product-400w.webp    ~25 KB  (target <25KB)

# Total WebP bundle: ~780 KB (desktop), ~400 KB (tablet), ~95 KB (mobile)
```

**If file sizes exceed targets**:
```bash
# Reduce quality for oversized images
cwebp -q 75 portrait-1200w.jpg -o preview-portrait.webp    # Try quality 75
cwebp -q 70 portrait-1200w.jpg -o preview-portrait.webp    # Or quality 70
```

**Visual quality check**:
- Open each WebP file in browser at 100% zoom
- Check for pixelation, compression artifacts
- Verify colors look natural (not oversaturated)

---

## Step 4: Update CREDITS.md

Edit `web/images/CREDITS.md` and fill in the photographer details:

```markdown
## Portrait

- **Image**: Preview portrait reference
- **Photographer**: John Doe
- **Source**: Unsplash
- **License**: CC0 (Public Domain)
- **Original URL**: https://unsplash.com/photos/abc123
- **Selection Criteria**: Natural lighting, visible skin tones, neutral background, soft expression
```

Repeat for Landscape and Product sections.

---

## Step 5: Clean Up Temporary Files

```bash
# Remove intermediate files (keep only optimized versions)
rm portrait-original.jpg portrait-1200w.jpg portrait-800w.jpg portrait-400w.jpg
rm landscape-original.jpg landscape-1200w.jpg landscape-800w.jpg landscape-400w.jpg
rm product-original.jpg product-1200w.jpg product-800w.jpg product-400w.jpg
```

**Final directory structure**:
```
web/images/
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
├── preview-product.jpg           (Fallback, ~140 KB)
├── CREDITS.md                    (Attribution)
└── IMAGE_PREPARATION_GUIDE.md    (This file)
```

---

## Step 6: Visual Validation

### Test with CSS Filters (Story 11.1 Integration)

Create a test HTML file to preview filters:

```html
<!DOCTYPE html>
<html>
<head>
<style>
  img { max-width: 100%; }
  .exposure-high { filter: brightness(150%); }
  .exposure-low { filter: brightness(50%); }
  .contrast-high { filter: contrast(150%); }
  .saturation-high { filter: saturate(150%); }
  .hue-shift { filter: hue-rotate(30deg); }
</style>
</head>
<body>
  <h1>Portrait Reference Image</h1>
  <img src="preview-portrait.webp" alt="Original">
  <img src="preview-portrait.webp" class="exposure-high" alt="Bright">
  <img src="preview-portrait.webp" class="contrast-high" alt="High Contrast">
  <img src="preview-portrait.webp" class="saturation-high" alt="Saturated">

  <!-- Repeat for landscape and product -->
</body>
</html>
```

**Validation Checklist**:
- [ ] `brightness(150%)` - Image should not blow out highlights
- [ ] `brightness(50%)` - Image should not crush blacks
- [ ] `contrast(150%)` - Image should show improved separation
- [ ] `saturate(150%)` - Image should look more colorful (not garish)
- [ ] `hue-rotate(30deg)` - Image should show subtle color shift

---

## Step 7: Browser Compatibility Testing

Test WebP support:
1. Chrome (latest): Open `preview-portrait.webp` - Should load WebP
2. Firefox (latest): Open `preview-portrait.webp` - Should load WebP
3. Safari (latest): Open `preview-portrait.webp` - Should load WebP (Safari 14+)
4. Edge (latest): Open `preview-portrait.webp` - Should load WebP

Test JPEG fallback:
- Safari 13 or IE11 (if available): Should load JPEG fallback

---

## Troubleshooting

### cwebp command not found
- macOS: `brew install webp`
- Windows: Add WebP bin directory to PATH
- Linux: `sudo apt-get install webp`

### convert command not found
- Install ImageMagick (see Prerequisites section)

### File sizes too large
- Reduce WebP quality: Try `-q 75` or `-q 70`
- Check original image resolution (should be 2000px+, but not excessive like 6000px+)

### Images look pixelated
- Increase WebP quality: Try `-q 85` or `-q 90`
- Verify original images are high resolution

### Colors look washed out
- Check for color profile issues: Use `-strip` flag with ImageMagick
- Verify original images have good tonal range

---

## AC Validation Checklist

Before marking Story 11.2 complete, verify:

**AC-1: Three Reference Images Provided**
- [ ] Portrait image selected (person, natural lighting, skin tones)
- [ ] Landscape image selected (outdoor, blue sky, greenery)
- [ ] Product image selected (inanimate, studio lighting, neutral background)
- [ ] Images representative of 80%+ typical use cases
- [ ] Images have neutral starting point (not extreme colors/exposure)

**AC-2: Image Optimization for Web**
- [ ] WebP format primary (quality 80-85%)
- [ ] JPEG fallback (quality 85-90%)
- [ ] File sizes: Portrait <150KB, Landscape <180KB, Product <120KB
- [ ] Total bundle <600KB (desktop), <400KB (tablet), <200KB (mobile)
- [ ] Dimensions: Desktop 1200px, Tablet 800px, Mobile 400px
- [ ] EXIF metadata stripped

**AC-3: Licensing and Attribution**
- [ ] All images CC0 (public domain)
- [ ] CREDITS.md updated with photographer details
- [ ] No licensing conflicts

**AC-4: Images Embedded in Web Bundle**
- [ ] All images in `web/images/` directory
- [ ] No external requests (self-hosted)
- [ ] Correct filenames (preview-portrait.webp, etc.)

**AC-5: Images Work Well with Typical Presets**
- [ ] Neutral starting point verified
- [ ] Tested with exposure ±1.0
- [ ] Tested with contrast ±0.5
- [ ] Tested with saturation ±0.5
- [ ] Tested with hue ±30°

**AC-6: Responsive Image Handling**
- [ ] Multiple sizes generated (1200px, 800px, 400px)
- [ ] Ready for `<picture>` element integration (Story 11.3)

**AC-7: Image Selection and Validation**
- [ ] Visual quality acceptable at 100% zoom
- [ ] CSS filters render correctly
- [ ] Browser compatibility verified
