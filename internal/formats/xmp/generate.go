// Package xmp provides functionality for generating Adobe Lightroom CC XMP preset files.
//
// The XMP generator creates valid XML files from UniversalRecipe data, enabling
// conversion from other preset formats (NP3, lrtemplate) to Lightroom CC XMP format.
//
// Generation Strategy:
//   - Validate UniversalRecipe input (nil check, range validation)
//   - Map UniversalRecipe fields to XMP attributes
//   - Build XML structure using encoding/xml
//   - Format output with proper indentation and namespace declarations
//
// Performance:
//   - Target: <30ms for single file generation
//   - Achieved through efficient XML marshaling and minimal allocations
//
// Round-trip Compatibility:
//   - Generated XMP files parse back to identical UniversalRecipe (±1 tolerance)
//   - Validation: XMP → Parse → Generate → Parse → Compare
package xmp

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/justin/recipe/internal/models"
)

// Generate creates a valid Adobe Lightroom CC XMP preset file from a UniversalRecipe.
//
// The function validates the input recipe, maps all parameters to XMP attributes,
// and generates a well-formed XML document with required namespace declarations.
//
// Parameters:
//   - recipe: UniversalRecipe containing photo editing parameters
//
// Returns:
//   - []byte: XML-formatted XMP file content
//   - error: ConversionError if generation fails
//
// The generated XMP file includes:
//   - XML declaration: <?xml version="1.0" encoding="UTF-8"?>
//   - Adobe XMP namespace declarations (x, crs, rdf)
//   - All 50+ photo editing parameters as XML attributes
//   - Proper indentation (2 spaces) for human readability
//
// Example:
//
//	recipe := &models.UniversalRecipe{
//	    Exposure: 1.5,
//	    Contrast: 25,
//	    Saturation: 20,
//	}
//	xmpData, err := Generate(recipe)
//	if err != nil {
//	    // Handle error
//	}
//	// xmpData contains valid XMP XML
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
	// Validation: Recipe must not be nil (fail fast)
	if recipe == nil {
		return nil, &ConversionError{
			Operation: "generate",
			Format:    "xmp",
			Cause:     fmt.Errorf("recipe is nil"),
		}
	}

	// Validate parameter ranges before generating
	if err := validateRecipe(recipe); err != nil {
		return nil, err
	}

	// Build XMP document structure
	doc := buildXMPDocument(recipe)

	// Marshal to XML with indentation
	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, &ConversionError{
			Operation: "generate",
			Format:    "xmp",
			Cause:     fmt.Errorf("xml marshal failed: %w", err),
		}
	}

	// Add XML declaration header
	header := []byte(xml.Header)
	result := append(header, xmlData...)

	return result, nil
}

// validateRecipe validates that all UniversalRecipe parameters are within valid XMP ranges.
// Returns ConversionError with field context if validation fails.
func validateRecipe(recipe *models.UniversalRecipe) error {
	// Validate Exposure (-5.0 to +5.0)
	if recipe.Exposure < -5.0 || recipe.Exposure > 5.0 {
		return &ConversionError{
			Operation: "generate",
			Format:    "xmp",
			Field:     "Exposure",
			Cause:     fmt.Errorf("value %.2f out of range [-5.0, +5.0]", recipe.Exposure),
		}
	}

	// Validate basic adjustments (-100 to +100)
	if err := validateRange(recipe.Contrast, -100, 100, "Contrast"); err != nil {
		return err
	}
	if err := validateRange(recipe.Highlights, -100, 100, "Highlights"); err != nil {
		return err
	}
	if err := validateRange(recipe.Shadows, -100, 100, "Shadows"); err != nil {
		return err
	}
	if err := validateRange(recipe.Whites, -100, 100, "Whites"); err != nil {
		return err
	}
	if err := validateRange(recipe.Blacks, -100, 100, "Blacks"); err != nil {
		return err
	}

	// Validate color parameters
	if err := validateRange(recipe.Saturation, -100, 100, "Saturation"); err != nil {
		return err
	}
	if err := validateRange(recipe.Vibrance, -100, 100, "Vibrance"); err != nil {
		return err
	}
	if err := validateRange(recipe.Clarity, -100, 100, "Clarity"); err != nil {
		return err
	}
	if err := validateRange(recipe.Sharpness, 0, 150, "Sharpness"); err != nil {
		return err
	}

	// Validate Temperature (nullable) - XMP uses integer Kelvin values
	if recipe.Temperature != nil {
		if *recipe.Temperature < 2000 || *recipe.Temperature > 50000 {
			return &ConversionError{
				Operation: "generate",
				Format:    "xmp",
				Field:     "Temperature",
				Cause:     fmt.Errorf("value %d out of range [2000, 50000]", *recipe.Temperature),
			}
		}
	}

	if err := validateRange(recipe.Tint, -150, 150, "Tint"); err != nil {
		return err
	}

	// Validate HSL adjustments for all colors (-100 to +100 for Hue, Saturation, Luminance)
	colors := []struct {
		name       string
		adjustment models.ColorAdjustment
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

	for _, color := range colors {
		if err := validateRange(color.adjustment.Hue, -100, 100, fmt.Sprintf("Hue%s", color.name)); err != nil {
			return err
		}
		if err := validateRange(color.adjustment.Saturation, -100, 100, fmt.Sprintf("Saturation%s", color.name)); err != nil {
			return err
		}
		if err := validateRange(color.adjustment.Luminance, -100, 100, fmt.Sprintf("Luminance%s", color.name)); err != nil {
			return err
		}
	}

	// Validate Split Toning
	if err := validateRange(recipe.SplitShadowHue, 0, 360, "SplitShadowHue"); err != nil {
		return err
	}
	if err := validateRange(recipe.SplitShadowSaturation, 0, 100, "SplitShadowSaturation"); err != nil {
		return err
	}
	if err := validateRange(recipe.SplitHighlightHue, 0, 360, "SplitHighlightHue"); err != nil {
		return err
	}
	if err := validateRange(recipe.SplitHighlightSaturation, 0, 100, "SplitHighlightSaturation"); err != nil {
		return err
	}
	if err := validateRange(recipe.SplitBalance, -100, 100, "SplitBalance"); err != nil {
		return err
	}

	return nil
}

// validateRange validates that an integer value is within the specified range.
// Returns ConversionError if validation fails.
func validateRange(value, min, max int, fieldName string) error {
	if value < min || value > max {
		return &ConversionError{
			Operation: "generate",
			Format:    "xmp",
			Field:     fieldName,
			Cause:     fmt.Errorf("value %d out of range [%d, %d]", value, min, max),
		}
	}
	return nil
}

// buildXMPDocument constructs the complete XMP XML structure from a UniversalRecipe.
// Maps all recipe parameters to XMP attributes with proper formatting and namespace declarations.
func buildXMPDocument(recipe *models.UniversalRecipe) *xmpDocWrapper {
	// Build Description with all parameters
	desc := descriptionWrapper{
		XMLNS: nsCameraRaw,

		// Basic Adjustments (formatted with appropriate precision)
		Exposure2012:   formatFloat(recipe.Exposure),
		Contrast2012:   formatInt(recipe.Contrast),
		Highlights2012: formatInt(recipe.Highlights),
		Shadows2012:    formatInt(recipe.Shadows),
		Whites2012:     formatInt(recipe.Whites),
		Blacks2012:     formatInt(recipe.Blacks),

		// Color Parameters
		Saturation:  formatInt(recipe.Saturation),
		Vibrance:    formatInt(recipe.Vibrance),
		Clarity2012: formatInt(recipe.Clarity),
		Sharpness:   formatInt(recipe.Sharpness),
		Tint:        formatInt(recipe.Tint),

		// Temperature (nullable - handle nil)
		Temperature: formatTemperature(recipe.Temperature),

		// HSL Adjustments - Red
		HueRed:        formatInt(recipe.Red.Hue),
		SaturationRed: formatInt(recipe.Red.Saturation),
		LuminanceRed:  formatInt(recipe.Red.Luminance),

		// HSL Adjustments - Orange
		HueOrange:        formatInt(recipe.Orange.Hue),
		SaturationOrange: formatInt(recipe.Orange.Saturation),
		LuminanceOrange:  formatInt(recipe.Orange.Luminance),

		// HSL Adjustments - Yellow
		HueYellow:        formatInt(recipe.Yellow.Hue),
		SaturationYellow: formatInt(recipe.Yellow.Saturation),
		LuminanceYellow:  formatInt(recipe.Yellow.Luminance),

		// HSL Adjustments - Green
		HueGreen:        formatInt(recipe.Green.Hue),
		SaturationGreen: formatInt(recipe.Green.Saturation),
		LuminanceGreen:  formatInt(recipe.Green.Luminance),

		// HSL Adjustments - Aqua
		HueAqua:        formatInt(recipe.Aqua.Hue),
		SaturationAqua: formatInt(recipe.Aqua.Saturation),
		LuminanceAqua:  formatInt(recipe.Aqua.Luminance),

		// HSL Adjustments - Blue
		HueBlue:        formatInt(recipe.Blue.Hue),
		SaturationBlue: formatInt(recipe.Blue.Saturation),
		LuminanceBlue:  formatInt(recipe.Blue.Luminance),

		// HSL Adjustments - Purple
		HuePurple:        formatInt(recipe.Purple.Hue),
		SaturationPurple: formatInt(recipe.Purple.Saturation),
		LuminancePurple:  formatInt(recipe.Purple.Luminance),

		// HSL Adjustments - Magenta
		HueMagenta:        formatInt(recipe.Magenta.Hue),
		SaturationMagenta: formatInt(recipe.Magenta.Saturation),
		LuminanceMagenta:  formatInt(recipe.Magenta.Luminance),

		// Split Toning
		SplitToningShadowHue:           formatInt(recipe.SplitShadowHue),
		SplitToningShadowSaturation:    formatInt(recipe.SplitShadowSaturation),
		SplitToningHighlightHue:        formatInt(recipe.SplitHighlightHue),
		SplitToningHighlightSaturation: formatInt(recipe.SplitHighlightSaturation),
		SplitToningBalance:             formatInt(recipe.SplitBalance),

		// Color Grading (Phase 2)
		ColorGradeHighlightHue: formatColorGradingZoneHue(recipe.ColorGrading, "highlights"),
		ColorGradeHighlightSat: formatColorGradingZoneChroma(recipe.ColorGrading, "highlights"),
		ColorGradeHighlightLum: formatColorGradingZoneBrightness(recipe.ColorGrading, "highlights"),
		ColorGradeMidtoneHue:   formatColorGradingZoneHue(recipe.ColorGrading, "midtone"),
		ColorGradeMidtoneSat:   formatColorGradingZoneChroma(recipe.ColorGrading, "midtone"),
		ColorGradeMidtoneLum:   formatColorGradingZoneBrightness(recipe.ColorGrading, "midtone"),
		ColorGradeShadowHue:    formatColorGradingZoneHue(recipe.ColorGrading, "shadows"),
		ColorGradeShadowSat:    formatColorGradingZoneChroma(recipe.ColorGrading, "shadows"),
		ColorGradeShadowLum:    formatColorGradingZoneBrightness(recipe.ColorGrading, "shadows"),
		ColorGradeBlending:     formatColorGradingBlending(recipe.ColorGrading),

		// Tone Curve
		ToneCurve: formatToneCurve(recipe.PointCurve),
	}

	// Construct complete XMP document with namespace declarations
	return &xmpDocWrapper{
		XMLNS:   nsAdobeMeta,
		XMPTool: "Adobe XMP Core 5.6-c140",
		RDF: rdfWrapper{
			XMLNS:       nsRDF,
			Description: desc,
		},
	}
}

// xmpDocWrapper is a wrapper struct for XML marshaling with proper namespace declarations
type xmpDocWrapper struct {
	XMLName xml.Name `xml:"x:xmpmeta"`
	XMLNS   string   `xml:"xmlns:x,attr"`
	XMPTool string   `xml:"x:xmptk,attr"`
	RDF     rdfWrapper
}

// rdfWrapper is a wrapper for RDF element with namespace declaration
type rdfWrapper struct {
	XMLName     xml.Name `xml:"rdf:RDF"`
	XMLNS       string   `xml:"xmlns:rdf,attr"`
	Description descriptionWrapper
}

// descriptionWrapper contains all photo editing parameters with crs: namespace prefix
type descriptionWrapper struct {
	XMLName xml.Name `xml:"rdf:Description"`
	XMLNS   string   `xml:"xmlns:crs,attr"`

	// Basic Adjustments
	Exposure2012   string `xml:"crs:Exposure2012,attr,omitempty"`
	Contrast2012   string `xml:"crs:Contrast2012,attr,omitempty"`
	Highlights2012 string `xml:"crs:Highlights2012,attr,omitempty"`
	Shadows2012    string `xml:"crs:Shadows2012,attr,omitempty"`
	Whites2012     string `xml:"crs:Whites2012,attr,omitempty"`
	Blacks2012     string `xml:"crs:Blacks2012,attr,omitempty"`

	// Color Parameters
	Saturation  string `xml:"crs:Saturation,attr,omitempty"`
	Vibrance    string `xml:"crs:Vibrance,attr,omitempty"`
	Clarity2012 string `xml:"crs:Clarity2012,attr,omitempty"`
	Sharpness   string `xml:"crs:Sharpness,attr,omitempty"`
	Temperature string `xml:"crs:Temperature,attr,omitempty"`
	Tint        string `xml:"crs:Tint,attr,omitempty"`

	// HSL Adjustments - Red
	HueRed        string `xml:"crs:HueRed,attr,omitempty"`
	SaturationRed string `xml:"crs:SaturationRed,attr,omitempty"`
	LuminanceRed  string `xml:"crs:LuminanceRed,attr,omitempty"`

	// HSL Adjustments - Orange
	HueOrange        string `xml:"crs:HueOrange,attr,omitempty"`
	SaturationOrange string `xml:"crs:SaturationOrange,attr,omitempty"`
	LuminanceOrange  string `xml:"crs:LuminanceOrange,attr,omitempty"`

	// HSL Adjustments - Yellow
	HueYellow        string `xml:"crs:HueYellow,attr,omitempty"`
	SaturationYellow string `xml:"crs:SaturationYellow,attr,omitempty"`
	LuminanceYellow  string `xml:"crs:LuminanceYellow,attr,omitempty"`

	// HSL Adjustments - Green
	HueGreen        string `xml:"crs:HueGreen,attr,omitempty"`
	SaturationGreen string `xml:"crs:SaturationGreen,attr,omitempty"`
	LuminanceGreen  string `xml:"crs:LuminanceGreen,attr,omitempty"`

	// HSL Adjustments - Aqua
	HueAqua        string `xml:"crs:HueAqua,attr,omitempty"`
	SaturationAqua string `xml:"crs:SaturationAqua,attr,omitempty"`
	LuminanceAqua  string `xml:"crs:LuminanceAqua,attr,omitempty"`

	// HSL Adjustments - Blue
	HueBlue        string `xml:"crs:HueBlue,attr,omitempty"`
	SaturationBlue string `xml:"crs:SaturationBlue,attr,omitempty"`
	LuminanceBlue  string `xml:"crs:LuminanceBlue,attr,omitempty"`

	// HSL Adjustments - Purple
	HuePurple        string `xml:"crs:HuePurple,attr,omitempty"`
	SaturationPurple string `xml:"crs:SaturationPurple,attr,omitempty"`
	LuminancePurple  string `xml:"crs:LuminancePurple,attr,omitempty"`

	// HSL Adjustments - Magenta
	HueMagenta        string `xml:"crs:HueMagenta,attr,omitempty"`
	SaturationMagenta string `xml:"crs:SaturationMagenta,attr,omitempty"`
	LuminanceMagenta  string `xml:"crs:LuminanceMagenta,attr,omitempty"`

	// Split Toning
	SplitToningShadowHue           string `xml:"crs:SplitToningShadowHue,attr,omitempty"`
	SplitToningShadowSaturation    string `xml:"crs:SplitToningShadowSaturation,attr,omitempty"`
	SplitToningHighlightHue        string `xml:"crs:SplitToningHighlightHue,attr,omitempty"`
	SplitToningHighlightSaturation string `xml:"crs:SplitToningHighlightSaturation,attr,omitempty"`
	SplitToningBalance             string `xml:"crs:SplitToningBalance,attr,omitempty"`

	// Color Grading (Phase 2) - Lightroom 2019+ Color Grading panel
	ColorGradeHighlightHue string `xml:"crs:ColorGradeHighlightHue,attr,omitempty"`
	ColorGradeHighlightSat string `xml:"crs:ColorGradeHighlightSat,attr,omitempty"`
	ColorGradeHighlightLum string `xml:"crs:ColorGradeHighlightLum,attr,omitempty"`
	ColorGradeMidtoneHue   string `xml:"crs:ColorGradeMidtoneHue,attr,omitempty"`
	ColorGradeMidtoneSat   string `xml:"crs:ColorGradeMidtoneSat,attr,omitempty"`
	ColorGradeMidtoneLum   string `xml:"crs:ColorGradeMidtoneLum,attr,omitempty"`
	ColorGradeShadowHue    string `xml:"crs:ColorGradeShadowHue,attr,omitempty"`
	ColorGradeShadowSat    string `xml:"crs:ColorGradeShadowSat,attr,omitempty"`
	ColorGradeShadowLum    string `xml:"crs:ColorGradeShadowLum,attr,omitempty"`
	ColorGradeBlending     string `xml:"crs:ColorGradeBlending,attr,omitempty"`
	ColorGradeGlobalHue    string `xml:"crs:ColorGradeGlobalHue,attr,omitempty"`
	ColorGradeGlobalSat    string `xml:"crs:ColorGradeGlobalSat,attr,omitempty"`
	ColorGradeGlobalLum    string `xml:"crs:ColorGradeGlobalLum,attr,omitempty"`

	// Tone Curve (stored as string, to be parsed separately if needed)
	ToneCurve string `xml:"crs:ToneCurve,attr,omitempty"`
}

// formatFloat formats a float64 value for XMP with 2 decimal places.
// Example: 1.5 → "1.50"
func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

// formatInt formats an integer value for XMP.
// Example: 25 → "25", -15 → "-15"
func formatInt(value int) string {
	return strconv.Itoa(value)
}

// formatTemperature formats a nullable temperature value for XMP.
// If temperature is nil, returns empty string (will be omitted from XML).
// Otherwise, formats as integer Kelvin value.
func formatTemperature(temp *int) string {
	if temp == nil {
		return ""
	}
	return strconv.Itoa(*temp)
}

// formatToneCurve formats a tone curve point array for XMP.
// Converts []ToneCurvePoint to comma-separated "input, output" pairs.
// Example: [{0,0}, {128,140}, {255,255}] → "0, 0 / 128, 140 / 255, 255"
// Returns empty string if points array is nil or empty.
func formatToneCurve(points []models.ToneCurvePoint) string {
	if len(points) == 0 {
		return ""
	}

	var result string
	for i, point := range points {
		if i > 0 {
			result += " / "
		}
		result += fmt.Sprintf("%d, %d", point.Input, point.Output)
	}
	return result
}

// formatColorGradingZoneHue formats the Hue value for a specific color grading zone.
// Returns empty string if ColorGrading is nil.
func formatColorGradingZoneHue(cg *models.ColorGrading, zone string) string {
	if cg == nil {
		return ""
	}
	switch zone {
	case "highlights":
		return formatInt(cg.Highlights.Hue)
	case "midtone":
		return formatInt(cg.Midtone.Hue)
	case "shadows":
		return formatInt(cg.Shadows.Hue)
	default:
		return ""
	}
}

// formatColorGradingZoneChroma formats the Chroma (Saturation) value for a specific color grading zone.
// Returns empty string if ColorGrading is nil.
func formatColorGradingZoneChroma(cg *models.ColorGrading, zone string) string {
	if cg == nil {
		return ""
	}
	switch zone {
	case "highlights":
		return formatInt(cg.Highlights.Chroma)
	case "midtone":
		return formatInt(cg.Midtone.Chroma)
	case "shadows":
		return formatInt(cg.Shadows.Chroma)
	default:
		return ""
	}
}

// formatColorGradingZoneBrightness formats the Brightness (Luminance) value for a specific color grading zone.
// Returns empty string if ColorGrading is nil.
func formatColorGradingZoneBrightness(cg *models.ColorGrading, zone string) string {
	if cg == nil {
		return ""
	}
	switch zone {
	case "highlights":
		return formatInt(cg.Highlights.Brightness)
	case "midtone":
		return formatInt(cg.Midtone.Brightness)
	case "shadows":
		return formatInt(cg.Shadows.Brightness)
	default:
		return ""
	}
}

// formatColorGradingBlending formats the Blending value for color grading.
// Returns empty string if ColorGrading is nil.
func formatColorGradingBlending(cg *models.ColorGrading) string {
	if cg == nil {
		return ""
	}
	return formatInt(cg.Blending)
}
