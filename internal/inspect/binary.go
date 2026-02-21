// Package inspect provides binary visualization for NP3 files.
package inspect

import (
	"fmt"
	"strings"
)

// np3FieldMap maps byte offsets to human-readable field names and descriptions.
// These offsets are based on reverse engineering documented in Epic 1.
// Source: internal/formats/np3/parse.go, internal/formats/np3/generate.go
var np3FieldMap = map[int]fieldInfo{
	// File structure (offsets 0-19)
	0x0000: {name: "Magic Bytes", formatter: formatMagicBytes},        // 'N'
	0x0001: {name: "Magic Bytes", formatter: formatMagicBytes},        // 'C'
	0x0002: {name: "Magic Bytes", formatter: formatMagicBytes},        // 'P'
	0x0003: {name: "Version (byte 1)", formatter: formatVersion},      // 4-byte version
	0x0004: {name: "Version (byte 2)", formatter: nil},                //
	0x0005: {name: "Version (byte 3)", formatter: nil},                //
	0x0006: {name: "Version (byte 4)", formatter: nil},                //
	0x0007: {name: "Reserved", formatter: nil},                        // Offsets 7-19 reserved
	0x0008: {name: "Reserved", formatter: nil},
	0x0009: {name: "Reserved", formatter: nil},
	0x000A: {name: "Reserved", formatter: nil},
	0x000B: {name: "Reserved", formatter: nil},
	0x000C: {name: "Reserved", formatter: nil},
	0x000D: {name: "Reserved", formatter: nil},
	0x000E: {name: "Reserved", formatter: nil},
	0x000F: {name: "Reserved", formatter: nil},
	0x0010: {name: "Reserved", formatter: nil},
	0x0011: {name: "Reserved", formatter: nil},
	0x0012: {name: "Reserved", formatter: nil},
	0x0013: {name: "Reserved", formatter: nil},

	// Preset name (offsets 20-59, 40 bytes max, null-terminated ASCII)
	0x0014: {name: "Preset Name (start)", formatter: formatPresetName}, // Offset 20

	// Reserved (offsets 60-63)
	0x003C: {name: "Reserved (post-name)", formatter: nil},
	0x003D: {name: "Reserved (post-name)", formatter: nil},
	0x003E: {name: "Reserved (post-name)", formatter: nil},
	0x003F: {name: "Reserved (post-name)", formatter: nil},

	// Raw parameter bytes (offsets 64-79) - Epic 1 discovered offsets
	0x0040: {name: "Raw Parameter (64)", formatter: formatRawParam},
	0x0041: {name: "Raw Parameter (65)", formatter: formatRawParam},
	// Sharpness: offsets 66-70 (5 bytes, same value)
	0x0042: {name: "Sharpness (66, byte 1/5)", formatter: formatSharpness},
	0x0043: {name: "Sharpness (67, byte 2/5)", formatter: formatSharpness},
	0x0044: {name: "Sharpness (68, byte 3/5)", formatter: formatSharpness},
	0x0045: {name: "Sharpness (69, byte 4/5)", formatter: formatSharpness},
	0x0046: {name: "Sharpness (70, byte 5/5)", formatter: formatSharpness},
	// Brightness: offsets 71-75 (5 bytes, same value)
	0x0047: {name: "Brightness (71, byte 1/5)", formatter: formatBrightness},
	0x0048: {name: "Brightness (72, byte 2/5)", formatter: formatBrightness},
	0x0049: {name: "Brightness (73, byte 3/5)", formatter: formatBrightness},
	0x004A: {name: "Brightness (74, byte 4/5)", formatter: formatBrightness},
	0x004B: {name: "Brightness (75, byte 5/5)", formatter: formatBrightness},
	// Hue: offsets 76-79 (4 bytes, same value)
	0x004C: {name: "Hue (76, byte 1/4)", formatter: formatHue},
	0x004D: {name: "Hue (77, byte 2/4)", formatter: formatHue},
	0x004E: {name: "Hue (78, byte 3/4)", formatter: formatHue},
	0x004F: {name: "Hue (79, byte 4/4)", formatter: formatHue},

	// Reserved (offsets 80-99)
	0x0050: {name: "Reserved (80-99)", formatter: nil},

	// Color data section (offsets 100-299, RGB triplets for saturation)
	// Parser: Extracts RGB triplets where at least one channel > 10
	0x0064: {name: "Color Data Section (start)", formatter: formatColorData}, // Offset 100

	// Tone curve section (offsets 150-499, paired values for contrast)
	// Parser: Extracts pairs where at least one byte is non-zero
	0x0096: {name: "Tone Curve Section (start)", formatter: formatToneCurve}, // Offset 150
}

// fieldInfo contains metadata for a known NP3 field.
type fieldInfo struct {
	name      string
	formatter func(offset int, b byte, data []byte) string
}

// BinaryDump generates an annotated hex dump of NP3 binary data.
// Returns a formatted string with byte offsets, hex values, and field annotations.
//
// Format: [0xOFFSET] HEX  FIELD_NAME (VALUE)
//
// Parameters:
//   - data: Raw bytes of the NP3 file
//   - format: Must be "np3" (binary mode only works with NP3 files)
//
// Returns:
//   - string: Formatted hex dump with annotations
//   - error: If format is not "np3" or if file structure is invalid
//
// Example output:
//
//	[0x0000] 4E  Magic Bytes ('N')
//	[0x0001] 43  Magic Bytes ('C')
//	[0x0002] 50  Magic Bytes ('P')
//	[0x0042] 80  Sharpness (raw: 128, normalized: 0)
func BinaryDump(data []byte, format string) (string, error) {
	// Validate format is NP3 (AC-3)
	if format != "np3" {
		return "", fmt.Errorf("--binary flag only works with NP3 files\n%s files are %s-based text files. View them with any text editor.\nUse 'recipe inspect %s' for JSON parameter output.",
			strings.ToUpper(format),
			getFormatType(format),
			getFormatType(format))
	}

	// Pre-allocate buffer for performance (AC-5)
	// Estimate: ~50-60 chars per line average, pre-allocate for better performance
	var builder strings.Builder
	builder.Grow(len(data) * 55) // Slightly more conservative estimate

	// Check for invalid magic bytes (graceful degradation, AC-6)
	if len(data) >= 3 {
		if data[0] != 'N' || data[1] != 'C' || data[2] != 'P' {
			builder.WriteString(fmt.Sprintf("⚠ Warning: Invalid magic bytes (expected 'NCP', got '%c%c%c')\n",
				data[0], data[1], data[2]))
			builder.WriteString("Proceeding with hex dump...\n\n")
		}
	}

	// Iterate through each byte and generate hex dump
	for offset := 0; offset < len(data); offset++ {
		b := data[offset]

		// Check if this offset has a known field annotation
		if field, exists := np3FieldMap[offset]; exists {
			// Format with field annotation
			fieldValue := ""
			if field.formatter != nil {
				fieldValue = field.formatter(offset, b, data)
			}

			if fieldValue != "" {
				builder.WriteString(fmt.Sprintf("[0x%04X] %02X  %s (%s)\n", offset, b, field.name, fieldValue))
			} else {
				builder.WriteString(fmt.Sprintf("[0x%04X] %02X  %s\n", offset, b, field.name))
			}
		} else {
			// Unknown byte - just show hex without annotation
			builder.WriteString(fmt.Sprintf("[0x%04X] %02X\n", offset, b))
		}
	}

	return builder.String(), nil
}

// getFormatType returns the human-readable type for error messages.
func getFormatType(format string) string {
	switch format {
	case "xmp":
		return "XML"
	case "lrtemplate":
		return "Lua"
	default:
		return "text"
	}
}

// Field formatters convert raw bytes to human-readable values.

func formatMagicBytes(offset int, b byte, data []byte) string {
	return fmt.Sprintf("'%c'", b)
}

func formatVersion(offset int, b byte, data []byte) string {
	if offset == 0x0003 && len(data) > 6 {
		// Read version as 4 bytes (little-endian)
		version := uint32(data[3]) | uint32(data[4])<<8 | uint32(data[5])<<16 | uint32(data[6])<<24
		return fmt.Sprintf("version %d", version)
	}
	return ""
}

func formatPresetName(offset int, b byte, data []byte) string {
	if offset == 0x0014 && len(data) >= 60 {
		// Extract preset name (null-terminated ASCII)
		nameBytes := data[20:60]
		nameEnd := 0
		for i, byte := range nameBytes {
			if byte == 0 {
				nameEnd = i
				break
			}
		}
		if nameEnd == 0 {
			nameEnd = len(nameBytes)
		}
		name := string(nameBytes[:nameEnd])
		return fmt.Sprintf("\"%s\"", name)
	}
	return ""
}

func formatRawParam(offset int, b byte, data []byte) string {
	// Convert to signed value (128 = neutral)
	signed := int(b) - 128
	return fmt.Sprintf("raw: %d, adjusted: %d", b, signed)
}

func formatSharpness(offset int, b byte, data []byte) string {
	// Sharpness stored at offsets 0x42-0x46 (5 bytes, same value)
	// Range: 0-9 (needs denormalization from raw 0-255)
	// Formula: sharpness = (raw / 255) * 9
	if b == 128 {
		return "raw: 128 (neutral/0)"
	}
	sharpness := float64(b) / 255.0 * 9.0
	return fmt.Sprintf("raw: %d, normalized: %.1f", b, sharpness)
}

func formatBrightness(offset int, b byte, data []byte) string {
	// Brightness stored at offsets 0x47-0x4B (5 bytes, same value)
	// Range: -1.0 to +1.0 (128 = neutral/0)
	// Formula: brightness = (raw - 128) / 128.0
	brightness := (float64(b) - 128.0) / 128.0
	if b == 128 {
		return "raw: 128 (neutral/0)"
	}
	return fmt.Sprintf("raw: %d, normalized: %.2f", b, brightness)
}

func formatHue(offset int, b byte, data []byte) string {
	// Hue stored at offsets 0x4C-0x4F (4 bytes, same value)
	// Range: -9° to +9° (128 = neutral/0)
	// Formula: hue = (raw - 128) * 9 / 128
	hue := (int(b) - 128) * 9 / 128
	if b == 128 {
		return "raw: 128 (neutral/0°)"
	}
	return fmt.Sprintf("raw: %d, normalized: %d°", b, hue)
}

func formatColorData(offset int, b byte, data []byte) string {
	// Color data starts at offset 0x64 (100 decimal)
	// RGB triplets for saturation analysis
	if offset == 0x0064 {
		return "RGB triplets for saturation"
	}
	return ""
}

func formatToneCurve(offset int, b byte, data []byte) string {
	// Tone curve starts at offset 0x96 (150 decimal)
	// Paired values for contrast analysis
	if offset == 0x0096 {
		return "paired values for contrast"
	}
	return ""
}
