# Split Toning Bug Fix

**Date**: 2025-12-18
**Issue**: Split Toning not being written to NP3 Color Grading zones
**Status**: ✅ Fixed in [internal/formats/np3/generate.go:371-385](internal/formats/np3/generate.go:371-385)

---

## Problem Summary

Split Toning values (Shadow Hue/Saturation, Highlight Hue/Saturation) were being parsed correctly from XMP but **not written** to the NP3 file's Color Grading zones (Highlights, Midtone, Shadows).

**Test Case**: Agfachrome RSX 200 preset
- Split Toning Shadow: Hue 215°, Saturation 13
- Split Toning Highlight: Hue 35°, Saturation 12

**Before Fix** (hexdump at offset 0x170):
```
00 00 80 80  00 dc 8a 80  00 00 80 80
[Shadows  ]  [Midtone  ]  [Highlights]
Hue 0, C 0   Hue220,C10  Hue 0, C 0
             ^^^^^^^^^^^
             Only midtone populated!
```

---

## Root Cause

The XMP parser creates a `ColorGrading` object whenever ANY Color Grading parameter exists. In the Agfachrome preset, only `ColorGradeMidtoneHue=220` and `ColorGradeMidtoneSat=10` exist, but the parser creates the FULL `ColorGrading` struct with **empty** `Highlights` and `Shadows` zones.

**Code Flow**:

1. **XMP Parser** ([internal/formats/xmp/parse.go:896-975](internal/formats/xmp/parse.go:896-975)):
   ```go
   func extractColorGrading(desc *Description) (*models.ColorGrading, error) {
       // Creates ColorGrading object if ANY parameter exists
       cg := &models.ColorGrading{}  // ← Creates empty zones!

       // Only fills Midtone (Hue 220, Chroma 10)
       cg.Midtone.Hue = 220
       cg.Midtone.Chroma = 10

       // Highlights and Shadows remain empty {Hue: 0, Chroma: 0}
       return cg, nil
   }
   ```

2. **NP3 Generator** ([internal/formats/np3/generate.go:372-378](internal/formats/np3/generate.go:372-378) - BEFORE FIX):
   ```go
   if recipe.ColorGrading != nil {
       params.highlightsZone = recipe.ColorGrading.Highlights  // ← Empty {0, 0}!
       params.midtoneZone = recipe.ColorGrading.Midtone        // ✅ Hue 220, Chroma 10
       params.shadowsZone = recipe.ColorGrading.Shadows        // ← Empty {0, 0}!
   }
   ```

3. **Split Toning Fallback** ([internal/formats/np3/generate.go:389-394](internal/formats/np3/generate.go:389-394)):
   ```go
   // This condition FAILS because zones are not zero - they're empty structs!
   if params.shadowsZone.Chroma == 0 && recipe.SplitShadowSaturation > 0 {
       params.shadowsZone.Hue = recipe.SplitShadowHue           // Never executes!
       params.shadowsZone.Chroma = recipe.SplitShadowSaturation // Never executes!
   }
   ```

**The Bug**: Empty `ColorGradingZone` structs (`{Hue: 0, Chroma: 0}`) were being assigned to `params.highlightsZone` and `params.shadowsZone`, overwriting the zero values that would have allowed the Split Toning fallback to trigger.

---

## Fix

Changed [internal/formats/np3/generate.go:371-385](internal/formats/np3/generate.go:371-385) to **only assign zones that have actual data** (non-zero chroma):

```go
// Map Color Grading if present
// Only assign zones that have actual data (non-zero chroma) to avoid overwriting Split Toning fallback
if recipe.ColorGrading != nil {
    if recipe.ColorGrading.Highlights.Chroma != 0 {
        params.highlightsZone = recipe.ColorGrading.Highlights
    }
    if recipe.ColorGrading.Midtone.Chroma != 0 {
        params.midtoneZone = recipe.ColorGrading.Midtone
    }
    if recipe.ColorGrading.Shadows.Chroma != 0 {
        params.shadowsZone = recipe.ColorGrading.Shadows
    }
    params.blending = recipe.ColorGrading.Blending
    params.balance = recipe.ColorGrading.Balance
}
```

**Logic**: If a Color Grading zone has `Chroma == 0`, it's empty and should NOT overwrite the zero value that allows Split Toning fallback to work.

---

## Verification

**After Fix** (hexdump at offset 0x170):
```
00 de 88 80  00 ea 85 80  00 1f 90 80
[Shadows  ]  [Midtone  ]  [Highlights]
Hue222,C 8   Hue234,C 5   Hue 31, C16
^^^^^^^^^^^                ^^^^^^^^^^^
Split Toning now applied!
```

**Decoded Values**:
- Shadows: Hue 222°, Chroma 8
- Midtone: Hue 234°, Chroma 5
- Highlights: Hue 31°, Chroma 16

**Why values differ from Split Toning input (215°/13 and 35°/12)?**

The `applyWhiteBalanceToColorGrading()` function ([internal/formats/np3/generate.go:1331-1413](internal/formats/np3/generate.go:1331-1413)) runs AFTER Split Toning fallback and adds white balance vectors to ALL color grading zones:

- IncrementalTemperature: 9 → Orange vector at 40° with chroma 4.5
- IncrementalTint: 3 → Magenta vector at 300° with chroma 1.5

These vectors are added using Cartesian coordinate math (lines 1374-1397), which modifies both hue and chroma. **This is correct behavior** - white balance should be additive with Split Toning.

**Verification**:
```bash
./bin/recipe inspect "testdata/visual-regression/test_output/Agfachrome RSX 200.np3"
```

Output confirms Color Grading zones are populated:
```json
"colorGrading": {
  "highlights": {
    "hue": 31,
    "chroma": 16
  },
  "midtone": {
    "hue": 234,
    "chroma": 5
  },
  "shadows": {
    "hue": 222,
    "chroma": 8
  },
  "blending": 50,
  "balance": -7
}
```

---

## Expected Visual Impact

**Color Grading zones add**:
- **Shadows** (Hue 222° ≈ Blue, Chroma 8): Cool blue tint in shadows
- **Highlights** (Hue 31° ≈ Orange, Chroma 16): Warm orange tint in highlights

This is the **classic film emulation look** of the Agfachrome RSX 200 preset - cool shadows, warm highlights.

**Expected Delta E Improvement**:
- Before fix: Shadows and Highlights had NO tint (Chroma 0)
- After fix: Shadows and Highlights have correct tinting
- Estimated improvement: **5-8 Delta E reduction** (from 19.40 to ~11-14)

**Note**: To measure actual improvement, the test image (`JMH_0079_nx.TIF`) would need to be **regenerated** in NX Studio using the NEW NP3 file. Current test image was created with the OLD (buggy) NP3 file.

---

## Related Issues

User observation: **"The blues and greens are much more vibrant in Lightroom"**

Split Toning adds:
- Blue hue (215°) to shadows at 13% saturation
- Warm orange (35°) to highlights at 12% saturation

Missing this would make:
- Shadows appear less saturated/vibrant (missing blue tint)
- Highlights appear less warm (missing orange tint)
- Overall image appear flatter/less colorful

**This fix directly addresses the user's observation!**

---

## Next Steps

1. ✅ **Fix implemented and verified** (Split Toning now writes to NP3)
2. ⏳ **Regenerate test image** (requires NX Studio on Windows with new NP3 file)
3. ⏳ **Measure Delta E improvement** (run visual comparison script after regeneration)
4. ⏳ **Investigate remaining vibrance/saturation issues** (if Delta E still >5 after fix)

---

## Technical Details

**NP3 Color Grading Format** (offset 0x170):
- Each zone: 4 bytes (2-byte hue, 1-byte chroma, 1-byte brightness)
- Hue: 12-bit value (0-360°) stored in 2 bytes with 0x0F mask
- Chroma: Signed8 encoding (value + 0x80), range 0-100
- Brightness: Signed8 encoding, range -100 to +100

**Encoding**:
```go
// Hue encoding (2 bytes)
byte1, byte2 := EncodeHue12(hue)  // Example: 222° = 0x00DE
data[offset] = byte1     // 0x00
data[offset+1] = byte2   // 0xDE

// Chroma encoding (1 byte, Signed8)
data[offset+2] = EncodeSigned8(chroma)  // Example: 8 → 0x88 (8 + 0x80)

// Brightness encoding (1 byte, Signed8)
data[offset+3] = EncodeSigned8(brightness)  // Example: 0 → 0x80
```

**Decoding**:
```go
hue := (int(byte1 & 0x0F) << 8) | int(byte2)  // 0x00DE = 222°
chroma := int(chromaByte) - 0x80              // 0x88 - 0x80 = 8
brightness := int(brightnessByte) - 0x80      // 0x80 - 0x80 = 0
```

---

## Conclusion

The Split Toning bug has been **fixed**. The NP3 generator now correctly writes Split Toning values to the Color Grading zones when Color Grading zones are empty.

**Root Cause**: XMP parser creates empty ColorGrading zones, which were overwriting the zero values needed for Split Toning fallback.

**Solution**: Check if zones have non-zero chroma before assigning, allowing Split Toning fallback to work correctly.

**Expected Improvement**: 5-8 Delta E reduction (user needs to regenerate test image in NX Studio to measure).
