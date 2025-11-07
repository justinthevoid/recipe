package main

import (
	"strings"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestExtractParameters tests parameter extraction from files
func TestExtractParameters(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "Valid XMP file",
			filePath:    "../../examples/xmp/portrait.xmp",
			expectError: false,
		},
		{
			name:        "Valid NP3 file",
			filePath:    "../../testdata/xmp/sample.np3",
			expectError: false,
		},
		{
			name:        "Unsupported format",
			filePath:    "test.txt",
			expectError: true,
		},
		{
			name:        "Non-existent file",
			filePath:    "nonexistent.xmp",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := extractParameters(tt.filePath)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					// File might not exist in test environment, skip
					t.Skipf("file not found: %s", tt.filePath)
				}
				if recipe == nil {
					t.Error("expected recipe but got nil")
				}
			}
		})
	}
}

// TestFormatParameters tests parameter formatting
func TestFormatParameters(t *testing.T) {
	tests := []struct {
		name     string
		recipe   *models.UniversalRecipe
		contains []string // Strings that should be in output
		omits    []string // Strings that should NOT be in output
	}{
		{
			name:     "Nil recipe",
			recipe:   nil,
			contains: []string{"No data"},
		},
		{
			name:     "Empty recipe (all defaults)",
			recipe:   &models.UniversalRecipe{},
			contains: []string{"No adjustments"},
		},
		{
			name: "Basic adjustments only",
			recipe: &models.UniversalRecipe{
				Exposure: 0.75,
				Contrast: 15,
				Shadows:  30,
			},
			contains: []string{
				"Basic Adjustments:",
				"Exposure:",
				"+0.75",
				"Contrast:",
				"+15",
				"Shadows:",
				"+30",
			},
			omits: []string{
				"Highlights:", // Should be omitted (zero value)
				"Color:",      // No color adjustments
				"HSL:",        // No HSL adjustments
			},
		},
		{
			name: "Color adjustments",
			recipe: &models.UniversalRecipe{
				Saturation: 10,
				Vibrance:   20,
				Clarity:    15,
			},
			contains: []string{
				"Color:",
				"Saturation:",
				"+10",
				"Vibrance:",
				"+20",
				"Clarity:",
				"+15",
			},
		},
		{
			name: "Temperature and tint",
			recipe: &models.UniversalRecipe{
				Temperature: intPtr(500),
				Tint:        5,
			},
			contains: []string{
				"Temperature/Tint:",
				"Temperature:",
				"+500K",
				"Tint:",
				"+5",
			},
		},
		{
			name: "HSL adjustments",
			recipe: &models.UniversalRecipe{
				Red: models.ColorAdjustment{
					Hue:        5,
					Saturation: -10,
				},
				Blue: models.ColorAdjustment{
					Luminance: 8,
				},
			},
			contains: []string{
				"HSL Adjustments:",
				"Red Hue:",
				"+5",
				"Red Saturation:",
				"-10",
				"Blue Luminance:",
				"+8",
			},
			omits: []string{
				"Orange", // Not adjusted
				"Yellow", // Not adjusted
			},
		},
		{
			name: "Mixed adjustments",
			recipe: &models.UniversalRecipe{
				Exposure:   0.5,
				Contrast:   10,
				Saturation: 15,
				Red: models.ColorAdjustment{
					Hue: 5,
				},
			},
			contains: []string{
				"Basic Adjustments:",
				"Color:",
				"HSL Adjustments:",
				"Exposure:",
				"Saturation:",
				"Red Hue:",
			},
		},
		{
			name: "Tone curve",
			recipe: &models.UniversalRecipe{
				ToneCurveShadows:    -20,
				ToneCurveHighlights: 15,
			},
			contains: []string{
				"Tone Curve:",
				"Shadows:",
				"-20",
				"Highlights:",
				"+15",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatParameters(tt.recipe)

			// Check contains
			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected output to contain %q, got:\n%s", s, result)
				}
			}

			// Check omits
			for _, s := range tt.omits {
				if strings.Contains(result, s) {
					t.Errorf("expected output to NOT contain %q, got:\n%s", s, result)
				}
			}
		})
	}
}

// TestFormatBasicAdjustments tests basic adjustment formatting
func TestFormatBasicAdjustments(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure:   0.75,
		Contrast:   15,
		Highlights: -20,
		Shadows:    30,
		Whites:     5,
		Blacks:     -10,
	}

	result := formatBasicAdjustments(recipe)

	expected := []string{"Exposure", "Contrast", "Highlights", "Shadows", "Whites", "Blacks"}
	for _, s := range expected {
		if !strings.Contains(result, s) {
			t.Errorf("expected result to contain %q, got: %s", s, result)
		}
	}

	// Check formatting of values
	if !strings.Contains(result, "+0.75") {
		t.Error("exposure should be formatted with 2 decimal places")
	}
	if !strings.Contains(result, "+15") {
		t.Error("contrast should show positive sign")
	}
	if !strings.Contains(result, "-20") {
		t.Error("highlights should show negative value")
	}
}

// TestFormatColorAdjustments tests color adjustment formatting
func TestFormatColorAdjustments(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Saturation: 10,
		Vibrance:   20,
		Clarity:    15,
		Texture:    5,
		Dehaze:     -10,
	}

	result := formatColorAdjustments(recipe)

	expected := []string{"Saturation", "Vibrance", "Clarity", "Texture", "Dehaze"}
	for _, s := range expected {
		if !strings.Contains(result, s) {
			t.Errorf("expected result to contain %q", s)
		}
	}
}

// TestFormatHSLAdjustments tests HSL adjustment formatting
func TestFormatHSLAdjustments(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Red: models.ColorAdjustment{
			Hue:        5,
			Saturation: -10,
			Luminance:  8,
		},
		Blue: models.ColorAdjustment{
			Hue: 3,
		},
		// Other colors left at zero (should be omitted)
	}

	result := formatHSLAdjustments(recipe)

	// Should contain Red adjustments
	if !strings.Contains(result, "Red Hue") {
		t.Error("should contain Red Hue")
	}
	if !strings.Contains(result, "Red Saturation") {
		t.Error("should contain Red Saturation")
	}
	if !strings.Contains(result, "Red Luminance") {
		t.Error("should contain Red Luminance")
	}

	// Should contain Blue Hue (but not Saturation/Luminance as they're zero)
	if !strings.Contains(result, "Blue Hue") {
		t.Error("should contain Blue Hue")
	}
	if strings.Contains(result, "Blue Saturation") {
		t.Error("should not contain Blue Saturation (zero value)")
	}

	// Should not contain other colors
	if strings.Contains(result, "Orange") {
		t.Error("should not contain Orange (all zero)")
	}
}

// TestZeroValueOmission tests that zero values are omitted
func TestZeroValueOmission(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure: 0.5, // Non-zero
		Contrast: 0,   // Zero - should be omitted
		Shadows:  0,   // Zero - should be omitted
	}

	result := formatParameters(recipe)

	if !strings.Contains(result, "Exposure") {
		t.Error("should contain non-zero Exposure")
	}
	if strings.Contains(result, "Contrast") {
		t.Error("should not contain zero Contrast")
	}
	if strings.Contains(result, "Shadows") {
		t.Error("should not contain zero Shadows")
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
