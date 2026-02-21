package costyle

import (
	"encoding/xml"
	"fmt"
	"math"

	"github.com/justin/recipe/internal/models"
)

// Parse parses a Capture One .costyle XML file and returns a UniversalRecipe.
//
// The .costyle format is an XML-based preset format used by Capture One photo editing software.
// It contains adjustment parameters similar to Adobe XMP, using RDF/Description structure.
//
// Parameters:
//   - data: Raw bytes of the .costyle XML file
//
// Returns:
//   - *models.UniversalRecipe: Populated recipe with all supported parameters
//   - error: Parsing error if XML is malformed or validation fails
//
// Example:
//   data, _ := os.ReadFile("portrait.costyle")
//   recipe, err := costyle.Parse(data)
//   if err != nil {
//       log.Fatal(err)
//   }
func Parse(data []byte) (*models.UniversalRecipe, error) {
	// Unmarshal XML structure
	var style CaptureOneStyle
	if err := xml.Unmarshal(data, &style); err != nil {
		return nil, fmt.Errorf("failed to parse .costyle XML: %w", err)
	}

	desc := style.RDF.Description

	// Initialize UniversalRecipe
	recipe := &models.UniversalRecipe{
		SourceFormat: "costyle",
		Metadata:     make(map[string]interface{}),
	}

	// Map metadata fields
	if desc.Name != "" {
		recipe.Name = desc.Name
		recipe.Metadata["costyle_name"] = desc.Name
	}
	if desc.Author != "" {
		recipe.Metadata["costyle_author"] = desc.Author
	}
	if desc.Description != "" {
		recipe.Metadata["costyle_description"] = desc.Description
	}

	// Map core adjustments
	// Exposure: -2.0 to +2.0 (direct map, clamp to UniversalRecipe range -5.0 to +5.0)
	recipe.Exposure = clampFloat64(desc.Exposure, -5.0, 5.0)

	// Contrast: -100 to +100 (direct map to UniversalRecipe)
	recipe.Contrast = clampInt(desc.Contrast, -100, 100)

	// Saturation: -100 to +100 (direct map to UniversalRecipe)
	recipe.Saturation = clampInt(desc.Saturation, -100, 100)

	// Temperature: -100 to +100 in Capture One (needs conversion to Kelvin for UniversalRecipe)
	// Capture One uses relative temperature scale, convert to absolute Kelvin
	// Assuming 0 = 5500K (daylight), -100 = ~2000K, +100 = ~10000K
	// UniversalRecipe.Temperature expects absolute Kelvin values in range [2000, 50000]
	if desc.Temperature != 0 {
		temp := clampInt(desc.Temperature, -100, 100)
		// Convert Capture One temperature (-100..+100) to absolute Kelvin
		// 0 → 5500K (baseline), -100 → 2000K, +100 → 10000K (approximate mapping)
		// Formula: K = 5500 + (temp * 35)
		kelvin := 5500 + int(float64(temp)*35.0)
		recipe.Temperature = &kelvin
		recipe.Metadata["costyle_temperature_relative"] = temp
	}

	// Tint: -100 to +100 in Capture One (scale to UniversalRecipe.Tint range -150 to +150)
	// Inverse of generate formula: C1 Tint * (150/100) = UR Tint
	// This ensures round-trip preservation: C1 → UR → C1
	if desc.Tint != 0 {
		recipe.Tint = clampInt(int(math.Round(float64(desc.Tint)*(150.0/100.0))), -150, 150)
	}

	// Clarity: -100 to +100 (direct map to UniversalRecipe)
	recipe.Clarity = clampInt(desc.Clarity, -100, 100)

	// Map color balance to split toning fields
	// Capture One has separate hue/saturation for shadows, midtones, highlights
	// UniversalRecipe has SplitShadowHue/Saturation and SplitHighlightHue/Saturation

	// Shadows: Map to SplitShadowHue/SplitShadowSaturation
	// Capture One stores hue in 0-360 range (directly compatible with UniversalRecipe)
	// Saturation uses -100 to +100 range (needs conversion to 0-100)
	if desc.ShadowsHue != 0 || desc.ShadowsSaturation != 0 {
		// Hue: Direct mapping 0-360 (Capture One range matches UniversalRecipe)
		recipe.SplitShadowHue = clampInt(desc.ShadowsHue, 0, 360)

		// Saturation: Convert from -100/+100 to 0/100
		// Formula: (sat + 100) / 2 maps -100→0, 0→50, +100→100
		shadowSat := (desc.ShadowsSaturation + 100) / 2
		recipe.SplitShadowSaturation = clampInt(shadowSat, 0, 100)

		recipe.Metadata["costyle_shadows_hue"] = desc.ShadowsHue
		recipe.Metadata["costyle_shadows_saturation"] = desc.ShadowsSaturation
	}

	// Highlights: Map to SplitHighlightHue/SplitHighlightSaturation
	// Same formulas as shadows
	if desc.HighlightsHue != 0 || desc.HighlightsSaturation != 0 {
		// Hue: Direct mapping 0-360
		recipe.SplitHighlightHue = clampInt(desc.HighlightsHue, 0, 360)

		// Saturation: Convert from -100/+100 to 0/100
		highlightSat := (desc.HighlightsSaturation + 100) / 2
		recipe.SplitHighlightSaturation = clampInt(highlightSat, 0, 100)

		recipe.Metadata["costyle_highlights_hue"] = desc.HighlightsHue
		recipe.Metadata["costyle_highlights_saturation"] = desc.HighlightsSaturation
	}

	// Midtones: Store in metadata (no direct UniversalRecipe field for midtone color)
	if desc.MidtonesHue != 0 || desc.MidtonesSaturation != 0 {
		recipe.Metadata["costyle_midtones_hue"] = desc.MidtonesHue
		recipe.Metadata["costyle_midtones_saturation"] = desc.MidtonesSaturation
	}

	return recipe, nil
}

// clampFloat64 clamps a float64 value to the given range.
func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// clampInt clamps an int value to the given range.
func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

