package dcp

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestGenerate_ValidRecipe tests DCP generation from a populated UniversalRecipe.
func TestGenerate_ValidRecipe(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure:   0.5,  // +0.5 EV
		Contrast:   30,   // +30
		Highlights: -20,  // -20
		Shadows:    10,   // +10
		Metadata: map[string]interface{}{
			"profile_name": "Portrait",
		},
	}

	dcpData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	if len(dcpData) == 0 {
		t.Fatal("Generate() returned empty data")
	}

	// Validate DNG magic bytes ("IIRC" = 0x49 0x49 0x52 0x43)
	if len(dcpData) < 4 {
		t.Fatal("DCP data too small to contain magic bytes")
	}
	if dcpData[0] != 0x49 || dcpData[1] != 0x49 {
		t.Errorf("Invalid TIFF byte order: got %#x %#x, expected II (0x49 0x49)", dcpData[0], dcpData[1])
	}
	if dcpData[2] != 0x52 || dcpData[3] != 0x43 {
		t.Errorf("Invalid DNG version: got %#x %#x, expected RC (0x52 0x43)", dcpData[2], dcpData[3])
	}

	// Validate IFD offset is reasonable (should be 8 for minimal DCP)
	ifdOffset := binary.LittleEndian.Uint32(dcpData[4:8])
	if ifdOffset != 8 {
		t.Errorf("IFD offset = %d, expected 8", ifdOffset)
	}
}

// TestGenerate_EmptyRecipe tests generation from a neutral recipe (all zeros).
func TestGenerate_EmptyRecipe(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure:   0,
		Contrast:   0,
		Highlights: 0,
		Shadows:    0,
	}

	dcpData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	if len(dcpData) == 0 {
		t.Fatal("Generate() returned empty data")
	}

	// Verify DNG magic bytes
	if dcpData[0] != 0x49 || dcpData[1] != 0x49 || dcpData[2] != 0x52 || dcpData[3] != 0x43 {
		t.Errorf("Invalid DNG magic bytes: %#x %#x %#x %#x", dcpData[0], dcpData[1], dcpData[2], dcpData[3])
	}
}

// TestGenerate_ExtremeValues tests generation with extreme parameters.
func TestGenerate_ExtremeValues(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure:   2.0,  // +2.0 EV (max)
		Contrast:   100,  // +100 (max)
		Highlights: 100,  // +100 (max)
		Shadows:    -100, // -100 (min)
	}

	dcpData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() with extreme values failed: %v", err)
	}

	if len(dcpData) == 0 {
		t.Fatal("Generate() returned empty data")
	}
}

// TestGenerate_NilRecipe tests that Generate returns error for nil recipe.
func TestGenerate_NilRecipe(t *testing.T) {
	_, err := Generate(nil)
	if err == nil {
		t.Fatal("Generate(nil) should return error")
	}
	if err.Error() != "recipe cannot be nil" {
		t.Errorf("Generate(nil) error = %q, expected %q", err.Error(), "recipe cannot be nil")
	}
}

// TestUniversalToToneCurve tests tone curve generation formulas.
func TestUniversalToToneCurve(t *testing.T) {
	tests := []struct {
		name     string
		recipe   *models.UniversalRecipe
		expected []ToneCurvePoint // Expected outputs (approximately)
	}{
		{
			name: "neutral (all zeros)",
			recipe: &models.UniversalRecipe{
				Exposure:   0,
				Contrast:   0,
				Highlights: 0,
				Shadows:    0,
			},
			expected: []ToneCurvePoint{
				{Input: 0.0, Output: 0.0},
				{Input: 0.25, Output: 0.25},
				{Input: 0.5, Output: 0.5},
				{Input: 0.75, Output: 0.75},
				{Input: 1.0, Output: 1.0},
			},
		},
		{
			name: "exposure +1.0",
			recipe: &models.UniversalRecipe{
				Exposure: 1.0, // Exposure is applied in contrast calculation: 0.5 + deviation + exposure/5.0
			},
			expected: []ToneCurvePoint{
				{Input: 0.0, Output: 0.2},   // 0.5 + (-0.5)*1.0 + 0.2 = 0.2
				{Input: 0.25, Output: 0.45}, // 0.5 + (-0.25)*1.0 + 0.2 = 0.45
				{Input: 0.5, Output: 0.7},   // 0.5 + 0*1.0 + 0.2 = 0.7
				{Input: 0.75, Output: 0.95}, // 0.5 + 0.25*1.0 + 0.2 = 0.95
				{Input: 1.0, Output: 1.0},   // 0.5 + 0.5*1.0 + 0.2 = 1.2 → clamped to 1.0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := universalToToneCurve(tt.recipe)

			if len(points) != 5 {
				t.Fatalf("universalToToneCurve() returned %d points, expected 5", len(points))
			}

			// Verify all outputs are in 0.0-1.0 range
			for i, pt := range points {
				if pt.Output < 0.0 || pt.Output > 1.0 {
					t.Errorf("Point %d output = %.3f, out of range [0.0, 1.0]", i, pt.Output)
				}
			}

			// Verify monotonic property (output[i] >= output[i-1])
			for i := 1; i < len(points); i++ {
				if points[i].Output < points[i-1].Output {
					t.Errorf("Non-monotonic curve: point %d output (%.3f) < point %d output (%.3f)",
						i, points[i].Output, i-1, points[i-1].Output)
				}
			}

			// Verify expected outputs (approximate match ±0.1)
			for i, expected := range tt.expected {
				diff := points[i].Output - expected.Output
				if diff < -0.1 || diff > 0.1 {
					t.Errorf("Point %d: output = %.3f, expected ~%.3f (diff: %.3f)",
						i, points[i].Output, expected.Output, diff)
				}
			}
		})
	}
}

// TestGenerateBinaryToneCurve tests conversion of ToneCurvePoint to binary format.
func TestGenerateBinaryToneCurve(t *testing.T) {
	points := []ToneCurvePoint{
		{Input: 0.0, Output: 0.0},
		{Input: 1.0, Output: 1.0},
	}

	binaryData := generateBinaryToneCurve(points)

	// Each point is 8 bytes (2 float32 values)
	expectedLen := 2 * 8
	if len(binaryData) != expectedLen {
		t.Fatalf("generateBinaryToneCurve() returned %d bytes, expected %d", len(binaryData), expectedLen)
	}

	// Verify binary format (little-endian float32)
	// First point: input=0.0, output=0.0
	input0 := binary.LittleEndian.Uint32(binaryData[0:4])
	output0 := binary.LittleEndian.Uint32(binaryData[4:8])

	// 0.0 as float32 = 0x00000000
	if input0 != 0 || output0 != 0 {
		t.Errorf("First point: input bits = %#x, output bits = %#x, expected both 0x00000000", input0, output0)
	}

	// Second point: input=1.0, output=1.0
	// 1.0 as float32 = 0x3F800000 (IEEE 754)
	input1 := binary.LittleEndian.Uint32(binaryData[8:12])
	output1 := binary.LittleEndian.Uint32(binaryData[12:16])

	expected := uint32(0x3F800000)
	if input1 != expected || output1 != expected {
		t.Errorf("Second point: input bits = %#x, output bits = %#x, expected both 0x3F800000", input1, output1)
	}
}

// TestGenerateColorMatrix tests calibrated Nikon Z f matrix generation (ColorMatrix1).
func TestGenerateColorMatrix(t *testing.T) {
	matrix := generateColorMatrix()

	// ColorMatrix has 9 SRationals
	if len(matrix) != 9 {
		t.Fatalf("generateColorMatrix() returned %d SRationals, expected 9", len(matrix))
	}

	// Expected Nikon Z f ColorMatrix1 (Standard Light A):
	// [1.3904, -0.7947,  0.0654]
	// [-0.432,  1.2105,  0.2497]
	// [-0.0235,  0.083,  0.9243]
	expected := []SRational{
		{13904, 10000}, {-7947, 10000}, {654, 10000},
		{-4320, 10000}, {12105, 10000}, {2497, 10000},
		{-235, 10000}, {830, 10000}, {9243, 10000},
	}

	for i, sr := range matrix {
		if sr.Numerator != expected[i].Numerator || sr.Denominator != expected[i].Denominator {
			t.Errorf("Matrix[%d] = {%d, %d}, expected {%d, %d}",
				i, sr.Numerator, sr.Denominator, expected[i].Numerator, expected[i].Denominator)
		}
	}
}

// TestSRationalArrayToBytes tests SRational to binary conversion.
func TestSRationalArrayToBytes(t *testing.T) {
	srs := []SRational{
		{Numerator: 1, Denominator: 1},   // 1.0
		{Numerator: 0, Denominator: 1},   // 0.0
		{Numerator: -2, Denominator: 3},  // -0.666...
	}

	binaryData := srationalArrayToBytes(srs)

	// Each SRational is 8 bytes (int32 numerator + int32 denominator)
	expectedLen := 3 * 8
	if len(binaryData) != expectedLen {
		t.Fatalf("srationalArrayToBytes() returned %d bytes, expected %d", len(binaryData), expectedLen)
	}

	// Verify first SRational: {1, 1}
	num0 := int32(binary.LittleEndian.Uint32(binaryData[0:4]))
	denom0 := int32(binary.LittleEndian.Uint32(binaryData[4:8]))
	if num0 != 1 || denom0 != 1 {
		t.Errorf("SRational[0] = {%d, %d}, expected {1, 1}", num0, denom0)
	}

	// Verify third SRational: {-2, 3}
	num2 := int32(binary.LittleEndian.Uint32(binaryData[16:20]))
	denom2 := int32(binary.LittleEndian.Uint32(binaryData[20:24]))
	if num2 != -2 || denom2 != 3 {
		t.Errorf("SRational[2] = {%d, %d}, expected {-2, 3}", num2, denom2)
	}
}

// TestDNGStructure tests that generated DCP has valid DNG structure.
func TestDNGStructure(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure: 0.5,
		Metadata: map[string]interface{}{
			"profile_name": "Test Profile",
		},
	}

	dcpData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Validate DNG magic bytes ("IIRC" for little-endian DNG)
	if len(dcpData) < 8 {
		t.Fatal("DCP data too small to validate")
	}

	// Bytes 0-1: "II" (little-endian)
	if !bytes.Equal(dcpData[0:2], []byte{0x49, 0x49}) {
		t.Errorf("Invalid byte order: got %#x %#x, expected II (0x49 0x49)", dcpData[0], dcpData[1])
	}

	// Bytes 2-3: "RC" (DNG version instead of TIFF 42)
	if !bytes.Equal(dcpData[2:4], []byte{0x52, 0x43}) {
		t.Errorf("Invalid DNG version: got %#x %#x, expected RC (0x52 0x43)", dcpData[2], dcpData[3])
	}

	// Bytes 4-7: IFD offset (should be 8 = right after header)
	ifdOffset := binary.LittleEndian.Uint32(dcpData[4:8])
	if ifdOffset != 8 {
		t.Errorf("IFD offset = %d, expected 8", ifdOffset)
	}

	// Validate IFD structure
	// IFD starts with 2-byte entry count
	if len(dcpData) < int(ifdOffset)+2 {
		t.Fatal("DCP data too small to contain IFD")
	}

	entryCount := binary.LittleEndian.Uint16(dcpData[ifdOffset : ifdOffset+2])
	if entryCount == 0 {
		t.Error("IFD has 0 entries")
	}

	// Should have at least 10 tags (6 standard TIFF + 4 required DCP tags)
	if entryCount < 10 {
		t.Errorf("IFD entry count = %d, expected at least 10 tags", entryCount)
	}
}

// TestRoundTrip_DCP tests Generate → Parse → Compare accuracy.
func TestRoundTrip_DCP(t *testing.T) {
	original := &models.UniversalRecipe{
		Exposure:   0.5,
		Contrast:   30,
		Highlights: -20,
		Shadows:    10,
		Metadata: map[string]interface{}{
			"profile_name": "Round Trip Test",
		},
	}

	// Generate DCP
	dcpData, err := Generate(original)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Parse DCP
	parsed, err := Parse(dcpData)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// Compare exposure (allow ±0.1 due to float32 precision)
	exposureDiff := parsed.Exposure - original.Exposure
	if exposureDiff < -0.1 || exposureDiff > 0.1 {
		t.Errorf("Exposure: parsed = %.3f, original = %.3f (diff: %.3f)",
			parsed.Exposure, original.Exposure, exposureDiff)
	}

	// Compare contrast (allow ±10 due to curve analysis approximation)
	// The round-trip involves: Generate (tone curve) → Parse (curve analysis)
	// Curve analysis is approximate since it extracts parameters from curve shape
	contrastDiff := parsed.Contrast - original.Contrast
	if contrastDiff < -10 || contrastDiff > 10 {
		t.Errorf("Contrast: parsed = %d, original = %d (diff: %d)",
			parsed.Contrast, original.Contrast, contrastDiff)
	}

	// Compare metadata (profile name)
	parsedName, ok := parsed.Metadata["profile_name"].(string)
	if !ok {
		t.Error("Parsed metadata missing profile_name")
	}
	originalName, _ := original.Metadata["profile_name"].(string)
	if parsedName != originalName {
		t.Errorf("Profile name: parsed = %q, original = %q", parsedName, originalName)
	}
}

// TestOptionalProfileName tests DCP generation with and without profile name.
func TestOptionalProfileName(t *testing.T) {
	// Test without profile name
	t.Run("without profile name", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			Exposure: 0.5,
			// No profile_name in metadata
		}

		dcpData, err := Generate(recipe)
		if err != nil {
			t.Fatalf("Generate() failed: %v", err)
		}

		// Should still generate valid DCP
		if len(dcpData) < 8 {
			t.Fatal("DCP data too small")
		}
		if dcpData[0] != 0x49 || dcpData[1] != 0x49 || dcpData[2] != 0x52 || dcpData[3] != 0x43 {
			t.Error("Invalid DNG magic bytes")
		}
	})

	// Test with profile name
	t.Run("with profile name", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			Exposure: 0.5,
			Metadata: map[string]interface{}{
				"profile_name": "Named Profile",
			},
		}

		dcpData, err := Generate(recipe)
		if err != nil {
			t.Fatalf("Generate() failed: %v", err)
		}

		// Parse and verify profile name
		parsed, err := Parse(dcpData)
		if err != nil {
			t.Fatalf("Parse() failed: %v", err)
		}

		parsedName, ok := parsed.Metadata["profile_name"].(string)
		if !ok || parsedName != "Named Profile" {
			t.Errorf("Profile name: parsed = %q, expected %q", parsedName, "Named Profile")
		}
	})
}

// ============================================================================
// PERFORMANCE BENCHMARKS (AC-4: Story 9-4 DCP Compatibility Validation)
// ============================================================================

// BenchmarkGenerate_Neutral benchmarks DCP generation with neutral preset (all zeros).
func BenchmarkGenerate_Neutral(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   0,
		Contrast:   0,
		Highlights: 0,
		Shadows:    0,
		Metadata: map[string]interface{}{
			"profile_name": "Neutral",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_Portrait benchmarks DCP generation with portrait preset.
func BenchmarkGenerate_Portrait(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   0.5,  // +0.5 stops
		Contrast:   30,   // +30
		Highlights: -20,  // -20
		Shadows:    0,    // 0
		Metadata: map[string]interface{}{
			"profile_name": "Portrait",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_Landscape benchmarks DCP generation with landscape preset.
func BenchmarkGenerate_Landscape(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   0.3,  // +0.3 stops
		Contrast:   0,    // 0
		Highlights: 0,    // 0
		Shadows:    20,   // +20
		Metadata: map[string]interface{}{
			"profile_name": "Landscape",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_HighContrast benchmarks DCP generation with high contrast.
func BenchmarkGenerate_HighContrast(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   0,   // 0
		Contrast:   100, // +100 (max)
		Highlights: 0,   // 0
		Shadows:    0,   // 0
		Metadata: map[string]interface{}{
			"profile_name": "High Contrast",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_HighExposure benchmarks DCP generation with high exposure.
func BenchmarkGenerate_HighExposure(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   2.0, // +2.0 stops (max)
		Contrast:   0,   // 0
		Highlights: 0,   // 0
		Shadows:    0,   // 0
		Metadata: map[string]interface{}{
			"profile_name": "High Exposure",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_LowExposure benchmarks DCP generation with low exposure.
func BenchmarkGenerate_LowExposure(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   -2.0, // -2.0 stops (min)
		Contrast:   0,    // 0
		Highlights: 0,    // 0
		Shadows:    0,    // 0
		Metadata: map[string]interface{}{
			"profile_name": "Low Exposure",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_MixedAdjustments benchmarks DCP generation with mixed parameters.
func BenchmarkGenerate_MixedAdjustments(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   0.8,  // +0.8 stops
		Contrast:   20,   // +20
		Highlights: -30,  // -30
		Shadows:    15,   // +15
		Metadata: map[string]interface{}{
			"profile_name": "Mixed Adjustments",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_ExtremeValues benchmarks DCP generation with extreme clamping.
func BenchmarkGenerate_ExtremeValues(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   2.0,  // +2.0 stops (max)
		Contrast:   100,  // +100 (max)
		Highlights: 100,  // +100 (max)
		Shadows:    -100, // -100 (min)
		Metadata: map[string]interface{}{
			"profile_name": "Extreme Values",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_MinimalParameters benchmarks DCP generation with exposure only.
func BenchmarkGenerate_MinimalParameters(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure: 0.5, // +0.5 stops only
		Metadata: map[string]interface{}{
			"profile_name": "Minimal Parameters",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerate_ComplexPreset benchmarks DCP generation with all parameters.
func BenchmarkGenerate_ComplexPreset(b *testing.B) {
	recipe := &models.UniversalRecipe{
		Exposure:   0.7,  // +0.7 stops
		Contrast:   40,   // +40
		Highlights: -25,  // -25
		Shadows:    30,   // +30
		Whites:     10,   // +10 (if supported)
		Blacks:     -5,   // -5 (if supported)
		Clarity:    20,   // +20 (if supported)
		Saturation: 15,   // +15 (if supported)
		Metadata: map[string]interface{}{
			"profile_name": "Complex Preset",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}
