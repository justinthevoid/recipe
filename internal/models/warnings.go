package models

// WarnLevel represents the severity of a conversion warning
type WarnLevel int

const (
	// WarnInfo is for informational warnings that don't affect output quality
	WarnInfo WarnLevel = iota
	// WarnAdvisory is for parameters that are approximated or partially supported
	WarnAdvisory
	// WarnCritical is for parameters that cannot be converted and will be lost
	WarnCritical
)

// String returns the human-readable name of the warning level
func (w WarnLevel) String() string {
	switch w {
	case WarnInfo:
		return "Info"
	case WarnAdvisory:
		return "Advisory"
	case WarnCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// ConversionWarning represents a warning generated during format conversion.
// It captures information about parameters that could not be fully converted.
type ConversionWarning struct {
	Level       WarnLevel // Severity level (Info, Advisory, Critical)
	Parameter   string    // Name of the parameter that triggered the warning
	Value       string    // The value that couldn't be converted (optional)
	Message     string    // Human-readable description of the issue
	Alternative string    // Suggested workaround or alternative approach
}

// ConversionResult wraps the result of a conversion operation with any warnings
type ConversionResult struct {
	Warnings []ConversionWarning
}

// AddWarning appends a warning to the result
func (r *ConversionResult) AddWarning(level WarnLevel, param, value, message, alt string) {
	r.Warnings = append(r.Warnings, ConversionWarning{
		Level:       level,
		Parameter:   param,
		Value:       value,
		Message:     message,
		Alternative: alt,
	})
}

// HasCritical returns true if there are any critical warnings
func (r *ConversionResult) HasCritical() bool {
	for _, w := range r.Warnings {
		if w.Level == WarnCritical {
			return true
		}
	}
	return false
}

// HasWarnings returns true if there are any warnings
func (r *ConversionResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// ByLevel returns all warnings at or above the specified level
func (r *ConversionResult) ByLevel(minLevel WarnLevel) []ConversionWarning {
	var result []ConversionWarning
	for _, w := range r.Warnings {
		if w.Level >= minLevel {
			result = append(result, w)
		}
	}
	return result
}

// UnsupportedXMPParameters lists XMP parameters that cannot be converted to NP3
var UnsupportedXMPParameters = map[string]struct {
	Level       WarnLevel
	Description string
	Alternative string
}{
	// Frequency-based effects (no equivalent in NP3)
	"Clarity": {
		Level:       WarnAdvisory,
		Description: "Clarity affects mid-tone contrast via frequency separation",
		Alternative: "Use NP3 Mid-Range Sharpening for a similar effect",
	},
	"Texture": {
		Level:       WarnAdvisory,
		Description: "Texture affects fine detail enhancement",
		Alternative: "Use NP3 Sharpening Detail setting",
	},
	"Dehaze": {
		Level:       WarnAdvisory,
		Description: "Dehaze removes atmospheric haze",
		Alternative: "Use increased Contrast and reduced Blacks",
	},

	// RGB Channel curves (NP3 only has master curve)
	"ToneCurvePV2012Red": {
		Level:       WarnCritical,
		Description: "NP3 only supports master tone curve, not per-channel curves",
		Alternative: "Use Color Blender Red adjustments instead",
	},
	"ToneCurvePV2012Green": {
		Level:       WarnCritical,
		Description: "NP3 only supports master tone curve, not per-channel curves",
		Alternative: "Use Color Blender Green adjustments instead",
	},
	"ToneCurvePV2012Blue": {
		Level:       WarnCritical,
		Description: "NP3 only supports master tone curve, not per-channel curves",
		Alternative: "Use Color Blender Cyan/Blue adjustments instead",
	},

	// Noise reduction (not supported in NP3)
	"LuminanceSmoothing": {
		Level:       WarnInfo,
		Description: "Noise reduction not available in NP3 Picture Control",
		Alternative: "Apply noise reduction in NX Studio post-import",
	},
	"ColorNoiseReduction": {
		Level:       WarnInfo,
		Description: "Color noise reduction not available in NP3",
		Alternative: "Apply noise reduction in NX Studio post-import",
	},

	// Lens corrections (not in Picture Control)
	"LensProfileEnable": {
		Level:       WarnInfo,
		Description: "Lens corrections are camera/lens specific",
		Alternative: "Enable lens correction in NX Studio separately",
	},
}

// CheckUnsupportedParameter returns a warning if the parameter is unsupported
func CheckUnsupportedParameter(paramName string, value interface{}) *ConversionWarning {
	if info, exists := UnsupportedXMPParameters[paramName]; exists {
		// Only warn if value is non-zero/non-empty
		if !isZeroValue(value) {
			valueStr := ""
			switch v := value.(type) {
			case int:
				valueStr = string(rune(v))
			case float64:
				// Use simple formatting
				valueStr = "non-zero"
			case string:
				valueStr = v
			}
			return &ConversionWarning{
				Level:       info.Level,
				Parameter:   paramName,
				Value:       valueStr,
				Message:     info.Description,
				Alternative: info.Alternative,
			}
		}
	}
	return nil
}

// isZeroValue checks if a value is the zero value for its type
func isZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case int:
		return val == 0
	case float64:
		return val == 0.0
	case string:
		return val == ""
	case []interface{}:
		return len(val) == 0
	default:
		return false
	}
}
