package lut

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"testing"

	"github.com/justin/recipe/internal/models"
)

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
	// Verify original Generate3DLUT is unchanged
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

	// Verify key fields survived round-trip
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
