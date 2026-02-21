// Package lut provides 3D LUT generation for Adobe XMP profiles.
// This enables high-accuracy color transformation from NP3 parametric adjustments
// to RGB lookup tables that match Adobe's processing pipeline.
package lut

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/justin/recipe/internal/models"
	"github.com/lucasb-eyer/go-colorful"
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

// applyColorTransform applies all NP3 color adjustments to an input RGB color.
// This mimics how NX Studio processes colors through the Picture Control pipeline.
func applyColorTransform(recipe *models.UniversalRecipe, r, g, b float64) (float64, float64, float64) {
	// Step 1: Apply exposure adjustment first (affects overall brightness)
	if recipe.Exposure != 0 {
		// Exposure: -5.0 to +5.0 maps to 0.03125x to 32x (2^-5 to 2^5)
		exposureFactor := math.Pow(2, recipe.Exposure)
		r = clamp(r*exposureFactor, 0, 1)
		g = clamp(g*exposureFactor, 0, 1)
		b = clamp(b*exposureFactor, 0, 1)
	}

	// Step 2: Convert RGB to HSL for parametric adjustments
	color := colorful.Color{R: r, G: g, B: b}
	h, s, l := color.Hsl()

	// Step 3: Apply tone curve adjustments (contrast, highlights, shadows, whites, blacks)
	l = applyToneCurve(recipe, l)

	// Step 4: Apply Color Blender adjustments (8-color HSL)
	// IMPORTANT: Only apply color-specific adjustments if saturation is meaningful
	// For neutral/near-neutral colors (s < 0.01), hue is undefined and should be ignored
	if s >= 0.01 {
		h, s, l = applyColorBlender(recipe, h, s, l)
	}

	// Step 5: Apply Color Grading zone adjustments (3-zone HSL)
	// Only apply hue shifts for colors with meaningful saturation
	if s >= 0.01 {
		h, s, l = applyColorGrading(recipe, h, s, l)
	} else {
		// For neutral colors, only apply brightness adjustments from color grading
		_, _, l = applyColorGrading(recipe, h, 0, l)
	}

	// Step 6: Apply global saturation adjustment
	if recipe.Saturation != 0 {
		// Saturation adjustment: -100 to +100 maps to 0% to 200%
		satScale := 1.0 + (float64(recipe.Saturation) / 100.0)
		s = clamp(s*satScale, 0, 1)
	}

	// Step 7: Convert back to RGB
	outColor := colorful.Hsl(h, s, l)
	return outColor.R, outColor.G, outColor.B
}

// applyToneCurve applies tone adjustments (contrast, highlights, shadows, whites, blacks)
// to the luminance value. This mimics Adobe's tone curve system.
func applyToneCurve(recipe *models.UniversalRecipe, l float64) float64 {
	// Apply contrast first (affects entire curve)
	if recipe.Contrast != 0 {
		// Contrast: -100 to +100
		// Pivots around midpoint (0.5), steepens or flattens the curve
		contrastFactor := 1.0 + (float64(recipe.Contrast) / 100.0)
		l = ((l - 0.5) * contrastFactor) + 0.5
		l = clamp(l, 0, 1)
	}

	// Apply highlights adjustment (affects upper range)
	if recipe.Highlights != 0 {
		// Highlights: -100 to +100
		// Affects luminance > 0.5, with stronger effect toward 1.0
		if l > 0.5 {
			weight := (l - 0.5) * 2.0 // 0.0 at l=0.5, 1.0 at l=1.0
			adjustment := float64(recipe.Highlights) / 100.0
			l += adjustment * weight * 0.5
			l = clamp(l, 0, 1)
		}
	}

	// Apply shadows adjustment (affects lower range)
	if recipe.Shadows != 0 {
		// Shadows: -100 to +100
		// Affects luminance < 0.5, with stronger effect toward 0.0
		if l < 0.5 {
			weight := (0.5 - l) * 2.0 // 1.0 at l=0.0, 0.0 at l=0.5
			adjustment := float64(recipe.Shadows) / 100.0
			l += adjustment * weight * 0.5
			l = clamp(l, 0, 1)
		}
	}

	// Apply whites adjustment (affects extreme highlights)
	if recipe.Whites != 0 {
		// Whites: -100 to +100
		// Affects luminance > 0.75, with strongest effect toward 1.0
		if l > 0.75 {
			weight := (l - 0.75) * 4.0 // 0.0 at l=0.75, 1.0 at l=1.0
			adjustment := float64(recipe.Whites) / 100.0
			l += adjustment * weight * 0.3
			l = clamp(l, 0, 1)
		}
	}

	// Apply blacks adjustment (affects extreme shadows)
	if recipe.Blacks != 0 {
		// Blacks: -100 to +100
		// Affects luminance < 0.25, with strongest effect toward 0.0
		if l < 0.25 {
			weight := (0.25 - l) * 4.0 // 1.0 at l=0.0, 0.0 at l=0.25
			adjustment := float64(recipe.Blacks) / 100.0
			l += adjustment * weight * 0.3
			l = clamp(l, 0, 1)
		}
	}

	return l
}

// applyColorBlender applies the 8-color HSL adjustments based on input hue.
// Each color channel affects a specific hue range with gradual falloff.
func applyColorBlender(recipe *models.UniversalRecipe, h, s, l float64) (float64, float64, float64) {
	// Define hue centers for each of the 8 colors (in degrees)
	colorCenters := map[string]float64{
		"red":     0,    // 0°
		"orange":  30,   // 30°
		"yellow":  60,   // 60°
		"green":   120,  // 120°
		"aqua":    180,  // 180°
		"blue":    240,  // 240°
		"purple":  270,  // 270°
		"magenta": 330,  // 330°
	}

	// Apply each color's adjustments with distance-based weighting
	hueAdj := 0.0
	satAdj := 0.0
	lumAdj := 0.0

	colors := []struct {
		name string
		adj  *models.ColorAdjustment
	}{
		{"red", &recipe.Red},
		{"orange", &recipe.Orange},
		{"yellow", &recipe.Yellow},
		{"green", &recipe.Green},
		{"aqua", &recipe.Aqua},
		{"blue", &recipe.Blue},
		{"purple", &recipe.Purple},
		{"magenta", &recipe.Magenta},
	}

	for _, color := range colors {
		center := colorCenters[color.name]
		weight := calculateHueWeight(h, center)

		if weight > 0 {
			hueAdj += float64(color.adj.Hue) * weight
			satAdj += float64(color.adj.Saturation) * weight
			lumAdj += float64(color.adj.Luminance) * weight
		}
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
// the input hue based on distance from the color's center hue.
// Uses a bell curve with 60° width for smooth transitions.
func calculateHueWeight(inputHue, centerHue float64) float64 {
	// Calculate angular distance (handle wraparound at 0/360)
	diff := math.Abs(inputHue - centerHue)
	if diff > 180 {
		diff = 360 - diff
	}

	// Bell curve with 60° width (covers ±30° from center)
	const width = 60.0
	if diff > width/2 {
		return 0
	}

	// Cosine-based smooth falloff
	return (math.Cos(diff*math.Pi/width) + 1) / 2
}

// applyColorGrading applies the 3-zone color grading adjustments.
// Zones are based on luminance: shadows (dark), midtones (medium), highlights (bright).
func applyColorGrading(recipe *models.UniversalRecipe, h, s, l float64) (float64, float64, float64) {
	if recipe.ColorGrading == nil {
		return h, s, l
	}

	cg := recipe.ColorGrading

	// Calculate zone weights based on luminance
	// Shadows: 0-0.33, Midtones: 0.33-0.67, Highlights: 0.67-1.0
	shadowWeight := calculateLuminanceWeight(l, 0, 0.33)
	midtoneWeight := calculateLuminanceWeight(l, 0.33, 0.67)
	highlightWeight := calculateLuminanceWeight(l, 0.67, 1.0)

	// Apply zone-specific adjustments with weights
	hueAdj := 0.0
	chromaAdj := 0.0
	brightnessAdj := 0.0

	if shadowWeight > 0 {
		// Apply Balance as intensity multiplier (see formatColorGradingZoneChroma)
		intensity := float64(cg.Balance+100) / 200.0
		hueAdj += float64(cg.Shadows.Hue) * shadowWeight
		chromaAdj += float64(cg.Shadows.Chroma) * intensity * shadowWeight
		brightnessAdj += float64(cg.Shadows.Brightness) * shadowWeight
	}

	if midtoneWeight > 0 {
		intensity := float64(cg.Balance+100) / 200.0
		hueAdj += float64(cg.Midtone.Hue) * midtoneWeight
		chromaAdj += float64(cg.Midtone.Chroma) * intensity * midtoneWeight
		brightnessAdj += float64(cg.Midtone.Brightness) * midtoneWeight
	}

	if highlightWeight > 0 {
		intensity := float64(cg.Balance+100) / 200.0
		hueAdj += float64(cg.Highlights.Hue) * highlightWeight
		chromaAdj += float64(cg.Highlights.Chroma) * intensity * highlightWeight
		brightnessAdj += float64(cg.Highlights.Brightness) * highlightWeight
	}

	// Apply adjustments
	h = math.Mod(h+hueAdj, 360)
	if h < 0 {
		h += 360
	}
	s = clamp(s+(chromaAdj/100.0), 0, 1)
	l = clamp(l+(brightnessAdj/100.0), 0, 1)

	return h, s, l
}

// calculateLuminanceWeight calculates the weight for a luminance zone
// using smooth transitions at zone boundaries.
func calculateLuminanceWeight(lum, zoneStart, zoneEnd float64) float64 {
	if lum < zoneStart || lum > zoneEnd {
		return 0
	}

	zoneMid := (zoneStart + zoneEnd) / 2
	zoneWidth := zoneEnd - zoneStart

	// Peak weight at zone center, smooth falloff at edges
	dist := math.Abs(lum - zoneMid)
	return 1.0 - (dist / (zoneWidth / 2))
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
	// Similar to standard Z85 but adjusted character set
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
