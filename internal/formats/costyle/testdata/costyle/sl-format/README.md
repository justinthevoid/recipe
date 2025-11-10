# Capture One Style Library Format (.costyle)

## Discovery

These .costyle files use a different XML format than the Adobe XMP-style format currently supported by Recipe.

**Root Element:** `<SL>` (Style Library) instead of `<xmpmeta>`
**Format Version:** Engine="1100" (Capture One version identifier)

## Example Structure

```xml
<?xml version="1.0" encoding="utf-8"?>
<SL Engine="1100">
 <E K="ICCProfile" V="agfa-optima-100-100%.icc"/>
 <E K="Name" V="Agfa Optima 100 100%"/>
 <E K="StyleSource" V="Styles"/>
 <E K="UUID" V="{C2238A8B-1793-4AEF-828C-03CAC9F68D46}"/>
 <E K="HighlightRecovery" V="50.0"/>
 <E K="Midtone" V="0;0;0;0.0500000007450581"/>
</SL>
```

## Parameters

The SL format uses key-value pairs (`<E K="..." V="..."/>`) for adjustments:
- **ICCProfile**: References external .icc color profile
- **Name**: Preset display name
- **StyleSource**: Source category (e.g., "Styles")
- **UUID**: Unique identifier for the preset
- **HighlightRecovery**: Highlight tone adjustment
- **Midtone**: Midtone curve adjustment (semicolon-separated values)

## Status

**NOT SUPPORTED** in current Recipe implementation (Story 8-4).

The current .costyle parser (`internal/formats/costyle/parse.go`) only supports Adobe XMP-style .costyle files with `<xmpmeta>` root element. Supporting the SL format would require:

1. Separate parser for `<SL>` XML structure
2. Mapping SL parameters to UniversalRecipe (different parameter model)
3. Understanding SL parameter ranges and semantics
4. Generating SL-format .costyle files from UniversalRecipe

## Future Work

**Epic 8 Scope:** XMP-style .costyle format only
**Future Epic:** Add support for SL-format .costyle files

This would be a new feature addition, not part of the round-trip testing story.

## Source

These 78 film emulation presets were provided by the user and appear to be from a Capture One preset pack:
- Agfa Optima 100/200 variants (13 files)
- Agfa Vista 100 (1 file)
- Fuji Natura 1600 variants (2 files)
- Fuji Pro 160NS variants (4 files)
- Kodak Ektar 100 variants (11 files)
- Kodak Gold 200 variants (6 files)
- Kodak Portra 160/400/800 variants (30 files)
- Kodak Ultramax 400 variants (3 files)
- Rollei Digibase CN 200 variants (3 files)
- ...and more

Each preset includes an accompanying .icc color profile file.

## Recommendations

For Story 8-4 (Round-Trip Testing):
- Use the existing synthetic XMP-style samples (sample1-portrait, sample2-minimal, sample3-landscape)
- Document SL format as unsupported in known-conversion-limitations.md
- Create issue for future SL format support

For users with SL-format .costyle files:
- These files cannot be converted by Recipe at this time
- Recipe only supports XMP-style .costyle exports from Capture One
- To use Recipe: Export presets from Capture One as XMP-style .costyle (if available)
