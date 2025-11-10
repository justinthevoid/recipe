package np3

import (
	"os"
	"path/filepath"
	"testing"
)

// TestContrastDebug examines contrast calculation for failing file
func TestContrastDebug(t *testing.T) {
	// Load the Life.np3 file that has contrast mismatch
	samplePath := filepath.Join("..", "..", "..", "examples", "np3", "Denis Zeqiri", "Life.np3")
	originalData, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("Read sample failed: %v", err)
	}

	// Parse original
	originalRecipe, err := Parse(originalData)
	if err != nil {
		t.Fatalf("Parse original failed: %v", err)
	}

	t.Logf("Original recipe:")
	t.Logf("  Contrast: %d", originalRecipe.Contrast)
	t.Logf("  params.contrast should be: %d", originalRecipe.Contrast/33)

	// Generate from original
	generatedData, err := Generate(originalRecipe)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check tone curve bytes in generated data
	t.Logf("\nGenerated tone curve bytes (200-210):")
	for i := 200; i <= 210; i++ {
		if i < len(generatedData) {
			t.Logf("  Offset %d: 0x%02X (%d)", i, generatedData[i], generatedData[i])
		}
	}

	// Parse generated
	roundTripRecipe, err := Parse(generatedData)
	if err != nil {
		t.Fatalf("Parse generated failed: %v", err)
	}

	t.Logf("\nRound-trip recipe:")
	t.Logf("  Contrast: %d (expected %d)", roundTripRecipe.Contrast, originalRecipe.Contrast)
}
