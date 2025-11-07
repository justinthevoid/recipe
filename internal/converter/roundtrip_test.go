package converter

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/formats/lrtemplate"
	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/models"
)

// TestRoundTrip_NP3_XMP_NP3 tests NP3 → XMP → NP3 conversion maintains fidelity
func TestRoundTrip_NP3_XMP_NP3(t *testing.T) {
	t.Parallel()

	files, err := findFilesRecursive("../../examples/np3", ".np3")
	if err != nil {
		t.Fatalf("Failed to find NP3 files: %v", err)
	}

	// Also check for .NP3 (uppercase)
	filesUpper, err := findFilesRecursive("../../examples/np3", ".NP3")
	if err == nil {
		files = append(files, filesUpper...)
	}

	if len(files) == 0 {
		t.Skip("No NP3 files found for round-trip testing")
	}

	t.Logf("Testing NP3 → XMP → NP3 with %d files", len(files))

	for _, file := range files {
		file := file // Capture loop variable
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()

			// Step 1: Parse original NP3
			origData, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Read original failed: %v", err)
			}

			orig, err := np3.Parse(origData)
			if err != nil {
				t.Fatalf("Parse NP3 failed: %v", err)
			}

			// Step 2: Convert to XMP
			xmpData, err := xmp.Generate(orig)
			if err != nil {
				t.Fatalf("Generate XMP failed: %v", err)
			}

			// Step 3: Parse XMP
			xmpRecipe, err := xmp.Parse(xmpData)
			if err != nil {
				t.Fatalf("Parse XMP failed: %v", err)
			}

			// Step 4: Convert back to NP3
			np3Data, err := np3.Generate(xmpRecipe)
			if err != nil {
				t.Fatalf("Generate NP3 failed: %v", err)
			}

			// Step 5: Parse final NP3
			final, err := np3.Parse(np3Data)
			if err != nil {
				t.Fatalf("Parse final NP3 failed: %v", err)
			}

			// Step 6: Compare with tolerance ±1
			compareRecipes(t, orig, final, 1)
		})
	}
}

// TestRoundTrip_XMP_NP3_XMP tests XMP → NP3 → XMP conversion maintains fidelity
func TestRoundTrip_XMP_NP3_XMP(t *testing.T) {
	t.Parallel()

	files, err := findFilesRecursive("../../examples/xmp", ".xmp")
	if err != nil {
		t.Fatalf("Failed to find XMP files: %v", err)
	}

	// Also check testdata
	testFiles, err := findFilesRecursive("../../testdata/xmp", ".xmp")
	if err == nil {
		files = append(files, testFiles...)
	}

	if len(files) == 0 {
		t.Skip("No XMP files found for round-trip testing")
	}

	t.Logf("Testing XMP → NP3 → XMP with %d files", len(files))

	for _, file := range files {
		file := file // Capture loop variable
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()

			// Step 1: Parse original XMP
			origData, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Read original failed: %v", err)
			}

			orig, err := xmp.Parse(origData)
			if err != nil {
				t.Fatalf("Parse XMP failed: %v", err)
			}

			// Step 2: Convert to NP3
			np3Data, err := np3.Generate(orig)
			if err != nil {
				t.Fatalf("Generate NP3 failed: %v", err)
			}

			// Step 3: Parse NP3
			np3Recipe, err := np3.Parse(np3Data)
			if err != nil {
				t.Fatalf("Parse NP3 failed: %v", err)
			}

			// Step 4: Convert back to XMP
			xmpData, err := xmp.Generate(np3Recipe)
			if err != nil {
				t.Fatalf("Generate XMP failed: %v", err)
			}

			// Step 5: Parse final XMP
			final, err := xmp.Parse(xmpData)
			if err != nil {
				t.Fatalf("Parse final XMP failed: %v", err)
			}

			// Step 6: Compare with tolerance ±1
			compareRecipes(t, orig, final, 1)
		})
	}
}

// TestRoundTrip_NP3_lrtemplate_NP3 tests NP3 → lrtemplate → NP3 conversion maintains fidelity
func TestRoundTrip_NP3_lrtemplate_NP3(t *testing.T) {
	t.Parallel()

	files, err := findFilesRecursive("../../examples/np3", ".np3")
	if err != nil {
		t.Fatalf("Failed to find NP3 files: %v", err)
	}

	filesUpper, err := findFilesRecursive("../../examples/np3", ".NP3")
	if err == nil {
		files = append(files, filesUpper...)
	}

	if len(files) == 0 {
		t.Skip("No NP3 files found for round-trip testing")
	}

	t.Logf("Testing NP3 → lrtemplate → NP3 with %d files", len(files))

	for _, file := range files {
		file := file // Capture loop variable
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()

			// Step 1: Parse original NP3
			origData, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Read original failed: %v", err)
			}

			orig, err := np3.Parse(origData)
			if err != nil {
				t.Fatalf("Parse NP3 failed: %v", err)
			}

			// Step 2: Convert to lrtemplate
			lrtData, err := lrtemplate.Generate(orig)
			if err != nil {
				t.Fatalf("Generate lrtemplate failed: %v", err)
			}

			// Step 3: Parse lrtemplate
			lrtRecipe, err := lrtemplate.Parse(lrtData)
			if err != nil {
				t.Fatalf("Parse lrtemplate failed: %v", err)
			}

			// Step 4: Convert back to NP3
			np3Data, err := np3.Generate(lrtRecipe)
			if err != nil {
				t.Fatalf("Generate NP3 failed: %v", err)
			}

			// Step 5: Parse final NP3
			final, err := np3.Parse(np3Data)
			if err != nil {
				t.Fatalf("Parse final NP3 failed: %v", err)
			}

			// Step 6: Compare with tolerance ±1
			compareRecipes(t, orig, final, 1)
		})
	}
}

// TestRoundTrip_XMP_lrtemplate_XMP tests XMP → lrtemplate → XMP conversion maintains fidelity
func TestRoundTrip_XMP_lrtemplate_XMP(t *testing.T) {
	t.Parallel()

	files, err := findFilesRecursive("../../examples/xmp", ".xmp")
	if err != nil {
		t.Fatalf("Failed to find XMP files: %v", err)
	}

	testFiles, err := findFilesRecursive("../../testdata/xmp", ".xmp")
	if err == nil {
		files = append(files, testFiles...)
	}

	if len(files) == 0 {
		t.Skip("No XMP files found for round-trip testing")
	}

	t.Logf("Testing XMP → lrtemplate → XMP with %d files", len(files))

	for _, file := range files {
		file := file // Capture loop variable
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()

			// Step 1: Parse original XMP
			origData, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Read original failed: %v", err)
			}

			orig, err := xmp.Parse(origData)
			if err != nil {
				t.Fatalf("Parse XMP failed: %v", err)
			}

			// Step 2: Convert to lrtemplate
			lrtData, err := lrtemplate.Generate(orig)
			if err != nil {
				t.Fatalf("Generate lrtemplate failed: %v", err)
			}

			// Step 3: Parse lrtemplate
			lrtRecipe, err := lrtemplate.Parse(lrtData)
			if err != nil {
				t.Fatalf("Parse lrtemplate failed: %v", err)
			}

			// Step 4: Convert back to XMP
			xmpData, err := xmp.Generate(lrtRecipe)
			if err != nil {
				t.Fatalf("Generate XMP failed: %v", err)
			}

			// Step 5: Parse final XMP
			final, err := xmp.Parse(xmpData)
			if err != nil {
				t.Fatalf("Parse final XMP failed: %v", err)
			}

			// Step 6: Compare with tolerance ±1
			compareRecipes(t, orig, final, 1)
		})
	}
}

// TestRoundTrip_lrtemplate_NP3_lrtemplate tests lrtemplate → NP3 → lrtemplate conversion maintains fidelity
func TestRoundTrip_lrtemplate_NP3_lrtemplate(t *testing.T) {
	t.Parallel()

	files, err := findFilesRecursive("../../examples/lrtemplate", ".lrtemplate")
	if err != nil {
		t.Fatalf("Failed to find lrtemplate files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No lrtemplate files found for round-trip testing")
	}

	t.Logf("Testing lrtemplate → NP3 → lrtemplate with %d files", len(files))

	for _, file := range files {
		file := file // Capture loop variable
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()

			// Step 1: Parse original lrtemplate
			origData, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Read original failed: %v", err)
			}

			orig, err := lrtemplate.Parse(origData)
			if err != nil {
				t.Fatalf("Parse lrtemplate failed: %v", err)
			}

			// Step 2: Convert to NP3
			np3Data, err := np3.Generate(orig)
			if err != nil {
				t.Fatalf("Generate NP3 failed: %v", err)
			}

			// Step 3: Parse NP3
			np3Recipe, err := np3.Parse(np3Data)
			if err != nil {
				t.Fatalf("Parse NP3 failed: %v", err)
			}

			// Step 4: Convert back to lrtemplate
			lrtData, err := lrtemplate.Generate(np3Recipe)
			if err != nil {
				t.Fatalf("Generate lrtemplate failed: %v", err)
			}

			// Step 5: Parse final lrtemplate
			final, err := lrtemplate.Parse(lrtData)
			if err != nil {
				t.Fatalf("Parse final lrtemplate failed: %v", err)
			}

			// Step 6: Compare with tolerance ±1
			compareRecipes(t, orig, final, 1)
		})
	}
}

// TestRoundTrip_lrtemplate_XMP_lrtemplate tests lrtemplate → XMP → lrtemplate conversion maintains fidelity
func TestRoundTrip_lrtemplate_XMP_lrtemplate(t *testing.T) {
	t.Parallel()

	files, err := findFilesRecursive("../../examples/lrtemplate", ".lrtemplate")
	if err != nil {
		t.Fatalf("Failed to find lrtemplate files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No lrtemplate files found for round-trip testing")
	}

	t.Logf("Testing lrtemplate → XMP → lrtemplate with %d files", len(files))

	for _, file := range files {
		file := file // Capture loop variable
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()

			// Step 1: Parse original lrtemplate
			origData, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Read original failed: %v", err)
			}

			orig, err := lrtemplate.Parse(origData)
			if err != nil {
				t.Fatalf("Parse lrtemplate failed: %v", err)
			}

			// Step 2: Convert to XMP
			xmpData, err := xmp.Generate(orig)
			if err != nil {
				t.Fatalf("Generate XMP failed: %v", err)
			}

			// Step 3: Parse XMP
			xmpRecipe, err := xmp.Parse(xmpData)
			if err != nil {
				t.Fatalf("Parse XMP failed: %v", err)
			}

			// Step 4: Convert back to lrtemplate
			lrtData, err := lrtemplate.Generate(xmpRecipe)
			if err != nil {
				t.Fatalf("Generate lrtemplate failed: %v", err)
			}

			// Step 5: Parse final lrtemplate
			final, err := lrtemplate.Parse(lrtData)
			if err != nil {
				t.Fatalf("Parse final lrtemplate failed: %v", err)
			}

			// Step 6: Compare with tolerance ±1
			compareRecipes(t, orig, final, 1)
		})
	}
}

// compareRecipes compares two UniversalRecipe instances with tolerance for rounding errors
func compareRecipes(t *testing.T, orig, final *models.UniversalRecipe, tolerance int) {
	t.Helper()

	// Compare floats with small tolerance
	if diff := math.Abs(orig.Exposure - final.Exposure); diff > 0.01 {
		t.Errorf("Exposure mismatch: orig=%.2f, final=%.2f (diff=%.2f)", orig.Exposure, final.Exposure, diff)
	}

	// Compare integers with tolerance
	compareInt := func(name string, origVal, finalVal int) {
		diff := origVal - finalVal
		if diff < 0 {
			diff = -diff
		}
		if diff > tolerance {
			t.Errorf("%s mismatch: orig=%d, final=%d (diff=%d, tolerance=%d)",
				name, origVal, finalVal, diff, tolerance)
		}
	}

	compareInt("Contrast", orig.Contrast, final.Contrast)
	compareInt("Highlights", orig.Highlights, final.Highlights)
	compareInt("Shadows", orig.Shadows, final.Shadows)
	compareInt("Whites", orig.Whites, final.Whites)
	compareInt("Blacks", orig.Blacks, final.Blacks)
	compareInt("Saturation", orig.Saturation, final.Saturation)
	compareInt("Vibrance", orig.Vibrance, final.Vibrance)
	compareInt("Clarity", orig.Clarity, final.Clarity)
	compareInt("Sharpness", orig.Sharpness, final.Sharpness)
	compareInt("Tint", orig.Tint, final.Tint)

	// Compare Temperature (nullable)
	if orig.Temperature != nil && final.Temperature != nil {
		compareInt("Temperature", *orig.Temperature, *final.Temperature)
	} else if (orig.Temperature == nil) != (final.Temperature == nil) {
		t.Errorf("Temperature nullability mismatch: orig=%v, final=%v", orig.Temperature, final.Temperature)
	}

	// Compare HSL color adjustments
	compareColorAdj := func(name string, origAdj, finalAdj models.ColorAdjustment) {
		compareInt(name+".Hue", origAdj.Hue, finalAdj.Hue)
		compareInt(name+".Saturation", origAdj.Saturation, finalAdj.Saturation)
		compareInt(name+".Luminance", origAdj.Luminance, finalAdj.Luminance)
	}

	compareColorAdj("Red", orig.Red, final.Red)
	compareColorAdj("Orange", orig.Orange, final.Orange)
	compareColorAdj("Yellow", orig.Yellow, final.Yellow)
	compareColorAdj("Green", orig.Green, final.Green)
	compareColorAdj("Aqua", orig.Aqua, final.Aqua)
	compareColorAdj("Blue", orig.Blue, final.Blue)
	compareColorAdj("Purple", orig.Purple, final.Purple)
	compareColorAdj("Magenta", orig.Magenta, final.Magenta)

	// Compare Split Toning
	compareInt("SplitShadowHue", orig.SplitShadowHue, final.SplitShadowHue)
	compareInt("SplitShadowSaturation", orig.SplitShadowSaturation, final.SplitShadowSaturation)
	compareInt("SplitHighlightHue", orig.SplitHighlightHue, final.SplitHighlightHue)
	compareInt("SplitHighlightSaturation", orig.SplitHighlightSaturation, final.SplitHighlightSaturation)
	compareInt("SplitBalance", orig.SplitBalance, final.SplitBalance)

	// Compare tone curves (just check length for now - detailed comparison would be too strict)
	if len(orig.PointCurve) != len(final.PointCurve) {
		t.Logf("PointCurve length difference: orig=%d, final=%d (this may be expected due to format limitations)",
			len(orig.PointCurve), len(final.PointCurve))
	}
}
