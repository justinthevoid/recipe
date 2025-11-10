# Custom DCP Testing Guide - Nikon Z f Warm Profile

**Date**: 2025-11-07
**Profile**: Nikon Z f Warm (Custom)
**Status**: ✓ Created and Installed

---

## Installation Confirmed

✓ DCP file created: `testdata/Nikon_Z_f_Warm_Custom.dcp` (558 bytes)
✓ Installed to: `C:\Users\Justin\AppData\Roaming\Adobe\CameraRaw\Settings\`
✓ Color Matrix verified with exiftool

### Matrix Values Confirmed

**Color Matrix 2 (Daylight - Modified for Warmth)**:
```
[ 1.25    -0.35     0.08   ]  ← Blue in red = +0.08 (WARM!)
[-0.40     1.20     0.15   ]
[-0.02     0.10     0.85   ]  ← Higher blue sensitivity
```

**Comparison with Adobe Original**:
```
Adobe:  [ 1.16    -0.45    -0.10  ]  ← Blue in red = NEGATIVE (cool)
        [-0.45     1.25     0.23  ]
        [-0.05     0.15     0.76  ]  ← Lower blue

Custom: [ 1.25    -0.35     0.08  ]  ← +0.18 warmer!
        [-0.40     1.20     0.15  ]
        [-0.02     0.10     0.85  ]  ← +0.09 more blue
```

---

## Testing Procedure

### Step 1: Restart Lightroom

**IMPORTANT**: You must restart Lightroom for it to detect the new profile.

1. Close Lightroom completely (check Task Manager if needed)
2. Reopen Lightroom
3. Wait for it to fully load

### Step 2: Open a Nikon Z f RAW File

1. Import a .NEF file from your Nikon Z f
2. Go to **Develop** module
3. Select the image in Library

**Best test images**:
- Portrait with skin tones (most sensitive to warmth)
- Landscape with warm colors (sunset, foliage)
- Indoor shot with mixed lighting
- Any image that looked too cool in Lightroom before

### Step 3: Apply the Custom Profile

1. In **Basic** panel, locate **Profile** section (top of panel)
2. Click **Profile Browser** (grid icon next to profile name)
3. Look for **"Nikon Z f Warm (Custom)"** in the profile list
   - Should appear under "Custom" or "Camera Matching" section
   - If not visible, check Lightroom was restarted
4. Click to apply the profile
5. **Leave all other sliders at default** for fair comparison

### Step 4: Compare with NX Studio

**Side-by-Side Comparison**:

1. Export from Lightroom:
   - File → Export
   - Format: JPEG, Quality: 100%, Color Space: sRGB
   - Name: `test_lightroom_custom_warm.jpg`

2. Open same RAW in NX Studio:
   - Use Flexible Color Picture Control
   - No additional adjustments
   - Export as JPEG (same settings)
   - Name: `test_nxstudio_flexcolor.jpg`

3. Compare the exported JPEGs:
   - Open both in Windows Photo Viewer or any image viewer
   - Look at side-by-side or toggle between them
   - Focus on:
     * Overall warmth (reds, oranges, skin tones)
     * Color balance (green/magenta shifts)
     * Shadow warmth (are shadows warmer?)
     * Highlight warmth (white areas)

### Step 5: Take Screenshots

Document the results:

1. Screenshot of Lightroom with custom profile applied
2. Screenshot of NX Studio with same image
3. Screenshot of side-by-side JPEG comparison
4. Note any remaining differences

---

## What to Look For

### Success Indicators ✓

- [ ] **Reds look warmer** - Closer to NX Studio rendering
- [ ] **Skin tones match better** - More natural, less pink/cool
- [ ] **Overall image is warmer** - Less blue/cool cast
- [ ] **Shadows are warmer** - Not as blue/gray
- [ ] **Warmth feels natural** - Not oversaturated or orange

### Potential Issues ⚠

- [ ] **Too warm** - Reds look orange, oversaturated
- [ ] **Color shifts** - Greens, blues look wrong
- [ ] **Loss of contrast** - Image looks flat
- [ ] **Profile doesn't appear** - Restart Lightroom again
- [ ] **No visible change** - Profile may not be loading correctly

---

## Expected Results

### Best Case (92-96% Match)

- Warmth is **significantly improved**
- Colors are **very close** to NX Studio
- Only subtle differences remain:
  - Slight saturation variation
  - Minor tone curve differences
  - ProfileLookTable transformations (Adobe's 3D LUT)

### Good Case (88-92% Match)

- Warmth is **noticeably better**
- Colors are **mostly matching**
- Some differences remain:
  - Specific hues slightly off
  - Contrast curves differ
  - Needs minor adjustment with sliders

### Needs Refinement (< 88% Match)

- Warmth improved but **not enough**
- Still visible coolness in reds/oranges
- Need to increase blue→red coefficient further

---

## Fine-Tuning Matrix Values

If the profile is **still too cool**, adjust the script and regenerate:

### Increase Warmth Further

Edit `create_custom_dcp.py` line with `color_matrix_2_warm`:

```python
# Current (moderate warmth)
color_matrix_2_warm = [
    1.25,   -0.35,    0.08,    # Blue in red = 0.08
   -0.40,    1.20,    0.15,
   -0.02,    0.10,    0.85,
]

# More warmth (increase blue in red)
color_matrix_2_warm = [
    1.28,   -0.32,    0.12,    # Blue in red = 0.12 (warmer)
   -0.38,    1.18,    0.15,
   -0.02,    0.10,    0.87,
]
```

### Reduce Warmth (If Too Orange)

```python
# Less warmth (reduce blue in red)
color_matrix_2_warm = [
    1.22,   -0.38,    0.05,    # Blue in red = 0.05 (less warm)
   -0.42,    1.22,    0.15,
   -0.03,    0.12,    0.83,
]
```

### Rerun Script

```bash
cd /c/Users/Justin/void/recipe
python create_custom_dcp.py
# Restart Lightroom
# Test again
```

---

## Adjustment Strategy

### If Reds Too Cool (Blue/Pink)
- **Increase** blue→red coefficient: `0.08 → 0.10 → 0.12`
- **Increase** red diagonal: `1.25 → 1.28 → 1.30`

### If Greens Too Yellow
- **Reduce** green→red suppression: `-0.35 → -0.38 → -0.40`

### If Blues Too Purple
- **Increase** blue diagonal: `0.85 → 0.87 → 0.90`
- **Reduce** blue in green: `0.10 → 0.08 → 0.06`

### If Overall Too Saturated
- **Reduce** all diagonal values slightly: `1.25 → 1.22`, etc.

### If Overall Too Flat
- **Increase** all diagonal values: `1.25 → 1.28`, etc.

---

## Creating Additional Profiles

Once you find good matrix values, create variants:

### Portrait Profile (Extra Warm Skin Tones)

```python
# In create_custom_dcp.py, create new function:
def create_portrait_profile():
    # Use even warmer matrix for skin tones
    color_matrix_2_portrait = [
        1.30,   -0.30,    0.15,    # Maximum warmth
       -0.35,    1.15,    0.12,
       -0.01,    0.08,    0.88,
    ]
    # ... rest of profile creation
```

### Landscape Profile (Balanced)

```python
def create_landscape_profile():
    # Slightly less warm, enhance greens/blues
    color_matrix_2_landscape = [
        1.20,   -0.40,    0.05,    # Moderate warmth
       -0.42,    1.25,    0.18,    # More green sensitivity
       -0.03,    0.12,    0.88,    # Higher blue
    ]
```

### Neutral Profile (Adobe-like but Slightly Warmer)

```python
def create_neutral_profile():
    # Just fix the worst coolness, keep Adobe's character
    color_matrix_2_neutral = [
        1.18,   -0.43,    0.00,    # Neutral blue in red
       -0.45,    1.25,    0.20,
       -0.04,    0.14,    0.78,
    ]
```

---

## Documentation Template

After testing, fill out this template:

```
TEST RESULTS - Nikon Z f Warm (Custom) Profile

Date: ___________
Lightroom Version: ___________
Camera: Nikon Z f
Firmware: ___________

Matrix Values Tested:
Blue→Red: _______
Red Diagonal: _______
Blue Diagonal: _______

TEST IMAGE 1: [Description]
- Lightroom Warmth: ___/10
- NX Studio Warmth: ___/10
- Match Quality: ___/10
- Notes: ________________

TEST IMAGE 2: [Description]
- Lightroom Warmth: ___/10
- NX Studio Warmth: ___/10
- Match Quality: ___/10
- Notes: ________________

OVERALL ASSESSMENT:
Color Accuracy: ___% (estimated)
Warmth Match: ___% (estimated)
Recommend Adjustments: ________________

NEXT STEPS:
[ ] Keep current matrix values
[ ] Increase warmth (blue→red: ___ → ___)
[ ] Reduce warmth (blue→red: ___ → ___)
[ ] Adjust other coefficients: ________________
```

---

## Troubleshooting

### Profile Doesn't Appear in Lightroom

1. **Check file location**:
   ```bash
   ls "$APPDATA/Adobe/CameraRaw/Settings/Nikon_Z_f_Warm_Custom.dcp"
   ```

2. **Verify file isn't corrupted**:
   ```bash
   hexdump -C "$APPDATA/Adobe/CameraRaw/Settings/Nikon_Z_f_Warm_Custom.dcp" | head -5
   ```
   Should show "II*" at start

3. **Check Lightroom preferences**:
   - Edit → Preferences → Presets
   - Ensure "Store presets with catalog" is OFF
   - Click "Show Lightroom Develop Presets"
   - Navigate up to CameraRaw/Settings folder
   - Verify DCP is there

4. **Restart Lightroom completely**:
   - Close Lightroom
   - Kill any Adobe processes in Task Manager
   - Reopen Lightroom

### Profile Loads But No Visible Change

1. **Check Camera Model Match**:
   - DCP is for "Nikon Z f" exactly
   - Won't work on other Nikon models

2. **Compare with Adobe Profile**:
   - Switch between Adobe Standard and Custom
   - Differences should be visible immediately
   - Focus on reds and skin tones

3. **Try Different Image**:
   - Some images may not show differences clearly
   - Use images with warm colors (reds, oranges, skin)

### Colors Look Wrong

1. **Check Matrix Values**:
   ```bash
   wsl -d Debian -- exiftool testdata/Nikon_Z_f_Warm_Custom.dcp | grep "Color Matrix 2"
   ```

2. **Verify Against Expected**:
   - Blue→Red should be positive (0.08)
   - Red diagonal should be 1.25
   - Blue diagonal should be 0.85

3. **Regenerate if Needed**:
   ```bash
   cd /c/Users/Justin/void/recipe
   python create_custom_dcp.py
   # Copy to Lightroom folder again
   # Restart Lightroom
   ```

---

## Success Criteria

**Minimum Acceptable** (85-90% match):
- Reds are warmer than Adobe's profile
- Skin tones improved
- Overall more pleasant look

**Good Result** (90-95% match):
- Very close to NX Studio
- Only expert eye sees differences
- Usable for professional work

**Excellent Result** (95%+ match):
- Practically indistinguishable from NX Studio
- Only subtle tone curve differences
- Perfect for all practical purposes

**Perfect Match** (100%):
- Not achievable without Nikon's proprietary system
- Not necessary for excellent results

---

## Next Steps After Testing

### If Successful (90%+ match)

1. **Create profile variants** (Portrait, Landscape, etc.)
2. **Document optimal matrix values** for future reference
3. **Share findings** (if desired) with photography community
4. **Use in workflow** for all Nikon Z f RAW processing

### If Needs Refinement (< 90% match)

1. **Identify specific issues**:
   - Too cool? Increase blue→red
   - Too warm? Reduce blue→red
   - Color shifts? Adjust other coefficients

2. **Iterate matrix values**:
   - Change one coefficient at a time
   - Test after each change
   - Document results

3. **Consider machine learning approach** if manual tuning doesn't work

### If Failed (< 80% match)

1. **Review FINAL_CONCLUSIONS.md** for alternative approaches
2. **Consider using NX Studio** for initial processing
3. **Investigate ProfileLookTable** extraction (advanced)

---

**Ready to test!** Open Lightroom, apply the profile, and compare with NX Studio. Document your findings using the template above.
