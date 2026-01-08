package np3

import (
	"fmt"
	"strings"
)

// ParseMetadata decodes a Nikon Picture Control (.np3) binary file into the Metadata struct.
// It uses formatting specific to Nikon sidecars rather than UniversalRecipe normalization.
func ParseMetadata(data []byte) (*Metadata, error) {
	// Validate file structure
	if err := validateFileStructure(data); err != nil {
		return nil, fmt.Errorf("metadata: %w", err)
	}

	// Reuse existing extraction logic
	params, err := extractParameters(data)
	if err != nil {
		return nil, fmt.Errorf("metadata: %w", err)
	}

	// Map internal params to exported Metadata
	m := &Metadata{
		RawBytes: data,
		Version:  string(data[3:7]),
		Label:    strings.TrimSpace(strings.TrimRight(params.name, "\x00")),
		BaseFlag: params.basePictureControlID,

		// Basic
		Sharpening:  params.sharpening,
		Clarity:     params.clarity,
		GrainAmount: params.grainAmount,

		// Advanced
		GrainSize:          params.grainSize,
		MidRangeSharpening: params.midRangeSharpening,
		Contrast:           params.contrast,
		Highlights:         params.highlights,
		Shadows:            params.shadows,
		WhiteLevel:         params.whiteLevel,
		BlackLevel:         params.blackLevel,
		Saturation:         params.saturation,

		// Heuristic/Legacy
		Hue:        params.hue,
		Brightness: int(params.brightness),

		ToneCurvePoints: params.toneCurve,
	}

	return m, nil
}
