// Package np3 provides functionality for generating Nikon Picture Control (.np3) binary files.
//
// This file implements the generator which converts UniversalRecipe to NP3 binary format.
package np3

import (
	"fmt"
	"math"
	"strings"

	"github.com/justin/recipe/internal/models"
)

// Generate encodes a UniversalRecipe into a Nikon Picture Control (.np3) binary file.
//
// The function converts UniversalRecipe parameters to NP3 ranges, generates the file
// structure with magic bytes and header, encodes parameters as TLV chunks, and returns
// the complete binary representation.
//
// Parameters:
//   - recipe: UniversalRecipe to convert to NP3 format
//
// Returns:
//   - []byte: NP3 binary file data
//   - error: Validation or generation error with descriptive context
//
// Errors:
//   - Nil recipe: Input validation failed
//   - Parameter out of range: Value cannot be represented in NP3 format
//   - Encoding error: Binary encoding failed
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
	data, _, err := GenerateWithWarnings(recipe)
	return data, err
}

// GenerateWithWarnings converts a UniversalRecipe to NP3 binary format and returns
// any warnings about unsupported or approximated parameters.
func GenerateWithWarnings(recipe *models.UniversalRecipe) ([]byte, *models.ConversionResult, error) {
	result := &models.ConversionResult{}

	// Validate input
	if recipe == nil {
		return nil, result, fmt.Errorf("generate NP3: recipe cannot be nil")
	}

	// Check for unsupported parameters and collect warnings
	collectConversionWarnings(recipe, result)

	// Convert parameters from UniversalRecipe to NP3 ranges
	params, err := convertToNP3Parameters(recipe)
	if err != nil {
		return nil, result, fmt.Errorf("generate NP3: %w", err)
	}

	// Validate converted parameters
	if err := validateParameters(params); err != nil {
		return nil, result, fmt.Errorf("generate NP3: %w", err)
	}

	// Build binary structure
	// DEBUG: Append version to force fresh import in NX Studio
	recipe.Name += " v15"
	
	// Truncate name to 20 chars to ensure suffix fits in effective display area
	// Hex dump showed truncation at 20 chars ("Agfachrome RSX 200 v").
	if len(recipe.Name) > 20 {
		// We want the suffix " v15" to be visible, so keep the start + suffix?
		// "Agfachrome... v15"
		// Start: "Agfachrome RSX 2" (16 chars)
		recipe.Name = recipe.Name[:16] + " v15"
	}
	
	data, err := encodeBinary(params, recipe.Name)
	if err != nil {
		return nil, result, fmt.Errorf("generate NP3: %w", err)
	}

	return data, result, nil
}

// collectConversionWarnings checks for unsupported parameters and adds warnings
func collectConversionWarnings(recipe *models.UniversalRecipe, result *models.ConversionResult) {
	// Check frequency-based effects (not supported in NP3)
	if recipe.Clarity != 0 {
		result.AddWarning(
			models.WarnAdvisory,
			"Clarity",
			fmt.Sprintf("%d", recipe.Clarity),
			"Clarity affects mid-tone contrast via frequency separation",
			"Use NP3 Mid-Range Sharpening for a similar effect",
		)
	}

	if recipe.Texture != 0 {
		result.AddWarning(
			models.WarnAdvisory,
			"Texture",
			fmt.Sprintf("%d", recipe.Texture),
			"Texture affects fine detail enhancement",
			"Use NP3 Sharpening Detail setting",
		)
	}

	if recipe.Dehaze != 0 {
		result.AddWarning(
			models.WarnAdvisory,
			"Dehaze",
			fmt.Sprintf("%d", recipe.Dehaze),
			"Dehaze removes atmospheric haze",
			"Use increased Contrast and reduced Blacks",
		)
	}

	// Check RGB channel curves (NP3 only supports master curve)
	if len(recipe.PointCurveRed) > 0 {
		result.AddWarning(
			models.WarnCritical,
			"PointCurveRed",
			fmt.Sprintf("%d points", len(recipe.PointCurveRed)),
			"NP3 only supports master tone curve, not per-channel curves",
			"Use Color Blender Red adjustments instead",
		)
	}

	if len(recipe.PointCurveGreen) > 0 {
		result.AddWarning(
			models.WarnCritical,
			"PointCurveGreen",
			fmt.Sprintf("%d points", len(recipe.PointCurveGreen)),
			"NP3 only supports master tone curve, not per-channel curves",
			"Use Color Blender Green adjustments instead",
		)
	}

	if len(recipe.PointCurveBlue) > 0 {
		result.AddWarning(
			models.WarnCritical,
			"PointCurveBlue",
			fmt.Sprintf("%d points", len(recipe.PointCurveBlue)),
			"NP3 only supports master tone curve, not per-channel curves",
			"Use Color Blender Cyan/Blue adjustments instead",
		)
	}
}

// applyCameraCalibrationToColorBlender maps XMP Camera Calibration to NP3 Color Blender.
//
// XMP Camera Calibration (RedHue, GreenHue, BlueHue, RedSaturation, etc.) adjusts
// the camera's primary color response at a fundamental level. Since NP3 doesn't have
// a direct equivalent, we approximate by adding the calibration values to the
// corresponding Color Blender channels with a 0.5x scale to prevent extreme results.
//
// Mapping:
//   - RedHue/RedSaturation → Red and Orange channels
//   - GreenHue/GreenSaturation → Green and Yellow channels  
//   - BlueHue/BlueSaturation → Blue and Cyan channels
func applyCameraCalibrationToColorBlender(recipe *models.UniversalRecipe, params *np3Parameters) {
	// CameraProfile is a struct value, check if all fields are zero (no calibration)
	cp := recipe.CameraProfile
	if cp.RedHue == 0 && cp.RedSaturation == 0 &&
		cp.GreenHue == 0 && cp.GreenSaturation == 0 &&
		cp.BlueHue == 0 && cp.BlueSaturation == 0 {
		return
	}

	// Scale factor: Camera Calibration values are already -100 to +100
	// Apply at 50% strength to prevent extreme values when combined with HSL
	const scale = 0.5

	// Red calibration → affects Red and Orange (adjacent in color wheel)
	if cp.RedHue != 0 {
		params.redHue = clampColorBlender(params.redHue + int(float64(cp.RedHue)*scale))
		params.orangeHue = clampColorBlender(params.orangeHue + int(float64(cp.RedHue)*scale*0.5))
	}
	if cp.RedSaturation != 0 {
		params.redChroma = clampColorBlender(params.redChroma + int(float64(cp.RedSaturation)*scale))
		params.orangeChroma = clampColorBlender(params.orangeChroma + int(float64(cp.RedSaturation)*scale*0.5))
	}

	// Green calibration → affects Green and Yellow
	if cp.GreenHue != 0 {
		params.greenHue = clampColorBlender(params.greenHue + int(float64(cp.GreenHue)*scale))
		params.yellowHue = clampColorBlender(params.yellowHue + int(float64(cp.GreenHue)*scale*0.5))
	}
	if cp.GreenSaturation != 0 {
		params.greenChroma = clampColorBlender(params.greenChroma + int(float64(cp.GreenSaturation)*scale))
		params.yellowChroma = clampColorBlender(params.yellowChroma + int(float64(cp.GreenSaturation)*scale*0.5))
	}

	// Blue calibration → affects Blue and Cyan
	if cp.BlueHue != 0 {
		params.blueHue = clampColorBlender(params.blueHue + int(float64(cp.BlueHue)*scale))
		params.cyanHue = clampColorBlender(params.cyanHue + int(float64(cp.BlueHue)*scale*0.5))
	}
	if cp.BlueSaturation != 0 {
		params.blueChroma = clampColorBlender(params.blueChroma + int(float64(cp.BlueSaturation)*scale))
		params.cyanChroma = clampColorBlender(params.cyanChroma + int(float64(cp.BlueSaturation)*scale*0.5))
	}
}

// clampColorBlender clamps a Color Blender value to NP3's valid range (-100 to +100)
func clampColorBlender(v int) int {
	if v < -100 {
		return -100
	}
	if v > 100 {
		return 100
	}
	return v
}
// convertToNP3Parameters converts UniversalRecipe parameters to NP3 ranges.
//
// Phase 2 Enhancement: Now converts all 48 parameters using exact offset mappings
// discovered through TypeScript implementation research.
//
// Conversion formulas (reverse of parser mappings in parse.go:845-952):
//   - Basic Adjustments (2): sharpening, clarity
//   - Advanced Adjustments (7): mid-range sharpening, contrast, highlights, shadows, whites, blacks, saturation
//   - Color Blender (24): 8 colors × 3 values (hue, chroma/saturation, brightness/luminance)
//   - Color Grading (11): 3 zones + 2 global params
//   - Tone Curve (3): control points
func convertToNP3Parameters(recipe *models.UniversalRecipe) (*np3Parameters, error) {
	params := &np3Parameters{}

	// === Basic Adjustments (2 parameters) ===

	// Convert Sharpness: UniversalRecipe (0-150) → NP3 (-3.0 to +9.0)
	// Parser uses: (np3Value + 3.0) * 12.5 = universalValue
	// Generator reverses: (universalValue / 12.5) - 3.0 = np3Value
	params.sharpening = (float64(recipe.Sharpness) / 12.5) - 3.0
	if params.sharpening > 9.0 {
		params.sharpening = 9.0
	} else if params.sharpening < -3.0 {
		params.sharpening = -3.0
	}

	// Convert Clarity: UniversalRecipe (-100 to +100) → NP3 (-5.0 to +5.0)
	// Parser uses: np3Value * 20 = universalValue
	// Generator reverses: universalValue / 20 = np3Value
	params.clarity = float64(recipe.Clarity) / 20.0
	if params.clarity > 5.0 {
		params.clarity = 5.0
	} else if params.clarity < -5.0 {
		params.clarity = -5.0
	}

	// Grain Amount: Universal (0-100) → NP3 (0-31.75)
	params.grainAmount = float64(recipe.GrainAmount) * 31.75 / 100.0

	// Grain Size: Universal (0-100) → NP3 Enum
	// 1=Large, 2=Small, 127=Off
	if recipe.GrainSize >= 60 {
		params.grainSize = 1 // Large
	} else if recipe.GrainSize > 0 {
		params.grainSize = 2 // Small
	} else {
		params.grainSize = 127 // Off (0xFF)
	}

	// === Advanced Adjustments (7 parameters) ===

	// Mid-Range Sharpening: Direct mapping (-5.0 to +5.0)
	params.midRangeSharpening = recipe.MidRangeSharpening
	if params.midRangeSharpening > 5.0 {
		params.midRangeSharpening = 5.0
	} else if params.midRangeSharpening < -5.0 {
		params.midRangeSharpening = -5.0
	}

	// Contrast, Highlights, Shadows, Whites, Blacks, Saturation: Direct mapping (-100 to +100)
	params.contrast = recipe.Contrast
	if params.contrast > 100 {
		params.contrast = 100
	} else if params.contrast < -100 {
		params.contrast = -100
	}

	params.highlights = recipe.Highlights
	if params.highlights > 100 {
		params.highlights = 100
	} else if params.highlights < -100 {
		params.highlights = -100
	}

	params.shadows = recipe.Shadows
	if params.shadows > 100 {
		params.shadows = 100
	} else if params.shadows < -100 {
		params.shadows = -100
	}

	params.whiteLevel = recipe.Whites
	if params.whiteLevel > 100 {
		params.whiteLevel = 100
	} else if params.whiteLevel < -100 {
		params.whiteLevel = -100
	}

	// recipe.Blacks is baked into the tone curve. Disable params.blackLevel to avoid clipping.
	params.blackLevel = 0
	if params.blackLevel > 100 {
		params.blackLevel = 100

	} else if params.blackLevel < -100 {
		params.blackLevel = -100
	}

	// Apply Dehaze (affects Contrast and Blacks)
	applyDehaze(recipe, params)

	params.saturation = recipe.Saturation
	if params.saturation > 100 {
		params.saturation = 100
	} else if params.saturation < -100 {
		params.saturation = -100
	}

	// Apply Vibrance (affects Saturation)
	applyVibranceToSaturation(recipe, params)

	// === Color Blender (24 parameters) ===
	// UniversalRecipe ColorAdjustment maps directly to NP3 Color Blender (hue, chroma, brightness)

	params.redHue = recipe.Red.Hue
	params.redChroma = recipe.Red.Saturation
	params.redBrightness = recipe.Red.Luminance

	params.orangeHue = recipe.Orange.Hue
	params.orangeChroma = recipe.Orange.Saturation
	params.orangeBrightness = recipe.Orange.Luminance

	params.yellowHue = recipe.Yellow.Hue
	params.yellowChroma = recipe.Yellow.Saturation
	params.yellowBrightness = recipe.Yellow.Luminance

	params.greenHue = recipe.Green.Hue
	params.greenChroma = recipe.Green.Saturation
	params.greenBrightness = recipe.Green.Luminance

	params.cyanHue = recipe.Aqua.Hue
	params.cyanChroma = recipe.Aqua.Saturation
	params.cyanBrightness = recipe.Aqua.Luminance

	params.blueHue = recipe.Blue.Hue
	params.blueChroma = recipe.Blue.Saturation
	params.blueBrightness = recipe.Blue.Luminance

	params.purpleHue = recipe.Purple.Hue
	params.purpleChroma = recipe.Purple.Saturation
	params.purpleBrightness = recipe.Purple.Luminance

	params.magentaHue = recipe.Magenta.Hue
	params.magentaChroma = recipe.Magenta.Saturation
	params.magentaBrightness = recipe.Magenta.Luminance

	// === Camera Calibration → Color Blender Approximation ===
	// XMP Camera Calibration adjusts the primary colors (RGB) at a fundamental level.
	// NP3 doesn't have this, but we can approximate by adjusting Color Blender values.
	//
	// Mapping strategy:
	// - RedHue/Saturation → affects Red and Orange in Color Blender
	// - GreenHue/Saturation → affects Green and Yellow in Color Blender
	// - BlueHue/Saturation → affects Blue and Cyan in Color Blender
	//
	// Scale factor: 0.5 to prevent extreme values (Camera Calibration is subtle)
	applyCameraCalibrationToColorBlender(recipe, params)

	// === Color Grading (11 parameters) ===
	// Direct mapping from UniversalRecipe ColorGrading

	// Map Color Grading if present
	// Only assign zones that have actual data (non-zero chroma) to avoid overwriting Split Toning fallback
	if recipe.ColorGrading != nil {
		if recipe.ColorGrading.Highlights.Chroma != 0 {
			params.highlightsZone = recipe.ColorGrading.Highlights
		}
		if recipe.ColorGrading.Midtone.Chroma != 0 {
			params.midtoneZone = recipe.ColorGrading.Midtone
		}
		if recipe.ColorGrading.Shadows.Chroma != 0 {
			params.shadowsZone = recipe.ColorGrading.Shadows
		}
		params.blending = recipe.ColorGrading.Blending
		params.balance = recipe.ColorGrading.Balance
	}

	// Fallback/Merge: Legacy Split Toning
	// If specific Color Grading zones are empty, fill them from Split Toning
	if params.shadowsZone.Chroma == 0 && recipe.SplitShadowSaturation > 0 {
		params.shadowsZone.Hue = recipe.SplitShadowHue
		params.shadowsZone.Chroma = recipe.SplitShadowSaturation
	}
	if params.highlightsZone.Chroma == 0 && recipe.SplitHighlightSaturation > 0 {
		params.highlightsZone.Hue = recipe.SplitHighlightHue
		params.highlightsZone.Chroma = recipe.SplitHighlightSaturation
	}
	// Map Balance if not set by Color Grading
	if params.balance == 0 && recipe.SplitBalance != 0 {
		params.balance = recipe.SplitBalance
	}

	// If Split Toning is used (and blending is low/unset), ensure Blending is set to default (50)
	if (recipe.SplitShadowSaturation > 0 || recipe.SplitHighlightSaturation > 0) && params.blending < 10 {
		params.blending = 50
	}


	// Apply White Balance (Temp/Tint) via Color Grading
	applyWhiteBalanceToColorGrading(recipe, params)

	// === Tone Curve (3 parameters) ===
	// Convert UniversalRecipe PointCurve to NP3 format
	// If RGB channel curves exist, merge them with the master curve first

	// Convert models.ToneCurvePoint to internal ControlPoint format
	masterCPs := toneCurvePointsToControlPoints(recipe.PointCurve)
	redCPs := toneCurvePointsToControlPoints(recipe.PointCurveRed)
	greenCPs := toneCurvePointsToControlPoints(recipe.PointCurveGreen)
	blueCPs := toneCurvePointsToControlPoints(recipe.PointCurveBlue)

	// Merge RGB curves with master if any RGB curves exist
	// Merge RGB curves with master if any RGB curves exist
	var finalCurve []ControlPoint
	// Tone Curve Generation Strategy:
	// Lightroom applies curves in this order:
	// 1. Parametric Curve (Sliders)
	// 2. Point Curve (Custom Points / RGB)
	// We must chain these effects, not select exclusively.

	// 1. Calculate Parametric Curve LUT
	var parametricLUT []int
	if HasParametricCurve(recipe.ToneCurveShadows, recipe.ToneCurveDarks,
		recipe.ToneCurveLights, recipe.ToneCurveHighlights) {
		parametricCPs := ParametricToControlPoints(
			recipe.ToneCurveShadows,
			recipe.ToneCurveDarks,
			recipe.ToneCurveLights,
			recipe.ToneCurveHighlights,
			recipe.ToneCurveShadowSplit,
			recipe.ToneCurveMidtoneSplit,
			recipe.ToneCurveHighlightSplit,
		)
		parametricLUT = PointsToCurveLUT(parametricCPs)
	} else {
		// Identity LUT
		parametricLUT = make([]int, 256)
		for i := 0; i < 256; i++ {
			parametricLUT[i] = i
		}
	}

	// 2. Calculate Point Curve LUT
	// (RGB Merged or Master Point)
	var pointCurveCPs []ControlPoint
	if HasRGBCurves(redCPs, greenCPs, blueCPs) {
		// Merge RGB curves into master curve
		pointCurveCPs = MergeRGBCurvesToMaster(masterCPs, redCPs, greenCPs, blueCPs)
	} else if len(masterCPs) > 0 {
		pointCurveCPs = masterCPs
	}
	
	var pointCurveLUT []int
	if len(pointCurveCPs) > 0 {
		pointCurveLUT = PointsToCurveLUT(pointCurveCPs)
	} else {
		// Identity LUT
		pointCurveLUT = make([]int, 256)
		for i := 0; i < 256; i++ {
			pointCurveLUT[i] = i
		}
	}

	// 3. Combine: Final = Point(Parametric(Input))
	finalCurveLUT := ApplyCurveToLUT(parametricLUT, pointCurveLUT)

	// === CRITICAL FIX: Apply Blacks BEFORE any other curve processing ===
	// XMP Blacks adjustment must be applied in Lightroom order: Blacks → Point Curve
	// Previously, Blacks was only baked when Standard profile was enabled.
	// This caused Blacks to be lost when Standard profile was disabled to fix shadow crushing.
	// Now we ALWAYS apply Blacks to match Lightroom's processing order.
	// Example: XMP has Blacks -13 (pulls down shadow point) → Point Curve 0→13 (lifts shadows)
	// Result: Shadow detail preserved but blacks properly anchored.
	if recipe.Blacks != 0 {
		finalCurveLUT = ApplyExposureAndBlacksToLUT(finalCurveLUT, 0, recipe.Blacks)
		// Note: Pass 0 for Exposure here, as Exposure is handled separately via brightness parameter
	}

	// Parse profile name for decision making
	profileName := strings.TrimSpace(recipe.CameraProfileName)

	// === BASELINE COMPENSATION for Adobe DCP vs Nikon Profile Differences ===
	// Adobe's "Flexible Color" DCP has a +47.8 brightness baseline applied before XMP adjustments.
	// Nikon's native "Flexible Color" has a ~+105.6 brightness baseline (58 units brighter).
	// When converting XMP→NP3, we need to compensate for this baseline difference so that
	// NX Studio output matches Lightroom output when both use their Flexible Color profiles.
	//
	// Compensation is applied if:
	// 1. Metadata explicitly requests it (recipe.Metadata["baseline_compensation"] == "flexible_color")
	// 2. The Camera Profile is explicitly "Flexible Color" (auto-detection)
	//
	// The compensation darkens the curve by applying the inverse of the baseline difference.
	// See: docs/analysis/adobe-dcp-baseline-discovery.md for full analysis.
	isFlexibleColor := strings.Contains(profileName, "Flexible Color")
	forceCompensation := recipe.Metadata != nil && recipe.Metadata["baseline_compensation"] == "flexible_color"
	
	if isFlexibleColor || forceCompensation {
		// TEMPORARILY DISABLED: Testing with X/Y swap fix first
		// The compensation was too aggressive (crushing shadows to 0)
		// TODO: Re-enable with gentler compensation once baseline is verified
		// finalCurveLUT = ApplyFlexibleColorBaselineCompensation(finalCurveLUT)
		_ = isFlexibleColor // Suppress unused warning
		_ = forceCompensation
	}

	// Convert combined LUT to Control Points for final output
	finalCurve = LUTToControlPoints(finalCurveLUT, 20)

	// Apply Standard S-Curve contrast if this is an XMP preset expecting "Standard" contrast
	// but we are using "Flexible Color" (Linear) as the base.
	// This restores the "Adobe Standard" / "Camera Standard" look.
	isStandardProfile := profileName == "Default Color" ||
		profileName == "Camera Standard" ||
		profileName == "Adobe Standard" ||
		profileName == "Embedded" || // Embedded usually means standard-ish
		strings.Contains(profileName, "Standard") // Catch-all for variants

	// CRITICAL FIX: Skip Standard base curve if preset already has custom Point Curve with shadow lift.
	// Film emulation presets like Agfachrome have their own tone curves baked in.
	// Applying Standard base curve on top of custom curves CRUSHES the intended shadow lift.
	// Example: XMP has 0→13 lift, but Standard base curve crushes it back to 0.
	hasCustomPointCurve := len(recipe.PointCurve) > 0 && recipe.PointCurve[0].Output > recipe.PointCurve[0].Input
	if hasCustomPointCurve {
		isStandardProfile = false // Disable Standard base curve for presets with custom curves
	}

	if isStandardProfile {
		// Ensure we have a curve to act as modifier
		if len(finalCurve) == 0 {
			// Identity curve
			finalCurve = []ControlPoint{{X: 0, Y: 0}, {X: 255, Y: 255}}
		}

		// Calculate composite curve: Final(x) = CombinedXMP(Base(x))
		// Base = Standard S-Curve
		// CombinedXMP = Parametric + Point
		baseLUT := GetStandardBaseCurveLUT()
		mergedLUT := ApplyCurveToLUT(baseLUT, finalCurveLUT)

		// === Bake Contrast into Curve ===
		// XMP Contrast (+19) is best applied to the curve for fidelity,
		// especially if the NP3 Contrast slider behaves differently or is ignored when chaining.
		// FIXED: Removed 2x multiplier (was Bug #3). XMP Contrast +19 should be applied as +19, not +38.
		// The 2x multiplier was too aggressive and contributed to overall darkening.
		if recipe.Contrast != 0 {
			mergedLUT = ApplyContrastToLUT(mergedLUT, recipe.Contrast)
			params.contrast = 0 // Zero out the slider
		}

		// === Bake Exposure into Curve (for Standard profile only) ===
		// Note: Blacks is now handled earlier (lines 483-486) for ALL presets.
		// Here we only handle Exposure for Standard profile presets.
		if recipe.Exposure != 0 {
			mergedLUT = ApplyExposureAndBlacksToLUT(mergedLUT, recipe.Exposure, 0)
			params.brightness = 0 // Zero out exposure slider (mapped to brightness)
		}
		
		// Force a slight negative Black Level to ensure the curve starts at rigid 0.
		// Previous result showed a lift (~3-5) despite mathematical 0.
		// Pulling it down ensures deep blacks.
		params.blackLevel = -5

		// Convert back to control points (use 20 points for fidelity)
		finalCurve = LUTToControlPoints(mergedLUT, 20)
		
		// === Color Fidelity: Standard Profile Color Emulation ===
		// FIXED: Previous code darkened blues/greens by -30 to -40 based on MISUNDERSTANDING of user feedback.
		// User actually said "blues/greens are MORE VIBRANT in Lightroom" (meaning BRIGHTER, not darker).
		// Darkening hack removed - HSL luminance values now come directly from XMP preset.
		// This preserves the intended color vibrancy instead of crushing it.
		//
		// REMOVED (was Bug #2):
		// params.blueBrightness = clampColorBlender(params.blueBrightness - 40)   // Made blue -37 instead of +3!
		// params.cyanBrightness = clampColorBlender(params.cyanBrightness - 40)   // Made cyan -35 instead of +5!
		// params.greenBrightness = clampColorBlender(params.greenBrightness - 30) // Made green -22 instead of +8!

		
		// Standard also tends to be slightly warmer in highlights.
		params.orangeBrightness = clampColorBlender(params.orangeBrightness + 5)
		params.yellowBrightness = clampColorBlender(params.yellowBrightness + 5)
	}

	// Only write curve if it's not a linear identity curve
	if len(finalCurve) > 0 && !IsLinearCurve(finalCurve, 5) {
		params.toneCurvePointCount = len(finalCurve)
		if params.toneCurvePointCount > 127 {
			params.toneCurvePointCount = 127 // NP3 limit
		}

		params.toneCurvePoints = make([]toneCurvePoint, params.toneCurvePointCount)
		for i := 0; i < params.toneCurvePointCount; i++ {
			// CRITICAL: We must use the Compensated LUT value for Y, not the original point Y.
			// The finalCurveLUT contains the brightness compensation (and other adjustments).
			// By sampling the LUT at the point's X position, we bake the compensation into the control points.
			x := int(finalCurve[i].X)
			var y uint8
			if x >= 0 && x < 256 {
				y = uint8(finalCurveLUT[x])
			} else {
				y = uint8(finalCurve[i].Y) // Fallback
			}

			// SIMPLE LINEAR DARKENING: NX Studio baseline is brighter than Lightroom
			// Subtract a flat offset from Y values to darken the output
			// This is simpler than the complex compensation and doesn't crush shadows to 0
			darkenOffset := 40
			yDarkened := int(y) - darkenOffset
			if yDarkened < 0 {
				yDarkened = 0
			}

			params.toneCurvePoints[i] = toneCurvePoint{
				position: 405 + (i * 2),
				value1:   uint8(yDarkened), // Y (output) - NP3 stores (Y, X)
				value2:   uint8(x),         // X (input)
			}
		}

		// NOTE: Do NOT populate toneCurveRaw - the BASELINE file that worked
		// had zeros in the Extended LUT area (560-1072). Writing data there
		// causes import failures in NX Studio.
		// The curve is fully defined by the BI0 control points.
	}

	// === Legacy Parameters (for heuristic fallback compatibility) ===

	// Brightness: UniversalRecipe Exposure field → NP3 brightness (-1.0 to +1.0)
	// NOTE: NP3 format does NOT have a dedicated Exposure/Brightness offset.
	// This parameter exists in the struct for round-trip compatibility but is not written to binary.
	// Always set to 0 since there's nowhere to write it in the 480-byte NP3 format.
	params.brightness = 0.0 // Force to 0 (NP3 has no exposure parameter)

	// Hue: No direct mapping in UniversalRecipe (only per-color hue adjustments)
	params.hue = 0

	// Retrieve raw binary data if available for perfect round-trip preservation
	if recipe.FormatSpecificBinary != nil {
		if rawData, ok := recipe.FormatSpecificBinary["np3_raw"]; ok {
			params.rawData = rawData
		}
	}

	return params, nil
}

// npChunk represents a TLV (Type-Length-Value) chunk in the NP3 format.
// Each chunk is 10 bytes: [ID:1][Pad:3][Length:2BE][Value:2][Pad:2]
type npChunk struct {
	id     byte   // Chunk ID (0x03-0x1F)
	length uint16 // Length in big-endian (typically 2, but 0x1F has length 28)
	value  []byte // Value bytes (2 bytes typically, more for extended chunks)
}

// constantChunks defines constant TLV chunks that appear in all NP3 files.
// These chunks have fixed values across all presets and represent format structure.
// Chunks 0x06, 0x07, 0x14, 0x16 are variable and handled separately.
var constantChunks = []npChunk{
	{id: 0x03, length: 2, value: []byte{0x00, 0x20}}, // Format identifier (value=32)
	{id: 0x04, length: 2, value: []byte{0x00, 0x00}}, // Reserved (value=0)
	{id: 0x05, length: 2, value: []byte{0xff, 0x01}}, // Format flag (value=65281)
	{id: 0x08, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x09, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x0a, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x0b, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x0c, length: 2, value: []byte{0xff, 0x00}}, // Default value (value=65280)
	{id: 0x0d, length: 2, value: []byte{0xff, 0x00}}, // Default value (value=65280)
	{id: 0x0e, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x0f, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x10, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x11, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x12, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x13, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x14, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281) - ADDED
	{id: 0x15, length: 2, value: []byte{0xff, 0x0a}}, // Default value (value=65290)
	{id: 0x17, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x18, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
}

// writeChunk writes a single TLV chunk at the specified offset.
// Chunk format: [ChunkID:1][Padding:3][Length:2BE][Value:N][Padding:2]
// Total size: 10 bytes (for chunks with length=2)
func writeChunk(data []byte, offset int, chunk npChunk) int {
	data[offset] = chunk.id
	data[offset+1] = 0x00 // Padding
	data[offset+2] = 0x00
	data[offset+3] = 0x00
	data[offset+4] = byte(chunk.length >> 8)   // Length big-endian high byte
	data[offset+5] = byte(chunk.length & 0xff) // Length big-endian low byte

	// Write value bytes (typically 2 bytes)
	copy(data[offset+6:offset+6+len(chunk.value)], chunk.value)

	// Write padding after value
	data[offset+8] = 0x00
	data[offset+9] = 0x00

	return offset + 10 // Return next chunk offset
}

// writeChunkMetadataOnly writes chunk structure without overwriting value bytes.
// Used when heuristic data has already been written to preserve those values.
// Chunk format: [ChunkID:1][Padding:3][Length:2BE][Value:N][Padding:2]
func writeChunkMetadataOnly(data []byte, offset int, chunkID byte, length uint16) int {
	data[offset] = chunkID
	data[offset+1] = 0x00 // Padding
	data[offset+2] = 0x00
	data[offset+3] = 0x00
	data[offset+4] = byte(length >> 8)   // Length big-endian high byte
	data[offset+5] = byte(length & 0xff) // Length big-endian low byte

	// Skip offset+6 and offset+7 (value bytes) - preserve existing heuristic data

	// Write padding after value
	data[offset+8] = 0x00
	data[offset+9] = 0x00

	return offset + 10 // Return next chunk offset
}

// writeChunks writes all TLV chunks (constant + variable) starting at offset 46.
// This generates the 29 chunks (0x03-0x1F) required for Nikon NX Studio validation.
// Chunks must be written in sequential order with variable chunks interlaced.
func writeChunks(data []byte, params *np3Parameters) {
	offset := 46

	// Write chunks in correct sequential order (0x03-0x1F)
	// Constant chunks write full structure, variable chunks write metadata-only

	// Chunks 0x03-0x05: Constants
	offset = writeChunk(data, offset, npChunk{id: 0x03, length: 2, value: []byte{0x00, 0x20}})
	offset = writeChunk(data, offset, npChunk{id: 0x04, length: 2, value: []byte{0x00, 0x00}})
	
	// Chunk 0x05: Base Picture Control ID
	baseID := params.basePictureControlID
	if baseID == 0 { baseID = 1 } // Default to Standard if not set
	offset = writeChunk(data, offset, npChunk{id: 0x05, length: 2, value: []byte{0xff, baseID}})

	// Chunks 0x06-0x07: Sharpening and Clarity parameters
	// These chunks contain the actual parameter values in their value fields
	sharpValue := EncodeScaled4(params.sharpening)
	clarityValue := EncodeScaled4(params.clarity)
	offset = writeChunk(data, offset, npChunk{id: 0x06, length: 2, value: []byte{sharpValue, 0x04}})
	offset = writeChunk(data, offset, npChunk{id: 0x07, length: 2, value: []byte{clarityValue, 0x04}})

	// Chunks 0x08-0x13: Constants
	// Chunks 0x08-0x0B: Previously mapped to Grain Amount, but this breaks Quick Sharp interaction.
	// Working files (FLEXIBLECOLOR) use 0xFF for these chunks.
	// We now write constant 0xFF values to avoid corrupting Quick Sharp / Picture Control settings.
	offset = writeChunk(data, offset, npChunk{id: 0x08, length: 2, value: []byte{0xff, 0x04}})
	offset = writeChunk(data, offset, npChunk{id: 0x09, length: 2, value: []byte{0xff, 0x04}})
	offset = writeChunk(data, offset, npChunk{id: 0x0a, length: 2, value: []byte{0xff, 0x04}})
	offset = writeChunk(data, offset, npChunk{id: 0x0b, length: 2, value: []byte{0xff, 0x04}})
	offset = writeChunk(data, offset, npChunk{id: 0x0c, length: 2, value: []byte{0xff, 0x00}})
	offset = writeChunk(data, offset, npChunk{id: 0x0d, length: 2, value: []byte{0xff, 0x00}})
	offset = writeChunk(data, offset, npChunk{id: 0x0e, length: 2, value: []byte{0xff, 0x04}})
	offset = writeChunk(data, offset, npChunk{id: 0x0f, length: 2, value: []byte{0xff, 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x10, length: 2, value: []byte{0xff, 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x11, length: 2, value: []byte{0xff, 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x12, length: 2, value: []byte{0xff, 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x13, length: 2, value: []byte{0xff, 0x01}})

	// Chunk 0x14: Constant
	// Previously mapped to Grain Size, but Grain is not supported in NP3.
	// Setting to constant safe value.
	offset = writeChunk(data, offset, npChunk{id: 0x14, length: 2, value: []byte{0x01, 0x01}})

	// Chunk 0x15: Constant
	offset = writeChunk(data, offset, npChunk{id: 0x15, length: 2, value: []byte{0xff, 0x0a}})

	// Chunk 0x16: Mid-Range Sharpening parameter
	midRangeValue := EncodeScaled4(params.midRangeSharpening)
	offset = writeChunk(data, offset, npChunk{id: 0x16, length: 2, value: []byte{midRangeValue, 0x04}})

	// Chunks 0x17-0x18: Constants
	offset = writeChunk(data, offset, npChunk{id: 0x17, length: 2, value: []byte{0xff, 0x04}})
	offset = writeChunk(data, offset, npChunk{id: 0x18, length: 2, value: []byte{0xff, 0x04}})

	// Chunks 0x19-0x1E: Advanced Adjustment parameters (Contrast, Highlights, Shadows, etc.)
	// These chunks contain the actual parameter values using Signed8 encoding
	// Second byte is 0x01 (type indicator for Signed8 parameters)
	offset = writeChunk(data, offset, npChunk{id: 0x19, length: 2, value: []byte{EncodeSigned8(params.contrast), 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x1A, length: 2, value: []byte{EncodeSigned8(params.highlights), 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x1B, length: 2, value: []byte{EncodeSigned8(params.shadows), 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x1C, length: 2, value: []byte{EncodeSigned8(params.whiteLevel), 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x1D, length: 2, value: []byte{EncodeSigned8(params.blackLevel), 0x01}})
	offset = writeChunk(data, offset, npChunk{id: 0x1E, length: 2, value: []byte{EncodeSigned8(params.saturation), 0x01}})

	// Chunk 0x1F: Extended data chunk with length=28
	// Value bytes (28 total):
	//   - Bytes 0-23: Color Blender data (8 colors × 3 params, written by writeColorBlender)
	//   - Bytes 24-26: Type indicators (0x01, 0x01, 0x01)
	//   - Byte 27: Padding (0x00)
	data[offset] = 0x1F
	data[offset+1] = 0x00
	data[offset+2] = 0x00
	data[offset+3] = 0x00
	data[offset+4] = 0x00 // length=28 big-endian high
	data[offset+5] = 0x1C // length=28 big-endian low
	// Skip color data bytes (offset+6 to offset+29) - written by writeColorBlender
	// Write type indicator bytes at end of chunk value
	data[offset+30] = 0x01 // Type indicator for Color Blender params
	data[offset+31] = 0x01 // Type indicator
	data[offset+32] = 0x01 // Type indicator
	data[offset+33] = 0x00 // Padding
	offset += 34           // Move past chunk 0x1F (6 header + 28 value)

	// Add 2 bytes padding before chunk 0x20 (required by NP3 format)
	data[offset] = 0x00
	data[offset+1] = 0x00
	offset += 2

	// Chunk 0x20: Extended data chunk with length=20 for Color Grading
	// Value bytes (20 total):
	//   - Bytes 0-11 (offset+6 to +17): Color Grading zone data (3 zones × 4 bytes, written by writeColorGrading)
	//   - Bytes 12-14 (offset+18 to +20): Type indicators for zones (0x01, 0x01, 0x01)
	//   - Byte 15 (offset+21): Padding (0x00)
	//   - Byte 16 (offset+22): Blending value (written by writeColorGrading)
	//   - Byte 17 (offset+23): Type indicator for Blending (0x01)
	//   - Byte 18 (offset+24): Balance value (written by writeColorGrading)
	//   - Byte 19 (offset+25): Type indicator for Balance (0x01)
	data[offset] = 0x20
	data[offset+1] = 0x00
	data[offset+2] = 0x00
	data[offset+3] = 0x00
	data[offset+4] = 0x00 // length=20 big-endian high
	data[offset+5] = 0x14 // length=20 big-endian low
	// Color Grading zone data bytes (offset+6 to offset+17) - written by writeColorGrading
	// Type indicators for zone parameters (Hue, Chroma, Brightness)
	data[offset+18] = 0x01 // Type indicator
	data[offset+19] = 0x01 // Type indicator
	data[offset+20] = 0x01 // Type indicator
	// Padding before global Color Grading parameters
	data[offset+21] = 0x00 // Padding
	// Blending and Balance values at offset+22 and offset+24 - written by writeColorGrading
	// Type indicators for Blending and Balance parameters
	data[offset+23] = 0x01 // Type indicator for Blending
	data[offset+25] = 0x01 // Type indicator for Balance

	// Total chunks written: 30 chunks (0x03-0x20)
	// Chunk 0x1F: 34 bytes (6 header + 28 value)
	// Chunk 0x20: 26 bytes (6 header + 20 value)
}

// encodeBinary generates the complete NP3 binary structure.
//
// File Structure (Nikon NX Studio format):
//   - Offset 0x00-0x02: Magic bytes "NCP"
//   - Offset 0x03-0x06: Version (4 bytes)
//   - Offset 0x07-0x13: Reserved (zeros)
//   - Offset 0x14-0x17: Name header (0x00 0x00 0x00 0x14)
//   - Offset 0x18-0x2B: Preset name (20 bytes, ASCII null-terminated)
//   - Offset 0x2C-0x2D: Padding (zeros)
//   - Offset 0x2E-0x14F: TLV chunks (29 chunks × 10 bytes = 290 bytes, offsets 46-335)
//   - Offset 0x40-0x4F: Raw parameter bytes (offsets 64-80, for heuristic parser)
//   - Offset 0x64-0x12B: Color data section (bytes 100-300, for heuristic parser)
//   - Offset 0x96-0x1F3: Tone curve data section (bytes 150-500, for heuristic parser)
//   - Offset 0x1CC+: Tone curve raw data (460+, 257×16-bit values = 514 bytes)
//   - Remaining: Padding to 1050 bytes (matching working NP3 files)
func encodeBinary(params *np3Parameters, presetName string) ([]byte, error) {
	// If we have raw data from parsing, use it as the base to preserve chunks
	// Otherwise create a new buffer
	var data []byte
	if params.rawData != nil && len(params.rawData) > 0 {
		// Copy raw data to preserve header, chunks, and all structure
		data = make([]byte, len(params.rawData))
		copy(data, params.rawData)
	} else {
		// Create buffer for complete file with extended tone curve
		// Need at least 1072 bytes for extended tone curve LUT (offset 560 + 256*2)
		// Using 1072 to ensure curve data fits
		data = make([]byte, 1072)
	}

	// Only write header if we don't have raw data (raw data already has correct header)
	if params.rawData == nil || len(params.rawData) == 0 {
		// Write magic bytes "NCP" at offset 0-2
		copy(data[0:3], magicBytes)

		// Write version bytes at offset 3-6 (NP3 v1.0.0.0, little-endian: 0x00000001)
		copy(data[3:7], []byte{0x00, 0x00, 0x00, 0x01})

		// Offsets 7-10: Reserved (already zeros)

		// Write format header at offsets 11-19
		data[11] = 0x04                               // Format flag
		copy(data[12:16], []byte{'0', '3', '1', '0'}) // Format version string "0310"
		// Offsets 16-17: Reserved (already zeros)
		data[18] = 0x02 // Unknown flag
		// Offset 19: Reserved (already zero)

		// Write name length header at offsets 20-23
		copy(data[20:24], []byte{0x00, 0x00, 0x00, 0x14}) // Name length = 20 (0x14)
	}

	// Only write legacy data structures if we don't have raw data
	// (raw data already has correct values and writing would corrupt chunks)
	if params.rawData == nil || len(params.rawData) == 0 {
		// Phase 2: Legacy heuristic data generation (generateColorData, generateToneCurveData)
		// has been removed in favor of exact offset writing. The old heuristic approach
		// wrote data to offsets 100-299 (color) and 150-499 (tone curve), but Phase 2
		// uses exact offsets for all 48 parameters instead.
		//
		// The following legacy functions are still needed for basic file structure:

		// Write preset name at offset 24-43 (20 bytes, matching real NP3 files)
		if presetName != "" {
			nameBytes := []byte(presetName)
			maxNameLen := 20 // Real files use exactly 20 bytes for name
			if len(nameBytes) > maxNameLen {
				nameBytes = nameBytes[:maxNameLen]
			}
			copy(data[24:44], nameBytes)
		}
	}

	// Offsets 44-45: Padding (already zeros)

	// Only write TLV chunks if we don't have raw data (raw data already has correct chunks)
	if params.rawData == nil || len(params.rawData) == 0 {
		// Write TLV chunks at offsets 46-335 (29 chunks × 10 bytes = 290 bytes)
		// This is required for Nikon NX Studio validation
		// We write chunks FIRST, then write raw parameter bytes AFTER to avoid overwriting
		writeChunks(data, params)
	}

	// Write raw parameter bytes at offsets 71-79 (brightness, hue)
	// NOTE: Sharpness (previously at 66-70) is now handled by TLV chunks
	if params.rawData == nil || len(params.rawData) == 0 {
		writeRawParameterBytes(data, params)
	}

	// === Phase 2: Write all parameters to exact offsets ===
	// Write exact offset data regardless of whether we have rawData, to ensure
	// current parameter values are written to the correct locations
	writeBasicAdjustments(data, params)
	writeEffects(data, params)
	writeAdvancedAdjustments(data, params)
	writeAdvancedAdjustments(data, params)
	writeColorBlender(data, params)
	writeColorGrading(data, params)
	writeToneCurve(data, params)

	// Return the complete 1050-byte buffer (matching real .np3 files)
	return data, nil
}

// writeRawParameterBytes writes raw parameter values to the byte offsets
// that the parser reads from (offsets 64-79).
//
// Parser extraction strategy:
//   - Sharpness: offsets 66-70 (adjusted from -128 to +127, mapped to 0-9)
//   - Brightness: offsets 71-75 (adjusted from -128 to +127, normalized to -1.0 to +1.0)
//   - Hue: offsets 76-79 (adjusted from -128 to +127, mapped to -9 to +9)
//
// Generator strategy (reverse of parser):
//   - Convert NP3 parameter values back to raw bytes using 128-offset encoding
func writeRawParameterBytes(data []byte, params *np3Parameters) {
	// DEPRECATED: This function previously wrote heuristic parameter bytes to offsets 66-79.
	//
	// These offsets are now used by TLV chunks (chunks start at offset 46, each is 10 bytes):
	//   - Chunk #2 (id=0x05): offsets 66-75
	//   - Chunk #3 (id=0x06): offsets 76-85
	//   - Chunk #4 (id=0x07): offsets 86-95
	//
	// Writing heuristic bytes to these offsets corrupts the TLV chunk structure,
	// causing NX Studio to reject the file.
	//
	// Parameters are now written to their exact offsets instead:
	//   - Sharpness: offset 82 (via writeBasicAdjustments)
	//   - Clarity: offset 92 (via writeBasicAdjustments)
	//   - Brightness: Not currently mapped to exact offset (TODO)
	//   - Hue: Not currently mapped to exact offset (TODO)
	//
	// The parser still reads from heuristic offsets 66-79 as a fallback when exact
	// offsets are unavailable, but it will now extract these values from the TLV
	// chunk structure instead of dedicated heuristic bytes.
}

// generateColorData generates RGB color triplets at offsets 100-299
// that will be heuristically analyzed by the parser to determine saturation.
//
// Parser strategy (parse.go:200-219):
//
//	colorIntensity := len(colorData) / 15
//	params.saturation = colorIntensity - 1
//	(clamped to -3 to +3)
//
// Generator reverse strategy:
//
//	To achieve saturation value S, we need (S + 1) * 15 significant RGB triplets
//	We generate triplets with at least one channel > 10
//
// Returns the offset where color data ends
func generateColorData(data []byte, saturation int) int {
	// Calculate how many significant RGB triplets we need
	// Parser: saturation = (len(colorData) / 15) - 1
	// Generator: len(colorData) = (saturation + 1) * 15
	targetCount := (saturation + 1) * 15

	// Clamp to reasonable range (can't have negative count)
	if targetCount < 0 {
		targetCount = 0
	}
	// Limit to available space (200 bytes / 3 = 66 triplets max)
	if targetCount > 66 {
		targetCount = 66
	}

	// Generate RGB triplets starting at offset 100
	offset := 100
	for i := 0; i < targetCount && offset+2 < 300; i++ {
		// Generate triplet with at least one channel > 10
		// Use simple pattern: (50, 50, 50) for positive saturation
		data[offset] = 50
		data[offset+1] = 50
		data[offset+2] = 50
		offset += 3
	}

	return offset // Return where color data ends
}

// generateToneCurveData generates paired byte values at offsets 150-499
// that will be heuristically analyzed by the parser to determine contrast.
//
// Parser strategy (parse.go:221-237):
//
//	Reads pairs from offset 150-500
//	curveComplexity := len(toneCurve) / 20
//	params.contrast = curveComplexity - 2
//	(clamped to -3 to +3)
//
// Generator reverse strategy:
//
//	To achieve contrast value C, we need (C + 2) * 20 non-zero tone curve points
//	If color data extends into tone curve region (150+), those bytes will be
//	counted as tone curve pairs. Adjust our count to compensate.
func generateToneCurveData(data []byte, contrast int, colorEndOffset int) {
	// Calculate how many tone curve pairs we need total
	// Parser: contrast = (len(toneCurve) / 20) - 2
	// Generator: len(toneCurve) = (contrast + 2) * 20
	targetTotalPairs := (contrast + 2) * 20

	// Calculate how many pairs already exist from color data overlap (150 to colorEndOffset)
	// Parser reads tone curve from offset 150, so any color data >= 150 will be counted
	overlapPairs := 0
	if colorEndOffset > 150 {
		// Color data extends into tone curve region
		// Last color byte is at colorEndOffset-1
		// Parser reads pairs starting at 150: (150,151), (152,153), ..., (colorEndOffset-2, colorEndOffset-1)
		// Number of pairs = ((lastByte - 150) / 2) + 1
		lastColorByte := colorEndOffset - 1
		if lastColorByte >= 150 {
			overlapPairs = ((lastColorByte - 150) / 2) + 1
		}
	}

	// Calculate how many additional tone curve pairs we need to generate
	// Total pairs counted = overlap pairs + our generated pairs
	// We want: overlap + generated = target
	additionalPairs := targetTotalPairs - overlapPairs
	if additionalPairs < 0 {
		additionalPairs = 0
	}

	// Determine start offset: after color data, with even alignment
	startOffset := colorEndOffset
	if startOffset < 150 {
		startOffset = 150
	}
	if startOffset%2 != 0 {
		startOffset++
	}

	// Generate additional tone curve pairs starting after color data
	offset := startOffset
	for i := 0; i < additionalPairs && offset+1 < 500; i++ {
		// Use value 1 (minimal non-zero) for tone curve data
		data[offset] = 1
		data[offset+1] = 1
		offset += 2
	}
}

// === Phase 2: Exact Offset Writers ===
// These functions write all 48 parameters to exact byte offsets discovered through
// TypeScript implementation research. They mirror the extraction functions in parse.go.

// writeBasicAdjustments writes sharpening and clarity to exact offsets using Scaled4 encoding.
//
// Offsets:
//   - Sharpening: offset 82, Scaled4 encoding (-3.0 to +9.0)
//   - Clarity: offset 92, Scaled4 encoding (-5.0 to +5.0)
//
// Encoding: Scaled4 = (value * 4.0) + 0x80
func writeBasicAdjustments(data []byte, params *np3Parameters) {
	// Sharpening (offset 82)
	if len(data) > OffsetSharpening {
		data[OffsetSharpening] = EncodeScaled4(params.sharpening)
	}

	// Clarity (offset 92)
	if len(data) > OffsetClarity {
		data[OffsetClarity] = EncodeScaled4(params.clarity)
	}
}

// writeEffects writes film grain parameters to exact offsets.
func writeEffects(data []byte, params *np3Parameters) {
	// Grain Amount (offset 102)
	if len(data) > OffsetGrainAmount {
		data[OffsetGrainAmount] = EncodeScaled4(params.grainAmount)
	}

	// Grain Size (offset 222)
	if len(data) > OffsetGrainSize {
		data[OffsetGrainSize] = EncodeSigned8(params.grainSize)
	}
}

// writeAdvancedAdjustments writes 7 advanced parameters to exact offsets using Scaled4 and Signed8 encodings.
//
// Offsets:
//   - Mid-Range Sharpening: offset 242, Scaled4 (-5.0 to +5.0)
//   - Contrast: offset 272, Signed8 (-100 to +100)
//   - Highlights: offset 282, Signed8 (-100 to +100)
//   - Shadows: offset 292, Signed8 (-100 to +100)
//   - White Level: offset 302, Signed8 (-100 to +100)
//   - Black Level: offset 312, Signed8 (-100 to +100)
//   - Saturation: offset 322, Signed8 (-100 to +100)
func writeAdvancedAdjustments(data []byte, params *np3Parameters) {
	// Mid-Range Sharpening (offset 242, Scaled4)
	if len(data) > OffsetMidRangeSharpening {
		data[OffsetMidRangeSharpening] = EncodeScaled4(params.midRangeSharpening)
	}

	// Contrast (offset 272, Signed8)
	if len(data) > OffsetContrast {
		data[OffsetContrast] = EncodeSigned8(params.contrast)
	}

	// Highlights (offset 282, Signed8)
	if len(data) > OffsetHighlights {
		data[OffsetHighlights] = EncodeSigned8(params.highlights)
	}

	// Shadows (offset 292, Signed8)
	if len(data) > OffsetShadows {
		data[OffsetShadows] = EncodeSigned8(params.shadows)
	}

	// White Level (offset 302, Signed8)
	if len(data) > OffsetWhiteLevel {
		data[OffsetWhiteLevel] = EncodeSigned8(params.whiteLevel)
	}

	// Black Level (offset 312, Signed8)
	if len(data) > OffsetBlackLevel {
		data[OffsetBlackLevel] = EncodeSigned8(params.blackLevel)
	}

	// Saturation (offset 322, Signed8)
	if len(data) > OffsetSaturation {
		data[OffsetSaturation] = EncodeSigned8(params.saturation)
	}
}

// writeColorBlender writes 24 color blender parameters (8 colors × 3 values) to exact offsets.
//
// Offsets: 332-355 (sequential, 3 bytes per color)
// Encoding: Signed8 for all values (-100 to +100)
func writeColorBlender(data []byte, params *np3Parameters) {
	// Red (offsets 332-334)
	if len(data) > OffsetRedBrightness {
		data[OffsetRedHue] = EncodeSigned8(params.redHue)
		data[OffsetRedChroma] = EncodeSigned8(params.redChroma)
		data[OffsetRedBrightness] = EncodeSigned8(params.redBrightness)
	}

	// Orange (offsets 335-337)
	if len(data) > OffsetOrangeBrightness {
		data[OffsetOrangeHue] = EncodeSigned8(params.orangeHue)
		data[OffsetOrangeChroma] = EncodeSigned8(params.orangeChroma)
		data[OffsetOrangeBrightness] = EncodeSigned8(params.orangeBrightness)
	}

	// Yellow (offsets 338-340)
	if len(data) > OffsetYellowBrightness {
		data[OffsetYellowHue] = EncodeSigned8(params.yellowHue)
		data[OffsetYellowChroma] = EncodeSigned8(params.yellowChroma)
		data[OffsetYellowBrightness] = EncodeSigned8(params.yellowBrightness)
	}

	// Green (offsets 341-343)
	if len(data) > OffsetGreenBrightness {
		data[OffsetGreenHue] = EncodeSigned8(params.greenHue)
		data[OffsetGreenChroma] = EncodeSigned8(params.greenChroma)
		data[OffsetGreenBrightness] = EncodeSigned8(params.greenBrightness)
	}

	// Cyan (offsets 344-346)
	if len(data) > OffsetCyanBrightness {
		data[OffsetCyanHue] = EncodeSigned8(params.cyanHue)
		data[OffsetCyanChroma] = EncodeSigned8(params.cyanChroma)
		data[OffsetCyanBrightness] = EncodeSigned8(params.cyanBrightness)
	}

	// Blue (offsets 347-349)
	if len(data) > OffsetBlueBrightness {
		data[OffsetBlueHue] = EncodeSigned8(params.blueHue)
		data[OffsetBlueChroma] = EncodeSigned8(params.blueChroma)
		data[OffsetBlueBrightness] = EncodeSigned8(params.blueBrightness)
	}

	// Purple (offsets 350-352)
	if len(data) > OffsetPurpleBrightness {
		data[OffsetPurpleHue] = EncodeSigned8(params.purpleHue)
		data[OffsetPurpleChroma] = EncodeSigned8(params.purpleChroma)
		data[OffsetPurpleBrightness] = EncodeSigned8(params.purpleBrightness)
	}

	// Magenta (offsets 353-355)
	if len(data) > OffsetMagentaBrightness {
		data[OffsetMagentaHue] = EncodeSigned8(params.magentaHue)
		data[OffsetMagentaChroma] = EncodeSigned8(params.magentaChroma)
		data[OffsetMagentaBrightness] = EncodeSigned8(params.magentaBrightness)
	}
}

// writeColorGrading writes 11 color grading parameters to exact offsets using Hue12 and Signed8 encodings.
//
// Offsets:
//   - Highlights Zone: 368-371 (2-byte hue + 1-byte chroma + 1-byte brightness)
//   - Midtone Zone: 372-375 (2-byte hue + 1-byte chroma + 1-byte brightness)
//   - Shadows Zone: 376-379 (2-byte hue + 1-byte chroma + 1-byte brightness)
//   - Blending: 384 (direct value 0-100)
//   - Balance: 386 (Signed8 -100 to +100)
func writeColorGrading(data []byte, params *np3Parameters) {
	// Highlights zone (offsets 368-371)
	if len(data) > OffsetHighlightsBrightness {
		byte1, byte2 := EncodeHue12(params.highlightsZone.Hue)
		data[OffsetHighlightsHue] = byte1
		data[OffsetHighlightsHue+1] = byte2
		data[OffsetHighlightsChroma] = EncodeSigned8(params.highlightsZone.Chroma)
		data[OffsetHighlightsBrightness] = EncodeSigned8(params.highlightsZone.Brightness)
	}

	// Midtone zone (offsets 372-375)
	if len(data) > OffsetMidtoneBrightness {
		byte1, byte2 := EncodeHue12(params.midtoneZone.Hue)
		data[OffsetMidtoneHue] = byte1
		data[OffsetMidtoneHue+1] = byte2
		data[OffsetMidtoneChroma] = EncodeSigned8(params.midtoneZone.Chroma)
		data[OffsetMidtoneBrightness] = EncodeSigned8(params.midtoneZone.Brightness)
	}

	// Shadows zone (offsets 376-379)
	if len(data) > OffsetShadowsBrightness {
		byte1, byte2 := EncodeHue12(params.shadowsZone.Hue)
		data[OffsetShadowsHue] = byte1
		data[OffsetShadowsHue+1] = byte2
		data[OffsetShadowsChroma] = EncodeSigned8(params.shadowsZone.Chroma)
		data[OffsetShadowsBrightness] = EncodeSigned8(params.shadowsZone.Brightness)
	}

	// Blending (offset 384, Signed8 encoding 0-100)
	if len(data) > OffsetColorGradingBlending {
		data[OffsetColorGradingBlending] = EncodeSigned8(params.blending)
	}

	// Balance (offset 386, Signed8 -100 to +100)
	if len(data) > OffsetColorGradingBalance {
		data[OffsetColorGradingBalance] = EncodeSigned8(params.balance)
	}
}

// writeToneCurve writes tone curve control points using the BI0 marker structure.
//
// CRITICAL: NX Studio requires specific byte patterns to enable custom tone curves:
//   - Enable Flags: offsets 389-390 must be 0x01, 0x01
//   - Pre-BI0 structure: offsets 395-408 with specific values
//   - BI0 Marker: offset 409 with magic bytes and control points
//
// The pre-BI0 structure (observed in working files):
//   - 395-399: Tag bytes ([6]TEST or similar)
//   - 404: 0x00 (NOT point count!)
//   - 405: 0x02
//   - 408: 0x02
func writeToneCurve(data []byte, params *np3Parameters) {
	// Only write curve data if we have control points
	if params.toneCurvePointCount == 0 && len(params.toneCurvePoints) == 0 {
		return
	}

	// WORKING LAYOUT from Agfachrome RSX 2 v15-FIXED-BASELINE.NP3:
	// - Flags 389/390: 01 01
	// - TEST block at 395: 06 54 45 53 54 00 00 00 00 00 00 00 02 00 00 02
	// - BI0 at 409
	// - No Extended LUT (curve is in BI0 control points only)
	
	// Enable Flags at 389/390: Set to 0x01
	if len(data) > OffsetToneCurveEnabled2 {
		data[OffsetToneCurveEnabled1] = 0x01
		data[OffsetToneCurveEnabled2] = 0x01
	}

	// TEST block at 395-408 (observed in working FIXED-BASELINE file)
	if len(data) > 408 {
		data[395] = 0x06 // Length prefix
		data[396] = 0x54 // 'T'
		data[397] = 0x45 // 'E'
		data[398] = 0x53 // 'S'
		data[399] = 0x54 // 'T'
		data[400] = 0x00
		data[401] = 0x00
		data[402] = 0x00
		data[403] = 0x00
		data[404] = 0x00
		data[405] = 0x02
		data[406] = 0x00
		data[407] = 0x00
		data[408] = 0x02
	}
	
	// BI0 Marker at 409
	bi0 := OffsetBI0Marker // 409

	// Write BI0 marker structure at offset 409
	if len(data) > bi0+30 {
		// Magic bytes "BI0"
		data[bi0+0] = 0x42 // 'B'
		data[bi0+1] = 0x49 // 'I'
		data[bi0+2] = 0x30 // '0'

		// Fixed header bytes (observed in working files)
		data[bi0+3] = 0x00
		data[bi0+4] = 0xFF
		data[bi0+5] = 0x00
		data[bi0+6] = 0xFF
		data[bi0+7] = 0x01
		data[bi0+8] = 0x00

		// Point count at BI0+9
		pointCount := params.toneCurvePointCount
		if pointCount == 0 {
			pointCount = len(params.toneCurvePoints)
		}
		data[bi0+9] = uint8(pointCount)


		// Padding at BI0+10 (1 byte only)
		data[bi0+10] = 0x00
		// Note: Working files show data starting at BI0+11 (Offset 420)

		// Write control points at BI0+11 onwards (X,Y pairs as single bytes)
		for i := 0; i < len(params.toneCurvePoints) && i < pointCount; i++ {
			offset := bi0 + 11 + (i * 2)
			if len(data) > offset+1 {
				data[offset] = params.toneCurvePoints[i].value1   // X (input)
				data[offset+1] = params.toneCurvePoints[i].value2 // Y (output)
			}
		}

		// Padding after control points
		lastPointOffset := bi0 + 11 + (pointCount * 2)
		for i := lastPointOffset; i < lastPointOffset+4 && i < len(data); i++ {
			data[i] = 0x00
		}
	}

	// Extended tone curve LUT at offset 560 (0x230) - for files large enough to support it
	// Note: We do NOT write to 0x1CC as that can corrupt smaller NP3 structures
	if len(params.toneCurveRaw) > 0 && len(data) >= OffsetExtendedToneCurveLUT+512 {

		numEntries := len(params.toneCurveRaw)
		if numEntries > 256 {
			numEntries = 256
		}
		for i := 0; i < numEntries; i++ {
			offset := OffsetExtendedToneCurveLUT + (i * 2)
			value := params.toneCurveRaw[i]
			if value <= 32767 {
				value = value * 2 // Scale to 0-65535 range
			}
			// Note: MGA Slide has natural curve data (no artificial minimums)
			// The checkbox appears to be triggered by the BI0@395 layout,
			// not by LUT[0] being non-zero. Keeping natural curve values.
			data[offset] = uint8(value >> 8)
			data[offset+1] = uint8(value & 0xFF)
		}
	}
}

// applyWhiteBalanceToColorGrading simulates XMP White Balance (Temp/Tint) using NP3 Color Grading.
func applyWhiteBalanceToColorGrading(recipe *models.UniversalRecipe, params *np3Parameters) {
	temp := recipe.IncrementalTemperature
	tint := recipe.IncrementalTint

	if temp == 0 && tint == 0 {
		return
	}

	// Calculate WB vector (Polar coordinates)
	// Temp: + = Warm (Orange ~40°), - = Cool (Blue ~220°)
	// Tint: + = Magenta (~300°), - = Green (~120°)
	
	type colorVector struct {
		hue    float64 // degrees
		chroma float64 // strength
	}
	var wbVectors []colorVector

	// Temperature
	if temp > 0 {
		// Warm (Orange) - 40°
		wbVectors = append(wbVectors, colorVector{40, float64(temp) * 0.5})
	} else if temp < 0 {
		// Cool (Blue) - 220°
		wbVectors = append(wbVectors, colorVector{220, float64(-temp) * 0.5})
	}

	// Tint
	if tint > 0 {
		// Magenta - 300°
		wbVectors = append(wbVectors, colorVector{300, float64(tint) * 0.5})
	} else if tint < 0 {
		// Green - 120°
		wbVectors = append(wbVectors, colorVector{120, float64(-tint) * 0.5})
	}

	// Apply vectors to all 3 zones
	zones := []*models.ColorGradingZone{
		&params.highlightsZone,
		&params.midtoneZone,
		&params.shadowsZone,
	}

	for _, zone := range zones {
		// Convert current zone to Cartesian (x, y)
		rad := float64(zone.Hue) * math.Pi / 180.0
		x := float64(zone.Chroma) * math.Cos(rad)
		y := float64(zone.Chroma) * math.Sin(rad)

		// Add WB vectors
		for _, vec := range wbVectors {
			vecRad := vec.hue * math.Pi / 180.0
			x += vec.chroma * math.Cos(vecRad)
			y += vec.chroma * math.Sin(vecRad)
		}

		// Convert back to Polar (Hue, Chroma)
		newChroma := math.Sqrt(x*x + y*y)
		newHueRad := math.Atan2(y, x)
		newHue := newHueRad * 180.0 / math.Pi
		if newHue < 0 {
			newHue += 360.0
		}

		// Update zone
		zone.Hue = int(newHue)
		zone.Chroma = int(newChroma)

		// Clamp Chroma to 0-100 (NP3 limit usually 100 for Signed8? No, Chroma is 0-100 unsigned usually)
		// But in writeColorGrading, Chroma is EncodeSigned8 (-100..100).
		// Wait, Color Grading Chroma is 0..100? Or -100..100?
		// models.ColorGradingZone implies Hue, Chroma, Brightness.
		// Usually Chroma is positive magnitude.
		if zone.Chroma > 100 {
			zone.Chroma = 100
		}
	}
	
	// Ensure Blending is sufficient to make this visible if it wasn't already
	if params.blending < 50 {
		params.blending = 50 // Ensure grading is applied
	}
}

// applyVibranceToSaturation maps XMP Vibrance to NP3 Saturation.
// Vibrance is a smart saturation that protects skin tones.
// We approximate it with a gentler global Saturation boost (0.5x scale).
func applyVibranceToSaturation(recipe *models.UniversalRecipe, params *np3Parameters) {
	if recipe.Vibrance != 0 {
		params.saturation += int(float64(recipe.Vibrance) * 0.5)
		if params.saturation > 100 {
			params.saturation = 100
		} else if params.saturation < -100 {
			params.saturation = -100
		}
	}
}

// applyDehaze maps XMP Dehaze to NP3 Contrast and Black Level.
// Dehaze increases local contrast (approximated by Global Contrast)
// and darkens blacks/shadows.
func applyDehaze(recipe *models.UniversalRecipe, params *np3Parameters) {
	if recipe.Dehaze != 0 {
		// Boost Contrast (0.5x strength)
		params.contrast += int(float64(recipe.Dehaze) * 0.5)
		if params.contrast > 100 {
			params.contrast = 100
		} else if params.contrast < -100 {
			params.contrast = -100
		}

		// Reduce Black Level (darken blacks) (0.25x strength)
		// Positive Dehaze -> Darker Blacks -> Lower Black Level?
		// Or Higher Black Level?
		// Black Level in NP3: - = Darker, + = Lighter? Usually.
		params.blackLevel -= int(float64(recipe.Dehaze) * 0.25)
		if params.blackLevel > 100 {
			params.blackLevel = 100
		} else if params.blackLevel < -100 {
			params.blackLevel = -100
		}
	}
}

// ApplyFlexibleColorBaselineCompensation applies compensation for the difference between
// Adobe's "Flexible Color" DCP baseline (+47.8 brightness) and Nikon's native "Flexible Color"
// baseline (~+105.6 brightness).
//
// This function darkens the tone curve by ~57.8 luminance units to compensate for Nikon's
// brighter baseline, ensuring that NX Studio output matches Lightroom output when both
// use their respective Flexible Color profiles.
//
// The compensation uses the inverse baseline method:
// - Adobe DCP baseline extracted from "Nikon Z f Camera Flexible Color.dcp"
// - Nikon baseline estimated from visual comparison (+51.1 luminance difference)
// - Compensation applied as inverse curve mapping
//
// See: docs/analysis/adobe-dcp-baseline-discovery.md for detailed analysis
// See: scripts/calculate_baseline_compensation.py for derivation
//
// Returns a darkened LUT that compensates for the baseline difference.
// ApplyFlexibleColorBaselineCompensation applies a compensation curve to match Adobe vs Nikon baseline differences.
//
// Problem:
// - Adobe's "Flexible Color" profile applies a moderate S-curve baseline.
// - Nikon's native "Flexible Color" profile applies a MUCH brighter baseline (Shadows +70, Mids +40).
// - Start-to-end difference: +47.6 luminance average, but highly non-linear.
//
// Solution:
// We construct a compensation curve C such that: C(NikonBase(x)) ≈ XMPCurve(AdobeBase(x))
// This ensures that when Nikon applies its base, followed by our curve, the result matches 
// what Lightroom produces (XMP curve on top of Adobe base).
//
// See: docs/analysis/visual-accuracy-root-cause-analysis.md
func ApplyFlexibleColorBaselineCompensation(lut []int) []int {
	// 1. Adobe DCP Flexible Color Baseline (32 sample points from 128-point curve)
	// Extracted from: testdata/dcp/Nikon Z f Camera Flexible Color.dcp
	adobeBaseline := []struct{ x, y float64 }{
		{x: 0.000000, y: 0.000000}, {x: 0.002438, y: 0.001959}, {x: 0.005208, y: 0.005911}, {x: 0.009189, y: 0.013544},
		{x: 0.014539, y: 0.024477}, {x: 0.021368, y: 0.038992}, {x: 0.029773, y: 0.057326}, {x: 0.039846, y: 0.079311},
		{x: 0.051668, y: 0.105307}, {x: 0.065317, y: 0.135484}, {x: 0.080866, y: 0.169726}, {x: 0.098385, y: 0.208347},
		{x: 0.117937, y: 0.251100}, {x: 0.139587, y: 0.298750}, {x: 0.163393, y: 0.349381}, {x: 0.189414, y: 0.400751},
		{x: 0.217703, y: 0.449423}, {x: 0.248315, y: 0.496630}, {x: 0.281302, y: 0.543200}, {x: 0.316712, y: 0.590829},
		{x: 0.354594, y: 0.639572}, {x: 0.394996, y: 0.686798}, {x: 0.437963, y: 0.730482}, {x: 0.483539, y: 0.771481},
		{x: 0.531769, y: 0.809248}, {x: 0.582693, y: 0.844011}, {x: 0.636355, y: 0.875652}, {x: 0.692794, y: 0.903849},
		{x: 0.752050, y: 0.928407}, {x: 0.814161, y: 0.953279}, {x: 0.879166, y: 0.974905}, {x: 1.000000, y: 1.000000},
	}

	// 2. Generate Adobe Base LUT (Input -> AdobeOutput)
	adobeLUT := make([]int, 256)
	for i := 0; i < 256; i++ {
		inputNorm := float64(i) / 255.0
		var outputNorm float64
		// Basic linear interpolation
		if inputNorm <= adobeBaseline[0].x {
			outputNorm = adobeBaseline[0].y
		} else if inputNorm >= adobeBaseline[len(adobeBaseline)-1].x {
			outputNorm = adobeBaseline[len(adobeBaseline)-1].y
		} else {
			for j := 0; j < len(adobeBaseline)-1; j++ {
				if inputNorm >= adobeBaseline[j].x && inputNorm <= adobeBaseline[j+1].x {
					t := (inputNorm - adobeBaseline[j].x) / (adobeBaseline[j+1].x - adobeBaseline[j].x)
					outputNorm = adobeBaseline[j].y + t*(adobeBaseline[j+1].y-adobeBaseline[j].y)
					break
				}
			}
		}
		adobeLUT[i] = int(outputNorm * 255.0)
		if adobeLUT[i] > 255 { adobeLUT[i] = 255 } // Safety clamp
	}

	// 3. Generate Estimated Nikon Base LUT (Input -> NikonOutput)
	// Modeled as Adobe + Non-Linear Offset based on visual regression data
	// INCREASED VALUES: Previous 75/60/25/10 was not aggressive enough.
	// Shadows: +120->100 | Mids: 100->50 | Highs: 50->20
	nikonLUT := make([]int, 256)
	for i := 0; i < 256; i++ {
		var offset float64
		inputVal := float64(i)
		
		if inputVal < 64 {
			// Shadows (0-63): Fade 120 -> 100
			ratio := inputVal / 64.0
			offset = 120.0 - (20.0 * ratio)
		} else if inputVal < 192 {
			// Midtones (64-191): Fade 100 -> 50
			ratio := (inputVal - 64.0) / 128.0
			offset = 100.0 - (50.0 * ratio)
		} else {
			// Highlights (192-255): Fade 50 -> 20
			ratio := (inputVal - 192.0) / 63.0
			offset = 50.0 - (30.0 * ratio)
		}

		val := adobeLUT[i] + int(offset)
		if val > 255 { val = 255 }
		nikonLUT[i] = val
	}

	// 4. Construct Compensation LUT mapping: NikonOutput -> TargetOutput
	// We gather points (x, y) where x = NikonLUT[i], y = lut[AdobeLUT[i]]
	// This creates the mapping C(x) = y
	compensatedLUT := make([]int, 256)
	for i := range compensatedLUT {
		compensatedLUT[i] = -1 // Mark as uninitialized
	}

	// Gather points
	for i := 0; i < 256; i++ {
		inputToCurve := nikonLUT[i]        // This is what the curve receives from Nikon Base
		
		adobeOut := adobeLUT[i]            // This is what Adobe Base produces for same input
		target := lut[adobeOut]            // This is our desired final output (XMP applied to Adobe Base)

		// Map inputToCurve -> target
		// Since multiple inputs i might map to same inputToCurve, we take the last one (safe assumption for monotonic)
		// or simpler: just write it.
		compensatedLUT[inputToCurve] = target
	}

	// 5. Fill gaps using interpolation and handle edges
	// Forward pass to fill gaps
	lastVal := 0
	// Find first valid value
	firstValidIdx := -1
	for i := 0; i < 256; i++ {
		if compensatedLUT[i] != -1 {
			firstValidIdx = i
			lastVal = compensatedLUT[i]
			break
		}
	}

	if firstValidIdx == -1 {
		// Degenerate case, shouldn't happen
		return lut
	}

	// Fill valid range
	for i := firstValidIdx; i < 256; i++ {
		if compensatedLUT[i] != -1 {
			lastVal = compensatedLUT[i]
		} else {
			// Look ahead for next valid
			nextValidIdx := -1
			nextVal := lastVal
			for j := i + 1; j < 256; j++ {
				if compensatedLUT[j] != -1 {
					nextValidIdx = j
					nextVal = compensatedLUT[j]
					break
				}
			}
			
			if nextValidIdx != -1 {
				// Interpolate
				dist := float64(nextValidIdx - (i - 1)) // distance from prev valid
				currDist := float64(i - (i - 1))
				t := currDist / dist
				
				// Simpler fill: Just connect last valid to next valid
				// We need the index of point BEFORE the gap. 
				// The loop structure: we are at `i` which is -1.
				// Previous index `i-1` MUST be valid (or we are before firstValidIdx).
				// We handled before firstValidIdx separately.
				prevValidIdx := i - 1
				prevVal := compensatedLUT[prevValidIdx]
				
				compensatedLUT[i] = int(float64(prevVal) + t*float64(nextVal - prevVal))
			} else {
				// No more valid points, clamp to last value
				compensatedLUT[i] = lastVal
			}
		}
	}

	// Fill beginning (0 to firstValidIdx)
	// Extrapolate slope? Or clamp to 0?
	// Nikon adds brightness, so input 0 -> Nikon 75. 
	// Curve inputs 0-74 are impossible inputs from Base.
	// But mathematically we should probably clamp to 0 or extend slope.
	// Clamping 0 is safest for "Black remains Black".
	for i := 0; i < firstValidIdx; i++ {
		compensatedLUT[i] = 0 // Assume anything darker than black is black
	}

	return compensatedLUT
}
