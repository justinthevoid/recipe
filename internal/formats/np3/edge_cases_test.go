package np3

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestGenerateBoundaryParameters tests valid extreme values at boundaries
func TestGenerateBoundaryParameters(t *testing.T) {
	tests := []struct {
		name       string
		sharpness  int
		contrast   int
		saturation int
		exposure   float64
	}{
		{
			name:       "Maximum valid values",
			sharpness:  150,
			contrast:   100,
			saturation: 100,
			exposure:   1.0,
		},
		{
			name:       "Minimum valid values",
			sharpness:  0,
			contrast:   -100,
			saturation: -100,
			exposure:   -1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := models.NewRecipeBuilder().
				WithName("BoundaryTest").
				WithSharpness(tt.sharpness).
				WithContrast(tt.contrast).
				WithSaturation(tt.saturation).
				WithExposure(tt.exposure).
				Build()
			if err != nil {
				t.Fatalf("Build recipe failed: %v", err)
			}

			data, err := Generate(recipe)
			if err != nil {
				t.Errorf("Generate failed: %v", err)
			}

			// Verify we can parse it back
			_, parseErr := Parse(data)
			if parseErr != nil {
				t.Errorf("Parse round-trip failed: %v", parseErr)
			}
		})
	}
}

// TestGenerateLongPresetName tests name truncation
func TestGenerateLongPresetName(t *testing.T) {
	// Create a name > 40 chars to trigger truncation (line 139-141)
	longName := "This is an extremely long preset name that exceeds forty characters"

	recipe, err := models.NewRecipeBuilder().
		WithName(longName).
		WithSharpness(50).
		Build()
	if err != nil {
		t.Fatalf("Build recipe failed: %v", err)
	}

	data, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Parse back and verify name was truncated
	parsedRecipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Name should be truncated to 40 chars
	if len(parsedRecipe.Name) > 40 {
		t.Errorf("Name not truncated: got %d chars, want ≤40", len(parsedRecipe.Name))
	}
}

// TestGenerateColorDataEdgeCases tests color data generation edge cases
func TestGenerateColorDataEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		saturation int
		desc       string
	}{
		{
			name:       "Negative saturation (triggers targetCount < 0)",
			saturation: -99,
			desc:       "Should handle negative saturation gracefully",
		},
		{
			name:       "Very high saturation (triggers max clamping)",
			saturation: 99,
			desc:       "Should clamp to maximum color triplets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := models.NewRecipeBuilder().
				WithName("ColorEdge").
				WithSaturation(tt.saturation).
				Build()
			if err != nil {
				t.Fatalf("Build recipe failed: %v", err)
			}

			data, err := Generate(recipe)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Should parse back successfully
			_, parseErr := Parse(data)
			if parseErr != nil {
				t.Errorf("Parse failed: %v", parseErr)
			}
		})
	}
}

// TestGenerateToneCurveEdgeCases tests tone curve generation edge cases
func TestGenerateToneCurveEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		contrast   int
		saturation int
		desc       string
	}{
		{
			name:       "Negative contrast with high saturation overlap",
			contrast:   -99,
			saturation: 99,
			desc:       "Should handle additionalPairs < 0",
		},
		{
			name:       "Very high contrast",
			contrast:   99,
			saturation: -99,
			desc:       "Should generate maximum tone curve pairs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := models.NewRecipeBuilder().
				WithName("ToneCurveEdge").
				WithContrast(tt.contrast).
				WithSaturation(tt.saturation).
				Build()
			if err != nil {
				t.Fatalf("Build recipe failed: %v", err)
			}

			data, err := Generate(recipe)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Should parse back successfully
			_, parseErr := Parse(data)
			if parseErr != nil {
				t.Errorf("Parse failed: %v", parseErr)
			}
		})
	}
}

// TestEstimateParametersEdgeCases tests heuristic estimation edge cases
func TestEstimateParametersEdgeCases(t *testing.T) {
	// Create data with extreme heuristic values to trigger all clamping paths
	data := make([]byte, 500)
	copy(data[0:3], magicBytes)
	copy(data[3:7], []byte{0x02, 0x10, 0x00, 0x00})

	// Fill sharpness bytes with extreme values (should trigger clamping in estimateParameters)
	for i := 66; i <= 70; i++ {
		data[i] = 255 // Max value
	}

	// Fill brightness bytes with extreme values
	for i := 71; i <= 75; i++ {
		data[i] = 255 // Max value
	}

	// Fill hue bytes with extreme values
	for i := 76; i <= 79; i++ {
		data[i] = 255 // Max value
	}

	// Fill color data region with extreme pattern
	for i := 100; i < 300; i += 3 {
		data[i] = 255
		data[i+1] = 255
		data[i+2] = 255
	}

	// Fill tone curve region with extreme pattern
	for i := 150; i < 500; i += 2 {
		data[i] = 255
		data[i+1] = 255
	}

	// Should parse without error despite extreme values
	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify clamping worked (should be within valid ranges)
	if recipe.Sharpness < 0 || recipe.Sharpness > 150 {
		t.Errorf("Sharpness out of range: %d", recipe.Sharpness)
	}
	if recipe.Contrast < -100 || recipe.Contrast > 100 {
		t.Errorf("Contrast out of range: %d", recipe.Contrast)
	}
	if recipe.Saturation < -100 || recipe.Saturation > 100 {
		t.Errorf("Saturation out of range: %d", recipe.Saturation)
	}
}

// TestBuildRecipeNullNameBytes tests name extraction with null bytes
func TestBuildRecipeNullNameBytes(t *testing.T) {
	// Create data with name containing null bytes mid-string
	data := make([]byte, 500)
	copy(data[0:3], magicBytes)
	copy(data[3:7], []byte{0x02, 0x10, 0x00, 0x00})

	// Write name with null byte in middle: "Test\x00More" (should stop at null)
	copy(data[20:60], []byte{'T', 'e', 's', 't', 0x00, 'M', 'o', 'r', 'e'})

	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Name should be "Test" (stops at null byte)
	if recipe.Name != "Test" {
		t.Errorf("Name parsing incorrect: got %q, want %q", recipe.Name, "Test")
	}
}
