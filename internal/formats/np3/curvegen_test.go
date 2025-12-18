package np3

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateTestCurve_Linear(t *testing.T) {
	lut, err := GenerateTestCurve(CurveLinear)
	if err != nil {
		t.Fatalf("GenerateTestCurve(Linear) failed: %v", err)
	}

	if len(lut) != 257 {
		t.Errorf("Expected 257 entries, got %d", len(lut))
	}

	// Linear curve: entry 0 should be 0, entry 256 should be 32767
	if lut[0] != 0 {
		t.Errorf("Expected lut[0]=0, got %d", lut[0])
	}
	if lut[256] != 32767 {
		t.Errorf("Expected lut[256]=32767, got %d", lut[256])
	}

	// Midpoint should be approximately half
	midpoint := lut[128]
	expectedMid := uint16(32767 / 2)
	tolerance := uint16(100)
	if midpoint < expectedMid-tolerance || midpoint > expectedMid+tolerance {
		t.Errorf("Expected lut[128]≈%d, got %d", expectedMid, midpoint)
	}

	// Curve should be monotonically increasing
	for i := 1; i < 257; i++ {
		if lut[i] < lut[i-1] {
			t.Errorf("Linear curve not monotonic at index %d: lut[%d]=%d < lut[%d]=%d",
				i, i, lut[i], i-1, lut[i-1])
		}
	}
}

func TestGenerateTestCurve_SCurve(t *testing.T) {
	lut, err := GenerateTestCurve(CurveSCurve)
	if err != nil {
		t.Fatalf("GenerateTestCurve(SCurve) failed: %v", err)
	}

	if len(lut) != 257 {
		t.Errorf("Expected 257 entries, got %d", len(lut))
	}

	// S-curve should compress shadows and highlights
	// Shadows should be pulled down (lower than linear)
	// Highlights should be pushed up (higher than linear) or compressed
	linearLut, _ := GenerateTestCurve(CurveLinear)

	// Check that shadows are different from linear
	shadowsDiff := int(lut[32]) - int(linearLut[32])
	t.Logf("Shadows diff (entry 32): %d", shadowsDiff)

	// Check that curve is still monotonically increasing
	for i := 1; i < 257; i++ {
		if lut[i] < lut[i-1] {
			t.Errorf("S-curve not monotonic at index %d", i)
		}
	}
}

func TestGenerateTestCurve_ShadowsBoost(t *testing.T) {
	lut, err := GenerateTestCurve(CurveShadowsBoost)
	if err != nil {
		t.Fatalf("GenerateTestCurve(ShadowsBoost) failed: %v", err)
	}

	linearLut, _ := GenerateTestCurve(CurveLinear)

	// Shadows boost should lift shadow values above linear
	if lut[32] <= linearLut[32] {
		t.Errorf("Shadows boost should lift shadows: lut[32]=%d should be > linear[32]=%d",
			lut[32], linearLut[32])
	}

	// Highlight values should be closer to linear (less affected)
	highlightDiff := int(lut[224]) - int(linearLut[224])
	shadowDiff := int(lut[32]) - int(linearLut[32])
	if shadowDiff <= highlightDiff {
		t.Errorf("Shadows should be affected more than highlights: shadowDiff=%d, highlightDiff=%d",
			shadowDiff, highlightDiff)
	}

	t.Logf("Shadow lift: +%d at entry 32, +%d at entry 224", shadowDiff, highlightDiff)
}

func TestGenerateTestCurve_HighlightsCompress(t *testing.T) {
	lut, err := GenerateTestCurve(CurveHighlightsCompress)
	if err != nil {
		t.Fatalf("GenerateTestCurve(HighlightsCompress) failed: %v", err)
	}

	linearLut, _ := GenerateTestCurve(CurveLinear)

	// Highlights compress should pull down highlight values
	if lut[224] >= linearLut[224] {
		t.Errorf("Highlights compress should reduce highlights: lut[224]=%d should be < linear[224]=%d",
			lut[224], linearLut[224])
	}

	// Shadow values should be closer to linear
	highlightDiff := int(linearLut[224]) - int(lut[224])
	shadowDiff := int(linearLut[32]) - int(lut[32])
	if highlightDiff <= shadowDiff {
		t.Errorf("Highlights should be affected more than shadows: highlightDiff=%d, shadowDiff=%d",
			highlightDiff, shadowDiff)
	}

	t.Logf("Highlight compression: -%d at entry 224, -%d at entry 32", highlightDiff, shadowDiff)
}

func TestGenerateTestCurve_Parametric(t *testing.T) {
	// Test shadows+darks boost
	lut, err := GenerateTestCurve(CurveShadowsDarks)
	if err != nil {
		t.Fatalf("GenerateTestCurve(ShadowsDarks) failed: %v", err)
	}

	linearLut, _ := GenerateTestCurve(CurveLinear)

	// Should lift values in shadows and darks zones
	if lut[32] <= linearLut[32] {
		t.Errorf("ShadowsDarks should lift shadows: lut[32]=%d should be > linear[32]=%d",
			lut[32], linearLut[32])
	}

	// Test lights+highlights compression
	lut2, err := GenerateTestCurve(CurveLightsHighlights)
	if err != nil {
		t.Fatalf("GenerateTestCurve(LightsHighlights) failed: %v", err)
	}

	// Should compress values in lights and highlights zones
	if lut2[224] >= linearLut[224] {
		t.Errorf("LightsHighlights should compress highlights: lut[224]=%d should be < linear[224]=%d",
			lut2[224], linearLut[224])
	}
}

func TestGenerateTestCurve_InvalidType(t *testing.T) {
	_, err := GenerateTestCurve("invalid-curve-type")
	if err == nil {
		t.Error("Expected error for invalid curve type, got nil")
	}
}

func TestGenerateTestNP3WithCurve(t *testing.T) {
	testCases := []struct {
		name      string
		curveType CurveType
	}{
		{"TestLinear", CurveLinear},
		{"TestSCurve", CurveSCurve},
		{"TestShadowsBoost", CurveShadowsBoost},
		{"TestHighlightsCompress", CurveHighlightsCompress},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := GenerateTestNP3WithCurve(tc.name, tc.curveType)
			if err != nil {
				t.Fatalf("GenerateTestNP3WithCurve failed: %v", err)
			}

			// Verify minimum file size (should be 1072 bytes for extended curve LUT)
			if len(data) < 1072 {
				t.Errorf("Expected at least 1072 bytes, got %d", len(data))
			}

			// Verify magic bytes
			if string(data[0:3]) != "NCP" {
				t.Errorf("Expected magic bytes 'NCP', got '%s'", string(data[0:3]))
			}

			// Verify tone curve data is present at offset 560 (OffsetExtendedToneCurveLUT)
			// Check that it's not all zeros
			allZeros := true
			for i := 560; i < 560+512 && i < len(data); i++ {
				if data[i] != 0 {
					allZeros = false
					break
				}
			}
			if allZeros {
				t.Error("Tone curve data at offset 560 is all zeros")
			}
		})
	}
}

// TestGenerateAndSaveTestNP3Files generates actual NP3 test files for manual NX Studio validation.
// This test creates files in the testdata directory that can be opened in NX Studio.
func TestGenerateAndSaveTestNP3Files(t *testing.T) {
	// Only run this test explicitly (not as part of normal test suite)
	if os.Getenv("GENERATE_TEST_FILES") != "1" {
		t.Skip("Skipping file generation test. Set GENERATE_TEST_FILES=1 to run.")
	}

	outputDir := filepath.Join("testdata", "curve_tests")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	testCases := []struct {
		filename  string
		curveType CurveType
	}{
		{"test_linear.np3", CurveLinear},
		{"test_scurve.np3", CurveSCurve},
		{"test_shadows_plus20.np3", CurveShadowsBoost},
		{"test_highlights_minus20.np3", CurveHighlightsCompress},
		{"test_shadows_darks_boost.np3", CurveShadowsDarks},
		{"test_lights_highlights_compress.np3", CurveLightsHighlights},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			data, err := GenerateTestNP3WithCurve(tc.filename[:len(tc.filename)-4], tc.curveType)
			if err != nil {
				t.Fatalf("Failed to generate NP3: %v", err)
			}

			outputPath := filepath.Join(outputDir, tc.filename)
			if err := os.WriteFile(outputPath, data, 0644); err != nil {
				t.Fatalf("Failed to write file: %v", err)
			}

			t.Logf("Generated: %s (%d bytes)", outputPath, len(data))
		})
	}

	t.Logf("Generated %d test NP3 files in %s", len(testCases), outputDir)
}
