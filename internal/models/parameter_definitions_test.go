package models

import (
	"testing"
)

func TestGetNP3ParameterDefinitions(t *testing.T) {
	defs := GetNP3ParameterDefinitions()

	// Actual counts:
	// Basic: 9 (sharpness, midRangeSharpening, clarity, contrast, highlights, shadows, whites, blacks, saturation)
	// Tone Curve: 7
	// Color Mixer: 24 (8 colors × 3)
	// Detail: 1 (grain)
	// System: 1 (processVersion)
	// Color Grading: 11 (3 zones × 3 + blending + balance)
	// Total: 53
	expectedCount := 53
	if len(defs) != expectedCount {
		t.Errorf("Expected %d parameter definitions, got %d", expectedCount, len(defs))
	}

	// Check grouping
	groups := make(map[string]int)
	for _, d := range defs {
		groups[d.Group]++
	}

	expectedGroups := map[string]int{
		"Basic":         9,
		"Tone Curve":    7,
		"Color Mixer":   24,
		"Detail":        1,
		"System":        1,
		"Color Grading": 11,
	}

	for g, count := range expectedGroups {
		if groups[g] != count {
			t.Errorf("Expected group %s to have %d parameters, got %d", g, count, groups[g])
		}
	}

	// Check for unique keys
	keys := make(map[string]bool)
	for _, d := range defs {
		if keys[d.Key] {
			t.Errorf("Duplicate parameter key found: %s", d.Key)
		}
		keys[d.Key] = true
	}

	// Verify all Basic parameters have Lane: "left"
	for _, d := range defs {
		if d.Group == "Basic" && d.Lane != "left" {
			t.Errorf("Basic parameter %s should have Lane 'left', got %q", d.Key, d.Lane)
		}
	}

	// Verify all Color Mixer parameters have Lane: "right"
	for _, d := range defs {
		if d.Group == "Color Mixer" && d.Lane != "right" {
			t.Errorf("Color Mixer parameter %s should have Lane 'right', got %q", d.Key, d.Lane)
		}
	}
}
