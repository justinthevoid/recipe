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
