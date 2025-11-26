package xmp

import (
	"strings"
	"testing"

	"github.com/justin/recipe/internal/models"
)

func TestGrainRoundTrip(t *testing.T) {
	tests := []struct {
		name           string
		grainAmount    int
		grainSize      int
		grainRoughness int
	}{
		{
			name:           "Grain Off",
			grainAmount:    0,
			grainSize:      0,
			grainRoughness: 0,
		},
		{
			name:           "Grain Typical",
			grainAmount:    25,
			grainSize:      30,
			grainRoughness: 50,
		},
		{
			name:           "Grain Max",
			grainAmount:    100,
			grainSize:      100,
			grainRoughness: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create recipe
			recipe, err := models.NewRecipeBuilder().
				WithGrain(tt.grainAmount, tt.grainSize, tt.grainRoughness).
				Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}

			// Generate XMP
			data, err := Generate(recipe)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Verify XML contains attributes
			xmlStr := string(data)
			if tt.grainAmount > 0 {
				if !strings.Contains(xmlStr, "crs:GrainAmount=") {
					t.Error("XML missing crs:GrainAmount")
				}
			}

			// Parse XMP
			parsed, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			// Verify Amount
			if parsed.GrainAmount != tt.grainAmount {
				t.Errorf("GrainAmount mismatch: got %d, want %d", parsed.GrainAmount, tt.grainAmount)
			}

			// Verify Size
			if parsed.GrainSize != tt.grainSize {
				t.Errorf("GrainSize mismatch: got %d, want %d", parsed.GrainSize, tt.grainSize)
			}

			// Verify Roughness (Frequency)
			if parsed.GrainRoughness != tt.grainRoughness {
				t.Errorf("GrainRoughness mismatch: got %d, want %d", parsed.GrainRoughness, tt.grainRoughness)
			}
		})
	}
}
