package lut

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// --- Existing regression tests ---

func TestGenerate3DLUTForPreview_OutputSize(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size 9", 9},
		{"size 17", 17},
		{"size 32", 32},
		{"size 33", 33},
	}

	recipe := &models.UniversalRecipe{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := generateAndValidate(t, recipe, tt.size)
			if err != nil {
				t.Fatalf("Generate3DLUTForPreview failed: %v", err)
			}

			expectedBytes := tt.size * tt.size * tt.size * 4 * 4 // RGBA float32
			if len(data) != expectedBytes {
				t.Errorf("expected %d bytes, got %d", expectedBytes, len(data))
			}
		})
	}
}

func TestGenerate3DLUTForPreview_FloatRange(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure:   1.5,
		Contrast:   50,
		Saturation: 80,
	}

	data, err := generateAndValidate(t, recipe, 9)
	if err != nil {
		t.Fatalf("Generate3DLUTForPreview failed: %v", err)
	}

	totalFloats := len(data) / 4
	for i := 0; i < totalFloats; i++ {
		val := math.Float32frombits(binary.LittleEndian.Uint32(data[i*4 : i*4+4]))

		channel := i % 4
		if channel == 3 {
			// Alpha must be exactly 1.0
			if val != 1.0 {
				t.Errorf("alpha at float index %d: expected 1.0, got %f", i, val)
			}
		} else {
			// RGB must be in [0, 1]
			if val < 0.0 || val > 1.0 {
				t.Errorf("RGB value at float index %d: %f out of [0,1] range", i, val)
			}
		}
	}
}

func TestGenerate3DLUTForPreview_IdentityRecipe(t *testing.T) {
	// Identity recipe: all defaults (zero values) should produce near-identity LUT
	recipe := &models.UniversalRecipe{}
	size := 17

	data, err := generateAndValidate(t, recipe, size)
	if err != nil {
		t.Fatalf("Generate3DLUTForPreview failed: %v", err)
	}

	maxDelta := float32(0.0)
	texelIdx := 0
	for bIdx := 0; bIdx < size; bIdx++ {
		for gIdx := 0; gIdx < size; gIdx++ {
			for rIdx := 0; rIdx < size; rIdx++ {
				offset := texelIdx * 16 // 4 floats × 4 bytes
				rOut := math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
				gOut := math.Float32frombits(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
				bOut := math.Float32frombits(binary.LittleEndian.Uint32(data[offset+8 : offset+12]))

				rExpected := float32(rIdx) / float32(size-1)
				gExpected := float32(gIdx) / float32(size-1)
				bExpected := float32(bIdx) / float32(size-1)

				deltaR := float32(math.Abs(float64(rOut - rExpected)))
				deltaG := float32(math.Abs(float64(gOut - gExpected)))
				deltaB := float32(math.Abs(float64(bOut - bExpected)))

				for _, d := range []float32{deltaR, deltaG, deltaB} {
					if d > maxDelta {
						maxDelta = d
					}
				}

				texelIdx++
			}
		}
	}

	// Identity recipe should produce near-identity output
	// Allow small epsilon for float64→float32 conversion and HSL round-trip
	if maxDelta > 0.01 {
		t.Errorf("identity recipe max delta too large: %f (expected < 0.01)", maxDelta)
	}
}

func TestGenerate3DLUTForPreview_InvalidInputs(t *testing.T) {
	t.Run("nil recipe", func(t *testing.T) {
		_, err := Generate3DLUTForPreview(nil, 17)
		if err == nil {
			t.Error("expected error for nil recipe")
		}
	})

	t.Run("size too small", func(t *testing.T) {
		_, err := Generate3DLUTForPreview(&models.UniversalRecipe{}, 1)
		if err == nil {
			t.Error("expected error for size < 2")
		}
	})

	t.Run("size too large", func(t *testing.T) {
		_, err := Generate3DLUTForPreview(&models.UniversalRecipe{}, 128)
		if err == nil {
			t.Error("expected error for size > 65")
		}
	})
}

func TestGenerate3DLUT_BackwardCompat(t *testing.T) {
	// Verify original Generate3DLUT output size is unchanged
	recipe := &models.UniversalRecipe{
		Exposure: 0.5,
		Contrast: 20,
	}

	data, err := Generate3DLUT(recipe)
	if err != nil {
		t.Fatalf("Generate3DLUT failed: %v", err)
	}

	// Original: 32³ × 3 × 4 = 393,216 bytes (RGB float32)
	expectedBytes := 32 * 32 * 32 * 3 * 4
	if len(data) != expectedBytes {
		t.Errorf("Generate3DLUT output size changed: expected %d, got %d", expectedBytes, len(data))
	}
}

func TestRecipeJSONRoundTrip(t *testing.T) {
	original := models.UniversalRecipe{
		Name:       "Test Recipe",
		Exposure:   1.5,
		Contrast:   30,
		Highlights: -20,
		Shadows:    40,
		Whites:     10,
		Blacks:     -15,
		Saturation: 25,
		Clarity:    15,
		Sharpness:  80,
		Red:        models.ColorAdjustment{Hue: 10, Saturation: 20, Luminance: -5},
		Orange:     models.ColorAdjustment{Hue: -15, Saturation: 0, Luminance: 10},
		ColorGrading: &models.ColorGrading{
			Highlights: models.ColorGradingZone{Hue: 45, Chroma: 30, Brightness: 10},
			Midtone:    models.ColorGradingZone{Hue: 200, Chroma: 20, Brightness: 0},
			Shadows:    models.ColorGradingZone{Hue: 220, Chroma: 50, Brightness: -10},
			Blending:   50,
			Balance:    -20,
		},
	}

	jsonBytes, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var restored models.UniversalRecipe
	if err := json.Unmarshal(jsonBytes, &restored); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if restored.Exposure != original.Exposure {
		t.Errorf("Exposure: got %f, want %f", restored.Exposure, original.Exposure)
	}
	if restored.Contrast != original.Contrast {
		t.Errorf("Contrast: got %d, want %d", restored.Contrast, original.Contrast)
	}
	if restored.Red.Hue != original.Red.Hue {
		t.Errorf("Red.Hue: got %d, want %d", restored.Red.Hue, original.Red.Hue)
	}
	if restored.ColorGrading == nil {
		t.Fatal("ColorGrading is nil after round-trip")
	}
	if restored.ColorGrading.Shadows.Hue != original.ColorGrading.Shadows.Hue {
		t.Errorf("ColorGrading.Shadows.Hue: got %d, want %d",
			restored.ColorGrading.Shadows.Hue, original.ColorGrading.Shadows.Hue)
	}
}

func generateAndValidate(t *testing.T, recipe *models.UniversalRecipe, size int) ([]byte, error) {
	t.Helper()
	return Generate3DLUTForPreview(recipe, size)
}

// --- Task 1: HSL + Gamma function tests ---

func TestRgbToHslKnownValues(t *testing.T) {
	tests := []struct {
		name       string
		r, g, b    float64
		wantH      float64
		wantS      float64
		wantL      float64
		hTolerance float64
	}{
		{"red", 1, 0, 0, 0, 1, 0.5, 1},
		{"green", 0, 1, 0, 120, 1, 0.5, 1},
		{"blue", 0, 0, 1, 240, 1, 0.5, 1},
		{"white", 1, 1, 1, 0, 0, 1, 1},
		{"black", 0, 0, 0, 0, 0, 0, 1},
		{"gray 50%", 0.5, 0.5, 0.5, 0, 0, 0.5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, s, l := rgbToHsl(tt.r, tt.g, tt.b)

			if math.Abs(h-tt.wantH) > tt.hTolerance {
				t.Errorf("hue: got %f, want %f", h, tt.wantH)
			}
			if math.Abs(s-tt.wantS) > 0.001 {
				t.Errorf("saturation: got %f, want %f", s, tt.wantS)
			}
			if math.Abs(l-tt.wantL) > 0.001 {
				t.Errorf("lightness: got %f, want %f", l, tt.wantL)
			}
		})
	}
}

func TestRgbToHslRedWraparound(t *testing.T) {
	// AC 5: RGB (1.0, 0.0, 0.1) should produce hue ~354°, not -6°
	h, s, _ := rgbToHsl(1.0, 0.0, 0.1)

	if h < 350 || h > 360 {
		t.Errorf("hue wraparound: got %f, expected ~354° (350-360 range)", h)
	}
	if s < 0.9 {
		t.Errorf("saturation too low: got %f, expected near 1.0", s)
	}
}

func TestHslRoundTrip(t *testing.T) {
	// AC 4: 100 random RGB colors → HSL → RGB, max delta < 0.001
	maxDelta := 0.0
	for i := 0; i < 100; i++ {
		// Deterministic "random" values
		r := float64(i%10) / 9.0
		g := float64((i*7)%10) / 9.0
		b := float64((i*13)%10) / 9.0

		h, s, l := rgbToHsl(r, g, b)
		rOut, gOut, bOut := hslToRgb(h, s, l)

		for _, d := range []float64{
			math.Abs(r - rOut),
			math.Abs(g - gOut),
			math.Abs(b - bOut),
		} {
			if d > maxDelta {
				maxDelta = d
			}
		}
	}

	if maxDelta > 0.001 {
		t.Errorf("HSL round-trip max delta: %f (expected < 0.001)", maxDelta)
	}
}

func TestSrgbGammaRoundTrip(t *testing.T) {
	maxDelta := 0.0
	for i := 0; i <= 100; i++ {
		x := float64(i) / 100.0
		rt := linearToSrgb(srgbToLinear(x))
		d := math.Abs(x - rt)
		if d > maxDelta {
			maxDelta = d
		}
	}

	if maxDelta > 0.001 {
		t.Errorf("sRGB gamma round-trip max delta: %f (expected < 0.001)", maxDelta)
	}
}

// --- Task 2: Exposure tests ---

func TestApplyExposure(t *testing.T) {
	t.Run("exposure +1.0 brightens", func(t *testing.T) {
		recipe := &models.UniversalRecipe{Exposure: 1.0}
		r, g, b := applyColorTransform(recipe, 0.5, 0.5, 0.5)

		// AC 6: output ≈ 0.684 (linearToSrgb(srgbToLinear(0.5) * 2))
		expected := math.Pow(math.Pow(0.5, 2.2)*2, 1.0/2.2)
		if math.Abs(r-expected) > 0.01 {
			t.Errorf("exposure +1.0: got %f, expected ~%f", r, expected)
		}
		if r <= 0.5 || g <= 0.5 || b <= 0.5 {
			t.Errorf("exposure +1.0 should brighten: got (%f, %f, %f)", r, g, b)
		}
		if r >= 1.0 {
			t.Errorf("exposure +1.0 on 0.5 should not clip: got %f", r)
		}
	})

	t.Run("exposure -1.0 darkens", func(t *testing.T) {
		recipe := &models.UniversalRecipe{Exposure: -1.0}
		r, _, _ := applyColorTransform(recipe, 0.5, 0.5, 0.5)
		if r >= 0.5 {
			t.Errorf("exposure -1.0 should darken: got %f", r)
		}
	})

	t.Run("exposure 0 is identity", func(t *testing.T) {
		recipe := &models.UniversalRecipe{Exposure: 0}
		r, g, b := applyColorTransform(recipe, 0.5, 0.5, 0.5)
		if math.Abs(r-0.5) > 0.001 || math.Abs(g-0.5) > 0.001 || math.Abs(b-0.5) > 0.001 {
			t.Errorf("exposure 0 should be identity: got (%f, %f, %f)", r, g, b)
		}
	})
}

// --- Task 3: Tone curve tests ---

func TestApplyToneCurve(t *testing.T) {
	t.Run("contrast +50 preserves midpoint", func(t *testing.T) {
		// AC 7: contrast=+50 at 0.5 → output ≈ 0.5
		recipe := &models.UniversalRecipe{Contrast: 50}
		out := applyToneCurve(recipe, 0.5)
		if math.Abs(out-0.5) > 0.01 {
			t.Errorf("contrast +50 at midpoint: got %f, want ~0.5", out)
		}
	})

	t.Run("contrast +50 darkens shadows", func(t *testing.T) {
		// AC 7: contrast=+50 at 0.25 → output < 0.25
		recipe := &models.UniversalRecipe{Contrast: 50}
		out := applyToneCurve(recipe, 0.25)
		if out >= 0.25 {
			t.Errorf("contrast +50 at 0.25: got %f, want < 0.25", out)
		}
	})

	t.Run("contrast +50 brightens highlights", func(t *testing.T) {
		recipe := &models.UniversalRecipe{Contrast: 50}
		out := applyToneCurve(recipe, 0.75)
		if out <= 0.75 {
			t.Errorf("contrast +50 at 0.75: got %f, want > 0.75", out)
		}
	})

	t.Run("negative contrast flattens", func(t *testing.T) {
		// AC 12: contrast=-50 at 0.25 → output > 0.25, at 0.75 → output < 0.75
		recipe := &models.UniversalRecipe{Contrast: -50}
		outLow := applyToneCurve(recipe, 0.25)
		outHigh := applyToneCurve(recipe, 0.75)
		if outLow <= 0.25 {
			t.Errorf("contrast -50 at 0.25: got %f, want > 0.25", outLow)
		}
		if outHigh >= 0.75 {
			t.Errorf("contrast -50 at 0.75: got %f, want < 0.75", outHigh)
		}
	})

	t.Run("highlights recovery", func(t *testing.T) {
		recipe := &models.UniversalRecipe{Highlights: -100}
		out := applyToneCurve(recipe, 0.9)
		if out >= 0.9 {
			t.Errorf("highlights -100 at 0.9: got %f, want < 0.9", out)
		}
	})

	t.Run("shadow lift", func(t *testing.T) {
		recipe := &models.UniversalRecipe{Shadows: 100}
		out := applyToneCurve(recipe, 0.1)
		if out <= 0.1 {
			t.Errorf("shadows +100 at 0.1: got %f, want > 0.1", out)
		}
	})

	t.Run("all zero identity", func(t *testing.T) {
		recipe := &models.UniversalRecipe{}
		for _, l := range []float64{0, 0.1, 0.25, 0.5, 0.75, 0.9, 1.0} {
			out := applyToneCurve(recipe, l)
			if math.Abs(out-l) > 0.001 {
				t.Errorf("identity at l=%f: got %f", l, out)
			}
		}
	})
}

// --- Task 4: Color blender tests ---

func TestApplyColorBlender(t *testing.T) {
	t.Run("red hue shift on red pixel", func(t *testing.T) {
		// Hue=30 on perceptual scale → 30 * 0.3 = 9° actual angular shift
		recipe := &models.UniversalRecipe{
			Red: models.ColorAdjustment{Hue: 30},
		}
		h, _, _ := applyColorBlender(recipe, 0, 0.8, 0.5)
		if math.Abs(h-9) > 2 {
			t.Errorf("red hue +30 on red pixel: got hue %f, want ~9 (30 * 0.3)", h)
		}
	})

	t.Run("red hue shift ignored on blue pixel", func(t *testing.T) {
		// AC 8: pixel at hue=240° should not be affected by red (center 0°, band 90°)
		recipe := &models.UniversalRecipe{
			Red: models.ColorAdjustment{Hue: 30},
		}
		h, _, _ := applyColorBlender(recipe, 240, 0.8, 0.5)
		if math.Abs(h-240) > 0.01 {
			t.Errorf("red hue +30 on blue pixel: got hue %f, want 240 (no change)", h)
		}
	})

	t.Run("band width covers 40 degrees from red center", func(t *testing.T) {
		// AC 8: pixel at hue=40° (within 90° band from center 0°) should shift
		recipe := &models.UniversalRecipe{
			Red: models.ColorAdjustment{Hue: 30},
		}
		h, _, _ := applyColorBlender(recipe, 40, 0.8, 0.5)
		if math.Abs(h-40) < 0.5 {
			t.Errorf("red hue +30 on pixel at 40°: expected shift, got hue %f", h)
		}
	})

	t.Run("pixel at 100 degrees not affected by red", func(t *testing.T) {
		// AC 8: pixel at hue=100° (outside 90° band from center 0°) should not shift
		recipe := &models.UniversalRecipe{
			Red: models.ColorAdjustment{Hue: 30},
		}
		h, _, _ := applyColorBlender(recipe, 100, 0.8, 0.5)
		if math.Abs(h-100) > 1 {
			t.Errorf("red hue +30 on pixel at 100°: got hue %f, want ~100 (no change)", h)
		}
	})

	t.Run("overlap between red and yellow at 45 degrees", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			Red:    models.ColorAdjustment{Hue: 10},
			Yellow: models.ColorAdjustment{Hue: 10},
		}
		h, _, _ := applyColorBlender(recipe, 45, 0.8, 0.5)
		// Both red (center 0°) and yellow (center 60°) should affect pixel at 45°
		if math.Abs(h-45) < 0.5 {
			t.Errorf("overlap test at 45°: expected shift, got hue %f", h)
		}
	})

	t.Run("overlap normalization at 5 degrees", func(t *testing.T) {
		// AC 16: Red+Orange+Magenta all +30 → normalized shift ≈9° (30*0.3), not 27° (90*0.3)
		recipe := &models.UniversalRecipe{
			Red:     models.ColorAdjustment{Hue: 30},
			Orange:  models.ColorAdjustment{Hue: 30},
			Magenta: models.ColorAdjustment{Hue: 30},
		}
		h, _, _ := applyColorBlender(recipe, 5, 0.8, 0.5)
		shift := h - 5
		if shift < 0 {
			shift += 360
		}
		// Should be ~9° (normalized 30*0.3), definitely not 27° (triple-counted 90*0.3)
		if shift > 15 {
			t.Errorf("overlap normalization: shift=%f, expected ~9 (not 27)", shift)
		}
	})

	t.Run("near-gray guard in pipeline", func(t *testing.T) {
		// Pixels with s < 0.02 should not be passed to blender (guard in applyColorTransform)
		// Verify by running through full pipeline: a near-gray pixel should not get hue shifted
		recipe := &models.UniversalRecipe{
			Red: models.ColorAdjustment{Hue: 90, Saturation: 50},
		}
		// Near-gray: r≈g≈b, very low saturation
		rOut, gOut, bOut := applyColorTransform(recipe, 0.5, 0.495, 0.505)
		// The output should still be near-gray (not shifted by the red hue adjustment)
		maxDiff := math.Max(math.Abs(rOut-gOut), math.Abs(gOut-bOut))
		if maxDiff > 0.1 {
			t.Errorf("near-gray pixel was color-shifted: RGB=(%f, %f, %f), maxChannelDiff=%f", rOut, gOut, bOut, maxDiff)
		}
	})

	t.Run("orange saturation increase", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			Orange: models.ColorAdjustment{Saturation: 50},
		}
		_, s, _ := applyColorBlender(recipe, 30, 0.5, 0.5)
		if s <= 0.5 {
			t.Errorf("orange sat +50 on orange pixel: got sat %f, want > 0.5", s)
		}
	})
}

// --- Task 5: Color grading tests ---

func TestApplyColorGrading(t *testing.T) {
	t.Run("THE critical test: shadows hue 220 tints toward blue", func(t *testing.T) {
		// AC 1: Shadows Hue=220 Chroma=50 on a red pixel at lum=0.2
		// Output hue should be TOWARD 220° (blue), NOT rotated by 220°
		recipe := &models.UniversalRecipe{
			ColorGrading: &models.ColorGrading{
				Shadows:  models.ColorGradingZone{Hue: 220, Chroma: 50},
				Blending: 50,
			},
		}
		h, _, _ := applyColorGrading(recipe, 0, 0.5, 0.2)
		// Output should be between 180° and 260° (blue range), NOT green
		if h < 180 || h > 260 {
			t.Errorf("CRITICAL: shadow tint hue=220 on red pixel: got hue=%f, expected 180-260 (blue range)", h)
		}
	})

	t.Run("chroma 0 is identity for hue", func(t *testing.T) {
		// AC 13: Chroma=0 → no hue change
		recipe := &models.UniversalRecipe{
			ColorGrading: &models.ColorGrading{
				Shadows:  models.ColorGradingZone{Hue: 220, Chroma: 0},
				Blending: 50,
			},
		}
		h, s, _ := applyColorGrading(recipe, 30, 0.5, 0.2)
		if math.Abs(h-30) > 0.01 {
			t.Errorf("chroma 0 should not change hue: got %f, want 30", h)
		}
		if math.Abs(s-0.5) > 0.01 {
			t.Errorf("chroma 0 should not change saturation: got %f, want 0.5", s)
		}
	})

	t.Run("brightness in shadows affects dark pixel", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			ColorGrading: &models.ColorGrading{
				Shadows:  models.ColorGradingZone{Brightness: 50, Chroma: 50, Hue: 200},
				Blending: 50,
			},
		}
		_, _, l := applyColorGrading(recipe, 0, 0.5, 0.1)
		if l <= 0.1 {
			t.Errorf("shadow brightness +50 on dark pixel: got lum=%f, want > 0.1", l)
		}
	})

	t.Run("highlight zone affects bright pixel not dark", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			ColorGrading: &models.ColorGrading{
				Highlights: models.ColorGradingZone{Hue: 100, Chroma: 80},
				Blending:   50,
			},
		}
		// Bright pixel should be affected
		hBright, _, _ := applyColorGrading(recipe, 0, 0.5, 0.9)
		// Dark pixel should NOT be affected (or minimally)
		hDark, _, _ := applyColorGrading(recipe, 0, 0.5, 0.1)

		brightShift := math.Abs(hBright - 0)
		if brightShift > 180 {
			brightShift = 360 - brightShift
		}
		darkShift := math.Abs(hDark - 0)
		if darkShift > 180 {
			darkShift = 360 - darkShift
		}

		if brightShift <= darkShift {
			t.Errorf("highlight grading: bright shift (%f) should be larger than dark shift (%f)", brightShift, darkShift)
		}
	})

	t.Run("achromatic tint", func(t *testing.T) {
		// AC 14: neutral gray pixel (s≈0) should get visible tint
		recipe := &models.UniversalRecipe{
			ColorGrading: &models.ColorGrading{
				Shadows:  models.ColorGradingZone{Hue: 200, Chroma: 80},
				Blending: 50,
			},
		}
		h, s, _ := applyColorGrading(recipe, 0, 0.0, 0.2)
		if s < 0.1 {
			t.Errorf("achromatic tint: saturation should be > 0.1, got %f", s)
		}
		// Hue should be near 200°
		hueDiff := math.Abs(h - 200)
		if hueDiff > 180 {
			hueDiff = 360 - hueDiff
		}
		if hueDiff > 30 {
			t.Errorf("achromatic tint: hue should be near 200, got %f", h)
		}
	})

	t.Run("blending 100 no crash", func(t *testing.T) {
		// AC 15: Blending=100 should not crash or produce NaN
		recipe := &models.UniversalRecipe{
			ColorGrading: &models.ColorGrading{
				Shadows:    models.ColorGradingZone{Hue: 200, Chroma: 50},
				Highlights: models.ColorGradingZone{Hue: 40, Chroma: 50},
				Blending:   100,
			},
		}
		h, s, l := applyColorGrading(recipe, 30, 0.5, 0.5)
		if math.IsNaN(h) || math.IsNaN(s) || math.IsNaN(l) {
			t.Errorf("blending 100 produced NaN: h=%f s=%f l=%f", h, s, l)
		}
	})

	t.Run("balance extremes no crash", func(t *testing.T) {
		// AC 17: Balance=-100 should not crash
		recipe := &models.UniversalRecipe{
			ColorGrading: &models.ColorGrading{
				Shadows:  models.ColorGradingZone{Hue: 200, Chroma: 50},
				Blending: 50,
				Balance:  -100,
			},
		}
		h, s, l := applyColorGrading(recipe, 30, 0.5, 0.2)
		if math.IsNaN(h) || math.IsNaN(s) || math.IsNaN(l) {
			t.Errorf("balance -100 produced NaN: h=%f s=%f l=%f", h, s, l)
		}
	})
}

// --- Task 6: Cross-validation test ---

func TestLUTCrossValidation(t *testing.T) {
	// AC 9: Both LUT generators should produce matching RGB values
	recipe := &models.UniversalRecipe{
		Exposure: 0.5,
		Contrast: 20,
		Red:      models.ColorAdjustment{Hue: 10},
		ColorGrading: &models.ColorGrading{
			Shadows:  models.ColorGradingZone{Hue: 200, Chroma: 30},
			Blending: 50,
		},
	}

	rgbData, err := Generate3DLUT(recipe)
	if err != nil {
		t.Fatalf("Generate3DLUT failed: %v", err)
	}

	rgbaData, err := Generate3DLUTForPreview(recipe, LUTSize)
	if err != nil {
		t.Fatalf("Generate3DLUTForPreview failed: %v", err)
	}

	maxDelta := float32(0.0)
	for i := 0; i < LUTPoints; i++ {
		// RGB from Generate3DLUT: 3 floats per texel
		rgbOff := i * 3 * 4
		rRGB := math.Float32frombits(binary.LittleEndian.Uint32(rgbData[rgbOff : rgbOff+4]))
		gRGB := math.Float32frombits(binary.LittleEndian.Uint32(rgbData[rgbOff+4 : rgbOff+8]))
		bRGB := math.Float32frombits(binary.LittleEndian.Uint32(rgbData[rgbOff+8 : rgbOff+12]))

		// RGBA from Generate3DLUTForPreview: 4 floats per texel
		rgbaOff := i * 4 * 4
		rRGBA := math.Float32frombits(binary.LittleEndian.Uint32(rgbaData[rgbaOff : rgbaOff+4]))
		gRGBA := math.Float32frombits(binary.LittleEndian.Uint32(rgbaData[rgbaOff+4 : rgbaOff+8]))
		bRGBA := math.Float32frombits(binary.LittleEndian.Uint32(rgbaData[rgbaOff+8 : rgbaOff+12]))

		for _, d := range []float32{
			float32(math.Abs(float64(rRGB - rRGBA))),
			float32(math.Abs(float64(gRGB - gRGBA))),
			float32(math.Abs(float64(bRGB - bRGBA))),
		} {
			if d > maxDelta {
				maxDelta = d
			}
		}
	}

	if maxDelta > 0.0001 {
		t.Errorf("cross-validation max delta: %f (expected < 0.0001)", maxDelta)
	}
}

// --- AC 3: Range safety with extreme values ---

func TestGenerate3DLUTForPreview_ExtremeValues(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure:   5,
		Contrast:   100,
		Highlights: 100,
		Shadows:    -100,
		Whites:     100,
		Blacks:     -100,
		Saturation: 100,
		Red:        models.ColorAdjustment{Hue: 180, Saturation: 100, Luminance: 100},
		Orange:     models.ColorAdjustment{Hue: -180, Saturation: -100, Luminance: -100},
		Yellow:     models.ColorAdjustment{Hue: 180, Saturation: 100, Luminance: 100},
		Green:      models.ColorAdjustment{Hue: -180, Saturation: -100, Luminance: -100},
		Aqua:       models.ColorAdjustment{Hue: 180, Saturation: 100, Luminance: 100},
		Blue:       models.ColorAdjustment{Hue: -180, Saturation: -100, Luminance: -100},
		Purple:     models.ColorAdjustment{Hue: 180, Saturation: 100, Luminance: 100},
		Magenta:    models.ColorAdjustment{Hue: -180, Saturation: -100, Luminance: -100},
		ColorGrading: &models.ColorGrading{
			Highlights: models.ColorGradingZone{Hue: 360, Chroma: 100, Brightness: 100},
			Midtone:    models.ColorGradingZone{Hue: 180, Chroma: -100, Brightness: -100},
			Shadows:    models.ColorGradingZone{Hue: 0, Chroma: 100, Brightness: 100},
			Blending:   100,
			Balance:    100,
		},
	}

	data, err := Generate3DLUTForPreview(recipe, 9)
	if err != nil {
		t.Fatalf("Generate3DLUTForPreview with extreme values failed: %v", err)
	}

	totalFloats := len(data) / 4
	for i := 0; i < totalFloats; i++ {
		val := math.Float32frombits(binary.LittleEndian.Uint32(data[i*4 : i*4+4]))
		if math.IsNaN(float64(val)) || math.IsInf(float64(val), 0) {
			t.Fatalf("NaN/Inf at float index %d", i)
		}
		channel := i % 4
		if channel == 3 {
			continue // alpha
		}
		if val < 0.0 || val > 1.0 {
			t.Errorf("RGB value at float index %d: %f out of [0,1] range", i, val)
		}
	}
}
