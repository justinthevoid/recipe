// Package dcp provides functionality for generating Adobe DNG Camera Profile (.dcp) files.
// DCP files are TIFF-based binary files that define camera color profiles for Adobe Lightroom.
//
// This package enables converting UniversalRecipe presets to DCP format, allowing
// presets from NP3, XMP, lrtemplate, and .costyle formats to be used as camera profiles
// in Lightroom, Camera Raw, and with DNG files.
package dcp

import (
	"fmt"

	"github.com/justin/recipe/internal/models"
)

// Generate creates a DNG Camera Profile (.dcp) file from a UniversalRecipe.
//
// DCP files are Adobe DNG (Digital Negative) containers with camera profiles
// stored as binary TIFF tags. This function generates:
//   - Tag 50708: Unique camera model ("Nikon Z f" or from recipe.Metadata["camera_model"])
//   - Tag 50721-50722: Calibrated color matrices (Nikon Z f calibration)
//   - Tag 50778-50779: Calibration illuminants (Standard Light A, D65)
//   - Tag 50932: Profile calibration signature ("com.adobe")
//   - Tag 50936: Profile name (from recipe.Metadata["profile_name"])
//   - Tag 50940: Profile tone curve (binary float32 array)
//   - Tag 50941: Profile embed policy (Allow Copying)
//   - Tag 50942: Profile copyright
//   - Tag 50964-50965: Calibrated forward matrices (Nikon Z f calibration)
//   - Tag 50981-50982: 3D color lookup table (90×16×16 identity HSV→RGB LUT)
//   - Tag 51108: Profile look table encoding (sRGB)
//   - Tag 51109: Baseline exposure offset (-0.15 EV for Nikon Z f)
//   - Tag 51110: Default black render (None)
//
// Tone curves are generated from UniversalRecipe exposure, contrast, highlights,
// and shadows parameters. All curves use 0.0-1.0 normalized values as required
// by the DNG specification.
//
// Color matrices are calibrated for Nikon Z f from Adobe Camera Raw profiles.
// This ensures compatibility with Adobe Lightroom's profile loading system.
//
// Example:
//
//	recipe := &models.UniversalRecipe{
//	    Exposure: 0.5,      // +0.5 EV
//	    Contrast: 30,       // +30
//	    Highlights: -20,    // -20
//	    Shadows: 10,        // +10
//	    Metadata: map[string]interface{}{
//	        "profile_name": "Portrait",
//	    },
//	}
//	dcpData, err := dcp.Generate(recipe)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("portrait.dcp", dcpData, 0644)
//
// Returns error if:
//   - recipe is nil
//   - Binary DNG structure generation fails
//   - TIFF IFD writing fails
//
// Supported DCP versions: DNG 1.0-1.6
// Performance target: <200ms (slower than other formats due to TIFF/DNG overhead)
func Generate(recipe *models.UniversalRecipe) ([]byte, error) {
	// Validate input
	if recipe == nil {
		return nil, fmt.Errorf("recipe cannot be nil")
	}

	// Step 1: Generate tone curve points (5-point piecewise linear curve)
	points := universalToToneCurve(recipe)

	// Step 2: Convert tone curve to binary float32 array (tag 50940)
	toneCurveBinary := generateBinaryToneCurve(points)

	// Step 3: Generate calibrated color matrices (tags 50721-50722)
	// Using Nikon Z f calibration from Adobe Camera Raw
	colorMatrix1Binary := srationalArrayToBytes(generateColorMatrix())

	// Use warm matrix if recipe metadata requests it
	// This enables custom warm DCP generation as documented in FINAL_CONCLUSIONS.md
	useWarmMatrix := false
	if recipe.Metadata != nil {
		if warmFlag, ok := recipe.Metadata["use_warm_matrix"].(bool); ok && warmFlag {
			useWarmMatrix = true
		}
	}

	var colorMatrix2Binary []byte
	if useWarmMatrix {
		colorMatrix2Binary = srationalArrayToBytes(generateColorMatrix2Warm())
	} else {
		colorMatrix2Binary = srationalArrayToBytes(generateColorMatrix2())
	}

	// Step 4: Generate calibrated forward matrices (tags 50964-50965)
	// Forward matrices transform XYZ → camera RGB (same for both illuminants)
	forwardMatrix := generateForwardMatrix()
	forwardMatrix1Binary := srationalArrayToBytes(forwardMatrix)
	forwardMatrix2Binary := srationalArrayToBytes(forwardMatrix) // Same matrix for both

	// Step 5: Generate baseline exposure offset (tag 51109)
	// Use Adobe's baseline (-0.15 EV) for Nikon Z f
	baselineExposureBinary := srationalArrayToBytes([]SRational{{Numerator: -15, Denominator: 100}})

	// Step 6: Get profile name from metadata (tag 52552 - OPTIONAL)
	profileName := ""
	if name, ok := recipe.Metadata["profile_name"]; ok {
		if nameStr, ok := name.(string); ok {
			profileName = nameStr
		}
	}

	// Step 7: Get camera model from metadata (tag 50708 - OPTIONAL)
	cameraModel := ""
	if model, ok := recipe.Metadata["camera_model"]; ok {
		if modelStr, ok := model.(string); ok {
			cameraModel = modelStr
		}
	}

	// Step 8: Generate identity 3D color lookup table (tags 50981-50982)
	// Required by Lightroom even for neutral/pass-through profiles
	lookTableBinary := generateIdentityLUT()

	// Step 9: Create DNG file with binary tags
	dngData, err := writeDNG(toneCurveBinary, colorMatrix1Binary, colorMatrix2Binary, forwardMatrix1Binary, forwardMatrix2Binary, profileName, cameraModel, baselineExposureBinary, lookTableBinary)
	if err != nil {
		return nil, fmt.Errorf("failed to write DNG: %w", err)
	}

	return dngData, nil
}
