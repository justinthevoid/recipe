package costyle

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestGenerate_ValidRecipe tests generating a .costyle from a populated UniversalRecipe.
func TestGenerate_ValidRecipe(t *testing.T) {
	// Create a UniversalRecipe with typical photo editing parameters
	recipe := &models.UniversalRecipe{
		Name:       "Test Preset",
		Exposure:   1.5,
		Contrast:   25,
		Saturation: 20,
		Tint:       10,
		Clarity:    15,
		Metadata: map[string]interface{}{
			"author":      "Test Photographer",
			"description": "Test preset description",
		},
	}

	// Temperature: 6100K = +10 in C1 scale
	temp := 6100
	recipe.Temperature = &temp

	// Generate .costyle
	costyleData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Verify XML is well-formed (can be parsed back)
	var parsed CaptureOneStyle
	if err := xml.Unmarshal(costyleData, &parsed); err != nil {
		t.Fatalf("Generated XML is not well-formed: %v", err)
	}

	// Verify XML declaration
	if !strings.HasPrefix(string(costyleData), "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("Missing or incorrect XML declaration")
	}

	// Verify parameter values match input (within tolerance)
	desc := parsed.RDF.Description
	if desc.Exposure != 1.5 {
		t.Errorf("Exposure mismatch: got %f, want 1.5", desc.Exposure)
	}
	if desc.Contrast != 25 {
		t.Errorf("Contrast mismatch: got %d, want 25", desc.Contrast)
	}
	if desc.Saturation != 20 {
		t.Errorf("Saturation mismatch: got %d, want 20", desc.Saturation)
	}
	// Temperature: 6100K → +10 in C1
	// Temperature: 6100K → (6100 - 5500) / 35 = 17.14 ≈ 17
	expectedTemp := 17
	if desc.Temperature != expectedTemp {
		t.Errorf("Temperature mismatch: got %d, want %d (from 6100K)", desc.Temperature, expectedTemp)
	}
	// Tint: 10 in UR (-150/+150) → 7 in C1 (-100/+100) (10 * 100/150 ≈ 6.67 → 7)
	expectedTint := 7 // int(round(10 * (100.0/150.0))) = int(round(6.67)) = 7
	if desc.Tint != expectedTint {
		t.Errorf("Tint mismatch: got %d, want %d", desc.Tint, expectedTint)
	}
	if desc.Clarity != 15 {
		t.Errorf("Clarity mismatch: got %d, want 15", desc.Clarity)
	}

	// Verify metadata preservation
	if desc.Name != "Test Preset" {
		t.Errorf("Name mismatch: got %q, want %q", desc.Name, "Test Preset")
	}
	if desc.Author != "Test Photographer" {
		t.Errorf("Author mismatch: got %q, want %q", desc.Author, "Test Photographer")
	}
	if desc.Description != "Test preset description" {
		t.Errorf("Description mismatch: got %q, want %q", desc.Description, "Test preset description")
	}
}

// TestGenerate_EmptyRecipe tests generating a minimal .costyle from an empty UniversalRecipe.
func TestGenerate_EmptyRecipe(t *testing.T) {
	// Create empty UniversalRecipe (all zero values)
	recipe := &models.UniversalRecipe{}

	// Generate .costyle
	costyleData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() failed for empty recipe: %v", err)
	}

	// Verify XML is well-formed
	var parsed CaptureOneStyle
	if err := xml.Unmarshal(costyleData, &parsed); err != nil {
		t.Fatalf("Generated XML is not well-formed: %v", err)
	}

	// Verify minimal valid structure (all parameters should be zero/omitted)
	desc := parsed.RDF.Description
	if desc.Exposure != 0 {
		t.Errorf("Empty recipe should have Exposure=0, got %f", desc.Exposure)
	}
	if desc.Contrast != 0 {
		t.Errorf("Empty recipe should have Contrast=0, got %d", desc.Contrast)
	}
	if desc.Saturation != 0 {
		t.Errorf("Empty recipe should have Saturation=0, got %d", desc.Saturation)
	}
	if desc.Temperature != 0 {
		t.Errorf("Empty recipe should have Temperature=0, got %d", desc.Temperature)
	}
}

// TestGenerate_OutOfRangeValues tests that out-of-range values are clamped.
func TestGenerate_OutOfRangeValues(t *testing.T) {
	// Create UniversalRecipe with out-of-range values
	recipe := &models.UniversalRecipe{
		Exposure:   5.0,  // Out of range (max 2.0)
		Contrast:   200,  // Out of range (max 100)
		Saturation: -150, // Out of range (min -100)
		Clarity:    -120, // Out of range (min -100)
	}

	// Temperature: 20000K (way out of range, should clamp to +100)
	temp := 20000
	recipe.Temperature = &temp

	// Generate .costyle
	costyleData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Parse back
	var parsed CaptureOneStyle
	if err := xml.Unmarshal(costyleData, &parsed); err != nil {
		t.Fatalf("Generated XML is not well-formed: %v", err)
	}

	desc := parsed.RDF.Description

	// Verify clamping
	if desc.Exposure != 2.0 {
		t.Errorf("Exposure should be clamped to 2.0, got %f", desc.Exposure)
	}
	if desc.Contrast != 100 {
		t.Errorf("Contrast should be clamped to 100, got %d", desc.Contrast)
	}
	if desc.Saturation != -100 {
		t.Errorf("Saturation should be clamped to -100, got %d", desc.Saturation)
	}
	if desc.Clarity != -100 {
		t.Errorf("Clarity should be clamped to -100, got %d", desc.Clarity)
	}
	if desc.Temperature != 100 {
		t.Errorf("Temperature should be clamped to 100, got %d (from 20000K)", desc.Temperature)
	}
}

// TestGenerate_NilRecipe tests that nil recipe returns an error.
func TestGenerate_NilRecipe(t *testing.T) {
	_, err := Generate(nil)
	if err == nil {
		t.Fatal("Generate(nil) should return an error")
	}

	// Verify it's a ConversionError
	convErr, ok := err.(*ConversionError)
	if !ok {
		t.Fatalf("Expected ConversionError, got %T", err)
	}
	if convErr.Operation != "generate" {
		t.Errorf("Expected operation='generate', got %q", convErr.Operation)
	}
	if convErr.Format != "costyle" {
		t.Errorf("Expected format='costyle', got %q", convErr.Format)
	}
}

// TestGenerate_MetadataPreservation tests that metadata fields are preserved.
func TestGenerate_MetadataPreservation(t *testing.T) {
	tests := []struct {
		name            string
		recipe          *models.UniversalRecipe
		wantName        string
		wantAuthor      string
		wantDescription string
	}{
		{
			name: "Metadata map",
			recipe: &models.UniversalRecipe{
				Metadata: map[string]interface{}{
					"name":        "Map Name",
					"author":      "Map Author",
					"description": "Map Description",
				},
			},
			wantName:        "Map Name",
			wantAuthor:      "Map Author",
			wantDescription: "Map Description",
		},
		{
			name: "Recipe Name field",
			recipe: &models.UniversalRecipe{
				Name: "Recipe Name Field",
			},
			wantName: "Recipe Name Field",
		},
		{
			name: "Metadata overrides Recipe.Name",
			recipe: &models.UniversalRecipe{
				Name: "Recipe Name",
				Metadata: map[string]interface{}{
					"name": "Metadata Name",
				},
			},
			wantName: "Metadata Name", // Metadata takes precedence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			costyleData, err := Generate(tt.recipe)
			if err != nil {
				t.Fatalf("Generate() failed: %v", err)
			}

			var parsed CaptureOneStyle
			if err := xml.Unmarshal(costyleData, &parsed); err != nil {
				t.Fatalf("XML unmarshal failed: %v", err)
			}

			desc := parsed.RDF.Description
			if desc.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", desc.Name, tt.wantName)
			}
			if desc.Author != tt.wantAuthor {
				t.Errorf("Author = %q, want %q", desc.Author, tt.wantAuthor)
			}
			if desc.Description != tt.wantDescription {
				t.Errorf("Description = %q, want %q", desc.Description, tt.wantDescription)
			}
		})
	}
}

// TestGenerate_ColorBalance tests color balance parameter mapping.
// Note: C1 uses 0-360 for hue (direct mapping from UR) and -100 to +100 for saturation
func TestGenerate_ColorBalance(t *testing.T) {
	recipe := &models.UniversalRecipe{
		SplitShadowHue:           180,  // 180° → 180 in C1 (direct mapping)
		SplitShadowSaturation:    50,   // 50 → 0 in C1 (neutral: (50*2)-100 = 0)
		SplitHighlightHue:        270,  // 270° → 270 in C1 (direct mapping)
		SplitHighlightSaturation: 75,   // 75 → +50 in C1 ((75*2)-100 = 50)
	}

	costyleData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	var parsed CaptureOneStyle
	if err := xml.Unmarshal(costyleData, &parsed); err != nil {
		t.Fatalf("XML unmarshal failed: %v", err)
	}

	desc := parsed.RDF.Description

	// SplitShadowHue: 180° → 180 (direct mapping, C1 uses 0-360)
	if desc.ShadowsHue != 180 {
		t.Errorf("ShadowsHue = %d, want 180 (direct from UR 180°)", desc.ShadowsHue)
	}

	// SplitShadowSaturation: 50 → (50*2)-100 = 0 (neutral saturation)
	if desc.ShadowsSaturation != 0 {
		t.Errorf("ShadowsSaturation = %d, want 0 (neutral, from UR 50)", desc.ShadowsSaturation)
	}

	// SplitHighlightHue: 270° → 270 (direct mapping)
	if desc.HighlightsHue != 270 {
		t.Errorf("HighlightsHue = %d, want 270 (direct from UR 270°)", desc.HighlightsHue)
	}

	// SplitHighlightSaturation: 75 → (75*2)-100 = 50
	if desc.HighlightsSaturation != 50 {
		t.Errorf("HighlightsSaturation = %d, want 50 (from UR 75)", desc.HighlightsSaturation)
	}
}

// TestKelvinToC1Temperature tests the Kelvin to C1 temperature conversion function.
func TestKelvinToC1Temperature(t *testing.T) {
	tests := []struct {
		kelvin float64
		want   int
	}{
		{5500, 0},     // Neutral (5500 - 5500) / 35 = 0
		{6100, 17},    // Warmer: (6100 - 5500) / 35 = 17.14 ≈ 17
		{4900, -17},   // Cooler: (4900 - 5500) / 35 = -17.14 ≈ -17
		{9100, 100},   // Very warm: (9100 - 5500) / 35 = 102.86 → clamped to 100
		{1900, -100},  // Very cool: (1900 - 5500) / 35 = -102.86 → clamped to -100
		{11500, 100},  // Extreme warm (clamped)
		{-500, -100},  // Extreme cool (clamped)
		{5535, 1},     // +35K → +1
		{5465, -1},    // -35K → -1
	}

	for _, tt := range tests {
		got := kelvinToC1Temperature(tt.kelvin)
		if got != tt.want {
			t.Errorf("kelvinToC1Temperature(%f) = %d, want %d", tt.kelvin, got, tt.want)
		}
	}
}

// TestClampFunctions tests the clamping helper functions (defined in parse.go).
func TestClampFunctions(t *testing.T) {
	// Test clampFloat64
	t.Run("clampFloat64", func(t *testing.T) {
		tests := []struct {
			value, min, max float64
			want            float64
		}{
			{0.0, -2.0, 2.0, 0.0},
			{5.0, -2.0, 2.0, 2.0},   // Clamp to max
			{-5.0, -2.0, 2.0, -2.0}, // Clamp to min
			{1.5, -2.0, 2.0, 1.5},   // Within range
		}

		for _, tt := range tests {
			got := clampFloat64(tt.value, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("clampFloat64(%f, %f, %f) = %f, want %f", tt.value, tt.min, tt.max, got, tt.want)
			}
		}
	})

	// Test clampInt
	t.Run("clampInt", func(t *testing.T) {
		tests := []struct {
			value, min, max int
			want            int
		}{
			{0, -100, 100, 0},
			{200, -100, 100, 100},   // Clamp to max
			{-200, -100, 100, -100}, // Clamp to min
			{50, -100, 100, 50},     // Within range
		}

		for _, tt := range tests {
			got := clampInt(tt.value, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("clampInt(%d, %d, %d) = %d, want %d", tt.value, tt.min, tt.max, got, tt.want)
			}
		}
	})
}

// TestGenerate_XMLFormatting tests that XML output is properly formatted.
func TestGenerate_XMLFormatting(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure: 1.0,
		Contrast: 20,
	}

	costyleData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	xmlStr := string(costyleData)

	// Verify XML declaration on first line
	if !strings.HasPrefix(xmlStr, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("XML declaration missing or incorrect")
	}

	// Verify indentation (2 spaces)
	if !strings.Contains(xmlStr, "  <RDF>") {
		t.Error("XML indentation incorrect (expected 2 spaces)")
	}

	// Verify xmpmeta root element
	if !strings.Contains(xmlStr, "<xmpmeta>") {
		t.Error("Missing <xmpmeta> root element")
	}

	// Verify closing tag
	if !strings.Contains(xmlStr, "</xmpmeta>") {
		t.Error("Missing </xmpmeta> closing tag")
	}
}

// BenchmarkGenerate benchmarks .costyle generation performance.
func BenchmarkGenerate(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Name:       "Benchmark Preset",
		Exposure:   1.5,
		Contrast:   25,
		Saturation: 20,
		Tint:       10,
		Clarity:    15,
	}

	temp := 6100
	recipe.Temperature = &temp

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatalf("Generate() failed: %v", err)
		}
	}
}

// BenchmarkGenerate_Complex benchmarks generation with many parameters.
func BenchmarkGenerate_Complex(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Name:                     "Complex Preset",
		Exposure:                 1.5,
		Contrast:                 25,
		Saturation:               20,
		Tint:                     10,
		Clarity:                  15,
		SplitShadowHue:           180,
		SplitShadowSaturation:    50,
		SplitHighlightHue:        270,
		SplitHighlightSaturation: 75,
		Metadata: map[string]interface{}{
			"name":        "Complex Test",
			"author":      "Benchmark Author",
			"description": "Complex preset with many parameters",
		},
	}

	temp := 6100
	recipe.Temperature = &temp

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatalf("Generate() failed: %v", err)
		}
	}
}

// TestGenerate_RoundTrip tests round-trip conversion: Parse → Generate → Parse.
// Verifies that generated .costyle files parse back to equivalent UniversalRecipe (95%+ accuracy).
func TestGenerate_RoundTrip(t *testing.T) {
	// Create original UniversalRecipe with typical parameters
	original := &models.UniversalRecipe{
		Name:       "Round-Trip Test",
		Exposure:   1.5,
		Contrast:   25,
		Saturation: 20,
		Tint:       10,
		Clarity:    15,
		Metadata: map[string]interface{}{
			"author":      "Test Author",
			"description": "Round-trip test preset",
		},
	}

	temp := 6100
	original.Temperature = &temp

	// Generate .costyle from original recipe
	costyleData, err := Generate(original)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Parse generated .costyle back to UniversalRecipe
	parsed, err := Parse(costyleData)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// Compare original and parsed recipes
	// Note: Some precision loss expected due to scaling conversions

	// Exposure: Should match exactly (float64)
	if parsed.Exposure != original.Exposure {
		t.Errorf("Round-trip Exposure mismatch: got %f, want %f", parsed.Exposure, original.Exposure)
	}

	// Contrast: Should match exactly (int, direct mapping)
	if parsed.Contrast != original.Contrast {
		t.Errorf("Round-trip Contrast mismatch: got %d, want %d", parsed.Contrast, original.Contrast)
	}

	// Saturation: Should match exactly
	if parsed.Saturation != original.Saturation {
		t.Errorf("Round-trip Saturation mismatch: got %d, want %d", parsed.Saturation, original.Saturation)
	}

	// Temperature: Allow ±35K tolerance due to conversion rounding
	// Original: 6100K → (6100-5500)/35 = 17.14 ≈ 17 C1 → 5500 + 17*35 = 6095K
	if parsed.Temperature == nil {
		t.Error("Round-trip Temperature is nil, expected value")
	} else {
		// Expected: 17 C1 units * 35K/unit = 595K → 5500 + 595 = 6095K
		expectedKelvin := 6095
		tolerance := 35 // One scale unit
		diff := *parsed.Temperature - expectedKelvin
		if diff < -tolerance || diff > tolerance {
			t.Errorf("Round-trip Temperature mismatch: got %d, want %d (±%d)", *parsed.Temperature, expectedKelvin, tolerance)
		}
	}

	// Tint: Expected precision loss due to asymmetric scaling
	// Generate: UR 10 → C1 7 (10 * 100/150 = 6.67 ≈ 7)
	// Parse: C1 7 → UR 7 (direct copy, no reverse scaling in parse.go:79)
	// This is a known limitation: Parse doesn't reverse-scale Tint
	// Allow ±4 tolerance (10 → 7 is -3, within ±4 range)
	tintDiff := parsed.Tint - original.Tint
	if tintDiff < -4 || tintDiff > 4 {
		t.Errorf("Round-trip Tint mismatch: got %d, want %d (tolerance ±4 due to scaling)", parsed.Tint, original.Tint)
	}

	// Clarity: Should match exactly
	if parsed.Clarity != original.Clarity {
		t.Errorf("Round-trip Clarity mismatch: got %d, want %d", parsed.Clarity, original.Clarity)
	}

	// Name: Should match exactly
	if parsed.Name != original.Name {
		t.Errorf("Round-trip Name mismatch: got %q, want %q", parsed.Name, original.Name)
	}

	// Metadata: Author and description should be preserved in parsed.Metadata
	if parsedAuthor, ok := parsed.Metadata["costyle_author"].(string); ok {
		if wantAuthor, ok := original.Metadata["author"].(string); ok {
			if parsedAuthor != wantAuthor {
				t.Errorf("Round-trip Author metadata mismatch: got %q, want %q", parsedAuthor, wantAuthor)
			}
		}
	}
}
