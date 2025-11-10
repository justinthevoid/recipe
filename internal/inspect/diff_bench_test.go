package inspect

import (
	"os"
	"testing"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/models"
)

// BenchmarkDiff_CompareOnly benchmarks just the comparison logic (no I/O)
func BenchmarkDiff_CompareOnly(b *testing.B) {
	recipe1 := &models.UniversalRecipe{
		Name:       "Test Recipe 1",
		Exposure:   0.5,
		Contrast:   15,
		Saturation: -10,
		Highlights: -50,
		Shadows:    30,
		Red: models.ColorAdjustment{
			Hue:        10,
			Saturation: 5,
			Luminance:  0,
		},
	}

	recipe2 := &models.UniversalRecipe{
		Name:       "Test Recipe 2",
		Exposure:   0.55,
		Contrast:   20,
		Saturation: -5,
		Highlights: -50, // Same
		Shadows:    35,
		Red: models.ColorAdjustment{
			Hue:        15,
			Saturation: 5,
			Luminance:  0,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Diff(recipe1, recipe2, 0.001)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDiff_EndToEnd benchmarks full diff including file I/O and parsing
func BenchmarkDiff_EndToEnd(b *testing.B) {
	// Use real sample files if available
	file1 := "./../../examples/np3/Denis Zeqiri/Classic Chrome.np3"
	file2 := "./../../examples/np3/Denis Zeqiri/Filmic.np3"

	// Check if files exist, skip if not
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		b.Skip("Sample files not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Read files
		data1, err := os.ReadFile(file1)
		if err != nil {
			b.Fatal(err)
		}

		data2, err := os.ReadFile(file2)
		if err != nil {
			b.Fatal(err)
		}

		// Parse
		recipe1, err := np3.Parse(data1)
		if err != nil {
			b.Fatal(err)
		}

		recipe2, err := np3.Parse(data2)
		if err != nil {
			b.Fatal(err)
		}

		// Diff
		_, err = Diff(recipe1, recipe2, 0.001)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDiff_CrossFormat benchmarks NP3 vs XMP diff
func BenchmarkDiff_CrossFormat(b *testing.B) {
	np3File := "./../../examples/lrtemplate/015. PRESETPRO - Emulation K/00. E - auto tone.np3"
	xmpFile := "./../../examples/lrtemplate/015. PRESETPRO - Emulation K/00. E - auto tone.xmp"

	// Check if files exist, skip if not
	if _, err := os.Stat(np3File); os.IsNotExist(err) {
		b.Skip("Sample files not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Read files
		np3Data, err := os.ReadFile(np3File)
		if err != nil {
			b.Fatal(err)
		}

		xmpData, err := os.ReadFile(xmpFile)
		if err != nil {
			b.Fatal(err)
		}

		// Parse
		recipe1, err := np3.Parse(np3Data)
		if err != nil {
			b.Fatal(err)
		}

		recipe2, err := xmp.Parse(xmpData)
		if err != nil {
			b.Fatal(err)
		}

		// Diff
		_, err = Diff(recipe1, recipe2, 0.001)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFormatDiff benchmarks text output formatting
func BenchmarkFormatDiff(b *testing.B) {
	results := []DiffResult{
		{Field: "Contrast", OldValue: int64(0), NewValue: int64(15), ChangeType: "modified", Significant: true},
		{Field: "Saturation", OldValue: int64(0), NewValue: int64(-10), ChangeType: "modified", Significant: true},
		{Field: "Exposure", OldValue: 0.5, NewValue: 0.51, ChangeType: "modified", Significant: false},
		{Field: "Vibrance", OldValue: nil, NewValue: int64(20), ChangeType: "added", Significant: false},
		{Field: "Highlights", OldValue: int64(-50), NewValue: int64(-50), ChangeType: "unchanged", Significant: false},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatDiff(results, false, false)
	}
}

// BenchmarkFormatDiff_Unified benchmarks unified mode formatting
func BenchmarkFormatDiff_Unified(b *testing.B) {
	results := make([]DiffResult, 0, 50)

	// Add many unchanged fields to simulate real scenario
	for i := 0; i < 45; i++ {
		results = append(results, DiffResult{
			Field:       "Field" + string(rune('A'+i)),
			OldValue:    0,
			NewValue:    0,
			ChangeType:  "unchanged",
			Significant: false,
		})
	}

	// Add a few changes
	results = append(results,
		DiffResult{Field: "Contrast", OldValue: int64(0), NewValue: int64(15), ChangeType: "modified", Significant: true},
		DiffResult{Field: "Saturation", OldValue: int64(0), NewValue: int64(-10), ChangeType: "modified", Significant: true},
		DiffResult{Field: "Exposure", OldValue: 0.5, NewValue: 0.51, ChangeType: "modified", Significant: false},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatDiff(results, true, false) // Unified mode
	}
}

// BenchmarkDiff_LargeRecipes benchmarks diff with all fields populated
func BenchmarkDiff_LargeRecipes(b *testing.B) {
	recipe1 := &models.UniversalRecipe{
		Name:          "Large Recipe 1",
		SourceFormat:  "np3",
		Exposure:      0.5,
		Contrast:      15,
		Highlights:    -50,
		Shadows:       30,
		Whites:        20,
		Blacks:        -20,
		Texture:       10,
		Clarity:       15,
		Dehaze:        5,
		Vibrance:      20,
		Saturation:    -10,
		Sharpness:     40,
		SharpnessRadius: 1.0,
		SharpnessDetail: 25,
		SharpnessMasking: 0,
		Tint:          5,
		Red:           models.ColorAdjustment{Hue: 10, Saturation: 5, Luminance: 0},
		Orange:        models.ColorAdjustment{Hue: 5, Saturation: 10, Luminance: -5},
		Yellow:        models.ColorAdjustment{Hue: 0, Saturation: 0, Luminance: 0},
		Green:         models.ColorAdjustment{Hue: -5, Saturation: 5, Luminance: 10},
		Aqua:          models.ColorAdjustment{Hue: 0, Saturation: 0, Luminance: 0},
		Blue:          models.ColorAdjustment{Hue: 10, Saturation: -5, Luminance: 5},
		Purple:        models.ColorAdjustment{Hue: 0, Saturation: 0, Luminance: 0},
		Magenta:       models.ColorAdjustment{Hue: 5, Saturation: 10, Luminance: 0},
		GrainAmount:   25,
		GrainSize:     50,
		GrainRoughness: 50,
		VignetteAmount: -20,
		VignetteMidpoint: 50,
		VignetteRoundness: 0,
		VignetteFeather: 50,
	}

	recipe2 := &models.UniversalRecipe{
		Name:          "Large Recipe 2",
		SourceFormat:  "xmp",
		Exposure:      0.55,
		Contrast:      20,
		Highlights:    -45,
		Shadows:       35,
		Whites:        20,
		Blacks:        -20,
		Texture:       15,
		Clarity:       20,
		Dehaze:        10,
		Vibrance:      25,
		Saturation:    -5,
		Sharpness:     45,
		SharpnessRadius: 1.2,
		SharpnessDetail: 30,
		SharpnessMasking: 5,
		Tint:          10,
		Red:           models.ColorAdjustment{Hue: 15, Saturation: 10, Luminance: 5},
		Orange:        models.ColorAdjustment{Hue: 10, Saturation: 15, Luminance: 0},
		Yellow:        models.ColorAdjustment{Hue: 5, Saturation: 5, Luminance: 5},
		Green:         models.ColorAdjustment{Hue: 0, Saturation: 10, Luminance: 15},
		Aqua:          models.ColorAdjustment{Hue: 5, Saturation: 5, Luminance: 5},
		Blue:          models.ColorAdjustment{Hue: 15, Saturation: 0, Luminance: 10},
		Purple:        models.ColorAdjustment{Hue: 5, Saturation: 5, Luminance: 5},
		Magenta:       models.ColorAdjustment{Hue: 10, Saturation: 15, Luminance: 5},
		GrainAmount:   30,
		GrainSize:     55,
		GrainRoughness: 55,
		VignetteAmount: -25,
		VignetteMidpoint: 55,
		VignetteRoundness: 5,
		VignetteFeather: 55,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Diff(recipe1, recipe2, 0.001)
		if err != nil {
			b.Fatal(err)
		}
	}
}
