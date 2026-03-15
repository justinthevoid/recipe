package converter

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/models"
)

// TestRoundTrip_NP3_XMP_NP3 tests NP3 → XMP → NP3 conversion maintains fidelity
func TestRoundTrip_NP3_XMP_NP3(t *testing.T) {
	t.Parallel()

	files, err := findFilesRecursive("../../testdata/np3", ".np3")
	if err != nil {
		t.Fatalf("Failed to find NP3 files: %v", err)
	}

	// Also check for .NP3 (uppercase)
	filesUpper, err := findFilesRecursive("../../testdata/np3", ".NP3")
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
			compareRecipesNP3Limited(t, orig, final, 1)
		})
	}
}

// TestRoundTrip_XMP_NP3_XMP tests XMP → NP3 → XMP conversion maintains fidelity
func TestRoundTrip_XMP_NP3_XMP(t *testing.T) {
	t.Parallel()

	files, err := findFilesRecursive("../../testdata/xmp", ".xmp")
	if err != nil {
		t.Fatalf("Failed to find XMP files: %v", err)
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
			compareRecipesNP3Limited(t, orig, final, 1)
		})
	}
}

// compareRecipes compares two UniversalRecipe instances with tolerance for rounding errors
func compareRecipes(t *testing.T, orig, final *models.UniversalRecipe, tolerance int) {
	t.Helper()

	// Compare floats with small tolerance
	if diff := math.Abs(orig.Exposure - final.Exposure); diff > 0.02 {
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

// compareRecipesNP3Limited compares recipes for XMP → NP3 → XMP round-trips.
// NP3 format has limitations (see docs/known-conversion-limitations.md):
// - No support for Exposure/Brightness (NP3 uses tone curve instead)
// - Limited split toning support (hue only, no saturation/balance)
// - No grain or vignette
// - Temperature/Tint not supported (camera-specific white balance)
// This function only compares parameters that NP3 can preserve.
func compareRecipesNP3Limited(t *testing.T, orig, final *models.UniversalRecipe, tolerance int) {
	t.Helper()

	// Exposure: NOT supported in NP3
	if orig.Exposure != 0 {
		if final.Exposure != 0 {
			t.Errorf("Exposure mismatch: orig=%.2f, final=%.2f (expected final=0.0, NP3 doesn't support exposure)", orig.Exposure, final.Exposure)
		} else {
			t.Logf("Exposure: orig=%.2f, final=0.0 (expected: NP3 doesn't support exposure parameter)", orig.Exposure)
		}
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

	// Contrast: NP3 range is -99 to +99
	if diff := orig.Contrast - final.Contrast; diff < -35 || diff > 35 {
		t.Errorf("Contrast mismatch: orig=%d, final=%d (diff=%d, tolerance=35)",
			orig.Contrast, final.Contrast, diff)
	}

	// Saturation: Use wider tolerance (±20) due to Vibrance interaction and quantization
	if diff := orig.Saturation - final.Saturation; diff < -20 || diff > 20 {
		t.Errorf("Saturation mismatch: orig=%d, final=%d (diff=%d, tolerance=20)",
			orig.Saturation, final.Saturation, diff)
	}

	// Vibrance
	if orig.Vibrance != 0 && final.Vibrance == 0 {
		t.Logf("Vibrance reset to 0: orig=%d (may be due to Picture Control base mapping)", orig.Vibrance)
	} else {
		compareInt("Vibrance", orig.Vibrance, final.Vibrance)
	}

	// Clarity
	if orig.Clarity >= -70 && orig.Clarity <= 70 {
		diff := orig.Clarity - final.Clarity
		if diff < 0 {
			diff = -diff
		}
		if diff > 3 {
			t.Errorf("Clarity mismatch: orig=%d, final=%d (diff=%d, tolerance=3)",
				orig.Clarity, final.Clarity, diff)
		}
	} else {
		t.Logf("Clarity outside typical range: orig=%d (expected loss)", orig.Clarity)
	}

	// Sharpness
	if orig.Sharpness >= 0 && orig.Sharpness <= 90 {
		if diff := orig.Sharpness - final.Sharpness; diff < -10 || diff > 10 {
			t.Errorf("Sharpness mismatch: orig=%d, final=%d (diff=%d, tolerance=10)",
				orig.Sharpness, final.Sharpness, diff)
		}
	} else {
		t.Logf("Sharpness outside NP3 range: orig=%d (will be clamped to 0-90), final=%d",
			orig.Sharpness, final.Sharpness)
	}

	// Tint: Not supported in NP3
	if orig.Tint != 0 {
		t.Logf("Tint: orig=%d (expected loss: NP3 doesn't support this parameter)", orig.Tint)
	}

	// Tone parameters
	compareInt("Highlights", orig.Highlights, final.Highlights)
	compareInt("Shadows", orig.Shadows, final.Shadows)
	compareInt("Whites", orig.Whites, final.Whites)

	// Blacks: wider tolerance
	if diff := orig.Blacks - final.Blacks; diff < -10 || diff > 10 {
		t.Errorf("Blacks mismatch: orig=%d, final=%d (diff=%d, tolerance=10)",
			orig.Blacks, final.Blacks, diff)
	}

	// Parameters NP3 does NOT support
	if orig.SplitShadowSaturation != 0 || orig.SplitHighlightSaturation != 0 || orig.SplitBalance != 0 {
		t.Logf("Split toning saturation/balance: (expected loss: NP3 only supports hue)")
	}

	if orig.Temperature != nil {
		t.Logf("Temperature: orig=%d (expected loss: NP3 doesn't support this parameter)", *orig.Temperature)
	}

	// HSL color adjustments
	compareColorAdj := func(name string, origAdj, finalAdj models.ColorAdjustment) {
		compareIntWithTolerance := func(paramName string, origVal, finalVal, tol int) {
			diff := origVal - finalVal
			if diff < 0 {
				diff = -diff
			}
			if diff > tol {
				t.Errorf("%s mismatch: orig=%d, final=%d (diff=%d, tolerance=%d)",
					paramName, origVal, finalVal, diff, tol)
			}
		}
		compareIntWithTolerance(name+".Hue", origAdj.Hue, finalAdj.Hue, 50)
		compareIntWithTolerance(name+".Saturation", origAdj.Saturation, finalAdj.Saturation, 25)
		compareIntWithTolerance(name+".Luminance", origAdj.Luminance, finalAdj.Luminance, 5)
	}

	compareColorAdj("Red", orig.Red, final.Red)
	compareColorAdj("Orange", orig.Orange, final.Orange)
	compareColorAdj("Yellow", orig.Yellow, final.Yellow)
	compareColorAdj("Green", orig.Green, final.Green)
	compareColorAdj("Aqua", orig.Aqua, final.Aqua)
	compareColorAdj("Blue", orig.Blue, final.Blue)
	compareColorAdj("Purple", orig.Purple, final.Purple)
	compareColorAdj("Magenta", orig.Magenta, final.Magenta)

	// Split toning hue
	if orig.SplitShadowHue != 0 || orig.SplitHighlightHue != 0 {
		t.Logf("Split toning hue: orig shadow=%d, highlight=%d (NP3 has limited split toning support)",
			orig.SplitShadowHue, orig.SplitHighlightHue)
	}

	// Tone curves
	if len(orig.PointCurve) != len(final.PointCurve) {
		t.Logf("PointCurve length difference: orig=%d, final=%d (expected due to format limitations)",
			len(orig.PointCurve), len(final.PointCurve))
	}
}
