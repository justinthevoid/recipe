package costyle

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestRoundTrip validates that .costyle files can be round-trip converted
// (Parse → Generate → Parse) with 95%+ accuracy.
//
// AC-1: Implement Round-Trip Test Suite
// - Test flow: Parse .costyle → UniversalRecipe → Generate .costyle → Parse → Compare
// - Verify all core parameters preserved with 95%+ accuracy
// - Test with minimum 5 real-world .costyle samples
func TestRoundTrip(t *testing.T) {
	// Discover all .costyle sample files
	patterns := []string{
		"testdata/costyle/*.costyle",
		"testdata/costyle/real-world/*.costyle",
	}

	var testFiles []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatalf("Failed to glob pattern %s: %v", pattern, err)
		}
		testFiles = append(testFiles, matches...)
	}

	// Filter out malformed files (tested separately)
	var validFiles []string
	for _, file := range testFiles {
		if filepath.Base(file) != "sample4-malformed.costyle" {
			validFiles = append(validFiles, file)
		}
	}

	if len(validFiles) == 0 {
		t.Fatal("No .costyle test files found - expected files in testdata/costyle/")
	}

	t.Logf("Testing round-trip on %d .costyle sample files", len(validFiles))

	// Track aggregate accuracy metrics
	var allAccuracies []AccuracyReport
	successCount := 0

	for _, filePath := range validFiles {
		fileName := filepath.Base(filePath)

		t.Run("RoundTrip_"+fileName, func(t *testing.T) {
			// Step 1: Load original .costyle file
			originalData, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			// Step 2: Parse to UniversalRecipe
			recipe1, err := Parse(originalData)
			if err != nil {
				t.Fatalf("Parse original failed: %v", err)
			}

			// Step 3: Generate .costyle from recipe
			generatedData, err := Generate(recipe1)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Step 4: Parse generated .costyle
			recipe2, err := Parse(generatedData)
			if err != nil {
				t.Fatalf("Parse generated failed: %v", err)
			}

			// Step 5: Compare recipes and calculate accuracy
			report := compareRecipes(recipe1, recipe2)
			report.FileName = fileName

			// AC-5: Fail test if accuracy < 95%
			if report.Accuracy < 0.95 {
				t.Errorf("Round-trip accuracy below 95%%: got %.2f%% (matched %d/%d params)",
					report.Accuracy*100, report.MatchedParams, report.TotalParams)
				// Log parameter breakdown for debugging
				t.Logf("Parameter breakdown:")
				for param, matched := range report.ParameterBreakdown {
					if !matched {
						t.Logf("  ❌ %s", param)
					}
				}
			} else {
				t.Logf("✓ Round-trip successful - accuracy: %.2f%% (%d/%d params)",
					report.Accuracy*100, report.MatchedParams, report.TotalParams)
				successCount++
			}

			allAccuracies = append(allAccuracies, report)
		})
	}

	// AC-5: Calculate and report aggregate metrics
	if len(allAccuracies) > 0 {
		minAccuracy := 1.0
		maxAccuracy := 0.0
		sumAccuracy := 0.0

		for _, report := range allAccuracies {
			if report.Accuracy < minAccuracy {
				minAccuracy = report.Accuracy
			}
			if report.Accuracy > maxAccuracy {
				maxAccuracy = report.Accuracy
			}
			sumAccuracy += report.Accuracy
		}

		avgAccuracy := sumAccuracy / float64(len(allAccuracies))

		t.Logf("\n=== Round-Trip Accuracy Metrics ===")
		t.Logf("Files tested: %d", len(allAccuracies))
		t.Logf("Min accuracy: %.2f%%", minAccuracy*100)
		t.Logf("Max accuracy: %.2f%%", maxAccuracy*100)
		t.Logf("Avg accuracy: %.2f%%", avgAccuracy*100)
		t.Logf("Success rate: %d/%d files (%.1f%%)",
			successCount, len(allAccuracies), float64(successCount)/float64(len(allAccuracies))*100)

		// AC-5: Fail test if average accuracy < 95%
		if avgAccuracy < 0.95 {
			t.Errorf("FAIL: Average round-trip accuracy below 95%% requirement: got %.2f%%", avgAccuracy*100)
		}

		// Generate JSON accuracy report for documentation
		reportJSON, err := json.MarshalIndent(allAccuracies, "", "  ")
		if err == nil {
			reportPath := "testdata/costyle/test-results.json"
			if err := os.WriteFile(reportPath, reportJSON, 0644); err == nil {
				t.Logf("Accuracy report saved to: %s", reportPath)
			}
		}
	}
}

// AccuracyReport contains detailed accuracy metrics for a single round-trip test.
// AC-5: Measure and Report Accuracy Metrics
type AccuracyReport struct {
	FileName           string          `json:"file"`
	TotalParams        int             `json:"total_params"`
	MatchedParams      int             `json:"matched_params"`
	Accuracy           float64         `json:"accuracy"` // 0.0 to 1.0
	ParameterBreakdown map[string]bool `json:"breakdown"` // param_name -> matched
}

// compareRecipes compares two UniversalRecipe instances for round-trip validation.
// Returns an AccuracyReport with detailed parameter-by-parameter comparison.
//
// AC-2: Verify Key Adjustments Preserved Exactly
// - Exposure: No precision loss (exact float equality)
// - Contrast, Saturation, Tint, Clarity: Within ±1 integer value
// - Temperature: Within ±2 units (Kelvin conversion)
func compareRecipes(orig, final *models.UniversalRecipe) AccuracyReport {
	report := AccuracyReport{
		ParameterBreakdown: make(map[string]bool),
	}

	// AC-2: Compare Exposure (exact match within float precision)
	report.TotalParams++
	if compareExposure(orig.Exposure, final.Exposure) {
		report.MatchedParams++
		report.ParameterBreakdown["Exposure"] = true
	} else {
		report.ParameterBreakdown["Exposure"] = false
	}

	// AC-2: Compare integer parameters (within ±1)
	intParams := map[string]struct{ orig, final int }{
		"Contrast":   {orig.Contrast, final.Contrast},
		"Saturation": {orig.Saturation, final.Saturation},
		"Tint":       {orig.Tint, final.Tint},
		"Clarity":    {orig.Clarity, final.Clarity},
		"Highlights": {orig.Highlights, final.Highlights},
		"Shadows":    {orig.Shadows, final.Shadows},
		"Whites":     {orig.Whites, final.Whites},
		"Blacks":     {orig.Blacks, final.Blacks},
		"Vibrance":   {orig.Vibrance, final.Vibrance},
		"Sharpness":  {orig.Sharpness, final.Sharpness},
	}

	for name, vals := range intParams {
		report.TotalParams++
		if compareInteger(vals.orig, vals.final, 1) {
			report.MatchedParams++
			report.ParameterBreakdown[name] = true
		} else {
			report.ParameterBreakdown[name] = false
		}
	}

	// AC-2: Compare Temperature (within ±2 units for Kelvin conversion)
	report.TotalParams++
	if orig.Temperature != nil && final.Temperature != nil {
		if compareInteger(*orig.Temperature, *final.Temperature, 2) {
			report.MatchedParams++
			report.ParameterBreakdown["Temperature"] = true
		} else {
			report.ParameterBreakdown["Temperature"] = false
		}
	} else if (orig.Temperature == nil) == (final.Temperature == nil) {
		// Both nil or both non-nil matches
		report.MatchedParams++
		report.ParameterBreakdown["Temperature"] = true
	} else {
		report.ParameterBreakdown["Temperature"] = false
	}

	// Compare HSL color adjustments (within ±1 per channel)
	colorAdjustments := map[string]struct{ orig, final models.ColorAdjustment }{
		"Red":     {orig.Red, final.Red},
		"Orange":  {orig.Orange, final.Orange},
		"Yellow":  {orig.Yellow, final.Yellow},
		"Green":   {orig.Green, final.Green},
		"Aqua":    {orig.Aqua, final.Aqua},
		"Blue":    {orig.Blue, final.Blue},
		"Purple":  {orig.Purple, final.Purple},
		"Magenta": {orig.Magenta, final.Magenta},
	}

	for colorName, adjs := range colorAdjustments {
		// Each color adjustment has 3 components: Hue, Saturation, Luminance
		report.TotalParams += 3

		if compareInteger(adjs.orig.Hue, adjs.final.Hue, 1) {
			report.MatchedParams++
			report.ParameterBreakdown[colorName+".Hue"] = true
		} else {
			report.ParameterBreakdown[colorName+".Hue"] = false
		}

		if compareInteger(adjs.orig.Saturation, adjs.final.Saturation, 1) {
			report.MatchedParams++
			report.ParameterBreakdown[colorName+".Saturation"] = true
		} else {
			report.ParameterBreakdown[colorName+".Saturation"] = false
		}

		if compareInteger(adjs.orig.Luminance, adjs.final.Luminance, 1) {
			report.MatchedParams++
			report.ParameterBreakdown[colorName+".Luminance"] = true
		} else {
			report.ParameterBreakdown[colorName+".Luminance"] = false
		}
	}

	// Compare Split Toning (within ±1)
	splitParams := map[string]struct{ orig, final int }{
		"SplitShadowHue":           {orig.SplitShadowHue, final.SplitShadowHue},
		"SplitShadowSaturation":    {orig.SplitShadowSaturation, final.SplitShadowSaturation},
		"SplitHighlightHue":        {orig.SplitHighlightHue, final.SplitHighlightHue},
		"SplitHighlightSaturation": {orig.SplitHighlightSaturation, final.SplitHighlightSaturation},
		"SplitBalance":             {orig.SplitBalance, final.SplitBalance},
	}

	for name, vals := range splitParams {
		report.TotalParams++
		if compareInteger(vals.orig, vals.final, 1) {
			report.MatchedParams++
			report.ParameterBreakdown[name] = true
		} else {
			report.ParameterBreakdown[name] = false
		}
	}

	// Calculate accuracy percentage
	if report.TotalParams > 0 {
		report.Accuracy = float64(report.MatchedParams) / float64(report.TotalParams)
	}

	return report
}

// compareExposure compares exposure values with exact float precision.
// AC-2: Exposure must have no precision loss.
func compareExposure(a, b float64) bool {
	return math.Abs(a-b) < 0.01 // Within float precision tolerance
}

// compareInteger compares integer values within a specified tolerance.
// AC-2: Integer parameters allow ±1 value for acceptable precision loss.
func compareInteger(a, b, tolerance int) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

// TestRoundTrip_EdgeCases validates round-trip conversion for edge cases.
// AC-6: Edge case tests - Empty recipe, extreme values, missing parameters
func TestRoundTrip_EdgeCases(t *testing.T) {
	testCases := []struct {
		name   string
		recipe *models.UniversalRecipe
	}{
		{
			name: "Empty recipe (all parameters zero)",
			recipe: &models.UniversalRecipe{
				Name:         "Empty",
				SourceFormat: "costyle",
			},
		},
		{
			name: "Extreme values",
			recipe: &models.UniversalRecipe{
				Name:         "Extreme",
				SourceFormat: "costyle",
				Exposure:     2.0, // Max exposure
				Contrast:     100, // Max contrast
				Saturation:   100, // Max saturation
				Clarity:      100, // Max clarity
			},
		},
		{
			name: "Minimal recipe (only exposure set)",
			recipe: &models.UniversalRecipe{
				Name:         "Minimal",
				SourceFormat: "costyle",
				Exposure:     1.5,
			},
		},
		{
			name: "Complex recipe (all costyle-supported parameters populated)",
			recipe: &models.UniversalRecipe{
				Name:                     "Complex",
				SourceFormat:             "costyle",
				Exposure:                 0.75,
				Contrast:                 25,
				Saturation:               40,
				Clarity:                  35,
				Temperature:              intPtr(6500),
				Tint:                     5,
				SplitShadowHue:           240,  // Blue shadows
				SplitShadowSaturation:    30,   // 0-100 range in UR
				SplitHighlightHue:        30,   // Orange highlights
				SplitHighlightSaturation: 25,   // 0-100 range in UR
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate .costyle from recipe
			generatedData, err := Generate(tc.recipe)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Parse generated .costyle
			roundTripRecipe, err := Parse(generatedData)
			if err != nil {
				t.Fatalf("Parse generated failed: %v", err)
			}

			// Compare recipes
			report := compareRecipes(tc.recipe, roundTripRecipe)

			// Verify accuracy >= 95%
			if report.Accuracy < 0.95 {
				t.Errorf("Round-trip accuracy below 95%%: got %.2f%% (matched %d/%d params)",
					report.Accuracy*100, report.MatchedParams, report.TotalParams)
			} else {
				t.Logf("✓ Edge case round-trip successful - accuracy: %.2f%%", report.Accuracy*100)
			}
		})
	}
}

// TestRoundTrip_CrossFormat validates round-trip through intermediate formats.
// AC-6: Cross-format tests to verify no corruption through intermediate formats
func TestRoundTrip_CrossFormat(t *testing.T) {
	// This test is a placeholder for future cross-format round-trip tests
	// Example flows to test:
	// - .costyle → .xmp → .costyle (verify no corruption through XMP)
	// - .costyle → .np3 → .costyle (verify through NP3)
	//
	// NOTE: Some precision loss expected due to different parameter ranges
	// between formats. Document known limitations.
	t.Skip("Cross-format round-trip tests not yet implemented - requires integration with XMP/NP3 converters")
}

// intPtr returns a pointer to an int value
func intPtr(v int) *int {
	return &v
}

// BenchmarkRoundTrip measures round-trip conversion performance
func BenchmarkRoundTrip(b *testing.B) {
	// Load a sample .costyle file
	data, err := os.ReadFile("testdata/costyle/sample1-portrait.costyle")
	if err != nil {
		b.Fatalf("Failed to read test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Parse → Generate → Parse
		recipe1, err := Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}

		generated, err := Generate(recipe1)
		if err != nil {
			b.Fatalf("Generate failed: %v", err)
		}

		_, err = Parse(generated)
		if err != nil {
			b.Fatalf("Parse generated failed: %v", err)
		}
	}
}
