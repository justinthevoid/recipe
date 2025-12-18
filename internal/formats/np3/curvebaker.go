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
