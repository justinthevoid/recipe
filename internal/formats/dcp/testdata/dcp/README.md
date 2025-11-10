# DCP Test Samples

This directory contains DNG Camera Profile (.dcp) files used for testing Recipe's DCP parser.

## Sample Files

### Synthetic Samples (Generated for Testing)

- **`minimal-linear.dcp`** - Minimal DCP with linear tone curve (no adjustments)
  - Tone curve: (0,0) → (255,255) linear
  - Color matrices: Identity matrices
  - ProfileName: "Minimal Linear Test Profile"

- **`portrait-adjusted.dcp`** - Portrait-style DCP with tone adjustments
  - Tone curve: Slight S-curve for portrait look
  - Exposure: +0.3 stops
  - Contrast: +15
  - Highlights: -10 (recover highlights)
  - Shadows: +10 (lift shadows)

### Known Limitations

Due to DCP files being proprietary Adobe format:
- Synthetic DCPs may not render identically in Lightroom/ACR
- Full DCP validation requires Adobe software (Story 9-4)
- Color calibration matrices not tested (MVP limitation)

## DCP File Structure

```
DCP (.dcp) = TIFF Container
├── TIFF Header (II or MM)
├── Image File Directory (IFD)
│   ├── Tag 50740: CameraProfile (XML data)
│   └── Other standard TIFF tags
```

### Tag 50740 XML Structure

```xml
<crs:CameraProfile xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
  <crs:ProfileName>Profile Name</crs:ProfileName>
  <crs:ToneCurve>
    <rdf:Seq>
      <rdf:li>0, 0</rdf:li>
      <rdf:li>128, 128</rdf:li>
      <rdf:li>255, 255</rdf:li>
    </rdf:Seq>
  </crs:ToneCurve>
  <crs:ColorMatrix1>
    <rdf:Seq>
      <rdf:li>1.0 0.0 0.0</rdf:li>
      <rdf:li>0.0 1.0 0.0</rdf:li>
      <rdf:li>0.0 0.0 1.0</rdf:li>
    </rdf:Seq>
  </crs:ColorMatrix1>
</crs:CameraProfile>
```

## DCP Version Support

Recipe supports DCP versions v1.0-v1.6 (all Adobe Camera Raw/Lightroom compatible versions).

## References

- [Adobe DNG Specification 1.6](https://helpx.adobe.com/camera-raw/digital-negative.html)
- [Adobe DNG SDK](https://www.adobe.com/support/downloads/dng/dng_sdk.html)
- [TIFF 6.0 Specification](http://partners.adobe.com/public/developer/en/tiff/TIFF6.pdf)
