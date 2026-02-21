// Package lrtemplate provides functionality for generating Adobe Lightroom Classic lrtemplate preset files.
//
// The lrtemplate format is a Lua table format used by Adobe Lightroom Classic to store photo editing
// presets. This package encodes the UniversalRecipe intermediate representation into lrtemplate files,
// enabling conversion from other preset formats.
//
// Format Structure:
//   - Lua table syntax starting with `s = {`
//   - Nested `value.settings` table containing all photo editing parameters
//   - Parameters stored as key-value pairs: `Exposure2012 = 1.5, Contrast2012 = 25`
//   - Arrays for tone curves: `ToneCurvePV2012 = { 0, 0, 255, 255 }`
//
// Supported Parameters (50+):
//   - Basic adjustments: Exposure, Contrast, Highlights, Shadows, Whites, Blacks
//   - Color: Saturation, Vibrance, Clarity, Sharpness, Temperature, Tint
//   - HSL adjustments: 8 colors × 3 properties (Hue, Saturation, Luminance)
//   - Advanced: Tone curves, Split toning, Grain, Vignette
//
// This generator uses bytes.Buffer for efficient string concatenation and achieves the
// <10ms performance target.
package lrtemplate

import (
	"bytes"
	"fmt"
	"math"

	"github.com/justin/recipe/internal/models"
)

// Generate converts a UniversalRecipe into a valid lrtemplate file format.
//
// The function produces a Lua table structure compatible with Adobe Lightroom Classic,
// following the format:
//
//	s = {
//	    id = "generated-uuid",
//	    internalName = "Generated Preset",
//	    title = "Generated Preset",
//	    type = "Develop",
//	    value = {
//	        settings = {
//	            Exposure2012 = 1.5,
//	            Contrast2012 = 25,
//	            ...
//	        },
//	        uuid = "settings-uuid",
//	    },
//	    version = 0,
//	}
//
// Parameters:
//   - recipe: The UniversalRecipe to convert. Must not be nil.
//
// Returns:
//   - []byte: The generated lrtemplate file content in Lua table format
//   - error: A ConversionError if generation fails, nil on success
//
// Error Conditions:
//   - nil recipe
//   - Invalid parameter values (handled via clamping with warnings in future)
//
// Performance:
//   - Target: <10ms for single file generation
//   - Implementation: Uses bytes.Buffer for efficient string building
//
// Example:
//
//	recipe := &models.UniversalRecipe{
//	    Name: "My Preset",
//	    Exposure: 1.5,
//	    Contrast: 25,
//	}
//	data, err := Generate(recipe)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("preset.lrtemplate", data, 0644)
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
	// Validate input
	if recipe == nil {
		return nil, &ConversionError{
			Operation: "generate",
			Format:    "lrtemplate",
			Cause:     fmt.Errorf("recipe is nil"),
		}
	}

	var buf bytes.Buffer

	// Generate unique IDs for the preset
	presetID := "00000000-0000-0000-0000-000000000000" // Placeholder UUID
	settingsUUID := "00000000-0000-0000-0000-000000000001"

	// Determine preset name
	presetName := recipe.Name
	if presetName == "" {
		presetName = "Generated Preset"
	}

	// Write outer Lua table structure
	buf.WriteString("s = {\n")
	buf.WriteString(fmt.Sprintf("\tid = \"%s\",\n", presetID))
	buf.WriteString(fmt.Sprintf("\tinternalName = \"%s\",\n", escapeString(presetName)))
	buf.WriteString(fmt.Sprintf("\ttitle = \"%s\",\n", escapeString(presetName)))
	buf.WriteString("\ttype = \"Develop\",\n")
	buf.WriteString("\tvalue = {\n")
	buf.WriteString("\t\tsettings = {\n")

	// Generate Process Version
	buf.WriteString("\t\t\tProcessVersion = \"10.0\",\n")

	// Generate Basic Adjustments
	if recipe.Exposure != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tExposure2012 = %.2f,\n", clampFloat(recipe.Exposure, -5.0, 5.0)))
	}
	if recipe.Contrast != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tContrast2012 = %d,\n", clampInt(recipe.Contrast, -100, 100)))
	}
	if recipe.Highlights != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tHighlights2012 = %d,\n", clampInt(recipe.Highlights, -100, 100)))
	}
	if recipe.Shadows != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tShadows2012 = %d,\n", clampInt(recipe.Shadows, -100, 100)))
	}
	if recipe.Whites != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tWhites2012 = %d,\n", clampInt(recipe.Whites, -100, 100)))
	}
	if recipe.Blacks != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tBlacks2012 = %d,\n", clampInt(recipe.Blacks, -100, 100)))
	}

	// Generate Color Parameters
	if recipe.Saturation != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tSaturation = %d,\n", clampInt(recipe.Saturation, -100, 100)))
	}
	if recipe.Vibrance != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tVibrance = %d,\n", clampInt(recipe.Vibrance, -100, 100)))
	}
	if recipe.Clarity != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tClarity2012 = %d,\n", clampInt(recipe.Clarity, -100, 100)))
	}
	if recipe.Sharpness != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tSharpness = %d,\n", clampInt(recipe.Sharpness, 0, 150)))
	}
	if recipe.Temperature != nil && *recipe.Temperature != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tTemperature = %d,\n", clampInt(*recipe.Temperature, 2000, 50000)))
	}
	if recipe.Tint != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tTint = %d,\n", clampInt(recipe.Tint, -150, 150)))
	}

	// Generate HSL Adjustments
	generateHSL(&buf, "Red", recipe.Red)
	generateHSL(&buf, "Orange", recipe.Orange)
	generateHSL(&buf, "Yellow", recipe.Yellow)
	generateHSL(&buf, "Green", recipe.Green)
	generateHSL(&buf, "Aqua", recipe.Aqua)
	generateHSL(&buf, "Blue", recipe.Blue)
	generateHSL(&buf, "Purple", recipe.Purple)
	generateHSL(&buf, "Magenta", recipe.Magenta)

	// Generate Split Toning
	if recipe.SplitShadowHue != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tSplitToningShadowHue = %d,\n", clampInt(recipe.SplitShadowHue, 0, 360)))
	}
	if recipe.SplitShadowSaturation != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tSplitToningShadowSaturation = %d,\n", clampInt(recipe.SplitShadowSaturation, 0, 100)))
	}
	if recipe.SplitHighlightHue != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tSplitToningHighlightHue = %d,\n", clampInt(recipe.SplitHighlightHue, 0, 360)))
	}
	if recipe.SplitHighlightSaturation != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tSplitToningHighlightSaturation = %d,\n", clampInt(recipe.SplitHighlightSaturation, 0, 100)))
	}
	if recipe.SplitBalance != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tSplitToningBalance = %d,\n", clampInt(recipe.SplitBalance, -100, 100)))
	}

	// Generate Tone Curve (Point Curve)
	if len(recipe.PointCurve) > 0 {
		buf.WriteString("\t\t\tToneCurvePV2012 = {\n")
		for _, point := range recipe.PointCurve {
			buf.WriteString(fmt.Sprintf("\t\t\t\t%d,\n\t\t\t\t%d,\n", clampInt(point.Input, 0, 255), clampInt(point.Output, 0, 255)))
		}
		buf.WriteString("\t\t\t},\n")
	}

	// Generate Tone Curve RGB channels
	if len(recipe.PointCurveRed) > 0 {
		buf.WriteString("\t\t\tToneCurvePV2012Red = {\n")
		for _, point := range recipe.PointCurveRed {
			buf.WriteString(fmt.Sprintf("\t\t\t\t%d,\n\t\t\t\t%d,\n", clampInt(point.Input, 0, 255), clampInt(point.Output, 0, 255)))
		}
		buf.WriteString("\t\t\t},\n")
	}
	if len(recipe.PointCurveGreen) > 0 {
		buf.WriteString("\t\t\tToneCurvePV2012Green = {\n")
		for _, point := range recipe.PointCurveGreen {
			buf.WriteString(fmt.Sprintf("\t\t\t\t%d,\n\t\t\t\t%d,\n", clampInt(point.Input, 0, 255), clampInt(point.Output, 0, 255)))
		}
		buf.WriteString("\t\t\t},\n")
	}
	if len(recipe.PointCurveBlue) > 0 {
		buf.WriteString("\t\t\tToneCurvePV2012Blue = {\n")
		for _, point := range recipe.PointCurveBlue {
			buf.WriteString(fmt.Sprintf("\t\t\t\t%d,\n\t\t\t\t%d,\n", clampInt(point.Input, 0, 255), clampInt(point.Output, 0, 255)))
		}
		buf.WriteString("\t\t\t},\n")
	}

	// Generate Camera Calibration
	if recipe.CameraProfile.RedHue != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tRedHue = %d,\n", clampInt(recipe.CameraProfile.RedHue, -100, 100)))
	}
	if recipe.CameraProfile.RedSaturation != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tRedSaturation = %d,\n", clampInt(recipe.CameraProfile.RedSaturation, -100, 100)))
	}
	if recipe.CameraProfile.GreenHue != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tGreenHue = %d,\n", clampInt(recipe.CameraProfile.GreenHue, -100, 100)))
	}
	if recipe.CameraProfile.GreenSaturation != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tGreenSaturation = %d,\n", clampInt(recipe.CameraProfile.GreenSaturation, -100, 100)))
	}
	if recipe.CameraProfile.BlueHue != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tBlueHue = %d,\n", clampInt(recipe.CameraProfile.BlueHue, -100, 100)))
	}
	if recipe.CameraProfile.BlueSaturation != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tBlueSaturation = %d,\n", clampInt(recipe.CameraProfile.BlueSaturation, -100, 100)))
	}

	// Generate Grain (Effects)
	if recipe.GrainAmount != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tGrainAmount = %d,\n", clampInt(recipe.GrainAmount, 0, 100)))
	}
	if recipe.GrainSize != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tGrainSize = %d,\n", clampInt(recipe.GrainSize, 0, 100)))
	}
	if recipe.GrainRoughness != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tGrainFrequency = %d,\n", clampInt(recipe.GrainRoughness, 0, 100)))
	}

	// Generate Vignette
	if recipe.VignetteAmount != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tVignetteAmount = %d,\n", clampInt(recipe.VignetteAmount, -100, 100)))
	}
	if recipe.VignetteMidpoint != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tVignetteMidpoint = %d,\n", clampInt(recipe.VignetteMidpoint, 0, 100)))
	}

	// Close settings table
	buf.WriteString("\t\t},\n")

	// Add settings UUID
	buf.WriteString(fmt.Sprintf("\t\tuuid = \"%s\",\n", settingsUUID))

	// Close value table
	buf.WriteString("\t},\n")

	// Add version
	buf.WriteString("\tversion = 0,\n")

	// Close outer table
	buf.WriteString("}\n")

	return buf.Bytes(), nil
}

// generateHSL writes HSL adjustment fields for a single color to the buffer.
// Only writes non-zero values to minimize output size.
func generateHSL(buf *bytes.Buffer, colorName string, adj models.ColorAdjustment) {
	if adj.Hue != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tHueAdjustment%s = %d,\n", colorName, clampInt(adj.Hue, -100, 100)))
	}
	if adj.Saturation != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tSaturationAdjustment%s = %d,\n", colorName, clampInt(adj.Saturation, -100, 100)))
	}
	if adj.Luminance != 0 {
		buf.WriteString(fmt.Sprintf("\t\t\tLuminanceAdjustment%s = %d,\n", colorName, clampInt(adj.Luminance, -100, 100)))
	}
}

// clampInt clamps an integer value to the specified range [min, max].
func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// clampFloat clamps a float64 value to the specified range [min, max].
func clampFloat(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// escapeString escapes special characters in strings for Lua syntax.
// Handles: backslashes, quotes, newlines, carriage returns, tabs.
func escapeString(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		switch r {
		case '\\':
			buf.WriteString("\\\\")
		case '"':
			buf.WriteString("\\\"")
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		default:
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
