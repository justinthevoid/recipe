package np3

// NP3 Binary File Structure Byte Offsets
// Based on research from: https://github.com/ssssota/nikon-flexible-color-picture-control
//
// This file defines exact byte offsets for all parameters in Nikon's NP3 (Picture Control) format.
// These offsets enable precise parameter extraction and encoding, replacing heuristic-based analysis.
//
// File Structure Overview:
// - Bytes 0-19: Header and metadata
// - Bytes 20-59: Picture Control name (40 bytes, but only 19 usable for alphanumeric)
// - Bytes 60-90: Basic adjustments and sharpening
// - Bytes 240-330: Advanced adjustments (mid-range sharpening, contrast, highlights, etc.)
// - Bytes 332-361: Color Blender (8 colors × 3 values)
// - Bytes 368-386: Color Grading (3 zones × 4 values + 2 global)
// - Bytes 404+: Tone Curve data (control points and raw curve)

// Header and Metadata Offsets
const (
	// OffsetName is the start of the Picture Control name field (40 bytes total, 19 chars usable).
	// Name must be 1-19 alphanumeric characters.
	OffsetName = 0x18 // 24 decimal
)

// Basic Adjustments Offsets
const (
	// OffsetSharpening is the sharpening amount parameter.
	// Formula: (byte - 0x80) / 4.0
	// Range: -3.0 to 9.0
	// Maps to: UniversalRecipe.Sharpness (scaled to 0-150)
	// Note: This is stored in chunk 0x06 value bytes
	OffsetSharpening = 0x52 // 82 decimal

	// OffsetClarity is the clarity parameter.
	// Formula: (byte - 0x80) / 4.0
	// Range: -5.0 to 5.0
	// Maps to: UniversalRecipe.Clarity (scaled to -100 to +100)
	// Note: This is stored in chunk 0x07 value bytes
	OffsetClarity = 0x5C // 92 decimal
)

// Advanced Adjustments Offsets
const (
	// OffsetMidRangeSharpening is the mid-range sharpening parameter (NP3-specific).
	// Formula: (byte - 0x80) / 4.0
	// Range: -5.0 to 5.0
	// Maps to: UniversalRecipe.MidRangeSharpening
	OffsetMidRangeSharpening = 0xF2 // 242 decimal

	// OffsetContrast is the contrast adjustment parameter.
	// Formula: byte - 0x80
	// Range: -100 to 100
	// Maps to: UniversalRecipe.Contrast
	OffsetContrast = 0x110 // 272 decimal

	// OffsetHighlights is the highlights adjustment parameter.
	// Formula: byte - 0x80
	// Range: -100 to 100
	// Maps to: UniversalRecipe.Highlights
	OffsetHighlights = 0x11A // 282 decimal

	// OffsetShadows is the shadows adjustment parameter.
	// Formula: byte - 0x80
	// Range: -100 to 100
	// Maps to: UniversalRecipe.Shadows
	OffsetShadows = 0x124 // 292 decimal

	// OffsetWhiteLevel is the white level adjustment parameter.
	// Formula: byte - 0x80
	// Range: -100 to 100
	// Maps to: UniversalRecipe.Whites
	OffsetWhiteLevel = 0x12E // 302 decimal

	// OffsetBlackLevel is the black level adjustment parameter.
	// Formula: byte - 0x80
	// Range: -100 to 100
	// Maps to: UniversalRecipe.Blacks
	OffsetBlackLevel = 0x138 // 312 decimal

	// OffsetSaturation is the saturation adjustment parameter.
	// Formula: byte - 0x80
	// Range: -100 to 100
	// Maps to: UniversalRecipe.Saturation
	OffsetSaturation = 0x142 // 322 decimal
)

// Color Blender Offsets (8 colors × 3 values = 24 bytes)
// Each color has 3 sequential bytes: Hue, Chroma, Brightness
// Formula for all: byte - 0x80
// Range: -100 to 100
const (
	// Red color adjustments (3 bytes starting at 332)
	OffsetRedHue        = 0x14C // 332 decimal - Maps to: UniversalRecipe.Red.Hue
	OffsetRedChroma     = 0x14D // 333 decimal - Maps to: UniversalRecipe.Red.Saturation (Chroma ≈ Saturation)
	OffsetRedBrightness = 0x14E // 334 decimal - Maps to: UniversalRecipe.Red.Luminance

	// Orange color adjustments (3 bytes starting at 335)
	OffsetOrangeHue        = 0x14F // 335 decimal - Maps to: UniversalRecipe.Orange.Hue
	OffsetOrangeChroma     = 0x150 // 336 decimal - Maps to: UniversalRecipe.Orange.Saturation
	OffsetOrangeBrightness = 0x151 // 337 decimal - Maps to: UniversalRecipe.Orange.Luminance

	// Yellow color adjustments (3 bytes starting at 338)
	OffsetYellowHue        = 0x152 // 338 decimal - Maps to: UniversalRecipe.Yellow.Hue
	OffsetYellowChroma     = 0x153 // 339 decimal - Maps to: UniversalRecipe.Yellow.Saturation
	OffsetYellowBrightness = 0x154 // 340 decimal - Maps to: UniversalRecipe.Yellow.Luminance

	// Green color adjustments (3 bytes starting at 341)
	OffsetGreenHue        = 0x155 // 341 decimal - Maps to: UniversalRecipe.Green.Hue
	OffsetGreenChroma     = 0x156 // 342 decimal - Maps to: UniversalRecipe.Green.Saturation
	OffsetGreenBrightness = 0x157 // 343 decimal - Maps to: UniversalRecipe.Green.Luminance

	// Cyan color adjustments (3 bytes starting at 344)
	// Note: NP3 uses "Cyan" but UniversalRecipe uses "Aqua"
	OffsetCyanHue        = 0x158 // 344 decimal - Maps to: UniversalRecipe.Aqua.Hue
	OffsetCyanChroma     = 0x159 // 345 decimal - Maps to: UniversalRecipe.Aqua.Saturation
	OffsetCyanBrightness = 0x15A // 346 decimal - Maps to: UniversalRecipe.Aqua.Luminance

	// Blue color adjustments (3 bytes starting at 347)
	OffsetBlueHue        = 0x15B // 347 decimal - Maps to: UniversalRecipe.Blue.Hue
	OffsetBlueChroma     = 0x15C // 348 decimal - Maps to: UniversalRecipe.Blue.Saturation
	OffsetBlueBrightness = 0x15D // 349 decimal - Maps to: UniversalRecipe.Blue.Luminance

	// Purple color adjustments (3 bytes starting at 350)
	OffsetPurpleHue        = 0x15E // 350 decimal - Maps to: UniversalRecipe.Purple.Hue
	OffsetPurpleChroma     = 0x15F // 351 decimal - Maps to: UniversalRecipe.Purple.Saturation
	OffsetPurpleBrightness = 0x160 // 352 decimal - Maps to: UniversalRecipe.Purple.Luminance

	// Magenta color adjustments (3 bytes starting at 353)
	OffsetMagentaHue        = 0x161 // 353 decimal - Maps to: UniversalRecipe.Magenta.Hue
	OffsetMagentaChroma     = 0x162 // 354 decimal - Maps to: UniversalRecipe.Magenta.Saturation
	OffsetMagentaBrightness = 0x163 // 355 decimal - Maps to: UniversalRecipe.Magenta.Luminance
)

// Color Grading Offsets (Flexible Color Picture Control - NP3-specific)
// 3 tonal zones (Highlights, Midtone, Shadows) × 4 bytes each + 2 global parameters
// Each zone: 2 bytes for Hue (12-bit), 1 byte for Chroma, 1 byte for Brightness
const (
	// Highlights zone (4 bytes starting at 368)
	// Hue formula: ((byte[0] & 0x0f) << 8) + byte[1]
	// Chroma/Brightness formula: byte - 0x80
	OffsetHighlightsHue        = 0x170 // 368 decimal - 2 bytes, 12-bit value (0-360 degrees)
	OffsetHighlightsChroma     = 0x172 // 370 decimal - 1 byte (-100 to 100)
	OffsetHighlightsBrightness = 0x173 // 371 decimal - 1 byte (-100 to 100)

	// Midtone zone (4 bytes starting at 372)
	OffsetMidtoneHue        = 0x174 // 372 decimal - 2 bytes, 12-bit value (0-360 degrees)
	OffsetMidtoneChroma     = 0x176 // 374 decimal - 1 byte (-100 to 100)
	OffsetMidtoneBrightness = 0x177 // 375 decimal - 1 byte (-100 to 100)

	// Shadows zone (4 bytes starting at 376)
	OffsetShadowsHue        = 0x178 // 376 decimal - 2 bytes, 12-bit value (0-360 degrees)
	OffsetShadowsChroma     = 0x17A // 378 decimal - 1 byte (-100 to 100)
	OffsetShadowsBrightness = 0x17B // 379 decimal - 1 byte (-100 to 100)

	// Global Color Grading parameters
	// OffsetColorGradingBlending controls transition smoothness between zones.
	// Formula: byte - 0x80 (Signed8 encoding, then clamped to 0-100)
	// Range: 0 to 100 (negative values after decoding are clamped to 0)
	OffsetColorGradingBlending = 0x180 // 384 decimal

	// OffsetColorGradingBalance shifts overall color balance.
	// Formula: byte - 0x80
	// Range: -100 to 100
	OffsetColorGradingBalance = 0x182 // 386 decimal
)

// Tone Curve Offsets
const (
	// OffsetToneCurvePointCount is the number of control points in the tone curve.
	// Range: 0 to 255
	OffsetToneCurvePointCount = 0x194 // 404 decimal

	// OffsetToneCurvePoints is the start of the control point data.
	// Each point is 2 bytes: x (input 0-255), y (output 0-255)
	// Total points = value at OffsetToneCurvePointCount
	OffsetToneCurvePoints = 0x195 // 405 decimal

	// OffsetToneCurveRaw is the start of the full 257-point tone curve.
	// Format: 257 × 16-bit big-endian values
	// Range: 0 to 32767 for each point
	// Maps input levels (0-256) to output levels
	OffsetToneCurveRaw = 0x1CC // 460 decimal
)

// File Size Constants
const (
	// MinNP3FileSize is the minimum valid NP3 file size (480 bytes).
	MinNP3FileSize = 480

	// ExpectedNP3FileSize is the standard NP3 file size.
	ExpectedNP3FileSize = 480

	// NameFieldSize is the total size of the name field in bytes.
	NameFieldSize = 40

	// MaxNameLength is the maximum length of the Picture Control name (alphanumeric only).
	MaxNameLength = 19
)

// Offset Validation Constants
const (
	// MaxValidOffset is the highest offset we should read from (tone curve raw data end).
	// 460 (start) + (257 points × 2 bytes) = 974, but file is only 480 bytes.
	// In practice, tone curve raw data is truncated or compressed in the 480-byte format.
	MaxValidOffset = ExpectedNP3FileSize - 1 // 479
)

// Parameter Encoding/Decoding Constants
const (
	// BiasValue is the offset used for signed parameter encoding (0x80 = 128).
	// Most NP3 parameters use: actualValue = byte - BiasValue
	BiasValue = 0x80

	// ScaleFactor4 is the division factor for fractional parameters.
	// Used for sharpening, mid-range sharpening, clarity: actualValue = (byte - 0x80) / 4.0
	ScaleFactor4 = 4.0
)

// Decoding Helper Functions

// DecodeSigned8 decodes a signed byte parameter (range: -128 to +127).
// Formula: byte - 0x80
// Example: 0x00 = -128, 0x80 = 0, 0xFF = 127
func DecodeSigned8(b byte) int {
	return int(b) - BiasValue
}

// EncodeSigned8 encodes a signed integer to a byte parameter.
// Formula: value + 0x80
// Clamps to 0-255 range for safety.
func EncodeSigned8(value int) byte {
	result := value + BiasValue
	if result < 0 {
		return 0
	}
	if result > 255 {
		return 255
	}
	return byte(result)
}

// DecodeScaled4 decodes a scaled fractional parameter.
// Formula: (byte - 0x80) / 4.0
// Used for: Sharpening, Clarity, Mid-range Sharpening
func DecodeScaled4(b byte) float64 {
	return float64(int(b)-BiasValue) / ScaleFactor4
}

// EncodeScaled4 encodes a fractional value to a scaled byte parameter.
// Formula: (value * 4.0) + 0x80
// Clamps to 0-255 range for safety.
func EncodeScaled4(value float64) byte {
	result := int(value*ScaleFactor4) + BiasValue
	if result < 0 {
		return 0
	}
	if result > 255 {
		return 255
	}
	return byte(result)
}

// DecodeHue12 decodes a 12-bit hue value from 2 bytes.
// Formula: ((byte[0] & 0x0F) << 8) + byte[1]
// Range: 0 to 360 degrees
// Used for: Color Grading zone hues
// Matches reference implementation: hue = ((b1 & 0x0F) << 8) + b2
// The high byte has 0x80 bias added during encoding, so we mask it off
func DecodeHue12(b1, b2 byte) int {
	return (int(b1&0x0F) << 8) | int(b2)
}

// EncodeHue12 encodes a hue value (0-360) to 2 bytes.
// Returns: (high byte, low byte)
// Formula: high byte = 0x80 + (hue >> 8), low byte = hue & 0xFF
// This matches the reference implementation's encoding.
func EncodeHue12(hue int) (byte, byte) {
	// Clamp to 0-360 range
	if hue < 0 {
		hue = 0
	}
	if hue > 360 {
		hue = 360
	}
	// Split into 12 bits with 0x80 bias on high byte
	b1 := byte(0x80 + (hue >> 8))
	b2 := byte(hue & 0xFF)
	return b1, b2
}

// ValidateOffset checks if an offset is within the valid NP3 file range.
func ValidateOffset(offset int) bool {
	return offset >= 0 && offset <= MaxValidOffset
}

// ValidateOffsetRange checks if a byte range [start, end) is valid.
func ValidateOffsetRange(start, end int) bool {
	return start >= 0 && end <= ExpectedNP3FileSize && start < end
}
