package np3

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestConvertToNP3ParametersDirectly tests parameter conversion with edge values
// by testing the internal function directly
func TestConvertToNP3ParametersDirectly(t *testing.T) {
	tests := []struct {
		name       string
		recipe     *models.UniversalRecipe
		wantErr    bool
		checkFunc  func(*testing.T, *np3Parameters)
	}{
		{
			name: "Sharpness at exact upper boundary (150 → 9)",
			recipe: &models.UniversalRecipe{
				Sharpness:  150,
				Contrast:   0,
				Saturation: 0,
				Exposure:   0.0,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, params *np3Parameters) {
				if params.sharpening != 9.0 && params.sharpening != 9 {
					t.Errorf("Sharpness: got %.2f, want 9", params.sharpening)
				}
			},
		},
		{
			name: "Contrast at upper boundary (100 → 100) - Phase 2 direct mapping",
			recipe: &models.UniversalRecipe{
				Sharpness:  0,
				Contrast:   100,
				Saturation: 0,
				Exposure:   0.0,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, params *np3Parameters) {
				if params.contrast != 100 {
					t.Errorf("Contrast: got %d, want 100", params.contrast)
				}
			},
		},
		{
			name: "Contrast at lower boundary (-100 → -100) - Phase 2 direct mapping",
			recipe: &models.UniversalRecipe{
				Sharpness:  0,
				Contrast:   -100,
				Saturation: 0,
				Exposure:   0.0,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, params *np3Parameters) {
				if params.contrast != -100 {
					t.Errorf("Contrast: got %d, want -100", params.contrast)
				}
			},
		},
		{
			name: "Saturation at upper boundary (100 → 100) - Phase 2 direct mapping",
			recipe: &models.UniversalRecipe{
				Sharpness:  0,
				Contrast:   0,
				Saturation: 100,
				Exposure:   0.0,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, params *np3Parameters) {
				if params.saturation != 100 {
					t.Errorf("Saturation: got %d, want 100", params.saturation)
				}
			},
		},
		{
			name: "Saturation at lower boundary (-100 → -100) - Phase 2 direct mapping",
			recipe: &models.UniversalRecipe{
				Sharpness:  0,
				Contrast:   0,
				Saturation: -100,
				Exposure:   0.0,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, params *np3Parameters) {
				if params.saturation != -100 {
					t.Errorf("Saturation: got %d, want -100", params.saturation)
				}
			},
		},
		{
			name: "Exposure at upper boundary (1.0)",
			recipe: &models.UniversalRecipe{
				Sharpness:  0,
				Contrast:   0,
				Saturation: 0,
				Exposure:   1.0,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, params *np3Parameters) {
				if params.brightness != 1.0 {
					t.Errorf("Brightness: got %f, want 1.0", params.brightness)
				}
			},
		},
		{
			name: "Exposure at lower boundary (-1.0)",
			recipe: &models.UniversalRecipe{
				Sharpness:  0,
				Contrast:   0,
				Saturation: 0,
				Exposure:   -1.0,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, params *np3Parameters) {
				if params.brightness != -1.0 {
					t.Errorf("Brightness: got %f, want -1.0", params.brightness)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := convertToNP3Parameters(tt.recipe)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToNP3Parameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, params)
			}
		})
	}
}

// TestEncodeBinaryMinFileSize tests that files are exactly 480 bytes (NP3 standard size)
func TestEncodeBinaryMinFileSize(t *testing.T) {
	params := &np3Parameters{
		name:       "",
		sharpening: 0,
		contrast:   0,
		saturation: 0,
		brightness: 0.0,
		hue:        0,
	}

	data, err := encodeBinary(params, "")
	if err != nil {
		t.Fatalf("encodeBinary failed: %v", err)
	}

	// Phase 2: NP3 files are 1050 bytes, matching working NP3 files (e.g., Filmic.np3)
	// Note: Real NP3 files vary in size (392-1050 bytes), but generator uses 1050 for compatibility
	if len(data) != 1050 {
		t.Errorf("File size: got %d, want 1050", len(data))
	}
}

// TestGenerateColorDataNegativeSaturation tests color data with negative saturation
func TestGenerateColorDataNegativeSaturation(t *testing.T) {
	data := make([]byte, 500)

	// Test with saturation = -3 (should generate 0 triplets, targetCount < 0)
	offset := generateColorData(data, -3)

	// With saturation=-3: targetCount = (-3 + 1) * 15 = -30 → clamped to 0
	// Should start at offset 100 and generate nothing
	if offset != 100 {
		t.Errorf("End offset: got %d, want 100 (no triplets generated)", offset)
	}

	// Verify no color data was written
	for i := 100; i < 300; i++ {
		if data[i] != 0 {
			t.Errorf("Expected no color data, but found byte %d at offset %d", data[i], i)
			break
		}
	}
}

// TestGenerateColorDataMaxSaturation tests color data at maximum saturation
func TestGenerateColorDataMaxSaturation(t *testing.T) {
	data := make([]byte, 500)

	// Test with saturation = +3 (should generate 60 triplets)
	offset := generateColorData(data, 3)

	// With saturation=3: targetCount = (3 + 1) * 15 = 60 triplets
	// 60 triplets * 3 bytes = 180 bytes → offset should be 100 + 180 = 280
	expectedOffset := 100 + (60 * 3)
	if offset != expectedOffset {
		t.Errorf("End offset: got %d, want %d", offset, expectedOffset)
	}

	// Verify color data was written
	tripletCount := 0
	for i := 100; i < offset; i += 3 {
		if data[i] == 50 && data[i+1] == 50 && data[i+2] == 50 {
			tripletCount++
		}
	}
	if tripletCount != 60 {
		t.Errorf("Color triplets: got %d, want 60", tripletCount)
	}
}

// TestGenerateToneCurveNegativeContrast tests tone curve with contrast that would
// result in negative additional pairs
func TestGenerateToneCurveNegativeContrast(t *testing.T) {
	data := make([]byte, 500)

	// First generate lots of color data that extends into tone curve region
	colorEndOffset := generateColorData(data, 3) // Ends at offset 280

	// With contrast=-3: targetTotalPairs = (-3 + 2) * 20 = -20
	// This is negative, but will be clamped in the function
	// overlapPairs from 150 to 279 = 65 pairs
	// additionalPairs = -20 - 65 = -85 → clamped to 0
	generateToneCurveData(data, -3, colorEndOffset)

	// No additional tone curve data should be written after colorEndOffset
	// Check that bytes after colorEndOffset are still zero
	for i := colorEndOffset; i < 500; i++ {
		if data[i] != 0 && data[i] != 50 { // Allow color data (50), but no tone curve (1)
			t.Errorf("Unexpected tone curve data at offset %d: %d", i, data[i])
		}
	}
}

// TestWriteRawParameterBytesSharpnessZero tests sharpness=0 special case
func TestWriteRawParameterBytesSharpnessZero(t *testing.T) {
	data := make([]byte, 500)
	params := &np3Parameters{
		sharpening: 0,
		brightness: 0.0,
		hue:        0,
	}

	writeRawParameterBytes(data, params)

	// Sharpness=0 should write byte value 1 (not 0) to avoid parser default
	for i := 66; i <= 70; i++ {
		if data[i] != 1 {
			t.Errorf("Sharpness byte at offset %d: got %d, want 1", i, data[i])
		}
	}
}

