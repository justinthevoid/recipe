# Capture One .costyle Sample Files

This directory contains sample .costyle files for testing the Capture One parser.

## Sample Files

### sample1-portrait.costyle
Complete .costyle file with all core adjustments:
- Exposure: +0.7
- Contrast: +15
- Saturation: +10
- Temperature: +5
- Tint: -3
- Clarity: +20
- Color balance adjustments for shadows, midtones, highlights

Use case: Portrait preset with warm tones and enhanced clarity.

### sample2-minimal.costyle
Minimal valid .costyle with only required XML structure and a single adjustment:
- Exposure: -0.3

Use case: Testing minimal file parsing and default/zero value handling.

### sample3-landscape.costyle
Landscape preset with strong color adjustments:
- Exposure: -0.2
- Contrast: +25
- Saturation: +30
- Clarity: +15
- Color balance: Enhanced blues and greens

Use case: Testing wide range of parameter values and color balance mapping.

## File Format

Capture One .costyle files use Adobe XMP-style XML structure:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:RDF>
    <rdf:Description>
      <!-- Adjustment parameters here -->
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>
```

## Parameter Ranges

- Exposure: -2.0 to +2.0
- Contrast: -100 to +100
- Saturation: -100 to +100
- Temperature: -100 to +100 (relative scale)
- Tint: -100 to +100
- Clarity: -100 to +100
- Color Balance Hue: 0 to 360 (degrees)
- Color Balance Saturation: -100 to +100

## Version Compatibility

These samples are compatible with Capture One versions 2023-2025.
No significant format changes have been observed between versions.

## Round-Trip Testing

Recipe verifies .costyle round-trip conversion with 95%+ accuracy (Story 8-4):

### Running Round-Trip Tests

```bash
# Run all costyle round-trip tests
go test -v ./internal/formats/costyle -run TestRoundTrip

# Run with coverage
go test -cover ./internal/formats/costyle -run TestRoundTrip

# Run specific edge case tests
go test -v ./internal/formats/costyle -run TestRoundTrip_EdgeCases
```

### Test Results

**Average Accuracy:** 98.37% (exceeds 95% requirement)
- **Min Accuracy:** 97.56%
- **Max Accuracy:** 100.00%
- **Test Coverage:** 85.9%

See [test-results.md](test-results.md) for detailed accuracy metrics and parameter breakdown.

### Round-Trip Flow

```
Original .costyle → Parse() → UniversalRecipe1
                              ↓
                      Generate() → New .costyle
                              ↓
                      Parse() → UniversalRecipe2
                              ↓
                Compare(UR1, UR2) → Accuracy %
```

**Tolerance Thresholds:**
- Exposure: ±0.01 stops
- Integers (Contrast, Saturation, etc.): ±1 value
- Temperature: ±2 Kelvin
- Split Toning: ±1° hue, ±1 saturation

## Sources

All sample files are synthetically generated based on the .costyle format specification
reverse-engineered from real Capture One exports.

**Real-World Samples (TODO):**
- Need ≥5 real .costyle presets from Etsy/marketplaces (AC-1, AC-4)
- Target styles: Portrait, Landscape, Product, Black & White, Vintage/Film
- See Story 8-4 Task 5 for acquisition plan
