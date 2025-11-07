package main

import (
	"os"
	"strings"
	"testing"

	"github.com/justin/recipe/internal/converter"
)

// TestDetectFormat tests extension-based format detection (AC-1)
func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		// Valid extensions
		{"np3 lowercase", "portrait.np3", "np3", false},
		{"np3 uppercase", "portrait.NP3", "np3", false},
		{"np3 mixed case", "portrait.Np3", "np3", false},
		{"xmp lowercase", "preset.xmp", "xmp", false},
		{"xmp uppercase", "preset.XMP", "xmp", false},
		{"xmp mixed case", "preset.Xmp", "xmp", false},
		{"lrtemplate lowercase", "classic.lrtemplate", "lrtemplate", false},
		{"lrtemplate uppercase", "classic.LRTEMPLATE", "lrtemplate", false},
		{"lrtemplate mixed case", "classic.LrTemplate", "lrtemplate", false},

		// Paths with directories
		{"/path/to/file.np3", "/path/to/file.np3", "np3", false},
		{"file with dots", "file.with.dots.xmp", "xmp", false},
		{"Windows path", `C:\Users\file.lrtemplate`, "lrtemplate", false},

		// Invalid extensions
		{"unknown extension", "unknown.txt", "", true},
		{"no extension", "no-extension", "", true},
		{"empty filename", "", "", true},
		{"jpg extension", "photo.jpg", "", true},
		{"pdf extension", "doc.pdf", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectFormat(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectFormat(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("detectFormat(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

// TestDetectFormatFromBytes tests content-based format detection (AC-2)
func TestDetectFormatFromBytes(t *testing.T) {
	// NP3: Magic bytes + minimum size
	np3Data := make([]byte, 300)
	copy(np3Data, []byte{'N', 'C', 'P'})

	// XMP: XML with crs namespace
	xmpData := []byte(`<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
    <rdf:RDF xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
    </rdf:RDF>
</x:xmpmeta>`)

	// lrtemplate: Lua table
	lrtemplateData := []byte(`s = {
	id = "12345678-1234-1234-1234-123456789012",
	internalName = "Preset Name",
}`)

	tests := []struct {
		name    string
		data    []byte
		want    string
		wantErr bool
	}{
		{"np3 magic bytes", np3Data, "np3", false},
		{"xmp xml structure", xmpData, "xmp", false},
		{"lrtemplate lua", lrtemplateData, "lrtemplate", false},
		{"unknown format", []byte("random data"), "", true},
		{"empty file", []byte{}, "", true},
		{"too small np3", []byte{'N', 'C', 'P'}, "", true}, // Less than 300 bytes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectFormatFromBytes(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectFormatFromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("detectFormatFromBytes() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestValidateFormat tests format validation (AC-3)
func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		// Valid formats
		{"np3", "np3", false},
		{"xmp", "xmp", false},
		{"lrtemplate", "lrtemplate", false},

		// Invalid formats
		{"uppercase NP3", "NP3", true},
		{"uppercase XMP", "XMP", true},
		{"pdf", "pdf", true},
		{"jpg", "jpg", true},
		{"empty string", "", true},
		{"unknown", "unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFormat(%q) error = %v, wantErr %v", tt.format, err, tt.wantErr)
			}
		})
	}
}

// TestFormatConstants verifies format constants exist and have correct values (AC-4)
func TestFormatConstants(t *testing.T) {
	tests := []struct {
		name  string
		got   string
		want  string
	}{
		{"FormatNP3", FormatNP3, "np3"},
		{"FormatXMP", FormatXMP, "xmp"},
		{"FormatLRTemplate", FormatLRTemplate, "lrtemplate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}

// TestErrorMessages verifies error messages are user-friendly (AC-6)
func TestErrorMessages(t *testing.T) {
	t.Run("extension detection error", func(t *testing.T) {
		_, err := detectFormat("test.txt")
		if err == nil || !strings.Contains(err.Error(), "unknown file format") {
			t.Errorf("Expected clear error message about unknown format, got: %v", err)
		}
		if err != nil && !strings.Contains(err.Error(), "expected") {
			t.Errorf("Expected error to suggest valid formats, got: %v", err)
		}
	})

	t.Run("missing extension error", func(t *testing.T) {
		_, err := detectFormat("noext")
		if err == nil || !strings.Contains(err.Error(), "no extension") {
			t.Errorf("Expected clear error message about missing extension, got: %v", err)
		}
	})

	t.Run("content detection error", func(t *testing.T) {
		_, err := detectFormatFromBytes([]byte("invalid"))
		if err == nil || !strings.Contains(err.Error(), "unable to detect") {
			t.Errorf("Expected clear error message about detection failure, got: %v", err)
		}
		if err != nil && !strings.Contains(err.Error(), "size:") {
			t.Errorf("Expected error to include file size, got: %v", err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		err := validateFormat("invalid")
		if err == nil || !strings.Contains(err.Error(), "unsupported format") {
			t.Errorf("Expected clear error message about unsupported format, got: %v", err)
		}
		if err != nil && !strings.Contains(err.Error(), "must be") {
			t.Errorf("Expected error to suggest valid formats, got: %v", err)
		}
	})
}

// TestDetectFormatUsesConverter verifies CLI detection wraps converter package (AC-5)
func TestDetectFormatUsesConverter(t *testing.T) {
	// Test that CLI format detection defers to converter package
	testFiles := []struct {
		path   string
		format string
	}{
		{"../../testdata/xmp/portrait.xmp", "xmp"},
		{"../../testdata/np3/sample.np3", "np3"},
	}

	for _, tt := range testFiles {
		t.Run(tt.path, func(t *testing.T) {
			// Skip if file doesn't exist
			data, err := os.ReadFile(tt.path)
			if err != nil {
				t.Skipf("Test file not found: %s", tt.path)
				return
			}

			// CLI detection should match converter detection
			cliFormat, err := detectFormatFromBytes(data)
			if err != nil {
				t.Fatalf("CLI detection failed: %v", err)
			}

			converterFormat, err := converter.DetectFormat(data)
			if err != nil {
				t.Fatalf("Converter detection failed: %v", err)
			}

			if cliFormat != converterFormat {
				t.Errorf("CLI format %q != converter format %q", cliFormat, converterFormat)
			}

			if cliFormat != tt.format {
				t.Errorf("Expected format %q, got %q", tt.format, cliFormat)
			}
		})
	}
}

// TestFormatDetectionIntegration tests format detection with real files (AC-5 Integration)
func TestFormatDetectionIntegration(t *testing.T) {
	testFiles := []struct {
		path   string
		format string
	}{
		{"../../testdata/xmp/portrait.xmp", "xmp"},
		{"../../testdata/np3/sample.np3", "np3"},
	}

	for _, tt := range testFiles {
		t.Run(tt.path, func(t *testing.T) {
			// Skip if file doesn't exist
			data, err := os.ReadFile(tt.path)
			if err != nil {
				t.Skipf("Test file not found: %s", tt.path)
				return
			}

			// Test extension-based detection
			format, err := detectFormat(tt.path)
			if err != nil {
				t.Errorf("detectFormat(%q) error: %v", tt.path, err)
			}
			if format != tt.format {
				t.Errorf("detectFormat(%q) = %q, want %q", tt.path, format, tt.format)
			}

			// Test content-based detection matches
			contentFormat, err := detectFormatFromBytes(data)
			if err != nil {
				t.Errorf("detectFormatFromBytes() error: %v", err)
			}
			if contentFormat != tt.format {
				t.Errorf("Content detection = %q, want %q", contentFormat, tt.format)
			}

			// Verify extension and content detection agree
			if format != contentFormat {
				t.Errorf("Extension detection (%q) != content detection (%q)", format, contentFormat)
			}
		})
	}
}
