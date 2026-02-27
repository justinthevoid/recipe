package models

// ParameterDefinition defines the metadata for a UI control.
type ParameterDefinition struct {
	Key          string  `json:"key"`          // Mapping key in UniversalRecipe
	Label        string  `json:"label"`        // Human-readable label
	Type         string  `json:"type"`         // 'continuous' or 'discrete'
	Min          float64 `json:"min"`          // Minimum value
	Max          float64 `json:"max"`          // Maximum value
	Step         float64 `json:"step"`         // Step increment
	DefaultValue float64 `json:"defaultValue"` // Default value
	Group        string  `json:"group"`        // UI group: Basic, Tone Curve, Detail, Color Mixer, Geometry
}

// GetNP3ParameterDefinitions returns the static list of all NP3 parameters.
func GetNP3ParameterDefinitions() []ParameterDefinition {
	return []ParameterDefinition{
		// Basic Adjustments
		{Key: "sharpness", Label: "Sharpness", Type: "continuous", Min: 0, Max: 150, Step: 1, DefaultValue: 0, Group: "Basic"},
		{Key: "midRangeSharpening", Label: "Mid-range Sharpening", Type: "continuous", Min: -5, Max: 5, Step: 0.25, DefaultValue: 0, Group: "Basic"},
		{Key: "clarity", Label: "Clarity", Type: "continuous", Min: -5, Max: 5, Step: 0.25, DefaultValue: 0, Group: "Basic"},
		{Key: "contrast", Label: "Contrast", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic"},
		{Key: "highlights", Label: "Highlights", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic"},
		{Key: "shadows", Label: "Shadows", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic"},
		{Key: "whites", Label: "Whites", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic"},
		{Key: "blacks", Label: "Blacks", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic"},
		{Key: "saturation", Label: "Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Basic"},

		// Tone Curve (Parametric)
		{Key: "toneCurveShadows", Label: "Shadows", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Tone Curve"},
		{Key: "toneCurveDarks", Label: "Darks", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Tone Curve"},
		{Key: "toneCurveLights", Label: "Lights", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Tone Curve"},
		{Key: "toneCurveHighlights", Label: "Highlights", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Tone Curve"},
		{Key: "toneCurveShadowSplit", Label: "Shadow Split", Type: "continuous", Min: 0, Max: 100, Step: 1, DefaultValue: 25, Group: "Tone Curve"},
		{Key: "toneCurveMidtoneSplit", Label: "Midtone Split", Type: "continuous", Min: 0, Max: 100, Step: 1, DefaultValue: 50, Group: "Tone Curve"},
		{Key: "toneCurveHighlightSplit", Label: "Highlight Split", Type: "continuous", Min: 0, Max: 100, Step: 1, DefaultValue: 75, Group: "Tone Curve"},

		// Color Mixer - Red
		{Key: "red.hue", Label: "Red Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "red.saturation", Label: "Red Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "red.luminance", Label: "Red Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},

		// Color Mixer - Orange
		{Key: "orange.hue", Label: "Orange Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "orange.saturation", Label: "Orange Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "orange.luminance", Label: "Orange Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},

		// Color Mixer - Yellow
		{Key: "yellow.hue", Label: "Yellow Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "yellow.saturation", Label: "Yellow Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "yellow.luminance", Label: "Yellow Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},

		// Color Mixer - Green
		{Key: "green.hue", Label: "Green Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "green.saturation", Label: "Green Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "green.luminance", Label: "Green Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},

		// Color Mixer - Aqua
		{Key: "aqua.hue", Label: "Aqua Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "aqua.saturation", Label: "Aqua Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "aqua.luminance", Label: "Aqua Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},

		// Color Mixer - Blue
		{Key: "blue.hue", Label: "Blue Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "blue.saturation", Label: "Blue Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "blue.luminance", Label: "Blue Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},

		// Color Mixer - Purple
		{Key: "purple.hue", Label: "Purple Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "purple.saturation", Label: "Purple Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "purple.luminance", Label: "Purple Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},

		// Color Mixer - Magenta
		{Key: "magenta.hue", Label: "Magenta Hue", Type: "continuous", Min: -180, Max: 180, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "magenta.saturation", Label: "Magenta Saturation", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},
		{Key: "magenta.luminance", Label: "Magenta Luminance", Type: "continuous", Min: -100, Max: 100, Step: 1, DefaultValue: 0, Group: "Color Mixer"},

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
