// Package np3 provides functionality for generating Nikon Picture Control (.np3) binary files.
//
// This file implements the generator which converts UniversalRecipe to NP3 binary format.
package np3

import (
	"fmt"

	"github.com/justin/recipe/internal/models"
)

// Generate encodes a UniversalRecipe into a Nikon Picture Control (.np3) binary file.
//
// The function converts UniversalRecipe parameters to NP3 ranges, generates the file
// structure with magic bytes and header, encodes parameters as TLV chunks, and returns
// the complete binary representation.
//
// Parameters:
//   - recipe: UniversalRecipe to convert to NP3 format
//
// Returns:
//   - []byte: NP3 binary file data
//   - error: Validation or generation error with descriptive context
//
// Errors:
//   - Nil recipe: Input validation failed
//   - Parameter out of range: Value cannot be represented in NP3 format
//   - Encoding error: Binary encoding failed
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
	// Validate input
	if recipe == nil {
		return nil, fmt.Errorf("generate NP3: recipe cannot be nil")
	}

	// Convert parameters from UniversalRecipe to NP3 ranges
	params, err := convertToNP3Parameters(recipe)
	if err != nil {
		return nil, fmt.Errorf("generate NP3: %w", err)
	}

	// Validate converted parameters
	if err := validateParameters(params); err != nil {
		return nil, fmt.Errorf("generate NP3: %w", err)
	}

	// Build binary structure
	data, err := encodeBinary(params, recipe.Name)
	if err != nil {
		return nil, fmt.Errorf("generate NP3: %w", err)
	}

	return data, nil
}

// convertToNP3Parameters converts UniversalRecipe parameters to NP3 ranges.
//
// Phase 2 Enhancement: Now converts all 48 parameters using exact offset mappings
// discovered through TypeScript implementation research.
//
// Conversion formulas (reverse of parser mappings in parse.go:845-952):
//   - Basic Adjustments (2): sharpening, clarity
//   - Advanced Adjustments (7): mid-range sharpening, contrast, highlights, shadows, whites, blacks, saturation
//   - Color Blender (24): 8 colors × 3 values (hue, chroma/saturation, brightness/luminance)
//   - Color Grading (11): 3 zones + 2 global params
//   - Tone Curve (3): control points
func convertToNP3Parameters(recipe *models.UniversalRecipe) (*np3Parameters, error) {
	params := &np3Parameters{}

	// === Basic Adjustments (2 parameters) ===

	// Convert Sharpness: UniversalRecipe (0-150) → NP3 (-3.0 to +9.0)
	// Parser uses: (np3Value + 3.0) * 12.5 = universalValue
	// Generator reverses: (universalValue / 12.5) - 3.0 = np3Value
	params.sharpening = (float64(recipe.Sharpness) / 12.5) - 3.0
	if params.sharpening > 9.0 {
		params.sharpening = 9.0
	} else if params.sharpening < -3.0 {
		params.sharpening = -3.0
	}

	// Convert Clarity: UniversalRecipe (-100 to +100) → NP3 (-5.0 to +5.0)
	// Parser uses: np3Value * 20 = universalValue
	// Generator reverses: universalValue / 20 = np3Value
	params.clarity = float64(recipe.Clarity) / 20.0
	if params.clarity > 5.0 {
		params.clarity = 5.0
	} else if params.clarity < -5.0 {
		params.clarity = -5.0
	}

	// === Advanced Adjustments (7 parameters) ===

	// Mid-Range Sharpening: Direct mapping (-5.0 to +5.0)
	params.midRangeSharpening = recipe.MidRangeSharpening
	if params.midRangeSharpening > 5.0 {
		params.midRangeSharpening = 5.0
	} else if params.midRangeSharpening < -5.0 {
		params.midRangeSharpening = -5.0
	}

	// Contrast, Highlights, Shadows, Whites, Blacks, Saturation: Direct mapping (-100 to +100)
	params.contrast = recipe.Contrast
	if params.contrast > 100 {
		params.contrast = 100
	} else if params.contrast < -100 {
		params.contrast = -100
	}

	params.highlights = recipe.Highlights
	if params.highlights > 100 {
		params.highlights = 100
	} else if params.highlights < -100 {
		params.highlights = -100
	}

	params.shadows = recipe.Shadows
	if params.shadows > 100 {
		params.shadows = 100
	} else if params.shadows < -100 {
		params.shadows = -100
	}

	params.whiteLevel = recipe.Whites
	if params.whiteLevel > 100 {
		params.whiteLevel = 100
	} else if params.whiteLevel < -100 {
		params.whiteLevel = -100
	}

	params.blackLevel = recipe.Blacks
	if params.blackLevel > 100 {
		params.blackLevel = 100
	} else if params.blackLevel < -100 {
		params.blackLevel = -100
	}

	params.saturation = recipe.Saturation
	if params.saturation > 100 {
		params.saturation = 100
	} else if params.saturation < -100 {
		params.saturation = -100
	}

	// === Color Blender (24 parameters) ===
	// UniversalRecipe ColorAdjustment maps directly to NP3 Color Blender (hue, chroma, brightness)

	params.redHue = recipe.Red.Hue
	params.redChroma = recipe.Red.Saturation
	params.redBrightness = recipe.Red.Luminance

	params.orangeHue = recipe.Orange.Hue
	params.orangeChroma = recipe.Orange.Saturation
	params.orangeBrightness = recipe.Orange.Luminance

	params.yellowHue = recipe.Yellow.Hue
	params.yellowChroma = recipe.Yellow.Saturation
	params.yellowBrightness = recipe.Yellow.Luminance

	params.greenHue = recipe.Green.Hue
	params.greenChroma = recipe.Green.Saturation
	params.greenBrightness = recipe.Green.Luminance

	params.cyanHue = recipe.Aqua.Hue
	params.cyanChroma = recipe.Aqua.Saturation
	params.cyanBrightness = recipe.Aqua.Luminance

	params.blueHue = recipe.Blue.Hue
	params.blueChroma = recipe.Blue.Saturation
	params.blueBrightness = recipe.Blue.Luminance

	params.purpleHue = recipe.Purple.Hue
	params.purpleChroma = recipe.Purple.Saturation
	params.purpleBrightness = recipe.Purple.Luminance

	params.magentaHue = recipe.Magenta.Hue
	params.magentaChroma = recipe.Magenta.Saturation
	params.magentaBrightness = recipe.Magenta.Luminance

	// === Color Grading (11 parameters) ===
	// Direct mapping from UniversalRecipe ColorGrading

	if recipe.ColorGrading != nil {
		params.highlightsZone = recipe.ColorGrading.Highlights
		params.midtoneZone = recipe.ColorGrading.Midtone
		params.shadowsZone = recipe.ColorGrading.Shadows
		params.blending = recipe.ColorGrading.Blending
		params.balance = recipe.ColorGrading.Balance
	}

	// === Tone Curve (3 parameters) ===
	// Convert UniversalRecipe PointCurve to NP3 format

	if len(recipe.PointCurve) > 0 {
		params.toneCurvePointCount = len(recipe.PointCurve)
		if params.toneCurvePointCount > 127 {
			params.toneCurvePointCount = 127 // NP3 limit
		}

		params.toneCurvePoints = make([]toneCurvePoint, params.toneCurvePointCount)
		for i := 0; i < params.toneCurvePointCount; i++ {
			params.toneCurvePoints[i] = toneCurvePoint{
				position: 405 + (i * 2), // Offset calculation
				value1:   uint8(recipe.PointCurve[i].Input),
				value2:   uint8(recipe.PointCurve[i].Output),
			}
		}
	}

	// === Legacy Parameters (for heuristic fallback compatibility) ===

	// Brightness: UniversalRecipe Exposure field → NP3 brightness (-1.0 to +1.0)
	params.brightness = recipe.Exposure
	if params.brightness > 1.0 {
		params.brightness = 1.0
	} else if params.brightness < -1.0 {
		params.brightness = -1.0
	}

	// Hue: No direct mapping in UniversalRecipe (only per-color hue adjustments)
	params.hue = 0

	// Retrieve raw binary data if available for perfect round-trip preservation
	if recipe.FormatSpecificBinary != nil {
		if rawData, ok := recipe.FormatSpecificBinary["np3_raw"]; ok {
			params.rawData = rawData
		}
	}

	return params, nil
}

// npChunk represents a TLV (Type-Length-Value) chunk in the NP3 format.
// Each chunk is 10 bytes: [ID:1][Pad:3][Length:2BE][Value:2][Pad:2]
type npChunk struct {
	id     byte    // Chunk ID (0x03-0x1F)
	length uint16  // Length in big-endian (typically 2, but 0x1F has length 28)
	value  []byte  // Value bytes (2 bytes typically, more for extended chunks)
}

// constantChunks defines all 18 constant TLV chunks found in real NP3 files.
// These chunks have the same values across all presets and represent format structure.
var constantChunks = []npChunk{
	{id: 0x03, length: 2, value: []byte{0x00, 0x20}}, // Format identifier (value=32)
	{id: 0x04, length: 2, value: []byte{0x00, 0x00}}, // Reserved (value=0)
	{id: 0x05, length: 2, value: []byte{0xff, 0x01}}, // Format flag (value=65281)
	{id: 0x08, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x09, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x0a, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x0b, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x0c, length: 2, value: []byte{0xff, 0x00}}, // Default value (value=65280)
	{id: 0x0d, length: 2, value: []byte{0xff, 0x00}}, // Default value (value=65280)
	{id: 0x0e, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x0f, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x10, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x11, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x12, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x13, length: 2, value: []byte{0xff, 0x01}}, // Default value (value=65281)
	{id: 0x15, length: 2, value: []byte{0xff, 0x0a}}, // Default value (value=65290)
	{id: 0x17, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
	{id: 0x18, length: 2, value: []byte{0xff, 0x04}}, // Default value (value=65284)
}

// writeChunk writes a single TLV chunk at the specified offset.
// Chunk format: [ChunkID:1][Padding:3][Length:2BE][Value:N][Padding:2]
// Total size: 10 bytes (for chunks with length=2)
func writeChunk(data []byte, offset int, chunk npChunk) int {
	data[offset] = chunk.id
	data[offset+1] = 0x00 // Padding
	data[offset+2] = 0x00
	data[offset+3] = 0x00
	data[offset+4] = byte(chunk.length >> 8)   // Length big-endian high byte
	data[offset+5] = byte(chunk.length & 0xff) // Length big-endian low byte

	// Write value bytes (typically 2 bytes)
	copy(data[offset+6:offset+6+len(chunk.value)], chunk.value)

	// Write padding after value
	data[offset+8] = 0x00
	data[offset+9] = 0x00

	return offset + 10 // Return next chunk offset
}

// writeChunkMetadataOnly writes chunk structure without overwriting value bytes.
// Used when heuristic data has already been written to preserve those values.
// Chunk format: [ChunkID:1][Padding:3][Length:2BE][Value:N][Padding:2]
func writeChunkMetadataOnly(data []byte, offset int, chunkID byte, length uint16) int {
	data[offset] = chunkID
	data[offset+1] = 0x00 // Padding
	data[offset+2] = 0x00
	data[offset+3] = 0x00
	data[offset+4] = byte(length >> 8)   // Length big-endian high byte
	data[offset+5] = byte(length & 0xff) // Length big-endian low byte

	// Skip offset+6 and offset+7 (value bytes) - preserve existing heuristic data

	// Write padding after value
	data[offset+8] = 0x00
	data[offset+9] = 0x00

	return offset + 10 // Return next chunk offset
}

// writeChunks writes all TLV chunks (constant + variable) starting at offset 46.
// This generates the 29 chunks (0x03-0x1F) required for Nikon NX Studio validation.
// Variable chunks use metadata-only writing to preserve heuristic data values.
func writeChunks(data []byte, params *np3Parameters) {
	offset := 46

	// Write all constant chunks (full chunk including values)
	for _, chunk := range constantChunks {
		offset = writeChunk(data, offset, chunk)
	}

	// Write variable chunks using METADATA-ONLY to preserve heuristic data
	// Phase 1: Heuristic data has already been written at offsets 64-79, 100-299, 150-500
	// We write chunk structure without overwriting those value bytes

	// Chunk 0x06: Possibly saturation-related
	offset = writeChunkMetadataOnly(data, offset, 0x06, 2)

	// Chunk 0x07: Possibly contrast-related
	offset = writeChunkMetadataOnly(data, offset, 0x07, 2)

	// Chunk 0x14: Possibly brightness-related
	offset = writeChunkMetadataOnly(data, offset, 0x14, 2)

	// Chunk 0x16: Unknown parameter
	offset = writeChunkMetadataOnly(data, offset, 0x16, 2)

	// Chunks 0x19-0x1E: Possibly color channels
	for chunkID := byte(0x19); chunkID <= 0x1E; chunkID++ {
		offset = writeChunkMetadataOnly(data, offset, chunkID, 2)
	}

	// Chunk 0x1F: Extended data chunk with length=28
	// For extended chunks, we still need to write some structure
	// Write the chunk header
	data[offset] = 0x1F
	data[offset+1] = 0x00
	data[offset+2] = 0x00
	data[offset+3] = 0x00
	data[offset+4] = 0x00 // length=28 big-endian high
	data[offset+5] = 0x1C // length=28 big-endian low
	// Skip value bytes (offset+6 to offset+33) - preserve heuristic data
	// No trailing padding needed as next chunk would be at offset+34

	// Total chunks written: 18 constant + 11 variable = 29 chunks
	// Total bytes: 29 * 10 = 290 bytes (offsets 46-335)
}

// encodeBinary generates the complete NP3 binary structure.
//
// File Structure (Nikon NX Studio format):
//   - Offset 0x00-0x02: Magic bytes "NCP"
//   - Offset 0x03-0x06: Version (4 bytes)
//   - Offset 0x07-0x13: Reserved (zeros)
//   - Offset 0x14-0x17: Name header (0x00 0x00 0x00 0x14)
//   - Offset 0x18-0x2B: Preset name (20 bytes, ASCII null-terminated)
//   - Offset 0x2C-0x2D: Padding (zeros)
//   - Offset 0x2E-0x14F: TLV chunks (29 chunks × 10 bytes = 290 bytes, offsets 46-335)
//   - Offset 0x40-0x4F: Raw parameter bytes (offsets 64-80, for heuristic parser)
//   - Offset 0x64-0x12B: Color data section (bytes 100-300, for heuristic parser)
//   - Offset 0x96-0x1F3: Tone curve data section (bytes 150-500, for heuristic parser)
//   - Remaining: Padding to minimum 480 bytes (matching sample.np3)
func encodeBinary(params *np3Parameters, presetName string) ([]byte, error) {
	// If we have raw data from parsing, use it as the base to preserve chunks
	// Otherwise create a new buffer
	var data []byte
	if params.rawData != nil && len(params.rawData) > 0 {
		// Copy raw data to preserve header, chunks, and all structure
		data = make([]byte, len(params.rawData))
		copy(data, params.rawData)
	} else {
		// Create buffer for complete file
		// Real files: 392-978+ bytes, using 480 as minimum (matching sample.np3)
		data = make([]byte, 480)
	}

	// Only write header if we don't have raw data (raw data already has correct header)
	if params.rawData == nil || len(params.rawData) == 0 {
		// Write magic bytes "NCP" at offset 0-2
		copy(data[0:3], magicBytes)

		// Write version bytes at offset 3-6 (NP3 v1.0.0.0, little-endian: 0x00000001)
		copy(data[3:7], []byte{0x00, 0x00, 0x00, 0x01})

		// Offsets 7-10: Reserved (already zeros)

		// Write format header at offsets 11-19
		data[11] = 0x04                                      // Format flag
		copy(data[12:16], []byte{'0', '3', '1', '0'})       // Format version string "0310"
		// Offsets 16-17: Reserved (already zeros)
		data[18] = 0x02                                      // Unknown flag
		// Offset 19: Reserved (already zero)

		// Write name length header at offsets 20-23
		copy(data[20:24], []byte{0x00, 0x00, 0x00, 0x14})    // Name length = 20 (0x14)
	}

	// Only write legacy data structures if we don't have raw data
	// (raw data already has correct values and writing would corrupt chunks)
	if params.rawData == nil || len(params.rawData) == 0 {
		// Phase 2: Legacy heuristic data generation (generateColorData, generateToneCurveData)
		// has been removed in favor of exact offset writing. The old heuristic approach
		// wrote data to offsets 100-299 (color) and 150-499 (tone curve), but Phase 2
		// uses exact offsets for all 48 parameters instead.
		//
		// The following legacy functions are still needed for basic file structure:

		// Write raw parameter bytes at offsets 64-79 (legacy structure)
		writeRawParameterBytes(data, params)

		// Write preset name at offset 24-43 (20 bytes, matching real NP3 files)
		if presetName != "" {
			nameBytes := []byte(presetName)
			maxNameLen := 20 // Real files use exactly 20 bytes for name
			if len(nameBytes) > maxNameLen {
				nameBytes = nameBytes[:maxNameLen]
			}
			copy(data[24:44], nameBytes)
		}
	}

	// Offsets 44-45: Padding (already zeros)

	// Only write TLV chunks if we don't have raw data (raw data already has correct chunks)
	if params.rawData == nil || len(params.rawData) == 0 {
		// Write TLV chunks at offsets 46-335 (29 chunks × 10 bytes = 290 bytes)
		// This is required for Nikon NX Studio validation
		// We write chunks LAST, after heuristic data, so the chunks contain the actual
		// parameter values in their value bytes (hybrid format)
		writeChunks(data, params)
	}

	// === Phase 2: Write all parameters to exact offsets ===
	// Write exact offset data regardless of whether we have rawData, to ensure
	// current parameter values are written to the correct locations
	writeBasicAdjustments(data, params)
	writeAdvancedAdjustments(data, params)
	writeColorBlender(data, params)
	writeColorGrading(data, params)
	writeToneCurve(data, params)

	// Return the complete 480-byte buffer (matching real .np3 files)
	return data, nil
}

// writeRawParameterBytes writes raw parameter values to the byte offsets
// that the parser reads from (offsets 64-79).
//
// Parser extraction strategy:
//   - Sharpness: offsets 66-70 (adjusted from -128 to +127, mapped to 0-9)
//   - Brightness: offsets 71-75 (adjusted from -128 to +127, normalized to -1.0 to +1.0)
//   - Hue: offsets 76-79 (adjusted from -128 to +127, mapped to -9 to +9)
//
// Generator strategy (reverse of parser):
//   - Convert NP3 parameter values back to raw bytes using 128-offset encoding
func writeRawParameterBytes(data []byte, params *np3Parameters) {
	// Sharpness: Write to offsets 66-70
	// Parser has default value of 5 when no non-zero bytes found
	// For sharpening=0, we need a non-zero byte that produces 0 after formula
	// Formula: params.sharpening = (avgSharpness + 128) * 9 / 255
	// For result=0: (adjusted + 128) * 9 / 255 = 0 → adjusted in range -128 to -100
	// Using byte value 1: adjusted = 1 - 128 = -127 → sharpening = (1 * 9) / 255 = 0
	var sharpnessRaw byte
	if params.sharpening == 0 {
		// Use byte value 1 to avoid default value of 5
		sharpnessRaw = 1
	} else {
		sharpnessAdjusted := (params.sharpening * 255 / 9) - 128
		sharpnessRaw = byte(sharpnessAdjusted + 128)
	}

	// Write same value to all 5 bytes (66-70) for consistency
	for i := 66; i <= 70; i++ {
		data[i] = sharpnessRaw
	}

	// Brightness: Write to offsets 71-75
	// Parser: params.brightness = float64(avgBrightness) / 128.0
	// Generator reverse: adjusted = brightness * 128.0
	brightnessAdjusted := int(params.brightness * 128.0)
	brightnessRaw := byte(brightnessAdjusted + 128)

	// Write same value to all 5 bytes (71-75) for consistency
	for i := 71; i <= 75; i++ {
		data[i] = brightnessRaw
	}

	// Hue: Write to offsets 76-79
	// Parser: params.hue = avgHue * 9 / 128
	// Generator reverse: adjusted = hue * 128 / 9
	hueAdjusted := params.hue * 128 / 9
	hueRaw := byte(hueAdjusted + 128)

	// Write same value to all 4 bytes (76-79) for consistency
	for i := 76; i <= 79; i++ {
		data[i] = hueRaw
	}
}

// generateColorData generates RGB color triplets at offsets 100-299
// that will be heuristically analyzed by the parser to determine saturation.
//
// Parser strategy (parse.go:200-219):
//   colorIntensity := len(colorData) / 15
//   params.saturation = colorIntensity - 1
//   (clamped to -3 to +3)
//
// Generator reverse strategy:
//   To achieve saturation value S, we need (S + 1) * 15 significant RGB triplets
//   We generate triplets with at least one channel > 10
//
// Returns the offset where color data ends
func generateColorData(data []byte, saturation int) int {
	// Calculate how many significant RGB triplets we need
	// Parser: saturation = (len(colorData) / 15) - 1
	// Generator: len(colorData) = (saturation + 1) * 15
	targetCount := (saturation + 1) * 15

	// Clamp to reasonable range (can't have negative count)
	if targetCount < 0 {
		targetCount = 0
	}
	// Limit to available space (200 bytes / 3 = 66 triplets max)
	if targetCount > 66 {
		targetCount = 66
	}

	// Generate RGB triplets starting at offset 100
	offset := 100
	for i := 0; i < targetCount && offset+2 < 300; i++ {
		// Generate triplet with at least one channel > 10
		// Use simple pattern: (50, 50, 50) for positive saturation
		data[offset] = 50
		data[offset+1] = 50
		data[offset+2] = 50
		offset += 3
	}

	return offset  // Return where color data ends
}

// generateToneCurveData generates paired byte values at offsets 150-499
// that will be heuristically analyzed by the parser to determine contrast.
//
// Parser strategy (parse.go:221-237):
//   Reads pairs from offset 150-500
//   curveComplexity := len(toneCurve) / 20
//   params.contrast = curveComplexity - 2
//   (clamped to -3 to +3)
//
// Generator reverse strategy:
//   To achieve contrast value C, we need (C + 2) * 20 non-zero tone curve points
//   If color data extends into tone curve region (150+), those bytes will be
//   counted as tone curve pairs. Adjust our count to compensate.
func generateToneCurveData(data []byte, contrast int, colorEndOffset int) {
	// Calculate how many tone curve pairs we need total
	// Parser: contrast = (len(toneCurve) / 20) - 2
	// Generator: len(toneCurve) = (contrast + 2) * 20
	targetTotalPairs := (contrast + 2) * 20

	// Calculate how many pairs already exist from color data overlap (150 to colorEndOffset)
	// Parser reads tone curve from offset 150, so any color data >= 150 will be counted
	overlapPairs := 0
	if colorEndOffset > 150 {
		// Color data extends into tone curve region
		// Last color byte is at colorEndOffset-1
		// Parser reads pairs starting at 150: (150,151), (152,153), ..., (colorEndOffset-2, colorEndOffset-1)
		// Number of pairs = ((lastByte - 150) / 2) + 1
		lastColorByte := colorEndOffset - 1
		if lastColorByte >= 150 {
			overlapPairs = ((lastColorByte - 150) / 2) + 1
		}
	}

	// Calculate how many additional tone curve pairs we need to generate
	// Total pairs counted = overlap pairs + our generated pairs
	// We want: overlap + generated = target
	additionalPairs := targetTotalPairs - overlapPairs
	if additionalPairs < 0 {
		additionalPairs = 0
	}

	// Determine start offset: after color data, with even alignment
	startOffset := colorEndOffset
	if startOffset < 150 {
		startOffset = 150
	}
	if startOffset%2 != 0 {
		startOffset++
	}

	// Generate additional tone curve pairs starting after color data
	offset := startOffset
	for i := 0; i < additionalPairs && offset+1 < 500; i++ {
		// Use value 1 (minimal non-zero) for tone curve data
		data[offset] = 1
		data[offset+1] = 1
		offset += 2
	}
}

// === Phase 2: Exact Offset Writers ===
// These functions write all 48 parameters to exact byte offsets discovered through
// TypeScript implementation research. They mirror the extraction functions in parse.go.

// writeBasicAdjustments writes sharpening and clarity to exact offsets using Scaled4 encoding.
//
// Offsets:
//   - Sharpening: offset 82, Scaled4 encoding (-3.0 to +9.0)
//   - Clarity: offset 92, Scaled4 encoding (-5.0 to +5.0)
//
// Encoding: Scaled4 = (value * 4.0) + 0x80
func writeBasicAdjustments(data []byte, params *np3Parameters) {
	// Sharpening (offset 82)
	if len(data) > OffsetSharpening {
		data[OffsetSharpening] = EncodeScaled4(params.sharpening)
	}

	// Clarity (offset 92)
	if len(data) > OffsetClarity {
		data[OffsetClarity] = EncodeScaled4(params.clarity)
	}
}

// writeAdvancedAdjustments writes 7 advanced parameters to exact offsets using Scaled4 and Signed8 encodings.
//
// Offsets:
//   - Mid-Range Sharpening: offset 242, Scaled4 (-5.0 to +5.0)
//   - Contrast: offset 272, Signed8 (-100 to +100)
//   - Highlights: offset 282, Signed8 (-100 to +100)
//   - Shadows: offset 292, Signed8 (-100 to +100)
//   - White Level: offset 302, Signed8 (-100 to +100)
//   - Black Level: offset 312, Signed8 (-100 to +100)
//   - Saturation: offset 322, Signed8 (-100 to +100)
func writeAdvancedAdjustments(data []byte, params *np3Parameters) {
	// Mid-Range Sharpening (offset 242, Scaled4)
	if len(data) > OffsetMidRangeSharpening {
		data[OffsetMidRangeSharpening] = EncodeScaled4(params.midRangeSharpening)
	}

	// Contrast (offset 272, Signed8)
	if len(data) > OffsetContrast {
		data[OffsetContrast] = EncodeSigned8(params.contrast)
	}

	// Highlights (offset 282, Signed8)
	if len(data) > OffsetHighlights {
		data[OffsetHighlights] = EncodeSigned8(params.highlights)
	}

	// Shadows (offset 292, Signed8)
	if len(data) > OffsetShadows {
		data[OffsetShadows] = EncodeSigned8(params.shadows)
	}

	// White Level (offset 302, Signed8)
	if len(data) > OffsetWhiteLevel {
		data[OffsetWhiteLevel] = EncodeSigned8(params.whiteLevel)
	}

	// Black Level (offset 312, Signed8)
	if len(data) > OffsetBlackLevel {
		data[OffsetBlackLevel] = EncodeSigned8(params.blackLevel)
	}

	// Saturation (offset 322, Signed8)
	if len(data) > OffsetSaturation {
		data[OffsetSaturation] = EncodeSigned8(params.saturation)
	}
}

// writeColorBlender writes 24 color blender parameters (8 colors × 3 values) to exact offsets.
//
// Offsets: 332-355 (sequential, 3 bytes per color)
// Encoding: Signed8 for all values (-100 to +100)
func writeColorBlender(data []byte, params *np3Parameters) {
	// Red (offsets 332-334)
	if len(data) > OffsetRedBrightness {
		data[OffsetRedHue] = EncodeSigned8(params.redHue)
		data[OffsetRedChroma] = EncodeSigned8(params.redChroma)
		data[OffsetRedBrightness] = EncodeSigned8(params.redBrightness)
	}

	// Orange (offsets 335-337)
	if len(data) > OffsetOrangeBrightness {
		data[OffsetOrangeHue] = EncodeSigned8(params.orangeHue)
		data[OffsetOrangeChroma] = EncodeSigned8(params.orangeChroma)
		data[OffsetOrangeBrightness] = EncodeSigned8(params.orangeBrightness)
	}

	// Yellow (offsets 338-340)
	if len(data) > OffsetYellowBrightness {
		data[OffsetYellowHue] = EncodeSigned8(params.yellowHue)
		data[OffsetYellowChroma] = EncodeSigned8(params.yellowChroma)
		data[OffsetYellowBrightness] = EncodeSigned8(params.yellowBrightness)
	}

	// Green (offsets 341-343)
	if len(data) > OffsetGreenBrightness {
		data[OffsetGreenHue] = EncodeSigned8(params.greenHue)
		data[OffsetGreenChroma] = EncodeSigned8(params.greenChroma)
		data[OffsetGreenBrightness] = EncodeSigned8(params.greenBrightness)
	}

	// Cyan (offsets 344-346)
	if len(data) > OffsetCyanBrightness {
		data[OffsetCyanHue] = EncodeSigned8(params.cyanHue)
		data[OffsetCyanChroma] = EncodeSigned8(params.cyanChroma)
		data[OffsetCyanBrightness] = EncodeSigned8(params.cyanBrightness)
	}

	// Blue (offsets 347-349)
	if len(data) > OffsetBlueBrightness {
		data[OffsetBlueHue] = EncodeSigned8(params.blueHue)
		data[OffsetBlueChroma] = EncodeSigned8(params.blueChroma)
		data[OffsetBlueBrightness] = EncodeSigned8(params.blueBrightness)
	}

	// Purple (offsets 350-352)
	if len(data) > OffsetPurpleBrightness {
		data[OffsetPurpleHue] = EncodeSigned8(params.purpleHue)
		data[OffsetPurpleChroma] = EncodeSigned8(params.purpleChroma)
		data[OffsetPurpleBrightness] = EncodeSigned8(params.purpleBrightness)
	}

	// Magenta (offsets 353-355)
	if len(data) > OffsetMagentaBrightness {
		data[OffsetMagentaHue] = EncodeSigned8(params.magentaHue)
		data[OffsetMagentaChroma] = EncodeSigned8(params.magentaChroma)
		data[OffsetMagentaBrightness] = EncodeSigned8(params.magentaBrightness)
	}
}

// writeColorGrading writes 11 color grading parameters to exact offsets using Hue12 and Signed8 encodings.
//
// Offsets:
//   - Highlights Zone: 368-371 (2-byte hue + 1-byte chroma + 1-byte brightness)
//   - Midtone Zone: 372-375 (2-byte hue + 1-byte chroma + 1-byte brightness)
//   - Shadows Zone: 376-379 (2-byte hue + 1-byte chroma + 1-byte brightness)
//   - Blending: 384 (direct value 0-100)
//   - Balance: 386 (Signed8 -100 to +100)
func writeColorGrading(data []byte, params *np3Parameters) {
	// Highlights zone (offsets 368-371)
	if len(data) > OffsetHighlightsBrightness {
		byte1, byte2 := EncodeHue12(params.highlightsZone.Hue)
		data[OffsetHighlightsHue] = byte1
		data[OffsetHighlightsHue+1] = byte2
		data[OffsetHighlightsChroma] = EncodeSigned8(params.highlightsZone.Chroma)
		data[OffsetHighlightsBrightness] = EncodeSigned8(params.highlightsZone.Brightness)
	}

	// Midtone zone (offsets 372-375)
	if len(data) > OffsetMidtoneBrightness {
		byte1, byte2 := EncodeHue12(params.midtoneZone.Hue)
		data[OffsetMidtoneHue] = byte1
		data[OffsetMidtoneHue+1] = byte2
		data[OffsetMidtoneChroma] = EncodeSigned8(params.midtoneZone.Chroma)
		data[OffsetMidtoneBrightness] = EncodeSigned8(params.midtoneZone.Brightness)
	}

	// Shadows zone (offsets 376-379)
	if len(data) > OffsetShadowsBrightness {
		byte1, byte2 := EncodeHue12(params.shadowsZone.Hue)
		data[OffsetShadowsHue] = byte1
		data[OffsetShadowsHue+1] = byte2
		data[OffsetShadowsChroma] = EncodeSigned8(params.shadowsZone.Chroma)
		data[OffsetShadowsBrightness] = EncodeSigned8(params.shadowsZone.Brightness)
	}

	// Blending (offset 384, direct value 0-100)
	if len(data) > OffsetColorGradingBlending {
		data[OffsetColorGradingBlending] = uint8(params.blending)
	}

	// Balance (offset 386, Signed8 -100 to +100)
	if len(data) > OffsetColorGradingBalance {
		data[OffsetColorGradingBalance] = EncodeSigned8(params.balance)
	}
}

// writeToneCurve writes tone curve control points to exact offsets.
//
// Offsets:
//   - Point Count: offset 404 (direct value 0-255)
//   - Control Points: offset 405 (2 bytes per point: input, output)
//   - Raw Curve: offset 460 (257 × 16-bit big-endian values)
func writeToneCurve(data []byte, params *np3Parameters) {
	// Point count (offset 404)
	if len(data) > OffsetToneCurvePointCount {
		data[OffsetToneCurvePointCount] = uint8(params.toneCurvePointCount)

		// Control points (offset 405, 2 bytes per point)
		if params.toneCurvePointCount > 0 && len(params.toneCurvePoints) > 0 {
			for i := 0; i < params.toneCurvePointCount && i < len(params.toneCurvePoints); i++ {
				offset := OffsetToneCurvePoints + (i * 2)
				if len(data) > offset+1 {
					data[offset] = params.toneCurvePoints[i].value1
					data[offset+1] = params.toneCurvePoints[i].value2
				}
			}
		}
	}

	// Raw curve (offset 460, 257 × 16-bit big-endian)
	// Note: This is optional and typically not written unless we have raw curve data
	if len(params.toneCurveRaw) == 257 && len(data) >= OffsetToneCurveRaw+514 {
		for i := 0; i < 257; i++ {
			offset := OffsetToneCurveRaw + (i * 2)
			data[offset] = uint8(params.toneCurveRaw[i] >> 8)   // High byte
			data[offset+1] = uint8(params.toneCurveRaw[i] & 0xFF) // Low byte
		}
	}
}
