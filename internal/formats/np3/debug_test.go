package np3

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDebugOriginalFiles examines raw bytes in original NP3 files
func TestDebugOriginalFiles(t *testing.T) {
	// Load one sample file
	samplePath := filepath.Join("..", "..", "..", "examples", "np3", "Denis Zeqiri", "Classic Chrome.np3")
	data, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("Read sample failed: %v", err)
	}

	t.Logf("File size: %d bytes", len(data))

	// Check bytes 66-70 (sharpness range)
	t.Logf("\nSharpness bytes (66-70):")
	for i := 66; i <= 70; i++ {
		if i < len(data) {
			t.Logf("  Offset %d: 0x%02X (%d)", i, data[i], data[i])
		}
	}

	// Check bytes 71-75 (brightness range)
	t.Logf("\nBrightness bytes (71-75):")
	for i := 71; i <= 75; i++ {
		if i < len(data) {
			t.Logf("  Offset %d: 0x%02X (%d)", i, data[i], data[i])
		}
	}

	// Check bytes 76-79 (hue range)
	t.Logf("\nHue bytes (76-79):")
	for i := 76; i <= 79; i++ {
		if i < len(data) {
			t.Logf("  Offset %d: 0x%02X (%d)", i, data[i], data[i])
		}
	}

	// Parse the file to see what parameters we get
	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	t.Logf("\nParsed UniversalRecipe values:")
	t.Logf("  Sharpness: %d", recipe.Sharpness)
	t.Logf("  Contrast: %d", recipe.Contrast)
	t.Logf("  Saturation: %d", recipe.Saturation)
	t.Logf("  Exposure (brightness): %f", recipe.Exposure)
}
