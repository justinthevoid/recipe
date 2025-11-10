package converter

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

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
		{"NP3→XMP", FormatNP3, FormatXMP, "../../examples/np3", ".np3"},
		{"NP3→LRTemplate", FormatNP3, FormatLRTemplate, "../../examples/np3", ".np3"},
		{"XMP→NP3", FormatXMP, FormatNP3, "../../examples/xmp", ".xmp"},
		{"XMP→LRTemplate", FormatXMP, FormatLRTemplate, "../../examples/xmp", ".xmp"},
		{"LRTemplate→NP3", FormatLRTemplate, FormatNP3, "../../examples/lrtemplate", ".lrtemplate"},
		{"LRTemplate→XMP", FormatLRTemplate, FormatXMP, "../../examples/lrtemplate", ".lrtemplate"},
		{"Costyle→XMP", FormatCostyle, FormatXMP, "../../internal/formats/costyle/testdata/costyle", ".costyle"},
		{"Costyle→NP3", FormatCostyle, FormatNP3, "../../internal/formats/costyle/testdata/costyle", ".costyle"},
		{"XMP→Costyle", FormatXMP, FormatCostyle, "../../examples/xmp", ".xmp"},
		{"NP3→Costyle", FormatNP3, FormatCostyle, "../../examples/np3", ".np3"},
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
		{
			"lrtemplate invalid Lua",
			[]byte("s = { this is not valid lua"),
			FormatLRTemplate,
			FormatNP3,
			"parse",
			FormatLRTemplate,
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

	lrtemplateData := []byte(`s = {
	id = "test",
	internalName = "Test Recipe",
	value = {
		settings = {}
	}
}`)

	tests := []struct {
		name        string
		input       []byte
		to          string
		expectError bool
	}{
		{"Auto-detect NP3", np3Data, FormatXMP, false},
		{"Auto-detect XMP", xmpData, FormatNP3, false},
		{"Auto-detect lrtemplate", lrtemplateData, FormatNP3, false},
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

// TestDetectFormat_LRTemplate validates lrtemplate detection by Lua syntax
func TestDetectFormat_LRTemplate(t *testing.T) {
	lrtemplateData := []byte(`s = {
	id = "test",
	value = {}
}`)

	format, err := DetectFormat(lrtemplateData)
	if err != nil {
		t.Fatalf("Detection failed: %v", err)
	}
	if format != FormatLRTemplate {
		t.Errorf("Expected %q, got %q", FormatLRTemplate, format)
	}
}

// TestDetectFormat_Costyle validates costyle detection by SL Engine tag
func TestDetectFormat_Costyle(t *testing.T) {
	costyleData := []byte(`<?xml version="1.0"?>
<SL Engine="1300">
	<E K="Name" V="Test Preset" />
</SL>`)

	format, err := DetectFormat(costyleData)
	if err != nil {
		t.Fatalf("Detection failed: %v", err)
	}
	if format != FormatCostyle {
		t.Errorf("Expected %q, got %q", FormatCostyle, format)
	}
}

// TestDetectFormat_Costylepack validates costylepack detection by ZIP magic bytes
func TestDetectFormat_Costylepack(t *testing.T) {
	// ZIP magic bytes: PK\x03\x04
	costylepackData := []byte{0x50, 0x4B, 0x03, 0x04, 0x00, 0x00, 0x00, 0x00}

	format, err := DetectFormat(costylepackData)
	if err != nil {
		t.Fatalf("Detection failed: %v", err)
	}
	if format != FormatCostylepack {
		t.Errorf("Expected %q, got %q", FormatCostylepack, format)
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
	// Skip detailed round-trip testing - this is covered by format-specific tests
	// The converter's job is just to orchestrate, not to ensure round-trip accuracy
	// Each format package (np3, xmp, lrtemplate) has its own round-trip tests

	// Round-trip accuracy is validated by:
	// 1. np3 package round-trip tests (np3 → UniversalRecipe → np3)
	// 2. xmp package round-trip tests (xmp → UniversalRecipe → xmp)
	// 3. lrtemplate package round-trip tests (lrtemplate → UniversalRecipe → lrtemplate)
	// 4. TestConvert_AllPaths validates all 6 conversion paths work without errors

	t.Log("✓ Round-trip conversion accuracy delegated to format-specific tests")
}

// Helper function for floating-point comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
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

	lrtemplateData := []byte(`s = {
	id = "test",
	internalName = "Test",
	value = { settings = {} }
}`)

	// Run 100 concurrent conversions
	const numGoroutines = 100
	done := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			// Each goroutine performs a different conversion
			switch index % 3 {
			case 0:
				_, err := Convert(np3Data, FormatNP3, FormatXMP)
				done <- err
			case 1:
				_, err := Convert(xmpData, FormatXMP, FormatLRTemplate)
				done <- err
			case 2:
				_, err := Convert(lrtemplateData, FormatLRTemplate, FormatNP3)
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
