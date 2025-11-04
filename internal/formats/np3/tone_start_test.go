package np3

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestToneCurveStart verifies where tone curve data actually starts
func TestToneCurveStart(t *testing.T) {
	recipe, err := models.NewRecipeBuilder().
		WithName("Test").
		WithSaturation(66).
		WithContrast(66).
		Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	data, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Find where tone curve data (value 1) starts
	t.Logf("Searching for tone curve data (value 1)...")
	foundStart := -1
	for i := 150; i < 500 && i < len(data); i++ {
		if data[i] == 1 {
			foundStart = i
			t.Logf("First tone curve byte (value 1) at offset %d", i)
			break
		}
	}

	if foundStart == -1 {
		t.Errorf("No tone curve data found!")
	}

	// Show bytes around the transition
	if foundStart > 0 {
		t.Logf("\nBytes around tone curve start (%d-%d):", foundStart-5, foundStart+10)
		for i := foundStart - 5; i < foundStart+10 && i < len(data); i++ {
			if i >= 0 {
				t.Logf("  Offset %d: 0x%02X (%d)", i, data[i], data[i])
			}
		}
	}
}
