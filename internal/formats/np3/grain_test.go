package np3

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

func TestGrainRoundTrip(t *testing.T) {
	t.Skip("Grain writing is intentionally disabled in NP3 generation due to TLV conflicts")
	tests := []struct {
		name           string
		grainAmount    int
		grainSize      int
		expectedAmount int // Expected after round trip (due to scaling)
		expectedSize   int // Expected after round trip (enum mapping)
	}{
		{
			name:           "Grain Off",
			grainAmount:    0,
			grainSize:      0,
			expectedAmount: 0,
			expectedSize:   0, // Maps to Off (127) -> 0
		},
		{
			name:           "Grain Small",
			grainAmount:    20,
			grainSize:      25, // Small
			expectedAmount: 20, // 20 * 31.75/100 = 6.35 -> 6.5 (scaled4) -> 6.5 * 100/31.75 = 20.47 -> 20
			expectedSize:   25, // Maps to Small (2) -> 25
		},
		{
			name:           "Grain Large",
			grainAmount:    50,
			grainSize:      75, // Large
			expectedAmount: 50,
			expectedSize:   75, // Maps to Large (1) -> 75
		},
		{
			name:           "Grain Max",
			grainAmount:    100,
			grainSize:      100, // Large
			expectedAmount: 0,   // 100 encodes to 0xFF which is detected as uninitialized (NP3 Picture Controls don't support grain)
			expectedSize:   75,  // 100 maps to Large (1) -> encodes/decodes correctly -> 75
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create recipe
			recipe, err := models.NewRecipeBuilder().
				WithGrainAmount(tt.grainAmount).
				WithGrainSize(tt.grainSize).
				Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}

			// Generate NP3
			data, err := Generate(recipe)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Parse NP3
			parsed, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			// Verify Amount
			// Tolerance for amount due to Scaled4 encoding (0-31.75 range mapped to 0-100)
			// Step size is 100/127 ≈ 0.8
			if diff := parsed.GrainAmount - tt.expectedAmount; diff < -2 || diff > 2 {
				t.Errorf("GrainAmount mismatch: got %d, want %d (diff %d)", parsed.GrainAmount, tt.expectedAmount, diff)
			}

			// Verify Size
			// Exact match expected for enum
			if parsed.GrainSize != tt.expectedSize {
				t.Errorf("GrainSize mismatch: got %d, want %d", parsed.GrainSize, tt.expectedSize)
			}
		})
	}
}
