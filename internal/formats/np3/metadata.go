package np3

// Metadata represents the raw structure and properties of a Nikon Picture Control file
// suitable for embedding in other formats like NKSC.
//
// Unlike UniversalRecipe, which normalizes data for conversion, Metadata retains
// the format-specific attributes needed for faithful reproduction or referencing.
type Metadata struct {
	// RawBytes contains the complete binary content of the NP3 file.
	RawBytes []byte

	// Version is the internal format version (e.g., "0200").
	Version string

	// Label is the user-defined name of the Picture Control.
	Label string

	// Base details
	BaseFlag uint8

	// Basic Parameters
	Sharpening  float64
	Clarity     float64
	GrainAmount float64

	// Advanced Parameters
	GrainSize          int
	MidRangeSharpening float64
	Contrast           int
	Highlights         int
	Shadows            int
	WhiteLevel         int
	BlackLevel         int
	Saturation         int

	// Heuristic/Legacy
	Hue        int
	Brightness int

	// Tone Curve
	ToneCurvePoints []ControlPoint
}
