package np3

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

func TestParametricToControlPoints_Linear(t *testing.T) {
	// All zeros should produce a nearly linear curve
	points := ParametricToControlPoints(0, 0, 0, 0, 25, 50, 75)

	if len(points) < 6 {
		t.Errorf("Expected at least 6 control points, got %d", len(points))
	}

	// Check that all points are on or near the diagonal (y ≈ x)
	for i, p := range points {
		diff := absInt(p.Y - p.X)
		if diff > 5 { // Allow small tolerance for rounding
			t.Errorf("Point %d: expected linear (y≈x), got (%d, %d), diff=%d", i, p.X, p.Y, diff)
		}
	}
}

func TestParametricToControlPoints_ShadowLift(t *testing.T) {
	// Shadows = +50 should lift shadow values (y > x in shadows)
	points := ParametricToControlPoints(50, 0, 0, 0, 25, 50, 75)

	if len(points) < 6 {
		t.Fatalf("Expected at least 6 control points, got %d", len(points))
	}

	// First point should show shadow lift
	if points[0].Y <= points[0].X {
		t.Errorf("Shadow lift: expected Y > X for first point, got (%d, %d)", points[0].X, points[0].Y)
	}

	// Last point should still be near (255, 255)
	lastPoint := points[len(points)-1]
	if lastPoint.X != 255 || lastPoint.Y != 255 {
		t.Errorf("Expected last point (255, 255), got (%d, %d)", lastPoint.X, lastPoint.Y)
	}
}

func TestParametricToControlPoints_HighlightCompress(t *testing.T) {
	// Highlights = -50 should compress highlights (y < x in highlights)
	points := ParametricToControlPoints(0, 0, 0, -50, 25, 50, 75)

	if len(points) < 6 {
		t.Fatalf("Expected at least 6 control points, got %d", len(points))
	}

	// Second to last point should show highlight compression
	nearEnd := points[len(points)-2]
	if nearEnd.Y >= nearEnd.X {
		t.Errorf("Highlight compression: expected Y < X near end, got (%d, %d)", nearEnd.X, nearEnd.Y)
	}
}

func TestParametricToControlPoints_SCurve(t *testing.T) {
	// Classic S-curve: lift shadows, compress highlights
	points := ParametricToControlPoints(30, 15, -15, -30, 25, 50, 75)

	if len(points) < 6 {
		t.Fatalf("Expected at least 6 control points, got %d", len(points))
	}

	// First point should show lift (y > x)
	if points[0].Y <= points[0].X {
		t.Errorf("S-curve shadow lift: expected Y > X, got (%d, %d)", points[0].X, points[0].Y)
	}

	// Near-end point should show compression (y < x)
	nearEnd := points[len(points)-2]
	if nearEnd.Y >= nearEnd.X {
		t.Errorf("S-curve highlight compress: expected Y < X, got (%d, %d)", nearEnd.X, nearEnd.Y)
	}
}

func TestParametricToControlPoints_DefaultSplits(t *testing.T) {
	// When split points are 0, defaults should be applied
	points := ParametricToControlPoints(25, 25, 25, 25, 0, 0, 0)

	if len(points) < 6 {
		t.Fatalf("Expected at least 6 control points, got %d", len(points))
	}

	// All points should have Y > X due to positive adjustments
	for i, p := range points[:len(points)-1] { // Exclude last point (255,255)
		if p.Y < p.X {
			t.Errorf("Point %d: expected lift (Y >= X) with positive adjustments, got (%d, %d)", i, p.X, p.Y)
		}
	}
}

func TestHasParametricCurve(t *testing.T) {
	tests := []struct {
		name       string
		s, d, l, h int
		want       bool
	}{
		{"All zeros", 0, 0, 0, 0, false},
		{"Shadows only", 10, 0, 0, 0, true},
		{"Darks only", 0, -20, 0, 0, true},
		{"Lights only", 0, 0, 30, 0, true},
		{"Highlights only", 0, 0, 0, -40, true},
		{"All non-zero", 10, 20, 30, 40, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasParametricCurve(tt.s, tt.d, tt.l, tt.h)
			if got != tt.want {
				t.Errorf("HasParametricCurve(%d,%d,%d,%d) = %v, want %v",
					tt.s, tt.d, tt.l, tt.h, got, tt.want)
			}
		})
	}
}

func TestIsLinearCurve(t *testing.T) {
	linearPoints := []ControlPoint{
		{X: 30, Y: 30},
		{X: 60, Y: 60},
		{X: 120, Y: 120},
		{X: 200, Y: 200},
		{X: 255, Y: 255},
	}

	if !IsLinearCurve(linearPoints, 3) {
		t.Error("Expected linear curve to be detected as linear")
	}

	nonLinearPoints := []ControlPoint{
		{X: 30, Y: 60}, // Y much greater than X
		{X: 60, Y: 60},
		{X: 120, Y: 120},
		{X: 200, Y: 180}, // Y less than X
		{X: 255, Y: 255},
	}

	if IsLinearCurve(nonLinearPoints, 3) {
		t.Error("Expected non-linear curve to NOT be detected as linear")
	}
}

func TestClampByte(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{-10, 0},
		{0, 0},
		{128, 128},
		{255, 255},
		{300, 255},
	}

	for _, tt := range tests {
		got := clampByte(tt.input)
		if got != tt.want {
			t.Errorf("clampByte(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

// TestParametricCurveIntegration tests the full flow from UniversalRecipe with
// parametric curve values to NP3 generation
func TestParametricCurveIntegration(t *testing.T) {
	// Create a recipe with parametric curve shadow lift (+50)
	recipe := &models.UniversalRecipe{
		ToneCurveShadows:    50,  // Lift shadows
		ToneCurveDarks:      25,  // Slight lift in darks
		ToneCurveLights:     0,   // No change in lights
		ToneCurveHighlights: -25, // Slight compression in highlights
	}

	params, err := convertToNP3Parameters(recipe)
	if err != nil {
		t.Fatalf("convertToNP3Parameters failed: %v", err)
	}

	// Verify that tone curve points were generated
	if params.toneCurvePointCount == 0 {
		t.Error("Expected tone curve points to be generated from parametric curve")
	}

	if len(params.toneCurvePoints) < 6 {
		t.Errorf("Expected at least 6 tone curve points, got %d", len(params.toneCurvePoints))
	}

	// First point should show shadow lift (value2 > value1)
	if len(params.toneCurvePoints) > 0 {
		firstPoint := params.toneCurvePoints[0]
		if firstPoint.value2 <= firstPoint.value1 {
			t.Errorf("Expected shadow lift (Y > X), got (%d, %d)", firstPoint.value1, firstPoint.value2)
		}
	}
}
