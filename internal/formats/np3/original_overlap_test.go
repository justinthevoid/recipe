package np3

import (
	"os"
	"path/filepath"
	"testing"
)

// TestOriginalOverlap examines what original files have in overlap region
func TestOriginalOverlap(t *testing.T) {
	// Load one original file
	samplePath := filepath.Join("..", "..", "..", "examples", "np3", "Denis Zeqiri", "Life.np3")
	data, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Parse to see parameters
	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	t.Logf("Original parameters:")
	t.Logf("  Saturation: %d → params: %d", recipe.Saturation, recipe.Saturation/33)
	t.Logf("  Contrast: %d → params: %d", recipe.Contrast, recipe.Contrast/33)

	// Show bytes in overlap region 150-235
	t.Logf("\nOriginal bytes in overlap region (150-235):")
	nonZeroCount := 0
	for i := 150; i <= 235 && i < len(data); i++ {
		if data[i] != 0 {
			if nonZeroCount < 20 {  // Show first 20 only
				t.Logf("  Offset %d: 0x%02X (%d)", i, data[i], data[i])
			}
			nonZeroCount++
		}
	}
	t.Logf("Total non-zero bytes in overlap: %d", nonZeroCount)
}
