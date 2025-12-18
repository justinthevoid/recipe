package np3

// curvebaker.go - Converts XMP parametric curve values to NP3 control points
//
// Lightroom's parametric curves have 4 zones:
// - Shadows: 0% to ShadowSplit (default 25%)
// - Darks: ShadowSplit to MidtoneSplit (default 50%)
// - Lights: MidtoneSplit to HighlightSplit (default 75%)
// - Highlights: HighlightSplit to 100%
//
// Each zone slider (-100 to +100) adjusts the output value for that zone.
// Positive = lift (output higher than input), Negative = compress (output lower)
//
// NP3 uses control points (x,y pairs where x=input, y=output) to define curves.
// Linear curve: x == y for all points
// S-curve: y > x in shadows (lift), y < x in highlights (compress)

import (
	"math"

	"github.com/justin/recipe/internal/models"
)

// ControlPoint represents a single point on the tone curve
type ControlPoint struct {
	X int // Input value (0-255)
	Y int // Output value (0-255)
}

// ParametricToControlPoints converts Lightroom parametric curve settings to NP3 control points.
//
// Parameters:
//   - shadows: Shadow zone adjustment (-100 to +100)
//   - darks: Dark zone adjustment (-100 to +100)
//   - lights: Light zone adjustment (-100 to +100)
//   - highlights: Highlight zone adjustment (-100 to +100)
//   - shadowSplit: Shadow/Dark boundary (0-100, default 25)
//   - midtoneSplit: Dark/Light boundary (0-100, default 50)
//   - highlightSplit: Light/Highlight boundary (0-100, default 75)
//
// Returns a slice of 6-8 control points suitable for NP3 format.
func ParametricToControlPoints(shadows, darks, lights, highlights int,
	shadowSplit, midtoneSplit, highlightSplit int) []ControlPoint {

	// Apply defaults if split points are 0 (not specified)
	if shadowSplit == 0 {
		shadowSplit = 25
	}
	if midtoneSplit == 0 {
		midtoneSplit = 50
	}
	if highlightSplit == 0 {
		highlightSplit = 75
	}

	// Convert split points from 0-100 to 0-255 range
	shadowX := (shadowSplit * 255) / 100
	midtoneX := (midtoneSplit * 255) / 100
	highlightX := (highlightSplit * 255) / 100

	// Create control points at zone boundaries
	// Each zone adjustment adds/subtracts from the linear diagonal (y=x)
	// Adjustment range: -100 to +100 maps to approximately -50 to +50 in output
	points := make([]ControlPoint, 0, 8)

	// Calculate Y offsets for each zone
	// Scale: -100 to +100 slider maps to approximately -15 to +15 pixel offset
	// This is a subtle adjustment - Lightroom parametric curves are relatively gentle
	// Previous scale of 50/100 was too aggressive
	shadowOffset := (shadows * 15) / 100
	darkOffset := (darks * 15) / 100
	lightOffset := (lights * 15) / 100
	highlightOffset := (highlights * 15) / 100

	// Start at 0,0
	points = append(points, ControlPoint{X: 0, Y: 0})

	// Point 1: In shadow zone (1/3 of shadow zone)
	x1 := shadowX / 3
	if x1 < 20 {
		x1 = 20 // Minimum x to avoid extreme black point
	}
	y1 := clampByte(x1 + shadowOffset)
	points = append(points, ControlPoint{X: x1, Y: y1})

	// Point 2: Shadow/Dark boundary
	// Blend shadow and dark offsets
	y2 := clampByte(shadowX + (shadowOffset+darkOffset)/2)
	points = append(points, ControlPoint{X: shadowX, Y: y2})

	// Point 3: Middle of dark zone
	darkMidX := (shadowX + midtoneX) / 2
	y3 := clampByte(darkMidX + darkOffset)
	points = append(points, ControlPoint{X: darkMidX, Y: y3})

	// Point 4: Dark/Light boundary (midtone)
	// Blend dark and light offsets
	y4 := clampByte(midtoneX + (darkOffset+lightOffset)/2)
	points = append(points, ControlPoint{X: midtoneX, Y: y4})

	// Point 5: Middle of light zone
	lightMidX := (midtoneX + highlightX) / 2
	y5 := clampByte(lightMidX + lightOffset)
	points = append(points, ControlPoint{X: lightMidX, Y: y5})

	// Point 6: Light/Highlight boundary
	// Blend light and highlight offsets
	y6 := clampByte(highlightX + (lightOffset+highlightOffset)/2)
	points = append(points, ControlPoint{X: highlightX, Y: y6})

	// Point 7: Near white point (maintain highlight adjustment)
	x7 := highlightX + (255-highlightX)/2
	if x7 > 240 {
		x7 = 240 // Keep some headroom
	}
	y7 := clampByte(x7 + highlightOffset)
	points = append(points, ControlPoint{X: x7, Y: y7})

	// Final point approaches (255, 255)
	points = append(points, ControlPoint{X: 255, Y: 255})

	return points
}

// HasParametricCurve checks if any parametric curve adjustments are non-zero
func HasParametricCurve(shadows, darks, lights, highlights int) bool {
	return shadows != 0 || darks != 0 || lights != 0 || highlights != 0
}

// clampByte clamps an integer to the valid byte range [0, 255]
func clampByte(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

// IsLinearCurve checks if the control points represent a linear (identity) curve
// within a tolerance. Used for optimization to skip writing curve data.
func IsLinearCurve(points []ControlPoint, tolerance int) bool {
	for _, p := range points {
		if absInt(p.Y-p.X) > tolerance {
			return false
		}
	}
	return true
}

// absInt returns the absolute value of an integer
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// MergeRGBCurvesToMaster merges RGB channel curves into a master curve.
//
// When RGB curves are identical (a common film emulation technique), they function
// as a luminosity adjustment and can be merged with the master curve.
//
// Strategy:
//  1. Convert point curves to 256-entry LUTs for smooth interpolation
//  2. Apply master curve first (if any)
//  3. If RGB curves are identical, apply them as a single luminosity curve
//  4. If RGB curves differ, average their effect (approximation)
//  5. Return merged curve as control points
//
// Parameters:
//   - master: Master tone curve points (may be nil or empty)
//   - red, green, blue: Per-channel curve points (may be nil or empty)
//
// Returns merged curve as control points suitable for NP3.
func MergeRGBCurvesToMaster(master, red, green, blue []ControlPoint) []ControlPoint {
	// If no RGB curves, just return master
	if len(red) == 0 && len(green) == 0 && len(blue) == 0 {
		return master
	}

	// Convert master to LUT (or identity if empty)
	masterLUT := PointsToCurveLUT(master)

	// Check if all RGB curves are identical
	if CurvesAreIdentical(red, green, blue, 2) {
		// Use red as the representative curve (they're all the same)
		rgbLUT := PointsToCurveLUT(red)

		// Compose: output = rgbLUT[masterLUT[input]]
		composedLUT := make([]int, 256)
		for i := 0; i < 256; i++ {
			intermediate := masterLUT[i]
			composedLUT[i] = clampByte(rgbLUT[intermediate])
		}

		return LUTToControlPoints(composedLUT, 8)
	}

	// RGB curves differ - average their effect
	// This is an approximation since we can't represent per-channel in NP3
	redLUT := PointsToCurveLUT(red)
	greenLUT := PointsToCurveLUT(green)
	blueLUT := PointsToCurveLUT(blue)

	// For each input, apply master first, then average RGB effects
	composedLUT := make([]int, 256)
	for i := 0; i < 256; i++ {
		intermediate := masterLUT[i]

		// Apply each channel curve to the intermediate value
		rOut := redLUT[clampByte(intermediate)]
		gOut := greenLUT[clampByte(intermediate)]
		bOut := blueLUT[clampByte(intermediate)]

		// Average the outputs (approximation of luminosity effect)
		composedLUT[i] = clampByte((rOut + gOut + bOut) / 3)
	}

	return LUTToControlPoints(composedLUT, 8)
}

// CurvesAreIdentical checks if multiple curves are identical within a tolerance.
// Returns true if all curves have the same points within the given tolerance,
// or if all curves are nil/empty.
func CurvesAreIdentical(red, green, blue []ControlPoint, tolerance int) bool {
	// All empty is considered identical
	if len(red) == 0 && len(green) == 0 && len(blue) == 0 {
		return true
	}

	// If any is empty while others aren't, not identical
	if len(red) == 0 || len(green) == 0 || len(blue) == 0 {
		return false
	}

	// All must have same length
	if len(red) != len(green) || len(green) != len(blue) {
		return false
	}

	// Compare each point
	for i := range red {
		if absInt(red[i].X-green[i].X) > tolerance ||
			absInt(red[i].X-blue[i].X) > tolerance ||
			absInt(red[i].Y-green[i].Y) > tolerance ||
			absInt(red[i].Y-blue[i].Y) > tolerance {
			return false
		}
	}

	return true
}

// PointsToCurveLUT converts control points to a 256-entry LUT using linear interpolation.
// Returns an identity curve (y=x) if points is empty.
func PointsToCurveLUT(points []ControlPoint) []int {
	lut := make([]int, 256)

	// Empty or nil returns identity curve
	if len(points) == 0 {
		for i := 0; i < 256; i++ {
			lut[i] = i
		}
		return lut
	}

	// Sort points by X (defensive, should already be sorted)
	sortedPoints := make([]ControlPoint, len(points))
	copy(sortedPoints, points)
	// Simple insertion sort since we typically have <10 points
	for i := 1; i < len(sortedPoints); i++ {
		for j := i; j > 0 && sortedPoints[j-1].X > sortedPoints[j].X; j-- {
			sortedPoints[j-1], sortedPoints[j] = sortedPoints[j], sortedPoints[j-1]
		}
	}

	// Linear interpolation between control points
	for i := 0; i < 256; i++ {
		// Find surrounding control points
		var p1, p2 ControlPoint

		// If at or below first point, use first point's Y value directly
		// This correctly handles cases where first point is at X=0
		if i <= sortedPoints[0].X {
			lut[i] = clampByte(sortedPoints[0].Y)
			continue
		} else if i >= sortedPoints[len(sortedPoints)-1].X {
			// If above last point, use line from last point to (255, 255)
			p1 = sortedPoints[len(sortedPoints)-1]
			p2 = ControlPoint{X: 255, Y: 255}
		} else {
			// Find bracketing points
			for j := 0; j < len(sortedPoints)-1; j++ {
				if i >= sortedPoints[j].X && i <= sortedPoints[j+1].X {
					p1 = sortedPoints[j]
					p2 = sortedPoints[j+1]
					break
				}
			}
		}

		// Linear interpolation
		if p2.X == p1.X {
			lut[i] = p1.Y
		} else {
			t := float64(i-p1.X) / float64(p2.X-p1.X)
			lut[i] = clampByte(p1.Y + int(t*float64(p2.Y-p1.Y)+0.5))
		}
	}

	return lut
}

// LUTToControlPoints samples a 256-entry LUT to create control points.
// numPoints specifies how many points to generate (typically 5-10).
func LUTToControlPoints(lut []int, numPoints int) []ControlPoint {
	if numPoints < 2 {
		numPoints = 2
	}
	if numPoints > 20 {
		numPoints = 20
	}

	points := make([]ControlPoint, 0, numPoints)

	// Always include first point
	points = append(points, ControlPoint{X: 0, Y: clampByte(lut[0])})

	// Sample evenly spaced points
	step := 255 / (numPoints - 1)
	for i := 1; i < numPoints-1; i++ {
		x := i * step
		if x > 254 {
			x = 254
		}
		points = append(points, ControlPoint{X: x, Y: clampByte(lut[x])})
	}

	// Always include last point
	points = append(points, ControlPoint{X: 255, Y: clampByte(lut[255])})

	return points
}

// HasRGBCurves checks if any per-channel RGB curves are defined.
func HasRGBCurves(red, green, blue []ControlPoint) bool {
	return len(red) > 0 || len(green) > 0 || len(blue) > 0
}

// toneCurvePointsToControlPoints converts models.ToneCurvePoint slice to ControlPoint slice.
// This enables interoperability between the UniversalRecipe format and internal curve processing.
func toneCurvePointsToControlPoints(points []models.ToneCurvePoint) []ControlPoint {
	if len(points) == 0 {
		return nil
	}

	result := make([]ControlPoint, len(points))
	for i, p := range points {
		result[i] = ControlPoint{
			X: p.Input,
			Y: p.Output,
		}
	}
	return result
}

// GetStandardBaseCurveLUT returns a LUT approximating the "Standard" Picture Control curve.
// This curve adds contrast (S-curve) relative to the linear "Flexible Color" base.
// Used when adapting XMP presets designed for "Adobe Standard" to "Flexible Color".
func GetStandardBaseCurveLUT() []int {
	// Aggressive S-Curve points to match Lightroom's deep shadows
	// Crushes shadows significantly (64->35) and lifts highlights (192->220)
	points := []ControlPoint{
		{X: 0, Y: 0},
		{X: 32, Y: 10},  // Was 15. Aggressively dark.
		{X: 64, Y: 35},  // Was 45. Much darker mids.
		{X: 128, Y: 128},
		{X: 192, Y: 220}, // Was 215. Slightly punchier highlights.
		{X: 224, Y: 245}, // Was 240.
		{X: 255, Y: 255},
	}
	return PointsToCurveLUT(points)
}

// ApplyCurveToLUT applies a modification curve (LUT) on top of a base curve (LUT).
// Result[i] = modifier[base[i]]
// This effectively chains the curves: modifier(base(input))
func ApplyCurveToLUT(baseLUT, modifierLUT []int) []int {
	result := make([]int, 256)
	for i := 0; i < 256; i++ {
		// Get output from base curve
		baseOut := clampByte(baseLUT[i])
		// Use as input to modifier curve
		finalOut := clampByte(modifierLUT[baseOut])
		result[i] = finalOut
	}
	return result
}

// ApplyExposureAndBlacksToLUT modifies a curve LUT to simulate Exposure and Blacks adjustments.
func ApplyExposureAndBlacksToLUT(lut []int, exposure float64, blacks int) []int {
	newLUT := make([]int, 256)
	
	// Exposure: gain on input index
	exposureMult := math.Pow(2, exposure)
	
	// Blacks: shift on OUTPUT (to lower lifted black points)
	// We want to darken the shadows without affecting highlights too much.
	// Scale factor 1.0 means Blacks -13 lowers output by ~13 levels at the bottom.
	blacksShift := float64(blacks) * 1.0

	for i := 0; i < 256; i++ {
		// 1. Exposure (Input Gain)
		input := float64(i)
		effectiveInput := input * exposureMult
		
		// Clamp input index
		idx := int(math.Round(effectiveInput))
		if idx < 0 {
			idx = 0
		} else if idx > 255 {
			idx = 255
		}
		
		// Get base output
		val := float64(lut[idx])
		
		// 2. Blacks (Output Shift)
		// Apply shift scaled by proximity to black (1.0 at black, 0.0 at white)
		// This prevents shifting the white point.
		weight := 1.0 - (val / 255.0)
		if weight < 0 {
			weight = 0
		}
		
		val += (blacksShift * weight)
		
		// Clamp output
		out := int(math.Round(val))
		if out < 0 {
			out = 0
		} else if out > 255 {
			out = 255
		}
		
		newLUT[i] = out
	}
	
	return newLUT
}

