// Package main provides file format detection utilities for the Recipe CLI.
//
// Format detection happens in two stages:
//  1. Extension-based: Fast detection from file extension (.np3, .xmp, .lrtemplate)
//  2. Content-based: Fallback to examining file content (magic bytes, XML structure, Lua syntax)
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
	FormatNP3        = "np3"
	FormatXMP        = "xmp"
	FormatLRTemplate = "lrtemplate"
	// FormatCostyle     = "costyle"     // DISABLED: costyle format support
	// FormatCostylepack = "costylepack" // DISABLED: costyle format support
)

// detectFormat returns format string based on file extension (fast path).
//
// Performs case-insensitive extension matching for all supported formats:
//   - .np3 → "np3"
//   - .xmp → "xmp"
//   - .lrtemplate → "lrtemplate"
//
// Returns error if extension is unrecognized or missing.
//
// Performance: <1μs (instant)
//
// Example:
//
//	format, err := detectFormat("portrait.xmp")
//	if err != nil {
//	    // Handle unknown extension
//	}
func detectFormat(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".np3":
		return FormatNP3, nil
	case ".xmp":
		return FormatXMP, nil
	case ".lrtemplate":
		return FormatLRTemplate, nil
	// DISABLED: costyle format support
	// case ".costyle":
	// 	return FormatCostyle, nil
	// case ".costylepack":
	// 	return FormatCostylepack, nil
	case "":
		return "", fmt.Errorf("file has no extension (expected .np3, .xmp, or .lrtemplate)")
	default:
		return "", fmt.Errorf("unknown file format: %s (expected .np3, .xmp, or .lrtemplate)", ext)
	}
}

// detectFormatFromBytes detects format from file content (fallback path).
//
// Wraps converter.DetectFormat() to maintain thin CLI layer architecture.
// Uses content-based detection via magic bytes, XML structure, and Lua syntax.
//
// Detection logic (from converter package):
//   - NP3: Magic bytes "NCP" (ASCII) at start + minimum 300 bytes
//   - DCP: Magic bytes "IIRC" (Adobe DNG Camera Profile)
//   - lrtemplate: Lua table syntax (s = { at start after trimming whitespace)
//   - XMP: XML with x:xmpmeta wrapper (Adobe Lightroom format)
//
// Returns error if content doesn't match any known format.
//
// Performance: <5ms for typical files (<50KB)
//
// Example:
//
//	data, _ := os.ReadFile("unknown_file")
//	format, err := detectFormatFromBytes(data)
//	if err != nil {
//	    // Handle unknown content
//	}
func detectFormatFromBytes(data []byte) (string, error) {
	// Defer to converter package for content-based detection
	// This maintains the thin CLI layer architecture pattern
	format, err := converter.DetectFormat(data)
	if err != nil {
		// Wrap error with CLI-appropriate message
		return "", fmt.Errorf("unable to detect format from file content (size: %d bytes)", len(data))
	}
	return format, nil
}

// validateFormat checks if a format string is supported.
//
// Accepts: "np3", "xmp", "lrtemplate"
// Rejects: Any other string with descriptive error
//
// Used by convert command to validate user input from --from and --to flags.
//
// Example:
//
//	if err := validateFormat(userFormat); err != nil {
//	    fmt.Println(err) // "Unsupported format: "pdf" (must be 'np3', 'xmp', or 'lrtemplate')"
//	}
func validateFormat(format string) error {
	switch format {
	case FormatNP3, FormatXMP, FormatLRTemplate: // DISABLED: FormatCostyle, FormatCostylepack
		return nil
	default:
		return fmt.Errorf("unsupported format: %q (must be %q, %q, or %q)",
			format, FormatNP3, FormatXMP, FormatLRTemplate)
	}
}
