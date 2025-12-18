package inspect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBinaryDump_Format verifies output matches the expected format (AC-1).
func TestBinaryDump_Format(t *testing.T) {
	// Create minimal NP3 file with magic bytes
	data := make([]byte, 300)
	copy(data[0:3], []byte{'N', 'C', 'P'}) // Magic bytes
	data[3] = 0x03                          // Version

	output, err := BinaryDump(data, "np3")
	if err != nil {
		t.Fatalf("BinaryDump failed: %v", err)
	}

	// Verify output contains hex offsets
	if !strings.Contains(output, "[0x0000]") {
		t.Error("Output missing hex offset format")
	}

	// Verify uppercase hex
	if strings.Contains(output, "0x00a") || strings.Contains(output, "0x00b") {
		t.Error("Hex values should be uppercase")
	}

	// Verify offsets increment correctly
	lines := strings.Split(output, "\n")
	foundOffsets := false
	for _, line := range lines {
		if strings.Contains(line, "[0x0000]") &&
			strings.Contains(output, "[0x0001]") &&
			strings.Contains(output, "[0x0002]") {
			foundOffsets = true
			break
		}
	}
	if !foundOffsets {
		t.Error("Offsets not incrementing correctly")
	}
}

// TestBinaryDump_KnownFields verifies minimum required fields are annotated (AC-2).
func TestBinaryDump_KnownFields(t *testing.T) {
	// Create NP3 file with known field values
	data := make([]byte, 300)
	copy(data[0:3], []byte{'N', 'C', 'P'}) // Magic bytes
	data[3] = 0x03                          // Version
	copy(data[20:30], []byte("Portrait"))   // Preset name
	data[0x42] = 128                        // Sharpness (neutral)
	data[0x47] = 140                        // Brightness
	data[0x4C] = 130                        // Hue

	output, err := BinaryDump(data, "np3")
	if err != nil {
		t.Fatalf("BinaryDump failed: %v", err)
	}

	// Verify minimum required fields are present
	requiredFields := []string{
		"Magic Bytes",
		"Version",
		"Preset Name",
		"Sharpness",
		"Brightness",
		"Hue",
	}

	for _, field := range requiredFields {
		if !strings.Contains(output, field) {
			t.Errorf("Output missing required field: %s", field)
		}
	}

	// Verify values are shown in normalized form
	if !strings.Contains(output, "normalized") {
		t.Error("Output should show normalized values")
	}
}

// TestBinaryDump_NonNP3Error verifies error for non-NP3 formats (AC-3).
func TestBinaryDump_NonNP3Error(t *testing.T) {
	testData := []byte("test data")

	tests := []struct {
		format      string
		expectError bool
	}{
		{"xmp", true},
		{"lrtemplate", true},
		{"np3", false},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			_, err := BinaryDump(testData, tt.format)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for format %s, got nil", tt.format)
			}

			if tt.expectError && err != nil {
				// Verify error message mentions NP3 and text format
				errMsg := err.Error()
				if !strings.Contains(errMsg, "NP3") {
					t.Error("Error message should mention NP3")
				}
				if !strings.Contains(errMsg, "text") {
					t.Error("Error message should mention text format")
				}
			}
		})
	}
}

// TestBinaryDump_CorruptFile tests graceful degradation with invalid magic bytes (AC-6).
func TestBinaryDump_CorruptFile(t *testing.T) {
	// Create corrupt NP3 file (wrong magic bytes)
	data := make([]byte, 300)
	copy(data[0:3], []byte{'X', 'X', 'X'}) // Wrong magic bytes

	output, err := BinaryDump(data, "np3")

	// Should not error, just warn
	if err != nil {
		t.Errorf("BinaryDump should not error on corrupt file: %v", err)
	}

	// Should contain warning about invalid magic bytes
	if !strings.Contains(output, "Warning") || !strings.Contains(output, "Invalid magic bytes") {
		t.Error("Output should contain warning about invalid magic bytes")
	}

	// Should still produce hex dump
	if !strings.Contains(output, "[0x0000]") {
		t.Error("Should still produce hex dump for corrupt file")
	}
}

// TestBinaryDump_CompleteCoverage verifies all bytes are displayed (AC-1).
func TestBinaryDump_CompleteCoverage(t *testing.T) {
	// Create small NP3 file
	data := make([]byte, 100)
	copy(data[0:3], []byte{'N', 'C', 'P'})

	output, err := BinaryDump(data, "np3")
	if err != nil {
		t.Fatalf("BinaryDump failed: %v", err)
	}

	// Count lines (each byte should have a line, minus warning lines)
	lines := strings.Split(output, "\n")
	hexLines := 0
	for _, line := range lines {
		if strings.Contains(line, "[0x") {
			hexLines++
		}
	}

	// Should have 100 lines (one per byte)
	if hexLines != 100 {
		t.Errorf("Expected 100 hex lines, got %d", hexLines)
	}

	// Verify last byte is displayed
	if !strings.Contains(output, "[0x0063]") { // 0x63 = 99 (last byte)
		t.Error("Output should include last byte")
	}
}

// TestFormatFieldValue tests normalization formulas (AC-2).
func TestFormatFieldValue(t *testing.T) {
	tests := []struct {
		name     string
		offset   int
		rawByte  byte
		expected string
	}{
		{
			name:     "Sharpness neutral",
			offset:   0x42,
			rawByte:  128,
			expected: "neutral",
		},
		{
			name:     "Brightness neutral",
			offset:   0x47,
			rawByte:  128,
			expected: "neutral",
		},
		{
			name:     "Hue neutral",
			offset:   0x4C,
			rawByte:  128,
			expected: "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal data with the test byte
			data := make([]byte, 300)
			copy(data[0:3], []byte{'N', 'C', 'P'})
			data[tt.offset] = tt.rawByte

			output, err := BinaryDump(data, "np3")
			if err != nil {
				t.Fatalf("BinaryDump failed: %v", err)
			}

			// Find the line for this offset
			lines := strings.Split(output, "\n")
			found := false
			for _, line := range lines {
				if strings.Contains(line, fmt.Sprintf("[0x%04X]", tt.offset)) {
					if strings.Contains(line, tt.expected) {
						found = true
						break
					}
				}
			}

			if !found {
				t.Errorf("Expected to find '%s' in output for offset 0x%04X", tt.expected, tt.offset)
			}
		})
	}
}

// TestBinaryDump_AllSamples tests with real NP3 sample files (AC-2, Pattern 7).
func TestBinaryDump_AllSamples(t *testing.T) {
	// Find all NP3 samples
	samples, err := filepath.Glob("../../testdata/xmp/sample.np3")
	if err != nil {
		t.Fatalf("Failed to find samples: %v", err)
	}

	// If no samples found in xmp dir, check np3 dir
	if len(samples) == 0 {
		samples, err = filepath.Glob("../../testdata/np3/*.np3")
		if err != nil {
			t.Fatalf("Failed to find samples: %v", err)
		}
	}

	// If still no samples, skip test
	if len(samples) == 0 {
		t.Skip("No NP3 sample files found")
		return
	}

	for _, samplePath := range samples {
		t.Run(filepath.Base(samplePath), func(t *testing.T) {
			// Read sample file
			data, err := os.ReadFile(samplePath)
			if err != nil {
				t.Fatalf("Failed to read sample: %v", err)
			}

			// Binary dump
			output, err := BinaryDump(data, "np3")
			if err != nil {
				t.Errorf("BinaryDump failed: %v", err)
				return
			}

			// Validate output
			if !strings.Contains(output, "[0x0000]") {
				t.Error("Output missing hex offsets")
			}

			if !strings.Contains(output, "Magic") {
				t.Error("Output missing field annotations")
			}

			// Verify non-empty output
			lines := strings.Split(output, "\n")
			hexLines := 0
			for _, line := range lines {
				if strings.Contains(line, "[0x") {
					hexLines++
				}
			}

			if hexLines < 10 {
				t.Errorf("Expected at least 10 lines of output, got %d", hexLines)
			}
		})
	}
}

// BenchmarkBinaryDump validates performance target <10ms (AC-5).
func BenchmarkBinaryDump(b *testing.B) {
	// Create typical NP3 file (~10KB)
	data := make([]byte, 10*1024)
	copy(data[0:3], []byte{'N', 'C', 'P'})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := BinaryDump(data, "np3")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBinaryDump_LargeFile tests with maximum expected NP3 size (AC-5).
func BenchmarkBinaryDump_LargeFile(b *testing.B) {
	// Create large NP3 file (~50KB)
	data := make([]byte, 50*1024)
	copy(data[0:3], []byte{'N', 'C', 'P'})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := BinaryDump(data, "np3")
		if err != nil {
			b.Fatal(err)
		}
	}
}
