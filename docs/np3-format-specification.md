# NP3 File Format Specification
## Deep Research Analysis - 2025-11-07

### Executive Summary

The Nikon Picture Control (.np3) format uses a **hybrid structure** combining:
1. **TLV chunks** (Type-Length-Value) - Required by Nikon NX Studio for validation
2. **Raw byte heuristics** - Used by legacy parsers for parameter extraction

**Critical Finding**: Our generator (Story 1-3) only generates the raw byte data, **not the TLV chunks**, causing Nikon NX Studio to reject files as invalid despite round-trip tests passing.

---

## File Structure Overview

```
Offset  Size  Description
------  ----  -----------
0-2     3     Magic bytes: "NCP" (0x4E 0x43 0x50)
3-6     4     Version: 0x00 0x00 0x00 0x01 (v1.0.0.0, little-endian)
7-19    13    Reserved/Unknown: 0x00 0x00 0x00 0x00 0x04 "0310" 0x00 0x00 0x02 0x00
20-23   4     Name header: 0x00 0x00 0x00 0x14
24-43   20    Preset name (ASCII, null-terminated)
44-45   2     Padding: 0x00 0x00
46-335  290   TLV Chunks (29 chunks × 10 bytes each)
336+    ?     Extended data / padding
```

### File Size Patterns
- **Smallest**: 392 bytes (Life.np3, Porta.np3)
- **Common**: 480 bytes (sample.np3, Classic Chrome.np3)
- **Standard**: 978 bytes (Armitage series, Goldmatic.np3)
- **Large**: 1050 bytes (Filmic series)

---

## TLV Chunk Structure (Critical for Nikon NX Studio)

### Chunk Format (10 bytes per chunk)
```
Offset  Size  Type        Description
------  ----  ----        -----------
0       1     uint8       Chunk ID (0x03 to 0x1F)
1-3     3     padding     Always 0x00 0x00 0x00
4-5     2     uint16-BE   Length (typically 2, but 0x1F chunk has length 28)
6-7     2     varies      Value (2 bytes, encoding varies by chunk)
8-9     2     padding     Always 0x00 0x00 (overlaps next chunk start)
```

### Chunk Sequence
Chunks start at **offset 46** and appear in sequential order:
- Chunk 0x03 @ offset 46
- Chunk 0x04 @ offset 56
- Chunk 0x05 @ offset 66
- ... (each 10 bytes apart)
- Chunk 0x1F @ offset 326

**Total**: 29 chunks (0x03 through 0x1F) = 290 bytes

---

## Chunk ID Mapping

### Analysis Method
Compared 4 real .np3 files to identify constant vs. variable chunks:

| Chunk ID | Type | Observed Values | Interpretation |
|----------|------|-----------------|----------------|
| **0x03** | CONST | 32 (0x0020) | Format identifier |
| **0x04** | CONST | 0 (0x0000) | Reserved/unused |
| **0x05** | CONST | 65281 (0xff01) | Format flag |
| **0x06** | VAR | 33540-35332 | **Parameter** (possibly saturation-related) |
| **0x07** | VAR | 31748-33284 | **Parameter** (possibly contrast-related) |
| 0x08-0x0B | CONST | 65284 (0xff04) | Default values |
| 0x0C-0x0D | CONST | 65280 (0xff00) | Default values |
| 0x0E | CONST | 65284 (0xff04) | Default value |
| 0x0F-0x13 | CONST | 65281 (0xff01) | Default values |
| **0x14** | VAR | 32769-65281 | **Parameter** (brightness?) |
| 0x15 | CONST | 65290 (0xff0a) | Default value |
| **0x16** | VAR | 32772-33796 | **Parameter** |
| 0x17-0x18 | CONST | 65284 (0xff04) | Default values |
| **0x19-0x1E** | VAR | Wide range | **Parameters** (color channels?) |
| **0x1F** | VAR | 32399-40576 | **Extended data** (length=28) |

### Key Observations

**Constant Chunks (18 total)**:
- Chunks with 0xFF prefix (like 0xFF01, 0xFF04): Likely default/neutral values
- Same across all presets → format structure, not user data

**Variable Chunks (11 total)**:
- Chunks 0x06, 0x07, 0x14, 0x16, 0x19-0x1F
- These contain the actual preset parameters
- ArmitageEktar100 shows 0x0101 (257) for chunks 0x19-0x1D → possible neutral/zero encoding

---

## Heuristic Data Regions (Used by Legacy Parsers)

**Our parser ignores chunks** and uses these byte offsets:

### Raw Parameter Bytes (offsets 64-80)
- 16 bytes of signed data
- Values >128 treated as negative
- Parser: "Extract raw parameter bytes (offsets 64-80)"

### Color Data (offsets 100-300)
- RGB triplets (3 bytes each)
- Only triplets where R>10 OR G>10 OR B>10 are significant
- Parser analyzes count: `saturation = (count / 15) - 1`

### Tone Curve Data (offsets 150-500)
- Paired byte values
- Only non-zero pairs counted
- Parser analyzes count: `contrast = (count / 20) - 2`

**Overlap Region**: Offsets 150-300 counted by both color and tone curve analysis

---

## Current Generator Issues

### What Our Generator Does (generate.go)
```
✅ Writes magic bytes "NCP"
❌ Writes version 0x02 0x10 0x00 0x00 (WRONG - should be 0x00 0x00 0x00 0x01)
✅ Writes preset name at offset 20
✅ Writes raw bytes at offsets 64-79 (heuristic data)
✅ Writes color data at offsets 100-299
✅ Writes tone curve data starting at offset 150
❌ NO TLV chunks generated at all!
❌ File size: 500 bytes (should be 392-978+ bytes)
```

### Why Round-Trip Tests Passed
```
Parse(Generate(x)) == x ✅
```
Because both parser and generator use the same heuristic byte offsets!

### Why Nikon NX Studio Fails
```
Nikon NX Studio validates TLV chunk structure ❌
Our files have NO chunks → REJECTED
```

---

## Required Generator Fixes

### Priority 1: Fix Version Bytes
```diff
- data[3:7] = []byte{0x02, 0x10, 0x00, 0x00}
+ data[3:7] = []byte{0x00, 0x00, 0x00, 0x01}
```

### Priority 2: Generate TLV Chunks (offset 46-335)

**Constant chunks** (copy from real files):
```go
chunks := []chunkDef{
    {id: 0x03, length: 2, value: []byte{0x00, 0x20}},
    {id: 0x04, length: 2, value: []byte{0x00, 0x00}},
    {id: 0x05, length: 2, value: []byte{0xff, 0x01}},
    // ... (copy all constant chunks)
}
```

**Variable chunks** (encode parameters):
```go
// NEEDS RESEARCH: Determine encoding formulas
chunks = append(chunks, chunkDef{
    id: 0x06,
    length: 2,
    value: encodeParameter(params.saturation, ???)
})
// ... (map remaining variable chunks)
```

### Priority 3: Extend File to Proper Size
Current: 500 bytes
Target: Match input file size (392-978+ bytes)

### Priority 4: Keep Heuristic Data
**IMPORTANT**: Don't remove offsets 64-80, 100-300, 150-500!
These must remain so our parser continues to work.

---

## Next Steps

1. **Reverse-engineer variable chunk encoding**:
   - Analyze more sample files
   - Correlate chunk values with parsed parameters
   - Derive encoding formulas for chunks 0x06, 0x07, 0x14, 0x16, 0x19-0x1F

2. **Implement dual-layer generation**:
   - Layer 1: TLV chunks (for Nikon NX Studio validation)
   - Layer 2: Heuristic bytes (for our parser compatibility)

3. **Test against Nikon NX Studio**:
   - Generate file with proper chunks
   - Validate it loads without errors
   - Verify parameters display correctly

4. **Update round-trip tests**:
   - Add Nikon NX Studio validation step
   - Don't rely solely on Parse→Generate→Parse

---

## Research Tools Used

- **Binary analysis**: `od`, `hexdump`, Python struct module
- **Comparison**: Cross-file chunk analysis (4 samples)
- **Parser analysis**: Read parse.go implementation
- **Generator testing**: Round-trip conversion with size comparison

## Files Analyzed

1. Life.np3 (392 bytes)
2. sample.np3 (480 bytes)
3. Porta.np3 (392 bytes)
4. ArmitageEktar100.NP3 (978 bytes)

---

## Appendix: Raw Hex Examples

### Life.np3 Header + First 3 Chunks
```
Offset 0-47:
4e 43 50 00 00 00 01 00 00 00 00 04 30 33 31 30
00 00 02 00 00 00 00 14 46 69 6c 6d 73 74 69 6c
6c 73 4c 69 66 65 00 00 00 00 00 00 00 00

Chunks starting at 46:
03 00 00 00 00 02 00 20 00 00   Chunk 0x03: Value=0x0020 (32)
04 00 00 00 00 02 00 00 00 00   Chunk 0x04: Value=0x0000 (0)
05 00 00 00 00 02 ff 01 00 00   Chunk 0x05: Value=0xff01 (65281)
```

### Our Generator Output
```
Offset 0-63:
4e 43 50 02 10 00 00 00 00 00 00 00 00 00 00 00
00 00 00 00 43 6c 61 73 73 69 63 20 43 68 72 6f
6d 65 20 00 00 00 00 00 00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00

NO CHUNKS - just zeros at offset 44-63!
```

---

**Document Status**: DRAFT - Deep research complete, reverse-engineering in progress
**Last Updated**: 2025-11-07
**Next Review**: After encoding formula discovery
