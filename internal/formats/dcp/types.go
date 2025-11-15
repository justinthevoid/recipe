package dcp

// DNG tag constants for Camera Profiles.
//
// Reference: Adobe DNG SDK 1.6 specification
const (
	TagColorMatrix1                = 50721 // SRational array (9 values)
	TagColorMatrix2                = 50722 // SRational array (9 values)
	TagCalibrationIlluminant1      = 50778 // Short (illuminant code)
	TagCalibrationIlluminant2      = 50779 // Short (illuminant code)
	TagProfileCalibrationSignature = 50932 // ASCII string (0xc6f4) - "com.adobe"
	TagProfileName                 = 50936 // ASCII string (0xc6f8) - profile display name
	TagProfileHueSatMapDims        = 50937 // Long array (3 values: hue, sat, val divisions)
	TagProfileHueSatMapData1       = 50938 // Float array
	TagProfileHueSatMapData2       = 50939 // Float array
	TagProfileToneCurve            = 50940 // Float array
	TagProfileEmbedPolicy          = 50941 // Long (0=allow copying, 3=no restrictions)
	TagProfileCopyright            = 50942 // ASCII string
	TagForwardMatrix1              = 50964 // SRational array (9 values)
	TagForwardMatrix2              = 50965 // SRational array (9 values)
	TagProfileLookTableDims        = 50981 // Long array (3 values: dimensions)
	TagProfileLookTableData        = 50982 // Float array
	TagProfileHueSatMapEncoding    = 51108 // Long (1=linear, 2=sRGB) - UNUSED in basic profiles
	TagProfileLookTableEncoding    = 51108 // Long (0xc7a4) - 1=sRGB (same as HueSat encoding)
	TagBaselineExposureOffset      = 51109 // SRational (0xc7a5)
	TagDefaultBlackRender          = 51110 // Long (0xc7a6) - 0=auto, 1=none
	TagProfileGroupName            = 52552 // ASCII string (0xcd48) - profile group name (optional)
)

// DNG illuminant codes (from DNG SDK)
const (
	IlluminantStandardLightA = 17 // Tungsten/Incandescent 2856K
	IlluminantD65            = 21 // Daylight 6504K
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
