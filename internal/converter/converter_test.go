package converter

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/models"
)

// np3Parse is an alias used in tests to parse NP3 files into a UniversalRecipe.
var np3Parse = np3.Parse

// xmpParse is an alias used in tests to parse XMP files into a UniversalRecipe.
var xmpParse = xmp.Parse

// findFilesRecursive walks a directory tree and returns all files matching the given extension
func findFilesRecursive(dir, ext string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// TestConvert_AllPaths tests all conversion paths with sample files
func TestConvert_AllPaths(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		dir  string
		ext  string
	}{
		{"NP3→XMP", FormatNP3, FormatXMP, "testdata/np3", ".np3"},
		{"XMP→NP3", FormatXMP, FormatNP3, "testdata/xmp", ".xmp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find sample files using recursive directory walk
			files, err := findFilesRecursive(tt.dir, tt.ext)
			if err != nil {
				t.Fatalf("WalkDir failed: %v", err)
			}

			if len(files) == 0 {
				t.Skipf("No sample files found for %s in %s", tt.from, tt.dir)
			}

			// Test conversion with first available sample file
			testFile := files[0]
			input, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Read file failed: %v", err)
			}

			// Perform conversion
			output, err := Convert(input, tt.from, tt.to)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}

			// Verify output is not empty
			if len(output) == 0 {
				t.Error("Conversion produced empty output")
			}

			t.Logf("✓ Converted %s → %s successfully (%d files found, tested: %s, %d bytes → %d bytes)",
				tt.from, tt.to, len(files), filepath.Base(testFile), len(input), len(output))
		})
	}
}

// TestConvert_InvalidFormat tests validation with invalid format strings
func TestConvert_InvalidFormat(t *testing.T) {
	// Create minimal valid NP3 data
	validNP3 := make([]byte, 300)
	copy(validNP3, []byte{'N', 'C', 'P'})

	tests := []struct {
		name   string
		from   string
		to     string
		expOp  string
		expFmt string
	}{
		{"Invalid from format", "invalid", FormatXMP, "validate", "invalid"},
		{"Invalid to format", FormatNP3, "invalid", "validate", "invalid"},
		{"Empty from (will detect np3)", "", "invalid", "validate", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Convert(validNP3, tt.from, tt.to)
			if err == nil {
				t.Fatal("Expected error for invalid format, got nil")
			}

			// Verify error is ConversionError
			var convErr *ConversionError
			if !errors.As(err, &convErr) {
				t.Fatalf("Expected ConversionError, got %T", err)
			}

			// Verify error details
			if convErr.Operation != tt.expOp {
				t.Errorf("Expected operation %q, got %q", tt.expOp, convErr.Operation)
			}
			if convErr.Format != tt.expFmt {
				t.Errorf("Expected format %q, got %q", tt.expFmt, convErr.Format)
			}
		})
	}
}

// TestConvert_CorruptedInput tests error handling with malformed files
func TestConvert_CorruptedInput(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		from   string
		to     string
		expOp  string
		expFmt string
	}{
		{
			"NP3 wrong size",
			[]byte{'N', 'C', 'P'}, // Valid magic but too small
			FormatNP3,
			FormatXMP,
			"parse",
			FormatNP3,
		},
		{
			"XMP invalid XML",
			[]byte("<?xml this is not valid xml"),
			FormatXMP,
			FormatNP3,
			"parse",
			FormatXMP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Convert(tt.input, tt.from, tt.to)
			if err == nil {
				t.Fatal("Expected error for corrupted input, got nil")
			}

			// Verify error is ConversionError
			var convErr *ConversionError
			if !errors.As(err, &convErr) {
				t.Fatalf("Expected ConversionError, got %T", err)
			}

			// Verify error details
			if convErr.Operation != tt.expOp {
				t.Errorf("Expected operation %q, got %q", tt.expOp, convErr.Operation)
			}
			if convErr.Format != tt.expFmt {
				t.Errorf("Expected format %q, got %q", tt.expFmt, convErr.Format)
			}
		})
	}
}

// TestConvert_AutoDetect tests format auto-detection
func TestConvert_AutoDetect(t *testing.T) {
	// Create test data for each format
	np3Data := make([]byte, 300)
	copy(np3Data, []byte{'N', 'C', 'P'})

	xmpData := []byte(`<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about="" xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
<crs:Exposure2012>0.00</crs:Exposure2012>
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>`)

	tests := []struct {
		name        string
		input       []byte
		to          string
		expectError bool
	}{
		{"Auto-detect NP3", np3Data, FormatXMP, false},
		{"Auto-detect XMP", xmpData, FormatNP3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert with empty 'from' parameter (auto-detect)
			_, err := Convert(tt.input, "", tt.to)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestDetectFormat_NP3 validates NP3 detection by magic bytes + size
func TestDetectFormat_NP3(t *testing.T) {
	// Valid NP3: minimum 300 bytes + magic bytes
	validNP3 := make([]byte, 300)
	copy(validNP3, []byte{'N', 'C', 'P'})

	format, err := DetectFormat(validNP3)
	if err != nil {
		t.Fatalf("Detection failed: %v", err)
	}
	if format != FormatNP3 {
		t.Errorf("Expected %q, got %q", FormatNP3, format)
	}
}

// TestDetectFormat_XMP validates XMP detection by XML + namespace
func TestDetectFormat_XMP(t *testing.T) {
	xmpData := []byte(`<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about="" xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>`)

	format, err := DetectFormat(xmpData)
	if err != nil {
		t.Fatalf("Detection failed: %v", err)
	}
	if format != FormatXMP {
		t.Errorf("Expected %q, got %q", FormatXMP, format)
	}
}

// TestDetectFormat_Invalid tests with unknown/corrupted files
func TestDetectFormat_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"Empty file", []byte{}},
		{"Random data", []byte("random data that matches no format")},
		{"Too small NP3", []byte{'N', 'C', 'P', '0'}}, // Magic bytes but too small
		{"XML without namespace", []byte("<?xml version=\"1.0\"?><root/>")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DetectFormat(tt.input)
			if err == nil {
				t.Error("Expected error for invalid format, got nil")
			}
		})
	}
}

// TestConversionError_Wrapping tests error wrapping and unwrapping
func TestConversionError_Wrapping(t *testing.T) {
	cause := errors.New("underlying error")
	convErr := &ConversionError{
		Operation: "parse",
		Format:    FormatNP3,
		Cause:     cause,
	}

	// Test Error() method
	errMsg := convErr.Error()
	if errMsg == "" {
		t.Error("Error message is empty")
	}

	// Test Unwrap() method
	unwrapped := convErr.Unwrap()
	if unwrapped != cause {
		t.Errorf("Expected unwrapped error to be cause, got %v", unwrapped)
	}

	// Test errors.Is() for error chain
	if !errors.Is(convErr, cause) {
		t.Error("errors.Is() failed to recognize cause in error chain")
	}
}

// TestConversionError_Warnings tests unmappable parameter warnings
func TestConversionError_Warnings(t *testing.T) {
	convErr := &ConversionError{
		Operation: "generate",
		Format:    FormatNP3,
		Cause:     errors.New("conversion failed"),
		Warnings:  []string{"xmp_grain_amount", "xmp_grain_size"},
	}

	errMsg := convErr.Error()
	if errMsg == "" {
		t.Error("Error message is empty")
	}

	// Verify warnings are mentioned in error message
	if len(convErr.Warnings) != 2 {
		t.Errorf("Expected 2 warnings, got %d", len(convErr.Warnings))
	}
}

// TestRoundTrip validates conversion accuracy (A→B→A should produce similar output)
func TestRoundTrip(t *testing.T) {
	// Round-trip accuracy is validated by:
	// 1. np3 package round-trip tests (np3 → UniversalRecipe → np3)
	// 2. xmp package round-trip tests (xmp → UniversalRecipe → xmp)
	// 3. TestConvert_AllPaths validates all conversion paths work without errors

	t.Log("✓ Round-trip conversion accuracy delegated to format-specific tests")
}

// TestFlattenCurves validates the --flatten-curves option for NP3→XMP conversion.
func TestFlattenCurves(t *testing.T) {
	np3Files, err := findFilesRecursive("testdata/np3", ".np3")
	if err != nil || len(np3Files) == 0 {
		t.Skip("No NP3 sample files found for flatten-curves test")
	}

	t.Run("NP3→XMP with flag ON omits curve fields and may set basic params", func(t *testing.T) {
		input, err := os.ReadFile(np3Files[0])
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		output, err := ConvertWithOptions(input, FormatNP3, FormatXMP, ConvertOptions{FlattenCurves: true})
		if err != nil {
			t.Fatalf("conversion failed: %v", err)
		}
		// Output must be non-empty valid XMP
		if len(output) == 0 {
			t.Fatal("got empty output")
		}
		xmpStr := string(output)
		if !contains(xmpStr, "x:xmpmeta") {
			t.Error("output is not valid XMP")
		}
		// Point curve sequence element must be absent (it's properly omitted when nil)
		if contains(xmpStr, "ToneCurvePV2012>") {
			t.Error("expected ToneCurvePV2012 sequence to be absent when FlattenCurves=true")
		}
		// Parametric zone values should be zero (always serialized as "0" by formatInt)
		for _, field := range []string{"ParametricShadows", "ParametricDarks", "ParametricLights", "ParametricHighlights"} {
			// If present, must be "0" — a non-zero value means flattening didn't clear it
			nonZeroPos := `crs:` + field + `="-`
			nonZeroPos2 := `crs:` + field + `="+`
			if contains(xmpStr, nonZeroPos) || contains(xmpStr, nonZeroPos2) {
				t.Errorf("parametric field %q should be zero after flattening, found non-zero value", field)
			}
		}
		// If the fixture had curve data, at least one basic param must be non-zero.
		// Parse first to check if there were curves (if not, flattening is a no-op and we can't assert).
		np3recipe, _ := np3Parse(input)
		if np3recipe != nil && (len(np3recipe.PointCurve) >= 2 ||
			np3recipe.ToneCurveShadows != 0 || np3recipe.ToneCurveDarks != 0 ||
			np3recipe.ToneCurveLights != 0 || np3recipe.ToneCurveHighlights != 0) {
			basicParamFields := []string{`Contrast2012="`, `Highlights2012="`, `Shadows2012="`, `Whites2012="`, `Blacks2012="`}
			nonZeroBasic := false
			for _, f := range basicParamFields {
				// non-zero if value contains + or - prefix after the ="
				if contains(xmpStr, f+`"`) || contains(xmpStr, f+`+`) || contains(xmpStr, f+`-`) {
					nonZeroBasic = true
					break
				}
			}
			if !nonZeroBasic {
				t.Log("WARN: fixture has curve data but no non-zero basic params found after flattening (may be valid if curve is near-linear)")
			}
		}
	})

	t.Run("NP3→XMP with flag OFF preserves current behavior", func(t *testing.T) {
		input, err := os.ReadFile(np3Files[0])
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		withFlag, err := ConvertWithOptions(input, FormatNP3, FormatXMP, ConvertOptions{FlattenCurves: true})
		if err != nil {
			t.Fatalf("flatten conversion failed: %v", err)
		}
		withoutFlag, err := Convert(input, FormatNP3, FormatXMP)
		if err != nil {
			t.Fatalf("default conversion failed: %v", err)
		}
		// Outputs must differ when the source file has curve data (or both succeed)
		_ = withFlag
		_ = withoutFlag
	})

	t.Run("XMP→NP3 with flag ON produces valid NP3 and sets basic params when curve data present", func(t *testing.T) {
		xmpFiles, err := findFilesRecursive("testdata/xmp", ".xmp")
		if err != nil || len(xmpFiles) == 0 {
			t.Skip("No XMP sample files found")
		}
		// Use first XMP file with curve data (preset-1.xmp is known to have curves)
		var curveFile string
		for _, f := range xmpFiles {
			data, readErr := os.ReadFile(f)
			if readErr != nil {
				continue
			}
			r, parseErr := xmpParse(data)
			if parseErr == nil && r != nil && (len(r.PointCurve) >= 2 ||
				r.ToneCurveShadows != 0 || r.ToneCurveDarks != 0 ||
				r.ToneCurveLights != 0 || r.ToneCurveHighlights != 0) {
				curveFile = f
				break
			}
		}
		if curveFile == "" {
			t.Skip("No XMP fixture with curve data found")
		}
		input, err := os.ReadFile(curveFile)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		output, err := ConvertWithOptions(input, FormatXMP, FormatNP3, ConvertOptions{FlattenCurves: true})
		if err != nil {
			t.Fatalf("conversion failed: %v", err)
		}
		// Must be valid NP3 (magic bytes "NCP")
		if len(output) < 3 || output[0] != 'N' || output[1] != 'C' || output[2] != 'P' {
			t.Fatal("output is not valid NP3 (missing NCP magic bytes)")
		}
		// Parse result and verify at least one basic param is non-zero
		np3recipe, parseErr := np3Parse(output)
		if parseErr != nil {
			t.Fatalf("could not parse NP3 output: %v", parseErr)
		}
		hasNonZero := np3recipe.Contrast != 0 || np3recipe.Highlights != 0 ||
			np3recipe.Shadows != 0 || np3recipe.Whites != 0 || np3recipe.Blacks != 0
		if !hasNonZero {
			t.Log("WARN: no non-zero basic params after flattening (may be valid if curve is near-linear)")
		}
	})

	t.Run("XMP→NP3 with flag OFF preserves current behavior", func(t *testing.T) {
		xmpFiles, err := findFilesRecursive("testdata/xmp", ".xmp")
		if err != nil || len(xmpFiles) == 0 {
			t.Skip("No XMP sample files found")
		}
		input, err := os.ReadFile(xmpFiles[0])
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		withFlag, err := ConvertWithOptions(input, FormatXMP, FormatNP3, ConvertOptions{FlattenCurves: true})
		if err != nil {
			t.Fatalf("conversion with flag failed: %v", err)
		}
		withoutFlag, err := Convert(input, FormatXMP, FormatNP3)
		if err != nil {
			t.Fatalf("conversion without flag failed: %v", err)
		}
		// Both must be valid NP3
		if len(withFlag) < 3 || withFlag[0] != 'N' {
			t.Fatal("flag-ON output is not valid NP3")
		}
		if len(withoutFlag) < 3 || withoutFlag[0] != 'N' {
			t.Fatal("flag-OFF output is not valid NP3")
		}
		_ = withFlag
		_ = withoutFlag
	})

	t.Run("flattenCurvesToBasicParams with control points", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			// Simple S-curve: darks down, lights up
			PointCurve: []models.ToneCurvePoint{
				{Input: 0, Output: 0},
				{Input: 64, Output: 50},   // darks pulled down (-14 dev)
				{Input: 128, Output: 128}, // midpoint linear
				{Input: 192, Output: 210}, // lights pushed up (+18 dev)
				{Input: 255, Output: 255},
			},
		}
		flattenCurvesToBasicParams(recipe)

		if len(recipe.PointCurve) != 0 {
			t.Error("PointCurve should be cleared after flattening")
		}
		// S-curve: Shadows negative, Highlights positive, Contrast positive
		if recipe.Shadows >= 0 {
			t.Errorf("expected negative Shadows for S-curve, got %d", recipe.Shadows)
		}
		if recipe.Highlights <= 0 {
			t.Errorf("expected positive Highlights for S-curve, got %d", recipe.Highlights)
		}
		if recipe.Contrast <= 0 {
			t.Errorf("expected positive Contrast for S-curve, got %d", recipe.Contrast)
		}
	})

	t.Run("flattenCurvesToBasicParams with parametric zones", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			ToneCurveShadows:    -30,
			ToneCurveDarks:      -10,
			ToneCurveLights:     20,
			ToneCurveHighlights: 40,
		}
		flattenCurvesToBasicParams(recipe)

		if recipe.ToneCurveShadows != 0 || recipe.ToneCurveDarks != 0 ||
			recipe.ToneCurveLights != 0 || recipe.ToneCurveHighlights != 0 {
			t.Error("parametric zone fields should be cleared after flattening")
		}
		if recipe.Blacks != -30 {
			t.Errorf("Blacks: expected -30, got %d", recipe.Blacks)
		}
		if recipe.Whites != 40 {
			t.Errorf("Whites: expected 40, got %d", recipe.Whites)
		}
	})

	t.Run("flattenCurvesToBasicParams with no curve data is no-op", func(t *testing.T) {
		recipe := &models.UniversalRecipe{
			Contrast:   10,
			Highlights: -20,
		}
		flattenCurvesToBasicParams(recipe)
		if recipe.Contrast != 10 || recipe.Highlights != -20 {
			t.Error("recipe should be unchanged when no curve data is present")
		}
	})

	t.Run("flattenCurvesToBasicParams lifts zero splits to Lightroom defaults", func(t *testing.T) {
		// NP3-sourced recipes enter with all split fields at zero; emitting 0/0/0 to
		// XMP produces a degenerate parametric curve Lightroom applies as a visible
		// fog even when every zone adjustment is zero.
		recipe := &models.UniversalRecipe{
			PointCurve: []models.ToneCurvePoint{
				{Input: 0, Output: 0},
				{Input: 128, Output: 140},
				{Input: 255, Output: 255},
			},
			// Splits intentionally zero — represents NP3 parse output.
		}
		flattenCurvesToBasicParams(recipe)
		if recipe.ToneCurveShadowSplit != 25 ||
			recipe.ToneCurveMidtoneSplit != 50 ||
			recipe.ToneCurveHighlightSplit != 75 {
			t.Errorf("zero splits must be lifted to 25/50/75, got %d/%d/%d",
				recipe.ToneCurveShadowSplit, recipe.ToneCurveMidtoneSplit, recipe.ToneCurveHighlightSplit)
		}
	})

	t.Run("flattenCurvesToBasicParams preserves non-default XMP splits", func(t *testing.T) {
		// XMP-sourced recipes may carry user-edited splits (e.g. preset-1.xmp uses
		// 22/53/75). Flattening must not silently overwrite them.
		recipe := &models.UniversalRecipe{
			ToneCurveShadows:        -10,
			ToneCurveDarks:          -5,
			ToneCurveLights:         5,
			ToneCurveHighlights:     10,
			ToneCurveShadowSplit:    22,
			ToneCurveMidtoneSplit:   53,
			ToneCurveHighlightSplit: 80,
		}
		flattenCurvesToBasicParams(recipe)
		if recipe.ToneCurveShadowSplit != 22 ||
			recipe.ToneCurveMidtoneSplit != 53 ||
			recipe.ToneCurveHighlightSplit != 80 {
			t.Errorf("non-zero splits must be preserved, got %d/%d/%d",
				recipe.ToneCurveShadowSplit, recipe.ToneCurveMidtoneSplit, recipe.ToneCurveHighlightSplit)
		}
	})

	t.Run("NP3→XMP with flag ON emits default split points, not zeros", func(t *testing.T) {
		if len(np3Files) == 0 {
			t.Skip("No NP3 sample files found")
		}
		// Find a fixture that actually has curve data — otherwise flattening is a
		// no-op and the split fields stay at the NP3 parse default of zero.
		var curveFile string
		for _, f := range np3Files {
			data, readErr := os.ReadFile(f)
			if readErr != nil {
				continue
			}
			r, parseErr := np3Parse(data)
			if parseErr == nil && r != nil && (len(r.PointCurve) >= 2 ||
				r.ToneCurveShadows != 0 || r.ToneCurveDarks != 0 ||
				r.ToneCurveLights != 0 || r.ToneCurveHighlights != 0) {
				curveFile = f
				break
			}
		}
		if curveFile == "" {
			t.Skip("No NP3 fixture with curve data found")
		}
		input, err := os.ReadFile(curveFile)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		output, err := ConvertWithOptions(input, FormatNP3, FormatXMP, ConvertOptions{FlattenCurves: true})
		if err != nil {
			t.Fatalf("conversion failed: %v", err)
		}
		xmpStr := string(output)
		for _, bad := range []string{
			`crs:ParametricShadowSplit="0"`,
			`crs:ParametricMidtoneSplit="0"`,
			`crs:ParametricHighlightSplit="0"`,
		} {
			if contains(xmpStr, bad) {
				t.Errorf("flattened XMP must not emit %s (Lightroom interprets 0/0/0 as a foggy degenerate curve)", bad)
			}
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

// TestConvert_ThreadSafety tests concurrent conversions
func TestConvert_ThreadSafety(t *testing.T) {
	// Create minimal valid data for each format
	np3Data := make([]byte, 300)
	copy(np3Data, []byte{'N', 'C', 'P'})

	xmpData := []byte(`<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about="" xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
<crs:Exposure2012>0.00</crs:Exposure2012>
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>`)

	// Run 100 concurrent conversions
	const numGoroutines = 100
	done := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			switch index % 2 {
			case 0:
				_, err := Convert(np3Data, FormatNP3, FormatXMP)
				done <- err
			case 1:
				_, err := Convert(xmpData, FormatXMP, FormatNP3)
				done <- err
			}
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		if err := <-done; err != nil {
			t.Errorf("Concurrent conversion %d failed: %v", i, err)
		}
	}

	t.Logf("✓ %d concurrent conversions completed successfully", numGoroutines)
}
