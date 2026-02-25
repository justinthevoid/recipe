package np3

// Metadata represents the raw structure and properties of a Nikon Picture Control file
// suitable for embedding in other formats like NKSC.
//
// Unlike UniversalRecipe, which normalizes data for conversion, Metadata retains
// the format-specific attributes needed for faithful reproduction or referencing.
type Metadata struct {
	// RawBytes contains the complete binary content of the NP3 file.
	RawBytes []byte `json:"rawBytes"`

	// Version is the internal format version (e.g., "0200").
	Version string `json:"version"`

	// Label is the user-defined name of the Picture Control.
	Label string `json:"label"`

	// Base details
	BaseFlag uint8 `json:"baseFlag"`

	// Basic Parameters
	Sharpening  float64 `json:"sharpening"`
	Clarity     float64 `json:"clarity"`
	GrainAmount float64 `json:"grainAmount"`

	// Advanced Parameters
	GrainSize          int     `json:"grainSize"`
	MidRangeSharpening float64 `json:"midRangeSharpening"`
	Contrast           int     `json:"contrast"`
	Highlights         int     `json:"highlights"`
	Shadows            int     `json:"shadows"`
	WhiteLevel         int     `json:"whiteLevel"`
	BlackLevel         int     `json:"blackLevel"`
	Saturation         int     `json:"saturation"`

	// Heuristic/Legacy
	Hue        int `json:"hue"`
	Brightness int `json:"brightness"`

	// Tone Curve
	ToneCurvePoints []ControlPoint `json:"toneCurvePoints"`
}
