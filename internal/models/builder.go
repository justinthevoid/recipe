package models

import (
	"fmt"
)

// RecipeBuilder provides a fluent API for constructing UniversalRecipe instances
// with validation and default value initialization.
type RecipeBuilder struct {
	recipe UniversalRecipe
	errors []error
}

// NewRecipeBuilder creates a new RecipeBuilder with sensible defaults.
func NewRecipeBuilder() *RecipeBuilder {
	return &RecipeBuilder{
		recipe: UniversalRecipe{
			SharpnessRadius:         1.0, // Default sharpness radius
			ToneCurveShadowSplit:    25,  // Default shadow split
			ToneCurveMidtoneSplit:   50,  // Default midtone split
			ToneCurveHighlightSplit: 75,  // Default highlight split
		},
		errors: []error{},
	}
}

// WithName sets the recipe name.
func (b *RecipeBuilder) WithName(name string) *RecipeBuilder {
	b.recipe.Name = name
	return b
}

// WithDescription sets the recipe description.
func (b *RecipeBuilder) WithDescription(description string) *RecipeBuilder {
	b.recipe.Description = description
	return b
}

// WithSourceFormat sets the source format.
func (b *RecipeBuilder) WithSourceFormat(format string) *RecipeBuilder {
	b.recipe.SourceFormat = format
	return b
}

// WithExposure sets the exposure value with validation.
func (b *RecipeBuilder) WithExposure(value float64) *RecipeBuilder {
	if err := ValidateExposure(value); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Exposure = value
	}
	return b
}

// WithContrast sets the contrast value with validation.
func (b *RecipeBuilder) WithContrast(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "contrast"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Contrast = value
	}
	return b
}

// WithHighlights sets the highlights value with validation.
func (b *RecipeBuilder) WithHighlights(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "highlights"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Highlights = value
	}
	return b
}

// WithShadows sets the shadows value with validation.
func (b *RecipeBuilder) WithShadows(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "shadows"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Shadows = value
	}
	return b
}

// WithWhites sets the whites value with validation.
func (b *RecipeBuilder) WithWhites(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "whites"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Whites = value
	}
	return b
}

// WithBlacks sets the blacks value with validation.
func (b *RecipeBuilder) WithBlacks(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "blacks"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Blacks = value
	}
	return b
}

// WithTexture sets the texture value with validation.
func (b *RecipeBuilder) WithTexture(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "texture"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Texture = value
	}
	return b
}

// WithClarity sets the clarity value with validation.
func (b *RecipeBuilder) WithClarity(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "clarity"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Clarity = value
	}
	return b
}

// WithDehaze sets the dehaze value with validation.
func (b *RecipeBuilder) WithDehaze(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "dehaze"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Dehaze = value
	}
	return b
}

// WithVibrance sets the vibrance value with validation.
func (b *RecipeBuilder) WithVibrance(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "vibrance"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Vibrance = value
	}
	return b
}

// WithSaturation sets the saturation value with validation.
func (b *RecipeBuilder) WithSaturation(value int) *RecipeBuilder {
	if err := ValidatePercentage(value, "saturation"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Saturation = value
	}
	return b
}

// WithSharpness sets the sharpness value.
func (b *RecipeBuilder) WithSharpness(value int) *RecipeBuilder {
	b.recipe.Sharpness = value
	return b
}

// WithSharpnessRadius sets the sharpness radius with validation.
func (b *RecipeBuilder) WithSharpnessRadius(value float64) *RecipeBuilder {
	if err := ValidateSharpnessRadius(value); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.SharpnessRadius = value
	}
	return b
}

// WithSharpnessDetail sets the sharpness detail with validation.
func (b *RecipeBuilder) WithSharpnessDetail(value int) *RecipeBuilder {
	if err := ValidatePositivePercentage(value, "sharpness detail"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.SharpnessDetail = value
	}
	return b
}

// WithSharpnessMasking sets the sharpness masking with validation.
func (b *RecipeBuilder) WithSharpnessMasking(value int) *RecipeBuilder {
	if err := ValidatePositivePercentage(value, "sharpness masking"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.SharpnessMasking = value
	}
	return b
}

// WithMidRangeSharpening sets the mid-range sharpening value with validation (NP3-specific).
func (b *RecipeBuilder) WithMidRangeSharpening(value float64) *RecipeBuilder {
	if err := ValidateMidRangeSharpening(value); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.MidRangeSharpening = value
	}
	return b
}

// WithTemperature sets the temperature value (nullable).
func (b *RecipeBuilder) WithTemperature(value int) *RecipeBuilder {
	// Temperature is nullable - only set if value is non-zero (within valid range)
	// A value of 0 means "not set" and should remain nil
	if value != 0 {
		if err := ValidateTemperature(value); err != nil {
			b.errors = append(b.errors, err)
		} else {
			b.recipe.Temperature = &value
		}
	}
	return b
}

// WithTint sets the tint value with validation.
func (b *RecipeBuilder) WithTint(value int) *RecipeBuilder {
	if err := ValidateTint(value); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Tint = value
	}
	return b
}

// WithRedHSL sets the red color HSL adjustments with validation.
func (b *RecipeBuilder) WithRedHSL(hue, saturation, luminance int) *RecipeBuilder {
	if err := ValidateHSL(hue, saturation, luminance, "red"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Red = ColorAdjustment{Hue: hue, Saturation: saturation, Luminance: luminance}
	}
	return b
}

// WithOrangeHSL sets the orange color HSL adjustments with validation.
func (b *RecipeBuilder) WithOrangeHSL(hue, saturation, luminance int) *RecipeBuilder {
	if err := ValidateHSL(hue, saturation, luminance, "orange"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Orange = ColorAdjustment{Hue: hue, Saturation: saturation, Luminance: luminance}
	}
	return b
}

// WithYellowHSL sets the yellow color HSL adjustments with validation.
func (b *RecipeBuilder) WithYellowHSL(hue, saturation, luminance int) *RecipeBuilder {
	if err := ValidateHSL(hue, saturation, luminance, "yellow"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Yellow = ColorAdjustment{Hue: hue, Saturation: saturation, Luminance: luminance}
	}
	return b
}

// WithGreenHSL sets the green color HSL adjustments with validation.
func (b *RecipeBuilder) WithGreenHSL(hue, saturation, luminance int) *RecipeBuilder {
	if err := ValidateHSL(hue, saturation, luminance, "green"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Green = ColorAdjustment{Hue: hue, Saturation: saturation, Luminance: luminance}
	}
	return b
}

// WithAquaHSL sets the aqua color HSL adjustments with validation.
func (b *RecipeBuilder) WithAquaHSL(hue, saturation, luminance int) *RecipeBuilder {
	if err := ValidateHSL(hue, saturation, luminance, "aqua"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Aqua = ColorAdjustment{Hue: hue, Saturation: saturation, Luminance: luminance}
	}
	return b
}

// WithBlueHSL sets the blue color HSL adjustments with validation.
func (b *RecipeBuilder) WithBlueHSL(hue, saturation, luminance int) *RecipeBuilder {
	if err := ValidateHSL(hue, saturation, luminance, "blue"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Blue = ColorAdjustment{Hue: hue, Saturation: saturation, Luminance: luminance}
	}
	return b
}

// WithPurpleHSL sets the purple color HSL adjustments with validation.
func (b *RecipeBuilder) WithPurpleHSL(hue, saturation, luminance int) *RecipeBuilder {
	if err := ValidateHSL(hue, saturation, luminance, "purple"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Purple = ColorAdjustment{Hue: hue, Saturation: saturation, Luminance: luminance}
	}
	return b
}

// WithMagentaHSL sets the magenta color HSL adjustments with validation.
func (b *RecipeBuilder) WithMagentaHSL(hue, saturation, luminance int) *RecipeBuilder {
	if err := ValidateHSL(hue, saturation, luminance, "magenta"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.Magenta = ColorAdjustment{Hue: hue, Saturation: saturation, Luminance: luminance}
	}
	return b
}

// WithColorGrading sets Nikon Flexible Color Picture Control color grading with validation (NP3-specific).
// Takes pre-constructed ColorGradingZone structs for each tonal zone plus global blending and balance parameters.
func (b *RecipeBuilder) WithColorGrading(highlights, midtone, shadows ColorGradingZone, blending, balance int) *RecipeBuilder {
	// Validate highlights zone
	if err := ValidateColorGradingHue(highlights.Hue); err != nil {
		b.errors = append(b.errors, fmt.Errorf("highlights: %w", err))
	}
	if err := ValidateColorGradingChroma(highlights.Chroma); err != nil {
		b.errors = append(b.errors, fmt.Errorf("highlights: %w", err))
	}
	if err := ValidateColorGradingBrightness(highlights.Brightness); err != nil {
		b.errors = append(b.errors, fmt.Errorf("highlights: %w", err))
	}

	// Validate midtone zone
	if err := ValidateColorGradingHue(midtone.Hue); err != nil {
		b.errors = append(b.errors, fmt.Errorf("midtone: %w", err))
	}
	if err := ValidateColorGradingChroma(midtone.Chroma); err != nil {
		b.errors = append(b.errors, fmt.Errorf("midtone: %w", err))
	}
	if err := ValidateColorGradingBrightness(midtone.Brightness); err != nil {
		b.errors = append(b.errors, fmt.Errorf("midtone: %w", err))
	}

	// Validate shadows zone
	if err := ValidateColorGradingHue(shadows.Hue); err != nil {
		b.errors = append(b.errors, fmt.Errorf("shadows: %w", err))
	}
	if err := ValidateColorGradingChroma(shadows.Chroma); err != nil {
		b.errors = append(b.errors, fmt.Errorf("shadows: %w", err))
	}
	if err := ValidateColorGradingBrightness(shadows.Brightness); err != nil {
		b.errors = append(b.errors, fmt.Errorf("shadows: %w", err))
	}

	// Validate global parameters
	if err := ValidateColorGradingBlending(blending); err != nil {
		b.errors = append(b.errors, err)
	}
	if err := ValidateColorGradingBalance(balance); err != nil {
		b.errors = append(b.errors, err)
	}

	// Set the color grading if no errors
	if len(b.errors) == 0 {
		b.recipe.ColorGrading = &ColorGrading{
			Highlights: highlights,
			Midtone:    midtone,
			Shadows:    shadows,
			Blending:   blending,
			Balance:    balance,
		}
	}

	return b
}

// WithToneCurve sets parametric tone curve parameters with validation.
func (b *RecipeBuilder) WithToneCurve(shadows, darks, lights, highlights int) *RecipeBuilder {
	if err := ValidatePercentage(shadows, "tone curve shadows"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.ToneCurveShadows = shadows
	}

	if err := ValidatePercentage(darks, "tone curve darks"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.ToneCurveDarks = darks
	}

	if err := ValidatePercentage(lights, "tone curve lights"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.ToneCurveLights = lights
	}

	if err := ValidatePercentage(highlights, "tone curve highlights"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.ToneCurveHighlights = highlights
	}

	return b
}

// WithPointCurve sets the master point curve.
func (b *RecipeBuilder) WithPointCurve(points []ToneCurvePoint) *RecipeBuilder {
	for _, point := range points {
		if err := ValidateToneCurvePoint(point.Input, point.Output); err != nil {
			b.errors = append(b.errors, err)
			return b
		}
	}
	b.recipe.PointCurve = points
	return b
}

// WithSplitToning sets split toning parameters with validation.
func (b *RecipeBuilder) WithSplitToning(shadowHue, shadowSat, highlightHue, highlightSat, balance int) *RecipeBuilder {
	if err := ValidateHue360(shadowHue, "split shadow hue"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.SplitShadowHue = shadowHue
	}

	if err := ValidatePositivePercentage(shadowSat, "split shadow saturation"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.SplitShadowSaturation = shadowSat
	}

	if err := ValidateHue360(highlightHue, "split highlight hue"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.SplitHighlightHue = highlightHue
	}

	if err := ValidatePositivePercentage(highlightSat, "split highlight saturation"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.SplitHighlightSaturation = highlightSat
	}

	if err := ValidatePercentage(balance, "split balance"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.SplitBalance = balance
	}

	return b
}

// WithCameraProfile sets camera calibration parameters.
func (b *RecipeBuilder) WithCameraProfile(profile CameraProfile) *RecipeBuilder {
	b.recipe.CameraProfile = profile
	return b
}

// WithCameraProfileName sets the camera profile name.
func (b *RecipeBuilder) WithCameraProfileName(name string) *RecipeBuilder {
	b.recipe.CameraProfileName = name
	return b
}

// WithGrain sets grain effect parameters with validation.
func (b *RecipeBuilder) WithGrain(amount, size, roughness int) *RecipeBuilder {
	if err := ValidatePositivePercentage(amount, "grain amount"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.GrainAmount = amount
	}

	if err := ValidatePositivePercentage(size, "grain size"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.GrainSize = size
	}

	if err := ValidatePositivePercentage(roughness, "grain roughness"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.GrainRoughness = roughness
	}

	return b
}

// WithGrainAmount sets the grain amount.
func (b *RecipeBuilder) WithGrainAmount(value int) *RecipeBuilder {
	if err := ValidatePositivePercentage(value, "grain amount"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.GrainAmount = value
	}
	return b
}

// WithGrainSize sets the grain size.
func (b *RecipeBuilder) WithGrainSize(value int) *RecipeBuilder {
	if err := ValidatePositivePercentage(value, "grain size"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.GrainSize = value
	}
	return b
}

// WithGrainRoughness sets the grain roughness.
func (b *RecipeBuilder) WithGrainRoughness(value int) *RecipeBuilder {
	if err := ValidatePositivePercentage(value, "grain roughness"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.GrainRoughness = value
	}
	return b
}

// WithVignette sets vignette parameters with validation.
func (b *RecipeBuilder) WithVignette(amount, midpoint, roundness, feather int) *RecipeBuilder {
	if err := ValidatePercentage(amount, "vignette amount"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.VignetteAmount = amount
	}

	if err := ValidatePositivePercentage(midpoint, "vignette midpoint"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.VignetteMidpoint = midpoint
	}

	if err := ValidatePercentage(roundness, "vignette roundness"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.VignetteRoundness = roundness
	}

	if err := ValidatePositivePercentage(feather, "vignette feather"); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.recipe.VignetteFeather = feather
	}

	return b
}

// Build validates all parameters and returns an immutable UniversalRecipe.
// Returns an error if any validation errors were accumulated during building.
func (b *RecipeBuilder) Build() (*UniversalRecipe, error) {
	// Return any errors accumulated during building
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("validation errors: %v", b.errors)
	}

	// Perform final validation of the complete recipe
	if err := b.recipe.Validate(); err != nil {
		return nil, fmt.Errorf("final validation failed: %w", err)
	}

	// Return a copy of the recipe (immutability)
	recipe := b.recipe
	return &recipe, nil
}
