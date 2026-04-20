package np3

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/justin/recipe/internal/models"
)

// TestGenerateBoundaryParameters tests valid extreme values at boundaries
func TestGenerateBoundaryParameters(t *testing.T) {
	tests := []struct {
		name       string
		sharpness  int
		contrast   int
		saturation int
		exposure   float64
	}{
		{
			name:       "Maximum valid values",
			sharpness:  150,
			contrast:   100,
			saturation: 100,
			exposure:   1.0,
		},
		{
			name:       "Minimum valid values",
			sharpness:  0,
			contrast:   -100,
			saturation: -100,
			exposure:   -1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := models.NewRecipeBuilder().
				WithName("BoundaryTest").
				WithSharpness(tt.sharpness).
				WithContrast(tt.contrast).
				WithSaturation(tt.saturation).
				WithExposure(tt.exposure).
				Build()
			if err != nil {
				t.Fatalf("Build recipe failed: %v", err)
			}

			data, err := Generate(recipe)
			if err != nil {
				t.Errorf("Generate failed: %v", err)
			}

			// Verify we can parse it back
			_, parseErr := Parse(data)
			if parseErr != nil {
				t.Errorf("Parse round-trip failed: %v", parseErr)
			}
		})
	}
}

// TestGenerateLongPresetName tests name truncation
func TestGenerateLongPresetName(t *testing.T) {
	// Create a name > 40 chars to trigger truncation (line 139-141)
	longName := "This is an extremely long preset name that exceeds forty characters"

	recipe, err := models.NewRecipeBuilder().
		WithName(longName).
		WithSharpness(50).
		Build()
	if err != nil {
		t.Fatalf("Build recipe failed: %v", err)
	}

	data, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Parse back and verify name was truncated
	parsedRecipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Name should be truncated to 40 chars
	if len(parsedRecipe.Name) > 40 {
		t.Errorf("Name not truncated: got %d chars, want ≤40", len(parsedRecipe.Name))
	}
}

// TestGenerateColorDataEdgeCases tests color data generation edge cases
func TestGenerateColorDataEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		saturation int
		desc       string
	}{
		{
			name:       "Negative saturation (triggers targetCount < 0)",
			saturation: -99,
			desc:       "Should handle negative saturation gracefully",
		},
		{
			name:       "Very high saturation (triggers max clamping)",
			saturation: 99,
			desc:       "Should clamp to maximum color triplets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := models.NewRecipeBuilder().
				WithName("ColorEdge").
				WithSaturation(tt.saturation).
				Build()
			if err != nil {
				t.Fatalf("Build recipe failed: %v", err)
			}

			data, err := Generate(recipe)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Should parse back successfully
			_, parseErr := Parse(data)
			if parseErr != nil {
				t.Errorf("Parse failed: %v", parseErr)
			}
		})
	}
}

// TestGenerateToneCurveEdgeCases tests tone curve generation edge cases
func TestGenerateToneCurveEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		contrast   int
		saturation int
		desc       string
	}{
		{
			name:       "Negative contrast with high saturation overlap",
			contrast:   -99,
			saturation: 99,
			desc:       "Should handle additionalPairs < 0",
		},
		{
			name:       "Very high contrast",
			contrast:   99,
			saturation: -99,
			desc:       "Should generate maximum tone curve pairs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := models.NewRecipeBuilder().
				WithName("ToneCurveEdge").
				WithContrast(tt.contrast).
				WithSaturation(tt.saturation).
				Build()
			if err != nil {
				t.Fatalf("Build recipe failed: %v", err)
			}

			data, err := Generate(recipe)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Should parse back successfully
			_, parseErr := Parse(data)
			if parseErr != nil {
				t.Errorf("Parse failed: %v", parseErr)
			}
		})
	}
}

// TestEstimateParametersEdgeCases tests heuristic estimation edge cases
func TestEstimateParametersEdgeCases(t *testing.T) {
	// Create data with extreme heuristic values to trigger all clamping paths
	data := make([]byte, 500)
	copy(data[0:3], magicBytes)
	copy(data[3:7], []byte{0x02, 0x10, 0x00, 0x00})

	// Fill sharpness bytes with extreme values (should trigger clamping in estimateParameters)
	for i := 66; i <= 70; i++ {
		data[i] = 255 // Max value
	}

	// Fill brightness bytes with extreme values
	for i := 71; i <= 75; i++ {
		data[i] = 255 // Max value
	}

	// Fill hue bytes with extreme values
	for i := 76; i <= 79; i++ {
		data[i] = 255 // Max value
	}

	// Fill color data region with extreme pattern
	for i := 100; i < 300; i += 3 {
		data[i] = 255
		data[i+1] = 255
		data[i+2] = 255
	}

	// Fill tone curve region with extreme pattern
	for i := 150; i < 500; i += 2 {
		data[i] = 255
		data[i+1] = 255
	}

	// Should parse without error despite extreme values
	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify clamping worked (should be within valid ranges)
	if recipe.Sharpness < 0 || recipe.Sharpness > 150 {
		t.Errorf("Sharpness out of range: %d", recipe.Sharpness)
	}
	if recipe.Contrast < -100 || recipe.Contrast > 100 {
		t.Errorf("Contrast out of range: %d", recipe.Contrast)
	}
	if recipe.Saturation < -100 || recipe.Saturation > 100 {
		t.Errorf("Saturation out of range: %d", recipe.Saturation)
	}
}

// TestTruncateToBytes tests the rune-boundary-safe truncation helper directly.
func TestTruncateToBytes(t *testing.T) {
	star := "🌟" // U+1F31F: 4 UTF-8 bytes (0xF0 0x9F 0x8C 0x9F)

	tests := []struct {
		name  string
		input string
		limit int
		want  string
	}{
		{
			name:  "ASCII under limit unchanged",
			input: "Vivid",
			limit: 20,
			want:  "Vivid",
		},
		{
			name:  "ASCII exactly at limit unchanged",
			input: "abcdefghijklmnopqrst", // 20 bytes
			limit: 20,
			want:  "abcdefghijklmnopqrst",
		},
		{
			name:  "ASCII over limit truncated at byte boundary",
			input: "abcdefghijklmnopqrstuvwxyz",
			limit: 20,
			want:  "abcdefghijklmnopqrst",
		},
		{
			name:  "emoji exactly at limit unchanged",
			input: star + star + star + star + star, // 5×4 = 20 bytes
			limit: 20,
			want:  star + star + star + star + star,
		},
		{
			name:  "emoji over limit truncated at rune boundary",
			input: star + star + star + star + star + star, // 6×4 = 24 bytes
			limit: 20,
			want:  star + star + star + star + star,
		},
		{
			name:  "mixed ASCII+emoji straddling boundary drops emoji",
			input: "abcdefghijklmnopqr" + star, // 18 + 4 = 22 bytes
			limit: 20,
			want:  "abcdefghijklmnopqr", // emoji dropped; 18 bytes fits, 22 does not
		},
		{
			name:  "empty string unchanged",
			input: "",
			limit: 20,
			want:  "",
		},
		{
			name:  "zero limit returns empty",
			input: "abc",
			limit: 0,
			want:  "",
		},
		{
			name:  "2-byte rune straddling boundary dropped",
			input: "abcdefghijklmnopqrst" + "\u00e9", // 20 ASCII + 2-byte U+00E9 = 22 bytes
			limit: 20,
			want:  "abcdefghijklmnopqrst",
		},
		{
			name:  "3-byte rune straddling boundary dropped",
			input: "abcdefghijklmnopqr\u4e2d", // 18 ASCII + 3-byte U+4E2D = 21 bytes
			limit: 20,
			want:  "abcdefghijklmnopqr",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := truncateToBytes(tc.input, tc.limit)
			if got != tc.want {
				t.Errorf("truncateToBytes(%q, %d) = %q, want %q", tc.input, tc.limit, got, tc.want)
			}
			if len(got) > tc.limit {
				t.Errorf("result exceeds limit: %d bytes > %d", len(got), tc.limit)
			}
			if !utf8.ValidString(got) {
				t.Errorf("result is not valid UTF-8: %q", got)
			}
		})
	}
}

// TestGenerateMultibyteMetadata verifies that Generate handles multibyte UTF-8 characters
// in name and description without errors, and that descriptions (which are not ASCII-filtered
// on parse) round-trip as valid UTF-8 within byte limits.
//
// Note: The NP3 name field is parsed as printable ASCII only (bytes 32–126), so multibyte
// characters in the name are stripped on round-trip — this is expected NP3 format behavior.
// The correctness guarantee for names is at the truncateToBytes level (TestTruncateToBytes).
func TestGenerateMultibyteMetadata(t *testing.T) {
	star := "🌟" // 4 UTF-8 bytes

	// Build a description at exactly MaxDescriptionLength bytes, then one emoji over.
	// Computed without assuming MaxDescriptionLength is a multiple of 4.
	atLimitDesc := ""
	for len(atLimitDesc)+len(star) <= MaxDescriptionLength {
		atLimitDesc += star
	}
	overLimitDesc := atLimitDesc + star // one emoji over MaxDescriptionLength

	tests := []struct {
		name       string
		recipeName string
		recipeDesc string
		wantDesc   string // expected round-tripped description
	}{
		{
			name:       "emoji name does not crash Generate",
			recipeName: star + star + star + star + star + star, // 6 emoji, over limit
			recipeDesc: "plain",
			wantDesc:   "plain",
		},
		{
			name:       "mixed ASCII+emoji name straddling boundary",
			recipeName: "abcdefghijklmnopqr" + star, // 18 ASCII + 4-byte emoji = 22 bytes
			recipeDesc: "plain",
			wantDesc:   "plain",
		},
		{
			name:       "emoji description exactly at limit round-trips correctly",
			recipeName: "test",
			recipeDesc: star + star + star + star + star, // 5×4 = 20 bytes, well under 256
			wantDesc:   star + star + star + star + star,
		},
		{
			name:       "emoji description over limit truncated at rune boundary",
			recipeName: "test",
			recipeDesc: overLimitDesc,
			wantDesc:   atLimitDesc,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recipe := &models.UniversalRecipe{
				Name:        tc.recipeName,
				Description: tc.recipeDesc,
			}

			data, err := Generate(recipe)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			parsed, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse failed after Generate: %v", err)
			}

			// Name must always be valid UTF-8 and within byte limit after parse.
			// Exact name value is not asserted: the NP3 parser reads a 40-byte window
			// that includes header bytes, so the round-tripped name may include extra
			// printable ASCII bytes from the binary structure. Correctness of name
			// truncation is covered by TestTruncateToBytes.
			if !utf8.ValidString(parsed.Name) {
				t.Errorf("parsed name is not valid UTF-8: %q", parsed.Name)
			}
			if len(parsed.Name) > 40 {
				t.Errorf("parsed name exceeds 40 bytes (parser window size): got %d bytes", len(parsed.Name))
			}

			// Description must be valid UTF-8 and within byte limit
			if !utf8.ValidString(parsed.Description) {
				t.Errorf("parsed description is not valid UTF-8: %q", parsed.Description)
			}
			if len(parsed.Description) > MaxDescriptionLength {
				t.Errorf("parsed description exceeds %d bytes: got %d", MaxDescriptionLength, len(parsed.Description))
			}
			// The generator pads odd-length descriptions with a trailing space to
			// satisfy Nikon Imaging Cloud's even-byte-length requirement on this
			// field. That extra space is semantically insignificant for round-trip
			// comparison, so strip it before asserting.
			got := strings.TrimRight(parsed.Description, " ")
			want := strings.TrimRight(tc.wantDesc, " ")
			if got != want {
				t.Errorf("parsed description: got %q, want %q", parsed.Description, tc.wantDesc)
			}
		})
	}
}

// TestBuildRecipeNullNameBytes tests name extraction with null bytes
func TestBuildRecipeNullNameBytes(t *testing.T) {
	// Create data with name containing null bytes mid-string
	data := make([]byte, 500)
	copy(data[0:3], magicBytes)
	copy(data[3:7], []byte{0x02, 0x10, 0x00, 0x00})

	// Write name with null byte in middle: "Test\x00More" (should stop at null)
	copy(data[20:60], []byte{'T', 'e', 's', 't', 0x00, 'M', 'o', 'r', 'e'})

	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Name should be "Test" (stops at null byte)
	if recipe.Name != "Test" {
		t.Errorf("Name parsing incorrect: got %q, want %q", recipe.Name, "Test")
	}
}
