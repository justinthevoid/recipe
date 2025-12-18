//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/models"
)

// TestCase defines a parametric curve configuration for testing
type TestCase struct {
	Name       string
	Shadows    int
	Darks      int
	Lights     int
	Highlights int
}

func main() {
	// Define test cases for Phase 2 validation
	testCases := []TestCase{
		{Name: "Parametric_Shadow_Lift_50", Shadows: 50, Darks: 0, Lights: 0, Highlights: 0},
		{Name: "Parametric_Highlight_Compress", Shadows: 0, Darks: 0, Lights: 0, Highlights: -50},
		{Name: "Parametric_Classic_SCurve", Shadows: 30, Darks: 15, Lights: -15, Highlights: -30},
		{Name: "Parametric_All_Lift_25", Shadows: 25, Darks: 25, Lights: 25, Highlights: 25},
		{Name: "Parametric_Dramatic_SCurve", Shadows: 50, Darks: 25, Lights: -25, Highlights: -50},
	}

	outputDir := "internal/formats/np3/testdata/curve_tests"

	fmt.Println("=== Phase 2 Parametric Curve Test Generator ===")
	fmt.Println()

	for _, tc := range testCases {
		// Create UniversalRecipe with parametric curve values
		recipe := &models.UniversalRecipe{
			Name:                tc.Name,
			ToneCurveShadows:    tc.Shadows,
			ToneCurveDarks:      tc.Darks,
			ToneCurveLights:     tc.Lights,
			ToneCurveHighlights: tc.Highlights,
			// Default other params
			Sharpness:  50,
			Contrast:   0,
			Saturation: 0,
		}

		// Generate NP3 with warnings
		data, warnings, err := np3.GenerateWithWarnings(recipe)
		if err != nil {
			fmt.Printf("❌ ERROR generating %s: %v\n", tc.Name, err)
			continue
		}

		// Show what control points were generated
		points := np3.ParametricToControlPoints(tc.Shadows, tc.Darks, tc.Lights, tc.Highlights, 0, 0, 0)
		fmt.Printf("📊 %s\n", tc.Name)
		fmt.Printf("   Parametric: Shadows=%+d, Darks=%+d, Lights=%+d, Highlights=%+d\n",
			tc.Shadows, tc.Darks, tc.Lights, tc.Highlights)
		fmt.Printf("   Control Points (%d):\n", len(points))
		for i, p := range points {
			delta := p.Y - p.X
			indicator := "="
			if delta > 0 {
				indicator = "↑"
			} else if delta < 0 {
				indicator = "↓"
			}
			fmt.Printf("      [%d] X=%3d → Y=%3d  (%s%+d)\n", i, p.X, p.Y, indicator, delta)
		}

		// Save the file
		filename := filepath.Join(outputDir, tc.Name+".np3")
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Printf("   ❌ Failed to write: %v\n", err)
			continue
		}
		fmt.Printf("   ✅ Saved: %s (%d bytes)\n", filename, len(data))

		// Show any warnings
		if warnings.HasWarnings() {
			for _, w := range warnings.Warnings {
				fmt.Printf("   ⚠️  Warning: %s - %s\n", w.Parameter, w.Message)
			}
		}
		fmt.Println()
	}

	fmt.Println("=== Test Files Ready for NX Studio Validation ===")
	fmt.Println()
	fmt.Println("Open each .np3 file in NX Studio Picture Control Utility 2 and verify:")
	fmt.Println("1. Custom Tone Curve checkbox is enabled")
	fmt.Println("2. Curve shape matches expected adjustments:")
	fmt.Println("   - Shadow_Lift_50: Curve should lift shadow region (left side up)")
	fmt.Println("   - Highlight_Compress: Curve should compress highlights (right side down)")
	fmt.Println("   - Classic_SCurve: S-shaped curve (shadows up, highlights down)")
	fmt.Println("   - All_Lift_25: Curve parallel above diagonal")
	fmt.Println("   - Dramatic_SCurve: Pronounced S-curve for high contrast look")
}
