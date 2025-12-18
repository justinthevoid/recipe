package inspect

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/formats/lrtemplate"
	"github.com/justin/recipe/internal/models"
)

// BenchmarkToJSON measures JSON serialization performance (AC-7)
// Target: <5ms for typical UniversalRecipe
func BenchmarkToJSON(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Name:          "Benchmark Test",
		Exposure:      0.5,
		Contrast:      15,
		Highlights:    -20,
		Shadows:       10,
		Whites:        5,
		Blacks:        -5,
		Texture:       3,
		Clarity:       10,
		Dehaze:        5,
		Vibrance:      15,
		Saturation:    -8,
		Sharpness:     30,
		SharpnessRadius: 1.2,
		Tint:          -3,
		Red:           models.ColorAdjustment{Hue: 5, Saturation: 10, Luminance: 3},
		Orange:        models.ColorAdjustment{Hue: 10, Saturation: 15, Luminance: 5},
		Yellow:        models.ColorAdjustment{Hue: 0, Saturation: 5, Luminance: 2},
		Green:         models.ColorAdjustment{Hue: -5, Saturation: 8, Luminance: 1},
		Aqua:          models.ColorAdjustment{Hue: 3, Saturation: 12, Luminance: 4},
		Blue:          models.ColorAdjustment{Hue: -3, Saturation: 6, Luminance: 2},
		Purple:        models.ColorAdjustment{Hue: 7, Saturation: 14, Luminance: 3},
		Magenta:       models.ColorAdjustment{Hue: 2, Saturation: 9, Luminance: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ToJSON(recipe)
		if err != nil {
			b.Fatalf("ToJSON failed: %v", err)
		}
	}
}

// BenchmarkToJSONWithMetadata measures full output generation (AC-7)
// Target: <5ms including metadata wrapper
func BenchmarkToJSONWithMetadata(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Name:       "Benchmark Test",
		Exposure:   0.5,
		Contrast:   15,
		Highlights: -20,
		Shadows:    10,
		Saturation: -5,
		Vibrance:   10,
		Clarity:    5,
		Sharpness:  25,
	}

	sourceFile := "portrait.np3"
	format := "np3"
	version := "2.0.0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ToJSONWithMetadata(recipe, sourceFile, format, version)
		if err != nil {
			b.Fatalf("ToJSONWithMetadata failed: %v", err)
		}
	}
}

// BenchmarkInspectEndToEnd measures complete workflow (AC-7)
// Target: <50ms for file read + parse + JSON (actual target is <100ms total)
func BenchmarkInspectEndToEnd_NP3(b *testing.B) {
	// Find a sample NP3 file
	testFile := findTestFile(b, "../../testdata/xmp/sample.np3")
	if testFile == "" {
		b.Skip("No sample NP3 file found for benchmark")
	}

	// Read file once for all iterations
	data, err := os.ReadFile(testFile)
	if err != nil {
		b.Fatalf("Failed to read test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Parse
		recipe, err := np3.Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}

		// Convert to JSON
		_, err = ToJSONWithMetadata(recipe, testFile, "np3", "2.0.0")
		if err != nil {
			b.Fatalf("ToJSONWithMetadata failed: %v", err)
		}
	}
}

// BenchmarkInspectEndToEnd_XMP measures XMP format performance
func BenchmarkInspectEndToEnd_XMP(b *testing.B) {
	testFile := findTestFile(b, "../../testdata/xmp/*.xmp")
	if testFile == "" {
		b.Skip("No sample XMP file found for benchmark")
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		b.Fatalf("Failed to read test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recipe, err := xmp.Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}

		_, err = ToJSONWithMetadata(recipe, testFile, "xmp", "2.0.0")
		if err != nil {
			b.Fatalf("ToJSONWithMetadata failed: %v", err)
		}
	}
}

// BenchmarkInspectEndToEnd_LRTemplate measures lrtemplate format performance
func BenchmarkInspectEndToEnd_LRTemplate(b *testing.B) {
	testFile := findTestFile(b, "../../testdata/lrtemplate/*.lrtemplate")
	if testFile == "" {
		b.Skip("No sample lrtemplate file found for benchmark")
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		b.Fatalf("Failed to read test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recipe, err := lrtemplate.Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}

		_, err = ToJSONWithMetadata(recipe, testFile, "lrtemplate", "2.0.0")
		if err != nil {
			b.Fatalf("ToJSONWithMetadata failed: %v", err)
		}
	}
}

// findTestFile finds first matching file for benchmarks
func findTestFile(b *testing.B, pattern string) string {
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}
	return matches[0]
}
