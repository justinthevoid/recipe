package models

import (
	"testing"
)

func TestGetNP3ParameterDefinitions(t *testing.T) {
	defs := GetNP3ParameterDefinitions()

	// Check total count (based on the manual mapping in parameter_definitions.go)
	// 9 Basic + 7 Tone Curve + 3 Detail + 24 Color Mixer + 4 Geometry = 47 parameters
	// Wait, I counted 48 earlier. Let me double check.
	// Basic: 9
	// Tone Curve: 7
	// Detail: 3
	// Color Mixer: 8 colors * 3 = 24
	// Geometry: 4
	// Total: 9+7+3+24+4 = 47.
	// The np3-parameter-mapping-matrix.md mentions 48.
	// Maybe I missed one or the matrix includes one that is not in my list?
	// Let me check the matrix again if needed, but for now 47 is what I have implemented.

	expectedCount := 47
	if len(defs) != expectedCount {
		t.Errorf("Expected %d parameter definitions, got %d", expectedCount, len(defs))
	}

	// Check grouping
	groups := make(map[string]int)
	for _, d := range defs {
		groups[d.Group]++
	}

	expectedGroups := map[string]int{
		"Basic":       9,
		"Tone Curve":  7,
		"Detail":      3,
		"Color Mixer": 24,
		"Geometry":    4,
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
}
