package np3

import (
	"testing"
)

// Test offset constant values match TypeScript implementation
func TestOffsetConstants(t *testing.T) {
	tests := []struct {
		name     string
		offset   int
		expected int
	}{
		// Header and Metadata
		{"Name", OffsetName, 24},

		// Basic Adjustments
		{"Sharpening", OffsetSharpening, 82},
		{"Clarity", OffsetClarity, 92},

		// Advanced Adjustments
		{"MidRangeSharpening", OffsetMidRangeSharpening, 242},
		{"Contrast", OffsetContrast, 272},
		{"Highlights", OffsetHighlights, 282},
		{"Shadows", OffsetShadows, 292},
		{"WhiteLevel", OffsetWhiteLevel, 302},
		{"BlackLevel", OffsetBlackLevel, 312},
		{"Saturation", OffsetSaturation, 322},

		// Color Blender - Red
		{"RedHue", OffsetRedHue, 332},
		{"RedChroma", OffsetRedChroma, 333},
		{"RedBrightness", OffsetRedBrightness, 334},

		// Color Blender - Orange
		{"OrangeHue", OffsetOrangeHue, 335},
		{"OrangeChroma", OffsetOrangeChroma, 336},
		{"OrangeBrightness", OffsetOrangeBrightness, 337},

		// Color Blender - Yellow
		{"YellowHue", OffsetYellowHue, 338},
		{"YellowChroma", OffsetYellowChroma, 339},
		{"YellowBrightness", OffsetYellowBrightness, 340},

		// Color Blender - Green
		{"GreenHue", OffsetGreenHue, 341},
		{"GreenChroma", OffsetGreenChroma, 342},
		{"GreenBrightness", OffsetGreenBrightness, 343},

		// Color Blender - Cyan
		{"CyanHue", OffsetCyanHue, 344},
		{"CyanChroma", OffsetCyanChroma, 345},
		{"CyanBrightness", OffsetCyanBrightness, 346},

		// Color Blender - Blue
		{"BlueHue", OffsetBlueHue, 347},
		{"BlueChroma", OffsetBlueChroma, 348},
		{"BlueBrightness", OffsetBlueBrightness, 349},

		// Color Blender - Purple
		{"PurpleHue", OffsetPurpleHue, 350},
		{"PurpleChroma", OffsetPurpleChroma, 351},
		{"PurpleBrightness", OffsetPurpleBrightness, 352},

		// Color Blender - Magenta
		{"MagentaHue", OffsetMagentaHue, 353},
		{"MagentaChroma", OffsetMagentaChroma, 354},
		{"MagentaBrightness", OffsetMagentaBrightness, 355},

		// Color Grading - Shadows (FIRST in binary order at offset 368)
		{"ShadowsHue", OffsetShadowsHue, 368},
		{"ShadowsChroma", OffsetShadowsChroma, 370},
		{"ShadowsBrightness", OffsetShadowsBrightness, 371},

		// Color Grading - Midtone (SECOND in binary order at offset 372)
		{"MidtoneHue", OffsetMidtoneHue, 372},
		{"MidtoneChroma", OffsetMidtoneChroma, 374},
		{"MidtoneBrightness", OffsetMidtoneBrightness, 375},

		// Color Grading - Highlights (THIRD in binary order at offset 376)
		{"HighlightsHue", OffsetHighlightsHue, 376},
		{"HighlightsChroma", OffsetHighlightsChroma, 378},
		{"HighlightsBrightness", OffsetHighlightsBrightness, 379},

		// Color Grading - Global
		{"ColorGradingBlending", OffsetColorGradingBlending, 384},
		{"ColorGradingBalance", OffsetColorGradingBalance, 386},

		// Tone Curve
		{"ToneCurvePointCount", OffsetToneCurvePointCount, 404},
		{"ToneCurvePoints", OffsetToneCurvePoints, 405},
		{"ToneCurveRaw", OffsetToneCurveRaw, 460},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.offset != tt.expected {
				t.Errorf("Offset%s = %d (0x%X), want %d (0x%X)",
					tt.name, tt.offset, tt.offset, tt.expected, tt.expected)
			}
		})
	}
}

// Test DecodeSigned8 function
func TestDecodeSigned8(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected int
	}{
		{"zero", 0x80, 0},
		{"positive max", 0xFF, 127},
		{"negative max", 0x00, -128},
		{"positive mid", 0xC0, 64},
		{"negative mid", 0x40, -64},
		{"positive 100", 0xE4, 100},
		{"negative 100", 0x1C, -100},
		{"positive 1", 0x81, 1},
		{"negative 1", 0x7F, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeSigned8(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeSigned8(0x%02X) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// Test EncodeSigned8 function
func TestEncodeSigned8(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected byte
	}{
		{"zero", 0, 0x80},
		{"positive max", 127, 0xFF},
		{"negative max", -128, 0x00},
		{"positive mid", 64, 0xC0},
		{"negative mid", -64, 0x40},
		{"positive 100", 100, 0xE4},
		{"negative 100", -100, 0x1C},
		{"positive 1", 1, 0x81},
		{"negative 1", -1, 0x7F},
		{"clamp high", 200, 0xFF}, // Should clamp to 255
		{"clamp low", -200, 0x00}, // Should clamp to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeSigned8(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeSigned8(%d) = 0x%02X, want 0x%02X", tt.input, result, tt.expected)
			}
		})
	}
}

// Test round-trip encoding/decoding for signed values
func TestSigned8RoundTrip(t *testing.T) {
	for value := -100; value <= 100; value++ {
		encoded := EncodeSigned8(value)
		decoded := DecodeSigned8(encoded)
		if decoded != value {
			t.Errorf("Round-trip failed for %d: encoded to 0x%02X, decoded to %d",
				value, encoded, decoded)
		}
	}
}

// Test DecodeScaled4 function
func TestDecodeScaled4(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected float64
	}{
		{"zero", 0x80, 0.0},
		{"positive max (sharpening)", 0xA4, 9.0},  // (164-128)/4 = 9.0
		{"negative max (sharpening)", 0x74, -3.0}, // (116-128)/4 = -3.0
		{"positive max (clarity)", 0x94, 5.0},     // (148-128)/4 = 5.0
		{"negative max (clarity)", 0x6C, -5.0},    // (108-128)/4 = -5.0
		{"positive mid", 0x8A, 2.5},               // (138-128)/4 = 2.5
		{"negative mid", 0x76, -2.5},              // (118-128)/4 = -2.5
		{"positive quarter", 0x81, 0.25},          // (129-128)/4 = 0.25
		{"negative quarter", 0x7F, -0.25},         // (127-128)/4 = -0.25
		{"mid-range sharp positive", 0x8A, 2.5},   // (138-128)/4 = 2.5
		{"mid-range sharp negative", 0x76, -2.5},  // (118-128)/4 = -2.5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeScaled4(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeScaled4(0x%02X) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

// Test EncodeScaled4 function
func TestEncodeScaled4(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected byte
	}{
		{"zero", 0.0, 0x80},
		{"positive max (sharpening)", 9.0, 0xA4},
		{"negative max (sharpening)", -3.0, 0x74},
		{"positive max (clarity)", 5.0, 0x94},
		{"negative max (clarity)", -5.0, 0x6C},
		{"positive mid", 2.5, 0x8A},
		{"negative mid", -2.5, 0x76},
		{"positive quarter", 0.25, 0x81},
		{"negative quarter", -0.25, 0x7F},
		{"clamp high", 50.0, 0xFF}, // Should clamp to 255
		{"clamp low", -50.0, 0x00}, // Should clamp to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeScaled4(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeScaled4(%f) = 0x%02X, want 0x%02X", tt.input, result, tt.expected)
			}
		})
	}
}

// Test round-trip encoding/decoding for scaled values
func TestScaled4RoundTrip(t *testing.T) {
	// Test in 0.25 increments from -5.0 to +5.0
	for i := -20; i <= 20; i++ {
		value := float64(i) * 0.25
		encoded := EncodeScaled4(value)
		decoded := DecodeScaled4(encoded)
		if decoded != value {
			t.Errorf("Round-trip failed for %f: encoded to 0x%02X, decoded to %f",
				value, encoded, decoded)
		}
	}
}

// Test DecodeHue12 function
func TestDecodeHue12(t *testing.T) {
	tests := []struct {
		name     string
		b1, b2   byte
		expected int
	}{
		{"zero", 0x00, 0x00, 0},
		{"max", 0x01, 0x68, 360}, // (1 << 8) + 104 = 360
		{"mid", 0x00, 0xB4, 180}, // (0 << 8) + 180 = 180
		{"90", 0x00, 0x5A, 90},   // (0 << 8) + 90 = 90
		{"270", 0x01, 0x0E, 270}, // (1 << 8) + 14 = 270
		{"1", 0x00, 0x01, 1},
		{"359", 0x01, 0x67, 359},                // (1 << 8) + 103 = 359
		{"high nibble masked", 0xF1, 0x68, 360}, // Top 4 bits ignored
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeHue12(tt.b1, tt.b2)
			if result != tt.expected {
				t.Errorf("DecodeHue12(0x%02X, 0x%02X) = %d, want %d",
					tt.b1, tt.b2, result, tt.expected)
			}
		})
	}
}

// Test EncodeHue12 function
func TestEncodeHue12(t *testing.T) {
	tests := []struct {
		name       string
		input      int
		expectedB1 byte
		expectedB2 byte
	}{
		{"zero", 0, 0x80, 0x00},         // 0x80 bias on high byte
		{"max", 360, 0x81, 0x68},        // 0x80 | 0x01 = 0x81
		{"mid", 180, 0x80, 0xB4},        // 0x80 | 0x00 = 0x80
		{"90", 90, 0x80, 0x5A},          // 0x80 | 0x00 = 0x80
		{"270", 270, 0x81, 0x0E},        // 0x80 | 0x01 = 0x81
		{"1", 1, 0x80, 0x01},            // 0x80 | 0x00 = 0x80
		{"359", 359, 0x81, 0x67},        // 0x80 | 0x01 = 0x81
		{"clamp high", 500, 0x81, 0x68}, // Should clamp to 360
		{"clamp low", -10, 0x80, 0x00},  // Should clamp to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b1, b2 := EncodeHue12(tt.input)
			if b1 != tt.expectedB1 || b2 != tt.expectedB2 {
				t.Errorf("EncodeHue12(%d) = (0x%02X, 0x%02X), want (0x%02X, 0x%02X)",
					tt.input, b1, b2, tt.expectedB1, tt.expectedB2)
			}
		})
	}
}

// Test round-trip encoding/decoding for 12-bit hue values
func TestHue12RoundTrip(t *testing.T) {
	for hue := 0; hue <= 360; hue++ {
		b1, b2 := EncodeHue12(hue)
		decoded := DecodeHue12(b1, b2)
		if decoded != hue {
			t.Errorf("Round-trip failed for %d: encoded to (0x%02X, 0x%02X), decoded to %d",
				hue, b1, b2, decoded)
		}
	}
}

// Test ValidateOffset function
func TestValidateOffset(t *testing.T) {
	tests := []struct {
		name     string
		offset   int
		expected bool
	}{
		{"zero", 0, true},
		{"valid mid", 240, true},
		{"max valid", 479, true},
		{"at file size", 480, false},
		{"beyond file size", 500, false},
		{"negative", -1, false},
		{"name offset", 24, true},
		{"contrast offset", 272, true},
		{"tone curve offset", 460, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateOffset(tt.offset)
			if result != tt.expected {
				t.Errorf("ValidateOffset(%d) = %v, want %v", tt.offset, result, tt.expected)
			}
		})
	}
}

// Test ValidateOffsetRange function
func TestValidateOffsetRange(t *testing.T) {
	tests := []struct {
		name       string
		start, end int
		expected   bool
	}{
		{"valid range", 0, 100, true},
		{"entire file", 0, 480, true},
		{"name field", 24, 64, true},
		{"invalid start negative", -1, 100, false},
		{"invalid end beyond file", 0, 481, false},
		{"invalid start >= end", 100, 100, false},
		{"invalid start > end", 100, 50, false},
		{"single byte", 100, 101, true},
		{"color grading zone", 368, 372, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateOffsetRange(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("ValidateOffsetRange(%d, %d) = %v, want %v",
					tt.start, tt.end, result, tt.expected)
			}
		})
	}
}

// Test file size constants
func TestFileSizeConstants(t *testing.T) {
	if MinNP3FileSize != 480 {
		t.Errorf("MinNP3FileSize = %d, want 480", MinNP3FileSize)
	}
	if ExpectedNP3FileSize != 480 {
		t.Errorf("ExpectedNP3FileSize = %d, want 480", ExpectedNP3FileSize)
	}
	if NameFieldSize != 40 {
		t.Errorf("NameFieldSize = %d, want 40", NameFieldSize)
	}
	if MaxNameLength != 19 {
		t.Errorf("MaxNameLength = %d, want 19", MaxNameLength)
	}
	if MaxValidOffset != 479 {
		t.Errorf("MaxValidOffset = %d, want 479", MaxValidOffset)
	}
}

// Test encoding constants
func TestEncodingConstants(t *testing.T) {
	if BiasValue != 0x80 {
		t.Errorf("BiasValue = 0x%02X, want 0x80", BiasValue)
	}
	if ScaleFactor4 != 4.0 {
		t.Errorf("ScaleFactor4 = %f, want 4.0", ScaleFactor4)
	}
}

// Test that color blender offsets are sequential
func TestColorBlenderSequential(t *testing.T) {
	colorOffsets := []int{
		OffsetRedHue, OffsetRedChroma, OffsetRedBrightness,
		OffsetOrangeHue, OffsetOrangeChroma, OffsetOrangeBrightness,
		OffsetYellowHue, OffsetYellowChroma, OffsetYellowBrightness,
		OffsetGreenHue, OffsetGreenChroma, OffsetGreenBrightness,
		OffsetCyanHue, OffsetCyanChroma, OffsetCyanBrightness,
		OffsetBlueHue, OffsetBlueChroma, OffsetBlueBrightness,
		OffsetPurpleHue, OffsetPurpleChroma, OffsetPurpleBrightness,
		OffsetMagentaHue, OffsetMagentaChroma, OffsetMagentaBrightness,
	}

	// Check each offset is exactly 1 byte after the previous
	for i := 1; i < len(colorOffsets); i++ {
		if colorOffsets[i] != colorOffsets[i-1]+1 {
			t.Errorf("Color offset %d is not sequential: %d (expected %d)",
				i, colorOffsets[i], colorOffsets[i-1]+1)
		}
	}

	// Check the sequence starts at the correct offset
	if OffsetRedHue != 332 {
		t.Errorf("OffsetRedHue = %d, want 332", OffsetRedHue)
	}

	// Check the sequence ends at the correct offset
	if OffsetMagentaBrightness != 355 {
		t.Errorf("OffsetMagentaBrightness = %d, want 355", OffsetMagentaBrightness)
	}
}

// Test that all defined offsets are within valid file bounds
func TestAllOffsetsValid(t *testing.T) {
	offsets := []struct {
		name   string
		offset int
	}{
		{"Name", OffsetName},
		{"Sharpening", OffsetSharpening},
		{"Clarity", OffsetClarity},
		{"MidRangeSharpening", OffsetMidRangeSharpening},
		{"Contrast", OffsetContrast},
		{"Highlights", OffsetHighlights},
		{"Shadows", OffsetShadows},
		{"WhiteLevel", OffsetWhiteLevel},
		{"BlackLevel", OffsetBlackLevel},
		{"Saturation", OffsetSaturation},
		{"RedHue", OffsetRedHue},
		{"MagentaBrightness", OffsetMagentaBrightness},
		{"HighlightsHue", OffsetHighlightsHue},
		{"ShadowsBrightness", OffsetShadowsBrightness},
		{"ColorGradingBlending", OffsetColorGradingBlending},
		{"ColorGradingBalance", OffsetColorGradingBalance},
		{"ToneCurvePointCount", OffsetToneCurvePointCount},
	}

	for _, tt := range offsets {
		t.Run(tt.name, func(t *testing.T) {
			if !ValidateOffset(tt.offset) {
				t.Errorf("Offset%s (%d) is outside valid file bounds", tt.name, tt.offset)
			}
		})
	}
}
