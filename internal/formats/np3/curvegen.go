// Package np3 provides functionality for generating test tone curves for NX Studio validation.
//
// This file contains test curve generation utilities to validate that our NP3 curve
// format is correctly interpreted by Nikon NX Studio / Picture Control Utility 2.

package np3

import (
	"fmt"
	"math"
)

// CurveType defines the type of test curve to generate
type CurveType string

const (
	// CurveLinear generates an identity curve (y = x)
	CurveLinear CurveType = "linear"

	// CurveSCurve generates a sigmoid contrast boost curve
	CurveSCurve CurveType = "s-curve"

	// CurveShadowsBoost lifts shadows by approximately 20%
	CurveShadowsBoost CurveType = "shadows-boost"

	// CurveHighlightsCompress compresses highlights by approximately 20%
	CurveHighlightsCompress CurveType = "highlights-compress"

	// CurveShadowsDarks tests parametric Shadows+Darks combined
	CurveShadowsDarks CurveType = "shadows-darks"

	// CurveLightsHighlights tests parametric Lights+Highlights combined
	CurveLightsHighlights CurveType = "lights-highlights"
)

// GenerateTestCurve creates a 257-entry tone curve LUT for testing.
// Returns an array of 257 uint16 values (0-32767 range) suitable for writing to NP3 offset 460.
//
// The curve entries map:
//   - Entry 0 = black (input 0)
//   - Entry 128 = midtone (input 128)
//   - Entry 256 = white (input 255)
//
// NP3 uses 16-bit big-endian values where 0 = black, 32767 = white.
func GenerateTestCurve(curveType CurveType) ([]uint16, error) {
	lut := make([]uint16, 257)

	switch curveType {
	case CurveLinear:
		generateLinearCurve(lut)
	case CurveSCurve:
		generateSCurve(lut, 1.5) // Medium contrast boost
	case CurveShadowsBoost:
		generateShadowsBoostCurve(lut, 0.20) // +20% shadow lift
	case CurveHighlightsCompress:
		generateHighlightsCompressCurve(lut, 0.20) // -20% highlight compression
	case CurveShadowsDarks:
		generateParametricCurve(lut, 30, 20, 0, 0) // Shadows+30, Darks+20
	case CurveLightsHighlights:
		generateParametricCurve(lut, 0, 0, -20, -30) // Lights-20, Highlights-30
	default:
		return nil, fmt.Errorf("unknown curve type: %s", curveType)
	}

	return lut, nil
}

// generateLinearCurve creates an identity curve (y = x)
func generateLinearCurve(lut []uint16) {
	for i := 0; i <= 256; i++ {
		// Map 0-256 input range to 0-32767 output range
		lut[i] = uint16(float64(i) / 256.0 * 32767.0)
	}
}

// generateSCurve creates a sigmoid contrast curve
// steepness controls the S-curve intensity (1.0 = mild, 2.0 = strong)
func generateSCurve(lut []uint16, steepness float64) {
	for i := 0; i <= 256; i++ {
		// Normalize input to -1 to +1 range
		x := (float64(i)/256.0)*2.0 - 1.0

		// Apply sigmoid function: y = x / (1 + |x|^steepness)
		// This creates an S-curve that compresses shadows and highlights
		y := x * math.Pow(math.Abs(x), steepness-1) / (1 + math.Pow(math.Abs(x), steepness))

		// Handle the midpoint separately to avoid division issues
		if i == 128 {
			y = 0.0
		}

		// Convert back to 0-1 range
		output := (y + 1.0) / 2.0

		// Clamp and convert to 16-bit
		if output < 0 {
			output = 0
		} else if output > 1 {
			output = 1
		}
		lut[i] = uint16(output * 32767.0)
	}
}

// generateShadowsBoostCurve lifts shadow tones
// amount is the percentage to lift (0.2 = 20%)
func generateShadowsBoostCurve(lut []uint16, amount float64) {
	for i := 0; i <= 256; i++ {
		input := float64(i) / 256.0

		// Apply shadow lift using curve: y = x + amount * (1-x) * (1-x) for shadows
		// This affects lower tones more than higher tones
		shadowWeight := (1.0 - input) * (1.0 - input)
		output := input + amount*shadowWeight

		// Clamp
		if output > 1.0 {
			output = 1.0
		}

		lut[i] = uint16(output * 32767.0)
	}
}

// generateHighlightsCompressCurve compresses highlight tones
// amount is the percentage to compress (0.2 = 20%)
func generateHighlightsCompressCurve(lut []uint16, amount float64) {
	for i := 0; i <= 256; i++ {
		input := float64(i) / 256.0

		// Apply highlight compression: y = x - amount * x * x for highlights
		// This affects higher tones more than lower tones
		highlightWeight := input * input
		output := input - amount*highlightWeight

		// Clamp
		if output < 0 {
			output = 0
		}

		lut[i] = uint16(output * 32767.0)
	}
}

// generateParametricCurve creates a curve from Lightroom-style parametric adjustments.
// Values are in -100 to +100 range like Lightroom.
//
// Zone boundaries (default):
//   - Shadows: 0-25% (0-64 in 0-255 space)
//   - Darks: 25-50% (64-128)
//   - Lights: 50-75% (128-192)
//   - Highlights: 75-100% (192-255)
func generateParametricCurve(lut []uint16, shadows, darks, lights, highlights int) {
	// Zone boundaries (normalized 0-1)
	const shadowEnd = 0.25
	const darksEnd = 0.50
	const lightsEnd = 0.75
	// highlightsEnd = 1.0

	for i := 0; i <= 256; i++ {
		input := float64(i) / 256.0
		output := input // Start with linear

		// Determine zone and apply adjustment
		var adjustment float64
		var zoneWeight float64

		if input < shadowEnd {
			// Shadows zone (0-25%)
			zoneWeight = 1.0 - (input / shadowEnd) // Strongest at black, fades to 0 at shadowEnd
			adjustment = float64(shadows) / 100.0 * 0.3 * zoneWeight
		} else if input < darksEnd {
			// Darks zone (25-50%)
			zonePos := (input - shadowEnd) / (darksEnd - shadowEnd)
			// Bell curve weight - strongest in middle of zone
			zoneWeight = 1.0 - math.Abs(zonePos-0.5)*2.0
			adjustment = float64(darks) / 100.0 * 0.2 * zoneWeight
		} else if input < lightsEnd {
			// Lights zone (50-75%)
			zonePos := (input - darksEnd) / (lightsEnd - darksEnd)
			zoneWeight = 1.0 - math.Abs(zonePos-0.5)*2.0
			adjustment = float64(lights) / 100.0 * 0.2 * zoneWeight
		} else {
			// Highlights zone (75-100%)
			zoneWeight = (input - lightsEnd) / (1.0 - lightsEnd) // Strongest at white
			adjustment = float64(highlights) / 100.0 * 0.3 * zoneWeight
		}

		output = input + adjustment

		// Clamp
		if output < 0 {
			output = 0
		} else if output > 1 {
			output = 1
		}

		lut[i] = uint16(output * 32767.0)
	}
}

// GenerateTestNP3WithCurve creates an NP3 file bytes with a specific test curve.
// This is a convenience function for generating test files.
func GenerateTestNP3WithCurve(name string, curveType CurveType) ([]byte, error) {
	// Generate the curve
	curve, err := GenerateTestCurve(curveType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate curve: %w", err)
	}

	// Generate control points from the LUT to enable custom curve in NX Studio
	// NX Studio requires toneCurvePointCount > 0 to enable the "Use Custom Tone Curve" checkbox
	// We sample 5 key points from the LUT: 0%, 25%, 50%, 75%, 100%
	controlPoints := make([]toneCurvePoint, 5)
	pointIndices := []int{0, 64, 128, 192, 256}

	for i, idx := range pointIndices {
		// Convert 16-bit LUT value (0-32767) back to 8-bit (0-255) for control points
		inputVal := uint8(idx)
		if idx == 256 {
			inputVal = 255
		}
		outputVal := uint8(float64(curve[idx]) / 32767.0 * 255.0)

		controlPoints[i] = toneCurvePoint{
			position: OffsetToneCurvePoints + (i * 2),
			value1:   inputVal,  // Input (x)
			value2:   outputVal, // Output (y)
		}
	}

	// Create NP3 parameters with both control points AND raw curve
	params := &np3Parameters{
		toneCurvePointCount: len(controlPoints), // Must be > 0 for NX Studio to enable custom curve
		toneCurvePoints:     controlPoints,
		toneCurveRaw:        curve,
	}

	// Generate the NP3 binary
	data, err := encodeBinary(params, name)
	if err != nil {
		return nil, fmt.Errorf("failed to encode NP3: %w", err)
	}

	return data, nil
}
