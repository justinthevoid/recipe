package models

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

// TestUniversalRecipeJSONMarshaling tests JSON serialization.
func TestUniversalRecipeJSONMarshaling(t *testing.T) {
	recipe := UniversalRecipe{
		Name:         "Test Recipe",
		SourceFormat: "np3",
		Exposure:     2.5,
		Contrast:     50,
		Highlights:   -25,
		Shadows:      30,
	}

	// Marshal to JSON
	data, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	// Unmarshal from JSON
	var decoded UniversalRecipe
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	// Verify fields
	if decoded.Name != recipe.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, recipe.Name)
	}
	if decoded.Exposure != recipe.Exposure {
		t.Errorf("Exposure mismatch: got %v, want %v", decoded.Exposure, recipe.Exposure)
	}
	if decoded.Contrast != recipe.Contrast {
		t.Errorf("Contrast mismatch: got %d, want %d", decoded.Contrast, recipe.Contrast)
	}
}

// TestUniversalRecipeXMLMarshaling tests XML serialization.
func TestUniversalRecipeXMLMarshaling(t *testing.T) {
	recipe := UniversalRecipe{
		Name:         "Test Recipe",
		SourceFormat: "xmp",
		Exposure:     1.5,
		Contrast:     25,
		Vibrance:     15,
	}

	// Marshal to XML
	data, err := xml.Marshal(recipe)
	if err != nil {
		t.Fatalf("Failed to marshal to XML: %v", err)
	}

	// Unmarshal from XML
	var decoded UniversalRecipe
	err = xml.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal from XML: %v", err)
	}

	// Verify fields
	if decoded.Name != recipe.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, recipe.Name)
	}
	if decoded.Exposure != recipe.Exposure {
		t.Errorf("Exposure mismatch: got %v, want %v", decoded.Exposure, recipe.Exposure)
	}
	if decoded.Vibrance != recipe.Vibrance {
		t.Errorf("Vibrance mismatch: got %d, want %d", decoded.Vibrance, recipe.Vibrance)
	}
}

// TestJSONRoundTrip tests round-trip JSON encoding/decoding.
func TestJSONRoundTrip(t *testing.T) {
	original := UniversalRecipe{
		Name:       "Round Trip Test",
		Exposure:   -2.0,
		Contrast:   -50,
		Highlights: 75,
		Shadows:    -30,
		Red:        ColorAdjustment{Hue: 10, Saturation: 20, Luminance: -15},
		PointCurve: []ToneCurvePoint{
			{Input: 0, Output: 0},
			{Input: 128, Output: 140},
			{Input: 255, Output: 255},
		},
	}

	// Encode
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Decode
	var decoded UniversalRecipe
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Compare
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch after round trip")
	}
	if decoded.Exposure != original.Exposure {
		t.Errorf("Exposure mismatch after round trip")
	}
	if decoded.Red.Hue != original.Red.Hue {
		t.Errorf("Red hue mismatch after round trip")
	}
	if len(decoded.PointCurve) != len(original.PointCurve) {
		t.Errorf("Point curve length mismatch: got %d, want %d", len(decoded.PointCurve), len(original.PointCurve))
	}
}

// TestXMLRoundTrip tests round-trip XML encoding/decoding.
func TestXMLRoundTrip(t *testing.T) {
	original := UniversalRecipe{
		Name:              "XML Round Trip",
		Exposure:          3.5,
		Temperature:       intPtr(5500),
		Tint:              10,
		SplitShadowHue:    240,
		SplitShadowSaturation: 50,
	}

	// Encode
	data, err := xml.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Decode
	var decoded UniversalRecipe
	err = xml.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Compare
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch after round trip")
	}
	if decoded.Exposure != original.Exposure {
		t.Errorf("Exposure mismatch after round trip")
	}
	if *decoded.Temperature != *original.Temperature {
		t.Errorf("Temperature mismatch after round trip")
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}

// TestValidateExposure tests exposure validation.
func TestValidateExposure(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{"Valid minimum", -5.0, false},
		{"Valid maximum", 5.0, false},
		{"Valid zero", 0.0, false},
		{"Valid positive", 2.5, false},
		{"Valid negative", -3.2, false},
		{"Invalid too low", -5.1, true},
		{"Invalid too high", 5.1, true},
		{"Invalid far too low", -10.0, true},
		{"Invalid far too high", 10.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExposure(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExposure(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

// TestValidatePercentage tests percentage validation.
func TestValidatePercentage(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		fieldName string
		wantErr   bool
	}{
		{"Valid minimum", -100, "test", false},
		{"Valid maximum", 100, "test", false},
		{"Valid zero", 0, "test", false},
		{"Valid positive", 50, "test", false},
		{"Valid negative", -75, "test", false},
		{"Invalid too low", -101, "test", true},
		{"Invalid too high", 101, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePercentage(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePercentage(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

// TestValidateHSL tests HSL validation.
func TestValidateHSL(t *testing.T) {
	tests := []struct {
		name       string
		hue        int
		saturation int
		luminance  int
		colorName  string
		wantErr    bool
	}{
		{"Valid all zeros", 0, 0, 0, "red", false},
		{"Valid max hue", 180, 0, 0, "red", false},
		{"Valid min hue", -180, 0, 0, "red", false},
		{"Valid max sat/lum", 0, 100, 100, "red", false},
		{"Valid min sat/lum", 0, -100, -100, "red", false},
		{"Invalid hue too high", 181, 0, 0, "red", true},
		{"Invalid hue too low", -181, 0, 0, "red", true},
		{"Invalid saturation too high", 0, 101, 0, "red", true},
		{"Invalid luminance too low", 0, 0, -101, "red", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHSL(tt.hue, tt.saturation, tt.luminance, tt.colorName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHSL(%d,%d,%d) error = %v, wantErr %v", tt.hue, tt.saturation, tt.luminance, err, tt.wantErr)
			}
		})
	}
}

// TestValidateNP3Ranges tests NP3-specific validation functions.
func TestValidateNP3Ranges(t *testing.T) {
	t.Run("Sharpening", func(t *testing.T) {
		tests := []struct {
			value   int
			wantErr bool
		}{
			{0, false},
			{5, false},
			{9, false},
			{-1, true},
			{10, true},
		}
		for _, tt := range tests {
			err := ValidateNP3Sharpening(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNP3Sharpening(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})

	t.Run("Contrast", func(t *testing.T) {
		tests := []struct {
			value   int
			wantErr bool
		}{
			{-3, false},
			{0, false},
			{3, false},
			{-4, true},
			{4, true},
		}
		for _, tt := range tests {
			err := ValidateNP3Contrast(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNP3Contrast(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})

	t.Run("Brightness", func(t *testing.T) {
		tests := []struct {
			value   float64
			wantErr bool
		}{
			{-1.0, false},
			{0.0, false},
			{1.0, false},
			{-1.1, true},
			{1.1, true},
		}
		for _, tt := range tests {
			err := ValidateNP3Brightness(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNP3Brightness(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})

	t.Run("Saturation", func(t *testing.T) {
		tests := []struct {
			value   int
			wantErr bool
		}{
			{-3, false},
			{0, false},
			{3, false},
			{-4, true},
			{4, true},
		}
		for _, tt := range tests {
			err := ValidateNP3Saturation(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNP3Saturation(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})

	t.Run("Hue", func(t *testing.T) {
		tests := []struct {
			value   int
			wantErr bool
		}{
			{-9, false},
			{0, false},
			{9, false},
			{-10, true},
			{10, true},
		}
		for _, tt := range tests {
			err := ValidateNP3Hue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNP3Hue(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})
}

// TestBuilderBasicUsage tests basic builder usage with method chaining.
func TestBuilderBasicUsage(t *testing.T) {
	recipe, err := NewRecipeBuilder().
		WithName("Test Recipe").
		WithSourceFormat("np3").
		WithExposure(2.0).
		WithContrast(50).
		WithHighlights(-25).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if recipe.Name != "Test Recipe" {
		t.Errorf("Name mismatch: got %s, want Test Recipe", recipe.Name)
	}
	if recipe.Exposure != 2.0 {
		t.Errorf("Exposure mismatch: got %v, want 2.0", recipe.Exposure)
	}
	if recipe.Contrast != 50 {
		t.Errorf("Contrast mismatch: got %d, want 50", recipe.Contrast)
	}
}

// TestBuilderValidation tests that builder validates parameters.
func TestBuilderValidation(t *testing.T) {
	_, err := NewRecipeBuilder().
		WithExposure(10.0). // Invalid: out of range
		Build()

	if err == nil {
		t.Error("Expected validation error for invalid exposure, got nil")
	}
}

// TestBuilderHSL tests HSL color adjustments via builder.
func TestBuilderHSL(t *testing.T) {
	recipe, err := NewRecipeBuilder().
		WithRedHSL(10, 20, -15).
		WithBlueHSL(-30, 40, 25).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if recipe.Red.Hue != 10 || recipe.Red.Saturation != 20 || recipe.Red.Luminance != -15 {
		t.Errorf("Red HSL mismatch")
	}
	if recipe.Blue.Hue != -30 || recipe.Blue.Saturation != 40 || recipe.Blue.Luminance != 25 {
		t.Errorf("Blue HSL mismatch")
	}
}

// TestBuilderToneCurve tests tone curve settings via builder.
func TestBuilderToneCurve(t *testing.T) {
	recipe, err := NewRecipeBuilder().
		WithToneCurve(-20, -10, 10, 20).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if recipe.ToneCurveShadows != -20 {
		t.Errorf("ToneCurveShadows mismatch")
	}
	if recipe.ToneCurveHighlights != 20 {
		t.Errorf("ToneCurveHighlights mismatch")
	}
}

// TestBuilderPointCurve tests point curve validation.
func TestBuilderPointCurve(t *testing.T) {
	points := []ToneCurvePoint{
		{Input: 0, Output: 0},
		{Input: 128, Output: 140},
		{Input: 255, Output: 255},
	}

	recipe, err := NewRecipeBuilder().
		WithPointCurve(points).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(recipe.PointCurve) != 3 {
		t.Errorf("Point curve length mismatch: got %d, want 3", len(recipe.PointCurve))
	}
}

// TestBuilderInvalidPointCurve tests that invalid point curves are rejected.
func TestBuilderInvalidPointCurve(t *testing.T) {
	points := []ToneCurvePoint{
		{Input: 0, Output: 0},
		{Input: 300, Output: 140}, // Invalid: input > 255
	}

	_, err := NewRecipeBuilder().
		WithPointCurve(points).
		Build()

	if err == nil {
		t.Error("Expected validation error for invalid point curve, got nil")
	}
}

// TestBuilderSplitToning tests split toning via builder.
func TestBuilderSplitToning(t *testing.T) {
	recipe, err := NewRecipeBuilder().
		WithSplitToning(240, 50, 60, 30, -10).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if recipe.SplitShadowHue != 240 {
		t.Errorf("SplitShadowHue mismatch")
	}
	if recipe.SplitHighlightSaturation != 30 {
		t.Errorf("SplitHighlightSaturation mismatch")
	}
}

// TestBuilderGrain tests grain effects via builder.
func TestBuilderGrain(t *testing.T) {
	recipe, err := NewRecipeBuilder().
		WithGrain(50, 25, 75).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if recipe.GrainAmount != 50 || recipe.GrainSize != 25 || recipe.GrainRoughness != 75 {
		t.Errorf("Grain parameters mismatch")
	}
}

// TestBuilderVignette tests vignette via builder.
func TestBuilderVignette(t *testing.T) {
	recipe, err := NewRecipeBuilder().
		WithVignette(-50, 40, 30, 60).
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if recipe.VignetteAmount != -50 {
		t.Errorf("VignetteAmount mismatch")
	}
	if recipe.VignetteFeather != 60 {
		t.Errorf("VignetteFeather mismatch")
	}
}

// TestBuilderDefaults tests that builder initializes defaults.
func TestBuilderDefaults(t *testing.T) {
	recipe, err := NewRecipeBuilder().Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if recipe.SharpnessRadius != 1.0 {
		t.Errorf("Default SharpnessRadius mismatch: got %v, want 1.0", recipe.SharpnessRadius)
	}
	if recipe.ToneCurveShadowSplit != 25 {
		t.Errorf("Default ToneCurveShadowSplit mismatch: got %d, want 25", recipe.ToneCurveShadowSplit)
	}
	if recipe.ToneCurveMidtoneSplit != 50 {
		t.Errorf("Default ToneCurveMidtoneSplit mismatch: got %d, want 50", recipe.ToneCurveMidtoneSplit)
	}
}

// TestUniversalRecipeValidate tests the Validate method.
func TestUniversalRecipeValidate(t *testing.T) {
	t.Run("Valid recipe", func(t *testing.T) {
		recipe := UniversalRecipe{
			Exposure:   2.0,
			Contrast:   50,
			Highlights: -25,
			Red:        ColorAdjustment{Hue: 10, Saturation: 20, Luminance: -15},
		}

		err := recipe.Validate()
		if err != nil {
			t.Errorf("Validate failed for valid recipe: %v", err)
		}
	})

	t.Run("Invalid exposure", func(t *testing.T) {
		recipe := UniversalRecipe{
			Exposure: 10.0, // Invalid
		}

		err := recipe.Validate()
		if err == nil {
			t.Error("Expected validation error for invalid exposure")
		}
	})

	t.Run("Invalid HSL", func(t *testing.T) {
		recipe := UniversalRecipe{
			Red: ColorAdjustment{Hue: 200, Saturation: 20, Luminance: -15}, // Invalid hue
		}

		err := recipe.Validate()
		if err == nil {
			t.Error("Expected validation error for invalid HSL")
		}
	})
}

// TestEdgeCases tests edge cases.
func TestEdgeCases(t *testing.T) {
	t.Run("Zero values", func(t *testing.T) {
		recipe := UniversalRecipe{}
		err := recipe.Validate()
		if err != nil {
			t.Errorf("Validate failed for zero values: %v", err)
		}
	})

	t.Run("Min values", func(t *testing.T) {
		recipe := UniversalRecipe{
			Exposure:   -5.0,
			Contrast:   -100,
			Highlights: -100,
		}
		err := recipe.Validate()
		if err != nil {
			t.Errorf("Validate failed for min values: %v", err)
		}
	})

	t.Run("Max values", func(t *testing.T) {
		recipe := UniversalRecipe{
			Exposure:   5.0,
			Contrast:   100,
			Highlights: 100,
		}
		err := recipe.Validate()
		if err != nil {
			t.Errorf("Validate failed for max values: %v", err)
		}
	})

	t.Run("Nil Temperature", func(t *testing.T) {
		recipe := UniversalRecipe{
			Temperature: nil,
		}
		err := recipe.Validate()
		if err != nil {
			t.Errorf("Validate failed for nil temperature: %v", err)
		}
	})

	t.Run("Empty arrays", func(t *testing.T) {
		recipe := UniversalRecipe{
			PointCurve:     []ToneCurvePoint{},
			NP3ColorData:   []map[string]interface{}{},
			NP3RawParams:   []map[string]interface{}{},
		}
		err := recipe.Validate()
		if err != nil {
			t.Errorf("Validate failed for empty arrays: %v", err)
		}
	})
}

// TestComprehensiveValidation tests all validation paths.
func TestComprehensiveValidation(t *testing.T) {
	t.Run("All basic adjustments", func(t *testing.T) {
		recipe := UniversalRecipe{
			Exposure:   2.5,
			Contrast:   50,
			Highlights: -25,
			Shadows:    30,
			Whites:     40,
			Blacks:     -30,
			Texture:    20,
			Clarity:    15,
			Dehaze:     10,
			Vibrance:   25,
			Saturation: 35,
		}
		if err := recipe.Validate(); err != nil {
			t.Errorf("Validation failed: %v", err)
		}
	})

	t.Run("All sharpening parameters", func(t *testing.T) {
		recipe := UniversalRecipe{
			Sharpness:       75,
			SharpnessRadius: 1.5,
			SharpnessDetail: 50,
			SharpnessMasking: 30,
		}
		if err := recipe.Validate(); err != nil {
			t.Errorf("Validation failed: %v", err)
		}
	})

	t.Run("All tone curve parameters", func(t *testing.T) {
		recipe := UniversalRecipe{
			ToneCurveShadows:        -20,
			ToneCurveDarks:          -10,
			ToneCurveLights:         10,
			ToneCurveHighlights:     20,
			ToneCurveShadowSplit:    25,
			ToneCurveMidtoneSplit:   50,
			ToneCurveHighlightSplit: 75,
		}
		if err := recipe.Validate(); err != nil {
			t.Errorf("Validation failed: %v", err)
		}
	})

	t.Run("All color HSL adjustments", func(t *testing.T) {
		recipe := UniversalRecipe{
			Red:     ColorAdjustment{Hue: 10, Saturation: 20, Luminance: -15},
			Orange:  ColorAdjustment{Hue: -5, Saturation: 15, Luminance: 10},
			Yellow:  ColorAdjustment{Hue: 8, Saturation: -10, Luminance: 5},
			Green:   ColorAdjustment{Hue: -12, Saturation: 25, Luminance: -8},
			Aqua:    ColorAdjustment{Hue: 6, Saturation: -15, Luminance: 12},
			Blue:    ColorAdjustment{Hue: -20, Saturation: 30, Luminance: -20},
			Purple:  ColorAdjustment{Hue: 15, Saturation: -5, Luminance: 8},
			Magenta: ColorAdjustment{Hue: -8, Saturation: 18, Luminance: -12},
		}
		if err := recipe.Validate(); err != nil {
			t.Errorf("Validation failed: %v", err)
		}
	})

	t.Run("Point curves all channels", func(t *testing.T) {
		points := []ToneCurvePoint{{Input: 0, Output: 0}, {Input: 128, Output: 140}, {Input: 255, Output: 255}}
		recipe := UniversalRecipe{
			PointCurve:      points,
			PointCurveRed:   points,
			PointCurveGreen: points,
			PointCurveBlue:  points,
		}
		if err := recipe.Validate(); err != nil {
			t.Errorf("Validation failed: %v", err)
		}
	})

	t.Run("Split toning parameters", func(t *testing.T) {
		recipe := UniversalRecipe{
			SplitShadowHue:           240,
			SplitShadowSaturation:    50,
			SplitHighlightHue:        60,
			SplitHighlightSaturation: 30,
			SplitBalance:             -10,
		}
		if err := recipe.Validate(); err != nil {
			t.Errorf("Validation failed: %v", err)
		}
	})

	t.Run("Grain and vignette effects", func(t *testing.T) {
		recipe := UniversalRecipe{
			GrainAmount:      50,
			GrainSize:        25,
			GrainRoughness:   75,
			VignetteAmount:   -50,
			VignetteMidpoint: 40,
			VignetteRoundness: 30,
			VignetteFeather:  60,
		}
		if err := recipe.Validate(); err != nil {
			t.Errorf("Validation failed: %v", err)
		}
	})

	t.Run("Invalid sharpness radius", func(t *testing.T) {
		recipe := UniversalRecipe{
			SharpnessRadius: 5.0, // Invalid
		}
		if err := recipe.Validate(); err == nil {
			t.Error("Expected validation error for invalid sharpness radius")
		}
	})

	t.Run("Invalid tint", func(t *testing.T) {
		recipe := UniversalRecipe{
			Tint: 200, // Invalid
		}
		if err := recipe.Validate(); err == nil {
			t.Error("Expected validation error for invalid tint")
		}
	})

	t.Run("Invalid point curve input", func(t *testing.T) {
		recipe := UniversalRecipe{
			PointCurveRed: []ToneCurvePoint{{Input: 300, Output: 100}}, // Invalid input
		}
		if err := recipe.Validate(); err == nil {
			t.Error("Expected validation error for invalid point curve")
		}
	})

	t.Run("Invalid point curve output", func(t *testing.T) {
		recipe := UniversalRecipe{
			PointCurveGreen: []ToneCurvePoint{{Input: 100, Output: 300}}, // Invalid output
		}
		if err := recipe.Validate(); err == nil {
			t.Error("Expected validation error for invalid point curve")
		}
	})

	t.Run("Invalid split toning hue", func(t *testing.T) {
		recipe := UniversalRecipe{
			SplitShadowHue: 400, // Invalid
		}
		if err := recipe.Validate(); err == nil {
			t.Error("Expected validation error for invalid split shadow hue")
		}
	})

	t.Run("Invalid split toning saturation", func(t *testing.T) {
		recipe := UniversalRecipe{
			SplitHighlightSaturation: 150, // Invalid
		}
		if err := recipe.Validate(); err == nil {
			t.Error("Expected validation error for invalid split highlight saturation")
		}
	})

	t.Run("Invalid grain amount", func(t *testing.T) {
		recipe := UniversalRecipe{
			GrainAmount: 150, // Invalid
		}
		if err := recipe.Validate(); err == nil {
			t.Error("Expected validation error for invalid grain amount")
		}
	})

	t.Run("Invalid vignette amount", func(t *testing.T) {
		recipe := UniversalRecipe{
			VignetteAmount: 150, // Invalid
		}
		if err := recipe.Validate(); err == nil {
			t.Error("Expected validation error for invalid vignette amount")
		}
	})
}

// TestBuilderAllSetters tests all builder setter methods.
func TestBuilderAllSetters(t *testing.T) {
	t.Run("All color HSL setters", func(t *testing.T) {
		recipe, err := NewRecipeBuilder().
			WithRedHSL(10, 20, -15).
			WithOrangeHSL(-5, 15, 10).
			WithYellowHSL(8, -10, 5).
			WithGreenHSL(-12, 25, -8).
			WithAquaHSL(6, -15, 12).
			WithBlueHSL(-20, 30, -20).
			WithPurpleHSL(15, -5, 8).
			WithMagentaHSL(-8, 18, -12).
			Build()

		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		if recipe.Orange.Saturation != 15 {
			t.Errorf("Orange saturation mismatch")
		}
		if recipe.Purple.Hue != 15 {
			t.Errorf("Purple hue mismatch")
		}
	})

	t.Run("All basic adjustment setters", func(t *testing.T) {
		recipe, err := NewRecipeBuilder().
			WithExposure(2.5).
			WithContrast(50).
			WithHighlights(-25).
			WithShadows(30).
			WithWhites(40).
			WithBlacks(-30).
			WithTexture(20).
			WithClarity(15).
			WithDehaze(10).
			WithVibrance(25).
			WithSaturation(35).
			Build()

		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		if recipe.Whites != 40 {
			t.Errorf("Whites mismatch")
		}
		if recipe.Texture != 20 {
			t.Errorf("Texture mismatch")
		}
	})

	t.Run("All sharpening setters", func(t *testing.T) {
		recipe, err := NewRecipeBuilder().
			WithSharpness(75).
			WithSharpnessRadius(1.5).
			WithSharpnessDetail(50).
			WithSharpnessMasking(30).
			Build()

		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		if recipe.Sharpness != 75 {
			t.Errorf("Sharpness mismatch")
		}
	})

	t.Run("Temperature and tint setters", func(t *testing.T) {
		recipe, err := NewRecipeBuilder().
			WithTemperature(5500).
			WithTint(10).
			Build()

		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		if *recipe.Temperature != 5500 {
			t.Errorf("Temperature mismatch")
		}
	})

	t.Run("Camera profile setter", func(t *testing.T) {
		profile := CameraProfile{
			RedHue:         5,
			RedSaturation:  10,
			GreenHue:       -3,
			GreenSaturation: 8,
			BlueHue:        2,
			BlueSaturation: -5,
		}

		recipe, err := NewRecipeBuilder().
			WithCameraProfile(profile).
			Build()

		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		if recipe.CameraProfile.RedHue != 5 {
			t.Errorf("CameraProfile RedHue mismatch")
		}
	})

	t.Run("Invalid builder parameters accumulate errors", func(t *testing.T) {
		_, err := NewRecipeBuilder().
			WithExposure(10.0). // Invalid
			WithContrast(200).  // Invalid
			WithHighlights(150). // Invalid
			Build()

		if err == nil {
			t.Error("Expected validation errors to accumulate")
		}
	})
}

// TestBuilderValidationErrors tests all builder validation error paths.
func TestBuilderValidationErrors(t *testing.T) {
	t.Run("WithExposure invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithExposure(6.0).Build()
		if err == nil {
			t.Error("Expected error for invalid exposure")
		}
	})

	t.Run("WithContrast invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithContrast(101).Build()
		if err == nil {
			t.Error("Expected error for invalid contrast")
		}
	})

	t.Run("WithHighlights invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithHighlights(-101).Build()
		if err == nil {
			t.Error("Expected error for invalid highlights")
		}
	})

	t.Run("WithShadows invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithShadows(101).Build()
		if err == nil {
			t.Error("Expected error for invalid shadows")
		}
	})

	t.Run("WithWhites invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithWhites(-101).Build()
		if err == nil {
			t.Error("Expected error for invalid whites")
		}
	})

	t.Run("WithBlacks invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithBlacks(101).Build()
		if err == nil {
			t.Error("Expected error for invalid blacks")
		}
	})

	t.Run("WithTexture invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithTexture(-101).Build()
		if err == nil {
			t.Error("Expected error for invalid texture")
		}
	})

	t.Run("WithClarity invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithClarity(101).Build()
		if err == nil {
			t.Error("Expected error for invalid clarity")
		}
	})

	t.Run("WithDehaze invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithDehaze(-101).Build()
		if err == nil {
			t.Error("Expected error for invalid dehaze")
		}
	})

	t.Run("WithVibrance invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithVibrance(101).Build()
		if err == nil {
			t.Error("Expected error for invalid vibrance")
		}
	})

	t.Run("WithSaturation invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSaturation(-101).Build()
		if err == nil {
			t.Error("Expected error for invalid saturation")
		}
	})

	t.Run("WithSharpnessRadius invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSharpnessRadius(3.5).Build()
		if err == nil {
			t.Error("Expected error for invalid sharpness radius")
		}
	})

	t.Run("WithSharpnessDetail invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSharpnessDetail(101).Build()
		if err == nil {
			t.Error("Expected error for invalid sharpness detail")
		}
	})

	t.Run("WithSharpnessMasking invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSharpnessMasking(-1).Build()
		if err == nil {
			t.Error("Expected error for invalid sharpness masking")
		}
	})

	t.Run("WithTint invalid", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithTint(200).Build()
		if err == nil {
			t.Error("Expected error for invalid tint")
		}
	})

	t.Run("WithRedHSL invalid hue", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithRedHSL(200, 0, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid red hue")
		}
	})

	t.Run("WithOrangeHSL invalid saturation", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithOrangeHSL(0, 101, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid orange saturation")
		}
	})

	t.Run("WithYellowHSL invalid luminance", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithYellowHSL(0, 0, -101).Build()
		if err == nil {
			t.Error("Expected error for invalid yellow luminance")
		}
	})

	t.Run("WithGreenHSL invalid hue", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithGreenHSL(-181, 0, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid green hue")
		}
	})

	t.Run("WithAquaHSL invalid saturation", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithAquaHSL(0, -101, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid aqua saturation")
		}
	})

	t.Run("WithBlueHSL invalid luminance", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithBlueHSL(0, 0, 101).Build()
		if err == nil {
			t.Error("Expected error for invalid blue luminance")
		}
	})

	t.Run("WithPurpleHSL invalid hue", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithPurpleHSL(181, 0, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid purple hue")
		}
	})

	t.Run("WithMagentaHSL invalid saturation", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithMagentaHSL(0, 101, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid magenta saturation")
		}
	})

	t.Run("WithToneCurve invalid shadows", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithToneCurve(-101, 0, 0, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid tone curve shadows")
		}
	})

	t.Run("WithToneCurve invalid darks", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithToneCurve(0, 101, 0, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid tone curve darks")
		}
	})

	t.Run("WithToneCurve invalid lights", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithToneCurve(0, 0, -101, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid tone curve lights")
		}
	})

	t.Run("WithToneCurve invalid highlights", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithToneCurve(0, 0, 0, 101).Build()
		if err == nil {
			t.Error("Expected error for invalid tone curve highlights")
		}
	})

	t.Run("WithPointCurve invalid input", func(t *testing.T) {
		points := []ToneCurvePoint{{Input: -1, Output: 128}}
		_, err := NewRecipeBuilder().WithPointCurve(points).Build()
		if err == nil {
			t.Error("Expected error for invalid point curve input")
		}
	})

	t.Run("WithPointCurve invalid output", func(t *testing.T) {
		points := []ToneCurvePoint{{Input: 128, Output: 256}}
		_, err := NewRecipeBuilder().WithPointCurve(points).Build()
		if err == nil {
			t.Error("Expected error for invalid point curve output")
		}
	})

	t.Run("WithSplitToning invalid shadow hue", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSplitToning(361, 0, 0, 0, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid split shadow hue")
		}
	})

	t.Run("WithSplitToning invalid shadow saturation", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSplitToning(0, 101, 0, 0, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid split shadow saturation")
		}
	})

	t.Run("WithSplitToning invalid highlight hue", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSplitToning(0, 0, -1, 0, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid split highlight hue")
		}
	})

	t.Run("WithSplitToning invalid highlight saturation", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSplitToning(0, 0, 0, -1, 0).Build()
		if err == nil {
			t.Error("Expected error for invalid split highlight saturation")
		}
	})

	t.Run("WithSplitToning invalid balance", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithSplitToning(0, 0, 0, 0, 101).Build()
		if err == nil {
			t.Error("Expected error for invalid split balance")
		}
	})

	t.Run("WithGrain invalid amount", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithGrain(101, 50, 50).Build()
		if err == nil {
			t.Error("Expected error for invalid grain amount")
		}
	})

	t.Run("WithGrain invalid size", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithGrain(50, -1, 50).Build()
		if err == nil {
			t.Error("Expected error for invalid grain size")
		}
	})

	t.Run("WithGrain invalid roughness", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithGrain(50, 50, 101).Build()
		if err == nil {
			t.Error("Expected error for invalid grain roughness")
		}
	})

	t.Run("WithVignette invalid amount", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithVignette(-101, 50, 0, 50).Build()
		if err == nil {
			t.Error("Expected error for invalid vignette amount")
		}
	})

	t.Run("WithVignette invalid midpoint", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithVignette(0, 101, 0, 50).Build()
		if err == nil {
			t.Error("Expected error for invalid vignette midpoint")
		}
	})

	t.Run("WithVignette invalid roundness", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithVignette(0, 50, 101, 50).Build()
		if err == nil {
			t.Error("Expected error for invalid vignette roundness")
		}
	})

	t.Run("WithVignette invalid feather", func(t *testing.T) {
		_, err := NewRecipeBuilder().WithVignette(0, 50, 0, -1).Build()
		if err == nil {
			t.Error("Expected error for invalid vignette feather")
		}
	})
}

// TestValidationHelpers tests validation helper functions.
func TestValidationHelpers(t *testing.T) {
	t.Run("ValidatePositivePercentage", func(t *testing.T) {
		tests := []struct {
			value   int
			wantErr bool
		}{
			{0, false},
			{50, false},
			{100, false},
			{-1, true},
			{101, true},
		}
		for _, tt := range tests {
			err := ValidatePositivePercentage(tt.value, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositivePercentage(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})

	t.Run("ValidateHue360", func(t *testing.T) {
		tests := []struct {
			value   int
			wantErr bool
		}{
			{0, false},
			{180, false},
			{360, false},
			{-1, true},
			{361, true},
		}
		for _, tt := range tests {
			err := ValidateHue360(tt.value, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHue360(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})

	t.Run("ValidateTint", func(t *testing.T) {
		tests := []struct {
			value   int
			wantErr bool
		}{
			{-150, false},
			{0, false},
			{150, false},
			{-151, true},
			{151, true},
		}
		for _, tt := range tests {
			err := ValidateTint(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTint(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})

	t.Run("ValidateToneCurvePoint", func(t *testing.T) {
		tests := []struct {
			input   int
			output  int
			wantErr bool
		}{
			{0, 0, false},
			{128, 140, false},
			{255, 255, false},
			{-1, 100, true},
			{100, 256, true},
		}
		for _, tt := range tests {
			err := ValidateToneCurvePoint(tt.input, tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToneCurvePoint(%d, %d) error = %v, wantErr %v", tt.input, tt.output, err, tt.wantErr)
			}
		}
	})

	t.Run("ValidateSharpnessRadius", func(t *testing.T) {
		tests := []struct {
			value   float64
			wantErr bool
		}{
			{0.5, false},
			{1.5, false},
			{3.0, false},
			{0.4, true},
			{3.1, true},
		}
		for _, tt := range tests {
			err := ValidateSharpnessRadius(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSharpnessRadius(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		}
	})
}

// TestValidateAllErrorPaths tests all error paths in the Validate() method.
func TestValidateAllErrorPaths(t *testing.T) {
	tests := []struct{
		name string
		recipe UniversalRecipe
		expectError bool
	}{
		{"Valid empty recipe", UniversalRecipe{}, false},
		{"Invalid contrast", UniversalRecipe{Contrast: 150}, true},
		{"Invalid highlights", UniversalRecipe{Highlights: -150}, true},
		{"Invalid shadows", UniversalRecipe{Shadows: 150}, true},
		{"Invalid whites", UniversalRecipe{Whites: -150}, true},
		{"Invalid blacks", UniversalRecipe{Blacks: 150}, true},
		{"Invalid texture", UniversalRecipe{Texture: -150}, true},
		{"Invalid clarity", UniversalRecipe{Clarity: 150}, true},
		{"Invalid dehaze", UniversalRecipe{Dehaze: -150}, true},
		{"Invalid vibrance", UniversalRecipe{Vibrance: 150}, true},
		{"Invalid saturation", UniversalRecipe{Saturation: -150}, true},
		{"Invalid sharpness radius", UniversalRecipe{SharpnessRadius: 5.0}, true},
		{"Invalid sharpness detail", UniversalRecipe{SharpnessDetail: 150}, true},
		{"Invalid sharpness masking", UniversalRecipe{SharpnessMasking: -10}, true},
		{"Invalid tint", UniversalRecipe{Tint: 200}, true},
		{"Invalid red hue", UniversalRecipe{Red: ColorAdjustment{Hue: 200}}, true},
		{"Invalid orange saturation", UniversalRecipe{Orange: ColorAdjustment{Saturation: 200}}, true},
		{"Invalid yellow luminance", UniversalRecipe{Yellow: ColorAdjustment{Luminance: -200}}, true},
		{"Invalid green hue", UniversalRecipe{Green: ColorAdjustment{Hue: -200}}, true},
		{"Invalid aqua saturation", UniversalRecipe{Aqua: ColorAdjustment{Saturation: -200}}, true},
		{"Invalid blue luminance", UniversalRecipe{Blue: ColorAdjustment{Luminance: 200}}, true},
		{"Invalid purple hue", UniversalRecipe{Purple: ColorAdjustment{Hue: 200}}, true},
		{"Invalid magenta saturation", UniversalRecipe{Magenta: ColorAdjustment{Saturation: 200}}, true},
		{"Invalid tone curve shadows", UniversalRecipe{ToneCurveShadows: 150}, true},
		{"Invalid tone curve darks", UniversalRecipe{ToneCurveDarks: -150}, true},
		{"Invalid tone curve lights", UniversalRecipe{ToneCurveLights: 150}, true},
		{"Invalid tone curve highlights", UniversalRecipe{ToneCurveHighlights: -150}, true},
		{"Invalid tone curve shadow split", UniversalRecipe{ToneCurveShadowSplit: 150}, true},
		{"Invalid tone curve midtone split", UniversalRecipe{ToneCurveMidtoneSplit: -10}, true},
		{"Invalid tone curve highlight split", UniversalRecipe{ToneCurveHighlightSplit: 150}, true},
		{"Invalid point curve input", UniversalRecipe{PointCurve: []ToneCurvePoint{{Input: -1, Output: 128}}}, true},
		{"Invalid point curve output", UniversalRecipe{PointCurve: []ToneCurvePoint{{Input: 128, Output: 300}}}, true},
		{"Invalid point curve red", UniversalRecipe{PointCurveRed: []ToneCurvePoint{{Input: 300, Output: 128}}}, true},
		{"Invalid point curve green", UniversalRecipe{PointCurveGreen: []ToneCurvePoint{{Input: 128, Output: -1}}}, true},
		{"Invalid point curve blue", UniversalRecipe{PointCurveBlue: []ToneCurvePoint{{Input: -1, Output: 128}}}, true},
		{"Invalid split shadow hue", UniversalRecipe{SplitShadowHue: 400}, true},
		{"Invalid split shadow saturation", UniversalRecipe{SplitShadowSaturation: 150}, true},
		{"Invalid split highlight hue", UniversalRecipe{SplitHighlightHue: -10}, true},
		{"Invalid split highlight saturation", UniversalRecipe{SplitHighlightSaturation: -10}, true},
		{"Invalid split balance", UniversalRecipe{SplitBalance: 150}, true},
		{"Invalid grain amount", UniversalRecipe{GrainAmount: 150}, true},
		{"Invalid grain size", UniversalRecipe{GrainSize: -10}, true},
		{"Invalid grain roughness", UniversalRecipe{GrainRoughness: 150}, true},
		{"Invalid vignette amount", UniversalRecipe{VignetteAmount: 150}, true},
		{"Invalid vignette midpoint", UniversalRecipe{VignetteMidpoint: -10}, true},
		{"Invalid vignette roundness", UniversalRecipe{VignetteRoundness: 150}, true},
		{"Invalid vignette feather", UniversalRecipe{VignetteFeather: -10}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.recipe.Validate()
			if tt.expectError && err == nil {
				t.Errorf("Expected validation error for %s", tt.name)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected validation error for %s: %v", tt.name, err)
			}
		})
	}
}
