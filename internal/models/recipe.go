// Package models defines the core data structures for photo editing recipe conversion.
package models

// ColorAdjustment represents HSL adjustments for a single color.
type ColorAdjustment struct {
	Hue        int `json:"hue" xml:"hue"`               // Hue shift: -180 to +180
	Saturation int `json:"saturation" xml:"saturation"` // Saturation: -100 to +100
	Luminance  int `json:"luminance" xml:"luminance"`   // Luminance: -100 to +100
}

// ToneCurvePoint represents a single point on a tone curve.
type ToneCurvePoint struct {
	Input  int `json:"input" xml:"input"`   // Input value: 0-255
	Output int `json:"output" xml:"output"` // Output value: 0-255
}

// SplitToning represents split toning/color grading parameters.
type SplitToning struct {
	HighlightHue        int `json:"highlightHue" xml:"highlightHue"`               // Highlight hue: 0-360
	HighlightSaturation int `json:"highlightSaturation" xml:"highlightSaturation"` // Highlight saturation: 0-100
	ShadowHue           int `json:"shadowHue" xml:"shadowHue"`                     // Shadow hue: 0-360
	ShadowSaturation    int `json:"shadowSaturation" xml:"shadowSaturation"`       // Shadow saturation: 0-100
	Balance             int `json:"balance,omitempty" xml:"balance,omitempty"`     // Balance: -100 to +100
}

// CameraProfile represents camera calibration settings.
type CameraProfile struct {
	RedHue         int `json:"redHue,omitempty" xml:"redHue,omitempty"`               // Red hue adjustment
	RedSaturation  int `json:"redSaturation,omitempty" xml:"redSaturation,omitempty"` // Red saturation adjustment
	GreenHue       int `json:"greenHue,omitempty" xml:"greenHue,omitempty"`           // Green hue adjustment
	GreenSaturation int `json:"greenSaturation,omitempty" xml:"greenSaturation,omitempty"` // Green saturation adjustment
	BlueHue        int `json:"blueHue,omitempty" xml:"blueHue,omitempty"`             // Blue hue adjustment
	BlueSaturation int `json:"blueSaturation,omitempty" xml:"blueSaturation,omitempty"` // Blue saturation adjustment
}

// UniversalRecipe is the universal data structure that represents all photo editing parameters
// from any supported format. It serves as the hub in the hub-and-spoke conversion pattern,
// eliminating N² conversion complexity by providing a single intermediate representation.
//
// All parsers convert: Format → UniversalRecipe
// All generators convert: UniversalRecipe → Format
type UniversalRecipe struct {
	// Metadata
	Name         string `json:"name,omitempty" xml:"name,omitempty"`               // Recipe name
	SourceFormat string `json:"sourceFormat,omitempty" xml:"sourceFormat,omitempty"` // Origin format: "np3", "xmp", "lrtemplate"

	// Basic Adjustments
	Exposure   float64 `json:"exposure,omitempty" xml:"exposure,omitempty"`     // Exposure: -5.0 to +5.0
	Contrast   int     `json:"contrast,omitempty" xml:"contrast,omitempty"`     // Contrast: -100 to +100
	Highlights int     `json:"highlights,omitempty" xml:"highlights,omitempty"` // Highlights: -100 to +100
	Shadows    int     `json:"shadows,omitempty" xml:"shadows,omitempty"`       // Shadows: -100 to +100
	Whites     int     `json:"whites,omitempty" xml:"whites,omitempty"`         // Whites: -100 to +100
	Blacks     int     `json:"blacks,omitempty" xml:"blacks,omitempty"`         // Blacks: -100 to +100

	// Presence
	Texture  int `json:"texture,omitempty" xml:"texture,omitempty"`   // Texture: -100 to +100
	Clarity  int `json:"clarity,omitempty" xml:"clarity,omitempty"`   // Clarity: -100 to +100
	Dehaze   int `json:"dehaze,omitempty" xml:"dehaze,omitempty"`     // Dehaze: -100 to +100
	Vibrance int `json:"vibrance,omitempty" xml:"vibrance,omitempty"` // Vibrance: -100 to +100
	Saturation int `json:"saturation,omitempty" xml:"saturation,omitempty"` // Saturation: -100 to +100

	// Sharpening
	Sharpness       int     `json:"sharpness,omitempty" xml:"sharpness,omitempty"`             // Sharpness amount: 0-150
	SharpnessRadius float64 `json:"sharpnessRadius,omitempty" xml:"sharpnessRadius,omitempty"` // Sharpness radius: 0.5-3.0
	SharpnessDetail int     `json:"sharpnessDetail,omitempty" xml:"sharpnessDetail,omitempty"` // Sharpness detail: 0-100
	SharpnessMasking int    `json:"sharpnessMasking,omitempty" xml:"sharpnessMasking,omitempty"` // Sharpness masking: 0-100

	// White Balance
	Temperature *int `json:"temperature,omitempty" xml:"temperature,omitempty"` // Temperature in Kelvin (nullable)
	Tint        int  `json:"tint,omitempty" xml:"tint,omitempty"`               // Tint: -150 to +150

	// HSL Adjustments (8 colors)
	Red     ColorAdjustment `json:"red,omitempty" xml:"red,omitempty"`         // Red HSL adjustments
	Orange  ColorAdjustment `json:"orange,omitempty" xml:"orange,omitempty"`   // Orange HSL adjustments
	Yellow  ColorAdjustment `json:"yellow,omitempty" xml:"yellow,omitempty"`   // Yellow HSL adjustments
	Green   ColorAdjustment `json:"green,omitempty" xml:"green,omitempty"`     // Green HSL adjustments
	Aqua    ColorAdjustment `json:"aqua,omitempty" xml:"aqua,omitempty"`       // Aqua HSL adjustments
	Blue    ColorAdjustment `json:"blue,omitempty" xml:"blue,omitempty"`       // Blue HSL adjustments
	Purple  ColorAdjustment `json:"purple,omitempty" xml:"purple,omitempty"`   // Purple HSL adjustments
	Magenta ColorAdjustment `json:"magenta,omitempty" xml:"magenta,omitempty"` // Magenta HSL adjustments

	// Tone Curve (Parametric)
	ToneCurveShadows      int `json:"toneCurveShadows,omitempty" xml:"toneCurveShadows,omitempty"`           // Tone curve shadows: -100 to +100
	ToneCurveDarks        int `json:"toneCurveDarks,omitempty" xml:"toneCurveDarks,omitempty"`               // Tone curve darks: -100 to +100
	ToneCurveLights       int `json:"toneCurveLights,omitempty" xml:"toneCurveLights,omitempty"`             // Tone curve lights: -100 to +100
	ToneCurveHighlights   int `json:"toneCurveHighlights,omitempty" xml:"toneCurveHighlights,omitempty"`     // Tone curve highlights: -100 to +100
	ToneCurveShadowSplit  int `json:"toneCurveShadowSplit,omitempty" xml:"toneCurveShadowSplit,omitempty"`   // Shadow split point: 0-100
	ToneCurveMidtoneSplit int `json:"toneCurveMidtoneSplit,omitempty" xml:"toneCurveMidtoneSplit,omitempty"` // Midtone split point: 0-100
	ToneCurveHighlightSplit int `json:"toneCurveHighlightSplit,omitempty" xml:"toneCurveHighlightSplit,omitempty"` // Highlight split point: 0-100

	// Point Curve (RGB)
	PointCurve      []ToneCurvePoint `json:"pointCurve,omitempty" xml:"pointCurve,omitempty"`           // Master curve points
	PointCurveRed   []ToneCurvePoint `json:"pointCurveRed,omitempty" xml:"pointCurveRed,omitempty"`     // Red channel curve points
	PointCurveGreen []ToneCurvePoint `json:"pointCurveGreen,omitempty" xml:"pointCurveGreen,omitempty"` // Green channel curve points
	PointCurveBlue  []ToneCurvePoint `json:"pointCurveBlue,omitempty" xml:"pointCurveBlue,omitempty"`   // Blue channel curve points

	// Split Toning / Color Grading
	SplitShadowHue        int `json:"splitShadowHue,omitempty" xml:"splitShadowHue,omitempty"`               // Shadow hue: 0-360
	SplitShadowSaturation int `json:"splitShadowSaturation,omitempty" xml:"splitShadowSaturation,omitempty"` // Shadow saturation: 0-100
	SplitHighlightHue     int `json:"splitHighlightHue,omitempty" xml:"splitHighlightHue,omitempty"`         // Highlight hue: 0-360
	SplitHighlightSaturation int `json:"splitHighlightSaturation,omitempty" xml:"splitHighlightSaturation,omitempty"` // Highlight saturation: 0-100
	SplitBalance          int `json:"splitBalance,omitempty" xml:"splitBalance,omitempty"`                   // Split balance: -100 to +100

	// Camera Calibration
	CameraProfile CameraProfile `json:"cameraProfile,omitempty" xml:"cameraProfile,omitempty"` // Camera calibration settings

	// Effects
	GrainAmount    int `json:"grainAmount,omitempty" xml:"grainAmount,omitempty"`       // Grain amount: 0-100
	GrainSize      int `json:"grainSize,omitempty" xml:"grainSize,omitempty"`           // Grain size: 0-100
	GrainRoughness int `json:"grainRoughness,omitempty" xml:"grainRoughness,omitempty"` // Grain roughness: 0-100

	// Vignette
	VignetteAmount   int `json:"vignetteAmount,omitempty" xml:"vignetteAmount,omitempty"`     // Vignette amount: -100 to +100
	VignetteMidpoint int `json:"vignetteMidpoint,omitempty" xml:"vignetteMidpoint,omitempty"` // Vignette midpoint: 0-100
	VignetteRoundness int `json:"vignetteRoundness,omitempty" xml:"vignetteRoundness,omitempty"` // Vignette roundness: -100 to +100
	VignetteFeather  int `json:"vignetteFeather,omitempty" xml:"vignetteFeather,omitempty"`   // Vignette feather: 0-100

	// Format-specific data (preserved for round-trip conversions)
	NP3ColorData    []map[string]interface{} `json:"np3ColorData,omitempty" xml:"np3ColorData,omitempty"`       // NP3 color data (raw)
	NP3RawParams    []map[string]interface{} `json:"np3RawParams,omitempty" xml:"np3RawParams,omitempty"`       // NP3 raw parameters
	NP3ToneCurveRaw []map[string]interface{} `json:"np3ToneCurveRaw,omitempty" xml:"np3ToneCurveRaw,omitempty"` // NP3 tone curve raw data

	// Generic metadata for unmappable parameters
	Metadata map[string]interface{} `json:"metadata,omitempty" xml:"-"` // Generic metadata for format-specific unmappable parameters

	// Complete raw binary data (for perfect round-trip of binary formats like NP3)
	// Not serialized to JSON/XML - used only for in-memory conversions
	FormatSpecificBinary map[string][]byte `json:"-" xml:"-"` // Raw binary data keyed by format (e.g., "np3_raw")
}
