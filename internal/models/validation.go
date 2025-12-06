package models

import (
	"fmt"
)

// ValidateExposure validates that the exposure value is within the acceptable range (-5.0 to +5.0).
func ValidateExposure(value float64) error {
	if value < -5.0 || value > 5.0 {
		return fmt.Errorf("exposure value %v is out of range (must be between -5.0 and +5.0)", value)
	}
	return nil
}

// ValidatePercentage validates that an integer value is within the -100 to +100 range.
// This is used for most photo editing parameters like contrast, highlights, shadows, etc.
func ValidatePercentage(value int, fieldName string) error {
	if value < -100 || value > 100 {
		return fmt.Errorf("%s value %d is out of range (must be between -100 and +100)", fieldName, value)
	}
	return nil
}

// ValidatePositivePercentage validates that an integer value is within the 0 to +100 range.
// This is used for parameters that only accept positive values.
func ValidatePositivePercentage(value int, fieldName string) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("%s value %d is out of range (must be between 0 and 100)", fieldName, value)
	}
	return nil
}

// ValidateHSL validates HSL color adjustment parameters.
// Hue: -180 to +180
// Saturation: -100 to +100
// Luminance: -100 to +100
func ValidateHSL(hue, saturation, luminance int, colorName string) error {
	if hue < -180 || hue > 180 {
		return fmt.Errorf("%s hue value %d is out of range (must be between -180 and +180)", colorName, hue)
	}
	if saturation < -100 || saturation > 100 {
		return fmt.Errorf("%s saturation value %d is out of range (must be between -100 and +100)", colorName, saturation)
	}
	if luminance < -100 || luminance > 100 {
		return fmt.Errorf("%s luminance value %d is out of range (must be between -100 and +100)", colorName, luminance)
	}
	return nil
}

// ValidateNP3Sharpening validates NP3-specific sharpening parameter (0-9).
func ValidateNP3Sharpening(value int) error {
	if value < 0 || value > 9 {
		return fmt.Errorf("NP3 sharpening value %d is out of range (must be between 0 and 9)", value)
	}
	return nil
}

// ValidateNP3Contrast validates NP3-specific contrast parameter (-3 to +3).
func ValidateNP3Contrast(value int) error {
	if value < -3 || value > 3 {
		return fmt.Errorf("NP3 contrast value %d is out of range (must be between -3 and +3)", value)
	}
	return nil
}

// ValidateNP3Brightness validates NP3-specific brightness parameter (-1 to +1).
func ValidateNP3Brightness(value float64) error {
	if value < -1.0 || value > 1.0 {
		return fmt.Errorf("NP3 brightness value %v is out of range (must be between -1.0 and +1.0)", value)
	}
	return nil
}

// ValidateNP3Saturation validates NP3-specific saturation parameter (-3 to +3).
func ValidateNP3Saturation(value int) error {
	if value < -3 || value > 3 {
		return fmt.Errorf("NP3 saturation value %d is out of range (must be between -3 and +3)", value)
	}
	return nil
}

// ValidateNP3Hue validates NP3-specific hue parameter (-9 to +9).
func ValidateNP3Hue(value int) error {
	if value < -9 || value > 9 {
		return fmt.Errorf("NP3 hue value %d is out of range (must be between -9 and +9)", value)
	}
	return nil
}

// ValidateToneCurvePoint validates a tone curve point's input and output values (0-255).
func ValidateToneCurvePoint(input, output int) error {
	if input < 0 || input > 255 {
		return fmt.Errorf("tone curve input value %d is out of range (must be between 0 and 255)", input)
	}
	if output < 0 || output > 255 {
		return fmt.Errorf("tone curve output value %d is out of range (must be between 0 and 255)", output)
	}
	return nil
}

// ValidateHue360 validates hue values in 0-360 degree range (used for split toning).
func ValidateHue360(value int, fieldName string) error {
	if value < 0 || value > 360 {
		return fmt.Errorf("%s value %d is out of range (must be between 0 and 360)", fieldName, value)
	}
	return nil
}

// ValidateSharpnessRadius validates sharpness radius (0.5-3.0).
func ValidateSharpnessRadius(value float64) error {
	if value < 0.5 || value > 3.0 {
		return fmt.Errorf("sharpness radius value %v is out of range (must be between 0.5 and 3.0)", value)
	}
	return nil
}

// ValidateTint validates tint parameter (-150 to +150).
func ValidateTint(value int) error {
	if value < -150 || value > 150 {
		return fmt.Errorf("tint value %d is out of range (must be between -150 and +150)", value)
	}
	return nil
}

// ValidateTemperature validates temperature value (nullable, in Kelvin).
// XMP/Lightroom uses absolute Kelvin values in range [2000, 50000].
func ValidateTemperature(value int) error {
	if value < 2000 || value > 50000 {
		return fmt.Errorf("temperature value %d is out of range (must be between 2000 and 50000 Kelvin)", value)
	}
	return nil
}

// Validate validates all parameters in a UniversalRecipe.
// Returns the first validation error encountered, or nil if all parameters are valid.
func (r *UniversalRecipe) Validate() error {
	// Validate exposure
	if err := ValidateExposure(r.Exposure); err != nil {
		return err
	}

	// Validate basic adjustments
	if err := ValidatePercentage(r.Contrast, "contrast"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.Highlights, "highlights"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.Shadows, "shadows"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.Whites, "whites"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.Blacks, "blacks"); err != nil {
		return err
	}

	// Validate presence
	if err := ValidatePercentage(r.Texture, "texture"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.Clarity, "clarity"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.Dehaze, "dehaze"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.Vibrance, "vibrance"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.Saturation, "saturation"); err != nil {
		return err
	}

	// Validate sharpening
	if err := ValidateSharpnessRadius(r.SharpnessRadius); err != nil && r.SharpnessRadius != 0 {
		return err
	}
	if err := ValidatePositivePercentage(r.SharpnessDetail, "sharpness detail"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.SharpnessMasking, "sharpness masking"); err != nil {
		return err
	}

	// Validate tint
	if err := ValidateTint(r.Tint); err != nil {
		return err
	}

	// Validate HSL adjustments
	if err := ValidateHSL(r.Red.Hue, r.Red.Saturation, r.Red.Luminance, "red"); err != nil {
		return err
	}
	if err := ValidateHSL(r.Orange.Hue, r.Orange.Saturation, r.Orange.Luminance, "orange"); err != nil {
		return err
	}
	if err := ValidateHSL(r.Yellow.Hue, r.Yellow.Saturation, r.Yellow.Luminance, "yellow"); err != nil {
		return err
	}
	if err := ValidateHSL(r.Green.Hue, r.Green.Saturation, r.Green.Luminance, "green"); err != nil {
		return err
	}
	if err := ValidateHSL(r.Aqua.Hue, r.Aqua.Saturation, r.Aqua.Luminance, "aqua"); err != nil {
		return err
	}
	if err := ValidateHSL(r.Blue.Hue, r.Blue.Saturation, r.Blue.Luminance, "blue"); err != nil {
		return err
	}
	if err := ValidateHSL(r.Purple.Hue, r.Purple.Saturation, r.Purple.Luminance, "purple"); err != nil {
		return err
	}
	if err := ValidateHSL(r.Magenta.Hue, r.Magenta.Saturation, r.Magenta.Luminance, "magenta"); err != nil {
		return err
	}

	// Validate tone curve parameters
	if err := ValidatePercentage(r.ToneCurveShadows, "tone curve shadows"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.ToneCurveDarks, "tone curve darks"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.ToneCurveLights, "tone curve lights"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.ToneCurveHighlights, "tone curve highlights"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.ToneCurveShadowSplit, "tone curve shadow split"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.ToneCurveMidtoneSplit, "tone curve midtone split"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.ToneCurveHighlightSplit, "tone curve highlight split"); err != nil {
		return err
	}

	// Validate point curve points
	for _, point := range r.PointCurve {
		if err := ValidateToneCurvePoint(point.Input, point.Output); err != nil {
			return err
		}
	}
	for _, point := range r.PointCurveRed {
		if err := ValidateToneCurvePoint(point.Input, point.Output); err != nil {
			return err
		}
	}
	for _, point := range r.PointCurveGreen {
		if err := ValidateToneCurvePoint(point.Input, point.Output); err != nil {
			return err
		}
	}
	for _, point := range r.PointCurveBlue {
		if err := ValidateToneCurvePoint(point.Input, point.Output); err != nil {
			return err
		}
	}

	// Validate split toning
	if err := ValidateHue360(r.SplitShadowHue, "split shadow hue"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.SplitShadowSaturation, "split shadow saturation"); err != nil {
		return err
	}
	if err := ValidateHue360(r.SplitHighlightHue, "split highlight hue"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.SplitHighlightSaturation, "split highlight saturation"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.SplitBalance, "split balance"); err != nil {
		return err
	}

	// Validate grain effects
	if err := ValidatePositivePercentage(r.GrainAmount, "grain amount"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.GrainSize, "grain size"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.GrainRoughness, "grain roughness"); err != nil {
		return err
	}

	// Validate vignette
	if err := ValidatePercentage(r.VignetteAmount, "vignette amount"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.VignetteMidpoint, "vignette midpoint"); err != nil {
		return err
	}
	if err := ValidatePercentage(r.VignetteRoundness, "vignette roundness"); err != nil {
		return err
	}
	if err := ValidatePositivePercentage(r.VignetteFeather, "vignette feather"); err != nil {
		return err
	}

	return nil
}

// ValidateMidRangeSharpening validates mid-range sharpening parameter (-5.0 to +5.0)
func ValidateMidRangeSharpening(value float64) error {
	if value < -5.0 || value > 5.0 {
		return fmt.Errorf("mid-range sharpening must be -5.0 to 5.0, got %.1f", value)
	}
	return nil
}

// ValidateColorGradingHue validates color grading hue (0-360 degrees)
func ValidateColorGradingHue(hue int) error {
	if hue < 0 || hue > 360 {
		return fmt.Errorf("color grading hue must be 0-360, got %d", hue)
	}
	return nil
}

// ValidateColorGradingChroma validates color grading chroma (-100 to +100)
func ValidateColorGradingChroma(chroma int) error {
	if chroma < -100 || chroma > 100 {
		return fmt.Errorf("color grading chroma must be -100 to 100, got %d", chroma)
	}
	return nil
}

// ValidateColorGradingBrightness validates color grading brightness (-100 to +100)
func ValidateColorGradingBrightness(brightness int) error {
	if brightness < -100 || brightness > 100 {
		return fmt.Errorf("color grading brightness must be -100 to 100, got %d", brightness)
	}
	return nil
}

// ValidateColorGradingBlending validates color grading blending (0-100)
func ValidateColorGradingBlending(blending int) error {
	if blending < 0 || blending > 100 {
		return fmt.Errorf("color grading blending must be 0-100, got %d", blending)
	}
	return nil
}

// ValidateColorGradingBalance validates color grading balance (-100 to +100)
func ValidateColorGradingBalance(balance int) error {
	if balance < -100 || balance > 100 {
		return fmt.Errorf("color grading balance must be -100 to 100, got %d", balance)
	}
	return nil
}
