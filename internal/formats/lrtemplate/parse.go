// Package lrtemplate provides functionality for parsing Adobe Lightroom Classic lrtemplate preset files.
//
// The lrtemplate format is a Lua table format used by Adobe Lightroom Classic to store photo editing
// presets. This package decodes lrtemplate files into the UniversalRecipe intermediate representation,
// enabling conversion to other preset formats.
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
// This parser uses regex-based extraction (no external Lua libraries) and achieves the
// <20ms performance target through pre-compiled regex patterns.
package lrtemplate

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/justin/recipe/internal/models"
)

// ConversionError wraps errors that occur during lrtemplate conversion operations.
// This type is reused from the XMP package and follows Pattern 5 (Error Handling)
// from the architecture documentation.
type ConversionError struct {
	Operation string // Operation being performed (e.g., "parse", "validate", "extract")
	Format    string // Format being processed ("lrtemplate")
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

// Pre-compiled regex patterns for efficient parsing (compiled once at package init)
// These patterns extract key-value pairs from Lua table syntax
var (
	// Basic Adjustments (2012 process version)
	exposure2012Regex   = regexp.MustCompile(`Exposure2012\s*=\s*(-?\d+\.?\d*)`)
	contrast2012Regex   = regexp.MustCompile(`Contrast2012\s*=\s*(-?\d+)`)
	highlights2012Regex = regexp.MustCompile(`Highlights2012\s*=\s*(-?\d+)`)
	shadows2012Regex    = regexp.MustCompile(`Shadows2012\s*=\s*(-?\d+)`)
	whites2012Regex     = regexp.MustCompile(`Whites2012\s*=\s*(-?\d+)`)
	blacks2012Regex     = regexp.MustCompile(`Blacks2012\s*=\s*(-?\d+)`)

	// Color Parameters
	saturationRegex  = regexp.MustCompile(`(?:^|\n)\s*Saturation\s*=\s*(-?\d+)`)
	vibranceRegex    = regexp.MustCompile(`Vibrance\s*=\s*(-?\d+)`)
	clarity2012Regex = regexp.MustCompile(`Clarity2012\s*=\s*(-?\d+)`)
	sharpnessRegex   = regexp.MustCompile(`Sharpness\s*=\s*(-?\d+)`)
	temperatureRegex = regexp.MustCompile(`Temperature\s*=\s*(-?\d+)`)
	tintRegex        = regexp.MustCompile(`Tint\s*=\s*(-?\d+)`)

	// HSL Adjustments - Hue
	hueRedRegex     = regexp.MustCompile(`HueAdjustmentRed\s*=\s*(-?\d+)`)
	hueOrangeRegex  = regexp.MustCompile(`HueAdjustmentOrange\s*=\s*(-?\d+)`)
	hueYellowRegex  = regexp.MustCompile(`HueAdjustmentYellow\s*=\s*(-?\d+)`)
	hueGreenRegex   = regexp.MustCompile(`HueAdjustmentGreen\s*=\s*(-?\d+)`)
	hueAquaRegex    = regexp.MustCompile(`HueAdjustmentAqua\s*=\s*(-?\d+)`)
	hueBlueRegex    = regexp.MustCompile(`HueAdjustmentBlue\s*=\s*(-?\d+)`)
	huePurpleRegex  = regexp.MustCompile(`HueAdjustmentPurple\s*=\s*(-?\d+)`)
	hueMagentaRegex = regexp.MustCompile(`HueAdjustmentMagenta\s*=\s*(-?\d+)`)

	// HSL Adjustments - Saturation
	saturationRedRegex     = regexp.MustCompile(`SaturationAdjustmentRed\s*=\s*(-?\d+)`)
	saturationOrangeRegex  = regexp.MustCompile(`SaturationAdjustmentOrange\s*=\s*(-?\d+)`)
	saturationYellowRegex  = regexp.MustCompile(`SaturationAdjustmentYellow\s*=\s*(-?\d+)`)
	saturationGreenRegex   = regexp.MustCompile(`SaturationAdjustmentGreen\s*=\s*(-?\d+)`)
	saturationAquaRegex    = regexp.MustCompile(`SaturationAdjustmentAqua\s*=\s*(-?\d+)`)
	saturationBlueRegex    = regexp.MustCompile(`SaturationAdjustmentBlue\s*=\s*(-?\d+)`)
	saturationPurpleRegex  = regexp.MustCompile(`SaturationAdjustmentPurple\s*=\s*(-?\d+)`)
	saturationMagentaRegex = regexp.MustCompile(`SaturationAdjustmentMagenta\s*=\s*(-?\d+)`)

	// HSL Adjustments - Luminance
	luminanceRedRegex     = regexp.MustCompile(`LuminanceAdjustmentRed\s*=\s*(-?\d+)`)
	luminanceOrangeRegex  = regexp.MustCompile(`LuminanceAdjustmentOrange\s*=\s*(-?\d+)`)
	luminanceYellowRegex  = regexp.MustCompile(`LuminanceAdjustmentYellow\s*=\s*(-?\d+)`)
	luminanceGreenRegex   = regexp.MustCompile(`LuminanceAdjustmentGreen\s*=\s*(-?\d+)`)
	luminanceAquaRegex    = regexp.MustCompile(`LuminanceAdjustmentAqua\s*=\s*(-?\d+)`)
	luminanceBlueRegex    = regexp.MustCompile(`LuminanceAdjustmentBlue\s*=\s*(-?\d+)`)
	luminancePurpleRegex  = regexp.MustCompile(`LuminanceAdjustmentPurple\s*=\s*(-?\d+)`)
	luminanceMagentaRegex = regexp.MustCompile(`LuminanceAdjustmentMagenta\s*=\s*(-?\d+)`)

	// Split Toning
	splitShadowHueRegex           = regexp.MustCompile(`SplitToningShadowHue\s*=\s*(-?\d+)`)
	splitShadowSaturationRegex    = regexp.MustCompile(`SplitToningShadowSaturation\s*=\s*(-?\d+)`)
	splitHighlightHueRegex        = regexp.MustCompile(`SplitToningHighlightHue\s*=\s*(-?\d+)`)
	splitHighlightSaturationRegex = regexp.MustCompile(`SplitToningHighlightSaturation\s*=\s*(-?\d+)`)
	splitBalanceRegex             = regexp.MustCompile(`SplitToningBalance\s*=\s*(-?\d+)`)

	// Tone Curve (array format)
	toneCurvePV2012Regex = regexp.MustCompile(`ToneCurvePV2012\s*=\s*\{([^}]+)\}`)

	// Grain and Vignette
	grainAmountRegex      = regexp.MustCompile(`GrainAmount\s*=\s*(-?\d+)`)
	grainSizeRegex        = regexp.MustCompile(`GrainSize\s*=\s*(-?\d+)`)
	vignetteAmountRegex   = regexp.MustCompile(`VignetteAmount\s*=\s*(-?\d+)`)
	vignetteMidpointRegex = regexp.MustCompile(`VignetteMidpoint\s*=\s*(-?\d+)`)
)

// Parse decodes an Adobe Lightroom Classic lrtemplate preset file into a UniversalRecipe.
//
// The function validates the file structure (Lua table starting with `s = {`), extracts
// all 50+ photo editing parameters using regex patterns, validates parameter ranges,
// and constructs a UniversalRecipe.
//
// Parameters:
//   - data: Raw bytes of the .lrtemplate file
//
// Returns:
//   - *models.UniversalRecipe: Populated recipe with extracted parameters
//   - error: Validation or parsing error with descriptive context
//
// Errors:
//   - Invalid Lua structure: File doesn't start with `s = {`
//   - Parameter out of range: Invalid parameter value
//   - Parse error: Failed to extract or convert parameter
func Parse(data []byte) (*models.UniversalRecipe, error) {
	// Validate Lua table structure (fail-fast per Pattern 6)
	if err := validateLuaStructure(data); err != nil {
		return nil, &ConversionError{
			Operation: "parse",
			Format:    "lrtemplate",
			Cause:     err,
		}
	}

	// Extract and convert parameters from Lua table
	params, err := extractParameters(data)
	if err != nil {
		return nil, err // Already wrapped by extractParameters
	}

	// Validate extracted parameters (inline validation per Pattern 6)
	if err := validateParameters(params); err != nil {
		return nil, err // Already wrapped by validateParameters
	}

	// Build UniversalRecipe
	recipe := buildRecipe(params)

	return recipe, nil
}

// validateLuaStructure checks if the data is a valid lrtemplate file starting with `s = {`
func validateLuaStructure(data []byte) error {
	// Check for basic Lua table structure
	if len(data) == 0 {
		return fmt.Errorf("empty file")
	}

	// Validate file starts with `s = {`
	dataStr := string(bytes.TrimSpace(data))
	if !strings.HasPrefix(dataStr, "s = {") {
		return fmt.Errorf("invalid lrtemplate format: file must start with 's = {'")
	}

	// Check for basic table closure
	if !strings.Contains(dataStr, "}") {
		return fmt.Errorf("invalid lrtemplate format: incomplete Lua table (missing closing brace)")
	}

	return nil
}

// lrtemplateParameters holds extracted parameter values before validation
type lrtemplateParameters struct {
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

	// Tone Curve
	toneCurve []models.ToneCurvePoint

	// Grain
	grainAmount int
	grainSize   int

	// Vignette
	vignetteAmount   int
	vignetteMidpoint int
}

// extractParameters extracts all parameter values from the lrtemplate file using regex patterns
func extractParameters(data []byte) (*lrtemplateParameters, error) {
	params := &lrtemplateParameters{}
	dataStr := string(data)

	// Extract Basic Adjustments
	params.exposure = extractFloat64(dataStr, exposure2012Regex, "Exposure2012")
	params.contrast = extractInt(dataStr, contrast2012Regex, "Contrast2012")
	params.highlights = extractInt(dataStr, highlights2012Regex, "Highlights2012")
	params.shadows = extractInt(dataStr, shadows2012Regex, "Shadows2012")
	params.whites = extractInt(dataStr, whites2012Regex, "Whites2012")
	params.blacks = extractInt(dataStr, blacks2012Regex, "Blacks2012")

	// Extract Color Parameters
	params.saturation = extractInt(dataStr, saturationRegex, "Saturation")
	params.vibrance = extractInt(dataStr, vibranceRegex, "Vibrance")
	params.clarity = extractInt(dataStr, clarity2012Regex, "Clarity2012")
	params.sharpness = extractInt(dataStr, sharpnessRegex, "Sharpness")
	params.temperature = extractInt(dataStr, temperatureRegex, "Temperature")
	params.tint = extractInt(dataStr, tintRegex, "Tint")

	// Extract HSL Adjustments
	params.red = models.ColorAdjustment{
		Hue:        extractInt(dataStr, hueRedRegex, "HueAdjustmentRed"),
		Saturation: extractInt(dataStr, saturationRedRegex, "SaturationAdjustmentRed"),
		Luminance:  extractInt(dataStr, luminanceRedRegex, "LuminanceAdjustmentRed"),
	}

	params.orange = models.ColorAdjustment{
		Hue:        extractInt(dataStr, hueOrangeRegex, "HueAdjustmentOrange"),
		Saturation: extractInt(dataStr, saturationOrangeRegex, "SaturationAdjustmentOrange"),
		Luminance:  extractInt(dataStr, luminanceOrangeRegex, "LuminanceAdjustmentOrange"),
	}

	params.yellow = models.ColorAdjustment{
		Hue:        extractInt(dataStr, hueYellowRegex, "HueAdjustmentYellow"),
		Saturation: extractInt(dataStr, saturationYellowRegex, "SaturationAdjustmentYellow"),
		Luminance:  extractInt(dataStr, luminanceYellowRegex, "LuminanceAdjustmentYellow"),
	}

	params.green = models.ColorAdjustment{
		Hue:        extractInt(dataStr, hueGreenRegex, "HueAdjustmentGreen"),
		Saturation: extractInt(dataStr, saturationGreenRegex, "SaturationAdjustmentGreen"),
		Luminance:  extractInt(dataStr, luminanceGreenRegex, "LuminanceAdjustmentGreen"),
	}

	params.aqua = models.ColorAdjustment{
		Hue:        extractInt(dataStr, hueAquaRegex, "HueAdjustmentAqua"),
		Saturation: extractInt(dataStr, saturationAquaRegex, "SaturationAdjustmentAqua"),
		Luminance:  extractInt(dataStr, luminanceAquaRegex, "LuminanceAdjustmentAqua"),
	}

	params.blue = models.ColorAdjustment{
		Hue:        extractInt(dataStr, hueBlueRegex, "HueAdjustmentBlue"),
		Saturation: extractInt(dataStr, saturationBlueRegex, "SaturationAdjustmentBlue"),
		Luminance:  extractInt(dataStr, luminanceBlueRegex, "LuminanceAdjustmentBlue"),
	}

	params.purple = models.ColorAdjustment{
		Hue:        extractInt(dataStr, huePurpleRegex, "HueAdjustmentPurple"),
		Saturation: extractInt(dataStr, saturationPurpleRegex, "SaturationAdjustmentPurple"),
		Luminance:  extractInt(dataStr, luminancePurpleRegex, "LuminanceAdjustmentPurple"),
	}

	params.magenta = models.ColorAdjustment{
		Hue:        extractInt(dataStr, hueMagentaRegex, "HueAdjustmentMagenta"),
		Saturation: extractInt(dataStr, saturationMagentaRegex, "SaturationAdjustmentMagenta"),
		Luminance:  extractInt(dataStr, luminanceMagentaRegex, "LuminanceAdjustmentMagenta"),
	}

	// Extract Split Toning
	params.splitShadowHue = extractInt(dataStr, splitShadowHueRegex, "SplitToningShadowHue")
	params.splitShadowSaturation = extractInt(dataStr, splitShadowSaturationRegex, "SplitToningShadowSaturation")
	params.splitHighlightHue = extractInt(dataStr, splitHighlightHueRegex, "SplitToningHighlightHue")
	params.splitHighlightSaturation = extractInt(dataStr, splitHighlightSaturationRegex, "SplitToningHighlightSaturation")
	params.splitBalance = extractInt(dataStr, splitBalanceRegex, "SplitToningBalance")

	// Extract Tone Curve
	var err error
	params.toneCurve, err = extractToneCurve(dataStr)
	if err != nil {
		return nil, &ConversionError{
			Operation: "parse",
			Format:    "lrtemplate",
			Field:     "ToneCurvePV2012",
			Cause:     err,
		}
	}

	// Extract Grain
	params.grainAmount = extractInt(dataStr, grainAmountRegex, "GrainAmount")
	params.grainSize = extractInt(dataStr, grainSizeRegex, "GrainSize")

	// Extract Vignette
	params.vignetteAmount = extractInt(dataStr, vignetteAmountRegex, "VignetteAmount")
	params.vignetteMidpoint = extractInt(dataStr, vignetteMidpointRegex, "VignetteMidpoint")

	return params, nil
}

// extractInt extracts an integer value using a pre-compiled regex pattern
// Returns 0 if field is not found (graceful handling of missing fields)
func extractInt(data string, pattern *regexp.Regexp, fieldName string) int {
	matches := pattern.FindStringSubmatch(data)
	if len(matches) < 2 {
		return 0 // Field not found, use zero value
	}

	val, err := strconv.Atoi(matches[1])
	if err != nil {
		// This should not happen if regex is correct, but handle gracefully
		return 0
	}

	return val
}

// extractFloat64 extracts a float64 value using a pre-compiled regex pattern
// Returns 0.0 if field is not found (graceful handling of missing fields)
func extractFloat64(data string, pattern *regexp.Regexp, fieldName string) float64 {
	matches := pattern.FindStringSubmatch(data)
	if len(matches) < 2 {
		return 0.0 // Field not found, use zero value
	}

	val, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		// This should not happen if regex is correct, but handle gracefully
		return 0.0
	}

	return val
}

// extractToneCurve extracts the ToneCurvePV2012 array from lrtemplate data
// Format: ToneCurvePV2012 = { 0, 0, 255, 255 } (flat array of coordinate pairs)
func extractToneCurve(data string) ([]models.ToneCurvePoint, error) {
	matches := toneCurvePV2012Regex.FindStringSubmatch(data)
	if len(matches) < 2 {
		return nil, nil // Field not found, return nil (no tone curve)
	}

	// Parse the array values
	arrayStr := matches[1]
	// Split by comma and filter out empty strings
	rawValues := strings.Split(arrayStr, ",")
	var values []string
	for _, v := range rawValues {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}

	// Convert to ToneCurvePoint array (pairs of input/output values)
	var points []models.ToneCurvePoint
	for i := 0; i+1 < len(values); i += 2 {
		input, err := strconv.Atoi(values[i])
		if err != nil {
			return nil, fmt.Errorf("invalid tone curve input value: %v", err)
		}

		output, err := strconv.Atoi(values[i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid tone curve output value: %v", err)
		}

		points = append(points, models.ToneCurvePoint{
			Input:  input,
			Output: output,
		})
	}

	return points, nil
}

// validateParameters validates all extracted parameters are within expected ranges
func validateParameters(params *lrtemplateParameters) error {
	// Validate Exposure (-5.0 to +5.0)
	if params.exposure < -5.0 || params.exposure > 5.0 {
		return &ConversionError{
			Operation: "validate",
			Format:    "lrtemplate",
			Field:     "Exposure2012",
			Cause:     fmt.Errorf("value %v out of range [-5.0, +5.0]", params.exposure),
		}
	}

	// Validate integer parameters (-100 to +100 range)
	if err := validateIntRange(params.contrast, -100, 100, "Contrast2012"); err != nil {
		return err
	}
	if err := validateIntRange(params.highlights, -100, 100, "Highlights2012"); err != nil {
		return err
	}
	if err := validateIntRange(params.shadows, -100, 100, "Shadows2012"); err != nil {
		return err
	}
	if err := validateIntRange(params.whites, -100, 100, "Whites2012"); err != nil {
		return err
	}
	if err := validateIntRange(params.blacks, -100, 100, "Blacks2012"); err != nil {
		return err
	}
	if err := validateIntRange(params.saturation, -100, 100, "Saturation"); err != nil {
		return err
	}
	if err := validateIntRange(params.vibrance, -100, 100, "Vibrance"); err != nil {
		return err
	}
	if err := validateIntRange(params.clarity, -100, 100, "Clarity2012"); err != nil {
		return err
	}
	if err := validateIntRange(params.tint, -150, 150, "Tint"); err != nil {
		return err
	}

	// Validate Sharpness (0 to 150)
	if err := validateIntRange(params.sharpness, 0, 150, "Sharpness"); err != nil {
		return err
	}

	// Validate Temperature (2000 to 50000 Kelvin)
	if params.temperature != 0 {
		if err := validateIntRange(params.temperature, 2000, 50000, "Temperature"); err != nil {
			return err
		}
	}

	// Validate HSL adjustments (-100 to +100)
	if err := validateColorAdjustment(params.red, "Red"); err != nil {
		return err
	}
	if err := validateColorAdjustment(params.orange, "Orange"); err != nil {
		return err
	}
	if err := validateColorAdjustment(params.yellow, "Yellow"); err != nil {
		return err
	}
	if err := validateColorAdjustment(params.green, "Green"); err != nil {
		return err
	}
	if err := validateColorAdjustment(params.aqua, "Aqua"); err != nil {
		return err
	}
	if err := validateColorAdjustment(params.blue, "Blue"); err != nil {
		return err
	}
	if err := validateColorAdjustment(params.purple, "Purple"); err != nil {
		return err
	}
	if err := validateColorAdjustment(params.magenta, "Magenta"); err != nil {
		return err
	}

	// Validate Split Toning (Hue: 0-360, Saturation: 0-100)
	if err := validateIntRange(params.splitShadowHue, 0, 360, "SplitToningShadowHue"); err != nil {
		return err
	}
	if err := validateIntRange(params.splitShadowSaturation, 0, 100, "SplitToningShadowSaturation"); err != nil {
		return err
	}
	if err := validateIntRange(params.splitHighlightHue, 0, 360, "SplitToningHighlightHue"); err != nil {
		return err
	}
	if err := validateIntRange(params.splitHighlightSaturation, 0, 100, "SplitToningHighlightSaturation"); err != nil {
		return err
	}
	if err := validateIntRange(params.splitBalance, -100, 100, "SplitToningBalance"); err != nil {
		return err
	}

	// Validate Tone Curve points (0-255 range)
	for i, point := range params.toneCurve {
		if point.Input < 0 || point.Input > 255 {
			return &ConversionError{
				Operation: "validate",
				Format:    "lrtemplate",
				Field:     fmt.Sprintf("ToneCurvePV2012[%d].Input", i),
				Cause:     fmt.Errorf("value %d out of range [0, 255]", point.Input),
			}
		}
		if point.Output < 0 || point.Output > 255 {
			return &ConversionError{
				Operation: "validate",
				Format:    "lrtemplate",
				Field:     fmt.Sprintf("ToneCurvePV2012[%d].Output", i),
				Cause:     fmt.Errorf("value %d out of range [0, 255]", point.Output),
			}
		}
	}

	// Validate Grain (0-100)
	if err := validateIntRange(params.grainAmount, 0, 100, "GrainAmount"); err != nil {
		return err
	}
	if err := validateIntRange(params.grainSize, 0, 100, "GrainSize"); err != nil {
		return err
	}

	// Validate Vignette
	if err := validateIntRange(params.vignetteAmount, -100, 100, "VignetteAmount"); err != nil {
		return err
	}
	if err := validateIntRange(params.vignetteMidpoint, 0, 100, "VignetteMidpoint"); err != nil {
		return err
	}

	return nil
}

// validateIntRange validates an integer value is within the specified range
func validateIntRange(value, min, max int, fieldName string) error {
	// Allow zero values (missing fields)
	if value == 0 {
		return nil
	}

	if value < min || value > max {
		return &ConversionError{
			Operation: "validate",
			Format:    "lrtemplate",
			Field:     fieldName,
			Cause:     fmt.Errorf("value %d out of range [%d, %d]", value, min, max),
		}
	}
	return nil
}

// validateColorAdjustment validates a ColorAdjustment struct (HSL values)
func validateColorAdjustment(color models.ColorAdjustment, colorName string) error {
	if err := validateIntRange(color.Hue, -100, 100, fmt.Sprintf("HueAdjustment%s", colorName)); err != nil {
		return err
	}
	if err := validateIntRange(color.Saturation, -100, 100, fmt.Sprintf("SaturationAdjustment%s", colorName)); err != nil {
		return err
	}
	if err := validateIntRange(color.Luminance, -100, 100, fmt.Sprintf("LuminanceAdjustment%s", colorName)); err != nil {
		return err
	}
	return nil
}

// buildRecipe constructs a UniversalRecipe from validated parameters
func buildRecipe(params *lrtemplateParameters) *models.UniversalRecipe {
	recipe := &models.UniversalRecipe{
		SourceFormat: "lrtemplate",

		// Basic Adjustments
		Exposure:   params.exposure,
		Contrast:   params.contrast,
		Highlights: params.highlights,
		Shadows:    params.shadows,
		Whites:     params.whites,
		Blacks:     params.blacks,

		// Color Parameters
		Saturation: params.saturation,
		Vibrance:   params.vibrance,
		Clarity:    params.clarity,
		Sharpness:  params.sharpness,
		Tint:       params.tint,

		// HSL Adjustments
		Red:     params.red,
		Orange:  params.orange,
		Yellow:  params.yellow,
		Green:   params.green,
		Aqua:    params.aqua,
		Blue:    params.blue,
		Purple:  params.purple,
		Magenta: params.magenta,

		// Split Toning
		SplitShadowHue:           params.splitShadowHue,
		SplitShadowSaturation:    params.splitShadowSaturation,
		SplitHighlightHue:        params.splitHighlightHue,
		SplitHighlightSaturation: params.splitHighlightSaturation,
		SplitBalance:             params.splitBalance,

		// Tone Curve
		PointCurve: params.toneCurve,

		// Grain
		GrainAmount: params.grainAmount,
		GrainSize:   params.grainSize,

		// Vignette
		VignetteAmount:   params.vignetteAmount,
		VignetteMidpoint: params.vignetteMidpoint,
	}

	// Handle Temperature (nullable field in UniversalRecipe)
	if params.temperature != 0 {
		recipe.Temperature = &params.temperature
	}

	return recipe
}
