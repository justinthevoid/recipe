package models

import (
	"testing"
)

func TestBuilderWithMidRangeSharpening(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{"valid min", -5.0, false},
		{"valid max", 5.0, false},
		{"valid zero", 0.0, false},
		{"valid mid", 2.5, false},
		{"invalid too low", -5.1, true},
		{"invalid too high", 5.1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewRecipeBuilder().WithMidRangeSharpening(tt.value)
			recipe, err := builder.Build()

			if tt.wantErr {
				if err == nil {
					t.Errorf("WithMidRangeSharpening(%f) expected error, got nil", tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("WithMidRangeSharpening(%f) unexpected error: %v", tt.value, err)
				}
				if recipe.MidRangeSharpening != tt.value {
					t.Errorf("WithMidRangeSharpening(%f) = %f, want %f", tt.value, recipe.MidRangeSharpening, tt.value)
				}
			}
		})
	}
}

func TestBuilderWithColorGrading(t *testing.T) {
	// Valid color grading setup
	validHighlights := ColorGradingZone{Hue: 180, Chroma: 20, Brightness: -10}
	validMidtone := ColorGradingZone{Hue: 90, Chroma: -15, Brightness: 5}
	validShadows := ColorGradingZone{Hue: 270, Chroma: 30, Brightness: -20}
	validBlending := 50
	validBalance := 10

	t.Run("valid color grading", func(t *testing.T) {
		recipe, err := NewRecipeBuilder().
			WithColorGrading(validHighlights, validMidtone, validShadows, validBlending, validBalance).
			Build()

		if err != nil {
			t.Fatalf("Valid color grading returned error: %v", err)
		}
		if recipe.ColorGrading == nil {
			t.Fatal("ColorGrading is nil")
		}

		// Verify highlights
		if recipe.ColorGrading.Highlights.Hue != validHighlights.Hue {
			t.Errorf("Highlights Hue = %d, want %d", recipe.ColorGrading.Highlights.Hue, validHighlights.Hue)
		}
		if recipe.ColorGrading.Highlights.Chroma != validHighlights.Chroma {
			t.Errorf("Highlights Chroma = %d, want %d", recipe.ColorGrading.Highlights.Chroma, validHighlights.Chroma)
		}
		if recipe.ColorGrading.Highlights.Brightness != validHighlights.Brightness {
			t.Errorf("Highlights Brightness = %d, want %d", recipe.ColorGrading.Highlights.Brightness, validHighlights.Brightness)
		}

		// Verify midtone
		if recipe.ColorGrading.Midtone.Hue != validMidtone.Hue {
			t.Errorf("Midtone Hue = %d, want %d", recipe.ColorGrading.Midtone.Hue, validMidtone.Hue)
		}
		if recipe.ColorGrading.Midtone.Chroma != validMidtone.Chroma {
			t.Errorf("Midtone Chroma = %d, want %d", recipe.ColorGrading.Midtone.Chroma, validMidtone.Chroma)
		}
		if recipe.ColorGrading.Midtone.Brightness != validMidtone.Brightness {
			t.Errorf("Midtone Brightness = %d, want %d", recipe.ColorGrading.Midtone.Brightness, validMidtone.Brightness)
		}

		// Verify shadows
		if recipe.ColorGrading.Shadows.Hue != validShadows.Hue {
			t.Errorf("Shadows Hue = %d, want %d", recipe.ColorGrading.Shadows.Hue, validShadows.Hue)
		}
		if recipe.ColorGrading.Shadows.Chroma != validShadows.Chroma {
			t.Errorf("Shadows Chroma = %d, want %d", recipe.ColorGrading.Shadows.Chroma, validShadows.Chroma)
		}
		if recipe.ColorGrading.Shadows.Brightness != validShadows.Brightness {
			t.Errorf("Shadows Brightness = %d, want %d", recipe.ColorGrading.Shadows.Brightness, validShadows.Brightness)
		}

		// Verify global parameters
		if recipe.ColorGrading.Blending != validBlending {
			t.Errorf("Blending = %d, want %d", recipe.ColorGrading.Blending, validBlending)
		}
		if recipe.ColorGrading.Balance != validBalance {
			t.Errorf("Balance = %d, want %d", recipe.ColorGrading.Balance, validBalance)
		}
	})

	t.Run("invalid highlights hue", func(t *testing.T) {
		invalidHighlights := ColorGradingZone{Hue: 361, Chroma: 0, Brightness: 0}
		_, err := NewRecipeBuilder().
			WithColorGrading(invalidHighlights, validMidtone, validShadows, validBlending, validBalance).
			Build()
		if err == nil {
			t.Error("Expected error for invalid highlights hue, got nil")
		}
	})

	t.Run("invalid midtone chroma", func(t *testing.T) {
		invalidMidtone := ColorGradingZone{Hue: 180, Chroma: 101, Brightness: 0}
		_, err := NewRecipeBuilder().
			WithColorGrading(validHighlights, invalidMidtone, validShadows, validBlending, validBalance).
			Build()
		if err == nil {
			t.Error("Expected error for invalid midtone chroma, got nil")
		}
	})

	t.Run("invalid shadows brightness", func(t *testing.T) {
		invalidShadows := ColorGradingZone{Hue: 180, Chroma: 0, Brightness: 101}
		_, err := NewRecipeBuilder().
			WithColorGrading(validHighlights, validMidtone, invalidShadows, validBlending, validBalance).
			Build()
		if err == nil {
			t.Error("Expected error for invalid shadows brightness, got nil")
		}
	})

	t.Run("invalid blending", func(t *testing.T) {
		invalidBlending := 101
		_, err := NewRecipeBuilder().
			WithColorGrading(validHighlights, validMidtone, validShadows, invalidBlending, validBalance).
			Build()
		if err == nil {
			t.Error("Expected error for invalid blending, got nil")
		}
	})

	t.Run("invalid balance", func(t *testing.T) {
		invalidBalance := 101
		_, err := NewRecipeBuilder().
			WithColorGrading(validHighlights, validMidtone, validShadows, validBlending, invalidBalance).
			Build()
		if err == nil {
			t.Error("Expected error for invalid balance, got nil")
		}
	})
}

func TestBuilderWithColorGradingEdgeCases(t *testing.T) {
	t.Run("all zeros", func(t *testing.T) {
		zeros := ColorGradingZone{Hue: 0, Chroma: 0, Brightness: 0}
		recipe, err := NewRecipeBuilder().
			WithColorGrading(zeros, zeros, zeros, 0, 0).
			Build()

		if err != nil {
			t.Fatalf("All zeros returned error: %v", err)
		}
		if recipe.ColorGrading == nil {
			t.Fatal("ColorGrading is nil")
		}
	})

	t.Run("maximum values", func(t *testing.T) {
		maxHue := ColorGradingZone{Hue: 360, Chroma: 100, Brightness: 100}
		recipe, err := NewRecipeBuilder().
			WithColorGrading(maxHue, maxHue, maxHue, 100, 100).
			Build()

		if err != nil {
			t.Fatalf("Maximum values returned error: %v", err)
		}
		if recipe.ColorGrading == nil {
			t.Fatal("ColorGrading is nil")
		}
	})

	t.Run("minimum values", func(t *testing.T) {
		minVals := ColorGradingZone{Hue: 0, Chroma: -100, Brightness: -100}
		recipe, err := NewRecipeBuilder().
			WithColorGrading(minVals, minVals, minVals, 0, -100).
			Build()

		if err != nil {
			t.Fatalf("Minimum values returned error: %v", err)
		}
		if recipe.ColorGrading == nil {
			t.Fatal("ColorGrading is nil")
		}
	})
}

func TestBuilderColorGradingAndMidRangeSharpening(t *testing.T) {
	// Test using both new features together
	highlights := ColorGradingZone{Hue: 180, Chroma: 20, Brightness: -10}
	midtone := ColorGradingZone{Hue: 90, Chroma: -15, Brightness: 5}
	shadows := ColorGradingZone{Hue: 270, Chroma: 30, Brightness: -20}

	recipe, err := NewRecipeBuilder().
		WithName("Test Recipe").
		WithMidRangeSharpening(2.5).
		WithColorGrading(highlights, midtone, shadows, 50, 10).
		WithExposure(0.5).
		Build()

	if err != nil {
		t.Fatalf("Combined features returned error: %v", err)
	}

	if recipe.Name != "Test Recipe" {
		t.Errorf("Name = %s, want Test Recipe", recipe.Name)
	}
	if recipe.MidRangeSharpening != 2.5 {
		t.Errorf("MidRangeSharpening = %f, want 2.5", recipe.MidRangeSharpening)
	}
	if recipe.ColorGrading == nil {
		t.Fatal("ColorGrading is nil")
	}
	if recipe.Exposure != 0.5 {
		t.Errorf("Exposure = %f, want 0.5", recipe.Exposure)
	}
}
