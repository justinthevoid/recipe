package dcp

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/justin/recipe/internal/models"
)

// analyzeToneCurve analyzes a piecewise linear tone curve to extract
// exposure, contrast, highlights, and shadows adjustments.
//
// Algorithm (adapted for 0.0-1.0 normalized values):
//   - Exposure: Midpoint (input=0.5) vertical shift from linear
//   - Contrast: Slope difference (top point 0.75 - bottom point 0.25)
//   - Highlights: Top-end curve shape (point 1.0 deviation)
//   - Shadows: Bottom-end curve shape (point 0.0 deviation)
//
// All values are normalized to UniversalRecipe ranges:
//   - Exposure: -2.0 to +2.0
//   - Contrast, Highlights, Shadows: -1.0 to +1.0
//
// Returns zero values if curve is empty or linear.
func analyzeToneCurve(points []ToneCurvePoint) (exposure, contrast, highlights, shadows float64) {
	if len(points) == 0 {
		return 0, 0, 0, 0 // Linear curve
	}

	// Find key points (exact match or interpolate)
	midpoint := findPoint(points, 0.5)
	topPoint := findPoint(points, 0.75)
	bottomPoint := findPoint(points, 0.25)
	highlightsPoint := findPoint(points, 1.0)
	shadowsPoint := findPoint(points, 0.0)

	// Exposure = vertical shift from linear at midpoint (0.5 → X)
	// Normalize to -5.0/+5.0 range (matching universalToToneCurve which uses ÷5)
	exposure = (midpoint.Output - 0.5) * 5.0

	// Contrast = slope difference (top - bottom)
	// Linear slope = 1.0, steeper = positive, flatter = negative
	slopeDiff := (topPoint.Output - bottomPoint.Output) / 0.5
	contrast = slopeDiff - 1.0

	// Highlights = top-end deviation (1.0 → X)
	// Normalize to -1.0/+1.0 range
	highlights = (highlightsPoint.Output - 1.0) / 0.25

	// Shadows = bottom-end deviation (0.0 → X)
	// Normalize to -1.0/+1.0 range
	shadows = shadowsPoint.Output / 0.25

	return exposure, contrast, highlights, shadows
}

// findPoint finds an exact tone curve point or interpolates between adjacent points.
func findPoint(points []ToneCurvePoint, input float64) ToneCurvePoint {
	// Look for exact match
	for _, p := range points {
		if math.Abs(p.Input-input) < 0.0001 {
			return p
		}
	}

	// Find adjacent points for interpolation
	var before, after *ToneCurvePoint
	for i := range points {
		if points[i].Input < input {
			if before == nil || points[i].Input > before.Input {
				before = &points[i]
			}
		}
		if points[i].Input > input {
			if after == nil || points[i].Input < after.Input {
				after = &points[i]
			}
		}
	}

	// Interpolate
	if before != nil && after != nil {
		// Linear interpolation
		ratio := (input - before.Input) / (after.Input - before.Input)
		output := before.Output + ratio*(after.Output-before.Output)
		return ToneCurvePoint{Input: input, Output: output}
	}

	// Fallback: use linear (input = output)
	return ToneCurvePoint{Input: input, Output: input}
}

// profileToUniversal converts a parsed DCP CameraProfile to UniversalRecipe.
//
// Steps:
//  1. Analyze tone curve to extract exposure/contrast/highlights/shadows
//  2. Clamp values to valid UniversalRecipe ranges
//  3. Store profile metadata (name, matrices, baseline exposure)
//  4. Log warning if non-identity color matrices detected
func profileToUniversal(profile *CameraProfile) *models.UniversalRecipe {
	recipe := &models.UniversalRecipe{
		Metadata: make(map[string]interface{}),
	}

	// Analyze tone curve (if present)
	if len(profile.ToneCurve) > 0 {
		exposure, contrast, highlights, shadows := analyzeToneCurve(profile.ToneCurve)

		// Clamp and convert to UniversalRecipe ranges
		// Exposure is float64 (-5.0 to +5.0 in Recipe, we use -2.0 to +2.0 for DCP)
		recipe.Exposure = clampFloat64(exposure, -2.0, 2.0)

		// Note: BaselineExposureOffset is camera-specific calibration, not part of preset
		// It's stored in metadata but not added to exposure for round-trip consistency

		// Contrast, Highlights, Shadows are int (-100 to +100 in Recipe)
		// Convert from normalized -1.0/+1.0 to -100/+100
		recipe.Contrast = clampInt(int(contrast*100), -100, 100)
		recipe.Highlights = clampInt(int(highlights*100), -100, 100)
		recipe.Shadows = clampInt(int(shadows*100), -100, 100)
	}

	// Store profile name in metadata (even if empty - caller can use filename as fallback)
	recipe.Metadata["profile_name"] = profile.ProfileName

	// Check for non-identity color matrices (log warning if present)
	if profile.ColorMatrix1 != nil && !isIdentityMatrix(profile.ColorMatrix1) {
		// Store matrix in metadata for future use
		recipe.Metadata["color_matrix_1"] = profile.ColorMatrix1
		recipe.Metadata["color_calibration_warning"] = "DCP contains color calibration matrices (not supported in MVP)"
	}

	if profile.ColorMatrix2 != nil && !isIdentityMatrix(profile.ColorMatrix2) {
		recipe.Metadata["color_matrix_2"] = profile.ColorMatrix2
	}

	// Store baseline exposure in metadata
	if profile.BaselineExposure != 0.0 {
		recipe.Metadata["baseline_exposure_offset"] = profile.BaselineExposure
	}

	return recipe
}

// isIdentityMatrix checks if a color matrix is an identity matrix.
//
// Identity matrix:
//
//	1.0 0.0 0.0
//	0.0 1.0 0.0
//	0.0 0.0 1.0
//
// Tolerance: ±0.001
func isIdentityMatrix(matrix *Matrix) bool {
	if matrix == nil {
		return false
	}

	expected := [3][3]float64{
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{0.0, 0.0, 1.0},
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(matrix.Rows[i][j]-expected[i][j]) > 0.001 {
				return false
			}
		}
	}

	return true
}

// clampFloat64 clamps a float64 value to a specified range.
func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// clampInt clamps an int value to a specified range.
func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// universalToToneCurve generates a DCP tone curve from UniversalRecipe adjustments.
//
// Generates a 5-point piecewise linear curve normalized to 0.0-1.0 range:
//   - Point 1: (0.0, shadows-adjusted)
//   - Point 2: (0.25, shadows+contrast-adjusted)
//   - Point 3: (0.5, exposure+contrast midpoint)
//   - Point 4: (0.75, highlights+contrast-adjusted)
//   - Point 5: (1.0, highlights-adjusted)
//
// Algorithm:
//  1. Start with linear curve: {0.0→0.0}, {0.25→0.25}, {0.5→0.5}, {0.75→0.75}, {1.0→1.0}
//  2. Apply exposure: Vertical shift of all points (exposure/5.0 normalized)
//  3. Apply contrast: Steepen/flatten slope around midpoint (contrast/100.0 factor)
//  4. Apply highlights: Adjust top-end points 0.75-1.0 (highlights/100.0 × 0.125)
//  5. Apply shadows: Adjust bottom-end points 0.0-0.25 (shadows/100.0 × 0.125)
//  6. Clamp all outputs to 0.0-1.0 and ensure monotonic property
//
// All UniversalRecipe values are denormalized from their storage format:
//   - Exposure: -5.0 to +5.0 (float64)
//   - Contrast: -100 to +100 (int, converted to -1.0/+1.0)
//   - Highlights: -100 to +100 (int, converted to -1.0/+1.0)
//   - Shadows: -100 to +100 (int, converted to -1.0/+1.0)
//
// Returns array of 5 ToneCurvePoints with Input/Output values in 0.0-1.0 range.
func universalToToneCurve(recipe *models.UniversalRecipe) []ToneCurvePoint {
	// Start with linear 5-point curve (0.0-1.0 normalized)
	points := []ToneCurvePoint{
		{Input: 0.0, Output: 0.0},
		{Input: 0.25, Output: 0.25},
		{Input: 0.5, Output: 0.5},
		{Input: 0.75, Output: 0.75},
		{Input: 1.0, Output: 1.0},
	}

	// Extract parameters with normalization
	exposure := recipe.Exposure    // -5.0 to +5.0
	contrast := float64(recipe.Contrast) / 100.0   // -1.0 to +1.0
	highlights := float64(recipe.Highlights) / 100.0 // -1.0 to +1.0
	shadows := float64(recipe.Shadows) / 100.0       // -1.0 to +1.0

	// Normalize exposure to DCP range (÷5 = -1.0 to +1.0)
	exposureShift := exposure / 5.0

	// Step 1: Apply contrast (steepen/flatten slope around midpoint 0.5)
	// contrastFactor: 0.0 = flat, 1.0 = linear, 2.0 = double slope
	contrastFactor := 1.0 + contrast
	for i := range points {
		deviation := points[i].Input - 0.5
		points[i].Output = 0.5 + deviation*contrastFactor + exposureShift
	}

	// Step 2: Apply highlights (adjust top-end points 0.75-1.0)
	// Scale factor: ±0.125 max shift (highlights range -1.0 to +1.0)
	highlightsShift := highlights * 0.125
	points[3].Output += highlightsShift
	points[4].Output += highlightsShift

	// Step 3: Apply shadows (adjust bottom-end points 0.0-0.25)
	// Scale factor: ±0.125 max shift (shadows range -1.0 to +1.0)
	shadowsShift := shadows * 0.125
	points[0].Output += shadowsShift
	points[1].Output += shadowsShift

	// Step 4: Clamp all outputs to 0.0-1.0 range
	for i := range points {
		points[i].Output = clampFloat64(points[i].Output, 0.0, 1.0)
	}

	// Step 5: Ensure monotonic property (output[i] >= output[i-1])
	for i := 1; i < len(points); i++ {
		if points[i].Output < points[i-1].Output {
			points[i].Output = points[i-1].Output
		}
	}

	return points
}

// generateBinaryToneCurve converts ToneCurvePoint array to binary float32 format.
//
// Each point becomes 8 bytes (4-byte float32 input + 4-byte float32 output).
// Uses little-endian byte order to match DNG specification.
//
// Binary format:
//   [float32 input₁][float32 output₁][float32 input₂][float32 output₂]...
//
// For 5 points: 5 × 8 bytes = 40 bytes total
//
// Returns binary data ready to write to TIFF tag 50940 (ProfileToneCurve).
func generateBinaryToneCurve(points []ToneCurvePoint) []byte {
	buf := new(bytes.Buffer)

	for _, pt := range points {
		// Convert float64 to float32 (DCP uses 32-bit floats)
		binary.Write(buf, binary.LittleEndian, float32(pt.Input))
		binary.Write(buf, binary.LittleEndian, float32(pt.Output))
	}

	return buf.Bytes()
}

// SRational represents a signed rational number (numerator/denominator).
// Used for color matrices and baseline exposure in DCP files.
type SRational struct {
	Numerator   int32
	Denominator int32
}

// generateColorMatrix generates ColorMatrix1 for Standard Light A illuminant.
//
// Uses Nikon Z f calibrated color matrices from Adobe Camera Raw.
// These matrices transform camera native RGB to CIE XYZ color space.
//
// ColorMatrix1 (Standard Light A illuminant):
//   1.3904 -0.7947  0.0654
//  -0.432   1.2105  0.2497
//  -0.0235  0.083   0.9243
//
// Returns array of 9 SRational values ready for tag 50721 (ColorMatrix1).
func generateColorMatrix() []SRational {
	return []SRational{
		// Row 1: [1.3904, -0.7947, 0.0654]
		{Numerator: 13904, Denominator: 10000}, {Numerator: -7947, Denominator: 10000}, {Numerator: 654, Denominator: 10000},
		// Row 2: [-0.432, 1.2105, 0.2497]
		{Numerator: -4320, Denominator: 10000}, {Numerator: 12105, Denominator: 10000}, {Numerator: 2497, Denominator: 10000},
		// Row 3: [-0.0235, 0.083, 0.9243]
		{Numerator: -235, Denominator: 10000}, {Numerator: 830, Denominator: 10000}, {Numerator: 9243, Denominator: 10000},
	}
}

// generateColorMatrix2 generates ColorMatrix2 for D65 illuminant.
//
// ColorMatrix2 (D65 illuminant):
//   1.1607 -0.4491 -0.0977
//  -0.4522  1.246   0.2304
//  -0.0458  0.1519  0.7616
//
// Returns array of 9 SRational values ready for tag 50722 (ColorMatrix2).
func generateColorMatrix2() []SRational {
	return []SRational{
		// Row 1: [1.1607, -0.4491, -0.0977]
		{Numerator: 11607, Denominator: 10000}, {Numerator: -4491, Denominator: 10000}, {Numerator: -977, Denominator: 10000},
		// Row 2: [-0.4522, 1.246, 0.2304]
		{Numerator: -4522, Denominator: 10000}, {Numerator: 12460, Denominator: 10000}, {Numerator: 2304, Denominator: 10000},
		// Row 3: [-0.0458, 0.1519, 0.7616]
		{Numerator: -458, Denominator: 10000}, {Numerator: 1519, Denominator: 10000}, {Numerator: 7616, Denominator: 10000},
	}
}

// generateForwardMatrix generates ForwardMatrix (XYZ → camera RGB).
//
// Forward matrices are the same for both illuminants in Nikon Z f profiles.
//
// ForwardMatrix:
//   0.7978  0.1352  0.0313
//   0.288   0.7119  0.0001
//   0       0       0.8251
//
// Returns array of 9 SRational values ready for tags 50964/50965.
func generateForwardMatrix() []SRational {
	return []SRational{
		// Row 1: [0.7978, 0.1352, 0.0313]
		{Numerator: 7978, Denominator: 10000}, {Numerator: 1352, Denominator: 10000}, {Numerator: 313, Denominator: 10000},
		// Row 2: [0.288, 0.7119, 0.0001]
		{Numerator: 2880, Denominator: 10000}, {Numerator: 7119, Denominator: 10000}, {Numerator: 1, Denominator: 10000},
		// Row 3: [0, 0, 0.8251]
		{Numerator: 0, Denominator: 10000}, {Numerator: 0, Denominator: 10000}, {Numerator: 8251, Denominator: 10000},
	}
}

// srationalArrayToBytes converts SRational array to binary format.
//
// Each SRational becomes 8 bytes:
//   - Bytes 0-3: Signed int32 numerator (little-endian)
//   - Bytes 4-7: Signed int32 denominator (little-endian)
//
// For 9 SRationals (color matrix): 9 × 8 bytes = 72 bytes total
//
// Binary format example for SRational{1, 1}:
//   [0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00]
//    └─────numerator = 1────┘ └────denominator = 1───┘
//
// Returns binary data ready to write to TIFF tags 50721/50722 (ColorMatrix1/2).
func srationalArrayToBytes(srs []SRational) []byte {
	buf := new(bytes.Buffer)

	for _, sr := range srs {
		binary.Write(buf, binary.LittleEndian, sr.Numerator)
		binary.Write(buf, binary.LittleEndian, sr.Denominator)
	}

	return buf.Bytes()
}

// generateIdentityLUT generates an identity 3D color lookup table.
//
// Dimensions: 90 (hue) × 16 (saturation) × 16 (value) = 23,040 entries
// Each entry: RGB triplet (3 float32 values) = 69,120 float values total
//
// Identity LUT: Each HSV coordinate maps to its RGB equivalent without transformation.
// This creates a pass-through LUT that doesn't modify colors.
//
// Returns binary data ready to write to TIFF tag 50982 (ProfileLookTableData).
func generateIdentityLUT() []byte {
	const (
		hueSteps   = 90 // 360° / 4° per step
		satSteps   = 16
		valSteps   = 16
	)

	buf := new(bytes.Buffer)

	// Iterate through HSV space in order: value, saturation, hue
	// (DNG spec requires this specific ordering)
	for v := 0; v < valSteps; v++ {
		for s := 0; s < satSteps; s++ {
			for h := 0; h < hueSteps; h++ {
				// Normalize coordinates to 0.0-1.0
				hNorm := float32(h) / float32(hueSteps-1)
				sNorm := float32(s) / float32(satSteps-1)
				vNorm := float32(v) / float32(valSteps-1)

				// Convert HSV to RGB (identity transformation)
				r, g, b := hsvToRGB(hNorm, sNorm, vNorm)

				// Write RGB triplet as float32 (little-endian)
				binary.Write(buf, binary.LittleEndian, r)
				binary.Write(buf, binary.LittleEndian, g)
				binary.Write(buf, binary.LittleEndian, b)
			}
		}
	}

	return buf.Bytes()
}

// hsvToRGB converts HSV color space to RGB.
//
// Input ranges:
//   - h: 0.0-1.0 (hue, 0° = 0.0, 360° = 1.0)
//   - s: 0.0-1.0 (saturation)
//   - v: 0.0-1.0 (value/brightness)
//
// Output ranges:
//   - r, g, b: 0.0-1.0
//
// Algorithm: Standard HSV→RGB conversion
func hsvToRGB(h, s, v float32) (r, g, b float32) {
	if s == 0 {
		// Achromatic (gray)
		return v, v, v
	}

	h = h * 6.0 // sector 0 to 5
	i := int(math.Floor(float64(h)))
	f := h - float32(i) // fractional part
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	switch i % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	default: // case 5:
		return v, p, q
	}
}
