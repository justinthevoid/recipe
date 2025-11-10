package np3

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestGeneratorOutput examines what bytes the generator produces
func TestGeneratorOutput(t *testing.T) {
	// Create a simple recipe with Sharpness=0, Saturation=33
	recipe, err := models.NewRecipeBuilder().
		WithName("Test").
		WithSharpness(0).
		WithSaturation(33).
		Build()
	if err != nil {
		t.Fatalf("Build recipe failed: %v", err)
	}

	// Generate NP3 data
	data, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	t.Logf("Generated file size: %d bytes", len(data))

	// Check sharpness bytes (66-70)
	t.Logf("\nGenerated Sharpness bytes (66-70):")
	for i := 66; i <= 70; i++ {
		if i < len(data) {
			t.Logf("  Offset %d: 0x%02X (%d)", i, data[i], data[i])
		}
	}

	// Check color data bytes (100-110 sample)
	t.Logf("\nGenerated Color data sample (100-110):")
	for i := 100; i <= 110; i++ {
		if i < len(data) {
			t.Logf("  Offset %d: 0x%02X (%d)", i, data[i], data[i])
		}
	}

	// Now parse it back
	parsedRecipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	t.Logf("\nParsed values:")
	t.Logf("  Sharpness: %d (expected 0)", parsedRecipe.Sharpness)
	t.Logf("  Saturation: %d (expected 33)", parsedRecipe.Saturation)
}
