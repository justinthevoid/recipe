package dcp

import (
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
	// Normalize to -2.0/+2.0 range
	exposure = (midpoint.Output - 0.5) / 0.25

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

		// Add baseline exposure offset
		recipe.Exposure += profile.BaselineExposure

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
