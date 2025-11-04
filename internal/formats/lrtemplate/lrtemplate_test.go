package lrtemplate

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestParse tests parsing of lrtemplate files using table-driven tests
func TestParse(t *testing.T) {
	// Discover all lrtemplate sample files
	files, err := filepath.Glob("../../../examples/lrtemplate/*.lrtemplate")
	if err != nil {
		t.Fatalf("Failed to discover lrtemplate files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No lrtemplate sample files found in examples/lrtemplate/")
	}

	t.Logf("Found %d lrtemplate sample files", len(files))

	// Test each file
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			recipe, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			// Verify SourceFormat is set
			if recipe.SourceFormat != "lrtemplate" {
				t.Errorf("Expected SourceFormat='lrtemplate', got '%s'", recipe.SourceFormat)
			}
		})
	}
}

// TestParseInvalidPrefix tests error handling when file doesn't start with `s = {`
func TestParseInvalidPrefix(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "empty file",
			data: "",
		},
		{
			name: "wrong prefix",
			data: "table = { Exposure2012 = 1.5 }",
		},
		{
			name: "missing opening brace",
			data: "s = Exposure2012 = 1.5 }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data))
			if err == nil {
				t.Error("Expected error for invalid prefix, got nil")
			}

			// Verify error is ConversionError
			convErr, ok := err.(*ConversionError)
			if !ok {
				t.Errorf("Expected ConversionError, got %T", err)
			} else {
				if convErr.Operation != "parse" {
					t.Errorf("Expected Operation='parse', got '%s'", convErr.Operation)
				}
				if convErr.Format != "lrtemplate" {
					t.Errorf("Expected Format='lrtemplate', got '%s'", convErr.Format)
				}
			}
		})
	}
}

// TestParseBasicParameters tests extraction of core adjustment parameters
func TestParseBasicParameters(t *testing.T) {
	data := `s = {
		value = {
			settings = {
				Exposure2012 = 1.5,
				Contrast2012 = 25,
				Highlights2012 = -50,
				Shadows2012 = 40,
				Whites2012 = -20,
				Blacks2012 = 15,
				Saturation = 10,
				Vibrance = 20,
				Clarity2012 = 30,
				Sharpness = 50,
				Temperature = 5500,
				Tint = -10,
			}
		}
	}`

	recipe, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify basic adjustments
	if recipe.Exposure != 1.5 {
		t.Errorf("Expected Exposure=1.5, got %v", recipe.Exposure)
	}
	if recipe.Contrast != 25 {
		t.Errorf("Expected Contrast=25, got %d", recipe.Contrast)
	}
	if recipe.Highlights != -50 {
		t.Errorf("Expected Highlights=-50, got %d", recipe.Highlights)
	}
	if recipe.Shadows != 40 {
		t.Errorf("Expected Shadows=40, got %d", recipe.Shadows)
	}
	if recipe.Whites != -20 {
		t.Errorf("Expected Whites=-20, got %d", recipe.Whites)
	}
	if recipe.Blacks != 15 {
		t.Errorf("Expected Blacks=15, got %d", recipe.Blacks)
	}

	// Verify color parameters
	if recipe.Saturation != 10 {
		t.Errorf("Expected Saturation=10, got %d", recipe.Saturation)
	}
	if recipe.Vibrance != 20 {
		t.Errorf("Expected Vibrance=20, got %d", recipe.Vibrance)
	}
	if recipe.Clarity != 30 {
		t.Errorf("Expected Clarity=30, got %d", recipe.Clarity)
	}
	if recipe.Sharpness != 50 {
		t.Errorf("Expected Sharpness=50, got %d", recipe.Sharpness)
	}
	if recipe.Temperature == nil || *recipe.Temperature != 5500 {
		if recipe.Temperature == nil {
			t.Error("Expected Temperature=5500, got nil")
		} else {
			t.Errorf("Expected Temperature=5500, got %d", *recipe.Temperature)
		}
	}
	if recipe.Tint != -10 {
		t.Errorf("Expected Tint=-10, got %d", recipe.Tint)
	}
}

// TestParseHSLColors tests extraction of HSL color adjustments
func TestParseHSLColors(t *testing.T) {
	data := `s = {
		value = {
			settings = {
				HueAdjustmentRed = 10,
				SaturationAdjustmentRed = -20,
				LuminanceAdjustmentRed = 30,

				HueAdjustmentOrange = 5,
				SaturationAdjustmentOrange = -10,
				LuminanceAdjustmentOrange = 15,

				HueAdjustmentYellow = -5,
				SaturationAdjustmentYellow = 10,
				LuminanceAdjustmentYellow = -15,

				HueAdjustmentGreen = 20,
				SaturationAdjustmentGreen = -25,
				LuminanceAdjustmentGreen = 35,
			}
		}
	}`

	recipe, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify Red HSL
	if recipe.Red.Hue != 10 {
		t.Errorf("Expected Red.Hue=10, got %d", recipe.Red.Hue)
	}
	if recipe.Red.Saturation != -20 {
		t.Errorf("Expected Red.Saturation=-20, got %d", recipe.Red.Saturation)
	}
	if recipe.Red.Luminance != 30 {
		t.Errorf("Expected Red.Luminance=30, got %d", recipe.Red.Luminance)
	}

	// Verify Orange HSL
	if recipe.Orange.Hue != 5 {
		t.Errorf("Expected Orange.Hue=5, got %d", recipe.Orange.Hue)
	}
	if recipe.Orange.Saturation != -10 {
		t.Errorf("Expected Orange.Saturation=-10, got %d", recipe.Orange.Saturation)
	}
	if recipe.Orange.Luminance != 15 {
		t.Errorf("Expected Orange.Luminance=15, got %d", recipe.Orange.Luminance)
	}

	// Verify Yellow HSL
	if recipe.Yellow.Hue != -5 {
		t.Errorf("Expected Yellow.Hue=-5, got %d", recipe.Yellow.Hue)
	}
	if recipe.Yellow.Saturation != 10 {
		t.Errorf("Expected Yellow.Saturation=10, got %d", recipe.Yellow.Saturation)
	}
	if recipe.Yellow.Luminance != -15 {
		t.Errorf("Expected Yellow.Luminance=-15, got %d", recipe.Yellow.Luminance)
	}

	// Verify Green HSL
	if recipe.Green.Hue != 20 {
		t.Errorf("Expected Green.Hue=20, got %d", recipe.Green.Hue)
	}
	if recipe.Green.Saturation != -25 {
		t.Errorf("Expected Green.Saturation=-25, got %d", recipe.Green.Saturation)
	}
	if recipe.Green.Luminance != 35 {
		t.Errorf("Expected Green.Luminance=35, got %d", recipe.Green.Luminance)
	}
}

// TestParseToneCurve tests parsing of ToneCurvePV2012 array
func TestParseToneCurve(t *testing.T) {
	data := `s = {
		value = {
			settings = {
				ToneCurvePV2012 = {
					0, 0,
					64, 58,
					128, 135,
					192, 196,
					255, 255,
				},
			}
		}
	}`

	recipe, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify tone curve points
	expectedPoints := []models.ToneCurvePoint{
		{Input: 0, Output: 0},
		{Input: 64, Output: 58},
		{Input: 128, Output: 135},
		{Input: 192, Output: 196},
		{Input: 255, Output: 255},
	}

	if len(recipe.PointCurve) != len(expectedPoints) {
		t.Fatalf("Expected %d tone curve points, got %d", len(expectedPoints), len(recipe.PointCurve))
	}

	for i, expected := range expectedPoints {
		actual := recipe.PointCurve[i]
		if actual.Input != expected.Input || actual.Output != expected.Output {
			t.Errorf("Point %d: expected (%d, %d), got (%d, %d)",
				i, expected.Input, expected.Output, actual.Input, actual.Output)
		}
	}
}

// TestParseSplitToning tests extraction of split toning parameters
func TestParseSplitToning(t *testing.T) {
	data := `s = {
		value = {
			settings = {
				SplitToningShadowHue = 240,
				SplitToningShadowSaturation = 25,
				SplitToningHighlightHue = 60,
				SplitToningHighlightSaturation = 30,
				SplitToningBalance = -10,
			}
		}
	}`

	recipe, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if recipe.SplitShadowHue != 240 {
		t.Errorf("Expected SplitShadowHue=240, got %d", recipe.SplitShadowHue)
	}
	if recipe.SplitShadowSaturation != 25 {
		t.Errorf("Expected SplitShadowSaturation=25, got %d", recipe.SplitShadowSaturation)
	}
	if recipe.SplitHighlightHue != 60 {
		t.Errorf("Expected SplitHighlightHue=60, got %d", recipe.SplitHighlightHue)
	}
	if recipe.SplitHighlightSaturation != 30 {
		t.Errorf("Expected SplitHighlightSaturation=30, got %d", recipe.SplitHighlightSaturation)
	}
	if recipe.SplitBalance != -10 {
		t.Errorf("Expected SplitBalance=-10, got %d", recipe.SplitBalance)
	}
}

// TestParseDataTypes tests correct parsing of different data types
func TestParseDataTypes(t *testing.T) {
	data := `s = {
		value = {
			settings = {
				Exposure2012 = -2.5,
				Contrast2012 = -50,
				Highlights2012 = 0,
				Sharpness = 100,
			}
		}
	}`

	recipe, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify negative float
	if recipe.Exposure != -2.5 {
		t.Errorf("Expected Exposure=-2.5, got %v", recipe.Exposure)
	}

	// Verify negative integer
	if recipe.Contrast != -50 {
		t.Errorf("Expected Contrast=-50, got %d", recipe.Contrast)
	}

	// Verify zero value
	if recipe.Highlights != 0 {
		t.Errorf("Expected Highlights=0, got %d", recipe.Highlights)
	}

	// Verify positive integer
	if recipe.Sharpness != 100 {
		t.Errorf("Expected Sharpness=100, got %d", recipe.Sharpness)
	}
}

// TestParseMissingFields tests graceful handling of missing fields
func TestParseMissingFields(t *testing.T) {
	data := `s = {
		value = {
			settings = {
				Exposure2012 = 1.0,
			}
		}
	}`

	recipe, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify Exposure is set
	if recipe.Exposure != 1.0 {
		t.Errorf("Expected Exposure=1.0, got %v", recipe.Exposure)
	}

	// Verify missing fields default to zero
	if recipe.Contrast != 0 {
		t.Errorf("Expected Contrast=0 (missing), got %d", recipe.Contrast)
	}
	if recipe.Highlights != 0 {
		t.Errorf("Expected Highlights=0 (missing), got %d", recipe.Highlights)
	}
	if recipe.Saturation != 0 {
		t.Errorf("Expected Saturation=0 (missing), got %d", recipe.Saturation)
	}
}

// TestParseInvalidSyntax tests error handling for malformed Lua syntax
func TestParseInvalidSyntax(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "no closing brace at all",
			data: "s = { value = settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data))
			if err == nil {
				t.Error("Expected error for invalid syntax, got nil")
			}
		})
	}
}

// TestParseValidationErrors tests parameter range validation
func TestParseValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantField string
	}{
		{
			name: "exposure out of range",
			data: `s = {
				value = {
					settings = {
						Exposure2012 = 10.0,
					}
				}
			}`,
			wantField: "Exposure2012",
		},
		{
			name: "contrast out of range",
			data: `s = {
				value = {
					settings = {
						Contrast2012 = 150,
					}
				}
			}`,
			wantField: "Contrast2012",
		},
		{
			name: "sharpness out of range",
			data: `s = {
				value = {
					settings = {
						Sharpness = 200,
					}
				}
			}`,
			wantField: "Sharpness",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data))
			if err == nil {
				t.Error("Expected validation error, got nil")
			}

			convErr, ok := err.(*ConversionError)
			if !ok {
				t.Errorf("Expected ConversionError, got %T", err)
			} else {
				if convErr.Operation != "validate" {
					t.Errorf("Expected Operation='validate', got '%s'", convErr.Operation)
				}
				if convErr.Field != tt.wantField {
					t.Errorf("Expected Field='%s', got '%s'", tt.wantField, convErr.Field)
				}
			}
		})
	}
}

// TestParseEdgeCases tests edge cases like empty strings, special characters
func TestParseEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "whitespace variations",
			data: "s = {\n\tvalue = {\n\t\tsettings = {\n\t\t\tExposure2012=1.5,\n\t\t}\n\t}\n}",
		},
		{
			name: "single line format",
			data: "s = { value = { settings = { Exposure2012 = 1.5, Contrast2012 = 25 } } }",
		},
		{
			name: "minimal valid file",
			data: "s = { value = { settings = {} } }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data))
			if err != nil {
				t.Errorf("Parse failed for valid edge case: %v", err)
			}
		})
	}
}

// TestParseAdditionalValidation tests additional validation scenarios
func TestParseAdditionalValidation(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantField string
	}{
		{
			name: "highlights out of range",
			data: `s = {
				value = {
					settings = {
						Highlights2012 = -150,
					}
				}
			}`,
			wantField: "Highlights2012",
		},
		{
			name: "shadows out of range",
			data: `s = {
				value = {
					settings = {
						Shadows2012 = 150,
					}
				}
			}`,
			wantField: "Shadows2012",
		},
		{
			name: "whites out of range",
			data: `s = {
				value = {
					settings = {
						Whites2012 = -150,
					}
				}
			}`,
			wantField: "Whites2012",
		},
		{
			name: "blacks out of range",
			data: `s = {
				value = {
					settings = {
						Blacks2012 = 150,
					}
				}
			}`,
			wantField: "Blacks2012",
		},
		{
			name: "saturation out of range",
			data: `s = {
				value = {
					settings = {
						Saturation = -150,
					}
				}
			}`,
			wantField: "Saturation",
		},
		{
			name: "vibrance out of range",
			data: `s = {
				value = {
					settings = {
						Vibrance = 150,
					}
				}
			}`,
			wantField: "Vibrance",
		},
		{
			name: "clarity out of range",
			data: `s = {
				value = {
					settings = {
						Clarity2012 = -150,
					}
				}
			}`,
			wantField: "Clarity2012",
		},
		{
			name: "temperature out of range low",
			data: `s = {
				value = {
					settings = {
						Temperature = 1000,
					}
				}
			}`,
			wantField: "Temperature",
		},
		{
			name: "temperature out of range high",
			data: `s = {
				value = {
					settings = {
						Temperature = 60000,
					}
				}
			}`,
			wantField: "Temperature",
		},
		{
			name: "tint out of range",
			data: `s = {
				value = {
					settings = {
						Tint = -200,
					}
				}
			}`,
			wantField: "Tint",
		},
		{
			name: "hsl hue out of range",
			data: `s = {
				value = {
					settings = {
						HueAdjustmentRed = 150,
					}
				}
			}`,
			wantField: "HueAdjustmentRed",
		},
		{
			name: "hsl saturation out of range",
			data: `s = {
				value = {
					settings = {
						SaturationAdjustmentOrange = -150,
					}
				}
			}`,
			wantField: "SaturationAdjustmentOrange",
		},
		{
			name: "hsl luminance out of range",
			data: `s = {
				value = {
					settings = {
						LuminanceAdjustmentYellow = 150,
					}
				}
			}`,
			wantField: "LuminanceAdjustmentYellow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data))
			if err == nil {
				t.Fatal("Expected validation error, got nil")
			}

			convErr, ok := err.(*ConversionError)
			if !ok {
				t.Errorf("Expected ConversionError, got %T", err)
			} else {
				if convErr.Operation != "validate" {
					t.Errorf("Expected Operation='validate', got '%s'", convErr.Operation)
				}
				if convErr.Field != tt.wantField {
					t.Errorf("Expected Field='%s', got '%s'", tt.wantField, convErr.Field)
				}
			}
		})
	}
}

// TestConversionErrorMethods tests the Error() and Unwrap() methods
func TestConversionErrorMethods(t *testing.T) {
	baseErr := errors.New("base error")

	// Test with Field set
	convErr := &ConversionError{
		Operation: "parse",
		Format:    "lrtemplate",
		Field:     "Exposure2012",
		Cause:     baseErr,
	}

	errorMsg := convErr.Error()
	if !strings.Contains(errorMsg, "parse") {
		t.Errorf("Error message should contain operation: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "lrtemplate") {
		t.Errorf("Error message should contain format: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "Exposure2012") {
		t.Errorf("Error message should contain field: %s", errorMsg)
	}

	// Test Unwrap() method
	unwrapped := errors.Unwrap(convErr)
	if unwrapped != baseErr {
		t.Errorf("Unwrap should return base error, got %v", unwrapped)
	}

	// Test error wrapping with errors.Is
	if !errors.Is(convErr, baseErr) {
		t.Error("errors.Is should find the base error")
	}

	// Test Error() method without Field
	convErrNoField := &ConversionError{
		Operation: "parse",
		Format:    "lrtemplate",
		Cause:     baseErr,
	}
	errorMsg2 := convErrNoField.Error()
	if !strings.Contains(errorMsg2, "parse") {
		t.Errorf("Error message should contain operation: %s", errorMsg2)
	}
	if !strings.Contains(errorMsg2, "lrtemplate") {
		t.Errorf("Error message should contain format: %s", errorMsg2)
	}
}

// TestParseMoreValidationCases tests additional validation scenarios
func TestParseMoreValidationCases(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantField string
	}{
		{
			name: "split toning highlight hue out of range",
			data: `s = {
				value = {
					settings = {
						SplitToningHighlightHue = 400,
					}
				}
			}`,
			wantField: "SplitToningHighlightHue",
		},
		{
			name: "split toning shadow saturation out of range",
			data: `s = {
				value = {
					settings = {
						SplitToningShadowSaturation = 150,
					}
				}
			}`,
			wantField: "SplitToningShadowSaturation",
		},
		{
			name: "hsl aqua hue out of range",
			data: `s = {
				value = {
					settings = {
						HueAdjustmentAqua = -150,
					}
				}
			}`,
			wantField: "HueAdjustmentAqua",
		},
		{
			name: "hsl blue saturation out of range",
			data: `s = {
				value = {
					settings = {
						SaturationAdjustmentBlue = 150,
					}
				}
			}`,
			wantField: "SaturationAdjustmentBlue",
		},
		{
			name: "hsl purple luminance out of range",
			data: `s = {
				value = {
					settings = {
						LuminanceAdjustmentPurple = -150,
					}
				}
			}`,
			wantField: "LuminanceAdjustmentPurple",
		},
		{
			name: "hsl magenta hue out of range",
			data: `s = {
				value = {
					settings = {
						HueAdjustmentMagenta = 150,
					}
				}
			}`,
			wantField: "HueAdjustmentMagenta",
		},
		{
			name: "hsl green saturation out of range",
			data: `s = {
				value = {
					settings = {
						SaturationAdjustmentGreen = -150,
					}
				}
			}`,
			wantField: "SaturationAdjustmentGreen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data))
			if err == nil {
				t.Fatal("Expected validation error, got nil")
			}

			convErr, ok := err.(*ConversionError)
			if !ok {
				t.Errorf("Expected ConversionError, got %T", err)
			} else {
				if convErr.Operation != "validate" {
					t.Errorf("Expected Operation='validate', got '%s'", convErr.Operation)
				}
				if convErr.Field != tt.wantField {
					t.Errorf("Expected Field='%s', got '%s'", tt.wantField, convErr.Field)
				}
			}
		})
	}
}

// TestParseAllParameters tests comprehensive parameter extraction
func TestParseAllParameters(t *testing.T) {
	data := `s = {
		value = {
			settings = {
				Exposure2012 = 1.5,
				Contrast2012 = 25,
				Highlights2012 = -30,
				Shadows2012 = 40,
				Whites2012 = 15,
				Blacks2012 = -10,
				Saturation = 20,
				Vibrance = 15,
				Clarity2012 = 30,
				Sharpness = 50,
				Temperature = 5500,
				Tint = 10,
				SplitToningHighlightHue = 45,
				SplitToningHighlightSaturation = 20,
				SplitToningShadowHue = 220,
				SplitToningShadowSaturation = 15,
				HueAdjustmentRed = 10,
				HueAdjustmentOrange = -5,
				SaturationAdjustmentYellow = 15,
				LuminanceAdjustmentGreen = 20,
			}
		}
	}`

	recipe, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify all parameters are extracted
	if recipe.Exposure != 1.5 {
		t.Errorf("Expected Exposure=1.5, got %v", recipe.Exposure)
	}
	if recipe.Contrast != 25 {
		t.Errorf("Expected Contrast=25, got %d", recipe.Contrast)
	}
	if recipe.Highlights != -30 {
		t.Errorf("Expected Highlights=-30, got %d", recipe.Highlights)
	}
	if recipe.Shadows != 40 {
		t.Errorf("Expected Shadows=40, got %d", recipe.Shadows)
	}
	if recipe.Clarity != 30 {
		t.Errorf("Expected Clarity=30, got %d", recipe.Clarity)
	}
	if recipe.Sharpness != 50 {
		t.Errorf("Expected Sharpness=50, got %d", recipe.Sharpness)
	}
	if recipe.Temperature == nil || *recipe.Temperature != 5500 {
		t.Errorf("Expected Temperature=5500, got %v", recipe.Temperature)
	}
	if recipe.Tint != 10 {
		t.Errorf("Expected Tint=10, got %d", recipe.Tint)
	}
}

// BenchmarkParse benchmarks parsing performance to validate <20ms target
func BenchmarkParse(t *testing.B) {
	// Load a sample lrtemplate file
	files, err := filepath.Glob("../../../examples/lrtemplate/*.lrtemplate")
	if err != nil || len(files) == 0 {
		t.Skip("No sample files found for benchmark")
	}

	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("Failed to read sample file: %v", err)
	}

	// Run benchmark
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_, err := Parse(data)
		if err != nil {
			t.Fatalf("Parse failed during benchmark: %v", err)
		}
	}
}
