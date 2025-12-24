//go:build ignore
// +build ignore

package np3

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

// ARCHIVED: These tests are for curve generation functionality that is no longer used.
// See sprint-change-proposal-2025-12-24.md for rationale.

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

// ===== RGB Curve Merging Tests =====

func TestMergeRGBCurvesToMaster_IdenticalCurves(t *testing.T) {
	// Simulates Agfachrome RSX 200: identical RGB curves (0,10), (71,64), (188,194), (255,242)
	rgbCurve := []ControlPoint{
		{X: 0, Y: 10},
		{X: 71, Y: 64},
		{X: 188, Y: 194},
		{X: 255, Y: 242},
	}

	masterCurve := []ControlPoint{
		{X: 0, Y: 13},
		{X: 65, Y: 65},
		{X: 255, Y: 255},
	}

	result := MergeRGBCurvesToMaster(masterCurve, rgbCurve, rgbCurve, rgbCurve)

	if len(result) == 0 {
		t.Fatal("Expected merged curve to have points")
	}

	// Verify first point is lifted (RGB curve lifts blacks from 0 to 10)
	if result[0].Y < 5 {
		t.Errorf("Expected lifted blacks at first point, got Y=%d", result[0].Y)
	}

	// Verify last point is compressed (RGB curve compresses whites to 242)
	lastIdx := len(result) - 1
	if result[lastIdx].Y > 250 {
		t.Errorf("Expected compressed whites at last point, got Y=%d", result[lastIdx].Y)
	}

	t.Logf("Merged curve has %d points:", len(result))
	for i, p := range result {
		t.Logf("  Point %d: (%d, %d)", i, p.X, p.Y)
	}
}

func TestMergeRGBCurvesToMaster_NoRGBCurves(t *testing.T) {
	masterCurve := []ControlPoint{
		{X: 0, Y: 0},
		{X: 128, Y: 140},
		{X: 255, Y: 255},
	}

	result := MergeRGBCurvesToMaster(masterCurve, nil, nil, nil)

	if len(result) != len(masterCurve) {
		t.Errorf("Expected %d points, got %d", len(masterCurve), len(result))
	}
}

func TestCurvesAreIdentical(t *testing.T) {
	curve1 := []ControlPoint{{X: 0, Y: 10}, {X: 128, Y: 140}, {X: 255, Y: 255}}
	curve2 := []ControlPoint{{X: 0, Y: 10}, {X: 128, Y: 140}, {X: 255, Y: 255}}
	curve3 := []ControlPoint{{X: 0, Y: 11}, {X: 128, Y: 140}, {X: 255, Y: 255}} // Slightly different

	if !CurvesAreIdentical(curve1, curve2, curve2, 2) {
		t.Error("Expected identical curves to be recognized")
	}

	if !CurvesAreIdentical(curve1, curve2, curve3, 2) {
		t.Error("Expected curves to be identical within tolerance=2")
	}

	if CurvesAreIdentical(curve1, curve2, curve3, 0) {
		t.Error("Expected curves to be different with tolerance=0")
	}
}

func TestPointsToCurveLUT(t *testing.T) {
	points := []ControlPoint{{X: 0, Y: 10}, {X: 255, Y: 240}}

	lut := PointsToCurveLUT(points)

	if lut[0] != 10 {
		t.Errorf("Expected lut[0]=10, got %d", lut[0])
	}

	if lut[255] != 240 {
		t.Errorf("Expected lut[255]=240, got %d", lut[255])
	}

	expectedMid := (10 + 240) / 2
	if absInt(lut[127]-expectedMid) > 2 {
		t.Errorf("Expected lut[127]≈%d, got %d", expectedMid, lut[127])
	}
}

func TestToneCurvePointsToControlPoints(t *testing.T) {
	points := []models.ToneCurvePoint{
		{Input: 0, Output: 10},
		{Input: 128, Output: 140},
		{Input: 255, Output: 255},
	}

	result := toneCurvePointsToControlPoints(points)

	if len(result) != 3 {
		t.Fatalf("Expected 3 points, got %d", len(result))
	}

	if result[0].X != 0 || result[0].Y != 10 {
		t.Errorf("First point mismatch: got (%d,%d)", result[0].X, result[0].Y)
	}
}

func TestRGBCurveIntegration(t *testing.T) {
	// Test the full flow: XMP recipe with RGB curves → NP3 with merged curve
	recipe := &models.UniversalRecipe{
		PointCurve: []models.ToneCurvePoint{
			{Input: 0, Output: 13},
			{Input: 65, Output: 65},
			{Input: 255, Output: 255},
		},
		PointCurveRed: []models.ToneCurvePoint{
			{Input: 0, Output: 10},
			{Input: 71, Output: 64},
			{Input: 188, Output: 194},
			{Input: 255, Output: 242},
		},
		PointCurveGreen: []models.ToneCurvePoint{
			{Input: 0, Output: 10},
			{Input: 71, Output: 64},
			{Input: 188, Output: 194},
			{Input: 255, Output: 242},
		},
		PointCurveBlue: []models.ToneCurvePoint{
			{Input: 0, Output: 10},
			{Input: 71, Output: 64},
			{Input: 188, Output: 194},
			{Input: 255, Output: 242},
		},
	}

	params, err := convertToNP3Parameters(recipe)
	if err != nil {
		t.Fatalf("convertToNP3Parameters failed: %v", err)
	}

	// Verify that tone curve was generated
	if params.toneCurvePointCount == 0 {
		t.Error("Expected tone curve points from merged RGB curves")
	}

	t.Logf("Generated %d tone curve points from RGB curve merging", params.toneCurvePointCount)

	// First point should have lifted blacks
	if len(params.toneCurvePoints) > 0 {
		first := params.toneCurvePoints[0]
		t.Logf("First point: (%d, %d)", first.value1, first.value2)
	}
}
