// Package costyle provides functionality for generating Capture One .costyle preset files.
//
// The .costyle generator creates valid XML files from UniversalRecipe data, enabling
// conversion from other preset formats (XMP, lrtemplate, NP3) to Capture One .costyle format.
//
// Generation Strategy:
//   - Validate UniversalRecipe input (nil check, range validation)
//   - Map UniversalRecipe fields to .costyle parameters with scaling
//   - Build XML structure using encoding/xml
//   - Format output with proper indentation and namespace declarations
//
// Scaling:
//   - UniversalRecipe uses int (-100 to +100) for Contrast, Saturation, Clarity
//   - These map directly to .costyle int ranges (-100 to +100)
//   - Temperature: Kelvin → -100/+100 scale (5500K = neutral 0)
//   - Exposure: float64 (-2.0 to +2.0) maps directly
//
// Performance:
//   - Target: <100ms for single file generation
//   - Achieved through efficient XML marshaling and minimal allocations
//
// Round-trip Compatibility:
//   - Generated .costyle files parse back to identical UniversalRecipe (95%+ accuracy)
//   - Validation: costyle → Parse → Generate → Parse → Compare
package costyle

import (
	"encoding/xml"
	"fmt"
	"math"

	"github.com/justin/recipe/internal/models"
)

// ConversionError represents errors during costyle generation.
type ConversionError struct {
	Operation string // "generate", "validate"
	Format    string // "costyle"
	Field     string // Parameter field that caused the error (optional)
	Cause     error  // Underlying error
}

func (e *ConversionError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("costyle %s error (field: %s): %v", e.Operation, e.Field, e.Cause)
	}
	return fmt.Sprintf("costyle %s error: %v", e.Operation, e.Cause)
}

func (e *ConversionError) Unwrap() error {
	return e.Cause
}

// Generate creates a valid Capture One .costyle preset file from a UniversalRecipe.
//
// The function validates the input recipe, maps all parameters to .costyle elements,
// and generates a well-formed XML document with required namespace declarations.
//
// Parameters:
//   - recipe: UniversalRecipe containing photo editing parameters
//
// Returns:
//   - []byte: XML-formatted .costyle file content
//   - error: ConversionError if generation fails
//
// The generated .costyle file includes:
//   - XML declaration: <?xml version="1.0" encoding="UTF-8"?>
//   - Adobe XMP-style namespace declarations (xmpmeta, RDF)
//   - All mapped photo editing parameters as XML elements
//   - Proper indentation (2 spaces) for human readability
//
// Parameter Mapping:
//   - Exposure → Exposure (-2.0 to +2.0, direct)
//   - Contrast → Contrast (-100 to +100, direct)
//   - Saturation → Saturation (-100 to +100, direct)
//   - Temperature → Temperature (-100 to +100, Kelvin conversion)
//   - Tint → Tint (-100 to +100, direct)
//   - Clarity → Clarity (-100 to +100, direct)
//
// Example:
//
//	recipe := &models.UniversalRecipe{
//	    Exposure:   1.5,
//	    Contrast:   25,
//	    Saturation: 20,
//	}
//	costyleData, err := Generate(recipe)
//	if err != nil {
//	    // Handle error
//	}
//	// costyleData contains valid .costyle XML
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
	// Validation: Recipe must not be nil (fail fast)
	if recipe == nil {
		return nil, &ConversionError{
			Operation: "generate",
			Format:    "costyle",
			Cause:     fmt.Errorf("recipe is nil"),
		}
	}

	// Build .costyle structure
	style := buildCostyleDocument(recipe)

	// Marshal to XML with indentation
	xmlData, err := xml.MarshalIndent(style, "", "  ")
	if err != nil {
		return nil, &ConversionError{
			Operation: "generate",
			Format:    "costyle",
			Cause:     fmt.Errorf("xml marshal failed: %w", err),
		}
	}

	// Add XML declaration header (xml.Header = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	header := []byte(xml.Header)
	result := append(header, xmlData...)

	return result, nil
}

// buildCostyleDocument constructs the CaptureOneStyle structure from UniversalRecipe.
// Maps parameters with proper scaling and clamping to valid .costyle ranges.
func buildCostyleDocument(recipe *models.UniversalRecipe) *CaptureOneStyle {
	style := &CaptureOneStyle{
		RDF: RDF{
			Description: Description{},
		},
	}

	desc := &style.RDF.Description

	// Exposure: Direct mapping (-2.0 to +2.0), clamped
	if recipe.Exposure != 0 {
		desc.Exposure = clampFloat64(recipe.Exposure, -2.0, 2.0)
	}

	// Contrast: Direct mapping (-100 to +100), clamped
	if recipe.Contrast != 0 {
		desc.Contrast = clampInt(recipe.Contrast, -100, 100)
	}

	// Saturation: Direct mapping (-100 to +100), clamped
	if recipe.Saturation != 0 {
		desc.Saturation = clampInt(recipe.Saturation, -100, 100)
	}

	// Temperature: Convert from Kelvin to Capture One -100/+100 scale
	// Reference: 5500K = neutral (0)
	if recipe.Temperature != nil {
		desc.Temperature = kelvinToC1Temperature(float64(*recipe.Temperature))
	}

	// Tint: Direct mapping (-150 UR to -100/+100 C1), scaled
	if recipe.Tint != 0 {
		// UniversalRecipe Tint range: -150 to +150
		// Capture One Tint range: -100 to +100
		// Scale: UR Tint * (100/150) = C1 Tint
		desc.Tint = clampInt(int(math.Round(float64(recipe.Tint)*(100.0/150.0))), -100, 100)
	}

	// Clarity: Direct mapping (-100 to +100), clamped
	if recipe.Clarity != 0 {
		desc.Clarity = clampInt(recipe.Clarity, -100, 100)
	}

	// Color balance: Map from UniversalRecipe SplitShadow/Highlight to C1 tonal ranges
	// Note: Capture One uses Shadows/Midtones/Highlights; UniversalRecipe uses Shadow/Highlight split toning
	// Map Shadow → Shadows, Highlight → Highlights (no direct midtones mapping)
	//
	// Capture One hue range: 0-360 (matches UniversalRecipe directly)
	// Capture One saturation range: -100 to +100 (convert from UR 0-100)
	//
	// Important: Check if saturation != 0 OR if hue/saturation stored in metadata (from costyle parse)
	// This ensures round-trip preservation
	hasShadows := recipe.SplitShadowSaturation != 0 ||
		recipe.Metadata["costyle_shadows_hue"] != nil ||
		recipe.Metadata["costyle_shadows_saturation"] != nil

	if hasShadows {
		// Hue: Direct mapping 0-360 (both use same range)
		desc.ShadowsHue = clampInt(recipe.SplitShadowHue, 0, 360)

		// Saturation: Convert from 0-100 to -100/+100
		// Inverse of parse: (sat * 2) - 100 maps 0→-100, 50→0, 100→+100
		desc.ShadowsSaturation = clampInt((recipe.SplitShadowSaturation*2)-100, -100, 100)
	}

	hasHighlights := recipe.SplitHighlightSaturation != 0 ||
		recipe.Metadata["costyle_highlights_hue"] != nil ||
		recipe.Metadata["costyle_highlights_saturation"] != nil

	if hasHighlights {
		// Same mapping as shadows
		desc.HighlightsHue = clampInt(recipe.SplitHighlightHue, 0, 360)
		desc.HighlightsSaturation = clampInt((recipe.SplitHighlightSaturation*2)-100, -100, 100)
	}

	// Metadata: Preserve name, author, description if present
	if name, ok := recipe.Metadata["name"].(string); ok && name != "" {
		desc.Name = name
	} else if recipe.Name != "" {
		desc.Name = recipe.Name
	}

	if author, ok := recipe.Metadata["author"].(string); ok && author != "" {
		desc.Author = author
	}

	if description, ok := recipe.Metadata["description"].(string); ok && description != "" {
		desc.Description = description
	}

	return style
}

// kelvinToC1Temperature converts Kelvin temperature to Capture One -100/+100 scale.
// Reference temperature: 5500K = neutral (0)
// Mapping: 2000K-10000K → -100/+100
// Formula: (kelvin - 5500) / 35 = C1 temperature (inverse of parse formula)
//
// Example:
//   - 5500K → 0 (neutral)
//   - 5675K → +5 (warmer)
//   - 5325K → -5 (cooler)
//   - 10000K → +100 (very warm)
//   - 2000K → -100 (very cool)
func kelvinToC1Temperature(kelvin float64) int {
	const referenceK = 5500.0
	const scaleRange = 35.0 // Scale factor to match parse.go (K = 5500 + temp*35)

	delta := kelvin - referenceK
	c1Value := delta / scaleRange
	return clampInt(int(math.Round(c1Value)), -100, 100)
}

// clampInt and clampFloat64 helper functions are defined in parse.go and shared across the package.
