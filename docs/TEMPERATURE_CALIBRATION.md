# Temperature Calibration Guide

## The Problem

Even with "Camera Flexible Color" as the base profile, Lightroom and NX Studio render Nikon RAW files with different color temperatures. Your screenshot shows:

- **NX Studio (left)**: Warmer, more golden tones in sky and sand
- **Lightroom (right)**: Cooler, more cyan/blue tones

This is because:
1. Adobe's "Camera Flexible Color" profile uses Adobe's interpretation of Nikon's color science
2. NX Studio uses Nikon's native color engine
3. The baseline color matrices are fundamentally different

## The Solution

Add a **Temperature offset** to the XMP profile that compensates for this baseline difference.

### Step 1: Measure the Shift

Looking at your screenshot:
- The sky should be warm (golden hour tones)
- Lightroom shows cooler/bluer tones
- Estimated shift needed: **+800 to +1200K warmer**

### Step 2: Add Temperature to Profile

We need to modify `GenerateProfileWithLUT()` to include a Temperature adjustment:

```go
// Add temperature compensation for baseline profile difference
Temperature2012: "+1000",  // Adjust this value to match NX Studio
```

### Step 3: Test and Iterate

1. Generate profile with Temperature offset
2. Apply to same RAW file in Lightroom
3. Compare with NX Studio side-by-side
4. Adjust Temperature value up/down by 100K increments
5. Regenerate and test until colors match

### Recommended Starting Values

- Start with: **+1000K** (makes Lightroom warmer)
- If still too cool: increase to +1200K
- If too warm: decrease to +800K

## Alternative Approach

If temperature alone doesn't solve it, we may need to:

1. **Use a different base profile**: Try "Camera Standard" or "Adobe Standard"
2. **Add Tint adjustment**: Compensate for green/magenta shifts
3. **Create a custom DCP file**: Build a true camera profile instead of XMP
4. **Apply calibration adjustments**: Use CameraCalibration tags for primary color shifts

## Next Steps

I can implement the temperature offset now. What temperature shift would you like me to try first?

Recommended: **+1000K** based on your screenshot comparison.
