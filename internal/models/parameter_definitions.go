package models

// ParameterOption defines a single choice for a discrete parameter.
type ParameterOption struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// ParameterDefinition defines the metadata for a UI control.
type ParameterDefinition struct {
	Key          string            `json:"key"`            // Mapping key in UniversalRecipe
	Label        string            `json:"label"`          // Human-readable label
	Type         string            `json:"type"`           // 'continuous' or 'discrete'
	Min          float64           `json:"min"`            // Minimum value
	Max          float64           `json:"max"`            // Maximum value
	Step         float64           `json:"step"`           // Step increment
	DefaultValue float64           `json:"defaultValue"`   // Default value
	Group        string            `json:"group"`          // UI group: Basic, Tone Curve, Detail, Color Mixer, Geometry
	Lane         string            `json:"lane,omitempty"` // UI lane hint: "left" or "right"
	Options      []ParameterOption `json:"options,omitempty"`
}

// GetNP3ParameterDefinitions returns the static list of all NP3 parameters.
func GetNP3ParameterDefinitions() []ParameterDefinition {
	return []ParameterDefinition{
		// Basic Adjustments
		{Key: "sharpness", Label: "Sharpness", Type: "continuous", Min: 0, Max: 150, Step: 1, DefaultValue: 0, Group: "Basic", Lane: "left"},
		{Key: "midRangeSharpening", Label: "Mid-range Sharpening", Type: "continuous", Min: -5, Max: 5, Step: 0.25, DefaultValue: 0, Group: "Basic", Lane: "left"},
		{Key: "clarity", Label: "Clarity", Type: "continuous", Min: -5, Max: 5, Step: 0.25, DefaultValue: 0, Group: "Basic", Lane: "left"},
		{Key: "contrast", Label: "Contrast", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic", Lane: "left"},
		{Key: "highlights", Label: "Highlights", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic", Lane: "left"},
		{Key: "shadows", Label: "Shadows", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic", Lane: "left"},
		{Key: "whites", Label: "Whites", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic", Lane: "left"},
		{Key: "blacks", Label: "Blacks", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic", Lane: "left"},
		{Key: "saturation", Label: "Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic", Lane: "left"},

		// Tone Curve (Parametric)
		{Key: "toneCurveShadows", Label: "Shadows", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Tone Curve", Lane: "left"},
		{Key: "toneCurveDarks", Label: "Darks", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Tone Curve", Lane: "left"},
		{Key: "toneCurveLights", Label: "Lights", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Tone Curve", Lane: "left"},
		{Key: "toneCurveHighlights", Label: "Highlights", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Tone Curve", Lane: "left"},
		{Key: "toneCurveShadowSplit", Label: "Shadow Split", Type: "continuous", Min: 0, Max: 100, Step: 1, DefaultValue: 25, Group: "Tone Curve", Lane: "left"},
		{Key: "toneCurveMidtoneSplit", Label: "Midtone Split", Type: "continuous", Min: 0, Max: 100, Step: 1, DefaultValue: 50, Group: "Tone Curve", Lane: "left"},
		{Key: "toneCurveHighlightSplit", Label: "Highlight Split", Type: "continuous", Min: 0, Max: 100, Step: 1, DefaultValue: 75, Group: "Tone Curve", Lane: "left"},

		// Color Mixer - Red
		{Key: "red.hue", Label: "Red Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "red.saturation", Label: "Red Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "red.luminance", Label: "Red Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},

		// Color Mixer - Orange
		{Key: "orange.hue", Label: "Orange Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "orange.saturation", Label: "Orange Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "orange.luminance", Label: "Orange Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},

		// Color Mixer - Yellow
		{Key: "yellow.hue", Label: "Yellow Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "yellow.saturation", Label: "Yellow Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "yellow.luminance", Label: "Yellow Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},

		// Color Mixer - Green
		{Key: "green.hue", Label: "Green Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "green.saturation", Label: "Green Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "green.luminance", Label: "Green Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},

		// Color Mixer - Aqua
		{Key: "aqua.hue", Label: "Aqua Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "aqua.saturation", Label: "Aqua Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "aqua.luminance", Label: "Aqua Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},

		// Color Mixer - Blue
		{Key: "blue.hue", Label: "Blue Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "blue.saturation", Label: "Blue Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "blue.luminance", Label: "Blue Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},

		// Detail
		{Key: "grain", Label: "Grain", Type: "continuous", Min: 0, Max: 100, Step: 1, DefaultValue: 0, Group: "Detail", Lane: "right"},

		// System
		{
			Key:          "processVersion",
			Label:        "Process Version",
			Type:         "discrete",
			Min:          1,
			Max:          2,
			DefaultValue: 1,
			Group:        "System",
			Lane:         "right",
			Options: []ParameterOption{
				{Label: "Version 1", Value: 1},
				{Label: "Version 2", Value: 2},
			},
		},

		// Color Mixer - Purple
		{Key: "purple.hue", Label: "Purple Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "purple.saturation", Label: "Purple Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "purple.luminance", Label: "Purple Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},

		// Color Mixer - Magenta
		{Key: "magenta.hue", Label: "Magenta Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "magenta.saturation", Label: "Magenta Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},
		{Key: "magenta.luminance", Label: "Magenta Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer", Lane: "right"},

		// Color Grading - Highlights
		{Key: "colorGrading.highlights.hue", Label: "Highlights Hue", Type: "continuous", Min: 0, Max: 360, Step: 1, DefaultValue: 0, Group: "Color Grading"},
		{Key: "colorGrading.highlights.chroma", Label: "Highlights Chroma", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Grading"},
		{Key: "colorGrading.highlights.brightness", Label: "Highlights Brightness", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Grading"},

		// Color Grading - Midtone
		{Key: "colorGrading.midtone.hue", Label: "Midtone Hue", Type: "continuous", Min: 0, Max: 360, Step: 1, DefaultValue: 0, Group: "Color Grading"},
		{Key: "colorGrading.midtone.chroma", Label: "Midtone Chroma", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Grading"},
		{Key: "colorGrading.midtone.brightness", Label: "Midtone Brightness", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Grading"},

		// Color Grading - Shadows
		{Key: "colorGrading.shadows.hue", Label: "Shadows Hue", Type: "continuous", Min: 0, Max: 360, Step: 1, DefaultValue: 0, Group: "Color Grading"},
		{Key: "colorGrading.shadows.chroma", Label: "Shadows Chroma", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Grading"},
		{Key: "colorGrading.shadows.brightness", Label: "Shadows Brightness", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Grading"},

		// Color Grading - Global
		{Key: "colorGrading.blending", Label: "Blending", Type: "continuous", Min: 0, Max: 100, Step: 1, DefaultValue: 50, Group: "Color Grading"},
		{Key: "colorGrading.balance", Label: "Balance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Grading"},
	}
}
