// Package dcp provides parsing and generation of DNG Camera Profile (.dcp) files.
// DCPs are TIFF-based DNG files containing binary camera profile data in
// TIFF tags 50700-52600. Recipe supports tone curve adjustments (exposure,
// contrast, highlights, shadows) extracted from the binary profile tone curve.
package dcp

import (
	"fmt"

	"github.com/justin/recipe/internal/models"
)

// Parse parses a DNG Camera Profile (.dcp) file and returns a UniversalRecipe.
//
// DCP files are Adobe DNG (Digital Negative) containers with camera profiles
// stored as binary TIFF tags. This function extracts:
//   - Tag 52552: Profile name
//   - Tag 50940: Profile tone curve (binary float32 array)
//   - Tag 50721-50722: Color matrices (optional, stored in metadata)
//   - Tag 50730: Baseline exposure offset
//
// Tone curves are analyzed to extract exposure, contrast, highlights, and
// shadows adjustments compatible with Lightroom/NX Studio.
//
// Supported DCP versions: DNG 1.0-1.6
// Unsupported features: Full camera calibration (HSV tables, LUTs, dual illuminant)
//
// Example:
//
//	data, _ := os.ReadFile("nikon-standard.dcp")
//	recipe, err := dcp.Parse(data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Exposure: %.2f\n", recipe.Exposure)
//
// Returns error if:
//   - File is not a valid TIFF/DNG (invalid magic bytes)
//   - Tag 52552 (ProfileName) is missing
//   - Binary tag data is malformed
//   - TIFF/DNG structure is corrupt
func Parse(data []byte) (*models.UniversalRecipe, error) {
	// Step 1: Read and validate TIFF/DNG structure
	tiffFile, err := readTIFF(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read DCP file: %w", err)
	}

	// Get first IFD (DCP files have profile data in first IFD)
	ifds := tiffFile.IFDs()
	if len(ifds) == 0 {
		return nil, fmt.Errorf("DCP file has no IFDs")
	}
	ifd := ifds[0]

	// Step 2: Extract camera profile data from binary tags
	profile := &CameraProfile{}

	// Extract profile name (optional - not all DCP files have it)
	profile.ProfileName, err = extractProfileName(ifd)
	if err != nil {
		return nil, fmt.Errorf("failed to extract profile name: %w", err)
	}
	// If no profile name, we'll just use an empty string (caller can use filename)

	// Extract tone curve (optional)
	profile.ToneCurve, err = extractToneCurve(ifd)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tone curve: %w", err)
	}

	// Extract color matrices (optional)
	profile.ColorMatrix1, err = extractColorMatrix(ifd, TagColorMatrix1)
	if err != nil {
		return nil, fmt.Errorf("failed to extract color matrix 1: %w", err)
	}

	profile.ColorMatrix2, err = extractColorMatrix(ifd, TagColorMatrix2)
	if err != nil {
		return nil, fmt.Errorf("failed to extract color matrix 2: %w", err)
	}

	// Extract baseline exposure (optional)
	profile.BaselineExposure, err = extractBaselineExposure(ifd)
	if err != nil {
		return nil, fmt.Errorf("failed to extract baseline exposure: %w", err)
	}

	// Step 3: Convert to UniversalRecipe
	recipe := profileToUniversal(profile)

	return recipe, nil
}
