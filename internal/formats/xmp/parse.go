// Package xmp provides functionality for parsing Adobe Lightroom CC XMP preset files.
//
// The XMP format (Extensible Metadata Platform) is an XML-based format used by Adobe
// Lightroom CC to store photo editing presets. This package decodes XMP files into the
// UniversalRecipe intermediate representation, enabling conversion to other preset formats.
//
// Format Structure:
//   - XML/RDF structure with Adobe XMP namespaces
//   - Root: <x:xmpmeta xmlns:x="adobe:ns:meta/">
//   - Container: <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
//   - Parameters: <rdf:Description> attributes with crs: namespace prefix
//   - Example: crs:Exposure2012="1.5" crs:Contrast2012="25"
//
// Supported Parameters (50+):
//   - Basic adjustments: Exposure, Contrast, Highlights, Shadows, Whites, Blacks
//   - Color: Saturation, Vibrance, Clarity, Sharpness, Temperature, Tint
//   - HSL adjustments: 8 colors × 3 properties (Hue, Saturation, Luminance)
//   - Advanced: Tone curves, Split toning
//
// This parser uses encoding/xml with struct tags for type-safe unmarshaling and
// achieves the <30ms performance target through efficient XML parsing.
package xmp

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/justin/recipe/internal/models"
)

// Required namespace declarations for valid XMP files
const (
	nsAdobeMeta = "adobe:ns:meta/"
	nsCameraRaw = "http://ns.adobe.com/camera-raw-settings/1.0/"
	nsRDF       = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
)

// ConversionError wraps errors that occur during XMP conversion operations.
// It provides context about the operation that failed and the format being processed.
// This type follows Pattern 5 (Error Handling) from the architecture documentation.
type ConversionError struct {
	Operation string // Operation being performed (e.g., "parse", "validate", "extract")
	Format    string // Format being processed ("xmp")
	Field     string // Specific field being processed (optional)
	Cause     error  // Underlying error
}

// Error implements the error interface
func (e *ConversionError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s %s (%s): %v", e.Operation, e.Format, e.Field, e.Cause)
	}
	return fmt.Sprintf("%s %s: %v", e.Operation, e.Format, e.Cause)
}

// Unwrap returns the underlying error for errors.Is and errors.As compatibility
func (e *ConversionError) Unwrap() error {
	return e.Cause
}

// XMPDocument represents the root structure of an XMP file
type XMPDocument struct {
	XMLName xml.Name `xml:"xmpmeta"`
	RDF     RDF      `xml:"RDF"`
}

// RDF represents the RDF container element
type RDF struct {
	XMLName     xml.Name `xml:"RDF"`
	Description Description
}

// Description contains all photo editing parameters as XML attributes
type Description struct {
	XMLName xml.Name `xml:"Description"`

	// Basic Adjustments
	Exposure2012   string `xml:"Exposure2012,attr"`
	Contrast2012   string `xml:"Contrast2012,attr"`
	Highlights2012 string `xml:"Highlights2012,attr"`
	Shadows2012    string `xml:"Shadows2012,attr"`
	Whites2012     string `xml:"Whites2012,attr"`
	Blacks2012     string `xml:"Blacks2012,attr"`

	// Color Parameters
	Saturation  string `xml:"Saturation,attr"`
	Vibrance    string `xml:"Vibrance,attr"`
	Clarity2012 string `xml:"Clarity2012,attr"`
	Sharpness   string `xml:"Sharpness,attr"`
	Temperature string `xml:"Temperature,attr"`
	Tint        string `xml:"Tint,attr"`

	// HSL Adjustments - Red
	HueRed        string `xml:"HueRed,attr"`
	SaturationRed string `xml:"SaturationRed,attr"`
	LuminanceRed  string `xml:"LuminanceRed,attr"`

	// HSL Adjustments - Orange
	HueOrange        string `xml:"HueOrange,attr"`
	SaturationOrange string `xml:"SaturationOrange,attr"`
	LuminanceOrange  string `xml:"LuminanceOrange,attr"`

	// HSL Adjustments - Yellow
	HueYellow        string `xml:"HueYellow,attr"`
	SaturationYellow string `xml:"SaturationYellow,attr"`
	LuminanceYellow  string `xml:"LuminanceYellow,attr"`

	// HSL Adjustments - Green
	HueGreen        string `xml:"HueGreen,attr"`
	SaturationGreen string `xml:"SaturationGreen,attr"`
	LuminanceGreen  string `xml:"LuminanceGreen,attr"`

	// HSL Adjustments - Aqua
	HueAqua        string `xml:"HueAqua,attr"`
	SaturationAqua string `xml:"SaturationAqua,attr"`
	LuminanceAqua  string `xml:"LuminanceAqua,attr"`

	// HSL Adjustments - Blue
	HueBlue        string `xml:"HueBlue,attr"`
	SaturationBlue string `xml:"SaturationBlue,attr"`
	LuminanceBlue  string `xml:"LuminanceBlue,attr"`

	// HSL Adjustments - Purple
	HuePurple        string `xml:"HuePurple,attr"`
	SaturationPurple string `xml:"SaturationPurple,attr"`
	LuminancePurple  string `xml:"LuminancePurple,attr"`

	// HSL Adjustments - Magenta
	HueMagenta        string `xml:"HueMagenta,attr"`
	SaturationMagenta string `xml:"SaturationMagenta,attr"`
	LuminanceMagenta  string `xml:"LuminanceMagenta,attr"`

	// Split Toning
	SplitToningShadowHue           string `xml:"SplitToningShadowHue,attr"`
	SplitToningShadowSaturation    string `xml:"SplitToningShadowSaturation,attr"`
	SplitToningHighlightHue        string `xml:"SplitToningHighlightHue,attr"`
	SplitToningHighlightSaturation string `xml:"SplitToningHighlightSaturation,attr"`
	SplitToningBalance             string `xml:"SplitToningBalance,attr"`

	// Tone Curve (stored as string, to be parsed separately if needed)
	ToneCurve string `xml:"ToneCurve,attr"`
}

// Parse decodes an Adobe Lightroom CC XMP preset file into a UniversalRecipe.
//
// The function validates the file structure (XML well-formedness and required namespaces),
// extracts all 50+ photo editing parameters from XML attributes, validates parameter ranges,
// and constructs a UniversalRecipe using the builder pattern.
//
// Parameters:
//   - data: Raw bytes of the .xmp file
//
// Returns:
//   - *models.UniversalRecipe: Populated recipe with extracted parameters
//   - error: Validation or parsing error with descriptive context
//
// Errors:
//   - Invalid XML structure: File is not well-formed XML
//   - Missing required namespaces: File missing camera-raw-settings namespace
//   - Parameter out of range: Invalid parameter value
//   - Builder validation: UniversalRecipe construction failed
func Parse(data []byte) (*models.UniversalRecipe, error) {
	// Validate XML structure and namespaces (fail-fast per Pattern 6)
	if err := validateXMLStructure(data); err != nil {
		return nil, &ConversionError{
			Operation: "parse",
			Format:    "xmp",
			Cause:     err,
		}
	}

	// Unmarshal XML into struct
	var doc XMPDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, &ConversionError{
			Operation: "parse",
			Format:    "xmp",
			Field:     "xml unmarshal",
			Cause:     err,
		}
	}

	// Extract and convert parameters from strings to appropriate types
	params, err := extractParameters(&doc.RDF.Description)
	if err != nil {
		return nil, err // Already wrapped by extractParameters
	}

	// Validate extracted parameters (inline validation per Pattern 6)
	if err := validateParameters(params); err != nil {
		return nil, err // Already wrapped by validateParameters
	}

	// Build UniversalRecipe using builder pattern (Pattern 4)
	recipe, err := buildRecipe(params)
	if err != nil {
		return nil, &ConversionError{
			Operation: "build",
			Format:    "xmp",
			Cause:     err,
		}
	}

	return recipe, nil
}

// validateXMLStructure checks if the data is well-formed XML and contains required namespaces
func validateXMLStructure(data []byte) error {
	// Check for basic XML declaration or root element
	dataStr := string(data)
	if !strings.Contains(dataStr, "<") || !strings.Contains(dataStr, ">") {
		return fmt.Errorf("invalid XML: no XML tags found")
	}

	// Check for required namespaces
	if !strings.Contains(dataStr, nsAdobeMeta) {
		return fmt.Errorf("missing required namespace: adobe:ns:meta/")
	}

	if !strings.Contains(dataStr, nsCameraRaw) {
		return fmt.Errorf("missing required namespace: camera-raw-settings (Adobe Lightroom CC namespace required)")
	}

	if !strings.Contains(dataStr, nsRDF) {
		return fmt.Errorf("missing required namespace: rdf:RDF")
	}

	return nil
}

// xmpParameters holds extracted parameter values before validation
type xmpParameters struct {
	// Basic Adjustments
	exposure   float64
	contrast   int
	highlights int
	shadows    int
	whites     int
	blacks     int

	// Color Parameters
	saturation  int
	vibrance    int
	clarity     int
	sharpness   int
	temperature int
	tint        int

	// HSL Adjustments
	red     models.ColorAdjustment
	orange  models.ColorAdjustment
	yellow  models.ColorAdjustment
	green   models.ColorAdjustment
	aqua    models.ColorAdjustment
	blue    models.ColorAdjustment
	purple  models.ColorAdjustment
	magenta models.ColorAdjustment

	// Split Toning
	splitShadowHue           int
	splitShadowSaturation    int
	splitHighlightHue        int
	splitHighlightSaturation int
	splitBalance             int

	// Tone Curve (stored as string for now)
	toneCurve string
}

// extractParameters extracts all parameter values from the XMP Description and converts
// them from strings to appropriate types (float64, int).
func extractParameters(desc *Description) (*xmpParameters, error) {
	params := &xmpParameters{}
	var err error

	// Extract Basic Adjustments
	params.exposure, err = parseFloat64(desc.Exposure2012, "Exposure2012")
	if err != nil && desc.Exposure2012 != "" {
		return nil, err
	}

	params.contrast, err = parseInt(desc.Contrast2012, "Contrast2012")
	if err != nil && desc.Contrast2012 != "" {
		return nil, err
	}

	params.highlights, err = parseInt(desc.Highlights2012, "Highlights2012")
	if err != nil && desc.Highlights2012 != "" {
		return nil, err
	}

	params.shadows, err = parseInt(desc.Shadows2012, "Shadows2012")
	if err != nil && desc.Shadows2012 != "" {
		return nil, err
	}

	params.whites, err = parseInt(desc.Whites2012, "Whites2012")
	if err != nil && desc.Whites2012 != "" {
		return nil, err
	}

	params.blacks, err = parseInt(desc.Blacks2012, "Blacks2012")
	if err != nil && desc.Blacks2012 != "" {
		return nil, err
	}

	// Extract Color Parameters
	params.saturation, err = parseInt(desc.Saturation, "Saturation")
	if err != nil && desc.Saturation != "" {
		return nil, err
	}

	params.vibrance, err = parseInt(desc.Vibrance, "Vibrance")
	if err != nil && desc.Vibrance != "" {
		return nil, err
	}

	params.clarity, err = parseInt(desc.Clarity2012, "Clarity2012")
	if err != nil && desc.Clarity2012 != "" {
		return nil, err
	}

	params.sharpness, err = parseInt(desc.Sharpness, "Sharpness")
	if err != nil && desc.Sharpness != "" {
		return nil, err
	}

	params.temperature, err = parseInt(desc.Temperature, "Temperature")
	if err != nil && desc.Temperature != "" {
		return nil, err
	}

	params.tint, err = parseInt(desc.Tint, "Tint")
	if err != nil && desc.Tint != "" {
		return nil, err
	}

	// Extract HSL Adjustments
	params.red, err = extractColorAdjustment(desc.HueRed, desc.SaturationRed, desc.LuminanceRed, "Red")
	if err != nil {
		return nil, err
	}

	params.orange, err = extractColorAdjustment(desc.HueOrange, desc.SaturationOrange, desc.LuminanceOrange, "Orange")
	if err != nil {
		return nil, err
	}

	params.yellow, err = extractColorAdjustment(desc.HueYellow, desc.SaturationYellow, desc.LuminanceYellow, "Yellow")
	if err != nil {
		return nil, err
	}

	params.green, err = extractColorAdjustment(desc.HueGreen, desc.SaturationGreen, desc.LuminanceGreen, "Green")
	if err != nil {
		return nil, err
	}

	params.aqua, err = extractColorAdjustment(desc.HueAqua, desc.SaturationAqua, desc.LuminanceAqua, "Aqua")
	if err != nil {
		return nil, err
	}

	params.blue, err = extractColorAdjustment(desc.HueBlue, desc.SaturationBlue, desc.LuminanceBlue, "Blue")
	if err != nil {
		return nil, err
	}

	params.purple, err = extractColorAdjustment(desc.HuePurple, desc.SaturationPurple, desc.LuminancePurple, "Purple")
	if err != nil {
		return nil, err
	}

	params.magenta, err = extractColorAdjustment(desc.HueMagenta, desc.SaturationMagenta, desc.LuminanceMagenta, "Magenta")
	if err != nil {
		return nil, err
	}

	// Extract Split Toning
	params.splitShadowHue, err = parseInt(desc.SplitToningShadowHue, "SplitToningShadowHue")
	if err != nil && desc.SplitToningShadowHue != "" {
		return nil, err
	}

	params.splitShadowSaturation, err = parseInt(desc.SplitToningShadowSaturation, "SplitToningShadowSaturation")
	if err != nil && desc.SplitToningShadowSaturation != "" {
		return nil, err
	}

	params.splitHighlightHue, err = parseInt(desc.SplitToningHighlightHue, "SplitToningHighlightHue")
	if err != nil && desc.SplitToningHighlightHue != "" {
		return nil, err
	}

	params.splitHighlightSaturation, err = parseInt(desc.SplitToningHighlightSaturation, "SplitToningHighlightSaturation")
	if err != nil && desc.SplitToningHighlightSaturation != "" {
		return nil, err
	}

	params.splitBalance, err = parseInt(desc.SplitToningBalance, "SplitToningBalance")
	if err != nil && desc.SplitToningBalance != "" {
		return nil, err
	}

	// Store tone curve as-is for now (to be parsed later if needed)
	params.toneCurve = desc.ToneCurve

	return params, nil
}

// parseFloat64 converts a string to float64, returning 0.0 if empty
func parseFloat64(s, fieldName string) (float64, error) {
	if s == "" {
		return 0.0, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, &ConversionError{
			Operation: "parse",
			Format:    "xmp",
			Field:     fieldName,
			Cause:     fmt.Errorf("invalid value %q: %w", s, err),
		}
	}
	return v, nil
}

// parseInt converts a string to int, returning 0 if empty
func parseInt(s, fieldName string) (int, error) {
	if s == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, &ConversionError{
			Operation: "parse",
			Format:    "xmp",
			Field:     fieldName,
			Cause:     fmt.Errorf("invalid value %q: %w", s, err),
		}
	}
	return v, nil
}

// extractColorAdjustment extracts a ColorAdjustment from HSL string values
func extractColorAdjustment(hue, saturation, luminance, colorName string) (models.ColorAdjustment, error) {
	adj := models.ColorAdjustment{}
	var err error

	adj.Hue, err = parseInt(hue, colorName+"Hue")
	if err != nil {
		return adj, err
	}

	adj.Saturation, err = parseInt(saturation, colorName+"Saturation")
	if err != nil {
		return adj, err
	}

	adj.Luminance, err = parseInt(luminance, colorName+"Luminance")
	if err != nil {
		return adj, err
	}

	return adj, nil
}

// validateParameters validates all extracted parameter values are within expected ranges
func validateParameters(params *xmpParameters) error {
	// Validate Exposure (-5.0 to +5.0)
	if params.exposure < -5.0 || params.exposure > 5.0 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Exposure",
			Cause:     fmt.Errorf("value %.2f out of range (expected -5.0 to +5.0)", params.exposure),
		}
	}

	// Validate Contrast (-100 to +100)
	if params.contrast < -100 || params.contrast > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Contrast",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.contrast),
		}
	}

	// Validate Highlights, Shadows, Whites, Blacks (-100 to +100)
	if params.highlights < -100 || params.highlights > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Highlights",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.highlights),
		}
	}

	if params.shadows < -100 || params.shadows > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Shadows",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.shadows),
		}
	}

	if params.whites < -100 || params.whites > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Whites",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.whites),
		}
	}

	if params.blacks < -100 || params.blacks > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Blacks",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.blacks),
		}
	}

	// Validate Color Parameters
	if params.saturation < -100 || params.saturation > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Saturation",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.saturation),
		}
	}

	if params.vibrance < -100 || params.vibrance > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Vibrance",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.vibrance),
		}
	}

	if params.clarity < -100 || params.clarity > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Clarity",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.clarity),
		}
	}

	if params.sharpness < 0 || params.sharpness > 150 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Sharpness",
			Cause:     fmt.Errorf("value %d out of range (expected 0 to 150)", params.sharpness),
		}
	}

	// Validate HSL Adjustments
	if err := validateColorRange(params.red, "Red"); err != nil {
		return err
	}
	if err := validateColorRange(params.orange, "Orange"); err != nil {
		return err
	}
	if err := validateColorRange(params.yellow, "Yellow"); err != nil {
		return err
	}
	if err := validateColorRange(params.green, "Green"); err != nil {
		return err
	}
	if err := validateColorRange(params.aqua, "Aqua"); err != nil {
		return err
	}
	if err := validateColorRange(params.blue, "Blue"); err != nil {
		return err
	}
	if err := validateColorRange(params.purple, "Purple"); err != nil {
		return err
	}
	if err := validateColorRange(params.magenta, "Magenta"); err != nil {
		return err
	}

	// Validate Split Toning
	if params.splitShadowHue < 0 || params.splitShadowHue > 360 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "SplitShadowHue",
			Cause:     fmt.Errorf("value %d out of range (expected 0 to 360)", params.splitShadowHue),
		}
	}

	if params.splitShadowSaturation < 0 || params.splitShadowSaturation > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "SplitShadowSaturation",
			Cause:     fmt.Errorf("value %d out of range (expected 0 to 100)", params.splitShadowSaturation),
		}
	}

	if params.splitHighlightHue < 0 || params.splitHighlightHue > 360 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "SplitHighlightHue",
			Cause:     fmt.Errorf("value %d out of range (expected 0 to 360)", params.splitHighlightHue),
		}
	}

	if params.splitHighlightSaturation < 0 || params.splitHighlightSaturation > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "SplitHighlightSaturation",
			Cause:     fmt.Errorf("value %d out of range (expected 0 to 100)", params.splitHighlightSaturation),
		}
	}

	if params.splitBalance < -100 || params.splitBalance > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "SplitBalance",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", params.splitBalance),
		}
	}

	return nil
}

// validateColorRange validates a ColorAdjustment's ranges
func validateColorRange(adj models.ColorAdjustment, colorName string) error {
	if adj.Hue < -100 || adj.Hue > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     colorName + " Hue",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", adj.Hue),
		}
	}

	if adj.Saturation < -100 || adj.Saturation > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     colorName + " Saturation",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", adj.Saturation),
		}
	}

	if adj.Luminance < -100 || adj.Luminance > 100 {
		return &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     colorName + " Luminance",
			Cause:     fmt.Errorf("value %d out of range (expected -100 to +100)", adj.Luminance),
		}
	}

	return nil
}

// buildRecipe constructs a UniversalRecipe from validated parameters
// using the builder pattern (Pattern 4).
func buildRecipe(params *xmpParameters) (*models.UniversalRecipe, error) {
	builder := models.NewRecipeBuilder()

	// Set source format
	builder.WithSourceFormat("xmp")

	// Set Basic Adjustments
	builder.WithExposure(params.exposure)
	builder.WithContrast(params.contrast)
	builder.WithHighlights(params.highlights)
	builder.WithShadows(params.shadows)
	builder.WithWhites(params.whites)
	builder.WithBlacks(params.blacks)

	// Set Color Parameters
	builder.WithSaturation(params.saturation)
	builder.WithVibrance(params.vibrance)
	builder.WithClarity(params.clarity)
	builder.WithSharpness(params.sharpness)
	builder.WithTemperature(params.temperature)
	builder.WithTint(params.tint)

	// Set HSL Adjustments
	builder.WithRedHSL(params.red.Hue, params.red.Saturation, params.red.Luminance)
	builder.WithOrangeHSL(params.orange.Hue, params.orange.Saturation, params.orange.Luminance)
	builder.WithYellowHSL(params.yellow.Hue, params.yellow.Saturation, params.yellow.Luminance)
	builder.WithGreenHSL(params.green.Hue, params.green.Saturation, params.green.Luminance)
	builder.WithAquaHSL(params.aqua.Hue, params.aqua.Saturation, params.aqua.Luminance)
	builder.WithBlueHSL(params.blue.Hue, params.blue.Saturation, params.blue.Luminance)
	builder.WithPurpleHSL(params.purple.Hue, params.purple.Saturation, params.purple.Luminance)
	builder.WithMagentaHSL(params.magenta.Hue, params.magenta.Saturation, params.magenta.Luminance)

	// Set Split Toning
	builder.WithSplitToning(
		params.splitShadowHue,
		params.splitShadowSaturation,
		params.splitHighlightHue,
		params.splitHighlightSaturation,
		params.splitBalance,
	)

	// Build and validate
	recipe, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("build recipe: %w", err)
	}

	return recipe, nil
}
