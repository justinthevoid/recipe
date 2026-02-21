package np3

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestOverlapRegion examines bytes in the overlap region (150-235)
func TestOverlapRegion(t *testing.T) {
	// Create recipe with Saturation=66, Contrast=66
	recipe, err := models.NewRecipeBuilder().
		WithName("Test").
		WithSaturation(66).
		WithContrast(66).
		Build()
	if err != nil {
		t.Fatalf("Build recipe failed: %v", err)
	}

	t.Logf("Recipe parameters:")
	t.Logf("  Saturation: %d → params.saturation: %d", recipe.Saturation, recipe.Saturation/33)
	t.Logf("  Contrast: %d → params.contrast: %d", recipe.Contrast, recipe.Contrast/33)

	// Calculate expected data
	satParams := recipe.Saturation / 33
	colorTriplets := (satParams + 1) * 15
	colorEndOffset := 100 + (colorTriplets * 3)
	
	contrastParams := recipe.Contrast / 33
	toneCurvePairs := (contrastParams + 2) * 20
	
	t.Logf("\nExpected generation:")
	t.Logf("  Color triplets: %d (offset 100-%d)", colorTriplets, colorEndOffset-1)
	t.Logf("  Tone curve pairs: %d (offset 200-%d)", toneCurvePairs, 200+(toneCurvePairs*2)-1)

	// Generate
	data, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Examine overlap region 150-235
	t.Logf("\nBytes in overlap region (150-235):")
	nonZeroCount := 0
	for i := 150; i <= 235 && i < len(data); i++ {
		if data[i] != 0 {
			t.Logf("  Offset %d: 0x%02X (%d)", i, data[i], data[i])
			nonZeroCount++
		}
	}
	t.Logf("Total non-zero bytes in overlap: %d", nonZeroCount)
	
	// Count tone curve pairs as parser would
	toneCurveCount := 0
	for i := 150; i < 500 && i+1 < len(data); i += 2 {
		if data[i] != 0 || data[i+1] != 0 {
			toneCurveCount++
		}
	}
	t.Logf("\nTone curve pairs counted by parser: %d", toneCurveCount)
	t.Logf("Expected tone curve pairs: %d", toneCurvePairs)
	
	curveComplexity := toneCurveCount / 20
	parsedContrast := curveComplexity - 2
	t.Logf("\nParsed params.contrast: %d (expected %d)", parsedContrast, contrastParams)
	t.Logf("Final Contrast value: %d (expected %d)", parsedContrast*33, recipe.Contrast)
}
