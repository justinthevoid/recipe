# DCP Test Samples

This directory contains DNG Camera Profile (.dcp) files used for testing Recipe's DCP parser.

## Sample Files

### Synthetic Samples (Generated for Testing)

- **`minimal-linear.dcp`** - Minimal DCP with linear tone curve (no adjustments)
  - Tone curve: (0.0, 0.0) → (1.0, 1.0) linear (normalized 0.0-1.0 range)
  - Color matrices: Identity matrices
  - ProfileName: "Minimal Linear Test Profile"
  - BaselineExposure: 0.0

- **`portrait-adjusted.dcp`** - Portrait-style DCP with tone adjustments
  - Tone curve: Slight S-curve for portrait look (normalized 0.0-1.0 range)
  - Exposure: +0.3 stops
  - Contrast: +15
  - Highlights: -10 (recover highlights)
  - Shadows: +10 (lift shadows)
  - BaselineExposure: +0.3

### Known Limitations

Due to DCP files being proprietary Adobe format:
- Synthetic DCPs may not render identically in Lightroom/ACR
- Full DCP validation requires Adobe software (Story 9-4)
- Color calibration matrices not tested (MVP limitation)

## DCP File Structure

**⚠️ IMPORTANT**: Real Adobe DCP files use **binary TIFF tags (50700-52600)**, NOT XML in tag 50740. This was discovered during implementation through analysis of actual Adobe DCP files from Nikon, Hasselblad, and Leica cameras.

```
DCP (.dcp) = TIFF/DNG Container
├── TIFF/DNG Header (II/MM for TIFF, IIRC/MMCR for DNG)
├── Image File Directory (IFD)
│   ├── Tag 52552: ProfileName (ASCII string, OPTIONAL)
│   ├── Tag 50940: ProfileToneCurve (Float32 array of input/output pairs)
│   ├── Tag 50721: ColorMatrix1 (9 SRational values)
│   ├── Tag 50722: ColorMatrix2 (9 SRational values)
│   ├── Tag 50730: BaselineExposureOffset (SRational)
│   └── Other standard TIFF/DNG tags (50700-52600 range)
```

### Binary Data Formats

**ProfileToneCurve (Tag 50940)**:
- Array of float32 pairs (input, output), normalized to 0.0-1.0 range
- Each pair is 8 bytes (2 × float32 = 4 bytes each)
- Example: `{0.0, 0.0}, {0.5, 0.5}, {1.0, 1.0}` = linear curve
- Total size: N points × 8 bytes (where N is typically 5-17)

**ColorMatrix1/2 (Tags 50721-50722)**:
- 9 SRational values in row-major order: `[R_R, R_G, R_B, G_R, G_G, G_B, B_R, B_G, B_B]`
- Each SRational is 8 bytes (signed int32 numerator + signed int32 denominator)
- Total size: 72 bytes (9 × 8 bytes)
- Identity matrix: `{1,1}, {0,1}, {0,1}, {0,1}, {1,1}, {0,1}, {0,1}, {0,1}, {1,1}` (diagonal 1.0, off-diagonal 0.0)

**ProfileName (Tag 52552)**:
- ASCII string, null-terminated
- **OPTIONAL** - many DCP files do not include this tag
- If missing, Recipe uses empty string (caller can use filename instead)

**BaselineExposureOffset (Tag 50730)**:
- SRational value (signed int32 numerator / signed int32 denominator)
- Exposure compensation offset in stops
- Added to tone curve exposure calculation

### DNG Version Conversion

DNG files use magic bytes **IIRC** (little-endian) or **MMCR** (big-endian) instead of TIFF version 42. Recipe converts these to version 42 for compatibility with `github.com/google/tiff` library:

```go
// Little-endian DNG: IIRC (0x49 0x49 0x52 0x43) → II + version 42 (0x49 0x49 0x2A 0x00)
// Big-endian DNG:    MMCR (0x4D 0x4D 0x43 0x52) → MM + version 42 (0x4D 0x4D 0x00 0x2A)
```

## DCP Version Support

Recipe supports DCP versions v1.0-v1.6 (all Adobe Camera Raw/Lightroom compatible versions).

## Format Discovery Documentation

For complete details on how the binary DNG format was discovered during implementation, see:
- **[9-1-dcp-parser-FORMAT-PIVOT.md](../../../../docs/stories/9-1-dcp-parser-FORMAT-PIVOT.md)** - Complete format discovery documentation

This documents the critical discovery that real Adobe DCP files use binary TIFF tags (50700-52600), NOT the XML format initially assumed to be in tag 50740.

## References

- [Adobe DNG Specification 1.6](https://helpx.adobe.com/camera-raw/digital-negative.html) - Binary tag definitions (tags 50700-52600)
- [Adobe DNG SDK](https://www.adobe.com/support/downloads/dng/dng_sdk.html) - Reference implementation
- [TIFF 6.0 Specification](http://partners.adobe.com/public/developer/en/tiff/TIFF6.pdf) - TIFF file structure
- [github.com/google/tiff](https://github.com/google/tiff) - Go TIFF library used by Recipe
