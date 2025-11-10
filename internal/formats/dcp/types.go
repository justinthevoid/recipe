package dcp

// DNG tag constants for Camera Profiles.
//
// Reference: Adobe DNG SDK 1.6 specification
const (
	TagProfileName            = 52552 // ASCII string
	TagColorMatrix1           = 50721 // SRational array (9 values)
	TagColorMatrix2           = 50722 // SRational array (9 values)
	TagProfileToneCurve       = 50940 // Float array
	TagProfileHueSatMapDims   = 50937 // Long array (3 values: hue, sat, val divisions)
	TagProfileHueSatMapData1  = 50938 // Float array
	TagProfileHueSatMapData2  = 50939 // Float array
	TagProfileLookTableDims   = 50981 // Long array (3 values: dimensions)
	TagProfileLookTableData   = 50982 // Float array
	TagBaselineExposureOffset = 50730 // SRational
)

// CameraProfile represents a binary DNG Camera Profile.
//
// Unlike older XML-based DCP files, modern Adobe DCP files store
// profile data as binary TIFF tags in the DNG container.
type CameraProfile struct {
	ProfileName      string
	ToneCurve        []ToneCurvePoint
	ColorMatrix1     *Matrix
	ColorMatrix2     *Matrix
	BaselineExposure float64
}

// ToneCurvePoint represents a single point on the tone curve.
//
// In binary DNG format, tone curves are stored as arrays of floats
// where each pair represents (input, output) normalized to 0.0-1.0.
type ToneCurvePoint struct {
	Input  float64 // 0.0-1.0
	Output float64 // 0.0-1.0
}

// Matrix represents a 3x3 color transformation matrix.
//
// Used for color calibration (ColorMatrix1, ColorMatrix2).
// Stored as 9 SRational values in row-major order.
type Matrix struct {
	Rows [3][3]float64
}
