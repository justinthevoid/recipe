package costyle

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestParse validates the .costyle parser against all sample files following Pattern 7
func TestParse(t *testing.T) {
	files, err := filepath.Glob("testdata/costyle/*.costyle")
	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatal("no test files found in testdata/costyle/")
	}

	// Filter out malformed file (tested separately)
	var validFiles []string
	for _, file := range files {
		if filepath.Base(file) != "sample4-malformed.costyle" {
			validFiles = append(validFiles, file)
		}
	}

	for _, file := range validFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("failed to read %s: %v", file, err)
			}

			recipe, err := Parse(data)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			// Validate SourceFormat set correctly
			if recipe.SourceFormat != "costyle" {
				t.Errorf("SourceFormat = %s, want 'costyle'", recipe.SourceFormat)
			}

			// Validate critical field ranges
			if recipe.Exposure < -5.0 || recipe.Exposure > 5.0 {
				t.Errorf("Exposure out of range: %.2f", recipe.Exposure)
			}

			if recipe.Contrast < -100 || recipe.Contrast > 100 {
				t.Errorf("Contrast out of range: %d", recipe.Contrast)
			}

			if recipe.Saturation < -100 || recipe.Saturation > 100 {
				t.Errorf("Saturation out of range: %d", recipe.Saturation)
			}

			if recipe.Tint < -150 || recipe.Tint > 150 {
				t.Errorf("Tint out of range: %d", recipe.Tint)
			}

			if recipe.Clarity < -100 || recipe.Clarity > 100 {
				t.Errorf("Clarity out of range: %d", recipe.Clarity)
			}

			// Validate split toning ranges
			if recipe.SplitShadowHue < 0 || recipe.SplitShadowHue > 360 {
				t.Errorf("SplitShadowHue out of range: %d", recipe.SplitShadowHue)
			}

			if recipe.SplitShadowSaturation < 0 || recipe.SplitShadowSaturation > 100 {
				t.Errorf("SplitShadowSaturation out of range: %d", recipe.SplitShadowSaturation)
			}

			if recipe.SplitHighlightHue < 0 || recipe.SplitHighlightHue > 360 {
				t.Errorf("SplitHighlightHue out of range: %d", recipe.SplitHighlightHue)
			}

			if recipe.SplitHighlightSaturation < 0 || recipe.SplitHighlightSaturation > 100 {
				t.Errorf("SplitHighlightSaturation out of range: %d", recipe.SplitHighlightSaturation)
			}

			// Validate metadata map exists
			if recipe.Metadata == nil {
				t.Error("Metadata map is nil")
			}
		})
	}
}

// TestParse_ValidFile tests parsing a complete .costyle file with all parameters
func TestParse_ValidFile(t *testing.T) {
	data, err := os.ReadFile("testdata/costyle/sample1-portrait.costyle")
	if err != nil {
		t.Fatalf("failed to read sample1-portrait.costyle: %v", err)
	}

	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify core adjustments extracted correctly
	if recipe.Exposure != 0.7 {
		t.Errorf("Exposure = %.2f, want 0.7", recipe.Exposure)
	}

	if recipe.Contrast != 15 {
		t.Errorf("Contrast = %d, want 15", recipe.Contrast)
	}

	if recipe.Saturation != 10 {
		t.Errorf("Saturation = %d, want 10", recipe.Saturation)
	}

	if recipe.Clarity != 20 {
		t.Errorf("Clarity = %d, want 20", recipe.Clarity)
	}

	// Tint: C1 -3 * (150/100) = -4.5 ≈ -5 (rounded)
	if recipe.Tint != -5 {
		t.Errorf("Tint = %d, want -5 (scaled from C1 -3)", recipe.Tint)
	}

	// Verify temperature converted to Kelvin offset
	if recipe.Temperature == nil {
		t.Error("Temperature is nil, want non-nil")
	} else {
		// Temperature 5 should map to ~175K offset (5 * 35)
		expectedTemp := 175
		if *recipe.Temperature != expectedTemp {
			t.Errorf("Temperature = %d, want %d", *recipe.Temperature, expectedTemp)
		}
	}

	// Verify metadata preserved
	if recipe.Name != "Portrait Warm" {
		t.Errorf("Name = %s, want 'Portrait Warm'", recipe.Name)
	}

	if recipe.Metadata["costyle_author"] != "Recipe Test Suite" {
		t.Errorf("costyle_author = %v, want 'Recipe Test Suite'", recipe.Metadata["costyle_author"])
	}

	if recipe.Metadata["costyle_description"] != "Warm portrait preset with enhanced clarity" {
		t.Errorf("costyle_description = %v", recipe.Metadata["costyle_description"])
	}

	// Verify color balance extracted (shadows, highlights)
	// ShadowsHue: C1 30 → UR 30 (direct mapping, C1 uses 0-360 like UR)
	if recipe.SplitShadowHue != 30 {
		t.Errorf("SplitShadowHue = %d, want 30 (C1 hue directly compatible)", recipe.SplitShadowHue)
	}

	// ShadowsSaturation: C1 10 → UR (10 + 100) / 2 = 55
	if recipe.SplitShadowSaturation != 55 {
		t.Errorf("SplitShadowSaturation = %d, want 55 (scaled from C1 10)", recipe.SplitShadowSaturation)
	}

	// HighlightsHue: C1 60 → UR 60 (direct mapping)
	if recipe.SplitHighlightHue != 60 {
		t.Errorf("SplitHighlightHue = %d, want 60 (C1 hue directly compatible)", recipe.SplitHighlightHue)
	}

	// HighlightsSaturation: C1 8 → UR (8 + 100) / 2 = 54
	if recipe.SplitHighlightSaturation != 54 {
		t.Errorf("SplitHighlightSaturation = %d, want 54 (scaled from C1 8)", recipe.SplitHighlightSaturation)
	}

	// Verify midtones stored in metadata
	if recipe.Metadata["costyle_midtones_hue"] != 45 {
		t.Errorf("costyle_midtones_hue = %v, want 45", recipe.Metadata["costyle_midtones_hue"])
	}
}

// TestParse_MinimalFile tests parsing a minimal .costyle with only one adjustment
func TestParse_MinimalFile(t *testing.T) {
	data, err := os.ReadFile("testdata/costyle/sample2-minimal.costyle")
	if err != nil {
		t.Fatalf("failed to read sample2-minimal.costyle: %v", err)
	}

	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify only exposure set, others default to zero
	if recipe.Exposure != -0.3 {
		t.Errorf("Exposure = %.2f, want -0.3", recipe.Exposure)
	}

	if recipe.Contrast != 0 {
		t.Errorf("Contrast = %d, want 0", recipe.Contrast)
	}

	if recipe.Saturation != 0 {
		t.Errorf("Saturation = %d, want 0", recipe.Saturation)
	}

	if recipe.Clarity != 0 {
		t.Errorf("Clarity = %d, want 0", recipe.Clarity)
	}

	// Verify metadata exists (even if empty)
	if recipe.Metadata == nil {
		t.Error("Metadata map is nil")
	}
}

// TestParse_MalformedXML tests error handling for malformed XML
func TestParse_MalformedXML(t *testing.T) {
	data, err := os.ReadFile("testdata/costyle/sample4-malformed.costyle")
	if err != nil {
		t.Fatalf("failed to read sample4-malformed.costyle: %v", err)
	}

	_, err = Parse(data)
	if err == nil {
		t.Error("Parse() error = nil, want error for malformed XML")
	}

	// Verify error message includes context
	if err != nil && err.Error() == "" {
		t.Error("error message is empty")
	}
}

// TestParse_OutOfRangeValues tests parameter clamping for out-of-range values
func TestParse_OutOfRangeValues(t *testing.T) {
	tests := []struct {
		name      string
		costyle   string
		checkFunc func(*testing.T, *models.UniversalRecipe)
	}{
		{
			name: "Exposure too high",
			costyle: `<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:RDF>
    <rdf:Description>
      <Exposure>10.0</Exposure>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`,
			checkFunc: func(t *testing.T, r *models.UniversalRecipe) {
				if r.Exposure != 5.0 {
					t.Errorf("Exposure = %.2f, want 5.0 (clamped)", r.Exposure)
				}
			},
		},
		{
			name: "Exposure too low",
			costyle: `<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:RDF>
    <rdf:Description>
      <Exposure>-10.0</Exposure>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`,
			checkFunc: func(t *testing.T, r *models.UniversalRecipe) {
				if r.Exposure != -5.0 {
					t.Errorf("Exposure = %.2f, want -5.0 (clamped)", r.Exposure)
				}
			},
		},
		{
			name: "Contrast too high",
			costyle: `<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:RDF>
    <rdf:Description>
      <Contrast>200</Contrast>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`,
			checkFunc: func(t *testing.T, r *models.UniversalRecipe) {
				if r.Contrast != 100 {
					t.Errorf("Contrast = %d, want 100 (clamped)", r.Contrast)
				}
			},
		},
		{
			name: "Saturation too low",
			costyle: `<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:RDF>
    <rdf:Description>
      <Saturation>-200</Saturation>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`,
			checkFunc: func(t *testing.T, r *models.UniversalRecipe) {
				if r.Saturation != -100 {
					t.Errorf("Saturation = %d, want -100 (clamped)", r.Saturation)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := Parse([]byte(tt.costyle))
			if err != nil {
				t.Fatalf("Parse() error = %v, want nil", err)
			}

			tt.checkFunc(t, recipe)
		})
	}
}

// TestParse_ColorBalance tests color balance parameter extraction and mapping
func TestParse_ColorBalance(t *testing.T) {
	data, err := os.ReadFile("testdata/costyle/sample3-landscape.costyle")
	if err != nil {
		t.Fatalf("failed to read sample3-landscape.costyle: %v", err)
	}

	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify shadows color balance mapped
	if recipe.SplitShadowHue != 200 {
		t.Errorf("SplitShadowHue = %d, want 200", recipe.SplitShadowHue)
	}

	// Verify highlights color balance mapped
	if recipe.SplitHighlightHue != 210 {
		t.Errorf("SplitHighlightHue = %d, want 210", recipe.SplitHighlightHue)
	}

	// Verify midtones stored in metadata (no direct UniversalRecipe field)
	if recipe.Metadata["costyle_midtones_hue"] != 120 {
		t.Errorf("costyle_midtones_hue = %v, want 120", recipe.Metadata["costyle_midtones_hue"])
	}

	if recipe.Metadata["costyle_midtones_saturation"] != 15 {
		t.Errorf("costyle_midtones_saturation = %v, want 15", recipe.Metadata["costyle_midtones_saturation"])
	}
}

// BenchmarkParse measures parsing performance (target <100ms per AC)
func BenchmarkParse(b *testing.B) {
	data, err := os.ReadFile("testdata/costyle/sample1-portrait.costyle")
	if err != nil {
		b.Fatalf("failed to read sample1-portrait.costyle: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Fatalf("Parse() error = %v", err)
		}
	}
}

// TestParse_RoundTrip_Tint validates tint scaling preserves values through round-trip conversion.
// This test specifically validates the fix for Critical Issue #1 from code review.
func TestParse_RoundTrip_Tint(t *testing.T) {
	testCases := []int{-100, -75, -50, -25, 0, 25, 50, 75, 100}

	for _, originalTint := range testCases {
		t.Run(fmt.Sprintf("Tint_%d", originalTint), func(t *testing.T) {
			// Create .costyle XML with specific tint value
			costyleXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:RDF>
    <rdf:Description>
      <Tint>%d</Tint>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`, originalTint)

			// Parse → Generate → Parse round-trip
			recipe, err := Parse([]byte(costyleXML))
			if err != nil {
				t.Fatalf("Parse() failed: %v", err)
			}

			generatedData, err := Generate(recipe)
			if err != nil {
				t.Fatalf("Generate() failed: %v", err)
			}

			var finalStyle CaptureOneStyle
			if err := xml.Unmarshal(generatedData, &finalStyle); err != nil {
				t.Fatalf("Unmarshal generated XML failed: %v", err)
			}

			// Verify tint preserved (allow ±1 tolerance for rounding)
			finalTint := finalStyle.RDF.Description.Tint
			diff := abs(finalTint - originalTint)

			if diff > 1 {
				t.Errorf("Tint round-trip failed: %d → (UR: %d) → %d (error: %d)",
					originalTint, recipe.Tint, finalTint, diff)
			}
		})
	}
}

// TestParse_RoundTrip_ColorBalance validates color balance formulas preserve values through round-trip.
// This test specifically validates the fix for Critical Issue #2 from code review.
// Note: C1 uses 0-360 for hue (like UR), and -100 to +100 for saturation
func TestParse_RoundTrip_ColorBalance(t *testing.T) {
	testCases := []struct {
		name string
		hue  int // C1 hue: 0-360
		sat  int // C1 saturation: -100 to +100
	}{
		{"Min", 0, -100},      // 0° red, minimum saturation
		{"Mid", 180, 0},        // 180° cyan, neutral saturation
		{"Max", 360, 100},      // 360° red (wraps to 0), maximum saturation
		{"Negative Sat", 90, -50},  // 90° yellow-green, low saturation
		{"Positive Sat", 270, 50},  // 270° blue-violet, high saturation
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test shadows
			t.Run("Shadows", func(t *testing.T) {
				costyleXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:RDF>
    <rdf:Description>
      <ShadowsHue>%d</ShadowsHue>
      <ShadowsSaturation>%d</ShadowsSaturation>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`, tc.hue, tc.sat)

				// Parse → Generate → Parse round-trip
				recipe, err := Parse([]byte(costyleXML))
				if err != nil {
					t.Fatalf("Parse() failed: %v", err)
				}

				generatedData, err := Generate(recipe)
				if err != nil {
					t.Fatalf("Generate() failed: %v", err)
				}

				var finalStyle CaptureOneStyle
				if err := xml.Unmarshal(generatedData, &finalStyle); err != nil {
					t.Fatalf("Unmarshal generated XML failed: %v", err)
				}

				// Verify values preserved (allow ±2 tolerance for rounding)
				finalHue := finalStyle.RDF.Description.ShadowsHue
				finalSat := finalStyle.RDF.Description.ShadowsSaturation

				hueErr := abs(finalHue - tc.hue)
				satErr := abs(finalSat - tc.sat)

				if hueErr > 2 {
					t.Errorf("ShadowsHue round-trip failed: %d → (UR: %d) → %d (error: %d)",
						tc.hue, recipe.SplitShadowHue, finalHue, hueErr)
				}

				if satErr > 2 {
					t.Errorf("ShadowsSaturation round-trip failed: %d → (UR: %d) → %d (error: %d)",
						tc.sat, recipe.SplitShadowSaturation, finalSat, satErr)
				}
			})

			// Test highlights
			t.Run("Highlights", func(t *testing.T) {
				costyleXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:RDF>
    <rdf:Description>
      <HighlightsHue>%d</HighlightsHue>
      <HighlightsSaturation>%d</HighlightsSaturation>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`, tc.hue, tc.sat)

				// Parse → Generate → Parse round-trip
				recipe, err := Parse([]byte(costyleXML))
				if err != nil {
					t.Fatalf("Parse() failed: %v", err)
				}

				generatedData, err := Generate(recipe)
				if err != nil {
					t.Fatalf("Generate() failed: %v", err)
				}

				var finalStyle CaptureOneStyle
				if err := xml.Unmarshal(generatedData, &finalStyle); err != nil {
					t.Fatalf("Unmarshal generated XML failed: %v", err)
				}

				// Verify values preserved (allow ±2 tolerance for rounding)
				finalHue := finalStyle.RDF.Description.HighlightsHue
				finalSat := finalStyle.RDF.Description.HighlightsSaturation

				hueErr := abs(finalHue - tc.hue)
				satErr := abs(finalSat - tc.sat)

				if hueErr > 2 {
					t.Errorf("HighlightsHue round-trip failed: %d → (UR: %d) → %d (error: %d)",
						tc.hue, recipe.SplitHighlightHue, finalHue, hueErr)
				}

				if satErr > 2 {
					t.Errorf("HighlightsSaturation round-trip failed: %d → (UR: %d) → %d (error: %d)",
						tc.sat, recipe.SplitHighlightSaturation, finalSat, satErr)
				}
			})
		})
	}
}

// abs returns the absolute value of an int.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
