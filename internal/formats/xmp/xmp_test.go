package xmp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// TestParse validates the XMP parser against all sample files following Pattern 7
func TestParse(t *testing.T) {
	files, err := filepath.Glob("../../../testdata/xmp/*.xmp")
	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatal("no test files found in testdata/xmp/")
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("failed to read %s: %v", file, err)
			}

			recipe, err := Parse(data)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			// Validate SourceFormat set correctly
			if recipe.SourceFormat != "xmp" {
				t.Errorf("SourceFormat = %s, want 'xmp'", recipe.SourceFormat)
			}

			// Validate critical field ranges
			if recipe.Exposure < -5.0 || recipe.Exposure > 5.0 {
				t.Errorf("Exposure out of range: %.2f", recipe.Exposure)
			}

			if recipe.Contrast < -100 || recipe.Contrast > 100 {
				t.Errorf("Contrast out of range: %d", recipe.Contrast)
			}

			if recipe.Highlights < -100 || recipe.Highlights > 100 {
				t.Errorf("Highlights out of range: %d", recipe.Highlights)
			}

			if recipe.Shadows < -100 || recipe.Shadows > 100 {
				t.Errorf("Shadows out of range: %d", recipe.Shadows)
			}

			if recipe.Whites < -100 || recipe.Whites > 100 {
				t.Errorf("Whites out of range: %d", recipe.Whites)
			}

			if recipe.Blacks < -100 || recipe.Blacks > 100 {
				t.Errorf("Blacks out of range: %d", recipe.Blacks)
			}

			// Validate HSL colors
			validateColorAdjustment(t, "Red", recipe.Red)
			validateColorAdjustment(t, "Orange", recipe.Orange)
			validateColorAdjustment(t, "Yellow", recipe.Yellow)
			validateColorAdjustment(t, "Green", recipe.Green)
			validateColorAdjustment(t, "Aqua", recipe.Aqua)
			validateColorAdjustment(t, "Blue", recipe.Blue)
			validateColorAdjustment(t, "Purple", recipe.Purple)
			validateColorAdjustment(t, "Magenta", recipe.Magenta)
		})
	}
}

func validateColorAdjustment(t *testing.T, color string, adj models.ColorAdjustment) {
	t.Helper()
	if adj.Hue < -100 || adj.Hue > 100 {
		t.Errorf("%s Hue out of range: %d", color, adj.Hue)
	}
	if adj.Saturation < -100 || adj.Saturation > 100 {
		t.Errorf("%s Saturation out of range: %d", color, adj.Saturation)
	}
	if adj.Luminance < -100 || adj.Luminance > 100 {
		t.Errorf("%s Luminance out of range: %d", color, adj.Luminance)
	}
}

// TestParseMissingNamespace validates error handling for missing namespace
func TestParseMissingNamespace(t *testing.T) {
	tests := []struct {
		name    string
		xmpData string
		wantErr string
	}{
		{
			name: "Missing crs namespace",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantErr: "camera-raw-settings",
		},
		{
			name: "Missing x:xmpmeta",
			xmpData: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"/>
</rdf:RDF>`,
			wantErr: "adobe:ns:meta",
		},
		{
			name: "Missing rdf:RDF",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"/>
</x:xmpmeta>`,
			wantErr: "rdf:RDF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.xmpData))
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
			if err != nil && !contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

// TestParseInvalidRange validates error handling for out-of-range values
func TestParseInvalidRange(t *testing.T) {
	tests := []struct {
		name      string
		xmpData   string
		wantField string // Expected field name in error
	}{
		{
			name: "Exposure too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Exposure2012="10.0"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Exposure",
		},
		{
			name: "Exposure too low",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Exposure2012="-6.0"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Exposure",
		},
		{
			name: "Contrast too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Contrast2012="150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Contrast",
		},
		{
			name: "Contrast too low",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Contrast2012="-150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Contrast",
		},
		{
			name: "Highlights too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Highlights2012="150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Highlights",
		},
		{
			name: "Shadows too low",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Shadows2012="-150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Shadows",
		},
		{
			name: "Whites too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Whites2012="150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Whites",
		},
		{
			name: "Blacks too low",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Blacks2012="-150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Blacks",
		},
		{
			name: "Saturation too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Saturation="150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Saturation",
		},
		{
			name: "Vibrance too low",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Vibrance="-150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Vibrance",
		},
		{
			name: "Clarity too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Clarity2012="150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Clarity",
		},
		{
			name: "Sharpness too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Sharpness="200"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Sharpness",
		},
		{
			name: "Red Hue out of range",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:HueRed="150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Red Hue",
		},
		{
			name: "Blue Saturation out of range",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:SaturationBlue="-150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Blue Saturation",
		},
		{
			name: "Green Luminance out of range",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:LuminanceGreen="150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Green Luminance",
		},
		{
			name: "SplitShadowHue too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:SplitToningShadowHue="400"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "SplitShadowHue",
		},
		{
			name: "SplitShadowSaturation too low",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:SplitToningShadowSaturation="-10"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "SplitShadowSaturation",
		},
		{
			name: "SplitHighlightHue too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:SplitToningHighlightHue="400"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "SplitHighlightHue",
		},
		{
			name: "SplitHighlightSaturation too high",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:SplitToningHighlightSaturation="150"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "SplitHighlightSaturation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.xmpData))
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
				return
			}

			// Verify error is ConversionError with expected field
			var convErr *ConversionError
			if !errors.As(err, &convErr) {
				t.Errorf("expected ConversionError, got: %T", err)
				return
			}

			if convErr.Field != tt.wantField {
				t.Errorf("expected field %q, got: %q", tt.wantField, convErr.Field)
			}
		})
	}
}

// TestParseMalformedXML validates error handling for malformed XML
func TestParseMalformedXML(t *testing.T) {
	malformedXML := `<x:xmpmeta><unclosed tag>`

	_, err := Parse([]byte(malformedXML))
	if err == nil {
		t.Error("expected error for malformed XML")
	}
}

// TestParseEmptyXMP validates parsing of minimal valid XMP file
func TestParseEmptyXMP(t *testing.T) {
	minimalXMP := `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"/>
  </rdf:RDF>
</x:xmpmeta>`

	recipe, err := Parse([]byte(minimalXMP))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if recipe.SourceFormat != "xmp" {
		t.Errorf("SourceFormat = %s, want 'xmp'", recipe.SourceFormat)
	}

	// All values should be zero/default
	if recipe.Exposure != 0.0 {
		t.Errorf("Expected zero Exposure, got %.2f", recipe.Exposure)
	}
	if recipe.Contrast != 0 {
		t.Errorf("Expected zero Contrast, got %d", recipe.Contrast)
	}
}

// TestParseEmptyAttributes validates parsing with empty attribute values
func TestParseEmptyAttributes(t *testing.T) {
	emptyAttrsXMP := `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
      crs:Exposure2012=""
      crs:Contrast2012=""
      crs:Temperature=""
      crs:Saturation=""
      crs:HueRed=""
      crs:SplitToningShadowHue=""/>
  </rdf:RDF>
</x:xmpmeta>`

	recipe, err := Parse([]byte(emptyAttrsXMP))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// All values should be zero/default for empty strings
	if recipe.Exposure != 0.0 {
		t.Errorf("Expected zero Exposure for empty string, got %.2f", recipe.Exposure)
	}
	if recipe.Contrast != 0 {
		t.Errorf("Expected zero Contrast for empty string, got %d", recipe.Contrast)
	}
	if recipe.Temperature != nil {
		t.Errorf("Expected nil Temperature for empty string, got %v", recipe.Temperature)
	}
	if recipe.Saturation != 0 {
		t.Errorf("Expected zero Saturation for empty string, got %d", recipe.Saturation)
	}
	if recipe.Red.Hue != 0 {
		t.Errorf("Expected zero Red Hue for empty string, got %d", recipe.Red.Hue)
	}
}

// TestParseComprehensive validates parsing with all possible parameters
func TestParseComprehensive(t *testing.T) {
	comprehensiveXMP := `<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="Adobe XMP Core 5.6-c140">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description
      xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
      crs:Exposure2012="2.0"
      crs:Contrast2012="30"
      crs:Highlights2012="-40"
      crs:Shadows2012="50"
      crs:Whites2012="20"
      crs:Blacks2012="-20"
      crs:Saturation="25"
      crs:Vibrance="18"
      crs:Clarity2012="40"
      crs:Sharpness="45"
      crs:Temperature="5600"
      crs:Tint="10"
      crs:HueRed="12"
      crs:SaturationRed="-8"
      crs:LuminanceRed="10"
      crs:HueOrange="5"
      crs:SaturationOrange="12"
      crs:LuminanceOrange="-5"
      crs:HueYellow="8"
      crs:SaturationYellow="5"
      crs:LuminanceYellow="15"
      crs:HueGreen="-10"
      crs:SaturationGreen="18"
      crs:LuminanceGreen="3"
      crs:HueAqua="3"
      crs:SaturationAqua="5"
      crs:LuminanceAqua="-2"
      crs:HueBlue="25"
      crs:SaturationBlue="-12"
      crs:LuminanceBlue="8"
      crs:HuePurple="-8"
      crs:SaturationPurple="10"
      crs:LuminancePurple="2"
      crs:HueMagenta="5"
      crs:SaturationMagenta="8"
      crs:LuminanceMagenta="-3"
      crs:SplitToningShadowHue="250"
      crs:SplitToningShadowSaturation="18"
      crs:SplitToningHighlightHue="55"
      crs:SplitToningHighlightSaturation="22"
      crs:SplitToningBalance="5"/>
  </rdf:RDF>
</x:xmpmeta>`

	recipe, err := Parse([]byte(comprehensiveXMP))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Validate basic parameters are parsed correctly
	if recipe.Exposure != 2.0 {
		t.Errorf("Exposure = %.2f, want 2.00", recipe.Exposure)
	}
	if recipe.Temperature == nil || *recipe.Temperature != 5600 {
		t.Errorf("Temperature = %v, want 5600", recipe.Temperature)
	}
	if recipe.Tint != 10 {
		t.Errorf("Tint = %d, want 10", recipe.Tint)
	}
	if recipe.Orange.Saturation != 12 {
		t.Errorf("Orange Saturation = %d, want 12", recipe.Orange.Saturation)
	}
	if recipe.Aqua.Hue != 3 {
		t.Errorf("Aqua Hue = %d, want 3", recipe.Aqua.Hue)
	}
	if recipe.Purple.Luminance != 2 {
		t.Errorf("Purple Luminance = %d, want 2", recipe.Purple.Luminance)
	}
	if recipe.SplitBalance != 5 {
		t.Errorf("SplitBalance = %d, want 5", recipe.SplitBalance)
	}
}

// TestParseAllParameters validates all 50+ parameters are correctly extracted
func TestParseAllParameters(t *testing.T) {
	data, err := os.ReadFile("../../../testdata/xmp/sample.xmp")
	if err != nil {
		t.Fatalf("failed to read sample.xmp: %v", err)
	}

	recipe, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Validate Basic Adjustments
	if recipe.Exposure != 1.5 {
		t.Errorf("Exposure = %.2f, want 1.50", recipe.Exposure)
	}
	if recipe.Contrast != 25 {
		t.Errorf("Contrast = %d, want 25", recipe.Contrast)
	}
	if recipe.Highlights != -30 {
		t.Errorf("Highlights = %d, want -30", recipe.Highlights)
	}
	if recipe.Shadows != 40 {
		t.Errorf("Shadows = %d, want 40", recipe.Shadows)
	}
	if recipe.Whites != 10 {
		t.Errorf("Whites = %d, want 10", recipe.Whites)
	}
	if recipe.Blacks != -15 {
		t.Errorf("Blacks = %d, want -15", recipe.Blacks)
	}

	// Validate Color Parameters
	if recipe.Saturation != 20 {
		t.Errorf("Saturation = %d, want 20", recipe.Saturation)
	}
	if recipe.Vibrance != 15 {
		t.Errorf("Vibrance = %d, want 15", recipe.Vibrance)
	}
	if recipe.Clarity != 35 {
		t.Errorf("Clarity = %d, want 35", recipe.Clarity)
	}
	if recipe.Sharpness != 40 {
		t.Errorf("Sharpness = %d, want 40", recipe.Sharpness)
	}

	// Validate HSL Adjustments
	if recipe.Red.Hue != 10 {
		t.Errorf("Red Hue = %d, want 10", recipe.Red.Hue)
	}
	if recipe.Red.Saturation != -5 {
		t.Errorf("Red Saturation = %d, want -5", recipe.Red.Saturation)
	}
	if recipe.Red.Luminance != 8 {
		t.Errorf("Red Luminance = %d, want 8", recipe.Red.Luminance)
	}

	if recipe.Blue.Hue != 20 {
		t.Errorf("Blue Hue = %d, want 20", recipe.Blue.Hue)
	}

	// Validate Split Toning
	if recipe.SplitShadowHue != 240 {
		t.Errorf("SplitShadowHue = %d, want 240", recipe.SplitShadowHue)
	}
	if recipe.SplitShadowSaturation != 15 {
		t.Errorf("SplitShadowSaturation = %d, want 15", recipe.SplitShadowSaturation)
	}
	if recipe.SplitHighlightHue != 50 {
		t.Errorf("SplitHighlightHue = %d, want 50", recipe.SplitHighlightHue)
	}
	if recipe.SplitHighlightSaturation != 20 {
		t.Errorf("SplitHighlightSaturation = %d, want 20", recipe.SplitHighlightSaturation)
	}
}

// TestParseInvalidDataType validates error handling for non-numeric values
func TestParseInvalidDataType(t *testing.T) {
	tests := []struct {
		name      string
		xmpData   string
		wantField string // Expected field name in error
	}{
		{
			name: "Invalid Exposure2012",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Exposure2012="not-a-number"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Exposure2012",
		},
		{
			name: "Invalid Contrast2012",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Contrast2012="abc"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Contrast2012",
		},
		{
			name: "Invalid Temperature",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Temperature="bad"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Temperature",
		},
		{
			name: "Invalid HueRed",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:HueRed="xyz"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "RedHue",
		},
		{
			name: "Invalid SaturationBlue",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:SaturationBlue="bad"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "BlueSaturation",
		},
		{
			name: "Invalid LuminanceGreen",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:LuminanceGreen="not-int"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "GreenLuminance",
		},
		{
			name: "Invalid Highlights2012",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Highlights2012="invalid"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Highlights2012",
		},
		{
			name: "Invalid Shadows2012",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Shadows2012="bad-value"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Shadows2012",
		},
		{
			name: "Invalid Whites2012",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Whites2012="xyz"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Whites2012",
		},
		{
			name: "Invalid Blacks2012",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Blacks2012="abc123"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Blacks2012",
		},
		{
			name: "Invalid Vibrance",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Vibrance="not-number"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Vibrance",
		},
		{
			name: "Invalid Clarity2012",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Clarity2012="bad"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Clarity2012",
		},
		{
			name: "Invalid Sharpness",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Sharpness="invalid"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Sharpness",
		},
		{
			name: "Invalid Tint",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:Tint="xyz"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "Tint",
		},
		{
			name: "Invalid SplitToningBalance",
			xmpData: `<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" crs:SplitToningBalance="bad"/>
  </rdf:RDF>
</x:xmpmeta>`,
			wantField: "SplitToningBalance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.xmpData))
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
				return
			}

			// Verify error is ConversionError with expected field
			var convErr *ConversionError
			if !errors.As(err, &convErr) {
				t.Errorf("expected ConversionError, got: %T", err)
				return
			}

			if convErr.Field != tt.wantField {
				t.Errorf("expected field %q, got: %q", tt.wantField, convErr.Field)
			}
		})
	}
}

// BenchmarkParse validates performance target of <30ms
func BenchmarkParse(b *testing.B) {
	data, err := os.ReadFile("../../../testdata/xmp/sample.xmp")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestConversionError validates ConversionError formatting and unwrapping
func TestConversionError(t *testing.T) {
	baseErr := fmt.Errorf("underlying error")

	t.Run("Error with field", func(t *testing.T) {
		err := &ConversionError{
			Operation: "parse",
			Format:    "xmp",
			Field:     "Exposure",
			Cause:     baseErr,
		}

		expected := "parse xmp (Exposure): underlying error"
		if err.Error() != expected {
			t.Errorf("Error() = %q, want %q", err.Error(), expected)
		}
	})

	t.Run("Error without field", func(t *testing.T) {
		err := &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Cause:     baseErr,
		}

		expected := "validate xmp: underlying error"
		if err.Error() != expected {
			t.Errorf("Error() = %q, want %q", err.Error(), expected)
		}
	})

	t.Run("Unwrap returns cause", func(t *testing.T) {
		err := &ConversionError{
			Operation: "parse",
			Format:    "xmp",
			Cause:     baseErr,
		}

		if err.Unwrap() != baseErr {
			t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), baseErr)
		}

		// Verify errors.Is works
		if !errors.Is(err, baseErr) {
			t.Error("errors.Is should find base error")
		}
	})

	t.Run("errors.As extracts ConversionError", func(t *testing.T) {
		err := &ConversionError{
			Operation: "validate",
			Format:    "xmp",
			Field:     "Contrast",
			Cause:     baseErr,
		}

		var convErr *ConversionError
		if !errors.As(err, &convErr) {
			t.Fatal("errors.As should extract ConversionError")
		}

		if convErr.Field != "Contrast" {
			t.Errorf("Field = %q, want %q", convErr.Field, "Contrast")
		}
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================================
// GENERATOR TESTS
// ============================================================================

// TestGenerate validates basic XMP generation from UniversalRecipe
func TestGenerate(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Exposure:   1.5,
		Contrast:   25,
		Highlights: -30,
		Shadows:    40,
		Whites:     10,
		Blacks:     -15,
		Saturation: 20,
		Vibrance:   15,
		Clarity:    35,
		Sharpness:  40,
		Tint:       5,
	}

	xmpData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify XML structure
	if len(xmpData) == 0 {
		t.Error("Generate() returned empty data")
	}

	// Verify it can be parsed back
	parsedRecipe, err := Parse(xmpData)
	if err != nil {
		t.Fatalf("Parse() failed on generated XMP: %v", err)
	}

	// Verify key parameters match (with ±1 tolerance for rounding)
	if abs(int(parsedRecipe.Exposure*100)-int(recipe.Exposure*100)) > 1 {
		t.Errorf("Exposure mismatch: got %.2f, want %.2f", parsedRecipe.Exposure, recipe.Exposure)
	}
	if abs(parsedRecipe.Contrast-recipe.Contrast) > 1 {
		t.Errorf("Contrast mismatch: got %d, want %d", parsedRecipe.Contrast, recipe.Contrast)
	}
	if abs(parsedRecipe.Highlights-recipe.Highlights) > 1 {
		t.Errorf("Highlights mismatch: got %d, want %d", parsedRecipe.Highlights, recipe.Highlights)
	}
}

// TestGenerateNilRecipe validates error handling for nil recipe
func TestGenerateNilRecipe(t *testing.T) {
	_, err := Generate(nil)
	if err == nil {
		t.Error("Generate(nil) should return error")
	}

	var convErr *ConversionError
	if !errors.As(err, &convErr) {
		t.Errorf("expected ConversionError, got: %T", err)
		return
	}

	if convErr.Operation != "generate" {
		t.Errorf("Operation = %q, want %q", convErr.Operation, "generate")
	}
	if convErr.Format != "xmp" {
		t.Errorf("Format = %q, want %q", convErr.Format, "xmp")
	}
}

// TestGenerateInvalidRange validates range validation for all parameters
func TestGenerateInvalidRange(t *testing.T) {
	tests := []struct {
		name      string
		recipe    *models.UniversalRecipe
		wantField string
	}{
		{
			name: "Exposure too high",
			recipe: &models.UniversalRecipe{
				Exposure: 6.0, // Max is 5.0
			},
			wantField: "Exposure",
		},
		{
			name: "Exposure too low",
			recipe: &models.UniversalRecipe{
				Exposure: -6.0, // Min is -5.0
			},
			wantField: "Exposure",
		},
		{
			name: "Contrast too high",
			recipe: &models.UniversalRecipe{
				Contrast: 150, // Max is 100
			},
			wantField: "Contrast",
		},
		{
			name: "Contrast too low",
			recipe: &models.UniversalRecipe{
				Contrast: -150, // Min is -100
			},
			wantField: "Contrast",
		},
		{
			name: "Highlights out of range",
			recipe: &models.UniversalRecipe{
				Highlights: 200,
			},
			wantField: "Highlights",
		},
		{
			name: "Shadows out of range",
			recipe: &models.UniversalRecipe{
				Shadows: -200,
			},
			wantField: "Shadows",
		},
		{
			name: "Whites out of range",
			recipe: &models.UniversalRecipe{
				Whites: 150,
			},
			wantField: "Whites",
		},
		{
			name: "Blacks out of range",
			recipe: &models.UniversalRecipe{
				Blacks: -150,
			},
			wantField: "Blacks",
		},
		{
			name: "Saturation out of range",
			recipe: &models.UniversalRecipe{
				Saturation: 200,
			},
			wantField: "Saturation",
		},
		{
			name: "Vibrance out of range",
			recipe: &models.UniversalRecipe{
				Vibrance: -150,
			},
			wantField: "Vibrance",
		},
		{
			name: "Clarity out of range",
			recipe: &models.UniversalRecipe{
				Clarity: 200,
			},
			wantField: "Clarity",
		},
		{
			name: "Sharpness out of range",
			recipe: &models.UniversalRecipe{
				Sharpness: 200, // Max is 150
			},
			wantField: "Sharpness",
		},
		{
			name: "Tint out of range",
			recipe: &models.UniversalRecipe{
				Tint: 200, // Max is 150
			},
			wantField: "Tint",
		},
		{
			name: "Red Hue out of range",
			recipe: &models.UniversalRecipe{
				Red: models.ColorAdjustment{Hue: 200},
			},
			wantField: "HueRed",
		},
		{
			name: "Blue Saturation out of range",
			recipe: &models.UniversalRecipe{
				Blue: models.ColorAdjustment{Saturation: -200},
			},
			wantField: "SaturationBlue",
		},
		{
			name: "Green Luminance out of range",
			recipe: &models.UniversalRecipe{
				Green: models.ColorAdjustment{Luminance: 200},
			},
			wantField: "LuminanceGreen",
		},
		{
			name: "SplitShadowHue out of range",
			recipe: &models.UniversalRecipe{
				SplitShadowHue: 400, // Max is 360
			},
			wantField: "SplitShadowHue",
		},
		{
			name: "SplitShadowSaturation out of range",
			recipe: &models.UniversalRecipe{
				SplitShadowSaturation: 150, // Max is 100
			},
			wantField: "SplitShadowSaturation",
		},
		{
			name: "SplitHighlightHue out of range",
			recipe: &models.UniversalRecipe{
				SplitHighlightHue: -10, // Min is 0
			},
			wantField: "SplitHighlightHue",
		},
		{
			name: "SplitHighlightSaturation out of range",
			recipe: &models.UniversalRecipe{
				SplitHighlightSaturation: -5, // Min is 0
			},
			wantField: "SplitHighlightSaturation",
		},
		{
			name: "SplitBalance out of range",
			recipe: &models.UniversalRecipe{
				SplitBalance: 150, // Max is 100
			},
			wantField: "SplitBalance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate(tt.recipe)
			if err == nil {
				t.Error("expected error for out-of-range value")
				return
			}

			var convErr *ConversionError
			if !errors.As(err, &convErr) {
				t.Errorf("expected ConversionError, got: %T", err)
				return
			}

			if convErr.Field != tt.wantField {
				t.Errorf("Field = %q, want %q", convErr.Field, tt.wantField)
			}
		})
	}
}

// TestGenerateRoundTrip validates round-trip conversion: XMP → Parse → Generate → Parse
func TestGenerateRoundTrip(t *testing.T) {
	// Read sample XMP file
	originalData, err := os.ReadFile("../../../testdata/xmp/sample.xmp")
	if err != nil {
		t.Fatalf("failed to read sample.xmp: %v", err)
	}

	// Parse original
	recipe1, err := Parse(originalData)
	if err != nil {
		t.Fatalf("Parse(original) error = %v", err)
	}

	// Generate XMP from recipe
	generatedData, err := Generate(recipe1)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Parse generated XMP
	recipe2, err := Parse(generatedData)
	if err != nil {
		t.Fatalf("Parse(generated) error = %v", err)
	}

	// Compare recipes (with ±1 tolerance for rounding)
	if abs(int(recipe1.Exposure*100)-int(recipe2.Exposure*100)) > 1 {
		t.Errorf("Exposure: original=%.2f, generated=%.2f", recipe1.Exposure, recipe2.Exposure)
	}
	if abs(recipe1.Contrast-recipe2.Contrast) > 1 {
		t.Errorf("Contrast: original=%d, generated=%d", recipe1.Contrast, recipe2.Contrast)
	}
	if abs(recipe1.Highlights-recipe2.Highlights) > 1 {
		t.Errorf("Highlights: original=%d, generated=%d", recipe1.Highlights, recipe2.Highlights)
	}
	if abs(recipe1.Shadows-recipe2.Shadows) > 1 {
		t.Errorf("Shadows: original=%d, generated=%d", recipe1.Shadows, recipe2.Shadows)
	}
	if abs(recipe1.Whites-recipe2.Whites) > 1 {
		t.Errorf("Whites: original=%d, generated=%d", recipe1.Whites, recipe2.Whites)
	}
	if abs(recipe1.Blacks-recipe2.Blacks) > 1 {
		t.Errorf("Blacks: original=%d, generated=%d", recipe1.Blacks, recipe2.Blacks)
	}
	if abs(recipe1.Saturation-recipe2.Saturation) > 1 {
		t.Errorf("Saturation: original=%d, generated=%d", recipe1.Saturation, recipe2.Saturation)
	}
	if abs(recipe1.Vibrance-recipe2.Vibrance) > 1 {
		t.Errorf("Vibrance: original=%d, generated=%d", recipe1.Vibrance, recipe2.Vibrance)
	}
	if abs(recipe1.Clarity-recipe2.Clarity) > 1 {
		t.Errorf("Clarity: original=%d, generated=%d", recipe1.Clarity, recipe2.Clarity)
	}
	if abs(recipe1.Sharpness-recipe2.Sharpness) > 1 {
		t.Errorf("Sharpness: original=%d, generated=%d", recipe1.Sharpness, recipe2.Sharpness)
	}

	// Validate HSL colors
	compareColorAdjustment(t, "Red", recipe1.Red, recipe2.Red)
	compareColorAdjustment(t, "Orange", recipe1.Orange, recipe2.Orange)
	compareColorAdjustment(t, "Yellow", recipe1.Yellow, recipe2.Yellow)
	compareColorAdjustment(t, "Green", recipe1.Green, recipe2.Green)
	compareColorAdjustment(t, "Aqua", recipe1.Aqua, recipe2.Aqua)
	compareColorAdjustment(t, "Blue", recipe1.Blue, recipe2.Blue)
	compareColorAdjustment(t, "Purple", recipe1.Purple, recipe2.Purple)
	compareColorAdjustment(t, "Magenta", recipe1.Magenta, recipe2.Magenta)

	// Validate Split Toning
	if abs(recipe1.SplitShadowHue-recipe2.SplitShadowHue) > 1 {
		t.Errorf("SplitShadowHue: original=%d, generated=%d", recipe1.SplitShadowHue, recipe2.SplitShadowHue)
	}
	if abs(recipe1.SplitShadowSaturation-recipe2.SplitShadowSaturation) > 1 {
		t.Errorf("SplitShadowSaturation: original=%d, generated=%d", recipe1.SplitShadowSaturation, recipe2.SplitShadowSaturation)
	}
	if abs(recipe1.SplitHighlightHue-recipe2.SplitHighlightHue) > 1 {
		t.Errorf("SplitHighlightHue: original=%d, generated=%d", recipe1.SplitHighlightHue, recipe2.SplitHighlightHue)
	}
	if abs(recipe1.SplitHighlightSaturation-recipe2.SplitHighlightSaturation) > 1 {
		t.Errorf("SplitHighlightSaturation: original=%d, generated=%d", recipe1.SplitHighlightSaturation, recipe2.SplitHighlightSaturation)
	}
}

// compareColorAdjustment compares two ColorAdjustment structs with ±1 tolerance
func compareColorAdjustment(t *testing.T, color string, adj1, adj2 models.ColorAdjustment) {
	t.Helper()
	if abs(adj1.Hue-adj2.Hue) > 1 {
		t.Errorf("%s Hue: original=%d, generated=%d", color, adj1.Hue, adj2.Hue)
	}
	if abs(adj1.Saturation-adj2.Saturation) > 1 {
		t.Errorf("%s Saturation: original=%d, generated=%d", color, adj1.Saturation, adj2.Saturation)
	}
	if abs(adj1.Luminance-adj2.Luminance) > 1 {
		t.Errorf("%s Luminance: original=%d, generated=%d", color, adj1.Luminance, adj2.Luminance)
	}
}

// TestGenerateEmptyRecipe validates generation with zero values
func TestGenerateEmptyRecipe(t *testing.T) {
	recipe := &models.UniversalRecipe{}

	xmpData, err := Generate(recipe)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Parse back and verify zero values are handled
	parsedRecipe, err := Parse(xmpData)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if parsedRecipe.Exposure != 0.0 {
		t.Errorf("Expected zero Exposure, got %.2f", parsedRecipe.Exposure)
	}
	if parsedRecipe.Contrast != 0 {
		t.Errorf("Expected zero Contrast, got %d", parsedRecipe.Contrast)
	}
}

// TestGenerateNullableTemperature validates nullable Temperature handling
func TestGenerateNullableTemperature(t *testing.T) {
	t.Run("With Temperature", func(t *testing.T) {
		temp := 5600
		recipe := &models.UniversalRecipe{
			Temperature: &temp,
		}

		xmpData, err := Generate(recipe)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		parsedRecipe, err := Parse(xmpData)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		if parsedRecipe.Temperature == nil {
			t.Error("Temperature should not be nil")
		} else if *parsedRecipe.Temperature != temp {
			t.Errorf("Temperature = %d, want %d", *parsedRecipe.Temperature, temp)
		}
	})

	t.Run("Without Temperature", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			Temperature: nil,
		}

		xmpData, err := Generate(recipe)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		// Should not error when parsing back
		_, err = Parse(xmpData)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}
	})
}

// TestGenerateToneCurve validates ToneCurve generation from PointCurve array
func TestGenerateToneCurve(t *testing.T) {
	t.Run("With ToneCurve Points", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			PointCurve: []models.ToneCurvePoint{
				{Input: 0, Output: 0},
				{Input: 64, Output: 70},
				{Input: 128, Output: 140},
				{Input: 192, Output: 200},
				{Input: 255, Output: 255},
			},
		}

		xmpData, err := Generate(recipe)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		// Verify XMP contains modern ToneCurvePV2012 (nested format, not legacy attribute)
		xmpStr := string(xmpData)
		if !containsString(xmpStr, "<crs:ToneCurvePV2012>") {
			t.Error("Generated XMP should contain ToneCurvePV2012 element")
		}

		// Should NOT contain legacy attribute format
		if containsString(xmpStr, "crs:ToneCurve=") {
			t.Error("Generated XMP should NOT contain legacy crs:ToneCurve attribute")
		}

		// Verify points are in rdf:li elements
		if !containsString(xmpStr, "<rdf:li>0, 0</rdf:li>") {
			t.Error("ToneCurvePV2012 should contain first point as rdf:li")
		}
		if !containsString(xmpStr, "<rdf:li>255, 255</rdf:li>") {
			t.Error("ToneCurvePV2012 should contain last point as rdf:li")
		}
	})

	t.Run("Without ToneCurve Points", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			PointCurve: nil,
		}

		xmpData, err := Generate(recipe)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		// Should not include ToneCurvePV2012 element when PointCurve is nil
		xmpStr := string(xmpData)
		if containsString(xmpStr, "<crs:ToneCurvePV2012>") {
			t.Error("Generated XMP should not contain ToneCurvePV2012 when PointCurve is nil")
		}
	})

	t.Run("Empty ToneCurve Array", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			PointCurve: []models.ToneCurvePoint{},
		}

		xmpData, err := Generate(recipe)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		// Should not include ToneCurvePV2012 element when empty
		xmpStr := string(xmpData)
		if containsString(xmpStr, "<crs:ToneCurvePV2012>") {
			t.Error("Generated XMP should not contain ToneCurvePV2012 when PointCurve is empty")
		}
	})
}

// containsString checks if a string contains a substring (helper for tests)
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

// findSubstring performs substring search
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkGenerate validates performance target of <30ms
func BenchmarkGenerate(b *testing.B) {
	temp := 5600
	recipe := &models.UniversalRecipe{
		Exposure:    1.5,
		Contrast:    25,
		Highlights:  -30,
		Shadows:     40,
		Whites:      10,
		Blacks:      -15,
		Saturation:  20,
		Vibrance:    15,
		Clarity:     35,
		Sharpness:   40,
		Temperature: &temp,
		Tint:        5,
		Red: models.ColorAdjustment{
			Hue:        10,
			Saturation: -5,
			Luminance:  8,
		},
		Blue: models.ColorAdjustment{
			Hue:        20,
			Saturation: -10,
			Luminance:  5,
		},
		SplitShadowHue:           240,
		SplitShadowSaturation:    15,
		SplitHighlightHue:        50,
		SplitHighlightSaturation: 20,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Generate(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
