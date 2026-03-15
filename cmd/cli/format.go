// Package main provides file format detection utilities for the Recipe CLI.
//
// Format detection happens in two stages:
//  1. Extension-based: Fast detection from file extension (.np3, .xmp)
//  2. Content-based: Fallback to examining file content (magic bytes, XML structure)
//
// Extension-based detection is preferred for performance (sub-microsecond).
// Content-based detection is used when extension is ambiguous or missing.
//
// All format detection defers to internal/converter for content analysis,
// maintaining the thin CLI layer architecture pattern.
package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/justin/recipe/internal/converter"
)

// Supported format identifiers (exported for use in other CLI modules)
const (
	FormatNP3 = "np3"
	FormatXMP = "xmp"
)

// detectFormat returns format string based on file extension (fast path).
//
// Performs case-insensitive extension matching for supported formats:
//   - .np3 → "np3"
//   - .xmp → "xmp"
//
// Returns error if extension is unrecognized or missing.
func detectFormat(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".np3":
		return FormatNP3, nil
	case ".xmp":
		return FormatXMP, nil
	case "":
		return "", fmt.Errorf("file has no extension (expected .np3 or .xmp)")
	default:
		return "", fmt.Errorf("unknown file format: %s (expected .np3 or .xmp)", ext)
	}
}

// detectFormatFromBytes detects format from file content (fallback path).
//
// Wraps converter.DetectFormat() to maintain thin CLI layer architecture.
func detectFormatFromBytes(data []byte) (string, error) {
	format, err := converter.DetectFormat(data)
	if err != nil {
		return "", fmt.Errorf("unable to detect format from file content (size: %d bytes)", len(data))
	}
	return format, nil
}

// validateFormat checks if a format string is supported.
//
// Accepts: "np3", "xmp"
// Rejects: Any other string with descriptive error
func validateFormat(format string) error {
	switch format {
	case FormatNP3, FormatXMP:
		return nil
	default:
		return fmt.Errorf("unsupported format: %q (must be %q or %q)",
			format, FormatNP3, FormatXMP)
	}
}
