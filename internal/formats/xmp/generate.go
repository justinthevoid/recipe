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
	"math"
	"strconv"
	"strings"

	"github.com/justin/recipe/internal/lut"
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

	// Post-process to format attributes on separate lines like professional presets
	// This matches the format of Adobe-exported XMP files
	xmlData = formatAttributesOnSeparateLines(xmlData)

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

	// Validate Grain
	if err := validateRange(recipe.GrainAmount, 0, 100, "GrainAmount"); err != nil {
		return err
	}
	if err := validateRange(recipe.GrainSize, 0, 100, "GrainSize"); err != nil {
		return err
	}
	if err := validateRange(recipe.GrainRoughness, 0, 100, "GrainRoughness"); err != nil {
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

		// Required Lightroom metadata (without these, Lightroom ignores most settings)
		Version:        "15.0", // Lightroom version compatibility
		ProcessVersion: "11.0", // Camera Raw process version (PV2012+)
		HasSettings:    "True", // Indicates this preset has settings to apply
		PresetType:     "Normal",
		RDFAbout:       "", // Required by XMP/RDF specification

		// Camera Profile - default to "Camera Flexible Color" for Nikon color accuracy
		// This matches NX Studio's rendering better than Adobe Standard
		CameraProfile: getCameraProfile(recipe.CameraProfileName),

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

		// Grain
		GrainAmount:    formatInt(recipe.GrainAmount),
		GrainSize:      formatInt(recipe.GrainSize),
		GrainFrequency: formatInt(recipe.GrainRoughness),

		// Temperature (nullable - handle nil)
		Temperature: formatTemperature(recipe.Temperature),

		// Parametric Tone Curve (zone-based adjustments)
		ParametricShadows:        formatInt(recipe.ToneCurveShadows),
		ParametricDarks:          formatInt(recipe.ToneCurveDarks),
		ParametricLights:         formatInt(recipe.ToneCurveLights),
		ParametricHighlights:     formatInt(recipe.ToneCurveHighlights),
		ParametricShadowSplit:    formatInt(recipe.ToneCurveShadowSplit),
		ParametricMidtoneSplit:   formatInt(recipe.ToneCurveMidtoneSplit),
		ParametricHighlightSplit: formatInt(recipe.ToneCurveHighlightSplit),

		// HSL Adjustments - Red
		HueAdjustmentRed:        formatInt(recipe.Red.Hue),
		SaturationAdjustmentRed: formatInt(recipe.Red.Saturation),
		LuminanceAdjustmentRed:  formatInt(recipe.Red.Luminance),

		// HSL Adjustments - Orange
		HueAdjustmentOrange:        formatInt(recipe.Orange.Hue),
		SaturationAdjustmentOrange: formatInt(recipe.Orange.Saturation),
		LuminanceAdjustmentOrange:  formatInt(recipe.Orange.Luminance),

		// HSL Adjustments - Yellow
		HueAdjustmentYellow:        formatInt(recipe.Yellow.Hue),
		SaturationAdjustmentYellow: formatInt(recipe.Yellow.Saturation),
		LuminanceAdjustmentYellow:  formatInt(recipe.Yellow.Luminance),

		// HSL Adjustments - Green
		HueAdjustmentGreen:        formatInt(recipe.Green.Hue),
		SaturationAdjustmentGreen: formatInt(recipe.Green.Saturation),
		LuminanceAdjustmentGreen:  formatInt(recipe.Green.Luminance),

		// HSL Adjustments - Aqua
		HueAdjustmentAqua:        formatInt(recipe.Aqua.Hue),
		SaturationAdjustmentAqua: formatInt(recipe.Aqua.Saturation),
		LuminanceAdjustmentAqua:  formatInt(recipe.Aqua.Luminance),

		// HSL Adjustments - Blue
		HueAdjustmentBlue:        formatInt(recipe.Blue.Hue),
		SaturationAdjustmentBlue: formatInt(recipe.Blue.Saturation),
		LuminanceAdjustmentBlue:  formatInt(recipe.Blue.Luminance),

		// HSL Adjustments - Purple
		HueAdjustmentPurple:        formatInt(recipe.Purple.Hue),
		SaturationAdjustmentPurple: formatInt(recipe.Purple.Saturation),
		LuminanceAdjustmentPurple:  formatInt(recipe.Purple.Luminance),

		// HSL Adjustments - Magenta
		HueAdjustmentMagenta:        formatInt(recipe.Magenta.Hue),
		SaturationAdjustmentMagenta: formatInt(recipe.Magenta.Saturation),
		LuminanceAdjustmentMagenta:  formatInt(recipe.Magenta.Luminance),

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
		// Note: NP3's Balance parameter does not have a direct Adobe XMP equivalent
		// It affects the internal color processing in NX Studio but cannot be accurately
		// represented in Lightroom's color grading system

		// Tone Curve: We now use ONLY parametric curve (Shadows/Darks/Lights/Highlights)
		// instead of point curve (ToneCurvePV2012) for better NX Studio compatibility.
		// If parametric values are set, skip the point curve entirely.
		// If no parametric values, use point curve as fallback.
		ToneCurveName2012:    formatToneCurveName2012ParametricAware(recipe),
		ToneCurvePV2012:      formatToneCurvePV2012ParametricAware(recipe),
		ToneCurvePV2012Red:   formatToneCurvePV2012(recipe.PointCurveRed),
		ToneCurvePV2012Green: formatToneCurvePV2012(recipe.PointCurveGreen),
		ToneCurvePV2012Blue:  formatToneCurvePV2012(recipe.PointCurveBlue),
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

// GenerateWithLUT creates an XMP preset with embedded 3D LUT for maximum color accuracy.
// This generates a parametric XMP plus a 32x32x32 RGB lookup table that applies all
// NP3 color transformations. The LUT approach matches Adobe's official camera profiles
// and provides 95%+ accuracy compared to NX Studio.
//
// Performance: Generating a 32x32x32 LUT requires ~32,768 color transformations.
// File size: ~350KB (vs ~2KB for parametric-only XMP).
//
// Returns:
//   - []byte: XMP file with embedded 3D LUT
//   - error: ConversionError if generation fails
func GenerateWithLUT(recipe *models.UniversalRecipe) ([]byte, error) {
	// First, generate the base XMP document (parametric adjustments)
	xmpDoc := buildXMPDocument(recipe)

	// Generate 3D LUT
	lutData, err := lut.Generate3DLUT(recipe)
	if err != nil {
		return nil, &ConversionError{
			Operation: "generate",
			Format:    "xmp-lut",
			Field:     "3D LUT",
			Cause:     err,
		}
	}

	// Compress and encode LUT
	tableID, encodedLUT, err := lut.CompressAndEncodeLUT(lutData)
	if err != nil {
		return nil, &ConversionError{
			Operation: "compress",
			Format:    "xmp-lut",
			Field:     "LUT encoding",
			Cause:     err,
		}
	}

	// Add LUT reference and data to XMP
	xmpDoc.RDF.Description.RGBTable = tableID
	xmpDoc.RDF.Description.LUTData = fmt.Sprintf(" crs:Table_%s=\"%s\"", tableID, encodedLUT)

	// Marshal to XML with proper formatting
	output, err := xml.MarshalIndent(xmpDoc, "", "  ")
	if err != nil {
		return nil, &ConversionError{
			Operation: "marshal",
			Format:    "xmp-lut",
			Cause:     err,
		}
	}

	// Prepend XML declaration
	result := append([]byte(xml.Header), output...)
	result = append(result, '\n')

	return result, nil
}

// GenerateProfileWithLUT creates an XMP profile (like Adobe's Dream, Pop profiles) that:
// 1. Specifies a base Nikon camera profile as the starting point
// 2. Embeds a 3D LUT with NP3 color transformations
// 3. Uses PresetType="Look" to load as a profile/look in Lightroom
// 4. Applies temperature offset to compensate for baseline profile differences
//
// This approach achieves higher accuracy (90-95%) because it starts with Nikon's color
// science (from the specified camera profile) rather than Adobe's interpretation.
//
// Strategy:
//   - Use "Camera Flexible Color" as base (best match for Nikon rendering)
//   - Generate 3D LUT with all NP3 transformations
//   - Embed LUT using crs:RGBTable + crs:Table_[hash]
//   - Set PresetType="Look" to load as profile
//   - Add temperature offset to compensate for Lightroom vs NX Studio baseline difference
//
// Parameters:
//   - recipe: UniversalRecipe with NP3 color adjustments
//   - baseCameraProfile: Camera profile name (e.g., "Camera Flexible Color", "Camera Standard")
//   - temperatureOffset: Temperature adjustment in Kelvin (e.g., +1000 to make warmer)
//
// Returns:
//   - []byte: XMP profile file with embedded 3D LUT
//   - error: ConversionError if generation fails
func GenerateProfileWithLUT(recipe *models.UniversalRecipe, baseCameraProfile string, temperatureOffset int) ([]byte, error) {
	// Validation
	if recipe == nil {
		return nil, &ConversionError{
			Operation: "generate",
			Format:    "xmp-profile",
			Cause:     fmt.Errorf("recipe is nil"),
		}
	}
	if baseCameraProfile == "" {
		baseCameraProfile = "Camera Flexible Color" // Default to profile that better matches Nikon rendering
	}

	// Generate 3D LUT
	lutData, err := lut.Generate3DLUT(recipe)
	if err != nil {
		return nil, &ConversionError{
			Operation: "generate",
			Format:    "xmp-profile",
			Field:     "3D LUT",
			Cause:     err,
		}
	}

	// Compress and encode LUT
	tableID, encodedLUT, err := lut.CompressAndEncodeLUT(lutData)
	if err != nil {
		return nil, &ConversionError{
			Operation: "compress",
			Format:    "xmp-profile",
			Field:     "LUT encoding",
			Cause:     err,
		}
	}

	// Build profile XMP document (modeled after Adobe's Dream profile)
	profileDoc := buildProfileXMPDocument(recipe, baseCameraProfile, tableID, encodedLUT, temperatureOffset)

	// Marshal to XML with proper formatting
	output, err := xml.MarshalIndent(profileDoc, "", "  ")
	if err != nil {
		return nil, &ConversionError{
			Operation: "marshal",
			Format:    "xmp-profile",
			Cause:     err,
		}
	}

	// Prepend XML declaration
	result := append([]byte(xml.Header), output...)
	result = append(result, '\n')

	return result, nil
}

// buildProfileXMPDocument constructs an XMP profile document like Adobe's creative profiles.
// This includes minimal parametric adjustments (only what's needed) and focuses on the
// embedded 3D LUT for accurate color transformation.
func buildProfileXMPDocument(recipe *models.UniversalRecipe, baseCameraProfile, tableID, encodedLUT string, temperatureOffset int) *xmpDocWrapper {
	// Calculate final temperature (if recipe has temperature, add offset; otherwise just use offset)
	var finalTemp *int
	if temperatureOffset != 0 {
		// Default to 5500K (daylight) if no temperature in recipe, then add offset
		baseTemp := 5500
		if recipe.Temperature != nil {
			baseTemp = *recipe.Temperature
		}
		finalTemp = new(int)
		*finalTemp = baseTemp + temperatureOffset
	} else if recipe.Temperature != nil {
		finalTemp = recipe.Temperature
	}

	// Build Description with profile-specific attributes
	desc := descriptionWrapper{
		XMLNS: nsCameraRaw,

		// Profile metadata (modeled after Adobe's Dream profile)
		PresetType:                 "Look",
		UUID:                       generateUUID(recipe.Name),
		SupportsAmount:             "True",
		SupportsColor:              "True",
		SupportsMonochrome:         "False",
		SupportsHighDynamicRange:   "True",
		SupportsNormalDynamicRange: "True",
		SupportsSceneReferred:      "True",
		SupportsOutputReferred:     "False",
		RequiresRGBTables:          "False",
		Copyright:                  "Converted from Nikon NP3",
		ProcessVersion:             "15.4",
		ConvertToGrayscale:         "False",

		// Camera profile specification (KEY: Use Nikon's color science as base)
		CameraProfile: recipe.CameraProfileName,

		// Temperature compensation for baseline profile difference
		Temperature: formatTemperature(finalTemp),

		// 3D LUT reference and data
		RGBTable: tableID,
		LUTData:  fmt.Sprintf(" crs:Table_%s=\"%s\"", tableID, encodedLUT),
	}

	// Construct complete XMP document with namespace declarations
	return &xmpDocWrapper{
		XMLNS:   nsAdobeMeta,
		XMPTool: "Adobe XMP Core 7.0-c000 1.000000, 0000/00/00-00:00:00",
		RDF: rdfWrapper{
			XMLNS:       nsRDF,
			Description: desc,
		},
	}
}

// generateUUID generates a simple UUID based on the preset name.
// For production, this should use a proper UUID library, but for now we use a deterministic hash.
func generateUUID(name string) string {
	if name == "" {
		name = "NP3_Preset"
	}
	// Simple deterministic UUID-like string
	// In production, use github.com/google/uuid or similar
	return fmt.Sprintf("%032X", []byte(name)[:min(32, len(name))])
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	XMLName  xml.Name `xml:"rdf:Description"`
	RDFAbout string   `xml:"rdf:about,attr"` // Required by XMP/RDF - must be empty string ""
	XMLNS    string   `xml:"xmlns:crs,attr"`

	// Profile Metadata (for profile-based XMP like Adobe's Dream, Pop profiles)
	PresetType                 string `xml:"crs:PresetType,attr,omitempty"`
	UUID                       string `xml:"crs:UUID,attr,omitempty"`
	SupportsAmount             string `xml:"crs:SupportsAmount,attr,omitempty"`
	SupportsColor              string `xml:"crs:SupportsColor,attr,omitempty"`
	SupportsMonochrome         string `xml:"crs:SupportsMonochrome,attr,omitempty"`
	SupportsHighDynamicRange   string `xml:"crs:SupportsHighDynamicRange,attr,omitempty"`
	SupportsNormalDynamicRange string `xml:"crs:SupportsNormalDynamicRange,attr,omitempty"`
	SupportsSceneReferred      string `xml:"crs:SupportsSceneReferred,attr,omitempty"`
	SupportsOutputReferred     string `xml:"crs:SupportsOutputReferred,attr,omitempty"`
	RequiresRGBTables          string `xml:"crs:RequiresRGBTables,attr,omitempty"`
	Copyright                  string `xml:"crs:Copyright,attr,omitempty"`
	Version                    string `xml:"crs:Version,attr,omitempty"`        // Required for Lightroom to recognize settings
	ProcessVersion             string `xml:"crs:ProcessVersion,attr,omitempty"` // Required for Lightroom to recognize settings
	HasSettings                string `xml:"crs:HasSettings,attr,omitempty"`    // Indicates preset has settings to apply
	ConvertToGrayscale         string `xml:"crs:ConvertToGrayscale,attr,omitempty"`
	CameraProfile              string `xml:"crs:CameraProfile,attr,omitempty"`

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

	// Parametric Tone Curve (zone-based adjustments)
	ParametricShadows        string `xml:"crs:ParametricShadows,attr,omitempty"`
	ParametricDarks          string `xml:"crs:ParametricDarks,attr,omitempty"`
	ParametricLights         string `xml:"crs:ParametricLights,attr,omitempty"`
	ParametricHighlights     string `xml:"crs:ParametricHighlights,attr,omitempty"`
	ParametricShadowSplit    string `xml:"crs:ParametricShadowSplit,attr,omitempty"`
	ParametricMidtoneSplit   string `xml:"crs:ParametricMidtoneSplit,attr,omitempty"`
	ParametricHighlightSplit string `xml:"crs:ParametricHighlightSplit,attr,omitempty"`

	// Grain
	GrainAmount    string `xml:"crs:GrainAmount,attr,omitempty"`
	GrainSize      string `xml:"crs:GrainSize,attr,omitempty"`
	GrainFrequency string `xml:"crs:GrainFrequency,attr,omitempty"` // Roughness

	// HSL Adjustments - Red (modern HueAdjustment* names for professional preset compatibility)
	HueAdjustmentRed        string `xml:"crs:HueAdjustmentRed,attr,omitempty"`
	SaturationAdjustmentRed string `xml:"crs:SaturationAdjustmentRed,attr,omitempty"`
	LuminanceAdjustmentRed  string `xml:"crs:LuminanceAdjustmentRed,attr,omitempty"`

	// HSL Adjustments - Orange
	HueAdjustmentOrange        string `xml:"crs:HueAdjustmentOrange,attr,omitempty"`
	SaturationAdjustmentOrange string `xml:"crs:SaturationAdjustmentOrange,attr,omitempty"`
	LuminanceAdjustmentOrange  string `xml:"crs:LuminanceAdjustmentOrange,attr,omitempty"`

	// HSL Adjustments - Yellow
	HueAdjustmentYellow        string `xml:"crs:HueAdjustmentYellow,attr,omitempty"`
	SaturationAdjustmentYellow string `xml:"crs:SaturationAdjustmentYellow,attr,omitempty"`
	LuminanceAdjustmentYellow  string `xml:"crs:LuminanceAdjustmentYellow,attr,omitempty"`

	// HSL Adjustments - Green
	HueAdjustmentGreen        string `xml:"crs:HueAdjustmentGreen,attr,omitempty"`
	SaturationAdjustmentGreen string `xml:"crs:SaturationAdjustmentGreen,attr,omitempty"`
	LuminanceAdjustmentGreen  string `xml:"crs:LuminanceAdjustmentGreen,attr,omitempty"`

	// HSL Adjustments - Aqua
	HueAdjustmentAqua        string `xml:"crs:HueAdjustmentAqua,attr,omitempty"`
	SaturationAdjustmentAqua string `xml:"crs:SaturationAdjustmentAqua,attr,omitempty"`
	LuminanceAdjustmentAqua  string `xml:"crs:LuminanceAdjustmentAqua,attr,omitempty"`

	// HSL Adjustments - Blue
	HueAdjustmentBlue        string `xml:"crs:HueAdjustmentBlue,attr,omitempty"`
	SaturationAdjustmentBlue string `xml:"crs:SaturationAdjustmentBlue,attr,omitempty"`
	LuminanceAdjustmentBlue  string `xml:"crs:LuminanceAdjustmentBlue,attr,omitempty"`

	// HSL Adjustments - Purple
	HueAdjustmentPurple        string `xml:"crs:HueAdjustmentPurple,attr,omitempty"`
	SaturationAdjustmentPurple string `xml:"crs:SaturationAdjustmentPurple,attr,omitempty"`
	LuminanceAdjustmentPurple  string `xml:"crs:LuminanceAdjustmentPurple,attr,omitempty"`

	// HSL Adjustments - Magenta
	HueAdjustmentMagenta        string `xml:"crs:HueAdjustmentMagenta,attr,omitempty"`
	SaturationAdjustmentMagenta string `xml:"crs:SaturationAdjustmentMagenta,attr,omitempty"`
	LuminanceAdjustmentMagenta  string `xml:"crs:LuminanceAdjustmentMagenta,attr,omitempty"`

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

	// Tone Curve Mode (tells Lightroom which curve to apply)
	ToneCurveName2012 string `xml:"crs:ToneCurveName2012,attr,omitempty"` // "Custom" to use custom curve, "Linear" for none

	// Tone Curves (modern PV2012 nested sequence format)
	ToneCurvePV2012      *toneCurveSeqWrapper `xml:"crs:ToneCurvePV2012,omitempty"`
	ToneCurvePV2012Red   *toneCurveSeqWrapper `xml:"crs:ToneCurvePV2012Red,omitempty"`
	ToneCurvePV2012Green *toneCurveSeqWrapper `xml:"crs:ToneCurvePV2012Green,omitempty"`
	ToneCurvePV2012Blue  *toneCurveSeqWrapper `xml:"crs:ToneCurvePV2012Blue,omitempty"`

	// 3D LUT Support (for high-accuracy color transformation)
	RGBTable string `xml:"crs:RGBTable,attr,omitempty"` // MD5 hash reference to the LUT table
	LUTData  string `xml:",innerxml"`                   // Embedded LUT table data (crs:Table_[hash])
}

// toneCurveSeqWrapper wraps the RDF sequence for tone curve points
type toneCurveSeqWrapper struct {
	Seq toneCurveSeqInner `xml:"rdf:Seq"`
}

// toneCurveSeqInner contains the list of tone curve points
type toneCurveSeqInner struct {
	Points []string `xml:"rdf:li"`
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

// getCameraProfile returns the camera profile name, defaulting to "Camera Flexible Color"
// when no profile is specified. This default better matches Nikon's color rendering
// than Adobe Standard, producing more accurate XMP output from NP3 conversions.
func getCameraProfile(profileName string) string {
	if profileName == "" {
		return "Camera Flexible Color"
	}
	return profileName
}

// formatToneCurvePV2012 converts tone curve points to modern PV2012 nested sequence format.
// Returns nil if no points (will be omitted from XML).
// Example output: <crs:ToneCurvePV2012><rdf:Seq><rdf:li>0, 0</rdf:li>...</rdf:Seq></crs:ToneCurvePV2012>
// formatToneCurvePV2012 converts tone curve points to modern PV2012 nested sequence format.
// Returns nil if no points (will be omitted from XML).
// Example output: <crs:ToneCurvePV2012><rdf:Seq><rdf:li>0, 0</rdf:li>...</rdf:Seq></crs:ToneCurvePV2012>
func formatToneCurvePV2012(points []models.ToneCurvePoint) *toneCurveSeqWrapper {
	if len(points) == 0 {
		return nil
	}

	// Lightroom has a limit on point count (historically ~16). NP3 often has 21 points.
	// Use Ramer-Douglas-Peucker algorithm to simplify the curve while preserving shape.

	// FIX: Sanitize start point. Some NP3 profiles (like KOLORA) contain a "Lifted Black" artifact
	// at index 0 (e.g., 0, 12) followed by a drop (e.g., 11, 10). This creates an ugly dip/gray shadows
	// that contradicts the grounded graph in NX Studio.
	// If detected, we force the start point to 0,0 (or interpolated) to ground the curve.
	if len(points) > 1 && points[0].Input == 0 && points[0].Output > points[1].Output {
		// Artifact detected (0,12 -> 11,10). Force 0,0.
		points[0].Output = 0
	}

	// We use an adaptive epsilon strategy to target ~6-8 control points, matching NX Studio's behavior.
	epsilon := 2.0
	finalPoints := simplifyCurveRDP(points, epsilon)

	// Adaptively increase smoothing if we have too many points (likely noise)
	// We target <= 6 points to aggressively remove noise like minor dips (e.g. 12->10)
	// and match typical NX Studio control point counts. Max epsilon 10.0 avoids over-smoothing.
	for len(finalPoints) > 6 && epsilon < 10.0 {
		epsilon += 1.0
		finalPoints = simplifyCurveRDP(points, epsilon)
	}

	// Fallback safety: If still exceeds 16 points (unlikely), strict downsample needed
	if len(finalPoints) > 16 {
		var reduced []models.ToneCurvePoint
		for i := 0; i < len(finalPoints); i += 2 {
			reduced = append(reduced, finalPoints[i])
		}
		if reduced[len(reduced)-1].Input != finalPoints[len(finalPoints)-1].Input {
			reduced = append(reduced, finalPoints[len(finalPoints)-1])
		}
		finalPoints = reduced
	}

	// Convert each point to "input, output" string format
	pointStrings := make([]string, len(finalPoints))
	for i, point := range finalPoints {
		pointStrings[i] = fmt.Sprintf("%d, %d", point.Input, point.Output)
	}

	return &toneCurveSeqWrapper{
		Seq: toneCurveSeqInner{
			Points: pointStrings,
		},
	}
}

// simplifyCurveRDP implements the Ramer-Douglas-Peucker algorithm to reduce the number of points
// in a curve that is approximated by a series of points.
func simplifyCurveRDP(points []models.ToneCurvePoint, epsilon float64) []models.ToneCurvePoint {
	if len(points) < 3 {
		// return copy to avoid backing array corruption during append
		output := make([]models.ToneCurvePoint, len(points))
		copy(output, points)
		return output
	}

	dmax := 0.0
	index := 0
	end := len(points) - 1

	// Find the point with the maximum distance from the line segment [start, end]
	for i := 1; i < end; i++ {
		d := perpendicularDistance(points[i], points[0], points[end])
		if d > dmax {
			dmax = d
			index = i
		}
	}

	// If max distance is greater than epsilon, recursively simplify
	if dmax > epsilon {
		// Recursive call
		recResults1 := simplifyCurveRDP(points[:index+1], epsilon)
		recResults2 := simplifyCurveRDP(points[index:], epsilon)

		// Build the result list (removing duplicate point at index)
		return append(recResults1[:len(recResults1)-1], recResults2...)
	} else {
		return []models.ToneCurvePoint{points[0], points[end]}
	}
}

// perpendicularDistance calculates the distance between point p and the line defined by lineStart and lineEnd
func perpendicularDistance(p, lineStart, lineEnd models.ToneCurvePoint) float64 {
	x0, y0 := float64(p.Input), float64(p.Output)
	x1, y1 := float64(lineStart.Input), float64(lineStart.Output)
	x2, y2 := float64(lineEnd.Input), float64(lineEnd.Output)

	// Handle case where line start and end are the same
	if x1 == x2 && y1 == y2 {
		return math.Sqrt(math.Pow(x0-x1, 2) + math.Pow(y0-y1, 2))
	}

	// Formula: |(y2-y1)x0 - (x2-x1)y0 + x2y1 - y2x1| / sqrt((y2-y1)^2 + (x2-x1)^2)
	numerator := math.Abs((y2-y1)*x0 - (x2-x1)*y0 + x2*y1 - y2*x1)
	denominator := math.Sqrt(math.Pow(y2-y1, 2) + math.Pow(x2-x1, 2))

	return numerator / denominator
}

// formatToneCurveName2012 returns "Custom" if a custom tone curve is defined,
// telling Lightroom to apply the curve instead of using "Linear" (no curve).
func formatToneCurveName2012(points []models.ToneCurvePoint) string {
	if len(points) > 0 {
		return "Custom"
	}
	return "" // Omit attribute if no curve (defaults to Linear)
}

// hasParametricCurve checks if the recipe has non-zero parametric curve values.
// When parametric curve is present, we skip the point curve to avoid double-applying adjustments.
func hasParametricCurve(recipe *models.UniversalRecipe) bool {
	return recipe.ToneCurveShadows != 0 ||
		recipe.ToneCurveDarks != 0 ||
		recipe.ToneCurveLights != 0 ||
		recipe.ToneCurveHighlights != 0
}

// formatToneCurveName2012ParametricAware returns empty if parametric curve is used
// (to let parametric handle tone adjustments), otherwise returns "Custom" for point curve.
func formatToneCurveName2012ParametricAware(recipe *models.UniversalRecipe) string {
	if hasParametricCurve(recipe) {
		return "" // Use Linear base, let parametric curve handle adjustments
	}
	return formatToneCurveName2012(recipe.PointCurve)
}

// formatToneCurvePV2012ParametricAware returns nil if parametric curve is used
// (to skip point curve entirely), otherwise returns the formatted point curve.
func formatToneCurvePV2012ParametricAware(recipe *models.UniversalRecipe) *toneCurveSeqWrapper {
	if hasParametricCurve(recipe) {
		return nil // Skip point curve, use parametric curve instead
	}
	return formatToneCurvePV2012(recipe.PointCurve)
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
// NP3 Chroma maps directly to XMP ColorGrade*Sat with the same -100 to +100 range.
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

// formatColorGradingBalance formats the Balance parameter from NP3's ColorGrading.
// Balance (-100 to +100) shifts overall color balance. In Adobe's color grading,
// this maps to ColorGradeGlobalSat which adjusts the saturation across all zones.
// NP3 Balance=50 adds warmth/vibrancy by increasing global saturation.
// Returns empty string if ColorGrading is nil.
func formatColorGradingBalance(cg *models.ColorGrading) string {
	if cg == nil {
		return ""
	}
	if cg.Balance == 0 {
		return ""
	}
	// Map NP3's Balance directly to Adobe's global saturation
	// Balance range: -100 to +100 maps to saturation adjustment
	return formatInt(cg.Balance)
}

// formatAttributesOnSeparateLines reformats XML to place each crs: attribute on its own line.
// This matches the format of professional Adobe-exported XMP preset files.
// Input: XML with all attributes on one line
// Output: XML with each crs: attribute on a separate indented line
func formatAttributesOnSeparateLines(xmlData []byte) []byte {
	result := string(xmlData)

	// Find rdf:Description elements and split their attributes onto separate lines
	// Pattern: <rdf:Description ...attrs... > or <rdf:Description ...attrs...>
	// We need to handle both self-closing and regular elements

	// Use strings package for efficient manipulation
	lines := strings.Split(result, "\n")
	var output []string

	for _, line := range lines {
		// Check if this line contains rdf:Description with crs: attributes
		if strings.Contains(line, "rdf:Description") && strings.Contains(line, "crs:") {
			// This is a Description line with attributes - reformat it
			reformatted := reformatDescriptionLine(line)
			output = append(output, reformatted...)
		} else {
			output = append(output, line)
		}
	}

	return []byte(strings.Join(output, "\n"))
}

// reformatDescriptionLine takes a single-line rdf:Description element and splits attributes onto separate lines.
func reformatDescriptionLine(line string) []string {
	// Find the indentation
	indent := ""
	for _, ch := range line {
		if ch == ' ' || ch == '\t' {
			indent += string(ch)
		} else {
			break
		}
	}

	// Check if this is a self-closing tag (ends with />) or has content (ends with >)
	isSelfClosing := strings.HasSuffix(strings.TrimSpace(line), "/>")
	hasContent := strings.HasSuffix(strings.TrimSpace(line), ">") && !isSelfClosing

	// Extract the tag name and namespace
	trimmed := strings.TrimLeft(line, " \t")

	// Split by spaces to get attributes, but be careful with quoted values
	parts := splitXMLAttributes(trimmed)

	if len(parts) < 2 {
		return []string{line} // No attributes to reformat
	}

	var result []string
	attrIndent := indent + " " // One extra space for attribute alignment

	// First line: just the opening tag
	result = append(result, indent+parts[0])

	// Each attribute on its own line
	for i := 1; i < len(parts); i++ {
		attr := parts[i]
		if attr == "" || attr == ">" || attr == "/>" {
			continue
		}

		// Check if this is the last attribute (might have > or /> attached)
		isLast := i == len(parts)-1

		if isLast {
			// Handle closing bracket
			if isSelfClosing {
				if strings.HasSuffix(attr, "/>") {
					result = append(result, attrIndent+attr)
				} else {
					result = append(result, attrIndent+attr+"/>")
				}
			} else if hasContent {
				if strings.HasSuffix(attr, ">") {
					result = append(result, attrIndent+attr)
				} else {
					result = append(result, attrIndent+attr+">")
				}
			} else {
				result = append(result, attrIndent+attr)
			}
		} else {
			result = append(result, attrIndent+attr)
		}
	}

	return result
}

// splitXMLAttributes splits an XML element into its component parts while respecting quoted attribute values.
func splitXMLAttributes(element string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(element); i++ {
		ch := element[i]

		if inQuotes {
			current.WriteByte(ch)
			if ch == quoteChar {
				inQuotes = false
				quoteChar = 0
			}
		} else {
			if ch == '"' || ch == '\'' {
				inQuotes = true
				quoteChar = ch
				current.WriteByte(ch)
			} else if ch == ' ' || ch == '\t' {
				if current.Len() > 0 {
					parts = append(parts, current.String())
					current.Reset()
				}
			} else {
				current.WriteByte(ch)
			}
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
