package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/justin/recipe/internal/formats/lrtemplate"
	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/models"
)

// extractParameters detects file format and extracts parameters to UniversalRecipe
func extractParameters(filePath string) (*models.UniversalRecipe, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Detect format from extension
	ext := strings.ToLower(filepath.Ext(filePath))

	var recipe *models.UniversalRecipe

	switch ext {
	case ".np3":
		recipe, err = np3.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("NP3 parse error: %w", err)
		}

	case ".xmp":
		recipe, err = xmp.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("XMP parse error: %w", err)
		}

	case ".lrtemplate":
		recipe, err = lrtemplate.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("lrtemplate parse error: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	return recipe, nil
}

// formatParameters formats UniversalRecipe parameters into display string
// Groups parameters by category and omits zero/default values
func formatParameters(recipe *models.UniversalRecipe) string {
	if recipe == nil {
		return "  (No data)"
	}

	var sections []string

	// Basic Adjustments
	basic := formatBasicAdjustments(recipe)
	if basic != "" {
		sections = append(sections, "Basic Adjustments:", basic)
	}

	// Color/Presence
	color := formatColorAdjustments(recipe)
	if color != "" {
		sections = append(sections, "Color:", color)
	}

	// Temperature/Tint
	temp := formatTemperatureTint(recipe)
	if temp != "" {
		sections = append(sections, "Temperature/Tint:", temp)
	}

	// HSL Adjustments
	hsl := formatHSLAdjustments(recipe)
	if hsl != "" {
		sections = append(sections, "HSL Adjustments:", hsl)
	}

	// Sharpening
	sharp := formatSharpening(recipe)
	if sharp != "" {
		sections = append(sections, "Sharpening:", sharp)
	}

	// Tone Curve
	tone := formatToneCurve(recipe)
	if tone != "" {
		sections = append(sections, "Tone Curve:", tone)
	}

	// Split Toning
	split := formatSplitToning(recipe)
	if split != "" {
		sections = append(sections, "Split Toning:", split)
	}

	// If all sections empty, show "no adjustments"
	if len(sections) == 0 {
		return "  (No adjustments - all defaults)"
	}

	return strings.Join(sections, "\n\n")
}

// formatBasicAdjustments formats basic adjustment parameters
func formatBasicAdjustments(recipe *models.UniversalRecipe) string {
	var lines []string

	if recipe.Exposure != 0 {
		lines = append(lines, fmt.Sprintf("  Exposure:    %+.2f", recipe.Exposure))
	}
	if recipe.Contrast != 0 {
		lines = append(lines, fmt.Sprintf("  Contrast:    %+d", recipe.Contrast))
	}
	if recipe.Highlights != 0 {
		lines = append(lines, fmt.Sprintf("  Highlights:  %+d", recipe.Highlights))
	}
	if recipe.Shadows != 0 {
		lines = append(lines, fmt.Sprintf("  Shadows:     %+d", recipe.Shadows))
	}
	if recipe.Whites != 0 {
		lines = append(lines, fmt.Sprintf("  Whites:      %+d", recipe.Whites))
	}
	if recipe.Blacks != 0 {
		lines = append(lines, fmt.Sprintf("  Blacks:      %+d", recipe.Blacks))
	}

	return strings.Join(lines, "\n")
}

// formatColorAdjustments formats color/presence parameters
func formatColorAdjustments(recipe *models.UniversalRecipe) string {
	var lines []string

	if recipe.Saturation != 0 {
		lines = append(lines, fmt.Sprintf("  Saturation:  %+d", recipe.Saturation))
	}
	if recipe.Vibrance != 0 {
		lines = append(lines, fmt.Sprintf("  Vibrance:    %+d", recipe.Vibrance))
	}
	if recipe.Clarity != 0 {
		lines = append(lines, fmt.Sprintf("  Clarity:     %+d", recipe.Clarity))
	}
	if recipe.Texture != 0 {
		lines = append(lines, fmt.Sprintf("  Texture:     %+d", recipe.Texture))
	}
	if recipe.Dehaze != 0 {
		lines = append(lines, fmt.Sprintf("  Dehaze:      %+d", recipe.Dehaze))
	}

	return strings.Join(lines, "\n")
}

// formatTemperatureTint formats white balance parameters
func formatTemperatureTint(recipe *models.UniversalRecipe) string {
	var lines []string

	if recipe.Temperature != nil && *recipe.Temperature != 0 {
		lines = append(lines, fmt.Sprintf("  Temperature: %+dK", *recipe.Temperature))
	}
	if recipe.Tint != 0 {
		lines = append(lines, fmt.Sprintf("  Tint:        %+d", recipe.Tint))
	}

	return strings.Join(lines, "\n")
}

// formatHSLAdjustments formats HSL color adjustments
func formatHSLAdjustments(recipe *models.UniversalRecipe) string {
	var lines []string

	colors := []struct {
		name   string
		adjust models.ColorAdjustment
	}{
		{"Red", recipe.Red},
		{"Orange", recipe.Orange},
		{"Yellow", recipe.Yellow},
		{"Green", recipe.Green},
		{"Aqua", recipe.Aqua},
		{"Blue", recipe.Blue},
		{"Purple", recipe.Purple},
		{"Magenta", recipe.Magenta},
	}

	for _, c := range colors {
		hasAdjustment := c.adjust.Hue != 0 || c.adjust.Saturation != 0 || c.adjust.Luminance != 0
		if hasAdjustment {
			if c.adjust.Hue != 0 {
				lines = append(lines, fmt.Sprintf("  %s Hue:        %+d", c.name, c.adjust.Hue))
			}
			if c.adjust.Saturation != 0 {
				lines = append(lines, fmt.Sprintf("  %s Saturation: %+d", c.name, c.adjust.Saturation))
			}
			if c.adjust.Luminance != 0 {
				lines = append(lines, fmt.Sprintf("  %s Luminance:  %+d", c.name, c.adjust.Luminance))
			}
		}
	}

	return strings.Join(lines, "\n")
}

// formatSharpening formats sharpening parameters
func formatSharpening(recipe *models.UniversalRecipe) string {
	var lines []string

	if recipe.Sharpness != 0 {
		lines = append(lines, fmt.Sprintf("  Sharpness:   %d", recipe.Sharpness))
	}
	if recipe.SharpnessRadius != 0 {
		lines = append(lines, fmt.Sprintf("  Radius:      %.1f", recipe.SharpnessRadius))
	}
	if recipe.SharpnessDetail != 0 {
		lines = append(lines, fmt.Sprintf("  Detail:      %d", recipe.SharpnessDetail))
	}
	if recipe.SharpnessMasking != 0 {
		lines = append(lines, fmt.Sprintf("  Masking:     %d", recipe.SharpnessMasking))
	}

	return strings.Join(lines, "\n")
}

// formatToneCurve formats tone curve parameters
func formatToneCurve(recipe *models.UniversalRecipe) string {
	var lines []string

	if recipe.ToneCurveShadows != 0 {
		lines = append(lines, fmt.Sprintf("  Shadows:     %+d", recipe.ToneCurveShadows))
	}
	if recipe.ToneCurveDarks != 0 {
		lines = append(lines, fmt.Sprintf("  Darks:       %+d", recipe.ToneCurveDarks))
	}
	if recipe.ToneCurveLights != 0 {
		lines = append(lines, fmt.Sprintf("  Lights:      %+d", recipe.ToneCurveLights))
	}
	if recipe.ToneCurveHighlights != 0 {
		lines = append(lines, fmt.Sprintf("  Highlights:  %+d", recipe.ToneCurveHighlights))
	}

	// Point curves
	if len(recipe.PointCurve) > 0 {
		lines = append(lines, fmt.Sprintf("  Point Curve: %d points", len(recipe.PointCurve)))
	}

	return strings.Join(lines, "\n")
}

// formatSplitToning formats split toning parameters
func formatSplitToning(recipe *models.UniversalRecipe) string {
	var lines []string

	if recipe.SplitShadowHue != 0 {
		lines = append(lines, fmt.Sprintf("  Shadow Hue:  %d°", recipe.SplitShadowHue))
	}
	if recipe.SplitShadowSaturation != 0 {
		lines = append(lines, fmt.Sprintf("  Shadow Sat:  %d", recipe.SplitShadowSaturation))
	}
	if recipe.SplitHighlightHue != 0 {
		lines = append(lines, fmt.Sprintf("  Highlight Hue: %d°", recipe.SplitHighlightHue))
	}
	if recipe.SplitHighlightSaturation != 0 {
		lines = append(lines, fmt.Sprintf("  Highlight Sat: %d", recipe.SplitHighlightSaturation))
	}

	return strings.Join(lines, "\n")
}
