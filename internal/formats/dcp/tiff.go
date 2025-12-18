package dcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/google/tiff"
)

// readTIFF reads and validates a TIFF/DNG file structure.
//
// Validates TIFF/DNG magic bytes:
//   - II (0x4949) = little-endian TIFF
//   - MM (0x4D4D) = big-endian TIFF
//   - IIRC (0x49494352) = little-endian DNG (used by DCP files)
//   - MMCR (0x4D4D4352) = big-endian DNG (used by DCP files)
//
// DCP files are actually DNG (Digital Negative) files, which use
// "CR" (0x4352) or "RC" (0x5243) as the version identifier instead
// of standard TIFF version 42 (0x002A).
//
// Returns error if:
//   - File is too small (<8 bytes)
//   - Magic bytes are invalid
//   - TIFF/DNG structure is corrupt
func readTIFF(data []byte) (tiff.TIFF, error) {
	// Validate file size
	if len(data) < 8 {
		return nil, fmt.Errorf("file too small to be a TIFF/DNG (got %d bytes, need at least 8)", len(data))
	}

	// Validate TIFF/DNG magic bytes (first 2 bytes)
	magicII := []byte{0x49, 0x49} // "II" little-endian
	magicMM := []byte{0x4D, 0x4D} // "MM" big-endian

	if !bytes.Equal(data[:2], magicII) && !bytes.Equal(data[:2], magicMM) {
		return nil, fmt.Errorf("invalid TIFF/DNG magic bytes: got %#x %#x, expected II (0x49 0x49) or MM (0x4D 0x4D)", data[0], data[1])
	}

	// For DNG files (DCP), replace "CR" version with standard TIFF version 42
	// This allows google/tiff to parse DNG files as standard TIFF
	modifiedData := make([]byte, len(data))
	copy(modifiedData, data)

	// Check if this is a DNG file (version bytes are "RC" or "CR")
	isDNG := false
	if bytes.Equal(data[:2], magicII) && bytes.Equal(data[2:4], []byte{0x52, 0x43}) { // "IIRC"
		isDNG = true
		modifiedData[2] = 0x2A // Version 42 (little-endian)
		modifiedData[3] = 0x00
	} else if bytes.Equal(data[:2], magicMM) && bytes.Equal(data[2:4], []byte{0x43, 0x52}) { // "MMCR"
		isDNG = true
		modifiedData[2] = 0x00 // Version 42 (big-endian)
		modifiedData[3] = 0x2A
	}

	// Create reader wrapper
	reader := tiff.NewReadAtReadSeeker(bytes.NewReader(modifiedData))

	// Parse TIFF structure using google/tiff library
	// Use default tag space and field type space
	tiffFile, err := tiff.Parse(reader, tiff.DefaultTagSpace, tiff.DefaultFieldTypeSpace)
	if err != nil {
		// Provide more context if it's a DNG file
		if isDNG {
			return nil, fmt.Errorf("failed to parse DNG structure: %w", err)
		}
		return nil, fmt.Errorf("failed to parse TIFF structure: %w", err)
	}

	return tiffFile, nil
}

// extractProfileName extracts the profile name from TIFF tag 52552.
// Returns empty string if tag is not present (tag is optional in some DCP files).
func extractProfileName(ifd tiff.IFD) (string, error) {
	if !ifd.HasField(TagProfileName) {
		// ProfileName is optional - some DCP files don't have it
		return "", nil
	}

	field := ifd.GetField(TagProfileName)
	if field == nil {
		return "", nil
	}

	data := field.Value().Bytes()
	if len(data) == 0 {
		return "", nil
	}

	// ASCII strings may be null-terminated
	name := string(data)
	if idx := bytes.IndexByte(data, 0); idx >= 0 {
		name = string(data[:idx])
	}

	return name, nil
}

// extractToneCurve extracts the tone curve from TIFF tag 50940.
//
// Tone curves are stored as arrays of 32-bit floats where each
// consecutive pair represents (input, output) normalized to 0.0-1.0.
func extractToneCurve(ifd tiff.IFD) ([]ToneCurvePoint, error) {
	if !ifd.HasField(TagProfileToneCurve) {
		// Tone curve is optional - return empty curve (linear)
		return nil, nil
	}

	field := ifd.GetField(TagProfileToneCurve)
	if field == nil {
		return nil, fmt.Errorf("tone curve field is nil")
	}

	data := field.Value().Bytes()
	if len(data) == 0 {
		return nil, nil
	}

	// Each point is 2 floats (input, output) = 8 bytes
	if len(data)%8 != 0 {
		return nil, fmt.Errorf("invalid tone curve data length: %d (must be multiple of 8)", len(data))
	}

	numPoints := len(data) / 8
	points := make([]ToneCurvePoint, numPoints)

	for i := 0; i < numPoints; i++ {
		offset := i * 8

		// Read as little-endian floats
		inputBits := binary.LittleEndian.Uint32(data[offset : offset+4])
		outputBits := binary.LittleEndian.Uint32(data[offset+4 : offset+8])

		points[i] = ToneCurvePoint{
			Input:  float64(bitsToFloat32(inputBits)),
			Output: float64(bitsToFloat32(outputBits)),
		}
	}

	return points, nil
}

// extractColorMatrix extracts a 3x3 color matrix from a TIFF tag.
//
// Color matrices are stored as 9 SRational values (numerator/denominator pairs)
// in row-major order.
func extractColorMatrix(ifd tiff.IFD, tag int) (*Matrix, error) {
	if !ifd.HasField(uint16(tag)) {
		// Color matrices are optional
		return nil, nil
	}

	field := ifd.GetField(uint16(tag))
	if field == nil {
		return nil, fmt.Errorf("color matrix tag %d field is nil", tag)
	}

	data := field.Value().Bytes()
	if len(data) != 72 { // 9 SRationals * 8 bytes each
		return nil, fmt.Errorf("invalid color matrix data length: %d (expected 72)", len(data))
	}

	matrix := &Matrix{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			offset := (i*3 + j) * 8

			// SRational = signed int32 numerator + signed int32 denominator
			num := int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
			denom := int32(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))

			if denom == 0 {
				return nil, fmt.Errorf("color matrix has zero denominator at [%d][%d]", i, j)
			}

			matrix.Rows[i][j] = float64(num) / float64(denom)
		}
	}

	return matrix, nil
}

// extractBaselineExposure extracts baseline exposure offset from tag 50730.
func extractBaselineExposure(ifd tiff.IFD) (float64, error) {
	if !ifd.HasField(TagBaselineExposureOffset) {
		// Baseline exposure is optional - default to 0.0
		return 0.0, nil
	}

	field := ifd.GetField(TagBaselineExposureOffset)
	if field == nil {
		return 0.0, nil
	}

	data := field.Value().Bytes()
	if len(data) != 8 { // 1 SRational = 8 bytes
		return 0.0, fmt.Errorf("invalid baseline exposure data length: %d (expected 8)", len(data))
	}

	num := int32(binary.LittleEndian.Uint32(data[0:4]))
	denom := int32(binary.LittleEndian.Uint32(data[4:8]))

	if denom == 0 {
		return 0.0, nil
	}

	return float64(num) / float64(denom), nil
}

// bitsToFloat32 converts uint32 bits to float32.
func bitsToFloat32(bits uint32) float32 {
	return *(*float32)(unsafe.Pointer(&bits))
}

// writeDNG creates a DNG file with binary profile tags.
//
// DNG files are TIFF files with special version bytes:
//   - Standard TIFF: II (little-endian) + 0x002A (version 42)
//   - DNG format: II (little-endian) + 0x5243 ("CR" = DNG version)
//
// Structure:
//   Bytes 0-1: "II" (0x49 0x49) - little-endian byte order
//   Bytes 2-3: "CR" (0x52 0x43) - DNG version identifier
//   Bytes 4-7: IFD offset (uint32) - offset to first Image File Directory
//   Bytes 8+:  IFD with profile tags
//
// This function uses github.com/google/tiff to write TIFF structure,
// then patches the version bytes to DNG format ("CR" instead of standard TIFF 42).
//
// Returns DNG file as []byte ready to write to disk.
func writeDNG(toneCurve []byte, colorMatrix1 []byte, colorMatrix2 []byte, forwardMatrix1 []byte, forwardMatrix2 []byte, profileName string, cameraModel string, baselineExposure []byte, lookTable []byte) ([]byte, error) {
	// Create a minimal TIFF file using google/tiff
	// We'll use the library to generate the structure, then patch it for DNG

	// Create TIFF buffer
	buf := new(bytes.Buffer)

	// Write TIFF header manually (we'll convert to DNG later)
	// II (little-endian)
	buf.Write([]byte{0x49, 0x49})
	// Version 42 (we'll patch this to "CR" later for DNG)
	buf.Write([]byte{0x2A, 0x00})
	// IFD offset (IFD starts right after 8-byte header)
	binary.Write(buf, binary.LittleEndian, uint32(8))

	// Build IFD (Image File Directory) with profile tags
	ifd := buildProfileIFD(toneCurve, colorMatrix1, colorMatrix2, forwardMatrix1, forwardMatrix2, profileName, cameraModel, baselineExposure, lookTable)

	// Write IFD to buffer
	if err := writeIFD(buf, ifd); err != nil {
		return nil, fmt.Errorf("failed to write IFD: %w", err)
	}

	// Patch TIFF version to DNG version ("CR" instead of 0x002A)
	// This makes it a DNG file instead of standard TIFF
	data := buf.Bytes()
	data[2] = 0x52 // 'R'
	data[3] = 0x43 // 'C'  (together: "RC" in little-endian = 0x4352)

	return data, nil
}

// ProfileIFDEntry represents a single TIFF tag in the IFD.
type ProfileIFDEntry struct {
	Tag      uint16
	Type     uint16
	Count    uint32
	ValueOrOffset []byte // 4 bytes for value, or offset to data if >4 bytes
	Data     []byte // Actual data if Count*TypeSize > 4
}

// buildProfileIFD creates an IFD with DCP profile tags.
//
// Required tags:
//   - 256 (ImageWidth): 1 pixel (minimal image for profile-only DCP)
//   - 257 (ImageLength): 1 pixel
//   - 258 (BitsPerSample): [8, 8, 8] for RGB
//   - 259 (Compression): 1 (uncompressed)
//   - 262 (PhotometricInterpretation): 2 (RGB)
//   - 277 (SamplesPerPixel): 3 (RGB)
//
// Profile tags:
//   - 50708 (UniqueCameraModel): "Recipe Converted Camera" (ASCII)
//   - 50721 (ColorMatrix1): 72 bytes (9 SRationals)
//   - 50722 (ColorMatrix2): 72 bytes (9 SRationals)
//   - 50730 (BaselineExposureOffset): 8 bytes (1 SRational)
//   - 50940 (ProfileToneCurve): N*8 bytes (N float32 pairs)
//   - 52552 (ProfileName): ASCII string (OPTIONAL - only if non-empty)
func buildProfileIFD(toneCurve, colorMatrix1, colorMatrix2, forwardMatrix1, forwardMatrix2 []byte, profileName, cameraModel string, baselineExposure []byte, lookTable []byte) []ProfileIFDEntry {
	// Use provided camera model or default
	if cameraModel == "" {
		cameraModel = "Recipe Converted Camera"
	}
	cameraModelBytes := append([]byte(cameraModel), 0) // Null-terminated

	// Build entries in strict ascending tag order
	entries := []ProfileIFDEntry{
		// DCP profile tags (MUST be in ascending order by tag number - TIFF requirement!)
		// NOTE: Unlike standard TIFF files, DCPs should NOT have ImageWidth, ImageHeight, etc.
		{Tag: 50708, Type: 2, Count: uint32(len(cameraModelBytes)), Data: cameraModelBytes}, // UniqueCameraModel (0xc614)
		{Tag: 50721, Type: 10, Count: 9, Data: colorMatrix1},                      // ColorMatrix1 (0xc621)
		{Tag: 50722, Type: 10, Count: 9, Data: colorMatrix2},                      // ColorMatrix2 (0xc622)
		{Tag: 50778, Type: 3, Count: 1, ValueOrOffset: uint16ToBytes(17)},           // CalibrationIlluminant1 (0xc65a)
		{Tag: 50779, Type: 3, Count: 1, ValueOrOffset: uint16ToBytes(21)},           // CalibrationIlluminant2 (0xc65b)
		{Tag: 50932, Type: 2, Count: uint32(len("com.adobe") + 1), Data: []byte("com.adobe\x00")}, // ProfileCalibrationSignature (0xc6f4)
	}

	// Add ProfileName (0xc6f8) in correct position if non-empty
	if profileName != "" {
		nameBytes := append([]byte(profileName), 0) // Null-terminated ASCII
		entries = append(entries, ProfileIFDEntry{
			Tag:   50936, // ProfileName (0xc6f8) - MUST come before 50940!
			Type:  2,     // ASCII
			Count: uint32(len(nameBytes)),
			Data:  nameBytes,
		})
	}

	// Continue with remaining tags in order
	entries = append(entries, []ProfileIFDEntry{
		{Tag: 50940, Type: 11, Count: uint32(len(toneCurve) / 4), Data: toneCurve},  // ProfileToneCurve (0xc6fc)
		{Tag: 50941, Type: 4, Count: 1, ValueOrOffset: uint32ToBytes(0)},            // ProfileEmbedPolicy (0xc6fd)
		{Tag: 50942, Type: 2, Count: uint32(len("Recipe DCP Generator") + 1), Data: []byte("Recipe DCP Generator\x00")}, // ProfileCopyright (0xc6fe)
		{Tag: 50964, Type: 10, Count: 9, Data: forwardMatrix1},                      // ForwardMatrix1 (0xc714)
		{Tag: 50965, Type: 10, Count: 9, Data: forwardMatrix2},                      // ForwardMatrix2 (0xc715)
		{Tag: 50981, Type: 4, Count: 3, Data: []byte{90, 0, 0, 0, 16, 0, 0, 0, 16, 0, 0, 0}}, // ProfileLookTableDims (0xc7b5): 90×16×16
		{Tag: 50982, Type: 11, Count: uint32(len(lookTable) / 4), Data: lookTable},  // ProfileLookTableData (0xc7b6)
		{Tag: 51108, Type: 4, Count: 1, ValueOrOffset: uint32ToBytes(1)},            // ProfileLookTableEncoding (0xc7a4): 1 = sRGB
		{Tag: 51109, Type: 10, Count: 1, Data: baselineExposure},                    // BaselineExposureOffset (0xc7a5)
		{Tag: 51110, Type: 4, Count: 1, ValueOrOffset: uint32ToBytes(1)},            // DefaultBlackRender (0xc7a6): 1 = None
	}...)

	return entries
}

// writeIFD writes an IFD (Image File Directory) to a buffer.
//
// IFD structure:
//   - 2 bytes: Number of entries (uint16)
//   - N × 12 bytes: Tag entries
//   - 4 bytes: Next IFD offset (0 = no more IFDs)
//   - Variable: Tag data (for tags with Count*TypeSize > 4)
//
// Each tag entry (12 bytes):
//   - Bytes 0-1: Tag ID (uint16)
//   - Bytes 2-3: Data type (uint16)
//   - Bytes 4-7: Count (uint32)
//   - Bytes 8-11: Value or offset to data (uint32)
func writeIFD(buf *bytes.Buffer, entries []ProfileIFDEntry) error {
	// Calculate IFD start offset (current buffer position)
	ifdStart := buf.Len()

	// Write number of entries
	binary.Write(buf, binary.LittleEndian, uint16(len(entries)))

	// Calculate offset for tag data (after all entries + next IFD offset)
	dataOffset := ifdStart + 2 + (len(entries) * 12) + 4

	// Write tag entries
	for _, entry := range entries {
		// Tag ID
		binary.Write(buf, binary.LittleEndian, entry.Tag)
		// Type
		binary.Write(buf, binary.LittleEndian, entry.Type)
		// Count
		binary.Write(buf, binary.LittleEndian, entry.Count)

		// Determine if value fits in 4 bytes or needs offset
		dataSize := len(entry.Data)
		if dataSize == 0 {
			// Value fits in ValueOrOffset (≤4 bytes)
			buf.Write(entry.ValueOrOffset)
			// Pad to 4 bytes if needed
			for len(entry.ValueOrOffset) < 4 {
				buf.WriteByte(0)
			}
		} else if dataSize <= 4 {
			// Small data - write inline
			buf.Write(entry.Data)
			// Pad to 4 bytes
			for i := dataSize; i < 4; i++ {
				buf.WriteByte(0)
			}
		} else {
			// Large data - write offset
			binary.Write(buf, binary.LittleEndian, uint32(dataOffset))
			dataOffset += dataSize
			// Account for padding to even byte boundary
			if dataSize%2 == 1 {
				dataOffset++
			}
		}
	}

	// Write next IFD offset (0 = last IFD)
	binary.Write(buf, binary.LittleEndian, uint32(0))

	// Write tag data (for tags with Data > 4 bytes)
	// IMPORTANT: All data must be aligned to even byte boundaries (TIFF requirement)
	for _, entry := range entries {
		if len(entry.Data) > 4 {
			buf.Write(entry.Data)
			// Pad to even byte boundary if odd
			if len(entry.Data)%2 == 1 {
				buf.WriteByte(0)
			}
		}
	}

	return nil
}

// uint32ToBytes converts uint32 to 4-byte little-endian array.
func uint32ToBytes(val uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, val)
	return buf
}

// uint16ToBytes converts uint16 to 2-byte little-endian array (padded to 4 bytes).
func uint16ToBytes(val uint16) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint16(buf, val)
	return buf
}
