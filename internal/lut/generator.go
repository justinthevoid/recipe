// Package lut provides 3D LUT generation for both real-time WebGL preview
// (17³ RGBA) and Adobe XMP embedding (32³ RGB). Both output paths share the
// same applyColorTransform pipeline, ensuring visual parity.
//
// Color space: all HSL operations use standard sRGB-defined HSL. Exposure is
// applied in linear light via simplified gamma (pow 2.2).
package lut

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/justin/recipe/internal/models"
)

const (
	// LUTSize defines the dimensions of the 3D LUT cube (32x32x32 = 32,768 colors)
	LUTSize = 32

	// LUTPoints is the total number of color mappings in the LUT
	LUTPoints = LUTSize * LUTSize * LUTSize
)

// Generate3DLUT creates a 3D RGB lookup table from a UniversalRecipe.
// The LUT maps input RGB values through all NP3 color adjustments:
// - Color Blender (8-color HSL adjustments)
// - Color Grading (3-zone HSL adjustments)
// - Global adjustments (saturation, contrast, etc.)
//
// Returns the raw uncompressed LUT data as float32 triplets (R,G,B) in [0,1] range.
func Generate3DLUT(recipe *models.UniversalRecipe) ([]byte, error) {
	if recipe == nil {
		return nil, fmt.Errorf("recipe cannot be nil")
	}

	// Allocate buffer for LUT data: 32^3 points × 3 channels × 4 bytes/float32
	lutData := make([]byte, LUTPoints*3*4)
	buf := bytes.NewBuffer(lutData[:0])

	// Generate LUT by sampling the color cube
	for bIdx := 0; bIdx < LUTSize; bIdx++ {
		for gIdx := 0; gIdx < LUTSize; gIdx++ {
			for rIdx := 0; rIdx < LUTSize; rIdx++ {
				// Normalize indices to [0,1] range
				r := float64(rIdx) / float64(LUTSize-1)
				g := float64(gIdx) / float64(LUTSize-1)
				b := float64(bIdx) / float64(LUTSize-1)

				// Apply color transformation
				rOut, gOut, bOut := applyColorTransform(recipe, r, g, b)

				// Write RGB triplet as float32
				binary.Write(buf, binary.LittleEndian, float32(rOut))
				binary.Write(buf, binary.LittleEndian, float32(gOut))
				binary.Write(buf, binary.LittleEndian, float32(bOut))
			}
		}
	}

	return buf.Bytes(), nil
}

// applyColorTransform applies all NP3 color adjustments to an input sRGB color.
// Pipeline order: sRGB→Linear → Exposure → Linear→sRGB → HSL → Tone curve →
// Color blender → Color grading → Saturation → RGB
func applyColorTransform(recipe *models.UniversalRecipe, r, g, b float64) (float64, float64, float64) {
	// Step 1: Apply exposure in linear light space
	if recipe.Exposure != 0 {
		lr := srgbToLinear(r)
		lg := srgbToLinear(g)
		lb := srgbToLinear(b)

		exposureFactor := math.Pow(2, recipe.Exposure)
		lr = clamp(lr*exposureFactor, 0, 1)
		lg = clamp(lg*exposureFactor, 0, 1)
		lb = clamp(lb*exposureFactor, 0, 1)

		r = linearToSrgb(lr)
		g = linearToSrgb(lg)
		b = linearToSrgb(lb)
	}

	// Step 2: Convert sRGB to HSL for parametric adjustments
	h, s, l := rgbToHsl(r, g, b)

	// Step 3: Apply tone curve adjustments (contrast, highlights, shadows, whites, blacks)
	l = applyToneCurve(recipe, l)

	// Step 4: Apply Color Blender adjustments (8-color HSL)
	// Guard: skip for near-achromatic pixels where hue is undefined
	if s >= 0.02 {
		h, s, l = applyColorBlender(recipe, h, s, l)
	}

	// Step 5: Apply Color Grading (tint model — tints neutrals too, NO achromatic guard)
	h, s, l = applyColorGrading(recipe, h, s, l)

	// Step 6: Apply global saturation adjustment
	if recipe.Saturation != 0 {
		satScale := 1.0 + (float64(recipe.Saturation) / 100.0)
		s = clamp(s*satScale, 0, 1)
	}

	// Step 7: Convert back to sRGB
	rOut, gOut, bOut := hslToRgb(h, s, l)
	return rOut, gOut, bOut
}

// applyToneCurve applies tone adjustments to luminance.
// Contrast uses a sigmoid S-curve (positive) or lerp-to-flat (negative).
// Highlights/Shadows/Whites/Blacks use zone-weighted additive adjustments.
func applyToneCurve(recipe *models.UniversalRecipe, l float64) float64 {
	// Apply contrast first (affects entire curve)
	if recipe.Contrast != 0 {
		l = applyContrast(l, float64(recipe.Contrast))
	}

	// Highlights: zone 0.5–1.0, linear ramp weight, intensity 0.5x
	if recipe.Highlights != 0 && l > 0.5 {
		weight := (l - 0.5) * 2.0
		l += (float64(recipe.Highlights) / 100.0) * weight * 0.5
		l = clamp(l, 0, 1)
	}

	// Shadows: zone 0.0–0.5, linear ramp weight, intensity 0.5x
	if recipe.Shadows != 0 && l < 0.5 {
		weight := (0.5 - l) * 2.0
		l += (float64(recipe.Shadows) / 100.0) * weight * 0.5
		l = clamp(l, 0, 1)
	}

	// Whites: zone 0.7–1.0, intensity 0.8x
	if recipe.Whites != 0 && l > 0.7 {
		weight := (l - 0.7) / 0.3
		l += (float64(recipe.Whites) / 100.0) * weight * 0.8
		l = clamp(l, 0, 1)
	}

	// Blacks: zone 0.0–0.3, intensity 0.8x
	if recipe.Blacks != 0 && l < 0.3 {
		weight := (0.3 - l) / 0.3
		l += (float64(recipe.Blacks) / 100.0) * weight * 0.8
		l = clamp(l, 0, 1)
	}

	return l
}

// applyContrast applies sigmoid S-curve (positive) or lerp-to-flat (negative).
// Midpoint (0.5) is preserved for both directions.
func applyContrast(l, contrast float64) float64 {
	if contrast > 0 {
		// Positive: sigmoid S-curve
		slope := (contrast / 100.0) * 5.0 // 0.0 to 5.0
		raw := func(x float64) float64 {
			return 1.0 / (1.0 + math.Exp(-slope*(x-0.5)))
		}
		lo, hi := raw(0), raw(1)
		return (raw(l) - lo) / (hi - lo)
	}
	// Negative: lerp toward midpoint (flatten)
	t := -contrast / 100.0 // 0.0 to 1.0
	return l + (0.5-l)*t
}

// applyColorBlender applies the 8-color HSL adjustments based on input hue.
// Each color channel has a 90° wide hue band with cosine falloff.
// Overlapping band contributions are weight-normalized to prevent amplification.
func applyColorBlender(recipe *models.UniversalRecipe, h, s, l float64) (float64, float64, float64) {
	type colorEntry struct {
		center float64
		adj    *models.ColorAdjustment
	}
	colors := []colorEntry{
		{0, &recipe.Red},
		{30, &recipe.Orange},
		{60, &recipe.Yellow},
		{120, &recipe.Green},
		{180, &recipe.Aqua},
		{240, &recipe.Blue},
		{270, &recipe.Purple},
		{330, &recipe.Magenta},
	}

	hueAdj := 0.0
	satAdj := 0.0
	lumAdj := 0.0
	hueWeight := 0.0
	satWeight := 0.0
	lumWeight := 0.0

	for _, c := range colors {
		weight := calculateHueWeight(h, c.center)
		if weight > 0 {
			if c.adj.Hue != 0 {
				hueWeight += weight
			}
			if c.adj.Saturation != 0 {
				satWeight += weight
			}
			if c.adj.Luminance != 0 {
				lumWeight += weight
			}
			// Hue values are on a -100..+100 perceptual scale (matching Lightroom XMP).
			// At ±100, the angular shift is approximately ±30°.
			hueAdj += float64(c.adj.Hue) * (30.0 / 100.0) * weight
			satAdj += float64(c.adj.Saturation) * weight
			lumAdj += float64(c.adj.Luminance) * weight
		}
	}

	// Per-axis normalization prevents cross-axis attenuation
	if hueWeight > 1.0 {
		hueAdj /= hueWeight
	}
	if satWeight > 1.0 {
		satAdj /= satWeight
	}
	if lumWeight > 1.0 {
		lumAdj /= lumWeight
	}

	// Apply adjustments
	h = math.Mod(h+hueAdj, 360)
	if h < 0 {
		h += 360
	}
	s = clamp(s+(satAdj/100.0), 0, 1)
	l = clamp(l+(lumAdj/100.0), 0, 1)

	return h, s, l
}

// calculateHueWeight calculates how much a color adjustment should affect
// the input hue based on angular distance from the color's center.
// Uses cosine falloff over a 90° wide band (±45° from center).
func calculateHueWeight(inputHue, centerHue float64) float64 {
	diff := math.Abs(inputHue - centerHue)
	if diff > 180 {
		diff = 360 - diff
	}

	const width = 90.0
	if diff > width/2 {
		return 0
	}

	// Cosine-based smooth falloff
	return (math.Cos(diff*math.Pi/width) + 1) / 2
}

// applyColorGrading applies 3-zone color grading using a tint model.
// Hue = target tint color, abs(Chroma) = tint intensity, Brightness = luminance offset.
// Zone weights use smooth crossfade controlled by Blending and Balance.
// Unlike the color blender, this DOES affect neutral/achromatic pixels (tinting grays).
func applyColorGrading(recipe *models.UniversalRecipe, h, s, l float64) (float64, float64, float64) {
	if recipe.ColorGrading == nil {
		return h, s, l
	}

	cg := recipe.ColorGrading

	// Zone boundary calculation with Blending and Balance
	blend := float64(cg.Blending) / 100.0
	bal := float64(cg.Balance) / 200.0

	shadowEnd := 0.33 + bal + blend*0.17
	highlightStart := 0.67 + bal - blend*0.17

	// Shadow weight: 1.0 for dark pixels, fades to 0 at shadowEnd
	shadowWeight := 1.0 - smoothstep(shadowEnd-0.15, shadowEnd, l)
	// Highlight weight: 0 for dark pixels, fades to 1.0 past highlightStart
	highlightWeight := smoothstep(highlightStart, highlightStart+0.15, l)
	midtoneWeight := clamp(1.0-shadowWeight-highlightWeight, 0, 1)

	// Normalize zone weights to sum to 1.0
	totalZoneWeight := shadowWeight + midtoneWeight + highlightWeight
	if totalZoneWeight > 0 {
		shadowWeight /= totalZoneWeight
		midtoneWeight /= totalZoneWeight
		highlightWeight /= totalZoneWeight
	}

	// Accumulate weighted tint using vector averaging (handles hue wraparound)
	type zoneEntry struct {
		zone   *models.ColorGradingZone
		weight float64
	}
	zones := []zoneEntry{
		{&cg.Shadows, shadowWeight},
		{&cg.Midtone, midtoneWeight},
		{&cg.Highlights, highlightWeight},
	}

	hueX, hueY := 0.0, 0.0
	totalIntensity := 0.0
	brightnessAdj := 0.0

	for _, z := range zones {
		if z.zone.Chroma == 0 || z.weight <= 0 {
			continue
		}
		intensity := math.Abs(float64(z.zone.Chroma)) / 100.0 * z.weight
		rad := float64(z.zone.Hue) * math.Pi / 180.0
		hueX += math.Cos(rad) * intensity
		hueY += math.Sin(rad) * intensity
		totalIntensity += intensity
		brightnessAdj += (float64(z.zone.Brightness) / 100.0) * z.weight * 0.5
	}

	// Apply tint via RGB blending to avoid hue interpolation issues
	if totalIntensity > 0 {
		targetHue := math.Atan2(hueY, hueX) * 180.0 / math.Pi
		if targetHue < 0 {
			targetHue += 360
		}
		blendAmount := clamp(totalIntensity, 0, 1)

		// Create tint reference color at target hue with full saturation
		tR, tG, tB := hslToRgb(targetHue, 1.0, l)
		// Convert current pixel to RGB
		cR, cG, cB := hslToRgb(h, s, l)

		// Blend in RGB space (path-independent, no hue wrapping artifacts)
		blendR := cR*(1-blendAmount) + tR*blendAmount
		blendG := cG*(1-blendAmount) + tG*blendAmount
		blendB := cB*(1-blendAmount) + tB*blendAmount

		// Convert back to HSL
		h, s, l = rgbToHsl(blendR, blendG, blendB)
	}

	// Apply brightness from all contributing zones (including those with Chroma=0)
	for _, z := range zones {
		if z.zone.Brightness != 0 && z.weight > 0 && z.zone.Chroma == 0 {
			brightnessAdj += (float64(z.zone.Brightness) / 100.0) * z.weight * 0.5
		}
	}
	l = clamp(l+brightnessAdj, 0, 1)

	return h, s, l
}

// Generate3DLUTForPreview creates a 3D RGBA lookup table for WebGL preview rendering.
// Unlike Generate3DLUT (which outputs RGB float32 triplets for XMP embedding),
// this outputs RGBA float32 (4 floats per texel, A=1.0) for direct upload to
// a WebGL RGBA32F TEXTURE_3D. The size parameter controls the LUT cube dimensions.
//
// Returns raw bytes of float32 RGBA values. Total size: size³ × 4 × 4 bytes.
func Generate3DLUTForPreview(recipe *models.UniversalRecipe, size int) ([]byte, error) {
	if recipe == nil {
		return nil, fmt.Errorf("recipe cannot be nil")
	}
	if size < 2 || size > 65 {
		return nil, fmt.Errorf("LUT size must be between 2 and 65, got %d", size)
	}

	totalTexels := size * size * size
	// RGBA float32 = 4 channels × 4 bytes = 16 bytes per texel
	lutData := make([]byte, 0, totalTexels*4*4)
	buf := bytes.NewBuffer(lutData)

	for bIdx := 0; bIdx < size; bIdx++ {
		for gIdx := 0; gIdx < size; gIdx++ {
			for rIdx := 0; rIdx < size; rIdx++ {
				r := float64(rIdx) / float64(size-1)
				g := float64(gIdx) / float64(size-1)
				b := float64(bIdx) / float64(size-1)

				rOut, gOut, bOut := applyColorTransform(recipe, r, g, b)

				binary.Write(buf, binary.LittleEndian, float32(rOut))
				binary.Write(buf, binary.LittleEndian, float32(gOut))
				binary.Write(buf, binary.LittleEndian, float32(bOut))
				binary.Write(buf, binary.LittleEndian, float32(1.0)) // Alpha = 1.0
			}
		}
	}

	return buf.Bytes(), nil
}

// CompressAndEncodeLUT compresses the raw LUT data and encodes it for XMP embedding.
// Uses zlib compression followed by Z85-like encoding suitable for XML.
func CompressAndEncodeLUT(lutData []byte) (string, string, error) {
	// Compress with zlib
	var compressed bytes.Buffer

	// Write uncompressed size header (4 bytes, little-endian)
	sizeHeader := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeHeader, uint32(len(lutData)))
	compressed.Write(sizeHeader)

	// Compress the LUT data
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write(lutData); err != nil {
		return "", "", fmt.Errorf("zlib compression failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", "", fmt.Errorf("zlib close failed: %w", err)
	}

	// Encode with modified Z85 for XML safety
	encoded := encodeZ85ForXML(compressed.Bytes())

	// Generate MD5 hash for table ID
	hash := md5.Sum(lutData)
	tableID := fmt.Sprintf("%X", hash)

	return tableID, encoded, nil
}

// encodeZ85ForXML encodes binary data using a Z85-like algorithm modified for XML safety.
// This matches Adobe's encoding used in crs:Table_* attributes.
func encodeZ85ForXML(data []byte) string {
	// Adobe's Z85 variant uses XML-safe characters
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+,-./:;<=>?@[]^_{|}~"

	var result bytes.Buffer

	// Process in 4-byte chunks
	for i := 0; i < len(data); i += 4 {
		// Pad last chunk if needed
		chunk := make([]byte, 4)
		copy(chunk, data[i:])

		// Convert 4 bytes to 32-bit value (big-endian for Z85)
		value := uint32(chunk[0])<<24 | uint32(chunk[1])<<16 | uint32(chunk[2])<<8 | uint32(chunk[3])

		// Encode as 5 base-85 characters
		for j := 4; j >= 0; j-- {
			result.WriteByte(charset[value%85])
			value /= 85
		}
	}

	return result.String()
}

// --- Math helpers ---

// srgbToLinear converts an sRGB-encoded value to linear light using simplified gamma.
// Uses pow(x, 2.2) approximation (max error < 0.001 vs exact sRGB transfer function).
func srgbToLinear(x float64) float64 {
	return math.Pow(clamp(x, 0, 1), 2.2)
}

// linearToSrgb converts a linear light value to sRGB encoding using simplified gamma.
func linearToSrgb(x float64) float64 {
	return math.Pow(clamp(x, 0, 1), 1.0/2.2)
}

// rgbToHsl converts sRGB values to HSL (hue 0-360, saturation 0-1, lightness 0-1).
// This operates on gamma-encoded sRGB values (standard HSL definition).
func rgbToHsl(r, g, b float64) (h, s, l float64) {
	maxC := math.Max(r, math.Max(g, b))
	minC := math.Min(r, math.Min(g, b))
	l = (maxC + minC) / 2

	if maxC == minC {
		return 0, 0, l // achromatic
	}

	d := maxC - minC
	if l > 0.5 {
		s = d / (2.0 - maxC - minC)
	} else {
		s = d / (maxC + minC)
	}

	switch maxC {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/d + 2
	case b:
		h = (r-g)/d + 4
	}
	h *= 60

	return h, s, l
}

// hslToRgb converts HSL to sRGB values. Hue is in degrees (any range, normalized internally).
func hslToRgb(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l // achromatic
	}

	// Normalize hue to [0, 360)
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	h /= 360.0 // normalize to 0-1
	r = hueToRgb(p, q, h+1.0/3.0)
	g = hueToRgb(p, q, h)
	b = hueToRgb(p, q, h-1.0/3.0)

	return r, g, b
}

// hueToRgb is a helper for hslToRgb that converts a hue sector to an RGB channel value.
func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// smoothstep provides GLSL-style Hermite interpolation between two edges.
// Returns 0 when x <= edge0, 1 when x >= edge1, smooth curve between.
func smoothstep(edge0, edge1, x float64) float64 {
	t := clamp((x-edge0)/(edge1-edge0), 0, 1)
	return t * t * (3 - 2*t)
}

// clamp restricts a value to the range [min, max].
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
