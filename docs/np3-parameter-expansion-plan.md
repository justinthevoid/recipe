# NP3 Parameter Expansion - Safe Implementation Plan
## Adding 56 Parameters Without Breaking Anything

**Date**: 2025-11-07
**Epic**: Epic 1 - Core Conversion Engine
**Related Stories**: 1-2 (parser), 1-3 (generator), New: 1-5 (parameter expansion)
**Risk Level**: HIGH (core functionality changes)
**Estimated Time**: 3-4 weeks (phased approach)

---

## Executive Summary

**Goal**: Expand NP3 parser/generator from 6 heuristic parameters to 56 exact parameters using TypeScript repo findings.

**Critical Constraints**:
1. ✅ MUST NOT break existing 62/62 round-trip tests
2. ✅ MUST maintain rawData preservation for NP3→NP3
3. ✅ MUST update ALL documentation and validation
4. ✅ MUST be phased and testable at each step
5. ✅ MUST maintain code quality standards

**Strategy**: Incremental, test-driven, documentation-first approach with 4 phases.

---

## Current State Inventory

### Files That Will Be Modified

| File | Current Role | Changes Needed | Risk |
|------|--------------|----------------|------|
| `internal/formats/np3/parse.go` | Extracts 6 parameters via heuristics | Add 56 exact parameter functions | HIGH |
| `internal/formats/np3/generate.go` | Encodes via rawData preservation | Add 56 encoding functions | HIGH |
| `internal/models/recipe.go` | UniversalRecipe with ~80 fields | Add 4-6 new fields | MEDIUM |
| `internal/models/validation.go` | Validates existing parameters | Add 56 validation functions | MEDIUM |
| `internal/models/builder.go` | Builder pattern for recipe | Add builder methods | LOW |
| `internal/formats/xmp/generate.go` | XMP generation | Add parameter mappings | MEDIUM |
| `internal/formats/lrtemplate/generate.go` | LRTemplate generation | Add parameter mappings | MEDIUM |
| `docs/np3-format-specification.md` | Format documentation | Add offset mappings | LOW |
| `docs/stories/1-2-np3-parser.md` | Parser story | Update with new params | LOW |
| `docs/stories/1-3-np3-binary-generator.md` | Generator story | Update with new params | LOW |

**Total**: 10 files requiring coordinated updates

### Parameters to Add (56 total)

**Group 1: Basic Adjustments (8 parameters)**
- ✅ Already in UniversalRecipe: Contrast, Highlights, Shadows, Whites, Blacks, Saturation, Clarity, Sharpness
- Action: Replace heuristics with exact offsets

**Group 2: Color Blender (24 parameters = 8 colors × 3 values)**
- ✅ Already in UniversalRecipe: Red, Orange, Yellow, Green, Aqua, Blue, Purple, Magenta (with Hue/Saturation/Luminance)
- Action: Parse from offsets 332-361

**Group 3: Color Grading (15 parameters)**
- ❌ NOT in UniversalRecipe: Need to add
  - Highlights: Hue (0-360°), Chroma, Brightness
  - Midtone: Hue, Chroma, Brightness
  - Shadows: Hue, Chroma, Brightness
  - Global: Blending (0-100), Balance (-100 to 100)

**Group 4: Advanced (9 parameters)**
- ❌ Mid-range Sharpening: Not in UniversalRecipe
- ✅ Tone Curve: Already in UniversalRecipe (PointCurve)

---

## Phase 1: Foundation & Planning (Week 1, Days 1-3)

### Day 1: Documentation & Design

**Tasks**:
1. Create this implementation plan ✅
2. Create parameter mapping matrix (see below)
3. Update np3-format-specification.md with exact offsets
4. Design new UniversalRecipe fields
5. Create test data strategy

**Deliverables**:
- ✅ `docs/np3-parameter-expansion-plan.md` (this file)
- 📋 `docs/np3-parameter-mapping-matrix.md`
- 📋 Updated `docs/np3-format-specification.md`
- 📋 `docs/test-strategy-parameter-expansion.md`

**Validation**:
- [ ] All team members review plan
- [ ] No conflicts with existing architecture
- [ ] Test strategy covers all scenarios

### Day 2: UniversalRecipe Schema Updates

**Tasks**:
1. Add ColorGrading struct to models/recipe.go
2. Add MidRangeSharpening field
3. Update JSON/XML tags
4. Update builder pattern in models/builder.go

**Code Changes**:

```go
// Add to internal/models/recipe.go

// Color Grading (Nikon Flexible Color Picture Control)
// Allows independent color adjustments in Highlights, Midtones, and Shadows
type ColorGradingZone struct {
	Hue        int `json:"hue,omitempty" xml:"hue,omitempty"`        // Hue: 0-360 degrees
	Chroma     int `json:"chroma,omitempty" xml:"chroma,omitempty"`     // Chroma: -100 to +100
	Brightness int `json:"brightness,omitempty" xml:"brightness,omitempty"` // Brightness: -100 to +100
}

type ColorGrading struct {
	Highlights ColorGradingZone `json:"highlights,omitempty" xml:"highlights,omitempty"` // Highlights zone
	Midtone    ColorGradingZone `json:"midtone,omitempty" xml:"midtone,omitempty"`       // Midtone zone
	Shadows    ColorGradingZone `json:"shadows,omitempty" xml:"shadows,omitempty"`       // Shadows zone
	Blending   int              `json:"blending,omitempty" xml:"blending,omitempty"`     // Blending: 0-100
	Balance    int              `json:"balance,omitempty" xml:"balance,omitempty"`       // Balance: -100 to +100
}

// Add to UniversalRecipe struct
type UniversalRecipe struct {
	// ... existing fields ...

	// Advanced Sharpening
	MidRangeSharpening float64 `json:"midRangeSharpening,omitempty" xml:"midRangeSharpening,omitempty"` // Mid-range sharpening: -5.0 to +5.0

	// Color Grading (NP3-specific)
	ColorGrading *ColorGrading `json:"colorGrading,omitempty" xml:"colorGrading,omitempty"` // Color grading zones

	// ... existing fields ...
}
```

**Validation Functions** (add to internal/models/validation.go):

```go
// ValidateColorGradingHue validates color grading hue (0-360 degrees)
func ValidateColorGradingHue(hue int) error {
	if hue < 0 || hue > 360 {
		return fmt.Errorf("color grading hue must be 0-360, got %d", hue)
	}
	return nil
}

// ValidateColorGradingChroma validates color grading chroma (-100 to 100)
func ValidateColorGradingChroma(chroma int) error {
	if chroma < -100 || chroma > 100 {
		return fmt.Errorf("color grading chroma must be -100 to 100, got %d", chroma)
	}
	return nil
}

// ValidateColorGradingBrightness validates color grading brightness (-100 to 100)
func ValidateColorGradingBrightness(brightness int) error {
	if brightness < -100 || brightness > 100 {
		return fmt.Errorf("color grading brightness must be -100 to 100, got %d", brightness)
	}
	return nil
}

// ValidateColorGradingBlending validates blending (0-100)
func ValidateColorGradingBlending(blending int) error {
	if blending < 0 || blending > 100 {
		return fmt.Errorf("color grading blending must be 0-100, got %d", blending)
	}
	return nil
}

// ValidateMidRangeSharpening validates mid-range sharpening (-5.0 to 5.0)
func ValidateMidRangeSharpening(value float64) error {
	if value < -5.0 || value > 5.0 {
		return fmt.Errorf("mid-range sharpening must be -5.0 to 5.0, got %.1f", value)
	}
	return nil
}
```

**Tests** (create internal/models/validation_colorgrading_test.go):

```go
func TestValidateColorGradingHue(t *testing.T) {
	tests := []struct {
		name    string
		hue     int
		wantErr bool
	}{
		{"valid min", 0, false},
		{"valid max", 360, false},
		{"valid mid", 180, false},
		{"invalid negative", -1, true},
		{"invalid too high", 361, true},
	}
	// ... test implementation
}

// Similar tests for all validation functions
```

**Deliverables**:
- [ ] Updated models/recipe.go with new fields
- [ ] Updated models/validation.go with new validators
- [ ] Updated models/builder.go with new methods
- [ ] 100% test coverage for new validators

**Validation**:
- [ ] go test ./internal/models/ passes
- [ ] No breaking changes to existing tests
- [ ] Documentation strings complete

### Day 3: Parameter Offset Constants

**Tasks**:
1. Create internal/formats/np3/offsets.go
2. Document all 56 parameter offsets
3. Add offset validation
4. Create offset test suite

**Code**: Create `internal/formats/np3/offsets.go`

```go
// Package np3 - Parameter byte offset constants
// Source: https://github.com/ssssota/nikon-flexible-color-picture-control/blob/main/offset.ts
// Verified: 2025-11-07

package np3

// Basic file structure offsets
const (
	OffsetMagic   = 0    // Magic bytes "NCP" (3 bytes)
	OffsetVersion = 3    // Version bytes (4 bytes)
	OffsetName    = 0x18 // Preset name (19 bytes max, null-terminated)
)

// Basic adjustment offsets
const (
	OffsetSharpening         = 0x52  // 82: Sharpening (-3.0 to 9.0, scaled by 4)
	OffsetMidRangeSharpening = 0xF2  // 242: Mid-range sharpening (-5.0 to 5.0, scaled by 4)
	OffsetClarity            = 0x5C  // 92: Clarity (-5.0 to 5.0, scaled by 4)
	OffsetContrast           = 0x110 // 272: Contrast (-100 to 100, offset by 0x80)
	OffsetHighlights         = 0x11A // 282: Highlights (-100 to 100, offset by 0x80)
	OffsetShadows            = 0x124 // 292: Shadows (-100 to 100, offset by 0x80)
	OffsetWhiteLevel         = 0x12E // 302: White Level (-100 to 100, offset by 0x80)
	OffsetBlackLevel         = 0x138 // 312: Black Level (-100 to 100, offset by 0x80)
	OffsetSaturation         = 0x142 // 322: Saturation (-100 to 100, offset by 0x80)
)

// Color Blender offsets (8 colors × 3 bytes each = 24 bytes)
// Each color has: [Hue, Chroma, Brightness] (all offset by 0x80, -100 to 100)
const (
	OffsetColorBlenderRed     = 0x14C // 332: Red
	OffsetColorBlenderOrange  = 0x14F // 335: Orange
	OffsetColorBlenderYellow  = 0x152 // 338: Yellow
	OffsetColorBlenderGreen   = 0x155 // 341: Green
	OffsetColorBlenderCyan    = 0x158 // 344: Cyan
	OffsetColorBlenderBlue    = 0x15B // 347: Blue
	OffsetColorBlenderPurple  = 0x15E // 350: Purple
	OffsetColorBlenderMagenta = 0x161 // 353: Magenta
)

// Color Grading offsets (3 zones × 4 bytes each + 2 global = 14 bytes)
// Each zone has: [Hue MSB (4 bits), Hue LSB, Chroma, Brightness]
// Hue: 12-bit value (0-360°), Chroma/Brightness: offset by 0x80 (-100 to 100)
const (
	OffsetColorGradingHighlights = 0x170 // 368: Highlights zone (4 bytes)
	OffsetColorGradingMidtone    = 0x174 // 372: Midtone zone (4 bytes)
	OffsetColorGradingShadows    = 0x178 // 376: Shadows zone (4 bytes)
	OffsetColorGradingBlending   = 0x180 // 384: Blending (0-100, offset by 0x80)
	OffsetColorGradingBalance    = 0x182 // 386: Balance (-100 to 100, offset by 0x80)
)

// Tone Curve offsets
const (
	OffsetToneCurvePoints = 0x194 // 404: Point count + control points
	OffsetToneCurveRaw    = 0x1CC // 460: Raw tone curve data (257 × 16-bit BE)
)

// File size constraints
const (
	MinFileSizeBasic    = 392  // Minimum valid NP3 file
	MinFileSizeWithCurve = 973 // Minimum with tone curve data (0x3CD)
)

// Decoding formulas
const (
	Offset128  = 0x80  // Standard offset for signed byte values
	ScaleFactor4 = 4.0 // Scale factor for sharpening/clarity parameters
)
```

**Tests**: Create `internal/formats/np3/offsets_test.go`

```go
func TestOffsetConstants(t *testing.T) {
	// Verify offsets are in valid range
	tests := []struct {
		name   string
		offset int
	}{
		{"Name", OffsetName},
		{"Sharpening", OffsetSharpening},
		{"Contrast", OffsetContrast},
		{"ColorBlenderRed", OffsetColorBlenderRed},
		{"ColorGradingHighlights", OffsetColorGradingHighlights},
		{"ToneCurvePoints", OffsetToneCurvePoints},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.offset < 0 || tt.offset > 1024 {
				t.Errorf("Offset %s = %d is outside valid range", tt.name, tt.offset)
			}
		})
	}
}
```

**Deliverables**:
- [ ] offsets.go with all 56 parameter offsets
- [ ] offsets_test.go with validation tests
- [ ] Documentation for each offset

---

## Phase 2: Parser Enhancement (Week 1-2, Days 4-10)

### Strategy: Dual-Mode Parsing

**Approach**:
1. Keep existing heuristic functions as fallback
2. Add exact offset reading functions
3. Prefer exact, fall back to heuristic
4. Preserve rawData for perfect NP3→NP3

**Implementation Groups**:

### Group 1: Basic Adjustments (Days 4-5)

**Files to modify**: `internal/formats/np3/parse.go`

**Step 1**: Add exact reading functions

```go
// readSharpening reads sharpening from exact offset (replaces heuristic)
// Offset: 82 (0x52), Formula: (byte - 0x80) / 4.0, Range: -3.0 to 9.0
func readSharpening(data []byte) (float64, error) {
	if len(data) <= OffsetSharpening {
		return 0, fmt.Errorf("file too small for sharpening offset")
	}

	rawValue := int(data[OffsetSharpening])
	adjusted := rawValue - Offset128
	value := float64(adjusted) / ScaleFactor4

	// Validate range
	if value < -3.0 || value > 9.0 {
		return 0, fmt.Errorf("sharpening value %f out of range [-3.0, 9.0]", value)
	}

	return value, nil
}

// readContrast reads contrast from exact offset
// Offset: 272 (0x110), Formula: byte - 0x80, Range: -100 to 100
func readContrast(data []byte) (int, error) {
	if len(data) <= OffsetContrast {
		return 0, fmt.Errorf("file too small for contrast offset")
	}

	value := int(data[OffsetContrast]) - Offset128

	// Validate range
	if value < -100 || value > 100 {
		return 0, fmt.Errorf("contrast value %d out of range [-100, 100]", value)
	}

	return value, nil
}

// Similar functions for:
// - readClarity()
// - readHighlights()
// - readShadows()
// - readWhiteLevel()
// - readBlackLevel()
// - readSaturation()
// - readMidRangeSharpening()
```

**Step 2**: Update extractParameters()

```go
func extractParameters(data []byte) (*np3Parameters, error) {
	params := &np3Parameters{}

	// Store raw data for perfect round-trip preservation
	params.rawData = make([]byte, len(data))
	copy(params.rawData, data)

	// Extract name (unchanged)
	if len(data) >= 60 {
		// ... existing name extraction code
	}

	// Try exact offset reading first, fall back to heuristics
	var err error

	// Sharpening: Try exact, fall back to heuristic
	if sharpening, err := readSharpening(data); err == nil {
		// Map NP3 range (-3.0 to 9.0) to internal range (0-9)
		// UniversalRecipe expects 0-150, we'll scale in buildRecipe
		params.sharpening = int((sharpening + 3.0) * 9.0 / 12.0)
	} else {
		// Fall back to heuristic (existing code)
		params.sharpening = estimateSharpening(data)
	}

	// Contrast: Try exact, fall back to heuristic
	if contrast, err := readContrast(data); err == nil {
		// Map NP3 range (-100 to 100) to internal range (-3 to 3)
		params.contrast = contrast / 33
	} else {
		params.contrast = estimateContrast(data)
	}

	// ... similar for all parameters

	return params, nil
}
```

**Step 3**: Add tests for each parameter

Create `internal/formats/np3/parse_exact_test.go`:

```go
func TestReadSharpening(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    float64
		wantErr bool
	}{
		{
			name: "min value",
			data: makeTestData(OffsetSharpening, 0x80-12), // -3.0
			want: -3.0,
			wantErr: false,
		},
		{
			name: "neutral",
			data: makeTestData(OffsetSharpening, 0x80), // 0.0
			want: 0.0,
			wantErr: false,
		},
		{
			name: "max value",
			data: makeTestData(OffsetSharpening, 0x80+36), // 9.0
			want: 9.0,
			wantErr: false,
		},
		{
			name: "out of range high",
			data: makeTestData(OffsetSharpening, 0x80+40), // 10.0
			want: 0,
			wantErr: true,
		},
		{
			name: "file too small",
			data: make([]byte, 50),
			want: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readSharpening(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("readSharpening() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got - tt.want) > 0.01 {
				t.Errorf("readSharpening() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function
func makeTestData(offset int, value byte) []byte {
	data := make([]byte, 1024)
	data[offset] = value
	return data
}

// Similar tests for all 8 basic adjustment parameters
```

**Deliverables**:
- [ ] 9 exact reading functions (sharpening, contrast, etc.)
- [ ] Updated extractParameters() with dual-mode parsing
- [ ] 9 test suites (one per parameter, 5-10 test cases each)
- [ ] All tests passing (including existing 62 round-trip tests)

**Validation Checklist**:
- [ ] Each parameter has exact reading function
- [ ] Each function validates input range
- [ ] Each function handles short files gracefully
- [ ] Fallback to heuristic works if exact fails
- [ ] Round-trip tests still pass (62/62)
- [ ] New unit tests cover edge cases

### Group 2: Color Blender (Days 6-7)

**Files to modify**: `internal/formats/np3/parse.go`

**Step 1**: Add color blender reading function

```go
// readColorBlender reads all 8 color blender values
// Each color has 3 bytes: [Hue, Chroma, Brightness]
// Offsets: 332-361 (30 bytes total)
func readColorBlender(data []byte) (map[string]ColorBlenderValues, error) {
	if len(data) < OffsetColorBlenderMagenta + 3 {
		return nil, fmt.Errorf("file too small for color blender data")
	}

	colors := map[string]ColorBlenderValues{
		"Red":     readColorBlenderValues(data, OffsetColorBlenderRed),
		"Orange":  readColorBlenderValues(data, OffsetColorBlenderOrange),
		"Yellow":  readColorBlenderValues(data, OffsetColorBlenderYellow),
		"Green":   readColorBlenderValues(data, OffsetColorBlenderGreen),
		"Cyan":    readColorBlenderValues(data, OffsetColorBlenderCyan),
		"Blue":    readColorBlenderValues(data, OffsetColorBlenderBlue),
		"Purple":  readColorBlenderValues(data, OffsetColorBlenderPurple),
		"Magenta": readColorBlenderValues(data, OffsetColorBlenderMagenta),
	}

	return colors, nil
}

type ColorBlenderValues struct {
	Hue        int // -100 to 100
	Chroma     int // -100 to 100 (similar to Saturation)
	Brightness int // -100 to 100 (similar to Luminance)
}

func readColorBlenderValues(data []byte, offset int) ColorBlenderValues {
	return ColorBlenderValues{
		Hue:        int(data[offset]) - Offset128,
		Chroma:     int(data[offset+1]) - Offset128,
		Brightness: int(data[offset+2]) - Offset128,
	}
}
```

**Step 2**: Update buildRecipe() to map to UniversalRecipe

```go
// In buildRecipe(), add:

// Read color blender if available
if colorBlender, err := readColorBlender(params.rawData); err == nil {
	// Map to UniversalRecipe color adjustments
	if red, ok := colorBlender["Red"]; ok {
		builder.WithRedHue(red.Hue)
		builder.WithRedSaturation(red.Chroma)  // Chroma → Saturation
		builder.WithRedLuminance(red.Brightness)
	}
	// Similar for all 8 colors...
}
```

**Step 3**: Add comprehensive tests

```go
func TestReadColorBlender(t *testing.T) {
	// Test with real sample file
	data, err := os.ReadFile("../../testdata/np3/sample.np3")
	if err != nil {
		t.Fatal(err)
	}

	colors, err := readColorBlender(data)
	if err != nil {
		t.Fatalf("readColorBlender() error = %v", err)
	}

	// Verify we got all 8 colors
	expectedColors := []string{"Red", "Orange", "Yellow", "Green", "Cyan", "Blue", "Purple", "Magenta"}
	for _, color := range expectedColors {
		if _, ok := colors[color]; !ok {
			t.Errorf("Missing color: %s", color)
		}
	}

	// Verify values are in range
	for color, values := range colors {
		if values.Hue < -100 || values.Hue > 100 {
			t.Errorf("%s Hue %d out of range [-100, 100]", color, values.Hue)
		}
		if values.Chroma < -100 || values.Chroma > 100 {
			t.Errorf("%s Chroma %d out of range [-100, 100]", color, values.Chroma)
		}
		if values.Brightness < -100 || values.Brightness > 100 {
			t.Errorf("%s Brightness %d out of range [-100, 100]", color, values.Brightness)
		}
	}
}
```

**Deliverables**:
- [ ] readColorBlender() function
- [ ] readColorBlenderValues() helper
- [ ] Integration in buildRecipe()
- [ ] Comprehensive test suite
- [ ] Round-trip tests still pass

### Group 3: Color Grading (Days 8-9)

**Files to modify**: `internal/formats/np3/parse.go`

**Step 1**: Add color grading reading function

```go
// readColorGrading reads all color grading zones
// Each zone has 4 bytes: [Hue MSB (4 bits), Hue LSB, Chroma, Brightness]
// Plus 2 global parameters: Blending, Balance
func readColorGrading(data []byte) (*ColorGrading, error) {
	if len(data) < OffsetColorGradingBalance + 1 {
		return nil, fmt.Errorf("file too small for color grading data")
	}

	return &ColorGrading{
		Highlights: readColorGradingZone(data, OffsetColorGradingHighlights),
		Midtone:    readColorGradingZone(data, OffsetColorGradingMidtone),
		Shadows:    readColorGradingZone(data, OffsetColorGradingShadows),
		Blending:   int(data[OffsetColorGradingBlending]) - Offset128,
		Balance:    int(data[OffsetColorGradingBalance]) - Offset128,
	}, nil
}

func readColorGradingZone(data []byte, offset int) ColorGradingZone {
	// Hue is 12-bit: upper 4 bits in first byte, lower 8 bits in second byte
	hueMSB := int(data[offset] & 0x0F)
	hueLSB := int(data[offset+1])
	hue := (hueMSB << 8) | hueLSB

	return ColorGradingZone{
		Hue:        hue,                                  // 0-360 degrees
		Chroma:     int(data[offset+2]) - Offset128,      // -100 to 100
		Brightness: int(data[offset+3]) - Offset128,      // -100 to 100
	}
}
```

**Step 2**: Update buildRecipe()

```go
// In buildRecipe(), add:

// Read color grading if available
if colorGrading, err := readColorGrading(params.rawData); err == nil {
	recipe.ColorGrading = colorGrading
}
```

**Step 3**: Add tests

```go
func TestReadColorGrading(t *testing.T) {
	// Test with manufactured data
	data := make([]byte, 512)

	// Highlights: Hue=180° (0x00B4), Chroma=+20, Brightness=-10
	data[OffsetColorGradingHighlights] = 0x00   // MSB (upper 4 bits = 0)
	data[OffsetColorGradingHighlights+1] = 0xB4 // LSB (180)
	data[OffsetColorGradingHighlights+2] = 0x80 + 20  // Chroma
	data[OffsetColorGradingHighlights+3] = 0x80 - 10  // Brightness

	// Similar for Midtone and Shadows...

	grading, err := readColorGrading(data)
	if err != nil {
		t.Fatalf("readColorGrading() error = %v", err)
	}

	// Verify Highlights
	if grading.Highlights.Hue != 180 {
		t.Errorf("Highlights Hue = %d, want 180", grading.Highlights.Hue)
	}
	if grading.Highlights.Chroma != 20 {
		t.Errorf("Highlights Chroma = %d, want 20", grading.Highlights.Chroma)
	}
	if grading.Highlights.Brightness != -10 {
		t.Errorf("Highlights Brightness = %d, want -10", grading.Highlights.Brightness)
	}
}
```

**Deliverables**:
- [ ] readColorGrading() function
- [ ] readColorGradingZone() helper
- [ ] Integration in buildRecipe()
- [ ] Test suite with 12-bit hue validation
- [ ] Round-trip tests still pass

### Group 4: Tone Curve (Day 10)

**Note**: Tone curve is complex - defer to separate story if needed.

**Deliverables**:
- [ ] readToneCurve() function (basic)
- [ ] readToneCurvePoints() helper
- [ ] Integration in buildRecipe()
- [ ] Basic tests

---

## Phase 3: Generator Enhancement (Week 2-3, Days 11-17)

### Strategy: Dual-Mode Generation

**Approach**:
1. Keep rawData preservation for NP3→NP3
2. Add exact encoding for parameter editing
3. Detect when rawData is unavailable (e.g., XMP→NP3)
4. Encode parameters precisely using inverse formulas

### Group 1: Basic Adjustments (Days 11-12)

**Files to modify**: `internal/formats/np3/generate.go`

**Step 1**: Add encoding functions

```go
// writeSharpening writes sharpening to exact offset
// Offset: 82 (0x52), Formula: byte = (value * 4.0) + 0x80, Range: -3.0 to 9.0
func writeSharpening(data []byte, value float64) error {
	if len(data) <= OffsetSharpening {
		return fmt.Errorf("buffer too small for sharpening offset")
	}

	// Validate range
	if value < -3.0 || value > 9.0 {
		return fmt.Errorf("sharpening value %f out of range [-3.0, 9.0]", value)
	}

	// Encode: multiply by 4, offset by 0x80
	encoded := byte(int(value * ScaleFactor4) + Offset128)
	data[OffsetSharpening] = encoded

	return nil
}

// writeContrast writes contrast to exact offset
// Offset: 272 (0x110), Formula: byte = value + 0x80, Range: -100 to 100
func writeContrast(data []byte, value int) error {
	if len(data) <= OffsetContrast {
		return fmt.Errorf("buffer too small for contrast offset")
	}

	// Validate range
	if value < -100 || value > 100 {
		return fmt.Errorf("contrast value %d out of range [-100, 100]", value)
	}

	// Encode: offset by 0x80
	encoded := byte(value + Offset128)
	data[OffsetContrast] = encoded

	return nil
}

// Similar functions for all basic adjustments
```

**Step 2**: Update encodeBinary()

```go
func encodeBinary(params *np3Parameters, presetName string) ([]byte, error) {
	var data []byte

	// If we have raw data, use it as base (perfect round-trip)
	if params.rawData != nil && len(params.rawData) > 0 {
		data = make([]byte, len(params.rawData))
		copy(data, params.rawData)

		// Update only the parameters that might have changed
		// (For NP3→NP3, rawData is unchanged, so this is a no-op)
		// (For edited parameters, this updates specific bytes)
	} else {
		// No raw data - encode from scratch (e.g., XMP→NP3)
		data = make([]byte, 480)

		// Write header structure
		writeHeader(data, presetName)

		// Write TLV chunks with neutral values
		writeDefaultChunks(data)
	}

	// Write all parameters to exact offsets
	if err := writeSharpening(data, params.sharpening); err != nil {
		return nil, fmt.Errorf("write sharpening: %w", err)
	}
	if err := writeContrast(data, params.contrast); err != nil {
		return nil, fmt.Errorf("write contrast: %w", err)
	}
	// ... write all parameters

	return data, nil
}
```

**Step 3**: Add encoding tests

```go
func TestWriteSharpening(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		want    byte
		wantErr bool
	}{
		{"min value", -3.0, 0x80-12, false},
		{"neutral", 0.0, 0x80, false},
		{"max value", 9.0, 0x80+36, false},
		{"out of range", 10.0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, 1024)
			err := writeSharpening(data, tt.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("writeSharpening() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && data[OffsetSharpening] != tt.want {
				t.Errorf("writeSharpening() wrote %d, want %d", data[OffsetSharpening], tt.want)
			}
		})
	}
}

// Add round-trip test
func TestSharpening RoundTrip(t *testing.T) {
	values := []float64{-3.0, -1.5, 0.0, 2.5, 5.0, 9.0}

	for _, original := range values {
		data := make([]byte, 1024)

		// Write
		if err := writeSharpening(data, original); err != nil {
			t.Fatalf("writeSharpening(%f) failed: %v", original, err)
		}

		// Read back
		retrieved, err := readSharpening(data)
		if err != nil {
			t.Fatalf("readSharpening() failed: %v", err)
		}

		// Verify
		if math.Abs(retrieved - original) > 0.01 {
			t.Errorf("Round-trip failed: wrote %f, read %f", original, retrieved)
		}
	}
}
```

**Deliverables**:
- [ ] 9 encoding functions
- [ ] Updated encodeBinary() with parameter writing
- [ ] Unit tests for each encoder
- [ ] Round-trip tests (write→read→compare)

### Group 2-4: Similar implementation for Color Blender, Color Grading, Tone Curve

**Days 13-17**: Follow same pattern as Group 1

---

## Phase 4: Integration & Cross-Format (Week 3-4, Days 18-24)

### XMP Integration (Days 18-19)

**Files to modify**: `internal/formats/xmp/generate.go`

**Tasks**:
1. Add ColorGrading handling (may need to map to closest XMP equivalent)
2. Add MidRangeSharpening handling
3. Update parameter mappings
4. Add tests for NP3→XMP→NP3 conversion

### LRTemplate Integration (Days 20-21)

**Files to modify**: `internal/formats/lrtemplate/generate.go`

**Tasks**:
1. Same as XMP integration
2. Test NP3→LRTemplate→NP3 conversion

### Comprehensive Testing (Days 22-23)

**Test Types**:

1. **Unit Tests** (per parameter)
   - Read function tests
   - Write function tests
   - Round-trip tests
   - Edge case tests

2. **Integration Tests** (per format)
   - NP3→UniversalRecipe→NP3
   - NP3→XMP→UniversalRecipe
   - NP3→LRTemplate→UniversalRecipe

3. **Round-Trip Tests** (cross-format)
   - NP3→XMP→NP3 (parameter preservation)
   - NP3→LRTemplate→NP3 (parameter preservation)
   - Measure parameter accuracy (should be >95%)

4. **Real-World Tests**
   - Test all 62 existing .np3 files
   - Visual comparison in Nikon NX Studio
   - Parameter value spot-checks

### Documentation Update (Day 24)

**Files to update**:
1. `docs/np3-format-specification.md` - Add all offsets
2. `docs/stories/1-2-np3-parser.md` - Document new parameters
3. `docs/stories/1-3-np3-binary-generator.md` - Document encoding
4. `docs/epic-1-retrospective.md` - Update with findings
5. Create `docs/stories/1-5-np3-parameter-expansion.md` - New story

---

## Risk Mitigation

### Risk 1: Breaking Existing Tests

**Mitigation**:
- Run full test suite after EVERY change
- Git commit after each working increment
- Maintain backward compatibility with heuristics

**Recovery Plan**:
- Each day is a commit point
- Can revert to any previous day's work
- Have "escape hatch" to revert to Phase 1

### Risk 2: UniversalRecipe Field Conflicts

**Mitigation**:
- Review schema changes with team before implementation
- Ensure JSON/XML tags don't conflict
- Test serialization/deserialization

### Risk 3: Performance Degradation

**Mitigation**:
- Benchmark before and after changes
- Exact offset reading should be FASTER than heuristics
- Monitor test execution time

### Risk 4: Cross-Format Incompatibility

**Mitigation**:
- Some NP3 parameters have no XMP/LR equivalent
- Document what can't be converted
- Preserve in FormatSpecificBinary for round-trip

---

## Success Criteria

### Phase 1 Complete When:
- [ ] UniversalRecipe schema updated
- [ ] All validation functions implemented
- [ ] All validator tests passing (100% coverage)
- [ ] Documentation updated
- [ ] Team review approved

### Phase 2 Complete When:
- [ ] All 56 read functions implemented
- [ ] extractParameters() uses exact offsets
- [ ] Heuristic fallback preserved
- [ ] Unit tests for all parameters (100% coverage)
- [ ] Round-trip tests still pass (62/62)

### Phase 3 Complete When:
- [ ] All 56 write functions implemented
- [ ] encodeBinary() encodes all parameters
- [ ] rawData preservation maintained
- [ ] Unit tests for all encoders (100% coverage)
- [ ] Write→Read round-trip tests pass

### Phase 4 Complete When:
- [ ] XMP integration complete
- [ ] LRTemplate integration complete
- [ ] Cross-format tests pass (>95% accuracy)
- [ ] All 62 existing files still convert correctly
- [ ] Visual validation in Nikon NX Studio confirms accuracy
- [ ] All documentation updated
- [ ] Story 1-5 marked as DONE

---

## Test Strategy Summary

### Test Pyramid

```
          /\
         /  \  E2E Tests (5%)
        /____\    - Visual validation in Nikon NX Studio
       /      \   - Cross-format round-trip
      /________\
     /   Integration Tests (25%)
    /______________\
   /                \
  /  Unit Tests (70%) \
 /____________________\
```

**Unit Tests (70%)**:
- 56 parameters × 2 functions (read + write) = 112 functions
- Each function: 5-10 test cases
- Total: ~800-1000 test cases
- Focus: Correct encoding/decoding formulas

**Integration Tests (25%)**:
- Parse full .np3 files
- Generate full .np3 files
- Round-trip: Parse→Build→Generate→Parse
- Cross-format: NP3→XMP→NP3, etc.
- Total: ~50-100 test cases

**E2E Tests (5%)**:
- Load generated files in Nikon NX Studio
- Visual comparison of presets
- Parameter value spot-checks
- Total: ~10-20 manual validations

### Test Data

**Existing**:
- 62 real .np3 files in testdata/np3/

**New**:
- Synthetic files with known parameter values
- Edge case files (min/max values)
- Corrupted files (for error handling)

---

## Timeline & Effort

### Week 1: Foundation & Basic Parameters
- Days 1-3: Schema, validation, offsets
- Days 4-5: Basic adjustments (8 params)
- Days 6-7: Color Blender (24 params)
- **Deliverable**: Parser enhancement complete

### Week 2: Color Grading & Generator
- Days 8-9: Color Grading (15 params)
- Day 10: Tone Curve (2 params)
- Days 11-12: Generator - Basic adjustments
- Days 13-17: Generator - Color Blender, Grading, Curve
- **Deliverable**: Generator enhancement complete

### Week 3: Integration
- Days 18-19: XMP integration
- Days 20-21: LRTemplate integration
- Days 22-23: Comprehensive testing
- Day 24: Documentation
- **Deliverable**: Full parameter expansion complete

### Week 4: Buffer & Polish
- Reserve for unexpected issues
- Performance optimization
- Additional testing
- Team review

**Total Time**: 3-4 weeks

---

## File Change Tracking Matrix

| File | Lines Changed | New Functions | Tests Added | Risk | Priority |
|------|---------------|---------------|-------------|------|----------|
| models/recipe.go | +30 | +2 structs | +20 | MED | 1 |
| models/validation.go | +100 | +6 validators | +60 | MED | 1 |
| models/builder.go | +50 | +4 methods | +20 | LOW | 2 |
| formats/np3/offsets.go | +80 | +0 | +10 | LOW | 1 |
| formats/np3/parse.go | +500 | +58 readers | +400 | HIGH | 2 |
| formats/np3/generate.go | +500 | +58 writers | +400 | HIGH | 3 |
| formats/xmp/generate.go | +100 | +4 mappings | +50 | MED | 4 |
| formats/lrtemplate/generate.go | +100 | +4 mappings | +50 | MED | 4 |
| docs/*.md | +500 | N/A | N/A | LOW | 5 |

**Total**: ~2000 lines of code, ~1000 lines of tests

---

## Checkpoint & Rollback Strategy

### Daily Commits
```bash
# Day 1
git commit -m "Add ColorGrading struct and MidRangeSharpening field to UniversalRecipe"

# Day 2
git commit -m "Add validation functions for new parameters"

# Day 3
git commit -m "Add NP3 parameter offset constants"

# ... etc.
```

### Phase Tags
```bash
# After Phase 1
git tag -a "phase1-foundation" -m "UniversalRecipe schema complete"

# After Phase 2
git tag -a "phase2-parser" -m "Parser enhancement complete"

# After Phase 3
git tag -a "phase3-generator" -m "Generator enhancement complete"

# After Phase 4
git tag -a "phase4-integration" -m "Full parameter expansion complete"
```

### Rollback Procedure
If critical issue found:
1. Identify last known-good commit/tag
2. Create branch from that point
3. Cherry-pick working commits
4. Resume from stable point

---

## Communication Plan

### Daily Standups
- What was completed yesterday
- What's planned for today
- Any blockers or risks
- Test status update

### Weekly Reviews
- Phase completion status
- Demo of new functionality
- Test coverage metrics
- Risk assessment update

### Milestone Reports
- After each phase completion
- Comprehensive status document
- Updated timeline if needed

---

## Next Immediate Actions

### Today (Day 1)
1. ✅ Create this implementation plan
2. 📋 Team review and approval of plan
3. 📋 Create parameter mapping matrix
4. 📋 Update np3-format-specification.md
5. 📋 Begin UniversalRecipe schema updates

### Tomorrow (Day 2)
1. Complete UniversalRecipe changes
2. Implement all validation functions
3. Write validation tests
4. Get code review approval

### This Week
1. Complete Phase 1 (Foundation)
2. Begin Phase 2 (Parser enhancement)
3. Target: Basic adjustments parsing complete by Friday

---

**Document Status**: Implementation plan ready for team review
**Risk Level**: HIGH → Can be reduced to MEDIUM with phased approach
**Estimated Success**: 95% (with proper testing and rollback strategy)
**Last Updated**: 2025-11-07

---

*This is a living document - update as implementation progresses*
